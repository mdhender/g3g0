// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// The small direct-stream cases run normally; the callback schedule requires
// the tool's source overlay. Neither generates expected dates from Go.
func TestCalendarMatchesC(t *testing.T) { runCalendarCases(t, "calendar") }

func TestDailyScheduleMatchesC(t *testing.T) {
	if os.Getenv("G3G0_SCHEDULE_TRACE") != "1" {
		t.Skip("schedule callbacks require tools/turn-driver-compat/run.sh overlay")
	}
	runCalendarCases(t, "schedule")
}

func runCalendarCases(t *testing.T, prefix string) {
	t.Helper()
	bin, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	section := 0
	for line := range strings.SplitSeq(readTestdata(t, prefix+"-cases.txt"), "\n") {
		if !strings.HasPrefix(line, "section ") {
			continue
		}
		// Capture the actual child process streams and status, not a fabricated
		// success status or a test-framework transcript. The shared supervisor
		// bounds every child even in ordinary default go tests.
		cmd := exec.Command("bash", "../../tools/lib/run-limited.sh", bin, "-test.run=^TestCalendarProcess$")
		cmd.Env = append(os.Environ(), "G3G0_CALENDAR_CHILD="+prefix, "G3G0_CALENDAR_SECTION="+strconv.Itoa(section))
		var stdout, stderr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		status := 0
		if err := cmd.Run(); err != nil {
			if exit, ok := err.(*exec.ExitError); ok {
				status = exit.ExitCode()
			} else {
				t.Fatal(err)
			}
		}
		fmt.Fprintf(&out, "stdout:\n%sstderr:\n%s\nexit %d\n", stdout.String(), stderr.String(), status)
		section++
	}
	if want := readTestdata(t, prefix+"-c-output.txt"); out.String() != want {
		t.Fatalf("C recording mismatch:\n got %q\nwant %q", out.String(), want)
	}
}

// TestCalendarProcess is a test-only process entry point. No engine mode or
// production hook is added; os.Exit avoids contaminating the raw C comparison
// with the testing package's PASS line. Temporary files are explicitly cleaned.
func TestCalendarProcess(t *testing.T) {
	prefix := os.Getenv("G3G0_CALENDAR_CHILD")
	if prefix == "" {
		t.Skip("child process entry point")
	}
	if prefix != "calendar" && prefix != "schedule" {
		t.Fatal("invalid calendar suite")
	}
	want, err := strconv.Atoi(os.Getenv("G3G0_CALENDAR_SECTION"))
	if err != nil {
		t.Fatal(err)
	}
	dir, err := os.MkdirTemp("", "g3g0-calendar-")
	if err != nil {
		t.Fatal(err)
	}
	seedPath := filepath.Join(dir, "seed")
	w := newWorld()
	section := -1
	for line := range strings.SplitSeq(readTestdata(t, prefix+"-cases.txt"), "\n") {
		if line == "" || line[0] == '#' {
			continue
		}
		if strings.HasPrefix(line, "section ") {
			section++
			if section == want {
				fmt.Fprintln(w.stdout, line)
			}
			continue
		}
		if section != want {
			continue
		}
		fmt.Fprintf(w.stdout, "> %s\n", line)
		c := &caseLine{rest: line}
		switch op := c.word(); op {
		case "clock":
			w.sysclock.turn, w.sysclock.day = int16(c.num(t)), int16(c.num(t))
		case "seed":
			var seed [16]byte
			first := c.num(t)
			for i := range seed {
				seed[i] = byte(first + i)
			}
			if err := os.WriteFile(seedPath, seed[:], 0600); err != nil {
				t.Fatal(err)
			}
			if err := w.rnd.Load(seedPath); err != nil {
				t.Fatal(err)
			}
		case "cal":
			w.scheduleCalDay(c.num(t), c.num(t), c.num(t))
		case "shuffle":
			l := make([]int, c.num(t))
			for i := range l {
				l[i] = i + 1
			}
			w.scheduleShuffle(l)
		case "rnd":
			fmt.Fprintf(w.stdout, "rnd -> %d\n", w.rnd.Rnd(c.num(t), c.num(t)))
		case "streams":
			m := w.rnd.MasterSeed()
			fmt.Fprintf(w.stdout, "master: %08x %08x %08x %08x\n", m[0], m[1], m[2], m[3])
			// Read private stream state only in the test; do not add production
			// observers or copy/reimplement the source-of-truth RNG helpers.
			for _, name := range []string{"calendar", "weather"} {
				p := reflect.ValueOf(w.streams).Elem().FieldByName(name)
				s := p.FieldByName("s")
				seed := s.FieldByName("Seed")
				fmt.Fprintf(w.stdout, "%s: turn=%d seed=%08x %08x %08x %08x counter=%d\n",
					name, p.FieldByName("turn").Int(), seed.Index(0).Uint(), seed.Index(1).Uint(),
					seed.Index(2).Uint(), seed.Index(3).Uint(), s.FieldByName("Counter").Uint())
			}
		case "daily_events":
			w.dailyEvents()
		case "days":
			first, last := c.num(t), c.num(t)
			for d := first; d <= last; d++ {
				w.sysclock.day = int16(d)
				w.dailyEvents()
			}
		default:
			t.Fatalf("unknown operation %q", op)
		}
	}
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	os.Exit(0)
}

func (w *world) scheduleCallback(name string) {
	fmt.Fprintf(w.stdout, "callback %s: %d %d\n", name, w.sysclock.turn, w.sysclock.day)
}

func (w *world) scheduleCalDay(pick, low, high int) int {
	n := w.streams.CalDay(pick, low, high)
	fmt.Fprintf(w.stdout, "cal_day %d %d %d -> %d\n", pick, low, high, n)
	return n
}

func (w *world) scheduleList(name string, l []int) {
	fmt.Fprintf(w.stdout, "%s:", name)
	for _, n := range l {
		fmt.Fprintf(w.stdout, " %d", n)
	}
	fmt.Fprintln(w.stdout)
}

func (w *world) scheduleShuffle(l []int) {
	w.scheduleList("shuffle input", l)
	w.streams.WthrShuffle(l)
	w.scheduleList("shuffle output", l)
}

func (w *world) scheduleObserve() {
	fmt.Fprintf(w.stdout, "schedule: curse=%d faery=%d dog=%d index=%d\n",
		w.curseErodeDay, w.faeryDay, w.dogBarkDay, w.wdayIndex)
	w.scheduleList("weather days", w.weatherDays)
}
