package main

import (
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"

	"net"
	"fmt"
	"bufio"
	"strconv"
)

// Variable globale qui contient le nombre qu'on doit deviner,
// celle-ci n'exitera plus dans la version finale du jeu
var answer int

func handleRead(reader *bufio.NewReader, c chan string){
	for {
		msg,err := reader.ReadString("\n")
		if err != nil {
			log.Println("erreur lecture message")
		}
		c <- msg
	}
}

func handleWrite(writer *bufio.NewWriter, c chan string){
	for {
		msg := <- c 
		_,err := writer.WriteString(msg)
		if err != nil {
			log.Println("erreur envoie message")
		}
		err = writer.Flush()
		if err != nil{
			log.Println("erreur envoie de message")
		}
	}
}

var gameReader chan string
var gameWriter chan string

// Initialisation et lancement du jeu
// Cette fonction est à modifier, mais seulement à l'endroit indiqué
func main() {

	// Initialisation de la connection
	conn, err := net.Dial("tcp",":8082")
	if err != nil {
		fmt.Println("Erreur de connexion:",err)
	}
	defer conn.Close()

	gameReader = make(chan string,1)
	gameWriter = make(chan string,1)
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	go handleRead(reader, gameReader)
	go handleWrite(writer, gameWriter)

	// Initialisation du jeu et de la fenêtre, ne pas toucher
	theGame := initGame()
	ebiten.SetWindowTitle("Le jeu du DS")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowSize(2*screenWidth, 2*screenHeight)

	// À partir d'ici vous pouvez modifier la fonction main (y compris en ajoutant des lignes)
	answer = rand.IntN(100)
	// À partir d'ici vous ne pouvez plus modifier la fonction main

	// Lancement du jeu, ne pas toucher
	if err := ebiten.RunGame(&theGame); err != nil {
		log.Fatal(err)
	}

}

// Envoi d'une proposition, cette méthode est appelée quand on clique sur le bouton "Proposer"
// (à condition que g.state soit égal à stateGuess)
// Pour l'instant cette fonction ne fait qu'enregistrer la proposition dans l'historique et changer le champs
// state de g pour indiquer qu'on attend de savoir si la proposition est bonne (stateRes)
// Cette fonction est à modifier
func (g *game) sendGuess() {
	g.history = append(g.history, guess{
		v:   g.guessVal,
		res: guessPending,
	})
	gameWriter <- strconv.Atoi(g.guessVal)

	g.state = stateRes
}

// Vérification d'une proposition, cette méthode est appelée 60 fois par seconde (par la méthode Update)
// à partir du moment où g.state vaut stateRes
// Pour l'instant la validité de la proposition est vérifiée immédiatement à partir de la variable globale
// answer et le dernier enregistrement de l'historique (qui représente la dernière proposition faite) est
// mis à jour en fonction du résultat de la comparaison
// En fonction du résultat de cette comparaison g.state est aussi mis à jour, soit pour indiquer que la partie
// est terminée (le joueur à gagné) avec stateWon, soit pour indiquer qu'il faut faire une autre proposition
// avec stateGuess
// Cette fonction est à modifier
func (g *game) checkGuess() {

	res := guessGood
	
	m := <- reader

	select{
	case m == "TropPetit":
		res = guessSmall
	case m == "TropGrand":
		res = guessBig
	case m == "Bravo":
		g.state = stateWon
	default:
	}

	if g.history[len(g.history)-1].v > answer {
		res = guessBig
	} else if g.history[len(g.history)-1].v < answer {
		res = guessSmall
	}

	g.history[len(g.history)-1].res = res

	if res == guessGood {
		g.state = stateWon
	} else {
		g.guessVal = 0
		g.state = stateGuess
	}

}
