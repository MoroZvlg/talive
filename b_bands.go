package talive

// BBands is a Bollinger Bands indicator.
type BBands struct {
	Period         int
	DevUp, DevDown float64
	MaType         MaType
	valueNumber    int
	ma             MA
	basicDeviation *StdDev
	out            []float64
}

// NewBBands creates a new Bollinger Bands indicator with the given parameters.
func NewBBands(period int, devUp, devDown float64, maType MaType) (*BBands, error) {
	ma, err := NewMa(period, maType)
	if err != nil {
		return nil, err
	}
	basicDeviation, err := NewStdDev(period, 1.0)
	if err != nil {
		return nil, err
	}
	return &BBands{
		Period:         period,
		DevUp:          devUp,
		DevDown:        devDown,
		MaType:         maType,
		valueNumber:    0,
		ma:             ma,
		basicDeviation: basicDeviation,
		out:            make([]float64, 3),
	}, nil
}

func (bb *BBands) Next(candle ICandle) []float64 {
	value := candle.Close()
	bb.valueNumber++
	ma := bb.ma.next(value)
	devBase := bb.basicDeviation.next(value)

	if bb.IsIdle() {
		bb.out[0] = 0.0
		bb.out[1] = 0.0
		bb.out[2] = 0.0
		return bb.out
	}

	bb.out[0] = ma + (devBase * bb.DevUp)
	bb.out[1] = ma
	bb.out[2] = ma - (devBase * bb.DevDown)
	return bb.out
}

func (bb *BBands) Current(candle ICandle) []float64 {
	value := candle.Close()
	bb.valueNumber++
	ma := bb.ma.current(value)
	devBase := bb.basicDeviation.current(value)

	if bb.IsIdle() {
		bb.valueNumber--
		bb.out[0] = 0.0
		bb.out[1] = 0.0
		bb.out[2] = 0.0
		return bb.out
	}

	bb.out[0] = ma + (devBase * bb.DevUp)
	bb.out[1] = ma
	bb.out[2] = ma - (devBase * bb.DevDown)
	bb.valueNumber--
	return bb.out
}

func (bb *BBands) IsIdle() bool {
	return bb.ma.IsIdle()
}

func (bb *BBands) IdlePeriod() int {
	return bb.ma.IdlePeriod()
}

func (bb *BBands) IsWarmedUp() bool {
	return bb.ma.IsWarmedUp()
}

func (bb *BBands) WarmUpPeriod() int {
	return bb.ma.WarmUpPeriod()
}
