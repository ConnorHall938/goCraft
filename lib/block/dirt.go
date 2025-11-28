package block

import "goCraft/lib/atlas"

var Dirt = BlockType{
	Name: "dirt",

	FaceTextures: [6]uint32{
		atlas.DIRT, // front
		atlas.DIRT, // back
		atlas.DIRT, // top
		atlas.DIRT, // bottom
		atlas.DIRT, // right
		atlas.DIRT, // left
	},

	FaceTints: [6][3]float32{
		{1.0, 1.0, 1.0}, // front
		{1.0, 1.0, 1.0},
		{1.0, 1.0, 1.0}, // top
		{1.0, 1.0, 1.0}, // bottom
		{1.0, 1.0, 1.0},
		{1.0, 1.0, 1.0},
	},
}
