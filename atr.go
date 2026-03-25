package talive

import (
	"fmt"
	"math"
)

// ATR is an Average True Range indicator.
type ATR struct {
	Period      int
	valueNumber int

	prevClose float64
	smma      MA

	out []float64
}

// NewATR creates a new ATR indicator with the given period.
func NewATR(period int) (*ATR, error) {
	if period < 2 {
		return nil, fmt.Errorf("period should be greater than 1")
	}
	smma, _ := NewSMMA(period)
	return &ATR{
		Period: period,
		smma:   smma,
		out:    make([]float64, 1),
	}, nil
}

func (atr *ATR) Next(candle ICandle) []float64 {
	atr.valueNumber++

	var trueRange float64
	if atr.valueNumber == 1 {
		trueRange = candle.High() - candle.Low()
	} else {
		highLow := candle.High() - candle.Low()
		highPrevClose := math.Abs(candle.High() - atr.prevClose)
		lowPrevClose := math.Abs(candle.Low() - atr.prevClose)
		trueRange = max(highLow, max(highPrevClose, lowPrevClose))
	}

	atr.prevClose = candle.Close()

	atrV := atr.smma.next(trueRange)

	if atr.smma.IsIdle() {
		return atr.out
	}

	atr.out[0] = atrV
	return atr.out
}

func (atr *ATR) Current(candle ICandle) []float64 {
	if atr.IsIdle() {
		return atr.out
	}

	highLow := candle.High() - candle.Low()
	highPrevClose := math.Abs(candle.High() - atr.prevClose)
	lowPrevClose := math.Abs(candle.Low() - atr.prevClose)
	trueRange := max(highLow, max(highPrevClose, lowPrevClose))

	atr.out[0] = atr.smma.current(trueRange)
	return atr.out
}

func (atr *ATR) IsIdle() bool {
	return atr.smma.IsIdle()
}

func (atr *ATR) IdlePeriod() int {
	return atr.smma.IdlePeriod()
}

func (atr *ATR) IsWarmedUp() bool {
	return atr.smma.IsWarmedUp()
}

func (atr *ATR) WarmUpPeriod() int {
	return atr.smma.WarmUpPeriod()
}
