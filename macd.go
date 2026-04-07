package talive

import "fmt"

// MACD is a Moving Average Convergence Divergence indicator.
type MACD struct {
	FastPeriod     int
	SlowPeriod     int
	SignalPeriod   int
	FastSourceFunc SourceFunc
	SlowSourceFunc SourceFunc
	valueNumber    int
	fastEMA        MA
	slowEMA        MA
	signalEMA      MA
	out            []float64
}

// NewMACD creates a new MACD indicator with the given periods.
func NewMACD(fastPeriod int, slowPeriod int, signalPeriod int, fastSource, slowSource SourceFunc) (*MACD, error) {
	if fastPeriod < 2 || slowPeriod < 2 || signalPeriod < 2 {
		return nil, fmt.Errorf("fastPeriod, slowPeriod, signalPeriod should be greater than 1")
	}
	if fastSource == nil {
		fastSource = SourceClose
	}
	if slowSource == nil {
		slowSource = SourceClose
	}
	fastEMA, errFast := NewEMA(fastPeriod, fastSource)
	slowEMA, errSlow := NewEMA(slowPeriod, slowSource)
	signalEMA, errSignal := NewEMA(signalPeriod, nil)
	if errFast != nil || errSlow != nil || errSignal != nil {
		return nil, fmt.Errorf("error creating EMA: fast: %w, slow: %w, signal: %w", errFast, errSlow, errSignal)
	}

	return &MACD{
		FastPeriod:     fastPeriod,
		SlowPeriod:     slowPeriod,
		SignalPeriod:   signalPeriod,
		FastSourceFunc: fastSource,
		SlowSourceFunc: slowSource,
		fastEMA:        fastEMA,
		slowEMA:        slowEMA,
		signalEMA:      signalEMA,
		out:            make([]float64, 3),
	}, nil
}

func (macd *MACD) String() string {
	return fmt.Sprintf("MACD(%d,%d,%d)", macd.FastPeriod, macd.SlowPeriod, macd.SignalPeriod)
}

func (macd *MACD) Next(candle ICandle) []float64 {
	macd.valueNumber++
	outMACD := macd.fastEMA.Next(candle)[0] - macd.slowEMA.Next(candle)[0]

	if macd.slowEMA.IsIdle() {
		macd.out[0] = 0.0
		macd.out[1] = 0.0
		macd.out[2] = 0.0
		return macd.out
	}

	outMACDSignal := macd.signalEMA.next(outMACD)
	if macd.signalEMA.IsIdle() {
		macd.out[0] = outMACD
		macd.out[1] = 0.0
		macd.out[2] = 0.0
		return macd.out
	}

	macd.out[0] = outMACD
	macd.out[1] = outMACDSignal
	macd.out[2] = outMACD - outMACDSignal
	return macd.out
}

func (macd *MACD) Current(candle ICandle) []float64 {
	outMACD := macd.fastEMA.Current(candle)[0] - macd.slowEMA.Current(candle)[0]

	if macd.slowEMA.IsIdle() {
		macd.out[0] = 0.0
		macd.out[1] = 0.0
		macd.out[2] = 0.0
		return macd.out
	}

	outMACDSignal := macd.signalEMA.current(outMACD)
	if macd.signalEMA.IsIdle() {
		macd.out[0] = outMACD
		macd.out[1] = 0.0
		macd.out[2] = 0.0
		return macd.out
	}

	macd.out[0] = outMACD
	macd.out[1] = outMACDSignal
	macd.out[2] = outMACD - outMACDSignal
	return macd.out
}

func (macd *MACD) IsIdle() bool {
	return macd.signalEMA.IsIdle()
}

func (macd *MACD) IdlePeriod() int {
	return macd.slowEMA.IdlePeriod() + macd.signalEMA.IdlePeriod()
}

func (macd *MACD) IsWarmedUp() bool {
	return macd.valueNumber > macd.WarmUpPeriod()
}

func (macd *MACD) WarmUpPeriod() int {
	return macd.IdlePeriod() + macd.SlowPeriod*6
}
