#!/usr/bin/env python3
"""Test-only callbacks around real C and Go orchestration, never engine stubs.

No orchestration body is modeled here. C definitions are extracted unchanged;
Go's build overlay replaces only leaf phase calls inside the corresponding
production functions. The normal build has no callback switch or interface.
"""
import json
import pathlib
import re
import sys

root, work = map(pathlib.Path, sys.argv[1:])
g3 = work / "g3/olympia"
functions = {
    "input.c": ["daily_command_loop", "process_orders"],
    "day.c": ["post_month", "daily_events"],
}
definitions = []
for file, names in functions.items():
    source = (g3 / file).read_text()
    for name in names:
        m = re.search(rf"^(?:static )?void\n{name}\(.*?^\}}\n", source, re.M | re.S)
        if not m:
            sys.exit(f"missing {file}:{name}")
        line = source.count("\n", 0, m.start()) + 1
        definitions.append(f'#line {line} "olympia/{file}"\n{m[0]}')

real = "\n".join(definitions)
callbacks = sorted(set(re.findall(r"^\s+(\w+)\(\);", real, re.M)) -
                   {"daily_events", "daily_command_loop", "clear_orders_sent"})
source = (g3 / "day.c").read_text()
clear = re.search(r"^static void\nclear_orders_sent\(.*?^\}\n", source, re.M | re.S)[0]

def camel(name):
    first, *rest = name.split("_")
    return first + "".join(word.title() for word in rest)

# This fixture supplies deliberately colliding calendar picks and a deliberately
# unsorted weather prefix. It tests order/guards, not caln/wthr RNG behavior.
prefix = r'''
#include <stdio.h>
#include <stdlib.h>
#include "z.h"
#include "oly.h"
#define phase_ran(x) ((void)0)
#define phase_acted(x) ((void)0)
extern int cmd_wait, auto_attack_flag, month_done;
static void trace(const char *name) {
    printf("phase %s: %d %d %d auto=%d month=%d post=%d\n", name,
           sysclock.turn, sysclock.day, sysclock.days_since_epoch,
           auto_attack_flag, month_done, post_has_been_run);
}
int find_command(char *s) { trace(s); return 42; }
static int cal_day(int pick, int low, int high) {
    printf("cal_day %d %d %d\n", pick, low, high);
    return pick == 0 ? 7 : 15;
}
void wthr_shuffle(ilist days) {
    int prefix[] = {15, 7, 1, 3}, i, j = 4;
    trace("wthr_shuffle");
    for (i = 0; i < 4; i++) days[i] = prefix[i];
    for (i = 1; i <= MONTH_DAYS; i++)
        if (i != 15 && i != 7 && i != 1 && i != 3) days[j++] = i;
}
int int_comp(const void *a, const void *b) { return *(const int *)a - *(const int *)b; }
void safe_qsort(void *base, size_t n, size_t size, int (*cmp)(const void *, const void *)) {
    qsort(base, n, size, cmp);
}
'''
for name in callbacks:
    storage = "static " if name in {"process_interrupted_units", "process_player_orders", "clear_orders_sent", "evening_phase"} else ""
    # start/evening aren't in the selected definitions' parameterless calls
    # unless daily_command_loop was extracted too (as required).
    prefix += f'{storage}void {name}(void) {{ trace("{name}"); }}\n'
source = (g3 / "z.c").read_text()
prefix += re.search(r"^void\nasfail\(.*?^\}\n", source, re.M | re.S)[0]
(work / "phases.c").write_text(prefix + clear + real)

# Replace calls only in the four orchestration bodies, not implementations.
# Every callback is generated from a real C call, making an absent/misspelled
# Go call a byte-for-byte mismatch rather than a generated expected trace.
mapping = {camel(n): f'w.phaseTrace("{n}")' for n in callbacks}
overlay = {}
for filename, names in {"input.go": ["processOrders", "dailyCommandLoop"],
                        "day.go": ["postMonth", "dailyEvents"]}.items():
    path = root / "internal/olympia" / filename
    source = path.read_text()
    for name in names:
        pattern = rf"^func \(w \*world\) {name}\(\) \{{.*?^\}}"
        m = re.search(pattern, source, re.M | re.S)
        if not m:
            sys.exit(f"missing Go {name}")
        body = m[0]
        for method in ("findCommand", "calDay", "wthrShuffle"):
            # Preserve the production arguments, not merely their call order.
            instrumented = "phase" + method[0].upper() + method[1:]
            body = body.replace(f"w.{method}(", f"w.{instrumented}(")
        for go, replacement in mapping.items():
            body = re.sub(rf"w\.{go}\(\)", lambda _: replacement, body)
        source = source[:m.start()] + body + source[m.end():]
    # Do not create a Go package under tmp that go test ./... would discover.
    dest = work / (filename + ".overlay")
    dest.write_text(source)
    overlay[str(path)] = str(dest)
(work / "overlay.json").write_text(json.dumps({"Replace": overlay}))
