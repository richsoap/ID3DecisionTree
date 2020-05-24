package builder

import (
	"math/rand"

	"github.com/richsoap/ID3Tree/adapter"
)

type DataSetProducer struct {
	Data   []*adapter.Adapter
	Weight []float64
}

func MakeDataSetProducer(data []*adapter.Adapter) *DataSetProducer {
	var d DataSetProducer
	d.Data = make([]*adapter.Adapter, len(data), len(data))
	copy(d.Data, data)
	d.Weight = make([]float64, len(data), len(data))
	for i := range d.Weight {
		d.Weight[i] = 1.0
	}
	return &d
}

func (d *DataSetProducer) UpdateWeight(weight []float64) {
	copy(d.Weight, weight)
}

func SelectRoutine(sums []float64, resChan chan int) {
	length := len(sums)
	val := rand.Float64() * sums[length-1]
	L := 0
	R := length - 1
	for R > L {
		mid := (L + R) / 2
		if sums[mid] > val {
			R = mid
		} else {
			L = mid + 1
		}
	}
	resChan <- L
}

func (d *DataSetProducer) ProduceDataSet(length int) []*adapter.Adapter {
	sum := 0.0
	sums := make([]float64, len(d.Data), len(d.Data))
	for i := range d.Data {
		sum += d.Weight[0]
		sums[i] = sum
	}
	result := make([]*adapter.Adapter, length, length)
	resChan := make(chan int)
	for i := 0; i < length; i++ {
		go SelectRoutine(sums, resChan)
	}
	for i := range result {
		index := <-resChan
		result[i] = d.Data[index]
	}
	return result
}

func (d *DataSetProducer) ResetWeight() {
	for i := range d.Weight {
		d.Weight[i] = 1.0
	}
}
