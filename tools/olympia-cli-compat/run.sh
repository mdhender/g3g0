#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: the command line of G3's olympia-g3 (C, the reference)
# against cmd/olympia (Go).
#
# Exports the pinned G3 revision, builds its olympia-g3, builds the Go port,
# then runs every argument vector in cases.txt with both. Each run starts in
# an empty directory with stdin from /dev/null and argv[0] set to
# olympia-g3, and runs twice: once with stdout and stderr in separate files,
# and once with both in one file (2>&1), which shows their order. For a
# "same" case, Go must match C byte for byte on all three files and on the
# exit status. For an "nyi" case (a mode the port has not implemented yet),
# Go must exit nonzero, say "not yet implemented" on stderr, and print the
# start of what C prints on stdout.
#
# Finally it runs G3's rng-selftest gate, unchanged, from the export against
# both engines.
#
# Usage: tools/olympia-cli-compat/run.sh
#
# Environment:
#   G3_REPO      G3 clone that contains G3_REV
#                (default ~/Software/playbymail/olympia-g3)
#   G3_REV       G3 revision to test against (default: the pinned revision)
#   WORK         scratch directory, emptied first (default tmp/olympia-cli-compat)
#   RUN_TIMEOUT  seconds allowed per program run (default 60)
#
# Exits 0 only if every check passes. Results are written to
# $WORK/results.tsv; per-run outputs are kept under $WORK/runs.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"

G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-aa68d137729564cfcd2e91c808c38f3e76d65744}"
WORK="${WORK:-${ROOT}/tmp/olympia-cli-compat}"
RUN_TIMEOUT="${RUN_TIMEOUT:-60}"

die() { echo "error: $*" >&2; exit 2; }
sha() { shasum -a 256 "$1" | cut -d' ' -f1; }

# run_as BIN ARGS... runs BIN with argv[0] set to olympia-g3, killing it
# after RUN_TIMEOUT seconds.
run_as() {
	perl -e 'alarm shift; my $bin = shift; exec { $bin } "olympia-g3", @ARGV or die "exec: $!"' \
		"${RUN_TIMEOUT}" "$@"
}

for tool in git cmake go perl python3 shasum cmp diff; do
	command -v "${tool}" >/dev/null || die "missing required tool: ${tool}"
done
git -C "${G3_REPO}" cat-file -e "${G3_REV}^{commit}" 2>/dev/null ||
	die "${G3_REPO} does not contain ${G3_REV}; clone https://github.com/playbymail/olympia-g3 and set G3_REPO"

if [ -d "${WORK}" ]; then chmod -R u+w "${WORK}"; fi
rm -rf "${WORK}"
mkdir -p "${WORK}/g3" "${WORK}/bin" "${WORK}/runs"
LOG="${WORK}/build.log"

echo "G3 revision: ${G3_REV}"
echo "G3 source:   ${G3_REPO}"
echo "work dir:    ${WORK}"
echo

# ---------------------------------------------------------------------------
# Build the reference from a clean export of the pinned revision, and the port.
# ---------------------------------------------------------------------------
echo "exporting and building G3 (log: ${LOG})"
git -C "${G3_REPO}" archive "${G3_REV}" | tar -x -C "${WORK}/g3"
(
	cd "${WORK}/g3"
	cmake --preset debug
	cmake --build --preset debug --target olympia-g3
) >"${LOG}" 2>&1 || die "G3 build failed; see ${LOG}"
C_BIN="${WORK}/g3/build/debug/olympia-g3"

echo "building the Go port"
(
	cd "${ROOT}"
	go build -o "${WORK}/bin/olympia-g3" ./cmd/olympia
) >>"${LOG}" 2>&1 || die "Go build failed; see ${LOG}"
GO_BIN="${WORK}/bin/olympia-g3"

# ---------------------------------------------------------------------------
# Runs.
# ---------------------------------------------------------------------------

# run_case BIN DIR ARGS... runs one argument vector twice, each time in a
# fresh empty directory: DIR/{stdout,stderr,status} with the streams apart,
# and DIR/{merged,merged-status} with them together.
run_case() {
	local bin="$1" dir="$2" ec
	shift 2
	mkdir -p "${dir}/cwd1" "${dir}/cwd2"
	ec=0
	(cd "${dir}/cwd1" && run_as "${bin}" ${1+"$@"} </dev/null >"${dir}/stdout" 2>"${dir}/stderr") || ec=$?
	echo "${ec}" >"${dir}/status"
	ec=0
	(cd "${dir}/cwd2" && run_as "${bin}" ${1+"$@"} </dev/null >"${dir}/merged" 2>&1) || ec=$?
	echo "${ec}" >"${dir}/merged-status"
	rm -rf "${dir}/cwd1" "${dir}/cwd2"
}

# differs A B FILE reports a difference in FILE between two run directories.
differs() {
	cmp -s "$1/$3" "$2/$3" && return 1
	echo "    $3 differs"
	diff -u "$1/$3" "$2/$3" | head -40 | sed 's/^/      /' || true
	return 0
}

failures=0
fail() { echo "  FAIL: $*"; failures=$((failures + 1)); }

RESULTS="${WORK}/results.tsv"
printf 'case\texpect\tc_status\tgo_status\tstdout_sha256\tstderr_sha256\tmerged_sha256\tresult\targs\n' >"${RESULTS}"

echo
echo "argument vectors"
labels=" "
while read -r label expect args; do
	[[ -z "${label}" || "${label}" == \#* ]] && continue
	case "${labels}" in *" ${label} "*) die "cases.txt: duplicate label ${label}" ;; esac
	labels="${labels}${label} "
	base="${WORK}/runs/${label}"
	eval "set -- ${args}"
	run_case "${C_BIN}" "${base}/c" ${1+"$@"}
	run_case "${GO_BIN}" "${base}/go" ${1+"$@"}

	result=PASS
	c_status="$(cat "${base}/c/status")"
	go_status="$(cat "${base}/go/status")"
	# perl's alarm kills a run that takes longer than RUN_TIMEOUT with
	# SIGALRM, which the shell reports as 142.
	if grep -qx 142 "${base}"/*/status "${base}"/*/merged-status; then
		fail "${label}: a run timed out after ${RUN_TIMEOUT}s"
		result=FAIL
	fi
	case "${expect}" in
	same)
		for f in stdout stderr merged status merged-status; do
			if differs "${base}/c" "${base}/go" "${f}"; then
				fail "${label}: ${f} differs from C"
				result=FAIL
			fi
		done
		;;
	nyi)
		if [ "${go_status}" = 0 ]; then
			fail "${label}: Go exited 0"
			result=FAIL
		fi
		if ! grep -q 'not yet implemented' "${base}/go/stderr"; then
			fail "${label}: Go did not say \"not yet implemented\""
			result=FAIL
		fi
		n="$(wc -c <"${base}/go/stdout" | tr -d ' ')"
		if ! head -c "${n}" "${base}/c/stdout" | cmp -s - "${base}/go/stdout"; then
			fail "${label}: Go stdout is not the start of C's"
			diff -u "${base}/c/stdout" "${base}/go/stdout" | head -20 | sed 's/^/      /' || true
			result=FAIL
		fi
		;;
	*) die "cases.txt: unknown expectation '${expect}' for ${label}" ;;
	esac

	printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' "${label}" "${expect}" "${c_status}" "${go_status}" \
		"$(sha "${base}/c/stdout")" "$(sha "${base}/c/stderr")" "$(sha "${base}/c/merged")" \
		"${result}" "${args}" >>"${RESULTS}"
	printf '  %-4s %-5s %-18s C %-3s Go %-3s %s\n' "${result}" "${expect}" "${label}" \
		"${c_status}" "${go_status}" "${args}"
done <"${HERE}/cases.txt"

# ---------------------------------------------------------------------------
# G3's rng-selftest gate, unchanged. It runs the C engine by default and
# another engine when OLYMPIA_ENGINE_BIN names it (tests/olympia/engine.sh).
# It checks -R on a pipe against EXPECT.txt and on a pseudo-terminal.
# ---------------------------------------------------------------------------
echo
echo "G3 rng-selftest gate"
check="${WORK}/g3/tests/olympia/regress/rng-selftest/check.sh"
for impl in c go; do
	if [ "${impl}" = go ]; then
		out="$(OLYMPIA_ENGINE_BIN="${GO_BIN}" bash "${check}" 2>&1)" || true
	else
		out="$(env -u OLYMPIA_ENGINE_BIN OLYMPIA_PRESET=debug bash "${check}" 2>&1)" || true
	fi
	printf '  %-14s %-3s %s\n' rng-selftest "${impl}" "$(head -1 <<<"${out}")"
	[ "$(head -1 <<<"${out}")" = YES ] || fail "gate rng-selftest (${impl}): ${out}"
done

echo
echo "results: ${RESULTS}"
if [ "${failures}" -ne 0 ]; then
	echo "${failures} check(s) FAILED"
	exit 1
fi
echo "all checks passed"
