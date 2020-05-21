package adapter

import "fmt"

type Adapter struct {
	Name    string
	Data    map[string]string
	Class   string
	UsedKey map[string]void
}

type void struct{}

func MakeAdapter() *Adapter {
	var res Adapter
	res.Name = ""
	res.Data = make(map[string]string)
	res.Class = ""
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

func (a *Adapter) Remove(key string) {
	if _, existed := a.Data[key]; existed {
		delete(a.Data, key)
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
		if _, existed := a.UsedKey[key]; !existed {
			res = append(res, key)
		}
	}
	return res
}

func (a *Adapter) ToString() string {
	return fmt.Sprintf("%v", a.Data)
}
