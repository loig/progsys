package main

/*
Vous ne devez pas modifier ce fichier
*/

func main() {

	var c chan int = make(chan int, 1000)

	go consumer(c)
	producer(c)

}
