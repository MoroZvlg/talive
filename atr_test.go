package talive_test

import (
	"reflect"
	"testing"

	"github.com/MoroZvlg/talive"
)

func TestAtrDefault(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/atr/output_default.csv", []int{1}, 7)
	indicator, _ := talive.NewATR(14)
	result := make([]float64, len(candles))
	for i, candle := range candles {
		result[i] = roundFloat(indicator.Next(candle)[0], 7)
	}
	if !(reflect.DeepEqual(result, expectedParsedData[0])) {
		t.Fatal(`[ATR(14)] values didn't match `, difference(result, expectedParsedData[0]))
	}
}

func TestAtrMin(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/atr/output_min.csv", []int{1}, 7)
	indicator, _ := talive.NewATR(2)
	result := make([]float64, len(candles))
	for i, candle := range candles {
		result[i] = roundFloat(indicator.Next(candle)[0], 7)
	}
	if !(reflect.DeepEqual(result, expectedParsedData[0])) {
		t.Fatal(`[ATR(2)] values didn't match `, difference(result, expectedParsedData[0]))
	}
}

func TestAtrIdle(t *testing.T) {
	indicator, _ := talive.NewATR(3)
	var result []string
	for i := 0; i < 4; i++ {
		indicator.Next(&testCandle{high: float64(i + 2), low: float64(i), close: float64(i + 1)})
		if indicator.IsIdle() {
			result = append(result, "true")
		} else {
			result = append(result, "false")
		}
	}
	if !reflect.DeepEqual(result, []string{"true", "true", "false", "false"}) {
		t.Fatal(`[ATR(3)] wrong idle value `, result)
	}
	trueCount := 0
	for _, v := range result {
		if v == "true" {
			trueCount++
		}
	}
	if trueCount != indicator.IdlePeriod() {
		t.Fatalf("[ATR(3)] IdlePeriod() = %d, but IsIdle() was true %d times", indicator.IdlePeriod(), trueCount)
	}
}

func TestAtrCurrentValue(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/atr/output_default.csv", []int{1}, 8)
	indicator, _ := talive.NewATR(14)
	for i := 0; i < 14; i++ {
		indicator.Next(candles[i])
	}
	currentValue := roundFloat(indicator.Current(candles[14])[0], 8)
	expectedValue := roundFloat(expectedParsedData[0][14], 8)
	if currentValue != expectedValue {
		t.Fatalf("[ATR(14)] wrong Current value %f, expected %f", currentValue, expectedValue)
	}
	nextValue := roundFloat(indicator.Next(candles[14])[0], 8)
	if nextValue != currentValue {
		t.Fatalf("[ATR(14)] Current value call broke Next value %f, expected %f", nextValue, expectedValue)
	}
}

var atrDummy *talive.ATR

func Benchmark_ATR_Init_Allocations(b *testing.B) {
	b.Run("ATR(2)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			atrDummy, _ = talive.NewATR(2)
		}
	})
	b.Run("ATR(50)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			atrDummy, _ = talive.NewATR(50)
		}
	})
}

func Benchmark_ATR_Next_Allocations(b *testing.B) {
	candles, _ := readCandles("test_data/input_data2.csv")
	dataLen := len(candles)

	b.Run("ATR(2)", func(b *testing.B) {
		atrDummy, _ = talive.NewATR(2)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = atrDummy.Next(candles[dataIndex])
		}
	})
	b.Run("ATR(50)", func(b *testing.B) {
		atrDummy, _ = talive.NewATR(50)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = atrDummy.Next(candles[dataIndex])
		}
	})
}

func Benchmark_ATR_Current_Allocations(b *testing.B) {
	candles, _ := readCandles("test_data/input_data2.csv")
	dataLen := len(candles)

	b.Run("ATR(2)", func(b *testing.B) {
		atrDummy, _ = talive.NewATR(2)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = atrDummy.Current(candles[dataIndex])
		}
	})
	b.Run("ATR(50)", func(b *testing.B) {
		atrDummy, _ = talive.NewATR(50)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = atrDummy.Current(candles[dataIndex])
		}
	})
}
