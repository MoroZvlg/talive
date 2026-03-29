package worker

import (
	"context"
	"fmt"
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
	MASignals  []signal.Signaler
	OscSignals []signal.Signaler
}

func NewScreener() *Screener {
	var maSignals []signal.Signaler
	var oscSignals []signal.Signaler

	// ================
	// MA Signals (15)
	// ================

	// SMA(10, 20, 30, 50, 100, 200)
	for _, period := range []int{10, 20, 30, 50, 100, 200} {
		ma, _ := talive.NewSMA(period)
		maSignals = append(maSignals, signal.NewMASignal(ma))
	}

	// EMA(10, 20, 30, 50, 100, 200)
	for _, period := range []int{10, 20, 30, 50, 100, 200} {
		ma, _ := talive.NewEMA(period)
		maSignals = append(maSignals, signal.NewMASignal(ma))
	}

	// HMA(9)
	hma, _ := talive.NewHMA(9)
	maSignals = append(maSignals, signal.NewMASignal(hma))

	// VWMA(20)
	vwma, _ := talive.NewVWMA(20)
	maSignals = append(maSignals, signal.NewMASignal(vwma))

	// Ichimoku (shift=27 so leadA/leadB match Pine's lead1[26]/lead2[26])
	ich, _ := talive.NewIchimoku(9, 26, 52, 27)
	maSignals = append(maSignals, signal.NewIchimokuSignal(ich))

	// =====================
	// Oscillator Signals (11)
	// =====================

	// RSI(14): buy if rsi < 30 AND rising, sell if rsi > 70 AND falling
	rsi, _ := talive.NewRSI(14)
	oscSignals = append(oscSignals, signal.NewThresholdTrendSignal(rsi, 30, 70))

	// Stochastic(14, 3, 3): buy if K<20 & D<20 & K>D, sell if K>80 & D>80 & K<D
	stoch, _ := talive.NewStochastic(14, 3, 3)
	oscSignals = append(oscSignals, signal.NewStochasticSignal(stoch))

	// CCI(20): buy if cci < -100 AND rising, sell if cci > 100 AND falling
	cci, _ := talive.NewCCI(20)
	oscSignals = append(oscSignals, signal.NewThresholdTrendSignal(cci, -100, 100))

	// DMI/ADX(14): buy if adx>20 & adx rising & +DI>-DI, sell opposite
	dmi, _ := talive.NewDMI(14)
	oscSignals = append(oscSignals, signal.NewDMISignal(dmi))

	// AO: zero-line crossover + saucer pattern
	ao, _ := talive.NewAO()
	oscSignals = append(oscSignals, signal.NewAOSignal(ao))

	// Momentum(10): buy if mom > mom[1], sell if mom < mom[1]
	momentum, _ := talive.NewMomentum(10)
	oscSignals = append(oscSignals, signal.NewMomentumSignal(momentum))

	// MACD(12, 26, 9): buy if macd > signal, sell if macd < signal
	macd, _ := talive.NewMACD(12, 26, 9)
	oscSignals = append(oscSignals, signal.NewMACDSignal(macd))

	// StochRSI(14, 14, 3, 3) with EMA(50) trend filter
	stochRSI, _ := talive.NewStochasticRSI(14, 14, 3, 3)
	stochRSITrend, _ := talive.NewEMA(50)
	oscSignals = append(oscSignals, signal.NewStochRSISignal(stochRSI, stochRSITrend))

	// Williams %R(14): buy if wr < -80 AND rising, sell if wr > -20 AND falling
	williams, _ := talive.NewWilliams(14)
	oscSignals = append(oscSignals, signal.NewWilliamsSignal(williams))

	// Bull/Bear Power: EMA(13) for power, EMA(50) for trend
	bbpEma, _ := talive.NewEMA(13)
	bbpTrend, _ := talive.NewEMA(50)
	oscSignals = append(oscSignals, signal.NewBullBearPowerSignal(bbpEma, bbpTrend))

	// Ultimate Oscillator(7, 14, 28): buy if uo > 70, sell if uo < 30
	uo, _ := talive.NewUO(7, 14, 28)
	oscSignals = append(oscSignals, signal.NewThresholdSignal(uo, 30, 70, true))

	return &Screener{
		MASignals:  maSignals,
		OscSignals: oscSignals,
	}
}

func (s *Screener) Next(kline *entity.Kline) float64 {
	maSum := 0.0
	for _, sig := range s.MASignals {
		maSum += float64(sig.Next(kline))
	}
	maRating := maSum / float64(len(s.MASignals))

	oscSum := 0.0
	for _, sig := range s.OscSignals {
		oscSum += float64(sig.Next(kline))
	}
	oscRating := oscSum / float64(len(s.OscSignals))

	total := (maRating + oscRating) / 2
	fmt.Printf("MA: %.4f | Oscillators: %.4f | Total: %.4f\n", maRating, oscRating, total)
	return total
}

func (s *Screener) MaxWarmUp() int {
	result := 0
	for _, sig := range s.MASignals {
		result = max(result, sig.MaxWarmUp())
	}
	for _, sig := range s.OscSignals {
		result = max(result, sig.MaxWarmUp())
	}
	return result
}
