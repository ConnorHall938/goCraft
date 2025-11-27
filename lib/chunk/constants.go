package chunk

// Chunk dimensions
const (
	ChunkWidth  = 32         // X dimension
	ChunkHeight = 128        // Y dimension
	ChunkDepth  = ChunkWidth // Z dimension
)

// World generation constants
const (
	DefaultSurfaceLevel = 64 // Default surface height for flat terrain
	DirtLayerDepth      = 4  // Depth of dirt layer below grass
)
