#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: G3's entity sections, load_box() and save_box()
# (olympia/io.c) and the command parser they use (oly_parse_cmd() in
# olympia/input.c), C being the reference, against internal/olympia (Go).
#
# Exports the pinned G3 revision, then builds a C driver against its
# olympia/z.c, olympia/code.c, olympia/glob.c and lib/ (with stubs for the
# routines they call that are not under test). load_box() and the section
# scanners are static, so the driver includes a copy of io.c from its first
# line through save_box(), taken from the export, along with copies of
# sout() (olympia/sout.c), set_bit() and test_bit() (olympia/u.c),
# int_comp() (olympia/seed.c), and parse_line(), remove_comment(),
# remove_ctrl_chars(), find_command() and oly_parse_cmd() (olympia/input.c).
#
# It runs two checks:
#
# 1. internal/olympia/testdata/box-cases.txt, one section at a time, each
#    in a fresh process and a fresh, empty directory, since a section may
#    crash. The driver runs twice, to check determinism, and its output
#    must equal internal/olympia/testdata/box-c-output.txt, which the Go
#    test (TestBoxMatchesC in internal/olympia) holds the port to.
#
# 2. Real saved games: for each game's lib directory, the driver allocates
#    every box from the box files' "n kind subkind" lines, loads every box
#    file with load_box(), and saves every box with save_box(), plain and
#    with pretty_data_files set. TestBoxWorlds does the same in Go, and
#    the outputs must be equal, byte for byte. The games are world.tgz from
#    the export and, if tools/olympia-roundtrip has been run, every input
#    it has cached under tmp/olympia-roundtrip/cases/*/input/lib.
#
# Usage: tools/box-compat/run.sh [--update]
#
#   --update  write the C output of check 1 to
#             internal/olympia/testdata/box-c-output.txt instead of
#             comparing with it. Expected output only ever comes from the C
#             build; never edit that file by hand.
#
# Environment:
#   G3_REPO  G3 clone that contains G3_REV
#            (default ~/Software/playbymail/olympia-g3)
#   G3_REV   G3 revision to test against (default: the pinned revision)
#   WORK     scratch directory, emptied first (default tmp/box-compat)
#   GAMES    directory of cached games (default tmp/olympia-roundtrip/cases)
#   CC       C compiler (default cc)
#
# Exits 0 only if the C output matches the recorded output, the Go test
# passes and every saved game gives the same output in C and Go.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"

G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-6bbcd63055d051ade6f745f368281628441a2c27}"
WORK="${WORK:-${ROOT}/tmp/box-compat}"
GAMES="${GAMES:-${ROOT}/tmp/olympia-roundtrip/cases}"
CC="${CC:-cc}"
CASES="${ROOT}/internal/olympia/testdata/box-cases.txt"
EXPECT="${ROOT}/internal/olympia/testdata/box-c-output.txt"

UPDATE=0
case "${1:-}" in
"") ;;
--update) UPDATE=1 ;;
*) echo "usage: $0 [--update]" >&2; exit 2 ;;
esac

die() { echo "error: $*" >&2; exit 2; }

for tool in git go python3 tar cmp diff "${CC}"; do
	command -v "${tool}" >/dev/null || die "missing required tool: ${tool}"
done
git -C "${G3_REPO}" cat-file -e "${G3_REV}^{commit}" 2>/dev/null ||
	die "${G3_REPO} does not contain ${G3_REV}; clone https://github.com/playbymail/olympia-g3 and set G3_REPO"

rm -rf "${WORK}"
mkdir -p "${WORK}/g3" "${WORK}/gen" "${WORK}/run1" "${WORK}/run2" "${WORK}/world" "${WORK}/games"

echo "G3 revision: ${G3_REV}"
echo "G3 source:   ${G3_REPO}"
echo "work dir:    ${WORK}"
echo

git -C "${G3_REPO}" archive "${G3_REV}" lib olympia tests/olympia/fixtures/world.tgz | tar -x -C "${WORK}/g3"

# Generate, from the exported sources:
#   stubs.c      a stub for every command routine cmd_tbl names (glob.c),
#                and for the routines code.c calls only from
#                display_name(); sout() from sout.c; set_bit() and
#                test_bit() from u.c; int_comp() from seed.c; and the
#                command parser of input.c
#   io_boxes.c   io.c up to open_write_fp(), just past save_box()
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
            '#include <string.h>\n#include "z.h"\n#include "oly.h"\n\n')
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
    f.write("\n")
    f.write(func("input.c", r"^char \*\*\nparse_line\(char \*\*l, char \*s\)\n\{.*?^\}\n", "parse_line()"))
    f.write("\n")
    f.write(func("input.c", r"^void\nremove_comment\(char \*s\)\n\{.*?^\}\n", "remove_comment()"))
    f.write("\n")
    f.write(func("input.c", r"^void\nremove_ctrl_chars\(char \*s\)\n\{.*?^\}\n", "remove_ctrl_chars()"))
    f.write("\n")
    f.write(func("input.c", r"^int fuzzy_find;.*?^int\nfind_command\(char \*s\)\n\{.*?^\}\n", "find_command()"))
    f.write("\n")
    f.write(func("input.c", r"^int\noly_parse_cmd\(struct command \*c, char \*s\)\n\{.*?^\}\n", "oly_parse_cmd()"))

io = open(f"{src}/io.c").read()
end = io.find("static FILE *\nopen_write_fp(")
if end < 0 or io.find("\nsave_box(FILE *fp, int n)\n") > end:
    sys.exit("open_write_fp() not found after save_box() in io.c")
with open(f"{out}/io_boxes.c", "w") as f:
    # so that assert() names io.c and its lines, as the engine's build does
    f.write('#line 1 "olympia/io.c"\n')
    f.write(io[:end])
EOF

cat >"${WORK}/driver.c" <<'EOF'
/*
 *  Differential driver for tools/box-compat/run.sh.
 *
 *  usage: driver CASES SECTION   runs one section of box-cases.txt
 *                                (SECTION counts from 0); see that file
 *         driver -w LIBDIR P     loads and saves every box of a saved game
 *                                with pretty_data_files = P: diagnostics
 *                                to stderr, saved boxes to stdout
 */
#include <dirent.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

/* io.c's top includes z.h and oly.h, which have no include guards */
#include "io_boxes.c"

void glob_init(void);
void init_lower(void);
int oly_parse_cmd(struct command *c, char *s);
extern int fuzzy_find;

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
		printf("%s [%s]\n", what, esc(s));
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

static void
show_cmd(struct command *c)
{
	int i;

	printf("cmd %d fuzzy %d conditional %d who %d use_skill %d\n",
	       c->cmd, c->fuzzy, c->conditional, c->who, c->use_skill);
	show_str("line", c->line);
	show_str("parsed_line", c->parsed_line);
	printf("parse:");
	if (c->parse == NULL)
		printf(" NULL");
	for (i = 0; i < cstrings_len(c->parse); i++)
		printf(" [%s]", esc(c->parse[i]));
	printf("\n");
}

/* load_box() for each box line from line on, as read_boxes() does, but a
   line that is not a box number ends the file, where read_boxes() would
   loop forever on it */
static void
load_all(void)
{
	int n;

	while (line)
	{
		n = atoi(line);
		if (n > 0)
			load_box(n);
		else
		{
			printf("stop [%s]\n", esc(line));
			break;
		}
	}
}

static const char *box_files[] = {
	"loc", "item", "skill", "gate", "road", "ship", "unform", "misc", NULL
};

static int
str_comp(const void *a, const void *b)
{
	return strcmp(*(char *const *) a, *(char *const *) b);
}

/* the box files of LIBDIR, in read_all_boxes() order, with fact/ sorted */
static int
game_files(const char *libdir, char ***files)
{
	static char *f[100000];
	static char buf[LEN];
	int n = 0, i, first;
	DIR *d;
	struct dirent *e;

	for (i = 0; box_files[i]; i++)
		f[n++] = strdup(box_files[i]);

	snprintf(buf, sizeof(buf), "%s/fact", libdir);
	if ((d = opendir(buf)) != NULL)
	{
		first = n;
		while ((e = readdir(d)) != NULL)
			if (e->d_name[0] >= '0' && e->d_name[0] <= '9')
			{
				snprintf(buf, sizeof(buf), "fact/%s", e->d_name);
				f[n++] = strdup(buf);
			}
		closedir(d);
		qsort(f + first, n - first, sizeof(char *), str_comp);
	}

	*files = f;
	return n;
}

/* alloc_box() for every "n kind subkind" line of a box file: the first
   line, and each line after a blank one */
static void
alloc_file(const char *path)
{
	static char buf[65536];
	FILE *fp = fopen(path, "r");
	int header = 1, n, k, sk;
	char *s, *t;

	if (fp == NULL)
		return;
	while (fgets(buf, sizeof(buf), fp) != NULL)
	{
		buf[strcspn(buf, "\n")] = '\0';
		if (buf[0] == '\0')
		{
			header = 1;
			continue;
		}
		if (!header || buf[0] == '#')
			continue;
		header = 0;
		n = atoi(buf);
		s = strchr(buf, ' ');
		if (n <= 0 || s == NULL)
			continue;
		s++;
		t = strchr(s, ' ');
		if (t == NULL)
			continue;
		*t++ = '\0';
		for (k = 0; kind_s[k] && strcmp(kind_s[k], s) != 0; k++)
			;
		if (strcmp(t, "0") == 0)
			sk = 0;
		else
			for (sk = 0; subkind_s[sk] && strcmp(subkind_s[sk], t) != 0; sk++)
				;
		if (kind_s[k] == NULL || subkind_s[sk] == NULL)
			continue;
		alloc_box(n, k, sk);
	}
	fclose(fp);
}

static int
world(const char *dir, int pretty)
{
	static char path[LEN];
	static char outbuf[1 << 20];
	char **files;
	int n, i;

	setvbuf(stdout, outbuf, _IOFBF, sizeof(outbuf));
	setvbuf(stderr, NULL, _IONBF, 0);
	pretty_data_files = pretty;
	libdir = (char *) dir;

	n = game_files(dir, &files);
	for (i = 0; i < n; i++)
	{
		snprintf(path, sizeof(path), "%s/%s", dir, files[i]);
		alloc_file(path);
	}
	for (i = 0; i < n; i++)
	{
		snprintf(path, sizeof(path), "%s/%s", dir, files[i]);
		if (!readfile(path))
			continue;
		advance;
		load_all();
		closefile(path);
	}
	for (i = 0; i < MAX_BOXES; i++)
		if (bx[i] != NULL)
			save_box(stdout, i);
	fflush(stdout);
	return 0;
}

int
main(int argc, char **argv)
{
	static char buf[65536];
	static struct command C;
	FILE *fp;
	int want, section = -1;

	/* the part of call_init_routines() in main.c that these need */
	init_lower();
	glob_init();

	if (argc == 4 && strcmp(argv[1], "-w") == 0)
		return world(argv[2], atoi(argv[3]));

	setvbuf(stdout, NULL, _IONBF, 0);
	dup2(1, 2);

	if (argc != 3 || (fp = fopen(argv[1], "r")) == NULL)
	{
		fprintf(stderr, "usage: driver CASES SECTION\n");
		return 2;
	}
	want = atoi(argv[2]);

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
		else if (strcmp(op, "pretty") == 0)
			pretty_data_files = num();
		else if (strcmp(op, "libdir") == 0)
			libdir = strdup(text(NULL));
		else if (strcmp(op, "put") == 0 || strcmp(op, "putl") == 0)
		{
			FILE *out;
			char name[256];

			strcpy(name, word());
			s = text(&len);
			out = fopen(name, "ab");
			fwrite(s, 1, (size_t) len, out);
			if (strcmp(op, "putl") == 0)
				fputc('\n', out);
			fclose(out);
		}
		else if (strcmp(op, "readfile") == 0)
			printf("readfile %d\n", readfile(word()));
		else if (strcmp(op, "advance") == 0)
			advance;
		else if (strcmp(op, "line") == 0)
			show_str("line", line);
		else if (strcmp(op, "load") == 0)
		{
			load_box(num());
			show_str("line", line);
		}
		else if (strcmp(op, "loadall") == 0)
			load_all();
		else if (strcmp(op, "save") == 0)
		{
			n = num();
			save_box(stdout, n);
			if (n > 0 && n < MAX_BOXES && bx[n] != NULL)
				printf("temp %d\n", bx[n]->temp);
		}
		else if (strcmp(op, "saveall") == 0)
		{
			for (n = 0; n < MAX_BOXES; n++)
				if (bx[n] != NULL)
					save_box(stdout, n);
		}
		else if (strcmp(op, "cmd") == 0)
		{
			n = num();	/* p_command() evaluates n more than once */
			show_cmd(p_command(n));
		}
		else if (strcmp(op, "parse") == 0)
		{
			s = text(NULL);
			k = oly_parse_cmd(&C, s);
			printf("oly_parse_cmd %d fuzzy_find %d\n", k, fuzzy_find);
			show_str("s", s);
			show_cmd(&C);
		}
		else if (strcmp(op, "find") == 0)
		{
			k = find_command(text(NULL));
			printf("find_command %d fuzzy_find %d\n", k, fuzzy_find);
		}
		else if (strcmp(op, "p") == 0)
		{
			char what[64];

			strcpy(what, word());
			n = num();
			if (strcmp(what, "char") == 0)
				p_char(n);
			else if (strcmp(what, "magic") == 0)
				p_magic(n);
			else if (strcmp(what, "loc") == 0)
				p_loc(n);
			else if (strcmp(what, "subloc") == 0)
				p_subloc(n);
			else if (strcmp(what, "item") == 0)
				p_item(n);
			else if (strcmp(what, "item_magic") == 0)
				p_item_magic(n);
			else if (strcmp(what, "player") == 0)
				p_player(n);
			else if (strcmp(what, "skill") == 0)
				p_skill(n);
			else if (strcmp(what, "gate") == 0)
				p_gate(n);
			else if (strcmp(what, "misc") == 0)
				p_misc(n);
			else if (strcmp(what, "disp") == 0)
				p_disp(n);
			else if (strcmp(what, "command") == 0)
				p_command(n);
			else
			{
				printf("unknown p\n");
				return 2;
			}
		}
		else if (strcmp(op, "poke") == 0)
		{
			char what[64];

			strcpy(what, word());
			n = num();
			k = num();
			if (strcmp(what, "loc_next") == 0)
				p_loc(n)->next = k;
			else if (strcmp(what, "prov_dest_empty") == 0)
			{
				ilist_append(&p_loc(n)->prov_dest, k);
				ilist_clear(&p_loc(n)->prov_dest);
			}
			else if (strcmp(what, "mage_worked") == 0)
				p_magic(n)->mage_worked = (char) k;
			else if (strcmp(what, "one_turn_use") == 0)
				p_item_magic(n)->one_turn_use = (char) k;
			else if (strcmp(what, "opium_double") == 0)
				p_misc(n)->opium_double = (char) k;
			else if (strcmp(what, "cmd") == 0)
				p_command(n)->cmd = k;
			else if (strcmp(what, "where") == 0)
				p_loc_info(n)->where = k;
			else
			{
				printf("unknown poke\n");
				return 2;
			}
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
	lib/admits.c lib/cstrings.c lib/item_ents.c lib/req_ents.c lib/skill_ents.c lib/trades.c \
	-o "${WORK}/box-driver") >"${WORK}/build.log" 2>&1 ||
	die "C driver build failed; see ${WORK}/build.log"

# Check 1: the recorded cases.

sections=$(grep -c '^section ' "${CASES}")
# The driver sends its stderr to its stdout. Bash reports each crash on its
# own stderr, which is dropped. Each section runs in a fresh directory, for
# the files it writes.
for run in run1 run2; do
	for ((i = 0; i < sections; i++)); do
		status=0
		mkdir "${WORK}/${run}/${i}"
		(cd "${WORK}/${run}/${i}" && exec "${WORK}/box-driver" "${CASES}" "${i}" >>../out.txt) || status=$?
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
(cd "${ROOT}" && go test -count=1 -run 'TestBoxMatchesC' ./internal/olympia)

# Check 2: the saved games.

tar -xzf "${WORK}/g3/tests/olympia/fixtures/world.tgz" -C "${WORK}/world"
libs=("${WORK}/world/lib")
if [ -d "${GAMES}" ]; then
	for d in "${GAMES}"/*/input/lib; do
		[ -d "${d}" ] && libs+=("${d}")
	done
fi

list="${WORK}/games/list.txt"
: >"${list}"
k=0
for lib in "${libs[@]}"; do
	for pretty in 0 1; do
		out="${WORK}/games/${k}.${pretty}"
		status=0
		"${WORK}/box-driver" -w "${lib}" "${pretty}" >"${out}.c.out" 2>"${out}.c.err" || status=$?
		[ "${status}" -eq 0 ] || die "C driver failed (${status}) on ${lib}"
		printf '%s\t%s\t%s\n' "${lib}" "${pretty}" "${out}.go" >>"${list}"
	done
	k=$((k + 1))
done

boxes=$(grep -c '^[0-9][0-9]* [a-z]' "${WORK}/games/0.0.c.out" || true)
total=$(cat "${WORK}"/games/*.0.c.out | grep -c '^[0-9][0-9]* [a-z]' || true)
echo "saved games: ${#libs[@]} (world.tgz: ${boxes} boxes; all games: ${total} boxes), plain and pretty"

(cd "${ROOT}" && BOX_COMPAT_GAMES="${list}" go test -count=1 -run 'TestBoxWorlds' ./internal/olympia)

fail=0
k=0
for lib in "${libs[@]}"; do
	for pretty in 0 1; do
		out="${WORK}/games/${k}.${pretty}"
		for s in out err; do
			if ! cmp -s "${out}.c.${s}" "${out}.go.${s}"; then
				echo "FAIL: ${lib} pretty ${pretty}: std${s} differs"
				diff -u "${out}.c.${s}" "${out}.go.${s}" | head -20
				fail=1
			fi
		done
	done
	k=$((k + 1))
done
[ "${fail}" -eq 0 ] || exit 1
echo "saved games: C and Go output identical"
echo "PASS"
