package cube

import "goCraft/lib/atlas"

// A face is 4 vertices in counter-clockwise order (looking at the face)
type Face [4][3]float32

// Cube faces in object space, CCW winding
var cubeFaces = [6]Face{
	// FRONT (+Z)
	{
		{-0.5, -0.5, 0.5}, // 0 bottom-left
		{0.5, -0.5, 0.5},  // 1 bottom-right
		{0.5, 0.5, 0.5},   // 2 top-right
		{-0.5, 0.5, 0.5},  // 3 top-left
	},
	// BACK (-Z)
	{
		{0.5, -0.5, -0.5},  // 0 bottom-left
		{-0.5, -0.5, -0.5}, // 1 bottom-right
		{-0.5, 0.5, -0.5},  // 2 top-right
		{0.5, 0.5, -0.5},   // 3 top-left
	},
	// TOP (+Y)
	{
		{-0.5, 0.5, 0.5},  // 0 bottom-left (towards +Z)
		{0.5, 0.5, 0.5},   // 1 bottom-right
		{0.5, 0.5, -0.5},  // 2 top-right (towards -Z)
		{-0.5, 0.5, -0.5}, // 3 top-left
	},
	// BOTTOM (-Y)
	{
		{-0.5, -0.5, -0.5}, // 0 bottom-left (towards -Z)
		{0.5, -0.5, -0.5},  // 1 bottom-right
		{0.5, -0.5, 0.5},   // 2 top-right (towards +Z)
		{-0.5, -0.5, 0.5},  // 3 top-left
	},
	// RIGHT (+X)
	{
		{0.5, -0.5, 0.5},  // 0 bottom-left (towards +Z)
		{0.5, -0.5, -0.5}, // 1 bottom-right
		{0.5, 0.5, -0.5},  // 2 top-right
		{0.5, 0.5, 0.5},   // 3 top-left
	},
	// LEFT (-X)
	{
		{-0.5, -0.5, -0.5}, // 0 bottom-left (towards -Z)
		{-0.5, -0.5, 0.5},  // 1 bottom-right
		{-0.5, 0.5, 0.5},   // 2 top-right
		{-0.5, 0.5, -0.5},  // 3 top-left
	},
}

// Append one face (4 vertices) with correct UVs for a tile from the atlas.
//
// IMPORTANT ASSUMPTION:
// atlas.UVRect(textureID) returns (u0,v0,u1,v1) where:
//   u0,v0 = bottom-left of tile
//   u1,v1 = top-right  of tile
func addFace(vertices *[]float32, face Face, textureID int) {
	u0, v0, u1, v1 := atlas.UVRect(textureID)

	// bottom-left, bottom-right, top-right, top-left
	*vertices = append(*vertices,
		face[0][0], face[0][1], face[0][2], u0, v0, // bottom-left
		face[1][0], face[1][1], face[1][2], u1, v0, // bottom-right
		face[2][0], face[2][1], face[2][2], u1, v1, // top-right
		face[3][0], face[3][1], face[3][2], u0, v1, // top-left
	)
}

// MakeCube returns vertices + indices for a cube where:
//   topID    = texture for top face
//   bottomID = texture for bottom face
//   sideID   = texture for all side faces
func MakeCube(topID, bottomID, sideID int) ([]float32, []uint32) {
	vertices := []float32{}
	indices := []uint32{}

	// faces: front, back, top, bottom, right, left
	faceTex := [6]int{
		sideID,   // front
		sideID,   // back
		topID,    // top
		bottomID, // bottom
		sideID,   // right
		sideID,   // left
	}

	for i := 0; i < 6; i++ {
		addFace(&vertices, cubeFaces[i], faceTex[i])

		offset := uint32(i * 4)
		indices = append(indices,
			offset+0, offset+1, offset+2,
			offset+2, offset+3, offset+0,
		)
	}

	return vertices, indices
}
