// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"os"
	"sort"
)

// coverage.c / oly.h: row order is execution order, not alphabetical order.
const (
	phSavageSpawn = iota
	phHadesSpawn
	phHadesUnit
	phNpcUnsworn
	phNpcUndead
	phNpcSavage
	phNpcMob
	phNpcBandit
	phNpcToken
	phGarrGuard
	phTokenMelt
	phTokenAdd
	phStopInterrupt
	phStopIdle
	phTradeMatch
	phGarrPay
	phCoastDamage
	phCoastOmen
	phLocMine
	phLocInn
	phOceanPlunge
	phHealHeal
	phHealSicken
	phHealRecover
	phPrisonerEscape
	phCurseErode
	phRingAura
	phDogsBark
	phWeatherStorm
	phFaerySpawn
	phFaeryUnit
	phSpecialOpen
	phSpecialClose
	phCityGold
	phGarrTax
	phGarrCastle
	phGarrOwner
	phTaxes
	phClaim
	phNp
	phUnformed
	phAura
	phAbilityShroud
	phRegionShroud
	phMedHinder
	phBarrierDec
	phBarrierGone
	phLoyDecay
	phLoyRelease
	phLoyDesert
	phPillage
	phRelicDecay
	phRelicVanish
	phHideMage
	phInnIncome
	phInnTraveller
	phTemple
	phMaintPaid
	phMaintStarve
	phAnimal
	phGhost
	phCorpse
	phBody
	phStormWeaken
	phStormGone
	phStormMove
	phMineDec
	phMineClear
	phProduction
	phAutoDrop
	phLinkDec
	phLinkOpen
	phQuest
	phRanks
	phMax
)

var phRows = [phMax]struct{ phase, routine, event string }{
	phSavageSpawn:    {"day0", "init_savage_attacks", "spawn"},
	phHadesSpawn:     {"day0", "auto_hades", "spawn"},
	phHadesUnit:      {"day0", "auto_hades", "unit"},
	phNpcUnsworn:     {"day0", "queue_npc_orders", "unsworn"},
	phNpcUndead:      {"day0", "queue_npc_orders", "undead"},
	phNpcSavage:      {"day0", "queue_npc_orders", "savage"},
	phNpcMob:         {"day0", "queue_npc_orders", "mob"},
	phNpcBandit:      {"day0", "queue_npc_orders", "bandit"},
	phNpcToken:       {"day0", "queue_npc_orders", "token"},
	phGarrGuard:      {"day0", "ping_garrisons", "guard"},
	phTokenMelt:      {"day0", "check_token_units", "melt"},
	phTokenAdd:       {"day0", "check_token_units", "add"},
	phStopInterrupt:  {"day0", "process_interrupted_units", "interrupt"},
	phStopIdle:       {"day0", "process_interrupted_units", "idle"},
	phTradeMatch:     {"day1", "match_all_trades", "trade"},
	phGarrPay:        {"daily", "default_garrison_pay", "pay"},
	phCoastDamage:    {"daily", "ship_coastal_damage", "damage"},
	phCoastOmen:      {"daily", "ship_coastal_damage", "omen"},
	phLocMine:        {"daily", "random_loc_damage", "mine"},
	phLocInn:         {"daily", "random_loc_damage", "inn"},
	phOceanPlunge:    {"daily", "check_ocean_chars", "plunge"},
	phHealHeal:       {"daily", "heal_characters", "heal"},
	phHealSicken:     {"daily", "heal_characters", "sicken"},
	phHealRecover:    {"daily", "heal_characters", "recover"},
	phPrisonerEscape: {"daily", "weekly_prisoner_escape_check", "escape"},
	phCurseErode:     {"daily", "noncreator_curse_erode", "erode"},
	phRingAura:       {"daily", "increment_stone_ring_aura", "aura"},
	phDogsBark:       {"daily", "dogs_bark_at_hidden_chars", "bark"},
	phWeatherStorm:   {"daily", "natural_weather", "storm"},
	phFaerySpawn:     {"daily", "auto_faery", "spawn"},
	phFaeryUnit:      {"daily", "auto_faery", "unit"},
	phSpecialOpen:    {"month", "special_locs_open", "open"},
	phSpecialClose:   {"month", "special_locs_close", "close"},
	phCityGold:       {"month", "move_city_gold", "move"},
	phGarrTax:        {"month", "garrison_gold", "tax"},
	phGarrCastle:     {"month", "garrison_gold", "castle"},
	phGarrOwner:      {"month", "garrison_gold", "owner"},
	phTaxes:          {"month", "collect_taxes", "collect"},
	phClaim:          {"month", "add_claim_gold", "grant"},
	phNp:             {"month", "add_noble_points", "point"},
	phUnformed:       {"month", "add_unformed", "mint"},
	phAura:           {"month", "increment_current_aura", "aura"},
	phAbilityShroud:  {"month", "decrement_ability_shroud", "decrement"},
	phRegionShroud:   {"month", "decrement_region_shroud", "decrement"},
	phMedHinder:      {"month", "decrement_meditation_hinder", "decrement"},
	phBarrierDec:     {"month", "decrement_loc_barrier", "decrement"},
	phBarrierGone:    {"month", "decrement_loc_barrier", "dissipate"},
	phLoyDecay:       {"month", "loyalty_decay", "decay"},
	phLoyRelease:     {"month", "loyalty_decay", "release"},
	phLoyDesert:      {"month", "loyalty_decay", "desert"},
	phPillage:        {"month", "pillage_decay", "decay"},
	phRelicDecay:     {"month", "relic_decay", "decay"},
	phRelicVanish:    {"month", "relic_decay", "vanish"},
	phHideMage:       {"month", "hide_mage_decay", "decrement"},
	phInnIncome:      {"month", "inn_income", "income"},
	phInnTraveller:   {"month", "inn_income", "traveller"},
	phTemple:         {"month", "temple_income", "offering"},
	phMaintPaid:      {"month", "charge_maint_costs", "paid"},
	phMaintStarve:    {"month", "charge_maint_costs", "starve"},
	phAnimal:         {"month", "animal_deaths", "death"},
	phGhost:          {"month", "ghost_warrior_decay", "evaporate"},
	phCorpse:         {"month", "corpse_decay", "decompose"},
	phBody:           {"month", "dead_body_rot", "rot"},
	phStormWeaken:    {"month", "storm_decay", "weaken"},
	phStormGone:      {"month", "storm_decay", "dissipate"},
	phStormMove:      {"month", "storm_move", "move"},
	phMineDec:        {"month", "collapsed_mine_decay", "decrement"},
	phMineClear:      {"month", "collapsed_mine_decay", "clear"},
	phProduction:     {"month", "post_production", "run"},
	phAutoDrop:       {"month", "auto_drop", "drop"},
	phLinkDec:        {"month", "link_decay", "decrement"},
	phLinkOpen:       {"month", "link_decay", "open"},
	phQuest:          {"month", "quest_decay", "decrement"},
	phRanks:          {"month", "determine_noble_ranks", "rank"},
}

func (w *world) phaseRan(row int) {
	if w.phaseCoveragePath == nil || w.immediate {
		return
	}
	assert(row >= 0 && row < phMax, "olympia/coverage.c", 405, "row >= 0 && row < PH_MAX")
	w.phRan[row]++
}

func (w *world) phaseActed(row int) {
	if w.phaseCoveragePath == nil || w.immediate {
		return
	}
	assert(row >= 0 && row < phMax, "olympia/coverage.c", 417, "row >= 0 && row < PH_MAX")
	w.phActed[row]++
}

func (w *world) writePhaseCoverage() {
	if w.phaseCoveragePath == nil {
		return
	}
	fp, err := os.Create(*w.phaseCoveragePath)
	if err != nil {
		fmt.Fprintf(w.stderr, "Can't write %s!\n", *w.phaseCoveragePath)
		return
	}
	defer fp.Close()
	fmt.Fprint(fp, "phase\troutine\tevent\tran\tacted\n")
	first := 0
	for i, row := range phRows {
		// The fixed C table asserts every slot was initialized. All original
		// routine names are nonempty; an omitted Go initializer is its NULL.
		assert(row.routine != "", "olympia/coverage.c", 442, "ph_rows[i].routine != NULL")
		if row.routine != phRows[first].routine {
			first = i
		}
		fmt.Fprintf(fp, "%s\t%s\t%s\t%d\t%d\n", row.phase, row.routine, row.event, w.phRan[first], w.phActed[i])
	}
}

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
