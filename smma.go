package talive

import "fmt"

// SMMA is a Smoothed Moving Average (Wilder's Moving Average) indicator.
type SMMA struct {
	Period       int
	Alpha        float64
	valuesNumber int
	prev         float64
	out          []float64
}

// NewSMMA creates a new SMMA indicator with the given period.
func NewSMMA(period int) (MA, error) {
	if period < 2 {
		return nil, fmt.Errorf("period should be greater than 1")
	}
	return &SMMA{
		Period: period,
		Alpha:  1.0 / float64(period),
		out:    make([]float64, 1),
	}, nil
}

func (smma *SMMA) String() string {
	return fmt.Sprintf("SMMA(%d)", smma.Period)
}

func (smma *SMMA) next(value float64) float64 {
	smma.valuesNumber++
	if smma.IsIdle() {
		smma.prev += value
		return 0.0
	}
	if smma.valuesNumber == smma.Period {
		smma.prev = (smma.prev + value) / float64(smma.Period)
		return smma.prev
	}

	current := value*smma.Alpha + smma.prev*(1-smma.Alpha)
	smma.prev = current
	return current
}

func (smma *SMMA) current(value float64) float64 {
	if smma.IsIdle() {
		return 0.0
	}
	if smma.valuesNumber+1 == smma.Period {
		return (smma.prev + value) / float64(smma.Period)
	}
	return value*smma.Alpha + smma.prev*(1-smma.Alpha)
}

func (smma *SMMA) Next(candle ICandle) []float64 {
	smma.out[0] = smma.next(candle.Close())
	return smma.out
}

func (smma *SMMA) Current(candle ICandle) []float64 {
	smma.out[0] = smma.current(candle.Close())
	return smma.out
}

func (smma *SMMA) IsIdle() bool {
	return smma.valuesNumber < smma.Period
}

func (smma *SMMA) IdlePeriod() int {
	return smma.Period - 1
}

func (smma *SMMA) IsWarmedUp() bool {
	return smma.valuesNumber > smma.WarmUpPeriod()
}

func (smma *SMMA) WarmUpPeriod() int {
	return smma.IdlePeriod() + smma.Period*6
}
