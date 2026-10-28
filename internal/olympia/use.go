// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Skill lookup and report helpers from olympia/use.c.

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
