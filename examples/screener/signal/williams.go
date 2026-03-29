package signal

import (
	"fmt"
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

// WilliamsSignal generates buy/sell based on Williams %R with trend confirmation.
// Pine Script:
//
//	buy  = wr < -80 AND wr > wr[1] (oversold and rising)
//	sell = wr > -20 AND wr < wr[1] (overbought and falling)
type WilliamsSignal struct {
	indicator  *talive.Williams
	lastResult float64
	hasLast    bool
}

func NewWilliamsSignal(will *talive.Williams) *WilliamsSignal {
	return &WilliamsSignal{
		indicator: will,
	}
}

func (s *WilliamsSignal) Next(kline *entity.Kline) int {
	result := s.indicator.Next(kline)
	value := result[0]
	signal := 0

	if s.indicator.IsWarmedUp() && s.hasLast {
		if value < -80 && value > s.lastResult {
			signal = 1
		}
		if value > -20 && value < s.lastResult {
			signal = -1
		}
	}

	fmt.Printf("[%s] %d (%f)\n", s.indicator, signal, value)
	s.lastResult = value
	s.hasLast = true
	return signal
}

func (s *WilliamsSignal) MaxWarmUp() int {
	return s.indicator.WarmUpPeriod()
}
