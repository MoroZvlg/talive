package talive

import "fmt"

// VWMA is a Volume Weighted Moving Average indicator.
// Uses raw ring buffers instead of SMA to optimize calculations.
// 1 division instead of 3 (Sum1 / Sum2 instead of (Sum1/period) / (Sum2 / period))
type VWMA struct {
	Period         int
	valueNumber    int
	closeVolBuffer *ringBuffer
	volBuffer      *ringBuffer
	out            []float64
}

// NewVWMA creates a new VWMA indicator with the given period.
func NewVWMA(period int) (*VWMA, error) {
	if period < 2 {
		return nil, fmt.Errorf("period should be greater than 1")
	}
	return &VWMA{
		Period:         period,
		valueNumber:    0,
		closeVolBuffer: newRingBuffer(period),
		volBuffer:      newRingBuffer(period),
		out:            make([]float64, 1),
	}, nil
}

func (v *VWMA) Next(candle ICandle) []float64 {
	v.valueNumber++

	cv := candle.Close() * candle.Volume()
	vol := candle.Volume()
	v.closeVolBuffer.Push(cv)
	v.volBuffer.Push(vol)

	if v.IsIdle() {
		v.out[0] = 0.0
		return v.out
	}

	v.out[0] = v.closeVolBuffer.Sum / v.volBuffer.Sum
	return v.out
}

func (v *VWMA) Current(candle ICandle) []float64 {
	if v.IsIdle() {
		v.out[0] = 0.0
		return v.out
	}

	cv := candle.Close() * candle.Volume()
	vol := candle.Volume()
	cvSum := v.closeVolBuffer.SumExceptLast() + cv
	vSum := v.volBuffer.SumExceptLast() + vol

	v.out[0] = cvSum / vSum
	return v.out
}

func (v *VWMA) IsIdle() bool {
	return v.valueNumber < v.Period
}

func (v *VWMA) IdlePeriod() int {
	return v.Period - 1
}

func (v *VWMA) IsWarmedUp() bool {
	return !v.IsIdle()
}

func (v *VWMA) WarmUpPeriod() int {
	return v.IdlePeriod()
}
