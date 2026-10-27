// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"io"
	"os"

	"github.com/mdhender/g3g0/internal/lists"
	"github.com/mdhender/g3g0/internal/rnd"
)

// world holds the C engine's global state: the box table bx[], the heads of
// the kind and subkind chains, and the globals that loading and saving read
// or write. In C they are file-scope globals; here each world is
// independent, so tests can build as many as they like.
//
// Every C int is an int32 here, as in the entity structures; box numbers
// passed to the accessors are Go ints.
type world struct {
	// glob.c

	bx      []*box        // all possible entities; len maxBoxes
	boxHead [tMax]int32   // heads of x_next_kind chain
	subHead [subMax]int32 // heads of x_next_sub chain

	fromHost       string
	replyHost      string
	gmAddress      string
	gameTitle      string
	gameURL        string
	rulesURL       string
	timesURL       string
	htpasswdLoc    string
	garrisonPay    int32
	armySlowFactor int32
	autoQuitTurns  int32

	sysclock       olytime
	indepPlayer    int32 // independent player
	gmPlayer       int32 // The Fates
	skillPlayer    int32 // skill listing
	eatPl          int32 // Order scanner
	npcPl          int32 // Subloc monster player
	garrPl         int32 // Garrison unit owner
	combatPl       int32 // Combat log
	postHasBeenRun int32 // FALSE
	garrisonMagic  int32
	mountOlympus   int32

	// libdir is libdir, the data directory (-l), which load_box()'s
	// diagnostic names. C starts it "lib".
	libdir string

	// Saved in libdir/system, but defined outside glob.c: seed_has_been_run
	// in io.c, faery_region and faery_player in faery.c, tunnel_region and
	// under_region in tunnel.c, hades_region, hades_pit and hades_player in
	// hades.c, nowhere_region and nowhere_loc in quest.c, and cloud_region in
	// cloud.c. All start at zero (FALSE).

	seedHasBeenRun int32 // FALSE
	faeryRegion    int32
	faeryPlayer    int32
	tunnelRegion   int32
	underRegion    int32
	hadesRegion    int32
	hadesPit       int32 // Pit of Hades
	hadesPlayer    int32
	nowhereRegion  int32
	nowhereLoc     int32
	cloudRegion    int32

	// maxMapRow, maxMapCol and maxMapInit are max_map_row, max_map_col and
	// max_map_init in dir.c, which determine_map_edges() sets.
	maxMapRow  int32
	maxMapCol  int32
	maxMapInit bool

	// The per-game trade-route buyer secret of buy.c (G3 issue #46):
	// trade_route_have_seed, trade_route_seed and trade_route_secret, its
	// statics, which init_trade_routes() sets and save_trade_routes()
	// writes, and trade_route_seed_pending and trade_route_seed_value,
	// which main() sets from -G before load_db().
	tradeRouteHaveSeed    bool
	tradeRouteSeed        uint64
	tradeRouteSecret      uint32 // tradeRouteDefaultSecret
	tradeRouteSeedPending bool
	tradeRouteSeedValue   uint64

	// newPlayers is new_players in add.c: the players added this turn,
	// whose Join-g3 files rename_act_join_files() renames.
	newPlayers lists.IList

	// oceanChars is the static ocean_chars of move.c: the characters at
	// sea, which init_ocean_chars() fills on load.
	oceanChars lists.IList

	// tradesToCheck is trades_to_check in buy.c: the boxes whose pending
	// trades investigate_possible_trade() wants matched.
	tradesToCheck lists.IList

	// The static cache of add_next_chain() and add_sub_chain() in code.c.

	nextCacheLast int32
	nextCacheKind int32
	subCacheLast  int32
	subCacheKind  int32 // -1

	// prettyDataFiles is pretty_data_files in main.c, which if_name() in
	// io.c reads. C starts it FALSE and -p toggles it.
	prettyDataFiles bool

	// lr is the state of the readfile()/readlin() line reader in z.c.
	lr lineReader

	// line and extBoxnum are the statics of the same names in io.c: the
	// line the scanners are looking at (nil for NULL, at end of file),
	// which they move past its first byte, and the box load_box() is
	// reading, for box_scan()'s diagnostic.
	line      *string
	extBoxnum int32

	// fuzzyFind is fuzzy_find in input.c: whether the last find_command()
	// matched a command name only fuzzily.
	fuzzyFind bool

	// loadQ and runQ are the statics load_q and run_q of input.c: the units
	// whose commands are loaded, by priority, and running.
	loadQ [maxPri]lists.IList
	runQ  lists.IList

	// input.c statics. cmdWait stays -1 until process_orders resolves it;
	// the wait command and full turn orchestration remain unported.
	curPri         int32
	cmdWait        int32
	immediate      bool
	showDay        bool // glob.c show_day; default formatting for command output
	evening        int32
	autoAttackFlag bool
	waitList       lists.IList // wait/collect population and polling: #61
	secondWaitList lists.IList // combat.c second_wait_list
	flags          []*flagEnt  // c1.c signals, not saved

	// monthDone is month_done in input.c, set after the daily loop.
	monthDone int32

	// dotCount is print_dot()'s static dot_count in u.c.
	dotCount int32

	// stderr is C's stderr, which print_box_usage() and alloc_box() write
	// to. newWorld sets it to os.Stderr; the engine sets it to its own.
	stderr io.Writer

	// stdout is C's stdout, which load_seed() writes to. newWorld sets it
	// to os.Stdout; the engine sets it to its buffered cStdout.
	stdout io.Writer

	// rnd is the state of the MD5 generator of olympia/rnd.c (digest and
	// master_seed), which load_seed() reads from libdir/randseed and
	// save_seed() writes back.
	rnd rnd.Source

	// testUseConstReportDate is test_use_const_report_date in main.c,
	// which times_masthead() reads.
	testUseConstReportDate bool

	// timeSelf is time_self in main.c (-T), which stage() reads, and
	// stageOld and stageFirst are stage()'s statics old and first.
	timeSelf   bool
	stageOld   int64
	stageFirst int64
}

// newWorld returns a world as olympia-g3 has it after glob_init(): the
// static initializers of the globals, an empty bx[] of maxBoxes entries and
// empty chains.
//
// The C flags (post_has_been_run, seed_has_been_run) stay int32, not bool:
// load_system() sets them with atoi() and save_system() prints them with
// %d, so any value round-trips.
func newWorld() *world {
	return &world{
		bx:      make([]*box, maxBoxes),
		cmdWait: -1,

		fromHost:       "g3@olympia.v-labs.be",
		replyHost:      "g3@olympia.v-labs.be",
		gmAddress:      "c@dvdc.be",
		gameTitle:      "Olympia G3",
		gameURL:        "http://olympia.v-labs.be/g3",
		rulesURL:       "http://oly.v-labs.be/g3/node/11",
		timesURL:       "http://oly.v-labs.be/g3/node/409",
		htpasswdLoc:    "/u/oly/g2/lib/ht-passwords",
		garrisonPay:    2,
		armySlowFactor: 20,
		autoQuitTurns:  0,

		indepPlayer:   100,
		gmPlayer:      200,
		skillPlayer:   202,
		eatPl:         203,
		npcPl:         206,
		garrPl:        207,
		combatPl:      0,
		garrisonMagic: 999,
		mountOlympus:  0,

		libdir: "lib",

		tradeRouteSecret: tradeRouteDefaultSecret,

		subCacheKind: -1,

		stderr: os.Stderr,
		stdout: os.Stdout,
	}
}

// kind is kind(): the kind of box n, or tDeleted if n is out of range or
// has no box.
func (w *world) kind(n int) int {
	if n > 0 && n < maxBoxes && w.bx[n] != nil {
		return int(w.bx[n].kind)
	}
	return tDeleted
}

// subkind is subkind(): the subkind of box n, or 0 if it has no box. Like
// C, it does not check n against the table; C's out-of-range read is
// undefined, and here it panics.
func (w *world) subkind(n int) int {
	if w.bx[n] != nil {
		return int(w.bx[n].skind)
	}
	return 0
}

// validBox is valid_box().
func (w *world) validBox(n int) bool {
	return w.kind(n) != tDeleted
}

// kindFirst is kind_first(): the first box of kind k, or 0.
func (w *world) kindFirst(k int) int { return int(w.boxHead[k]) }

// kindNext is kind_next(): the box after n in the chain of its kind, or 0.
// Box n must exist.
func (w *world) kindNext(n int) int { return int(w.bx[n].xNextKind) }

// subFirst is sub_first(): the first box of subkind sk, or 0.
func (w *world) subFirst(sk int) int { return int(w.subHead[sk]) }

// subNext is sub_next(): the box after n in the chain of its subkind, or 0.
// Box n must exist.
func (w *world) subNext(n int) int { return int(w.bx[n].xNextSub) }

// loc is loc(): where box n is. Box n must exist.
func (w *world) loc(n int) int { return int(w.rpLocInfo(n).where) }

// The p_xxx accessors are the malloc-on-demand references of oly.h: each
// returns box n's substructure, allocating a zeroed one first if it has
// none (if_malloc). The rp_xxx accessors return the "raw" pointer, which may
// be nil, and never allocate. Box n must exist, as in C, where a NULL box
// is dereferenced; here it panics.

// pLocInfo is p_loc_info(). The loc_info is part of the box, so there is
// nothing to allocate.
func (w *world) pLocInfo(n int) *locInfo { return &w.bx[n].xLocInfo }

// pChar is p_char().
func (w *world) pChar(n int) *entityChar {
	b := w.bx[n]
	if b.xChar == nil {
		b.xChar = &entityChar{}
	}
	return b.xChar
}

// pLoc is p_loc().
func (w *world) pLoc(n int) *entityLoc {
	b := w.bx[n]
	if b.xLoc == nil {
		b.xLoc = &entityLoc{}
	}
	return b.xLoc
}

// pSubloc is p_subloc().
func (w *world) pSubloc(n int) *entitySubloc {
	b := w.bx[n]
	if b.xSubloc == nil {
		b.xSubloc = &entitySubloc{}
	}
	return b.xSubloc
}

// pItem is p_item().
func (w *world) pItem(n int) *entityItem {
	b := w.bx[n]
	if b.xItem == nil {
		b.xItem = &entityItem{}
	}
	return b.xItem
}

// pPlayer is p_player().
func (w *world) pPlayer(n int) *entityPlayer {
	b := w.bx[n]
	if b.xPlayer == nil {
		b.xPlayer = &entityPlayer{}
	}
	return b.xPlayer
}

// pSkill is p_skill().
func (w *world) pSkill(n int) *entitySkill {
	b := w.bx[n]
	if b.xSkill == nil {
		b.xSkill = &entitySkill{}
	}
	return b.xSkill
}

// pGate is p_gate().
func (w *world) pGate(n int) *entityGate {
	b := w.bx[n]
	if b.xGate == nil {
		b.xGate = &entityGate{}
	}
	return b.xGate
}

// pMisc is p_misc().
func (w *world) pMisc(n int) *entityMisc {
	b := w.bx[n]
	if b.xMisc == nil {
		b.xMisc = &entityMisc{}
	}
	return b.xMisc
}

// pDisp is p_disp().
func (w *world) pDisp(n int) *attEnt {
	b := w.bx[n]
	if b.xDisp == nil {
		b.xDisp = &attEnt{}
	}
	return b.xDisp
}

// pCommand is p_command().
func (w *world) pCommand(n int) *command {
	b := w.bx[n]
	if b.cmd == nil {
		b.cmd = &command{}
	}
	return b.cmd
}

// pMagic is p_magic(): it allocates the entity_char first, if need be.
func (w *world) pMagic(n int) *charMagic {
	p := w.pChar(n)
	if p.xCharMagic == nil {
		p.xCharMagic = &charMagic{}
	}
	return p.xCharMagic
}

// pItemMagic is p_item_magic(): it allocates the entity_item first, if need
// be.
func (w *world) pItemMagic(n int) *itemMagic {
	p := w.pItem(n)
	if p.xItemMagic == nil {
		p.xItemMagic = &itemMagic{}
	}
	return p.xItemMagic
}

// rpLocInfo is rp_loc_info(), the same as p_loc_info().
func (w *world) rpLocInfo(n int) *locInfo { return &w.bx[n].xLocInfo }

// rpChar is rp_char().
func (w *world) rpChar(n int) *entityChar { return w.bx[n].xChar }

// rpLoc is rp_loc().
func (w *world) rpLoc(n int) *entityLoc { return w.bx[n].xLoc }

// rpSubloc is rp_subloc().
func (w *world) rpSubloc(n int) *entitySubloc { return w.bx[n].xSubloc }

// rpItem is rp_item().
func (w *world) rpItem(n int) *entityItem { return w.bx[n].xItem }

// rpPlayer is rp_player().
func (w *world) rpPlayer(n int) *entityPlayer { return w.bx[n].xPlayer }

// rpSkill is rp_skill().
func (w *world) rpSkill(n int) *entitySkill { return w.bx[n].xSkill }

// rpGate is rp_gate().
func (w *world) rpGate(n int) *entityGate { return w.bx[n].xGate }

// rpMisc is rp_misc().
func (w *world) rpMisc(n int) *entityMisc { return w.bx[n].xMisc }

// rpDisp is rp_disp().
func (w *world) rpDisp(n int) *attEnt { return w.bx[n].xDisp }

// rpCommand is rp_command().
func (w *world) rpCommand(n int) *command { return w.bx[n].cmd }

// rpMagic is rp_magic(): nil if box n has no entity_char.
func (w *world) rpMagic(n int) *charMagic {
	if p := w.rpChar(n); p != nil {
		return p.xCharMagic
	}
	return nil
}

// rpItemMagic is rp_item_magic(): nil if box n has no entity_item.
func (w *world) rpItemMagic(n int) *itemMagic {
	if p := w.rpItem(n); p != nil {
		return p.xItemMagic
	}
	return nil
}

// The attribute macros of oly.h that check_db() uses. Each reads a
// substructure through its rp_xxx accessor and is 0 (FALSE) if the box has
// none. Box n must exist.

// nobleItem is noble_item(): the item unit n is made of, or 0.
func (w *world) nobleItem(n int) int {
	if p := w.rpChar(n); p != nil {
		return int(p.unitItem)
	}
	return 0
}

// charMoving is char_moving(): the day unit n began to move, or 0.
func (w *world) charMoving(n int) int {
	if p := w.rpChar(n); p != nil {
		return int(p.moving)
	}
	return 0
}

// isPrisoner is is_prisoner(): whether unit n is a prisoner.
func (w *world) isPrisoner(n int) bool {
	if p := w.rpChar(n); p != nil {
		return p.prisoner != 0
	}
	return false
}

// itemUnique is item_unique(): who has unique item n, or 0 if it is not
// unique.
func (w *world) itemUnique(n int) int {
	if p := w.rpItem(n); p != nil {
		return int(p.whoHas)
	}
	return 0
}

// learnTime is learn_time(): the days of study skill n takes, or 0.
func (w *world) learnTime(n int) int {
	if p := w.rpSkill(n); p != nil {
		return int(p.timeToLearn)
	}
	return 0
}

// reqSkill is req_skill(): the skill n requires (its parent), or 0.
func (w *world) reqSkill(n int) int {
	if p := w.rpSkill(n); p != nil {
		return int(p.requiredSkill)
	}
	return 0
}
