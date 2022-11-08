package main

import "sync"

var x, y int
var mx, my sync.Mutex
var wg sync.WaitGroup

func routine1() {
	mx.Lock()
	defer mx.Unlock()
	x--
	my.Lock()
	defer my.Unlock()
	y = x
	wg.Done()
}

func routine2() {
	my.Lock()
	defer my.Unlock()
	y++
	mx.Lock()
	defer mx.Unlock()
	x = y
	wg.Done()
}

func main() {
	for {
		wg.Add(2)
		go routine1()
		go routine2()
		wg.Wait()
	}
}
