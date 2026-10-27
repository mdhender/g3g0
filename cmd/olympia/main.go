// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Command olympia is a port of G3's olympia-g3 game engine.
//
// It takes the C engine's command line. So far -R (the random number
// generator self-test), -t (the ilist test), the usage errors, and loading
// an existing game and saving it again (-l dir and -S) are ported; every
// other mode stops with a "not yet implemented" message and a nonzero exit
// status.
package main

import (
	"os"

	"golang.org/x/term"

	"github.com/mdhender/g3g0/internal/olympia"
)

func main() {
	os.Exit(olympia.Run(os.Args, os.Stdout, os.Stderr, term.IsTerminal(int(os.Stdout.Fd()))))
}
