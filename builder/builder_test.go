package builder

import (
	"fmt"
	"testing"

	"github.com/richsoap/ID3Tree/adapter"
)

func TestBuild(t *testing.T) {
	data := make([]*adapter.Adapter, 0)
	for i := 0; i < 5; i++ {
		a := adapter.MakeAdapter()
		a.Add("Key", fmt.Sprintf("%v", i))
		a.Add("SubKey", fmt.Sprintf("%v", i%2))
		data = append(data, a)
		t.Logf("data[%v]: %v", i, a.ToString())
	}
	a := adapter.MakeAdapter()
	a.Add("Key", "0")
	a.Add("SubKey", "1")
}
