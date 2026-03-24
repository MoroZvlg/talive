package talive

import (
	"fmt"
)

// MaType defines the type of moving average.
type MaType int

// Supported moving average types.
const (
	SMAtype MaType = iota
	EMAtype
	SMMAtype
	WMAtype
	// VWMAtype — requires volume, doesn't fit next(float64) interface

)

// MA is the common interface for moving average indicators.
type MA interface {
	IIndicator
	next(float64) float64
	current(float64) float64
}

// NewMa creates a moving average indicator of the given type.
func NewMa(period int, maType MaType) (MA, error) {
	switch maType {
	case SMAtype:
		return NewSMA(period)
	case EMAtype:
		return NewEMA(period)
	case SMMAtype:
		return NewSMMA(period)
	case WMAtype:
		return NewWMA(period)
	}
	return nil, fmt.Errorf("invalid ma type")
}
