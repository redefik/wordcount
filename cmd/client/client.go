package main

import (
	"flag"
	"fmt"
	"github.com/redefik/wordcount/config"
	"github.com/redefik/wordcount/mapreduce"
	"log"
	"net/rpc"
	"os"
)

var configurationFile = flag.String("config", "config/config.json", "Location of the config file.")

func main() {

	flag.Parse()
	// Parse input filename.
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: ./wordcount <filename> <filename> ...<filename>\n")
		os.Exit(1)
	}

	// Get system configuration
	config, err := config.GetConfiguration(*configurationFile)
	if err != nil {
		log.Fatal("Couldn't retrieve system configuration:", err)
	}
	// Connect to the master node.
	if err != nil {
		log.Fatal("Couldn't get master node address:", err)
	}
	masterAddress := config.Master[0]
	client, err := rpc.Dial("tcp", masterAddress)
	if err != nil {
		log.Fatal("Couldn't connect to master:", err)
	}
	defer client.Close()

	// Invoke the master RPC method to get the word count for the files given.
	input := mapreduce.WordCountInput{Files:args[1:], Config:config}
	var outputFiles []string
	err = client.Call("Master.WordCount", &input, &outputFiles)
	if err != nil {
		log.Fatal("Word count failed:", err)
	}

	// Give user location of output files
	fmt.Println("Work done!")
	fmt.Println("You can find results in:")
	for h := 0; h < len(outputFiles); h++ {
		fmt.Println(outputFiles[h])
	}
}
