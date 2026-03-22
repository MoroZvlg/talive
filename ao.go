package talive

// AO is an Awesome Oscillator indicator (Bill Williams).
type AO struct {
	fastSma MA
	slowSma MA
	out     []float64
}

// NewAO creates a new Awesome Oscillator indicator.
func NewAO() *AO {
	fastSma, _ := NewSMA(5)
	slowSma, _ := NewSMA(34)
	return &AO{
		fastSma: fastSma,
		slowSma: slowSma,
		out:     make([]float64, 1),
	}
}

func (ao *AO) Next(candle ICandle) []float64 {
	hl2 := (candle.High() + candle.Low()) / 2
	fast := ao.fastSma.next(hl2)
	slow := ao.slowSma.next(hl2)

	if ao.IsIdle() {
		return ao.out
	}

	ao.out[0] = fast - slow
	return ao.out
}

func (ao *AO) Current(candle ICandle) []float64 {
	if ao.IsIdle() {
		return ao.out
	}

	hl2 := (candle.High() + candle.Low()) / 2
	fast := ao.fastSma.current(hl2)
	slow := ao.slowSma.current(hl2)

	ao.out[0] = fast - slow
	return ao.out
}

func (ao *AO) IsIdle() bool {
	return ao.slowSma.IsIdle()
}

func (ao *AO) IdlePeriod() int {
	return ao.slowSma.IdlePeriod()
}

func (ao *AO) IsWarmedUp() bool {
	return !ao.IsIdle()
}

func (ao *AO) WarmUpPeriod() int {
	return ao.IdlePeriod()
}
