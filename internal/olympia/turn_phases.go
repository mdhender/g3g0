// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Explicit unfinished world phases in input.c/day.c order. These are entry
// boundaries for #66–#89, NOT implementations, even on an empty world.
func (w *world) matchAllTrades() { panic(notImplemented("match_all_trades")) }

// #74/#75 own the actual caln/wthr draws. Do not draw from a substitute RNG.
func (w *world) calDay(pick, low, high int) int { panic(notImplemented("cal_day")) }
func (w *world) wthrShuffle(days []int)         { panic(notImplemented("wthr_shuffle")) }
func (w *world) shipCoastalDamage()             { panic(notImplemented("ship_coastal_damage")) }
func (w *world) randomLocDamage()               { panic(notImplemented("random_loc_damage")) }
func (w *world) checkOceanChars()               { panic(notImplemented("check_ocean_chars")) }
func (w *world) healCharacters()                { panic(notImplemented("heal_characters")) }
func (w *world) weeklyPrisonerEscapeCheck()     { panic(notImplemented("weekly_prisoner_escape_check")) }
func (w *world) noncreatorCurseErode()          { panic(notImplemented("noncreator_curse_erode")) }
func (w *world) incrementStoneRingAura()        { panic(notImplemented("increment_stone_ring_aura")) }
func (w *world) dogsBarkAtHiddenChars()         { panic(notImplemented("dogs_bark_at_hidden_chars")) }
func (w *world) naturalWeather()                { panic(notImplemented("natural_weather")) }
func (w *world) autoFaery()                     { panic(notImplemented("auto_faery")) }

func (w *world) specialLocsOpen()           { panic(notImplemented("special_locs_open")) }
func (w *world) specialLocsClose()          { panic(notImplemented("special_locs_close")) }
func (w *world) moveCityGold()              { panic(notImplemented("move_city_gold")) }
func (w *world) garrisonGold()              { panic(notImplemented("garrison_gold")) }
func (w *world) collectTaxes()              { panic(notImplemented("collect_taxes")) }
func (w *world) addClaimGold()              { panic(notImplemented("add_claim_gold")) }
func (w *world) addNoblePoints()            { panic(notImplemented("add_noble_points")) }
func (w *world) addUnformed()               { panic(notImplemented("add_unformed")) }
func (w *world) incrementCurrentAura()      { panic(notImplemented("increment_current_aura")) }
func (w *world) decrementAbilityShroud()    { panic(notImplemented("decrement_ability_shroud")) }
func (w *world) decrementRegionShroud()     { panic(notImplemented("decrement_region_shroud")) }
func (w *world) decrementMeditationHinder() { panic(notImplemented("decrement_meditation_hinder")) }
func (w *world) decrementLocBarrier()       { panic(notImplemented("decrement_loc_barrier")) }
func (w *world) loyaltyDecay()              { panic(notImplemented("loyalty_decay")) }
func (w *world) pillageDecay()              { panic(notImplemented("pillage_decay")) }
func (w *world) relicDecay()                { panic(notImplemented("relic_decay")) }
func (w *world) hideMageDecay()             { panic(notImplemented("hide_mage_decay")) }
func (w *world) innIncome()                 { panic(notImplemented("inn_income")) }
func (w *world) templeIncome()              { panic(notImplemented("temple_income")) }
func (w *world) chargeMaintCosts()          { panic(notImplemented("charge_maint_costs")) }
func (w *world) animalDeaths()              { panic(notImplemented("animal_deaths")) }
func (w *world) ghostWarriorDecay()         { panic(notImplemented("ghost_warrior_decay")) }
func (w *world) corpseDecay()               { panic(notImplemented("corpse_decay")) }
func (w *world) deadBodyRot()               { panic(notImplemented("dead_body_rot")) }
func (w *world) stormDecay()                { panic(notImplemented("storm_decay")) }
func (w *world) stormMove()                 { panic(notImplemented("storm_move")) }
func (w *world) collapsedMineDecay()        { panic(notImplemented("collapsed_mine_decay")) }
func (w *world) postProduction()            { panic(notImplemented("post_production")) }
func (w *world) autoDrop()                  { panic(notImplemented("auto_drop")) }
func (w *world) linkDecay()                 { panic(notImplemented("link_decay")) }
func (w *world) questDecay()                { panic(notImplemented("quest_decay")) }
func (w *world) determineNobleRanks()       { panic(notImplemented("determine_noble_ranks")) }
