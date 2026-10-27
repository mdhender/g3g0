// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Command island is a port of G3's island-g3 island generator.
//
// Like the C program it reads a map on stdin, adds one random island, and
// prints the new map on stdout. It reads and advances randseed in the working
// directory. The flags are -b border_width, -c continental_shelf_width and
// -s size.
package main

import (
	"os"

	"github.com/mdhender/g3g0/internal/island"
)

func main() {
	os.Exit(island.Run(".", os.Args, os.Stdin, os.Stdout, os.Stderr))
}
