package signal

import (
	"fmt"
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

// BullBearPowerSignal generates buy/sell using separate bull/bear power with EMA trend filter.
// Pine Script:
//
//	bullPower = high - ema(close, 13)
//	bearPower = low  - ema(close, 13)
//	buy  = close > ema50 (uptrend)   AND bearPower < 0 AND bearPower > bearPower[1]
//	sell = close < ema50 (downtrend) AND bullPower > 0 AND bullPower < bullPower[1]
type BullBearPowerSignal struct {
	ema      talive.MA
	trendMA  talive.MA
	lastBull float64
	lastBear float64
	hasLast  bool
}

func NewBullBearPowerSignal(ema talive.MA, trendMA talive.MA) *BullBearPowerSignal {
	return &BullBearPowerSignal{
		ema:     ema,
		trendMA: trendMA,
	}
}

func (s *BullBearPowerSignal) Next(kline *entity.Kline) int {
	emaResult := s.ema.Next(kline)
	trendResult := s.trendMA.Next(kline)
	emaVal := emaResult[0]

	bullPower := kline.High() - emaVal
	bearPower := kline.Low() - emaVal

	signal := 0
	if s.ema.IsWarmedUp() && s.trendMA.IsWarmedUp() && s.hasLast {
		priceAvg := trendResult[0]
		upTrend := kline.Close() > priceAvg
		downTrend := kline.Close() < priceAvg

		if upTrend && bearPower < 0 && bearPower > s.lastBear {
			signal = 1
		}
		if downTrend && bullPower > 0 && bullPower < s.lastBull {
			signal = -1
		}
	}

	fmt.Printf("[BullBearPower(13)] %d (%f, %f)\n", signal, bullPower, bearPower)
	s.lastBull = bullPower
	s.lastBear = bearPower
	s.hasLast = true
	return signal
}

func (s *BullBearPowerSignal) MaxWarmUp() int {
	return max(s.ema.WarmUpPeriod(), s.trendMA.WarmUpPeriod())
}
