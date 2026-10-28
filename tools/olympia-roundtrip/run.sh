#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: load and save with no turn, G3's olympia-g3 (C, the
# reference) against cmd/olympia (Go).
#
# Each case is a saved game: G3's shared base world (world.tgz), every
# regress fixture's saved game (scenario.patch over world.tgz, unpacked with
# G3's own scenario_unpack, or the fixture's scenario.tgz), and each of those
# after one successful C turn (-r -S). Pinned expected-abort turns are
# checked separately; they cannot supply saved post-turn games. Both engines run
#
#     olympia-g3 -l ./lib -S test-use-const-report-date
#
# on identical copies of the saved game, and the whole directory (every file
# under it), stdout, stderr and the exit status must match byte for byte.
# test-use-const-report-date keeps times_0's masthead date fixed, so the C
# results do not depend on the day they were made.
#
# The C side is cached: the export, the C build, the case inputs and C's
# results are kept under $WORK and reused while G3_REV and the C binary are
# unchanged. Only the Go engine is rebuilt and rerun each time.
#
# When it makes a case's C result, C runs a second time on a few fixed
# cases (DETERMINISM_SAMPLE), and the case is an error if the two runs
# differ. --determinism-all does this for every selected case, cached or not.
#
# Usage: tools/olympia-roundtrip/run.sh [--fresh] [--determinism-all] [--list] [pattern...]
#
#   pattern   run only the cases whose name matches one of these shell
#             patterns, e.g. world, move-go, 'move-*' or '*+turn'. A case is
#             named after its saved game (world, or a fixture's directory);
#             the same game after a C turn adds +turn.
#   --fresh   discard the cached export, C build and C results first
#   --determinism-all
#             check that C gives the same result twice on every selected
#             case, not just on DETERMINISM_SAMPLE, including cases already
#             cached without that check
#   --list    print available case names and exit (no builds or engine runs)
#
# Environment:
#   G3_REPO      G3 clone that contains G3_REV
#                (default ~/Software/playbymail/olympia-g3)
#   G3_REV       G3 revision to test against (default: the pinned revision)
#   WORK         scratch directory (default tmp/olympia-roundtrip)
#   RUN_TIMEOUT  seconds allowed per program run (default 120)
#
# Exits 0 only if every selected case/control passes, 1 if any fails or could
# not be prepared, and 2 if the harness itself fails. Inputs, C outputs and
# diffs are kept under $WORK/cases/<case>, with Go outputs only on failure.
# Expected-abort captures are under $WORK/expected-aborts; the summary is
# written to $WORK/results.tsv.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"
# shellcheck source=../lib/timeout.sh
. "${HERE}/../lib/timeout.sh"

PINNED_REV=e733acb63bacb21167a1f2589d8dcdc023364896
G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-${PINNED_REV}}"
WORK="${WORK:-${ROOT}/tmp/olympia-roundtrip}"
RUN_TIMEOUT="${RUN_TIMEOUT:-120}"

# What both engines run on each case, and the C turn that makes a +turn case.
RT_ARGS=(-l ./lib -S test-use-const-report-date)
TURN_ARGS=(-r -l ./lib -S test-use-const-report-date)
# Bump when the way inputs or C results are made changes, to drop the cache.
CACHE_FORMAT=1
# The cases whose C round trip runs twice when it is made, to check that C
# is deterministic: the base world, a scenario.patch game, the scenario.tgz
# game, and each after a turn.
DETERMINISM_SAMPLE=(world world+turn move-go move-go+turn trade-routes trade-routes+turn)

die() { echo "error: $*" >&2; exit 2; }
sha() { shasum -a 256 "$1" | cut -d' ' -f1; }

# run_as BIN ARGS... runs BIN with argv[0] set to olympia-g3 and stdin from
# /dev/null, killing it after RUN_TIMEOUT seconds (exit status TIMEOUT_STATUS;
# see tools/lib/timeout.sh).
run_as() {
	run_limited "${RUN_TIMEOUT}" olympia-g3 "$@" </dev/null
}

fresh=0 list=0 determinism_all=0
patterns=()
for a in "$@"; do
	case "${a}" in
	--update) die "this harness never re-baselines; expected output only comes from the C engine" ;;
	--fresh) fresh=1 ;;
	--determinism-all) determinism_all=1 ;;
	--list) list=1 ;;
	-*) die "unknown option ${a} (usage: $0 [--fresh] [--determinism-all] [--list] [pattern...])" ;;
	*) patterns+=("${a}") ;;
	esac
done

for tool in git cmake go perl shasum patch diff cmp tar; do
	command -v "${tool}" >/dev/null || die "missing required tool: ${tool}"
done
[ -d "${G3_REPO}" ] || die "no G3 clone at ${G3_REPO}; clone https://github.com/playbymail/olympia-g3 and set G3_REPO"
git -C "${G3_REPO}" cat-file -e "${G3_REV}^{commit}" 2>/dev/null ||
	die "${G3_REPO} does not contain ${G3_REV}; fetch it or set G3_REPO"
G3_REV="$(git -C "${G3_REPO}" rev-parse "${G3_REV}^{commit}")"

mkdir -p "${WORK}"
WORK="$(cd "${WORK}" && pwd -P)"
G3="${WORK}/g3"
CASES="${WORK}/cases"
LOG="${WORK}/build.log"

if [ "${fresh}" -eq 1 ] || [ "$(cat "${WORK}/g3.rev" 2>/dev/null)" != "${G3_REV}" ]; then
	chmod -R u+w "${WORK}"
	rm -rf "${G3}" "${CASES}" "${WORK}/g3.rev" "${WORK}/c.stamp"
fi
mkdir -p "${WORK}/bin" "${CASES}"

# The reference comes from a clean export of the revision, never from the
# clone's working tree or build, and is kept for later runs.
if [ ! -f "${WORK}/g3.rev" ]; then
	echo "exporting G3 to ${G3}" >&2
	mkdir -p "${G3}"
	git -C "${G3_REPO}" archive "${G3_REV}" | tar -x -C "${G3}"
	echo "${G3_REV}" >"${WORK}/g3.rev"
fi

# ---------------------------------------------------------------------------
# Cases.
# ---------------------------------------------------------------------------
REGRESS="${G3}/tests/olympia/regress"
# shellcheck source=/dev/null
. "${G3}/tests/olympia/fixtures/scenario.sh"

games=(world)
for d in "${REGRESS}"/*/; do
	n="$(basename "${d}")"
	if [ -f "${d}scenario.patch" ] || [ -f "${d}scenario.tgz" ]; then games+=("${n}"); fi
done
all=()
# Exact-pin contract, not a general allowance for failed turns. Source:
# e733acb tests/olympia/regress/release-swear-np/{check.sh,README.md}
# (G3 #357): status 134, swear.c assertion, no sanitizer report. At older
# or override pins keep the ordinary +turn path, requiring status 0.
unavailable=""
if [ "${G3_REV}" = e733acb63bacb21167a1f2589d8dcdc023364896 ]; then
	unavailable=release-swear-np+turn
fi
for g in "${games[@]}"; do
	all+=("${g}")
	[ "${g}+turn" = "${unavailable}" ] || all+=("${g}+turn")
done

selected=()
for c in "${all[@]}"; do
	if [ ${#patterns[@]} -eq 0 ]; then selected+=("${c}"); continue; fi
	for p in "${patterns[@]}"; do
		# shellcheck disable=SC2053
		if [[ "${c}" == ${p} ]]; then selected+=("${c}"); break; fi
	done
done
# Selecting its ordinary game also validates the absent snapshot's contract.
abort_selected=0
if [ -n "${unavailable}" ]; then
	if [ ${#patterns[@]} -eq 0 ]; then abort_selected=1; fi
	for p in "${patterns[@]}"; do
		# shellcheck disable=SC2053
		if [[ "${unavailable}" == ${p} ]] || [[ release-swear-np == ${p} ]]; then abort_selected=1; fi
	done
	if [ "${abort_selected}" -eq 1 ]; then
		echo "unavailable: ${unavailable}: pinned C intentionally aborts; no saved post-turn game" >&2
	fi
fi
[ ${#selected[@]} -gt 0 ] || die "no available case matches: ${patterns[*]} (see --list)"
if [ "${list}" -eq 1 ]; then
	printf '%s\n' "${selected[@]+"${selected[@]}"}"
	exit 0
fi

echo "G3 revision: ${G3_REV}"
[ "${G3_REV}" = "${PINNED_REV}" ] || echo "             (not the pinned ${PINNED_REV})"
echo "G3 source:   ${G3_REPO}"
echo "work dir:    ${WORK}"
echo

# ---------------------------------------------------------------------------
# Build the reference (once per export) and the port (every run).
# ---------------------------------------------------------------------------
: >"${LOG}"
C_BIN="${G3}/build/debug/olympia-g3"
if [ ! -x "${C_BIN}" ]; then
	echo "building G3 (log: ${LOG})"
	(
		cd "${G3}"
		cmake --preset debug
		cmake --build --preset debug --target olympia-g3
	) >>"${LOG}" 2>&1 || die "G3 build failed; see ${LOG}"
else
	echo "reusing the G3 build in ${G3}/build/debug"
fi

# The cached inputs and C results are valid for one C binary and one way of
# making them.
stamp="$(printf '%s\n' "format ${CACHE_FORMAT}" "rev ${G3_REV}" "c $(sha "${C_BIN}")" \
	"rt ${RT_ARGS[*]}" "turn ${TURN_ARGS[*]}")"
if [ "$(cat "${WORK}/c.stamp" 2>/dev/null)" != "${stamp}" ]; then
	rm -rf "${CASES}"
	mkdir -p "${CASES}"
	printf '%s\n' "${stamp}" >"${WORK}/c.stamp"
fi

echo "building the Go port"
(
	cd "${ROOT}"
	go build -o "${WORK}/bin/olympia-g3" ./cmd/olympia
) >>"${LOG}" 2>&1 || die "Go build failed; see ${LOG}"
GO_BIN="${WORK}/bin/olympia-g3"

# unpack_game GAME DIR extracts a saved game into the empty directory DIR.
unpack_game() {
	local g="$1" dir="$2"
	if [ "${g}" = world ]; then
		tar zxf "${G3}/tests/olympia/fixtures/world.tgz" -C "${dir}"
	elif [ -f "${REGRESS}/${g}/scenario.patch" ]; then
		scenario_unpack "${REGRESS}/${g}" "${dir}"
	else
		tar zxf "${REGRESS}/${g}/scenario.tgz" -C "${dir}"
	fi
}

# run_engine BIN SRC OUT ARGS... copies SRC to OUT/cwd and runs BIN ARGS...
# there, with OUT/{stdout,stderr,status}.
run_engine() {
	local bin="$1" src="$2" out="$3" ec=0
	shift 3
	rm -rf "${out}"
	mkdir -p "${out}"
	cp -a "${src}" "${out}/cwd"
	chmod -R u+w "${out}/cwd"
	(cd "${out}/cwd" && run_as "${bin}" "$@" >"${out}/stdout" 2>"${out}/stderr") || ec=$?
	echo "${ec}" >"${out}/status"
}

# same_run A B succeeds if two run_engine outputs match exactly.
same_run() {
	cmp -s "$1/status" "$2/status" && cmp -s "$1/stdout" "$2/stdout" &&
		cmp -s "$1/stderr" "$2/stderr" && diff -r "$1/cwd" "$2/cwd" >/dev/null 2>&1
}

# Require the precise pinned assertion, not just SIGABRT. Reject missing
# dependency markers before accepting a failure, even alongside the assertion.
expected_abort() {
	[ "$(cat "$1/status")" = 134 ] &&
		grep -q 'assertion failure: .*/swear\.c (180): FALSE$' "$1/stderr" &&
		! grep -Ei 'AddressSanitizer|runtime error|out-of-scope|not yet implemented|missing.*dependency|C stub' "$1/stdout" "$1/stderr" >/dev/null
}

# Always execute twice: failed turns are controls, never cached saved games.
# Keep raw diagnostics but remove the matching partial working directories.
check_expected_abort() {
	local dir="${WORK}/expected-aborts/release-swear-np+turn"
	mkdir -p "${dir}"
	rm -f "${dir}/error" "${dir}/deterministic"
	if ! prepare release-swear-np; then
		echo "the saved game release-swear-np could not be prepared" >"${dir}/error"
		return 1
	fi
	run_engine "${C_BIN}" "${CASES}/release-swear-np/input" "${dir}/turn" "${TURN_ARGS[@]}"
	if ! expected_abort "${dir}/turn"; then
		echo "expected C status 134 and swear.c (180): FALSE without dependency/sanitizer errors; got $(cat "${dir}/turn/status"); see ${dir}/turn" >"${dir}/error"
		return 1
	fi
	run_engine "${C_BIN}" "${CASES}/release-swear-np/input" "${dir}/turn-again" "${TURN_ARGS[@]}"
	if ! expected_abort "${dir}/turn-again" || ! same_run "${dir}/turn" "${dir}/turn-again"; then
		echo "the expected C abort is not deterministic; compare ${dir}/turn and ${dir}/turn-again" >"${dir}/error"
		return 1
	fi
	rm -rf "${dir}/turn/cwd" "${dir}/turn-again/cwd"
	touch "${dir}/deterministic"
}

# check_determinism CASE says whether CASE's C round trip should run twice.
check_determinism() {
	local s
	[ "${determinism_all}" -eq 1 ] && return 0
	for s in "${DETERMINISM_SAMPLE[@]}"; do
		[ "$1" = "${s}" ] && return 0
	done
	return 1
}

# same_c CASE runs C's round trip on CASES/CASE/input again and succeeds if
# it matches CASES/CASE/c, marking the case "deterministic". If not, it
# writes the reason to CASES/CASE/error and fails.
same_c() {
	local dir="${CASES}/$1"
	run_engine "${C_BIN}" "${dir}/input" "${dir}/c-again" "${RT_ARGS[@]}"
	if ! same_run "${dir}/c" "${dir}/c-again"; then
		echo "the C round trip is not deterministic; compare ${dir}/c and ${dir}/c-again" >"${dir}/error"
		rm -f "${dir}/ready"
		return 1
	fi
	rm -rf "${dir}/c-again"
	touch "${dir}/deterministic"
}

# prepare CASE makes CASES/CASE/input (the saved game) and CASES/CASE/c (the
# C round trip) unless they are cached, and checks C's determinism if
# check_determinism says to and it has not been checked. On failure it
# writes the reason to CASES/CASE/error and fails.
prepare() {
	local c="$1" dir="${CASES}/$1" g
	[ -f "${dir}/error" ] && return 1
	if [ -f "${dir}/ready" ]; then
		if check_determinism "${c}" && [ ! -f "${dir}/deterministic" ]; then
			same_c "${c}" || return 1
		fi
		return 0
	fi
	rm -rf "${dir}"
	mkdir -p "${dir}"
	case "${c}" in
	*+turn)
		g="${c%+turn}"
		if ! prepare "${g}"; then
			echo "the saved game ${g} could not be prepared" >"${dir}/error"
			return 1
		fi
		run_engine "${C_BIN}" "${CASES}/${g}/input" "${dir}/turn" "${TURN_ARGS[@]}"
		if [ "$(cat "${dir}/turn/status")" != 0 ]; then
			echo "the C turn exited $(cat "${dir}/turn/status"); see ${dir}/turn" >"${dir}/error"
			return 1
		fi
		mv "${dir}/turn/cwd" "${dir}/input"
		;;
	*)
		mkdir -p "${dir}/input"
		if ! unpack_game "${c}" "${dir}/input" >"${dir}/unpack.log" 2>&1; then
			{ echo "unpacking failed:"; cat "${dir}/unpack.log"; } >"${dir}/error"
			return 1
		fi
		rm -f "${dir}/unpack.log"
		;;
	esac
	run_engine "${C_BIN}" "${dir}/input" "${dir}/c" "${RT_ARGS[@]}"
	# A difference between C's own runs would make every Go comparison
	# meaningless. Checking a few cases is enough to catch it, unless
	# --determinism-all asks for every case.
	if check_determinism "${c}"; then
		same_c "${c}" || return 1
	fi
	touch "${dir}/ready"
}

RESULTS="${WORK}/results.tsv"
printf 'case\tresult\tc_status\tgo_status\tstdout\tstderr\tfiles_differing\tfirst_files\n' >"${RESULTS}"

npass=0 nfail=0 nerror=0 nabort=0
start=${SECONDS}
if [ "${abort_selected}" -eq 1 ]; then
	if check_expected_abort; then
		nabort=1
		echo "  EXPECTED_ABORT ${unavailable}: C status 134, swear.c (180): FALSE; deterministic; no snapshot"
		printf '%s\tEXPECTED_ABORT\t134\t\t\t\t\tno saved snapshot\n' "${unavailable}" >>"${RESULTS}"
	else
		nerror=$((nerror + 1))
		reason="$(head -1 "${WORK}/expected-aborts/${unavailable}/error")"
		echo "  ERROR ${unavailable}: ${reason}"
		abort_status="$(cat "${WORK}/expected-aborts/${unavailable}/turn/status" 2>/dev/null || true)"
		printf '%s\tERROR\t%s\t\t\t\t\t%s\n' "${unavailable}" "${abort_status}" "${reason}" >>"${RESULTS}"
	fi
fi
echo
echo "round trip: olympia-g3 ${RT_ARGS[*]}  (${#selected[@]} cases)"
for c in "${selected[@]}"; do
	dir="${CASES}/${c}"
	if ! prepare "${c}"; then
		nerror=$((nerror + 1))
		printf '  ERROR %s: %s\n' "${c}" "$(head -1 "${dir}/error")"
		printf '%s\tERROR\t\t\t\t\t\t%s\n' "${c}" "$(head -1 "${dir}/error")" >>"${RESULTS}"
		continue
	fi
	run_engine "${GO_BIN}" "${dir}/input" "${dir}/go" "${RT_ARGS[@]}"

	c_status="$(cat "${dir}/c/status")"
	go_status="$(cat "${dir}/go/status")"
	what=()
	[ "${c_status}" = "${go_status}" ] || what+=("status C ${c_status} Go ${go_status}")
	if [ "${go_status}" = "${TIMEOUT_STATUS}" ]; then what+=("Go timed out after ${RUN_TIMEOUT}s"); fi
	so=same se=same
	cmp -s "${dir}/c/stdout" "${dir}/go/stdout" || { so=differs; what+=(stdout); }
	cmp -s "${dir}/c/stderr" "${dir}/go/stderr" || { se=differs; what+=(stderr); }
	# Files that differ, or are only on one side, relative to the run dir.
	files=()
	while IFS= read -r line; do
		files+=("${line}")
	done < <(diff -rq "${dir}/c/cwd" "${dir}/go/cwd" 2>&1 |
		sed -e "s|^Files ${dir}/c/cwd/\(.*\) and ${dir}/go/cwd/.* differ\$|\1|" \
			-e "s|^Only in ${dir}/c/cwd/\(.*\): \(.*\)\$|\1/\2 (only C)|" \
			-e "s|^Only in ${dir}/c/cwd: \(.*\)\$|\1 (only C)|" \
			-e "s|^Only in ${dir}/go/cwd/\(.*\): \(.*\)\$|\1/\2 (only Go)|" \
			-e "s|^Only in ${dir}/go/cwd: \(.*\)\$|\1 (only Go)|" |
		LC_ALL=C sort)
	first=""
	if [ ${#files[@]} -gt 0 ]; then first="${files[*]:0:5}"; fi
	if [ ${#files[@]} -gt 5 ]; then first="${first} +$((${#files[@]} - 5)) more"; fi
	[ ${#files[@]} -eq 0 ] || what+=("${#files[@]} file(s): ${first}")

	{
		echo "case ${c}: status C ${c_status} Go ${go_status}"
		diff -u --label c/stdout --label go/stdout "${dir}/c/stdout" "${dir}/go/stdout" || true
		diff -u --label c/stderr --label go/stderr "${dir}/c/stderr" "${dir}/go/stderr" || true
		(cd "${dir}" && diff -ru c/cwd go/cwd) || true
	} >"${dir}/diff.txt"

	if [ ${#what[@]} -eq 0 ]; then
		result=PASS
		npass=$((npass + 1))
		printf '  PASS  %s\n' "${c}"
		rm -rf "${dir}/go"
	else
		result=FAIL
		nfail=$((nfail + 1))
		printf '  FAIL  %s: ' "${c}"
		(IFS=';'; printf '%s\n' "${what[*]}" | sed 's/;/; /g')
	fi
	printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' "${c}" "${result}" "${c_status}" "${go_status}" \
		"${so}" "${se}" "${#files[@]}" "${first}" >>"${RESULTS}"
done

echo
echo "pass ${npass}  fail ${nfail}  error ${nerror}  expected-abort ${nabort}  ($((SECONDS - start))s)"
echo "results: ${RESULTS}"
echo "per case: ${CASES}/<case>/{input,c,go,diff.txt}"
[ "${nfail}" -eq 0 ] && [ "${nerror}" -eq 0 ]
