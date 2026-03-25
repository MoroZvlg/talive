package talive

// Stochastic is a Stochastic Oscillator indicator.
type Stochastic struct {
	KLen        int
	KSmooth     int
	DSmooth     int
	valueNumber int
	lowest      *ringBuffer
	highest     *ringBuffer
	kSMA        MA
	dSMA        MA
	out         []float64
}

// NewStochastic creates a new Stochastic Oscillator indicator with the given params.
func NewStochastic(kLen, kSmooth, dSmooth int) (*Stochastic, error) {
	kSMA, _ := NewSMA(kSmooth)
	dSMA, _ := NewSMA(dSmooth)
	return &Stochastic{
		KLen:        kLen,
		KSmooth:     kSmooth,
		DSmooth:     dSmooth,
		valueNumber: 0,
		lowest:      newRingBuffer(kLen),
		highest:     newRingBuffer(kLen),
		kSMA:        kSMA,
		dSMA:        dSMA,
		out:         make([]float64, 2),
	}, nil
}

func (stoch *Stochastic) Next(candle ICandle) []float64 {
	stoch.valueNumber++

	stoch.lowest.Push(candle.Low())
	stoch.highest.Push(candle.High())

	if stoch.valueNumber < stoch.KLen {
		return stoch.out
	}

	lowestLow := stoch.lowest.Min()
	value := (candle.Close() - lowestLow) / (stoch.highest.Max() - lowestLow) * 100.0
	kSmooth := stoch.kSMA.next(value)
	dSmooth := stoch.dSMA.next(kSmooth)
	stoch.out[0] = kSmooth
	stoch.out[1] = dSmooth

	return stoch.out
}

func (stoch *Stochastic) Current(candle ICandle) []float64 {
	if stoch.valueNumber < stoch.KLen {
		return stoch.out
	}

	lowestV := min(stoch.lowest.MinExceptLast(), candle.Low())
	highestV := max(stoch.highest.MinExceptLast(), candle.High())

	value := (candle.Close() - lowestV) / (highestV - lowestV) * 100.0
	kSmooth := stoch.kSMA.current(value)
	dSmooth := stoch.dSMA.current(kSmooth)
	stoch.out[0] = kSmooth
	stoch.out[1] = dSmooth

	return stoch.out
}

func (stoch *Stochastic) IsIdle() bool {
	return stoch.dSMA.IsIdle()
}

func (stoch *Stochastic) IsWarmedUp() bool {
	return stoch.valueNumber > stoch.WarmUpPeriod()
}

func (stoch *Stochastic) IdlePeriod() int {
	return stoch.KLen - 1 + stoch.dSMA.IdlePeriod()
}

func (stoch *Stochastic) WarmUpPeriod() int {
	return stoch.IdlePeriod()
}
