# Paso 02. Recursos empotrados.
En este paso vamos a tener un conjunto de emojis y una fuente empotradas para ser utilizadas en la demostración.

Los emojis los descargamos del proyecto [OpenMoji](https://openmoji.org), diseñados por la [HfG Schwäbisch Gmünd](https://www.hfg-gmuend.de/) y licenciados bajo [Creative Commons Attribution-ShareAlike 4.0 International (CC BY-SA 4.0)](https://creativecommons.org/licenses/by-sa/4.0/).

La fuente la descargamos de Goolge Fonts y es [Luckiest Guy](https://fonts.google.com/specimen/Luckiest+Guy) diseñada por [Astigmatic](https://fonts.google.com/?query=Astigmatic) y con licencia [Apache 2.0](https://fonts.google.com/specimen/Luckiest+Guy/license).

## Punto de entrada para desktop.
Vamos a añadir un punto de entrada para desktop que nos permite desarrollar la demo sin tener que estar compilando la librería y el apk para cada cambio o ajuste.

Añadimos un fichero **main.go** en el raiz del proyecto:

~~~go
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/programatta/demoandroid/game"
)

func main() {
	ebiten.SetWindowSize(420, 760)
	ebiten.SetWindowTitle("Demo Android")

	game := game.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
~~~

## Preparación de recursos.
Bajo el directorio **game** creamos los directorios **assets/fonts** y **assets/images/emojis**:

~~~shell
cd game
mkdir -p assets/fonts
mkdir -p assets/images/emojis
~~~ 

Bajo **assets/fonts** colocamos el fichero **LuckiestGuy-Regular.ttf** y bajo **assets/images/emojis** colocamos 28 emojis.

Cremos los ficheros que van a empotrar los recursos en tiempo de compilación:

~~~shell
touch assets/fonts/fonts.go
touch assets/images/images.go
~~~

### assets/fonts/fonts.go
~~~go
package fonts

import _ "embed"

//go:embed *.ttf
var FontFiles []byte
~~~

### assets/images/images.go
~~~go
package images

import (
	"embed"
)

//go:embed emojis/*.png
var EmojisDataFS embed.FS
~~~

En este punto la estructura del proyecto la tenemos de la siguiente forma:
~~~shell
.
├── bin
│   ├── android
│   │   ├── ...
│   └── android-libs
│       ├── game.aar
│       └── game-sources.jar
├── game
│   ├── assets
│   │   ├── fonts
│   │   │   ├── fonts.go
│   │   │   └── LuckiestGuy-Regular.ttf
│   │   └── images
│   │       ├── emojis
│   │       │   ├── 1F311.png
│   │       │   ├── 1F312.png
│   │       │   ├── 1F313.png
│   │       │   ├── 1F314.png
│   │       │   ├── 1F920.png
│   │       │   ├── 1F921.png
│   │       │   ├── 1F922.png
│   │       │   ├── 1F923.png
│   │       │   ├── 1F98A.png
│   │       │   ├── 1F98B.png
│   │       │   ├── 1F98C.png
│   │       │   ├── 1F98D.png
│   │       │   ├── 2600.png
│   │       │   ├── 2601.png
│   │       │   ├── 2602.png
│   │       │   ├── 2603.png
│   │       │   ├── 2614.png
│   │       │   ├── 2615.png
│   │       │   ├── 2618.png
│   │       │   ├── 2620.png
│   │       │   ├── E040.png
│   │       │   ├── E042.png
│   │       │   ├── E044.png
│   │       │   ├── E045.png
│   │       │   ├── E151.png
│   │       │   ├── E152.png
│   │       │   ├── E153.png
│   │       │   └── E154.png
│   │       └── images.go
│   └── game.go
├── go.mod
├── go.sum
├── main.go
├── mobile
│   └── main.go
└── README.md
~~~

## Adaptación de game.go.
Vamos a añadir código para hacer uso de la fuente y de los emojis para ser presentados en una rejilla. Para ello creamos un paquete **utils** desde **game** para tener funcionalidad para cargar la fuente y los emojis empotrados:

~~~shell
mkdir utils
touch utils/utils.go
~~~

### utils/utils.go
~~~go
package utils

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/programatta/demoandroid/game/assets/fonts"
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
~~~

Actualizamos las dependencias desde el raiz del proyecto para usar las fuentes:
~~~shell
go mod tidy
~~~

Actualizamos **game/game.go** para cargar la fuente, los emojis y hacer la representación en pantalla.
### game/game.go
~~~go
package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/programatta/demoandroid/game/utils"
)

type Game struct {
	textFace *text.GoTextFace
	emojis   []*ebiten.Image
}

func NewGame() *Game {
	textFace := utils.LoadEmbeddedFont(32)
	emojis := utils.LoadEmojiImages()

	return &Game{
		textFace: textFace,
		emojis:   emojis,
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
	g.drawText(screen)
	g.drawEmojis(screen)
}

// Layout determina el tamaño del canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
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
~~~

Al compilar y ejecutar para ambas plataformas (**desktop** y **Android**) tendremos la siguiente salida:
![Salida ejecución demo](./images//paso_02_demo-desktop-android.png)

Podemos observar el detalle de que se corta la imagen:
![Salida ejecución demo cortada](./images/paso_02_demo-desktop-android-02.png)

## Mejora en la visión del juego.
El método encargado de determinar el tamaño del canvas es **Layout**. Actualmente está devolviendo el tamaño indicado:
* en **desktop** a través de la función `ebiten.SetWindowSize(420, 760)` del fichero **main.go**
* en **Android** a través del sistema que detecta el alto y ancho del dispositivo y luego lo escala.

Para que se vea de la misma forma tanto en **desktop** como en **Android** vamos a devolver en **Layout** el tamaño con el que está diseñado nuestro juego, es decir, **420** y **760**.

Creamos un paquete **config** en el raiz del proyecto y dentro de él, un fichero **config.go**:
~~~shell
mkdir config
touch config/config.go
~~~

### config/config.go
~~~go
package config

const GameWindowWidth int = 420
const GameWindowHeight int = 760
~~~

### main.go
~~~go
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
~~~

### game/game.go
~~~go
package game

import (
  ...
  "github.com/programatta/demoandroid/config"
  ...
)
...
// Layout determina el tamaño del canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.GameWindowWidth, config.GameWindowHeight
}
...
~~~

Con este pequeño cambio, ya tenemos una vista igual en **desktop** y **Android**, como se observa en la imagen a continuación:
![Salida ejecución demo arreglada](./images/paso_02_demo-desktop-android-03.png)

A destacar que nos aparecen las dos bandas negras de ajueste de imagen, pero en los siguientes pasos vemos como resolver esto.

Puedes consultar el código de este paso en la rama [paso-02](https://github.com/programatta/demoandroid/tree/paso-02).