package render

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

const (
	TEXTURE_MAX_ANISOTROPY_EXT     = 0x84FE
	MAX_TEXTURE_MAX_ANISOTROPY_EXT = 0x84FF
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

		// copy texture into its cell
		x := (i % 4) * 256
		y := (i / 4) * 256
		draw.Draw(rgba, image.Rect(x, y, x+256, y+256), img, image.Point{}, draw.Src)

		file.Close()
	}

	// --- Upload to OpenGL ---
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

	// ----------------------------------------
	// ✅ MIPMAPS ENABLED
	// ----------------------------------------
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// MINIFICATION: use trilinear filtering (best quality)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)

	// MAGNIFICATION: linear smoothing
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// ----------------------------------------
	// ✅ CLAMP to prevent bleeding between tiles
	// ----------------------------------------
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// ----------------------------------------
	// ✅ ANISOTROPIC FILTERING
	// ----------------------------------------
	var maxAniso float32
	gl.GetFloatv(MAX_TEXTURE_MAX_ANISOTROPY_EXT, &maxAniso)

	if maxAniso > 0 {
		gl.TexParameterf(gl.TEXTURE_2D, TEXTURE_MAX_ANISOTROPY_EXT, maxAniso)
		fmt.Printf("Anisotropic filtering enabled: %.1fx\n", maxAniso)
	} else {
		fmt.Println("Anisotropic filtering not supported.")
	}

	return tex
}
