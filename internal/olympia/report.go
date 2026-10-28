// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// charRepLocation is report.c:357–407. Death logging needs this helper;
// report-body assembly remains separate stage-2 work.
func (w *world) charRepLocation(who int) string {
	where := w.subloc(who)
	if where == 0 {
		return "nowhere"
	}
	s := ""
	for w.locDepth(where) > locProvince {
		if s != "" {
			s = sout("%s, in %s", s, w.boxName(where))
		} else {
			s = w.boxName(where)
		}
		where = w.loc(where)
	}
	if w.subkind(w.province(where)) == subOcean {
		if s != "" {
			s = sout("%s, in %s", s, w.boxName(where))
		} else {
			s = w.boxName(where)
		}
	} else {
		if s != "" {
			s = sout("%s, in province %s", s, w.boxName(where))
		} else {
			s = w.boxName(where)
		}
	}
	if regName := w.name(w.region(where)); regName != "" {
		s = sout("%s, in %s", s, regName)
	}
	return s
}
