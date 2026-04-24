package talive_test

import (
	"reflect"
	"testing"

	"github.com/MoroZvlg/talive"
)

func TestOBVDefault(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/obv/output_default.csv", []int{1}, 6)
	indicator, _ := talive.NewOBV()
	result := make([]float64, len(candles))
	for i, candle := range candles {
		result[i] = roundFloat(indicator.Next(candle)[0], 6)
	}
	if !(reflect.DeepEqual(result, expectedParsedData[0])) {
		t.Fatal(`[OBV] values didn't match `, difference(result, expectedParsedData[0]))
	}
}

func TestOBVIdle(t *testing.T) {
	indicator, _ := talive.NewOBV()
	var result []string
	for i := 0; i < 5; i++ {
		indicator.Next(&testCandle{close: float64(i), volume: 1.0})
		if indicator.IsIdle() {
			result = append(result, "true")
		} else {
			result = append(result, "false")
		}
	}
	if !reflect.DeepEqual(result, []string{"true", "false", "false", "false", "false"}) {
		t.Fatal(`[OBV] wrong idle value `, result)
	}
	trueCount := 0
	for _, v := range result {
		if v == "true" {
			trueCount++
		}
	}
	if trueCount != indicator.IdlePeriod() {
		t.Fatalf("[OBV] IdlePeriod() = %d, but IsIdle() was true %d times", indicator.IdlePeriod(), trueCount)
	}
}

func TestOBVCurrentVal(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/obv/output_default.csv", []int{1}, 6)
	indicator, _ := talive.NewOBV()
	for i := 0; i < 11; i++ {
		indicator.Next(candles[i])
	}
	CurrentVal := roundFloat(indicator.Current(candles[11])[0], 6)
	expectedValue := roundFloat(expectedParsedData[0][11], 6)
	if CurrentVal != expectedValue {
		t.Fatalf("[OBV] wrong Current value %f, expected %f", CurrentVal, expectedValue)
	}
	NextVal := roundFloat(indicator.Next(candles[11])[0], 6)
	if NextVal != CurrentVal {
		t.Fatalf("[OBV] Current value call broke Next value %f, expected %f", NextVal, expectedValue)
	}
}

func Benchmark_OBV_Init_Allocations(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		benchSink, _ = talive.NewOBV()
	}
}

func Benchmark_OBV_Next_Allocations(benchmark *testing.B) {
	candles, _ := readCandles("test_data/input_data2.csv")
	dataLen := len(candles)
	indicator, _ := talive.NewOBV()
	dataIndex := 0
	benchmark.ResetTimer()
	for i := 0; i < benchmark.N; i++ {
		dataIndex = limitedDataIndex(dataIndex, dataLen)
		sliceDummy = indicator.Next(candles[dataIndex])
	}
}

func Benchmark_OBV_Current_Allocations(benchmark *testing.B) {
	candles, _ := readCandles("test_data/input_data2.csv")
	dataLen := len(candles)
	indicator, _ := talive.NewOBV()
	dataIndex := 0
	benchmark.ResetTimer()
	for i := 0; i < benchmark.N; i++ {
		dataIndex = limitedDataIndex(dataIndex, dataLen)
		sliceDummy = indicator.Current(candles[dataIndex])
	}
}
