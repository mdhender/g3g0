// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode"
)

// TestEntityMatchesC prints the entity model the way the C driver in
// tools/entity-compat/run.sh does and requires its output, line for line.
// testdata/entity-c-output.txt is that driver's output from the pinned G3
// sources; regenerate it only with run.sh --update.
//
// The C struct layouts it records are those of the host that recorded them
// (LP64, macOS arm64), so on a host with other pointer sizes the test is
// skipped.
func TestEntityMatchesC(t *testing.T) {
	if cPtrSize != 8 {
		t.Skip("testdata records an LP64 host")
	}
	want := readTestdata(t, "entity-c-output.txt")
	wantLines := strings.Split(strings.TrimSuffix(want, "\n"), "\n")

	var out strings.Builder
	p := func(format string, a ...any) { fmt.Fprintf(&out, format+"\n", a...) }

	// Constants, in the order the C driver prints them: every numeric
	// #define in oly.h must be in ohConstants with the same value.
	seen := map[string]bool{}
	for _, line := range wantLines {
		var name string
		var v int
		if _, err := fmt.Sscanf(line, "const %s %d", &name, &v); err != nil {
			continue
		}
		got, ok := ohConstants[name]
		if !ok {
			t.Errorf("oly.h constant %s (%d) is not ported", name, v)
			got = v
		}
		seen[name] = true
		p("const %s %d", name, got)
	}
	for name := range ohConstants {
		if !seen[name] {
			t.Errorf("ported constant %s is not in oly.h", name)
		}
	}

	// Struct layouts, derived from the Go types.
	for _, s := range ohStructs {
		size, align := cLayout(s.t)
		p("struct %s size %d align %d", s.label, size, align)
		var off uintptr
		for f := range s.t.Fields() {
			fs, fa := cLayout(f.Type)
			off = (off + fa - 1) &^ (fa - 1)
			p("\t%s offset %d size %d", snake(f.Name), off, fs)
			off += fs
		}
	}

	// glob.c tables.
	for i, s := range kindS {
		p("kind_s[%d] [%s]", i, s)
	}
	p("kind_s[%d] NULL", len(kindS))
	for i, s := range subkindS {
		p("subkind_s[%d] [%s]", i, s)
	}
	p("subkind_s[%d] NULL", len(subkindS))
	for i, c := range cmdTbl {
		p("cmd_tbl[%d] [%s] [%s] time %d poll %d pri %d", i, c.allow, c.name, c.time, c.poll, c.pri)
	}
	p("cmd_tbl[%d] NULL", len(cmdTbl))

	// glob.c globals after glob_init().
	w := newWorld()
	p("from_host [%s]", w.fromHost)
	p("reply_host [%s]", w.replyHost)
	p("gm_address [%s]", w.gmAddress)
	p("game_title [%s]", w.gameTitle)
	p("game_url [%s]", w.gameURL)
	p("rules_url [%s]", w.rulesURL)
	p("times_url [%s]", w.timesURL)
	p("htpasswd_loc [%s]", w.htpasswdLoc)
	p("garrison_pay = %d", w.garrisonPay)
	p("army_slow_factor = %d", w.armySlowFactor)
	p("auto_quit_turns = %d", w.autoQuitTurns)
	p("sysclock %d %d %d", w.sysclock.day, w.sysclock.turn, w.sysclock.daysSinceEpoch)
	p("indep_player = %d", w.indepPlayer)
	p("gm_player = %d", w.gmPlayer)
	p("skill_player = %d", w.skillPlayer)
	p("eat_pl = %d", w.eatPl)
	p("npc_pl = %d", w.npcPl)
	p("garr_pl = %d", w.garrPl)
	p("combat_pl = %d", w.combatPl)
	p("post_has_been_run = %d", w.postHasBeenRun)
	p("garrison_magic = %d", w.garrisonMagic)
	p("mount_olympus = %d", w.mountOlympus)
	n := 0
	for _, b := range w.bx {
		if b != nil {
			n++
		}
	}
	p("bx non-NULL %d", n)
	n = 0
	for _, h := range w.boxHead {
		if h != 0 {
			n++
		}
	}
	p("box_head non-zero %d", n)
	n = 0
	for _, h := range w.subHead {
		if h != 0 {
			n++
		}
	}
	p("sub_head non-zero %d", n)

	// Accessors.
	labels := map[any]int{}
	ptr := func(v any) string {
		if reflect.ValueOf(v).IsNil() {
			return "NULL"
		}
		if _, ok := labels[v]; !ok {
			labels[v] = len(labels) + 1
		}
		return fmt.Sprintf("#%d", labels[v])
	}
	showI := func(expr string, v int) { p("%s = %d", expr, v) }
	showP := func(expr string, v any) { p("%s = %s", expr, ptr(v)) }
	showZ := func(expr string, v any) {
		p("%s = %s zero %d", expr, ptr(v), b2i(reflect.ValueOf(v).Elem().IsZero()))
	}
	newBox := func(n, k, sk int) { w.bx[n] = &box{kind: int8(k), skind: int8(sk)} }

	newBox(1001, tChar, 0)
	newBox(1002, tItem, subMagic)
	newBox(1003, tChar, subNi)
	newBox(1004, tItem, subArtifact)
	newBox(maxBoxes-1, tLoc, subOcean)
	newBox(1, tItem, 0)

	showI("kind(-1)", w.kind(-1))
	showI("kind(0)", w.kind(0))
	showI("kind(1)", w.kind(1))
	showI("kind(2)", w.kind(2))
	showI("kind(1001)", w.kind(1001))
	showI("kind(1002)", w.kind(1002))
	showI("kind(1005)", w.kind(1005))
	showI("kind(MAX_BOXES - 1)", w.kind(maxBoxes-1))
	showI("kind(MAX_BOXES)", w.kind(maxBoxes))
	showI("kind(MAX_BOXES + 1)", w.kind(maxBoxes+1))
	showI("valid_box(0)", b2i(w.validBox(0)))
	showI("valid_box(1001)", b2i(w.validBox(1001)))
	showI("valid_box(1005)", b2i(w.validBox(1005)))
	showI("valid_box(MAX_BOXES)", b2i(w.validBox(maxBoxes)))
	showI("subkind(0)", w.subkind(0))
	showI("subkind(1001)", w.subkind(1001))
	showI("subkind(1002)", w.subkind(1002))
	showI("subkind(1003)", w.subkind(1003))
	showI("subkind(1005)", w.subkind(1005))
	showI("subkind(MAX_BOXES - 1)", w.subkind(maxBoxes-1))

	showP("p_loc_info(1001)", w.pLocInfo(1001))
	showP("rp_loc_info(1001)", w.rpLocInfo(1001))
	showI("zero(p_loc_info(1001), sizeof(struct loc_info))", b2i(reflect.ValueOf(*w.pLocInfo(1001)).IsZero()))

	// accessor runs P() of the C driver: rp, p twice, rp.
	accessor := func(name string, n int, rp, pp func(int) any) {
		showP(fmt.Sprintf("rp_%s(%d)", name, n), rp(n))
		showZ(fmt.Sprintf("p_%s(%d)", name, n), pp(n))
		showP(fmt.Sprintf("p_%s(%d)", name, n), pp(n))
		showP(fmt.Sprintf("rp_%s(%d)", name, n), rp(n))
	}

	accessor("char", 1001, func(n int) any { return w.rpChar(n) }, func(n int) any { return w.pChar(n) })
	showP("rp_magic(1001)", w.rpMagic(1001))
	showZ("p_magic(1001)", w.pMagic(1001))
	showP("p_magic(1001)", w.pMagic(1001))
	showP("rp_magic(1001)", w.rpMagic(1001))
	showP("p_char(1001)", w.pChar(1001))

	showP("rp_magic(1003)", w.rpMagic(1003))
	showP("rp_char(1003)", w.rpChar(1003))
	showZ("p_magic(1003)", w.pMagic(1003))
	showZ("rp_char(1003)", w.rpChar(1003))
	showP("p_char(1003)", w.pChar(1003))

	accessor("item", 1002, func(n int) any { return w.rpItem(n) }, func(n int) any { return w.pItem(n) })
	showP("rp_item_magic(1002)", w.rpItemMagic(1002))
	showZ("p_item_magic(1002)", w.pItemMagic(1002))
	showP("p_item_magic(1002)", w.pItemMagic(1002))
	showP("rp_item_magic(1002)", w.rpItemMagic(1002))

	showP("rp_item_magic(1004)", w.rpItemMagic(1004))
	showP("rp_item(1004)", w.rpItem(1004))
	showZ("p_item_magic(1004)", w.pItemMagic(1004))
	showZ("rp_item(1004)", w.rpItem(1004))

	accessor("loc", 1001, func(n int) any { return w.rpLoc(n) }, func(n int) any { return w.pLoc(n) })
	accessor("subloc", 1001, func(n int) any { return w.rpSubloc(n) }, func(n int) any { return w.pSubloc(n) })
	accessor("player", 1001, func(n int) any { return w.rpPlayer(n) }, func(n int) any { return w.pPlayer(n) })
	accessor("skill", 1001, func(n int) any { return w.rpSkill(n) }, func(n int) any { return w.pSkill(n) })
	accessor("gate", 1001, func(n int) any { return w.rpGate(n) }, func(n int) any { return w.pGate(n) })
	accessor("misc", 1001, func(n int) any { return w.rpMisc(n) }, func(n int) any { return w.pMisc(n) })
	accessor("disp", 1001, func(n int) any { return w.rpDisp(n) }, func(n int) any { return w.pDisp(n) })
	accessor("command", 1001, func(n int) any { return w.rpCommand(n) }, func(n int) any { return w.pCommand(n) })

	showP("rp_char(1002)", w.rpChar(1002))
	showP("rp_item(1001)", w.rpItem(1001))
	showP("rp_command(1003)", w.rpCommand(1003))

	gotLines := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
	for i := 0; i < len(wantLines) || i < len(gotLines); i++ {
		var w, g string
		if i < len(wantLines) {
			w = wantLines[i]
		}
		if i < len(gotLines) {
			g = gotLines[i]
		}
		if w != g {
			t.Fatalf("line %d:\n got %q\nwant %q", i+1, g, w)
		}
	}
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// snake turns a Go field name back into the C one: xLocInfo is x_loc_info.
func snake(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsUpper(r) {
			b.WriteByte('_')
			r = unicode.ToLower(r)
		}
		b.WriteRune(r)
	}
	return b.String()
}

// ohStructs lists the structures of oly.h in the order the header defines
// them, with the label the C driver prints.
var ohStructs = []struct {
	label string
	t     reflect.Type
}{
	{"loc_info", reflect.TypeFor[locInfo]()},
	{"box", reflect.TypeFor[box]()},
	{"entity_player", reflect.TypeFor[entityPlayer]()},
	{"order_list", reflect.TypeFor[orderList]()},
	{"accept_ent", reflect.TypeFor[acceptEnt]()},
	{"att_ent", reflect.TypeFor[attEnt]()},
	{"entity_char", reflect.TypeFor[entityChar]()},
	{"char_magic", reflect.TypeFor[charMagic]()},
	{"skill_ent", reflect.TypeFor[skillEnt]()},
	{"item_ent", reflect.TypeFor[itemEnt]()},
	{"entity_loc", reflect.TypeFor[entityLoc]()},
	{"entity_subloc", reflect.TypeFor[entitySubloc]()},
	{"entity_item", reflect.TypeFor[entityItem]()},
	{"item_magic", reflect.TypeFor[itemMagic]()},
	{"entity_skill", reflect.TypeFor[entitySkill]()},
	{"req_ent", reflect.TypeFor[reqEnt]()},
	{"entity_gate", reflect.TypeFor[entityGate]()},
	{"entity_misc", reflect.TypeFor[entityMisc]()},
	{"wait_arg", reflect.TypeFor[waitArg]()},
	{"command", reflect.TypeFor[command]()},
	{"cmd_tbl_ent", reflect.TypeFor[cmdTblEnt]()},
	{"trade", reflect.TypeFor[trade]()},
	{"admit", reflect.TypeFor[admit]()},
	{"olytime", reflect.TypeFor[olytime]()},
}

// ohConstants maps each numeric #define of oly.h to its port.
var ohConstants = map[string]int{
	"MAX_BOXES":                maxBoxes,
	"MONTH_DAYS":               monthDays,
	"NUM_MONTHS":               numMonths,
	"T_deleted":                tDeleted,
	"T_player":                 tPlayer,
	"T_char":                   tChar,
	"T_loc":                    tLoc,
	"T_item":                   tItem,
	"T_skill":                  tSkill,
	"T_gate":                   tGate,
	"T_road":                   tRoad,
	"T_deadchar":               tDeadchar,
	"T_ship":                   tShip,
	"T_post":                   tPost,
	"T_storm":                  tStorm,
	"T_unform":                 tUnform,
	"T_lore":                   tLore,
	"T_MAX":                    tMax,
	"sub_ocean":                subOcean,
	"sub_forest":               subForest,
	"sub_plain":                subPlain,
	"sub_mountain":             subMountain,
	"sub_desert":               subDesert,
	"sub_swamp":                subSwamp,
	"sub_under":                subUnder,
	"sub_faery_hill":           subFaeryHill,
	"sub_island":               subIsland,
	"sub_stone_cir":            subStoneCir,
	"sub_mallorn_grove":        subMallornGrove,
	"sub_bog":                  subBog,
	"sub_cave":                 subCave,
	"sub_city":                 subCity,
	"sub_lair":                 subLair,
	"sub_graveyard":            subGraveyard,
	"sub_ruins":                subRuins,
	"sub_battlefield":          subBattlefield,
	"sub_ench_forest":          subEnchForest,
	"sub_rocky_hill":           subRockyHill,
	"sub_tree_circle":          subTreeCircle,
	"sub_pits":                 subPits,
	"sub_pasture":              subPasture,
	"sub_oasis":                subOasis,
	"sub_yew_grove":            subYewGrove,
	"sub_sand_pit":             subSandPit,
	"sub_sacred_grove":         subSacredGrove,
	"sub_poppy_field":          subPoppyField,
	"sub_temple":               subTemple,
	"sub_galley":               subGalley,
	"sub_roundship":            subRoundship,
	"sub_castle":               subCastle,
	"sub_galley_notdone":       subGalleyNotdone,
	"sub_roundship_notdone":    subRoundshipNotdone,
	"sub_ghost_ship":           subGhostShip,
	"sub_temple_notdone":       subTempleNotdone,
	"sub_inn":                  subInn,
	"sub_inn_notdone":          subInnNotdone,
	"sub_castle_notdone":       subCastleNotdone,
	"sub_mine":                 subMine,
	"sub_mine_notdone":         subMineNotdone,
	"sub_scroll":               subScroll,
	"sub_magic":                subMagic,
	"sub_palantir":             subPalantir,
	"sub_auraculum":            subAuraculum,
	"sub_tower":                subTower,
	"sub_tower_notdone":        subTowerNotdone,
	"sub_pl_system":            subPlSystem,
	"sub_pl_regular":           subPlRegular,
	"sub_region":               subRegion,
	"sub_pl_savage":            subPlSavage,
	"sub_pl_npc":               subPlNpc,
	"sub_mine_collapsed":       subMineCollapsed,
	"sub_ni":                   subNi,
	"sub_undead":               subUndead,
	"sub_dead_body":            subDeadBody,
	"sub_fog":                  subFog,
	"sub_wind":                 subWind,
	"sub_rain":                 subRain,
	"sub_hades_pit":            subHadesPit,
	"sub_artifact":             subArtifact,
	"sub_pl_silent":            subPlSilent,
	"sub_npc_token":            subNpcToken,
	"sub_garrison":             subGarrison,
	"sub_cloud":                subCloud,
	"sub_raft":                 subRaft,
	"sub_raft_notdone":         subRaftNotdone,
	"sub_suffuse_ring":         subSuffuseRing,
	"sub_relic":                subRelic,
	"sub_tunnel":               subTunnel,
	"sub_sewer":                subSewer,
	"sub_chamber":              subChamber,
	"sub_tradegood":            subTradegood,
	"SUB_MAX":                  subMax,
	"item_gold":                itemGold,
	"item_peasant":             itemPeasant,
	"item_worker":              itemWorker,
	"item_soldier":             itemSoldier,
	"item_archer":              itemArcher,
	"item_knight":              itemKnight,
	"item_elite_guard":         itemEliteGuard,
	"item_pikeman":             itemPikeman,
	"item_blessed_soldier":     itemBlessedSoldier,
	"item_ghost_warrior":       itemGhostWarrior,
	"item_sailor":              itemSailor,
	"item_swordsman":           itemSwordsman,
	"item_crossbowman":         itemCrossbowman,
	"item_elite_arch":          itemEliteArch,
	"item_angry_peasant":       itemAngryPeasant,
	"item_pirate":              itemPirate,
	"item_elf":                 itemElf,
	"item_spirit":              itemSpirit,
	"item_corpse":              itemCorpse,
	"item_savage":              itemSavage,
	"item_skeleton":            itemSkeleton,
	"item_barbarian":           itemBarbarian,
	"item_wild_horse":          itemWildHorse,
	"item_riding_horse":        itemRidingHorse,
	"item_warmount":            itemWarmount,
	"item_pegasus":             itemPegasus,
	"item_nazgul":              itemNazgul,
	"item_flotsam":             itemFlotsam,
	"item_battering_ram":       itemBatteringRam,
	"item_catapult":            itemCatapult,
	"item_siege_tower":         itemSiegeTower,
	"item_ratspider_venom":     itemRatspiderVenom,
	"item_lana_bark":           itemLanaBark,
	"item_avinia_leaf":         itemAviniaLeaf,
	"item_spiny_root":          itemSpinyRoot,
	"item_farrenstone":         itemFarrenstone,
	"item_yew":                 itemYew,
	"item_elfstone":            itemElfstone,
	"item_mallorn_wood":        itemMallornWood,
	"item_pretus_bones":        itemPretusBones,
	"item_longbow":             itemLongbow,
	"item_plate":               itemPlate,
	"item_longsword":           itemLongsword,
	"item_pike":                itemPike,
	"item_ox":                  itemOx,
	"item_lumber":              itemLumber,
	"item_stone":               itemStone,
	"item_iron":                itemIron,
	"item_leather":             itemLeather,
	"item_ratspider":           itemRatspider,
	"item_mithril":             itemMithril,
	"item_gate_crystal":        itemGateCrystal,
	"item_blank_scroll":        itemBlankScroll,
	"item_crossbow":            itemCrossbow,
	"item_fish":                itemFish,
	"item_opium":               itemOpium,
	"item_basket":              itemBasket,
	"item_pot":                 itemPot,
	"item_tax_cookie":          itemTaxCookie,
	"item_drum":                itemDrum,
	"item_hide":                itemHide,
	"item_mob_cookie":          itemMobCookie,
	"item_lead":                itemLead,
	"item_glue":                itemGlue,
	"item_centaur":             itemCentaur,
	"item_minotaur":            itemMinotaur,
	"item_undead_cookie":       itemUndeadCookie,
	"item_fog_cookie":          itemFogCookie,
	"item_wind_cookie":         itemWindCookie,
	"item_rain_cookie":         itemRainCookie,
	"item_mage_menial":         itemMageMenial,
	"item_spider":              itemSpider,
	"item_rat":                 itemRat,
	"item_lion":                itemLion,
	"item_bird":                itemBird,
	"item_lizard":              itemLizard,
	"item_bandit":              itemBandit,
	"item_chimera":             itemChimera,
	"item_harpie":              itemHarpie,
	"item_dragon":              itemDragon,
	"item_orc":                 itemOrc,
	"item_gorgon":              itemGorgon,
	"item_wolf":                itemWolf,
	"item_orb":                 itemOrb,
	"item_cyclops":             itemCyclops,
	"item_giant":               itemGiant,
	"item_faery":               itemFaery,
	"item_petty_thief":         itemPettyThief,
	"item_hound":               itemHound,
	"lore_skeleton_npc_token":  loreSkeletonNpcToken,
	"lore_orc_npc_token":       loreOrcNpcToken,
	"lore_undead_npc_token":    loreUndeadNpcToken,
	"lore_savage_npc_token":    loreSavageNpcToken,
	"lore_barbarian_npc_token": loreBarbarianNpcToken,
	"lore_orb":                 loreOrb,
	"lore_faery_stone":         loreFaeryStone,
	"lore_barbarian_kill":      loreBarbarianKill,
	"lore_savage_kill":         loreSavageKill,
	"lore_undead_kill":         loreUndeadKill,
	"lore_orc_kill":            loreOrcKill,
	"lore_skeleton_kill":       loreSkeletonKill,
	"sk_shipcraft":             skShipcraft,
	"sk_pilot_ship":            skPilotShip,
	"sk_shipbuilding":          skShipbuilding,
	"sk_fishing":               skFishing,
	"sk_combat":                skCombat,
	"sk_survive_fatal":         skSurviveFatal,
	"sk_fight_to_death":        skFightToDeath,
	"sk_make_catapult":         skMakeCatapult,
	"sk_defense":               skDefense,
	"sk_archery":               skArchery,
	"sk_swordplay":             skSwordplay,
	"sk_weaponsmith":           skWeaponsmith,
	"sk_stealth":               skStealth,
	"sk_petty_thief":           skPettyThief,
	"sk_spy_inv":               skSpyInv,
	"sk_spy_skills":            skSpySkills,
	"sk_spy_lord":              skSpyLord,
	"sk_hide_lord":             skHideLord,
	"sk_find_rich":             skFindRich,
	"sk_torture":               skTorture,
	"sk_hide_self":             skHideSelf,
	"sk_sneak_build":           skSneakBuild,
	"sk_beast":                 skBeast,
	"sk_bird_spy":              skBirdSpy,
	"sk_capture_beasts":        skCaptureBeasts,
	"sk_use_beasts":            skUseBeasts,
	"sk_breed_beasts":          skBreedBeasts,
	"sk_catch_horse":           skCatchHorse,
	"sk_train_wild":            skTrainWild,
	"sk_train_warmount":        skTrainWarmount,
	"sk_summon_savage":         skSummonSavage,
	"sk_keep_savage":           skKeepSavage,
	"sk_breed_hound":           skBreedHound,
	"sk_persuasion":            skPersuasion,
	"sk_bribe_noble":           skBribeNoble,
	"sk_persuade_oath":         skPersuadeOath,
	"sk_raise_mob":             skRaiseMob,
	"sk_rally_mob":             skRallyMob,
	"sk_incite_mob":            skInciteMob,
	"sk_train_angry":           skTrainAngry,
	"sk_construction":          skConstruction,
	"sk_make_siege":            skMakeSiege,
	"sk_quarry_stone":          skQuarryStone,
	"sk_alchemy":               skAlchemy,
	"sk_brew_heal":             skBrewHeal,
	"sk_record_skill":          skRecordSkill,
	"sk_extract_venom":         skExtractVenom,
	"sk_brew_slave":            skBrewSlave,
	"sk_collect_elem":          skCollectElem,
	"sk_brew_death":            skBrewDeath,
	"sk_lead_to_gold":          skLeadToGold,
	"sk_forestry":              skForestry,
	"sk_make_ram":              skMakeRam,
	"sk_harvest_lumber":        skHarvestLumber,
	"sk_harvest_yew":           skHarvestYew,
	"sk_collect_foliage":       skCollectFoliage,
	"sk_harvest_mallorn":       skHarvestMallorn,
	"sk_harvest_opium":         skHarvestOpium,
	"sk_improve_opium":         skImproveOpium,
	"sk_mining":                skMining,
	"sk_mine_iron":             skMineIron,
	"sk_mine_gold":             skMineGold,
	"sk_mine_mithril":          skMineMithril,
	"sk_trade":                 skTrade,
	"sk_cloak_trade":           skCloakTrade,
	"sk_find_sell":             skFindSell,
	"sk_find_buy":              skFindBuy,
	"sk_religion":              skReligion,
	"sk_reveal_vision":         skRevealVision,
	"sk_last_rites":            skLastRites,
	"sk_pray":                  skPray,
	"sk_resurrect":             skResurrect,
	"sk_remove_bless":          skRemoveBless,
	"sk_vision_protect":        skVisionProtect,
	"sk_basic":                 skBasic,
	"sk_meditate":              skMeditate,
	"sk_mage_menial":           skMageMenial,
	"sk_appear_common":         skAppearCommon,
	"sk_view_aura":             skViewAura,
	"sk_heal":                  skHeal,
	"sk_write_basic":           skWriteBasic,
	"sk_reveal_mage":           skRevealMage,
	"sk_tap_health":            skTapHealth,
	"sk_shroud_abil":           skShroudAbil,
	"sk_detect_abil":           skDetectAbil,
	"sk_dispel_abil":           skDispelAbil,
	"sk_adv_med":               skAdvMed,
	"sk_hinder_med":            skHinderMed,
	"sk_weather":               skWeather,
	"sk_fierce_wind":           skFierceWind,
	"sk_bind_storm":            skBindStorm,
	"sk_write_weather":         skWriteWeather,
	"sk_summon_wind":           skSummonWind,
	"sk_summon_rain":           skSummonRain,
	"sk_summon_fog":            skSummonFog,
	"sk_direct_storm":          skDirectStorm,
	"sk_dissipate":             skDissipate,
	"sk_renew_storm":           skRenewStorm,
	"sk_lightning":             skLightning,
	"sk_seize_storm":           skSeizeStorm,
	"sk_death_fog":             skDeathFog,
	"sk_scry":                  skScry,
	"sk_scry_region":           skScryRegion,
	"sk_write_scry":            skWriteScry,
	"sk_shroud_region":         skShroudRegion,
	"sk_dispel_region":         skDispelRegion,
	"sk_bar_loc":               skBarLoc,
	"sk_unbar_loc":             skUnbarLoc,
	"sk_locate_char":           skLocateChar,
	"sk_detect_scry":           skDetectScry,
	"sk_proj_cast":             skProjCast,
	"sk_save_proj":             skSaveProj,
	"sk_banish_corpses":        skBanishCorpses,
	"sk_gate":                  skGate,
	"sk_detect_gates":          skDetectGates,
	"sk_jump_gate":             skJumpGate,
	"sk_write_gate":            skWriteGate,
	"sk_seal_gate":             skSealGate,
	"sk_unseal_gate":           skUnsealGate,
	"sk_notify_unseal":         skNotifyUnseal,
	"sk_rem_seal":              skRemSeal,
	"sk_reveal_key":            skRevealKey,
	"sk_notify_jump":           skNotifyJump,
	"sk_teleport":              skTeleport,
	"sk_rev_jump":              skRevJump,
	"sk_artifact":              skArtifact,
	"sk_forge_aura":            skForgeAura,
	"sk_write_art":             skWriteArt,
	"sk_forge_weapon":          skForgeWeapon,
	"sk_forge_armor":           skForgeArmor,
	"sk_forge_bow":             skForgeBow,
	"sk_curse_noncreat":        skCurseNoncreat,
	"sk_show_art_creat":        skShowArtCreat,
	"sk_show_art_reg":          skShowArtReg,
	"sk_destroy_art":           skDestroyArt,
	"sk_cloak_creat":           skCloakCreat,
	"sk_cloak_reg":             skCloakReg,
	"sk_rem_art_cloak":         skRemArtCloak,
	"sk_forge_palantir":        skForgePalantir,
	"sk_necromancy":            skNecromancy,
	"sk_raise_corpses":         skRaiseCorpses,
	"sk_summon_ghost":          skSummonGhost,
	"sk_write_necro":           skWriteNecro,
	"sk_undead_lord":           skUndeadLord,
	"sk_renew_undead":          skRenewUndead,
	"sk_banish_undead":         skBanishUndead,
	"sk_eat_dead":              skEatDead,
	"sk_aura_blast":            skAuraBlast,
	"sk_absorb_blast":          skAbsorbBlast,
	"sk_transcend_death":       skTranscendDeath,
	"sk_adv_sorcery":           skAdvSorcery,
	"sk_trance":                skTrance,
	"sk_teleport_item":         skTeleportItem,
	"sk_quick_cast":            skQuickCast,
	"sk_save_quick":            skSaveQuick,
	"sk_add_ram":               skAddRam,
	"PROG_bandit":              progBandit,
	"PROG_subloc_monster":      progSublocMonster,
	"PROG_npc_token":           progNpcToken,
	"PROG_faery_bandit":        progFaeryBandit,
	"PROG_hades_bandit":        progHadesBandit,
	"use_death_potion":         useDeathPotion,
	"use_heal_potion":          useHealPotion,
	"use_slave_potion":         useSlavePotion,
	"use_palantir":             usePalantir,
	"use_proj_cast":            useProjCast,
	"use_quick_cast":           useQuickCast,
	"use_drum":                 useDrum,
	"use_faery_stone":          useFaeryStone,
	"use_orb":                  useOrb,
	"use_barbarian_kill":       useBarbarianKill,
	"use_savage_kill":          useSavageKill,
	"use_corpse_kill":          useCorpseKill,
	"use_orc_kill":             useOrcKill,
	"use_skeleton_kill":        useSkeletonKill,
	"use_bta_skull":            useBtaSkull,
	"USE_KEY_MAX":              useKeyMax,
	"DIR_N":                    dirN,
	"DIR_E":                    dirE,
	"DIR_S":                    dirS,
	"DIR_W":                    dirW,
	"DIR_UP":                   dirUp,
	"DIR_DOWN":                 dirDown,
	"DIR_IN":                   dirIn,
	"DIR_OUT":                  dirOut,
	"MAX_DIR":                  maxDir,
	"LOC_region":               locRegion,
	"LOC_province":             locProvince,
	"LOC_subloc":               locSubloc,
	"LOC_build":                locBuild,
	"LOY_UNCHANGED":            loyUnchanged,
	"LOY_unsworn":              loyUnsworn,
	"LOY_contract":             loyContract,
	"LOY_oath":                 loyOath,
	"LOY_fear":                 loyFear,
	"LOY_npc":                  loyNpc,
	"LOY_summon":               loySummon,
	"exp_novice":               expNovice,
	"exp_journeyman":           expJourneyman,
	"exp_teacher":              expTeacher,
	"exp_master":               expMaster,
	"exp_grand":                expGrand,
	"ATT_NONE":                 attNone,
	"NEUTRAL":                  attNeutral,
	"HOSTILE":                  attHostile,
	"DEFEND":                   attDefend,
	"SKILL_dont":               skillDont,
	"SKILL_learning":           skillLearning,
	"SKILL_know":               skillKnow,
	"REQ_NO":                   reqNo,
	"REQ_YES":                  reqYes,
	"REQ_OR":                   reqOr,
	"STATE_DONE":               stateDone,
	"STATE_LOAD":               stateLoad,
	"STATE_RUN":                stateRun,
	"STATE_ERROR":              stateError,
	"BUY":                      tradeBuy,
	"SELL":                     tradeSell,
	"PRODUCE":                  tradeProduce,
	"CONSUME":                  tradeConsume,
	"COV_REJECTED":             covRejected,
	"COV_UNIMPL":               covUnimpl,
	"COV_OK":                   covOk,
	"COV_FAIL":                 covFail,
	"COV_MAX":                  covMax,
	"UCOV_CHECKS":              ucovChecks,
	"UCOV_OK":                  ucovOk,
	"UCOV_FAIL":                ucovFail,
	"UCOV_MAX":                 ucovMax,
	"CHAR_FIELD":               charField,
	"MAX_POST":                 maxPost,
	"STYLE_TEXT":               styleText,
	"STYLE_HTML":               styleHTML,
	"STYLE_PREV":               stylePrev,
	"RELIC_THRONE":             relicThrone,
	"RELIC_CROWN":              relicCrown,
	"RELIC_BTA_SKULL":          relicBtaSkull,
}
