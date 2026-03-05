package talive

import (
	"fmt"
)

type MaType int

const (
	SMAtype MaType = iota
	EMAtype
)

type MA interface {
	IIndicator
	next(float64) float64
	current(float64) float64
}

func NewMa(period int, maType MaType) (MA, error) {
	switch maType {
	case SMAtype:
		return NewSMA(period)
	case EMAtype:
		return NewEMA(period)
	}
	return nil, fmt.Errorf("invalid ma type")
}
