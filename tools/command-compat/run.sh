#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"
. "${HERE}/../lib/timeout.sh"
RUN_TIMEOUT="${RUN_TIMEOUT:-10}"
G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-2567e6a5dd0ba0d95516d8bed9e1bafb87a51985}"
WORK="${WORK:-${ROOT}/tmp/command-compat}"
CC="${CC:-cc}"
CASES="${ROOT}/internal/olympia/testdata/command-cases.txt"
EXPECT="${ROOT}/internal/olympia/testdata/command-c-output.txt"
update=0
case "${1:-}" in
"") ;;
--update) update=1 ;;
*) echo "usage: $0 [--update]" >&2; exit 2 ;;
esac
die() { echo "error: $*" >&2; exit 2; }
git -C "${G3_REPO}" cat-file -e "${G3_REV}^{commit}" || die "missing G3 revision"
rm -rf "${WORK}"
mkdir -p "${WORK}/g3" "${WORK}/run1" "${WORK}/run2"
git -C "${G3_REPO}" archive "${G3_REV}" lib olympia | tar -x -C "${WORK}/g3"
echo "G3 revision: ${G3_REV}"

python3 -I - "${WORK}/g3/olympia" "${WORK}" <<'PY'
import re, sys
src, dst = sys.argv[1:]
functions = {
    'input.c': ['parse_line', 'remove_comment', 'remove_ctrl_chars', 'find_command',
                'oly_parse_cmd', 'parse_arg', 'oly_parse', 'cmd_shift', 'check_allow',
                'get_command', 'set_state', 'load_command', 'command_done',
                'finish_command', 'do_command', 'init_load_sup'],
    'order.c': ['p_order_head', 'rp_order_head', 'top_order', 'pop_order'],
    'swear.c': ['lord', 'player'],
    'stack.c': ['stack_parent', 'stack_leader'],
    'loc.c': ['add_char_here', 'all_stack', 'subloc', 'first_character'],
    'garr.c': ['garrison_here'],
    'u.c': ['min', 'cap', 'set_bit', 'test_bit', 'count_stack_any', 'count_any', 'check_skill'],
    'use.c': ['rp_skill_ent', 'exp_level', 'has_skill'],
    'c2.c': ['v_stop', 'v_format', 'v_notab'],
    'immed.c': ['v_xyzzy', 'v_plugh'],
    'stealth.c': ['v_hide', 'd_hide'],
}
real = {n for names in functions.values() for n in names}
with open(f'{dst}/helpers.c', 'w') as out:
    out.write('''#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "z.h"
#include "oly.h"
#define MAX_PRI 5
#define MAX_BUFFER 50
#define BUF_LEN LEN
ilist load_q[MAX_PRI], run_q;
int cur_pri, cmd_wait = -1, fuzzy_find;
''')
    for file, names in functions.items():
        source = open(f'{src}/{file}').read()
        matches = []
        for name in names:
            m = re.search(rf'^(?:static )?[^\n(]*\n{name}\(.*?^\}}\n', source, re.M | re.S)
            if not m:
                sys.exit(f'no definition of {name} in {file}')
            line = source.count('\n', 0, m.start()) + 1
            matches.append((m.start(), f'#line {line} "olympia/{file}"\n{m.group(0)}\n'))
        for _, text in sorted(matches):
            out.write(text)
    out.write('''
int drv_load(struct command *c) { return load_command(c); }
void drv_state(struct command *c, int s) { set_state(c, s, c->pri); }
struct order_list *drv_orders(int pl, int who) { return p_order_head(pl, who); }
''')
glob = open(f'{src}/glob.c').read()
table = glob[glob.index('struct cmd_tbl_ent cmd_tbl[]'):glob.index('char *kind_s[]')]
routines = set()
for m in re.finditer(r'^\{[^,]*,[^,]*,\s*(\w+),\s*(\w+),\s*(\w+),', table, re.M):
    routines.update(n for n in m.groups() if n != 'NULL')
with open(f'{dst}/stubs.c', 'w') as out:
    out.write('''#include <stdio.h>
#include <stdlib.h>
#include "z.h"
#include "oly.h"
#define STUB(n) { fprintf(stderr, "out-of-scope C dependency: %s\\n", n); abort(); }
''')
    for name in sorted(routines - real):
        out.write(f'int {name}(struct command *c) STUB("{name}")\n')
    out.write('''
/* Only unreachable reporting routes and unrelated code.c dependencies. */
int is_port_city(int n) STUB("is_port_city")
char *nice_num(int n) STUB("nice_num")
ilist players_who_rule_here(int n) STUB("players_who_rule_here")
int loc_depth(int n) STUB("loc_depth")
int viewloc(int n) STUB("viewloc")
int loc_hidden(int n) STUB("loc_hidden")
int test_known(int who, int n) STUB("test_known")
void check_validated_trades(void) STUB("check_validated_trades")
/* Instrumentation only; coverage artifacts belong to stage 3. */
void coverage_count(int cmd, int result) { (void)cmd; (void)result; }
int immediate, flush_always = 1;
ilist trades_to_check;
''')
PY

case "$(uname -s)" in
 Darwin) ldflags=(-Wl,-dead_strip) ;;
 *) ldflags=(-Wl,--gc-sections) ;;
esac
(cd "${WORK}/g3" && "${CC}" -std=c11 -w -ffunction-sections -fdata-sections "${ldflags[@]}" \
 -I olympia -I lib \
 "${HERE}/driver.c" "${WORK}/helpers.c" "${WORK}/stubs.c" \
 olympia/{code,glob,z,sout}.c \
 lib/{checked_alloc,ilist,rng,rnd,cstrings,item_ents,skill_ents,orders}.c \
 -o "${WORK}/driver") >"${WORK}/build.log" 2>&1 || die "C build failed: ${WORK}/build.log"
sections=$(grep -c '^section ' "${CASES}")
for run in run1 run2; do
 for ((i=0; i<sections; i++)); do
  mkdir -p "${WORK}/${run}/${i}/log"
  status=0
  (run_limited "${RUN_TIMEOUT}" command-driver "${WORK}/driver" "${CASES}" "$i" "${WORK}/${run}/${i}") \
   >"${WORK}/${run}/stdout" 2>"${WORK}/${run}/stderr" || status=$?
  [ "$status" != 124 ] || die "C section $i timed out"
  ! grep -q 'out-of-scope C dependency:' "${WORK}/${run}/stderr" || die "section $i reached a C stub"
  case "$status" in
   0|134) ;;
   *) die "unexpected C status $status in section $i; see ${WORK}/${run}" ;;
  esac
  {
   printf 'stdout:\n'; cat "${WORK}/${run}/stdout"
   printf 'stderr:\n'; cat "${WORK}/${run}/stderr"
   printf '\nexit %d\n' "$status"
  } >>"${WORK}/${run}/out.txt"
 done
done
cmp "${WORK}/run1/out.txt" "${WORK}/run2/out.txt" || die "C is not deterministic"
echo "C driver: ${sections} sections, $(wc -l <"${WORK}/run1/out.txt") lines, deterministic"
if [ "$update" = 1 ]; then
 cp "${WORK}/run1/out.txt" "${EXPECT}"
else
 diff -u "${EXPECT}" "${WORK}/run1/out.txt"
fi
(cd "${ROOT}" && RUN_TIMEOUT="${RUN_TIMEOUT}" go test -exec "bash ${HERE}/../lib/run-limited.sh" \
 -count=1 -run '^TestCommand' ./internal/olympia)
echo PASS
