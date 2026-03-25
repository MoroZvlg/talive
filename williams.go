package talive

// Williams is a Williams %R indicator.
type Williams struct {
	Period      int
	valueNumber int
	lowest      *ringBuffer
	highest     *ringBuffer
	out         []float64
}

// NewWilliams creates a new Williams %R indicator with the given period.
func NewWilliams(period int) (*Williams, error) {
	return &Williams{
		Period:      period,
		valueNumber: 0,
		lowest:      newRingBuffer(period),
		highest:     newRingBuffer(period),
		out:         make([]float64, 1),
	}, nil
}

func (w *Williams) Next(candle ICandle) []float64 {
	w.valueNumber++

	w.lowest.Push(candle.Low())
	w.highest.Push(candle.High())

	if w.valueNumber < w.Period {
		return w.out
	}

	highestV := w.highest.Max()
	w.out[0] = (highestV - candle.Close()) / (highestV - w.lowest.Min()) * -100.0

	return w.out
}

func (w *Williams) Current(candle ICandle) []float64 {
	if w.valueNumber < w.Period {
		return w.out
	}

	lowestV := min(w.lowest.MinExceptLast(), candle.Low())
	highestV := max(w.highest.MaxExceptLast(), candle.High())

	w.out[0] = (highestV - candle.Close()) / (highestV - lowestV) * -100.0

	return w.out
}

func (w *Williams) IsIdle() bool {
	return w.valueNumber < w.Period
}

func (w *Williams) IsWarmedUp() bool {
	return w.valueNumber > w.WarmUpPeriod()
}

func (w *Williams) IdlePeriod() int {
	return w.Period - 1
}

func (w *Williams) WarmUpPeriod() int {
	return w.IdlePeriod()
}
