package tree

import (
	"fmt"
	"math/rand"
	"time"
)

type BaseNode struct {
	UID string
}

func MakeBaseNode() BaseNode {
	rand.Seed(time.Now().UnixNano())
	return BaseNode{fmt.Sprintf("%v", rand.Int31())}
}

func (b *BaseNode) GetUID() string {
	return b.UID
}
