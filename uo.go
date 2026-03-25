package talive

import (
	"fmt"
	"math"
)

// UO is an Ultimate Oscillator indicator.
type UO struct {
	PeriodMin   int
	PeriodMid   int
	PeriodMax   int
	valueNumber int
	prevClose   float64
	bpBuf       [3]*ringBuffer
	trBuf       [3]*ringBuffer
	out         []float64
}

// NewUO creates a new Ultimate Oscillator indicator.
func NewUO(PeriodMin, PeriodMid, PeriodMax int) (*UO, error) {
	if PeriodMin < 1 || PeriodMid < 1 || PeriodMax < 1 {
		return nil, fmt.Errorf("periods should be greater than 0")
	}
	return &UO{
		PeriodMin: PeriodMin,
		PeriodMid: PeriodMid,
		PeriodMax: PeriodMax,
		bpBuf:     [3]*ringBuffer{newRingBuffer(PeriodMin), newRingBuffer(PeriodMid), newRingBuffer(PeriodMax)},
		trBuf:     [3]*ringBuffer{newRingBuffer(PeriodMin), newRingBuffer(PeriodMid), newRingBuffer(PeriodMax)},
		out:       make([]float64, 1),
	}, nil
}

func (u *UO) Next(candle ICandle) []float64 {
	u.valueNumber++

	if u.valueNumber == 1 {
		u.prevClose = candle.Close()
		return u.out
	}

	bp := candle.Close() - math.Min(candle.Low(), u.prevClose)
	tr := math.Max(candle.High(), u.prevClose) - math.Min(candle.Low(), u.prevClose)
	u.prevClose = candle.Close()

	for i := 0; i < 3; i++ {
		u.bpBuf[i].Push(bp)
		u.trBuf[i].Push(tr)
	}

	if u.IsIdle() {
		return u.out
	}

	avgMin := u.bpBuf[0].Sum / u.trBuf[0].Sum
	avgMid := u.bpBuf[1].Sum / u.trBuf[1].Sum
	avgMax := u.bpBuf[2].Sum / u.trBuf[2].Sum
	u.out[0] = 100 * (4*avgMin + 2*avgMid + avgMax) / 7
	return u.out
}

func (u *UO) Current(candle ICandle) []float64 {
	if u.IsIdle() {
		return u.out
	}

	bp := candle.Close() - math.Min(candle.Low(), u.prevClose)
	tr := math.Max(candle.High(), u.prevClose) - math.Min(candle.Low(), u.prevClose)

	avgMin := (u.bpBuf[0].SumExceptLast() + bp) / (u.trBuf[0].SumExceptLast() + tr)
	avgMid := (u.bpBuf[1].SumExceptLast() + bp) / (u.trBuf[1].SumExceptLast() + tr)
	avgMax := (u.bpBuf[2].SumExceptLast() + bp) / (u.trBuf[2].SumExceptLast() + tr)

	u.out[0] = 100 * (4*avgMin + 2*avgMid + avgMax) / 7
	return u.out
}

func (u *UO) IsIdle() bool {
	maxPeriod := max(u.PeriodMin, max(u.PeriodMid, u.PeriodMax))
	return u.valueNumber <= maxPeriod
}

func (u *UO) IdlePeriod() int {
	return max(u.PeriodMin, max(u.PeriodMid, u.PeriodMax))
}

func (u *UO) IsWarmedUp() bool {
	return !u.IsIdle()
}

func (u *UO) WarmUpPeriod() int {
	return u.IdlePeriod()
}
