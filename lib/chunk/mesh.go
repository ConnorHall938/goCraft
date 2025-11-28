package chunk

import (
	"goCraft/lib/atlas"
	"goCraft/lib/block"
)

type ChunkMesh struct {
	Vertices []float32
	Indices  []uint32
}

func (c *Chunk) BuildMesh(atlas atlas.Atlas) ChunkMesh {
	mesh := ChunkMesh{}
	verts := []float32{}
	inds := []uint32{}
	indexOffset := uint32(0)

	for x := 0; x < 32; x++ {
		for y := 0; y < 128; y++ {
			for z := 0; z < 32; z++ {

				blockID := c.Blocks[x][y][z]
				if blockID == 0 {
					continue // air
				}

				bt := block.Registry[blockID]

				// Check each of 6 faces
				neighbors := [6]bool{
					c.isAir(x, y, z+1), // front
					c.isAir(x, y, z-1), // back
					c.isAir(x, y+1, z), // top
					c.isAir(x, y-1, z), // bottom
					c.isAir(x+1, y, z), // right
					c.isAir(x-1, y, z), // left
				}

				for faceIdx := 0; faceIdx < 6; faceIdx++ {
					if !neighbors[faceIdx] {
						continue
					}

					fv, fi := block.MakeFaceMesh(atlas, bt, faceIdx, x, y, z)

					// append vertices
					verts = append(verts, fv...)

					// remap indices
					for _, idx := range fi {
						inds = append(inds, indexOffset+idx)
					}

					indexOffset += 4 // 4 vertices per face
				}
			}
		}
	}

	mesh.Vertices = verts
	mesh.Indices = inds
	return mesh
}
