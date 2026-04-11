package signal

import (
	"fmt"
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

type StochasticSignal struct {
	stoch talive.Indicator
}

func NewStochasticSignal(stoch talive.Indicator) *StochasticSignal {
	return &StochasticSignal{
		stoch: stoch,
	}
}

func (s *StochasticSignal) Next(kline *entity.Kline) int {
	result := s.stoch.Next(kline)
	k := result[0]
	d := result[1]
	signal := 0
	if s.stoch.IsWarmedUp() {
		if k < 20 && d < 20 && k > d {
			signal = 1
		}
		if k > 80 && d > 80 && k < d {
			signal = -1
		}
	}
	fmt.Printf("[%s] %d (%.10f, %.10f)\n", s.stoch, signal, k, d)
	return signal
}

func (s *StochasticSignal) MaxWarmUp() int {
	return s.stoch.WarmUpPeriod()
}
