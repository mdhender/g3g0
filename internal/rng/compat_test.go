// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package rng_test

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/mdhender/g3g0/internal/rnd"
	"github.com/mdhender/g3g0/internal/rng"
)

// TestMatchesC runs testdata/cases.txt the way the C driver in
// tools/rng-compat/run.sh does and requires its output, line for line.
// testdata/c-output.txt is that driver's output from the pinned G3 build;
// regenerate it only with run.sh --update.
func TestMatchesC(t *testing.T) {
	cases, err := os.ReadFile(filepath.Join("testdata", "cases.txt"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "c-output.txt"))
	if err != nil {
		t.Fatal(err)
	}
	got := runCases(t, cases)

	wantLines := strings.Split(string(want), "\n")
	gotLines := strings.Split(string(got), "\n")
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

// runCases is the Go side of the driver: same operations, same output.
func runCases(t *testing.T, cases []byte) []byte {
	t.Helper()
	dir := t.TempDir()
	var out bytes.Buffer

	src := rnd.New([16]byte{})
	var turn, day int
	z := rng.NewStreams(src.MasterSeed, func() (int, int) { return turn, day })
	var cur rng.Stream

	stream := func(s rng.Stream) string {
		return fmt.Sprintf("seed %08x %08x %08x %08x ctr %d", s.Seed[0], s.Seed[1], s.Seed[2], s.Seed[3], s.Counter)
	}
	master := func() string {
		m := src.MasterSeed()
		return fmt.Sprintf("master %08x %08x %08x %08x", m[0], m[1], m[2], m[3])
	}
	list := func(l []int) string {
		var b strings.Builder
		for i, v := range l {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprint(&b, v)
		}
		return b.String()
	}
	makeList := func(n int) []int {
		l := make([]int, n)
		for i := range l {
			l[i] = 100 + i
		}
		return l
	}
	helpers := streamHelpers(z, makeList, list)

	sc := bufio.NewScanner(bytes.NewReader(cases))
	for sc.Scan() {
		line, _, _ := strings.Cut(sc.Text(), "#")
		f := strings.Fields(line)
		if len(f) == 0 {
			continue
		}
		// a holds the integer fields as C's (int)strtol(s, NULL, 0) reads
		// them; u holds them as (uint32_t)strtoul(s, NULL, 0).
		a := make([]int, len(f))
		u := make([]uint32, len(f))
		for i := 1; i < len(f); i++ {
			if v, err := strconv.ParseInt(f[i], 0, 64); err == nil {
				a[i-1] = int(int32(v))
				u[i-1] = uint32(v)
			} else if v, err := strconv.ParseUint(f[i], 0, 64); err == nil {
				a[i-1] = int(int32(v))
				u[i-1] = uint32(v)
			}
		}

		var res string
		switch f[0] {
		case "master":
			var b []byte
			if len(f) > 1 {
				var err error
				if b, err = hex.DecodeString(f[1]); err != nil {
					t.Fatalf("master %q: %v", f[1], err)
				}
			}
			path := filepath.Join(dir, "randseed")
			if err := os.WriteFile(path, b, 0o644); err != nil {
				t.Fatal(err)
			}
			if err := src.Load(path); err != nil {
				t.Fatal(err)
			}
			res = master()
		case "loadmissing":
			res = ""
			if err := src.Load(filepath.Join(dir, "missing-randseed")); err != nil {
				res = "missing-randseed could not be opened.\n"
			}
			res += master()
		case "rnd":
			res = strconv.Itoa(src.Rnd(a[0], a[1]))
		case "md5int":
			res = strconv.Itoa(rnd.MD5Int(a[0], a[1], a[2], a[3]))
		case "clock":
			turn, day = int(int16(a[0])), int(int16(a[1]))
			res = "ok"
		case "root":
			cur = rng.Seed(src.MasterSeed())
			res = stream(cur)
		case "of":
			cur = cur.Of(a[0], u[1])
			res = stream(cur)
		case "draw":
			v := cur.Draw(a[0], a[1])
			res = strconv.Itoa(v) + " " + stream(cur)
		case "keyed":
			res = strconv.Itoa(cur.Keyed(a[0], a[1], u[2], a[3], a[4]))
		case "shuffle":
			l := makeList(a[0])
			rng.Shuffle(l, &cur)
			res = fmt.Sprintf("ctr %d: %s", cur.Counter, list(l))
		case "call":
			h, ok := helpers[f[1]]
			if !ok {
				t.Fatalf("unknown helper %q", f[1])
			}
			res = h(a[1:])
		default:
			t.Fatalf("unknown op %q", f[0])
		}
		fmt.Fprintf(&out, "%s -> %s\n", strings.Join(f, " "), res)
	}
	return out.Bytes()
}

// streamHelpers maps each C helper name the driver accepts to its port.
func streamHelpers(z *rng.Streams, makeList func(int) []int, list func([]int) string) map[string]func(x []int) string {
	i := strconv.Itoa
	shuffle := func(f func([]int)) func([]int) string {
		return func(x []int) string { l := makeList(x[0]); f(l); return list(l) }
	}
	return map[string]func(x []int) string{
		"begin_battle":  func(x []int) string { z.BeginBattle(x[0], x[1], x[2]); return "ok" },
		"crnd":          func(x []int) string { return i(z.Crnd(x[0], x[1])) },
		"begin_pillage": func(x []int) string { z.BeginPillage(x[0], x[1]); return "ok" },
		"prnd":          func(x []int) string { return i(z.Prnd(x[0], x[1])) },
		"pillage_mob":   func(x []int) string { return i(z.Pillage().Draw(1, 2)) },

		"begin_economy": func(x []int) string { z.BeginEconomy(x[0]); return "ok" },
		"econ_pick":     func(x []int) string { return i(z.EconPick(x[0], x[1], x[2], x[3])) },
		"econ_stock":    func(x []int) string { return i(z.EconStock(x[0], x[1], x[2], x[3])) },
		"econ_qty":      func(x []int) string { return i(z.EconQty(x[0], x[1], x[2], x[3])) },
		"econ_cost":     func(x []int) string { return i(z.EconCost(x[0], x[1], x[2], x[3])) },
		"econ_expire":   func(x []int) string { return i(z.EconExpire(x[0], x[1], x[2], x[3])) },
		"econ_mine":     func(x []int) string { return i(z.EconMine(x[0], x[1], x[2], x[3])) },
		"econ_harvest":  func(x []int) string { return i(z.EconHarvest(x[0], x[1], x[2], x[3])) },
		"econ_ring":     func(x []int) string { return i(z.EconRing(x[0], x[1], x[2], x[3])) },

		"begin_npc":    func(x []int) string { z.BeginNPC(x[0]); return "ok" },
		"npc_spawn":    func(x []int) string { return i(z.NPCSpawn(x[0], x[1], x[2], x[3])) },
		"npc_qty":      func(x []int) string { return i(z.NPCQty(x[0], x[1], x[2], x[3])) },
		"npc_behavior": func(x []int) string { return i(z.NPCBehavior(x[0], x[1], x[2], x[3])) },
		"npc_victim":   func(x []int) string { return i(z.NPCVictim(x[0], x[1])) },

		"begin_quest":   func(x []int) string { z.BeginQuest(x[0], x[1]); return "ok" },
		"qrnd":          func(x []int) string { return i(z.Qrnd(x[0], x[1])) },
		"quest_shuffle": shuffle(func(l []int) { rng.Shuffle(l, z.Quest()) }),

		"begin_cloud_build":  func(x []int) string { z.BeginCloudBuild(x[0]); return "ok" },
		"cseq_rnd":           func(x []int) string { return i(z.CseqRnd(x[0], x[1])) },
		"cseq_shuffle":       shuffle(z.CseqShuffle),
		"begin_faery_build":  func(x []int) string { z.BeginFaeryBuild(x[0]); return "ok" },
		"fseq_rnd":           func(x []int) string { return i(z.FseqRnd(x[0], x[1])) },
		"fseq_shuffle":       shuffle(z.FseqShuffle),
		"begin_hades_build":  func(x []int) string { z.BeginHadesBuild(x[0]); return "ok" },
		"hseq_rnd":           func(x []int) string { return i(z.HseqRnd(x[0], x[1])) },
		"hseq_shuffle":       shuffle(z.HseqShuffle),
		"begin_worldgen_loc": func(x []int) string { z.BeginWorldgenLoc(x[0]); return "ok" },
		"wseq_rnd":           func(x []int) string { return i(z.WseqRnd(x[0], x[1])) },
		"wseq_shuffle":       shuffle(z.WseqShuffle),

		"magc_scry":   func(x []int) string { return i(z.MagcScry(x[0], x[1])) },
		"magc_piety":  func(x []int) string { return i(z.MagcPiety(x[0], x[1])) },
		"magc_eat":    func(x []int) string { return i(z.MagcEat(x[0], x[1], x[2])) },
		"magc_learn":  func(x []int) string { return i(z.MagcLearn(x[0], x[1])) },
		"magc_potion": func(x []int) string { return i(z.MagcPotion(x[0], x[1], x[2], x[3], x[4])) },
		"magc_med":    func(x []int) string { return i(z.MagcMed(x[0], x[1])) },
		"magc_omen":   func(x []int) string { return i(z.MagcOmen(x[0], x[1])) },
		"magc_heal":   func(x []int) string { return i(z.MagcHeal(x[0], x[1])) },
		"magc_forge":  func(x []int) string { return i(z.MagcForge(x[0], x[1])) },
		"magc_orb":    func(x []int) string { return i(z.MagcOrb(x[0], x[1])) },
		"magc_ring":   func(x []int) string { return i(z.MagcRing(x[0], x[1])) },

		"expl_find":   func(x []int) string { return i(z.ExplFind(x[0], x[1])) },
		"expl_gate":   func(x []int) string { return i(z.ExplGate(x[0], x[1])) },
		"expl_flavor": func(x []int) string { return i(z.ExplFlavor(x[0], x[1])) },
		"expl_pick":   func(x []int) string { return i(z.ExplPick(x[0], x[1], x[2])) },
		"expl_seek":   func(x []int) string { return i(z.ExplSeek(x[0], x[1])) },
		"expl_detect": func(x []int) string { return i(z.ExplDetect(x[0], x[1])) },

		"mint_alloc":    func(x []int) string { return i(z.MintAlloc(x[0], x[1])) },
		"mint_password": func(x []int) string { return i(z.MintPassword(x[0], x[1], x[2])) },
		"mint_city":     func(x []int) string { return i(z.MintCity()) },
		"mint_shuffle":  shuffle(z.MintShuffle),

		"exitdir_shuffle": shuffle(func(l []int) { z.BeginExitdir(); rng.Shuffle(l, z.Exitdir()) }),

		"up_heal":   func(x []int) string { return i(z.UpHeal(x[0], x[1], x[2], x[3])) },
		"up_loyal":  func(x []int) string { return i(z.UpLoyal(x[0], x[1], x[2], x[3])) },
		"up_starve": func(x []int) string { return i(z.UpStarve(x[0], x[1], x[2], x[3])) },
		"up_animal": func(x []int) string { return i(z.UpAnimal(x[0], x[1], x[2], x[3], x[4])) },
		"up_corpse": func(x []int) string { return i(z.UpCorpse(x[0], x[1], x[2])) },
		"up_income": func(x []int) string { return i(z.UpIncome(x[0], x[1], x[2], x[3])) },
		"cal_day":   func(x []int) string { return i(z.CalDay(x[0], x[1], x[2])) },

		"faer_hunt_loc":    func(x []int) string { return i(z.FaerHuntLoc(x[0], x[1], x[2], x[3])) },
		"faer_hunt_qty":    func(x []int) string { return i(z.FaerHuntQty(x[0])) },
		"faer_bandit_kind": func(x []int) string { return i(z.FaerBanditKind(x[0], x[1])) },
		"faer_bandit_qty":  func(x []int) string { return i(z.FaerBanditQty(x[0])) },
		"faer_bandit_gold": func(x []int) string { return i(z.FaerBanditGold(x[0])) },
		"faer_ambush":      func(x []int) string { return i(z.FaerAmbush(x[0], x[1])) },
		"faer_retal":       func(x []int) string { return i(z.FaerRetal(x[0], x[1])) },

		"hads_nasty":       func(x []int) string { return i(z.HadsNasty(x[0], x[1], x[2], x[3])) },
		"hads_bandit_kind": func(x []int) string { return i(z.HadsBanditKind(x[0], x[1])) },
		"hads_bandit_qty":  func(x []int) string { return i(z.HadsBanditQty(x[0])) },
		"hads_ambush":      func(x []int) string { return i(z.HadsAmbush(x[0], x[1])) },
		"hads_retal":       func(x []int) string { return i(z.HadsRetal(x[0], x[1])) },
		"hads_random_loc":  func(x []int) string { return i(z.HadsRandomLoc(x[0], x[1])) },

		"wgen_prom":  func(x []int) string { return i(z.WgenProm(x[0])) },
		"wgen_skill": func(x []int) string { return i(z.WgenSkill(x[0], x[1], x[2], x[3])) },
		"wgen_garr":  func(x []int) string { return i(z.WgenGarr(x[0])) },
		"wgen_gate":  func(x []int) string { return i(z.WgenGate(x[0])) },

		"wthr_shuffle": shuffle(z.WthrShuffle),
		"wthr_storm":   func(x []int) string { return i(z.WthrStorm(x[0], x[1], x[2], x[3])) },
		"wthr_wreck":   func(x []int) string { return i(z.WthrWreck(x[0], x[1], x[2], x[3])) },

		"soc_gift":     func(x []int) string { return i(z.SocGift(x[0], x[1])) },
		"soc_bribe":    func(x []int) string { return i(z.SocBribe(x[0], x[1], x[2], x[3], x[4])) },
		"soc_terror":   func(x []int) string { return i(z.SocTerror(x[0], x[1])) },
		"soc_incite":   func(x []int) string { return i(z.SocIncite(x[0], x[1], x[2], x[3], x[4])) },
		"soc_persuade": func(x []int) string { return i(z.SocPersuade(x[0], x[1])) },
		"soc_breed":    func(x []int) string { return i(z.SocBreed(x[0], x[1], x[2], x[3], x[4])) },

		"skil_crit":          func(x []int) string { return i(z.SkilCrit(x[0], x[1])) },
		"skil_bonus":         func(x []int) string { return i(z.SkilBonus(x[0], x[1], x[2], x[3])) },
		"skil_study":         func(x []int) string { return i(z.SkilStudy(x[0], x[1])) },
		"skil_research":      func(x []int) string { return i(z.SkilResearch(x[0], x[1])) },
		"skil_research_pick": func(x []int) string { return i(z.SkilResearchPick(x[0], x[1], x[2])) },
		"skil_torture":       func(x []int) string { return i(z.SkilTorture(x[0], x[1])) },
		"skil_petty":         func(x []int) string { return i(z.SkilPetty(x[0], x[1], x[2], x[3], x[4])) },

		"ent_prisoner":  func(x []int) string { return i(z.EntPrisoner(x[0], x[1])) },
		"ent_drop":      func(x []int) string { return i(z.EntDrop(x[0], x[1], x[2], x[3], x[4])) },
		"ent_take":      func(x []int) string { return i(z.EntTake(x[0], x[1], x[2], x[3], x[4])) },
		"ent_sick":      func(x []int) string { return i(z.EntSick(x[0], x[1])) },
		"ent_land_dir":  func(x []int) string { return i(z.EntLandDir(x[0], x[1])) },
		"ent_land_pick": func(x []int) string { return i(z.EntLandPick(x[0], x[1])) },
		"ent_bark":      func(x []int) string { return i(z.EntBark(x[0], x[1])) },
		"ent_build":     func(x []int) string { return i(z.EntBuild(x[0])) },
		"ent_menial":    func(x []int) string { return i(z.EntMenial(x[0])) },
	}
}
