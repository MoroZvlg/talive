package talive

import "math"

// StdDev is a Standard Deviation indicator.
type StdDev struct {
	SourceFunc SourceFunc
	period     int
	deviation  float64
	variance   *Variance
	out        []float64
}

// NewStdDev creates a new Standard Deviation indicator.
func NewStdDev(period int, deviation float64, source SourceFunc) (*StdDev, error) {
	// TODO: add validations
	if source == nil {
		source = SourceClose
	}
	variance, err := NewVariance(period, nil)
	if err != nil {
		return nil, err
	}
	return &StdDev{
		SourceFunc: source,
		period:     period,
		deviation:  deviation,
		variance:   variance,
		out:        make([]float64, 1),
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
	stdDev.out[0] = stdDev.next(stdDev.SourceFunc(candle))
	return stdDev.out
}

func (stdDev *StdDev) Current(candle ICandle) []float64 {
	stdDev.out[0] = stdDev.current(stdDev.SourceFunc(candle))
	return stdDev.out
}

func (stdDev *StdDev) IsIdle() bool {
	return stdDev.variance.IsIdle()
}

func (stdDev *StdDev) IdlePeriod() int {
	return stdDev.variance.IdlePeriod()
}

func (stdDev *StdDev) IsWarmedUp() bool {
	return !stdDev.IsIdle()
}

func (stdDev *StdDev) WarmUpPeriod() int {
	return stdDev.IdlePeriod()
}

// Variance is a Variance indicator.
type Variance struct {
	SourceFunc      SourceFunc
	valueNumber     int
	period          int
	buffer          *ringBuffer
	quadraticBuffer *ringBuffer
	out             []float64
}

// NewVariance creates a new Variance indicator with the given period.
func NewVariance(period int, source SourceFunc) (*Variance, error) {
	// TODO: add validations
	if source == nil {
		source = SourceClose
	}
	return &Variance{
		SourceFunc:      source,
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
	variance.out[0] = variance.next(variance.SourceFunc(candle))
	return variance.out
}

func (variance *Variance) Current(candle ICandle) []float64 {
	variance.out[0] = variance.current(variance.SourceFunc(candle))
	return variance.out
}

func (variance *Variance) IsIdle() bool {
	return variance.valueNumber < variance.period
}

func (variance *Variance) IdlePeriod() int {
	return variance.period - 1
}

func (variance *Variance) IsWarmedUp() bool {
	return !variance.IsIdle()
}

func (variance *Variance) WarmUpPeriod() int {
	return variance.IdlePeriod()
}
