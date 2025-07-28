package utils

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/programatta/demoandroid/game/assets/fonts"
	"github.com/programatta/demoandroid/game/assets/images"
)

func LoadEmbeddedFont(size float64) *text.GoTextFace {
	faceSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.FontFiles))
	if err != nil {
		log.Fatal(err)
	}

	return &text.GoTextFace{
		Source: faceSource,
		Size:   size,
	}
}

func LoadEmojiImages() []*ebiten.Image {
	var emojis []*ebiten.Image

	dirEntries, errDirEntries := images.EmojisDataFS.ReadDir("emojis")
	if errDirEntries != nil {
		panic(errDirEntries)
	}

	for _, dirEntry := range dirEntries {
		emojiName := fmt.Sprintf("emojis/%s", dirEntry.Name())
		emojiBytes, err := images.EmojisDataFS.ReadFile(emojiName)
		if err != nil {
			panic(err)
		}
		emojis = append(emojis, GenerateImage(94, 94, emojiBytes))
	}
	return emojis
}

func GenerateImage(width, heigh int, emojiBytes []byte) *ebiten.Image {
	imgTmp, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(emojiBytes))
	if err != nil {
		panic(err)
	}

	opTmp := &ebiten.DrawImageOptions{}
	opTmp.GeoM.Translate(float64(width)/2-float64(imgTmp.Bounds().Dx())/2, float64(heigh)/2-float64(imgTmp.Bounds().Dy())/2)

	img := ebiten.NewImage(width, heigh)
	img.Fill(color.White)
	img.DrawImage(imgTmp, opTmp)

	return img
}
