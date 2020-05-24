package loader

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/richsoap/ID3Tree/tree"
)

func LoadTreeFromStrings(lines []string) (tree.Node, error) {
	nodeBuffer := make(map[string]tree.Node)
	structBuffer := make(map[string]map[string]string)
	childrenRecord := make(map[string]int)
	for _, line := range lines {
		cols := strings.Split(line, ",")
		for i := range cols {
			cols[i] = strings.Trim(cols[i], " ")
		}
		if cols[1] == "judgenode" {
			node := tree.MakeJudgeNode(cols[2])
			node.Base.UID = cols[0]
			children := make(map[string]string)
			for i := 3; i < len(cols); i++ {
				child := strings.Split(cols[i], ":")
				children[child[0]] = child[1]
				childrenRecord[child[1]] = 0
			}
			nodeBuffer[cols[0]] = node
			structBuffer[cols[0]] = children
		} else {
			node := tree.MakeLeafNode(cols[2])
			node.Base.UID = cols[0]
			nodeBuffer[cols[0]] = node
		}
	}
	rootUID := ""
	for key := range nodeBuffer {
		if _, existed := childrenRecord[key]; !existed {
			rootUID = key
			break
		}
	}
	if rootUID == "" {
		return nil, errors.New("Cannot find root")
	}
	return LinkNodes(rootUID, nodeBuffer, structBuffer), nil
}

func LoadTreeFromFile(datapath string) (tree.Node, error) {
	modeFile, err := os.Open(datapath)
	if err != nil {
		log.Printf("open modelfile %v error: %v", datapath, err)
		return nil, err
	}
	defer modeFile.Close()

	scanner := bufio.NewScanner(modeFile)
	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return LoadTreeFromStrings(lines)
}

func LinkNodes(uid string, nodeBuffer map[string]tree.Node, structBuffer map[string]map[string]string) tree.Node {
	result := nodeBuffer[uid]
	if children, existed := structBuffer[uid]; existed {
		for key := range children {
			result.AddNode(key, LinkNodes(children[key], nodeBuffer, structBuffer))
		}
	}
	return result
}
