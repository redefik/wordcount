// Author: Federico Viglietta
package main

import (
	"fmt"
	"github.com/redefik/wordcount/mapreduce"
	"log"
	"net"
	"net/rpc"
	"os"
)

func main() {

	// Listening address is required.
	args := os.Args
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <address>\n", args[0])
		os.Exit(1)
	}

	// Register RPC service.
	server := rpc.NewServer()
	service := new(mapreduce.Master)
	err := server.RegisterName("Master", service)
	if err != nil {
		log.Fatal(err)
	}
	address := args[1]
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for RPC requests incoming.
	server.Accept(listener)

}
