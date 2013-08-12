package main

import (
	"flag"
	"fmt"
	"github.com/kch42/gomcmap/mcmap"
	"image"
	"image/png"
	"os"
)

func main() {
	path := flag.String("path", "", "Path to region directory")
	output := flag.String("output", "map.png", "File to write image to")
	flag.Parse()

	if *path == "" {
		flag.Usage()
		os.Exit(1)
	}

	region, err := mcmap.OpenRegion(*path, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open region: %s\n", err)
		os.Exit(1)
	}

	xmin, xmax, zmin, zmax := region.MaxDims()
	w := (xmax - xmin) * 16
	h := (zmax - zmin) * 16
	img := image.NewRGBA(image.Rect(0, 0, w, h))

chunkLoop:
	for chunkPos := range region.AllChunks() {
		cx, cz := chunkPos.X, chunkPos.Z
		chunk, err := region.Chunk(cx, cz)
		switch err {
		case nil:
		case mcmap.NotAvailable:
			continue chunkLoop
		default:
			fmt.Fprintf(os.Stderr, "Error while getting chunk (%d, %d): %s\n", cx, cz, err)
			os.Exit(1)
		}

		for x := 0; x < 16; x++ {
		scanZ:
			for z := 0; z < 16; z++ {
				ax, az := mcmap.ChunkToBlock(cx, cz, x, z)
				for y := 255; y >= 0; y-- {
					blk := chunk.Block(x, y, z)
					c, ok := colors[blk.ID]
					if ok {
						img.Set(ax-(xmin*16), az-(zmin*16), c)
						continue scanZ
					}
				}
				img.Set(ax-(xmin*16), az-(zmin*16), rgb(0x000000))
			}
		}

		if err := region.UnloadChunk(cx, cz); err != nil {
			fmt.Fprintf(os.Stderr, "Error while unloading chunk %d, %d: %s\n", cx, cz, err)
			os.Exit(1)
		}
	}

	f, err := os.Create(*output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not write image: %s", err)
		os.Exit(1)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		fmt.Fprintf(os.Stderr, "Could not write image: %s", err)
		os.Exit(1)
	}
}
