package chunk

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

func (c *Chunk) Render() {
	// Placeholder for rendering logic
}
