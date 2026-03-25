package talive

import "fmt"

// StochasticRSI is a Stochastic RSI indicator.
type StochasticRSI struct {
	RSIPeriod int
	StochLen  int
	KSmooth   int
	DSmooth   int

	valueNumber int
	rsi         *RSI
	buffer      *ringBuffer
	kSMA        MA
	dSMA        MA
	out         []float64
}

// NewStochasticRSI creates a new Stochastic RSI indicator.
func NewStochasticRSI(rsiPeriod, stochLen, kSmooth, dSmooth int) (*StochasticRSI, error) {
	if rsiPeriod < 2 || stochLen < 2 || kSmooth < 1 || dSmooth < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}
	rsi, _ := NewRSI(rsiPeriod)
	kSMA, _ := NewSMA(kSmooth)
	dSMA, _ := NewSMA(dSmooth)
	return &StochasticRSI{
		RSIPeriod: rsiPeriod,
		StochLen:  stochLen,
		KSmooth:   kSmooth,
		DSmooth:   dSmooth,
		rsi:       rsi,
		buffer:    newRingBuffer(stochLen),
		kSMA:      kSMA,
		dSMA:      dSMA,
		out:       make([]float64, 2),
	}, nil
}

func (sr *StochasticRSI) Next(candle ICandle) []float64 {
	sr.valueNumber++
	rsiValue := sr.rsi.Next(candle)[0]
	if sr.rsi.IsIdle() {
		return sr.out
	}

	sr.buffer.Push(rsiValue)

	// we need to skip iteration with 1 value in buffer (min = max).
	if sr.buffer.Len() < 2 {
		return sr.out
	}

	minV, maxV := sr.buffer.MinMax()
	raw := sr.stochValue(rsiValue, minV, maxV)

	k := sr.kSMA.next(raw)
	if sr.kSMA.IsIdle() {
		return sr.out
	}
	sr.out[0] = k

	d := sr.dSMA.next(k)
	if !sr.dSMA.IsIdle() {
		sr.out[1] = d
	}

	return sr.out
}

func (sr *StochasticRSI) Current(candle ICandle) []float64 {
	if sr.IsIdle() {
		return sr.out
	}

	rsiValue := sr.rsi.Current(candle)[0]

	minV, maxV := sr.buffer.MinMaxExceptLast()
	minV = min(minV, rsiValue)
	maxV = max(maxV, rsiValue)

	raw := sr.stochValue(rsiValue, minV, maxV)
	k := sr.kSMA.current(raw)
	d := sr.dSMA.current(k)
	sr.out[0] = k
	sr.out[1] = d

	return sr.out
}

func (sr *StochasticRSI) stochValue(value, minV, maxV float64) float64 {
	if maxV == minV {
		return 0
	}
	return (value - minV) / (maxV - minV) * 100
}

func (sr *StochasticRSI) IsIdle() bool {
	return sr.dSMA.IsIdle()
}

func (sr *StochasticRSI) IdlePeriod() int {
	return sr.rsi.IdlePeriod() + 1 + sr.kSMA.IdlePeriod() + sr.dSMA.IdlePeriod()
}

func (sr *StochasticRSI) IsWarmedUp() bool {
	return sr.valueNumber > sr.WarmUpPeriod()
}

func (sr *StochasticRSI) WarmUpPeriod() int {
	// StochLen*2 because Min/Max is sensitive to small RSI errorsr. Subject to further clarification.
	return sr.rsi.WarmUpPeriod() + sr.StochLen*2 + sr.kSMA.WarmUpPeriod() + sr.dSMA.WarmUpPeriod()
}
