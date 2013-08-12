package mcmap

const (
	ChunkSizeXZ = 16
	ChunkSizeY  = 256
	ChunkRectXZ = ChunkSizeXZ * ChunkSizeXZ
	ChunkSize   = ChunkRectXZ * ChunkSizeY
)

const superchunkSizeXZ = 32
const chunkSectionSize = ChunkRectXZ * 16
