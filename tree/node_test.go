package tree

import (
	"fmt"
	"testing"

	"github.com/richsoap/ID3Tree/adapter"
	"github.com/richsoap/ID3Tree/utils"
)

func TestJudge(t *testing.T) {
	data := make([]*adapter.Adapter, 0)
	target := make([]string, 0)
	for i := 0; i < 5; i++ {
		a := adapter.MakeAdapter()
		a.Add("Key", fmt.Sprintf("%v", i))
		a.Add("SubKey", fmt.Sprintf("%v", i%2))
		data = append(data, a)
		t.Logf("data[%v]: %v", i, a.ToString())
	}

	tree := MakeJudgeNode("Key")
	for i := 0; i < 4; i++ {
		k := fmt.Sprintf("%v", i)
		n := MakeLeafNode(k)
		tree.AddNode(k, n)
		target = append(target, k)
	}
	subTree := MakeJudgeNode("SubKey")
	for i := 0; i < 2; i++ {
		leaf := MakeLeafNode(fmt.Sprintf("s%v", i))
		subTree.AddNode(fmt.Sprintf("%v", i), leaf)
	}

	tree.AddNode("4", subTree)
	target = append(target, "s0")
	t.Logf("Tree %v", tree.ToString())

	result := tree.Judge(data...)
	if !utils.ComapreStringSlice(result, target) {
		t.Errorf("judge error want %v, get %v", target, result)
	}
}
