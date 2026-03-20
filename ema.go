package talive

import "fmt"

// EMA is an Exponential Moving Average indicator.
type EMA struct {
	Period       int
	Alpha        float64
	valuesNumber int
	prevEma      float64
	out          []float64
}

// NewEMA creates a new EMA indicator with the given period.
func NewEMA(period int) (MA, error) {
	if period < 2 {
		return nil, fmt.Errorf("period should be greater than 1")
	}
	return &EMA{
		Period:       period,
		Alpha:        2.0 / float64(period+1),
		valuesNumber: 0,
		prevEma:      0.0,
		out:          make([]float64, 1),
	}, nil
}

func (ema *EMA) next(value float64) float64 {
	ema.valuesNumber++
	if ema.IsIdle() {
		// first EMA value = avg of close prices. We need to save them
		ema.prevEma += value
		return 0.0
	}
	if ema.valuesNumber == ema.Period {
		ema.prevEma = (ema.prevEma + value) / float64(ema.Period)
		return ema.prevEma
	}

	currentEma := value*ema.Alpha + ema.prevEma*(1-ema.Alpha)
	ema.prevEma = currentEma
	return currentEma
}

func (ema *EMA) current(value float64) float64 {
	if ema.IsIdle() {
		return 0.0
	}
	if ema.valuesNumber+1 == ema.Period {
		result := (ema.prevEma + value) / float64(ema.Period)
		return result
	}
	result := value*ema.Alpha + ema.prevEma*(1-ema.Alpha)
	return result
}

func (ema *EMA) Next(candle ICandle) []float64 {
	ema.out[0] = ema.next(candle.Close())
	return ema.out
}

func (ema *EMA) Current(candle ICandle) []float64 {
	ema.out[0] = ema.current(candle.Close())
	return ema.out
}

func (ema *EMA) IsIdle() bool {
	return ema.valuesNumber < ema.Period
}

func (ema *EMA) IdlePeriod() int {
	return ema.Period - 1
}

func (ema *EMA) IsWarmedUp() bool {
	return ema.valuesNumber > ema.WarmUpPeriod()
}

func (ema *EMA) WarmUpPeriod() int {
	return ema.IdlePeriod() + ema.Period*2
}
