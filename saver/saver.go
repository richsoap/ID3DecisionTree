package saver

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/richsoap/ID3Tree/tree"
	"github.com/richsoap/ID3Tree/utils"
)

const (
	NODE_TEMPLATE         = "%v [label = \"%v\"]\n"     //UID, key/result
	NODE_TEMPLATE_WITHUID = "%v [label = \"%v[%v]\"]\n" //UID, key/result, UID
	LINK_TEMPLATE         = "%v->%v\n"                  // UID, child UID
	LINK_TEMPLATE_WITHLAB = "%v->%v [label = \"%v\"]\n" //UID, child UID, key
)

func SaveForest(forest *tree.Forest, filepath string) {
	forestString := forest.Serialize()
	SaveString(forestString, filepath)
}

// This func will panic, if error occurs
func SaveTree(node tree.Node, filepath string) {
	modeString := node.Serialize()
	SaveString(modeString, filepath)
}

// This func will panic, if error occurs
func SaveResult(result []string, filepath string) {
	var sb strings.Builder
	for i := range result {
		sb.WriteString(result[i])
		sb.WriteString(" ")
	}
	SaveString(sb.String(), filepath)
}

// This func will panic, if error occurs
func SaveString(str string, filepath string) {
	f, err := os.Create(filepath)
	utils.CheckError(err)
	n, err := io.WriteString(f, str)
	utils.CheckError(err)
	if n != len(str) {
		log.Fatalf("Serialized string contains %v chars, only write %v chars", len(str), n)
	}
}
