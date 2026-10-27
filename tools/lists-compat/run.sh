#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: G3's reallocing arrays (C, the reference) against
# internal/lists (Go).
#
# Exports the pinned G3 revision, then builds a small C driver against its
# lib/ilist.c, the fourteen template list modules (lib/accept_ents.c through
# lib/wait_args.c, exit_views.c included), lib/rnd.c and lib/rng.c. The driver
# interprets internal/lists/testdata/cases.txt and prints one result line per
# operation. It runs twice, to check determinism, and its output must equal
# internal/lists/testdata/c-output.txt, which the Go test (TestMatchesC in
# internal/lists) holds the port to.
#
# Usage: tools/lists-compat/run.sh [--update]
#
#   --update  write the C output to internal/lists/testdata/c-output.txt
#             instead of comparing with it. Expected output only ever comes
#             from the C build; never edit that file by hand.
#
# Environment:
#   G3_REPO  G3 clone that contains G3_REV
#            (default ~/Software/playbymail/olympia-g3)
#   G3_REV   G3 revision to test against (default: the pinned revision)
#   WORK     scratch directory, emptied first (default tmp/lists-compat)
#   CC       C compiler (default cc)
#
# Exits 0 only if the C output matches the recorded output and the Go test
# passes.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"

G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-47535543d83fbc5649d54b63b49e578f2b11ab1c}"
WORK="${WORK:-${ROOT}/tmp/lists-compat}"
CC="${CC:-cc}"
CASES="${ROOT}/internal/lists/testdata/cases.txt"
EXPECT="${ROOT}/internal/lists/testdata/c-output.txt"

UPDATE=0
case "${1:-}" in
"") ;;
--update) UPDATE=1 ;;
*) echo "usage: $0 [--update]" >&2; exit 2 ;;
esac

die() { echo "error: $*" >&2; exit 2; }

for tool in git go cmp diff "${CC}"; do
	command -v "${tool}" >/dev/null || die "missing required tool: ${tool}"
done
git -C "${G3_REPO}" cat-file -e "${G3_REV}^{commit}" 2>/dev/null ||
	die "${G3_REPO} does not contain ${G3_REV}; clone https://github.com/playbymail/olympia-g3 and set G3_REPO"

rm -rf "${WORK}"
mkdir -p "${WORK}/g3" "${WORK}/run1" "${WORK}/run2"

echo "G3 revision: ${G3_REV}"
echo "G3 source:   ${G3_REPO}"
echo "work dir:    ${WORK}"
echo

git -C "${G3_REPO}" archive "${G3_REV}" lib | tar -x -C "${WORK}/g3"

# The template modules. Each name is both the function prefix and, with
# _list, the list type; ELEM is the element type the module stores.
MODULES=(
	"accept_ents:struct accept_ent *"
	"admits:struct admit *"
	"cstrings:char *"
	"exit_views:struct exit_view *"
	"fights:struct fight *"
	"flag_ents:struct flag_ent *"
	"item_ents:struct item_ent *"
	"orders:struct order_list *"
	"req_ents:struct req_ent *"
	"roads:struct road *"
	"skill_ents:struct skill_ent *"
	"tiles:struct tile *"
	"trades:struct trade *"
	"wait_args:struct wait_arg *"
)
SOURCES=()
{
	for m in "${MODULES[@]}"; do
		name="${m%%:*}"
		elem="${m#*:}"
		printf 'POINTER_LIST(%s, %s)\n' "${name}" "${elem}"
		SOURCES+=("${WORK}/g3/lib/${name}.c")
	done
} >"${WORK}/modules.h"
{
	printf '#define MODULE_TABLE \\\n'
	for m in "${MODULES[@]}"; do
		printf '\tM(%s) \\\n' "${m%%:*}"
	done
	printf '\n'
} >>"${WORK}/modules.h"

cat >"${WORK}/driver.c" <<'EOF'
/*
 *  Differential driver for tools/lists-compat/run.sh. Reads operations from
 *  stdin (internal/lists/testdata/cases.txt) and prints "<op> -> <result>"
 *  for each. See that file for the operations.
 */
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "rnd.h"
#include "rng.h"
#include "lists.h"

/* Elements of the pointer lists: element k is &objs[k]. */
static char objs[256];

static rng_stream cur = {{0, 0, 0, 0}, 0};

#define NSLOT 3

static void
print_state(int len, int is_null, int *v)
{
	int i;

	if (is_null) {
		printf("NULL\n");
		return;
	}
	printf("len %d:", len);
	for (i = 0; i < len; i++)
		printf(" %d", v[i]);
	printf("\n");
}

/*
 *  LIST_OPS defines do_P(slot, op, a, nf) for list module P whose list type
 *  is LT and element type ET; MK makes an element from an int and UN turns
 *  one back. ADD, REMUNIQ and SHUFRNG add the operations only some modules
 *  have (ilist_add, ilist_rem_value_uniq, and the *_shuffle_rng pair), or are
 *  NO_OP.
 */
#define LIST_OPS(P, LT, ET, MK, UN, ADD, REMUNIQ, SHUFRNG)			\
static LT P##_slots[NSLOT];							\
static int									\
do_##P(int slot, const char *op, int *a, int na)				\
{										\
	LT *l = &P##_slots[slot];						\
	int i, n, vals[1024];							\
										\
	if (!strcmp(op, "append")) P##_append(l, MK(a[0]));			\
	else if (!strcmp(op, "fill")) {						\
		for (i = 0; i < a[0]; i++) P##_append(l, MK(a[1] + i));		\
	}									\
	else if (!strcmp(op, "prepend")) P##_prepend(l, MK(a[0]));		\
	else if (!strcmp(op, "clear")) P##_clear(l);				\
	else if (!strcmp(op, "copy")) *l = P##_copy(P##_slots[a[0]]);		\
	else if (!strcmp(op, "delete")) P##_delete(l, a[0]);			\
	else if (!strcmp(op, "len")) printf("%d ", P##_len(*l));		\
	else if (!strcmp(op, "lookup")) printf("%d ", P##_lookup(*l, MK(a[0])));	\
	else if (!strcmp(op, "reclaim")) P##_reclaim(l);			\
	else if (!strcmp(op, "remvalue")) P##_rem_value(l, MK(a[0]));		\
	else if (!strcmp(op, "scramble")) P##_scramble(*l);			\
	else if (!strcmp(op, "shuffle")) P##_shuffle(*l);			\
	ADD(P, l, op, a)							\
	REMUNIQ(P, l, op, a)							\
	SHUFRNG(P, l, op, a)							\
	else return 0;								\
	n = P##_len(*l);							\
	for (i = 0; i < n; i++) vals[i] = UN((*l)[i]);				\
	print_state(n, *l == NULL, vals);					\
	return 1;								\
}

#define NO_OP(P, l, op, a)
#define OP_ADD(P, l, op, a)	else if (!strcmp(op, "add")) P##_add(l, MK_INT(a[0]));
#define OP_REMUNIQ(P, l, op, a)	else if (!strcmp(op, "remuniq")) P##_rem_value_uniq(l, MK_INT(a[0]));
#define OP_SHUFRNG(P, l, op, a)	else if (!strcmp(op, "shufflerng")) { P##_shuffle_rng(*l, &cur); printf("ctr %u ", cur.counter); }

#define MK_INT(v)	(v)
#define UN_INT(x)	(x)
LIST_OPS(ilist, ilist, int, MK_INT, UN_INT, OP_ADD, OP_REMUNIQ, OP_SHUFRNG)

#define MK_PTR(v)	((void *)&objs[(v)])
#define UN_PTR(x)	((int)((char *)(x) - objs))
#define POINTER_LIST(P, ET)							\
	LIST_OPS(P, P##_list, ET, MK_PTR, UN_PTR, NO_OP, NO_OP, SHUF_##P)

#define SHUF_exit_views		OP_SHUFRNG
#define SHUF_accept_ents	NO_OP
#define SHUF_admits		NO_OP
#define SHUF_cstrings		NO_OP
#define SHUF_fights		NO_OP
#define SHUF_flag_ents		NO_OP
#define SHUF_item_ents		NO_OP
#define SHUF_orders		NO_OP
#define SHUF_req_ents		NO_OP
#define SHUF_roads		NO_OP
#define SHUF_skill_ents		NO_OP
#define SHUF_tiles		NO_OP
#define SHUF_trades		NO_OP
#define SHUF_wait_args		NO_OP

#include "modules.h"

int
main(void)
{
	char line[1024];

	setvbuf(stdout, NULL, _IOLBF, 0);

	while (fgets(line, sizeof(line), stdin)) {
		char *f[16];
		int nf = 0, i, a[16] = {0}, slot, ok = 0;
		char *p, *op;

		if ((p = strchr(line, '#')) != NULL)
			*p = '\0';
		for (p = strtok(line, " \t\r\n"); p && nf < 16; p = strtok(NULL, " \t\r\n"))
			f[nf++] = p;
		if (nf == 0)
			continue;
		for (i = 0; i < nf; i++)
			printf("%s%s", i ? " " : "", f[i]);
		printf(" -> ");
		op = f[0];

		if (!strcmp(op, "master")) {
			unsigned char b[64];
			size_t n = 0;
			FILE *fp;
			const char *h = nf > 1 ? f[1] : "";

			while (h[0] && h[1] && n < sizeof(b)) {
				unsigned v;
				sscanf(h, "%2x", &v);
				b[n++] = (unsigned char)v;
				h += 2;
			}
			fp = fopen("randseed", "wb");
			fwrite(b, 1, n, fp);
			fclose(fp);
			load_seed("randseed");
			printf("ok\n");
			continue;
		}
		if (!strcmp(op, "rnd")) {
			printf("%d\n", rnd((int)strtol(f[1], NULL, 0), (int)strtol(f[2], NULL, 0)));
			continue;
		}
		if (!strcmp(op, "root")) {
			uint32_t m[4];
			rng_master_seed(m);
			cur = rng_seed(m);
			printf("ctr %u\n", cur.counter);
			continue;
		}
		if (!strcmp(op, "of")) {
			cur = rng_stream_of(&cur, (int)strtol(f[1], NULL, 0), (uint32_t)strtoul(f[2], NULL, 0));
			printf("ctr %u\n", cur.counter);
			continue;
		}

		/* <module> <slot> <op> [args] */
		if (nf < 3 || f[1][0] < 'a' || f[1][0] >= 'a' + NSLOT || f[1][1]) {
			printf("bad line\n");
			fprintf(stderr, "bad line\n");
			return 1;
		}
		slot = f[1][0] - 'a';
		for (i = 3; i < nf; i++)
			a[i - 3] = (int)strtol(f[i], NULL, 0);
		if (!strcmp(f[2], "copy"))
			a[0] = f[3][0] - 'a';

		if (!strcmp(op, "ilist"))
			ok = do_ilist(slot, f[2], a, nf - 3);
#define M(P)	else if (!strcmp(op, #P)) ok = do_##P(slot, f[2], a, nf - 3);
		MODULE_TABLE
#undef M
		if (!ok) {
			printf("unknown op\n");
			fprintf(stderr, "unknown op %s %s\n", op, f[2]);
			return 1;
		}
	}
	return 0;
}
EOF

echo "building the C driver (log: ${WORK}/build.log)"
"${CC}" -std=c11 -w -I "${WORK}" -I "${WORK}/g3/lib" \
	"${WORK}/driver.c" "${WORK}/g3/lib/ilist.c" "${SOURCES[@]}" \
	"${WORK}/g3/lib/rnd.c" "${WORK}/g3/lib/rng.c" \
	-o "${WORK}/lists-driver" >"${WORK}/build.log" 2>&1 ||
	die "C driver build failed; see ${WORK}/build.log"

for run in run1 run2; do
	(cd "${WORK}/${run}" && "${WORK}/lists-driver" <"${CASES}" >out.txt) ||
		die "C driver failed in ${run}"
done
cmp -s "${WORK}/run1/out.txt" "${WORK}/run2/out.txt" || die "C driver is not deterministic"
echo "C driver: $(wc -l <"${WORK}/run1/out.txt" | tr -d ' ') results, deterministic"

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
(cd "${ROOT}" && go test -count=1 -run 'TestMatchesC' ./internal/lists)
echo "PASS"
