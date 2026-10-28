// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"

	"github.com/mdhender/g3g0/internal/lists"
)

// Ported from olympia/garr.c: the garrison ownership/output closure.

// rankS is garr.c:603–624 rank_s; ranks are the garr.h constants.
func (w *world) rankS(who int) string {
	var rank int8
	if p := w.rpChar(who); p != nil {
		rank = p.rank
	}
	switch rank {
	case 0:
		return ""
	case 10:
		return ", lord"
	case 20:
		return ", knight"
	case 30:
		return ", baron"
	case 40:
		return ", count"
	case 50:
		return ", earl"
	case 60:
		return ", marquess"
	case 70:
		return ", duke"
	case 80:
		return ", king"
	default:
		assert(false, "olympia/garr.c", 621, "FALSE")
	}
	return ""
}

// garrisonHere is garrison_here(): the first unit in where's here list if
// it is a garrison, or 0.
func (w *world) garrisonHere(where int) int {
	n := w.firstCharacter(where)

	if n != 0 && w.subkind(n) == subGarrison {
		return n
	}

	return 0
}

func (w *world) provinceAdmin(n int) int {
	var garr int
	if w.kind(n) == tLoc {
		assert(w.locDepth(n) == locProvince, "olympia/garr.c", 81, "loc_depth(n) == LOC_province")
		garr = w.garrisonHere(n)
		if garr == 0 {
			return 0
		}
	} else {
		assert(w.subkind(n) == subGarrison, "olympia/garr.c", 88, "subkind(n) == sub_garrison")
		garr = n
	}
	castle := 0
	if p := w.rpMisc(garr); p != nil {
		castle = int(p.garrCastle)
	}
	if !w.validBox(castle) {
		return 0
	}
	return w.buildingOwner(castle)
}

func (w *world) playersWhoRuleHere(where int) lists.IList {
	var l lists.IList
	loopCheck := 5000
	for i := w.provinceAdmin(where); i > 0; {
		// loop_loc_owner reads the next pledge before entering its body.
		next := 0
		if p := w.rpMagic(i); p != nil {
			next = int(p.pledge)
		}
		pl := w.player(i)
		if pl != 0 && l.Lookup(pl) < 0 {
			l.Append(pl)
		}
		if loopCheck <= 0 {
			for j, pl := range l {
				fmt.Fprintf(w.stderr, "l[%d] = %d\n", j, pl)
			}
			fmt.Fprintf(w.stderr, "where = %d, i = %d\n", where, i)
		}
		assert(loopCheck > 0, "olympia/garr.c", 195, "loop_check-- > 0")
		loopCheck--
		i = next
	}
	return l
}
