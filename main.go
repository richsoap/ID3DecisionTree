package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/richsoap/ID3Tree/builder"
	"github.com/richsoap/ID3Tree/loader"
	"github.com/richsoap/ID3Tree/saver"
	"github.com/richsoap/ID3Tree/tree"
	"github.com/richsoap/ID3Tree/utils"
)

var (
	build     = flag.String("build", "", "Build a new dataset. Should pass dataset file path in")
	load      = flag.String("load", "", "Load a built model. Should pass model file path in")
	run       = flag.String("run", "", "Run data with build model or load model. Should pass dataset file path in.")
	optimize  = flag.String("optimize", "", "Use test dataset to optimize. Should pass dataset file path in.")
	output    = flag.String("output", "", "Output decicision result to file. Otherwise, stdout")
	save      = flag.String("save", "", "Save model to file. Otherwise, abort")
	scoreFunc = flag.String("func", "IG", "use IG or IGR as scoreFunc")
	depth     = flag.Int("depth", -1, "Max depth for precut")
	minNode   = flag.Int("min", -1, "The min number of data pieces, for pre cut")
)

func main() {
	flag.Parse()
	if *build != "" && *load != "" {
		log.Fatal("Build and Load cannot be used at the same time")
	}
	if *build == "" && *load == "" {
		log.Fatal("One of Build and Load should be uesd")
	}
	var decisionTree tree.Node
	if *build != "" {
		decisionTree = BuildTreeFromDataset()
	} else {
		decisionTree = LoadTreeFromFile()
	}
	if *optimize != "" {
		dataset, err := loader.LoaderData(*optimize)
		utils.CheckError(err)
		decisionTree.Optimize(dataset...)
	}
	if *save != "" {
		saver.SaveMode(decisionTree, *save)
	}
	if *run != "" {
		dataset, err := loader.LoaderData(*run)
		utils.CheckError(err)
		result := decisionTree.Judge(dataset...)
		fmt.Sprintf("decision result: %v", result)
		if *output != "" {
			saver.SaveResult(result, *output)
		}
	}
}

func BuildTreeFromDataset() tree.Node {
	dataset, err := loader.LoaderData(*build)
	if err != nil {
		log.Fatalf("%v", err)
	}
	b := builder.MakeBuilder()
	b.MaxDepth = *depth
	b.MinNode = *minNode
	if *scoreFunc == "IG" {
		b.ScoreFunc = utils.IG
	} else if *scoreFunc == "IGR" {
		b.ScoreFunc = utils.IGR
	} else {
		log.Fatalf("ScoreFunc should be IG or IGR")
	}
	return b.BuildTree(dataset)
}

func LoadTreeFromFile() tree.Node {
	node, err := loader.LoadeModel(*load)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return node
}
