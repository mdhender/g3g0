// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Explicit unfinished world phases in input.c/day.c order. These are entry
// boundaries for #66–#89, NOT implementations, even on an empty world.
func (w *world) garrisonGold()        { panic(notImplemented("garrison_gold")) }
func (w *world) chargeMaintCosts()    { panic(notImplemented("charge_maint_costs")) }
func (w *world) animalDeaths()        { panic(notImplemented("animal_deaths")) }
func (w *world) ghostWarriorDecay()   { panic(notImplemented("ghost_warrior_decay")) }
func (w *world) corpseDecay()         { panic(notImplemented("corpse_decay")) }
func (w *world) deadBodyRot()         { panic(notImplemented("dead_body_rot")) }
func (w *world) collapsedMineDecay()  { panic(notImplemented("collapsed_mine_decay")) }
func (w *world) autoDrop()            { panic(notImplemented("auto_drop")) }
func (w *world) linkDecay()           { panic(notImplemented("link_decay")) }
func (w *world) questDecay()          { panic(notImplemented("quest_decay")) }
func (w *world) determineNobleRanks() { panic(notImplemented("determine_noble_ranks")) }
