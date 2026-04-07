package talive

import (
	"fmt"
)

// RSI is a Relative Strength Index indicator.
type RSI struct {
	Period      int
	SourceFunc  SourceFunc
	valueNumber int
	prevPrice   float64
	gainSmma    MA
	lossSmma    MA
	out         []float64
}

// NewRSI creates a new RSI indicator with the given period.
func NewRSI(period int, source SourceFunc) (*RSI, error) {
	if period < 2 {
		return nil, fmt.Errorf("period should be greater than 1")
	}
	if source == nil {
		source = SourceClose
	}
	gainSmma, _ := NewSMMA(period, nil)
	lossSmma, _ := NewSMMA(period, nil)
	return &RSI{
		Period:     period,
		SourceFunc: source,
		gainSmma:   gainSmma,
		lossSmma:   lossSmma,
		out:        make([]float64, 1),
	}, nil
}

func (rsi *RSI) String() string {
	return fmt.Sprintf("RSI(%d)", rsi.Period)
}

func (rsi *RSI) Next(candle ICandle) []float64 {
	rsi.valueNumber++

	price := rsi.SourceFunc(candle)
	if rsi.valueNumber == 1 {
		rsi.prevPrice = price
		return rsi.out
	}

	gain, loss := rsi.gainLoss(price)
	rsi.prevPrice = price

	avgGain := rsi.gainSmma.next(gain)
	avgLoss := rsi.lossSmma.next(loss)

	if rsi.IsIdle() {
		return rsi.out
	}

	rsi.out[0] = 100.0 * avgGain / (avgGain + avgLoss)
	return rsi.out
}

func (rsi *RSI) Current(candle ICandle) []float64 {
	if rsi.IsIdle() {
		return rsi.out
	}

	gain, loss := rsi.gainLoss(rsi.SourceFunc(candle))
	avgGain := rsi.gainSmma.current(gain)
	avgLoss := rsi.lossSmma.current(loss)

	rsi.out[0] = 100.0 * avgGain / (avgGain + avgLoss)
	return rsi.out
}

func (rsi *RSI) IsIdle() bool {
	return rsi.valueNumber <= rsi.Period
}

func (rsi *RSI) IsWarmedUp() bool {
	return rsi.valueNumber > rsi.WarmUpPeriod()
}

func (rsi *RSI) IdlePeriod() int {
	return rsi.Period
}

func (rsi *RSI) WarmUpPeriod() int {
	return rsi.IdlePeriod() + rsi.Period*6
}

func (rsi *RSI) gainLoss(price float64) (gain, loss float64) {
	change := price - rsi.prevPrice
	if change > 0 {
		gain = change
	} else {
		loss = -change
	}
	return gain, loss
}
