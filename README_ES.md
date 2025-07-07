# Proyecto Demo Go + Ebiten para Android.
Este documento guía paso a paso cómo generar una aplicación Android (**APK**) usando **Go + Ebiten**, generando un archivo **.aar** sin necesidad de **Android Studio**.

> 📢 Este paso a paso está basado en el tutorial de [Saffron Dionysius, Can You Really Develop Android Apps Without Android Studio?](https://medium.com/@sdiony/can-you-really-develop-android-apps-without-android-studio-cdd9b951de65)

## ✅ Requisitos.
Debemos tener en el entorno de desarrollo configurado lo siguiente:
Herramienta | Versión utilizada
------------|------------------
Golang | 1.24
Ebiten | github.com/hajimehoshi/ebiten/v2
Java SDK | 17
Android |  Commandline tools: [commandlinetools-linux-11076708_latest.zip](https://dl.google.com/android/repository/commandlinetools-linux-11076708_latest.zip)
Gradle  | [gradle-8.14.2-bin.zip](https://services.gradle.org/distributions/gradle-8.14.2-bin.zip)

>💡 Desarrollo realizado en **Debian 12** usando **VSCode** junto con un [**devcontainer personalizado**](https://github.com/programatta/devcontainers/tree/master/goebitendevcontainer) y un docker auxiliar con los requisitos indicados [**go-android**](https://github.com/programatta/toolscontainers/tree/master/go-android).


## 🚀 Estructura general del proyecto.
La estructura que va a presentar el proyecto es la siguiente:
~~~shell
.
├── bin
│   ├── android-libs      # Librerías generadas (.aar, .jar)
│   └── android           # Proyecto Android (Gradle)
├── game                  # Código del juego (Ebiten)
├── mobile                # Entrada para gomobile/ebitenmobile
├── go.mod
└── go.sum
~~~

La vamos a ir construyendo paso a paso.

## 🕹️ Creación del proyecto Go con Ebiten.
Creamos un modulo de android de la forma habitual y añadimos la librería ebiten.

~~~shell
mkdir demoandroid
cd demoandroid
go mod init github.com/programatta/demoandroid
go get github.com/hajimehoshi/ebiten/v2
~~~

Vamos a crear un pequeño programa que pinte la pantalla de un color morado claro y un texto con `ebitenutil.DebugPrint()`.

Creamos un paquete **game** que va acontener lo básico del juego, que es la pantalla morada.

~~~shell
mkdir game
~~~

y dentro del paquete creamos el fichero **game.go**:

~~~go
package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func NewGame() *Game {
	return &Game{}
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
  ebitenutil.DebugPrint(screen, "Hola Android desde Go!")
}

// Layout determina el tamaño del canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
~~~

Ahora solo tenemos implementada la interfáz de Ebiten Game, donde en la función **Draw()** establecemos el color morado claro.

Para llevar nuestro juego a Android, vamos a creamos un paquete **mobile** que va a contener el punto de entrada del código golang para android cuando se generen las librerías nativas __*.so__.

~~~shell
mkdir mobile
~~~

y dentro de ese paquete creamos fichero **main.go**:

~~~go
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
~~~

Aquí la principal diferencia es que no llamamos a **ebiten.RunGame()**, sino a **mobile.SetGame()** y además se encuentra en una función **init()** no en una función **main()**.

Hasta este momento, la estructura de directorios del proyecto es la siguiente:

~~~shell
.
├── game
│   └── game.go
├── mobile
│   └── main.go
├── go.mod
├── go.sum
└── README.md
~~~

## 📦 Creación de la librería de Android.
Para llevar nuestro juego a Android, el código **Go** debe ser compilado y transformado en librerías dinámicas **.so** y cargadas por **Java/Kotlin** a través de unas clases auxiliares, y este conjunto de ficheros son almacenados en un fichero **.aar** (android arquive).

Este proceso se realiza gracias a la utilidad **ebitenmobile** que se apoya en **gomobile**, por lo que necesitamos instalar **ebitenmobile**:

~~~shell
go install github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile@latest
~~~

Una vez instalado, vamos a crear un directorio **bin/android-libs** donde dejaremos las librerías generadas por **ebitenmobile** para luego incluirlas en un proyecto android. La estructura nos va quedando así:

~~~shell
.
├── bin
│   └── android-libs
├── game
│   └── game.go
├── mobile
│   └── main.go
├── go.mod
├── go.sum
└── README.md
~~~

Ejecutamos desde el raiz del proyecto **ebitenmobile** indicandole que va a ser para android, el identificador del paquete que va a tener la librería, donde lo vamos a dejar y de que paquete tomamos el punto de entrada:

~~~shell
ebitenmobile bind -target android -javapkg com.programatta.games.demoandroid.corelib -o bin/android-libs/game.aar github.com/programatta/demoandroid/mobile
~~~

Nota a tener en cuenta, es que no llamemos con el mismo nombre de paquete a la librería creada y al proyecto android que vamos a crear. En un principio si se puede, pero pueden aparecer conflictos por lo que en este ejemplo, la librería va a llevar el nombre de paquete **com.programatta.games.demoandroid.corelib** y el futuro proyecto android **com.programatta.games.demoandroid**.

Tras la ejecución, si ha ido todo bien, nos aparecerán dos ficheros en **bin/android-libs**:
* game.aar
* game-sources.jar

El fichero importante es **game.aar**, va a tener el código **Go** compilado en ficheros **.so** para las plataformas soportadas por android (arm64-v8a, armeabi-v7a, x86 y x86_64) bajo el directorio **jni**. Aparte contiene un fichero **classess.jar** (que es el **game-sources.jar** compilado) con las clases auxiliares que van a parmitir acceder al juego cargando el correspondiente **.so** según plataforma.

## 🛠️ Creación del proyecto Android (Gradle).
Bajo el directorio **bin** vamos a crear uno nuevo llamado **android**:

~~~shell
mkdir bin/android
~~~

Con lo que tenemos la estructura siguiente:
~~~shell
.
├── bin
│   ├── android
│   └── android-libs
│       ├── game.aar
│       └── game-sources.jar
├── game
│   └── game.go
├── mobile
│   └── main.go
├── go.mod
├── go.sum
└── README.md
~~~

nos situamos en el directorio **bin/android** y creamos un projecto **Java** (no llega a ser **Android** ya que necesitaremos ir añadiendo directorios, archivos y alguna que otra modificación para que sea un proyecto **Android**) a través de **gradle**:

~~~shell
gradle init --type java-application --dsl groovy --package com.programatta.games.demoandroid --project-name "android" --no-split-project --java-version 17 --use-defaults
~~~

Al realizar esto, la estructura del directorio **android** se presenta así:
~~~shell
.
├── app
│   ├── build.gradle
│   └── src
│       ├── main
│       │   ├── java
│       │   │   └── com
│       │   │       └── programatta
│       │   │           └── games
│       │   │               └── demoandroid
│       │   │                   └── App.java
│       │   └── resources
│       └── test
│           ├── java
│           │   └── com
│           │       └── programatta
│           │           └── games
│           │               └── demoandroid
│           │                   └── AppTest.java
│           └── resources
├── gradle
│   ├── libs.versions.toml
│   └── wrapper
│       ├── gradle-wrapper.jar
│       └── gradle-wrapper.properties
├── gradle.properties
├── gradlew
├── gradlew.bat
└── settings.gradle
~~~

### 🪛 Cambios que necesitamos.
Creamos el directorio **libs** dentro de **app** donde colocaremos la librería creada en el paso anterior **game.aar** y que usaremos para ir montando la vista de **Android**:

~~~shell
mkdir app/libs
cp ../android-libs/game.aar app/libs
~~~

Renombramos el fichero **App.java** que se encuentra en **app/src/main/java/com/programatta/games/demoandroid** como **MainActivity.java** y lo mimso hacemos con el directorio **resources** que se encuentra en **app/src/main** y lo llamamos **res**:

~~~shell
mv app/src/main/java/com/programatta/games/demoandroid/App.java app/src/main/java/com/programatta/games/demoandroid/MainActivity.java
mv app/src/main/resources app/src/main/res
~~~

Vamos a crear una vista que va a contener el juego y que va a heredar de **EbitenView**, creamos el fichero **EbitenViewWithHerrorHandling.java** en el mismo paquete que está **MainActivity.java**:

~~~shell
touch app/src/main/java/com/programatta/games/demoandroid/EbitenViewWithErrorHandling.java
~~~

Bajo el directorio **res** crearemos otros dos, directorios **layout** y **values** para la definición de la vista y colores:

~~~shell
mkdir app/src/main/res/layout
mkdir app/src/main/res/values
~~~

En **layout** creamos el fichero **activity_main.xml** y en **values** creamos **colors.xml** y **styles.xml**:

~~~shell
touch app/src/main/res/layout/activity_main.xml
touch app/src/main/res/values/colors.xml       
touch app/src/main/res/values/styles.xml
~~~

Añadimos el fichero **AndroidManifest.xml** en **app/src/main**:

~~~shell
touch app/src/main/AndroidManifest.xml
~~~

Añadimos el fichero **build.gradle** en el raiz del proyecto android:
~~~shell
touch build.gradle
~~~

Y finalmente  eliminamos el directorio de **test**, no lo vamos a necesitar por ahora:

~~~shell
rm -rf app/src/test 
~~~

Con todo esto, tenemos la siguiente estructura del proyecto android:
~~~shell
.
├── app
│   ├── build.gradle
│   ├── libs
│   │   └── game.aar
│   └── src
│       └── main
│           ├── AndroidManifest.xml
│           ├── java
│           │   └── com
│           │       └── programatta
│           │           └── games
│           │               └── demoandroid
│           │                   ├── EbitenViewWithErrorHandling.java
│           │                   └── MainActivity.java
│           └── res
│               ├── layout
│               │   └── activity_main.xml
│               └── values
│                   ├── colors.xml
│                   └── styles.xml
├── build.gradle
├── gradle
│   ├── libs.versions.toml
│   └── wrapper
│       ├── gradle-wrapper.jar
│       └── gradle-wrapper.properties
├── gradle.properties
├── gradlew
├── gradlew.bat
└── settings.gradle
~~~

Una vez visto esto, vamos rellenando/modificando los siguientes ficheros:

#### 📝 gradle.properties
Este fichero lo dejamos de esta manera:
~~~properties
# This file was generated by the Gradle 'init' task.
# https://docs.gradle.org/current/userguide/build_environment.html#sec:gradle_configuration_properties

org.gradle.configuration-cache=true

# Para el uso moderno de dependencias de Android (AndroidX)
android.enableJetifier=true
android.useAndroidX=true
~~~

#### 📝 build.gradle
~~~gradle
/*
 * This file was generated by the Gradle 'init' task.
 *
 * This is a general purpose Gradle build.
 * Learn more about Gradle by exploring our Samples at https://docs.gradle.org/8.14.2/samples
 */
buildscript {
  ext {
    agp_version = '8.10.1' // Versión del Android Gradle Plugin
  }
  repositories {
    google()
    mavenCentral()
  }
  dependencies {
    classpath "com.android.tools.build:gradle:$agp_version"
  }
}

allprojects {
  repositories {
    google()
    mavenCentral()
  }
}
~~~

#### 📝 app/build.gradle
~~~gradle
/*
 * This file was generated by the Gradle 'init' task.
 *
 * This generated file contains a sample Java application project to get you started.
 * For more details on building Java & JVM projects, please refer to https://docs.gradle.org/8.14.2/userguide/building_java_projects.html in the Gradle documentation.
 */

plugins {
    // Apply the application plugin to add support for building a CLI application in Java.
    id 'com.android.application'
}

android {
  namespace 'com.programatta.games.demoandroid'
  compileSdk 35
  defaultConfig {
    applicationId "com.programatta.games.demoandroid"
    minSdk 24
    targetSdk 35
    versionCode 1
    versionName "1.0"
  }

  buildTypes {
    release {
      minifyEnabled false
      proguardFiles getDefaultProguardFile('proguard-android-optimize.txt'), 'proguard-rules.pro'
    }
  }
}

dependencies {
  // Dependencias estándar de Android
  implementation 'androidx.appcompat:appcompat:1.7.0'
 
  // ¡Aquí es donde añades tu AAR!
  implementation files('libs/game.aar')
 
  // This line is needed to resolve a mysterious compilation error.
  // https://stackoverflow.com/questions/75263047/duplicate-class-in-kotlin-android
  implementation platform("org.jetbrains.kotlin:kotlin-bom:1.8.0")
}

// Apply a specific Java toolchain to ease working on different environments.
java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(17)
    }
}
~~~

#### 📝 app/src/main/AndroidManifest.xml
~~~xml
<manifest xmlns:android="http://schemas.android.com/apk/res/android">
  <uses-feature android:glEsVersion="0x00020000" android:required="true" />
  <application
    android:supportsRtl="true"
    android:allowBackup="true"
    android:label="Demo Android"
    android:theme="@style/AppTheme">
    <activity 
      android:exported="true"
      android:name=".MainActivity"
      android:label="Demo Android"
      android:screenOrientation="portrait"
      android:launchMode="singleInstance">
      <intent-filter>
        <action android:name="android.intent.action.MAIN" />
        <category android:name="android.intent.category.LAUNCHER" />
      </intent-filter>
    </activity>
  </application>
</manifest>
~~~

#### 📝 app/src/main/java/com/programatta/games/demoandroid/MainActivity.java
~~~java
/*
 * This source file was generated by the Gradle 'init' task
 */
package com.programatta.games.demoandroid;

import android.os.Bundle;
import android.util.Log;
import androidx.appcompat.app.AppCompatActivity;
import go.Seq;
import com.programatta.games.demoandroid.corelib.mobile.EbitenView;


public class MainActivity extends AppCompatActivity {
  @Override
  protected void onCreate(Bundle savedInstanceState) {
    super.onCreate(savedInstanceState);
    setContentView(R.layout.activity_main);
    Seq.setContext(getApplicationContext());
  }

  @Override
  protected void onPause() {
    super.onPause();
    this.getEbitenView().suspendGame();
  }

  @Override
  protected void onResume() {
    super.onResume();
    this.getEbitenView().resumeGame();
  }

  private EbitenView getEbitenView() {
    return (EbitenView)this.findViewById(R.id.ebitenview);
  }
}
~~~

#### 📝 app/src/main/java/com/programatta/games/demoandroid/EbitenViewWithErrorHandling.java
~~~java
package com.programatta.games.demoandroid;

import android.content.Context;
import android.util.AttributeSet;
import com.programatta.games.demoandroid.corelib.mobile.EbitenView;


class EbitenViewWithErrorHandling extends EbitenView {
  public EbitenViewWithErrorHandling(Context context) {
    super(context);
  }

  public EbitenViewWithErrorHandling(Context context, AttributeSet attributeSet) {
    super(context, attributeSet);
  }

  @Override
  protected void onErrorOnGameUpdate(Exception e) {
    // You can define your own error handling e.g., using Crashlytics.
    // e.g., Crashlytics.logException(e);
    super.onErrorOnGameUpdate(e);
  }
}
~~~


#### 📝 app/src/main/layout/activity_main.xml
~~~xml
<?xml version="1.0" encoding="utf-8"?>
<RelativeLayout xmlns:android="http://schemas.android.com/apk/res/android"
  xmlns:tools="http://schemas.android.com/tools"
  android:layout_width="match_parent"
  android:layout_height="match_parent"
  android:keepScreenOn="true"
  tools:context=".MainActivity">
  <com.programatta.games.demoandroid.EbitenViewWithErrorHandling
    android:id="@+id/ebitenview"
    android:layout_width="match_parent"
    android:layout_height="match_parent"
    android:focusable="true" />
</RelativeLayout>
~~~

#### 📝 app/src/main/values/colors.xml
~~~xml
<?xml version="1.0" encoding="utf-8"?>
<resources>
  <color name="colorPrimary">#3F51B5</color>
  <color name="colorPrimaryDark">#303F9F</color>
  <color name="colorAccent">#FF4081</color>
</resources>
~~~

#### 📝 app/src/main/values/styles.xml
~~~xml
<resources>
  <!-- Base application theme. -->
  <style name="AppTheme" parent="Theme.AppCompat.Light.DarkActionBar">
    <!-- Customize your theme here. -->
    <item name="colorPrimary">@color/colorPrimary</item>
    <item name="colorPrimaryDark">@color/colorPrimaryDark</item>
    <item name="colorAccent">@color/colorAccent</item>
    <item name="windowNoTitle">true</item>
    <item name="android:windowFullscreen">true</item>
    <item name="android:windowContentOverlay">@null</item>
  </style>
</resources>
~~~

Una vez finalizado todos estos añadidos y cambios, pedimos a **cradle** que complete la configuración que necesita ejecutando el script **gradlew**:
~~~shell
./gradlew tasks
~~~

Una vez que finalice el proceso, si no ha habido ningún error, ya podemos crear el **apk** tanto de debug, como de release (aquí siempre y cuando tengamos una llave).
Para debug:
~~~shell
./gradlew assembleDebug
~~~

Para Release:
~~~shell
./gradlew assembleRelease
~~~

El **apk** generado lo encontramos en **app/build/outputs/apk/debug** o **app/build/outputs/apk/release**.

