package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	listener, _ := net.Listen("tcp", ":8081")

	conn, _ := listener.Accept()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Étape 1
	writer.WriteString("Salut\n")
	writer.Flush()

	msg, _ := reader.ReadString('\n')

	if msg != "Bonjour\n" {
		fmt.Println("Échec à l'étape 1")
		return
	}
	fmt.Println("Étape 1 passée")

	// Étape 2
	req := rand.Intn(25) + 5
	writer.WriteString(fmt.Sprint(req) + "\n")
	writer.Flush()

	for count := 0; count < req; count++ {
		msg, _ := reader.ReadString('\n')
		if msg != fmt.Sprint(count)+"\n" {
			fmt.Println("Échec à l'étape 2, tour", count)
			return
		}
	}
	fmt.Println("Étape 2 passée")

	// Étape 3
	turn := rand.Intn(5) + 3
	for i := 0; i < turn; i++ {
		writer.WriteString("Hip\n")
		writer.Flush()
		msg, _ := reader.ReadString('\n')
		if msg != "Hop\n" {
			fmt.Println("Échec à l'étape 3, tour", i)
			return
		}
	}
	conn.Close()

	conn, _ = listener.Accept()
	writer = bufio.NewWriter(conn)
	reader = bufio.NewReader(conn)

	writer.WriteString("De retour\n")
	writer.Flush()

	msg, _ = reader.ReadString('\n')
	if msg != "Enfin !\n" {
		fmt.Println("Échec à l'étape 3, reconnexion")
		return
	}
	fmt.Println("Étape 3 passée")

	// Bravo
	fmt.Println("Bravo !")

}
