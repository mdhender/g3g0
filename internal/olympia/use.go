// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// Ported from olympia/use.c: skill_school(), which check_db() uses, and
// the read-only skill lookup used by the minimal hide command handlers.

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
