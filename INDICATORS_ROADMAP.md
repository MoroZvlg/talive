# Indicators Roadmap

## Already Implemented

| # | Indicator | Category | Usage |
|---|-----------|----------|-------|
| 1 | SMA | Moving Average | 10/10 |
| 2 | EMA | Moving Average | 10/10 |
| 3 | WMA | Moving Average | 6/10 |
| 4 | SMMA (Wilder's) | Moving Average | 5/10 |
| 5 | HMA | Moving Average | 7/10 |
| 6 | VWMA | Moving Average | 7/10 |
| 7 | RSI | Oscillator | 10/10 |
| 8 | MACD | Trend/Momentum | 10/10 |
| 9 | Bollinger Bands | Volatility | 10/10 |
| 10 | Stochastic | Oscillator | 9/10 |
| 11 | Stochastic RSI | Oscillator | 8/10 |
| 12 | ATR | Volatility | 9/10 |
| 13 | ADX | Trend | 8/10 |
| 14 | DMI (+DI/-DI) | Trend | 8/10 |
| 15 | CCI | Oscillator | 7/10 |
| 16 | Williams %R | Oscillator | 7/10 |
| 17 | MFI | Volume/Oscillator | 7/10 |
| 18 | Ichimoku Cloud | Trend/Complex | 8/10 |
| 19 | Parabolic SAR | Trend | 8/10 |
| 20 | AO (Awesome Osc) | Momentum | 6/10 |
| 21 | UO (Ultimate Osc) | Momentum | 5/10 |
| 22 | Momentum | Momentum | 6/10 |
| 23 | Bull Bear Power | Momentum | 5/10 |

## To Build

Sorted by priority. Usage is approximate popularity among retail/algo traders.

| Priority | Indicator | Category | Usage | Comment |
|----------|-----------|----------|-------|---------|
| 1 | OBV (On Balance Volume) | Volume | 9/10 | Most popular volume indicator. Cumulative volume flow, confirms trends. Every screener has it |
| 2 | VWAP | Volume | 9/10 | Standard for intraday trading. Institutional benchmark price. Resets daily |
| 3 | Supertrend | Trend | 9/10 | ATR-based trend follower. Huge in crypto/forex. Simple buy/sell signals |
| 4 | Keltner Channel | Volatility | 8/10 | EMA + ATR channel. The "TTM Squeeze" (Keltner inside BBands) is a famous setup |
| 5 | DEMA | Moving Average | 7/10 | Double EMA, reduces lag. Popular in crossover strategies |
| 6 | TEMA | Moving Average | 7/10 | Triple EMA, even less lag. Often paired with DEMA |
| 7 | Aroon | Trend | 7/10 | Measures time since highest high / lowest low. Good for detecting new trends |
| 8 | ROC (Rate of Change) | Momentum | 7/10 | Percentage price change over N periods. Simple, widely used in screening |
| 9 | Pivot Points | Support/Resistance | 8/10 | Standard, Fibonacci, Camarilla, Woodie. Floor traders' tool, still very relevant |
| 10 | Donchian Channel | Volatility | 7/10 | Highest high / lowest low channel. Turtle Traders made it famous. Basis for Supertrend-like strategies |
| 11 | CMF (Chaikin Money Flow) | Volume | 7/10 | Money flow over N periods. Shows buying/selling pressure. Common screener filter |
| 12 | ADL (Accumulation/Distribution) | Volume | 7/10 | Cumulative volume-price indicator. Foundation for CMF. Divergence with price = strong signal |
| 13 | TRIX | Momentum | 6/10 | Triple-smoothed EMA rate of change. Good noise filter, used for signal line crossovers |
| 14 | EOM (Ease of Movement) | Volume | 5/10 | Relates price change to volume. Useful for detecting low-volume breakouts |
| 15 | PPO (Percentage Price Osc) | Momentum | 6/10 | Normalized MACD (percentage). Allows comparison across different price levels |
| 16 | DPO (Detrended Price Osc) | Momentum | 5/10 | Removes trend to isolate cycles. Niche but used in cycle analysis |
| 17 | Chande Momentum Osc (CMO) | Momentum | 5/10 | Like RSI but unsmoothed. Used in adaptive MA systems (Chande's VIDYA) |
| 18 | VIDYA | Moving Average | 5/10 | Volatility-adaptive MA using CMO. Adjusts speed based on market conditions |
| 19 | KAMA (Kaufman Adaptive MA) | Moving Average | 6/10 | Adapts to noise level. Less whipsaw than EMA in ranging markets |
| 20 | Coppock Curve | Momentum | 4/10 | Long-term momentum for monthly charts. Originally for S&P 500 bottom detection |
| 21 | KST (Know Sure Thing) | Momentum | 5/10 | Weighted sum of 4 ROC smoothed values. Martin Pring's trend confirmation tool |
| 22 | Mass Index | Volatility | 4/10 | Detects trend reversals via range narrowing/widening ("reversal bulge") |
| 23 | TSI (True Strength Index) | Momentum | 5/10 | Double-smoothed momentum. Cleaner than RSI for divergences |
| 24 | Connors RSI | Oscillator | 5/10 | Composite of RSI + streak RSI + ROC percentile. Popular in mean-reversion algos |
| 25 | Chaikin Oscillator | Volume | 5/10 | MACD applied to ADL. Fast/slow ADL crossover for volume momentum |
| 26 | Force Index | Volume | 5/10 | Elder's price change * volume. Measures the force behind moves |
| 27 | ATR Trailing Stop | Volatility | 7/10 | ATR-based dynamic stop-loss. Used by almost everyone who uses ATR |
| 28 | Stc (Schaff Trend Cycle) | Trend | 5/10 | MACD through stochastic filter. Faster signals than MACD alone |
| 29 | McGinley Dynamic | Moving Average | 4/10 | Self-adjusting MA that tracks price better in fast markets |
| 30 | ALMA (Arnaud Legoux MA) | Moving Average | 5/10 | Gaussian-weighted MA. Smooth, low lag. Gaining popularity on TradingView |
| 31 | Klinger Oscillator | Volume | 4/10 | Volume force oscillator. Detects long-term money flow vs short-term |
| 32 | Vortex Indicator | Trend | 5/10 | +VI/-VI based on true range. Identifies trend start and direction |
| 33 | RVI (Relative Vigor Index) | Momentum | 4/10 | Compares close-open to high-low range. Measures conviction of moves |
| 34 | Elder Ray (Bull/Bear Power) | Momentum | 4/10 | You have Bull Bear Power already; this is the full Elder Ray system with EMA |
| 35 | Chop Index (CHOP) | Volatility | 5/10 | Measures if market is choppy (ranging) or trending. Good as a filter |
| 36 | ZigZag | Pattern | 5/10 | Connects swing highs/lows filtering noise by %. Not predictive but useful for pattern detection |
| 37 | Fibonacci Retracement | Support/Resistance | 8/10 | Not a traditional indicator, but worth having as a utility. Every trader knows Fib levels |
| 38 | Linear Regression Channel | Statistical | 5/10 | Regression line + stddev bands. Used for mean-reversion and trend channels |
| 39 | Standard Deviation | Statistical | 6/10 | Building block for BBands and other indicators. Useful standalone for volatility screening |
| 40 | ADR (Average Daily Range) | Volatility | 6/10 | Average of daily high-low range. Day traders use it to gauge how much room a move has left |
