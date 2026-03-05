package talive

import "math"

type StdDev struct {
	period    int
	deviation float64
	variance  *Variance
	out       []float64
}

func NewStdDev(period int, deviation float64) (*StdDev, error) {
	// TODO: add validations
	variance, err := NewVariance(period)
	if err != nil {
		return nil, err
	}
	return &StdDev{
		period:    period,
		deviation: deviation,
		variance:  variance,
		out:       make([]float64, 1),
	}, nil
}

func (stdDev *StdDev) next(value float64) float64 {
	variance := stdDev.variance.next(value)
	return math.Sqrt(variance) * stdDev.deviation
}

func (stdDev *StdDev) current(value float64) float64 {
	variance := stdDev.variance.current(value)
	return math.Sqrt(variance) * stdDev.deviation
}

func (stdDev *StdDev) Next(candle ICandle) []float64 {
	stdDev.out[0] = stdDev.next(candle.Close())
	return stdDev.out
}

func (stdDev *StdDev) Current(candle ICandle) []float64 {
	stdDev.out[0] = stdDev.current(candle.Close())
	return stdDev.out
}

func (stdDev *StdDev) IsIdle() bool {
	return stdDev.variance.IsIdle()
}

type Variance struct {
	valueNumber     int
	period          int
	buffer          *ringBuffer
	quadraticBuffer *ringBuffer
	out             []float64
}

func NewVariance(period int) (*Variance, error) {
	// TODO: add validations
	return &Variance{
		valueNumber:     0,
		period:          period,
		buffer:          newRingBuffer(period),
		quadraticBuffer: newRingBuffer(period),
		out:             make([]float64, 1),
	}, nil
}

func (variance *Variance) next(value float64) float64 {
	variance.valueNumber++
	variance.buffer.Push(value)
	variance.quadraticBuffer.Push(value * value)
	if variance.IsIdle() {
		return 0.0
	}
	meanValue := variance.buffer.Sum / float64(variance.period)
	meanQuadroValue := variance.quadraticBuffer.Sum / float64(variance.period)
	return meanQuadroValue - meanValue*meanValue
}

func (variance *Variance) current(value float64) float64 {
	variance.valueNumber++
	if variance.IsIdle() {
		variance.valueNumber--
		return 0.0
	}
	meanValue := (variance.buffer.SumExceptLast() + value) / float64(variance.period)
	meanQuadroValue := (variance.quadraticBuffer.SumExceptLast() + value*value) / float64(variance.period)
	result := meanQuadroValue - meanValue*meanValue
	variance.valueNumber--
	return result
}

func (variance *Variance) Next(candle ICandle) []float64 {
	variance.out[0] = variance.next(candle.Close())
	return variance.out
}

func (variance *Variance) Current(candle ICandle) []float64 {
	variance.out[0] = variance.current(candle.Close())
	return variance.out
}

func (variance *Variance) IsIdle() bool {
	return variance.valueNumber < variance.period
}
