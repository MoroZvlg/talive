package talive

import "fmt"

// DonchianChannel is a Donchian Channel indicator.
//
// Output layout: [Upper, Mid, Lower].
type DonchianChannel struct {
	Period      int
	valueNumber int
	highest     *ringBuffer
	lowest      *ringBuffer
	out         []float64
}

// NewDonchianChannel creates a new Donchian Channel indicator with the given period.
func NewDonchianChannel(period int) (*DonchianChannel, error) {
	if period < 1 {
		return nil, fmt.Errorf("period should be positive")
	}
	return &DonchianChannel{
		Period:  period,
		highest: newRingBuffer(period),
		lowest:  newRingBuffer(period),
		out:     make([]float64, 3),
	}, nil
}

func (dc *DonchianChannel) String() string {
	return fmt.Sprintf("DonchianChannel(%d)", dc.Period)
}

func (dc *DonchianChannel) Next(candle OHLCV) []float64 {
	dc.valueNumber++
	dc.highest.Push(candle.High())
	dc.lowest.Push(candle.Low())
	if dc.IsIdle() {
		return dc.out
	}

	upper := dc.highest.Max()
	lower := dc.lowest.Min()
	dc.out[0] = upper
	dc.out[1] = (upper + lower) / 2
	dc.out[2] = lower
	return dc.out
}

func (dc *DonchianChannel) Current(candle OHLCV) []float64 {
	if dc.IsIdle() {
		return dc.out
	}

	upper := max(dc.highest.MaxExceptLast(), candle.High())
	lower := min(dc.lowest.MinExceptLast(), candle.Low())
	dc.out[0] = upper
	dc.out[1] = (upper + lower) / 2
	dc.out[2] = lower
	return dc.out
}

func (dc *DonchianChannel) IsIdle() bool {
	return dc.valueNumber < dc.Period
}

func (dc *DonchianChannel) IdlePeriod() int {
	return dc.Period - 1
}

func (dc *DonchianChannel) IsWarmedUp() bool {
	return !dc.IsIdle()
}

func (dc *DonchianChannel) WarmUpPeriod() int {
	return dc.IdlePeriod()
}
