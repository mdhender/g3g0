// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// newStorm is storm.c:279–323. Allocation, placement and garrison output
// deliberately use the same machinery as other minted entities.
func (w *world) newStorm(n, sk, aura, where int) int {
	assert(sk == subRain || sk == subWind || sk == subFog, "olympia/storm.c", 284, "sk == sub_rain || sk == sub_wind || sk == sub_fog")
	assert(w.locDepth(where) == locProvince, "olympia/storm.c", 285, "loc_depth(where) == LOC_province")
	before := w.weatherHere(where, sk)
	if n == 0 {
		n = w.newEnt(tStorm, sk)
		if n <= 0 {
			return -1
		}
	}
	w.pMisc(n).stormStr = int16(aura)
	w.setWhere(n, where)
	w.showToGarrison = true
	if before == 0 {
		switch sk {
		case subRain:
			w.commandOut(where, "It has begun to rain.")
		case subWind:
			w.commandOut(where, "It has become quite windy.")
		case subFog:
			w.commandOut(where, "It has become quite foggy.")
		}
	}
	w.showToGarrison = false
	return 0
}

// createSomeStorms is storm.c:1525–1554. Even num == 0 shuffles all
// candidates; special-region provinces still count toward nprovinces().
func (w *world) createSomeStorms(num, sk int) {
	var l []int
	for i := range w.loopProvince() {
		if w.greaterRegion(i) != 0 || w.weatherHere(i, sk) != 0 {
			continue
		}
		l = append(l, i)
	}
	w.streams.WthrShuffle(l)
	for i := 0; i < len(l) && i < num; i++ {
		w.newStorm(0, sk, w.streams.WthrStorm(l[i], int(w.sysclock.day)*256+sk, 2, 3), l[i])
	}
}

// naturalWeather is storm.c:1558–1613, sharing #74's persistent wthr stream.
// Movement/decay (#76) and unrelated weather magic remain explicit boundaries.
func (w *world) naturalWeather() {
	n := w.nprovinces() / 4 / 2 / 4
	switch olyMonth(w.sysclock) {
	case 0:
		w.createSomeStorms(n, subFog)
		w.createSomeStorms(n, subWind)
	case 1:
		w.createSomeStorms(n, subFog)
		w.createSomeStorms(n, subRain)
	case 2, 5:
	case 3:
		w.createSomeStorms(n, subRain)
	case 4:
		w.createSomeStorms(n*2, subRain)
	case 6:
		w.createSomeStorms(n, subRain)
		w.createSomeStorms(n, subFog)
		w.createSomeStorms(n, subRain)
	case 7:
		w.createSomeStorms(n, subWind)
	default:
		assert(false, "olympia/storm.c", 1611, "FALSE")
	}
}

func (w *world) vBindStorm(c *command) int {
	storm, who := int(c.a), int(c.who)
	ship := w.subloc(who)
	if w.kind(storm) != tStorm || w.npcSummoner(storm) != who {
		w.commandOut(who, "%s doesn't control any storm %s.", w.boxName(who), w.boxCode(storm))
		return 0
	}
	if !w.isShip(ship) {
		w.commandOut(who, "%s must be on a ship to bind the storm to.", w.boxName(who))
		return 0
	}
	if w.province(storm) != w.province(ship) {
		w.commandOut(who, "Storm must be in the same province as the ship it is to be bound to.")
		return 0
	}
	if !w.checkAura(who, 3) {
		return 0
	}
	return 1
}

func (w *world) dBindStorm(c *command) int {
	storm, who := int(c.a), int(c.who)
	ship := w.subloc(who)
	if w.kind(storm) != tStorm || w.npcSummoner(storm) != who {
		w.commandOut(who, "%s doesn't control storm %s anymore.", w.boxName(who), w.boxCode(storm))
		return 0
	}
	if !w.isShip(ship) {
		w.commandOut(who, "%s is no longer on a ship.", w.boxName(who))
		return 0
	}
	if w.province(storm) != w.province(ship) {
		w.commandOut(who, "Storm is no longer in the same province as the ship it is to be bound to.")
		return 0
	}
	if !w.chargeAura(who, 3) {
		return 0
	}
	if old := int(w.pMisc(storm).bindStorm); old != 0 {
		if p := w.rpSubloc(old); p != nil {
			p.boundStorms.RemValue(storm)
		}
	}
	w.pMisc(storm).bindStorm = int32(ship)
	w.pSubloc(ship).boundStorms.Append(storm)
	w.commandOut(who, "Bound %s to %s.", w.boxName(storm), w.boxName(ship))
	return 1
}

func (w *world) npcSummoner(n int) int {
	if p := w.rpMisc(n); p != nil {
		return int(p.summonedBy)
	}
	return 0
}

// weatherHere is weather_here(): buildings suppress weather; elsewhere
// sum matching storms directly in the province, in here-list order.
func (w *world) weatherHere(where, sk int) int32 {
	if w.locDepth(where) == locBuild {
		return 0
	}
	where = w.province(where)
	var sum int32
	assert(w.validBox(where), "olympia/storm.c", 439, "valid_box(where)")
	for i := range w.loopHere(where) {
		if w.kind(i) == tStorm && w.subkind(i) == sk {
			sum += w.stormStrength(i)
		}
	}
	return sum
}

func (w *world) updateWeatherViewLocSup(who, where int) {
	pl := w.player(who)
	assert(w.validBox(pl), "olympia/storm.c", 1622, "valid_box(pl)")
	setBit(&w.pPlayer(pl).weatherSeen, where)
}

// updateWeatherViewLocs is storm.c:1628–1644, including container roots.
func (w *world) updateWeatherViewLocs(stack, where int) {
	where = w.province(where)
	mage := func(who int) bool {
		p := w.rpMagic(who)
		return p != nil && p.knowsWeather != 0
	}
	if w.kind(stack) == tChar && mage(stack) {
		w.updateWeatherViewLocSup(stack, where)
	}
	for _, i := range w.allCharHere(stack) {
		if !w.isPrisoner(i) && mage(i) {
			w.updateWeatherViewLocSup(i, where)
		}
	}
}

func (w *world) initWeatherViews() {
	for who := range w.loopKind(tChar) {
		if p := w.rpMagic(who); p != nil && p.knowsWeather != 0 {
			w.updateWeatherViewLocSup(who, w.province(who))
		}
	}
}
