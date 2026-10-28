// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Explicit unfinished world phases in input.c/day.c order. These are entry
// boundaries for #66–#89, NOT implementations, even on an empty world.
func (w *world) autoDrop()            { panic(notImplemented("auto_drop")) }
func (w *world) linkDecay()           { panic(notImplemented("link_decay")) }
func (w *world) questDecay()          { panic(notImplemented("quest_decay")) }
func (w *world) determineNobleRanks() { panic(notImplemented("determine_noble_ranks")) }
