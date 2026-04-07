package talive

import "fmt"

// SMA is a Simple Moving Average indicator.
type SMA struct {
	Period      int
	SourceFunc  SourceFunc
	valueNumber int
	buffer      *ringBuffer
	out         []float64
}

// NewSMA creates a new SMA indicator with the given period.
func NewSMA(period int, source SourceFunc) (MA, error) {
	if source == nil {
		source = SourceClose
	}
	return &SMA{
		Period:      period,
		SourceFunc:  source,
		valueNumber: 0,
		buffer:      newRingBuffer(period),
		out:         make([]float64, 1),
	}, nil
}

func (sma *SMA) String() string {
	return fmt.Sprintf("SMA(%d)", sma.Period)
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
	if sma.IsIdle() {
		return 0.0
	}
	result := (sma.buffer.SumExceptLast() + value) / float64(sma.Period)
	return result
}

func (sma *SMA) Next(candle ICandle) []float64 {
	sma.out[0] = sma.next(sma.SourceFunc(candle))
	return sma.out
}

func (sma *SMA) Current(candle ICandle) []float64 {
	sma.out[0] = sma.current(sma.SourceFunc(candle))
	return sma.out
}

func (sma *SMA) IsIdle() bool {
	return sma.valueNumber < sma.Period
}

func (sma *SMA) IdlePeriod() int {
	return sma.Period - 1
}

func (sma *SMA) IsWarmedUp() bool {
	return !sma.IsIdle()
}

func (sma *SMA) WarmUpPeriod() int {
	return sma.IdlePeriod()
}
