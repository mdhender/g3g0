#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: the host C library's sscanf("%d") (the reference, on
# macOS) against clib.Scanner (Go).
#
# Builds a small C driver that runs each case of
# internal/clib/testdata/scanf-cases.txt through sscanf() and prints the
# return value, the converted values and where scanning stopped. It runs
# twice, to check determinism, and its output must equal
# internal/clib/testdata/scanf-c-output.txt, which the Go test
# (TestScannerMatchesC in internal/clib) holds the port to.
#
# Nothing here comes from G3: the behavior under test is the C library's.
# Record expectations only on the reference platform (macOS).
#
# Usage: tools/scanf-compat/run.sh [--update]
#
#   --update  write the C output to internal/clib/testdata/scanf-c-output.txt
#             instead of comparing with it. Expected output only ever comes
#             from the C build; never edit that file by hand.
#
# Environment:
#   WORK     scratch directory, emptied first (default tmp/scanf-compat)
#   CC       C compiler (default cc)
#
# Exits 0 only if the C output matches the recorded output and the Go test
# passes.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"

WORK="${WORK:-${ROOT}/tmp/scanf-compat}"
CC="${CC:-cc}"
CASES="${ROOT}/internal/clib/testdata/scanf-cases.txt"
EXPECT="${ROOT}/internal/clib/testdata/scanf-c-output.txt"

UPDATE=0
case "${1:-}" in
"") ;;
--update) UPDATE=1 ;;
*) echo "usage: $0 [--update]" >&2; exit 2 ;;
esac

die() { echo "error: $*" >&2; exit 2; }

for tool in go cmp diff "${CC}"; do
	command -v "${tool}" >/dev/null || die "missing required tool: ${tool}"
done
[[ "$(uname -s)" == Darwin ]] ||
	echo "warning: the reference platform is macOS; this C library may differ" >&2

rm -rf "${WORK}"
mkdir -p "${WORK}"

echo "work dir:    ${WORK}"
echo

cat >"${WORK}/driver.c" <<'CEOF'
/*
 *  Differential driver for tools/scanf-compat/run.sh. Reads cases from
 *  stdin (internal/clib/testdata/scanf-cases.txt) and prints one line for
 *  each: the case as written, then "-> r=R a=A b=B c=C n=N". Targets the
 *  conversion does not reach keep their initial values (7, 7, 7, -1).
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static int hexval(int c) {
	if (c >= '0' && c <= '9') return c - '0';
	if (c >= 'a' && c <= 'f') return c - 'a' + 10;
	if (c >= 'A' && c <= 'F') return c - 'A' + 10;
	fprintf(stderr, "bad hex digit %c\n", c);
	exit(2);
}

/* one character of input at *p, escapes decoded */
static int decode1(const char **p) {
	const char *s = *p;
	int c;
	if (s[0] == '\\' && s[1] == 'x') {
		c = hexval(s[2]) * 16 + hexval(s[3]);
		s += 4;
	} else if (s[0] == '\\' && s[1] == '\\') {
		c = '\\';
		s += 2;
	} else {
		c = (unsigned char)*s++;
	}
	*p = s;
	return c;
}

static char buf[1 << 16];

static void decode(const char *s) {
	size_t n = 0;
	while (*s) {
		if (*s == '{') {
			char *end;
			long count = strtol(s + 1, &end, 10);
			if (*end != ':') { fprintf(stderr, "bad repeat\n"); exit(2); }
			s = end + 1;
			int c = decode1(&s);
			if (*s != '}') { fprintf(stderr, "bad repeat\n"); exit(2); }
			s++;
			while (count-- > 0) buf[n++] = (char)c;
		} else {
			buf[n++] = (char)decode1(&s);
		}
	}
	buf[n] = 0;
}

int main(void) {
	static char line[1 << 16];
	while (fgets(line, sizeof line, stdin)) {
		line[strcspn(line, "\n")] = 0;
		if (line[0] == 0 || line[0] == '#')
			continue;
		char *tab = strchr(line, '\t');
		if (!tab) { fprintf(stderr, "no tab: %s\n", line); return 2; }
		*tab = 0;
		const char *fmt = line;
		decode(tab + 1);

		int a = 7, b = 7, n = -1, r;
		char c = 7;
		if (strcmp(fmt, "d") == 0)
			r = sscanf(buf, "%d%n", &a, &n);
		else if (strcmp(fmt, "dd") == 0)
			r = sscanf(buf, "%d%d%n", &a, &b, &n);
		else if (strcmp(fmt, "d,d") == 0)
			r = sscanf(buf, "%d,%d%n", &a, &b, &n);
		else if (strcmp(fmt, "d,d c") == 0)
			r = sscanf(buf, "%d,%d %c%n", &a, &b, &c, &n);
		else { fprintf(stderr, "unknown format %s\n", fmt); return 2; }
		printf("%s\t%s\t-> r=%d a=%d b=%d c=%d n=%d\n", fmt, tab + 1, r, a, b, c, n);
	}
	return 0;
}
CEOF

echo "building the C driver (log: ${WORK}/build.log)"
"${CC}" -std=c11 -Wall -o "${WORK}/scanf-driver" "${WORK}/driver.c" >"${WORK}/build.log" 2>&1 ||
	die "C driver build failed; see ${WORK}/build.log"

echo "running the C driver twice"
"${WORK}/scanf-driver" <"${CASES}" >"${WORK}/c-output-1.txt"
"${WORK}/scanf-driver" <"${CASES}" >"${WORK}/c-output-2.txt"
cmp -s "${WORK}/c-output-1.txt" "${WORK}/c-output-2.txt" ||
	die "the C driver is not deterministic; compare ${WORK}/c-output-1.txt and ${WORK}/c-output-2.txt"
echo "C cases:     $(wc -l <"${WORK}/c-output-1.txt" | tr -d ' ')"

if [[ ${UPDATE} -eq 1 ]]; then
	cp "${WORK}/c-output-1.txt" "${EXPECT}"
	echo "wrote ${EXPECT}"
elif ! cmp -s "${WORK}/c-output-1.txt" "${EXPECT}"; then
	diff -u "${EXPECT}" "${WORK}/c-output-1.txt" | head -40 >&2
	die "C output differs from ${EXPECT}"
else
	echo "C output matches ${EXPECT}"
fi

echo
echo "running the Go test"
cd "${ROOT}"
go test -count=1 -run 'TestScannerMatchesC' ./internal/clib
