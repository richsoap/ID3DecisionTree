package utils

import (
	"log"
	"math"
	"sync"

	"github.com/richsoap/ID3Tree/adapter"
)

// basic score func for builder
func IG(data []*adapter.Adapter, key string) float64 {
	return H(data, data[0].Class) - Rem(data, key)
}

func IGR(data []*adapter.Adapter, key string) float64 {
	return IG(data, key) / H(data, data[0].Class)
}

// return bit entropy of specific property
func H(data []*adapter.Adapter, key string) float64 {
	proportion := P(data, key)
	result := 0.0
	for i := range proportion {
		result += -proportion[i] * math.Log2(proportion[i])
	}
	return result
}

// GroupBy return a map
// Key is groupby property
// value is a sclice of adapter
func GroupBy(data []*adapter.Adapter, key string) map[string][]*adapter.Adapter {
	result := make(map[string][]*adapter.Adapter)
	for i := range data {
		index := data[i].Data[key]
		if _, ok := result[index]; !ok {
			result[index] = make([]*adapter.Adapter, 0)
		}
		result[index] = append(result[index], data[i])
	}
	return result
}

// GroupCount is like GroupBy, but just return the numbers of each kind
func GroupCount(data []*adapter.Adapter, key string) map[string]int {
	result := make(map[string]int)
	for i := range data {
		index := data[i].Data[key]
		if _, existed := result[index]; !existed {
			result[index] = 0
		}
		result[index]++
	}
	return result
}

// Get Majority value and numbers for specific kind
func GetMajority(data []*adapter.Adapter, key string) (string, int) {
	count := GroupCount(data, key)
	resKey := ""
	resCount := -1
	for k := range count {
		if count[k] > resCount {
			resCount = count[k]
			resKey = k
		}
	}
	return resKey, resCount
}

// return proportion of data for the specific key
func P(data []*adapter.Adapter, key string) map[string]float64 {
	grouped := GroupBy(data, key)
	result := make(map[string]float64)
	for key := range grouped {
		result[key] = float64(len(grouped[key])) / float64(len(data))
	}
	return result
}

func Rem(data []*adapter.Adapter, key string) float64 {
	group := GroupBy(data, key)
	result := 0.0
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(group))
	resChan := make(chan float64)
	defer close(resChan)
	for key := range group {
		go func(data []*adapter.Adapter, resChan chan float64, p float64) {
			result := H(data, data[0].Class)
			resChan <- result * p
		}(group[key], resChan, float64(len(group[key]))/float64(len(data)))
	}

	for i := 0; i < len(group); i++ {
		if res, ok := <-resChan; !ok {
			log.Fatal("chan was close unexpected")
			break
		} else {
			result += res
		}
	}
	return result
}

func CheckError(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}
