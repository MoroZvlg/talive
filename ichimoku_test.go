package talive_test

import (
	"reflect"
	"testing"

	"github.com/MoroZvlg/talive"
)

func TestIchimokuDefault(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/ichimoku/output_default.csv", []int{1, 2, 4, 5}, 7)
	indicator, _ := talive.NewIchimoku(9, 26, 52, 26)
	results := make([][]float64, 4)
	for i := range results {
		results[i] = make([]float64, len(candles))
	}
	for i, candle := range candles {
		out := indicator.Next(candle)
		for j := 0; j < 4; j++ {
			results[j][i] = roundFloat(out[j], 7)
		}
	}
	labels := []string{"Conversion Line", "Base Line", "Leading Span A", "Leading Span B"}
	for j := 0; j < 4; j++ {
		if !reflect.DeepEqual(results[j], expectedParsedData[j]) {
			t.Fatalf("[Ichimoku(9,26,52,26)] %s values didn't match %v", labels[j], difference(results[j], expectedParsedData[j]))
		}
	}
}

func TestIchimokuMin(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/ichimoku/output_min.csv", []int{1, 2, 4, 5}, 7)
	indicator, _ := talive.NewIchimoku(2, 3, 4, 5)
	results := make([][]float64, 4)
	for i := range results {
		results[i] = make([]float64, len(candles))
	}
	for i, candle := range candles {
		out := indicator.Next(candle)
		for j := 0; j < 4; j++ {
			results[j][i] = roundFloat(out[j], 7)
		}
	}
	labels := []string{"Conversion Line", "Base Line", "Leading Span A", "Leading Span B"}
	for j := 0; j < 4; j++ {
		if !reflect.DeepEqual(results[j], expectedParsedData[j]) {
			t.Fatalf("[Ichimoku(2,3,4,5)] %s values didn't match %v", labels[j], difference(results[j], expectedParsedData[j]))
		}
	}
}

func TestIchimokuIdle(t *testing.T) {
	indicator, _ := talive.NewIchimoku(2, 3, 4, 5)
	// IdlePeriod = max(3,4) + 5 - 2 = 7, first non-idle at bar 7
	var result []string
	for i := 0; i < 10; i++ {
		indicator.Next(&testCandle{high: float64(i + 2), low: float64(i), close: float64(i + 1)})
		if indicator.IsIdle() {
			result = append(result, "true")
		} else {
			result = append(result, "false")
		}
	}
	if !reflect.DeepEqual(result, []string{"true", "true", "true", "true", "true", "true", "true", "false", "false", "false"}) {
		t.Fatal(`[Ichimoku(2,3,4,5)] wrong idle value `, result)
	}
	trueCount := 0
	for _, v := range result {
		if v == "true" {
			trueCount++
		}
	}
	if trueCount != indicator.IdlePeriod() {
		t.Fatalf("[Ichimoku(2,3,4,5)] IdlePeriod() = %d, but IsIdle() was true %d times", indicator.IdlePeriod(), trueCount)
	}
}

func TestIchimokuCurrentValue(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/ichimoku/output_default.csv", []int{1, 2, 4, 5}, 8)
	indicator, _ := talive.NewIchimoku(9, 26, 52, 26)
	// Process past the idle period (76)
	for i := 0; i < 80; i++ {
		indicator.Next(candles[i])
	}
	currentOut := indicator.Current(candles[80])
	for j := 0; j < 4; j++ {
		currentValue := roundFloat(currentOut[j], 8)
		expectedValue := roundFloat(expectedParsedData[j][80], 8)
		if currentValue != expectedValue {
			labels := []string{"Conversion Line", "Base Line", "Leading Span A", "Leading Span B"}
			t.Fatalf("[Ichimoku(9,26,52,26)] wrong Current %s value %f, expected %f", labels[j], currentValue, expectedValue)
		}
	}
	nextOut := indicator.Next(candles[80])
	for j := 0; j < 4; j++ {
		nextValue := roundFloat(nextOut[j], 8)
		currentValue := roundFloat(currentOut[j], 8)
		if nextValue != currentValue {
			labels := []string{"Conversion Line", "Base Line", "Leading Span A", "Leading Span B"}
			t.Fatalf("[Ichimoku(9,26,52,26)] Current call broke Next %s value %f, expected %f", labels[j], nextValue, currentValue)
		}
	}
}

var ichimokuDummy *talive.Ichimoku

func Benchmark_Ichimoku_Init_Allocations(benchmark *testing.B) {
	benchmark.Run("Ichimoku 9,26,52,26", func(benchmark *testing.B) {
		for i := 0; i < benchmark.N; i++ {
			ichimokuDummy, _ = talive.NewIchimoku(9, 26, 52, 26)
		}
	})
	benchmark.Run("Ichimoku 2,3,4,5", func(benchmark *testing.B) {
		for i := 0; i < benchmark.N; i++ {
			ichimokuDummy, _ = talive.NewIchimoku(2, 3, 4, 5)
		}
	})
}

func Benchmark_Ichimoku_Next_Allocations(benchmark *testing.B) {
	candles, _ := readCandles("test_data/input_data2.csv")
	dataLen := len(candles)
	benchmark.Run("Ichimoku 9,26,52,26", func(benchmark *testing.B) {
		indicator, _ := talive.NewIchimoku(9, 26, 52, 26)
		dataIndex := 0
		benchmark.ResetTimer()
		for i := 0; i < benchmark.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = indicator.Next(candles[dataIndex])
		}
	})
	benchmark.Run("Ichimoku 2,3,4,5", func(benchmark *testing.B) {
		indicator, _ := talive.NewIchimoku(2, 3, 4, 5)
		dataIndex := 0
		benchmark.ResetTimer()
		for i := 0; i < benchmark.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = indicator.Next(candles[dataIndex])
		}
	})
}

func Benchmark_Ichimoku_Current_Allocations(benchmark *testing.B) {
	candles, _ := readCandles("test_data/input_data2.csv")
	dataLen := len(candles)
	benchmark.Run("Ichimoku 9,26,52,26", func(benchmark *testing.B) {
		indicator, _ := talive.NewIchimoku(9, 26, 52, 26)
		dataIndex := 0
		benchmark.ResetTimer()
		for i := 0; i < benchmark.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = indicator.Current(candles[dataIndex])
		}
	})
	benchmark.Run("Ichimoku 2,3,4,5", func(benchmark *testing.B) {
		indicator, _ := talive.NewIchimoku(2, 3, 4, 5)
		dataIndex := 0
		benchmark.ResetTimer()
		for i := 0; i < benchmark.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = indicator.Current(candles[dataIndex])
		}
	})
}
