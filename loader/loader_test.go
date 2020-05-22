package loader

import (
	"testing"

	"github.com/richsoap/ID3Tree/adapter"
	"github.com/richsoap/ID3Tree/tree"
)

func TestLoadData(t *testing.T) {
	filePath := "../data/example/example.data"
	readData, err := LoaderData(filePath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	targetData := adapter.GetExampleAdapterSlice()
	if !adapter.CompareAdapterSlice(readData, targetData) {
		t.Error("Load data is different from target")
		t.Errorf("read class %v, target class %v", readData[0].Class, targetData[0].Class)
		for i := range readData {
			t.Errorf("%v, %v", readData[i].ToString(), targetData[i].ToString())
		}
	}
}

func TestLoadModel(t *testing.T) {
	filePath := "../data/example/example.model"
	readData, err := LoadeModel(filePath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	targetData := tree.GetExampleTree()
	if !tree.CompareTree(readData, targetData) {
		t.Error("Load model is different from target")
		t.Error("read")
		t.Error(readData.Serialize())
		t.Error("target")
		t.Error(targetData.Serialize())
	}
}
