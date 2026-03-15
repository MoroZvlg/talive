# talive

[![CI](https://github.com/MoroZvlg/talive/actions/workflows/ci.yml/badge.svg)](https://github.com/MoroZvlg/talive/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/MoroZvlg/talive)](https://goreportcard.com/report/github.com/MoroZvlg/talive)
[![codecov](https://codecov.io/gh/MoroZvlg/talive/branch/main/graph/badge.svg)](https://codecov.io/gh/MoroZvlg/talive)
[![Go Reference](https://pkg.go.dev/badge/github.com/MoroZvlg/talive.svg)](https://pkg.go.dev/github.com/MoroZvlg/talive)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Streaming technical analysis indicators for Go — zero-allocation, bar-by-bar.

Unlike batch TA libraries (ta-lib, go-talib) that recalculate over a full history slice on every update, `talive` maintains running state so each new candle costs **O(1)** regardless of history length. Built for live market feeds, algo trading bots, and any latency-sensitive system.

## Install

```bash
go get github.com/MoroZvlg/talive
```

## Quick Start

```go
import "github.com/MoroZvlg/talive"

// Create indicator
rsi, _ := talive.NewRSI(14)

// Feed candles one by one
for _, candle := range candles {
    values := rsi.Next(candle) // updates state, returns []float64
}

// Peek without updating state
values := rsi.Current(candle)

// Check if indicator has enough data
if rsi.IsIdle() {
    // still warming up
}
```

Candles must implement the `ICandle` interface:

```go
type ICandle interface {
    Open() float64
    High() float64
    Low() float64
    Close() float64
    Volume() float64
}
```

## Indicators

| Indicator | Constructor | Output |
|-----------|-------------|--------|
| EMA | `NewEMA(period)` | `[ema]` |
| SMA | `NewSMA(period)` | `[sma]` |
| RSI | `NewRSI(period)` | `[rsi]` |
| MACD | `NewMACD(fast, slow, signal)` | `[macd, signal, hist]` |
| Bollinger Bands | `NewBBands(period, upMult, downMult, maType)` | `[upper, mid, lower]` |
| MFI | `NewMFI(period)` | `[mfi]` |
| StdDev | `NewStdDev(period, ddof)` | `[stddev]` |
| Variance | `NewVariance(period)` | `[variance]` |

## Benchmarks

Measured on Apple M3 Pro:

```
BenchmarkEMANext      ~2.7 ns/op    0 B/op    0 allocs/op
BenchmarkSMANext      ~4.2 ns/op    0 B/op    0 allocs/op
BenchmarkRSINext      ~7.6 ns/op    0 B/op    0 allocs/op
BenchmarkMACDNext     ~6.8 ns/op    0 B/op    0 allocs/op
BenchmarkBBandsNext   ~9.3 ns/op    0 B/op    0 allocs/op
BenchmarkMFINext     ~11.0 ns/op    0 B/op    0 allocs/op
```

Zero allocations in the hot path — output slices are pre-allocated in constructors and reused on every call.

## Why talive

- **Streaming-first** — designed for live feeds, not batch recalculation
- **Zero-allocation hot path** — output slices pre-allocated in constructors, reused on every call
- **Simple interface** — `Next`, `Current`, `IsIdle` — that's it
- **No dependencies** — pure Go, nothing to break

## License

MIT
