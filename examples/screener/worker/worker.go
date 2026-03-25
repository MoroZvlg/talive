package worker

import (
	"context"
	"log/slog"
	"screener/binance"
	"screener/domain/entity"
	"screener/signal"
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
	Screener          *Screener
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
		Screener:        NewScreener(),
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
	klines, err := w.httpClient.LastKlines(w.ctx, w.Symbol, w.Screener.MaxWarmUp()+1)
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
	if kline.IsClosed {
		result := w.Screener.Next(&kline)
		w.log.Info("Screener result", "symbols", w.Symbol, "result", result, "receive->processed", time.Since(kline.TimeReceived))
	} else {
		w.log.Debug("Do not calcualte screener on open kline", "symbol", w.Symbol)
	}
}

type Screener struct {
	IndicatorsWeight map[signal.Signaler]float64
}

func NewScreener() *Screener {
	signalers := make(map[signal.Signaler]float64)

	// Oscillators Rating is calculated on the following oscillators:
	//Stochastic (14, 3, 3),
	//CCI (20),
	//ADX (14, 14),
	//AO,
	//Momentum (10),
	//Stochastic RSI (3, 3, 14, 14),
	//Williams %R (14),
	//Bulls and Bears Power and UO (7,14,28).
	rsiI, _ := talive.NewRSI(14)
	rsiSignaler := signal.NewRSISignal(rsiI)
	signalers[rsiSignaler] = 2.0

	macd, _ := talive.NewMACD(12, 26, 9)
	macdSignaler := signal.NewMACDSignal(macd)
	signalers[macdSignaler] = 2.0

	//the Ichimoku Cloud (9, 26, 52), VWMA (20), and HullMA (9).
	periods := []int{10, 20, 30, 50, 100, 200}
	for _, period := range periods {
		ma, _ := talive.NewEMA(period)
		signaler := signal.NewMASignal(ma)
		signalers[signaler] = 0.5
	}

	return &Screener{IndicatorsWeight: signalers}
}

func (s *Screener) Next(kline *entity.Kline) float64 {
	result := 0.0
	for signaler, weight := range s.IndicatorsWeight {
		result += float64(signaler.Next(kline)) * weight
	}
	return result / float64(len(s.IndicatorsWeight))
}

func (s *Screener) MaxWarmUp() int {
	result := 0
	for signaler := range s.IndicatorsWeight {
		result = max(result, signaler.MaxWarmUp())
	}
	return result
}
