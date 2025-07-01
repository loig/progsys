package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
)

func main() {

	listener, _ := net.Listen("tcp", ":8081")

	conn, _ := listener.Accept()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	var msgSeparator byte = '\n'

	// Étape 1
	writer.WriteString(fmt.Sprint("Bonjour", string(msgSeparator)))
	writer.Flush()

	reader.ReadString(msgSeparator)

	fmt.Println("Étape 1 passée")

	// Étape 2
	pass := []string{
		"zéro",
		"un",
		"deux",
		"trois",
		"quatre",
		"cinq",
	}
	req := rand.Intn(5) + 1
	writer.WriteString(fmt.Sprint(req, string(msgSeparator)))
	writer.Flush()

	for count := 0; count < req; count++ {
		msg, _ := reader.ReadString(msgSeparator)
		log.Print("(", msg, ")")
		log.Print("(", pass[req]+string(msgSeparator), ")")
		if msg != pass[req]+string(msgSeparator) {
			fmt.Println("Échec à l'étape 2, tour", count)
			return
		}
	}
	fmt.Println("Étape 2 passée")

	// Étape 3
	turn := rand.Intn(5) + 3
	for i := 0; i < turn; i++ {
		writer.WriteString(fmt.Sprint("Ping", string(msgSeparator)))
		writer.Flush()
		msg, _ := reader.ReadString(msgSeparator)
		if msg != fmt.Sprint("Pong", string(msgSeparator)) {
			fmt.Println("Échec à l'étape 3, tour", i)
			return
		}
	}
	conn.Close()

	conn, _ = listener.Accept()
	writer = bufio.NewWriter(conn)
	reader = bufio.NewReader(conn)

	writer.WriteString(fmt.Sprint("De retour", string(msgSeparator)))
	writer.Flush()

	msg, _ := reader.ReadString(msgSeparator)
	if msg != fmt.Sprint("Ouf", string(msgSeparator)) {
		fmt.Println("Échec à l'étape 3, reconnexion")
		return
	}
	fmt.Println("Étape 3 passée")

	// Bravo
	fmt.Println("Bravo !")

}
