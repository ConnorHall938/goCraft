package chunk

import (
	"goCraft/lib/block"
)

type ChunkMesh struct {
	Vertices []float32
	Indices  []uint32
}

func (c *Chunk) BuildMesh() ChunkMesh {
	mesh := ChunkMesh{}

	for x := 0; x < 32; x++ {
		for y := 0; y < 128; y++ {
			for z := 0; z < 32; z++ {
				id := c.Blocks[x][y][z]
				if id == 0 {
					continue // air
				}

				bt := block.Registry[id]
				bt.BasePosition = [3]float32{float32(x), float32(y), float32(z)}

				cube := block.BuildCubeMesh(bt)

				// Append vertices & fix index offset
				startIndex := uint32(len(mesh.Vertices) / 8)

				mesh.Vertices = append(mesh.Vertices, cube.Vertices...)

				for _, idx := range cube.Indices {
					mesh.Indices = append(mesh.Indices, startIndex+idx)
				}
			}
		}
	}

	return mesh
}
