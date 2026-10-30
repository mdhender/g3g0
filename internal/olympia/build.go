// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Construction: build.c:8–811. Repair/improve and raze remain separate owners.
func fortDefaultDefense(sk int) int {
	switch sk {
	case subCastle:
		return 50
	case subTower:
		return 40
	case subGalley:
		return 20
	case subRoundship, subTemple, subMine, subInn:
		return 10
	}
	return 0
}

func (w *world) shipLocOkay(c *command, where int) int {
	if w.hasOceanAccess(where) == 0 {
		w.commandOut(int(c.who), "%s is not an ocean port location.", w.boxName(where))
		return 0
	}
	if w.subkind(where) != subCity {
		w.commandOut(int(c.who), "Ships may only be built in cities.")
		return 0
	}
	return 1
}

func (w *world) templeLocOkay(c *command, where int) int {
	if w.safeHaven(where) != 0 {
		w.commandOut(int(c.who), "Building is not permitted in safe havens.")
		return 0
	}
	if w.locDepth(where) != locBuild {
		return 1
	}
	w.commandOut(int(c.who), "A temple may not be built inside another building.")
	return 0
}

func (w *world) towerLocOkay(c *command, where int) int {
	ld := w.locDepth(where)
	// build.c:72–78 disables the tower safe-haven check with #if 0.
	if ld != locProvince && ld != locSubloc && w.subkind(where) != subCastle && w.subkind(where) != subCastleNotdone {
		w.commandOut(int(c.who), "A tower may not be built here.")
		return 0
	}
	if ld == locBuild && w.countLocStructures(where, subTower, subTowerNotdone) >= 6 {
		w.commandOut(int(c.who), "Six towers at most can be built within a %s.", subkindS[w.subkind(where)])
		return 0
	}
	return 1
}

func (w *world) mineLocOkay(c *command, where int) int {
	if w.subkind(where) != subMountain && w.subkind(where) != subRockyHill {
		w.commandOut(int(c.who), "Mines may only be built in mountain provinces and rocky hills.")
		return 0
	}
	if w.safeHaven(where) != 0 {
		w.commandOut(int(c.who), "Building is not permitted in safe havens.")
		return 0
	}
	if w.countLocStructures(where, subMine, subMineNotdone) != 0 {
		w.commandOut(int(c.who), "A location may not have more than one mine.")
		return 0
	}
	if w.countLocStructures(where, subMineCollapsed, 0) != 0 {
		w.commandOut(int(c.who), "Another mine may not be built here until the collapsed mine vanishes.")
		return 0
	}
	return 1
}

func (w *world) innLocOkay(c *command, where int) int {
	if w.safeHaven(where) != 0 {
		w.commandOut(int(c.who), "Building is not permitted in safe havens.")
		return 0
	}
	if w.locDepth(where) != locProvince && w.subkind(where) != subCity {
		w.commandOut(int(c.who), "Inns may only be built in cities and provinces.")
		return 0
	}
	return 1
}

func (w *world) castleLocOkay(c *command, where int) int {
	if w.safeHaven(where) != 0 {
		w.commandOut(int(c.who), "Building is not permitted in safe havens.")
		return 0
	}
	if w.locDepth(where) != locProvince && w.subkind(where) != subCity {
		w.commandOut(int(c.who), "A castle must be built in a province or a city.")
		return 0
	}
	if w.provinceSubloc(where, subCastle) != 0 || w.provinceSubloc(where, subCastleNotdone) != 0 {
		w.commandOut(int(c.who), "This province already contains a castle.  Another may not be built here.")
		return 0
	}
	return 1
}

type buildEnt struct {
	what                               string
	skillReq, skill2                   int
	kind                               int
	locOK                              func(*world, *command, int) int
	unfinishedSubkind, finishedSubkind int
	minWorkers, workerDays, minDays    int
	reqItem, reqQty, capacity          int
	defaultName                        string
}

// Every original field and indirect callback, in original matching order.
var buildTbl = [...]buildEnt{
	{"galley", skShipbuilding, 0, tShip, (*world).shipLocOkay, subGalleyNotdone, subGalley, 3, 250, 1, itemLumber, 10, 5000, "New galley"},
	{"roundship", skShipbuilding, 0, tShip, (*world).shipLocOkay, subRoundshipNotdone, subRoundship, 3, 500, 1, itemLumber, 20, 25000, "New roundship"},
	{"raft", 0, 0, tShip, (*world).shipLocOkay, subRaftNotdone, subRaft, 0, 45, 1, itemFlotsam, 5, 2500, "New raft"},
	{"temple", skConstruction, 0, tLoc, (*world).templeLocOkay, subTempleNotdone, subTemple, 3, 1000, 1, itemStone, 10, 0, "New temple"},
	{"inn", skConstruction, 0, tLoc, (*world).innLocOkay, subInnNotdone, subInn, 3, 300, 1, itemLumber, 15, 0, "New inn"},
	{"castle", skConstruction, 0, tLoc, (*world).castleLocOkay, subCastleNotdone, subCastle, 3, 10000, 1, itemStone, 100, 0, "New castle"},
	{"mine", skConstruction, skMining, tLoc, (*world).mineLocOkay, subMineNotdone, subMine, 3, 500, 1, itemLumber, 5, 0, "New mine"},
	{"tower", skConstruction, 0, tLoc, (*world).towerLocOkay, subTowerNotdone, subTower, 3, 2000, 1, itemStone, 20, 0, "New tower"},
	{}, // C's NULL sentinel.
}

func (w *world) findBuild(s string) *buildEnt {
	w.fuzzyBuildMatch = false
	for i := range buildTbl {
		if buildTbl[i].what == "" {
			break
		}
		if iStrcmp(buildTbl[i].what, s) == 0 {
			return &buildTbl[i]
		}
	}
	w.fuzzyBuildMatch = true
	for i := range buildTbl {
		if buildTbl[i].what == "" {
			break
		}
		if fuzzyStrcmp(buildTbl[i].what, s) {
			return &buildTbl[i]
		}
	}
	return nil
}

func (w *world) buildMaterialsCheck(c *command, bi *buildEnt) int {
	who := int(c.who)
	if bi.skillReq != 0 {
		if bi.skill2 != 0 {
			if w.hasSkill(who, bi.skillReq) < 1 && w.hasSkill(who, bi.skill2) < 1 {
				w.commandOut(who, "Building a %s requires either %s or %s.", bi.what, w.boxName(bi.skillReq), w.boxName(bi.skill2))
				return 0
			}
		} else if w.hasSkill(who, bi.skillReq) < 1 {
			w.commandOut(who, "Building a %s requires %s.", bi.what, w.boxName(bi.skillReq))
			return 0
		}
	}
	if bi.reqItem > 0 && w.hasItem(who, bi.reqItem) < int32(bi.reqQty) {
		w.commandOut(who, "Need %s to start.", w.boxNameQty(bi.reqItem, bi.reqQty))
		return 0
	}
	if w.hasItem(who, itemWorker) < int32(bi.minWorkers) {
		w.commandOut(who, "Need at least %s for construction.", w.boxNameQty(itemWorker, bi.minWorkers))
		return 0
	}
	return 1
}

// create_new_building, build.c:430–457. Use the shared subtype/cache,
// keyed entity stream, mine stock floor and skill-experience owners.
func (w *world) createNewBuilding(c *command, bi *buildEnt, where int) {
	w.changeBoxSubkind(where, bi.finishedSubkind)
	p := w.pSubloc(where)
	p.effortGiven, p.effortRequired = 0, 0
	p.defense = int32(fortDefaultDefense(bi.finishedSubkind))
	p.damage, p.buildMaterials = 0, 0
	if bi.finishedSubkind == subMine {
		p.shaftDepth = 3
		if w.streams.EntBuild(where) == 1 {
			w.genItem(where, itemGateCrystal, 1)
		}
		w.mineProduction(where)
	}
	if bi.finishedSubkind == subTemple {
		p.teaches.Append(skReligion)
	}
	if bi.skillReq != 0 {
		w.addSkillExperience(int(c.who), bi.skillReq)
	}
}

func (w *world) startBuild(c *command, bi *buildEnt, new int) int {
	pl := w.player(int(c.who))
	where := w.subloc(int(c.who))
	// instant_build is always FALSE. Its enabling code is #if 0 in C:
	// neither instant notices nor create_new_building can execute. Do not
	// substitute a successful completion callback for that disabled branch.
	if w.buildMaterialsCheck(c, bi) == 0 {
		return 0
	}
	if new == 0 {
		new = w.newEnt(tUnform, 0)
	}
	if new < 0 {
		w.commandOut(int(c.who), "Out of new entity codes.")
		return 0
	}
	w.pPlayer(pl).unformed.RemValue(new)
	w.changeBoxKind(new, bi.kind)
	w.changeBoxSubkind(new, bi.unfinishedSubkind)
	w.setWhere(new, where)
	name := bi.defaultName
	if c.parse.Len() >= 3 && c.parse[2] != "" {
		name = c.parse[2]
	}
	if len(name) > 25 {
		w.commandOut(int(c.who), "The name you gave is too long.  Place names must be 25 characters or less.  Please use the NAME order to set a shorter name next turn.")
		name = bi.defaultName
	}
	w.setName(new, name)
	p := w.pSubloc(new)
	p.effortRequired = int32(bi.workerDays) * 100
	p.damage, p.buildMaterials, p.capacity = 0, 0, int32(bi.capacity)
	w.commandOut(int(c.who), "Created %s.", w.boxNameKind(new))
	w.showToGarrison = true
	w.commandOut(where, "%s began construction of %s in %s.", w.boxName(int(c.who)), w.boxNameKind(new), w.boxName(where))
	w.showToGarrison = false
	w.moveStack(int(c.who), new)
	return 1
}

// daily_build, build.c:563–663. A failed later material increment keeps
// earlier consumption, but does not commit today's effort. The <= threshold
// intentionally charges the first fifth even when fifth is zero.
func (w *world) dailyBuild(c *command, bi *buildEnt) int {
	who := int(c.who)
	inside := w.subloc(who)
	if w.subkind(inside) == bi.finishedSubkind {
		w.commandOut(who, "%s is finished!", w.boxName(inside))
		c.wait = 0
		return 1
	}
	if w.subkind(inside) != bi.unfinishedSubkind {
		w.commandOut(who, "%s is no longer in a %s.  Construction halts.", w.justName(who), bi.what)
		return 0
	}
	nworkers := int32(1)
	if bi.minWorkers != 0 {
		nworkers = w.hasItem(who, itemWorker)
	}
	if nworkers <= 0 {
		w.commandOut(who, "%s has no workers.  Construction halts.", w.justName(who))
		return 0
	}
	p := w.pSubloc(inside)
	bonus := 5 * c.useExp * nworkers
	effortGiven := nworkers*100 + bonus
	if bi.reqItem > 0 {
		fifth := (p.effortGiven + effortGiven) * 5 / p.effortRequired
		for p.buildMaterials < 5 && int32(p.buildMaterials) <= fifth {
			if !w.consumeItem(who, bi.reqItem, bi.reqQty) {
				w.commandOut(who, "Need another %s to continue work.  Construction halted.", w.boxNameQty(bi.reqItem, bi.reqQty))
				return 0
			}
			w.commandOut(who, "Used %s in construction of %s.", w.boxNameQty(bi.reqItem, bi.reqQty), w.boxName(inside))
			p.buildMaterials++
		}
	}
	p.effortGiven += effortGiven
	if p.effortGiven < p.effortRequired || c.daysExecuting < int32(bi.minDays) {
		return 1
	}
	w.createNewBuilding(c, bi, inside)
	w.commandOut(who, "%s is finished!", w.boxName(inside))
	w.showToGarrison = true
	w.commandOut(w.subloc(who), "%s completed construction of %s in %s.", w.boxName(who), w.boxNameKind(inside), w.boxName(w.loc(inside)))
	w.showToGarrison = false
	c.wait = 0
	return 1
}

func (w *world) buildStructure(c *command, bi *buildEnt, new int) int {
	who := int(c.who)
	where := w.subloc(who)
	if w.locDepth(where) == locBuild {
		if w.subkind(where) == bi.finishedSubkind {
			w.commandOut(who, "%s is already finished.", w.boxName(where))
			return 0
		}
		if w.subkind(where) == bi.unfinishedSubkind {
			w.commandOut(who, "Continuing work on %s.", w.boxName(where))
			return 1
		}
	}
	if w.subkind(where) == subOcean {
		w.commandOut(who, "Construction may not take place at sea.")
		return 0
	}
	if w.subkind(where) == subTunnel {
		w.commandOut(who, "Building is not permitted underground.")
		return 0
	}
	if bi.locOK(w, c, where) == 0 {
		return 0
	}
	return w.startBuild(c, bi, new)
}

func (w *world) unfinishedBuilding(who int) *string {
	// Deliberately no raft: C's bare BUILD shorthand omits it.
	switch w.subkind(w.subloc(who)) {
	case subCastleNotdone:
		return strSave("castle")
	case subTowerNotdone:
		return strSave("tower")
	case subTempleNotdone:
		return strSave("temple")
	case subGalleyNotdone:
		return strSave("galley")
	case subRoundshipNotdone:
		return strSave("roundship")
	case subInnNotdone:
		return strSave("inn")
	case subMineNotdone:
		return strSave("mine")
	}
	return nil
}

func (w *world) vBuild(c *command) int {
	days, new := c.c, int(c.d)
	if c.parse.Len() < 2 {
		if s := w.unfinishedBuilding(int(c.who)); s != nil {
			w.commandOut(int(c.who), "(assuming you meant 'build %s')", *s)
			line := sout("build %s", *s)
			ret := w.olyParse(c, &line)
			assert(ret, "olympia/build.c", 742, "ret")
			return w.vBuild(c)
		}
	}
	if c.parse.Len() < 2 {
		w.commandOut(int(c.who), "Must specify what to build.")
		return 0
	}
	t := w.findBuild(c.parse[1])
	if t == nil {
		w.commandOut(int(c.who), "Don't know how to build '%s'.", c.parse[1])
		return 0
	}
	if w.fuzzyBuildMatch {
		w.commandOut(int(c.who), "(assuming you meant 'build %s')", t.what)
	}
	if days != 0 {
		c.wait = days
	}
	p := w.pPlayer(w.player(int(c.who)))
	if new != 0 && (w.kind(new) != tUnform || p.unformed.Lookup(new) < 0) {
		w.commandOut(int(c.who), "%s is not a valid unformed entity code.", w.boxCode(new))
		new = 0
	}
	// C's automatic choice from p->unformed is also disabled by #if 0.
	return w.buildStructure(c, t, new)
}

// d_build, build.c:796–811. Saved command text is reparsed by the scheduler;
// lookup still happens at completion rather than retaining a table pointer.
func (w *world) dBuild(c *command) int {
	t := w.findBuild(c.parse[1])
	if t == nil {
		w.logWrite(10, "d_build: t is NULL (%s)", c.parse[1])
		w.commandOut(int(c.who), "Internal error.")
		return 0
	}
	return w.dailyBuild(c, t)
}
