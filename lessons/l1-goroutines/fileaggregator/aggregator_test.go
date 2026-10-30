package main

import (
	_ "embed"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/ronna-s/golab2026/lessons/l1-goroutines/fileaggregator/fileservice"
	"github.com/stretchr/testify/assert"
)

var delay = time.Second

const file = "012345"

type FakeClient struct {
	t     *testing.T
	parts []*fileservice.Part
}

func NewFakeClient() *FakeClient {
	s := make([]*fileservice.Part, len(file))
	for i := range s {
		s[len(s)-1-i] = &fileservice.Part{}
	}
	return &FakeClient{parts: s}
}

func (c FakeClient) GetParts() []*fileservice.Part {
	return c.parts
}

func (c FakeClient) GetPart(p *fileservice.Part) ([]byte, int) {
	idx := slices.Index(c.parts, p)
	if idx < 0 {
		c.t.Fatalf("part doesn't exist")
		return nil, -1
	}
	time.Sleep(delay)
	return []byte{file[idx]}, idx
}

func TestAggergateFile(t *testing.T) {
	client := NewFakeClient()
	expectedWait := delay
	done := false
	var output string

	synctest.Test(t, func(t *testing.T) {
		go func() {
			ts := time.Now()
			output = AggergateFile(client)
			assert.LessOrEqual(t, time.Since(ts), expectedWait, "waited too long...")
			done = true
		}()

		time.Sleep(time.Hour)
		synctest.Wait()
		if !done {
			t.Error("waited an hour, but AggergateFile didn't finish :(")
		}
	})

	assert.Equal(t, file, output)
}
