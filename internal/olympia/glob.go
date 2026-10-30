// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// The constant tables of olympia/glob.c that loading, saving and check_db()
// read. The scalar globals of glob.c are fields of world (world.go).
//
// The C tables end with a NULL entry, which C loops test for. The slices
// here have no terminator: an index past the end, which in C reads the
// NULL, panics, and loops run to len() instead. check_glob() in check.c
// asserts that kind_s and subkind_s have T_MAX and SUB_MAX entries before
// the NULL; with no terminator that is len(kindS) == tMax and
// len(subkindS) == subMax.
//
// libdir is options.libdir (olympia.go) and world.libdir; show_day is
// world.showDay. See cmdTbl for the command-handler boundaries.

// The build-depth entry is NULL in C: json_str renders it as empty. Retain
// the original missing closing bracket in the no-depth label as well.
var locDepthS = []string{"<no depth", "region", "province", "subloc", ""}

var shortDirS = []string{"<no dir>", "n", "e", "s", "w", "u", "d", "i", "o"}
var fullDirS = []string{"<no dir>", "north", "east", "south", "west", "up", "down", "in", "out"}

// exitOpposite is the original glob.c table, including its non-cardinal order.
var exitOpposite = []int{0, dirS, dirW, dirN, dirE, dirOut, dirIn, dirDown, dirUp, 0}

// kindS is kind_s, the name of each entity kind, indexed by T_xxx.
var kindS = []string{
	"deleted",  // T_deleted
	"player",   // T_player
	"char",     // T_char
	"loc",      // T_loc
	"item",     // T_item
	"skill",    // T_skill
	"gate",     // T_gate
	"road",     // T_road
	"deadchar", // T_deadchar
	"ship",     // T_ship
	"post",     // T_post
	"storm",    // T_storm
	"unform",   // T_unform
	"lore",     // T_lore
}

// monthNames is month_names, the name of each month of the Olympian
// year, which times_masthead() prints.
var monthNames = []string{
	"Fierce winds",
	"Snowmelt",
	"Blossom bloom",
	"Sunsear",
	"Thunder and rain",
	"Harvest",
	"Waning days",
	"Dark night",
}

// subkindS is subkind_s, the name of each subkind, indexed by sub_xxx.
var subkindS = []string{
	"<no subkind>",
	"ocean",                 // sub_ocean
	"forest",                // sub_forest
	"plain",                 // sub_plain
	"mountain",              // sub_mountain
	"desert",                // sub_desert
	"swamp",                 // sub_swamp
	"underground",           // sub_under
	"faery hill",            // sub_faery_hill
	"island",                // sub_island
	"ring of stones",        // sub_stone_cir
	"mallorn grove",         // sub_mallorn_grove
	"bog",                   // sub_bog
	"cave",                  // sub_cave
	"city",                  // sub_city
	"lair",                  // sub_lair
	"graveyard",             // sub_graveyard
	"ruins",                 // sub_ruins
	"battlefield",           // sub_battlefield
	"enchanted forest",      // sub_ench_forest
	"rocky hill",            // sub_rocky_hill
	"circle of trees",       // sub_tree_cir
	"pits",                  // sub_pits
	"pasture",               // sub_pasture
	"oasis",                 // sub_oasis
	"yew grove",             // sub_yew_grove
	"sand pit",              // sub_sand_pit
	"sacred grove",          // sub_sacred_grove
	"poppy field",           // sub_poppy_field
	"temple",                // sub_temple
	"galley",                // sub_galley
	"roundship",             // sub_roundship
	"castle",                // sub_castle
	"galley-in-progress",    // sub_galley_notdone
	"roundship-in-progress", // sub_roundship_notdone
	"ghost ship",            // sub_ghost_ship
	"temple-in-progress",    // sub_temple_notdone
	"inn",                   // sub_inn
	"inn-in-progress",       // sub_inn_notdone
	"castle-in-progress",    // sub_castle_notdone
	"mine",                  // sub_mine
	"mine-in-progress",      // sub_mine_notdone
	"scroll",                // sub_scroll
	"magic",                 // sub_magic
	"palantir",              // sub_palantir
	"auraculum",             // sub_auraculum
	"tower",                 // sub_tower
	"tower-in-progress",     // sub_tower_notdone
	"pl_system",             // sub_pl_system
	"pl_regular",            // sub_pl_regular
	"region",                // sub_region
	"pl_savage",             // sub_pl_savage
	"pl_npc",                // sub_pl_npc
	"collapsed mine",        // sub_mine_collapsed
	"ni",                    // sub_ni
	"demon lord",            // sub_undead
	"dead body",             // sub_dead_body
	"fog",                   // sub_fog
	"wind",                  // sub_wind
	"rain",                  // sub_rain
	"pit",                   // sub_hades_pit
	"artifact",              // sub_artifact
	"pl_silent",             // sub_pl_silent
	"npc_token",             // sub_npc_token
	"garrison",              // sub_garrison
	"cloud",                 // sub_cloud
	"raft",                  // sub_raft
	"raft-in-progress",      // sub_raft_notdone
	"suffuse_ring",          // sub_suffuse_ring
	"relic",                 // sub_relic
	"tunnel",                // sub_tunnel
	"sewer",                 // sub_sewer
	"chamber",               // sub_chamber
	"tradegood",             // sub_tradegood
}

// cmdTbl is cmd_tbl, the command table. A command is its index here.
//
// init installs the handlers below, preserving NULL versus an unported
// routine. The C routine names are in each row's comment, in that order,
// with - for NULL. Never reorder rows: command indexes are saved in the DB.
//
// Allow field:
//
//	c	character
//	p	player entity
//	i	immediate mode only (debugging/maintenance)
//	r	restricted -- for npc units under control
//	g	garrison
//	m	Gamemaster only
var cmdTbl = []cmdTblEnt{
	// allow, name, start, finish, interrupt, time, poll, pri
	{"", "", nil, nil, nil, 0, 0, 3},            // - - -
	{"cpr", "accept", nil, nil, nil, 0, 0, 0},   // v_accept - -
	{"cpr", "admit", nil, nil, nil, 0, 0, 0},    // v_admit - -
	{"cr", "attack", nil, nil, nil, 1, 0, 3},    // v_attack - -
	{"cr", "banner", nil, nil, nil, 0, 0, 1},    // v_banner - -
	{"cr", "behind", nil, nil, nil, 0, 0, 1},    // v_behind - -
	{"c", "bind", nil, nil, nil, 7, 0, 3},       // v_bind_storm d_bind_storm -
	{"c", "board", nil, nil, nil, 0, 0, 2},      // v_board - -
	{"c", "breed", nil, nil, nil, 7, 0, 3},      // v_breed d_breed -
	{"c", "bribe", nil, nil, nil, 7, 0, 3},      // v_bribe d_bribe -
	{"c", "build", nil, nil, nil, -1, 1, 3},     // v_build d_build -
	{"c", "buy", nil, nil, nil, 0, 0, 1},        // v_buy - -
	{"c", "catch", nil, nil, nil, -1, 1, 3},     // v_catch - -
	{"c", "claim", nil, nil, nil, 0, 0, 1},      // v_claim - -
	{"c", "collect", nil, nil, nil, -1, 1, 3},   // v_collect d_collect i_collect
	{"cr", "contact", nil, nil, nil, 0, 0, 0},   // v_contact - -
	{"m", "credit", nil, nil, nil, 0, 0, 0},     // v_credit - -
	{"c", "decree", nil, nil, nil, 0, 0, 0},     // v_decree - -
	{"cpr", "default", nil, nil, nil, 0, 0, 0},  // v_att_clear - -
	{"cpr", "defend", nil, nil, nil, 0, 0, 0},   // v_defend - -
	{"c", "die", nil, nil, nil, 0, 0, 1},        // v_die - -
	{"c", "discard", nil, nil, nil, 0, 0, 1},    // v_discard - -
	{"cr", "drop", nil, nil, nil, 0, 0, 1},      // v_discard - -
	{"m", "emote", nil, nil, nil, 0, 0, 1},      // v_emote - -
	{"cr", "execute", nil, nil, nil, 0, 0, 1},   // v_execute - -
	{"c", "explore", nil, nil, nil, 7, 0, 3},    // v_explore d_explore -
	{"c", "fee", nil, nil, nil, 0, 0, 1},        // v_fee - -
	{"c", "ferry", nil, nil, nil, 0, 0, 1},      // v_ferry - -
	{"c", "fish", nil, nil, nil, -1, 1, 3},      // v_fish - -
	{"cr", "flag", nil, nil, nil, 0, 0, 1},      // v_flag - -
	{"c", "fly", nil, nil, nil, -1, 0, 2},       // v_fly d_fly -
	{"c", "forget", nil, nil, nil, 0, 0, 1},     // v_forget - -
	{"c", "form", nil, nil, nil, 7, 0, 3},       // v_form d_form -
	{"cp", "format", nil, nil, nil, 0, 0, 1},    // v_format - -
	{"c", "garrison", nil, nil, nil, 1, 0, 3},   // v_garrison - -
	{"cr", "get", nil, nil, nil, 0, 0, 1},       // v_get - -
	{"cr", "give", nil, nil, nil, 0, 0, 1},      // v_give - -
	{"cr", "go", nil, nil, nil, -1, 0, 2},       // v_move d_move -
	{"c", "guard", nil, nil, nil, 0, 0, 1},      // v_guard - -
	{"c", "hide", nil, nil, nil, 3, 0, 3},       // v_hide d_hide -
	{"c", "honor", nil, nil, nil, 1, 0, 3},      // v_honor - -
	{"c", "honour", nil, nil, nil, 1, 0, 3},     // v_honor - -
	{"cpr", "hostile", nil, nil, nil, 0, 0, 0},  // v_hostile - -
	{"c", "improve", nil, nil, nil, -1, 1, 3},   // v_improve d_improve -
	{"c", "incite", nil, nil, nil, 7, 0, 3},     // v_incite - -
	{"c", "make", nil, nil, nil, -1, 1, 3},      // v_make d_make i_make
	{"c", "mallorn", nil, nil, nil, -1, 1, 3},   // v_mallorn - -
	{"cp", "message", nil, nil, nil, 1, 0, 3},   // v_message - -
	{"cr", "move", nil, nil, nil, -1, 0, 2},     // v_move d_move -
	{"cpr", "name", nil, nil, nil, 0, 0, 1},     // v_name - -
	{"cpr", "neutral", nil, nil, nil, 0, 0, 0},  // v_neutral - -
	{"cp", "notab", nil, nil, nil, 0, 0, 1},     // v_notab - -
	{"c", "oath", nil, nil, nil, 1, 0, 3},       // v_oath - -
	{"c", "opium", nil, nil, nil, -1, 1, 3},     // v_opium - -
	{"cr", "pay", nil, nil, nil, 0, 0, 1},       // v_pay - -
	{"cr", "pillage", nil, nil, nil, 7, 0, 3},   // v_pillage d_pillage -
	{"c", "pledge", nil, nil, nil, 0, 0, 1},     // v_pledge - -
	{"cr", "plugh", nil, nil, nil, 0, 0, 3},     // v_plugh - -
	{"c", "post", nil, nil, nil, 1, 0, 3},       // v_post - -
	{"cp", "press", nil, nil, nil, 0, 0, 1},     // v_press - -
	{"cr", "promote", nil, nil, nil, 0, 0, 1},   // v_promote - -
	{"cp", "public", nil, nil, nil, 0, 0, 1},    // v_public - -
	{"c", "quarry", nil, nil, nil, -1, 1, 3},    // v_quarry - -
	{"c", "quest", nil, nil, nil, 7, 0, 3},      // v_quest d_quest -
	{"p", "quit", nil, nil, nil, 0, 0, 1},       // v_quit - -
	{"c", "raise", nil, nil, nil, 7, 0, 3},      // v_raise d_raise -
	{"c", "rally", nil, nil, nil, 7, 0, 3},      // v_rally d_rally -
	{"cr", "raze", nil, nil, nil, -1, 1, 3},     // v_raze d_raze -
	{"cpr", "realname", nil, nil, nil, 0, 0, 1}, // v_fullname - -
	{"c", "reclaim", nil, nil, nil, 0, 0, 1},    // v_reclaim - -
	{"c", "recruit", nil, nil, nil, -1, 1, 3},   // v_recruit - -
	{"c", "repair", nil, nil, nil, -1, 1, 3},    // v_repair d_repair i_repair
	{"c", "research", nil, nil, nil, 7, 0, 3},   // v_research d_research -
	{"cp", "rumor", nil, nil, nil, 0, 0, 1},     // v_rumor - -
	{"c", "sail", nil, nil, nil, -1, 0, 4},      // v_sail d_sail i_sail
	{"c", "sell", nil, nil, nil, 0, 0, 1},       // v_sell - -
	{"cr", "seek", nil, nil, nil, 7, 1, 3},      // v_seek d_seek -
	{"c", "sneak", nil, nil, nil, 3, 0, 3},      // v_sneak d_sneak -
	{"cp", "split", nil, nil, nil, 0, 0, 1},     // v_split - -
	{"cr", "stack", nil, nil, nil, 0, 0, 1},     // v_stack - -
	{"c", "stone", nil, nil, nil, -1, 1, 3},     // v_quarry - -
	{"c", "study", nil, nil, nil, 7, 1, 3},      // v_study d_study -
	{"c", "surrender", nil, nil, nil, 1, 0, 1},  // v_surrender - -
	{"c", "swear", nil, nil, nil, 0, 0, 1},      // v_swear - -
	{"cr", "take", nil, nil, nil, 0, 0, 1},      // v_get - -
	{"cp", "times", nil, nil, nil, 0, 0, 1},     // v_times - -
	{"c", "train", nil, nil, nil, -1, 1, 3},     // v_make d_make i_make
	{"c", "trance", nil, nil, nil, 28, 0, 3},    // v_trance d_trance -
	{"cr", "terrorize", nil, nil, nil, 7, 0, 3}, // v_terrorize d_terrorize -
	{"c", "torture", nil, nil, nil, 7, 0, 3},    // v_torture d_torture -
	{"c", "unload", nil, nil, nil, 0, 0, 3},     // v_unload - -
	{"c", "ungarrison", nil, nil, nil, 1, 0, 3}, // v_ungarrison - -
	{"cr", "unstack", nil, nil, nil, 0, 0, 1},   // v_unstack - -
	{"c", "use", nil, nil, nil, -1, 1, 3},       // v_use d_use i_use
	{"crm", "wait", nil, nil, nil, -1, 1, 1},    // v_wait d_wait i_wait
	{"c", "wood", nil, nil, nil, -1, 1, 3},      // v_wood - -
	{"cr", "xyzzy", nil, nil, nil, 0, 0, 3},     // v_xyzzy - -
	{"c", "yew", nil, nil, nil, -1, 1, 3},       // v_yew - -

	{"cr", "north", nil, nil, nil, -1, 0, 2}, // v_north - -
	{"cr", "n", nil, nil, nil, -1, 0, 2},     // v_north - -
	{"cr", "s", nil, nil, nil, -1, 0, 2},     // v_south - -
	{"cr", "south", nil, nil, nil, -1, 0, 2}, // v_south - -
	{"cr", "east", nil, nil, nil, -1, 0, 2},  // v_east - -
	{"cr", "e", nil, nil, nil, -1, 0, 2},     // v_east - -
	{"cr", "west", nil, nil, nil, -1, 0, 2},  // v_west - -
	{"cr", "w", nil, nil, nil, -1, 0, 2},     // v_west - -
	{"cr", "enter", nil, nil, nil, -1, 0, 2}, // v_enter - -
	{"cr", "exit", nil, nil, nil, -1, 0, 2},  // v_exit - -
	{"cr", "in", nil, nil, nil, -1, 0, 2},    // v_enter - -
	{"cr", "out", nil, nil, nil, -1, 0, 2},   // v_exit - -

	{"", "begin", nil, nil, nil, 0, 0, 0},     // - - -
	{"", "unit", nil, nil, nil, 0, 0, 0},      // - - -
	{"", "email", nil, nil, nil, 0, 0, 0},     // - - -
	{"", "vis_email", nil, nil, nil, 0, 0, 0}, // - - -
	{"", "end", nil, nil, nil, 0, 0, 0},       // - - -
	{"", "flush", nil, nil, nil, 0, 0, 0},     // - - -
	{"", "lore", nil, nil, nil, 0, 0, 0},      // - - -
	{"", "passwd", nil, nil, nil, 0, 0, 0},    // - - -
	{"", "password", nil, nil, nil, 0, 0, 0},  // - - -
	{"", "players", nil, nil, nil, 0, 0, 0},   // - - -
	{"", "resend", nil, nil, nil, 0, 0, 0},    // - - -
	{"cpr", "stop", nil, nil, nil, 0, 0, 0},   // v_stop - -

	{"i", "look", nil, nil, nil, 0, 0, 1},       // v_look - -
	{"i", "l", nil, nil, nil, 0, 0, 1},          // v_look - -
	{"i", "ct", nil, nil, nil, 0, 0, 1},         // v_ct - -
	{"i", "be", nil, nil, nil, 0, 0, 1},         // v_be - -
	{"i", "additem", nil, nil, nil, 0, 0, 1},    // v_add_item - -
	{"i", "subitem", nil, nil, nil, 0, 0, 1},    // v_sub_item - -
	{"i", "h", nil, nil, nil, 0, 0, 1},          // v_listcmds - -
	{"i", "dump", nil, nil, nil, 0, 0, 1},       // v_dump - -
	{"i", "i", nil, nil, nil, 0, 0, 1},          // v_invent - -
	{"i", "fix", nil, nil, nil, 0, 0, 1},        // v_fix - -
	{"i", "fix2", nil, nil, nil, 0, 0, 1},       // v_fix2 - -
	{"i", "kill", nil, nil, nil, 0, 0, 1},       // v_kill - -
	{"i", "los", nil, nil, nil, 0, 0, 1},        // v_los - -
	{"m", "relore", nil, nil, nil, 0, 0, 1},     // v_relore - -
	{"i", "sk", nil, nil, nil, 0, 0, 1},         // v_skills - -
	{"i", "know", nil, nil, nil, 0, 0, 1},       // v_know - -
	{"i", "seed", nil, nil, nil, 0, 0, 1},       // v_seed - -
	{"i", "seedmarket", nil, nil, nil, 0, 0, 1}, // v_seedmarket - -
	{"i", "sheet", nil, nil, nil, 0, 0, 1},      // v_lore - -
	{"i", "poof", nil, nil, nil, 0, 0, 1},       // v_poof - -
	{"i", "postproc", nil, nil, nil, 0, 0, 1},   // v_postproc - -
	{"i", "save", nil, nil, nil, 0, 0, 1},       // v_save - -
	{"i", "seeall", nil, nil, nil, 0, 0, 1},     // v_see_all - -
	{"i", "tp", nil, nil, nil, 0, 0, 1},         // v_take_pris - -
	{"i", "makeloc", nil, nil, nil, 0, 0, 1},    // v_makeloc - -
	{"i", "remail", nil, nil, nil, 0, 0, 1},     // v_remail - -
}

func init() {
	// All nonempty allow rows have a start routine in the pinned cmd_tbl.
	// Keep real C NULLs nil, but never mistake an unported routine for NULL.
	for i := range cmdTbl {
		if cmdTbl[i].allow != "" {
			cmdTbl[i].start = unportedCommand(cmdTbl[i].name, "start")
		}
	}
	for _, name := range []string{
		"bind", "breed", "bribe", "build", "collect", "explore", "fly", "form",
		"go", "hide", "improve", "make", "move", "pillage", "quest", "raise",
		"rally", "raze", "repair", "research", "sail", "seek", "sneak", "study",
		"train", "trance", "terrorize", "torture", "use", "wait",
	} {
		for i := range cmdTbl {
			if cmdTbl[i].name == name {
				cmdTbl[i].finish = unportedCommand(name, "finish")
			}
		}
	}
	for _, name := range []string{"collect", "make", "repair", "sail", "train", "use", "wait"} {
		for i := range cmdTbl {
			if cmdTbl[i].name == name {
				cmdTbl[i].interrupt = unportedCommand(name, "interrupt")
			}
		}
	}
	for i := range cmdTbl {
		switch cmdTbl[i].name {
		case "raze":
			cmdTbl[i].start = (*world).vRaze
			cmdTbl[i].finish = (*world).dRaze
		case "repair":
			cmdTbl[i].start = (*world).vRepair
			cmdTbl[i].finish = (*world).dRepair
			cmdTbl[i].interrupt = (*world).iRepair
		case "improve":
			cmdTbl[i].start = (*world).vImprove
			cmdTbl[i].finish = (*world).dImprove
		case "build":
			cmdTbl[i].start = (*world).vBuild
			cmdTbl[i].finish = (*world).dBuild
		case "make", "train":
			cmdTbl[i].start = (*world).vMake
			cmdTbl[i].finish = (*world).dMake
			cmdTbl[i].interrupt = (*world).iMake
		case "collect":
			cmdTbl[i].start = (*world).vCollect
			cmdTbl[i].finish = (*world).dCollect
			cmdTbl[i].interrupt = (*world).iCollect
		case "recruit":
			cmdTbl[i].start = (*world).vRecruit
		case "quarry", "stone":
			cmdTbl[i].start = (*world).vQuarry
		case "fish":
			cmdTbl[i].start = (*world).vFish
		case "wood":
			cmdTbl[i].start = (*world).vWood
		case "opium":
			cmdTbl[i].start = (*world).vOpium
		case "mallorn":
			cmdTbl[i].start = (*world).vMallorn
		case "yew":
			cmdTbl[i].start = (*world).vYew
		case "catch":
			cmdTbl[i].start = (*world).vCatch
		case "honor", "honour":
			cmdTbl[i].start = (*world).vHonor
		case "oath":
			cmdTbl[i].start = (*world).vOath
		case "pledge":
			cmdTbl[i].start = (*world).vPledge
		case "swear":
			cmdTbl[i].start = (*world).vSwear
		case "raise":
			cmdTbl[i].start = (*world).vRaise
			cmdTbl[i].finish = (*world).dRaise
		case "rally":
			cmdTbl[i].start = (*world).vRally
			cmdTbl[i].finish = (*world).dRally
		case "bribe":
			cmdTbl[i].start = (*world).vBribe
			cmdTbl[i].finish = (*world).dBribe
		case "terrorize":
			cmdTbl[i].start = (*world).vTerrorize
			cmdTbl[i].finish = (*world).dTerrorize
		case "incite":
			cmdTbl[i].start = (*world).vIncite
		case "accept":
			cmdTbl[i].start = (*world).vAccept
		case "give":
			cmdTbl[i].start = (*world).vGive
		case "pay":
			cmdTbl[i].start = (*world).vPay
		case "get", "take":
			cmdTbl[i].start = (*world).vGet
		case "discard", "drop":
			cmdTbl[i].start = (*world).vDiscard
		case "claim":
			cmdTbl[i].start = (*world).vClaim
		case "reclaim":
			cmdTbl[i].start = (*world).vReclaim
		case "die":
			cmdTbl[i].start = (*world).vDie
		case "quit":
			cmdTbl[i].start = (*world).vQuit
		case "credit":
			cmdTbl[i].start = (*world).vCredit
		case "relore":
			cmdTbl[i].start = (*world).vRelore
		case "emote":
			cmdTbl[i].start = (*world).vEmote
		case "admit":
			cmdTbl[i].start = (*world).vAdmit
		case "form":
			cmdTbl[i].start = (*world).vForm
			cmdTbl[i].finish = (*world).dForm
		case "contact":
			cmdTbl[i].start = (*world).vContact
		case "default":
			cmdTbl[i].start = (*world).vAttClear
		case "defend":
			cmdTbl[i].start = (*world).vDefend
		case "hostile":
			cmdTbl[i].start = (*world).vHostile
		case "neutral":
			cmdTbl[i].start = (*world).vNeutral
		case "use":
			cmdTbl[i].start = (*world).vUse
			cmdTbl[i].finish = (*world).dUse
			cmdTbl[i].interrupt = (*world).iUse
		case "bind":
			cmdTbl[i].start = (*world).vBindStorm
			cmdTbl[i].finish = (*world).dBindStorm
		case "stop", "xyzzy", "plugh":
			// These three real C handlers do nothing except return TRUE.
			cmdTbl[i].start = func(*world, *command) int { return 1 }
		case "format":
			cmdTbl[i].start = (*world).vFormat
		case "notab":
			cmdTbl[i].start = (*world).vNotab
		case "name":
			cmdTbl[i].start = (*world).vName
		case "realname":
			cmdTbl[i].start = (*world).vFullname
		case "banner":
			cmdTbl[i].start = (*world).vBanner
		case "flag":
			cmdTbl[i].start = (*world).vFlag
		case "times":
			cmdTbl[i].start = (*world).vTimes
		case "post":
			cmdTbl[i].start = (*world).vPost
		case "message":
			cmdTbl[i].start = (*world).vMessage
		case "rumor":
			cmdTbl[i].start = (*world).vRumor
		case "press":
			cmdTbl[i].start = (*world).vPress
		case "public":
			cmdTbl[i].start = (*world).vPublic
		case "split":
			cmdTbl[i].start = (*world).vSplit
		case "guard":
			cmdTbl[i].start = (*world).vGuard
		case "stack":
			cmdTbl[i].start = (*world).vStack
		case "unstack":
			cmdTbl[i].start = (*world).vUnstack
		case "promote":
			cmdTbl[i].start = (*world).vPromote
		case "hide":
			cmdTbl[i].start = (*world).vHide
			cmdTbl[i].finish = (*world).dHide
		case "move", "go":
			cmdTbl[i].start = (*world).vMove
			cmdTbl[i].finish = (*world).dMove
		case "fly":
			cmdTbl[i].start = (*world).vFly
			cmdTbl[i].finish = (*world).dFly
		case "sail":
			cmdTbl[i].start = (*world).vSail
			cmdTbl[i].finish = (*world).dSail
			cmdTbl[i].interrupt = (*world).iSail
		case "fee":
			cmdTbl[i].start = (*world).vFee
		case "board":
			cmdTbl[i].start = (*world).vBoard
		case "unload":
			cmdTbl[i].start = (*world).vUnload
		case "ferry":
			cmdTbl[i].start = (*world).vFerry
		case "north", "n":
			cmdTbl[i].start = (*world).vNorth
		case "south", "s":
			cmdTbl[i].start = (*world).vSouth
		case "east", "e":
			cmdTbl[i].start = (*world).vEast
		case "west", "w":
			cmdTbl[i].start = (*world).vWest
		case "enter", "in":
			cmdTbl[i].start = (*world).vEnter
		case "exit", "out":
			cmdTbl[i].start = (*world).vExit
		case "wait":
			cmdTbl[i].start = (*world).vWait
			cmdTbl[i].finish = (*world).dWait
			cmdTbl[i].interrupt = (*world).iWait
		}
	}
}

func unportedCommand(name, phase string) func(*world, *command) int {
	return func(*world, *command) int {
		panic(notImplemented("command " + name + " " + phase))
	}
}
