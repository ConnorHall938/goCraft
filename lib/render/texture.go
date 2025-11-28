package render

import (
	"fmt"
	"goCraft/lib/atlas"
	"image"
	"image/draw"
	_ "image/png"
	"math"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

const (
	TEXTURE_MAX_ANISOTROPY_EXT     = 0x84FE
	MAX_TEXTURE_MAX_ANISOTROPY_EXT = 0x84FF
)

func LoadAtlas(paths []string) (atlas.Atlas, error) {
	// On the first image, find the size and create the atlas RGBA
	file, err := os.Open(paths[0])
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decoded %s\n", paths[0])

	image_width := img.Bounds().Dx()
	image_height := img.Bounds().Dy()
	// Calculate optimal square-ish atlas dimensions
	// e.g., 6 textures -> 3x2, 10 textures -> 4x3
	atlas_width := int(math.Ceil(math.Sqrt(float64(len(paths)))))
	atlas_height := int(math.Ceil(float64(len(paths)) / float64(atlas_width)))

	rgba := image.NewRGBA(image.Rect(0, 0, atlas_width*image_width, atlas_height*image_height))
	// copy first texture into its cell
	draw.Draw(rgba, image.Rect(0, 0, image_width, image_height), img, image.Point{}, draw.Src)

	file.Close()
	// Now load each image and copy it into the atlas RGBA

	// Iterate over the other images, skipping the first.
	for i := 1; i < len(paths); i++ {
		file, err := os.Open(paths[i])
		if err != nil {
			panic(err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Decoded %s\n", paths[i])

		// copy texture into its cell
		x := (i % atlas_width) * image_width
		y := (i / atlas_width) * image_height
		draw.Draw(rgba, image.Rect(x, y, x+image_width, y+image_height), img, image.Point{}, draw.Src)

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

	return atlas.Atlas{
		AtlasImageId: tex,
		Rows:         uint32(atlas_height),
		Columns:      uint32(atlas_width),
		ImageWidth:   uint32(image_width),
		ImageHeight:  uint32(image_height),
	}, nil
}
