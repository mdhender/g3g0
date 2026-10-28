#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"
. "${HERE}/../lib/timeout.sh"
PINNED_REV=e733acb63bacb21167a1f2589d8dcdc023364896
G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
WORK="${WORK:-${ROOT}/tmp/turn-compat}"
RUN_TIMEOUT="${RUN_TIMEOUT:-120}"
full=0 record=0
die() { echo "error: $*" >&2; exit 2; }
for arg in "$@"; do
	case "$arg" in
	--full) full=1 ;;
	--record) record=1 ;;
	*) die "usage: $0 [--full [--record]]" ;;
	esac
done
[ "$record" = 0 ] || [ "$full" = 1 ] || die "--record requires --full and two real pinned C runs"
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
	rm -rf "${out}"
	mkdir -p "${out}/cwd"
	cp -a "${input}/." "${out}/cwd/"
	chmod -R u+w "${out}/cwd"
	(cd "${out}/cwd" && run_limited "${RUN_TIMEOUT}" olympia-g3 "${bin}" "${ARGS[@]}" "$@" \
		</dev/null >"${out}/stdout" 2>"${out}/stderr") || status=$?
	echo "${status}" >"${out}/status"
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
	run_engine "${C_BIN}" "${WORK}/input" "${WORK}/c"
	run_engine "${C_BIN}" "${WORK}/input" "${WORK}/c-again"
	[ "$(cat "${WORK}/c/status")" = 0 ] || die "C turn failed (status $(cat "${WORK}/c/status")); see ${WORK}/c"
	# Check even deferred reports for reference determinism, not just state.
	diff -rq "${WORK}/c" "${WORK}/c-again" >"${WORK}/determinism.diff" ||
		die "C is nondeterministic; see ${WORK}/determinism.diff"
	{
		echo "pin ${PINNED_REV}"
		echo "args ${ARGS[*]}"
		echo "binary $(shasum -a 256 "${C_BIN}" | cut -d' ' -f1)"
		echo "world.tgz $(shasum -a 256 "${WORLD}" | cut -d' ' -f1)"
		echo "C repeats: identical status, stdout, stderr, entire cwd (including reports)"
		"${COMPARE}" --manifest "${WORK}/c"
	} >"${WORK}/c-determinism.txt"
	rm -rf "${WORK}/c-again"
	echo "PASS: two pinned C turns identical, including report artifacts"
	if [ "$record" = 1 ]; then
		run_engine "${C_BIN}" "${WORK}/empty" "${WORK}/rejected-c" '-?'
		run_engine "${C_BIN}" "${WORK}/empty" "${WORK}/rejected-c-again" '-?'
		diff -rq "${WORK}/rejected-c" "${WORK}/rejected-c-again" || die "small C case is nondeterministic"
		rm -rf "${HERE}/testdata/rejected-option" "${HERE}/testdata/world-sample"
		mkdir -p "${HERE}/testdata/world-sample/cwd/lib"
		cp -a "${WORK}/rejected-c" "${HERE}/testdata/rejected-option"
		# Small, real post-turn state, plus untouched raw process observations.
		cp "${WORK}/c/"{stdout,stderr,status} "${HERE}/testdata/world-sample/"
		cp "${WORK}/c/cwd/lib/"{system,randseed,ship} "${HERE}/testdata/world-sample/cwd/lib/"
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
(cd "${ROOT}" && RUN_TIMEOUT="${RUN_TIMEOUT}" go test -count=1 \
	-exec "bash ${HERE}/../lib/run-limited.sh" ./tools/turn-compat)
run_engine "${GO_BIN}" "${WORK}/empty" "${WORK}/rejected-go" '-?'
"${COMPARE}" "${HERE}/testdata/rejected-option" "${WORK}/rejected-go"
rm -rf "${WORK}/rejected-go"
[ "$full" = 1 ] || exit 0

run_engine "${GO_BIN}" "${WORK}/input" "${WORK}/go"
echo "turn: C status $(cat "${WORK}/c/status"), Go status $(cat "${WORK}/go/status")"
status=0
"${COMPARE}" "${WORK}/c" "${WORK}/go" >"${WORK}/diff.txt" || status=$?
{
	diff -u --label C/stdout --label Go/stdout "${WORK}/c/stdout" "${WORK}/go/stdout" || true
	diff -u --label C/stderr --label Go/stderr "${WORK}/c/stderr" "${WORK}/go/stderr" || true
} >>"${WORK}/diff.txt"
cat "${WORK}/diff.txt"
if [ "$status" = 0 ]; then
	rm -rf "${WORK}/c" "${WORK}/go"
else
	echo "FAIL: raw captures and state retained in ${WORK}/{c,go}; ${WORK}/diff.txt"
	[ "$(cat "${WORK}/go/status")" != "$TIMEOUT_STATUS" ] || echo "Go timed out after ${RUN_TIMEOUT}s"
fi
[ "$status" = 0 ] || exit "$status"

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
	run_engine "${C_BIN}" "${base}/input" "${base}/c"
	run_engine "${C_BIN}" "${base}/input" "${base}/c-again"
	[ "$(cat "${base}/c/status")" = "$expected_status" ] || die "${case_name}: unexpected C status; see ${base}"
	diff -rq "${base}/c" "${base}/c-again" >"${base}/determinism.diff" ||
		die "${case_name}: C is nondeterministic; see ${base}"
	rm -rf "${base}/c-again"
	run_engine "${GO_BIN}" "${base}/input" "${base}/go"
	if ! "${COMPARE}" "${base}/c" "${base}/go" >"${base}/diff.txt"; then
		cat "${base}/diff.txt"
		diff -u "${base}/c/stdout" "${base}/go/stdout" || true
		diff -u "${base}/c/stderr" "${base}/go/stderr" || true
		echo "FAIL: ${case_name}; raw captures retained in ${base}"
		exit 1
	fi
	printf '%s: C/C deterministic, C/Go status %s; ' "$case_name" "$expected_status" | tee -a "${WORK}/lifecycle.txt"
	cat "${base}/diff.txt" | tee -a "${WORK}/lifecycle.txt"
	rm -rf "${base}"
done
rm -rf "${WORK}/lifecycle" "${WORK}/input"
