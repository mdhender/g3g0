// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

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

func (w *world) initWeatherViews() {
	for who := range w.loopKind(tChar) {
		if p := w.rpMagic(who); p != nil && p.knowsWeather != 0 {
			w.updateWeatherViewLocSup(who, w.province(who))
		}
	}
}
