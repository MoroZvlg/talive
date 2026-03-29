package signal

import (
	"fmt"
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

// ----------------
// ---    MA    ---
// ----------------

type MASignal struct {
	ma talive.IIndicator
}

func NewMASignal(ma talive.IIndicator) *MASignal {
	return &MASignal{ma: ma}
}

func (s *MASignal) Next(kline *entity.Kline) int {
	maV := s.ma.Next(kline)
	signal := 0
	if s.ma.IsWarmedUp() {
		if maV[0] > kline.Close() {
			signal = -1
		} else if maV[0] < kline.Close() {
			signal = 1
		}
	}
	fmt.Printf("[%s] %d (%f)\n", s.ma, signal, maV[0])
	return signal
}

func (s *MASignal) MaxWarmUp() int {
	return s.ma.WarmUpPeriod()
}

// ----------------
// --- MA Cross ---
// ----------------

type MACrossSignal struct {
	fast talive.MA
	slow talive.MA
}

func NewMACrossSignal(fast, slow talive.MA) *MACrossSignal {
	return &MACrossSignal{fast: fast, slow: slow}
}

func (s *MACrossSignal) Next(kline *entity.Kline) int {
	fast := s.fast.Next(kline)
	slow := s.slow.Next(kline)
	signal := 0
	if s.fast.IsWarmedUp() && s.slow.IsWarmedUp() {
		if fast[0] > slow[0] {
			signal = 1
		} else if fast[0] < slow[0] {
			signal = -1
		}
	}
	fmt.Printf("[%s/%s] %d (%f, %f)\n", s.fast, s.slow, signal, fast[0], slow[0])
	return signal
}

func (s *MACrossSignal) MaxWarmUp() int {
	return max(s.fast.WarmUpPeriod(), s.slow.WarmUpPeriod())
}
