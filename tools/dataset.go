package main

import (
	"bufio"
	"flag"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/richsoap/ID3Tree/saver"
)

var (
	input = flag.String("input", "", "input filepath")
	train = flag.String("train", "", "train filepath")
	test  = flag.String("test", "", "test filepath")
	p     = flag.Float64("p", 0.25, "test/all data")
)

func main() {
	flag.Parse()
	if *input == "" || *train == "" || *test == "" {
		log.Fatal("see -h")
	}
	inputfile, err := os.Open(*input)
	if err != nil {
		log.Fatalf("open modelfile %v error: %v", *input, err)
	}
	defer inputfile.Close()

	scanner := bufio.NewScanner(inputfile)
	if !scanner.Scan() {
		log.Fatal("title scan error")
	}
	title := scanner.Text()
	trainLines := make([]string, 0)
	testLines := make([]string, 0)
	for scanner.Scan() {
		if rand.Float64() < *p {
			testLines = append(testLines, scanner.Text())
		} else {
			trainLines = append(trainLines, scanner.Text())
		}
	}
	SaveLines(title, *test, testLines)
	SaveLines(title, *train, trainLines)
}

func SaveLines(title, filepath string, lines []string) {
	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString("\n")
	for i := range lines {
		sb.WriteString(lines[i])
		sb.WriteString("\n")
	}
	saver.SaveString(sb.String(), filepath)
}
