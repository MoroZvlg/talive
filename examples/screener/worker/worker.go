package worker

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"screener/binance"
	"screener/domain/entity"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MoroZvlg/talive"
)

type ScreenerWorker struct {
	PID               uint
	log               *slog.Logger
	ctx               context.Context
	ready             atomic.Bool
	Symbol            string
	Indicators        []talive.IIndicator
	httpClient        *binance.HTTPClient
	wsLoadingBuffer   []entity.Kline
	lastProcessedTime time.Time
}

func NewScreenerWorker(
	ctx context.Context,
	log *slog.Logger,
	pid uint,
	symbol string,
	client *binance.HTTPClient,
) *ScreenerWorker {
	return &ScreenerWorker{
		log:             log.With("PID", pid, "symbol", symbol),
		ctx:             ctx,
		PID:             pid,
		Symbol:          symbol,
		Indicators:      GenerateRandomIndicators(),
		httpClient:      client,
		wsLoadingBuffer: make([]entity.Kline, 0),
	}
}

func (w *ScreenerWorker) Start(klineCh <-chan entity.Kline, wg *sync.WaitGroup) {
	go w.fetchKlinesHistory()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-w.ctx.Done():
				return
			case kline, ok := <-klineCh:
				if !ok {
					return
				}
				if w.ready.Load() {
					for _, bufKline := range w.wsLoadingBuffer {
						if bufKline.TimeStart.After(w.lastProcessedTime) {
							w.processKline(bufKline)
						}
					}
					w.wsLoadingBuffer = nil
					w.processKline(kline)
				} else {
					w.wsLoadingBuffer = append(w.wsLoadingBuffer, kline)
				}
			}
		}
	}()
}

func (w *ScreenerWorker) fetchKlinesHistory() {
	var maxWarmUp int
	for _, i := range w.Indicators {
		if i.WarmUpPeriod() > maxWarmUp {
			maxWarmUp = i.WarmUpPeriod()
		}
	}
	klines, err := w.httpClient.LastKlines(w.ctx, w.Symbol, maxWarmUp+1)
	if err != nil {
		w.log.Error("Error fetching LastKlines", "error", err)
		// NOTE: it's ok for example to go without history. Do not return
	}
	for _, kline := range klines {
		w.processKline(kline)
	}
	w.ready.Store(true)
}

func (w *ScreenerWorker) processKline(kline entity.Kline) {
	for _, indicator := range w.Indicators {
		var result []float64
		if kline.IsClosed {
			result = indicator.Next(&kline)
		} else {
			result = indicator.Current(&kline)
		}
		msg := fmt.Sprintf("%T", indicator)
		w.log.Debug(msg, "result", result, "closed", kline.IsClosed)
	}
	w.lastProcessedTime = kline.TimeStart
	if kline.IsClosed {
		w.log.Info(
			"Timings",
			"symbols",
			w.Symbol,
			"receive->processed",
			time.Since(kline.TimeReceived),
			"closed->processed",
			time.Since(kline.TimeStart.Add(time.Minute)),
		)
	} else {
		w.log.Debug("Timings", "symbols", w.Symbol, "receive->processed", time.Since(kline.TimeReceived))
	}
}

func GenerateRandomIndicators() []talive.IIndicator {
	rsi, _ := talive.NewRSI(rand.IntN(45) + 5)
	bband, _ := talive.NewBBands(rand.IntN(28)+2, rand.Float64()*2+1, rand.Float64()*2+1, talive.SMAtype)
	ema, _ := talive.NewEMA(rand.IntN(190) + 10)
	sma, _ := talive.NewSMA(rand.IntN(90) + 10)
	mfi, _ := talive.NewMFI(rand.IntN(45) + 5)
	slow := rand.IntN(45) + 5
	macd, _ := talive.NewMACD(slow/2, slow, 9)

	return []talive.IIndicator{rsi, bband, ema, sma, mfi, macd}
}
