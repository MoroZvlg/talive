package talive

type ICandle interface {
	Open() float64
	High() float64
	Low() float64
	Close() float64
	Volume() float64
}

// IIndicator is the common interface for all technical indicators.
type IIndicator interface {
	// Next feeds the next candle and advances the indicator state.
	// Returns zero values while IsIdle() is true.
	Next(candle ICandle) []float64

	// Current calculates the indicator value for a candle without advancing state.
	// Returns zero values while IsIdle() is true.
	Current(candle ICandle) []float64

	// IsIdle returns true while the indicator has not received enough candles
	// to produce meaningful output. All output values are zero during this phase.
	IsIdle() bool

	// IdlePeriod returns the number of candles that must be fed before the indicator
	// starts producing non-zero output.
	IdlePeriod() uint

	// IsWarmedUp returns true when the indicator has received enough candles
	// for its output to be considered reliable. This requires more candles than
	// IdlePeriod for indicators with exponential memory (like EMA, RSI, MACD),
	// because early non-zero outputs still carry bias from limited history.
	// For indicators backed by fixed-size buffers (like SMA, MFI, StdDev),
	// IsWarmedUp is equivalent to !IsIdle. (see warmup_analysis_test.go).
	IsWarmedUp() bool

	// WarmUpPeriod returns the total number of candles that must be fed before
	// the indicator output is reliable. This value always includes IdlePeriod.
	WarmUpPeriod() uint
}
