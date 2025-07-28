package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/programatta/demoandroid/config"
	"github.com/programatta/demoandroid/game"
)

func main() {
	ebiten.SetWindowSize(config.GameWindowWidth, config.GameWindowHeight)
	ebiten.SetWindowTitle("Demo Android")

	game := game.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
