package talive

import (
	"fmt"
)

type RSI struct {
	Period      int
	valueNumber int
	prevPrice   float64
	prevAvgGain float64
	prevAvgLoss float64
	out         []float64
}

func NewRSI(period int) (*RSI, error) {
	if period < 2 {
		return nil, fmt.Errorf("period should be greater than 1")
	}
	return &RSI{
		Period:      period,
		valueNumber: 0,
		prevPrice:   0.0,
		prevAvgGain: 0.0,
		prevAvgLoss: 0.0,
		out:         make([]float64, 1),
	}, nil
}

func (rsi *RSI) assignData(prevPrice, prevAvgGain, prevAvgLoss float64) {
	rsi.prevPrice = prevPrice
	rsi.prevAvgGain = prevAvgGain
	rsi.prevAvgLoss = prevAvgLoss
}

func (rsi *RSI) Next(candle ICandle) []float64 {
	value := candle.Close()
	rsi.valueNumber++

	if rsi.valueNumber == 1 {
		rsi.assignData(value, 0.0, 0.0)
		rsi.out[0] = 0.0
		return rsi.out
	}

	gain := 0.0
	loss := 0.0
	change := value - rsi.prevPrice
	if change > 0.0 {
		gain = change
	} else {
		loss = change
	}
	avgGain := 0.0
	avgLoss := 0.0

	if rsi.IsIdle() {
		prevGain := rsi.prevAvgGain * float64(rsi.valueNumber-2)
		prevLoss := rsi.prevAvgLoss * float64(rsi.valueNumber-2)
		avgGain = (prevGain + gain) / float64(rsi.valueNumber-1)
		avgLoss = (prevLoss - loss) / float64(rsi.valueNumber-1)
		rsi.assignData(value, avgGain, avgLoss)
		rsi.out[0] = 0.0
		return rsi.out
	}

	prevGain := rsi.prevAvgGain * float64(rsi.Period-1)
	prevLoss := rsi.prevAvgLoss * float64(rsi.Period-1)
	avgGain = (prevGain + gain) / float64(rsi.Period)
	avgLoss = (prevLoss - loss) / float64(rsi.Period)

	rsi.assignData(value, avgGain, avgLoss)
	rsi.out[0] = 100.0 * (avgGain / (avgGain + avgLoss))
	return rsi.out
}

func (rsi *RSI) Current(candle ICandle) []float64 {
	value := candle.Close()
	rsi.valueNumber++
	if rsi.IsIdle() {
		rsi.valueNumber--
		rsi.out[0] = 0.0
		return rsi.out
	}

	gain := 0.0
	loss := 0.0
	change := value - rsi.prevPrice
	if change > 0.0 {
		gain = change
	} else {
		loss = change
	}

	prevGain := rsi.prevAvgGain * float64(rsi.Period-1)
	prevLoss := rsi.prevAvgLoss * float64(rsi.Period-1)
	avgGain := (prevGain + gain) / float64(rsi.Period)
	avgLoss := (prevLoss - loss) / float64(rsi.Period)
	rsi.valueNumber--
	rsi.out[0] = 100.0 * (avgGain / (avgGain + avgLoss))
	return rsi.out
}

func (rsi *RSI) IsIdle() bool {
	return rsi.valueNumber <= rsi.Period
}

func (rsi *RSI) IsWarmedUp() bool {
	return rsi.valueNumber > int(rsi.WarmUpPeriod())
}

func (rsi *RSI) IdlePeriod() uint {
	return uint(rsi.Period)
}

func (rsi *RSI) WarmUpPeriod() uint {
	return rsi.IdlePeriod() + uint(rsi.Period*6)
}
