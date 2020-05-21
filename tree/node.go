package tree

import (
	"fmt"
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
	Serialize() string
	Optimize(data ...*adapter.Adapter) Node
	GetUID() string
}

type ResultEntry struct {
	Data   *adapter.Adapter
	Result string
}

type NodeEntry struct {
	Value string
	Node  Node
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

func CompareTree(a, b Node) bool {
	return ComapreJudgeNode(a, b) || CompareLeafNode(a, b)
}

func ComapreJudgeNode(ra, rb Node) bool {
	a, aok := ra.(*JudgeNode)
	b, bok := rb.(*JudgeNode)
	if !aok || !bok {
		return false
	}
	if a.Key != b.Key {
		return false
	}
	if len(a.Key) != len(b.Key) {
		return false
	}
	for key := range a.Children {
		if bchild, existed := b.Children[key]; existed {
			if !CompareTree(a.Children[key], bchild) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func CompareLeafNode(ra, rb Node) bool {
	a, aok := ra.(*LeafNode)
	b, bok := rb.(*LeafNode)
	if !aok || !bok {
		return false
	}
	return a.Result == b.Result
}

func GetExampleTree() Node {
	tree := MakeJudgeNode("Key")
	for i := 0; i < 4; i++ {
		k := fmt.Sprintf("%v", i)
		n := MakeLeafNode(k)
		tree.AddNode(k, n)
	}
	subTree := MakeJudgeNode("SubKey")
	for i := 0; i < 2; i++ {
		leaf := MakeLeafNode(fmt.Sprintf("s%v", i))
		subTree.AddNode(fmt.Sprintf("%v", i), leaf)
	}

	tree.AddNode("4", subTree)
	return tree
}
