// cube/mesh.go
package block

import "goCraft/lib/atlas"

type CubeMesh struct {
	Vertices   []float32
	Indices    []uint32
	FaceRanges [6][2]uint32 // {start, count}
}

func BuildCubeMesh(atlas atlas.Atlas, bt BlockType) CubeMesh {
	mesh := CubeMesh{}
	vertices := []float32{}
	indices := []uint32{}

	for face := 0; face < 6; face++ {
		tex := bt.FaceTextures[face]
		tint := bt.FaceTints[face]

		// UV coordinates
		u0, v0, u1, v1 := atlas.UVRect(tex)

		// 4 vertices
		f := cubeFaces[face]
		baseIndex := uint32(len(vertices) / 8)

		vertices = append(vertices,
			f[0][0]+bt.BasePosition[0], f[0][1]+bt.BasePosition[1], f[0][2]+bt.BasePosition[2], u0, v0, tint[0], tint[1], tint[2],
			f[1][0]+bt.BasePosition[0], f[1][1]+bt.BasePosition[1], f[1][2]+bt.BasePosition[2], u1, v0, tint[0], tint[1], tint[2],
			f[2][0]+bt.BasePosition[0], f[2][1]+bt.BasePosition[1], f[2][2]+bt.BasePosition[2], u1, v1, tint[0], tint[1], tint[2],
			f[3][0]+bt.BasePosition[0], f[3][1]+bt.BasePosition[1], f[3][2]+bt.BasePosition[2], u0, v1, tint[0], tint[1], tint[2],
		)

		// 2 triangles
		ranges := [2]uint32{uint32(len(indices)), 6}
		mesh.FaceRanges[face] = ranges

		indices = append(indices,
			baseIndex, baseIndex+1, baseIndex+2,
			baseIndex+2, baseIndex+3, baseIndex,
		)
	}

	mesh.Vertices = vertices
	mesh.Indices = indices
	return mesh
}

// MakeFaceMesh builds a *single face* of a cube.
// faceIndex = 0..5 (front, back, top, bottom, right, left)
// world position = (bx, by, bz)
func MakeFaceMesh(atlas atlas.Atlas, bt BlockType, faceIndex int, bx, by, bz int) ([]float32, []uint32) {
	tex := bt.FaceTextures[faceIndex]
	tint := bt.FaceTints[faceIndex]
	u0, v0, u1, v1 := atlas.UVRect(tex)

	f := cubeFaces[faceIndex]

	// Convert to world position
	offset := [3]float32{float32(bx), float32(by), float32(bz)}

	// 4 vertices (pos + uv + tint)
	verts := []float32{
		f[0][0] + offset[0], f[0][1] + offset[1], f[0][2] + offset[2], u0, v0, tint[0], tint[1], tint[2],
		f[1][0] + offset[0], f[1][1] + offset[1], f[1][2] + offset[2], u1, v0, tint[0], tint[1], tint[2],
		f[2][0] + offset[0], f[2][1] + offset[1], f[2][2] + offset[2], u1, v1, tint[0], tint[1], tint[2],
		f[3][0] + offset[0], f[3][1] + offset[1], f[3][2] + offset[2], u0, v1, tint[0], tint[1], tint[2],
	}

	inds := []uint32{0, 1, 2, 2, 3, 0}
	return verts, inds
}
