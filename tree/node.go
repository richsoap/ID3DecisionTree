package tree

import (
	"sync"

	"github.com/richsoap/ID3Tree/adapter"
)

const (
	NAK = "NaK"
)

type Node interface {
	JudgeOne(data *adapter.Adapter) string
	Judge(data ...*adapter.Adapter) []string
	ErrorNum(data ...*adapter.Adapter) int
	ErrorRate(data ...*adapter.Adapter) float64
	AddNode(key string, node Node) error
	ToString() string
	GetUID() int64
}

type ResultEntry struct {
	Data   *adapter.Adapter
	Result string
}

func IterateSubTree(subTree Node, data []*adapter.Adapter, resChan chan ResultEntry) {
	var result []string
	if subTree != nil {
		result = subTree.Judge(data...)
	} else {
		result = make([]string, len(data), len(data))
		for i := range result {
			result[i] = NAK
		}
	}
	for i := range result {
		resChan <- ResultEntry{data[i], result[i]}
	}
}

func MarkNakData(data []*adapter.Adapter, wg *sync.WaitGroup, resChan chan ResultEntry) {
	defer wg.Done()
	for i := range data {
		resChan <- ResultEntry{data[i], NAK}
	}
}

func CollectResult(resMap *map[*adapter.Adapter]string, resChan chan ResultEntry, wg *sync.WaitGroup) {
	for res := range resChan {
		(*resMap)[res.Data] = res.Result
		wg.Done()
	}
}
