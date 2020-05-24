package loader

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/richsoap/ID3Tree/tree"
)

func LoadForestFromFile(filepath string) (*tree.Forest, error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("open modelfile %v error: %v", filepath, err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var result *tree.Forest
	if scanner.Scan() {
		line := scanner.Text()
		if !(line == tree.BAGGING || line == tree.BOOSTING || line == tree.SINGLE_TREE) {
			return nil, errors.New("First line should be bagging/single/boosting")
		}
		result = tree.MakeForest(line)
	} else {
		return nil, errors.New("Cannot scan the first line")
	}
	lines := make([]string, 0)
	weight := -1.0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "weight") {
			if weight != -1 {
				node, err := LoadTreeFromStrings(lines)
				if err != nil {
					return nil, err
				}
				result.AddTree(node, weight)
			}
			lines = make([]string, 0)
			newWeight, err := strconv.ParseFloat(line[7:], 64)
			if err != nil {
				return nil, err
			}
			weight = newWeight
			continue
		}
		lines = append(lines, line)
	}
	node, err := LoadTreeFromStrings(lines)
	if err != nil {
		return nil, err
	}
	result.AddTree(node, weight)
	return result, nil
}
