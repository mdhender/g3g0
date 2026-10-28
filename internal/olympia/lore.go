// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"bytes"
	"fmt"

	"github.com/mdhender/g3g0/internal/clib"
)

// setKnown is u.c set_known. Invalid entities are silently ignored by C.
func (w *world) setKnown(who, i int) {
	if !w.validBox(who) || !w.validBox(i) {
		return
	}
	pl := w.player(who)
	if !w.validBox(pl) {
		return
	}
	setBit(&w.pPlayer(pl).known, i)
}

// queueLore preserves insertion order and anyway duplicates. Delivery sorts
// this same queue and skips adjacent duplicates, without removing them.
func (w *world) queueLore(who, num int, anyway bool) {
	pl := w.player(who)
	assert(w.kind(pl) == tPlayer, "olympia/lore.c", 236, "kind(pl) == T_player")
	if !anyway && w.testKnown(pl, num) {
		return
	}
	w.pPlayer(pl).deliverLore.Append(num)
	w.setKnown(pl, num)
}

func (w *world) scanCharItemLore() {
	for who := range w.loopKind(tChar) {
		for e := range w.loopInv(who) {
			p := w.rpItemMagic(int(e.item))
			if p != nil && p.lore != 0 && !w.testKnown(who, int(e.item)) {
				w.queueLore(who, int(e.item), false)
			}
		}
	}
}

// scanCharSkillLore is lore.c:373–387. loop_char_skill copies ALL entries,
// unlike loop_char_skill_known: learning/don't-know, prisoners and silent
// factions all participate. Reuse queueLore's assertions and known bits.
func (w *world) scanCharSkillLore() {
	for who := range w.loopKind(tChar) {
		if p := w.rpChar(who); p != nil {
			for _, e := range p.skills.Copy() {
				w.queueLore(who, int(e.skill), false)
			}
		}
	}
}

// deliverLoreSheetState is the input/diagnostic projection of lore.c:49–123.
// Text, @ substitutions and the two legal $ expansions only write deferred
// report artifacts. Reading the sheet and validating functions cannot be
// deferred: an unknown function emits stderr and aborts, even in a silent
// faction. ReadFile preserves fopen success followed by a read error (e.g.
// a directory), and getlin preserves C's line, NUL and final-line semantics.
func (w *world) deliverLoreSheetState(num int) {
	var path string
	if w.kind(num) == tSkill {
		path = sout("%s/lore/%d/%d", w.libdir, w.skillSchool(num), num)
	} else {
		path = sout("%s/lore/etc/%d", w.libdir, num)
	}
	data, err := clib.ReadFile(path)
	if err != nil {
		fmt.Fprintf(w.stderr, "can't open %s: ", path)
		clib.Perror(w.stderr, "", err)
		return
	}
	fp := bytes.NewReader(data)
	for line := getlin(fp); line != nil; line = getlin(fp) {
		if len(*line) == 0 || (*line)[0] != '$' {
			continue
		}
		switch fn := (*line)[1:]; fn {
		case "capturable_animals", "animal_fighters":
			// Only artifact text, no state/diagnostic effects.
		default:
			fmt.Fprintf(w.stderr, "bad lore sheet function: %s\n", fn)
			assert(false, "olympia/lore.c", 43, "FALSE")
		}
	}
}

func (w *world) deliverLoreState(num int) {
	if w.kind(num) == tItem {
		if p := w.rpItemMagic(num); p != nil {
			num = int(p.lore)
		} else {
			num = 0
		}
	}
	w.deliverLoreSheetState(num)
}

// showLoreSheetsState is lore.c:259–294's state and raw diagnostic effects.
// It sorts the original persistent queue; duplicate deliveries are skipped,
// but duplicates and the whole queue survive for later calls.
func (w *world) showLoreSheetsState() {
	w.stage("show_lore_sheets()")
	for pl := range w.loopKind(tPlayer) {
		p := w.rpPlayer(pl)
		if p == nil || len(p.deliverLore) == 0 {
			continue
		}
		safeQsort(p.deliverLore, intComp)
		for i, num := range p.deliverLore {
			if i > 0 && num == p.deliverLore[i-1] {
				continue
			}
			w.deliverLoreState(num)
		}
	}
}

// gmShowAllSkillsState preserves gm_show_all_skills' player allocation and
// reads in school-then-nonschool order. It does not grant knowledge or queue
// lore. The preceding listing and offered/research lines are stage-2 text.
func (w *world) gmShowAllSkillsState(pl int) {
	w.reportOutputState(pl)
	// The listing traversals happen BEFORE any sheet reads. Keep their
	// skill_school assertions in order, even though the text is deferred.
	for sk := range w.loopKind(tSkill) {
		w.skillSchool(sk)
	}
	for sk := range w.loopKind(tSkill) {
		if w.skillSchool(sk) == sk {
			for i := range w.loopKind(tSkill) {
				w.skillSchool(i)
			}
		}
	}
	for _, schools := range []bool{true, false} {
		for sk := range w.loopKind(tSkill) {
			if (w.skillSchool(sk) == sk) == schools {
				w.deliverLoreSheetState(sk)
			}
		}
	}
}
