package canvas

import (
	_ "embed"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed "Arial Unicode MS.ttf"
var ArialUnicodeMSRaw []byte

var ArialUnicodeMS *truetype.Font

func init() {
	if f, err := truetype.Parse(ArialUnicodeMSRaw); err != nil {
		panic(err)
	} else {
		ArialUnicodeMS = f
	}
}

func GetFace(size float64) font.Face {
	return truetype.NewFace(ArialUnicodeMS, &truetype.Options{
		Size: size,
	})
}
