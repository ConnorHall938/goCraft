package atlas

type Atlas struct {
	AtlasImageId uint32
	Columns      uint32
	Rows         uint32
	ImageWidth   uint32
	ImageHeight  uint32
}

const (
	padU = 0 //0.5 / float32(AtlasCols*256) // half pixel padding
	padV = 0 //0.5 / float32(AtlasRows*256)
)

// Returns the UV rectangle (u0,v0,u1,v1) for a texture index
func (atlas *Atlas) UVRect(index uint32) (float32, float32, float32, float32) {
	col := index % atlas.Columns
	row := index / atlas.Columns

	u0 := float32(col)/float32(atlas.Columns) + padU
	u1 := u0 + 1.0/float32(atlas.Columns) - padU

	v0 := float32(row)/float32(atlas.Rows) + padV
	v1 := v0 + 1.0/float32(atlas.Rows) - padV
	return u0, v1, u1, v0
}
