package block

import "goCraft/lib/atlas"

var Grass = BlockType{
	Name: "grass",

	FaceTextures: [6]int{
		atlas.GRASS_SIDE, // front
		atlas.GRASS_SIDE, // back
		atlas.GRASS_TOP,  // top
		atlas.DIRT,       // bottom
		atlas.GRASS_SIDE, // right
		atlas.GRASS_SIDE, // left
	},

	FaceTints: [6][3]float32{
		{1.0, 1.0, 1.0}, // front
		{1.0, 1.0, 1.0},
		{0.8, 1.0, 0.8}, // top
		{1.0, 1.0, 1.0}, // bottom
		{1.0, 1.0, 1.0},
		{1.0, 1.0, 1.0},
	},
}
