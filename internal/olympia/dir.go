// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"strings"

	"github.com/mdhender/g3g0/internal/rng"
)

// Ported from olympia/dir.c: map initialization and exit topology/selection.
//
// A province's box number gives its place on the map: 10101 is row 1,
// column 1, and [10199] row 1, column 99.

// regionRow is region_row(): the map row of location where.
func regionRow(where int) int {
	return (where / 100) % 100
}

// regionCol is region_col(): the map column of location where.
func regionCol(where int) int {
	return where % 100
}

// determineMapEdges is determine_map_edges(): it raises maxMapRow and
// maxMapCol to the largest row and column of every province, in chain
// order, that is not in Faery, Hades, the Cloudlands, the tunnels or
// Nowhere, and sets maxMapInit. A special region not yet created is ignored.
func (w *world) determineMapEdges() {
	for i := range w.loopProvince() {
		if (w.faeryRegion != 0 && w.region(i) == int(w.faeryRegion)) ||
			(w.hadesRegion != 0 && w.region(i) == int(w.hadesRegion)) ||
			(w.cloudRegion != 0 && w.region(i) == int(w.cloudRegion)) ||
			(w.tunnelRegion != 0 && w.region(i) == int(w.tunnelRegion)) ||
			(w.nowhereRegion != 0 && w.region(i) == int(w.nowhereRegion)) {
			continue
		}

		row := regionRow(i)
		col := regionCol(i)

		if row > int(w.maxMapRow) {
			w.maxMapRow = int32(row)
		}

		if col > int(w.maxMapCol) {
			w.maxMapCol = int32(col)
		}
	}

	w.maxMapInit = true
}

// losProvinceDistance is dir.c:128–169's directed graph search. The queue
// lives in the province records: even a same-province query resets/allocates
// every province, and finding the target leaves its frontier unexpanded.
func (w *world) losProvinceDistance(a, b int) int32 {
	a, b = w.province(a), w.province(b)
	for here := range w.loopProvince() {
		w.bx[here].temp = -1
		w.pLoc(here).next = 0
	}

	here := a
	w.bx[here].temp = 0
	end := here
	for {
		if here == b {
			return w.bx[here].temp
		}
		p := w.rpLoc(here)
		for _, dest := range p.provDest {
			if dest != 0 && w.bx[dest].temp == -1 {
				w.rpLoc(end).next = int32(dest)
				end = dest
				w.bx[end].temp = w.bx[here].temp + 1
			}
		}
		here = int(w.rpLoc(here).next)
		if here == 0 {
			break
		}
	}
	return -1
}

// exitView is dir.h's struct exit_view, in declaration order and C int widths.
type exitView struct {
	direction, destination, orig, distance                  int32
	impassable, destHidden, origHidden, hidden              int32
	inside, road, water, inTransit, magicBarrier, hadesCost int32
}

// LAND and WATER from dir.h (RAND is the nonzero random argument).
const (
	exitLand  = 1
	exitWater = 2
)

// These are the oly.h macros used by exit construction. In particular,
// is_ship_either tests subkind, and ship_gone is arithmetic, not a moving bit.
func (w *world) isShipEither(n int) bool { return w.isShip(n) || w.isShipNotdone(n) }
func (w *world) isLocOrShip(n int) bool  { return w.kind(n) == tLoc || w.kind(n) == tShip }
func (w *world) shipGone(n int) int32 {
	if p := w.rpSubloc(n); p != nil && p.moving != 0 {
		return w.sysclock.daysSinceEpoch - p.moving + w.evening
	}
	return 0
}
func (w *world) locSeaLane(n int) int8 {
	if p := w.rpLoc(n); p != nil {
		return p.seaLane
	}
	return 0
}
func (w *world) locBarrier(n int) int16 {
	if p := w.rpLoc(n); p != nil {
		return p.barrier
	}
	return 0
}
func (w *world) locLinkOpen(n int) int8 {
	if p := w.rpSubloc(n); p != nil {
		return p.linkOpen
	}
	return 0
}
func (w *world) uldim(n int) int8 {
	if p := w.rpSubloc(n); p != nil {
		return p.uldimFlag
	}
	return 0
}
func (w *world) summerbridge(n int) int8 {
	if p := w.rpSubloc(n); p != nil {
		return p.summerFlag
	}
	return 0
}
func (w *world) roadDest(n int) int {
	if p := w.rpGate(n); p != nil {
		return int(p.toLoc)
	}
	return 0
}
func (w *world) roadHidden(n int) int8 {
	if p := w.rpGate(n); p != nil {
		return p.roadHidden
	}
	return 0
}

func (w *world) locationDirection(where, dir int) int {
	dir--
	p := w.rpLoc(where)
	if p == nil || dir >= len(p.provDest) {
		return 0
	}
	return p.provDest[dir]
}

func (w *world) exitDistance(loc1, loc2 int) int32 {
	if w.subkind(loc1) == subHadesPit || w.subkind(loc2) == subHadesPit {
		return 28
	}
	if w.locDepth(loc1) > w.locDepth(loc2) {
		loc1, loc2 = loc2, loc1
	}
	wd, dd := w.locDepth(loc1), w.locDepth(loc2)
	if dd == locBuild {
		return 0
	}
	if dd == locSubloc {
		return 1
	}
	if w.subkind(loc1) == subOcean && w.subkind(loc2) != subOcean {
		return 2
	}
	if w.subkind(loc1) != subOcean && w.subkind(loc2) == subOcean {
		return 2
	}
	if w.province(loc1) != w.province(loc2) {
		loc1, loc2 = w.province(loc1), w.province(loc2)
	}
	switch w.subkind(loc2) {
	case subOcean:
		if w.locSeaLane(loc1) != 0 && w.locSeaLane(loc2) != 0 {
			return 2
		}
		return 3
	case subMountain:
		return 10
	case subForest, subDesert:
		return 8
	case subSwamp:
		return 14
	case subPlain, subUnder, subCloud:
		return 7
	case subTunnel, subChamber:
		return 5
	default:
		fmt.Fprintf(w.stderr, "exit_distance: subkind=%s, loc1=%d, loc2=%d, w_d=%d, d_d=%d\n", subkindS[w.subkind(loc2)], loc1, loc2, wd, dd)
		assert(false, "olympia/dir.c", 289, "FALSE")
	}
	return 0
}

func (w *world) isPortCity(where int) bool {
	if w.subkind(where) != subCity {
		return false
	}
	assert(w.locDepth(where) == locSubloc, "olympia/dir.c", 305, "loc_depth(where) == LOC_subloc")
	p := w.province(where)
	if w.subkind(p) == subMountain {
		return false
	}
	n, s := w.locationDirection(p, dirN), w.locationDirection(p, dirS)
	e, west := w.locationDirection(p, dirE), w.locationDirection(p, dirW)
	return (n != 0 && w.subkind(n) == subOcean) || (s != 0 && w.subkind(s) == subOcean) ||
		(e != 0 && w.subkind(e) == subOcean) || (west != 0 && w.subkind(west) == subOcean)
}

func (w *world) provinceHasPortCity(where int) int {
	assert(w.locDepth(where) == locProvince, "olympia/dir.c", 333, "loc_depth(where) == LOC_province")
	assert(w.validBox(where), "olympia/dir.c", 335, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.subkind(i) == subCity && w.isPortCity(i) {
			return i
		}
	}
	return 0
}

func (w *world) summerUldimOpenNow() bool {
	return (olyMonth(w.sysclock) >= 3 && olyMonth(w.sysclock) <= 6) || (olyMonth(w.sysclock) == 2 && w.monthDone != 0)
}

func (w *world) addProvinceExit(who, where, dest, dir int, l *[]*exitView) {
	assert(w.validBox(dest), "olympia/dir.c", 370, "valid_box(dest)")
	v := &exitView{}
	if (w.isShipEither(where) && w.shipGone(where) != 0) || (w.isShipEither(dest) && w.shipGone(dest) != 0) {
		v.inTransit = 1
	}
	if w.isShipEither(where) && w.subkind(dest) == subOcean {
		v.impassable = 1
	}
	if w.subkind(where) == subOcean && !w.isShipEither(dest) {
		v.water = 1
	}
	if w.subkind(dest) == subOcean {
		v.water = 1
	}
	if w.locDepth(where) == locProvince && w.subkind(dest) == subOcean && w.provinceHasPortCity(where) != 0 {
		v.impassable = 1
	}
	if w.subkind(dest) == subMineCollapsed {
		v.impassable = 1
	}
	if w.subkind(where) == subOcean && w.subkind(dest) != subOcean && w.subkind(dest) != subMountain && w.locDepth(dest) == locProvince {
		if n := w.provinceHasPortCity(dest); n != 0 {
			v.impassable = 1
			w.addProvinceExit(who, where, n, dir, l)
		}
	}
	if (w.subkind(where) == subMountain && w.subkind(dest) == subOcean) || (w.subkind(where) == subOcean && w.subkind(dest) == subMountain) {
		v.impassable = 1
	}
	if (dir == dirUp || dir == dirDown) && (w.subkind(where) == subCloud || w.subkind(dest) == subCloud) {
		v.impassable = 1
	}
	if (dir == dirN && w.uldim(where) == 1) || (dir == dirS && w.uldim(where) == 2) {
		v.impassable = 1
	}
	if (dir == dirN && (w.uldim(where) == 4 || w.summerbridge(where) == 1)) || (dir == dirS && (w.uldim(where) == 3 || w.summerbridge(where) == 2)) {
		if !w.summerUldimOpenNow() {
			v.impassable = 1
		}
	}
	v.orig, v.destination, v.direction = int32(where), int32(dest), int32(dir)
	v.distance = w.exitDistance(where, dest)
	if w.locHidden(where) != 0 {
		v.origHidden = 1
	}
	if w.locHidden(dest) != 0 {
		v.destHidden = 1
	}
	if w.locHidden(dest) != 0 && !w.testKnown(who, dest) && dir != dirOut {
		v.hidden = 1
	}
	if w.region(where) != w.region(dest) {
		v.inside = int32(w.region(dest))
		if w.region(where) != int(w.hadesRegion) && w.region(dest) == int(w.hadesRegion) {
			v.hadesCost = 100
		}
	}
	if w.locBarrier(dest) != 0 && dir != dirOut {
		v.impassable, v.magicBarrier = 1, 1
	}
	*l = append(*l, v)
}

func (w *world) extraRoutes(who, where int, l *[]*exitView) {
	assert(w.validBox(where), "olympia/dir.c", 515, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.kind(i) != tRoad {
			continue
		}
		dest := w.roadDest(i)
		assert(w.validBox(dest), "olympia/dir.c", 520, "valid_box(dest)")
		v := &exitView{orig: int32(where), destination: int32(dest), distance: w.exitDistance(where, dest), road: int32(i)}
		if (w.subkind(where) == subMountain && w.subkind(dest) == subCloud) || (w.subkind(where) == subCloud && w.subkind(dest) == subMountain) {
			v.impassable = 1
		}
		if w.roadHidden(i) != 0 {
			v.origHidden, v.destHidden = 1, 1
		}
		if w.roadHidden(i) != 0 && !w.testKnown(who, i) {
			v.hidden = 1
		}
		if w.region(where) != w.region(dest) {
			v.inside = int32(w.region(dest))
		}
		if w.subkind(where) == subOcean || w.subkind(dest) == subOcean {
			v.water = 1
		}
		*l = append(*l, v)
	}
}

func (w *world) provinceExits(who, where int, l *[]*exitView) {
	for dir := 1; dir <= dirDown; dir++ {
		if n := w.locationDirection(where, dir); n != 0 {
			w.addProvinceExit(who, where, n, dir, l)
		}
	}
}

func (w *world) provinceSubExits(who, where int, l *[]*exitView) {
	assert(w.validBox(where), "olympia/dir.c", 589, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.isLocOrShip(i) {
			w.addProvinceExit(who, where, i, dirIn, l)
		}
	}
	if p := w.rpSubloc(where); p != nil {
		for _, i := range p.linkFrom {
			if w.locLinkOpen(i) != 0 {
				w.addProvinceExit(who, where, i, dirIn, l)
			}
		}
	}
}

func (w *world) sublocExits(who, where int, l *[]*exitView) {
	if w.isPortCity(where) {
		p := w.province(where)
		for dir := 1; dir <= 4; dir++ {
			if n := w.locationDirection(p, dir); n != 0 && w.subkind(n) == subOcean {
				w.addProvinceExit(who, where, n, dir, l)
			}
		}
	}
	assert(w.validBox(where), "olympia/dir.c", 633, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.isLocOrShip(i) {
			w.addProvinceExit(who, where, i, dirIn, l)
		}
	}
	w.addProvinceExit(who, where, w.loc(where), dirOut, l)
	if p := w.rpSubloc(where); p != nil && w.locLinkOpen(where) != 0 {
		for _, i := range p.linkTo {
			w.addProvinceExit(who, where, i, 0, l)
		}
	}
}

func (w *world) shipExits(who, ship int, l *[]*exitView) {
	assert(w.isShipEither(ship), "olympia/dir.c", 665, "is_ship_either(ship)")
	outer := w.loc(ship)
	w.addProvinceExit(who, ship, outer, dirOut, l)
	assert(w.validBox(outer), "olympia/dir.c", 687, "valid_box(outer_loc)")
	for i := range w.loopHere(outer) {
		if i != ship && w.isShipEither(i) {
			w.addProvinceExit(who, ship, i, 0, l)
		}
	}
}

func (w *world) buildingExits(who, where int, l *[]*exitView) {
	if w.isShipEither(where) {
		w.shipExits(who, where, l)
	} else {
		w.addProvinceExit(who, where, w.loc(where), dirOut, l)
	}
	assert(w.validBox(where), "olympia/dir.c", 706, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.isLocOrShip(i) {
			w.addProvinceExit(who, where, i, dirIn, l)
		}
	}
}

func (w *world) exitsFromLoc(who, where int) []*exitView {
	var l []*exitView
	switch w.locDepth(where) {
	case locProvince:
		w.provinceExits(who, where, &l)
		w.provinceSubExits(who, where, &l)
	case locSubloc:
		w.sublocExits(who, where, &l)
	case locBuild:
		w.provinceExits(who, where, &l)
		w.buildingExits(who, where, &l)
	default:
		fmt.Fprintf(w.stderr, "where=%d, depth=%d\n", where, w.locDepth(where))
		assert(false, "olympia/dir.c", 747, "FALSE")
	}
	w.extraRoutes(who, where, &l)
	return l
}

func (w *world) exitsFromLocNSEW(who, where int) []*exitView {
	if w.locDepth(where) != locProvince {
		return nil
	}
	var l []*exitView
	w.provinceExits(who, where, &l)
	return l
}

func (w *world) exitsFromLocNSEWSelect(who, where, land, random int) []*exitView {
	if w.locDepth(where) != locProvince {
		return nil
	}
	var ret []*exitView
	for _, v := range w.exitsFromLocNSEW(who, where) {
		if (land&exitLand != 0 && v.water == 0) || (land&exitWater != 0 && v.water != 0) {
			ret = append(ret, v)
		}
	}
	if random != 0 {
		// Original begin_exitdir: one sequential exdr stream per turn,
		// rooted at (master, turn), not per location, actor or day.
		w.streams.BeginExitdir()
		rng.Shuffle(ret, w.streams.Exitdir())
	}
	return ret
}

func (w *world) hasOceanAccess(where int) int {
	ret := 0
	for _, v := range w.exitsFromLoc(0, where) {
		if v.water != 0 {
			if v.impassable != 0 {
				if ret == 0 {
					ret = 1
				}
			} else {
				ret = 2
			}
		}
	}
	return ret
}

// Route reporting is dir.c:839–1072. These callers deliberately use the
// unshuffled exitsFromLoc, including in-transit exits (the C skips are #if 0).
func (w *world) listExitExtras(who int, v *exitView) {
	if v.magicBarrier != 0 {
		w.indent += 3
		w.commandOut(who, "A magical barrier prevents entry.")
		w.indent -= 3
	}
	if v.hadesCost != 0 {
		w.indent += 3
		w.wiout(who, 1, "\"Notice to mortals, from the Gatekeeper Spirit of Hades: 100 gold/head is removed from any stack taking this road.\"")
		w.indent -= 3
	}
}

func (w *world) listExitsSup(who, where int, v *exitView, first *string) {
	if v.hidden != 0 && w.immedSeeAll == 0 {
		return
	}
	if first != nil && *first != "" {
		w.commandOut(who, "%s", *first)
		w.indent += 3
		*first = ""
	}
	var parts []string
	if v.direction > 0 {
		parts = append(parts, fullDirS[v.direction])
	}
	if w.name(int(v.destination)) != "" && !w.isShipEither(int(v.destination)) {
		parts = append(parts, subkindS[w.subkind(int(v.destination))])
	}
	parts = append(parts, sout("to %s", w.boxName(int(v.destination))))
	if v.inside != 0 {
		if s := w.name(int(v.inside)); s != "" {
			parts = append(parts, s)
		}
	}
	if v.destHidden != 0 {
		parts = append(parts, "hidden")
	}
	if v.impassable != 0 {
		parts = append(parts, "impassable")
	} else {
		parts = append(parts, sout("%d~day%s", v.distance, addS(int(v.distance))))
	}
	w.commandOut(who, "%s", capitalize(strings.Join(parts, ", ")))
	w.listExitExtras(who, v)
}

func (w *world) listRoadSup(who, where int, v *exitView, first *string) {
	if v.hidden != 0 && w.immedSeeAll == 0 {
		return
	}
	hid := ""
	if v.destHidden != 0 {
		hid = "hidden, "
	}
	if first != nil && *first != "" {
		w.commandOut(who, "")
		w.commandOut(who, "%s", *first)
		*first = ""
		w.indent += 3
	}
	dist := "impassable"
	if v.impassable == 0 && v.inTransit == 0 {
		dist = sout("%d~day%s", v.distance, addS(int(v.distance)))
	}
	w.commandOut(who, "%s, to %s, %s%s", w.justName(int(v.road)), w.boxName(int(v.destination)), hid, dist)
	w.listExitExtras(who, v)
}

func (w *world) listExits(who, where int) {
	l := w.exitsFromLoc(who, where)
	first := sout("Routes leaving %s: ", w.justName(where))
	for _, v := range l {
		if v.road == 0 && (v.direction != dirIn || w.immedSeeAll == 2) {
			w.listExitsSup(who, where, v, &first)
		}
	}
	for _, v := range l {
		if v.road != 0 {
			w.listRoadSup(who, where, v, &first)
		}
	}
	if first != "" {
		if w.isShipEither(where) {
			w.commandOut(who, "No current exits from %s", w.boxName(where))
		} else {
			w.commandOut(who, "No known routes leaving %s", w.boxName(where))
		}
	} else {
		w.indent -= 3
	}
	if (w.uldim(where) == 3 || w.uldim(where) == 4 || w.summerbridge(where) != 0) && !w.summerUldimOpenNow() {
		w.commandOut(who, "")
		switch {
		case w.uldim(where) == 3:
			w.commandOut(who, "Heavy snow blocks Uldim pass to the south.")
		case w.uldim(where) == 4:
			w.commandOut(who, "Heavy snow blocks Uldim pass to the north.")
		case w.summerbridge(where) == 1:
			w.commandOut(who, "Summerbridge to the north is impassable until the muds dry in the spring.")
		case w.summerbridge(where) == 2:
			w.commandOut(who, "Summerbridge to the south is impassable until the muds dry in the spring.")
		default:
			assert(false, "olympia/dir.c", 1029, "FALSE")
		}
	}
}

func (w *world) listSailableRoutes(who, ship int) {
	if !w.isShipEither(ship) {
		return
	}
	outer := w.loc(ship)
	l := w.exitsFromLoc(who, outer)
	first := "Ocean routes:"
	for _, v := range l {
		if v.direction > 0 && (v.direction != dirIn || w.immedSeeAll == 2) && v.water != 0 {
			w.listExitsSup(who, outer, v, &first)
		}
	}
	for _, v := range l {
		if v.road != 0 && v.water != 0 {
			w.listRoadSup(who, outer, v, &first)
		}
	}
	if first != "" {
		w.commandOut(who, "No visible sailable routes")
	} else {
		w.indent -= 3
	}
	w.commandOut(who, "")
}
