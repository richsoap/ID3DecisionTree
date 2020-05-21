package builder

import (
	"testing"

	"github.com/richsoap/ID3Tree/adapter"
	"github.com/richsoap/ID3Tree/tree"
)

func TestBuild(t *testing.T) {
	source := adapter.GetExampleAdapterSlice()
	targetTree := tree.GetExampleTree()
	b := MakeBuilder()
	buildTree := b.BuildTree(source)
	if !tree.CompareTree(targetTree, buildTree) {
		t.Error("target is different")
		t.Errorf("target:\n%v", targetTree.Serialize())
		t.Errorf("build:\n%v", buildTree.Serialize())
	}
}
