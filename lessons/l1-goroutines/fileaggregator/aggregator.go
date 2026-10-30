package main

import (
	"fmt"

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
	// your code goes here
	return ""
}
