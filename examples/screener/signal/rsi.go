package signal

import (
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

type RSISignal struct {
	rsi           *talive.RSI
	buyThreshold  float64
	sellThreshold float64
}

func NewRSISignal(rsi *talive.RSI) *RSISignal {
	return &RSISignal{
		rsi:           rsi,
		buyThreshold:  30,
		sellThreshold: 70,
	}
}

func (s *RSISignal) Next(kline *entity.Kline) int {
	result := s.rsi.Next(kline)
	if s.rsi.IsWarmedUp() {
		if result[0] > s.sellThreshold {
			return -1
		} else if result[0] < s.buyThreshold {
			return 1
		} else {
			return 0
		}
	} else {
		return 0
	}
}

func (s *RSISignal) MaxWarmUp() int {
	return s.rsi.WarmUpPeriod()
}
