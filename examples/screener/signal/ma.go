package signal

import (
	"screener/domain/entity"

	"github.com/MoroZvlg/talive"
)

// ----------------
// ---    MA    ---
// ----------------

type MASignal struct {
	ma talive.MA
}

func NewMASignal(ma talive.MA) *MASignal {
	return &MASignal{ma: ma}
}

func (s *MASignal) Next(kline *entity.Kline) int {
	maV := s.ma.Next(kline)
	if s.ma.IsWarmedUp() {
		if maV[0] > kline.Close() {
			return -1
		} else if maV[0] < kline.Close() {
			return 1
		}
	}
	return 0
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
	var result int
	if s.fast.IsWarmedUp() && s.slow.IsWarmedUp() {
		if fast[0] > slow[0] {
			result = 1
		} else if fast[0] < slow[0] {
			result = -1
		} else {
			result = 0
		}
	} else {
		result = 0
	}
	return result
}

func (s *MACrossSignal) MaxWarmUp() int {
	return max(s.fast.WarmUpPeriod(), s.slow.WarmUpPeriod())
}
