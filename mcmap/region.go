package mcmap

import (
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

var (
	NotAvailable = errors.New("Chunk or Superchunk not available")
	AlreadyThere = errors.New("Chunk is already there")
)

type superchunk struct {
	preChunks map[XZPos]*preChunk
	chunks    map[XZPos]*Chunk
	modified  bool
}

type Region struct {
	path             string
	autosave         bool
	superchunksAvail map[XZPos]bool
	superchunks      map[XZPos]*superchunk
}

var mcaRegex = regexp.MustCompile(`^r\.([0-9-]+)\.([0-9-]+)\.mca$`)

// OpenRegion opens a region directory. If autosave is true, mcmap will save modified and unloaded chunks automatically to reduce memory usage. You still have to call Save at the end.
//
// You can also use OpenRegion to create a new region. Yust make sure the path exists.
func OpenRegion(path string, autosave bool) (*Region, error) {
	rv := &Region{
		path:             path,
		superchunksAvail: make(map[XZPos]bool),
		superchunks:      make(map[XZPos]*superchunk),
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", path)
	}

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		match := mcaRegex.FindStringSubmatch(name)
		if len(match) == 3 {
			// We ignore the error here. The Regexp already ensures that the inputs are numbers.
			x, _ := strconv.ParseInt(match[1], 10, 32)
			z, _ := strconv.ParseInt(match[2], 10, 32)

			rv.superchunksAvail[XZPos{int(x), int(z)}] = true
		}
	}

	return rv, nil
}

// MaxDims calculates the approximate maximum x, z dimensions of this region in number of chunks. The actual maximum dimensions might be a bit smaller.
func (reg *Region) MaxDims() (xmin, xmax, zmin, zmax int) {
	if len(reg.superchunksAvail) == 0 {
		return 0, 0, 0, 0
	}

	xmin = math.MaxInt32
	zmin = math.MaxInt32
	xmax = math.MinInt32
	zmax = math.MinInt32

	for pos := range reg.superchunksAvail {
		if pos.X < xmin {
			xmin = pos.X
		}
		if pos.Z < zmin {
			zmin = pos.Z
		}
		if pos.X > xmax {
			xmax = pos.X
		}
		if pos.Z > zmax {
			zmax = pos.Z
		}
	}

	xmax++
	zmax++
	xmin *= superchunkSizeXZ
	xmax *= superchunkSizeXZ
	zmin *= superchunkSizeXZ
	zmax *= superchunkSizeXZ
	return
}

func chunkToSuperchunk(cx, cz int) (scx, scz, rx, rz int) {
	scx = cx >> 5
	scz = cz >> 5
	rx = ((cx % superchunkSizeXZ) + superchunkSizeXZ) % superchunkSizeXZ
	rz = ((cz % superchunkSizeXZ) + superchunkSizeXZ) % superchunkSizeXZ
	return
}

func superchunkToChunk(scx, scz, rx, rz int) (cx, cz int) {
	cx = scx*superchunkSizeXZ + rx
	cz = scz*superchunkSizeXZ + rz
	return
}

func (reg *Region) loadSuperchunk(pos XZPos) error {
	if !reg.superchunksAvail[pos] {
		return NotAvailable
	}
	fname := fmt.Sprintf("%s%cr.%d.%d.mca", reg.path, os.PathSeparator, pos.X, pos.Z)

	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	pcs, err := readRegionFile(f)
	if err != nil {
		return err
	}

	reg.superchunks[pos] = &superchunk{
		preChunks: pcs,
		chunks:    make(map[XZPos]*Chunk),
	}
	return nil
}

func (reg *Region) cleanSuperchunks(forceSave bool) error {
	del := make(map[XZPos]bool)

	for scPos, sc := range reg.superchunks {
		if len(sc.chunks) > 0 {
			continue
		}

		if sc.modified {
			if !(reg.autosave || forceSave) {
				continue
			}
			fn := fmt.Sprintf("%s%cr.%d.%d.mca", reg.path, os.PathSeparator, scPos.X, scPos.Z)
			f, err := os.Create(fn)
			if err != nil {
				return err
			}
			defer f.Close()

			if err := writeRegionFile(f, sc.preChunks); err != nil {
				return err
			}
		}

		del[scPos] = true
	}

	for scPos, _ := range del {
		delete(reg.superchunks, scPos)
	}

	return nil
}

func (sc *superchunk) loadChunk(reg *Region, rx, rz int) (*Chunk, error) {
	cPos := XZPos{rx, rz}

	if chunk, ok := sc.chunks[cPos]; ok {
		return chunk, nil
	}

	pc, ok := sc.preChunks[cPos]
	if !ok {
		return nil, NotAvailable
	}

	chunk, err := pc.toChunk(reg)
	if err != nil {
		return nil, err
	}
	sc.chunks[cPos] = chunk
	return chunk, nil
}

// Chunk returns the chunk at x, z. If no chunk could be found, the error NotAvailable will be returned. Other errors indicate an internal error (I/O error, file format violated, ...)
func (reg *Region) Chunk(x, z int) (*Chunk, error) {
	scx, scz, cx, cz := chunkToSuperchunk(x, z)
	scPos := XZPos{scx, scz}

	sc, ok := reg.superchunks[scPos]
	if !ok {
		if err := reg.loadSuperchunk(scPos); err != nil {
			return nil, err
		}
		sc = reg.superchunks[scPos]
	}

	chunk, err := sc.loadChunk(reg, cx, cz)
	if err != nil {
		return nil, err
	}

	if err := reg.cleanSuperchunks(false); err != nil {
		return nil, err
	}

	return chunk, nil
}

func (reg *Region) unloadChunk(x, z int) error {
	scx, scz, cx, cz := chunkToSuperchunk(x, z)
	scPos := XZPos{scx, scz}
	cPos := XZPos{cx, cz}

	sc, ok := reg.superchunks[scPos]
	if !ok {
		return nil
	}

	chunk, ok := sc.chunks[cPos]
	if !ok {
		return nil
	}

	if chunk.modified {
		pc, err := chunk.toPreChunk()
		if err != nil {
			return err
		}
		sc.preChunks[cPos] = pc

		chunk.modified = false
		sc.modified = true
	}

	delete(sc.chunks, cPos)

	return nil
}

// AllChunks returns a channel that will give you the positions of all possibly available chunks in an efficient order.
//
// Note the "possibly available", you still have to check, if the chunk could actually be loaded.
func (reg *Region) AllChunks() <-chan XZPos {
	ch := make(chan XZPos)
	go func(ch chan<- XZPos) {
		for spos, _ := range reg.superchunksAvail {
			scx, scz := spos.X, spos.Z
			for rx := 0; rx < superchunkSizeXZ; rx++ {
				for rz := 0; rz < superchunkSizeXZ; rz++ {
					cx, cz := superchunkToChunk(scx, scz, rx, rz)
					ch <- XZPos{cx, cz}
				}
			}
		}
		close(ch)
	}(ch)

	return ch
}

// NewChunk adds a new, blank chunk. If the Chunk is already there, error AlreadyThere will be returned.
// Other errors indicate internal errors.
func (reg *Region) NewChunk(cx, cz int) (*Chunk, error) {
	scx, scz, rx, rz := chunkToSuperchunk(cx, cz)

	scPos := XZPos{scx, scz}

	var sc *superchunk
	if reg.superchunksAvail[scPos] {
		var ok bool
		if sc, ok = reg.superchunks[scPos]; !ok {
			if err := reg.loadSuperchunk(scPos); err != nil {
				return nil, err
			}
			sc = reg.superchunks[scPos]
		}
	} else {
		sc = &superchunk{
			chunks:    make(map[XZPos]*Chunk),
			preChunks: make(map[XZPos]*preChunk),
			modified:  true,
		}
		reg.superchunksAvail[scPos] = true
		reg.superchunks[scPos] = sc
	}

	switch chunk, err := sc.loadChunk(reg, rx, rz); err {
	case nil:
		chunk.MarkUnused()
		return nil, AlreadyThere
	case NotAvailable:
	default:
		return nil, err
	}

	cPos := XZPos{rx, rz}
	chunk := newChunk(reg, cx, cz)

	pc, err := chunk.toPreChunk()
	if err != nil {
		return nil, err
	}

	sc.preChunks[cPos] = pc
	sc.chunks[cPos] = chunk
	sc.modified = true

	return chunk, nil
}

// Save saves modified and unused chunks.
func (reg *Region) Save() error {
	return reg.cleanSuperchunks(true)
}
