#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"
. "${HERE}/../lib/timeout.sh"
RUN_TIMEOUT="${RUN_TIMEOUT:-10}"
G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-6a25f2d252ce45bd5adbe3baaccb3be3eed56279}"
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
                'finish_command', 'do_command', 'init_load_sup', 'exec_precedence',
                'exec_comp', 'sort_load_queue', 'sort_run_queue', 'initial_command_load',
                'min_pri_ready', 'init_wait_list', 'check_all_waits', 'start_phase', 'evening_phase',
                'daily_command_loop', 'process_player_orders', 'interrupt_order', 'process_interrupted_units'],
    'order.c': ['p_order_head', 'rp_order_head', 'top_order', 'pop_order', 'is_stop_order', 'stop_order', 'flush_unit_orders', 'queue_order', 'queue', 'save_player_orders'],
    'swear.c': ['lord', 'player', 'set_loyal', 'reswear_all_sworn', 'set_lord'],
    'move.c': ['restore_stack_actions', 'touch_loc_after_move', 'move_stack'],
    'day.c': ['touch_loc_pl', 'touch_loc', 'storm_owner_touch_loc', 'init_locs_touched'],
    'seed.c': ['int_comp'],
    'stack.c': ['stack_parent', 'stack_leader', 'here_pos', 'stacked_beneath',
                'first_prisoner_pos', 'promote', 'unstack', 'leave_stack', 'stack', 'join_stack'],
    'combat.c': ['fort_covers', 'auto_attack', 'check_auto_attack_sup', 'check_all_auto_attacks', 'clear_second_waits'],
    'loc.c': ['add_char_here', 'all_stack', 'subloc', 'first_character',
              'in_safe_now', 'province', 'building_owner', 'viewloc', 'all_char_here',
              'in_here_list', 'add_to_here_list', 'remove_from_here_list', 'set_where',
              'region', 'mark_loc_stack_known'],
    'garr.c': ['garrison_here', 'province_admin', 'players_who_rule_here'],
    'u.c': ['min', 'cap', 'set_bit', 'test_bit', 'count_stack_any', 'count_any', 'check_skill',
            'contacted', 'char_where', 'char_here', 'test_known', 'loc_hidden', 'has_item', 'loc_depth',
            'lookup', 'gold_s', 'comma_num', 'nice_num', 'add_np', 'deduct_np',
            'new_char', 'beast_capturable', 'create_unique_item', 'create_unique_item_alloc',
            'destroy_unique_item', 'add_item', 'sub_item', 'set_known', 'gen_item',
            'can_pay', 'charge', 'stack_has_item', 'has_use_key', 'stack_has_use_key',
            'stack_sub_item', 'autocharge'],
    'npc.c': ['create_hades_bandit', 'create_faery_bandit', 'hades_attack_check', 'faery_attack_check'],
    'lore.c': ['queue_lore', 'scan_char_item_lore'],
    'buy.c': ['investigate_possible_trade', 'market_here', 'seller_comp',
              'seller_list', 'buyer_list', 'scan_trades', 'match_trades'],
    'c1.c': ['flag_raised', 'clear_wait_parse', 'parse_wait_args', 'wait_leader', 'check_wait_conditions',
             'v_wait', 'd_wait', 'i_wait'],
    'produce.c': ['init_collect_list'],
    'perm.c': ['is_hostile', 'clear_all_att'],
    'storm.c': ['weather_here', 'update_weather_view_loc_sup', 'init_weather_views',
                'update_weather_view_locs'],
    'use.c': ['rp_skill_ent', 'exp_level', 'has_skill', 'char_np_total'],
    'c2.c': ['v_stop', 'v_format', 'v_notab'],
    'immed.c': ['v_xyzzy', 'v_plugh'],
    'stealth.c': ['v_hide', 'd_hide', 'clear_contacts'],
}
real = {n for names in functions.values() for n in names}
with open(f'{dst}/helpers.c', 'w') as out:
    out.write('''#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include "z.h"
#include "oly.h"
#define MAX_PRI 5
#define MAX_BUFFER 50
#define BUF_LEN LEN
ilist load_q[MAX_PRI], run_q;
int cur_pri, cmd_wait = -1, fuzzy_find;
int evening, auto_attack_flag;
ilist wait_list, second_wait_list;
int pledge_backlinks;
#define phase_ran(x) ((void)0)
#define phase_acted(x) ((void)0)
''')
    for file, names in functions.items():
        source = open(f'{src}/{file}').read()
        # Original data declarations too: no copied or modeled tag/number table.
        if file == 'c1.c':
            start, end = source.index('struct flag_ent {'), source.index('static int\nflag_raised(')
            out.write(source[start:end])
            start, end = source.index('static char *wait_tags[]'), source.index('void\nclear_wait_parse(')
            out.write(source[start:end])
        if file == 'u.c':
            start, end = source.index('static char *num_s[]'), source.index('char *\nnice_num(')
            out.write(source[start:end])
        if file == 'produce.c':
            out.write(re.search(r'^static ilist collectors = NULL;\n', source, re.M).group(0))
        if file == 'move.c':
            out.write('#line 8 "olympia/move.c"\n')
            out.write(re.search(r'^static ilist ocean_chars = NULL;\n', source, re.M).group(0))
        if file == 'buy.c':
            out.write('static void attempt_trade(struct trade *b, struct trade *s) { fprintf(stderr, "out-of-scope C dependency: attempt_trade\\n"); fflush(NULL); abort(); }\n')
        matches = []
        for name in names:
            pattern = rf'^(?:static )?[^\n(]*\n{name}\(.*?^\}}\n'
            if name == 'queue':
                pattern = r'^void queue\(.*?^\}\n'
            m = re.search(pattern, source, re.M | re.S)
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
void drv_evening(void) { evening_phase(); }
void drv_daily(void) { daily_command_loop(); }
void drv_players(void) { process_player_orders(); }
void drv_interrupted(void) { process_interrupted_units(); }
int drv_sub_item(int who, int item, int qty) { return sub_item(who, item, qty); }
void drv_add_item(int who, int item, int qty) { add_item(who, item, qty); }
void drv_sort_load(int pri) { sort_load_queue(load_q[pri]); }
void drv_sort_run(void) { sort_run_queue(run_q); }
ilist drv_collectors(void) { return collectors; }
int drv_fort_covers(int n) { return fort_covers(n); }
int drv_bandit(int faery, int who, int where) {
 return faery ? create_faery_bandit(who, where) : create_hades_bandit(who, where);
}
int drv_market(int who) { return market_here(who); }
void drv_ocean(int who) { ilist_append(&ocean_chars, who); }
ilist drv_oceans(void) { return ocean_chars; }
trades_list drv_offers(int seller, int where, int except) {
 return seller ? seller_list(where, except) : buyer_list(where, except);
}
char *drv_wait_check(struct command *c) { return check_wait_conditions(c); }
void drv_signal(int who, char *flag) {
 struct flag_ent *p = my_malloc(sizeof(*p));
 p->who = who; p->flag = str_save(flag); flag_ents_append(&flags, p);
}
''')
# Region prefixes retain the original tags, structures, turn guards and whole
# keyed leaf functions. No replaced draw or reset; the same master feeds mint.
for region in ['hades', 'faery']:
    source = open(f'{src}/{region}.c').read()
    with open(f'{dst}/{region}-rng.c', 'w') as out:
        out.write(f'#line 1 "olympia/{region}.c"\n')
        out.write(source[:source.index(f'void\ncreate_{region}(void)')])
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
void check_validated_trades(void) STUB("check_validated_trades")
/* Instrumentation only; coverage artifacts belong to stage 3. */
void coverage_count(int cmd, int result) { (void)cmd; (void)result; }
int immediate, flush_always = 1, pretty_data_files;
ilist trades_to_check;
''')
PY

case "$(uname -s)" in
 Darwin) ldflags=(-Wl,-dead_strip) ;;
 *) ldflags=(-Wl,--gc-sections) ;;
esac
(cd "${WORK}/g3" && "${CC}" -std=c11 -w -ffunction-sections -fdata-sections "${ldflags[@]}" \
 -I olympia -I lib \
 "${HERE}/driver.c" "${WORK}/helpers.c" "${WORK}/stubs.c" "${WORK}/hades-rng.c" "${WORK}/faery-rng.c" \
 olympia/{code,glob,z,sout,io}.c \
 lib/{checked_alloc,ilist,rng,rnd,cstrings,item_ents,skill_ents,orders,wait_args,flag_ents,admits,trades,req_ents}.c \
 -o "${WORK}/driver") >"${WORK}/build.log" 2>&1 || die "C build failed: ${WORK}/build.log"
sections=$(grep -c '^section ' "${CASES}")
section_names=() # Bash 3 (Darwin) has no mapfile.
while IFS= read -r line; do section_names+=("$line"); done < <(grep '^section ' "${CASES}")
for run in run1 run2; do
 for ((i=0; i<sections; i++)); do
  mkdir -p "${WORK}/${run}/${i}/log" "${WORK}/${run}/${i}/orders"
  status=0
  (run_limited "${RUN_TIMEOUT}" command-driver "${WORK}/driver" "${CASES}" "$i" "${WORK}/${run}/${i}") \
   >"${WORK}/${run}/stdout" 2>"${WORK}/${run}/stderr" || status=$?
  [ "$status" != 124 ] || die "C section $i timed out"
  ! grep -q 'out-of-scope C dependency:' "${WORK}/${run}/stderr" || die "section $i reached a C stub"
  if [[ "${section_names[$i]}" == 'section arrival '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section arrival mismatched trade owner retains exact assertion'|\
    'section arrival move rejects noncharacter root'|\
    'section arrival city matching includes captive descendants'|\
    'section arrival unowned weather mage retains exact assertion') expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected arrival status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section combat boundary '* ]]; then
   [ "$status" = 134 ] && grep -qx 'combat boundary: regular_combat' "${WORK}/${run}/stderr" \
    || die "section $i did not reach actual combat"
  else
   ! grep -q 'combat boundary:' "${WORK}/${run}/stderr" || die "unexpected combat in section $i"
  fi
  case "$status" in
   0|134|139) ;;
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
