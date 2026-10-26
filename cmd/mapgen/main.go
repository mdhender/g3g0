// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Command mapgen is a port of G3's mapgen-g3 map generator.
//
// Like the C program it takes no arguments. It reads Map, Regions, Land,
// Cities and randseed from the working directory, writes loc, gate and road
// there, and replaces randseed with the advanced seed.
package main

import (
	"os"

	"github.com/mdhender/g3g0/internal/mapgen"
)

func main() {
	os.Exit(mapgen.Run(".", os.Stdout, os.Stderr))
}
