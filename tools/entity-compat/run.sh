#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: G3's entity model (olympia/oly.h and olympia/glob.c, C,
# the reference) against internal/olympia (Go).
#
# Exports the pinned G3 revision, then builds a small C driver against its
# olympia/oly.h and olympia/glob.c (with stubs for the command routines
# cmd_tbl names) and lib/checked_alloc.c. The driver prints:
#
#   - every numeric #define in oly.h, with its value;
#   - sizeof, _Alignof and the offset and size of every field, for every
#     struct oly.h defines (and the olytime typedef);
#   - the glob.c tables kind_s, subkind_s and cmd_tbl;
#   - the glob.c globals, after glob_init();
#   - the p_xxx and rp_xxx accessors, kind(), subkind() and valid_box() over
#     a few boxes.
#
# The struct fields and the #defines are read from oly.h by this script, so
# the C side lists everything the header has. The driver runs twice, to
# check determinism, and its output must equal
# internal/olympia/testdata/entity-c-output.txt, which the Go test
# (TestEntityMatchesC in internal/olympia) holds the port to.
#
# Usage: tools/entity-compat/run.sh [--update]
#
#   --update  write the C output to internal/olympia/testdata/entity-c-output.txt
#             instead of comparing with it. Expected output only ever comes
#             from the C build; never edit that file by hand.
#
# Environment:
#   G3_REPO  G3 clone that contains G3_REV
#            (default ~/Software/playbymail/olympia-g3)
#   G3_REV   G3 revision to test against (default: the pinned revision)
#   WORK     scratch directory, emptied first (default tmp/entity-compat)
#   CC       C compiler (default cc)
#
# Exits 0 only if the C output matches the recorded output and the Go test
# passes.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"

G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-6bbcd63055d051ade6f745f368281628441a2c27}"
WORK="${WORK:-${ROOT}/tmp/entity-compat}"
CC="${CC:-cc}"
EXPECT="${ROOT}/internal/olympia/testdata/entity-c-output.txt"

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
#   consts.h  one C(NAME) line per numeric #define in oly.h
#   structs.h one S(TYPE, LABEL) line per struct, followed by its F(TYPE, FIELD)s
#   stubs.c   a stub for every command routine cmd_tbl names
python3 -I - "${WORK}/g3/olympia" "${WORK}/gen" <<'EOF'
import re, sys

src, out = sys.argv[1], sys.argv[2]
oly = open(f"{src}/oly.h").read()

with open(f"{out}/consts.h", "w") as f:
    for m in re.finditer(r"^#define\s+(\w+)\s+\(?-?\d+\)?\s*(?:/\*.*?\*/)?\s*$", oly, re.M):
        f.write(f"C({m.group(1)})\n")

nocomment = re.sub(r"/\*.*?\*/", " ", oly, flags=re.S)
bodies = [(f"struct {m.group(1)}", m.group(1), m.group(2))
          for m in re.finditer(r"^struct\s+(\w+)\s*\{(.*?)^\};", nocomment, re.M | re.S)]
for m in re.finditer(r"^typedef\s+struct\s*\{(.*?)^\}\s*(\w+)\s*;", nocomment, re.M | re.S):
    bodies.append((m.group(2), m.group(2), m.group(1)))

with open(f"{out}/structs.h", "w") as f:
    for ctype, label, body in bodies:
        f.write(f"S({ctype}, {label})\n")
        for decl in body.split(";"):
            decl = " ".join(decl.split())
            if not decl:
                continue
            fp = re.match(r".*\(\s*\*\s*(\w+)\s*\)\s*\(.*\)$", decl)
            if fp:
                f.write(f"F({ctype}, {fp.group(1)})\n")
                continue
            parts = decl.split(",")
            names = [parts[0].split()[-1]] + parts[1:]
            for n in names:
                f.write(f"F({ctype}, {n.strip().lstrip('*')})\n")

glob = open(f"{src}/glob.c").read()
tbl = glob[glob.index("struct cmd_tbl_ent cmd_tbl[]"):glob.index("char *kind_s[]")]
routines = set()
for m in re.finditer(r"^\{[^,]*,[^,]*,\s*(\w+),\s*(\w+),\s*(\w+),", tbl, re.M):
    routines.update(r for r in m.groups() if r != "NULL")
with open(f"{out}/stubs.c", "w") as f:
    f.write('#include <stdio.h>\n#include "z.h"\n#include "oly.h"\n\n')
    for r in sorted(routines):
        f.write(f"int {r}(struct command *c) {{ return 0; }}\n")
EOF

cat >"${WORK}/driver.c" <<'EOF'
/*
 *  Differential driver for tools/entity-compat/run.sh. Prints the entity
 *  model of oly.h and glob.c; see run.sh.
 */
#include <stddef.h>
#include <stdio.h>
#include "z.h"
#include "oly.h"

void glob_init(void);

/*
 *  Pointers print as labels, #1, #2, ... in the order they are first seen,
 *  so the output does not depend on addresses.
 */
static void *seen[64];
static int nseen;

static const char *
ptr(void *p)
{
	static char buf[8][16];
	static int k;
	int i;

	if (p == NULL)
		return "NULL";
	for (i = 0; i < nseen && seen[i] != p; i++)
		;
	if (i == nseen)
		seen[nseen++] = p;
	k = (k + 1) % 8;
	snprintf(buf[k], sizeof(buf[k]), "#%d", i + 1);
	return buf[k];
}

/* zero_check() from io.c */
static int
zero(void *ptr, unsigned size)
{
	char *p = ptr;
	int sum = 0;

	while (size-- > 0)
		sum |= *p++;
	return sum == 0;
}

#define	SHOWI(e)	printf("%s = %d\n", #e, (int) (e))
#define	SHOWP(e)	printf("%s = %s\n", #e, ptr(e))
#define	SHOWZ(e)	printf("%s = %s zero %d\n", #e, ptr(e), zero((e), sizeof(*(e))))

static void
new_box(int n, int k, int sk)
{
	bx[n] = my_malloc(sizeof(struct box));
	bx[n]->kind = k;
	bx[n]->skind = sk;
}

/* P(name, n): rp_name, p_name twice, rp_name; the substructure must be zero. */
#define	P(x, n)	do {				\
	SHOWP(rp_##x(n));			\
	SHOWZ(p_##x(n));			\
	SHOWP(p_##x(n));			\
	SHOWP(rp_##x(n));			\
} while (0)

int
main(void)
{
	int i, n;

	setvbuf(stdout, NULL, _IOLBF, 0);

#define	C(x)	printf("const %s %d\n", #x, (int) (x));
#include "consts.h"
#undef C

#define	S(t, l)	printf("struct %s size %zu align %zu\n", #l, sizeof(t), _Alignof(t));
#define	F(t, f)	printf("\t%s offset %zu size %zu\n", #f, offsetof(t, f), sizeof(((t *) 0)->f));
#include "structs.h"
#undef S
#undef F

	for (i = 0; kind_s[i] != NULL; i++)
		printf("kind_s[%d] [%s]\n", i, kind_s[i]);
	printf("kind_s[%d] NULL\n", i);
	for (i = 0; subkind_s[i] != NULL; i++)
		printf("subkind_s[%d] [%s]\n", i, subkind_s[i]);
	printf("subkind_s[%d] NULL\n", i);
	for (i = 0; cmd_tbl[i].name != NULL; i++)
		printf("cmd_tbl[%d] [%s] [%s] time %d poll %d pri %d\n", i,
		       cmd_tbl[i].allow, cmd_tbl[i].name,
		       cmd_tbl[i].time, cmd_tbl[i].poll, cmd_tbl[i].pri);
	printf("cmd_tbl[%d] NULL\n", i);

	glob_init();

	printf("from_host [%s]\n", from_host);
	printf("reply_host [%s]\n", reply_host);
	printf("gm_address [%s]\n", gm_address);
	printf("game_title [%s]\n", game_title);
	printf("game_url [%s]\n", game_url);
	printf("rules_url [%s]\n", rules_url);
	printf("times_url [%s]\n", times_url);
	printf("htpasswd_loc [%s]\n", htpasswd_loc);
	SHOWI(garrison_pay);
	SHOWI(army_slow_factor);
	SHOWI(auto_quit_turns);
	printf("sysclock %d %d %d\n", sysclock.day, sysclock.turn, sysclock.days_since_epoch);
	SHOWI(indep_player);
	SHOWI(gm_player);
	SHOWI(skill_player);
	SHOWI(eat_pl);
	SHOWI(npc_pl);
	SHOWI(garr_pl);
	SHOWI(combat_pl);
	SHOWI(post_has_been_run);
	SHOWI(garrison_magic);
	SHOWI(mount_olympus);
	for (i = n = 0; i < MAX_BOXES; i++)
		n += bx[i] != NULL;
	printf("bx non-NULL %d\n", n);
	for (i = n = 0; i < T_MAX; i++)
		n += box_head[i] != 0;
	printf("box_head non-zero %d\n", n);
	for (i = n = 0; i < SUB_MAX; i++)
		n += sub_head[i] != 0;
	printf("sub_head non-zero %d\n", n);

	new_box(1001, T_char, 0);
	new_box(1002, T_item, sub_magic);
	new_box(1003, T_char, sub_ni);
	new_box(1004, T_item, sub_artifact);
	new_box(MAX_BOXES - 1, T_loc, sub_ocean);
	new_box(1, T_item, 0);

	SHOWI(kind(-1));
	SHOWI(kind(0));
	SHOWI(kind(1));
	SHOWI(kind(2));
	SHOWI(kind(1001));
	SHOWI(kind(1002));
	SHOWI(kind(1005));
	SHOWI(kind(MAX_BOXES - 1));
	SHOWI(kind(MAX_BOXES));
	SHOWI(kind(MAX_BOXES + 1));
	SHOWI(valid_box(0));
	SHOWI(valid_box(1001));
	SHOWI(valid_box(1005));
	SHOWI(valid_box(MAX_BOXES));
	SHOWI(subkind(0));
	SHOWI(subkind(1001));
	SHOWI(subkind(1002));
	SHOWI(subkind(1003));
	SHOWI(subkind(1005));
	SHOWI(subkind(MAX_BOXES - 1));

	SHOWP(p_loc_info(1001));
	SHOWP(rp_loc_info(1001));
	SHOWI(zero(p_loc_info(1001), sizeof(struct loc_info)));

	P(char, 1001);
	SHOWP(rp_magic(1001));
	SHOWZ(p_magic(1001));
	SHOWP(p_magic(1001));
	SHOWP(rp_magic(1001));
	SHOWP(p_char(1001));

	SHOWP(rp_magic(1003));
	SHOWP(rp_char(1003));
	SHOWZ(p_magic(1003));
	SHOWZ(rp_char(1003));
	SHOWP(p_char(1003));

	P(item, 1002);
	SHOWP(rp_item_magic(1002));
	SHOWZ(p_item_magic(1002));
	SHOWP(p_item_magic(1002));
	SHOWP(rp_item_magic(1002));

	SHOWP(rp_item_magic(1004));
	SHOWP(rp_item(1004));
	SHOWZ(p_item_magic(1004));
	SHOWZ(rp_item(1004));

	P(loc, 1001);
	P(subloc, 1001);
	P(player, 1001);
	P(skill, 1001);
	P(gate, 1001);
	P(misc, 1001);
	P(disp, 1001);
	P(command, 1001);

	SHOWP(rp_char(1002));
	SHOWP(rp_item(1001));
	SHOWP(rp_command(1003));

	return 0;
}
EOF

echo "building the C driver (log: ${WORK}/build.log)"
"${CC}" -std=c11 -w -I "${WORK}/gen" -I "${WORK}/g3/olympia" -I "${WORK}/g3/lib" \
	"${WORK}/driver.c" "${WORK}/g3/olympia/glob.c" "${WORK}/gen/stubs.c" \
	"${WORK}/g3/lib/checked_alloc.c" \
	-o "${WORK}/entity-driver" >"${WORK}/build.log" 2>&1 ||
	die "C driver build failed; see ${WORK}/build.log"

for run in run1 run2; do
	(cd "${WORK}/${run}" && "${WORK}/entity-driver" >out.txt) ||
		die "C driver failed in ${run}"
done
cmp -s "${WORK}/run1/out.txt" "${WORK}/run2/out.txt" || die "C driver is not deterministic"
echo "C driver: $(wc -l <"${WORK}/run1/out.txt" | tr -d ' ') lines, deterministic"

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
(cd "${ROOT}" && go test -count=1 -run 'TestEntityMatchesC' ./internal/olympia)
echo "PASS"
