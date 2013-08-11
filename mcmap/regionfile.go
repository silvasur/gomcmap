package mcmap

import (
	"bytes"
	"encoding/binary"
	"github.com/kch42/kagus"
	"io"
	"time"
)

const sectorSize = 4096

type chunkOffTs struct {
	offset, size int64
	ts           time.Time
}

func (cOff chunkOffTs) readPreChunk(r io.ReadSeeker) (*preChunk, error) {
	pc := preChunk{ts: cOff.ts}

	if _, err := r.Seek(cOff.offset, 0); err != nil {
		return nil, err
	}

	lr := io.LimitReader(r, cOff.size)

	var length uint32
	if err := binary.Read(lr, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	lr = io.LimitReader(lr, int64(length))

	compType, err := kagus.ReadByte(lr)
	if err != nil {
		return nil, err
	}
	pc.compression = compType

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, lr); err != nil {
		return nil, err
	}
	pc.data = buf.Bytes()

	return &pc, err

}

func readRegionFile(r io.ReadSeeker) (map[XZPos]*preChunk, error) {
	if _, err := r.Seek(0, 0); err != nil {
		return nil, err
	}

	offs := make(map[XZPos]*chunkOffTs)

	for z := 0; z < 32; z++ {
		for x := 0; x < 32; x++ {
			var location uint32
			if err := binary.Read(r, binary.BigEndian, &location); err != nil {
				return nil, err
			}

			if location == 0 {
				continue
			}

			offs[XZPos{x, z}] = &chunkOffTs{
				offset: int64((location >> 8) * sectorSize),
				size:   int64((location & 0xff) * sectorSize),
			}
		}
	}

	for z := 0; z < 32; z++ {
		for x := 0; x < 32; x++ {
			pos := XZPos{x, z}

			var ts int32
			if err := binary.Read(r, binary.BigEndian, &ts); err != nil {
				return nil, err
			}

			if _, ok := offs[pos]; !ok {
				continue
			}

			offs[pos].ts = time.Unix(int64(ts), 0)
		}
	}

	preChunks := make(map[XZPos]*preChunk)
	for pos, cOff := range offs {
		pc, err := cOff.readPreChunk(r)
		if err != nil {
			return nil, err
		}
		preChunks[pos] = pc
	}

	return preChunks, nil
}
