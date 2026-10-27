#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: G3's entity allocation, kind and subkind chains, codes
# and names (olympia/code.c, the loops of olympia/loop.h and if_name() in
# olympia/io.c, C, the reference) against internal/olympia (Go).
#
# Exports the pinned G3 revision, then builds a C driver against its
# olympia/code.c, olympia/glob.c, olympia/z.c and lib/ (with stubs for the
# few routines code.c calls that are not under test). if_name() is static in
# io.c, so the driver includes a copy of it taken from the export. The
# driver runs internal/olympia/testdata/code-cases.txt one section at a
# time, each in a fresh process, since a section may end in a failed
# assertion, which aborts. code.c is compiled as olympia/code.c, so asfail()
# prints that path, as the Go port does.
#
# The driver runs twice, to check determinism, and its output must equal
# internal/olympia/testdata/code-c-output.txt, which the Go test
# (TestCodeMatchesC in internal/olympia) holds the port to.
#
# Usage: tools/code-compat/run.sh [--update]
#
#   --update  write the C output to internal/olympia/testdata/code-c-output.txt
#             instead of comparing with it. Expected output only ever comes
#             from the C build; never edit that file by hand.
#
# Environment:
#   G3_REPO  G3 clone that contains G3_REV
#            (default ~/Software/playbymail/olympia-g3)
#   G3_REV   G3 revision to test against (default: the pinned revision)
#   WORK     scratch directory, emptied first (default tmp/code-compat)
#   CC       C compiler (default cc)
#
# Exits 0 only if the C output matches the recorded output and the Go test
# passes.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"

G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-66b03fb09b541560268e66b553179d06ea30de80}"
WORK="${WORK:-${ROOT}/tmp/code-compat}"
CC="${CC:-cc}"
CASES="${ROOT}/internal/olympia/testdata/code-cases.txt"
EXPECT="${ROOT}/internal/olympia/testdata/code-c-output.txt"

UPDATE=0
case "${1:-}" in
"") ;;
--update) UPDATE=1 ;;
*) echo "usage: $0 [--update]" >&2; exit 2 ;;
esac

die() { echo "error: $*" >&2; exit 2; }

for tool in git go python3 cmp diff "${CC}"; do
	command -v "${tool}" >/dev/null || die "missing required tool: ${tool}"
done
git -C "${G3_REPO}" cat-file -e "${G3_REV}^{commit}" 2>/dev/null ||
	die "${G3_REPO} does not contain ${G3_REV}; clone https://github.com/playbymail/olympia-g3 and set G3_REPO"

rm -rf "${WORK}"
mkdir -p "${WORK}/g3" "${WORK}/gen" "${WORK}/run1" "${WORK}/run2"

echo "G3 revision: ${G3_REV}"
echo "G3 source:   ${G3_REPO}"
echo "work dir:    ${WORK}"
echo

git -C "${G3_REPO}" archive "${G3_REV}" lib olympia | tar -x -C "${WORK}/g3"

# Generate, from the exported sources:
#   stubs.c   a stub for every command routine cmd_tbl names (glob.c), and
#             for the routines code.c calls that are not under test
#   if_name.c if_name() from io.c
python3 -I - "${WORK}/g3/olympia" "${WORK}/gen" <<'EOF'
import re, sys

src, out = sys.argv[1], sys.argv[2]

glob = open(f"{src}/glob.c").read()
tbl = glob[glob.index("struct cmd_tbl_ent cmd_tbl[]"):glob.index("char *kind_s[]")]
routines = set()
for m in re.finditer(r"^\{[^,]*,[^,]*,\s*(\w+),\s*(\w+),\s*(\w+),", tbl, re.M):
    routines.update(r for r in m.groups() if r != "NULL")
with open(f"{out}/stubs.c", "w") as f:
    f.write('#include <stdarg.h>\n#include <stdio.h>\n#include "z.h"\n#include "oly.h"\n\n')
    for r in sorted(routines):
        f.write(f"int {r}(struct command *c) {{ return 0; }}\n")
    # sout() as olympia/sout.c has it, less its dependencies: a ring of
    # buffers.
    f.write("""
char *
sout(const char *format, ...)
{
	static char ring[16][4096];
	static int k;
	va_list ap;

	k = (k + 1) % 16;
	va_start(ap, format);
	vsnprintf(ring[k], sizeof(ring[k]), format, ap);
	va_end(ap);
	return ring[k];
}

/* Called only by display_name() and its users, which are not under test. */
char *cap(char *s) { return s; }
int is_port_city(int where) { return 0; }
char *nice_num(int n) { return ""; }
""")

io = open(f"{src}/io.c").read()
m = re.search(r"^static char \*\nif_name\(int num\).*?^\}\n", io, re.M | re.S)
if not m:
    sys.exit("if_name() not found in io.c")
with open(f"{out}/if_name.c", "w") as f:
    f.write(m.group(0))
EOF

cat >"${WORK}/driver.c" <<'EOF'
/*
 *  Differential driver for tools/code-compat/run.sh. Runs one section of
 *  internal/olympia/testdata/code-cases.txt; see that file for the
 *  operations.
 *
 *  usage: driver CASES SECTION   (SECTION counts from 0)
 */
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include "z.h"
#include "oly.h"
#include "loop.h"

void glob_init(void);
void init_lower(void);

int pretty_data_files = FALSE;

#include "if_name.c"

/* the operation's text: decode \xHH and \\ in place */
static char *
decode(char *s)
{
	char *p = s, *q = s;

	while (*p)
	{
		if (p[0] == '\\' && p[1] == 'x' && p[2] && p[3])
		{
			char hex[3] = { p[2], p[3], 0 };
			*q++ = (char) strtol(hex, NULL, 16);
			p += 4;
		}
		else if (p[0] == '\\' && p[1] == '\\')
		{
			*q++ = '\\';
			p += 2;
		}
		else
			*q++ = *p++;
	}
	*q = '\0';
	return s;
}

/* s, printable: bytes outside ' '..'~' and '\' as \xHH */
static char *
esc(const char *s)
{
	static char buf[4][8192];
	static int k;
	char *q;

	k = (k + 1) % 4;
	q = buf[k];
	for (; *s; s++)
	{
		unsigned char c = (unsigned char) *s;
		if (c < ' ' || c > '~' || c == '\\')
			q += sprintf(q, "\\x%02x", c);
		else
			*q++ = (char) c;
	}
	*q = '\0';
	return buf[k];
}

static char *pos;

static int
num(void)
{
	return (int) strtol(pos, &pos, 0);
}

/* TEXT: the rest of the line after one space */
static char *
text(void)
{
	if (*pos == ' ')
		pos++;
	return decode(pos);
}

static void
chains(void)
{
	int i, k;

	for (k = 0; k < T_MAX; k++)
		if (kind_first(k))
		{
			printf("kind %d:", k);
			for (i = kind_first(k); i > 0; i = kind_next(i))
				printf(" %d", i);
			printf("\n");
		}
	for (k = 0; k < SUB_MAX; k++)
		if (sub_first(k))
		{
			printf("sub %d:", k);
			for (i = sub_first(k); i > 0; i = sub_next(i))
				printf(" %d", i);
			printf("\n");
		}
}

/* FNV-1a, 64 bits */
static uint64_t sum;

static void
hash(const char *fmt, ...)
{
	char buf[256];
	va_list ap;
	char *p;

	va_start(ap, fmt);
	vsnprintf(buf, sizeof(buf), fmt, ap);
	va_end(ap);
	for (p = buf; *p; p++)
	{
		sum ^= (unsigned char) *p;
		sum *= 1099511628211ULL;
	}
}

static void
codesum(void)
{
	static const char alpha[] = "abelzAEZ059[( -+\x80\xff";
	static const char shape1[] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
	static const char shape2[] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
	char s[8];
	int l, n, len, i, a, b, c, d;
	long idx[4];

	/* int_to_code over the whole box space and past it, and back */
	sum = 14695981039346656037ULL;
	n = 0;
	for (l = -10; l <= MAX_BOXES + 10; l++)
	{
		char *code = int_to_code(l);
		int back = code_to_int(code);

		hash("%d %s %d\n", l, code, back);
		if (back != l)
		{
			printf("int_to_code(%d) = [%s] code_to_int = %d\n", l, code, back);
			n++;
		}
	}
	printf("int_to_code %d..%d: %d do not convert back, sum %016llx\n",
	       -10, MAX_BOXES + 10, n, (unsigned long long) sum);

	/* code_to_int and scode over every string of up to four bytes of alpha */
	sum = 14695981039346656037ULL;
	n = 0;
	for (len = 0; len <= 4; len++)
	{
		int total = 1;

		for (i = 0; i < len; i++)
			total *= (int) strlen(alpha);
		for (l = 0; l < total; l++)
		{
			int v = l;

			for (i = 0; i < len; i++)
			{
				s[i] = alpha[v % strlen(alpha)];
				v /= strlen(alpha);
			}
			s[len] = '\0';
			hash("%d %d\n", code_to_int(s), scode(s));
			n++;
		}
	}
	printf("strings: %d, sum %016llx\n", n, (unsigned long long) sum);

	/* code_to_int over every letter-letter/digit-digit[-digit] shape */
	sum = 14695981039346656037ULL;
	n = 0;
	for (a = 0; shape1[a]; a++)
		for (b = 0; shape2[b]; b++)
			for (c = 0; c < 10; c++)
				for (d = -1; d < 10; d++)
				{
					s[0] = shape1[a];
					s[1] = shape2[b];
					s[2] = (char) ('0' + c);
					s[3] = d < 0 ? '\0' : (char) ('0' + d);
					s[4] = '\0';
					hash("%d\n", code_to_int(s));
					n++;
				}
	printf("shapes: %d, sum %016llx\n", n, (unsigned long long) sum);
}

static void
loop_op(void)
{
	char what[16], mode[16];
	int k, i, first = TRUE;
	int an = 0, ak = 0, ask = 0, del = FALSE;
	int nread;

	what[0] = mode[0] = '\0';
	sscanf(pos, "%15s%n", what, &nread);
	pos += nread;
	if (strcmp(what, "kind") == 0 || strcmp(what, "sub") == 0)
	{
		k = num();
		if (sscanf(pos, "%15s%n", mode, &nread) == 1)
		{
			pos += nread;
			if (strcmp(mode, "del") == 0)
				del = TRUE;
			else if (strcmp(mode, "add") == 0)
			{
				an = num();
				ak = num();
				ask = num();
			}
		}
	}

	printf("visit:");
	if (strcmp(what, "kind") == 0)
	{
		loop_kind(k, i)
		{
			printf(" %d", i);
			if (del)
				delete_box(i);
			if (an && first)
				alloc_box(an, ak, ask);
			first = FALSE;
		}
		next_kind;
	}
	else if (strcmp(what, "sub") == 0)
	{
		loop_subkind(k, i)
		{
			printf(" %d", i);
			if (del)
				delete_box(i);
			if (an && first)
				alloc_box(an, ak, ask);
			first = FALSE;
		}
		next_subkind;
	}
	else if (strcmp(what, "boxes") == 0)
	{
		loop_boxes(i)
		{
			printf(" %d", i);
		}
		next_box;
	}
	else if (strcmp(what, "locship") == 0)
	{
		loop_loc_or_ship(i)
		{
			printf(" %d", i);
		}
		next_loc_or_ship;
	}
	printf("\n");
}

int
main(int argc, char **argv)
{
	static char line[4096];
	FILE *fp;
	int want, section = -1;

	setvbuf(stdout, NULL, _IONBF, 0);
	dup2(1, 2);

	if (argc != 3 || (fp = fopen(argv[1], "r")) == NULL)
	{
		fprintf(stderr, "usage: driver CASES SECTION\n");
		return 2;
	}
	want = atoi(argv[2]);

	/* the part of call_init_routines() in main.c that code.c needs */
	init_lower();
	glob_init();

	while (fgets(line, sizeof(line), fp) != NULL)
	{
		char op[16];
		int nread, n, k, sk;

		line[strcspn(line, "\n")] = '\0';
		if (line[0] == '#' || line[0] == '\0')
			continue;
		if (strncmp(line, "section ", 8) == 0)
		{
			if (++section == want)
				printf("%s\n", line);
			continue;
		}
		if (section != want)
			continue;

		printf("> %s\n", line);
		op[0] = '\0';
		sscanf(line, "%15s%n", op, &nread);
		pos = line + nread;

		if (strcmp(op, "alloc") == 0)
		{
			n = num();
			k = num();
			sk = num();
			alloc_box(n, k, sk);
		}
		else if (strcmp(op, "delete") == 0)
			delete_box(num());
		else if (strcmp(op, "kind") == 0)
		{
			n = num();
			change_box_kind(n, num());
		}
		else if (strcmp(op, "sub") == 0)
		{
			n = num();
			change_box_subkind(n, num());
		}
		else if (strcmp(op, "chains") == 0)
			chains();
		else if (strcmp(op, "loop") == 0)
			loop_op();
		else if (strcmp(op, "usage") == 0)
			print_box_usage();
		else if (strcmp(op, "name") == 0)
			printf("name [%s]\n", esc(name(num())));
		else if (strcmp(op, "setname") == 0)
		{
			n = num();
			set_name(n, text());
		}
		else if (strcmp(op, "banner") == 0)
		{
			n = num();
			set_banner(n, text());
			if (rp_misc(n)->display)
				printf("banner [%s]\n", esc(rp_misc(n)->display));
			else
				printf("banner NULL\n");
		}
		else if (strcmp(op, "ifname") == 0)
		{
			n = num();
			pretty_data_files = num();
			printf("if_name [%s]\n", esc(if_name(n)));
		}
		else if (strcmp(op, "loc") == 0)
			printf("loc %d\n", loc(num()));
		else if (strcmp(op, "boxcode") == 0)
		{
			n = num();
			printf("box_code [%s] box_code_less [%s]\n", box_code(n), box_code_less(n));
		}
		else if (strcmp(op, "code") == 0)
		{
			char *s = int_to_code(num());
			printf("int_to_code [%s] code_to_int %d\n", s, code_to_int(s));
		}
		else if (strcmp(op, "int") == 0)
			printf("code_to_int %d\n", code_to_int(text()));
		else if (strcmp(op, "scode") == 0)
			printf("scode %d\n", scode(text()));
		else if (strcmp(op, "letter") == 0)
		{
			char *s = text();

			printf("letter_val:");
			for (; *s; s++)
				printf(" %d", letter_val(*s, "abcdfghjkmnpqrstvwxz"));
			printf("\n");
		}
		else if (strcmp(op, "codesum") == 0)
			codesum();
		else
		{
			printf("unknown operation\n");
			return 2;
		}
	}

	return 0;
}
EOF

echo "building the C driver (log: ${WORK}/build.log)"
(cd "${WORK}/g3" && "${CC}" -std=c11 -w -I "${WORK}/gen" -I olympia -I lib \
	"${WORK}/driver.c" olympia/code.c olympia/glob.c olympia/z.c \
	"${WORK}/gen/stubs.c" lib/checked_alloc.c lib/ilist.c lib/rng.c lib/rnd.c \
	-o "${WORK}/code-driver") >"${WORK}/build.log" 2>&1 ||
	die "C driver build failed; see ${WORK}/build.log"

sections=$(grep -c '^section ' "${CASES}")
# The driver sends its stderr to its stdout. Bash reports each abort on its
# own stderr, which is dropped.
for run in run1 run2; do
	for ((i = 0; i < sections; i++)); do
		status=0
		(cd "${WORK}/${run}" && exec "${WORK}/code-driver" "${CASES}" "${i}" >>out.txt) || status=$?
		echo "exit ${status}" >>"${WORK}/${run}/out.txt"
	done
done 2>/dev/null
cmp -s "${WORK}/run1/out.txt" "${WORK}/run2/out.txt" || die "C driver is not deterministic"
echo "C driver: ${sections} sections, $(wc -l <"${WORK}/run1/out.txt" | tr -d ' ') lines, deterministic"

if [ "${UPDATE}" -eq 1 ]; then
	cp "${WORK}/run1/out.txt" "${EXPECT}"
	echo "wrote ${EXPECT}"
elif ! cmp -s "${WORK}/run1/out.txt" "${EXPECT}"; then
	echo "FAIL: C output differs from ${EXPECT}"
	diff -u "${EXPECT}" "${WORK}/run1/out.txt" | head -40
	exit 1
else
	echo "C output matches ${EXPECT#"${ROOT}/"}"
fi

echo "running the Go port against it"
(cd "${ROOT}" && go test -count=1 -run 'TestCodeMatchesC' ./internal/olympia)
echo "PASS"
