#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: G3's line readers (olympia/z.c) and the field scan and
# print helpers of olympia/io.c (C, the reference) against internal/olympia
# (Go).
#
# Exports the pinned G3 revision, then builds a C driver against its
# olympia/z.c, olympia/code.c, olympia/glob.c and lib/ (with stubs for the
# routines they call that are not under test). The io.c helpers are static,
# so the driver includes a copy of the top of io.c, from its first line to
# just before print_loc_info(), taken from the export, along with copies of
# sout() (olympia/sout.c), set_bit() and test_bit() (olympia/u.c) and
# int_comp() (olympia/seed.c). The driver runs
# internal/olympia/testdata/io-cases.txt one section at a time, each in a
# fresh process and a fresh, empty directory, since a section may crash.
#
# The driver runs twice, to check determinism, and its output must equal
# internal/olympia/testdata/io-c-output.txt, which the Go test
# (TestIOMatchesC in internal/olympia) holds the port to.
#
# Usage: tools/io-compat/run.sh [--update]
#
#   --update  write the C output to internal/olympia/testdata/io-c-output.txt
#             instead of comparing with it. Expected output only ever comes
#             from the C build; never edit that file by hand.
#
# Environment:
#   G3_REPO  G3 clone that contains G3_REV
#            (default ~/Software/playbymail/olympia-g3)
#   G3_REV   G3 revision to test against (default: the pinned revision)
#   WORK     scratch directory, emptied first (default tmp/io-compat)
#   CC       C compiler (default cc)
#
# Exits 0 only if the C output matches the recorded output and the Go test
# passes.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"

. "${HERE}/../lib/timeout.sh"
RUN_TIMEOUT="${RUN_TIMEOUT:-10}"

G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-f73e48118805ffe8ff5e763bfa4f20ae00c3a705}"
WORK="${WORK:-${ROOT}/tmp/io-compat}"
CC="${CC:-cc}"
CASES="${ROOT}/internal/olympia/testdata/io-cases.txt"
EXPECT="${ROOT}/internal/olympia/testdata/io-c-output.txt"

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
#   stubs.c       a stub for every command routine cmd_tbl names (glob.c),
#                 and for the routines code.c calls only from
#                 display_name(); sout() from sout.c; set_bit() and
#                 test_bit() from u.c; int_comp() from seed.c
#   io_helpers.c  io.c up to print_loc_info()
python3 -I - "${WORK}/g3/olympia" "${WORK}/gen" <<'EOF'
import re, sys

src, out = sys.argv[1], sys.argv[2]

def func(path, pattern, what):
    text = open(f"{src}/{path}").read()
    m = re.search(pattern, text, re.M | re.S)
    if not m:
        sys.exit(f"{what} not found in {path}")
    return m.group(0)

glob = open(f"{src}/glob.c").read()
tbl = glob[glob.index("struct cmd_tbl_ent cmd_tbl[]"):glob.index("char *kind_s[]")]
routines = set()
for m in re.finditer(r"^\{[^,]*,[^,]*,\s*(\w+),\s*(\w+),\s*(\w+),", tbl, re.M):
    routines.update(r for r in m.groups() if r != "NULL")
with open(f"{out}/stubs.c", "w") as f:
    f.write('#include <stdarg.h>\n#include <stdio.h>\n#include <stdlib.h>\n'
            '#include "z.h"\n#include "oly.h"\n\n')
    for r in sorted(routines):
        f.write(f"int {r}(struct command *c) {{ return 0; }}\n")
    f.write("""
/* Called only by display_name() and its users, which are not under test. */
char *cap(char *s) { return s; }
int is_port_city(int where) { return 0; }
char *nice_num(int n) { return ""; }

#define		MAX_BUFFER	50
#define		BUF_LEN		LEN

""")
    f.write(func("sout.c", r"^char \*\nsout\(const char \*format, \.\.\.\)\n\{.*?^\}\n", "sout()"))
    f.write("\n")
    f.write(func("u.c", r"^int\ntest_bit\(sparse kr, int i\)\n\{.*?^\}\n", "test_bit()"))
    f.write("\n")
    f.write(func("u.c", r"^void\nset_bit\(sparse \*kr, int i\)\n\{.*?^\}\n", "set_bit()"))
    f.write("\n")
    f.write(func("seed.c", r"^int\nint_comp\(const void \*av, const void \*bv\)\n\{.*?^\}\n", "int_comp()"))

io = open(f"{src}/io.c").read()
end = io.find("static void\nprint_loc_info(")
if end < 0:
    sys.exit("print_loc_info() not found in io.c")
with open(f"{out}/io_helpers.c", "w") as f:
    f.write(io[:end])
EOF

cat >"${WORK}/driver.c" <<'EOF'
/*
 *  Differential driver for tools/io-compat/run.sh. Runs one section of
 *  internal/olympia/testdata/io-cases.txt; see that file for the
 *  operations.
 *
 *  usage: driver CASES SECTION   (SECTION counts from 0)
 */
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

/* io.c's top includes z.h and oly.h, which have no include guards */
#include "io_helpers.c"

void glob_init(void);
void init_lower(void);
int int_comp(const void *av, const void *bv);

int pretty_data_files = FALSE;

/* the operation's text: decode \xHH and \\ in place, and return its length */
static int
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
	return (int) (q - s);
}

/* s, printable: bytes outside ' '..'~' and '\' as \xHH */
static char *
esc(const char *s)
{
	static char buf[4][16384];
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

static void
show_str(const char *what, const char *s)
{
	if (s)
		printf("%s [%s] len %d\n", what, esc(s), (int) strlen(s));
	else
		printf("%s NULL\n", what);
}

static char *pos;

static int
num(void)
{
	return (int) strtol(pos, &pos, 0);
}

/* TEXT: the rest of the line after one space; its length in *len */
static char *
text(int *len)
{
	int n;

	if (*pos == ' ')
		pos++;
	n = decode(pos);
	if (len)
		*len = n;
	return pos;
}

/* NAME: the next word */
static char *
word(void)
{
	static char buf[256];
	int nread = 0;

	buf[0] = '\0';
	sscanf(pos, "%255s%n", buf, &nread);
	pos += nread;
	return buf;
}

/* the data the helpers scan into and print from */
static ilist L;
static sparse K;
static skill_ents_list SL;
static item_ents_list IL;
static trades_list TL;
static req_ents_list RL;
static struct entity_player P;
static olytime OT;
static FILE *gfp;

static void
show_ilist(const char *what, ilist l)
{
	int i;

	printf("%s", what);
	if (l == NULL)
		printf(" NULL");
	for (i = 0; i < ilist_len(l); i++)
		printf(" %d", l[i]);
	printf("\n");
}

static void
show_skills(void)
{
	int i;

	printf("skills:");
	for (i = 0; i < skill_ents_len(SL); i++)
		printf(" [%d %d %d %d]", SL[i]->skill, SL[i]->know,
		       SL[i]->days_studied, SL[i]->experience);
	printf("\n");
}

static void
show_items(void)
{
	int i;

	printf("items:");
	for (i = 0; i < item_ents_len(IL); i++)
		printf(" [%d %d]", IL[i]->item, IL[i]->qty);
	printf("\n");
}

static void
show_trades(void)
{
	int i;

	printf("trades:");
	for (i = 0; i < trades_len(TL); i++)
		printf(" [%d %d %d %d %d %d %d %d who %d]", TL[i]->kind,
		       TL[i]->item, TL[i]->qty, TL[i]->cost, TL[i]->cloak,
		       TL[i]->have_left, TL[i]->month_prod, TL[i]->expire,
		       TL[i]->who);
	printf("\n");
}

static void
show_reqs(void)
{
	int i;

	printf("reqs:");
	for (i = 0; i < req_ents_len(RL); i++)
		printf(" [%d %d %d]", RL[i]->item, RL[i]->qty, RL[i]->consume);
	printf("\n");
}

static void
show_admits(void)
{
	int i, j;

	printf("admits:");
	for (i = 0; i < admits_len(P.admits); i++)
	{
		printf(" [%d %d", P.admits[i]->targ, P.admits[i]->sense);
		for (j = 0; j < ilist_len(P.admits[i]->l); j++)
			printf(" %d", P.admits[i]->l[j]);
		printf("]");
	}
	printf("\n");
}

static void
show_olytime(void)
{
	printf("olytime: turn %d day %d days %d\n", OT.turn, OT.day,
	       OT.days_since_epoch);
}

/* t_string() of the current line, as the scanners pass it */
static char *
field(void)
{
	return t_string(line);
}

int
main(int argc, char **argv)
{
	static char buf[65536];
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

	/* the part of call_init_routines() in main.c that these need */
	init_lower();
	glob_init();

	while (fgets(buf, sizeof(buf), fp) != NULL)
	{
		char op[32];
		int nread = 0, n, k, sk, len;
		char *s;

		buf[strcspn(buf, "\n")] = '\0';
		if (buf[0] == '#' || buf[0] == '\0')
			continue;
		if (strncmp(buf, "section ", 8) == 0)
		{
			if (++section == want)
				printf("%s\n", buf);
			continue;
		}
		if (section != want)
			continue;

		printf("> %s\n", buf);
		op[0] = '\0';
		sscanf(buf, "%31s%n", op, &nread);
		pos = buf + nread;

		if (strcmp(op, "alloc") == 0)
		{
			n = num();
			k = num();
			sk = num();
			alloc_box(n, k, sk);
		}
		else if (strcmp(op, "setname") == 0)
		{
			n = num();
			set_name(n, text(NULL));
		}
		else if (strcmp(op, "namen") == 0)
		{
			static char big[8192];

			n = num();
			k = num();
			s = text(&len);
			big[0] = '\0';
			while (k-- > 0)
				strcat(big, s);
			set_name(n, big);
		}
		else if (strcmp(op, "pretty") == 0)
			pretty_data_files = num();
		else if (strcmp(op, "extbox") == 0)
			ext_boxnum = num();
		else if (strcmp(op, "put") == 0 || strcmp(op, "putn") == 0)
		{
			FILE *out;
			char name[256];

			strcpy(name, word());
			k = strcmp(op, "putn") == 0 ? num() : 1;
			s = text(&len);
			out = fopen(name, "ab");
			while (k-- > 0)
				fwrite(s, 1, (size_t) len, out);
			fclose(out);
		}
		else if (strcmp(op, "readfile") == 0)
			printf("readfile %d\n", readfile(word()));
		else if (strcmp(op, "closefile") == 0)
			closefile(word());
		else if (strcmp(op, "readlin") == 0)
			show_str("readlin", readlin());
		else if (strcmp(op, "readlin_ew") == 0)
			show_str("readlin_ew", readlin_ew());
		else if (strcmp(op, "skip") == 0)
		{
			for (k = num(); k > 0; k--)
				readlin();
		}
		else if (strcmp(op, "drain") == 0)
		{
			static char last[LEN];

			last[0] = '\0';
			for (k = 0; (s = readlin()) != NULL; k++)
				strcpy(last, s);
			printf("drain %d lines, last [%s]\n", k, esc(last));
		}
		else if (strcmp(op, "fopen") == 0)
		{
			gfp = fopen(word(), "r");
			printf("fopen %s\n", gfp ? "ok" : "NULL");
		}
		else if (strcmp(op, "getlin") == 0)
			show_str("getlin", getlin(gfp));
		else if (strcmp(op, "getlin_ew") == 0)
			show_str("getlin_ew", getlin_ew(gfp));
		else if (strcmp(op, "ew") == 0)
			show_str("ew", eat_leading_trailing_whitespace(text(NULL)));
		else if (strcmp(op, "hash") == 0)
		{
			s = text(NULL);
			printf("linehash %d t_string [%s]\n", linehash(s), esc(t_string(s)));
		}
		else if (strcmp(op, "next") == 0)
		{
			line = readlin();
			if (line && iswhite(*line))
				line++;
			if (line)
				printf("line [%s] hash %d t [%s]\n", esc(line),
				       linehash(line), esc(t_string(line)));
			else
				printf("line NULL\n");
		}
		else if (strcmp(op, "box_scan") == 0)
			printf("box_scan %d\n", box_scan(field()));
		else if (strcmp(op, "boxlist_scan") == 0)
		{
			boxlist_scan(field(), num(), &L);
			show_ilist("L:", L);
		}
		else if (strcmp(op, "boxlist0_scan") == 0)
		{
			boxlist0_scan(field(), num(), &L);
			show_ilist("L:", L);
		}
		else if (strcmp(op, "ilist_scan") == 0)
		{
			ilist_scan(field(), &L);
			show_ilist("L:", L);
		}
		else if (strcmp(op, "known_scan") == 0)
		{
			known_scan(field(), &K, num());
			show_ilist("K:", K);
		}
		else if (strcmp(op, "admit_scan") == 0)
		{
			admit_scan(field(), num(), &P);
			show_admits();
		}
		else if (strcmp(op, "skill_scan") == 0)
		{
			skill_list_scan(field(), &SL, num());
			show_skills();
		}
		else if (strcmp(op, "item_scan") == 0)
		{
			item_list_scan(field(), &IL, num());
			show_items();
		}
		else if (strcmp(op, "trade_scan") == 0)
		{
			trade_list_scan(field(), &TL, num());
			show_trades();
		}
		else if (strcmp(op, "req_scan") == 0)
		{
			req_list_scan(field(), &RL, num());
			show_reqs();
		}
		else if (strcmp(op, "olytime_scan") == 0)
		{
			olytime_scan(field(), &OT);
			show_olytime();
		}
		else if (strcmp(op, "box_print") == 0)
		{
			n = num();
			box_print(stdout, text(NULL), n);
		}
		else if (strcmp(op, "boxlist_print") == 0)
			boxlist_print(stdout, text(NULL), L);
		else if (strcmp(op, "boxlist0_print") == 0)
			boxlist0_print(stdout, text(NULL), L);
		else if (strcmp(op, "ilist_print") == 0)
			ilist_print(stdout, text(NULL), L);
		else if (strcmp(op, "known_print") == 0)
		{
			known_print(stdout, text(NULL), K);
			show_ilist("K:", K);
		}
		else if (strcmp(op, "admit_print") == 0)
			admit_print(stdout, &P);
		else if (strcmp(op, "skill_print") == 0)
			skill_list_print(stdout, text(NULL), SL);
		else if (strcmp(op, "item_print") == 0)
			item_list_print(stdout, text(NULL), IL);
		else if (strcmp(op, "trade_print") == 0)
			trade_list_print(stdout, text(NULL), TL);
		else if (strcmp(op, "req_print") == 0)
			req_list_print(stdout, text(NULL), RL);
		else if (strcmp(op, "olytime_print") == 0)
			olytime_print(stdout, text(NULL), &OT);
		else if (strcmp(op, "set") == 0 || strcmp(op, "setk") == 0)
		{
			ilist *l = strcmp(op, "set") == 0 ? &L : &K;
			char *end;

			ilist_clear(l);
			for (;;)
			{
				n = (int) strtol(pos, &end, 0);
				if (end == pos)
					break;
				pos = end;
				ilist_append(l, n);
			}
		}
		else if (strcmp(op, "skill") == 0)
		{
			struct skill_ent *e = my_malloc(sizeof(*e));

			e->skill = num();
			e->know = (char) num();
			e->days_studied = num();
			e->experience = (short) num();
			skill_ents_append(&SL, e);
		}
		else if (strcmp(op, "item") == 0)
		{
			struct item_ent *e = my_malloc(sizeof(*e));

			e->item = num();
			e->qty = num();
			item_ents_append(&IL, e);
		}
		else if (strcmp(op, "trade") == 0)
		{
			struct trade *e = my_malloc(sizeof(*e));

			e->kind = num();
			e->item = num();
			e->qty = num();
			e->cost = num();
			e->cloak = num();
			e->have_left = num();
			e->month_prod = num();
			e->expire = num();
			trades_append(&TL, e);
		}
		else if (strcmp(op, "req") == 0)
		{
			struct req_ent *e = my_malloc(sizeof(*e));

			e->item = num();
			e->qty = num();
			e->consume = (schar) num();
			req_ents_append(&RL, e);
		}
		else if (strcmp(op, "admit") == 0)
		{
			struct admit *e = my_malloc(sizeof(*e));
			char *end;

			e->targ = num();
			e->sense = num();
			for (;;)
			{
				n = (int) strtol(pos, &end, 0);
				if (end == pos)
					break;
				pos = end;
				ilist_append(&e->l, n);
			}
			admits_append(&P.admits, e);
		}
		else if (strcmp(op, "olytime") == 0)
		{
			OT.turn = (short) num();
			OT.day = (short) num();
			OT.days_since_epoch = num();
		}
		else if (strcmp(op, "zero") == 0)
		{
			struct loc_info li;

			memset(&li, 0, sizeof(li));
			li.where = num();
			if (num())
			{
				ilist_append(&li.here_list, 1);
				ilist_clear(&li.here_list);
			}
			printf("zero_check %d\n", zero_check(&li, sizeof(li)));
		}
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
	lib/admits.c lib/item_ents.c lib/req_ents.c lib/skill_ents.c lib/trades.c \
	-o "${WORK}/io-driver") >"${WORK}/build.log" 2>&1 ||
	die "C driver build failed; see ${WORK}/build.log"

sections=$(grep -c '^section ' "${CASES}")
# Reject invalid references before accepting statuses or updating recordings.
# stderr is merged by the driver itself, so one capture covers both streams.
valid_reference() {
	local status=0
	grep -Ei 'out-of-scope|not (yet )?implemented|missing.*dependency|NYI|stub|Sanitizer|runtime error:' "$1" >/dev/null 2>&1 || status=$?
	[ "$status" = 1 ]
}
# Negative controls include expected-failure captures and unreadable streams.
for marker in 'missing dependency' 'NYI' 'C stub' 'AddressSanitizer' 'runtime error:'; do
	printf '%s\nexit 134\n' "$marker" >"${WORK}/invalid-reference"
	if valid_reference "${WORK}/invalid-reference"; then die "accepted invalid C reference"; fi
done
rm "${WORK}/invalid-reference"
if valid_reference "${WORK}/invalid-reference"; then die "accepted missing C reference"; fi
# The driver sends its stderr to its stdout. Bash reports each crash on its
# own stderr, which is dropped. Each section runs in a fresh directory, for
# the files it writes.
for run in run1 run2; do
	for ((i = 0; i < sections; i++)); do
		status=0
		mkdir "${WORK}/${run}/${i}"
		(cd "${WORK}/${run}/${i}" && run_limited "${RUN_TIMEOUT}" io-driver "${WORK}/io-driver" "${CASES}" "${i}" >>../out.txt) || status=$?
		valid_reference "${WORK}/${run}/out.txt" || die "invalid C reference in section ${i}"
		[ "$status" != 124 ] || die "C section ${i} timed out"
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
(cd "${ROOT}" && RUN_TIMEOUT="${RUN_TIMEOUT}" go test -exec "bash ${HERE}/../lib/run-limited.sh" -count=1 -run 'TestIOMatchesC' ./internal/olympia)
echo "PASS"
