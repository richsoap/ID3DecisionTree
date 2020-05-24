package tree

import (
	"sort"
	"strings"

	"github.com/richsoap/ID3Tree/adapter"
	"github.com/richsoap/ID3Tree/utils"
)

const (
	SINGLE_TREE = "single"
	BOOSTING    = "boosting"
	BAGGING     = "bagging"
)

type Forest struct {
	Type  string
	Trees []Tree
}

type JudgeEntry struct {
	Weight float64
	Result []string
}

func MakeForest(Type string) *Forest {
	var f Forest
	f.Type = Type
	f.Trees = make([]Tree, 0)
	return &f
}

func (f *Forest) AddTree(tree Node, weight ...float64) {
	w := 1.0
	if len(weight) > 0 {
		w = weight[0]
	}
	f.Trees = append(f.Trees, MakeTree(tree, w))
	sort.Sort(TreeSlice(f.Trees))
}

func JudgeRoutine(data []*adapter.Adapter, tree Tree, resChan chan JudgeEntry) {
	resChan <- JudgeEntry{tree.Weight, tree.Root.Judge(data...)}
}

func (f *Forest) Judge(data ...*adapter.Adapter) []string {
	record := make([]map[string]float64, len(data), len(data))
	for i := range record {
		record[i] = make(map[string]float64)
	}
	resChan := make(chan JudgeEntry)
	defer close(resChan)
	for i := range f.Trees {
		go JudgeRoutine(data, f.Trees[i], resChan)
	}

	for i := 0; i < len(f.Trees); i++ {
		sResult := <-resChan
		judgeResult := sResult.Result
		weight := sResult.Weight
		for j := range data {
			if _, existed := record[j][judgeResult[j]]; !existed {
				record[j][judgeResult[j]] = 0
			}
			record[j][judgeResult[j]] += weight
		}
	}
	result := make([]string, len(data), len(data))
	for i := range record {
		result[i] = utils.GetMaxKey(record[i])
	}
	return result
}

func (f *Forest) ErrorNum(data []*adapter.Adapter) int {
	judgeResult := f.Judge(data...)
	result := 0
	class := data[0].Class
	for i := range data {
		if judgeResult[i] != data[i].Data[class] {
			result++
		}
	}
	return result
}

func (f *Forest) ErrorRate(data []*adapter.Adapter) float64 {
	return float64(f.ErrorNum(data)) / float64(len(data))
}

func (f *Forest) Optimize(data []*adapter.Adapter) {
	for i := range f.Trees {
		f.Trees[i].Root.Optimize(data)
	}
}

func CompareForest(a, b *Forest) bool {
	if a.Type != b.Type || len(a.Trees) != len(b.Trees) {
		return false
	}
	for i := range a.Trees {
		if !CompareTree(a.Trees[i], b.Trees[i]) {
			return false
		}
	}
	return true
}

func GetExampleForest() *Forest {
	result := MakeForest(BOOSTING)
	result.AddTree(GetExampleTree(), 0.5)
	result.AddTree(GetExampleTree(), 0.75)
	return result
}

func (f *Forest) Serialize() string {
	var sb strings.Builder
	sb.WriteString(f.Type)
	sb.WriteString("\n")
	for i := range f.Trees {
		sb.WriteString(f.Trees[i].Serialize())
	}
	return sb.String()
}
