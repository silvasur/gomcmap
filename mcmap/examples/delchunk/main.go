// delchunk deletes the chunk at 200, 200
package main

import (
	"flag"
	"fmt"
	"github.com/kch42/gomcmap/mcmap"
	"os"
)

func main() {
	path := flag.String("path", "", "Path to region directory")
	flag.Parse()

	if *path == "" {
		flag.Usage()
		os.Exit(1)
	}

	region, err := mcmap.OpenRegion(*path, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open region: %s\n", err)
		os.Exit(1)
	}

	chunk, err := region.Chunk(200, 200)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting chunk 200,200: %s\n", err)
		os.Exit(1)
	}

	chunk.MarkDeleted()
	if err := chunk.MarkUnused(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not mark chunk as unused: %s\n", err)
		os.Exit(1)
	}

	if err := region.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not save region: %s\n", err)
		os.Exit(1)
	}
}
