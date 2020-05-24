package saver

import (
	"fmt"
	"strings"

	"github.com/richsoap/ID3Tree/tree"
)

func SaveTreeAsDotFile(node tree.Node, filepath string, withUID bool) {
	var sb strings.Builder
	sb.WriteString("digraph action {") // Dot file head
	sb.WriteString(SprintTreeAsDotFile(node, withUID))
	sb.WriteString("}")
	SaveString(sb.String(), filepath)
}

func SaveForestAsDotFile(forest *tree.Forest, filepath string, withUID bool) {
	var sb strings.Builder
	sb.WriteString("digraph action {\nnode [shape = ellipse]\n") // Dot file head
	sb.WriteString(SprintForestAsDotFile(forest, withUID))
	sb.WriteString("}")
	SaveString(sb.String(), filepath)
}

func SprintForestAsDotFile(forest *tree.Forest, withUID bool) string {
	var sb strings.Builder
	if forest.Type == tree.SINGLE_TREE {
		sb.WriteString("weight=1\n")
		sb.WriteString(SprintTreeAsDotFile(forest.Trees[0].Root, withUID))
	} else if forest.Type == tree.BOOSTING {
		rootNode := tree.MakeBaseNode()
		sb.WriteString(fmt.Sprintf(NODE_TEMPLATE, rootNode.GetUID(), tree.BOOSTING))
		for _, t := range forest.Trees {
			sb.WriteString(fmt.Sprintf(LINK_TEMPLATE_WITHLAB, rootNode.GetUID(), t.Root.GetUID(), t.Weight))
			sb.WriteString(SprintTreeAsDotFile(t.Root, withUID))
		}
	} else if forest.Type == tree.BAGGING {
		rootNode := tree.MakeBaseNode()
		sb.WriteString(fmt.Sprintf(NODE_TEMPLATE, rootNode.GetUID(), tree.BAGGING))
		for _, t := range forest.Trees {
			sb.WriteString(fmt.Sprintf(LINK_TEMPLATE, rootNode.GetUID(), t.Root.GetUID()))
			sb.WriteString(SprintTreeAsDotFile(t.Root, withUID))
		}
	}
	return sb.String()
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
		sb.WriteString(fmt.Sprintf(LINK_TEMPLATE_WITHLAB, j.GetUID(), j.Children[key].GetUID(), key))
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
