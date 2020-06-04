package main

import (
	"flag"
	"log"
	"time"

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
	scoreFunc = flag.String("func", "IG", "Use IG or IGR as scoreFunc")
	depth     = flag.Int("depth", -1, "Max depth for precut, default is -1")
	minNode   = flag.Int("leafsize", 0, "The min number of data pieces, for pre cut, default is 0")
	dot       = flag.String("dot", "", "Save model as DotFile")
	withUID   = flag.Bool("withUID", false, "Print UID in DotFile")
	forest    = flag.String("forest", "single", "Different decision forest type, single, boosting, bagging")
	trees     = flag.Int("trees", 5, "Forest size")
	setsize   = flag.Float64("setsize", 0.2, "Sample data set size")
	autostop  = flag.Bool("autostop", true, "Stop boosting train, when epsilon is 0.5")
)

func main() {
	flag.Parse()
	startTS := time.Now().UnixNano()
	if *build != "" && *load != "" {
		log.Fatal("Build and Load cannot be used at the same time")
	}
	if *build == "" && *load == "" {
		log.Fatal("One of Build and Load should be uesd, use -h for help")
	}
	var decisionForest *tree.Forest
	if *build != "" {
		decisionForest = BuildForestFromDataset()
	} else {
		decisionForest = LoadForestFromFile()
	}
	if *optimize != "" {
		dataset, err := loader.LoaderData(*optimize)
		utils.CheckError(err)
		log.Printf("Before Optimization: Error Rate: %v", decisionForest.ErrorRate(dataset))
		decisionForest.Optimize(dataset)
		log.Printf("After Optimization: Error Rate: %v", decisionForest.ErrorRate(dataset))
	}
	if *save != "" {
		saver.SaveForest(decisionForest, *save)
	}
	if *run != "" {
		log.Printf("load run data from %v", *run)
		dataset, err := loader.LoaderData(*run)
		utils.CheckError(err)
		result := decisionForest.Judge(dataset...)
		log.Printf("decision result Error Rate: %v", decisionForest.ErrorRate(dataset))
		if *output != "" {
			saver.SaveResult(result, *output)
		}
	}
	if *dot != "" {
		log.Printf("Save dot file to %v", *dot)
		saver.SaveForestAsDotFile(decisionForest, *dot, *withUID)
	}
	endTS := time.Now().UnixNano()
	log.Printf("use time %v nano second", endTS-startTS)
}

func BuildForestFromDataset() *tree.Forest {
	log.Printf("load train data from %v", *build)
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
	if *forest != tree.SINGLE_TREE && *forest != tree.BAGGING && *forest != tree.BOOSTING {
		log.Fatalf("forest should be one of %v/%v/%v", tree.SINGLE_TREE, tree.BAGGING, tree.BOOSTING)
	}
	f := builder.MakeForestBuilder(b, *forest, *trees, *setsize, *autostop)
	decisionForest := f.BuildForest(dataset)
	log.Printf("Train DataSet Error Rate: %v", decisionForest.ErrorRate(dataset))
	return decisionForest
}

func LoadForestFromFile() *tree.Forest {
	log.Printf("load model from %v", *load)
	node, err := loader.LoadForestFromFile(*load)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return node
}
