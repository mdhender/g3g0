// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"os"
	"sort"
)

// coverageCount follows coverage.c: allocation precedes validation, but the
// disabled/immediate guards precede both. A non-nil empty path still counts.
func (w *world) coverageCount(cmd, outcome int) {
	if w.coveragePath == nil || w.immediate {
		return
	}
	if w.covCounts == nil {
		w.covCounts = make([][covMax]int32, len(cmdTbl))
	}
	if cmd < 0 || cmd >= len(w.covCounts) {
		cmd = 0
	}
	assert(outcome >= 0 && outcome < covMax, "olympia/coverage.c", 62, "outcome >= 0 && outcome < COV_MAX")
	w.covCounts[cmd][outcome]++
}

func (w *world) writeCoverage() {
	if w.coveragePath == nil {
		return
	}
	fp, err := os.Create(*w.coveragePath)
	if err != nil {
		fmt.Fprintf(w.stderr, "Can't write %s!\n", *w.coveragePath)
		return
	}
	defer fp.Close()
	fmt.Fprint(fp, "verb\trejected\tunimplemented\tok\tfail\n")
	// Go omits C's terminating NULL-name sentinel; every actual row, including
	// the empty-name placeholder and mail-only verbs, remains in cmdTbl.
	order := make([]int, len(cmdTbl))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return cmdTbl[order[i]].name < cmdTbl[order[j]].name })
	for _, cmd := range order {
		name := cmdTbl[cmd].name
		if cmd == 0 {
			name = "(unrecognized)"
		}
		var counts [covMax]int32
		if w.covCounts != nil {
			counts = w.covCounts[cmd]
		}
		fmt.Fprintf(fp, "%s\t%d\t%d\t%d\t%d\n", name, counts[covRejected], counts[covUnimpl], counts[covOk], counts[covFail])
	}
}
