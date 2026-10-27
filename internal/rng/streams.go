// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package rng

// Purpose tags. Each is four ASCII bytes read as a big-endian word, written
// in C as a hex literal with the string in a comment. Names follow the C
// macros; several files define the same value under their own name.
const (
	tagTurn = 0x7475726e // "turn"
	tagDay  = 0x64617920 // "day "

	// combat.c
	tagCombat  = 0x636f6d62 // "comb"
	tagPillage = 0x70696c6c // "pill"

	// buy.c
	tagEconomy  = 0x65636f6e // "econ"
	etagPick    = 0x7069636b // "pick"
	etagStock   = 0x73746f6b // "stok"
	etagQty     = 0x716e7479 // "qnty"
	etagCost    = 0x636f7374 // "cost"
	etagExpire  = 0x65787072 // "expr"
	etagMine    = 0x6d696e65 // "mine"
	etagHarvest = 0x68617276 // "harv"
	etagRing    = 0x72696e67 // "ring"

	// npc.c
	tagNPC      = 0x6e706373 // "npcs"
	ntagSpawn   = 0x7370776e // "spwn"
	ntagQty     = 0x716e7479 // "qnty"
	ntagBehave  = 0x62687672 // "bhvr"
	ntagVictim  = 0x76696374 // "vict"
	tagQuest    = 0x71657374 // "qest" (quest.c)
	tagCloud    = 0x636c7564 // "clud" (cloud.c)
	tagExitdir  = 0x65786472 // "exdr" (dir.c)
	tagMint     = 0x6d696e74 // "mint" (code.c)
	mtagPW      = 0x70617373 // "pass" (code.c)
	tagWorldgen = 0x7767656e // "wgen" (seed.c, tunnel.c)

	// seed.c
	wtagProm = 0x70726f6d // "prom"
	wtagTech = 0x74656368 // "tech"
	wtagGarr = 0x67617272 // "garr"
	wtagGate = 0x67617465 // "gate"

	// storm.c
	tagWeather = 0x77746872 // "wthr"
	wtagStorm  = 0x73746f6d // "stom"
	wtagWreck  = 0x7772636b // "wrck"

	// day.c
	tagUpkeep   = 0x75706b70 // "upkp"
	utagHeal    = 0x6865616c // "heal"
	utagLoyal   = 0x6c6f796c // "loyl"
	utagSicken  = 0x73636b6e // "sckn"
	utagAnimal  = 0x64696572 // "dier"
	utagCorpse  = 0x636f7270 // "corp"
	utagIncome  = 0x696e636f // "inco"
	tagCalendar = 0x63616c6e // "caln"
	ctagDay     = 0x64617920 // "day "

	// c1.c
	tagExplore    = 0x6578706c // "expl"
	xtagFind      = 0x66696e64 // "find"
	xtagGate      = 0x67617465 // "gate"
	xtagFlavor    = 0x666c6176 // "flav"
	xtagPick      = 0x7069636b // "pick"
	xtagSeek      = 0x7365656b // "seek"
	xtagDetect    = 0x64746374 // "dtct"
	tagSkills     = 0x736b696c // "skil" (use.c)
	stagCrit      = 0x63726974 // "crit"
	stagYield     = 0x7969656c // "yiel"
	stagStudy     = 0x73746479 // "stdy"
	stagResch     = 0x72736368 // "rsch"
	stagRPick     = 0x7270696b // "rpik"
	stagTort      = 0x746f7274 // "tort"
	stagPetty     = 0x70747479 // "ptty"
	tagMagic      = 0x6d616763 // "magc" (basic.c)
	mtagScry      = 0x73637279 // "scry"
	mtagPiety     = 0x70696574 // "piet"
	mtagEatD      = 0x65617464 // "eatd"
	mtagLearn     = 0x6c65726e // "lern"
	mtagMed       = 0x6d656469 // "medi"
	mtagOmen      = 0x6f6d656e // "omen"
	mtagHeal      = 0x6865616c // "heal"
	mtagPotn      = 0x706f746e // "potn"
	mtagForg      = 0x666f7267 // "forg"
	mtagOrb       = 0x6f726220 // "orb "
	mtagRing      = 0x72696e67 // "ring"
	tagSocial     = 0x736f636c // "socl" (swear.c)
	sotagGift     = 0x67696674 // "gift"
	sotagBribe    = 0x62726962 // "brib"
	sotagTerror   = 0x74657272 // "terr"
	sotagIncite   = 0x696e6369 // "inci"
	sotagOath     = 0x6f617468 // "oath"
	sotagBreed    = 0x62726564 // "bred"
	tagEntity     = 0x656e7479 // "enty" (u.c)
	entagPrison   = 0x70726973 // "pris"
	entagDrop     = 0x64726f70 // "drop"
	entagTake     = 0x74616b65 // "take"
	entagSick     = 0x7369636b // "sick"
	entagLandDir  = 0x6c646972 // "ldir"
	entagLandPick = 0x6c706b6b // "lpkk"
	entagBark     = 0x6261726b // "bark"
	entagBuild    = 0x626c6467 // "bldg"
	entagMenial   = 0x6d656e69 // "meni"

	// faery.c
	tagFaery = 0x66616572 // "faer"
	ftagHunt = 0x68756e74 // "hunt"
	ftagHQty = 0x68717479 // "hqty"
	ftagBKnd = 0x626b6e64 // "bknd"
	ftagBQty = 0x62717479 // "bqty"
	ftagBGld = 0x62676c64 // "bgld"
	ftagAmbs = 0x616d6273 // "ambs"
	ftagRtal = 0x7274616c // "rtal"
	tagHades = 0x68616473 // "hads" (hades.c)
	htagNast = 0x6e617374 // "nast"
	htagBKnd = 0x626b6e64 // "bknd"
	htagBQty = 0x62717479 // "bqty"
	htagAmbs = 0x616d6273 // "ambs"
	htagRtal = 0x7274616c // "rtal"
	htagRLoc = 0x726c6f63 // "rloc"
)

// Streams holds the engine's named streams: the file-scope rng_stream
// variables in olympia/*.c, with their begin_*() seeders and the leaf
// helpers that draw from them. Every stream is rooted at the master seed and
// the current turn, which the engine supplies.
//
// The zero value is not ready; use NewStreams. As in C, a helper that does
// not seed its own stream (the econ_*, npc_*, crnd/prnd, qrnd and *seq_*
// helpers) draws from whatever its begin_*() last left there, starting from
// an all-zero stream.
type Streams struct {
	master func() [4]uint32       // rng_master_seed()
	clock  func() (turn, day int) // sysclock.turn, sysclock.day

	// Reseeded by every call to their begin_*().
	economy, npc, combat, pillage, quest      Stream
	cloudSeq, faerySeq, hadesSeq, worldgenSeq Stream

	// Seeded at most once per turn, like the C statics guarded by a
	// <name>_rng_turn variable that starts at -1.
	magic, upkeep, calendar, faery, hades, mint perTurn
	exitdir, worldgen, weather, social, skills  perTurn
	entity                                      perTurn

	// Seeded at most once per (turn, day).
	explore perTurn
}

// perTurn is a turn-guarded stream: the stream and the turn (and, for
// explore, the day) it was last seeded for.
type perTurn struct {
	s         Stream
	turn, day int
}

// NewStreams returns the named streams. master is rng_master_seed(), for
// example (*rnd.Source).MasterSeed; clock returns sysclock's turn and day.
// Both are called at every seeding, as C reads them each time.
func NewStreams(master func() [4]uint32, clock func() (turn, day int)) *Streams {
	z := &Streams{master: master, clock: clock}
	for _, p := range []*perTurn{
		&z.magic, &z.upkeep, &z.calendar, &z.faery, &z.hades, &z.mint,
		&z.exitdir, &z.worldgen, &z.weather, &z.social, &z.skills, &z.entity,
		&z.explore,
	} {
		p.turn, p.day = -1, -1
	}
	return z
}

// turnRoot is the shared head of every begin_*(): the master seed's child
// for the current turn.
func (z *Streams) turnRoot() Stream {
	root := Seed(z.master())
	turn, _ := z.clock()
	return root.Of(turn, tagTurn)
}

// seedPerTurn is the body of a turn-guarded begin_*(): seed p as the turn
// root's child (0, tag), unless it was already seeded this turn.
func (z *Streams) seedPerTurn(p *perTurn, tag uint32) {
	turn, _ := z.clock()
	if p.turn == turn {
		return
	}
	root := z.turnRoot()
	p.s = root.Of(0, tag)
	p.turn = turn
}

// today is the <prefix>_today() helpers in basic.c, day.c, faery.c,
// hades.c, swear.c, use.c and u.c: a keyed leaf drawn from the per-day child
// of s, so a repeated action rolls afresh each day (#132, #167, #186).
func (z *Streams) today(s *Stream, k1, k2 int, tag uint32, low, high int) int {
	_, d := z.clock()
	day := s.Of(d, tag)
	return day.Keyed(k1, k2, tag, low, high)
}

// ---------------------------------------------------------------------------
// combat.c: per-battle and per-pillage sequential streams.
// ---------------------------------------------------------------------------

// BeginBattle is begin_battle(): reseed the combat stream for one battle.
func (z *Streams) BeginBattle(where, attacker, defender int) {
	turn := z.turnRoot()
	s := turn.Of(where, tagCombat)
	s = s.Of(attacker, tagCombat)
	s = s.Of(defender, tagCombat)
	_, day := z.clock()
	z.combat = s.Of(day, tagDay)
}

// Crnd is crnd(), a sequential draw from the combat stream.
func (z *Streams) Crnd(low, high int) int {
	return z.combat.Draw(low, high)
}

// BeginPillage is begin_pillage(): reseed the pillage stream.
func (z *Streams) BeginPillage(where, who int) {
	turn := z.turnRoot()
	s := turn.Of(where, tagPillage)
	s = s.Of(who, tagPillage)
	_, day := z.clock()
	z.pillage = s.Of(day, tagDay)
}

// Prnd is prnd(), a sequential draw from the pillage stream.
func (z *Streams) Prnd(low, high int) int {
	return z.pillage.Draw(low, high)
}

// Pillage is &pillage_rng, which d_pillage hands to create_peasant_mob().
func (z *Streams) Pillage() *Stream {
	return &z.pillage
}

// ---------------------------------------------------------------------------
// buy.c: the per-market economy stream. Not turn-guarded; the econ_*
// helpers do not seed it.
// ---------------------------------------------------------------------------

// BeginEconomy is begin_economy().
func (z *Streams) BeginEconomy(where int) {
	turn := z.turnRoot()
	z.economy = turn.Of(where, tagEconomy)
}

// EconPick is econ_pick().
func (z *Streams) EconPick(item, where, low, high int) int {
	return z.economy.Keyed(item, where, etagPick, low, high)
}

// EconStock is econ_stock().
func (z *Streams) EconStock(item, where, low, high int) int {
	return z.economy.Keyed(item, where, etagStock, low, high)
}

// EconQty is econ_qty().
func (z *Streams) EconQty(item, where, low, high int) int {
	return z.economy.Keyed(item, where, etagQty, low, high)
}

// EconCost is econ_cost().
func (z *Streams) EconCost(item, where, low, high int) int {
	return z.economy.Keyed(item, where, etagCost, low, high)
}

// EconExpire is econ_expire().
func (z *Streams) EconExpire(item, where, low, high int) int {
	return z.economy.Keyed(item, where, etagExpire, low, high)
}

// econUnitDay is econ_unit_day(): a leaf on a per-unit child of the market
// stream, keyed on the item and the day (#132).
func (z *Streams) econUnitDay(item, who int, tag uint32, low, high int) int {
	unit := z.economy.Of(who, tag)
	_, day := z.clock()
	return unit.Keyed(item, day, tag, low, high)
}

// EconMine is econ_mine().
func (z *Streams) EconMine(item, who, low, high int) int {
	return z.econUnitDay(item, who, etagMine, low, high)
}

// EconHarvest is econ_harvest().
func (z *Streams) EconHarvest(item, who, low, high int) int {
	return z.econUnitDay(item, who, etagHarvest, low, high)
}

// EconRing is econ_ring().
func (z *Streams) EconRing(item, where, low, high int) int {
	return z.economy.Keyed(item, where, etagRing, low, high)
}

// ---------------------------------------------------------------------------
// npc.c: the per-location NPC stream. Not turn-guarded; the npc_* helpers
// do not seed it.
// ---------------------------------------------------------------------------

// BeginNPC is begin_npc().
func (z *Streams) BeginNPC(where int) {
	turn := z.turnRoot()
	z.npc = turn.Of(where, tagNPC)
}

// NPCSpawn is npc_spawn().
func (z *Streams) NPCSpawn(k1, k2, low, high int) int {
	return z.npc.Keyed(k1, k2, ntagSpawn, low, high)
}

// NPCQty is npc_qty().
func (z *Streams) NPCQty(k1, k2, low, high int) int {
	return z.npc.Keyed(k1, k2, ntagQty, low, high)
}

// NPCBehavior is npc_behavior().
func (z *Streams) NPCBehavior(k1, k2, low, high int) int {
	return z.npc.Keyed(k1, k2, ntagBehave, low, high)
}

// NPCVictim is the auto_bandit victim pick, written inline in C: an index
// into a target list of length n, keyed on the bandit.
func (z *Streams) NPCVictim(who, n int) int {
	return z.npc.Keyed(who, 0, ntagVictim, 0, n-1)
}

// ---------------------------------------------------------------------------
// quest.c: a fresh sequential stream per quest scenario.
// ---------------------------------------------------------------------------

// BeginQuest is begin_quest().
func (z *Streams) BeginQuest(key, who int) {
	turn := z.turnRoot()
	s := turn.Of(key, tagQuest)
	s = s.Of(who, tagQuest)
	_, day := z.clock()
	z.quest = s.Of(day, tagDay)
}

// Qrnd is qrnd().
func (z *Streams) Qrnd(low, high int) int {
	return z.quest.Draw(low, high)
}

// Quest is &quest_rng, which quest.c shuffles candidate lists with.
func (z *Streams) Quest() *Stream {
	return &z.quest
}

// ---------------------------------------------------------------------------
// Region builds (cloud.c, faery.c, hades.c) and dungeon generation
// (tunnel.c): a fresh sequential stream per build.
// ---------------------------------------------------------------------------

// BeginCloudBuild is begin_cloud_build().
func (z *Streams) BeginCloudBuild(where int) {
	turn := z.turnRoot()
	z.cloudSeq = turn.Of(where, tagCloud)
}

// CseqRnd is cseq_rnd().
func (z *Streams) CseqRnd(low, high int) int { return z.cloudSeq.Draw(low, high) }

// CseqShuffle is cseq_shuffle().
func (z *Streams) CseqShuffle(l []int) { Shuffle(l, &z.cloudSeq) }

// BeginFaeryBuild is begin_faery_build().
func (z *Streams) BeginFaeryBuild(where int) {
	turn := z.turnRoot()
	z.faerySeq = turn.Of(where, tagFaery)
}

// FseqRnd is fseq_rnd().
func (z *Streams) FseqRnd(low, high int) int { return z.faerySeq.Draw(low, high) }

// FseqShuffle is fseq_shuffle().
func (z *Streams) FseqShuffle(l []int) { Shuffle(l, &z.faerySeq) }

// BeginHadesBuild is begin_hades_build().
func (z *Streams) BeginHadesBuild(where int) {
	turn := z.turnRoot()
	z.hadesSeq = turn.Of(where, tagHades)
}

// HseqRnd is hseq_rnd().
func (z *Streams) HseqRnd(low, high int) int { return z.hadesSeq.Draw(low, high) }

// HseqShuffle is hseq_shuffle().
func (z *Streams) HseqShuffle(l []int) { Shuffle(l, &z.hadesSeq) }

// BeginWorldgenLoc is begin_worldgen_loc() in tunnel.c.
func (z *Streams) BeginWorldgenLoc(where int) {
	turn := z.turnRoot()
	z.worldgenSeq = turn.Of(where, tagWorldgen)
}

// WseqRnd is wseq_rnd().
func (z *Streams) WseqRnd(low, high int) int { return z.worldgenSeq.Draw(low, high) }

// WseqShuffle is wseq_shuffle().
func (z *Streams) WseqShuffle(l []int) { Shuffle(l, &z.worldgenSeq) }

// ---------------------------------------------------------------------------
// basic.c: the per-turn magic stream.
// ---------------------------------------------------------------------------

func (z *Streams) beginMagic() { z.seedPerTurn(&z.magic, tagMagic) }

// MagcScry is magc_scry().
func (z *Streams) MagcScry(who, k2 int) int {
	z.beginMagic()
	return z.today(&z.magic.s, who, k2, mtagScry, 1, 100)
}

// MagcPiety is magc_piety().
func (z *Streams) MagcPiety(who, k2 int) int {
	z.beginMagic()
	return z.today(&z.magic.s, who, k2, mtagPiety, 1, 100)
}

// MagcEat is magc_eat().
func (z *Streams) MagcEat(who, body, sub int) int {
	z.beginMagic()
	return z.magic.s.Keyed(who, body*2+sub, mtagEatD, 1, 100)
}

// MagcLearn is magc_learn().
func (z *Streams) MagcLearn(who, skill int) int {
	z.beginMagic()
	return z.magic.s.Keyed(who, skill, mtagLearn, 1, 100)
}

// MagcPotion is magc_potion().
func (z *Streams) MagcPotion(who, item, sub, low, high int) int {
	z.beginMagic()
	return z.today(&z.magic.s, who, item*4+sub, mtagPotn, low, high)
}

// MagcMed is magc_med().
func (z *Streams) MagcMed(who, sub int) int {
	z.beginMagic()
	return z.today(&z.magic.s, who, sub, mtagMed, 1, 100)
}

// MagcOmen is magc_omen().
func (z *Streams) MagcOmen(who, other int) int {
	z.beginMagic()
	return z.magic.s.Keyed(who, other, mtagOmen, 1, 4)
}

// MagcHeal is magc_heal().
func (z *Streams) MagcHeal(who, target int) int {
	z.beginMagic()
	return z.today(&z.magic.s, who, target, mtagHeal, 1, 100)
}

// MagcForge is magc_forge().
func (z *Streams) MagcForge(who, sub int) int {
	z.beginMagic()
	return z.today(&z.magic.s, who, sub, mtagForg, 1, 3)
}

// MagcOrb is magc_orb().
func (z *Streams) MagcOrb(who, item int) int {
	z.beginMagic()
	return z.today(&z.magic.s, who, item, mtagOrb, 1, 3)
}

// MagcRing is magc_ring().
func (z *Streams) MagcRing(who, item int) int {
	z.beginMagic()
	return z.today(&z.magic.s, who, item, mtagRing, 1, 3)
}

// ---------------------------------------------------------------------------
// c1.c: the explore stream, seeded once per (turn, day) (#167).
// ---------------------------------------------------------------------------

func (z *Streams) beginExplore() {
	turn, day := z.clock()
	if z.explore.turn == turn && z.explore.day == day {
		return
	}
	root := z.turnRoot()
	s := root.Of(0, tagExplore)
	z.explore.s = s.Of(day, tagDay)
	z.explore.turn, z.explore.day = turn, day
}

// ExplFind is expl_find().
func (z *Streams) ExplFind(who, where int) int {
	z.beginExplore()
	return z.explore.s.Keyed(who, where, xtagFind, 1, 100)
}

// ExplGate is expl_gate().
func (z *Streams) ExplGate(who, where int) int {
	z.beginExplore()
	return z.explore.s.Keyed(who, where, xtagGate, 1, 100)
}

// ExplFlavor is expl_flavor().
func (z *Streams) ExplFlavor(who, where int) int {
	z.beginExplore()
	return z.explore.s.Keyed(who, where, xtagFlavor, 1, 4)
}

// ExplPick is expl_pick().
func (z *Streams) ExplPick(who, where, high int) int {
	z.beginExplore()
	return z.explore.s.Keyed(who, where, xtagPick, 1, high)
}

// ExplSeek is expl_seek().
func (z *Streams) ExplSeek(who, target int) int {
	z.beginExplore()
	return z.explore.s.Keyed(who, target, xtagSeek, 1, 10)
}

// ExplDetect is expl_detect().
func (z *Streams) ExplDetect(who, cand int) int {
	z.beginExplore()
	return z.explore.s.Keyed(who, cand, xtagDetect, 1, 100)
}

// ---------------------------------------------------------------------------
// code.c: the per-turn mint stream (entity ids, passwords).
// ---------------------------------------------------------------------------

func (z *Streams) beginMint() { z.seedPerTurn(&z.mint, tagMint) }

// MintAlloc is mint_alloc(), the draw under rnd_alloc_num().
func (z *Streams) MintAlloc(low, high int) int {
	z.beginMint()
	return z.mint.s.Draw(low, high)
}

// MintPassword is mint_password().
func (z *Streams) MintPassword(pl, i, n int) int {
	z.beginMint()
	return z.mint.s.Keyed(pl, i, mtagPW, 1, n)
}

// MintCity is mint_city().
func (z *Streams) MintCity() int {
	z.beginMint()
	return z.mint.s.Draw(1, 6)
}

// MintShuffle is mint_shuffle().
func (z *Streams) MintShuffle(l []int) {
	z.beginMint()
	Shuffle(l, &z.mint.s)
}

// ---------------------------------------------------------------------------
// dir.c: the per-turn exit-direction stream.
// ---------------------------------------------------------------------------

// BeginExitdir is begin_exitdir().
func (z *Streams) BeginExitdir() { z.seedPerTurn(&z.exitdir, tagExitdir) }

// Exitdir is &exitdir_rng, which exits_from_loc_nsew_select() shuffles its
// exit views with after BeginExitdir.
func (z *Streams) Exitdir() *Stream { return &z.exitdir.s }

// ---------------------------------------------------------------------------
// day.c: the per-turn upkeep and calendar streams.
// ---------------------------------------------------------------------------

func (z *Streams) beginUpkeep() { z.seedPerTurn(&z.upkeep, tagUpkeep) }

// UpHeal is up_heal().
func (z *Streams) UpHeal(who, sub, low, high int) int {
	z.beginUpkeep()
	return z.today(&z.upkeep.s, who, sub, utagHeal, low, high)
}

// UpLoyal is up_loyal().
func (z *Streams) UpLoyal(who, sub, low, high int) int {
	z.beginUpkeep()
	return z.upkeep.s.Keyed(who, sub, utagLoyal, low, high)
}

// UpStarve is up_starve().
func (z *Streams) UpStarve(who, item, low, high int) int {
	z.beginUpkeep()
	return z.upkeep.s.Keyed(who, item, utagSicken, low, high)
}

// UpAnimal is up_animal(); the item is folded into the tag.
func (z *Streams) UpAnimal(who, item, idx, low, high int) int {
	z.beginUpkeep()
	return z.upkeep.s.Keyed(who, idx, utagAnimal^uint32(item), low, high)
}

// UpCorpse is up_corpse().
func (z *Streams) UpCorpse(who, low, high int) int {
	z.beginUpkeep()
	return z.upkeep.s.Keyed(who, 0, utagCorpse, low, high)
}

// UpIncome is up_income().
func (z *Streams) UpIncome(inn, sub, low, high int) int {
	z.beginUpkeep()
	return z.upkeep.s.Keyed(inn, sub, utagIncome, low, high)
}

// CalDay is cal_day().
func (z *Streams) CalDay(which, low, high int) int {
	z.seedPerTurn(&z.calendar, tagCalendar)
	return z.calendar.s.Keyed(which, 0, ctagDay, low, high)
}

// ---------------------------------------------------------------------------
// faery.c: the per-turn autonomous faery stream.
// ---------------------------------------------------------------------------

func (z *Streams) beginFaery() { z.seedPerTurn(&z.faery, tagFaery) }

// FaerHuntLoc is faer_hunt_loc().
func (z *Streams) FaerHuntLoc(slot, retry, low, high int) int {
	z.beginFaery()
	return z.faery.s.Keyed(slot, retry, ftagHunt, low, high)
}

// FaerHuntQty is faer_hunt_qty().
func (z *Streams) FaerHuntQty(slot int) int {
	z.beginFaery()
	return z.faery.s.Keyed(slot, 0, ftagHQty, 25, 100)
}

// FaerBanditKind is faer_bandit_kind().
func (z *Streams) FaerBanditKind(who, where int) int {
	z.beginFaery()
	return z.today(&z.faery.s, who, where, ftagBKnd, 1, 3)
}

// FaerBanditQty is faer_bandit_qty().
func (z *Streams) FaerBanditQty(unit int) int {
	z.beginFaery()
	return z.faery.s.Keyed(unit, 0, ftagBQty, 4, 24)
}

// FaerBanditGold is faer_bandit_gold().
func (z *Streams) FaerBanditGold(unit int) int {
	z.beginFaery()
	return z.faery.s.Keyed(unit, 0, ftagBGld, 1, 25)
}

// FaerAmbush is faer_ambush().
func (z *Streams) FaerAmbush(who, where int) int {
	z.beginFaery()
	return z.today(&z.faery.s, who, where, ftagAmbs, 1, 100)
}

// FaerRetal is faer_retal().
func (z *Streams) FaerRetal(who, where int) int {
	z.beginFaery()
	return z.today(&z.faery.s, who, where, ftagRtal, 1, 2)
}

// ---------------------------------------------------------------------------
// hades.c: the per-turn autonomous hades stream.
// ---------------------------------------------------------------------------

func (z *Streams) beginHades() { z.seedPerTurn(&z.hades, tagHades) }

// HadsNasty is hads_nasty().
func (z *Streams) HadsNasty(slot, sub, low, high int) int {
	z.beginHades()
	return z.hades.s.Keyed(slot, sub, htagNast, low, high)
}

// HadsBanditKind is hads_bandit_kind().
func (z *Streams) HadsBanditKind(who, where int) int {
	z.beginHades()
	return z.today(&z.hades.s, who, where, htagBKnd, 1, 5)
}

// HadsBanditQty is hads_bandit_qty().
func (z *Streams) HadsBanditQty(unit int) int {
	z.beginHades()
	return z.hades.s.Keyed(unit, 0, htagBQty, 4, 24)
}

// HadsAmbush is hads_ambush().
func (z *Streams) HadsAmbush(who, where int) int {
	z.beginHades()
	return z.today(&z.hades.s, who, where, htagAmbs, 1, 100)
}

// HadsRetal is hads_retal().
func (z *Streams) HadsRetal(who, where int) int {
	z.beginHades()
	return z.today(&z.hades.s, who, where, htagRtal, 1, 2)
}

// HadsRandomLoc is the pick in random_hades_loc(), written inline in C: an
// index into a candidate list of length n, keyed on who.
func (z *Streams) HadsRandomLoc(who, n int) int {
	z.beginHades()
	return z.hades.s.Keyed(who, 0, htagRLoc, 0, n-1)
}

// ---------------------------------------------------------------------------
// seed.c: the per-turn worldgen stream (INIT-time city seeding).
// ---------------------------------------------------------------------------

func (z *Streams) beginWorldgen() { z.seedPerTurn(&z.worldgen, tagWorldgen) }

// WgenProm is wgen_prom().
func (z *Streams) WgenProm(where int) int {
	z.beginWorldgen()
	return z.worldgen.s.Keyed(where, 0, wtagProm, 1, 100)
}

// WgenSkill is wgen_skill().
func (z *Streams) WgenSkill(where, sub, low, high int) int {
	z.beginWorldgen()
	return z.worldgen.s.Keyed(where, sub, wtagTech, low, high)
}

// WgenGarr is wgen_garr().
func (z *Streams) WgenGarr(where int) int {
	z.beginWorldgen()
	return z.worldgen.s.Keyed(where, 0, wtagGarr, 25, 150)
}

// WgenGate is wgen_gate().
func (z *Streams) WgenGate(city int) int {
	z.beginWorldgen()
	return z.worldgen.s.Keyed(city, 0, wtagGate, 1, 2)
}

// ---------------------------------------------------------------------------
// storm.c: the per-turn weather stream. Its counter carries across the turn.
// ---------------------------------------------------------------------------

// BeginWeather is begin_weather().
func (z *Streams) BeginWeather() { z.seedPerTurn(&z.weather, tagWeather) }

// WthrShuffle is wthr_shuffle().
func (z *Streams) WthrShuffle(l []int) {
	z.BeginWeather()
	Shuffle(l, &z.weather.s)
}

// WthrStorm is wthr_storm().
func (z *Streams) WthrStorm(where, k2, low, high int) int {
	z.BeginWeather()
	return z.weather.s.Keyed(where, k2, wtagStorm, low, high)
}

// WthrWreck is wthr_wreck().
func (z *Streams) WthrWreck(where, sub, low, high int) int {
	z.BeginWeather()
	return z.weather.s.Keyed(where, sub, wtagWreck, low, high)
}

// ---------------------------------------------------------------------------
// swear.c: the per-turn social stream.
// ---------------------------------------------------------------------------

func (z *Streams) beginSocial() { z.seedPerTurn(&z.social, tagSocial) }

// SocGift is soc_gift().
func (z *Streams) SocGift(who, target int) int {
	z.beginSocial()
	return z.social.s.Keyed(who, target, sotagGift, 1, 3)
}

// SocBribe is soc_bribe().
func (z *Streams) SocBribe(who, target, sub, low, high int) int {
	z.beginSocial()
	return z.today(&z.social.s, who, target*2+sub, sotagBribe, low, high)
}

// SocTerror is soc_terror().
func (z *Streams) SocTerror(who, target int) int {
	z.beginSocial()
	return z.today(&z.social.s, who, target, sotagTerror, 1, 100)
}

// SocIncite is soc_incite().
func (z *Streams) SocIncite(who, target, sub, low, high int) int {
	z.beginSocial()
	return z.today(&z.social.s, who, target*2+sub, sotagIncite, low, high)
}

// SocPersuade is soc_persuade().
func (z *Streams) SocPersuade(who, target int) int {
	z.beginSocial()
	return z.today(&z.social.s, who, target, sotagOath, 1, 100)
}

// SocBreed is soc_breed().
func (z *Streams) SocBreed(who, item, sub, low, high int) int {
	z.beginSocial()
	return z.today(&z.social.s, who, item*2+sub, sotagBreed, low, high)
}

// ---------------------------------------------------------------------------
// use.c: the per-turn skills stream.
// ---------------------------------------------------------------------------

func (z *Streams) beginSkills() { z.seedPerTurn(&z.skills, tagSkills) }

// SkilCrit is skil_crit().
func (z *Streams) SkilCrit(who, wsk int) int {
	z.beginSkills()
	return z.today(&z.skills.s, who, wsk, stagCrit, 1, 100)
}

// SkilBonus is skil_bonus().
func (z *Streams) SkilBonus(who, wsk, low, high int) int {
	z.beginSkills()
	return z.today(&z.skills.s, who, wsk, stagYield, low, high)
}

// SkilStudy is skil_study().
func (z *Streams) SkilStudy(who, sk int) int {
	z.beginSkills()
	return z.today(&z.skills.s, who, sk, stagStudy, 1, 4)
}

// SkilResearch is skil_research().
func (z *Streams) SkilResearch(who, sk int) int {
	z.beginSkills()
	return z.today(&z.skills.s, who, sk, stagResch, 1, 100)
}

// SkilResearchPick is skil_research_pick().
func (z *Streams) SkilResearchPick(who, sk, high int) int {
	z.beginSkills()
	return z.today(&z.skills.s, who, sk, stagRPick, 0, high)
}

// SkilTorture is skil_torture().
func (z *Streams) SkilTorture(who, target int) int {
	z.beginSkills()
	return z.today(&z.skills.s, who, target, stagTort, 1, 100)
}

// SkilPetty is skil_petty().
func (z *Streams) SkilPetty(who, where, sub, low, high int) int {
	z.beginSkills()
	return z.today(&z.skills.s, who, where*16+sub, stagPetty, low, high)
}

// ---------------------------------------------------------------------------
// u.c: the per-turn entity stream.
// ---------------------------------------------------------------------------

func (z *Streams) beginEntity() { z.seedPerTurn(&z.entity, tagEntity) }

// EntPrisoner is ent_prisoner().
func (z *Streams) EntPrisoner(who, sub int) int {
	z.beginEntity()
	return z.today(&z.entity.s, who, sub, entagPrison, 1, 1000)
}

// EntDrop is ent_drop().
func (z *Streams) EntDrop(who, prisoner, sub, low, high int) int {
	z.beginEntity()
	return z.today(&z.entity.s, who, prisoner*2+sub, entagDrop, low, high)
}

// EntTake is ent_take().
func (z *Streams) EntTake(from, item, sub, low, high int) int {
	z.beginEntity()
	return z.entity.s.Keyed(from, item*2+sub, entagTake, low, high)
}

// EntSick is ent_sick().
func (z *Streams) EntSick(who, health int) int {
	z.beginEntity()
	return z.entity.s.Keyed(who, health, entagSick, 1, 100)
}

// EntLandDir is ent_land_dir().
func (z *Streams) EntLandDir(origWhere, retry int) int {
	z.beginEntity()
	return z.entity.s.Keyed(origWhere, retry, entagLandDir, 1, 4)
}

// EntLandPick is ent_land_pick().
func (z *Streams) EntLandPick(origWhere, n int) int {
	z.beginEntity()
	return z.entity.s.Keyed(origWhere, 0, entagLandPick, 0, n-1)
}

// EntBark is ent_bark().
func (z *Streams) EntBark(who, hound int) int {
	z.beginEntity()
	return z.today(&z.entity.s, who, hound, entagBark, 1, 2)
}

// EntBuild is ent_build().
func (z *Streams) EntBuild(where int) int {
	z.beginEntity()
	return z.entity.s.Keyed(where, 0, entagBuild, 1, 5)
}

// EntMenial is ent_menial().
func (z *Streams) EntMenial(who int) int {
	z.beginEntity()
	return z.entity.s.Keyed(who, 0, entagMenial, 1, 9)
}
