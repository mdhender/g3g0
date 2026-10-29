#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"
. "${HERE}/../lib/timeout.sh"
RUN_TIMEOUT="${RUN_TIMEOUT:-10}"
G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-9a03d126b329ed466c4af48425553c025eaca26f}"
WORK="${WORK:-${ROOT}/tmp/turn-driver-compat}"
CC="${CC:-cc}"
CASES="${ROOT}/internal/olympia/testdata/turn-driver-cases.txt"
EXPECT="${ROOT}/internal/olympia/testdata/turn-driver-c-output.txt"
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
echo "G3 revision: ${G3_REV}"
git -C "${G3_REPO}" archive "${G3_REV}" lib olympia | tar -x -C "${WORK}/g3"

# Like helper-compat, extract real definitions with their original locations.
# No hand-written substitutes for clock, progress or calendar behavior.
python3 -I - "${WORK}/g3/olympia/u.c" "${WORK}/helpers.c" <<'PY'
import re, sys
text = open(sys.argv[1]).read()
with open(sys.argv[2], "w") as out:
    out.write('#include <stdio.h>\n#include <time.h>\n#include "z.h"\n#include "oly.h"\n')
    # print_dot reads this file-scope initializer as well as its own body.
    names = [r"static int dot_count = [^;]+;\n"]
    names += [rf"^void\n{name}\(.*?^\}}\n" for name in
              ("olytime_increment", "olytime_turn_change", "print_dot", "stage")]
    for pattern in names:
        m = re.search(pattern, text, re.M | re.S)
        if not m:
            sys.exit(f"definition not found: {pattern}")
        line = text.count('\n', 0, m.start()) + 1
        out.write(f'#line {line} "olympia/u.c"\n{m.group(0)}\n')
PY
"${CC}" -std=c11 -Wall -Wextra -I "${WORK}/g3/olympia" \
	"${HERE}/driver.c" "${WORK}/helpers.c" -o "${WORK}/driver" \
	>"${WORK}/build.log" 2>&1 || die "C driver build failed; see ${WORK}/build.log"
sections=$(grep -c '^section ' "${CASES}")
for run in run1 run2; do
	for ((i=0; i<sections; i++)); do
		status=0
		run_limited "${RUN_TIMEOUT}" turn-driver "${WORK}/driver" "${CASES}" "$i" \
			>"${WORK}/${run}/stdout" 2>"${WORK}/${run}/stderr" || status=$?
		[ "$status" != 124 ] || die "C section $i timed out"
		! grep -q 'out-of-scope C dependency:' "${WORK}/${run}/stderr" || die "C section $i reached an unfinished dependency"
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
	-count=1 -run '^TestTurnDriverMatchesC$' ./internal/olympia)

# Isolated phase callbacks run the actual C orchestration bodies, and the
# actual Go bodies under a test-only source overlay. No production hooks.
python3 -I "${HERE}/prepare-phases.py" "${ROOT}" "${WORK}"
case "$(uname -s)" in
 Darwin) ldflags=(-Wl,-dead_strip) ;;
 *) ldflags=(-Wl,--gc-sections) ;;
esac
"${CC}" -std=c11 -Wall -Wextra -ffunction-sections -fdata-sections \
	-I "${WORK}/g3/olympia" -I "${WORK}/g3/lib" \
	"${HERE}/phase-driver.c" "${WORK}/phases.c" "${WORK}/helpers.c" \
	"${WORK}/g3/lib/ilist.c" "${ldflags[@]}" -o "${WORK}/phase-driver" \
	>"${WORK}/phase-build.log" 2>&1 || die "C phase driver build failed; see ${WORK}/phase-build.log"
CASES="${ROOT}/internal/olympia/testdata/turn-phase-cases.txt"
EXPECT="${ROOT}/internal/olympia/testdata/turn-phase-c-output.txt"
sections=$(grep -c '^section ' "${CASES}")
for run in run1 run2; do
	for ((i=0; i<sections; i++)); do
		status=0
		run_limited "${RUN_TIMEOUT}" turn-phases "${WORK}/phase-driver" "${CASES}" "$i" \
			>"${WORK}/${run}/stdout" 2>"${WORK}/${run}/stderr" || status=$?
		[ "$status" != 124 ] || die "C phase section $i timed out"
		! grep -q 'out-of-scope C dependency:' "${WORK}/${run}/stderr" || die "C phase section $i reached an unfinished dependency"
		[ "$status" = 0 ] || die "C phase section $i failed ($status)"
		{
			printf 'stdout:\n'; cat "${WORK}/${run}/stdout"
			printf 'stderr:\n'; cat "${WORK}/${run}/stderr"
			printf '\nexit %d\n' "$status"
		} >>"${WORK}/${run}/phases.txt"
	done
done
cmp "${WORK}/run1/phases.txt" "${WORK}/run2/phases.txt" || die "C phase trace is not deterministic"
echo "C phase callbacks: ${sections} sections, $(wc -l <"${WORK}/run1/phases.txt") lines, deterministic (NOT gameplay acceptance)"
if [ "$update" = 1 ]; then
	cp "${WORK}/run1/phases.txt" "${EXPECT}"
else
	diff -u "${EXPECT}" "${WORK}/run1/phases.txt"
fi
(cd "${ROOT}" && G3G0_PHASE_TRACE=1 RUN_TIMEOUT="${RUN_TIMEOUT}" go test \
	-overlay "${WORK}/overlay.json" -exec "bash ${ROOT}/tools/lib/run-limited.sh" \
	-count=1 -run '^TestTurnPhaseMatchesC$' ./internal/olympia)

# Real calendar/weather definitions and original static schedule storage, not
# the preceding synthetic phase fixture. Leaf gameplay callbacks ONLY trace.
"${CC}" -std=c11 -Wall -Wextra -ffunction-sections -fdata-sections \
	-I "${WORK}" -I "${WORK}/g3/olympia" -I "${WORK}/g3/lib" \
	"${HERE}/schedule-driver.c" "${WORK}/g3/olympia/glob.c" \
	"${WORK}/g3/olympia/rnd.c" "${WORK}/g3/lib/rng.c" \
	"${WORK}/g3/lib/ilist.c" "${ldflags[@]}" -o "${WORK}/schedule-driver" \
	>"${WORK}/schedule-build.log" 2>&1 || die "C schedule driver build failed; see ${WORK}/schedule-build.log"
for suite in calendar schedule; do
	CASES="${ROOT}/internal/olympia/testdata/${suite}-cases.txt"
	EXPECT="${ROOT}/internal/olympia/testdata/${suite}-c-output.txt"
	sections=$(grep -c '^section ' "${CASES}")
	for run in run1 run2; do
		for ((i=0; i<sections; i++)); do
			status=0
			(cd "${WORK}/${run}" && run_limited "${RUN_TIMEOUT}" turn-schedule \
				"${WORK}/schedule-driver" "${CASES}" "$i") \
				>"${WORK}/${run}/stdout" 2>"${WORK}/${run}/stderr" || status=$?
			[ "$status" != 124 ] || die "C $suite section $i timed out"
			! grep -q 'out-of-scope C dependency:' "${WORK}/${run}/stderr" || die "C $suite section $i reached an unfinished dependency"
			[ "$status" = 0 ] || die "C $suite section $i failed ($status)"
			{
				printf 'stdout:\n'; cat "${WORK}/${run}/stdout"
				printf 'stderr:\n'; cat "${WORK}/${run}/stderr"
				printf '\nexit %d\n' "$status"
			} >>"${WORK}/${run}/${suite}.txt"
		done
	done
	cmp "${WORK}/run1/${suite}.txt" "${WORK}/run2/${suite}.txt" || die "C $suite is not deterministic"
	echo "C $suite: ${sections} sections, $(wc -l <"${WORK}/run1/${suite}.txt") lines, deterministic (schedule ONLY, NOT gameplay acceptance)"
	if [ "$update" = 1 ]; then
		cp "${WORK}/run1/${suite}.txt" "${EXPECT}"
	else
		diff -u "${EXPECT}" "${WORK}/run1/${suite}.txt"
	fi
done
(cd "${ROOT}" && G3G0_SCHEDULE_TRACE=1 RUN_TIMEOUT="${RUN_TIMEOUT}" go test \
	-overlay "${WORK}/schedule-overlay.json" -exec "bash ${ROOT}/tools/lib/run-limited.sh" \
	-count=1 -run '^(TestCalendarMatchesC|TestDailyScheduleMatchesC)$' ./internal/olympia)
echo PASS
