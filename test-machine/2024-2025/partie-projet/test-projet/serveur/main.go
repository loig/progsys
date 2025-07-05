package main

import (
	"bufio"
	"log"
	"net"
	"strconv"

	"math/rand/v2"
)

func main() {

	answer := rand.IntN(100)

	log.Print("Le nombre à deviner est ", answer)

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal("Le serveur n'arrive pas à se mettre en écoute : ", err)
	}

	log.Print("Le serveur écoute")

	co, err := listener.Accept()
	if err != nil {
		log.Fatal("Le serveur n'arrive pas à accepter la connexion : ", err)
	}

	log.Print("Une connexion a été établie")

	reader := bufio.NewReader(co)
	writer := bufio.NewWriter(co)

	for {
		proposition, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Problème lors de la réception d'un message : ", err)
		}

		proposition = proposition[:len(proposition)-1]

		valProposition, err := strconv.Atoi(proposition)
		if err != nil {
			log.Fatal("La proposition reçue (", proposition, ") n'est pas un entier : ", err)
		}

		log.Print("J'ai reçu la proposition suivante : ", proposition)

		reponse := "Bravo"

		if valProposition < answer {
			reponse = "TropPetit"
		}

		if valProposition > answer {
			reponse = "TropGrand"
		}

		log.Print("Je vais donc répondre ", reponse)

		_, err = writer.WriteString(reponse + "\n")
		if err != nil {
			log.Fatal("Problème lors de l'écriture d'un message : ", err)
		}

		err = writer.Flush()
		if err != nil {
			log.Fatal("Problème lors de l'envoi d'un message : ", err)
		}

		log.Print("La réponse est envoyée")

		if valProposition == answer {
			log.Print("C'est la bonne réponse, le jeu doit se terminer")
			break
		}
	}

	_, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal("Problème lors de la réception du dernier message : ", err)
	}

	log.Print("C'est fini, bravo !")
}
