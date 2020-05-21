package tree

import (
	"errors"
	"fmt"
	"strings"

	"github.com/richsoap/ID3Tree/adapter"
)

type LeafNode struct {
	Base   BaseNode
	Result string
}

func MakeLeafNode(result string) *LeafNode {
	return &LeafNode{MakeBaseNode(), result}
}

func (l *LeafNode) JudgeOne(data *adapter.Adapter) string {
	return l.Result
}

func (l *LeafNode) Judge(data ...*adapter.Adapter) []string {
	result := make([]string, len(data), len(data))
	for i := range result {
		result[i] = l.Result
	}
	return result
}

func (l *LeafNode) ErrorNum(data ...*adapter.Adapter) int {
	result := 0
	for i := range data {
		if l.Result != data[i].Data[data[i].Class] {
			result++
		}
	}
	return result
}

func (l *LeafNode) ErrorRate(data ...*adapter.Adapter) float64 {
	return float64(l.ErrorNum()) / float64(len(data))
}

func (l *LeafNode) AddNode(key string, node Node) error {
	return errors.New("Add Node in leafNode is forbiden")
}

func (l *LeafNode) ToString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%v,leafnode,%v ", l.Base.UID, l.Result))
	return sb.String()
}

func (l *LeafNode) GetUID() string {
	return l.Base.GetUID()
}

func (l *LeafNode) Serialize() string {
	return l.ToString() + "\n"
}

func (l *LeafNode) Optimize(data ...*adapter.Adapter) Node {
	return l
}
