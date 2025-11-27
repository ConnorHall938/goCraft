package block

const (
	BlockAir   = 0
	BlockGrass = 1
	BlockDirt  = 2
	BlockStone = 3
)

var Registry = map[uint8]BlockType{
	BlockGrass: Grass,
	BlockDirt:  Dirt,
	BlockStone: Stone,
}
