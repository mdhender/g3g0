// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// hasAuraculum is art.c:7–18 has_auraculum.
func (w *world) hasAuraculum(who int) int {
	var ac int
	if p := w.rpMagic(who); p != nil {
		ac = int(p.auraculum)
	}
	if ac != 0 && w.validBox(ac) && w.hasItem(who, ac) > 0 {
		return ac
	}
	return 0
}

// maxEffAura is art.c:25–50 max_eff_aura. Bonuses are per positive
// inventory entry, not per item quantity; only innate aura is clamped.
func (w *world) maxEffAura(who int) int32 {
	var a int32
	if p := w.rpMagic(who); p != nil {
		a = p.maxAura
	}
	if a < 0 {
		a = 0
	}
	if ac := w.hasAuraculum(who); ac != 0 {
		if p := w.rpItemMagic(ac); p != nil {
			a += int32(p.aura)
		}
	}
	assert(w.validBox(who), "olympia/art.c", 41, "valid_box(who)")
	for e := range w.loopInv(who) {
		if p := w.rpItemMagic(int(e.item)); p != nil {
			a += int32(p.auraBonus)
		}
	}
	return a
}

// tokenPlayer is art.c:1016–1029. Nonregular owners share the independent
// token faction, including noncharacter holders and NPC players.
func (w *world) tokenPlayer(owner int) int {
	if w.kind(owner) != tChar {
		return int(w.indepPlayer)
	}
	pl := w.player(owner)
	if w.subkind(pl) != subPlRegular {
		return int(w.indepPlayer)
	}
	return pl
}

func (w *world) swearTokenUnits(item, target int) {
	w.logWrite(13, "%s got npc token %s", w.boxName(target), w.boxName(item))
	assert(w.subkind(item) == subNpcToken, "olympia/art.c", 1040, "subkind(item) == sub_npc_token")
	for i := range w.loopUnits(item) {
		if w.kind(i) == tChar {
			w.logWrite(13, "   swearing %s", w.boxNameKind(i))
			w.setLord(i, target, loyUnchanged, 0)
		}
	}
}

func (w *world) addTokenUnitSup(item int) {
	owner := w.itemUnique(item)
	assert(owner != 0, "olympia/art.c", 1092, "owner")
	where := w.province(owner)
	if w.subkind(where) == subOcean {
		where = w.subloc(owner)
	}
	ni := 0
	if p := w.rpItemMagic(item); p != nil {
		ni = int(p.tokenNi)
	}
	n := w.newChar(subNi, ni, where, -1, w.tokenPlayer(owner), loyNpc, 0, "")
	if n < 0 {
		w.logWrite(10, "  FAILed to add unit to token %s", boxCodeLess(item))
		return
	}
	w.logWrite(13, "  adding %s to %s", w.boxName(n), w.boxName(item))
	if w.beastCapturable(n) {
		w.pChar(n).breakPoint = 0
	}
	w.pMisc(n).cmdAllow = 'r'
	w.pMagic(n).token = int32(item)
	w.pPlayer(item).units.Append(n)
	w.commandOut(where, "%s appears.", w.boxName(n))
}

func (w *world) addTokenUnits(item int) {
	w.logWrite(13, "add_token_units(%s)", w.boxName(item))
	assert(w.subkind(item) == subNpcToken, "olympia/art.c", 1131, "subkind(item) == sub_npc_token")
	l := w.pPlayer(item).units.Len()
	// C advances the local count even when allocation fails; do not retry.
	for {
		num := 0
		if p := w.rpItemMagic(item); p != nil {
			num = int(p.tokenNum)
		}
		if l >= num {
			break
		}
		w.addTokenUnitSup(item)
		l++
	}
}

func (w *world) moveToken(item, from, to int) {
	toPl := w.tokenPlayer(to)
	w.logWrite(13, "Token %s moved from %s (%s) to %s (%s)", w.boxName(item), w.boxName(from), w.boxName(w.tokenPlayer(from)), w.boxName(to), w.boxName(toPl))
	if w.tokenPlayer(from) == int(w.indepPlayer) && w.pPlayer(item).units.Len() == 0 {
		num := int8(0)
		if p := w.rpItemMagic(item); p != nil {
			num = p.tokenNum
		}
		w.logWrite(11, "token %s from %d to 1.", boxCodeLess(item), num)
		w.pItemMagic(item).tokenNum = 1
	}
	if w.tokenPlayer(from) != toPl {
		w.swearTokenUnits(item, toPl)
		// The compiled C branch defers independent melting until reconciliation.
		if toPl != int(w.indepPlayer) {
			w.addTokenUnits(item)
		}
	}
}
