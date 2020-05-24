package builder

import (
	"log"
	"math"

	"github.com/richsoap/ID3Tree/adapter"
	"github.com/richsoap/ID3Tree/tree"
)

// Use for build a decision forest
type ForestBuilder struct {
	Builder  *TreeBuilder
	Type     string
	Number   int
	Size     float64
	AutoStop bool
}

func MakeForestBuilder(b *TreeBuilder, t string, num int, size float64, autostop bool) *ForestBuilder {
	return &ForestBuilder{b, t, num, size, autostop}
}

func (f *ForestBuilder) BuildForest(data []*adapter.Adapter) *tree.Forest {
	if f.Type == tree.SINGLE_TREE {
		return f.BuildSingleTree(data)
	} else if f.Type == tree.BOOSTING {
		return f.BuildBoosting(data)
	} else if f.Type == tree.BAGGING {
		return f.BuildBagging(data)
	} else {
		return nil
	}
}

func (f *ForestBuilder) BuildSingleTree(data []*adapter.Adapter) *tree.Forest {
	result := tree.MakeForest(tree.SINGLE_TREE)
	result.AddTree(f.Builder.BuildTree(data), 1)
	return result
}

func (f *ForestBuilder) BuildBoosting(data []*adapter.Adapter) *tree.Forest {
	dataProcuder := MakeDataSetProducer(data)
	dataSetSize := int(float64(len(data)) * f.Size)
	result := tree.MakeForest(tree.BOOSTING)
	class := data[0].Class
	for i := 0; i < f.Number; i++ {
		dataset := dataProcuder.ProduceDataSet(dataSetSize)
		t := f.Builder.BuildTree(dataset)
		judgeResult := t.Judge(data...)
		epsilon := 0.0
		for i := range dataProcuder.Weight {
			if judgeResult[i] != data[i].Data[class] {
				epsilon += dataProcuder.Weight[i]
			}
		}
		log.Printf("%v iter: epsilon=%v", i+1, epsilon)
		result.AddTree(t, 0.5*math.Log2((1-epsilon)/epsilon)/math.Log2E)
		for i := range dataProcuder.Weight {
			if judgeResult[i] != data[i].Data[class] {
				dataProcuder.Weight[i] = dataProcuder.Weight[i] / 2 / epsilon
			} else {
				dataProcuder.Weight[i] = dataProcuder.Weight[i] / 2 / (1 - epsilon)
			}
		}
		/*if math.Abs(0.5-epsilon) < 0.001 && f.AutoStop {
			log.Printf("epsilon(%v) is approximately equal to 0.5， stop training", epsilon)
			break
		}*/
	}
	return result
}

func (f *ForestBuilder) BuildBagging(data []*adapter.Adapter) *tree.Forest {
	dataProcuder := MakeDataSetProducer(data)
	dataSetSize := int(float64(len(data)) * f.Size)
	result := tree.MakeForest(tree.BAGGING)
	for i := 0; i < f.Number; i++ {
		dataset := dataProcuder.ProduceDataSet(dataSetSize)
		t := f.Builder.BuildTree(dataset)
		result.AddTree(t, 1)
	}
	return result
}
