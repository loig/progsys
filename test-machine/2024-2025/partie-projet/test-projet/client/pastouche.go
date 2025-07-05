package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Constantes utilisées par le jeu
// - taille de la zone de jeu en pixels
// - réponses possibles pour une proposition (guessXXX)
// - états du jeu (stateXXX)
// - identifiants des boutons autres que les chiffres (buttonXXX)
// - caractères pour l'affichage de la petite animation qui permet de contrôler que le jeu ne bloque pas
const (
	screenWidth  = 155
	screenHeight = 130
)

const (
	guessGood int = iota
	guessBig
	guessSmall
	guessPending
	stateGuess
	stateRes
	stateWon
)

const (
	buttonClear = -1
	buttonBack  = -2
	buttonSend  = -3
)

var progress []byte = []byte{'/', '-', '\\', '|'}

// Type pour définir un jeu :
// - history est l'historique des propositions faites et des réponses obtenues
// - buttons contient l'ensemble des boutons sur lesquels on peut cliquer
// - guessVal contient la proposition actuelle
// - state indique l'état du jeu (en train de faire une proposition, en attente de réponse, partie gagnée)
// - frame est un compteur d'appels de update utilisé pour afficher une petite animation qui permet de voir si le jeu se bloque
type game struct {
	history  []guess
	buttons  []button
	guessVal int
	state    int
	frame    int
}

// Type pour représenter une proposition et la réponse obtenue
type guess struct {
	v   int
	res int
}

// Type pour représenter un bouton du jeu avec
// - une position
// - une taille
// - un contenu à afficher (ex: "8")
// - une fonction à utiliser quand on clique sur le bouton
// - une info pour savoir si la souris est au dessus du bouton ou non
type button struct {
	x, y          int
	width, height int
	content       string
	onClick       func() int
	selected      bool
}

// Méthode Layout pour que *game implémente ebiten.Game
func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Méthode Update pour que *game implémente ebiten.Game
// pour rappel, cette méthode est appelée automatiquement 60 fois par seconde
func (g *game) Update() error {

	g.frame++
	g.frame = g.frame % 4

	mouseX, mouseY := ebiten.CursorPosition()

	if g.state == stateGuess {
		for i := 0; i < len(g.buttons); i++ {
			g.buttons[i].setSelected(mouseX, mouseY)

			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				if g.buttons[i].selected {
					g.updateState(g.buttons[i].onClick())
				}
			}
		}
	} else if g.state == stateRes {
		g.checkGuess()
	}

	return nil
}

// Méthode Draw pour que *game implémente ebiten.Game
// pour rappel, cette méthode est appelée automatiquement 60 fois par seconde
func (g *game) Draw(screen *ebiten.Image) {
	g.drawProgress(screen)
	g.drawHistory(screen)
	g.drawButtons(screen)
	g.drawCurrentGuess(screen)
}

// Fonction d'initialisation d'une partie (en particulier, mise en place des boutons)
func initGame() (g game) {

	g.state = stateGuess

	xPos := 5
	yPos := 5
	buttonSize := 20
	buttonMargin := 5

	for i := 1; i <= 9; i++ {
		g.buttons = append(g.buttons, button{
			content: fmt.Sprintf("%c", byte(i)+48),
			width:   buttonSize, height: buttonSize,
			x:       xPos + (buttonSize+buttonMargin)*((i-1)%3),
			y:       yPos - (buttonSize+buttonMargin)*((i-1)/3) + (buttonSize+buttonMargin)*2,
			onClick: func() int { return i },
		})
	}

	g.buttons = append(g.buttons, button{
		content: "C",
		width:   buttonSize, height: buttonSize,
		x:       xPos,
		y:       yPos + (buttonSize+buttonMargin)*3,
		onClick: func() int { return buttonClear },
	})

	g.buttons = append(g.buttons, button{
		content: "0",
		width:   buttonSize, height: buttonSize,
		x:       xPos + (buttonSize + buttonMargin),
		y:       yPos + (buttonSize+buttonMargin)*3,
		onClick: func() int { return 0 },
	})

	g.buttons = append(g.buttons, button{
		content: "<",
		width:   buttonSize, height: buttonSize,
		x:       xPos + (buttonSize+buttonMargin)*2,
		y:       yPos + (buttonSize+buttonMargin)*3,
		onClick: func() int { return buttonBack },
	})

	g.buttons = append(g.buttons, button{
		content: " Proposer",
		width:   3*buttonSize + 2*buttonMargin, height: buttonSize,
		x:       xPos,
		y:       yPos + (buttonSize+buttonMargin)*4,
		onClick: func() int { return buttonSend },
	})

	return
}

// Méthode utilitaire pour indiquer si la souris est au dessus d'un bouton
// utilisée par Update
func (b *button) setSelected(mouseX, mouseY int) {
	b.selected = mouseX >= b.x && mouseX <= b.x+b.width && mouseY >= b.y && mouseY <= b.y+b.height
}

// Méthode utilitaire appelée par Update permettant de mettre à jour l'état du jeu en fonction
// du click sur un bouton
func (g *game) updateState(update int) {

	switch update {
	case buttonBack:
		g.guessVal = g.guessVal / 10
	case buttonClear:
		g.guessVal = 0
	case buttonSend:
		g.sendGuess()
	default:
		if g.guessVal < 10 {
			g.guessVal = g.guessVal*10 + update
		}
	}
}

// Méthode utilitaire permettant d'afficher l'animation pour contrôler que le jeu ne bloque pas
// appelée par Draw
func (g *game) drawProgress(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%c", progress[g.frame]), 145, 112)
}

// Méthode utilitaire pour afficher l'historique des propositions faites depuis le début de la partie
// appelée par Draw
func (g *game) drawHistory(screen *ebiten.Image) {

	xstart := 85
	ystart := 44
	ystep := 12
	ymargin := 5
	maxHistory := 5

	x, y := xstart, ystart

	ebitenutil.DebugPrintAt(screen, "Historique", x, y)

	y += ymargin

	for i := len(g.history) - 1; i >= 0 && maxHistory > 0; i-- {
		maxHistory--
		y += ystep
		sign := '='
		switch g.history[i].res {
		case guessBig:
			sign = '>'
		case guessSmall:
			sign = '<'
		}
		display := fmt.Sprintf("%d %c x", g.history[i].v, sign)
		ebitenutil.DebugPrintAt(screen, display, x, y)
	}

}

// Méthode utilitaire pour afficher la prochaine proposition qui sera faite
// appelée par Draw
func (g *game) drawCurrentGuess(screen *ebiten.Image) {
	xstart := 85
	ystart := 5
	ystep := 17

	ebitenutil.DebugPrintAt(screen, "Proposition", xstart, ystart)

	disp := "Gagné !"
	if g.state != stateWon {
		disp = fmt.Sprintf("x = %d", g.guessVal)
	}
	ebitenutil.DebugPrintAt(screen, disp, xstart, ystart+ystep)

}

// Méthode utilitaire pour afficher les boutons
// appelée par Draw
func (g *game) drawButtons(screen *ebiten.Image) {
	for _, b := range g.buttons {
		b.draw(screen, g.state == stateWon)
	}
}

// Méthode utilitaire pour l'affichage d'un bouton donné
// appelée par Draw via drawButtons
func (b button) draw(screen *ebiten.Image, won bool) {

	xshift := 6
	yshift := 2

	if b.selected && !won {
		vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.width), float32(b.height), color.Gray{128}, false)
	}

	vector.StrokeRect(screen, float32(b.x), float32(b.y), float32(b.width), float32(b.height), 1, color.White, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s", b.content), b.x+xshift, b.y+yshift)
}
