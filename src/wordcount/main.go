package main

import (
	mr "mapreduce"
	"strings"
	"strconv"
	"fmt"
	"unicode"
)

// The mapping function is called once for each piece of the input. In this
// framework, the key is the name of the file that is being processed, and the
// value is the file's contents. The return value should be a slice of key/value
// pairs, each represented by a mapreduce.KeyValue.
func mapF(fileName string, contents string) (res []mr.KeyValue) {
	wcm := make(map[string]uint64)

	words := strings.FieldsFunc(contents, func(c rune) bool { return !unicode.IsLetter(c) })

	for _, word := range words {
		wcm[word] += 1
	}

	for k, v := range wcm {
		res = append(res, mr.KeyValue{k, strconv.FormatUint(v, 10)})
	}

	return
}

// The reduce function is called once for each key generated by Map, with a list
// of that key's string value (merged across all inputs). The return value
// should be a single output value for that key.
func reduceF(key string, values []string) string {
	total := uint64(0)

	for _, count := range values {
		x, err := strconv.ParseUint(count, 0, 64)

		if err != nil {
			fmt.Println("WC reduceF: failed to parse %s", count);
			panic(err)
		}

		total += x
	}

	return strconv.FormatUint(total, 10)
}

// Parses the command line arguments and runs the computation. "Running the
// computation" could mean that this node is a worker or a master, depending on
// the command line flags. See `mr/parse_cmd_line.go:parseCmdLine` for details.
func main() {
	mr.Run("wordcount", mapF, reduceF)
}
