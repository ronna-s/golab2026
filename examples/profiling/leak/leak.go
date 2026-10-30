package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func DoNoLeak() {
	go func() {
		time.Sleep(time.Minute)
	}()
}

func DoLeak() {
	go func() {
		ch := make(chan struct{})
		<-ch
	}()
}

func main() {
	DoLeak()
	DoNoLeak()
	log.Fatal(http.ListenAndServe("localhost:6061", nil))
}
