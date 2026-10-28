#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail
ulimit -c 0
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"
. "${HERE}/../lib/timeout.sh"
RUN_TIMEOUT="${RUN_TIMEOUT:-10}"
G3_REPO="${G3_REPO:-${ROOT}/tmp/upstream.git}"
G3_REV="${G3_REV:-e733acb63bacb21167a1f2589d8dcdc023364896}"
WORK="${WORK:-${ROOT}/tmp/move-compat}"
CC="${CC:-cc}"
CASES="${ROOT}/internal/olympia/testdata/move-cases.txt"
EXPECT="${ROOT}/internal/olympia/testdata/move-c-output.txt"
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

# Complete definitions, original declarations/macros and original source lines.
# No C gameplay stubs, modeled allocator, or replacement exit shuffle.
python3 -I - "${WORK}/g3/olympia" "${WORK}/helpers.c" "${WORK}/combat-structs.h" <<'PY'
import re, sys
src, dest, structs = sys.argv[1:]
combat = open(f'{src}/combat.c').read()
with open(structs, 'w') as out:
    start = combat.index('struct fight {')
    end = combat.index('\n};', combat.index('struct wield {')) + 3
    out.write('#line 110 "olympia/combat.c"\n' + combat[start:end] + '\n')
functions = {
    'dir.c': ['location_direction', 'exit_distance', 'is_port_city',
              'province_has_port_city', 'summer_uldim_open_now', 'add_province_exit',
              'extra_routes', 'province_exits', 'province_sub_exits', 'subloc_exits',
              'ship_exits', 'building_exits', 'exits_from_loc', 'exits_from_loc_nsew',
              'exits_from_loc_nsew_select', 'has_ocean_access'],
    'loc.c': ['region', 'province', 'subloc', 'viewloc', 'add_here', 'all_here',
              'add_char_here', 'all_char_here', 'all_stack', 'first_character',
              'building_owner', 'in_here_list', 'add_to_here_list',
              'remove_from_here_list', 'set_where'],
    'u.c': ['loc_depth', 'loc_hidden', 'test_known', 'test_bit', 'set_bit',
            'max', 'add_item_weight', 'determine_unit_weights', 'determine_stack_weights',
            'count_stack_move_nobles', 'count_stack_figures', 'count_any',
            'count_stack_units', 'count_stack_any',
            'has_item', 'sub_item', 'comma_num', 'nice_num', 'cap', 'lookup'],
    'swear.c': ['lord', 'player', 'set_loyal'],
    'move.c': ['move_exit_land', 'move_exit_fly', 'move_exit_water',
               'parse_exit_dir', 'save_v_array', 'restore_v_array',
               'suspend_stack_actions', 'restore_stack_actions', 'clear_guard_flag',
               'land_check', 'can_move_here', 'can_move_at_outer_level',
               'check_captain_loses_sailors', 'departure_message', 'arrival_message'],
    'input.c': ['find_command'],
    'perm.c': ['rp_admit', 'will_admit'],
    'stack.c': ['here_pos', 'stack_parent', 'stack_leader', 'promote',
                'unstack', 'leave_stack'],
    'storm.c': ['weather_here'],
    'garr.c': ['garrison_here', 'province_admin', 'players_who_rule_here', 'rank_s',
               'garrison_notices', 'garrison_spot_check'],
    'art.c': ['has_auraculum', 'max_eff_aura'],
    'combat.c': ['find_wield', 'wield_s'],
    'day.c': ['touch_loc_pl'],
    'use.c': ['exp_level', 'rp_skill_ent', 'has_skill'],
    'display.c': ['with_inventory_string', 'display_with', 'mage_s', 'priest_s',
                  'emperor', 'liner_desc_char', 'liner_desc', 'show_chars_below'],
    'seed.c': ['int_comp'],
}
with open(dest, 'w') as out:
    out.write('#include <stdlib.h>\nint fuzzy_find;\n')
    for file, names in functions.items():
        source = open(f'{src}/{file}').read()
        if file == 'dir.c':
            # Includes, tag constants, static stream/turn guard and begin_exitdir.
            out.write('#line 1 "olympia/dir.c"\n')
            out.write(source[:source.index('int max_map_row')])
        if file == 'u.c':
            start, end = source.index('static char *num_s[]'), source.index('char *\nnice_num(')
            out.write(source[start:end])
        if file == 'combat.c':
            out.write('#include "combat-structs.h"\n')
        if file == 'display.c':
            start = source.index('int show_display_string =')
            end = source.index('\n\n', start)
            out.write('#line 402 "olympia/display.c"\n' + source[start:end] + '\n')
            # Unrelated descriptions are explicit rejected boundaries, not
            # invented successful output. Keep the actual liner_desc switch.
            for name in ['liner_desc_ship', 'liner_desc_loc', 'liner_desc_road', 'liner_desc_storm']:
                out.write(f'static char *{name}(int n) {{ fprintf(stderr, "out-of-scope C dependency: {name}\\n"); abort(); }}\n')
        for name in names:
            m = re.search(rf'^(?:static )?[^\n(]+\n{name}\(.*?^\}}\n', source, re.M | re.S)
            if not m:
                sys.exit(f'no complete definition of {name} in {file}')
            line = source.count('\n', 0, m.start()) + 1
            out.write(f'#line {line} "olympia/{file}"\n{m.group(0)}\n')
    out.write('''
int drv_land(struct command *c, struct exit_view *v, int show) { return move_exit_land(c, v, show); }
int drv_fly(struct command *c, struct exit_view *v, int show) { return move_exit_fly(c, v, show); }
int drv_water(struct command *c, struct exit_view *v, int ship, int show) { return move_exit_water(c, v, ship, show); }
int drv_check(struct command *c, struct exit_view *v, int show) { return land_check(c, v, show); }
int drv_can(int where, struct command *c) { return can_move_here(where, c); }
int drv_outer(int where, struct command *c) { return can_move_at_outer_level(where, c); }
void drv_save(struct command *c, struct exit_view *v) { save_v_array(c, v); }
void drv_restore(struct command *c, struct exit_view *v) { restore_v_array(c, v); }
void drv_suspend(int who) { suspend_stack_actions(who); }
void drv_unstack(int who) { unstack(who); }
int drv_sub_item(int who, int item, int qty) { return sub_item(who, item, qty); }
void drv_arrival(int who, struct exit_view *v) { arrival_message(who, v); }
''')
# find_command reads the actual glob.c table. Its callbacks are deliberately
# unreachable here: never make an unfinished SAIL callback a successful stub.
glob = open(f'{src}/glob.c').read()
table = glob[glob.index('struct cmd_tbl_ent cmd_tbl[]'):glob.index('char *kind_s[]')]
routines = set()
for m in re.finditer(r'^\{[^,]*,[^,]*,\s*(\w+),\s*(\w+),\s*(\w+),', table, re.M):
    routines.update(n for n in m.groups() if n != 'NULL')
with open(dest, 'a') as out:
    for name in sorted(routines):
        out.write(f'int {name}(struct command *c) {{ fprintf(stderr, "out-of-scope C dependency: {name}\\n"); abort(); }}\n')
PY
case "$(uname -s)" in
 Darwin) ldflags=(-Wl,-dead_strip) ;;
 *) ldflags=(-Wl,--gc-sections) ;;
esac
(cd "${WORK}/g3" && "${CC}" -std=c11 -ffunction-sections -fdata-sections "${ldflags[@]}" \
 -I olympia -I lib -I "${WORK}" "${HERE}/driver.c" "${WORK}/helpers.c" \
 olympia/{code,glob,z,sout,io}.c lib/{checked_alloc,ilist,exit_views,rng,rnd,item_ents,admits,cstrings,skill_ents,trades,req_ents}.c \
 -o "${WORK}/driver") >"${WORK}/build.log" 2>&1 || die "C build failed: ${WORK}/build.log"
sections=$(grep -c '^section ' "${CASES}")
section_names=()
while IFS= read -r line; do section_names+=("$line"); done < <(grep '^section ' "${CASES}")
for run in run1 run2; do
 for ((i=0; i<sections; i++)); do
  mkdir -p "${WORK}/${run}/${i}/log"
  status=0
  (run_limited "${RUN_TIMEOUT}" move-driver "${WORK}/driver" "${CASES}" "$i" "${WORK}/${run}/${i}/randseed") \
   >"${WORK}/${run}/stdout" 2>"${WORK}/${run}/stderr" || status=$?
  [ "$status" != 124 ] || die "C section $i timed out"
  if grep -q 'out-of-scope C dependency:' "${WORK}/${run}/stderr"; then
   die "C section $i reached an unfinished dependency; see ${WORK}/${run}/stderr"
  fi
  expected_status=0
  case "${section_names[$i]}" in
   'section invalid '*|'section unsupported '*) expected_status=134 ;;
   'section unowned hidden route G3 353') expected_status=134 ;;
  esac
  [ "$status" = "$expected_status" ] || die "unexpected C status $status in section $i; expected $expected_status; see ${WORK}/${run}"
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
(cd "${ROOT}" && RUN_TIMEOUT="${RUN_TIMEOUT}" go test -exec "bash ${ROOT}/tools/lib/run-limited.sh" \
 -count=1 -run '^TestMove' ./internal/olympia)
echo PASS
