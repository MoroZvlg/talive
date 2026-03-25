package talive

import "math"

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

func (a *ADX) Next(candle ICandle) []float64 {
	a.valueNumber++

	if a.valueNumber == 1 {
		a.prevHigh = candle.High()
		a.prevLow = candle.Low()
		a.prevClose = candle.Close()
		return a.out
	}

	plusDM, minusDM, tr := a.computeDMTR(candle)

	a.prevHigh = candle.High()
	a.prevLow = candle.Low()
	a.prevClose = candle.Close()

	sPlusDM := a.plusDMSmma.next(plusDM)
	sMinusDM := a.minusDMSmma.next(minusDM)
	sTR := a.trSmma.next(tr)

	if a.trSmma.IsIdle() {
		return a.out
	}

	plusDI := 100 * sPlusDM / sTR
	minusDI := 100 * sMinusDM / sTR

	dx := 100 * math.Abs(plusDI-minusDI) / (plusDI + minusDI)
	adx := a.adxSmma.next(dx)

	if a.adxSmma.IsIdle() {
		return a.out
	}

	a.out[0] = adx
	return a.out
}

func (a *ADX) Current(candle ICandle) []float64 {
	if a.IsIdle() {
		return a.out
	}

	plusDM, minusDM, tr := a.computeDMTR(candle)

	sPlusDM := a.plusDMSmma.current(plusDM)
	sMinusDM := a.minusDMSmma.current(minusDM)
	sTR := a.trSmma.current(tr)

	plusDI := 100 * sPlusDM / sTR
	minusDI := 100 * sMinusDM / sTR

	dx := 100 * math.Abs(plusDI-minusDI) / (plusDI + minusDI)
	a.out[0] = a.adxSmma.current(dx)
	return a.out
}

func (a *ADX) computeDMTR(candle ICandle) (plusDM, minusDM, tr float64) {
	upMove := candle.High() - a.prevHigh
	downMove := a.prevLow - candle.Low()

	if upMove > downMove && upMove > 0 {
		plusDM = upMove
	}
	if downMove > upMove && downMove > 0 {
		minusDM = downMove
	}

	highLow := candle.High() - candle.Low()
	highPrevClose := math.Abs(candle.High() - a.prevClose)
	lowPrevClose := math.Abs(candle.Low() - a.prevClose)
	tr = max(highLow, max(highPrevClose, lowPrevClose))
	return plusDM, minusDM, tr
}

func (a *ADX) IsIdle() bool {
	return a.valueNumber < 2*a.Period
}

func (a *ADX) IdlePeriod() int {
	return 2*a.Period - 1
}

func (a *ADX) IsWarmedUp() bool {
	return a.valueNumber > a.WarmUpPeriod()
}

func (a *ADX) WarmUpPeriod() int {
	return a.IdlePeriod() + a.Period*9
}
