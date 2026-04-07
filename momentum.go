package talive

import "fmt"

// Momentum is a Momentum indicator.
type Momentum struct {
	Period      int
	SourceFunc  SourceFunc
	valueNumber int
	buffer      *ringBuffer
	out         []float64
}

// NewMomentum creates a new Momentum indicator with the given period.
func NewMomentum(period int, source SourceFunc) (*Momentum, error) {
	if period < 1 {
		return nil, fmt.Errorf("period should be greater than 0")
	}
	if source == nil {
		source = SourceClose
	}
	return &Momentum{
		Period:     period,
		SourceFunc: source,
		buffer:     newRingBuffer(period),
		out:        make([]float64, 1),
	}, nil
}

func (m *Momentum) String() string {
	return fmt.Sprintf("Momentum(%d)", m.Period)
}

func (m *Momentum) Next(candle ICandle) []float64 {
	m.valueNumber++

	value := m.SourceFunc(candle)
	oldest := m.buffer.Last()
	m.buffer.Push(value)

	if m.IsIdle() {
		return m.out
	}

	m.out[0] = value - oldest
	return m.out
}

func (m *Momentum) Current(candle ICandle) []float64 {
	if m.IsIdle() {
		return m.out
	}

	m.out[0] = m.SourceFunc(candle) - m.buffer.Last()
	return m.out
}

func (m *Momentum) IsIdle() bool {
	return m.valueNumber <= m.Period
}

func (m *Momentum) IdlePeriod() int {
	return m.Period
}

func (m *Momentum) IsWarmedUp() bool {
	return !m.IsIdle()
}

func (m *Momentum) WarmUpPeriod() int {
	return m.IdlePeriod()
}
