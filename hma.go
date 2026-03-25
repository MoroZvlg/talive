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

func (h *HMA) Next(candle ICandle) []float64 {
	h.valueNumber++
	halfVal := h.halfWma.next(candle.Close())
	fullVal := h.fullWma.next(candle.Close())

	if h.fullWma.IsIdle() {
		return h.out
	}

	diff := 2*halfVal - fullVal
	hma := h.sqrtWma.next(diff)

	if h.sqrtWma.IsIdle() {
		return h.out
	}

	h.out[0] = hma
	return h.out
}

func (h *HMA) Current(candle ICandle) []float64 {
	if h.IsIdle() {
		return h.out
	}

	halfVal := h.halfWma.current(candle.Close())
	fullVal := h.fullWma.current(candle.Close())
	diff := 2*halfVal - fullVal
	h.out[0] = h.sqrtWma.current(diff)
	return h.out
}

func (h *HMA) IsIdle() bool {
	return h.sqrtWma.IsIdle()
}

func (h *HMA) IdlePeriod() int {
	return h.fullWma.IdlePeriod() + h.sqrtWma.IdlePeriod()
}

func (h *HMA) IsWarmedUp() bool {
	return !h.IsIdle()
}

func (h *HMA) WarmUpPeriod() int {
	return h.IdlePeriod()
}
