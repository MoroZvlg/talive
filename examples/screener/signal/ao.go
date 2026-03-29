package signal

import (
	"fmt"
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

// AOSignal generates buy/sell based on AO zero-line crossover and saucer pattern.
// Pine Script:
//
//	buy  = crossover(ao, 0) OR (ao > 0 AND ao[1] > 0 AND ao > ao[1] AND ao[2] > ao[1])
//	sell = crossunder(ao, 0) OR (ao < 0 AND ao[1] < 0 AND ao < ao[1] AND ao[2] < ao[1])
type AOSignal struct {
	indicator *talive.AO
	prev      float64 // ao[1]
	prevPrev  float64 // ao[2]
	count     int
}

func NewAOSignal(indicator *talive.AO) *AOSignal {
	return &AOSignal{
		indicator: indicator,
	}
}

func (s *AOSignal) Next(kline *entity.Kline) int {
	result := s.indicator.Next(kline)
	ao := result[0]
	signal := 0

	if s.indicator.IsWarmedUp() && s.count >= 2 {
		crossover := ao > 0 && s.prev <= 0
		crossunder := ao < 0 && s.prev >= 0
		saucerBull := ao > 0 && s.prev > 0 && ao > s.prev && s.prevPrev > s.prev
		saucerBear := ao < 0 && s.prev < 0 && ao < s.prev && s.prevPrev < s.prev

		if crossover || saucerBull {
			signal = 1
		}
		if crossunder || saucerBear {
			signal = -1
		}
	}

	fmt.Printf("[%s] %d (%f)\n", s.indicator, signal, ao)
	s.prevPrev = s.prev
	s.prev = ao
	s.count++
	return signal
}

func (s *AOSignal) MaxWarmUp() int {
	return s.indicator.WarmUpPeriod()
}
