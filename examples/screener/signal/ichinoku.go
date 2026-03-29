package signal

import (
	"fmt"
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

type IchimokuSignal struct {
	ich *talive.Ichimoku
}

func NewIchimokuSignal(ich *talive.Ichimoku) *IchimokuSignal {
	return &IchimokuSignal{
		ich: ich,
	}
}

func (s *IchimokuSignal) Next(kline *entity.Kline) int {
	result := s.ich.Next(kline)
	conv := result[0]
	base := result[1]
	leadA := result[2]
	leadB := result[3]
	signal := 0
	if s.ich.IsWarmedUp() {
		if leadA > leadB && base > leadA && conv > base && kline.Close() > conv {
			signal = 1
		}
		if leadA < leadB && base < leadA && conv < base && kline.Close() < conv {
			signal = -1
		}
	}
	fmt.Printf("[%s] %d (%f, %f, %f, %f)\n", s.ich, signal, conv, base, leadA, leadB)
	return signal
}

func (s *IchimokuSignal) MaxWarmUp() int {
	return s.ich.WarmUpPeriod()
}
