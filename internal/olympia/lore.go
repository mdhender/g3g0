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

// loreFunction is complete lore.c:9–44.
func (w *world) loreFunction(who int, fn string) {
	switch fn {
	case "capturable_animals", "animal_fighters":
		for i := range w.loopKind(tItem) {
			p := w.rpItem(i)
			if p == nil || p.animal == 0 {
				continue
			}
			if fn == "capturable_animals" {
				if p.capturable != 0 {
					w.commandOut(who, "    %s", w.boxName(i))
				}
			} else if p.attack != 0 || p.defense != 0 || p.missile != 0 || i == itemGhostWarrior {
				w.commandOut(who, "    %-20s %s", w.boxName(i), sout("(%d,%d,%d)", p.attack, p.defense, p.missile))
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
	var path string
	if w.kind(num) == tSkill {
		path = sout("%s/lore/%d/%d", w.libdir, w.skillSchool(num), num)
	} else {
		path = sout("%s/lore/etc/%d", w.libdir, num)
	}
	data, err := clib.ReadFile(path)
	if err != nil {
		w.commandOut(who, "<lore sheet not available>")
		fmt.Fprintf(w.stderr, "can't open %s: ", path)
		clib.Perror(w.stderr, "", err)
		return
	}
	code := boxCodeLess(displayNum)
	fp := bytes.NewReader(data)
	for line := getlin(fp); line != nil; line = getlin(fp) {
		if strings.HasPrefix(*line, "#") {
			continue
		}
		if strings.HasPrefix(*line, "$") {
			w.loreFunction(who, (*line)[1:])
			continue
		}
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
		w.commandOut(who, "%s", text.String())
	}
}

// deliverLore is complete lore.c:202–221.
func (w *world) deliverLore(who, num int) {
	switch w.kind(num) {
	case tSkill:
		w.deliverSkillLore(who, num, false)
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

// outSkillLine and deliverSkillLore are lore.c:137–199.
func (w *world) outSkillLine(who, sk int) {
	w.commandOut(who, "%-*s  %-34s %s%s", charField, boxCodeLess(sk), w.justName(sk), weeks(w.learnTime(sk)), w.npReqS(sk))
}

func (w *world) deliverSkillLore(who, sk int, showResearch bool) {
	w.deliverLoreSheet(who, sk, sk)
	p := w.rpSkill(sk)
	if p != nil && len(p.offered) != 0 {
		w.commandOut(who, "")
		w.commandOut(who, "The following skills may be studied directly once %s is known:", w.justName(sk))
		w.commandOut(who, "")
		w.indent += 3
		w.commandOut(who, "%-*s  %-34s %13s", charField, "num", "skill", "time to learn")
		w.commandOut(who, "%-*s  %-34s %13s", charField, "---", "-----", "-------------")
		for _, i := range p.offered {
			w.outSkillLine(who, i)
		}
		w.indent -= 3
		w.commandOut(who, "")
		if len(p.research) != 0 {
			w.commandOut(who, "Further skills may be found through research.")
		}
	} else if p != nil && len(p.research) != 0 {
		w.commandOut(who, "")
		w.commandOut(who, "Skills within this category may be found through research.")
	}
	if showResearch && p != nil && len(p.research) != 0 {
		w.commandOut(who, "")
		w.indent += 3
		for _, i := range p.research {
			w.outSkillLine(who, i)
		}
		w.indent -= 3
	}
}

// showLoreSheets is lore.c:258–294. Sort the original queue and retain duplicates.
func (w *world) showLoreSheets() {
	w.stage("show_lore_sheets()")
	w.outPath, w.outAltWho = master, 2 // OUT_LORE
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
			w.deliverLore(pl, num)
		}
	}
	w.outPath, w.outAltWho = 0, 0
}

// gmShowAllSkills is lore.c:297–369. Complete the listing before any reads.
func (w *world) gmShowAllSkills(pl int) {
	w.outPath, w.outAltWho = master, 2 // OUT_LORE
	w.commandOut(pl, "")
	w.commandOut(pl, "Skill listing:")
	w.commandOut(pl, "--------------")
	w.commandOut(pl, "")
	w.commandOut(pl, "Skill schools:")
	w.commandOut(pl, "")
	for sk := range w.loopKind(tSkill) {
		if w.skillSchool(sk) == sk {
			w.outSkillLine(pl, sk)
		}
	}
	w.commandOut(pl, "")
	for sk := range w.loopKind(tSkill) {
		if w.skillSchool(sk) == sk {
			w.commandOut(pl, "%s", w.boxName(sk))
			w.indent += 3
			for i := range w.loopKind(tSkill) {
				if w.skillSchool(i) == sk && i != sk {
					w.outSkillLine(pl, i)
				}
			}
			w.indent -= 3
			w.commandOut(pl, "")
		}
	}
	for _, schools := range []bool{true, false} {
		for sk := range w.loopKind(tSkill) {
			if (w.skillSchool(sk) == sk) == schools {
				w.deliverSkillLore(pl, sk, true)
				w.commandOut(pl, "")
			}
		}
	}
	w.outPath, w.outAltWho = 0, 0
}
