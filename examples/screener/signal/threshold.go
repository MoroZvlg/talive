package signal

import (
	"fmt"
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

type ThresholdSignal struct {
	indicator     talive.Indicator
	buyThreshold  float64
	sellThreshold float64
	reverse       int
}

func NewThresholdSignal(indicator talive.Indicator, buyTh, sellTh float64, reverse bool) *ThresholdSignal {
	var r int
	if reverse {
		r = -1
	} else {
		r = 1
	}
	return &ThresholdSignal{
		indicator:     indicator,
		buyThreshold:  buyTh,
		sellThreshold: sellTh,
		reverse:       r,
	}
}

func (s *ThresholdSignal) Next(kline *entity.Kline) int {
	result := s.indicator.Next(kline)
	signal := 0
	if s.indicator.IsWarmedUp() {
		if result[0] > s.sellThreshold {
			signal = -1 * s.reverse
		}
		if result[0] < s.buyThreshold {
			signal = 1 * s.reverse
		}
	}
	fmt.Printf("[%s] %d (%f)\n", s.indicator, signal, result[0])
	return signal
}

func (s *ThresholdSignal) MaxWarmUp() int {
	return s.indicator.WarmUpPeriod()
}
