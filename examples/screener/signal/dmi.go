package signal

import (
	"fmt"
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

type DMISignal struct {
	dmi        *talive.DMI
	lastResult []float64
}

func NewDMISignal(dmi *talive.DMI) *DMISignal {
	return &DMISignal{
		dmi:        dmi,
		lastResult: make([]float64, 3),
	}
}

func (s *DMISignal) Next(kline *entity.Kline) int {
	result := s.dmi.Next(kline)
	adx := result[0]
	diPlus := result[1]
	diMinus := result[2]
	signal := 0

	if s.dmi.IsWarmedUp() {
		if diPlus > diMinus && adx > 20 && adx > s.lastResult[0] {
			signal = 1
		}
		if diPlus < diMinus && adx > 20 && adx > s.lastResult[0] {
			signal = -1
		}
	}
	fmt.Printf("[%s] %d (%f, %f, %f) <-> (%f, %f, %f)\n", s.dmi, signal, adx, diPlus, diMinus, s.lastResult[0], s.lastResult[1], s.lastResult[2])
	s.lastResult[0] = result[0]
	s.lastResult[1] = result[1]
	s.lastResult[2] = result[2]
	return signal
}

func (s *DMISignal) MaxWarmUp() int {
	return s.dmi.WarmUpPeriod()
}
