package mcmap

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/kch42/gonbt/nbt"
	"time"
)

const (
	_ = iota
	compressGZip
	compressZlib
)

type preChunk struct {
	ts          time.Time
	data        []byte
	compression byte
}

var (
	UnknownCompression = errors.New("Unknown chunk compression")
)

func halfbyte(b []byte, i int) byte {
	if i%2 == 0 {
		return b[i/2] & 0x0f
	}
	return (b[i/2] >> 4) & 0x0f
}

func setHalfbyte(b []byte, i int, v byte) {
	v &= 0xf
	if i%2 == 0 {
		b[i/2] |= v
	} else {
		b[i/2] |= v << 4
	}
}

func extractCoord(tc nbt.TagCompound) (x, y, z int, err error) {
	var _x, _y, _z int32
	if _x, err = tc.GetInt("x"); err != nil {
		return
	}
	if _y, err = tc.GetInt("y"); err != nil {
		return
	}
	_z, err = tc.GetInt("z")
	x, y, z = int(_x), int(_y), int(_z)
	return
}

func (pc *preChunk) getLevelTag() (nbt.TagCompound, error) {
	r := bytes.NewReader(pc.data)

	var root nbt.Tag
	var err error
	switch pc.compression {
	case compressGZip:
		root, _, err = nbt.ReadGzipdNamedTag(r)
	case compressZlib:
		root, _, err = nbt.ReadZlibdNamedTag(r)
	default:
		err = UnknownCompression
	}

	if err != nil {
		return nil, err
	}

	if root.Type != nbt.TAG_Compound {
		return nil, errors.New("Root tag is not a TAG_Compound")
	}

	lvl, err := root.Payload.(nbt.TagCompound).GetCompound("Level")
	if err != nil {
		return nil, fmt.Errorf("Could not read Level tag: %s", err)
	}

	return lvl, nil
}

func (pc *preChunk) toChunk(reg *Region) (*Chunk, error) {
	c := Chunk{ts: pc.ts, reg: reg}

	lvl, err := pc.getLevelTag()
	if err != nil {
		return nil, err
	}

	c.x, err = lvl.GetInt("xPos")
	if err != nil {
		return nil, fmt.Errorf("Could not read xPos tag: %s", err)
	}
	c.z, err = lvl.GetInt("zPos")
	if err != nil {
		return nil, fmt.Errorf("Could not read zPos tag: %s", err)
	}

	c.lastUpdate, err = lvl.GetLong("LastUpdate")
	if err != nil {
		return nil, fmt.Errorf("Could not read LastUpdate tag: %s", err)
	}

	populated, err := lvl.GetByte("TerrainPopulated")
	switch err {
	case nil:
	case nbt.NotFound:
		populated = 1
	default:
		return nil, fmt.Errorf("Could not read TerrainPopulated tag: %s", err)
	}
	c.populated = (populated == 1)

	c.inhabitedTime, err = lvl.GetLong("InhabitedTime")
	switch err {
	case nil:
	case nbt.NotFound:
		c.inhabitedTime = 0
	default:
		return nil, fmt.Errorf("Could not read InhabitatedTime tag: %s", err)
	}

	c.biomes = make([]Biome, 256)
	biomes, err := lvl.GetByteArray("Biomes")
	switch err {
	case nil:
		for i, bio := range biomes {
			c.biomes[i] = Biome(bio)
		}
	case nbt.NotFound:
		for i := 0; i < 256; i++ {
			c.biomes[i] = BioUncalculated
		}
	default:
		return nil, fmt.Errorf("Could not read Biomes tag: %s", err)
	}

	c.heightMap, err = lvl.GetIntArray("HeightMap")
	if err != nil {
		return nil, fmt.Errorf("Could not read HeightMap tag: %s", err)
	}

	ents, err := lvl.GetList("Entities")
	if err != nil {
		return nil, fmt.Errorf("Could not read Entities tag: %s", err)
	}
	if ents.Type != nbt.TAG_Compound {
		c.Entities = []nbt.TagCompound{}
	} else {
		c.Entities = make([]nbt.TagCompound, len(ents.Elems))
		for i, ent := range ents.Elems {
			c.Entities[i] = ent.(nbt.TagCompound)
		}
	}

	sections, err := lvl.GetList("Sections")
	if (err != nil) || (sections.Type != nbt.TAG_Compound) {
		return nil, fmt.Errorf("Could not read Section tag: %s", err)
	}

	c.blocks = make([]Block, 16*16*256)
	for _, _section := range sections.Elems {
		section := _section.(nbt.TagCompound)

		y, err := section.GetByte("Y")
		if err != nil {
			return nil, fmt.Errorf("Could not read Section -> Y tag: %s", err)
		}
		off := int(y) * 4096

		blocks, err := section.GetByteArray("Blocks")
		if err != nil {
			return nil, fmt.Errorf("Could not read Section -> Blocks tag: %s", err)
		}
		blocksAdd := make([]byte, 4096)
		add, err := section.GetByteArray("Add")
		switch err {
		case nil:
			for i := 0; i < 4096; i++ {
				blocksAdd[i] = halfbyte(add, i)
			}
		case nbt.NotFound:
		default:
			return nil, fmt.Errorf("Could not read Section -> Add tag: %s", err)
		}

		blkData, err := section.GetByteArray("Data")
		if err != nil {
			return nil, fmt.Errorf("Could not read Section -> Data tag: %s", err)
		}
		blockLight, err := section.GetByteArray("BlockLight")
		if err != nil {
			return nil, fmt.Errorf("Could not read Section -> BlockLight tag: %s", err)
		}
		skyLight, err := section.GetByteArray("SkyLight")
		if err != nil {
			return nil, fmt.Errorf("Could not read Section -> SkyLight tag: %s", err)
		}

		for i := 0; i < 4096; i++ {
			c.blocks[off+i] = Block{
				ID:         BlockID(uint16(blocks[i]) | (uint16(blocksAdd[i]) << 8)),
				Data:       halfbyte(blkData, i),
				BlockLight: halfbyte(blockLight, i),
				SkyLight:   halfbyte(skyLight, i)}
		}
	}

	tileEnts, err := lvl.GetList("TileEntities")
	if err != nil {
		return nil, fmt.Errorf("Could not read TileEntities tag: %s", err)
	}
	if tileEnts.Type == nbt.TAG_Compound {
		for _, _tEnt := range tileEnts.Elems {
			tEnt := _tEnt.(nbt.TagCompound)
			x, y, z, err := extractCoord(tEnt)
			if err != nil {
				return nil, fmt.Errorf("Could not Extract coords: %s", err)
			}

			_, _, x, z = BlockToChunk(x, z)

			c.blocks[calcBlockOffset(x, y, z)].TileEntity = tEnt
		}
	}

	tileTicks, err := lvl.GetList("TileTicks")
	if (err == nil) && (tileTicks.Type == nbt.TAG_Compound) {
		for _, _tTick := range tileTicks.Elems {
			tTick := _tTick.(nbt.TagCompound)
			x, y, z, err := extractCoord(tTick)
			if err != nil {
				return nil, fmt.Errorf("Could not Extract coords: %s", err)
			}

			_, _, x, z = BlockToChunk(x, z)

			x %= 16
			z %= 16

			tick := TileTick{}
			if tick.i, err = tTick.GetInt("i"); err != nil {
				return nil, fmt.Errorf("Could not read i of a TileTag tag: %s", err)
			}
			if tick.t, err = tTick.GetInt("t"); err != nil {
				return nil, fmt.Errorf("Could not read t of a TileTag tag: %s", err)
			}
			switch tick.p, err = tTick.GetInt("p"); err {
			case nil:
				tick.hasP = true
			case nbt.NotFound:
				tick.hasP = false
			default:
				return nil, fmt.Errorf("Could not read p of a TileTag tag: %s", err)
			}

			c.blocks[calcBlockOffset(x, y, z)].Tick = &tick
		}
	}

	return &c, nil
}

func (c *Chunk) toPreChunk() (*preChunk, error) {
	terraPopulated := byte(0)
	if c.populated {
		terraPopulated = 1
	}
	lvl := nbt.TagCompound{
		"xPos":             nbt.NewIntTag(c.x),
		"zPos":             nbt.NewIntTag(c.z),
		"LastUpdate":       nbt.NewLongTag(c.lastUpdate),
		"TerrainPopulated": nbt.NewByteTag(terraPopulated),
		"InhabitedTime":    nbt.NewLongTag(c.inhabitedTime),
		"HeightMap":        nbt.NewIntArrayTag(c.heightMap),
		"Entities":         nbt.Tag{nbt.TAG_Compound, c.Entities},
	}

	hasBiomes := false
	biomes := make([]byte, 16*16)
	for i, bio := range c.biomes {
		if bio != BioUncalculated {
			hasBiomes = true
			break
		}
		biomes[i] = byte(bio)
	}
	if hasBiomes {
		lvl["Biomes"] = nbt.NewByteArrayTag(biomes)
	}

	sections := make([]nbt.TagCompound, 0)
	tileEnts := make([]nbt.TagCompound, 0)
	tileTicks := make([]nbt.TagCompound, 0)

	for subchunk := 0; subchunk < 16; subchunk++ {
		off := subchunk * 4096

		blocks := make([]byte, 4096)
		add := make([]byte, 2048)
		data := make([]byte, 2048)
		blockLight := make([]byte, 2048)
		skyLight := make([]byte, 2048)

		allAir, addEmpty := true, true
		for i := 0; i < 4096; i++ {
			blk := c.blocks[i+off]
			id := blk.ID
			if id != BlkAir {
				allAir = false
			}

			blocks[i] = byte(id & 0xff)
			idH := byte(id >> 8)
			if idH != 0 {
				addEmpty = false
			}
			setHalfbyte(add, i, idH)

			setHalfbyte(data, i, blk.Data)
			setHalfbyte(blockLight, i, blk.BlockLight)
			setHalfbyte(skyLight, i, blk.SkyLight)

			x, y, z := offsetToPos(i + off)
			x, z = ChunkToBlock(int(c.x), int(c.z), x, z)

			if (blk.TileEntity != nil) && (len(blk.TileEntity) > 0) {
				// Fix coords
				blk.TileEntity["x"] = nbt.NewIntTag(int32(x))
				blk.TileEntity["y"] = nbt.NewIntTag(int32(y))
				blk.TileEntity["z"] = nbt.NewIntTag(int32(z))
				tileEnts = append(tileEnts, blk.TileEntity)
			}

			if blk.Tick != nil {
				tileTick := nbt.TagCompound{
					"x": nbt.NewIntTag(int32(x)),
					"y": nbt.NewIntTag(int32(y)),
					"z": nbt.NewIntTag(int32(z)),
					"i": nbt.NewIntTag(blk.Tick.i),
					"t": nbt.NewIntTag(blk.Tick.t),
				}
				if blk.Tick.hasP {
					tileTick["p"] = nbt.NewIntTag(blk.Tick.p)
				}
				tileTicks = append(tileTicks, tileTick)
			}
		}

		if !allAir {
			comp := nbt.TagCompound{
				"Y":          nbt.NewByteTag(byte(subchunk)),
				"Blocks":     nbt.NewByteArrayTag(blocks),
				"Data":       nbt.NewByteArrayTag(data),
				"BlockLight": nbt.NewByteArrayTag(blockLight),
				"SkyLight":   nbt.NewByteArrayTag(skyLight),
			}
			if !addEmpty {
				comp["Add"] = nbt.NewByteArrayTag(add)
			}
			sections = append(sections, comp)
		}
	}

	lvl["Sections"] = nbt.NewListTag(nbt.TAG_Compound, sections)

	if len(c.Entities) > 0 {
		lvl["Entities"] = nbt.NewListTag(nbt.TAG_Compound, c.Entities)
	} else {
		lvl["Entities"] = nbt.NewListTag(nbt.TAG_Byte, []byte{})
	}
	if len(tileEnts) > 0 {
		lvl["TileEntities"] = nbt.NewListTag(nbt.TAG_Compound, tileEnts)
	} else {
		lvl["TileEntities"] = nbt.NewListTag(nbt.TAG_Byte, []byte{})
	}
	if len(tileTicks) > 0 {
		lvl["TileTicks"] = nbt.NewListTag(nbt.TAG_Compound, tileTicks)
	}

	root := nbt.Tag{nbt.TAG_Compound, nbt.TagCompound{
		"Level": nbt.Tag{nbt.TAG_Compound, lvl},
	}}

	buf := new(bytes.Buffer)
	if err := nbt.WriteZlibdNamedTag(buf, "", root); err != nil {
		return nil, err
	}

	return &preChunk{
		ts:          c.ts,
		data:        buf.Bytes(),
		compression: compressZlib,
	}, nil
}
