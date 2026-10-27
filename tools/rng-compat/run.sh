#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: G3's addressable RNG layer (C, the reference) against
# internal/rng and internal/rnd (Go).
#
# Exports the pinned G3 revision, then builds a small C driver against its
# lib/rng.c, lib/rnd.c and lib/ilist.c. The driver also compiles the named
# stream code (the begin_*() seeders and their leaf helpers) cut verbatim from
# the pinned olympia/*.c files, at the line ranges listed in SNIPPETS below.
# The driver interprets internal/rng/testdata/cases.txt and prints one result
# line per operation. It runs twice, to check determinism, and its output must
# equal internal/rng/testdata/c-output.txt, which the Go test
# (TestMatchesC in internal/rng) holds the port to.
#
# Usage: tools/rng-compat/run.sh [--update]
#
#   --update  write the C output to internal/rng/testdata/c-output.txt
#             instead of comparing with it. Expected output only ever comes
#             from the C build; never edit that file by hand.
#
# Environment:
#   G3_REPO  G3 clone that contains G3_REV
#            (default ~/Software/playbymail/olympia-g3)
#   G3_REV   G3 revision to test against (default: the pinned revision)
#   WORK     scratch directory, emptied first (default tmp/rng-compat)
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
G3_REV="${G3_REV:-6a25f2d252ce45bd5adbe3baaccb3be3eed56279}"
WORK="${WORK:-${ROOT}/tmp/rng-compat}"
CC="${CC:-cc}"
CASES="${ROOT}/internal/rng/testdata/cases.txt"
EXPECT="${ROOT}/internal/rng/testdata/c-output.txt"

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

git -C "${G3_REPO}" archive "${G3_REV}" lib olympia | tar -x -C "${WORK}/g3"

# ---------------------------------------------------------------------------
# The named streams, cut from the pinned engine sources. Each range runs from
# the file's first stream declaration or tag through the closing brace of its
# last helper. They hold only tags, the static streams and the helpers, which
# need nothing from the engine but sysclock, the rng layer and ilist.
# ---------------------------------------------------------------------------
SNIPPETS=(
	buy.c:34:34 buy.c:215:312
	combat.c:32:90
	npc.c:24:62
	quest.c:32:56
	cloud.c:25:53
	faery.c:34:158
	hades.c:37:145
	tunnel.c:26:54
	basic.c:66:221
	c1.c:46:138
	code.c:729:796
	dir.c:21:42
	day.c:50:155 day.c:189:227
	seed.c:30:91
	storm.c:40:104
	swear.c:428:528
	use.c:1623:1740
	u.c:41:181
)
{
	for s in "${SNIPPETS[@]}"; do
		IFS=: read -r file first last <<<"${s}"
		printf '\n/* ---- olympia/%s lines %s-%s ---- */\n' "${file}" "${first}" "${last}"
		sed -n "${first},${last}p" "${WORK}/g3/olympia/${file}"
	done
} >"${WORK}/snippets.c"

cat >"${WORK}/driver.c" <<'EOF'
/*
 *  Differential driver for tools/rng-compat/run.sh. Reads operations from
 *  stdin (internal/rng/testdata/cases.txt) and prints "<op> -> <result>" for
 *  each. See that file for the operations.
 */
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "rnd.h"
#include "rng.h"
#include "lists.h"

typedef struct {
	short day;
	short turn;
	int days_since_epoch;
} olytime;

olytime sysclock;

extern void load_seed(char *fnam);
extern int rnd(int low, int high);
extern int md5_int(int a, int b, int c, int d);

void *my_malloc(unsigned size) { return calloc(1, size); }
void *my_realloc(void *p, unsigned size) { return realloc(p, size); }

#include "snippets.c"

static void
print_stream(const rng_stream *s)
{
	printf("seed %08x %08x %08x %08x ctr %u\n", s->seed[0], s->seed[1],
	       s->seed[2], s->seed[3], s->counter);
}

static void
print_master(void)
{
	uint32_t m[4];

	rng_master_seed(m);
	printf("master %08x %08x %08x %08x\n", m[0], m[1], m[2], m[3]);
}

static ilist
make_list(int n)
{
	ilist l = NULL;
	int i;

	for (i = 0; i < n; i++)
		ilist_append(&l, 100 + i);
	return l;
}

static void
print_list(ilist l)
{
	int i;

	for (i = 0; i < ilist_len(l); i++)
		printf("%s%d", i ? " " : "", l[i]);
	printf("\n");
	ilist_reclaim(&l);
}

#define MAXF 16

int
main(void)
{
	char line[1024];
	rng_stream cur = {{0, 0, 0, 0}, 0};

	setvbuf(stdout, NULL, _IOLBF, 0);

	while (fgets(line, sizeof(line), stdin)) {
		char *f[MAXF];
		int nf = 0, i;
		int a[MAXF];
		char *p, *op;

		if ((p = strchr(line, '#')) != NULL)
			*p = '\0';
		for (p = strtok(line, " \t\r\n"); p && nf < MAXF; p = strtok(NULL, " \t\r\n"))
			f[nf++] = p;
		if (nf == 0)
			continue;
		for (i = 0; i < nf; i++)
			printf("%s%s", i ? " " : "", f[i]);
		printf(" -> ");
		for (i = 1; i < nf; i++)
			a[i - 1] = (int)strtol(f[i], NULL, 0);
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
			print_master();
		} else if (!strcmp(op, "loadmissing")) {
			load_seed("missing-randseed");
			print_master();
		} else if (!strcmp(op, "rnd")) {
			printf("%d\n", rnd(a[0], a[1]));
		} else if (!strcmp(op, "md5int")) {
			printf("%d\n", md5_int(a[0], a[1], a[2], a[3]));
		} else if (!strcmp(op, "clock")) {
			sysclock.turn = (short)a[0];
			sysclock.day = (short)a[1];
			printf("ok\n");
		} else if (!strcmp(op, "root")) {
			uint32_t m[4];
			rng_master_seed(m);
			cur = rng_seed(m);
			print_stream(&cur);
		} else if (!strcmp(op, "of")) {
			cur = rng_stream_of(&cur, a[0], (uint32_t)strtoul(f[2], NULL, 0));
			print_stream(&cur);
		} else if (!strcmp(op, "draw")) {
			printf("%d ", rng_draw(&cur, a[0], a[1]));
			print_stream(&cur);
		} else if (!strcmp(op, "keyed")) {
			printf("%d\n", rng_keyed(&cur, a[0], a[1],
			    (uint32_t)strtoul(f[3], NULL, 0), a[3], a[4]));
		} else if (!strcmp(op, "shuffle")) {
			ilist l = make_list(a[0]);
			ilist_shuffle_rng(l, &cur);
			printf("ctr %u: ", cur.counter);
			print_list(l);
		} else if (!strcmp(op, "call")) {
			const char *fn = f[1];
			int *x = a + 1;

#define V(name, ...)	if (!strcmp(fn, #name)) { name(__VA_ARGS__); printf("ok\n"); continue; }
#define I(name, ...)	if (!strcmp(fn, #name)) { printf("%d\n", name(__VA_ARGS__)); continue; }
#define L(name, ...)	if (!strcmp(fn, #name)) { ilist l = make_list(x[0]); name(__VA_ARGS__); print_list(l); continue; }

			V(begin_battle, x[0], x[1], x[2])
			I(crnd, x[0], x[1])
			V(begin_pillage, x[0], x[1])
			I(prnd, x[0], x[1])
			/* create_peasant_mob()'s name draw on &pillage_rng */
			if (!strcmp(fn, "pillage_mob")) { printf("%d\n", rng_draw(&pillage_rng, 1, 2)); continue; }

			V(begin_economy, x[0])
			I(econ_pick, x[0], x[1], x[2], x[3])
			I(econ_stock, x[0], x[1], x[2], x[3])
			I(econ_qty, x[0], x[1], x[2], x[3])
			I(econ_cost, x[0], x[1], x[2], x[3])
			I(econ_expire, x[0], x[1], x[2], x[3])
			I(econ_mine, x[0], x[1], x[2], x[3])
			I(econ_harvest, x[0], x[1], x[2], x[3])
			I(econ_ring, x[0], x[1], x[2], x[3])

			V(begin_npc, x[0])
			I(npc_spawn, x[0], x[1], x[2], x[3])
			I(npc_qty, x[0], x[1], x[2], x[3])
			I(npc_behavior, x[0], x[1], x[2], x[3])
			/* auto_bandit's inline victim pick (npc.c:394) */
			if (!strcmp(fn, "npc_victim")) { printf("%d\n", rng_keyed(&npc_rng, x[0], 0, NTAG_VICTIM, 0, x[1] - 1)); continue; }

			V(begin_quest, x[0], x[1])
			I(qrnd, x[0], x[1])
			if (!strcmp(fn, "quest_shuffle")) { ilist l = make_list(x[0]); ilist_shuffle_rng(l, &quest_rng); print_list(l); continue; }

			V(begin_cloud_build, x[0])
			I(cseq_rnd, x[0], x[1])
			L(cseq_shuffle, l)
			V(begin_faery_build, x[0])
			I(fseq_rnd, x[0], x[1])
			L(fseq_shuffle, l)
			V(begin_hades_build, x[0])
			I(hseq_rnd, x[0], x[1])
			L(hseq_shuffle, l)
			V(begin_worldgen_loc, x[0])
			I(wseq_rnd, x[0], x[1])
			L(wseq_shuffle, l)

			I(magc_scry, x[0], x[1])
			I(magc_piety, x[0], x[1])
			I(magc_eat, x[0], x[1], x[2])
			I(magc_learn, x[0], x[1])
			I(magc_potion, x[0], x[1], x[2], x[3], x[4])
			I(magc_med, x[0], x[1])
			I(magc_omen, x[0], x[1])
			I(magc_heal, x[0], x[1])
			I(magc_forge, x[0], x[1])
			I(magc_orb, x[0], x[1])
			I(magc_ring, x[0], x[1])

			I(expl_find, x[0], x[1])
			I(expl_gate, x[0], x[1])
			I(expl_flavor, x[0], x[1])
			I(expl_pick, x[0], x[1], x[2])
			I(expl_seek, x[0], x[1])
			I(expl_detect, x[0], x[1])

			I(mint_alloc, x[0], x[1])
			I(mint_password, x[0], x[1], x[2])
			I(mint_city)
			L(mint_shuffle, l)

			/* dir.c:797-798: begin_exitdir(), then a shuffle on &exitdir_rng */
			if (!strcmp(fn, "exitdir_shuffle")) { ilist l = make_list(x[0]); begin_exitdir(); ilist_shuffle_rng(l, &exitdir_rng); print_list(l); continue; }

			I(up_heal, x[0], x[1], x[2], x[3])
			I(up_loyal, x[0], x[1], x[2], x[3])
			I(up_starve, x[0], x[1], x[2], x[3])
			I(up_animal, x[0], x[1], x[2], x[3], x[4])
			I(up_corpse, x[0], x[1], x[2])
			I(up_income, x[0], x[1], x[2], x[3])
			I(cal_day, x[0], x[1], x[2])

			I(faer_hunt_loc, x[0], x[1], x[2], x[3])
			I(faer_hunt_qty, x[0])
			I(faer_bandit_kind, x[0], x[1])
			I(faer_bandit_qty, x[0])
			I(faer_bandit_gold, x[0])
			I(faer_ambush, x[0], x[1])
			I(faer_retal, x[0], x[1])

			I(hads_nasty, x[0], x[1], x[2], x[3])
			I(hads_bandit_kind, x[0], x[1])
			I(hads_bandit_qty, x[0])
			I(hads_ambush, x[0], x[1])
			I(hads_retal, x[0], x[1])
			/* random_hades_loc's inline pick (hades.c:550-551) */
			if (!strcmp(fn, "hads_random_loc")) { begin_hades(); printf("%d\n", rng_keyed(&hades_rng, x[0], 0, HTAG_RLOC, 0, x[1] - 1)); continue; }

			I(wgen_prom, x[0])
			I(wgen_skill, x[0], x[1], x[2], x[3])
			I(wgen_garr, x[0])
			I(wgen_gate, x[0])

			L(wthr_shuffle, l)
			I(wthr_storm, x[0], x[1], x[2], x[3])
			I(wthr_wreck, x[0], x[1], x[2], x[3])

			I(soc_gift, x[0], x[1])
			I(soc_bribe, x[0], x[1], x[2], x[3], x[4])
			I(soc_terror, x[0], x[1])
			I(soc_incite, x[0], x[1], x[2], x[3], x[4])
			I(soc_persuade, x[0], x[1])
			I(soc_breed, x[0], x[1], x[2], x[3], x[4])

			I(skil_crit, x[0], x[1])
			I(skil_bonus, x[0], x[1], x[2], x[3])
			I(skil_study, x[0], x[1])
			I(skil_research, x[0], x[1])
			I(skil_research_pick, x[0], x[1], x[2])
			I(skil_torture, x[0], x[1])
			I(skil_petty, x[0], x[1], x[2], x[3], x[4])

			I(ent_prisoner, x[0], x[1])
			I(ent_drop, x[0], x[1], x[2], x[3], x[4])
			I(ent_take, x[0], x[1], x[2], x[3], x[4])
			I(ent_sick, x[0], x[1])
			I(ent_land_dir, x[0], x[1])
			I(ent_land_pick, x[0], x[1])
			I(ent_bark, x[0], x[1])
			I(ent_build, x[0])
			I(ent_menial, x[0])

			printf("unknown helper\n");
			fprintf(stderr, "unknown helper %s\n", fn);
			return 1;
		} else {
			printf("unknown op\n");
			fprintf(stderr, "unknown op %s\n", op);
			return 1;
		}
	}
	return 0;
}
EOF

echo "building the C driver (log: ${WORK}/build.log)"
"${CC}" -std=c11 -w -I "${WORK}" -I "${WORK}/g3/lib" \
	"${WORK}/driver.c" "${WORK}/g3/lib/rng.c" "${WORK}/g3/lib/rnd.c" "${WORK}/g3/lib/ilist.c" \
	-o "${WORK}/rng-driver" >"${WORK}/build.log" 2>&1 ||
	die "C driver build failed; see ${WORK}/build.log"

for run in run1 run2; do
	(cd "${WORK}/${run}" && run_limited "${RUN_TIMEOUT}" rng-driver "${WORK}/rng-driver" <"${CASES}" >out.txt) ||
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
(cd "${ROOT}" && RUN_TIMEOUT="${RUN_TIMEOUT}" go test -exec "bash ${HERE}/../lib/run-limited.sh" -count=1 -run 'TestMatchesC' ./internal/rng)
echo "PASS"
