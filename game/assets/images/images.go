package images

import (
	"embed"
)

//go:embed emojis/*.png
var EmojisDataFS embed.FS
