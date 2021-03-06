package tree

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/richsoap/ID3Tree/adapter"
	"github.com/richsoap/ID3Tree/utils"
)

type JudgeNode struct {
	Base     BaseNode
	Key      string
	Children map[string]Node
}

func MakeJudgeNode(key string) *JudgeNode {
	return &JudgeNode{MakeBaseNode(), key, make(map[string]Node)}
}

// return "NaK" if not fit any sub tree
func (j *JudgeNode) JudgeOne(data *adapter.Adapter) string {
	return j.Judge(data)[0]
}

// InValid Value will be "NaK" string
func (j *JudgeNode) Judge(data ...*adapter.Adapter) []string {
	group := utils.GroupBy(data, j.Key)
	resultMap := make(map[*adapter.Adapter]string)
	result := make([]string, len(data), len(data))
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(data))
	resChan := make(chan ResultEntry)
	defer close(resChan)

	for key := range group {
		subTree := j.Children[key]
		go IterateSubTree(subTree, group[key], resChan)
	}

	for i := 0; i < len(data); i++ {
		res, ok := <-resChan
		if !ok {
			log.Printf("channel closed unexpect")
		} else {
			resultMap[res.Data] = res.Result
		}
	}

	for i := range data {
		if res, ok := resultMap[data[i]]; ok {
			result[i] = res
		} else {
			log.Printf("can't found result for %v", i)
			result[i] = NAK
		}
	}
	return result
}

func (j *JudgeNode) IsMatched(data *adapter.Adapter) bool {
	if val, ok := data.Data[data.Class]; ok {
		return j.JudgeOne(data) == val
	}
	return false
}

func (j *JudgeNode) ErrorNum(data []*adapter.Adapter) int {
	judgeRes := j.Judge(data...)
	result := 0
	for i := range data {
		if judgeRes[i] != data[i].Data[data[i].Class] {
			result++
		}
	}
	return result
}

func (j *JudgeNode) ErrorRate(data []*adapter.Adapter) float64 {
	errorNum := j.ErrorNum(data)
	return float64(errorNum) / float64(len(data))
}

func (j *JudgeNode) AddNode(key string, node Node) error {
	if _, existed := j.Children[key]; existed {
		return errors.New("Node existd")
	}
	j.Children[key] = node
	return nil
}

func (j *JudgeNode) ToString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%v,judgenode,%v ", j.Base.UID, j.Key))
	for key := range j.Children {
		sb.WriteString(fmt.Sprintf(",%v:%v", key, j.Children[key].GetUID()))
	}
	return sb.String()
}

func (j *JudgeNode) GetUID() string {
	return j.Base.GetUID()
}

func (j *JudgeNode) Serialize() string {
	result := j.ToString() + "\n"
	for i := range j.Children {
		result += j.Children[i].Serialize()
	}
	return result
}

func (j *JudgeNode) Optimize(data []*adapter.Adapter) Node {
	if len(data) == 0 {
		return j
	}
	// Optimize subtree
	resChan := make(chan NodeEntry)
	defer close(resChan)
	group := utils.GroupBy(data, j.Key)
	for key := range j.Children {
		subSet, existed := group[key]
		if !existed {
			subSet = make([]*adapter.Adapter, 0, 0)
		}

		go func(n Node, data []*adapter.Adapter, key string, resChan chan NodeEntry) {
			resChan <- NodeEntry{key, n.Optimize(data)}
		}(j.Children[key], subSet, key, resChan)
	}
	for i := 0; i < len(j.Children); i++ {
		res, ok := <-resChan
		if !ok {
			log.Fatal("chan is closed")
		}
		j.Children[res.Value] = res.Node
	}
	// Optimize itself
	childError := j.ErrorNum(data)
	newKey, newMatch := utils.GetMajority(data, data[0].Class)
	newError := len(data) - newMatch
	if newError <= childError {
		return MakeLeafNode(newKey)
	} else {
		return j
	}
}
