package talive

import (
	"fmt"
	"math"
)

// HMA is a Hull Moving Average indicator.
type HMA struct {
	Period      int
	valueNumber int
	halfWma     MA
	fullWma     MA
	sqrtWma     MA
	out         []float64
}

// NewHMA creates a new Hull Moving Average indicator with the given period.
func NewHMA(period int) (*HMA, error) {
	if period < 2 {
		return nil, fmt.Errorf("period should be greater than 1")
	}
	halfPeriod := period / 2
	if halfPeriod < 1 {
		halfPeriod = 1
	}
	sqrtPeriod := int(math.Floor(math.Sqrt(float64(period))))
	if sqrtPeriod < 1 {
		sqrtPeriod = 1
	}
	halfWma, _ := NewWMA(halfPeriod)
	fullWma, _ := NewWMA(period)
	sqrtWma, _ := NewWMA(sqrtPeriod)
	return &HMA{
		Period:  period,
		halfWma: halfWma,
		fullWma: fullWma,
		sqrtWma: sqrtWma,
		out:     make([]float64, 1),
	}, nil
}

func (hma *HMA) String() string {
	return fmt.Sprintf("HullMA(%d)", hma.Period)
}

func (hma *HMA) Next(candle ICandle) []float64 {
	hma.valueNumber++
	halfVal := hma.halfWma.next(candle.Close())
	fullVal := hma.fullWma.next(candle.Close())

	if hma.fullWma.IsIdle() {
		return hma.out
	}

	diff := 2*halfVal - fullVal
	hmaV := hma.sqrtWma.next(diff)

	if hma.sqrtWma.IsIdle() {
		return hma.out
	}

	hma.out[0] = hmaV
	return hma.out
}

func (hma *HMA) Current(candle ICandle) []float64 {
	if hma.IsIdle() {
		return hma.out
	}

	halfVal := hma.halfWma.current(candle.Close())
	fullVal := hma.fullWma.current(candle.Close())
	diff := 2*halfVal - fullVal
	hma.out[0] = hma.sqrtWma.current(diff)
	return hma.out
}

func (hma *HMA) IsIdle() bool {
	return hma.sqrtWma.IsIdle()
}

func (hma *HMA) IdlePeriod() int {
	return hma.fullWma.IdlePeriod() + hma.sqrtWma.IdlePeriod()
}

func (hma *HMA) IsWarmedUp() bool {
	return !hma.IsIdle()
}

func (hma *HMA) WarmUpPeriod() int {
	return hma.IdlePeriod()
}
