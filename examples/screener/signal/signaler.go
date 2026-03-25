package signal

import "screener/domain/entity"

type Signaler interface {
	Next(kline *entity.Kline) int // buy = 1, sell = -1, hold = 0
	MaxWarmUp() int
}
