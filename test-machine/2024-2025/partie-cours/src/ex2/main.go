package main

// Vous ne devez pas modifier ce fichier

type timeStep struct {
	clockValue int
	newTasks   []task
}

type task struct {
	id       int
	duration int
	priority int // plus la priorité est proche de 0, plus la tâche est prioritaire
}

func main() {

	var c chan timeStep

	go scheduler(c)
	generator(c)

}
