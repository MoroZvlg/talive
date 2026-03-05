package talive

type ICandle interface {
	Open() float64
	High() float64
	Low() float64
	Close() float64
	Volume() float64
}

type IIndicator interface {
	Next(candle ICandle) []float64
	Current(candle ICandle) []float64
	IsIdle() bool
}
