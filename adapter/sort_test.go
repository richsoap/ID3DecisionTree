package adapter

import "testing"

func TestSort(t *testing.T) {
	inputSlice := make([]*Adapter, 5, 5)
	targetSlice := make([]AdapterWithOrder, 5, 5)
	targetIndex := []int{0, 4, 1, 2, 3}
	targetValue := []float64{0, 2, 3, 6, 9}
	for index := range inputSlice {
		inputSlice[index] = MakeAdapter()
		inputSlice[index].AddCtn("key", float64((index*3)%10))

		targetSlice[index].Index = targetIndex[index]
		targetSlice[index].Data = MakeAdapter()
		targetSlice[index].Data.AddCtn("key", targetValue[index])
	}
	result := MakeAnOrderSlice(inputSlice, "key")
	for index := range result {
		if result[index].Index != targetSlice[index].Index {
			t.Error("Index not match")
			return
		}
		if !CompareAdapter(result[index].Data, targetSlice[index].Data) {
			t.Error("Adapter not match")
			return
		}
	}
}
