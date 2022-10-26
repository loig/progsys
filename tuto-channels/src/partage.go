package main

import "fmt"

var turn1 bool

func routine1() {
	for {
		if turn1 {
			fmt.Println("Routine 1")
			turn1 = false
		}
	}
}

func routine2() {
	for {
		if !turn1 {
			fmt.Println("Routine 2")
			turn1 = true
		}
	}
}

func main() {
	go routine1()
	routine2()
}
