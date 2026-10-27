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
G3_REV="${G3_REV:-6a25f2d252ce45bd5adbe3baaccb3be3eed56279}"
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
python3 -I - "${WORK}/g3/olympia" "${WORK}/helpers.c" <<'PY'
import re, sys
src, dest = sys.argv[1:]
functions = {
    'dir.c': ['location_direction', 'exit_distance', 'is_port_city',
              'province_has_port_city', 'summer_uldim_open_now', 'add_province_exit',
              'extra_routes', 'province_exits', 'province_sub_exits', 'subloc_exits',
              'ship_exits', 'building_exits', 'exits_from_loc', 'exits_from_loc_nsew',
              'exits_from_loc_nsew_select', 'has_ocean_access'],
    'loc.c': ['region', 'province'],
    'u.c': ['loc_depth', 'loc_hidden', 'test_known', 'test_bit', 'set_bit'],
    'swear.c': ['lord', 'player'],
}
with open(dest, 'w') as out:
    for file, names in functions.items():
        source = open(f'{src}/{file}').read()
        if file == 'dir.c':
            # Includes, tag constants, static stream/turn guard and begin_exitdir.
            out.write('#line 1 "olympia/dir.c"\n')
            out.write(source[:source.index('int max_map_row')])
        for name in names:
            m = re.search(rf'^(?:static )?[^\n(]+\n{name}\(.*?^\}}\n', source, re.M | re.S)
            if not m:
                sys.exit(f'no complete definition of {name} in {file}')
            line = source.count('\n', 0, m.start()) + 1
            out.write(f'#line {line} "olympia/{file}"\n{m.group(0)}\n')
PY
case "$(uname -s)" in
 Darwin) ldflags=(-Wl,-dead_strip) ;;
 *) ldflags=(-Wl,--gc-sections) ;;
esac
(cd "${WORK}/g3" && "${CC}" -std=c11 -ffunction-sections -fdata-sections "${ldflags[@]}" \
 -I olympia -I lib "${HERE}/driver.c" "${WORK}/helpers.c" \
 olympia/{code,glob,z}.c lib/{checked_alloc,ilist,exit_views,rng,rnd}.c \
 -o "${WORK}/driver") >"${WORK}/build.log" 2>&1 || die "C build failed: ${WORK}/build.log"
sections=$(grep -c '^section ' "${CASES}")
section_names=()
while IFS= read -r line; do section_names+=("$line"); done < <(grep '^section ' "${CASES}")
for run in run1 run2; do
 for ((i=0; i<sections; i++)); do
  mkdir -p "${WORK}/${run}/${i}"
  status=0
  (run_limited "${RUN_TIMEOUT}" move-driver "${WORK}/driver" "${CASES}" "$i" "${WORK}/${run}/${i}/randseed") \
   >"${WORK}/${run}/stdout" 2>"${WORK}/${run}/stderr" || status=$?
  [ "$status" != 124 ] || die "C section $i timed out"
  expected_status=0
  case "${section_names[$i]}" in
   'section invalid '*|'section unsupported '*) expected_status=134 ;;
   'section unowned hidden route G3 353') expected_status=139 ;;
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
