#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: the helpers check_db() (olympia/check.c) calls, C being
# the reference, against internal/olympia (Go): the names (display_name(),
# box_name(), just_name(), plural_item_name() and cap()), the here lists and
# stacks, set_lord() and the command loading it does, skill_school(),
# loc_depth(), the attribute macros and the loops of olympia/loop.h.
#
# Exports the pinned G3 revision, then builds a C driver against its
# olympia/code.c, olympia/glob.c, olympia/z.c and lib/, and copies of the
# helper functions taken from the export, each under a #line directive that
# names its own file and line, so that assert() prints them as the engine's
# build does. The routines those helpers call only on paths the cases never
# take are stubs that print "stub NAME" and abort. The driver runs
# internal/olympia/testdata/helper-cases.txt one section at a time, each in
# a fresh process, since a section may end in a failed assertion, which
# aborts.
#
# The driver runs twice, to check determinism, and its output must equal
# internal/olympia/testdata/helper-c-output.txt, which the Go test
# (TestHelperMatchesC in internal/olympia) holds the port to.
#
# Usage: tools/helper-compat/run.sh [--update]
#
#   --update  write the C output to internal/olympia/testdata/helper-c-output.txt
#             instead of comparing with it. Expected output only ever comes
#             from the C build; never edit that file by hand.
#
# Environment:
#   G3_REPO  G3 clone that contains G3_REV
#            (default ~/Software/playbymail/olympia-g3)
#   G3_REV   G3 revision to test against (default: the pinned revision)
#   WORK     scratch directory, emptied first (default tmp/helper-compat)
#   CC       C compiler (default cc)
#
# Exits 0 only if the C output matches the recorded output and the Go test
# passes.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"

G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-47535543d83fbc5649d54b63b49e578f2b11ab1c}"
WORK="${WORK:-${ROOT}/tmp/helper-compat}"
CC="${CC:-cc}"
CASES="${ROOT}/internal/olympia/testdata/helper-cases.txt"
EXPECT="${ROOT}/internal/olympia/testdata/helper-c-output.txt"

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
#   stubs.c    a stub for every command routine cmd_tbl names (glob.c), for
#              the two routines code.c calls only from display_kind() and
#              just_name_qty(), and for the routines the helpers call only on
#              paths the cases do not take
#   helpers.c  the helpers, copied from their files in source order, each
#              under a #line directive
python3 -I - "${WORK}/g3/olympia" "${WORK}/gen" <<'EOF'
import re, sys

src, out = sys.argv[1], sys.argv[2]

def func(path, name, pre=""):
    """The definition of function name in path, from its return type line
    (and pre, a pattern for what comes before it) to its closing brace,
    under a #line directive, and its offset in the file."""
    text = open(f"{src}/{path}").read()
    m = re.search(rf"^{pre}(?:static )?[^\n(]*\n{name}\(.*?^\}}\n", text, re.M | re.S)
    if not m:
        sys.exit(f"{name}() not found in {path}")
    line = text.count("\n", 0, m.start()) + 1
    return m.start(), f'#line {line} "olympia/{path}"\n{m.group(0)}\n'

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
#define STUB(name) { printf("stub %s\\n", name); abort(); }

/* Called only by display_kind() and just_name_qty(), not under test. */
int is_port_city(int where) STUB("is_port_city")
char *nice_num(int n) STUB("nice_num")

/* Called by set_lord() and load_command() only on paths not ported. */
void flush_unit_orders(int pl, int who) STUB("flush_unit_orders")
void interrupt_order(int who) STUB("interrupt_order")
void clear_all_att(int who) STUB("clear_all_att")
void reswear_all_sworn(int who) STUB("reswear_all_sworn")
int char_np_total(int who) STUB("char_np_total")
void add_np(int pl, int num) STUB("add_np")
int deduct_np(int pl, int num) STUB("deduct_np")
void out(int who, const char *format, ...) STUB("out")
void log_write(int k, const char *format, ...) STUB("log_write")
void do_command(struct command *c) STUB("do_command")
int pledge_backlinks = FALSE;
""")

# Each file's helpers, in the order they appear in it.
helpers = {
    "sout.c": ["sout"],
    "u.c": ["loc_depth", "clear_temps", "min", "cap", "test_bit", "set_bit"],
    "seed.c": ["int_comp"],
    "loc.c": ["add_char_here", "all_stack", "subloc", "viewloc", "add_to_here_list",
              "remove_from_here_list", "in_here_list", "first_character"],
    "stack.c": ["stack_parent", "stack_leader"],
    "swear.c": ["lord", "player", "set_loyal", "set_lord", "is_unit"],
    "use.c": ["skill_school"],
    "move.c": ["restore_stack_actions"],
    "day.c": ["touch_loc_pl", "touch_loc"],
    "garr.c": ["garrison_here"],
    "order.c": ["p_order_head", "rp_order_head", "top_order", "pop_order"],
    "input.c": ["parse_line", "remove_comment", "remove_ctrl_chars", "parse_arg",
                "find_command", "oly_parse_cmd", "oly_parse", "get_command",
                "set_state", "load_command", "init_load_sup"],
}
with open(f"{out}/helpers.c", "w") as f:
    f.write('#include <ctype.h>\n#include <stdarg.h>\n#include <stdio.h>\n'
            '#include <stdlib.h>\n#include <string.h>\n#include "z.h"\n#include "oly.h"\n\n'
            "#define MAX_BUFFER 50\n#define BUF_LEN LEN\n#define MAX_PRI 5\n"
            "ilist load_q[MAX_PRI];\nilist run_q;\nint fuzzy_find;\n"
            "void reswear_all_sworn(int who);\n\n")
    for path, names in helpers.items():
        for _, text in sorted(func(path, n) for n in names):
            f.write(text)
    # the static helpers the driver calls
    f.write("struct order_list *drv_p_order_head(int pl, int who) { return p_order_head(pl, who); }\n"
            "int drv_parse_arg(int who, char *s) { return parse_arg(who, s); }\n")
EOF

cat >"${WORK}/driver.c" <<'EOF'
/*
 *  Differential driver for tools/helper-compat/run.sh. Runs one section of
 *  internal/olympia/testdata/helper-cases.txt; see that file for the
 *  operations.
 *
 *  usage: driver CASES SECTION   (SECTION counts from 0)
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include "z.h"
#include "oly.h"

void glob_init(void);
void init_lower(void);

#define MAX_PRI 5
extern ilist load_q[MAX_PRI];
extern ilist run_q;
struct order_list *drv_p_order_head(int pl, int who);
void set_loyal(int who, int k, int lev);
void touch_loc_pl(int pl, int where);
int garrison_here(int where);
int drv_parse_arg(int who, char *s);
int is_unit(int pl, int v);

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
	static char buf[4][16384];
	static int k;
	char *q;

	if (s == NULL)
		return "(null)";
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
list(const char *label, ilist l)
{
	int i;

	printf("%s:", label);
	for (i = 0; i < ilist_len(l); i++)
		printf(" %d", l[i]);
	printf("\n");
}

/* everything about box n that the helpers read or change */
static void
dump(int n)
{
	struct box *b = bx[n];
	int i, j;

	if (b == NULL)
	{
		printf("box %d: none\n", n);
		return;
	}
	printf("box %d: kind %d sub %d where %d temp %d\n", n, b->kind, b->skind,
	       b->x_loc_info.where, b->temp);
	list("  here", b->x_loc_info.here_list);
	if (b->x_char)
	{
		struct entity_char *p = b->x_char;

		printf("  char: lord %d prev %d loy %d/%d moving %d prisoner %d item %d\n",
		       p->unit_lord, p->prev_lord, p->loy_kind, p->loy_rate,
		       p->moving, p->prisoner, p->unit_item);
	}
	if (b->x_player)
	{
		struct entity_player *p = b->x_player;

		printf("  player: last_order_turn %d\n", p->last_order_turn);
		list("  units", p->units);
		list("  locs", p->locs);
		for (i = 0; i < orders_len(p->orders); i++)
		{
			printf("  orders %d:", p->orders[i]->unit);
			for (j = 0; j < cstrings_len(p->orders[i]->l); j++)
				printf(" [%s]", esc(p->orders[i]->l[j]));
			printf("\n");
		}
	}
	if (b->cmd)
	{
		struct command *c = b->cmd;

		printf("  cmd: who %d state %d pri %d wait %d poll %d days %d cmd %d cond %d fuzzy %d\n",
		       c->who, c->state, c->pri, c->wait, c->poll,
		       c->days_executing, c->cmd, c->conditional, c->fuzzy);
		printf("  args: %d %d %d %d %d %d %d %d line [%s]\n",
		       c->a, c->b, c->c, c->d, c->e, c->f, c->g, c->h, esc(c->line));
	}
}

static void
loop_op(void)
{
	char what[16];
	int nread, i, n;
	struct item_ent *e;

	what[0] = '\0';
	sscanf(pos, "%15s%n", what, &nread);
	pos += nread;
	n = num();

	printf("visit:");
	if (strcmp(what, "here") == 0)
	{
		loop_here(n, i)
		{
			printf(" %d", i);
		}
		next_here;
	}
	else if (strcmp(what, "stack") == 0)
	{
		loop_stack(n, i)
		{
			printf(" %d", i);
		}
		next_stack;
	}
	else if (strcmp(what, "inv") == 0)
	{
		loop_inv(n, e)
		{
			printf(" %d/%d", e->item, e->qty);
		}
		next_inv;
	}
	else if (strcmp(what, "kind") == 0)
	{
		loop_kind(n, i)
		{
			printf(" %d", i);
		}
		next_kind;
	}
	else if (strcmp(what, "units") == 0)
	{
		loop_units(n, i)
		{
			printf(" %d", i);
		}
		next_unit;
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

	init_lower();
	glob_init();

	while (fgets(line, sizeof(line), fp) != NULL)
	{
		char op[16];
		int nread, n, k, sk, l;

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

		/* setup */
		if (strcmp(op, "alloc") == 0)
		{
			n = num();
			k = num();
			sk = num();
			alloc_box(n, k, sk);
		}
		else if (strcmp(op, "delete") == 0)
			delete_box(num());
		else if (strcmp(op, "setname") == 0)
		{
			n = num();
			set_name(n, text());
		}
		else if (strcmp(op, "where") == 0)
		{
			n = num();
			p_loc_info(n)->where = num();
		}
		else if (strcmp(op, "here") == 0)
		{
			n = num();
			ilist_append(&p_loc_info(n)->here_list, num());
		}
		else if (strcmp(op, "lord") == 0)
		{
			n = num();
			p_char(n)->unit_lord = num();
		}
		else if (strcmp(op, "prevlord") == 0)
		{
			n = num();
			p_char(n)->prev_lord = num();
		}
		else if (strcmp(op, "units") == 0)
		{
			n = num();
			ilist_append(&p_player(n)->units, num());
		}
		else if (strcmp(op, "unititem") == 0)
		{
			n = num();
			p_char(n)->unit_item = num();
		}
		else if (strcmp(op, "plural") == 0)
		{
			n = num();
			p_item(n)->plural_name = str_save(text());
		}
		else if (strcmp(op, "whohas") == 0)
		{
			n = num();
			p_item(n)->who_has = num();
		}
		else if (strcmp(op, "req") == 0)
		{
			n = num();
			p_skill(n)->required_skill = num();
		}
		else if (strcmp(op, "learn") == 0)
		{
			n = num();
			p_skill(n)->time_to_learn = num();
		}
		else if (strcmp(op, "moving") == 0)
		{
			n = num();
			p_char(n)->moving = num();
		}
		else if (strcmp(op, "prisoner") == 0)
		{
			n = num();
			p_char(n)->prisoner = (schar) num();
		}
		else if (strcmp(op, "cmd") == 0)
		{
			struct command *c;

			n = num();
			c = p_command(n);
			c->who = num();
			c->state = (schar) num();
			c->pri = (schar) num();
		}
		else if (strcmp(op, "order") == 0)
		{
			n = num();
			k = num();
			cstrings_append(&drv_p_order_head(n, k)->l, str_save(text()));
		}
		else if (strcmp(op, "turn") == 0)
			sysclock.turn = (short) num();
		else if (strcmp(op, "temp") == 0)
		{
			n = num();
			bx[n]->temp = num();
		}
		else if (strcmp(op, "inv") == 0)
		{
			struct item_ent *e = my_malloc(sizeof(*e));

			n = num();
			e->item = num();
			e->qty = num();
			item_ents_append(&bx[n]->items, e);
		}
		/* queries and operations under test */
		else if (strcmp(op, "names") == 0)
		{
			n = num();
			printf("display_name [%s]\n", esc(display_name(n)));
			printf("box_name [%s]\n", esc(box_name(n)));
			printf("just_name [%s]\n", esc(just_name(n)));
		}
		else if (strcmp(op, "pluralname") == 0)
		{
			n = num();
			printf("plural_item_name [%s]\n", esc(plural_item_name(n, num())));
		}
		else if (strcmp(op, "cap") == 0)
			printf("cap [%s]\n", esc(cap(text())));
		else if (strcmp(op, "depth") == 0)
			printf("loc_depth %d\n", loc_depth(num()));
		else if (strcmp(op, "inhere") == 0)
		{
			n = num();
			printf("in_here_list %d\n", in_here_list(n, num()));
		}
		else if (strcmp(op, "addhere") == 0)
		{
			n = num();
			add_to_here_list(n, num());
		}
		else if (strcmp(op, "remhere") == 0)
		{
			n = num();
			remove_from_here_list(n, num());
		}
		else if (strcmp(op, "parent") == 0)
			printf("stack_parent %d\n", stack_parent(num()));
		else if (strcmp(op, "leader") == 0)
			printf("stack_leader %d\n", stack_leader(num()));
		else if (strcmp(op, "school") == 0)
			printf("skill_school %d\n", skill_school(num()));
		else if (strcmp(op, "isunit") == 0)
		{
			n = num();
			printf("is_unit %d\n", is_unit(n, num()));
		}
		else if (strcmp(op, "player") == 0)
			printf("player %d\n", player(num()));
		else if (strcmp(op, "subloc") == 0)
			printf("subloc %d\n", subloc(num()));
		else if (strcmp(op, "viewloc") == 0)
			printf("viewloc %d\n", viewloc(num()));
		else if (strcmp(op, "firstchar") == 0)
			printf("first_character %d\n", first_character(num()));
		else if (strcmp(op, "garrison") == 0)
			printf("garrison_here %d\n", garrison_here(num()));
		else if (strcmp(op, "parsearg") == 0)
		{
			n = num();
			printf("parse_arg %d\n", drv_parse_arg(n, text()));
		}
		else if (strcmp(op, "setlord") == 0)
		{
			n = num();
			l = num();
			k = num();
			set_lord(n, l, k, num());
		}
		else if (strcmp(op, "restore") == 0)
			restore_stack_actions(num());
		else if (strcmp(op, "touch") == 0)
			touch_loc(num());
		else if (strcmp(op, "initload") == 0)
			init_load_sup(num());
		else if (strcmp(op, "clear") == 0)
			clear_temps(num());
		else if (strcmp(op, "macros") == 0)
		{
			n = num();
			printf("noble_item %d char_moving %d is_prisoner %d item_unique %d "
			       "learn_time %d req_skill %d\n",
			       noble_item(n), char_moving(n), is_prisoner(n),
			       item_unique(n), learn_time(n), req_skill(n));
		}
		else if (strcmp(op, "sort") == 0)
		{
			ilist v = NULL;
			char *end;

			for (;;)
			{
				n = (int) strtol(pos, &end, 0);
				if (end == pos)
					break;
				pos = end;
				ilist_append(&v, n);
			}
			safe_qsort(v, (size_t)ilist_len(v), sizeof(int), int_comp);
			list("sorted", v);
		}
		else if (strcmp(op, "loop") == 0)
			loop_op();
		else if (strcmp(op, "dump") == 0)
			dump(num());
		else if (strcmp(op, "queues") == 0)
		{
			char label[16];

			for (k = 0; k < MAX_PRI; k++)
			{
				sprintf(label, "load_q %d", k);
				list(label, load_q[k]);
			}
			list("run_q", run_q);
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
	"${WORK}/driver.c" "${WORK}/gen/helpers.c" olympia/code.c olympia/glob.c olympia/z.c \
	"${WORK}/gen/stubs.c" lib/checked_alloc.c lib/ilist.c lib/rng.c lib/rnd.c \
	lib/cstrings.c lib/item_ents.c lib/orders.c \
	-o "${WORK}/helper-driver") >"${WORK}/build.log" 2>&1 ||
	die "C driver build failed; see ${WORK}/build.log"

sections=$(grep -c '^section ' "${CASES}")
# The driver sends its stderr to its stdout. Bash reports each abort on its
# own stderr, which is dropped.
for run in run1 run2; do
	for ((i = 0; i < sections; i++)); do
		status=0
		(cd "${WORK}/${run}" && exec "${WORK}/helper-driver" "${CASES}" "${i}" >>out.txt) || status=$?
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
(cd "${ROOT}" && go test -count=1 -run 'TestHelperMatchesC' ./internal/olympia)
echo "PASS"
