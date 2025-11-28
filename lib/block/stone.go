package block

import "goCraft/lib/atlas"

var Stone = BlockType{
	Name: "stone",

	FaceTextures: [6]uint32{
		atlas.STONE, // front
		atlas.STONE, // back
		atlas.STONE, // top
		atlas.STONE, // bottom
		atlas.STONE, // right
		atlas.STONE, // left
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
