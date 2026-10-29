// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "fmt"

// Skill lookup and report helpers from olympia/use.c.

// useTbl is the complete use.c:84–229 table. Index zero is a placeholder;
// the separate final zero row terminates walks starting at index one.
type useTblEnt struct {
	allow                    *string
	skill                    int
	start, finish, interrupt func(*world, *command) int
	time, poll               int32
}

var useAllow = "c"

// A non-NULL C callback must never become a nil Go callback merely because
// its gameplay is not ported. Only these four callbacks already exist in Go.
func useCallback(name string) func(*world, *command) int {
	switch name {
	case "v_hide":
		return (*world).vHide
	case "d_hide":
		return (*world).dHide
	case "v_bind_storm":
		return (*world).vBindStorm
	case "d_bind_storm":
		return (*world).dBindStorm
	}
	return unportedCommand("USE "+name, "callback")
}

var useTbl = []useTblEnt{
	{},
	{&useAllow, skMeditate, useCallback("v_meditate"), useCallback("d_meditate"), nil, 7, 0},
	{&useAllow, skDetectGates, useCallback("v_detect_gates"), useCallback("d_detect_gates"), nil, 7, 0},
	{&useAllow, skJumpGate, useCallback("v_jump_gate"), nil, nil, 1, 0},
	{&useAllow, skTeleport, useCallback("v_teleport"), nil, nil, 1, 0},
	{&useAllow, skSealGate, useCallback("v_seal_gate"), useCallback("d_seal_gate"), nil, 7, 0},
	{&useAllow, skUnsealGate, useCallback("v_unseal_gate"), useCallback("d_unseal_gate"), nil, 7, 0},
	{&useAllow, skNotifyUnseal, useCallback("v_notify_unseal"), useCallback("d_notify_unseal"), nil, 7, 0},
	{&useAllow, skRemSeal, useCallback("v_rem_seal"), useCallback("d_rem_seal"), nil, 7, 0},
	{&useAllow, skRevealKey, useCallback("v_reveal_key"), useCallback("d_reveal_key"), nil, 7, 0},
	{&useAllow, skNotifyJump, useCallback("v_notify_jump"), useCallback("d_notify_jump"), nil, 7, 0},
	{&useAllow, skHeal, useCallback("v_heal"), useCallback("d_heal"), nil, 7, 0},
	{&useAllow, skRevJump, useCallback("v_reverse_jump"), nil, nil, 1, 0},
	{&useAllow, skRevealMage, useCallback("v_reveal_mage"), useCallback("d_reveal_mage"), nil, 7, 0},
	{&useAllow, skViewAura, useCallback("v_view_aura"), useCallback("d_view_aura"), nil, 7, 0},
	{&useAllow, skShroudAbil, useCallback("v_shroud_abil"), useCallback("d_shroud_abil"), nil, 3, 0},
	{&useAllow, skDetectAbil, useCallback("v_detect_abil"), useCallback("d_detect_abil"), nil, 7, 0},
	{&useAllow, skScryRegion, useCallback("v_scry_region"), useCallback("d_scry_region"), nil, 7, 0},
	{&useAllow, skShroudRegion, useCallback("v_shroud_region"), useCallback("d_shroud_region"), nil, 3, 0},
	{&useAllow, skDetectScry, useCallback("v_detect_scry"), useCallback("d_detect_scry"), nil, 7, 0},
	{&useAllow, skDispelRegion, useCallback("v_dispel_region"), useCallback("d_dispel_region"), nil, 3, 0},
	{&useAllow, skDispelAbil, useCallback("v_dispel_abil"), useCallback("d_dispel_abil"), nil, 3, 0},
	{&useAllow, skAdvMed, useCallback("v_adv_med"), useCallback("d_adv_med"), nil, 7, 0},
	{&useAllow, skHinderMed, useCallback("v_hinder_med"), useCallback("d_hinder_med"), nil, 10, 0},
	{&useAllow, skProjCast, useCallback("v_proj_cast"), useCallback("d_proj_cast"), nil, 7, 0},
	{&useAllow, skLocateChar, useCallback("v_locate_char"), useCallback("d_locate_char"), nil, 10, 0},
	{&useAllow, skBarLoc, useCallback("v_bar_loc"), useCallback("d_bar_loc"), nil, 10, 0},
	{&useAllow, skUnbarLoc, useCallback("v_unbar_loc"), useCallback("d_unbar_loc"), nil, 7, 0},
	{&useAllow, skForgePalantir, useCallback("v_forge_palantir"), useCallback("d_forge_palantir"), nil, 10, 0},
	{&useAllow, skDestroyArt, useCallback("v_destroy_art"), useCallback("d_destroy_art"), nil, 7, 0},
	{&useAllow, skShowArtCreat, useCallback("v_show_art_creat"), useCallback("d_show_art_creat"), nil, 7, 0},
	{&useAllow, skShowArtReg, useCallback("v_show_art_reg"), useCallback("d_show_art_reg"), nil, 7, 0},
	{&useAllow, skSaveProj, useCallback("v_save_proj"), useCallback("d_save_proj"), nil, 7, 0},
	{&useAllow, skSaveQuick, useCallback("v_save_quick"), useCallback("d_save_quick"), nil, 7, 0},
	{&useAllow, skQuickCast, useCallback("v_quick_cast"), useCallback("d_quick_cast"), nil, 4, 0},
	{&useAllow, skRemArtCloak, useCallback("v_rem_art_cloak"), useCallback("d_rem_art_cloak"), nil, 10, 0},
	{&useAllow, skWriteBasic, useCallback("v_write_spell"), useCallback("d_write_spell"), nil, 7, 0},
	{&useAllow, skWriteWeather, useCallback("v_write_spell"), useCallback("d_write_spell"), nil, 7, 0},
	{&useAllow, skWriteScry, useCallback("v_write_spell"), useCallback("d_write_spell"), nil, 7, 0},
	{&useAllow, skWriteGate, useCallback("v_write_spell"), useCallback("d_write_spell"), nil, 7, 0},
	{&useAllow, skWriteArt, useCallback("v_write_spell"), useCallback("d_write_spell"), nil, 7, 0},
	{&useAllow, skWriteNecro, useCallback("v_write_spell"), useCallback("d_write_spell"), nil, 7, 0},
	{&useAllow, skCloakCreat, useCallback("v_cloak_creat"), useCallback("d_cloak_creat"), nil, 7, 0},
	{&useAllow, skCloakReg, useCallback("v_cloak_reg"), useCallback("d_cloak_reg"), nil, 7, 0},
	{&useAllow, skCurseNoncreat, useCallback("v_curse_noncreat"), useCallback("d_curse_noncreat"), nil, 14, 0},
	{&useAllow, skForgeAura, useCallback("v_forge_aura"), useCallback("d_forge_aura"), nil, 14, 0},
	{&useAllow, skShipbuilding, useCallback("v_shipbuild"), nil, nil, 0, 0},
	{&useAllow, skPilotShip, useCallback("v_sail"), useCallback("d_sail"), useCallback("i_sail"), -1, 0},
	{&useAllow, skTrainWild, useCallback("v_use_train_riding"), nil, nil, 7, 0},
	{&useAllow, skTrainWarmount, useCallback("v_use_train_war"), nil, nil, 14, 0},
	{&useAllow, skMakeRam, nil, nil, nil, 14, 0},
	{&useAllow, skMakeCatapult, nil, nil, nil, 14, 0},
	{&useAllow, skMakeSiege, nil, nil, nil, 14, 0},
	{&useAllow, skBrewSlave, useCallback("v_brew"), useCallback("d_brew_slave"), nil, 7, 0},
	{&useAllow, skBrewHeal, useCallback("v_brew"), useCallback("d_brew_heal"), nil, 7, 0},
	{&useAllow, skBrewDeath, useCallback("v_brew"), useCallback("d_brew_death"), nil, 10, 0},
	{&useAllow, skMineIron, useCallback("v_mine_iron"), useCallback("d_mine_iron"), nil, 7, 0},
	{&useAllow, skMineGold, useCallback("v_mine_gold"), useCallback("d_mine_gold"), nil, 7, 0},
	{&useAllow, skMineMithril, useCallback("v_mine_mithril"), useCallback("d_mine_mithril"), nil, 7, 0},
	{&useAllow, skQuarryStone, useCallback("v_quarry"), nil, nil, -1, 1},
	{&useAllow, skCatchHorse, useCallback("v_catch"), nil, nil, -1, 1},
	{&useAllow, skExtractVenom, nil, nil, nil, 7, 0},
	{&useAllow, skHarvestLumber, useCallback("v_wood"), nil, nil, -1, 1},
	{&useAllow, skHarvestYew, useCallback("v_yew"), nil, nil, -1, 1},
	{&useAllow, skAddRam, useCallback("v_add_ram"), useCallback("d_add_ram"), nil, 10, 0},
	{&useAllow, skSpyInv, useCallback("v_spy_inv"), useCallback("d_spy_inv"), nil, 7, 0},
	{&useAllow, skSpySkills, useCallback("v_spy_skills"), useCallback("d_spy_skills"), nil, 7, 0},
	{&useAllow, skSpyLord, useCallback("v_spy_lord"), useCallback("d_spy_lord"), nil, 7, 0},
	{&useAllow, skRecordSkill, useCallback("v_write_spell"), useCallback("d_write_spell"), nil, 7, 0},
	{&useAllow, skBribeNoble, useCallback("v_bribe"), useCallback("d_bribe"), nil, 7, 0},
	{&useAllow, skSummonSavage, useCallback("v_summon_savage"), nil, nil, 1, 0},
	{&useAllow, skKeepSavage, useCallback("v_keep_savage"), useCallback("d_keep_savage"), nil, 7, 0},
	{&useAllow, skImproveOpium, useCallback("v_improve_opium"), useCallback("d_improve_opium"), nil, 7, 0},
	{&useAllow, skRaiseMob, useCallback("v_raise"), useCallback("d_raise"), nil, 7, 0},
	{&useAllow, skRallyMob, useCallback("v_rally"), useCallback("d_rally"), nil, 7, 0},
	{&useAllow, skInciteMob, useCallback("v_incite"), useCallback("d_incite"), nil, 7, 0},
	{&useAllow, skBirdSpy, useCallback("v_bird_spy"), useCallback("d_bird_spy"), nil, 3, 0},
	{&useAllow, skLeadToGold, useCallback("v_lead_to_gold"), useCallback("d_lead_to_gold"), nil, 7, 0},
	{&useAllow, skRaiseCorpses, useCallback("v_raise_corpses"), nil, nil, -1, 1},
	{&useAllow, skUndeadLord, useCallback("v_undead_lord"), useCallback("d_undead_lord"), nil, 7, 0},
	{&useAllow, skBanishUndead, useCallback("v_banish_undead"), useCallback("d_banish_undead"), nil, 7, 0},
	{&useAllow, skRenewUndead, useCallback("v_keep_undead"), useCallback("d_keep_undead"), nil, 7, 0},
	{&useAllow, skEatDead, useCallback("v_eat_dead"), useCallback("d_eat_dead"), nil, 14, 0},
	{&useAllow, skAuraBlast, useCallback("v_aura_blast"), useCallback("d_aura_blast"), nil, 1, 0},
	{&useAllow, skAbsorbBlast, useCallback("v_aura_reflect"), nil, nil, 0, 0},
	{&useAllow, skSummonRain, useCallback("v_summon_rain"), useCallback("d_summon_rain"), nil, 7, 0},
	{&useAllow, skSummonWind, useCallback("v_summon_wind"), useCallback("d_summon_wind"), nil, 7, 0},
	{&useAllow, skSummonFog, useCallback("v_summon_fog"), useCallback("d_summon_fog"), nil, 7, 0},
	{&useAllow, skDirectStorm, useCallback("v_direct_storm"), nil, nil, 1, 0},
	{&useAllow, skRenewStorm, useCallback("v_renew_storm"), useCallback("d_renew_storm"), nil, 3, 0},
	{&useAllow, skDissipate, useCallback("v_dissipate"), useCallback("d_dissipate"), nil, 7, 0},
	{&useAllow, skLightning, useCallback("v_lightning"), useCallback("d_lightning"), nil, 7, 0},
	{&useAllow, skFierceWind, useCallback("v_fierce_wind"), useCallback("d_fierce_wind"), nil, 7, 0},
	{&useAllow, skSeizeStorm, useCallback("v_seize_storm"), useCallback("d_seize_storm"), nil, 7, 0},
	{&useAllow, skDeathFog, useCallback("v_death_fog"), useCallback("d_death_fog"), nil, 7, 0},
	{&useAllow, skBanishCorpses, useCallback("v_banish_corpses"), useCallback("d_banish_corpses"), nil, 7, 0},
	{&useAllow, skHideSelf, useCallback("v_hide"), useCallback("d_hide"), nil, 3, 0},
	{&useAllow, skSneakBuild, useCallback("v_sneak"), useCallback("d_sneak"), nil, 3, 0},
	{&useAllow, skMageMenial, useCallback("v_mage_menial"), nil, nil, -1, 1},
	{&useAllow, skPettyThief, useCallback("v_petty_thief"), useCallback("d_petty_thief"), nil, 7, 0},
	{&useAllow, skAppearCommon, useCallback("v_appear_common"), nil, nil, 1, 0},
	{&useAllow, skDefense, useCallback("v_defense"), useCallback("d_defense"), nil, 7, 0},
	{&useAllow, skArchery, useCallback("v_archery"), useCallback("d_archery"), nil, 7, 0},
	{&useAllow, skSwordplay, useCallback("v_swordplay"), useCallback("d_swordplay"), nil, 7, 0},
	{&useAllow, skRevealVision, useCallback("v_reveal_vision"), useCallback("d_reveal_vision"), nil, 10, 0},
	{&useAllow, skResurrect, useCallback("v_resurrect"), useCallback("d_resurrect"), nil, 10, 0},
	{&useAllow, skPray, useCallback("v_prep_ritual"), useCallback("d_prep_ritual"), nil, 3, 0},
	{&useAllow, skLastRites, useCallback("v_last_rites"), useCallback("d_last_rites"), nil, 10, 0},
	{&useAllow, skRemoveBless, useCallback("v_remove_bless"), useCallback("d_remove_bless"), nil, 10, 0},
	{&useAllow, skVisionProtect, useCallback("v_vision_protect"), useCallback("d_vision_protect"), nil, 10, 0},
	{&useAllow, skFindRich, useCallback("v_find_rich"), useCallback("d_find_rich"), nil, 7, 0},
	{&useAllow, skHarvestOpium, useCallback("v_implicit"), nil, nil, 0, 0},
	{&useAllow, skTrainAngry, useCallback("v_implicit"), nil, nil, 0, 0},
	{&useAllow, skWeaponsmith, useCallback("v_implicit"), nil, nil, 0, 0},
	{&useAllow, skHideLord, useCallback("v_implicit"), nil, nil, 0, 0},
	{&useAllow, skTranscendDeath, useCallback("v_implicit"), nil, nil, 0, 0},
	{&useAllow, skCollectFoliage, useCallback("v_implicit"), nil, nil, 0, 0},
	{&useAllow, skFishing, useCallback("v_fish"), nil, nil, 0, 0},
	{&useAllow, skSummonGhost, useCallback("v_implicit"), nil, nil, 0, 0},
	{&useAllow, skCaptureBeasts, useCallback("v_implicit"), nil, nil, 0, 0},
	{&useAllow, skUseBeasts, useCallback("v_implicit"), nil, nil, 0, 0},
	{&useAllow, skCollectElem, useCallback("v_implicit"), nil, nil, 0, 0},
	{&useAllow, skTorture, useCallback("v_torture"), useCallback("d_torture"), nil, 7, 0},
	{&useAllow, skFightToDeath, useCallback("v_fight_to_death"), nil, nil, 0, 0},
	{&useAllow, skBreedBeasts, useCallback("v_breed"), useCallback("d_breed"), nil, 7, 0},
	{&useAllow, skBreedHound, useCallback("v_breed_hound"), useCallback("d_breed_hound"), nil, 28, 0},
	{&useAllow, skPersuadeOath, useCallback("v_persuade_oath"), useCallback("d_persuade_oath"), nil, 7, 0},
	{&useAllow, skForgeWeapon, useCallback("v_forge_art_x"), useCallback("d_forge_art_x"), nil, 7, 0},
	{&useAllow, skForgeArmor, useCallback("v_forge_art_x"), useCallback("d_forge_art_x"), nil, 7, 0},
	{&useAllow, skForgeBow, useCallback("v_forge_art_x"), useCallback("d_forge_art_x"), nil, 7, 0},
	{&useAllow, skTrance, useCallback("v_trance"), useCallback("d_trance"), nil, 28, 0},
	{&useAllow, skTeleportItem, useCallback("v_teleport_item"), useCallback("d_teleport_item"), nil, 3, 0},
	{&useAllow, skTapHealth, useCallback("v_tap_health"), useCallback("d_tap_health"), nil, 7, 0},
	{&useAllow, skBindStorm, useCallback("v_bind_storm"), useCallback("d_bind_storm"), nil, 7, 0},
	{&useAllow, skFindSell, useCallback("v_find_sell"), useCallback("d_find_sell"), nil, 21, 0},
	{&useAllow, skFindBuy, useCallback("v_find_buy"), useCallback("d_find_buy"), nil, 14, 0},
	{},
}

func findUseEntry(skill int) int {
	for i := 1; useTbl[i].skill != 0; i++ {
		if useTbl[i].skill == skill {
			return i
		}
	}
	return -1
}

// Like C, callers must supply an index within the table, including its sentinel.
func useTblSkill(ent int) int { return useTbl[ent].skill }

func (w *world) mayUseSkill(who, sk int) int {
	if w.hasSkill(who, sk) > 0 {
		return sk
	}
	ret, scroll := 0, 0
	assert(w.validBox(who), "olympia/use.c", 303, "valid_box(who)")
	for e := range w.loopInv(who) {
		if p := w.rpItemMagic(int(e.item)); p != nil && p.mayUse.Lookup(sk) >= 0 {
			if w.subkind(int(e.item)) == subScroll {
				scroll = int(e.item)
			} else {
				ret = int(e.item)
			}
		}
	}
	if ret != 0 {
		return ret
	}
	return scroll
}

func (w *world) magicallySpeedCasting(c *command, sk int) {
	if w.subkind(w.skillSchool(sk)) != subMagic {
		return
	}
	if p := w.rpMagic(int(c.who)); p == nil || p.quickCast == 0 || sk == skSaveQuick || sk == skTrance {
		return
	}
	p := w.pMagic(int(c.who))
	n := int32(0)
	if c.wait == 0 {
		// Even a zero-duration cast logs, without consuming quick_cast.
	} else if int32(p.quickCast) < c.wait {
		n = int32(p.quickCast)
		c.wait -= n
		p.quickCast = 0
	} else {
		n = c.wait - 1
		p.quickCast = 0
		c.wait = 1
	}
	w.commandOut(int(c.who), "(speeded cast by %d day%s)", n, addS(int(n)))
}

func (w *world) correctUseItem(c *command) int {
	item := int(c.a)
	if p := w.rpItemMagic(item); p != nil && p.useKey != 0 {
		return item
	}
	if p := w.rpItemMagic(item); p != nil && len(p.mayUse) > 0 {
		c.a = int32(p.mayUse[0])
	}
	return int(c.a)
}

func (w *world) meetsRequirements(who, skill int) bool {
	p := w.rpSkill(skill)
	if p == nil {
		return true
	}
	l := p.req
	for i := 0; i < len(l); {
		for w.hasItem(who, int(l[i].item)) < l[i].qty && l[i].consume == reqOr {
			i++
			if i >= len(l) {
				fmt.Fprintf(w.stderr, "skill = %d\n", skill)
				assert(false, "olympia/use.c", 412, "FALSE")
			}
		}
		if w.hasItem(who, int(l[i].item)) < l[i].qty {
			w.commandOut(who, "%s does not have %s.", w.justName(who), w.boxNameQty(int(l[i].item), int(l[i].qty)))
			return false
		}
		// Preserve C's index-before-bound ordering; valid OR chains terminate
		// in REQ_YES or REQ_NO. Do not model undefined malformed-list reads.
		for l[i].consume == reqOr && i < len(l) {
			i++
		}
		i++
	}
	return true
}

func (w *world) consumeRequirements(who, skill int) {
	p := w.rpSkill(skill)
	if p == nil {
		return
	}
	l := p.req
	for i := 0; i < len(l); {
		for w.hasItem(who, int(l[i].item)) < l[i].qty && l[i].consume == reqOr {
			i++
			if i >= len(l) {
				fmt.Fprintf(w.stderr, "req list ends with REQ_OR, skill = %d\n", skill)
				assert(false, "olympia/use.c", 463, "FALSE")
			}
		}
		item, qty := l[i].item, l[i].qty
		for l[i].consume == reqOr && i < len(l) {
			i++
		}
		if l[i].consume == reqYes {
			w.consumeItem(who, int(item), int(qty))
		}
		i++
	}
}

func (w *world) consumeScroll(who, basis int) {
	if w.subkind(basis) == subScroll {
		w.commandOut(who, "%s vanishes.", w.boxName(basis))
		w.destroyUniqueItem(who, basis)
	}
}

func experienceUseSpeedup(c *command) {
	exp := max(c.useExp-1, 0)
	if exp != 0 && c.wait >= 7 {
		if c.wait >= 14 {
			c.wait -= exp
		} else if c.wait >= 10 {
			c.wait -= exp / 2
		} else if exp >= 2 {
			c.wait--
		}
	}
}

// Despite its comment, C increments only the supplied entry, irrespective of
// knowledge, no_exp or school; exp_this_month alone prevents another increment.
func (w *world) addSkillExperience(who, sk int) {
	if p := w.rpSkillEnt(who, sk); p != nil && p.expThisMonth == 0 {
		p.experience++
		p.expThisMonth = 1
	}
}

func expLevel(exp int) int {
	switch {
	case exp <= 4:
		return expNovice
	case exp <= 11:
		return expJourneyman
	case exp <= 20:
		return expTeacher
	case exp <= 34:
		return expMaster
	default:
		return expGrand
	}
}

// expS is complete use.c:977–994.
func expS(level int) string {
	switch level {
	case expNovice:
		return "apprentice"
	case expJourneyman:
		return "journeyman"
	case expTeacher:
		return "adept"
	case expMaster:
		return "master"
	case expGrand:
		return "grand master"
	default:
		assert(false, "olympia/use.c", 990, "FALSE")
	}
	return ""
}

// repSkillComp/flatSkillComp are use.c:1199–1238. The original comparator
// puts nonknown entries FIRST despite its comment; the renderer skips them.
func (w *world) repSkillComp(a, b *skillEnt) int {
	if a.know != skillKnow && b.know == skillKnow {
		return -1
	}
	if b.know != skillKnow && a.know == skillKnow {
		return 1
	}
	pa, pb := w.skillSchool(int(a.skill)), w.skillSchool(int(b.skill))
	if pa != pb {
		return int(int32(pa) - int32(pb))
	}
	return flatSkillComp(a, b)
}

func flatSkillComp(a, b *skillEnt) int { return int(a.skill - b.skill) }

// listSkillSup is the active definition at use.c:1262–1271.
func (w *world) listSkillSup(who int, e *skillEnt) {
	p := w.rpSkill(int(e.skill))
	if p != nil && p.noExp != 0 || w.skillSchool(int(e.skill)) == int(e.skill) {
		w.commandOut(who, "%s", w.boxName(int(e.skill)))
	} else {
		w.commandOut(who, "%s, %s", w.boxName(int(e.skill)), expS(expLevel(int(e.experience))))
	}
}

// listSkills is complete use.c:1276–1335. Only the temporary pointer list
// is sorted; the raw and serialized skill order must remain unchanged.
func (w *world) listSkills(who, num int) {
	assert(w.validBox(num), "olympia/use.c", 1283, "valid_box(num)")
	w.commandOut(who, "")
	w.commandOut(who, "Skills known:")
	w.indent += 3
	first := true
	if p := w.rpChar(num); p != nil && len(p.skills) > 0 {
		l := p.skills.Copy()
		safeQsort(l, w.repSkillComp)
		for _, e := range l {
			if e.know != skillKnow {
				continue
			}
			first = false
			if w.reqSkill(int(e.skill)) != 0 {
				w.indent += 6
				w.listSkillSup(who, e)
				w.indent -= 6
			} else {
				w.listSkillSup(who, e)
			}
		}
	}
	if first {
		w.commandOut(who, "none")
	}
	w.indent -= 3
}

func (w *world) fractionalSkillQualifier(p *skillEnt) string {
	assert(p.know != skillKnow, "olympia/use.c", 1349, "p->know != SKILL_know")
	if p.know == skillDont {
		return sout("0/%d%s", w.learnTime(int(p.skill)), w.npReqS(int(p.skill)))
	}
	assert(p.know == skillLearning, "olympia/use.c", 1358, "p->know == SKILL_learning")
	return sout("%d/%d", p.daysStudied, w.learnTime(int(p.skill)))
}

// listPartialSkills is complete use.c:1364–1406. Missing/empty/all-known
// lists produce no heading; assertion failures retain the first heading.
func (w *world) listPartialSkills(who, num int) {
	assert(w.validBox(num), "olympia/use.c", 1371, "valid_box(num)")
	p := w.rpChar(num)
	if p == nil || len(p.skills) < 1 {
		return
	}
	l := p.skills.Copy()
	safeQsort(l, flatSkillComp)
	first := true
	for _, e := range l {
		if e.know == skillKnow {
			continue
		}
		if first {
			w.commandOut(who, "")
			w.commandOut(who, "Partially known skills:")
			w.indent += 3
			first = false
		}
		w.wiout(who, 6, "%s, %s", w.boxName(int(e.skill)), w.fractionalSkillQualifier(e))
	}
	if !first {
		w.indent -= 3
	}
}

// charNPTotal counts every skill entry, not just known skills.
func (w *world) charNPTotal(who int) int {
	if w.isNPC(who) && w.subkind(who) != subDeadBody {
		return 0
	}
	sum := int32(1)
	assert(w.validBox(who), "olympia/use.c", 2126, "valid_box(who)")
	if p := w.rpChar(who); p != nil {
		if p.loyKind == loyOath {
			sum += p.loyRate
		}
		for _, e := range p.skills.Copy() {
			if sk := w.rpSkill(int(e.skill)); sk != nil {
				sum += sk.npReq
			}
		}
	}
	return int(sum)
}

func (w *world) rpSkillEnt(who, skill int) *skillEnt {
	if p := w.rpChar(who); p != nil {
		for _, e := range p.skills {
			if int(e.skill) == skill {
				return e
			}
		}
	}
	return nil
}

func (w *world) hasSkill(who, skill int) int {
	p := w.rpSkillEnt(who, skill)
	if p == nil || p.know != skillKnow {
		return 0
	}
	return expLevel(int(p.experience))
}

// forgetSkill removes an entry regardless of its learning state (use.c:1038–1067).
func (w *world) forgetSkill(who, skill int) bool {
	p := w.rpChar(who)
	if p == nil {
		return false
	}
	t := w.rpSkillEnt(who, skill)
	if t == nil {
		return false
	}
	p.skills.RemValue(t)
	if w.subkind(w.skillSchool(skill)) == subMagic {
		ch := w.pMagic(who)
		ch.maxAura--
		if ch.curAura > w.maxEffAura(who) && ch.curAura > 0 {
			ch.curAura--
		}
	}
	return true
}

// skillSchool is skill_school(): the school of skill sk, found by
// following its required skills up to one that requires none. A chain
// longer than 1000 (a loop) fails an assertion.
func (w *world) skillSchool(sk int) int {
	count := 0

	for {
		assert(count < 1000, "olympia/use.c", 1182, "count++ < 1000")
		count++

		n := w.reqSkill(sk)
		if n == 0 {
			return sk
		}
		sk = n
	}
}
