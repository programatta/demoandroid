package mobile

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/programatta/demoandroid/game"
)

func init() {
	mobile.SetGame(game.NewGame())
}

// At least one exported function is required by gomobile.
func Dummy() {}
