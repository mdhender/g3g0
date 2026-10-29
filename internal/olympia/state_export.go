// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/mdhender/g3g0/internal/lists"
)

// state_export.c at f73e481. Names are C byte strings, not Unicode: stop
// at NUL and encode each non-ASCII byte as Latin-1, including UTF-8 bytes.
// NULL and empty have the same JSON representation.
func jsonStr(fp io.Writer, s string) {
	fmt.Fprint(fp, "\"")
	for i := 0; i < len(s) && s[i] != 0; i++ {
		b := s[i]
		switch {
		case b == '"' || b == '\\':
			fmt.Fprintf(fp, "\\%c", b)
		case b < 0x20 || b >= 0x7f:
			fmt.Fprintf(fp, "\\u%04x", b)
		default:
			fmt.Fprintf(fp, "%c", b)
		}
	}
	fmt.Fprint(fp, "\"")
}

func sortIDs(l lists.IList) { slices.SortFunc(l, intComp) }

func (w *world) exportUnit(fp io.Writer, who int) {
	var items, skills lists.IList
	fmt.Fprintf(fp, "{\"id\": %d, \"name\": ", who)
	jsonStr(fp, w.name(who))
	fmt.Fprintf(fp, ", \"where\": %d, \"stacked_under\": %d", w.subloc(who), w.stackParent(who))
	assert(w.validBox(who), "olympia/state_export.c", 92, "valid_box(who)")
	for e := range w.loopInv(who) {
		items.Append(int(e.item))
	}
	sortIDs(items)
	fmt.Fprint(fp, ",\n      \"items\": [")
	for i, item := range items {
		if i != 0 {
			fmt.Fprint(fp, ", ")
		}
		fmt.Fprintf(fp, "{\"id\": %d, \"name\": ", item)
		jsonStr(fp, w.name(item))
		fmt.Fprintf(fp, ", \"qty\": %d}", w.hasItem(who, item))
	}
	fmt.Fprint(fp, "]")
	assert(w.validBox(who), "olympia/state_export.c", 110, "valid_box(who)")
	if p := w.rpChar(who); p != nil {
		// loop_char_skill_known copies the pointer list, not the entries.
		for _, s := range p.skills.Copy() {
			if s.know == skillKnow {
				skills.Append(int(s.skill))
			}
		}
	}
	sortIDs(skills)
	fmt.Fprint(fp, ",\n      \"skills\": [")
	for i, sk := range skills {
		if i != 0 {
			fmt.Fprint(fp, ", ")
		}
		fmt.Fprintf(fp, "{\"id\": %d, \"name\": ", sk)
		jsonStr(fp, w.name(sk))
		fmt.Fprint(fp, "}")
	}
	fmt.Fprint(fp, "]}")
}

// Location traversal is deliberately unfinished until #176. The CLI gate
// remains in place; direct faction/root calls must not require locations.
func (w *world) exportLocation(fp io.Writer, where, viewer int) {
	panic(notImplemented("export_location"))
}

func (w *world) exportFaction(fp io.Writer, pl int) {
	var units, wheres, viewers lists.IList
	// No additional character/ownership/prisoner filter: loop_units walks
	// the player's actual units list, even when it contains unusual entries.
	for who := range w.loopUnits(pl) {
		units.Append(who)
	}
	sortIDs(units)
	fmt.Fprintf(fp, "{\"id\": %d, \"name\": ", pl)
	jsonStr(fp, w.name(pl))
	fmt.Fprint(fp, ",\n   \"units\": [")
	for i, who := range units {
		if i != 0 {
			fmt.Fprint(fp, ",")
		}
		fmt.Fprint(fp, "\n     ")
		w.exportUnit(fp, who)
	}
	fmt.Fprint(fp, "]")
	for _, who := range units {
		where := w.subloc(who)
		if where > 0 && wheres.Lookup(where) < 0 {
			wheres.Append(where)
		}
	}
	sortIDs(wheres)
	for _, where := range wheres {
		for _, who := range units {
			if w.subloc(who) == where {
				viewers.Append(who)
				break
			}
		}
	}
	fmt.Fprint(fp, ",\n   \"locations\": [")
	for i, where := range wheres {
		if i != 0 {
			fmt.Fprint(fp, ",")
		}
		fmt.Fprint(fp, "\n     ")
		w.exportLocation(fp, where, viewers[i])
	}
	fmt.Fprint(fp, "]}")
}

func (w *world) writeStateExport() {
	if w.stateExportPath == nil {
		return
	}
	fp, err := os.Create(*w.stateExportPath)
	if err != nil {
		fmt.Fprintf(w.stderr, "Can't write %s!\n", *w.stateExportPath)
		return
	}
	defer fp.Close()
	var factions lists.IList
	for pl := range w.loopSubkind(subPlRegular) {
		factions.Append(pl)
	}
	sortIDs(factions)
	fmt.Fprintf(fp, "{\"turn\": %d, \"factions\": [", w.sysclock.turn)
	for i, pl := range factions {
		if i != 0 {
			fmt.Fprint(fp, ",")
		}
		fmt.Fprint(fp, "\n  ")
		w.exportFaction(fp, pl)
	}
	fmt.Fprint(fp, "]}\n")
}
