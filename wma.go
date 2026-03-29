package talive

import "fmt"

// WMA is a Weighted Moving Average indicator.
//
// O(1) update: each tick, every value ages by one position and loses 1 weight unit.
// Losing 1 unit from N values = subtracting their sum. The new value enters at full weight.
// So: weightedSum += newValue*period - buffer.Sum
type WMA struct {
	Period      int
	valueNumber int
	denominator float64
	weightedSum float64
	buffer      *ringBuffer
	out         []float64
}

// NewWMA creates a new WMA indicator with the given period.
func NewWMA(period int) (MA, error) {
	if period < 1 {
		return nil, fmt.Errorf("period should be greater than 0")
	}
	return &WMA{
		Period:      period,
		denominator: float64(period) * float64(period+1) / 2,
		buffer:      newRingBuffer(period),
		out:         make([]float64, 1),
	}, nil
}

func (wma *WMA) String() string {
	return fmt.Sprintf("WMA(%d)", wma.Period)
}

func (wma *WMA) next(value float64) float64 {
	wma.valueNumber++
	wma.weightedSum += value*float64(wma.Period) - wma.buffer.Sum
	wma.buffer.Push(value)

	if wma.IsIdle() {
		return 0.0
	}
	return wma.weightedSum / wma.denominator
}

func (wma *WMA) current(value float64) float64 {
	if wma.IsIdle() {
		return 0.0
	}
	ws := wma.weightedSum + value*float64(wma.Period) - wma.buffer.Sum
	return ws / wma.denominator
}

func (wma *WMA) Next(candle ICandle) []float64 {
	wma.out[0] = wma.next(candle.Close())
	return wma.out
}

func (wma *WMA) Current(candle ICandle) []float64 {
	wma.out[0] = wma.current(candle.Close())
	return wma.out
}

func (wma *WMA) IsIdle() bool {
	return wma.valueNumber < wma.Period
}

func (wma *WMA) IdlePeriod() int {
	return wma.Period - 1
}

func (wma *WMA) IsWarmedUp() bool {
	return !wma.IsIdle()
}

func (wma *WMA) WarmUpPeriod() int {
	return wma.IdlePeriod()
}
