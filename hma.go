package talive

import (
	"fmt"
	"math"
)

// HMA is a Hull Moving Average indicator.
type HMA struct {
	Period         int
	HalfSourceFunc SourceFunc
	FullSourceFunc SourceFunc
	valueNumber    int
	halfWma        MA
	fullWma        MA
	sqrtWma        MA
	out            []float64
}

// NewHMA creates a new Hull Moving Average indicator with the given period.
func NewHMA(period int, halfSource, fullSource SourceFunc) (*HMA, error) {
	if period < 2 {
		return nil, fmt.Errorf("period should be greater than 1")
	}
	if halfSource == nil {
		halfSource = SourceClose
	}
	if fullSource == nil {
		fullSource = SourceClose
	}
	halfPeriod := period / 2
	if halfPeriod < 1 {
		halfPeriod = 1
	}
	sqrtPeriod := int(math.Floor(math.Sqrt(float64(period))))
	if sqrtPeriod < 1 {
		sqrtPeriod = 1
	}
	halfWma, _ := NewWMA(halfPeriod, halfSource)
	fullWma, _ := NewWMA(period, fullSource)
	sqrtWma, _ := NewWMA(sqrtPeriod, nil)
	return &HMA{
		Period:         period,
		HalfSourceFunc: halfSource,
		FullSourceFunc: fullSource,
		halfWma:        halfWma,
		fullWma:        fullWma,
		sqrtWma:        sqrtWma,
		out:            make([]float64, 1),
	}, nil
}

func (hma *HMA) String() string {
	return fmt.Sprintf("HullMA(%d)", hma.Period)
}

func (hma *HMA) Next(candle ICandle) []float64 {
	hma.valueNumber++
	halfVal := hma.halfWma.Next(candle)[0]
	fullVal := hma.fullWma.Next(candle)[0]

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

	halfVal := hma.halfWma.Current(candle)[0]
	fullVal := hma.fullWma.Current(candle)[0]
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
