// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mdhender/g3g0/internal/clib"
)

// npReqS is complete lore.c:120–134; absent skill records have no NP cost.
func (w *world) npReqS(skill int) string {
	np := int32(0)
	if p := w.rpSkill(skill); p != nil {
		np = p.npReq
	}
	if np < 1 {
		return ""
	}
	if np == 1 {
		return ", 1 NP req'd"
	}
	return sout(", %d NP req'd", np)
}

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

// loreFunction is complete lore.c:9–44. A nil emit is used only by the
// retained stage-1 queue/GM projections, whose rendering belongs to #160.
func (w *world) loreFunction(fn string, emit func(string, ...any)) {
	switch fn {
	case "capturable_animals", "animal_fighters":
		if emit == nil {
			return
		}
		for i := range w.loopKind(tItem) {
			p := w.rpItem(i)
			if p == nil || p.animal == 0 {
				continue
			}
			if fn == "capturable_animals" {
				if p.capturable != 0 {
					emit("    %s", w.boxName(i))
				}
			} else if p.attack != 0 || p.defense != 0 || p.missile != 0 || i == itemGhostWarrior {
				emit("    %-20s %s", w.boxName(i), sout("(%d,%d,%d)", p.attack, p.defense, p.missile))
			}
		}
	default:
		fmt.Fprintf(w.stderr, "bad lore sheet function: %s\n", fn)
		assert(false, "olympia/lore.c", 43, "FALSE")
	}
}

// deliverLoreSheet is complete lore.c:47–117. Header output precedes path
// validation/open, and display-code conversion follows a successful open.
func (w *world) deliverLoreSheet(who, num, displayNum int) {
	w.commandOut(who, "")
	w.matchLines(who, w.boxName(displayNum))
	w.loreSheetInput(num, displayNum, func(format string, args ...any) {
		w.commandOut(who, format, args...)
	})
}

// loreSheetInput is the single reader shared with the stage-1 projections.
// They retain real input/errors until #160 replaces them; they do not invoke
// real delivery or manufacture output bits. No sheet is opened twice.
func (w *world) loreSheetInput(num, displayNum int, emit func(string, ...any)) {
	var path string
	if w.kind(num) == tSkill {
		path = sout("%s/lore/%d/%d", w.libdir, w.skillSchool(num), num)
	} else {
		path = sout("%s/lore/etc/%d", w.libdir, num)
	}
	data, err := clib.ReadFile(path)
	if err != nil {
		if emit != nil {
			emit("<lore sheet not available>")
		}
		fmt.Fprintf(w.stderr, "can't open %s: ", path)
		clib.Perror(w.stderr, "", err)
		return
	}
	code := ""
	if emit != nil {
		code = boxCodeLess(displayNum)
	}
	fp := bytes.NewReader(data)
	for line := getlin(fp); line != nil; line = getlin(fp) {
		if strings.HasPrefix(*line, "#") {
			continue
		}
		if strings.HasPrefix(*line, "$") {
			w.loreFunction((*line)[1:], emit)
			continue
		}
		if emit != nil {
			var text strings.Builder
			for p := 0; p < len(*line); {
				if (*line)[p] == '@' {
					text.WriteString(code)
					for p < len(*line) && (*line)[p] == '@' {
						p++
					}
				} else {
					text.WriteByte((*line)[p])
					p++
				}
			}
			emit("%s", text.String())
		}
	}
}

// deliverLore is complete lore.c:202–221, except for the explicit #160
// dependency: real skill delivery must not silently succeed without its
// offered/research rendering. The stage-1 readers below remain separate roots.
func (w *world) deliverLore(who, num int) {
	switch w.kind(num) {
	case tSkill:
		panic(notImplemented("deliver_skill_lore (#160)"))
	case tItem:
		lore := 0
		if p := w.rpItemMagic(num); p != nil {
			lore = int(p.lore)
		}
		w.deliverLoreSheet(who, lore, num)
	default:
		w.deliverLoreSheet(who, num, num)
	}
	w.commandOut(who, "")
}

func (w *world) deliverLoreState(num int) {
	if w.kind(num) == tItem {
		if p := w.rpItemMagic(num); p != nil {
			num = int(p.lore)
		} else {
			num = 0
		}
	}
	w.loreSheetInput(num, num, nil)
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
				w.loreSheetInput(sk, sk, nil)
			}
		}
	}
}
