# Paso 04. Eliminación de las bandas negras.

La aplicación al ejecutarse, deja unas bandas negras debido al ajuste realizado por **Ebitengine** al devolver las dimensiones lógicas que estamos utilizando para representar el area de juego.

Se pude crear una imagen intermedia llamada **offscreen** con las dimensiones lógicas de nuestro juego y que el método `Layout()` devuelva las dimensiones detectadas por el sistema.

En el método `Draw()` vamos a limpiar y pintar **offscreen** de verde y la vamos a ubicar en el centro de la pantalla.

### game/game.go
~~~go
package game
...
type Game struct {
  ...
  offscreen *ebiten.Image
}

func NewGame() *Game {
  ...
  offscreen := ebiten.NewImage(config.GameWindowWidth, config.GameWindowHeight)

	return &Game{
    ...
    offscreen: offscreen,
	}
}

...

// Draw dibuja el estado actual.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0xcf, 0xba, 0xf0, 0xff})

	g.offscreen.Clear()
	g.offscreen.Fill(color.NRGBA{0x00, 0xff, 0x00, 0xff})

	// g.drawText(screen)
	// g.drawEmojis(screen)

	offscreenOpt := &ebiten.DrawImageOptions{}
	offscreenOpt.GeoM.Translate(
		float64(screen.Bounds().Dx()-g.offscreen.Bounds().Dx())/2,
		float64(screen.Bounds().Dy()-g.offscreen.Bounds().Dy())/2,
	)
	screen.DrawImage(g.offscreen, offscreenOpt)
}

// Layout determina el tamaño del canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
  return outsideWidth, outsideHeight
}
...
~~~

![zona-dibujo](./images/paso_04_demo-android-zona-dibujo.png)

De esta forma, eliminamos las bandas negras y ya tenemos la pantalla cubierta de nuevo con el color de la aplicación. En la sección verde, es donde vamos a renderizar nuestro juego, por lo que pasamos a las funciones **g.drawText()** y **g.drawEmojis** la nueva imagen sobre la que pintar **g.offscreen**:

### game/game.go
~~~go
package game
...
type Game struct {
  ...
  offscreen *ebiten.Image
}

func NewGame() *Game {
  ...
  offscreen := ebiten.NewImage(config.GameWindowWidth, config.GameWindowHeight)

	return &Game{
    ...
    offscreen: offscreen,
	}
}

...

// Draw dibuja el estado actual.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0xcf, 0xba, 0xf0, 0xff})

	g.offscreen.Clear()
	g.offscreen.Fill(color.NRGBA{0x00, 0xff, 0x00, 0xff})

	g.drawText(g.offscreen)
	g.drawEmojis(g.offscreen)

	offscreenOpt := &ebiten.DrawImageOptions{}
	offscreenOpt.GeoM.Translate(
		float64(screen.Bounds().Dx()-g.offscreen.Bounds().Dx())/2,
		float64(screen.Bounds().Dy()-g.offscreen.Bounds().Dy())/2,
	)
	screen.DrawImage(g.offscreen, offscreenOpt)
}

// Layout determina el tamaño del canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
  return outsideWidth, outsideHeight
}
...
~~~

![zona-dibujo-corta](./images/paso_04_demo-android-zona-dibujo-corta.png)

Ya tenemos el juego renderizado en la zona verde, pero vemos que se sale por los extremos de la pantalla. Para solucionar esto, vamos a escalar las dimensiones recibidas por el sistema en el método `Layout()` con relación a las dimensiones lógicas de nuestro juego y las vamos a aplicar a la imagen **offscreen**.

Esta operación la vamos a realizar en el método `Layout()`.

### game/game.go
~~~go
package game
...
type Game struct {
  ...
	scaled    bool
	offscreenOpt *ebiten.DrawImageOptions
}
...
// Draw dibuja el estado actual.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0xcf, 0xba, 0xf0, 0xff})

	g.offscreen.Clear()
	g.offscreen.Fill(color.NRGBA{0x00, 0xff, 0x00, 0xff})

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
	}
	return outsideWidth, outsideHeight
}
...
~~~

![zona-dibujo-escalado-1](./images/paso_04_demo-android-zona-dibujo-escalada-01.png)

Con esta funcionalidad, ya disponemos de la vista de juego centrada, aunque los dibujos se ven un poco borrosos. Para solucionarlo, aplicamos el filtro `ebiten.FilterLinear` a las opciones de dibujo **g.offscreenOpt** establecidas en el método `Layout()`:

### game/game.go
~~~go
package game
...
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
...
~~~


![zona-dibujo-escalado-2](./images/paso_04_demo-android-zona-dibujo-escalada-02.png)

Eliminando el relleno de color verde de la imagen **g.offscreen** y asignando el color de fondo de **screen** ya tendríamos la aplicación lista.

### game/game.go
~~~go
package game
...
// Draw dibuja el estado actual.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0xcf, 0xba, 0xf0, 0xff})

	g.offscreen.Clear()
	g.offscreen.Fill(color.NRGBA{0xcf, 0xba, 0xf0, 0xff})
  ...
}
~~~

![zona-dibujo-escalado-3](./images/paso_04_demo-android-zona-dibujo-escalada-03.png)

Puedes consultar el código de este paso en la rama [paso-04-01](https://github.com/programatta/demoandroid/tree/paso-04-01).


## Completar el color en la zona del "notch".
Podemos observar, que en la zona donde se encuentra la barra de estado, o en la zona superior o **notch** se queda sin pintar. Para solventar esto, añadimos al estilo **SplashTheme** que estamos usando en el fichero **styles.xml** el siguiente elemento `android:windowLayoutInDisplayCutoutMode` que admite valores `default`, `shortEdges`, o `never` y le asignamos `shortEdges`.

### android/app/src/main/res/values/styles.xml
~~~xml
<resources>
  <!-- Base application theme. -->
  <style name="AppTheme" parent="Theme.AppCompat.Light.DarkActionBar">
    <!-- Customize your theme here. -->
    ...
  </style>

  <!-- Splash theme -->
  <style name="SplashTheme" parent="Theme.AppCompat.NoActionBar">
    ...
    <item name="android:windowLayoutInDisplayCutoutMode">shortEdges</item>
  </style>
</resources>
~~~

Una vez compilado el APK e instalado en el dispositivo / emulador la aplicación se muestra de la siguiente manera:

![zona-dibujo-final-ok](./images/paso_04_demo-android-final-ok.png)

> ⚠️ **Atención.**
>
> Un problema que se observa, es que al arrancar la aplicación, los iconos de **status bar** y de **navigation bar** comienzan en un color blanco, y al entrar en modo inmersivo, en **status bar** cambia los iconos a negro, mientras que en **navigation bar** queda en blanco, lo que produce un parpadeo extraño.
>
> No he encontrado una solución a esto, salvo oscureciendo el color del splash screen pasando de **#CFBAF0** a **#7A5B9C**. Con este cambio, al arrancar la aplicación, los iconos de **status bar** y de **navigation bar** comienzan en un color blanco y al entrar en modo inmersivo se quedan igual con lo que no se produce el parpadeo extraño.

Puedes consultar el código de este paso en la rama [paso-04-02](https://github.com/programatta/demoandroid/tree/paso-04-02).
