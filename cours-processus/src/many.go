package main

import (
	"fmt"
	"time"
)

const maxroutine = 10000

func echoID(ID int) {
	for {
		fmt.Println("Je suis la gouroutine numéro", ID)
		time.Sleep(time.Millisecond)
	}
}

func main() {
	for i := 1; i < maxroutine; i++ {
		go echoID(i)
	}

	echoID(maxroutine)
}
