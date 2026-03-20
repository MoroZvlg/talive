package talive_test

import (
	"reflect"
	"testing"

	"github.com/MoroZvlg/talive"
)

func TestStochasticDefault(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/stochastic/output_default.csv", []int{1, 2}, 8)
	indicator, _ := talive.NewStochastic(14, 1, 3)
	result := [][]float64{
		make([]float64, len(candles)),
		make([]float64, len(candles)),
	}
	for i, candle := range candles {
		res := indicator.Next(candle)
		result[0][i] = roundFloat(res[0], 8)
		result[1][i] = roundFloat(res[1], 8)
	}
	if !(reflect.DeepEqual(result[0], expectedParsedData[0])) {
		t.Fatal(`[Stochastic(14, 1, 3)] values didn't match `, difference(result[0], expectedParsedData[0]))
	}
	if !(reflect.DeepEqual(result[1], expectedParsedData[1])) {
		t.Fatal(`[Stochastic(14, 1, 3)] values didn't match `, difference(result[1], expectedParsedData[1]))
	}
}

func TestStochasticMin(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/stochastic/output_min.csv", []int{1, 2}, 7)
	indicator, _ := talive.NewStochastic(2, 1, 2)
	result := [][]float64{
		make([]float64, len(candles)),
		make([]float64, len(candles)),
	}
	for i, candle := range candles {
		res := indicator.Next(candle)
		result[0][i] = roundFloat(res[0], 7)
		result[1][i] = roundFloat(res[1], 7)
	}
	if !(reflect.DeepEqual(result[0], expectedParsedData[0])) {
		t.Fatal(`[Stochastic(14, 1, 3)] values didn't match `, difference(result[0], expectedParsedData[0]))
	}
	if !(reflect.DeepEqual(result[1], expectedParsedData[1])) {
		t.Fatal(`[Stochastic(14, 1, 3)] values didn't match `, difference(result[1], expectedParsedData[1]))
	}
}

func TestStochasticIdle(t *testing.T) {
	indicator, _ := talive.NewStochastic(3, 1, 2)
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
		t.Fatal(`[Stochastic(3, 1, 2)] wrong idle value `, result)
	}
}

func TestRsiCurrent(t *testing.T) {
	candles, _ := readCandles("test_data/input_data2.csv")
	expectedParsedData, _ := readData("test_data/stochastic/output_default.csv", []int{1, 2}, 8)
	indicator, _ := talive.NewStochastic(14, 1, 3)
	for i := 0; i < 14; i++ {
		indicator.Next(candles[i])
	}
	currentValueK := roundFloat(indicator.Current(candles[14])[0], 8)
	currentValueD := roundFloat(indicator.Current(candles[14])[1], 8)
	expectedValueK := roundFloat(expectedParsedData[0][14], 8)
	expectedValueD := roundFloat(expectedParsedData[1][14], 8)

	if currentValueK != expectedValueK {
		t.Fatalf("[Stochastic(14, 1 3)] wrong Current K value %f, expected %f", currentValueK, expectedValueK)
	}
	if currentValueD != expectedValueD {
		t.Fatalf("[Stochastic(14, 1 3)] wrong Current D value %f, expected %f", currentValueD, expectedValueD)
	}

	nextValues := indicator.Next(candles[14])
	nextK := roundFloat(nextValues[0], 8)
	nextD := roundFloat(nextValues[1], 8)
	if nextK != expectedValueK {
		t.Fatalf("[Stochastic(14, 1 3)] Current value broke Next K value %f, expected %f", currentValueK, expectedValueK)
	}
	if nextD != expectedValueD {
		t.Fatalf("[Stochastic(14, 1 3)] Current value broke  Next D value %f, expected %f", currentValueD, expectedValueD)
	}
}

var stochDummy *talive.Stochastic

func Benchmark_Stochastic_Init_Allocations(b *testing.B) {
	b.Run("Stochastic(2, 1, 2)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stochDummy, _ = talive.NewStochastic(2, 1, 2)
		}
	})
	b.Run("Stochastic(14, 1, 3)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stochDummy, _ = talive.NewStochastic(14, 1, 3)
		}
	})
	b.Run("Stochastic(50, 10, 20)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stochDummy, _ = talive.NewStochastic(50, 10, 20)
		}
	})
}

func Benchmark_Stochastic_Next_Allocations(b *testing.B) {
	candles, _ := readCandles("test_data/input_data2.csv")
	dataLen := len(candles)

	b.Run("Stochastic(2, 1, 2)", func(b *testing.B) {
		stochDummy, _ = talive.NewStochastic(2, 1, 2)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = stochDummy.Next(candles[dataIndex])
		}
	})
	b.Run("Stochastic(14, 1, 3)", func(b *testing.B) {
		stochDummy, _ = talive.NewStochastic(14, 1, 3)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = stochDummy.Next(candles[dataIndex])
		}
	})
	b.Run("Stochastic(50, 10, 20)", func(b *testing.B) {
		stochDummy, _ = talive.NewStochastic(50, 10, 20)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = stochDummy.Next(candles[dataIndex])
		}
	})
}

func Benchmark_Stochastic_Current_Allocations(b *testing.B) {
	candles, _ := readCandles("test_data/input_data2.csv")
	dataLen := len(candles)

	b.Run("Stochastic(2, 1, 2)", func(b *testing.B) {
		stochDummy, _ = talive.NewStochastic(2, 1, 2)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = stochDummy.Current(candles[dataIndex])
		}
	})
	b.Run("Stochastic(14, 1, 3)", func(b *testing.B) {
		stochDummy, _ = talive.NewStochastic(14, 1, 3)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = stochDummy.Current(candles[dataIndex])
		}
	})
	b.Run("Stochastic(50, 10, 20)", func(b *testing.B) {
		stochDummy, _ = talive.NewStochastic(50, 10, 20)
		dataIndex := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dataIndex = limitedDataIndex(dataIndex, dataLen)
			sliceDummy = stochDummy.Current(candles[dataIndex])
		}
	})
}
