package tree

import (
	"testing"

	"github.com/richsoap/ID3Tree/adapter"
	"github.com/richsoap/ID3Tree/utils"
)

func TestJudge(t *testing.T) {
	data := adapter.GetExampleAdapterSlice()
	tree := GetExampleTree()
	target := make([]string, 0)
	for i := range data {
		target = append(target, data[i].Data[data[i].Class])
	}
	t.Logf("Tree %v", tree.ToString())

	result := tree.Judge(data...)
	if !utils.CompareStringSlice(result, target) {
		t.Errorf("judge error want %v, get %v", target, result)
	}
}
