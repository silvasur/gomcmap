// emeraldfinder is a gomcmap demo program to find emerald ores in a Minecraft map.
package main

import (
	"flag"
	"fmt"
	"github.com/silvasur/gomcmap/mcmap"
	"os"
)

func main() {
	path := flag.String("path", "", "Path to region directory")
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

		chunk.Iter(func(x, y, z int, blk *mcmap.Block) {
			if blk.ID == mcmap.BlkEmeraldOre {
				absx, absz := mcmap.ChunkToBlock(cx, cz, x, z)
				fmt.Printf("%d, %d, %d\n", absx, y, absz)
			}
		})

		chunk.MarkUnused()
	}
}
