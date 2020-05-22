package loader

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/richsoap/ID3Tree/adapter"
)

// first rows is filedesc,
// colname1:class colname2 colname3:id ...
// data1 data2 data3 ...
func LoaderData(datapath string) ([]*adapter.Adapter, error) {
	modeFile, err := os.Open(datapath)
	if err != nil {
		log.Printf("open datafile %v error: %v", datapath, err)
		return nil, err
	}
	defer modeFile.Close()

	scanner := bufio.NewScanner(modeFile)
	firstLine := true

	colsNames := make([]string, 0)
	nameIndex := -1
	classname := ""
	result := make([]*adapter.Adapter, 0)
	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Split(line, " ")
		for i := range cols {
			cols[i] = strings.Trim(cols[i], " ")
		}

		if firstLine {
			firstLine = false
			for i := range cols {
				words := strings.Split(cols[i], ":")
				if len(words) > 1 {
					for j := range words {
						words[j] = strings.Trim(words[j], " ")
					}
					if words[1] == "class" {
						classname = words[0]
					} else if words[1] == "id" {
						nameIndex = i
					}
				}
				colsNames = append(colsNames, words[0])
			}
			continue
		}
		a := adapter.MakeAdapter(classname)
		for i := range cols {
			if i == nameIndex {
				a.SetName(cols[i])
			} else {
				a.Add(colsNames[i], cols[i])
			}
		}
		result = append(result, a)
	}
	return result, nil
}
