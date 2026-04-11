package signal

import (
	"fmt"
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

// StochRSISignal generates buy/sell based on StochRSI with EMA trend filter.
// Pine Script:
//
//	buy  = close < ema50 (downtrend) AND kStochRsi < 20 AND dStochRsi < 20 AND k > d
//	sell = close > ema50 (uptrend)   AND kStochRsi > 80 AND dStochRsi > 80 AND k < d
type StochRSISignal struct {
	stochRSI talive.Indicator
	trendMA  talive.Scalar
}

func NewStochRSISignal(stochRSI talive.Indicator, trendMA talive.Scalar) *StochRSISignal {
	return &StochRSISignal{
		stochRSI: stochRSI,
		trendMA:  trendMA,
	}
}

func (s *StochRSISignal) Next(kline *entity.Kline) int {
	result := s.stochRSI.Next(kline)
	trendResult := s.trendMA.Next(kline)
	k := result[0]
	d := result[1]
	signal := 0

	if s.stochRSI.IsWarmedUp() && s.trendMA.IsWarmedUp() {
		priceAvg := trendResult[0]
		downTrend := kline.Close() < priceAvg
		upTrend := kline.Close() > priceAvg

		if downTrend && k < 20 && d < 20 && k > d {
			signal = 1
		}
		if upTrend && k > 80 && d > 80 && k < d {
			signal = -1
		}
	}

	fmt.Printf("[%s] %d (%.10f, %.10f)\n", s.stochRSI, signal, k, d)
	return signal
}

func (s *StochRSISignal) MaxWarmUp() int {
	return max(s.stochRSI.WarmUpPeriod(), s.trendMA.WarmUpPeriod())
}
