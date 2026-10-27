#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: G3's mapgen-g3 (C, the reference) against cmd/mapgen (Go).
#
# Exports the pinned G3 revision, builds its mapgen-g3 and genesis-g3, builds
# the Go port, then runs both on every input set in inputs.txt with every seed
# in seeds.txt. Each run gets its own disposable directory with byte-identical
# copies of the inputs. Each program runs twice per case to check determinism.
# loc, gate, road, the final randseed, stdout, stderr and the exit status must
# match byte for byte. Finally it runs G3's own mapgen regression checks, with
# their assertions and expectations unchanged, against both binaries.
#
# Usage: tools/mapgen-compat/run.sh [--update]
# --update records failure statuses from deterministic C runs and results.tsv.
#
# Environment:
#   G3_REPO      G3 clone that contains G3_REV
#                (default ~/Software/playbymail/olympia-g3)
#   G3_REV       G3 revision to test against (default: the pinned revision)
#   WORK         scratch directory, emptied first (default tmp/mapgen-compat)
#   RUN_TIMEOUT  seconds allowed per program run (default 120)
#
# Exits 0 only if every check passes. Results are written to
# $WORK/results.tsv; per-run outputs are kept under $WORK/runs.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"
# shellcheck source=../lib/timeout.sh
. "${HERE}/../lib/timeout.sh"

G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-5a114320c05e37c6c521b4eb730af49017cb6fa8}"
WORK="${WORK:-${ROOT}/tmp/mapgen-compat}"
RUN_TIMEOUT="${RUN_TIMEOUT:-120}"

OUTPUTS=(loc gate road randseed stdout stderr status)
INPUTS=(Map Regions Land Cities)

die() { echo "error: $*" >&2; exit 2; }
sha() { shasum -a 256 "$1" | cut -d' ' -f1; }

UPDATE=0
if [ "${1:-}" = --update ]; then UPDATE=1; shift; fi
[ $# -eq 0 ] || die "usage: $0 [--update]"

# with_timeout BIN ARGS... runs BIN, killing it after RUN_TIMEOUT seconds
# (tools/lib/timeout.sh). A run that times out exits TIMEOUT_STATUS, for C
# and Go alike.
with_timeout() { run_limited "${RUN_TIMEOUT}" "$1" "$@"; }

# run_capture BIN runs BIN under the same timeout, writing stdout and stderr in
# the current directory.
run_capture() { with_timeout "$1" >stdout 2>stderr; }

for tool in git cmake go perl shasum xxd cmp diff; do
	command -v "${tool}" >/dev/null || die "missing required tool: ${tool}"
done
git -C "${G3_REPO}" cat-file -e "${G3_REV}^{commit}" 2>/dev/null ||
	die "${G3_REPO} does not contain ${G3_REV}; clone https://github.com/playbymail/olympia-g3 and set G3_REPO"

rm -rf "${WORK}"
mkdir -p "${WORK}/g3" "${WORK}/bin" "${WORK}/inputs" "${WORK}/seeds" "${WORK}/runs"
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
	cmake --build --preset debug --target mapgen-g3 genesis-g3
) >"${LOG}" 2>&1 || die "G3 build failed; see ${LOG}"
C_BIN="${WORK}/g3/build/debug/mapgen-g3"
GENESIS="${WORK}/g3/build/debug/genesis-g3"

echo "building the Go port"
(cd "${ROOT}" && go build -o "${WORK}/bin/mapgen-g3" ./cmd/mapgen) >>"${LOG}" 2>&1 ||
	die "Go build failed; see ${LOG}"
GO_BIN="${WORK}/bin/mapgen-g3"

# ---------------------------------------------------------------------------
# Input sets and seeds.
# ---------------------------------------------------------------------------
input_labels=()
while read -r label source args; do
	[[ -z "${label}" || "${label}" == \#* ]] && continue
	dir="${WORK}/inputs/${label}"
	mkdir -p "${dir}"
	case "${source}" in
	fixtures)
		for f in "${INPUTS[@]}"; do cp -p "${WORK}/g3/tests/mapgen/fixtures/${f}" "${dir}/"; done
		;;
	regress/*)
		# A G3 regression's own inputs; files it lacks come from the fixtures.
		for f in "${INPUTS[@]}"; do
			src="${WORK}/g3/tests/mapgen/${source}/${f}"
			[ -f "${src}" ] || src="${WORK}/g3/tests/mapgen/fixtures/${f}"
			cp -p "${src}" "${dir}/"
		done
		;;
	genesis)
		# shellcheck disable=SC2086 # args is a word list by design
		(cd "${dir}" && with_timeout "${GENESIS}" ${args} >genesis.out 2>&1) ||
			die "genesis-g3 ${args} failed; see ${dir}/genesis.out"
		rm -f "${dir}/randseed" "${dir}/genesis.out"
		;;
	*) die "inputs.txt: unknown source '${source}' for ${label}" ;;
	esac
	for f in "${INPUTS[@]}"; do [ -f "${dir}/${f}" ] || die "input set ${label} has no ${f}"; done
	input_labels+=("${label}")
done <"${HERE}/inputs.txt"

seed_labels=()
while read -r label hex; do
	[[ -z "${label}" || "${label}" == \#* ]] && continue
	if [ "${hex}" != "missing" ]; then
		printf '%s' "${hex}" | xxd -r -p >"${WORK}/seeds/${label}"
	fi
	seed_labels+=("${label}")
done <"${HERE}/seeds.txt"

# ---------------------------------------------------------------------------
# Runs.
# ---------------------------------------------------------------------------

# run_one BIN INPUT SEED DIR [SECONDS] runs BIN in a fresh DIR holding copies
# of the input set and the original seed, with a limit of SECONDS (default
# RUN_TIMEOUT).
run_one() {
	local bin="$1" input="$2" seed="$3" dir="$4" RUN_TIMEOUT="${5:-${RUN_TIMEOUT}}"
	mkdir -p "${dir}"
	for f in "${INPUTS[@]}"; do cp -p "${WORK}/inputs/${input}/${f}" "${dir}/"; done
	if [ -f "${WORK}/seeds/${seed}" ]; then cp -p "${WORK}/seeds/${seed}" "${dir}/randseed"; fi
	local ec=0
	(cd "${dir}" && run_capture "${bin}") 2>/dev/null || ec=$?
	echo "${ec}" >"${dir}/status"
	# A run that writes no randseed must still compare against one that does.
	[ -e "${dir}/randseed" ] || : >"${dir}/randseed.absent"
}

# same A B NAME reports whether two run directories agree on every output.
same() {
	local a="$1" b="$2" name="$3" ok=0 f
	for f in "${OUTPUTS[@]}" randseed.absent; do
		if [ ! -e "${a}/${f}" ] && [ ! -e "${b}/${f}" ]; then continue; fi
		if [ -f "${a}/${f}" ] && [ -f "${b}/${f}" ] && cmp -s "${a}/${f}" "${b}/${f}"; then continue; fi
		ok=1
		echo "    ${name}: ${f} differs"
		if [ -f "${a}/${f}" ] && [ -f "${b}/${f}" ]; then
			echo "      $(sha "${a}/${f}")  ${a}/${f}"
			echo "      $(sha "${b}/${f}")  ${b}/${f}"
			cmp "${a}/${f}" "${b}/${f}" 2>&1 | sed 's/^/      /' || true
			if [ "${f}" != randseed ]; then
				diff -u "${a}/${f}" "${b}/${f}" | head -40 | sed 's/^/      /' || true
			else
				echo "      $(xxd -p "${a}/${f}") vs $(xxd -p "${b}/${f}")"
			fi
		else
			echo "      present only in one of ${a} and ${b}"
		fi
	done
	return "${ok}"
}

failures=0
fail() { echo "  FAIL: $*"; failures=$((failures + 1)); }

RESULTS="${WORK}/results.tsv"
printf 'input\tseed\tseed_hex\tstatus\tloc_sha256\tgate_sha256\troad_sha256\trandseed_hex\tresult\n' >"${RESULTS}"

echo
echo "differential runs (${#input_labels[@]} input sets x ${#seed_labels[@]} seeds)"
for input in "${input_labels[@]}"; do
	echo
	echo "input set: ${input}"
	for seed in "${seed_labels[@]}"; do
		base="${WORK}/runs/${input}/${seed}"
		run_one "${C_BIN}" "${input}" "${seed}" "${base}/c1"
		run_one "${C_BIN}" "${input}" "${seed}" "${base}/c2"
		run_one "${GO_BIN}" "${input}" "${seed}" "${base}/go1"
		run_one "${GO_BIN}" "${input}" "${seed}" "${base}/go2"

		result=PASS
		status="$(cat "${base}/c1/status")"
		if [ "${status}" != 0 ]; then
			fail "${input}/${seed}: reference mapgen-g3 exited ${status}"
			result=FAIL
		fi
		same "${base}/c1" "${base}/c2" "C determinism" || { fail "${input}/${seed}: C runs disagree"; result=FAIL; }
		same "${base}/go1" "${base}/go2" "Go determinism" || { fail "${input}/${seed}: Go runs disagree"; result=FAIL; }
		same "${base}/c1" "${base}/go1" "C vs Go" || { fail "${input}/${seed}: Go output differs from C"; result=FAIL; }

		seed_hex=missing
		[ -f "${WORK}/seeds/${seed}" ] && seed_hex="$(xxd -p "${WORK}/seeds/${seed}")"
		final_hex=none
		[ -f "${base}/c1/randseed" ] && final_hex="$(xxd -p "${base}/c1/randseed")"
		printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' "${input}" "${seed}" "${seed_hex}" "${status}" \
			"$(sha "${base}/c1/loc")" "$(sha "${base}/c1/gate")" "$(sha "${base}/c1/road")" \
			"${final_hex}" "${result}" >>"${RESULTS}"
		printf '  %-6s %-18s loc %.12s  gate %.12s  road %.12s  randseed %s\n' "${result}" "${seed}" \
			"$(sha "${base}/c1/loc")" "$(sha "${base}/c1/gate")" "$(sha "${base}/c1/road")" "${final_hex}"
	done

	# Different full seeds must give different maps, or a port that ignored
	# the seed could pass. Seeds shorter than 16 bytes or missing are edge
	# cases of load_seed(), not distinct seed states, so they are excluded.
	hashes=()
	for seed in "${seed_labels[@]}"; do
		f="${WORK}/seeds/${seed}"
		[ -f "${f}" ] && [ "$(wc -c <"${f}" | tr -d ' ')" = 16 ] || continue
		d="${WORK}/runs/${input}/${seed}/c1"
		hashes+=("$(cat "${d}/loc" "${d}/gate" "${d}/road" | shasum -a 256 | cut -d' ' -f1)")
	done
	distinct="$(printf '%s\n' "${hashes[@]}" | sort -u | wc -l | tr -d ' ')"
	if [ "${distinct}" -ne "${#hashes[@]}" ]; then
		fail "${input}: only ${distinct} distinct reference outputs for ${#hashes[@]} seeds"
	else
		echo "  ${distinct} seeds gave ${distinct} distinct reference maps"
	fi
done

# ---------------------------------------------------------------------------
# Failure cases: inputs on which the C program fails. The reference must exit
# with the recorded status, and Go must match it on every output. A case may
# set its own time limit; a run that hits it has status TIMEOUT_STATUS. asfail()
# prints the compiler's __FILE__, the absolute path of the exported source,
# where the port prints mapgen/mapgen.c, so that prefix is removed from the C
# stderr before comparing.
# ---------------------------------------------------------------------------
G3_ABS="${WORK}/g3/"
G3_PHYS="$(cd "${WORK}/g3" && pwd -P)/"

# strip_file_prefix DIR removes the export's path from DIR's stderr.
strip_file_prefix() {
	G3_ABS="${G3_ABS}" G3_PHYS="${G3_PHYS}" perl -pi -e 's/\Q$ENV{G3_PHYS}\E//g; s/\Q$ENV{G3_ABS}\E//g' "$1/stderr"
}

echo
echo "failure cases"
: >"${WORK}/failures.txt"
while IFS= read -r row; do
	read -r label seed want limit <<<"${row}"
	if [[ -z "${label}" || "${label}" == \#* ]]; then
		printf '%s\n' "${row}" >>"${WORK}/failures.txt"
		continue
	fi
	src="${HERE}/failures/${label}"
	[ -d "${src}" ] || die "failures.txt: no directory failures/${label}"
	[ -n "${want}" ] || die "failures.txt: no status for ${label}"
	[ "${want}" != timeout ] || want="${TIMEOUT_STATUS}"
	recorded_limit="${limit}"
	limit="${limit:-${RUN_TIMEOUT}}"
	printf '%s\n' "${seed_labels[@]}" | grep -qx -- "${seed}" || die "failures.txt: unknown seed '${seed}' for ${label}"
	input="fail-${label}"
	if [ ! -d "${WORK}/inputs/${input}" ]; then
		mkdir -p "${WORK}/inputs/${input}"
		for f in "${INPUTS[@]}"; do
			from="${src}/${f}"
			[ -f "${from}" ] || from="${WORK}/g3/tests/mapgen/fixtures/${f}"
			cp -p "${from}" "${WORK}/inputs/${input}/"
		done
	fi
	base="${WORK}/runs/${input}/${seed}"
	run_one "${C_BIN}" "${input}" "${seed}" "${base}/c1" "${limit}"
	run_one "${C_BIN}" "${input}" "${seed}" "${base}/c2" "${limit}"
	run_one "${GO_BIN}" "${input}" "${seed}" "${base}/go1" "${limit}"
	run_one "${GO_BIN}" "${input}" "${seed}" "${base}/go2" "${limit}"
	strip_file_prefix "${base}/c1"
	strip_file_prefix "${base}/c2"

	result=PASS
	status="$(cat "${base}/c1/status")"
	if [ "${UPDATE}" -eq 1 ]; then want="${status}"; fi
	printf '%-18s %-8s %s%s\n' "${label}" "${seed}" "${want}" "${recorded_limit:+  ${recorded_limit}}" >>"${WORK}/failures.txt"
	if [ "${status}" != "${want}" ]; then
		fail "${input}/${seed}: reference mapgen-g3 exited ${status}, not ${want}"
		result=FAIL
	fi
	same "${base}/c1" "${base}/c2" "C determinism" || { fail "${input}/${seed}: C runs disagree"; result=FAIL; }
	same "${base}/go1" "${base}/go2" "Go determinism" || { fail "${input}/${seed}: Go runs disagree"; result=FAIL; }
	same "${base}/c1" "${base}/go1" "C vs Go" || { fail "${input}/${seed}: Go output differs from C"; result=FAIL; }

	seed_hex=missing
	[ -f "${WORK}/seeds/${seed}" ] && seed_hex="$(xxd -p "${WORK}/seeds/${seed}")"
	final_hex=none
	[ -f "${base}/c1/randseed" ] && final_hex="$(xxd -p "${base}/c1/randseed")"
	printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' "${input}" "${seed}" "${seed_hex}" "${status}" \
		"$(sha "${base}/c1/loc")" "$(sha "${base}/c1/gate")" "$(sha "${base}/c1/road")" \
		"${final_hex}" "${result}" >>"${RESULTS}"
	printf '  %-6s %-18s %-8s status %s  stdout %s bytes  %s\n' "${result}" "${label}" "${seed}" "${status}" \
		"$(wc -c <"${base}/c1/stdout" | tr -d ' ')" "$(tail -1 "${base}/c1/stderr")"
done <"${HERE}/failures.txt"
if [ "${UPDATE}" -eq 1 ] && [ "${failures}" -eq 0 ]; then
	cp "${WORK}/failures.txt" "${HERE}/failures.txt"
fi

# ---------------------------------------------------------------------------
# G3's mapgen regression checks, unchanged. They run the C reference by
# default and another map generator when MAPGEN_BIN names it
# (tests/mapgen/mapgen.sh).
# ---------------------------------------------------------------------------
echo
echo "G3 mapgen regression checks"
for check in "${WORK}"/g3/tests/mapgen/regress/*/check.sh; do
	name="$(basename "$(dirname "${check}")")"
	for impl in c go; do
		bin="${C_BIN}"
		[ "${impl}" != go ] || bin="${GO_BIN}"
		cat >"${WORK}/bin/fixture-engine" <<EOF
#!/usr/bin/env bash
exec bash "${HERE}/../lib/run-limited.sh" "${bin}" "\$@"
EOF
		chmod +x "${WORK}/bin/fixture-engine"
		out="$(MAPGEN_BIN="${WORK}/bin/fixture-engine" bash "${check}" 2>&1)" || true
		printf '  %-22s %-3s %s\n' "${name}" "${impl}" "$(head -1 <<<"${out}")"
		[ "$(head -1 <<<"${out}")" = YES ] || fail "regression ${name} (${impl}): ${out}"
	done
done

echo
echo "results: ${RESULTS}"
if [ "${failures}" -ne 0 ]; then
	echo "${failures} check(s) FAILED"
	exit 1
fi
[ "${UPDATE}" -ne 1 ] || cp "${RESULTS}" "${HERE}/results.tsv"
echo "all checks passed"
