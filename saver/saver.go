package saver

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/richsoap/ID3Tree/tree"
	"github.com/richsoap/ID3Tree/utils"
)

const (
	NODE_TEMPLATE         = "%v [shape = ellipse label = \"%v\"]\n"     //UID, key/result
	NODE_TEMPLATE_WITHUID = "%v [shape = ellipse label = \"%v[%v]\"]\n" //UID, key/result, UID
	LINK_TEMPLATE         = "%v->%v [label = \"%v\"]\n"                 //UID, child UID, key
)

// This func will panic, if error occurs
func SaveModel(node tree.Node, filepath string) {
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

func SaveModelAsDotFile(node tree.Node, filepath string, withUID bool) {
	var sb strings.Builder
	sb.WriteString("digraph action {") // Dot file head
	sb.WriteString(SprintTreeAsDotFile(node, withUID))
	sb.WriteString("}")
	SaveString(sb.String(), filepath)
}

func SprintTreeAsDotFile(node tree.Node, withUID bool) string {
	if res, ok := SprintJudgeNodeAsDotFile(node, withUID); ok {
		return res
	}
	if res, ok := SprintLeafNodeAsDotFile(node, withUID); ok {
		return res
	}
	return fmt.Sprintf("%v Conver Error", node.GetUID())
}

func SprintJudgeNodeAsDotFile(node tree.Node, withUID bool) (string, bool) {
	j, ok := node.(*tree.JudgeNode)
	if !ok {
		return "", false
	}
	var sb strings.Builder
	if withUID {
		sb.WriteString(fmt.Sprintf(NODE_TEMPLATE_WITHUID, j.GetUID(), j.Key, j.GetUID()))
	} else {
		sb.WriteString(fmt.Sprintf(NODE_TEMPLATE, j.GetUID(), j.Key))
	}
	for key := range j.Children {
		sb.WriteString(fmt.Sprintf(LINK_TEMPLATE, j.GetUID(), j.Children[key].GetUID(), key))
		sb.WriteString(SprintTreeAsDotFile(j.Children[key], withUID))
	}
	return sb.String(), true
}

func SprintLeafNodeAsDotFile(node tree.Node, withUID bool) (string, bool) {
	l, ok := node.(*tree.LeafNode)
	if !ok {
		return "", false
	}
	if withUID {
		return fmt.Sprintf(NODE_TEMPLATE_WITHUID, l.GetUID(), l.Result, l.GetUID()), true
	} else {
		return fmt.Sprintf(NODE_TEMPLATE, l.GetUID(), l.Result), true
	}
}
