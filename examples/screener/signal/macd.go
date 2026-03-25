package signal

import (
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

type MACDSignal struct {
	macd *talive.MACD
}

func NewMACDSignal(macd *talive.MACD) *MACDSignal {
	return &MACDSignal{
		macd: macd,
	}
}

func (s *MACDSignal) Next(kline *entity.Kline) int {
	result := s.macd.Next(kline)
	macdLine := result[0]
	signalLine := result[1]
	if s.macd.IsWarmedUp() {
		if macdLine > signalLine {
			return 1
		} else if macdLine < signalLine {
			return -1
		}
	}
	return 0
}

func (s *MACDSignal) MaxWarmUp() int {
	return s.macd.WarmUpPeriod()
}
