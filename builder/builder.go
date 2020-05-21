package builder

import (
	"log"

	"github.com/richsoap/ID3Tree/adapter"
	"github.com/richsoap/ID3Tree/tree"
	"github.com/richsoap/ID3Tree/utils"
)

type TreeBuilder struct {
	MaxDepth  int
	MinNode   int
	ScoreFunc func(data []*adapter.Adapter, key string) float64
}

type ScoreEntry struct {
	Key   string
	Score float64
}

type NodeEntry struct {
	Value string
	Node  tree.Node
}

func MakeBuilder() TreeBuilder {
	var b TreeBuilder
	b.MaxDepth = -1
	b.MinNode = -1
	b.ScoreFunc = utils.IG
	return b
}

// Build tree func, depth is used to precut, which is optional
func (b *TreeBuilder) BuildTree(data []*adapter.Adapter, depth ...int) tree.Node {
	currDep := 1
	if len(depth) != 0 {
		currDep = depth[0]
	}
	// precut depth
	if b.MaxDepth != -1 && currDep >= b.MaxDepth {
		return b.BuildLeafNode(data)
	}
	// precut minNode
	if b.MinNode != -1 && len(data) < b.MinNode {
		return b.BuildLeafNode(data)
	}
	// all the same
	if majority, num := utils.GetMajority(data, data[0].Class); num == len(data) {
		return tree.MakeLeafNode(majority)
	}
	// decide best bracn
	tryCount := 0
	resChan := make(chan ScoreEntry)
	defer close(resChan)
	for key := range data[0].Data {
		if key == data[0].Class { //result should not be used
			continue
		}
		if _, existed := data[0].UsedKey[key]; existed {
			continue
		}
		tryCount++
		go b.ScoreRoutine(data, key, resChan)
	}
	if tryCount == 0 {
		return b.BuildLeafNode(data) // All kind has been used
	}
	// Get the best score
	bestScore := ScoreEntry{"", -10000}
	for i := 0; i < tryCount; i++ {
		res, ok := <-resChan
		if !ok {
			log.Fatal("chan was closed")
		}
		if res.Score > bestScore.Score {
			bestScore = res
		}
	}
	for i := range data {
		data[i].AddUsedKey(bestScore.Key)
	}
	node := tree.MakeJudgeNode(bestScore.Key)
	group := utils.GroupBy(data, bestScore.Key)
	nodeChan := make(chan NodeEntry)
	defer close(nodeChan)
	for key := range group {
		go b.BuildTreeRoutine(group[key], key, currDep+1, nodeChan)
	}
	for i := 0; i < len(group); i++ {
		res, ok := <-nodeChan
		if !ok {
			log.Fatal("chan was closed")
		}
		node.AddNode(res.Value, res.Node)
	}
	return node
}

func (b *TreeBuilder) ScoreRoutine(data []*adapter.Adapter, key string, resChan chan ScoreEntry) {
	score := b.ScoreFunc(data, key)
	resChan <- ScoreEntry{key, score}
}

func (b *TreeBuilder) BuildTreeRoutine(data []*adapter.Adapter, key string, depth int, resChan chan NodeEntry) {
	resChan <- NodeEntry{key, b.BuildTree(data, depth)}
}

// Called by BuildTree func, to build a leafnode with majority
func (b *TreeBuilder) BuildLeafNode(data []*adapter.Adapter) tree.Node {
	majority, _ := utils.GetMajority(data, data[0].Class)
	return tree.MakeLeafNode(majority)
}
