package signal

import (
	"fmt"
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
	signal := 0
	if s.macd.IsWarmedUp() {
		if macdLine > signalLine {
			signal = 1
		} else if macdLine < signalLine {
			signal = -1
		}
	}
	fmt.Printf("[%s] %d (%f, %f)\n", s.macd, signal, macdLine, signalLine)
	return signal
}

func (s *MACDSignal) MaxWarmUp() int {
	return s.macd.WarmUpPeriod()
}
