package main

import (
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("unix", "test.sock")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Println("accept error:", err)
		return
	}
	defer conn.Close()
	log.Println("Le client s'est connecté")

	time.Sleep(10 * time.Second)

}
