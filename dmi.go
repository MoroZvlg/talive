package talive

import (
	"fmt"
	"math"
)

// DMI is a Directional Movement Index indicator.
// It returns ADX, +DI and -DI values.
type DMI struct {
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

// NewDMI creates a new DMI indicator with the given periodmi.
func NewDMI(period int) (*DMI, error) {
	plusDMSmma, _ := NewSMMA(period)
	minusDMSmma, _ := NewSMMA(period)
	trSmma, _ := NewSMMA(period)
	adxSmma, _ := NewSMMA(period)
	return &DMI{
		Period:      period,
		plusDMSmma:  plusDMSmma,
		minusDMSmma: minusDMSmma,
		trSmma:      trSmma,
		adxSmma:     adxSmma,
		out:         make([]float64, 3),
	}, nil
}

func (dmi *DMI) String() string {
	return fmt.Sprintf("DMI(%d)", dmi.Period)
}

func (dmi *DMI) Next(candle ICandle) []float64 {
	dmi.valueNumber++

	if dmi.valueNumber == 1 {
		dmi.prevHigh = candle.High()
		dmi.prevLow = candle.Low()
		dmi.prevClose = candle.Close()
		return dmi.out
	}

	plusDM, minusDM, tr := dmi.computeDMTR(candle)

	dmi.prevHigh = candle.High()
	dmi.prevLow = candle.Low()
	dmi.prevClose = candle.Close()

	sPlusDM := dmi.plusDMSmma.next(plusDM)
	sMinusDM := dmi.minusDMSmma.next(minusDM)
	sTR := dmi.trSmma.next(tr)

	if dmi.trSmma.IsIdle() {
		return dmi.out
	}

	plusDI := 100 * sPlusDM / sTR
	minusDI := 100 * sMinusDM / sTR

	dmi.out[1] = plusDI
	dmi.out[2] = minusDI

	dx := 100 * math.Abs(plusDI-minusDI) / (plusDI + minusDI)
	adx := dmi.adxSmma.next(dx)

	if dmi.adxSmma.IsIdle() {
		return dmi.out
	}

	dmi.out[0] = adx
	return dmi.out
}

func (dmi *DMI) Current(candle ICandle) []float64 {
	if dmi.IsIdle() {
		return dmi.out
	}

	plusDM, minusDM, tr := dmi.computeDMTR(candle)

	sPlusDM := dmi.plusDMSmma.current(plusDM)
	sMinusDM := dmi.minusDMSmma.current(minusDM)
	sTR := dmi.trSmma.current(tr)

	plusDI := 100 * sPlusDM / sTR
	minusDI := 100 * sMinusDM / sTR

	dmi.out[1] = plusDI
	dmi.out[2] = minusDI

	dx := 100 * math.Abs(plusDI-minusDI) / (plusDI + minusDI)
	dmi.out[0] = dmi.adxSmma.current(dx)
	return dmi.out
}

func (dmi *DMI) computeDMTR(candle ICandle) (plusDM, minusDM, tr float64) {
	upMove := candle.High() - dmi.prevHigh
	downMove := dmi.prevLow - candle.Low()

	if upMove > downMove && upMove > 0 {
		plusDM = upMove
	}
	if downMove > upMove && downMove > 0 {
		minusDM = downMove
	}

	highLow := candle.High() - candle.Low()
	highPrevClose := math.Abs(candle.High() - dmi.prevClose)
	lowPrevClose := math.Abs(candle.Low() - dmi.prevClose)
	tr = max(highLow, max(highPrevClose, lowPrevClose))
	return plusDM, minusDM, tr
}

func (dmi *DMI) IsIdle() bool {
	return dmi.valueNumber <= dmi.Period
}

func (dmi *DMI) IdlePeriod() int {
	return dmi.Period
}

func (dmi *DMI) IsWarmedUp() bool {
	return dmi.valueNumber > dmi.WarmUpPeriod()
}

func (dmi *DMI) WarmUpPeriod() int {
	return 2*dmi.Period - 1 + dmi.Period*9
}
