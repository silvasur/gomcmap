package mcmap

import (
	"errors"
	"github.com/kch42/gonbt/nbt"
	"time"
)

func calcBlockOffset(x, y, z int) int {
	if (x < 0) || (y < 0) || (z < 0) || (x >= 16) || (y >= 256) || (z >= 16) {
		return -1
	}

	return x | (z << 4) | (y << 8)
}

func offsetToPos(off int) (x, y, z int) {
	x = off & 0xf
	z = (off >> 4) & 0xf
	y = (off >> 8) & 0xff
	return
}

// BlockToChunk calculates the chunk (cx, cz) and the block position in this chunk(rbx, rbz) of a block position given global coordinates.
func BlockToChunk(bx, bz int) (cx, cz, rbx, rbz int) {
	cx = bx << 4
	cz = bz << 4
	rbx = ((cx % 16) + 16) % 16
	rbz = ((cz % 16) + 16) % 16
	return
}

// ChunkToBlock calculates the global position of a block, given the chunk position (cx, cz) and the plock position in that chunk (rbx, rbz).
func ChunkToBlock(cx, cz, rbx, rbz int) (bx, bz int) {
	bx = cx*16 + rbx
	bz = cz*16 + rbz
	return
}

// Chunk represents a 16*16*256 Chunk of the region.
type Chunk struct {
	Entities []nbt.TagCompound

	x, z int32

	lastUpdate    int64
	populated     bool
	inhabitedTime int64
	ts            time.Time

	heightMap []int32 // Ordered ZX

	modified bool
	blocks   []Block // Ordered YZX
	biomes   []Biome // Ordered XZ
}

// MarkModified needs to be called, if some data of the chunk was modified.
func (c *Chunk) MarkModified() { c.modified = true }

// Coords returns the Chunk's coordinates.
func (c *Chunk) Coords() (X, Z int32) { return c.x, c.z }

// Block gives you a reference to the Block located at x, y, z. If you modify the block data, you need to call the MarkModified() function of the chunk.
//
// x and z must be in [0, 15], y in [0, 255]. Otherwise a nil pointer is returned.
func (c *Chunk) Block(x, y, z int) *Block {
	off := calcBlockOffset(x, y, z)
	if off < 0 {
		return nil
	}

	return &(c.blocks[off])
}

// Height returns the height at x, z.
//
// x and z must be in [0, 15]. Height will panic, if this is violated!
func (c *Chunk) Height(x, z int) int {
	if (x < 0) || (x > 15) || (z < 0) || (z > 15) {
		panic(errors.New("x or z parameter was out of range"))
	}

	return int(c.heightMap[z*16+x])
}

// Biome gets the Biome at x,z.
func (c *Chunk) Biome(x, z int) Biome { return c.biomes[x*16+z] }

// SetBiome sets the biome at x,z.
func (c *Chunk) SetBiome(x, z int, bio Biome) { c.biomes[x*16+z] = bio }

// TODO: func (c *Chunk) RecalcHeightMap()
