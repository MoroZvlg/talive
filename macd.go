package talive

import "fmt"

// MACD is a Moving Average Convergence Divergence indicator.
type MACD struct {
	FastPeriod   int
	SlowPeriod   int
	SignalPeriod int
	valueNumber  int
	fastEMA      MA
	slowEMA      MA
	signalEMA    MA
	out          []float64
}

// NewMACD creates a new MACD indicator with the given periods.
func NewMACD(fastPeriod int, slowPeriod int, signalPeriod int) (*MACD, error) {
	if fastPeriod < 2 || slowPeriod < 2 || signalPeriod < 2 {
		return nil, fmt.Errorf("fastPeriod, slowPeriod, signalPeriod should be greater than 1")
	}
	fastEMA, errFast := NewEMA(fastPeriod)
	slowEMA, errSlow := NewEMA(slowPeriod)
	signalEMA, errSignal := NewEMA(signalPeriod)
	if errFast != nil || errSlow != nil || errSignal != nil {
		return nil, fmt.Errorf("error creating EMA: fast: %w, slow: %w, signal: %w", errFast, errSlow, errSignal)
	}

	return &MACD{
		FastPeriod:   fastPeriod,
		SlowPeriod:   slowPeriod,
		SignalPeriod: signalPeriod,
		fastEMA:      fastEMA,
		slowEMA:      slowEMA,
		signalEMA:    signalEMA,
		out:          make([]float64, 3),
	}, nil
}

func (macd *MACD) String() string {
	return fmt.Sprintf("MACD(%d,%d,%d)", macd.FastPeriod, macd.SlowPeriod, macd.SignalPeriod)
}

func (macd *MACD) Next(candle ICandle) []float64 {
	macd.valueNumber++
	value := candle.Close()
	outMACD := macd.fastEMA.next(value) - macd.slowEMA.next(value)

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
	value := candle.Close()
	outMACD := macd.fastEMA.current(value) - macd.slowEMA.current(value)

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
