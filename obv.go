package talive

// OBV is an On Balance Volume indicator.
type OBV struct {
	valueNumber int
	prevClose   float64
	value       float64
	out         []float64
}

// NewOBV creates a new OBV indicator.
func NewOBV() (*OBV, error) {
	return &OBV{
		out: make([]float64, 1),
	}, nil
}

func (o *OBV) String() string {
	return "OBV"
}

func (o *OBV) Next(candle OHLCV) []float64 {
	o.valueNumber++
	closeV := candle.Close()

	if o.valueNumber == 1 {
		o.prevClose = closeV
		o.out[0] = 0.0
		return o.out
	}

	switch {
	case closeV > o.prevClose:
		o.value += candle.Volume()
	case closeV < o.prevClose:
		o.value -= candle.Volume()
	}

	o.prevClose = closeV
	o.out[0] = o.value
	return o.out
}

func (o *OBV) Current(candle OHLCV) []float64 {
	o.valueNumber++
	if o.IsIdle() {
		o.valueNumber--
		o.out[0] = 0.0
		return o.out
	}

	closeV := candle.Close()
	value := o.value
	switch {
	case closeV > o.prevClose:
		value += candle.Volume()
	case closeV < o.prevClose:
		value -= candle.Volume()
	}
	o.valueNumber--
	o.out[0] = value
	return o.out
}

func (o *OBV) IsIdle() bool {
	return o.valueNumber <= 1
}

func (o *OBV) IdlePeriod() int {
	return 1
}

func (o *OBV) IsWarmedUp() bool {
	return !o.IsIdle()
}

func (o *OBV) WarmUpPeriod() int {
	return o.IdlePeriod()
}
