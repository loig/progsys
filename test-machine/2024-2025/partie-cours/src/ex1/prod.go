package main

/*
Vous ne devez pas modifier ce fichier
*/

import (
	"math/rand"
	"time"
)

func producer(c chan int) {

	for {
		time.Sleep(time.Duration(rand.Intn(2000)+1) * time.Millisecond)
		c <- rand.Intn(100)
	}

}
