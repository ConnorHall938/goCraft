package atlas

const (
	AtlasCols = 4 // 4 textures per row
	AtlasRows = 2 // 2 textures per column
)

// Map texture names to atlas index
var TextureIndex = map[string]int{
	"grass_top":  0,
	"grass_side": 1,
	"dirt":       2,
	"stone":      3,
	"wood_side":  4,
	"wood_end":   5,
}

// Returns the UV rectangle (u0,v0,u1,v1) for a texture index
func UVRect(index int) (float32, float32, float32, float32) {
	col := index % AtlasCols
	row := index / AtlasCols

	u0 := float32(col) / float32(AtlasCols)
	u1 := u0 + 1.0/float32(AtlasCols)

	v0 := float32(row) / float32(AtlasRows)
	v1 := v0 + 1.0/float32(AtlasRows)
	return u0, v1, u1, v0
}
