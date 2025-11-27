// block/blocktype.go
package block

type BlockType struct {
	Name         string
	FaceTextures [6]int        // texture index for each face
	FaceTints    [6][3]float32 // RGB tint per face
	BasePosition [3]float32    // Position in the chunk, default (0,0,0), up to 32, 32, 32
}
