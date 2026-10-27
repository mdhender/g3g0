#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: G3's island-g3 (C, the reference) against cmd/island (Go).
#
# Exports the pinned G3 revision, builds its island-g3 and genesis-g3, builds
# the Go port, then runs every case in cases.txt. Each case is a pipeline of
# one or more island runs from one input map and one starting randseed. Each
# run gets its own directory with a fresh copy of the seed, and the stages run
# one after another, each reading the stage before's stdout, which is what a
# pipe does: island reads stdin to EOF before it loads the seed. Each program
# runs the case twice to check determinism. Every stage's stdout, stderr and
# exit status, and the final randseed, must match byte for byte.
#
# Finally it runs G3's own island regression gates, with their assertions and
# expectations unchanged, against both binaries.
#
# Usage: tools/island-compat/run.sh
#
# Environment:
#   G3_REPO      G3 clone that contains G3_REV
#                (default ~/Software/playbymail/olympia-g3)
#   G3_REV       G3 revision to test against (default: the pinned revision)
#   WORK         scratch directory, emptied first (default tmp/island-compat)
#   RUN_TIMEOUT  seconds allowed per program run (default 60)
#
# Exits 0 only if every check passes. Results are written to
# $WORK/results.tsv and $WORK/gates.tsv; per-run outputs are kept under
# $WORK/runs.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"
# shellcheck source=../lib/timeout.sh
. "${HERE}/../lib/timeout.sh"

G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-6bbcd63055d051ade6f745f368281628441a2c27}"
WORK="${WORK:-${ROOT}/tmp/island-compat}"
RUN_TIMEOUT="${RUN_TIMEOUT:-60}"

die() { echo "error: $*" >&2; exit 2; }
sha() { shasum -a 256 "$1" | cut -d' ' -f1; }

# run_as BIN ARGS... runs BIN with argv[0] set to island-g3, killing it after
# RUN_TIMEOUT seconds.
run_as() {
	run_limited "${RUN_TIMEOUT}" island-g3 "$@"
}

for tool in git cmake go perl shasum xxd cmp diff; do
	command -v "${tool}" >/dev/null || die "missing required tool: ${tool}"
done
git -C "${G3_REPO}" cat-file -e "${G3_REV}^{commit}" 2>/dev/null ||
	die "${G3_REPO} does not contain ${G3_REV}; clone https://github.com/playbymail/olympia-g3 and set G3_REPO"

if [ -d "${WORK}" ]; then chmod -R u+w "${WORK}"; fi
rm -rf "${WORK}"
mkdir -p "${WORK}/g3" "${WORK}/bin" "${WORK}/runs" "${WORK}/inputs" "${WORK}/groups"
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
	cmake --build --preset debug --target island-g3 genesis-g3
) >"${LOG}" 2>&1 || die "G3 build failed; see ${LOG}"
C_BIN="${WORK}/g3/build/debug/island-g3"
C_GENESIS="${WORK}/g3/build/debug/genesis-g3"
FIXTURES="${WORK}/g3/tests/island/fixtures"

echo "building the Go port"
(
	cd "${ROOT}"
	go build -o "${WORK}/bin/island-g3" ./cmd/island
) >>"${LOG}" 2>&1 || die "Go build failed; see ${LOG}"
GO_BIN="${WORK}/bin/island-g3"

# ---------------------------------------------------------------------------
# Inputs.
# ---------------------------------------------------------------------------

# make_map SPEC writes the map SPEC names and prints its path.
#   fx:NAME          tests/island/fixtures/NAME from the export
#   file:NAME        maps/NAME here
#   gen:W:H:XX[:crlf] H rows of W copies of the byte with hex value XX
#   genesis:ARGS     the Map that the pinned genesis-g3 writes for ARGS, with
#                    commas in ARGS standing for blanks
make_map() {
	local spec="$1" out
	out="${WORK}/inputs/map-$(printf '%s' "${spec}" | shasum -a 256 | cut -c1-16)"
	if [ ! -f "${out}" ]; then
		case "${spec}" in
		fx:*) cp "${FIXTURES}/${spec#fx:}" "${out}" ;;
		file:*) cp "${HERE}/maps/${spec#file:}" "${out}" ;;
		gen:*)
			perl -e 'my ($w, $h, $c, $eol) = split /:/, shift;
				$eol = ($eol // "") eq "crlf" ? "\r\n" : "\n";
				binmode STDOUT; print((chr(hex $c) x $w) . $eol) for 1 .. $h' "${spec#gen:}" >"${out}"
			;;
		genesis:*)
			local args="${spec#genesis:}" dir="${out}.genesis"
			mkdir -p "${dir}"
			# shellcheck disable=SC2086
			(cd "${dir}" && "${C_GENESIS}" ${args//,/ } >/dev/null 2>&1) ||
				die "genesis-g3 ${args//,/ } failed"
			cp "${dir}/Map" "${out}"
			;;
		*) die "cases.txt: unknown map '${spec}'" ;;
		esac
	fi
	echo "${out}"
}

# put_seed SPEC DIR writes the starting randseed SPEC names into DIR.
#   fx:NAME  tests/island/fixtures/NAME from the export
#   hex:HEX  the bytes HEX; any length, so short and long seed files too
#   none     no randseed
#   dir      randseed is an empty directory
put_seed() {
	case "$1" in
	fx:*) cp "${FIXTURES}/${1#fx:}" "$2/randseed" ;;
	hex:*) printf '%s' "${1#hex:}" | xxd -r -p >"$2/randseed" ;;
	none) ;;
	dir) mkdir "$2/randseed" ;;
	*) die "cases.txt: unknown seed '$1'" ;;
	esac
}

# ---------------------------------------------------------------------------
# Runs.
# ---------------------------------------------------------------------------

# run_case BIN DIR MAP SEED SETUP STAGES... runs one pipeline in a fresh DIR.
# Stage K reads in.K and writes out.K, err.K and status.K; the next stage
# reads out.K. The seed is left in DIR/out/randseed.
run_case() {
	local bin="$1" dir="$2" map="$3" seed="$4" setup="$5"
	shift 5
	mkdir -p "${dir}/out"
	put_seed "${seed}" "${dir}/out"
	case "${setup}" in
	-) ;;
	ro) chmod a-w "${dir}/out/randseed" ;;
	*) die "cases.txt: unknown setup '${setup}'" ;;
	esac
	cp "${map}" "${dir}/in.1"
	local k=1 stage ec
	for stage in "$@"; do
		ec=0
		eval "set -- ${stage}"
		(cd "${dir}/out" && run_as "${bin}" ${1+"$@"} <"${dir}/in.${k}" >"${dir}/out.${k}" 2>"${dir}/err.${k}") ||
			ec=$?
		echo "${ec}" >"${dir}/status.${k}"
		k=$((k + 1))
		cp "${dir}/out.$((k - 1))" "${dir}/in.${k}"
	done
	rm -f "${dir}/in.${k}"
	if [ -e "${dir}/out/randseed" ]; then chmod -R u+w "${dir}/out/randseed"; fi
}

# same A B NAME reports whether two run directories agree on every file they
# contain.
same() {
	local a="$1" b="$2" name="$3" ok=0 f
	if ! diff <(cd "${a}" && find . | sort) <(cd "${b}" && find . | sort) >/dev/null; then
		ok=1
		echo "    ${name}: different files"
		diff <(cd "${a}" && find . | sort) <(cd "${b}" && find . | sort) | sed 's/^/      /' || true
	fi
	while IFS= read -r f; do
		[ -f "${a}/${f}" ] && [ -f "${b}/${f}" ] || continue
		cmp -s "${a}/${f}" "${b}/${f}" && continue
		ok=1
		echo "    ${name}: ${f} differs"
		echo "      $(sha "${a}/${f}")  ${a}/${f}"
		echo "      $(sha "${b}/${f}")  ${b}/${f}"
		cmp "${a}/${f}" "${b}/${f}" 2>&1 | sed 's/^/      /' || true
		case "${f}" in
		*randseed) echo "      $(xxd -p "${a}/${f}") vs $(xxd -p "${b}/${f}")" ;;
		*) diff -u "${a}/${f}" "${b}/${f}" | head -40 | sed 's/^/      /' || true ;;
		esac
	done < <(cd "${a}" && find . -type f | sort)
	return "${ok}"
}

failures=0
fail() { echo "  FAIL: $*"; failures=$((failures + 1)); }

RESULTS="${WORK}/results.tsv"
printf 'case\texpect\tstatuses\tstdout_sha256\trandseed_hex\tstderr_sha256\tresult\tmap\tseed\tstages\n' >"${RESULTS}"

echo
echo "differential runs"
while read -r label expect group mapspec seed setup stages; do
	[[ -z "${label}" || "${label}" == \#* ]] && continue
	map="$(make_map "${mapspec}")"
	# Stages are separated by '|'; an empty stage runs island with no arguments.
	IFS='|' read -r -a stage_list <<<"${stages}|"
	[ "${#stage_list[@]}" -gt 0 ] || stage_list=("")
	base="${WORK}/runs/${label}"
	for run in c1 c2 go1 go2; do
		bin="${C_BIN}"
		case "${run}" in go*) bin="${GO_BIN}" ;; esac
		run_case "${bin}" "${base}/${run}" "${map}" "${seed}" "${setup}" "${stage_list[@]}"
	done

	result=PASS
	statuses="$(cd "${base}/c1" && cat $(ls status.* | sort -t. -k2 -n) | paste -sd, -)"
	last="${statuses##*,}"
	# run_limited kills a run that takes longer than RUN_TIMEOUT and reports
	# TIMEOUT_STATUS (tools/lib/timeout.sh).
	if grep -qx "${TIMEOUT_STATUS}" "${base}"/*/status.*; then
		fail "${label}: a run timed out after ${RUN_TIMEOUT}s"
		result=FAIL
	fi
	case "${expect}" in
	ok)
		[[ "${statuses}" =~ ^0(,0)*$ ]] || { fail "${label}: reference stages exited ${statuses}"; result=FAIL; }
		;;
	fail)
		[ "${last}" != 0 ] || { fail "${label}: reference island-g3 exited 0"; result=FAIL; }
		;;
	*) die "cases.txt: unknown expectation '${expect}' for ${label}" ;;
	esac
	same "${base}/c1" "${base}/c2" "C determinism" || { fail "${label}: C runs disagree"; result=FAIL; }
	same "${base}/go1" "${base}/go2" "Go determinism" || { fail "${label}: Go runs disagree"; result=FAIL; }
	same "${base}/c1" "${base}/go1" "C vs Go" || { fail "${label}: Go output differs from C"; result=FAIL; }

	n="${#stage_list[@]}"
	out_sha="$(sha "${base}/c1/out.${n}")"
	err_sha="$(cat "${base}"/c1/err.* | shasum -a 256 | cut -d' ' -f1)"
	seed_hex=-
	[ -f "${base}/c1/out/randseed" ] && seed_hex="$(xxd -p "${base}/c1/out/randseed" | tr -d '\n')"
	[ -n "${seed_hex}" ] || seed_hex=empty
	printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' "${label}" "${expect}" "${statuses}" \
		"${out_sha}" "${seed_hex}" "${err_sha}" "${result}" "${mapspec}" "${seed}" "${stages}" >>"${RESULTS}"
	printf '  %-4s %-24s status %-9s stdout %.12s  randseed %s\n' "${result}" "${label}" \
		"${statuses}" "${out_sha}" "${seed_hex}"

	if [ "${expect}" = ok ] && [ "${group}" != - ]; then
		cat "${base}"/c1/out.* "${base}/c1/out/randseed" | shasum -a 256 | cut -d' ' -f1 >>"${WORK}/groups/${group}"
	fi
done <"${HERE}/cases.txt"

# Cases in a group differ only in their seed, so they must give different
# maps, or a port that ignored the seed could pass.
echo
echo "seed distinctness"
for file in "${WORK}"/groups/*; do
	[ -e "${file}" ] || continue
	group="$(basename "${file}")"
	total="$(wc -l <"${file}" | tr -d ' ')"
	distinct="$(sort -u "${file}" | wc -l | tr -d ' ')"
	if [ "${distinct}" -ne "${total}" ]; then
		fail "group ${group}: only ${distinct} distinct reference outputs for ${total} cases"
	else
		echo "  group ${group}: ${distinct} cases gave ${distinct} distinct reference outputs"
	fi
done

# ---------------------------------------------------------------------------
# G3's island regression gates, unchanged. They run the C reference by
# default and another island generator when ISLAND_BIN names it
# (tests/island/island.sh).
# ---------------------------------------------------------------------------
echo
echo "G3 island regression gates"
GATES="${WORK}/gates.tsv"
printf 'gate\timpl\tresult\n' >"${GATES}"
for check in "${WORK}"/g3/tests/island/regress/*/check.sh; do
	name="$(basename "$(dirname "${check}")")"
	for impl in c go; do
		if [ "${impl}" = go ]; then
			out="$(ISLAND_BIN="${GO_BIN}" bash "${check}" 2>&1)" || true
		else
			out="$(env -u ISLAND_BIN OLYMPIA_PRESET=debug bash "${check}" 2>&1)" || true
		fi
		printf '  %-20s %-3s %s\n' "${name}" "${impl}" "$(head -1 <<<"${out}")"
		printf '%s\t%s\t%s\n' "${name}" "${impl}" "$(head -1 <<<"${out}")" >>"${GATES}"
		[ "$(head -1 <<<"${out}")" = YES ] || fail "gate ${name} (${impl}): ${out}"
	done
done

echo
echo "results: ${RESULTS}"
if [ "${failures}" -ne 0 ]; then
	echo "${failures} check(s) FAILED"
	exit 1
fi
echo "all checks passed"
