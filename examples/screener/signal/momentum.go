package signal

import (
	"fmt"
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

// MomentumSignal generates buy/sell based on momentum direction.
// Pine Script: buy if mom > mom[1], sell if mom < mom[1].
type MomentumSignal struct {
	indicator  talive.Indicator
	lastResult float64
	hasLast    bool
}

func NewMomentumSignal(indicator talive.Indicator) *MomentumSignal {
	return &MomentumSignal{
		indicator: indicator,
	}
}

func (s *MomentumSignal) Next(kline *entity.Kline) int {
	result := s.indicator.Next(kline)
	value := result[0]
	signal := 0

	if s.indicator.IsWarmedUp() && s.hasLast {
		if value > s.lastResult {
			signal = 1
		} else if value < s.lastResult {
			signal = -1
		}
	}

	fmt.Printf("[%s] %d (%.10f) prev=%.10f\n", s.indicator, signal, value, s.lastResult)
	s.lastResult = value
	s.hasLast = true
	return signal
}

func (s *MomentumSignal) MaxWarmUp() int {
	return s.indicator.WarmUpPeriod()
}
