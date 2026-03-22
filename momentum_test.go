package talive_test

import (
	"reflect"
	"testing"

	"github.com/MoroZvlg/talive"
)

func TestMomentumDefault(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/momentum/output_default.csv", []int{1}, 7)
	indicator, _ := talive.NewMomentum(10)
	result := make([]float64, len(candles))
	for i, candle := range candles {
		result[i] = roundFloat(indicator.Next(candle)[0], 7)
	}
	if !(reflect.DeepEqual(result, expectedParsedData[0])) {
		t.Fatal(`[Momentum(10)] values didn't match `, difference(result, expectedParsedData[0]))
	}
}

func TestMomentumMin(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/momentum/output_min.csv", []int{1}, 7)
	indicator, _ := talive.NewMomentum(2)
	result := make([]float64, len(candles))
	for i, candle := range candles {
		result[i] = roundFloat(indicator.Next(candle)[0], 7)
	}
	if !(reflect.DeepEqual(result, expectedParsedData[0])) {
		t.Fatal(`[Momentum(2)] values didn't match `, difference(result, expectedParsedData[0]))
	}
}

func TestMomentumIdle(t *testing.T) {
	indicator, _ := talive.NewMomentum(3)
	var result []string
	for i := 0; i < 5; i++ {
		indicator.Next(&testCandle{close: float64(i)})
		if indicator.IsIdle() {
			result = append(result, "true")
		} else {
			result = append(result, "false")
		}
	}
	if !reflect.DeepEqual(result, []string{"true", "true", "true", "false", "false"}) {
		t.Fatal(`[Momentum(3)] wrong idle value `, result)
	}
}

func TestMomentumCurrentValue(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/momentum/output_default.csv", []int{1}, 8)
	indicator, _ := talive.NewMomentum(10)
	for i := 0; i < 11; i++ {
		indicator.Next(candles[i])
	}
	currentValue := roundFloat(indicator.Current(candles[11])[0], 8)
	expectedValue := roundFloat(expectedParsedData[0][11], 8)
	if currentValue != expectedValue {
		t.Fatalf("[Momentum(10)] wrong Current value %f, expected %f", currentValue, expectedValue)
	}
	nextValue := roundFloat(indicator.Next(candles[11])[0], 8)
	if nextValue != currentValue {
		t.Fatalf("[Momentum(10)] Current value call broke Next value %f, expected %f", nextValue, expectedValue)
	}
}

var momentumDummy *talive.Momentum

func Benchmark_Momentum_Init_Allocations(benchmark *testing.B) {
	benchmark.Run("Momentum 2", func(benchmark *testing.B) {
		for i := 0; i < benchmark.N; i++ {
			momentumDummy, _ = talive.NewMomentum(2)
		}
	})
	benchmark.Run("Momentum 50", func(benchmark *testing.B) {
		for i := 0; i < benchmark.N; i++ {
			momentumDummy, _ = talive.NewMomentum(50)
		}
	})
}

func Benchmark_Momentum_Next_Allocations(benchmark *testing.B) {
	candles, _ := readCandles("test_data/input_data2.csv")
	dataLen := len(candles)
	benchmark.Run("Momentum 2", func(benchmark *testing.B) {
		indicator, _ := talive.NewMomentum(2)
		dataIndex := 0
		benchmark.ResetTimer()
		for i := 0; i < benchmark.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = indicator.Next(candles[dataIndex])
		}
	})
	benchmark.Run("Momentum 50", func(benchmark *testing.B) {
		indicator, _ := talive.NewMomentum(50)
		dataIndex := 0
		benchmark.ResetTimer()
		for i := 0; i < benchmark.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = indicator.Next(candles[dataIndex])
		}
	})
}

func Benchmark_Momentum_Current_Allocations(benchmark *testing.B) {
	candles, _ := readCandles("test_data/input_data2.csv")
	dataLen := len(candles)
	benchmark.Run("Momentum 2", func(benchmark *testing.B) {
		indicator, _ := talive.NewMomentum(2)
		dataIndex := 0
		benchmark.ResetTimer()
		for i := 0; i < benchmark.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = indicator.Current(candles[dataIndex])
		}
	})
	benchmark.Run("Momentum 50", func(benchmark *testing.B) {
		indicator, _ := talive.NewMomentum(50)
		dataIndex := 0
		benchmark.ResetTimer()
		for i := 0; i < benchmark.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = indicator.Current(candles[dataIndex])
		}
	})
}
