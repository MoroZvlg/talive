package talive

import (
	"fmt"
	"math"
)

// ADX is an Average Directional Movement Index indicator.
type ADX struct {
	Period      int
	valueNumber int

	prevHigh  float64
	prevLow   float64
	prevClose float64

	plusDMSmma  MA
	minusDMSmma MA
	trSmma      MA
	adxSmma     MA

	out []float64
}

// NewADX creates a new ADX indicator with the given period.
func NewADX(period int) (*ADX, error) {
	plusDMSmma, _ := NewSMMA(period)
	minusDMSmma, _ := NewSMMA(period)
	trSmma, _ := NewSMMA(period)
	adxSmma, _ := NewSMMA(period)
	return &ADX{
		Period:      period,
		plusDMSmma:  plusDMSmma,
		minusDMSmma: minusDMSmma,
		trSmma:      trSmma,
		adxSmma:     adxSmma,
		out:         make([]float64, 1),
	}, nil
}

func (adx *ADX) String() string {
	return fmt.Sprintf("ADX(%d)", adx.Period)
}

func (adx *ADX) Next(candle ICandle) []float64 {
	adx.valueNumber++

	if adx.valueNumber == 1 {
		adx.prevHigh = candle.High()
		adx.prevLow = candle.Low()
		adx.prevClose = candle.Close()
		return adx.out
	}

	plusDM, minusDM, tr := adx.computeDMTR(candle)

	adx.prevHigh = candle.High()
	adx.prevLow = candle.Low()
	adx.prevClose = candle.Close()

	sPlusDM := adx.plusDMSmma.next(plusDM)
	sMinusDM := adx.minusDMSmma.next(minusDM)
	sTR := adx.trSmma.next(tr)

	if adx.trSmma.IsIdle() {
		return adx.out
	}

	plusDI := 100 * sPlusDM / sTR
	minusDI := 100 * sMinusDM / sTR

	dx := 100 * math.Abs(plusDI-minusDI) / (plusDI + minusDI)
	adxV := adx.adxSmma.next(dx)

	if adx.adxSmma.IsIdle() {
		return adx.out
	}

	adx.out[0] = adxV
	return adx.out
}

func (adx *ADX) Current(candle ICandle) []float64 {
	if adx.IsIdle() {
		return adx.out
	}

	plusDM, minusDM, tr := adx.computeDMTR(candle)

	sPlusDM := adx.plusDMSmma.current(plusDM)
	sMinusDM := adx.minusDMSmma.current(minusDM)
	sTR := adx.trSmma.current(tr)

	plusDI := 100 * sPlusDM / sTR
	minusDI := 100 * sMinusDM / sTR

	dx := 100 * math.Abs(plusDI-minusDI) / (plusDI + minusDI)
	adx.out[0] = adx.adxSmma.current(dx)
	return adx.out
}

func (adx *ADX) computeDMTR(candle ICandle) (plusDM, minusDM, tr float64) {
	upMove := candle.High() - adx.prevHigh
	downMove := adx.prevLow - candle.Low()

	if upMove > downMove && upMove > 0 {
		plusDM = upMove
	}
	if downMove > upMove && downMove > 0 {
		minusDM = downMove
	}

	highLow := candle.High() - candle.Low()
	highPrevClose := math.Abs(candle.High() - adx.prevClose)
	lowPrevClose := math.Abs(candle.Low() - adx.prevClose)
	tr = max(highLow, max(highPrevClose, lowPrevClose))
	return plusDM, minusDM, tr
}

func (adx *ADX) IsIdle() bool {
	return adx.valueNumber < 2*adx.Period
}

func (adx *ADX) IdlePeriod() int {
	return 2*adx.Period - 1
}

func (adx *ADX) IsWarmedUp() bool {
	return adx.valueNumber > adx.WarmUpPeriod()
}

func (adx *ADX) WarmUpPeriod() int {
	return adx.IdlePeriod() + adx.Period*9
}
