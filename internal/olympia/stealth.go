// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// addContact is the shared seek helper. CONTACT deliberately appends itself,
// as in C; it must not also call a helper that adds a second side effect.
func (w *world) addContact(a, b int) {
	assert(w.kind(a) == tChar, "olympia/stealth.c", 305, "kind(a) == T_char")
	w.pChar(a).contact.Append(b)
}

func (w *world) vContact(c *command) int {
	for c.parse.Len() > 1 {
		if k := w.kind(int(c.a)); k != tChar && k != tPlayer {
			w.commandOut(int(c.who), "%s is not a character or player entity.", c.parse[1])
		} else {
			w.pChar(int(c.who)).contact.Append(int(c.a))
			w.commandOut(int(c.a), "%s contacted us.", w.boxName(int(c.who)))
		}
		w.cmdShift(c)
	}
	return 1
}

func (w *world) cloakLord(n int) int { return w.hasSkill(n, skHideLord) }

// vHide and dHide are the complete hide handlers from stealth.c. They
// provide real failure, timed completion and zero-time inhibited completion
// paths for dispatch without implementing any part of wait or scheduling.
func (w *world) vHide(c *command) int {
	if !w.checkSkill(int(c.who), skHideSelf) {
		return 0
	}
	if c.a != 0 && (w.stackParent(int(c.who)) != 0 || w.countStackAny(int(c.who)) != 1) {
		w.commandOut(int(c.who), "Must be alone to hide.")
		return 0
	}
	if c.a == 0 {
		w.pMagic(int(c.who)).hideSelf = 0
		w.commandOut(int(c.who), "No longer hidden.")
		c.wait = 0
		c.inhibitFinish = 1
	}
	return 1
}

func (w *world) dHide(c *command) int {
	if w.stackParent(int(c.who)) != 0 || w.countStackAny(int(c.who)) != 1 {
		w.commandOut(int(c.who), "Must be alone to hide.")
		return 0
	}
	w.pMagic(int(c.who)).hideSelf = 1
	w.commandOut(int(c.who), "Now hidden.")
	return 1
}
