// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package lists_test

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

	"github.com/mdhender/g3g0/internal/lists"
	"github.com/mdhender/g3g0/internal/rnd"
	"github.com/mdhender/g3g0/internal/rng"
)

// TestMatchesC runs testdata/cases.txt the way the C driver in
// tools/lists-compat/run.sh does and requires its output, line for line.
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

// harness holds what the driver's operations share: the global rnd() state
// and the current stream.
type harness struct {
	src *rnd.Source
	cur rng.Stream
}

// objs gives the pointer lists distinct elements: element k is &objs[k], as
// in the C driver.
var objs [256]byte

// module is one list module of the driver: three list slots, the
// conversions between an element and the int written in the cases, and the
// operations only some modules have.
type module[T comparable] struct {
	slots      [3]lists.List[T]
	mk         func(int) T
	un         func(T) int
	ilistOnly  bool // ilist_add, ilist_rem_value_uniq
	shuffleRng bool // ilist_shuffle_rng, exit_views_shuffle_rng
}

func (m *module[T]) do(h *harness, slot int, op string, a []int) (string, bool) {
	l := &m.slots[slot]
	var res string
	switch op {
	case "append":
		l.Append(m.mk(a[0]))
	case "fill":
		for i := 0; i < a[0]; i++ {
			l.Append(m.mk(a[1] + i))
		}
	case "prepend":
		l.Prepend(m.mk(a[0]))
	case "clear":
		l.Clear()
	case "copy":
		*l = m.slots[a[0]].Copy()
	case "delete":
		l.Delete(a[0])
	case "len":
		res = strconv.Itoa(l.Len()) + " "
	case "lookup":
		res = strconv.Itoa(l.Lookup(m.mk(a[0]))) + " "
	case "reclaim":
		l.Reclaim()
	case "remvalue":
		l.RemValue(m.mk(a[0]))
	case "scramble":
		l.Scramble(h.src)
	case "shuffle":
		l.Shuffle(h.src)
	case "add":
		if !m.ilistOnly {
			return "", false
		}
		l.Add(m.mk(a[0]))
	case "remuniq":
		if !m.ilistOnly {
			return "", false
		}
		l.RemValueUniq(m.mk(a[0]))
	case "shufflerng":
		if !m.shuffleRng {
			return "", false
		}
		l.ShuffleRng(&h.cur)
		res = fmt.Sprintf("ctr %d ", h.cur.Counter)
	default:
		return "", false
	}
	if *l == nil {
		return res + "NULL", true
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%slen %d:", res, l.Len())
	for _, v := range *l {
		fmt.Fprintf(&b, " %d", m.un(v))
	}
	return b.String(), true
}

// elem stands in for the element pointer types of the template modules.
type elem = *byte

func ptrModule(shuffleRng bool) *module[elem] {
	return &module[elem]{
		mk: func(k int) elem { return &objs[k] },
		un: func(p elem) int {
			for k := range objs {
				if &objs[k] == p {
					return k
				}
			}
			panic("not an element")
		},
		shuffleRng: shuffleRng,
	}
}

// runCases is the Go side of the driver: same operations, same output.
func runCases(t *testing.T, cases []byte) []byte {
	t.Helper()
	dir := t.TempDir()
	var out bytes.Buffer
	h := &harness{src: rnd.New([16]byte{})}

	type doer func(h *harness, slot int, op string, a []int) (string, bool)
	modules := map[string]doer{
		"ilist": (&module[int]{
			mk:         func(v int) int { return v },
			un:         func(v int) int { return v },
			ilistOnly:  true,
			shuffleRng: true,
		}).do,
	}
	for _, name := range []string{"accept_ents", "admits", "cstrings", "exit_views", "fights",
		"flag_ents", "item_ents", "orders", "req_ents", "roads", "skill_ents", "tiles", "trades", "wait_args"} {
		modules[name] = ptrModule(name == "exit_views").do
	}

	sc := bufio.NewScanner(bytes.NewReader(cases))
	for sc.Scan() {
		line, _, _ := strings.Cut(sc.Text(), "#")
		f := strings.Fields(line)
		if len(f) == 0 {
			continue
		}
		num := func(s string) int {
			v, err := strconv.ParseInt(s, 0, 64)
			if err != nil {
				t.Fatalf("%q: %v", line, err)
			}
			return int(int32(v))
		}

		var res string
		switch f[0] {
		case "master":
			b, err := hex.DecodeString(f[1])
			if err != nil {
				t.Fatalf("master %q: %v", f[1], err)
			}
			path := filepath.Join(dir, "randseed")
			if err := os.WriteFile(path, b, 0o644); err != nil {
				t.Fatal(err)
			}
			if err := h.src.Load(path); err != nil {
				t.Fatal(err)
			}
			res = "ok"
		case "rnd":
			res = strconv.Itoa(h.src.Rnd(num(f[1]), num(f[2])))
		case "root":
			h.cur = rng.Seed(h.src.MasterSeed())
			res = fmt.Sprintf("ctr %d", h.cur.Counter)
		case "of":
			v, err := strconv.ParseUint(f[2], 0, 32)
			if err != nil {
				t.Fatalf("%q: %v", line, err)
			}
			h.cur = h.cur.Of(num(f[1]), uint32(v))
			res = fmt.Sprintf("ctr %d", h.cur.Counter)
		default:
			do, ok := modules[f[0]]
			if !ok || len(f) < 3 || len(f[1]) != 1 || f[1][0] < 'a' || f[1][0] > 'c' {
				t.Fatalf("bad line %q", line)
			}
			a := make([]int, 16)
			for i := 3; i < len(f); i++ {
				if f[2] == "copy" {
					a[0] = int(f[3][0] - 'a')
					break
				}
				a[i-3] = num(f[i])
			}
			if res, ok = do(h, int(f[1][0]-'a'), f[2], a); !ok {
				t.Fatalf("unknown op %q", line)
			}
		}
		fmt.Fprintf(&out, "%s -> %s\n", strings.Join(f, " "), res)
	}
	return out.Bytes()
}
