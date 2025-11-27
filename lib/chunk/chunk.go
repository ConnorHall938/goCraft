package chunk

import "goCraft/lib/block"

type Chunk struct {
	Blocks   [32][128][32]uint8 // 3D array of block IDs
	Location [3]int             // Chunk location in the world
}

func NewChunk() *Chunk {
	return &Chunk{}
}

func (c *Chunk) SetBlock(x, y, z int, blockID uint8) {
	if x < 0 || x >= 32 || y < 0 || y >= 128 || z < 0 || z >= 32 {
		return
	}
	c.Blocks[x][y][z] = blockID
}

func (c *Chunk) GetBlock(x, y, z int) uint8 {
	if x < 0 || x >= 32 || y < 0 || y >= 128 || z < 0 || z >= 32 {
		return 0 // Air
	}
	return c.Blocks[x][y][z]
}

func (c *Chunk) Fill(blockID uint8) {
	for x := 0; x < 32; x++ {
		for y := 0; y < 128; y++ {
			for z := 0; z < 32; z++ {
				c.Blocks[x][y][z] = blockID
			}
		}
	}
}

func (c *Chunk) GenerateFlat() {
	surface := 64

	for x := 0; x < 32; x++ {
		for z := 0; z < 32; z++ {
			for y := 0; y <= surface; y++ {

				switch {
				case y == surface:
					c.Blocks[x][y][z] = block.BlockGrass

				case y <= surface && y >= surface-5:
					c.Blocks[x][y][z] = block.BlockDirt

				default:
					c.Blocks[x][y][z] = block.BlockStone
				}
			}
		}
	}
}

func (c *Chunk) isAir(x, y, z int) bool {
	if x < 0 || x >= 32 || y < 0 || y >= 128 || z < 0 || z >= 32 {
		return true // Out of bounds is considered air
	}
	return c.Blocks[x][y][z] == block.BlockAir
}
