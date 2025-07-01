package main

import (
	"time"
)

/*
Vous ne devez pas modifier ce fichier
*/

func generator(c chan timeStep) {

	var newTasksByTime map[int][]task = make(map[int][]task)
	newTasksByTime[0] = []task{
		{id: 0, duration: 20, priority: 3},
		{id: 1, duration: 30, priority: 2},
	}
	newTasksByTime[10] = []task{
		{id: 2, duration: 15, priority: 4},
	}
	newTasksByTime[20] = []task{
		{id: 3, duration: 35, priority: 1},
		{id: 4, duration: 35, priority: 0},
		{id: 5, duration: 35, priority: 3},
	}
	newTasksByTime[40] = []task{
		{id: 6, duration: 20, priority: 1},
		{id: 7, duration: 15, priority: 2},
	}

	for clockValue := 0; ; clockValue++ {
		c <- timeStep{
			clockValue: clockValue,
			newTasks:   newTasksByTime[clockValue],
		}
		time.Sleep(time.Duration(10 * time.Millisecond))
	}

}
