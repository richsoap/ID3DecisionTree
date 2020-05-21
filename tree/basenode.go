package tree

import (
	"math/rand"
	"time"
)

type BaseNode struct {
	UID int64
}

func MakeBaseNode() BaseNode {
	rand.Seed(time.Now().UnixNano())
	return BaseNode{rand.Int63()}
}

func (b *BaseNode) GetUID() int64 {
	return b.UID
}
