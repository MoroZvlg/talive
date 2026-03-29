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

func (vwma *VWMA) String() string {
	return fmt.Sprintf("VWMA(%d)", vwma.Period)
}

func (vwma *VWMA) Next(candle ICandle) []float64 {
	vwma.valueNumber++

	cv := candle.Close() * candle.Volume()
	vol := candle.Volume()
	vwma.closeVolBuffer.Push(cv)
	vwma.volBuffer.Push(vol)

	if vwma.IsIdle() {
		vwma.out[0] = 0.0
		return vwma.out
	}

	vwma.out[0] = vwma.closeVolBuffer.Sum / vwma.volBuffer.Sum
	return vwma.out
}

func (vwma *VWMA) Current(candle ICandle) []float64 {
	if vwma.IsIdle() {
		vwma.out[0] = 0.0
		return vwma.out
	}

	cv := candle.Close() * candle.Volume()
	vol := candle.Volume()
	cvSum := vwma.closeVolBuffer.SumExceptLast() + cv
	vSum := vwma.volBuffer.SumExceptLast() + vol

	vwma.out[0] = cvSum / vSum
	return vwma.out
}

func (vwma *VWMA) IsIdle() bool {
	return vwma.valueNumber < vwma.Period
}

func (vwma *VWMA) IdlePeriod() int {
	return vwma.Period - 1
}

func (vwma *VWMA) IsWarmedUp() bool {
	return !vwma.IsIdle()
}

func (vwma *VWMA) WarmUpPeriod() int {
	return vwma.IdlePeriod()
}
