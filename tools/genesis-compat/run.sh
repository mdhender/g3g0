#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: G3's genesis-g3 (C, the reference) against cmd/genesis (Go).
#
# Exports the pinned G3 revision, builds its genesis-g3 and mapgen-g3, builds
# the Go ports of both, then runs every case in cases.txt. Each run gets its
# own empty directory, and both programs see the same argv, including argv[0].
# Each program runs twice per case to check determinism. Every file the run
# leaves behind (Map, Regions, Land, Cities, randseed, any BMP), stderr,
# stdout and the exit status must match byte for byte.
#
# Then each world the Go port generated is fed to the C mapgen-g3 and to the
# Go mapgen, which must both accept it and agree. Finally it runs G3's own
# genesis regression gates, with their assertions and expectations unchanged,
# against both binaries.
#
# Usage: tools/genesis-compat/run.sh
#
# Environment:
#   G3_REPO      G3 clone that contains G3_REV
#                (default ~/Software/playbymail/olympia-g3)
#   G3_REV       G3 revision to test against (default: the pinned revision)
#   WORK         scratch directory, emptied first (default tmp/genesis-compat)
#   RUN_TIMEOUT  seconds allowed per program run (default 120)
#
# Exits 0 only if every check passes. Results are written to
# $WORK/results.tsv and $WORK/gates.tsv; per-run outputs are kept under
# $WORK/runs.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"

G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-2322a0474de5fb38308d795f6cd72beb811652e9}"
WORK="${WORK:-${ROOT}/tmp/genesis-compat}"
RUN_TIMEOUT="${RUN_TIMEOUT:-120}"

WORLD=(Map Regions Land Cities randseed)

die() { echo "error: $*" >&2; exit 2; }
sha() { shasum -a 256 "$1" | cut -d' ' -f1; }

# run_as BIN ARGS... runs BIN with argv[0] set to genesis-g3, killing it after
# RUN_TIMEOUT seconds. genesis prints argv[0] in its usage message.
run_as() {
	perl -e 'alarm shift; my $bin = shift; exec { $bin } "genesis-g3", @ARGV or die "exec: $!"' \
		"${RUN_TIMEOUT}" "$@"
}
with_timeout() { perl -e 'alarm shift; exec @ARGV or die "exec: $!"' "${RUN_TIMEOUT}" "$@"; }

for tool in git cmake go perl shasum xxd cmp diff; do
	command -v "${tool}" >/dev/null || die "missing required tool: ${tool}"
done
git -C "${G3_REPO}" cat-file -e "${G3_REV}^{commit}" 2>/dev/null ||
	die "${G3_REPO} does not contain ${G3_REV}; clone https://github.com/playbymail/olympia-g3 and set G3_REPO"

if [ -d "${WORK}" ]; then chmod -R u+w "${WORK}"; fi
rm -rf "${WORK}"
mkdir -p "${WORK}/g3" "${WORK}/bin" "${WORK}/runs" "${WORK}/mapgen"
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
	cmake --build --preset debug --target genesis-g3 mapgen-g3
) >"${LOG}" 2>&1 || die "G3 build failed; see ${LOG}"
C_BIN="${WORK}/g3/build/debug/genesis-g3"
C_MAPGEN="${WORK}/g3/build/debug/mapgen-g3"

echo "building the Go ports"
(
	cd "${ROOT}"
	go build -o "${WORK}/bin/genesis-g3" ./cmd/genesis
	go build -o "${WORK}/bin/mapgen-g3" ./cmd/mapgen
) >>"${LOG}" 2>&1 || die "Go build failed; see ${LOG}"
GO_BIN="${WORK}/bin/genesis-g3"
GO_MAPGEN="${WORK}/bin/mapgen-g3"

# ---------------------------------------------------------------------------
# Runs.
# ---------------------------------------------------------------------------

# run_one BIN SETUP DIR ARGS... runs BIN in a fresh DIR.
run_one() {
	local bin="$1" setup="$2" dir="$3"
	shift 3
	mkdir -p "${dir}/out"
	local posix=""
	case "${setup}" in
	-) ;;
	dir:*) mkdir "${dir}/out/${setup#dir:}" ;;
	ro) chmod a-w "${dir}/out" ;;
	posix) posix=1 ;;
	*) die "cases.txt: unknown setup '${setup}'" ;;
	esac
	local ec=0
	(
		cd "${dir}/out"
		if [ -n "${posix}" ]; then export POSIXLY_CORRECT=1; fi
		run_as "${bin}" ${1+"$@"} >"${dir}/stdout" 2>"${dir}/stderr"
	) || ec=$?
	echo "${ec}" >"${dir}/status"
	chmod u+w "${dir}/out"
}

# same A B NAME reports whether two run directories agree on every file they
# contain, and on stdout, stderr and the exit status.
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
		*randseed | *.bmp) echo "      $(xxd -p -l 32 "${a}/${f}") vs $(xxd -p -l 32 "${b}/${f}")" ;;
		*) diff -u "${a}/${f}" "${b}/${f}" | head -40 | sed 's/^/      /' || true ;;
		esac
	done < <(cd "${a}" && find . -type f | sort)
	return "${ok}"
}

# world_hash DIR hashes the world files a run wrote.
world_hash() {
	local f
	for f in "${WORLD[@]}"; do cat "$1/out/${f}"; done | shasum -a 256 | cut -d' ' -f1
}

failures=0
fail() { echo "  FAIL: $*"; failures=$((failures + 1)); }

RESULTS="${WORK}/results.tsv"
printf 'case\texpect\tstatus\tMap_sha256\tRegions_sha256\tLand_sha256\tCities_sha256\trandseed_hex\tstderr_sha256\tresult\targs\n' >"${RESULTS}"

mkdir -p "${WORK}/groups"
ok_labels=()

echo
echo "differential runs"
while read -r label expect group setup args; do
	[[ -z "${label}" || "${label}" == \#* ]] && continue
	eval "set -- ${args}"
	base="${WORK}/runs/${label}"
	run_one "${C_BIN}" "${setup}" "${base}/c1" ${1+"$@"}
	run_one "${C_BIN}" "${setup}" "${base}/c2" ${1+"$@"}
	run_one "${GO_BIN}" "${setup}" "${base}/go1" ${1+"$@"}
	run_one "${GO_BIN}" "${setup}" "${base}/go2" ${1+"$@"}

	result=PASS
	status="$(cat "${base}/c1/status")"
	case "${expect}" in
	ok)
		[ "${status}" = 0 ] || { fail "${label}: reference genesis-g3 exited ${status}"; result=FAIL; }
		for f in "${WORLD[@]}"; do
			[ -f "${base}/c1/out/${f}" ] || { fail "${label}: reference wrote no ${f}"; result=FAIL; }
		done
		;;
	help)
		[ "${status}" = 0 ] || { fail "${label}: reference genesis-g3 exited ${status}"; result=FAIL; }
		[ -z "$(ls -A "${base}/c1/out")" ] || { fail "${label}: reference wrote files"; result=FAIL; }
		;;
	fail)
		[ "${status}" != 0 ] || { fail "${label}: reference genesis-g3 exited 0"; result=FAIL; }
		;;
	*) die "cases.txt: unknown expectation '${expect}' for ${label}" ;;
	esac
	same "${base}/c1" "${base}/c2" "C determinism" || { fail "${label}: C runs disagree"; result=FAIL; }
	same "${base}/go1" "${base}/go2" "Go determinism" || { fail "${label}: Go runs disagree"; result=FAIL; }
	same "${base}/c1" "${base}/go1" "C vs Go" || { fail "${label}: Go output differs from C"; result=FAIL; }

	hashes=()
	for f in Map Regions Land Cities; do
		if [ -f "${base}/c1/out/${f}" ]; then hashes+=("$(sha "${base}/c1/out/${f}")"); else hashes+=(-); fi
	done
	seed_hex=-
	[ -f "${base}/c1/out/randseed" ] && seed_hex="$(xxd -p "${base}/c1/out/randseed")"
	printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' "${label}" "${expect}" "${status}" \
		"${hashes[@]}" "${seed_hex}" "$(sha "${base}/c1/stderr")" "${result}" "${args}" >>"${RESULTS}"
	printf '  %-4s %-20s status %-3s Map %.12s  Regions %.12s  randseed %s\n' "${result}" "${label}" \
		"${status}" "${hashes[0]}" "${hashes[1]}" "${seed_hex}"

	if [ "${expect}" = ok ]; then
		ok_labels+=("${label}")
		if [ "${group}" != - ]; then
			world_hash "${base}/c1" >>"${WORK}/groups/${group}"
		fi
	fi
done <"${HERE}/cases.txt"

# Cases in a group differ only in their seed, so they must give different
# worlds, or a port that ignored the seed could pass.
echo
echo "seed distinctness"
for file in "${WORK}"/groups/*; do
	group="$(basename "${file}")"
	total="$(wc -l <"${file}" | tr -d ' ')"
	distinct="$(sort -u "${file}" | wc -l | tr -d ' ')"
	if [ "${distinct}" -ne "${total}" ]; then
		fail "group ${group}: only ${distinct} distinct reference worlds for ${total} cases"
	else
		echo "  group ${group}: ${distinct} cases gave ${distinct} distinct reference worlds"
	fi
done

# ---------------------------------------------------------------------------
# Feed each Go-generated world to mapgen: the C reference must accept it,
# and the Go mapgen must agree with it byte for byte.
# ---------------------------------------------------------------------------
echo
echo "mapgen acceptance (${#ok_labels[@]} Go-generated worlds)"
MAPGEN_OUT=(loc gate road randseed stdout stderr status)
for label in "${ok_labels[@]}"; do
	for impl in c go; do
		dir="${WORK}/mapgen/${label}/${impl}"
		mkdir -p "${dir}"
		for f in "${WORLD[@]}"; do cp -p "${WORK}/runs/${label}/go1/out/${f}" "${dir}/"; done
		bin="${C_MAPGEN}"
		[ "${impl}" = go ] && bin="${GO_MAPGEN}"
		ec=0
		(cd "${dir}" && with_timeout "${bin}" >stdout 2>stderr) || ec=$?
		echo "${ec}" >"${dir}/status"
	done
	c="${WORK}/mapgen/${label}/c"
	g="${WORK}/mapgen/${label}/go"
	result=PASS
	[ "$(cat "${c}/status")" = 0 ] || { fail "${label}: C mapgen-g3 exited $(cat "${c}/status")"; result=FAIL; }
	for f in "${MAPGEN_OUT[@]}"; do
		if ! cmp -s "${c}/${f}" "${g}/${f}"; then
			fail "${label}: Go mapgen ${f} differs from C"
			result=FAIL
		fi
	done
	printf '  %-4s %-20s loc %.12s  gate %.12s  road %.12s\n' "${result}" "${label}" \
		"$(sha "${c}/loc")" "$(sha "${c}/gate")" "$(sha "${c}/road")"
done

# ---------------------------------------------------------------------------
# G3's genesis regression gates, unchanged. They run the C reference by
# default and another world generator when GENESIS_BIN names it
# (tests/genesis/genesis.sh).
# ---------------------------------------------------------------------------
echo
echo "G3 genesis regression gates"
GATES="${WORK}/gates.tsv"
printf 'gate\timpl\tresult\n' >"${GATES}"
for check in "${WORK}"/g3/tests/genesis/regress/*/check.sh; do
	name="$(basename "$(dirname "${check}")")"
	for impl in c go; do
		if [ "${impl}" = go ]; then
			out="$(GENESIS_BIN="${GO_BIN}" bash "${check}" 2>&1)" || true
		else
			out="$(env -u GENESIS_BIN OLYMPIA_PRESET=debug bash "${check}" 2>&1)" || true
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
