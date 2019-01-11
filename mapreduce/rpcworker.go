// Author: Federico Viglietta
package mapreduce

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
	"strings"
	"unicode"
)

// RPC service.
type Worker int

type MapperArgs struct {
	InputFiles  []string
	OutputFiles []string
}

type ReducerArgs struct {
	InputFiles []string
	OutputFile string
}

// Hashing function used to distribute words among reducers.
func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

// Utility function used to parse file.
func isPunctuation(c rune) bool {
	return !unicode.IsLetter(c) && !unicode.IsNumber(c)
}

// Compute the count of words in a collection of files storing the results in a set of intermediate files.
func (w *Worker) Map(in MapperArgs,  _ *string) error {
	files := in.InputFiles
	freq := make(map[string]int) // stores word occurrences
	// Compute word count
	for idx := 0; idx < len(files); idx++ {
		f, err := os.Open(files[idx])
		if err != nil {
			return err
		}
		s := bufio.NewScanner(f)
		s.Split(bufio.ScanWords)
		for s.Scan() {
			w := s.Text()
			w = strings.ToLower(w)
			f := strings.FieldsFunc(w, isPunctuation)
			for h := 0; h < len(f); h++ {
				freq[f[h]]++
			}
		}
		f.Close()
	}
	// Open intermediate files
	var tempFiles []*os.File
	for idx := 0; idx < len(in.OutputFiles); idx++ {
		f, err := os.OpenFile(in.OutputFiles[idx], os.O_RDWR, 0666)
		if err != nil {
			return err
		}
		tempFiles = append(tempFiles, f)
	}
	// Write the previously computed word count to temporary files.
	for w, f := range freq {
		dst := hash(w) % len(in.OutputFiles)
		fmt.Fprintf(tempFiles[dst], "%s %d\n", w, f)
	}
	// Close temporary files.
	for idx := 0; idx < len(tempFiles); idx++ {
		tempFiles[idx].Close()
	}
	return nil
}

// Merge the results stored in intermediate files putting results in the output file indicated by the master node.
func (w *Worker) Reduce(in ReducerArgs, _ *string) error {
	files := in.InputFiles
	freq := make(map[string]int) // stores word occurrences
	// Merge word count
	for idx := 0; idx < len(files); idx++ {
		f, err := os.Open(files[idx])
		if err != nil {
			return err
		}
		s := bufio.NewScanner(f)
		for s.Scan() {
			line := s.Text()
			var word string
			var count int
			fmt.Sscanf(line, "%s%d", &word, &count)
			freq[word] += count
		}
		f.Close()
	}
	// Store results in the output file
	out, err := os.OpenFile(in.OutputFile, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	for w, f := range freq {
		fmt.Fprintf(out, "%s %d\n", w, f)
	}
	out.Close()

	return nil
}
