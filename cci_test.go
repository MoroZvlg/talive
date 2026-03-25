package talive_test

import (
	"reflect"
	"testing"

	"github.com/MoroZvlg/talive"
)

func TestCciDefault(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/cci/output_default.csv", []int{1}, 7)
	indicator, _ := talive.NewCCI(20)
	result := make([]float64, len(candles))
	for i, candle := range candles {
		result[i] = roundFloat(indicator.Next(candle)[0], 7)
	}
	if !(reflect.DeepEqual(result, expectedParsedData[0])) {
		t.Fatal(`[CCI(20)] values didn't match `, difference(result, expectedParsedData[0]))
	}
}

func TestCciMin(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/cci/output_min.csv", []int{1}, 6)
	indicator, _ := talive.NewCCI(2)
	result := make([]float64, len(candles))
	for i, candle := range candles {
		result[i] = roundFloat(indicator.Next(candle)[0], 6)
	}
	if !(reflect.DeepEqual(result, expectedParsedData[0])) {
		t.Fatal(`[CCI(2)] values didn't match `, difference(result, expectedParsedData[0]))
	}
}

func TestCciIdle(t *testing.T) {
	indicator, _ := talive.NewCCI(3)
	var result []string
	for i := 0; i < 4; i++ {
		indicator.Next(&testCandle{close: float64(i)})
		if indicator.IsIdle() {
			result = append(result, "true")
		} else {
			result = append(result, "false")
		}
	}
	if !reflect.DeepEqual(result, []string{"true", "true", "false", "false"}) {
		t.Fatal(`[CCI(3)] wrong idle value `, result)
	}
	trueCount := 0
	for _, v := range result {
		if v == "true" {
			trueCount++
		}
	}
	if trueCount != indicator.IdlePeriod() {
		t.Fatalf("[CCI(3)] IdlePeriod() = %d, but IsIdle() was true %d times", indicator.IdlePeriod(), trueCount)
	}
}

func TestCciCurrentValue(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/cci/output_default.csv", []int{1}, 8)
	indicator, _ := talive.NewCCI(20)
	for i := 0; i < 20; i++ {
		indicator.Next(candles[i])
	}
	currentValue := roundFloat(indicator.Current(candles[20])[0], 8)
	expectedValue := roundFloat(expectedParsedData[0][20], 8)
	if currentValue != expectedValue {
		t.Fatalf("[CCI(20)] wrong Current value %f, expected %f", currentValue, expectedValue)
	}
	nextValue := roundFloat(indicator.Next(candles[20])[0], 8)
	if nextValue != currentValue {
		t.Fatalf("[CCI(20)] Current value call broke Next value %f, expected %f", nextValue, expectedValue)
	}
}

var cciDummy *talive.CCI

func Benchmark_CCI_Init_Allocations(b *testing.B) {
	b.Run("CCI(2)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cciDummy, _ = talive.NewCCI(2)
		}
	})
	b.Run("CCI(20)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cciDummy, _ = talive.NewCCI(20)
		}
	})
	b.Run("CCI(50)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cciDummy, _ = talive.NewCCI(50)
		}
	})
}

func Benchmark_CCI_Next_Allocations(b *testing.B) {
	candles, _ := readCandles("test_data/input_data2.csv")
	dataLen := len(candles)

	b.Run("CCI(2)", func(b *testing.B) {
		cciDummy, _ = talive.NewCCI(2)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = cciDummy.Next(candles[dataIndex])
		}
	})
	b.Run("CCI(20)", func(b *testing.B) {
		cciDummy, _ = talive.NewCCI(20)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = cciDummy.Next(candles[dataIndex])
		}
	})
	b.Run("CCI(50)", func(b *testing.B) {
		cciDummy, _ = talive.NewCCI(50)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = cciDummy.Next(candles[dataIndex])
		}
	})
}

func Benchmark_CCI_Current_Allocations(b *testing.B) {
	candles, _ := readCandles("test_data/input_data2.csv")
	dataLen := len(candles)

	b.Run("CCI(2)", func(b *testing.B) {
		cciDummy, _ = talive.NewCCI(2)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = cciDummy.Current(candles[dataIndex])
		}
	})
	b.Run("CCI(20)", func(b *testing.B) {
		cciDummy, _ = talive.NewCCI(20)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = cciDummy.Current(candles[dataIndex])
		}
	})
	b.Run("CCI(50)", func(b *testing.B) {
		cciDummy, _ = talive.NewCCI(50)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = cciDummy.Current(candles[dataIndex])
		}
	})
}
