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

var useKeyNames = [...]string{
	"(invalid)", "death_potion", "heal_potion", "slave_potion", "palantir",
	"proj_cast", "quick_cast", "drum", "faery_stone", "orb", "barbarian_kill",
	"savage_kill", "corpse_kill", "orc_kill", "skeleton_kill", "bta_skull",
}

func (w *world) useCoverageSkill(ent, outcome int) {
	if w.useCoveragePath == nil || w.immediate {
		return
	}
	assert(outcome >= 0 && outcome < ucovMax, "olympia/coverage.c", 187, "outcome >= 0 && outcome < UCOV_MAX")
	if w.ucovSkills == nil {
		w.ucovSkills = make([][ucovMax]int32, len(useTbl)-1)
	}
	if ent <= 0 || ent >= len(w.ucovSkills) {
		w.ucovItems[0][outcome]++
	} else {
		w.ucovSkills[ent][outcome]++
	}
}

func (w *world) useCoverageItem(key, outcome int) {
	if w.useCoveragePath == nil || w.immediate {
		return
	}
	assert(outcome >= 0 && outcome < ucovMax, "olympia/coverage.c", 213, "outcome >= 0 && outcome < UCOV_MAX")
	if key < 0 || key > useKeyMax {
		key = 0
	}
	w.ucovItems[key][outcome]++
}

func (w *world) writeUseCoverage() {
	if w.useCoveragePath == nil {
		return
	}
	fp, err := os.Create(*w.useCoveragePath)
	if err != nil {
		fmt.Fprintf(w.stderr, "Can't write %s!\n", *w.useCoveragePath)
		return
	}
	defer fp.Close()
	fmt.Fprint(fp, "kind\tkey\tname\tchecks\tok\tfail\n")
	row := func(kind string, key int, name string, n [ucovMax]int32) {
		fmt.Fprintf(fp, "%s\t%d\t%s\t%d\t%d\t%d\n", kind, key, name, n[ucovChecks], n[ucovOk], n[ucovFail])
	}
	row("invalid", 0, useKeyNames[0], w.ucovItems[0])
	order := make([]int, len(useTbl)-2)
	for i := range order {
		order[i] = i + 1
	}
	sort.Slice(order, func(i, j int) bool { return useTblSkill(order[i]) < useTblSkill(order[j]) })
	for _, ent := range order {
		sk, name := useTblSkill(ent), ""
		if w.kind(sk) == tSkill {
			name = w.justName(sk)
		}
		var n [ucovMax]int32
		if w.ucovSkills != nil {
			n = w.ucovSkills[ent]
		}
		row("skill", sk, name, n)
	}
	for key := 1; key <= useKeyMax; key++ {
		row("item", key, useKeyNames[key], w.ucovItems[key])
	}
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
