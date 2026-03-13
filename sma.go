package talive

import (
	"fmt"
)

type SMA struct {
	Period      int
	valueNumber int
	buffer      *ringBuffer
	out         []float64
}

func NewSMA(period int) (MA, error) {
	if period < 2 {
		return nil, fmt.Errorf("period should be greater than 1")
	}
	return &SMA{
		Period:      period,
		valueNumber: 0,
		buffer:      newRingBuffer(period),
		out:         make([]float64, 1),
	}, nil
}

func (sma *SMA) next(value float64) float64 {
	sma.buffer.Push(value)
	sma.valueNumber++
	if sma.IsIdle() {
		return 0.0
	}
	return sma.buffer.Sum / float64(sma.Period)
}

func (sma *SMA) current(value float64) float64 {
	sma.valueNumber++
	if sma.IsIdle() {
		return 0.0
	}
	result := (sma.buffer.SumExceptLast() + value) / float64(sma.Period)
	sma.valueNumber--
	return result
}

func (sma *SMA) Next(candle ICandle) []float64 {
	sma.out[0] = sma.next(candle.Close())
	return sma.out
}

func (sma *SMA) Current(candle ICandle) []float64 {
	sma.out[0] = sma.current(candle.Close())
	return sma.out
}

func (sma *SMA) IsIdle() bool {
	return sma.valueNumber < sma.Period
}

func (sma *SMA) IdlePeriod() uint {
	return uint(sma.Period - 1)
}

func (sma *SMA) IsWarmedUp() bool {
	return !sma.IsIdle()
}

func (sma *SMA) WarmUpPeriod() uint {
	return sma.IdlePeriod()
}
