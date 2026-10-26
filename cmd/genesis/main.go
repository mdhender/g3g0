// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Command genesis is a port of G3's genesis-g3 world generator.
//
// It takes the same flags as the C program and writes a fresh mapgen input
// set (Map, Regions, Land, Cities and randseed) to the working directory.
// Run it with --help for the flags.
package main

import (
	"os"

	"github.com/mdhender/g3g0/internal/genesis"
)

func main() {
	os.Exit(genesis.Run(".", os.Args, os.Stderr))
}
