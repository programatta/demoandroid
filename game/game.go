package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/programatta/demoandroid/config"
	"github.com/programatta/demoandroid/game/utils"
)

type Game struct {
	textFace     *text.GoTextFace
	emojis       []*ebiten.Image
	offscreen    *ebiten.Image
	scaled       bool
	offscreenOpt *ebiten.DrawImageOptions
}

func NewGame() *Game {
	textFace := utils.LoadEmbeddedFont(32)
	emojis := utils.LoadEmojiImages()
	offscreen := ebiten.NewImage(config.GameWindowWidth, config.GameWindowHeight)

	return &Game{
		textFace:  textFace,
		emojis:    emojis,
		offscreen: offscreen,
	}
}

// ----------------------------------------------------------------------------
// Implementa Ebiten Game Interface
// ----------------------------------------------------------------------------

// Update realiza el cambio de estado si es necesario y permite procesar
// eventos y actualizar su lógica.
func (g *Game) Update() error {
	return nil
}

// Draw dibuja el estado actual.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0xcf, 0xba, 0xf0, 0xff})

	g.offscreen.Clear()
	g.offscreen.Fill(color.NRGBA{0xcf, 0xba, 0xf0, 0xff})

	g.drawText(g.offscreen)
	g.drawEmojis(g.offscreen)

	screen.DrawImage(g.offscreen, g.offscreenOpt)
}

// Layout determina el tamaño del canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if !g.scaled {
		g.scaled = true

		scaleX := float64(outsideWidth) / float64(config.GameWindowWidth)
		scaleY := float64(outsideHeight) / float64(config.GameWindowHeight)
		scale := math.Min(scaleX, scaleY)

		g.offscreenOpt = &ebiten.DrawImageOptions{}
		g.offscreenOpt.GeoM.Scale(scale, scale)
		g.offscreenOpt.GeoM.Translate(
			(float64(outsideWidth)-float64(g.offscreen.Bounds().Dx())*scale)/2,
			(float64(outsideHeight)-float64(g.offscreen.Bounds().Dy())*scale)/2,
		)
		g.offscreenOpt.Filter = ebiten.FilterLinear
	}
	return outsideWidth, outsideHeight
}

func (g *Game) drawText(screen *ebiten.Image) {
	messageText := "Hola Android desde Go"
	widthText, _ := text.Measure(messageText, g.textFace, 0)
	opText := &text.DrawOptions{}
	opText.GeoM.Translate(float64(screen.Bounds().Dx())/2-widthText/2, 10)
	opText.ColorScale.ScaleWithColor(color.NRGBA{0xff, 0x00, 0x00, 0xff})
	text.Draw(screen, messageText, g.textFace, opText)
}

func (g *Game) drawEmojis(screen *ebiten.Image) {
	for pos, emoji := range g.emojis {
		x := pos % 4
		y := int(pos / 4)

		emojiOp := &ebiten.DrawImageOptions{}
		emojiOp.GeoM.Translate(float64(x*101+12), float64(y*101+48))
		emojiOp.ColorScale.ScaleAlpha(0.5)
		screen.DrawImage(emoji, emojiOp)

		vector.StrokeRect(screen, float32(x*101)+9, float32(y*101)+45, 100, 100, 1, color.NRGBA{0xff, 0x00, 0x00, 0xff}, true)
	}
}
