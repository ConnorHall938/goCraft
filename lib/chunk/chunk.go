package chunk

import "goCraft/lib/block"

type Chunk struct {
	Blocks   [ChunkWidth][ChunkHeight][ChunkDepth]uint8 // 3D array of block IDs
	Location [3]int                                     // Chunk location in the world
}

func NewChunk() *Chunk {
	return &Chunk{}
}

func (c *Chunk) SetBlock(x, y, z int, blockID uint8) {
	if x < 0 || x >= ChunkWidth || y < 0 || y >= ChunkHeight || z < 0 || z >= ChunkDepth {
		return
	}
	c.Blocks[x][y][z] = blockID
}

func (c *Chunk) GetBlock(x, y, z int) uint8 {
	if x < 0 || x >= ChunkWidth || y < 0 || y >= 128 || z < 0 || z >= ChunkWidth {
		return 0 // Air
	}
	return c.Blocks[x][y][z]
}

func (c *Chunk) Fill(blockID uint8) {
	for x := 0; x < ChunkWidth; x++ {
		for y := 0; y < ChunkHeight; y++ {
			for z := 0; z < ChunkWidth; z++ {
				c.Blocks[x][y][z] = blockID
			}
		}
	}
}

func (c *Chunk) GenerateFlat() {
	for x := 0; x < ChunkWidth; x++ {
		for z := 0; z < ChunkWidth; z++ {
			for y := 0; y <= DefaultSurfaceLevel; y++ {

				switch {
				case y == DefaultSurfaceLevel:
					c.Blocks[x][y][z] = block.BlockGrass

				case y <= DefaultSurfaceLevel && y >= DefaultSurfaceLevel-DirtLayerDepth:
					c.Blocks[x][y][z] = block.BlockDirt

				default:
					c.Blocks[x][y][z] = block.BlockStone
				}
			}
		}
	}
}

func (c *Chunk) isAir(x, y, z int) bool {
	if x < 0 || x >= ChunkWidth || y < 0 || y >= ChunkHeight || z < 0 || z >= ChunkDepth {
		return true // Out of bounds is considered air
	}
	return c.Blocks[x][y][z] == block.BlockAir
}
