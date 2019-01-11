// Author: Federico Viglietta
package mapreduce

import (
	"fmt"
	"github.com/redefik/wordcount/config"
	"io/ioutil"
	"net/rpc"
	"os"
)

// RPC service.
type Master int

// List of file to process
type WordCountInput struct {
	Files []string
	Config config.Config
}

// Creates a collection of temporary files used to store the intermediate results of computation.
func createTempFiles(num int) ([]string, error) {
	var files []string
	for idx := 0; idx < num; idx++ {
		f, err := ioutil.TempFile(".", "temp")
		if err != nil {
			return nil, err
		}
		files = append(files, f.Name())
		f.Close()
	}
	return files, nil
}

// Deletes given temporary files.
func removeTempFiles(files [][]string) error {
	for h := 0; h < len(files); h++ {
		for k := 0; k < len(files[0]); k++ {
			err := os.Remove(files[h][k])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Create a collection of files used to store the results of computation.
func createOutputFiles(num int, outDir string) ([]string, error) {
	var files []string
	for idx := 0; idx < num; idx++ {
		name := fmt.Sprintf("%s/out%d.txt", outDir, idx)
		f,err := os.Create(name)
		if err != nil {
			return nil, err
		}
		files = append(files, f.Name())
		f.Close()
	}
	return files, nil
}

// Connects to the given list of address returning the list of established connection.
func connectTo(address []string) ([]*rpc.Client, error) {
	var connections []*rpc.Client
	for idx := 0; idx < len(address); idx++ {
		c, err := rpc.Dial("tcp", address[idx])
		if err != nil {
			return nil, err
		}
		connections = append(connections, c)
	}
	return connections, nil
}

// Close given RPC connections.
func disconnectFrom(connection []*rpc.Client) {
	for idx := 0; idx < len(connection); idx++ {
		connection[idx].Close()
	}
}

// Waits for workers completion
func waitWorkers(calls []*rpc.Call) error {
	for idx := 0; idx < len(calls); idx++ {
		c := <-calls[idx].Done
		if c.Error != nil {
			return c.Error
		}
	}
	return nil
}

// Computes the count of words in a collection of files given by the user.
// Results are stored into output files.
// The computation is done using a MapReduce approach.
func (master *Master) WordCount(in WordCountInput, out *[]string) error {

	var mapperConnection []*rpc.Client  // stores connection established with mapper nodes
	var reducerConnection []*rpc.Client // stores connection established with reducer nodes
	var mapperCalls []*rpc.Call // stores asynchronous call to mappers
	var reducerCalls []*rpc.Call // stores asynchronous call to reducers
	var m int // number of total mappers
	var r int // number of total reducers
	var activeMappers int    // number of working mappers
	var tempFiles [][]string // temporary files used to store the intermediate results of computation


	// Retrieve mappers address
	mappers := in.Config.Mapper
	m = len(mappers)

	//Retrieve reducers address
	reducers := in.Config.Reducer
	r = len(reducers)

	// Connect to mappers
	mapperConnection, err := connectTo(mappers)
	if err != nil {
		return err
	}

	// Assign map
	if len(in.Files) >= m {
		activeMappers = m
		for idx := 0; idx < m; idx++ {
			lower := idx * len(in.Files) / m
			upper := (idx + 1) * len(in.Files) / m
			temps, err := createTempFiles(r)
			if err != nil {
				removeTempFiles(tempFiles)
				return err
			}
			tempFiles = append(tempFiles, temps)
			input := MapperArgs{InputFiles: in.Files[lower:upper], OutputFiles: temps}
			call := mapperConnection[idx].Go("Worker.Map", &input, nil, nil)
			mapperCalls = append(mapperCalls, call)
		}
	} else {
		activeMappers = len(in.Files)
		for idx := 0; idx < activeMappers; idx++ {
			temps, err := createTempFiles(r)
			if err != nil {
				removeTempFiles(tempFiles)
				return err
			}
			tempFiles = append(tempFiles, temps)
			input := MapperArgs{InputFiles: in.Files[idx:idx + 1], OutputFiles:temps}
			call := mapperConnection[idx].Go("Worker.Map", &input, nil, nil)
			mapperCalls = append(mapperCalls, call)
		}
	}

	// Wait for mappers completion
	err = waitWorkers(mapperCalls)
	if err != nil {
		removeTempFiles(tempFiles)
		return err
	}

	// Connect to reducers
	reducerConnection, err = connectTo(reducers)
	if err != nil {
		removeTempFiles(tempFiles)
		return err
	}

	// Create output files
	outputFiles, err := createOutputFiles(r, in.Config.OutDir)
	if err != nil {
		removeTempFiles(tempFiles)
		return err
	}

	// Assign reduce
	for h := 0; h < r; h++ {
		var inputFiles []string
		for k := 0; k < activeMappers; k++ {
			inputFiles = append(inputFiles, tempFiles[k][h])
		}
		input := ReducerArgs{InputFiles:inputFiles, OutputFile:outputFiles[h]}
		call := reducerConnection[h].Go("Worker.Reduce", &input, nil, nil)
		reducerCalls = append(reducerCalls, call)
	}

	// Wait for reducers completion
	err = waitWorkers(reducerCalls)
	if err != nil {
		removeTempFiles(tempFiles)
		return err
	}

	// Remove intermediate files
	err = removeTempFiles(tempFiles)
	if err != nil {
		return err
	}

	// Close connections
	disconnectFrom(mapperConnection)
	disconnectFrom(reducerConnection)

	// Return output files to client
	*out = outputFiles

	return nil
}
