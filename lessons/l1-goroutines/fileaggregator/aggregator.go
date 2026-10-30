package main

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/ronna-s/golab2026/lessons/l1-goroutines/fileaggregator/fileservice"
)

type Client interface {
	// GetParts lists the file parts so that they can be fetched individually.
	// Parts are listed in a random order.
	GetParts() []*fileservice.Part
	// GetPart returns the individual part of the file (a slice of bytes) and its index.
	GetPart(p *fileservice.Part) ([]byte, int)
}

func main() {
	fmt.Println(AggergateFile(fileservice.NewClient()))
}

func AggergateFile(client Client) string {
	var wg sync.WaitGroup
	parts := client.GetParts()
	out := make([][]byte, len(parts))
	for _, part := range parts {
		wg.Go(func() {
			b, idx := client.GetPart(part)
			out[idx] = b
		})
	}
	wg.Wait()
	// your code goes here
	return string(bytes.Join(out, nil))
}
