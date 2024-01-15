package main

import (
	"fmt"
	"io"
	"os"

	"github.com/cespare/xxhash"
	"github.com/depot/orca/util/chunker"
)

func main() {
	file, _ := os.Open("../example.tar")
	defer file.Close()

	numChunks := 0
	chunker := chunker.NewChunker(file)

	for {
		data, err := chunker.Next()
		if err != nil && err != io.EOF {
			break
		}
		hash := xxhash.Sum64(data)
		fmt.Printf("len(data): %d %x\n", len(data), hash)
		numChunks += 1

		if err == io.EOF {
			break
		}
	}

	fmt.Printf("numChunks: %d\n", numChunks)
}
