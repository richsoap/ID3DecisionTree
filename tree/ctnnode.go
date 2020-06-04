package tree

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/richsoap/ID3Tree/adapter"
	"github.com/richsoap/ID3Tree/utils"
)

type CtnNode struct {
	Base       BaseNode
	Key        string
	DividValue float64
	Children   map[string]Node
}

func MakeCtnNode(key string, value float64) *CtnNode {
	return &CtnNode{MakeBaseNode(), key, value, make(map[string]Node)}
}

func (c *CtnNode) JudgeOne(data *adapter.Adapter) string {
	if data.CtnData[c.Key] < c.DividValue {
		return c.Children["Left"].Judge(data)[0]
	}
	return c.Children["Right"].Judge(data)[0]
}

type orderResultEntry struct {
	index  int
	result string
}

func judgeOrderResult(data []adapter.AdapterWithOrder, node Node, resChan chan orderResultEntry) {
	inputSlice := make([]*adapter.Adapter, len(data), len(data))
	result := node.Judge(inputSlice...)
	for index := range result {
		resChan <- orderResultEntry{data[index].Index, result[index]}
	}
}

func findMid(data []adapter.AdapterWithOrder, key string, value float64) int {
	for index := 0; index < len(data); index++ {
		if data[index].Data.CtnData[key] >= value {
			return index
		}
	}
	return len(data)
}

func (c *CtnNode) Judge(data ...*adapter.Adapter) []string {
	result := make([]string, len(data), len(data))
	sortSlice := adapter.MakeAnOrderSlice(data, c.Key)
	index := findMid(sortSlice, c.Key, c.DividValue)
	resChan := make(chan orderResultEntry)
	defer close(resChan)
	go judgeOrderResult(sortSlice[:index], c.Children["Left"], resChan)
	if index < len(data) {
		go judgeOrderResult(sortSlice[index:], c.Children["Right"], resChan)
	}
	for i := 0; i < len(data); i++ {
		res, ok := <-resChan
		if !ok {
			log.Printf("channel closed")
		} else {
			result[res.index] = res.result
		}
	}
	return result
}

func (c *CtnNode) IsMatched(data *adapter.Adapter) bool {
	if val, ok := data.Data[data.Class]; ok {
		return c.JudgeOne(data) == val
	}
	return false
}

func (c *CtnNode) ErrorNum(data []*adapter.Adapter) int {
	judgeRes := c.Judge(data...)
	result := 0
	for i := range data {
		if judgeRes[i] != data[i].Data[data[i].Class] {
			result++
		}
	}
	return result
}

func (c *CtnNode) ErrorRate(data []*adapter.Adapter) float64 {
	errorNum := c.ErrorNum(data)
	return float64(errorNum) / float64(len(data))
}

func (c *CtnNode) AddNode(key string, node Node) error {
	if key == "Left" {
		c.Children["Left"] = node
	} else if key == "Right" {
		c.Children["Right"] = node
	} else {
		return errors.New(fmt.Sprintf("Unknown key %v", key))
	}
	return nil
}

func (c *CtnNode) ToString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%v,ctnnode,%v:%v,left:%v,right:%v", c.Base.UID, c.Key, c.DividValue, c.Children["Left"].GetUID(), c.Children["Right"].GetUID()))
	return sb.String()
}

func (c *CtnNode) GetUID() string {
	return c.Base.GetUID()
}

func (c *CtnNode) Serialize() string {
	result := c.ToString() + "\n"
	for i := range c.Children {
		result += c.Children[i].Serialize()
	}
	return result
}

func (c *CtnNode) Optimize(data []*adapter.Adapter) Node {
	if len(data) == 0 {
		return c
	}
	beforeNum := c.ErrorNum(data)
	majority, afterNum := utils.GetMajority(data, data[0].Class)
	afterNum = len(data) - afterNum
	if beforeNum >= afterNum {
		return MakeLeafNode(majority)
	} else {
		return c
	}
}
