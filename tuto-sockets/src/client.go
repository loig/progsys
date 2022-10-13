package main

import (
	"log"
	"net"
)

func main() {

	conn, err := net.Dial("unix", "test.sock")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close()

	log.Println("Je suis connecté")

}
