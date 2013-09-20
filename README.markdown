# gomcmap

A library for reading and modifying Minecraft maps.

## Install

	go get github.com/kch42/gomcmap/mcmap

## Examples

See `mcmap/examples`

## Compatibility

Currently only the Anvil map format is supported. Tested with Minecraft 1.6.2, will probably work with older versions too, but I haven't tested it.

## WARNING

Although I tested the library with some maps, I can't guarantee that everything always works (especially if you use mods in your Minecraft installation). So make a backup of your maps, just in case!

## Wishlist / TODO

* Removing chunks.
* Recalculating light data.
* Reading and modifying level.dat and other files (currently only the region files are used).
* Test compatibility with older versions of Minecraft.
