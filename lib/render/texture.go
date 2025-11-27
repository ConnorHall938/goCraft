package render

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

func LoadAtlas(paths []string) uint32 {
	rgba := image.NewRGBA(image.Rect(0, 0, 256*4, 256*2))

	for i, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Decoded %s\n", path)

		img, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}
		x := (i % 4) * 256
		y := (i / 4) * 256
		draw.Draw(rgba, image.Rect(x, y, x+256, y+256), img, image.Point{}, draw.Src)
		file.Close()
	}

	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	size := rgba.Rect.Size()

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA8,
		int32(size.X),
		int32(size.Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)

	// 🔴 TEMP: no mipmaps
	// gl.GenerateMipmap(gl.TEXTURE_2D)

	// 🔴 TEMP: nearest sampling only
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// wrapping doesn’t really matter for an atlas, but this is fine:
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	return tex
}
