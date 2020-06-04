package adapter

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Adapter struct {
	Name    string
	Data    map[string]string
	CtnData map[string]float64
	Class   string
	UsedKey map[string]void
}

type void struct{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func MakeAdapter(classname ...string) *Adapter {
	var res Adapter
	res.Name = fmt.Sprintf("%v", rand.Int31())
	res.Data = make(map[string]string)
	res.CtnData = make(map[string]float64)
	res.Class = ""
	if len(classname) > 0 {
		res.Class = classname[0]
	}
	res.UsedKey = make(map[string]void)
	return &res
}

func (a *Adapter) SetName(name string) {
	a.Name = name
}

func (a *Adapter) SetClass(c string) {
	a.Class = c
}

func (a *Adapter) Add(key, value string) {
	a.Data[key] = value
}

func (a *Adapter) AddCtn(key string, value float64) {
	a.CtnData[key] = value
}

func (a *Adapter) Remove(key string) {
	if _, existed := a.Data[key]; existed {
		delete(a.Data, key)
	} else if _, existed := a.CtnData[key]; existed {
		delete(a.CtnData, key)
	}
}

func (a *Adapter) AddUsedKey(key string) {
	a.UsedKey[key] = void{}
}

func (a *Adapter) RemoveUsedKey(key string) {
	if _, existd := a.UsedKey[key]; existd {
		delete(a.UsedKey, key)
	}
}

func (a *Adapter) GetUnusedKeys() []string {
	res := make([]string, 0)
	for key := range a.Data {
		if key == a.Class {
			continue
		}
		if _, existed := a.UsedKey[key]; !existed {
			res = append(res, key)
		}
	}
	return res
}

func (a *Adapter) ToString() string {
	return fmt.Sprintf("%v", a.Data)
}

func (a *Adapter) ResetUsedKey() {
	a.UsedKey = make(map[string]void)
}

func (a *Adapter) IsUsedKey(key string) bool {
	_, existed := a.UsedKey[key]
	return existed
}

func GetExampleAdapterSlice() []*Adapter {
	data := make([]*Adapter, 0)
	for i := 0; i < 5; i++ {
		a := MakeAdapter("Target")
		a.Add("Key", fmt.Sprintf("%v", i))
		a.Add("SubKey", fmt.Sprintf("%v", i%2))
		a.Add("Target", fmt.Sprintf("%v", i))
		data = append(data, a)
	}
	data[len(data)-1].Data["Target"] = "s0"
	a := MakeAdapter("Target")
	a.Add("Key", "4")
	a.Add("SubKey", "1")
	a.Add("Target", "s1")
	data = append(data, a)
	return data
}

func CompareAdapter(a, b *Adapter) bool {
	if len(a.Data) != len(b.Data) || len(a.CtnData) != len(b.CtnData) {
		return false
	}
	if a.Class != b.Class {
		return false
	}
	for key := range a.Data {
		if bval, existed := b.Data[key]; existed {
			if bval != a.Data[key] {
				return false
			}
		} else {
			return false
		}
	}
	for key := range a.CtnData {
		if bval, existed := b.CtnData[key]; existed {
			if bval != a.CtnData[key] {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func CompareAdapterSlice(a, b []*Adapter) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !CompareAdapter(a[i], b[i]) {
			return false
		}
	}
	return true
}

type AdapterWithOrder struct {
	Data  *Adapter
	Index int
}

type adapterSortInterface struct {
	data []AdapterWithOrder
	key  string
}

func (s adapterSortInterface) Len() int {
	return len(s.data)
}

func (s adapterSortInterface) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

func (s adapterSortInterface) Less(i, j int) bool {
	return s.data[i].Data.CtnData[s.key] < s.data[j].Data.CtnData[s.key]
}

func MakeAnOrderSlice(data []*Adapter, key string) []AdapterWithOrder {
	result := make([]AdapterWithOrder, len(data), len(data))
	for index := range data {
		result[index] = AdapterWithOrder{data[index], index}
	}
	sortInterface := adapterSortInterface{result, key}
	sort.Sort(sortInterface)
	return result
}
