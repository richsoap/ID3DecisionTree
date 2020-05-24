package tree

import (
	"fmt"
	"strings"
)

type Tree struct {
	Root   Node
	Weight float64
}

func MakeTree(tree Node, weight float64) Tree {
	return Tree{tree, weight}
}

func (t *Tree) Serialize() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("weight=%v\n", t.Weight))
	sb.WriteString(t.Root.Serialize())
	return sb.String()
}

type TreeSlice []Tree

func (s TreeSlice) Len() int {
	return len(s)
}

func (s TreeSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s TreeSlice) Less(i, j int) bool {
	return s[i].Weight < s[j].Weight
}

func CompareTree(a, b Tree) bool {
	if a.Weight != b.Weight {
		return false
	}
	return CompareNode(a.Root, b.Root)
}
