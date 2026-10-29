#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Runs G3's olympia spec, unchanged, against the Go engine (cmd/olympia).
#
# Exports the pinned G3 revision, builds the Go port, and runs, from the
# export, with OLYMPIA_ENGINE_BIN set to the Go engine:
#
#   tests/olympia/regress/run-all.sh    every regress fixture's check.sh
#   run/olympia-g3.sh                   the full golden turn, in RUN_DIR
#   tests/olympia/golden_check.sh       that turn against the golden manifest
#
# The fixtures that pass, plus "golden" when golden_check.sh prints YES, are
# written to passing.txt next to this script, and any change from the list
# recorded there is reported. Expected output comes only from G3's committed
# files: nothing here re-baselines, and --update is refused.
#
# Usage: tools/olympia-spec/run.sh [fixture...]
#
#   With fixture names, runs only those fixtures (and not the golden turn)
#   and leaves passing.txt alone.
#
# Environment:
#   G3_REPO   G3 clone that contains G3_REV
#             (default ~/Software/playbymail/olympia-g3)
#   G3_REV    G3 revision to test against (default: the pinned revision).
#             passing.txt is written only for the pinned revision.
#   WORK      scratch directory, emptied first (default tmp/olympia-spec)
#   FORCE_RECORD=1  write passing.txt even when it loses an entry
#   RUN_DIR   OLYMPIA_RUN_DIR for the golden turn, emptied first; must be
#             outside this repository and the G3 clone
#             (default ${TMPDIR:-/tmp}/g3g0-olympia-spec-run)
#
# Exits 0 unless a fixture or gate recorded as passing no longer passes
# (exit 1; passing.txt is then left unchanged), or the harness itself fails
# (exit 2). Logs are kept under $WORK.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"

PINNED_REV=9a03d126b329ed466c4af48425553c025eaca26f
G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-${PINNED_REV}}"
WORK="${WORK:-${ROOT}/tmp/olympia-spec}"
TMP_BASE="${TMPDIR:-/tmp}"
RUN_DIR="${RUN_DIR:-${TMP_BASE%/}/g3g0-olympia-spec-run}"
RECORD="${HERE}/passing.txt"

die() { echo "error: $*" >&2; exit 2; }

for a in "$@"; do
	case "${a}" in
	--update) die "this harness never re-baselines; expected output only comes from the C engine" ;;
	-*) die "unknown option ${a} (usage: $0 [fixture...])" ;;
	esac
done

for tool in git go perl python3 shasum patch diff; do
	command -v "${tool}" >/dev/null || die "missing required tool: ${tool}"
done
[ -d "${G3_REPO}" ] || die "no G3 clone at ${G3_REPO}; clone https://github.com/playbymail/olympia-g3 and set G3_REPO"
G3_REPO="$(cd "${G3_REPO}" && pwd -P)"
git -C "${G3_REPO}" cat-file -e "${G3_REV}^{commit}" 2>/dev/null ||
	die "${G3_REPO} does not contain ${G3_REV}; fetch it or set G3_REPO"
G3_REV="$(git -C "${G3_REPO}" rev-parse "${G3_REV}^{commit}")"

# RUN_DIR must be outside both checkouts; the export under WORK is inside
# this repository, so that rule covers it too.
mkdir -p "${RUN_DIR}"
RUN_DIR="$(cd "${RUN_DIR}" && pwd -P)"
ROOT_REAL="$(cd "${ROOT}" && pwd -P)"
for d in "${ROOT_REAL}" "${G3_REPO}"; do
	case "${RUN_DIR}/" in "${d}/"*) die "RUN_DIR=${RUN_DIR} is inside ${d}" ;; esac
done

if [ -d "${WORK}" ]; then chmod -R u+w "${WORK}"; fi
rm -rf "${WORK}" "${RUN_DIR}"
mkdir -p "${WORK}/g3" "${WORK}/bin" "${RUN_DIR}"
LOG="${WORK}/build.log"

echo "G3 revision: ${G3_REV}"
[ "${G3_REV}" = "${PINNED_REV}" ] || echo "             (not the pinned ${PINNED_REV}; ${RECORD} is left alone)"
echo "G3 source:   ${G3_REPO}"
echo "work dir:    ${WORK}"
echo "run dir:     ${RUN_DIR}"
echo

# The spec runs from a clean export, so neither the clone's working tree nor
# its build is used, and nothing the gates write can reach the clone.
echo "exporting G3"
git -C "${G3_REPO}" archive "${G3_REV}" | tar -x -C "${WORK}/g3"
G3="${WORK}/g3"

echo "building the Go port (log: ${LOG})"
(
	cd "${ROOT}"
	go build -o "${WORK}/bin/olympia-g3" ./cmd/olympia
) >"${LOG}" 2>&1 || die "Go build failed; see ${LOG}"

# A script loses exec -a's argv[0] at the interpreter boundary (#99).
# This native launcher passes it unchanged to the existing run_limited.
cat >"${WORK}/engine.go" <<'EOF'
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	bash, err := exec.LookPath("bash")
	if err == nil {
		args := []string{"bash", "-c", `. "$1"; shift; run_limited "$@"`,
			"fixture-engine", os.Getenv("G3G0_TIMEOUT_SH"), os.Getenv("RUN_TIMEOUT"),
			os.Args[0], os.Getenv("G3G0_ENGINE_BIN")}
		err = syscall.Exec(bash, append(args, os.Args[1:]...), os.Environ())
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(127)
}
EOF
go build -o "${WORK}/bin/engine" "${WORK}/engine.go" >>"${LOG}" 2>&1 ||
	die "fixture launcher build failed; see ${LOG}"
export G3G0_TIMEOUT_SH="${HERE}/../lib/timeout.sh"
export G3G0_ENGINE_BIN="${WORK}/bin/olympia-g3"
export RUN_TIMEOUT="${RUN_TIMEOUT:-120}"
export OLYMPIA_ENGINE_BIN="${WORK}/bin/engine"
export OLYMPIA_RUN_DIR="${RUN_DIR}"

# ---------------------------------------------------------------------------
# The regress fixtures.
# ---------------------------------------------------------------------------
echo
echo "tests/olympia/regress/run-all.sh (log: ${WORK}/run-all.log)"
start=${SECONDS}
"${G3}/tests/olympia/regress/run-all.sh" "$@" >"${WORK}/run-all.log" 2>&1 || true
grep -E '^(PASS |FAIL |ERROR )' "${WORK}/run-all.log" | sed 's/^/  /' || true
summary="$(grep -E '^pass [0-9]+ ' "${WORK}/run-all.log" || true)"
[ -n "${summary}" ] || die "run-all.sh did not finish; see ${WORK}/run-all.log"
echo "  ${summary} ($((SECONDS - start))s)"

sed -n 's/^PASS  *//p' "${WORK}/run-all.log" >"${WORK}/passing.tmp"

# ---------------------------------------------------------------------------
# The golden turn, only on a full run.
# ---------------------------------------------------------------------------
if [ $# -eq 0 ]; then
	echo
	echo "run/olympia-g3.sh (log: ${WORK}/olympia-g3.log)"
	ec=0
	"${G3}/run/olympia-g3.sh" >"${WORK}/olympia-g3.log" 2>&1 || ec=$?
	echo "  exit ${ec}"
	echo "tests/olympia/golden_check.sh (log: ${WORK}/golden.log)"
	"${G3}/tests/olympia/golden_check.sh" >"${WORK}/golden.log" 2>&1 || true
	echo "  $(head -1 "${WORK}/golden.log")"
	if [ "${ec}" -eq 0 ] && [ "$(head -1 "${WORK}/golden.log")" = YES ]; then
		echo golden >>"${WORK}/passing.tmp"
	fi
fi

LC_ALL=C sort -u "${WORK}/passing.tmp" >"${WORK}/passing.txt"
rm -f "${WORK}/passing.tmp"

# ---------------------------------------------------------------------------
# Compare with the recorded list.
# ---------------------------------------------------------------------------
echo
if [ $# -ne 0 ] || [ "${G3_REV}" != "${PINNED_REV}" ]; then
	echo "passing: $(wc -l <"${WORK}/passing.txt" | tr -d ' ') (${WORK}/passing.txt; ${RECORD} not compared)"
	exit 0
fi

recorded="${WORK}/recorded.txt"
if [ -f "${RECORD}" ]; then grep -v '^#' "${RECORD}" | LC_ALL=C sort -u >"${recorded}"; else : >"${recorded}"; fi
lost="$(LC_ALL=C comm -23 "${recorded}" "${WORK}/passing.txt")"
gained="$(LC_ALL=C comm -13 "${recorded}" "${WORK}/passing.txt")"

if [ -n "${lost}" ] && [ "${FORCE_RECORD:-0}" != 1 ]; then
	echo "passing: $(wc -l <"${WORK}/passing.txt" | tr -d ' ') (${WORK}/passing.txt)"
	[ -z "${gained}" ] || { echo "newly passing:"; sed 's/^/  /' <<<"${gained}"; }
	echo "NO LONGER PASSING (${RECORD} left unchanged; FORCE_RECORD=1 records the loss):"
	sed 's/^/  /' <<<"${lost}"
	exit 1
fi

{
	echo "# Spec gates that pass against the Go engine (cmd/olympia), one per line."
	echo "# Written by tools/olympia-spec/run.sh; do not edit by hand."
	echo "# G3 revision: ${PINNED_REV}"
	echo "# A name is a tests/olympia/regress fixture; \"golden\" is golden_check.sh"
	echo "# after run/olympia-g3.sh."
	cat "${WORK}/passing.txt"
} >"${RECORD}"

echo "passing: $(wc -l <"${WORK}/passing.txt" | tr -d ' ') (recorded in ${RECORD})"
[ -z "${gained}" ] || { echo "newly passing:"; sed 's/^/  /' <<<"${gained}"; }
[ -z "${lost}" ] || { echo "no longer passing (recorded, FORCE_RECORD=1):"; sed 's/^/  /' <<<"${lost}"; }
[ -n "${gained}${lost}" ] || echo "no change from the recorded list"
