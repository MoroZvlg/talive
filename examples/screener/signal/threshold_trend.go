package signal

import (
	"fmt"
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

// ThresholdTrendSignal generates buy/sell based on threshold + rising/falling trend.
// Used for RSI and CCI in TradingView's technical rating.
// RSI: buy if rsi < 30 AND rsi rising, sell if rsi > 70 AND rsi falling.
// CCI: buy if cci < -100 AND cci rising, sell if cci > 100 AND cci falling.
type ThresholdTrendSignal struct {
	indicator  talive.IIndicator
	buyTh      float64
	sellTh     float64
	lastResult float64
	hasLast    bool
}

func NewThresholdTrendSignal(indicator talive.IIndicator, buyTh, sellTh float64) *ThresholdTrendSignal {
	return &ThresholdTrendSignal{
		indicator: indicator,
		buyTh:     buyTh,
		sellTh:    sellTh,
	}
}

func (s *ThresholdTrendSignal) Next(kline *entity.Kline) int {
	result := s.indicator.Next(kline)
	value := result[0]
	signal := 0

	if s.indicator.IsWarmedUp() && s.hasLast {
		if value < s.buyTh && value > s.lastResult {
			signal = 1
		}
		if value > s.sellTh && value < s.lastResult {
			signal = -1
		}
	}

	fmt.Printf("[%s] %d (%.10f) prev=%.10f\n", s.indicator, signal, value, s.lastResult)
	s.lastResult = value
	s.hasLast = true
	return signal
}

func (s *ThresholdTrendSignal) MaxWarmUp() int {
	return s.indicator.WarmUpPeriod()
}
