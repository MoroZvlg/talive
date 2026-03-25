package talive

import "math"

// SAR is a Parabolic SAR indicator.
type SAR struct {
	AfStart     float64
	AfIncrement float64
	AfMax       float64

	valueNumber  int
	inUpTrend    bool
	prevAf       float64 // acceleration factor
	prevEp       float64 // extreme price
	prevSar      float64
	prevHigh     float64
	prevLow      float64
	prevPrevHigh float64
	prevPrevLow  float64
	prevClose    float64
	out          []float64
}

// NewSAR creates a new Parabolic SAR indicator.
func NewSAR(start, increment, maxAF float64) *SAR {
	return &SAR{
		AfStart:     start,
		AfIncrement: increment,
		AfMax:       maxAF,
		out:         make([]float64, 1),
	}
}

func (s *SAR) Next(candle ICandle) []float64 {
	s.valueNumber++

	if s.valueNumber == 1 {
		s.prevHigh = candle.High()
		s.prevLow = candle.Low()
		s.prevClose = candle.Close()

		return s.out
	}

	var currSar float64
	ep := s.prevEp
	af := s.prevAf

	if s.valueNumber == 2 {
		if candle.Close() > s.prevClose {
			s.inUpTrend = true
			ep = math.Max(s.prevHigh, candle.High())
			currSar = math.Min(s.prevLow, candle.Low())
		} else {
			s.inUpTrend = false
			ep = math.Min(s.prevLow, candle.Low())
			currSar = math.Max(s.prevHigh, candle.High())
		}
		af = s.AfStart
		s.out[0] = currSar

		s.prevPrevHigh = s.prevHigh
		s.prevPrevLow = s.prevLow
		s.prevEp = ep
		s.prevSar = currSar
		s.prevAf = af
		s.prevHigh = candle.High()
		s.prevLow = candle.Low()
		s.prevClose = candle.Close()

		return s.out
	}

	currSar = s.prevSar + s.prevAf*(s.prevEp-s.prevSar)

	if s.inUpTrend {
		if candle.High() > s.prevEp {
			ep = candle.High()
			af = math.Min(s.prevAf+s.AfIncrement, s.AfMax)
		}
		if currSar > candle.Low() {
			s.inUpTrend = false
			currSar = ep
			ep = candle.Low()
			af = s.AfStart
		} else {
			currSar = math.Min(currSar, math.Min(s.prevLow, s.prevPrevLow))
		}
	} else {
		if candle.Low() < s.prevEp {
			ep = candle.Low()
			af = math.Min(s.prevAf+s.AfIncrement, s.AfMax)
		}
		if currSar < candle.High() {
			s.inUpTrend = true
			currSar = ep
			ep = candle.High()
			af = s.AfStart
		} else {
			currSar = math.Max(currSar, math.Max(s.prevHigh, s.prevPrevHigh))
		}
	}

	s.out[0] = currSar

	s.prevSar = currSar
	s.prevAf = af
	s.prevEp = ep
	s.prevPrevHigh = s.prevHigh
	s.prevPrevLow = s.prevLow
	s.prevHigh = candle.High()
	s.prevLow = candle.Low()
	s.prevClose = candle.Close()

	return s.out
}

func (s *SAR) Current(candle ICandle) []float64 {
	if s.IsIdle() {
		return s.out
	}

	currSar := s.prevSar + s.prevAf*(s.prevEp-s.prevSar)
	ep := s.prevEp

	if s.inUpTrend {
		if candle.High() > ep {
			ep = candle.High()
		}
		if currSar > candle.Low() {
			currSar = ep
		} else {
			currSar = math.Min(currSar, math.Min(s.prevLow, s.prevPrevLow))
		}
	} else {
		if candle.Low() < ep {
			ep = candle.Low()
		}
		if currSar < candle.High() {
			currSar = ep
		} else {
			currSar = math.Max(currSar, math.Max(s.prevHigh, s.prevPrevHigh))
		}
	}

	s.out[0] = currSar
	return s.out
}

func (s *SAR) IsIdle() bool {
	return s.valueNumber <= s.IdlePeriod()
}

func (s *SAR) IdlePeriod() int {
	return 1
}

func (s *SAR) IsWarmedUp() bool {
	return s.valueNumber > s.WarmUpPeriod()
}

// WarmUpPeriod returns a conservative warmup estimate.
// SAR warmup depends on data (trend length), not on formula parameters.
// Bigger AfStart and AfIncrement reduce warmup (faster trend flips), but the dependency is not linear.
func (s *SAR) WarmUpPeriod() int {
	return 100
}
