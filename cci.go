package talive

import (
	"fmt"
	"math"
)

// CCI is a Commodity Channel Index indicator.
type CCI struct {
	Period      int
	SourceFunc  SourceFunc
	valueNumber int
	buffer      *ringBuffer
	out         []float64
}

// NewCCI creates a new CCI indicator with the given period.
// source replaces the typical price used in the formula. Pass nil to default
// to the classical HLC3 typical price.
func NewCCI(period int, source SourceFunc) (*CCI, error) {
	if source == nil {
		source = SourceHLC3
	}
	return &CCI{
		Period:      period,
		SourceFunc:  source,
		valueNumber: 0,
		buffer:      newRingBuffer(period),
		out:         make([]float64, 1),
	}, nil
}

func (cci *CCI) String() string {
	return fmt.Sprintf("CCI(%d)", cci.Period)
}

func (cci *CCI) Next(candle ICandle) []float64 {
	cci.valueNumber++

	typicalPrice := cci.SourceFunc(candle)
	cci.buffer.Push(typicalPrice)

	if cci.IsIdle() {
		return cci.out
	}

	// SMA replacement for optimisation
	// TODO: can be moved to RingBuffer??
	avg := cci.buffer.Sum / float64(cci.Period)

	// mean deviation
	var devSum float64
	for _, v := range cci.buffer.buffer {
		devSum += math.Abs(v - avg)
	}
	dev := devSum / float64(cci.Period)

	cci.out[0] = (typicalPrice - avg) / (0.015 * dev)
	return cci.out
}

func (cci *CCI) Current(candle ICandle) []float64 {
	if cci.IsIdle() {
		return cci.out
	}

	typicalPrice := cci.SourceFunc(candle)

	// SMA replacement for optimisation
	// TODO: can be moved to RingBuffer??
	avg := (cci.buffer.SumExceptLast() + typicalPrice) / float64(cci.Period)

	var devSum float64
	for i, v := range cci.buffer.buffer {
		if i == cci.buffer.writeIdx {
			devSum += math.Abs(typicalPrice - avg)
		} else {
			devSum += math.Abs(v - avg)
		}
	}
	dev := devSum / float64(cci.Period)

	cci.out[0] = (typicalPrice - avg) / (0.015 * dev)
	return cci.out
}

func (cci *CCI) IsIdle() bool {
	return cci.valueNumber < cci.Period
}

func (cci *CCI) IdlePeriod() int {
	return cci.Period - 1
}

func (cci *CCI) IsWarmedUp() bool {
	return !cci.IsIdle()
}

func (cci *CCI) WarmUpPeriod() int {
	return cci.IdlePeriod()
}
