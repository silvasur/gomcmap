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

func (co chunkOffTs) calcLocationEntry() uint32 {
	return uint32((co.size>>12)&0xff) | (uint32(co.offset>>12) << 8)
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

	for z := 0; z < superchunkSizeXZ; z++ {
		for x := 0; x < superchunkSizeXZ; x++ {
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

	for z := 0; z < superchunkSizeXZ; z++ {
		for x := 0; x < superchunkSizeXZ; x++ {
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

func (pc *preChunk) writePreChunk(w io.Writer) error {
	length := uint32(len(pc.data) + 1)
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}
	if _, err := w.Write([]byte{pc.compression}); err != nil {
		return err
	}
	_, err := w.Write(pc.data)
	return err
}

func writeRegionFile(w io.Writer, pcs map[XZPos]*preChunk) error {
	offs := make(map[XZPos]chunkOffTs)
	buf := new(bytes.Buffer)
	pw := kagus.NewPaddedWriter(buf, sectorSize)

	for pos, pc := range pcs {
		off := buf.Len()
		if err := pc.writePreChunk(pw); err != nil {
			return err
		}
		if err := pw.Pad(); err != nil {
			return err
		}
		offs[pos] = chunkOffTs{
			offset: int64(2*sectorSize + off),
			size:   int64(buf.Len() - off),
			ts:     pc.ts,
		}
	}

	for z := 0; z < superchunkSizeXZ; z++ {
		for x := 0; x < superchunkSizeXZ; x++ {
			off := uint32(0)
			if cOff, ok := offs[XZPos{x, z}]; ok {
				off = cOff.calcLocationEntry()
			}

			if err := binary.Write(w, binary.BigEndian, off); err != nil {
				return err
			}
		}
	}

	for z := 0; z < superchunkSizeXZ; z++ {
		for x := 0; x < superchunkSizeXZ; x++ {
			ts := int32(0)
			if cOff, ok := offs[XZPos{x, z}]; ok {
				ts = int32(cOff.ts.Unix())
			}

			if err := binary.Write(w, binary.BigEndian, ts); err != nil {
				return err
			}
		}
	}

	_, err := io.Copy(w, buf)
	return err
}
