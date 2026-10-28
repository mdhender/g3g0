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
            body = body.replace(f"w.streams.{method[0].upper() + method[1:]}(",
                                f"w.{instrumented}(")
        for go, replacement in mapping.items():
            body = re.sub(rf"w\.{go}\(\)", lambda _: replacement, body)
        source = source[:m.start()] + body + source[m.end():]
    # Do not create a Go package under tmp that go test ./... would discover.
    dest = work / (filename + ".overlay")
    dest.write_text(source)
    overlay[str(path)] = str(dest)
(work / "overlay.json").write_text(json.dumps({"Replace": overlay}))

# Separate REAL caln/wthr evidence. The older synthetic collision fixture above
# remains unchanged. Keep the original static declarations and all executable
# statements; the sole insertion into daily_events is a read-only observation
# at its end, since its function-local static storage has no external accessor.
def extract(file, pattern):
    text = (g3 / file).read_text()
    m = re.search(pattern, text, re.M | re.S)
    if not m:
        sys.exit(f"missing schedule definition: {file} {pattern}")
    line = text.count("\n", 0, m.start()) + 1
    return f'#line {line} "olympia/{file}"\n{m[0]}\n'

schedule = '#include <stdio.h>\n#include <stdlib.h>\n#include "z.h"\n#include "oly.h"\n#include "rng.h"\n'
# Complete stream blocks, including original macros, initializers and guards.
schedule += extract("day.c", r"^#define\s+TAG_TURN[^\n]*\n")
schedule += extract("day.c", r"^static rng_stream calendar_rng;.*?^cal_day\(.*?^\}\n")
schedule += extract("storm.c", r"^static rng_stream weather_rng;.*?^wthr_wreck\(.*?^\}\n")
schedule += extract("seed.c", r"^int\nint_comp\(.*?^\}\n")
schedule += extract("z.c", r"^void\nsafe_qsort\(.*?^\}\n")
schedule += extract("z.c", r"^void\nasfail\(.*?^\}\n")
daily = next(d for d in definitions if "\ndaily_events(" in d)
daily_callbacks = sorted(set(re.findall(r"^\s+(\w+)\(\);", daily, re.M)))
for name in daily_callbacks:
    schedule += f'void {name}(void) {{ printf("callback {name}: %d %d\\n", sysclock.turn, sysclock.day); }}\n'
schedule += r'''
static void print_list(const char *name, ilist l) {
    printf("%s:", name);
    for (int i = 0; i < ilist_len(l); i++) printf(" %d", l[i]);
    printf("\n");
}
static int observed_cal_day(int pick, int low, int high) {
    int n = cal_day(pick, low, high);
    printf("cal_day %d %d %d -> %d\n", pick, low, high, n);
    return n;
}
static void observed_shuffle(ilist l) {
    print_list("shuffle input", l);
    wthr_shuffle(l);
    print_list("shuffle output", l);
}
#define cal_day observed_cal_day
#define wthr_shuffle observed_shuffle
'''
observe = '''
    printf("schedule: curse=%d faery=%d dog=%d index=%d\\n",
           curse_erode_day, faery_day, dog_bark_day, wday_index);
    print_list("weather days", weather_days);
'''
schedule += daily.rsplit("}", 1)[0] + observe + "}\n"
schedule += '#undef cal_day\n#undef wthr_shuffle\n'
(work / "schedule.c").write_text(schedule)

path = root / "internal/olympia/day.go"
source = path.read_text()
m = re.search(r"^func \(w \*world\) dailyEvents\(\) \{.*?^\}", source, re.M | re.S)
body = m[0]
for name in daily_callbacks:
    body = body.replace(f"w.{camel(name)}()", f'w.scheduleCallback("{name}")')
body = body.replace("w.streams.CalDay(", "w.scheduleCalDay(")
body = body.replace("w.streams.WthrShuffle(", "w.scheduleShuffle(")
body = body[:-1] + '\tw.scheduleObserve()\n}'
dest = work / "schedule-day.go.overlay"
dest.write_text(source[:m.start()] + body + source[m.end():])
(work / "schedule-overlay.json").write_text(json.dumps({"Replace": {str(path): str(dest)}}))
