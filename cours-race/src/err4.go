package main

import (
	"log"
	"sync"
)

var wg sync.WaitGroup
var mux sync.Mutex

var x int

func routine() {
	mux.Lock()
	for x < 5 {
		log.Print(x)
		x++
	}
	x = 0
	wg.Done()
}

func autre() {
	for {

	}
	wg.Done()
}

func main() {
	wg.Add(3)
	go routine()
	go routine()
	go autre()
	wg.Wait()
}
