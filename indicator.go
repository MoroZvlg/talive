// Package talive provides streaming technical analysis indicators with zero allocations.
package talive

import "fmt"

// OHLCV represents a single candlestick data point.
type OHLCV interface {
	Open() float64
	High() float64
	Low() float64
	Close() float64
	Volume() float64
}

// Indicator is the common interface for all technical indicators.
type Indicator interface {
	// Next feeds the next candle and advances the indicator state.
	// Returns zero values while IsIdle() is true.
	Next(candle OHLCV) []float64

	// Current calculates the indicator value for a candle without advancing state.
	// Returns zero values while IsIdle() is true.
	Current(candle OHLCV) []float64

	// IsIdle returns true while the indicator has not received enough candles
	// to produce meaningful output. All output values are zero during this phase.
	IsIdle() bool

	// IdlePeriod returns the number of candles that must be fed before the indicator
	// starts producing non-zero output.
	IdlePeriod() int

	// IsWarmedUp returns true when the indicator has received enough candles
	// for its output to be considered reliable. This requires more candles than
	// IdlePeriod for indicators with exponential memory (like EMA, RSI, MACD),
	// because early non-zero outputs still carry bias from limited history.
	// For indicators backed by fixed-size buffers (like SMA, MFI, StdDev),
	// IsWarmedUp is equivalent to !IsIdle. (see warmup_analysis_test.go).
	IsWarmedUp() bool

	// WarmUpPeriod returns the total number of candles that must be fed before
	// the indicator output is reliable. This value always includes IdlePeriod.
	WarmUpPeriod() int
}

// Scalar extends Indicator for single-value composable indicators.
// These operate on a raw float64 stream (pure math, no candle structure)
// and can be chained inside other indicators (e.g., RSI uses SMMA internally).
type Scalar interface {
	Indicator

	// NextVal feeds the next value and advances the indicator state.
	// Returns zero values while IsIdle() is true.
	NextVal(float64) float64

	// CurrentVal calculates the indicator value for a candle without advancing state.
	// Returns zero values while IsIdle() is true.
	CurrentVal(float64) float64
}

// SourceFunc selects which price value an indicator reads from a candle.
// Pass one of the predefined sources or a custom function to derive the price series.
type SourceFunc func(OHLCV) float64

// SourceClose returns Candle's Close price
func SourceClose(candle OHLCV) float64 { return candle.Close() }

// SourceOpen returns Candle's Open price
func SourceOpen(candle OHLCV) float64 { return candle.Open() }

// SourceHigh returns Candle's High price
func SourceHigh(candle OHLCV) float64 { return candle.High() }

// SourceLow returns Candle's Low price
func SourceLow(candle OHLCV) float64 { return candle.Low() }

// SourceHLC3 returns the candle's typical price (high + low + close) / 3.
func SourceHLC3(candle OHLCV) float64 { return (candle.High() + candle.Low() + candle.Close()) / 3 }

// SourceHL2 returns the candle's median price (high + low) / 2.
func SourceHL2(candle OHLCV) float64 { return (candle.High() + candle.Low()) / 2 }

// MaType defines the type of moving average.
type MaType int

func (mt MaType) String() string {
	switch mt {
	case UseSMA:
		return "SMA"
	case UseEMA:
		return "EMA"
	case UseSMMA:
		return "SMMA"
	case UseWMA:
		return "WMA"
	}
	return "UnknownMA"
}

// New creates a new Scalar moving average of this type with the given period.
func (mt MaType) New(period int) (Scalar, error) {
	switch mt {
	case UseSMA:
		return NewSMA(period)
	case UseEMA:
		return NewEMA(period)
	case UseSMMA:
		return NewSMMA(period)
	case UseWMA:
		return NewWMA(period)
	}
	return nil, fmt.Errorf("unknown MA type: %s", mt)
}

// Supported moving average types.
const (
	UseSMA MaType = iota
	UseEMA
	UseSMMA
	UseWMA
)
