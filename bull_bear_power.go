package talive

import "fmt"

// BullBearPower is an Elder's Bull Bear Power indicator.
type BullBearPower struct {
	Period int

	ema MA
	out []float64
}

// NewBullBearPower creates a new Bull Bear Power indicator with the given period.
func NewBullBearPower(period int) (*BullBearPower, error) {
	if period < 2 {
		return nil, fmt.Errorf("period should be greater than 1")
	}
	ema, _ := NewEMA(period)
	return &BullBearPower{
		Period: period,
		ema:    ema,
		out:    make([]float64, 1),
	}, nil
}

func (bbp *BullBearPower) Next(candle ICandle) []float64 {
	emaVal := bbp.ema.next(candle.Close())

	if bbp.ema.IsIdle() {
		return bbp.out
	}

	bbp.out[0] = candle.High() + candle.Low() - 2*emaVal
	return bbp.out
}

func (bbp *BullBearPower) Current(candle ICandle) []float64 {
	if bbp.IsIdle() {
		return bbp.out
	}

	emaVal := bbp.ema.current(candle.Close())
	bbp.out[0] = candle.High() + candle.Low() - 2*emaVal
	return bbp.out
}

func (bbp *BullBearPower) IsIdle() bool {
	return bbp.ema.IsIdle()
}

func (bbp *BullBearPower) IdlePeriod() int {
	return bbp.ema.IdlePeriod()
}

func (bbp *BullBearPower) IsWarmedUp() bool {
	return bbp.ema.IsWarmedUp()
}

func (bbp *BullBearPower) WarmUpPeriod() int {
	return bbp.ema.WarmUpPeriod()
}
