#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"
. "${HERE}/../lib/timeout.sh"
PINNED_REV=f73e48118805ffe8ff5e763bfa4f20ae00c3a705
G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
WORK="${WORK:-${ROOT}/tmp/turn-compat}"
RUN_TIMEOUT="${RUN_TIMEOUT:-120}"
full=0 record=0 command_coverage=0 use_coverage=0 phase_coverage=0 state_export=0
die() { echo "error: $*" >&2; exit 2; }
for arg in "$@"; do
	case "$arg" in
	--full) full=1 ;;
	--record) record=1 ;;
	--command-coverage) command_coverage=1 ;;
	--use-coverage) use_coverage=1 ;;
	--phase-coverage) phase_coverage=1 ;;
	--state-export) state_export=1 ;;
	*) die "usage: $0 [--full [--record] [--command-coverage] [--use-coverage] [--phase-coverage] [--state-export]]" ;;
	esac
done
[ "$record" = 0 ] || [ "$full" = 1 ] || die "--record requires --full and two real pinned C runs"
[ "$command_coverage" = 0 ] || [ "$full" = 1 ] || die "--command-coverage requires --full"
[ "$use_coverage" = 0 ] || [ "$full" = 1 ] || die "--use-coverage requires --full"
[ "$phase_coverage" = 0 ] || [ "$full" = 1 ] || die "--phase-coverage requires --full"
[ "$state_export" = 0 ] || [ "$full" = 1 ] || die "--state-export requires --full"
mkdir -p "${WORK}/bin"
WORK="$(cd "${WORK}" && pwd -P)"
(cd "${ROOT}" && go build -o "${WORK}/bin/compare" ./tools/turn-compat &&
	go build -o "${WORK}/bin/olympia-g3" ./cmd/olympia)
COMPARE="${WORK}/bin/compare"
GO_BIN="${WORK}/bin/olympia-g3"
ARGS=(-r -l ./lib -S test-use-const-report-date)
# A retained input file makes the diagnostic cwd representable in git.
# Both engines see this same input; -? rejects the command before DB load.
rm -rf "${WORK}/empty"
mkdir -p "${WORK}/empty"
printf 'Rejected option diagnostic input.\n' >"${WORK}/empty/input.txt"
run_engine() {
	local bin="$1" input="$2" out="$3" status=0
	shift 3
	[ -d "$input" ] && [ ! -L "$input" ] || die "engine input must be an actual directory, not a symlink: $input"
	rm -rf "${out}"
	mkdir -p "${out}/cwd"
	cp -a "${input}/." "${out}/cwd/"
	chmod -R u+w "${out}/cwd"
	(cd "${out}/cwd" && run_limited "${RUN_TIMEOUT}" olympia-g3 "${bin}" "${ARGS[@]}" "$@" \
		</dev/null >"${out}/stdout" 2>"${out}/stderr") || status=$?
	echo "${status}" >"${out}/status"
}

run_c_engine() {
	run_engine "$@"
	local out="$3" marker_status=0
	# Never bless a missing C dependency, even on an expected diagnostic/failure.
	grep -Ei 'out-of-scope|not (yet )?implemented|missing.*dependency|NYI|stub|Sanitizer|runtime error:' \
		"${out}/stdout" "${out}/stderr" >&2 || marker_status=$?
	[ "$marker_status" = 1 ] || die "invalid C reference or unreadable streams; see ${out}"
}

# Validate every entry (rejecting unsupported types), then use the shared
# portable lstat policy for all cwd modes, including directories and cwd itself.
# A failed listing/capture is a harness error (2), never an ordinary mismatch (1).
compare_captures() {
	local left="$1" right="$2" bytes_status=0 modes_status=0
	"${COMPARE}" "$left" "$right" || bytes_status=$?
	[ "$bytes_status" -le 1 ] || return 2
	python3 "${HERE}/../lib/tree-modes.py" "$left/cwd" "$right/cwd" || modes_status=$?
	[ "$modes_status" -le 1 ] || return 2
	[ "$bytes_status" = 0 ] && [ "$modes_status" = 0 ]
}

if [ "$full" = 1 ]; then
	# Reuse only an exact pinned archive/build cache, never a source checkout.
	# Otherwise make our own export and build with the upstream preset.
	C_CACHE="${C_CACHE:-${ROOT}/tmp/olympia-roundtrip}"
	if [ "$(cat "${C_CACHE}/g3.rev" 2>/dev/null)" != "$PINNED_REV" ] || \
		[ ! -x "${C_CACHE}/g3/build/debug/olympia-g3" ]; then
		C_CACHE="${WORK}/reference"
		if [ "$(cat "${C_CACHE}/g3.rev" 2>/dev/null)" != "$PINNED_REV" ]; then
			rm -rf "${C_CACHE}"
			mkdir -p "${C_CACHE}/g3"
			git -C "${G3_REPO}" archive "$PINNED_REV" | tar -x -C "${C_CACHE}/g3"
			echo "$PINNED_REV" >"${C_CACHE}/g3.rev"
		fi
		(cd "${C_CACHE}/g3" && cmake --preset debug && \
			cmake --build --preset debug --target olympia-g3) >"${WORK}/build.log" 2>&1 ||
			die "C build failed; see ${WORK}/build.log"
	fi
	C_CACHE="$(cd "${C_CACHE}" && pwd -P)"
	C_BIN="${C_CACHE}/g3/build/debug/olympia-g3"
	WORLD="${C_CACHE}/g3/tests/olympia/fixtures/world.tgz"
	input_stamp="${PINNED_REV} $(shasum -a 256 "${WORLD}" | cut -d' ' -f1)"
	if [ "$(cat "${WORK}/input.stamp" 2>/dev/null)" != "$input_stamp" ] || [ ! -d "${WORK}/input" ]; then
		rm -rf "${WORK}/input"
		mkdir -p "${WORK}/input"
		tar xzf "${WORLD}" -C "${WORK}/input"
		printf '%s\n' "$input_stamp" >"${WORK}/input.stamp"
	fi
	run_c_engine "${C_BIN}" "${WORK}/input" "${WORK}/c"
	run_c_engine "${C_BIN}" "${WORK}/input" "${WORK}/c-again"
	[ "$(cat "${WORK}/c/status")" = 0 ] || die "C turn failed (status $(cat "${WORK}/c/status")); see ${WORK}/c"
	compare_captures "${WORK}/c" "${WORK}/c-again" >"${WORK}/determinism.diff" ||
		die "C is nondeterministic; see ${WORK}/determinism.diff"
	{
		echo "pin ${PINNED_REV}"
		echo "args ${ARGS[*]}"
		echo "binary $(shasum -a 256 "${C_BIN}" | cut -d' ' -f1)"
		echo "world.tgz $(shasum -a 256 "${WORLD}" | cut -d' ' -f1)"
		echo "C repeats: identical status, stdout, stderr, entire cwd (all entry types/modes/bytes)"
		"${COMPARE}" --manifest "${WORK}/c"
	} >"${WORK}/c-determinism.txt"
	rm -rf "${WORK}/c-again"
	echo "PASS: two pinned C turns identical, including report artifacts"
	if [ "$record" = 1 ]; then
		run_c_engine "${C_BIN}" "${WORK}/empty" "${WORK}/rejected-c" '-?'
		run_c_engine "${C_BIN}" "${WORK}/empty" "${WORK}/rejected-c-again" '-?'
		compare_captures "${WORK}/rejected-c" "${WORK}/rejected-c-again" || die "small C case is nondeterministic"
		rm -rf "${HERE}/testdata/rejected-option" "${HERE}/testdata/world-sample"
		mkdir -p "${HERE}/testdata/world-sample/cwd/lib"
		cp -a "${WORK}/rejected-c" "${HERE}/testdata/rejected-option"
		# Small, real post-turn state, plus untouched raw process observations.
		cp -p "${WORK}/c/"{stdout,stderr,status} "${HERE}/testdata/world-sample/"
		cp -p "${WORK}/c/cwd/lib/"{system,randseed,ship,times_0} "${HERE}/testdata/world-sample/cwd/lib/"
		cp "${WORK}/c-determinism.txt" "${HERE}/testdata/"
		{
			echo "rejected-option C repeats: identical status, stdout, stderr, entire cwd"
			"${COMPARE}" --manifest "${WORK}/rejected-c"
		} >>"${HERE}/testdata/c-determinism.txt"
		rm -rf "${WORK}/rejected-c" "${WORK}/rejected-c-again"
	fi
fi

# Fast defaults: comparator contracts on C recordings, then a real Go
# diagnostic run compared byte-for-byte with the small pinned C recording.
python3 "${HERE}/../lib/tree-modes-test.py"
(cd "${ROOT}" && RUN_TIMEOUT="${RUN_TIMEOUT}" go test -count=1 \
	-exec "bash ${HERE}/../lib/run-limited.sh" ./tools/turn-compat)
run_engine "${GO_BIN}" "${WORK}/empty" "${WORK}/rejected-go" '-?'
compare_captures "${HERE}/testdata/rejected-option" "${WORK}/rejected-go"
rm -rf "${WORK}/rejected-go"
[ "$full" = 1 ] || exit 0

run_engine "${GO_BIN}" "${WORK}/input" "${WORK}/go"
echo "turn: C status $(cat "${WORK}/c/status"), Go status $(cat "${WORK}/go/status")"
status=0
compare_captures "${WORK}/c" "${WORK}/go" >"${WORK}/diff.txt" || status=$?
{
	diff -u --label C/stdout --label Go/stdout "${WORK}/c/stdout" "${WORK}/go/stdout" || true
	diff -u --label C/stderr --label Go/stderr "${WORK}/c/stderr" "${WORK}/go/stderr" || true
} >>"${WORK}/diff.txt"
cat "${WORK}/diff.txt"
if [ "$status" = 0 ]; then
	rm -rf "${WORK}/c" "${WORK}/go"
else
	echo "FAIL: raw captures and complete cwd retained in ${WORK}/{c,go}; ${WORK}/diff.txt"
	[ "$(cat "${WORK}/go/status")" != "$TIMEOUT_STATUS" ] || echo "Go timed out after ${RUN_TIMEOUT}s"
fi
[ "$status" = 0 ] || exit "$status"

for family in command use phase state; do
	flag=--coverage
	if [ "$family" = command ]; then
		[ "$command_coverage" = 1 ] || continue
	elif [ "$family" = use ]; then
		[ "$use_coverage" = 1 ] || continue
		flag=--use-coverage
	elif [ "$family" = phase ]; then
		[ "$phase_coverage" = 1 ] || continue
		flag=--phase-coverage
	else
		[ "$state_export" = 1 ] || continue
		flag=--state-export
	fi
	cases=(load load-save turn turn-nosave duplicate missing directory blocked saved-system report)
	if [ "$family" = state ]; then
		cases+=(saved-directory seed times)
		# Real nonempty factions/topology from the unchanged pinned scenario.
		# This exercises export at each save boundary, not just an empty root.
		mkdir -p "${WORK}/state-fixtures"
		git -C "${G3_REPO}" archive "$PINNED_REV" tests/olympia/fixtures \
			tests/olympia/regress/phase-month-open/scenario.patch | tar -x -C "${WORK}/state-fixtures"
		. "${WORK}/state-fixtures/tests/olympia/fixtures/scenario.sh"
	fi
	if [ "$family" = phase ]; then
		cases+=(command-combined use-combined all-combined command-shared use-shared all-shared command-failure use-failure prior-failures dash-path parse-missing parse-empty parse-equals)
	fi
	if [ "$family" = use ]; then
		cases+=(combined shared-path command-failure shared-checks hide)
		# Pristine pinned fixture inputs, never a previous engine/cache result.
		mkdir -p "${WORK}/use-fixtures"
		git -C "${G3_REPO}" archive "$PINNED_REV" tests/olympia/fixtures \
			tests/olympia/regress/use-checks/scenario.patch \
			tests/olympia/regress/use-stealth/scenario.patch | tar -x -C "${WORK}/use-fixtures"
		. "${WORK}/use-fixtures/tests/olympia/fixtures/scenario.sh"
		. "${WORK}/use-fixtures/tests/olympia/fixtures/coverage.sh"
	fi
	: >"${WORK}/${family}-coverage.txt"
	for case_name in "${cases[@]}"; do
		base="${WORK}/${family}-coverage/${case_name}"
		rm -rf "${base}"
		mkdir -p "${base}/input"
		if [ "$family" = state ]; then
			scenario_unpack "${WORK}/state-fixtures/tests/olympia/regress/phase-month-open" "${base}/input" || die "cannot unpack phase-month-open"
		else
			cp -a "${WORK}/input/." "${base}/input/"
		fi
		ARGS=(-l ./lib test-use-const-report-date)
		expected_status=0
		path=cov.tsv
		case "$case_name" in
		load-save) ARGS+=(-S) ;;
		load) ;;
		turn-nosave) ARGS+=(-r) ;;
		*) ARGS+=(-r -S) ;;
		esac
		case "$case_name" in
		missing) path=missing/cov.tsv ;;
		directory) path=lib/log ;;
		saved-directory) path=lib/fact ;;
		seed) path=lib/randseed ;;
		times) path=lib/times_0 ;;
		blocked) path=blocked/cov.tsv; printf 'not a directory\n' >"${base}/input/blocked" ;;
		saved-system) path=lib/system ;;
		report) path=lib/players ;;
		combined) ARGS+=(--coverage commands.tsv) ;;
		shared-path) ARGS+=(--coverage "$path") ;;
		command-failure) ARGS+=(--coverage lib/log) ;;
		command-combined) ARGS+=(--coverage commands.tsv) ;;
		use-combined) ARGS+=(--use-coverage use.tsv) ;;
		all-combined) ARGS+=(--coverage commands.tsv --use-coverage use.tsv) ;;
		command-shared) ARGS+=(--coverage "$path") ;;
		use-shared) ARGS+=(--use-coverage "$path") ;;
		all-shared) ARGS+=(--coverage "$path" --use-coverage "$path") ;;
		use-failure) ARGS+=(--use-coverage lib/log) ;;
		prior-failures) ARGS+=(--coverage lib/log --use-coverage missing/use.tsv) ;;
		dash-path) path=-R ;;
		shared-checks|hide)
			rm -rf "${base}/input"
			mkdir "${base}/input"
			fixture=use-checks
			[ "$case_name" != hide ] || fixture=use-stealth
			scenario_unpack "${WORK}/use-fixtures/tests/olympia/regress/${fixture}" "${base}/input" || die "cannot unpack ${fixture}"
			if [ "$case_name" = shared-checks ]; then
				# The ten implicit callbacks remain explicitly NYI. Select only
				# the unchanged fixture's real shared-check orders (301/302).
				rm "${base}/input/lib/orders/300"
			else
				# Select the existing shadow noble's implemented hide callback;
				# no missing thief/sneak/torture/form gameplay is substituted.
				rm "${base}/input/lib/orders/"*
				printf '6947:use 638 1\n6947:use 638 0\n' >"${base}/input/lib/orders/304"
			fi
			;;
		esac
		# Baseline controls use the exact same full-engine/comparator owners.
		# They establish no flag-free state/RNG/report/stream effect, not just
		# agreement between two implementations that both changed their state.
		case "$case_name" in
		load|load-save|turn|turn-nosave|shared-checks|hide)
			run_c_engine "${C_BIN}" "${base}/input" "${base}/plain-c"
			run_c_engine "${C_BIN}" "${base}/input" "${base}/plain-again"
			[ "$(cat "${base}/plain-c/status")" = 0 ] || die "${case_name}: flag-free C failed"
			compare_captures "${base}/plain-c" "${base}/plain-again" >"${base}/plain-determinism.diff" || die "${case_name}: flag-free C nondeterminism"
			rm -rf "${base}/plain-again"
			run_engine "${GO_BIN}" "${base}/input" "${base}/plain-go"
			compare_captures "${base}/plain-c" "${base}/plain-go" >"${base}/plain.diff" || die "${case_name}: flag-free Go differs"
			rm -rf "${base}/plain-go"
			;;
		esac
		case "$case_name" in shared-checks|hide) ARGS+=(--coverage commands.tsv) ;; esac
		[ "$case_name" != duplicate ] || ARGS+=("$flag" ignored.tsv)
		case "$case_name" in
		parse-missing) ARGS+=("$flag"); expected_status=1 ;;
		parse-empty) ARGS+=("$flag" ""); expected_status=1 ;;
		# This is not long-option syntax in C: getopt scans its characters.
		# The 't' in cov.tsv returns ilist_test before checking errflag.
		parse-equals) ARGS+=("$flag=$path") ;;
		*) ARGS+=("$flag" "$path") ;;
		esac
		run_c_engine "${C_BIN}" "${base}/input" "${base}/c"
		run_c_engine "${C_BIN}" "${base}/input" "${base}/c-again"
		[ "$(cat "${base}/c/status")" = "$expected_status" ] || die "${case_name}: unexpected coverage C status"
		compare_captures "${base}/c" "${base}/c-again" >"${base}/determinism.diff" || die "${case_name}: coverage C nondeterminism"
		rm -rf "${base}/c-again"
		run_engine "${GO_BIN}" "${base}/input" "${base}/go"
		compare_captures "${base}/c" "${base}/go" >"${base}/diff.txt" || die "${case_name}: coverage Go differs; see ${base}"
		printf '%s: C/C/Go status %s, exact streams and all entries/types/modes/bytes\n' "$case_name" "$expected_status" | tee -a "${WORK}/${family}-coverage.txt"
		if [ "$family" = phase ] && [ "$case_name" = turn ]; then
			# Keep the complete ordinary-turn TSV as evidence, not just selected
			# totals. Expected bytes remain the two unchanged full C engines.
			cp "${base}/c/cwd/cov.tsv" "${WORK}/phase-turn.tsv"
			printf 'complete ordinary-turn phase TSV: ' >>"${WORK}/phase-coverage.txt"
			shasum -a 256 "${WORK}/phase-turn.tsv" | cut -d' ' -f1 >>"${WORK}/phase-coverage.txt"
		fi
		if [ "$family" = state ] && [ "$case_name" = turn ]; then
			cp "${base}/c/cwd/cov.tsv" "${WORK}/state-turn.json"
			printf 'complete ordinary-turn JSON: ' >>"${WORK}/state-coverage.txt"
			shasum -a 256 "${WORK}/state-turn.json" | cut -d' ' -f1 >>"${WORK}/state-coverage.txt"
		fi
		case "$case_name" in
		shared-checks|hide)
			if [ "$case_name" = shared-checks ]; then
				use_cells "${base}/c/cwd/cov.tsv" invalid:0:checks skill:603:checks \
					skill:631:checks skill:613:checks skill:801:checks item:8:checks item:9:checks || die "selected shared checks were not reached"
			else
				use_cells "${base}/c/cwd/cov.tsv" skill:638:ok || die "selected hide callback was not reached"
			fi
			# Supplemental observations, never a replacement for the strict raw
			# comparisons above: retain actual C nonzero cells before cleanup.
			printf '%s: C USE bytes %s\n' "$case_name" "$(wc -c <"${base}/c/cwd/cov.tsv")" >>"${WORK}/${family}-coverage.txt"
			awk -F '\t' 'NR == 1 || $4 != 0 || $5 != 0 || $6 != 0' "${base}/c/cwd/cov.tsv" >>"${WORK}/${family}-coverage.txt"
			printf '%s: command dump ' "$case_name" >>"${WORK}/${family}-coverage.txt"
			shasum -a 256 "${base}/c/cwd/commands.tsv" | cut -d' ' -f1 >>"${WORK}/${family}-coverage.txt"
			rm "${base}/c/cwd/commands.tsv" "${base}/go/cwd/commands.tsv"
			;;
		esac
		case "$case_name" in
		load|load-save|turn|turn-nosave|shared-checks|hide)
			printf '%s: C dump ' "$case_name" >>"${WORK}/${family}-coverage.txt"
			shasum -a 256 "${base}/c/cwd/cov.tsv" | cut -d' ' -f1 >>"${WORK}/${family}-coverage.txt"
			# Only after strict C/C/Go byte comparison, remove the requested
			# dumps from disposable captures to check flag-free aftermath.
			rm "${base}/c/cwd/cov.tsv" "${base}/go/cwd/cov.tsv"
			compare_captures "${base}/plain-c" "${base}/c" >"${base}/effect-c.diff" || die "${case_name}: C coverage changed other observations"
			compare_captures "${base}/plain-c" "${base}/go" >"${base}/effect-go.diff" || die "${case_name}: Go coverage changed other observations"
			echo "${case_name}: flag-free aftermath identical except requested dump" | tee -a "${WORK}/${family}-coverage.txt"
			;;
		esac
		rm -rf "${base}"
	done
	rm -rf "${WORK}/${family}-coverage"
done
rm -rf "${WORK}/use-fixtures" "${WORK}/state-fixtures"

# Exercise main's actual composition as well as the component recordings in
# command-compat. These opt-in controls change only lifecycle inputs on the
# same pinned bare world; expected bytes always come from two real C runs.
# Keep at most one case's captures, deleting matching outputs immediately.
: >"${WORK}/lifecycle.txt"
for case_name in no-save empty-act dot-act partial-act act-file \
	times_press times_rumor times_0 players email totimes forward factions; do
	base="${WORK}/lifecycle/${case_name}"
	rm -rf "${base}"
	mkdir -p "${base}/input"
	cp -a "${WORK}/input/." "${base}/input/"
	ARGS=(-r -l ./lib -S test-use-const-report-date)
	expected_status=0
	case "$case_name" in
	no-save) ARGS=(-r -l ./lib test-use-const-report-date) ;;
	empty-act) mkdir "${base}/input/act" ;;
	dot-act)
		mkdir -p "${base}/input/act/.ignored"
		printf 'faction\ncharacter\ncity\nname\nemail\n' >"${base}/input/act/.ignored/Join-g3"
		;;
	partial-act)
		mkdir -p "${base}/input/act/aa0"
		printf 'faction\ncharacter\ncity\nname\n' >"${base}/input/act/aa0/Join-g3"
		;;
	act-file) printf 'not a directory\n' >"${base}/input/act" ;;
	*)
		rm -f "${base}/input/lib/${case_name}"
		mkdir "${base}/input/lib/${case_name}"
		case "$case_name" in times_*) expected_status=1 ;; esac
		;;
	esac
	run_c_engine "${C_BIN}" "${base}/input" "${base}/c"
	run_c_engine "${C_BIN}" "${base}/input" "${base}/c-again"
	[ "$(cat "${base}/c/status")" = "$expected_status" ] || die "${case_name}: unexpected C status; see ${base}"
	compare_captures "${base}/c" "${base}/c-again" >"${base}/determinism.diff" ||
		die "${case_name}: C is nondeterministic; see ${base}"
	rm -rf "${base}/c-again"
	run_engine "${GO_BIN}" "${base}/input" "${base}/go"
	status=0
	compare_captures "${base}/c" "${base}/go" >"${base}/diff.txt" || status=$?
	if [ "$status" != 0 ]; then
		cat "${base}/diff.txt"
		diff -u "${base}/c/stdout" "${base}/go/stdout" || true
		diff -u "${base}/c/stderr" "${base}/go/stderr" || true
		echo "FAIL: ${case_name}; raw captures retained in ${base}"
		exit "$status"
	fi
	printf '%s: C/C deterministic, C/Go status %s; ' "$case_name" "$expected_status" | tee -a "${WORK}/lifecycle.txt"
	cat "${base}/diff.txt" | tee -a "${WORK}/lifecycle.txt"
	rm -rf "${base}"
done
rm -rf "${WORK}/lifecycle" "${WORK}/input"
