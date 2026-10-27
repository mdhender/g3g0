// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Expected bytes come only from tools/move-compat's real pinned C helpers.
func TestMoveMatchesC(t *testing.T) {
	var sections [][]string
	for line := range strings.SplitSeq(readTestdata(t, "move-cases.txt"), "\n") {
		switch {
		case line == "" || strings.HasPrefix(line, "#"):
		case strings.HasPrefix(line, "section "):
			sections = append(sections, []string{line})
		default:
			sections[len(sections)-1] = append(sections[len(sections)-1], line)
		}
	}
	var out strings.Builder
	for _, section := range sections {
		out.WriteString(runMoveSection(t, section))
	}
	want := readTestdata(t, "move-c-output.txt")
	gotLines, wantLines := strings.Split(out.String(), "\n"), strings.Split(want, "\n")
	if len(gotLines) != len(wantLines) {
		t.Fatalf("line count = %d, want %d (including final empty lines)", len(gotLines), len(wantLines))
	}
	for i := range wantLines {
		if wantLines[i] != gotLines[i] {
			t.Fatalf("line %d:\n got %q\nwant %q", i+1, gotLines[i], wantLines[i])
		}
	}
}

func runMoveSection(t *testing.T, lines []string) (output string) {
	var stdout, stderr strings.Builder
	w := newWorld()
	w.stdout, w.stderr = &stdout, &stderr
	status := 0
	defer func() {
		switch r := recover().(type) {
		case nil:
		case assertion:
			stderr.WriteString(r.Error())
			status = 134
		case runtime.Error:
			// Only the separately C-recorded, known unowned-character NULL
			// dereference may discard buffered output (G3 #353).
			if lines[0] != "section unowned hidden route G3 353" {
				panic(r)
			}
			stdout.Reset()
			stderr.Reset()
			status = 139
		default:
			panic(r)
		}
		output = fmt.Sprintf("stdout:\n%sstderr:\n%s\nexit %d\n", stdout.String(), stderr.String(), status)
	}()
	dump := func(l []*exitView) {
		fmt.Fprintf(&stdout, "exits %d\n", len(l))
		for _, v := range l {
			fmt.Fprintf(&stdout, "exit %d %d %d %d %d %d %d %d %d %d %d %d %d %d\n",
				v.direction, v.destination, v.orig, v.distance, v.impassable, v.destHidden,
				v.origHidden, v.hidden, v.inside, v.road, v.water, v.inTransit, v.magicBarrier, v.hadesCost)
		}
	}
	fmt.Fprintln(&stdout, lines[0])
	for _, line := range lines[1:] {
		fmt.Fprintf(&stdout, "> %s\n", line)
		c := &caseLine{rest: line}
		op := c.word()
		var a [4]int
		for i := range a {
			if strings.TrimSpace(c.rest) != "" {
				a[i] = c.num(t)
			}
		}
		n, b, sk, d := a[0], a[1], a[2], a[3]
		switch op {
		case "box":
			w.allocBox(n, b, sk)
			w.pLocInfo(n).where = int32(d)
			if d != 0 {
				w.pLocInfo(d).hereList.Append(n)
			}
		case "dest":
			w.pLoc(n).provDest.Append(b)
		case "hidden":
			w.pLoc(n).hidden = int8(b)
		case "barrier":
			w.pLoc(n).barrier = int16(b)
		case "lane":
			w.pLoc(n).seaLane = int8(b)
		case "uldim":
			w.pSubloc(n).uldimFlag = int8(b)
		case "summer":
			w.pSubloc(n).summerFlag = int8(b)
		case "open":
			w.pSubloc(n).linkOpen = int8(b)
		case "from":
			w.pSubloc(n).linkFrom.Append(b)
		case "to":
			w.pSubloc(n).linkTo.Append(b)
		case "road":
			w.pGate(n).toLoc, w.pGate(n).roadHidden = int32(b), int8(sk)
		case "lord":
			w.pChar(n).unitLord = int32(b)
		case "known":
			setBit(&w.pPlayer(n).known, b)
		case "hades":
			w.hadesRegion = int32(n)
		case "clock":
			w.sysclock.turn, w.sysclock.daysSinceEpoch, w.evening, w.monthDone = int16(n), int32(b), int32(sk), int32(d)
		case "day":
			w.sysclock.day = int16(n)
		case "moving":
			w.pSubloc(n).moving = int32(b)
		case "master":
			seed := make([]byte, 16)
			for i := range seed {
				seed[i] = byte(n + i)
			}
			path := filepath.Join(t.TempDir(), "randseed")
			if err := os.WriteFile(path, seed, 0o600); err != nil {
				t.Fatal(err)
			}
			w.loadSeed(path)
		case "legacy":
			fmt.Fprintf(&stdout, "legacy %d\n", w.rnd.Rnd(n, b))
		case "direction":
			fmt.Fprintf(&stdout, "direction %d\n", w.locationDirection(n, b))
		case "distance":
			fmt.Fprintf(&stdout, "distance %d\n", w.exitDistance(n, b))
		case "port":
			fmt.Fprintf(&stdout, "port %d\n", b2i(w.isPortCity(n)))
		case "province_port":
			fmt.Fprintf(&stdout, "province_port %d\n", w.provinceHasPortCity(n))
		case "ocean":
			fmt.Fprintf(&stdout, "ocean %d\n", w.hasOceanAccess(n))
		case "exits":
			dump(w.exitsFromLoc(n, b))
		case "nsew":
			dump(w.exitsFromLocNSEW(n, b))
		case "select":
			dump(w.exitsFromLocNSEWSelect(n, b, sk, d))
		default:
			t.Fatalf("unknown move operation %q", op)
		}
	}
	return ""
}
