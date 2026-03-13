package talive

import "fmt"

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

func NewMACD(fastPeriod int, slowPeriod int, signalPeriod int) (*MACD, error) {
	if fastPeriod < 2 || slowPeriod < 2 || signalPeriod < 2 {
		return nil, fmt.Errorf("fastPeriod, slowPeriod, signalPeriod should be greater than 1")
	}
	fastEMA, errFast := NewEMA(fastPeriod)
	slowEMA, errSlow := NewEMA(slowPeriod)
	signalEMA, errSignal := NewEMA(signalPeriod)
	if errFast != nil || errSlow != nil || errSignal != nil {
		return nil, fmt.Errorf("error creating EMA: fast: %v, slow: %v, signal: %v", errFast, errSlow, errSignal)
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

func (macd *MACD) IdlePeriod() uint {
	return macd.slowEMA.IdlePeriod() + macd.signalEMA.IdlePeriod()
}

func (macd *MACD) IsWarmedUp() bool {
	return macd.valueNumber > int(macd.WarmUpPeriod())
}

func (macd *MACD) WarmUpPeriod() uint {
	return macd.IdlePeriod() + uint(macd.SlowPeriod*6)
}
