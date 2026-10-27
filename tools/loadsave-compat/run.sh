#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Differential test: the routines of G3's load_db() and save_db() that
# read or write files of their own or set up globals (#50), C being the
# reference, against internal/olympia (Go): load_system(), save_system()
# and cleanup_posts() (olympia/io.c), load_orders() and save_orders()
# (olympia/order.c), init_trade_routes() and save_trade_routes()
# (olympia/buy.c), determine_map_edges() (olympia/dir.c) and
# rename_act_join_files() (olympia/add.c) (#50); and create_relics()
# (olympia/quest.c), init_ocean_chars() (olympia/move.c) and
# delete_deadchars() (olympia/io.c) (#51); around the box-file passes and
# check_db() that tools/check-compat checks.
#
# Exports the pinned G3 revision, then builds a C driver as
# tools/check-compat does, against olympia/check.c, olympia/code.c,
# olympia/glob.c, olympia/z.c and lib/, with copies of the helpers check_db()
# calls and of the routines above and what they call (each under a #line
# directive naming its own file and line), io.c from its first line up to
# delete_deadchars() (so through load_system() and save_system()), its
# delete_deadchars(), cleanup_posts() and remove_dir(), and buy.c's
# trade-route block, which the driver includes so that it can print
# buy.c's statics.
#
# It runs two checks:
#
# 1. internal/olympia/testdata/loadsave-cases.txt, one section at a time,
#    each in a fresh process and a fresh, empty directory, since a section
#    may end in a failed assertion or a crash. Each section writes a small
#    game, runs the routines, prints the globals they set and lists the
#    directory's tree (every file's path, mode and contents). The driver
#    runs twice, to check determinism, and its output must equal
#    internal/olympia/testdata/loadsave-c-output.txt, which the Go test
#    (TestLoadSaveMatchesC in internal/olympia) holds the port to.
#
# 2. Real saved games, only with --full or --games: for each game's lib
#    directory, copied, the driver runs what load_db() runs of these
#    routines, in its order (init_trade_routes(), load_system(), the two
#    box-file passes, load_orders(), check_db(), determine_map_edges(),
#    create_relics(), init_ocean_chars() and delete_deadchars()),
#    prints the globals to stdout, and runs what save_db() runs of them
#    (cleanup_posts(), save_trade_routes(), save_system(),
#    write_all_boxes(), write_master(), save_orders() and
#    rename_act_join_files()); the seed file is left to #52.
#    TestLoadSaveWorlds does the same in Go on its own copy, and the lib
#    trees (which files exist, their modes and contents) and stdout and
#    stderr must be equal, byte for byte. The games are world.tgz from the
#    export and, if tools/olympia-roundtrip has been run, every input it
#    has cached under tmp/olympia-roundtrip/cases/*/input/lib. Each game's
#    copies and output are deleted as soon as they match; only failing
#    games' are kept.
#
# Usage: tools/loadsave-compat/run.sh [--update] [--full | --games N]
#
#   --update   write the C output of check 1 to
#              internal/olympia/testdata/loadsave-c-output.txt instead of
#              comparing with it. Expected output only ever comes from the
#              C build; never edit that file by hand.
#   --full     run check 2 on every game
#   --games N  run check 2 on N games, always the same ones: world.tgz and
#              cached games spread evenly over the sorted list
#
# Environment:
#   G3_REPO  G3 clone that contains G3_REV
#            (default ~/Software/playbymail/olympia-g3)
#   G3_REV   G3 revision to test against (default: the pinned revision)
#   WORK     scratch directory, emptied first (default tmp/loadsave-compat)
#   GAMES    directory of cached games (default tmp/olympia-roundtrip/cases)
#   CC       C compiler (default cc)
#   RUN_TIMEOUT seconds allowed per driver/test process (default 10)
#
# Exits 0 only if the C output matches the recorded output, the Go test
# passes and every saved game checked gives the same result in C and Go.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"
# shellcheck source=../lib/timeout.sh
. "${HERE}/../lib/timeout.sh"

G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-6a25f2d252ce45bd5adbe3baaccb3be3eed56279}"
WORK="${WORK:-${ROOT}/tmp/loadsave-compat}"
GAMES="${GAMES:-${ROOT}/tmp/olympia-roundtrip/cases}"
CC="${CC:-cc}"
RUN_TIMEOUT="${RUN_TIMEOUT:-10}"
CASES="${ROOT}/internal/olympia/testdata/loadsave-cases.txt"
EXPECT="${ROOT}/internal/olympia/testdata/loadsave-c-output.txt"

usage() { echo "usage: $0 [--update] [--full | --games N]" >&2; exit 2; }
UPDATE=0
NGAMES=0 # games for check 2: 0 none, -1 all
while [ $# -gt 0 ]; do
	case "$1" in
	--update) UPDATE=1 ;;
	--full) NGAMES=-1 ;;
	--games)
		[ $# -ge 2 ] && [[ "$2" =~ ^[1-9][0-9]*$ ]] || usage
		NGAMES="$2"
		shift
		;;
	*) usage ;;
	esac
	shift
done

die() { echo "error: $*" >&2; exit 2; }

for tool in git go perl python3 tar cmp diff find stat sort "${CC}"; do
	command -v "${tool}" >/dev/null || die "missing required tool: ${tool}"
done
git -C "${G3_REPO}" cat-file -e "${G3_REV}^{commit}" 2>/dev/null ||
	die "${G3_REPO} does not contain ${G3_REV}; clone https://github.com/playbymail/olympia-g3 and set G3_REPO"

# Sections may leave unreadable directories behind.
[ -d "${WORK}" ] && chmod -R u+rwx "${WORK}" 2>/dev/null || true
rm -rf "${WORK}"
mkdir -p "${WORK}/g3" "${WORK}/gen"

echo "G3 revision: ${G3_REV}"
echo "G3 source:   ${G3_REPO}"
echo "work dir:    ${WORK}"
echo

git -C "${G3_REPO}" archive "${G3_REV}" lib olympia tests/olympia/fixtures/world.tgz | tar -x -C "${WORK}/g3"

# Generate, from the exported sources:
#   stubs.c      a stub for every command routine cmd_tbl names (glob.c),
#                for the two routines code.c calls only from display_kind()
#                and just_name_qty(), and for the routines the helpers call
#                only on paths not ported (among them test_known() and
#                queue_lore(), which add_item() calls only for an item
#                with lore), each of which prints "stub NAME" and aborts
#   helpers.c    the helpers check_db() calls, as tools/helper-compat copies
#                them, lookup() and stage() from u.c, which io.c and
#                check_db() call, and the routines this check covers and
#                what they call, each under a #line directive; and
#                the globals of other files they use: max_map_row,
#                max_map_col and max_map_init (dir.c), new_players
#                (add.c), faery_region and faery_player (faery.c),
#                hades_region, hades_pit and hades_player (hades.c),
#                nowhere_region and nowhere_loc (quest.c), tunnel_region
#                and under_region (tunnel.c), cloud_region (cloud.c),
#                ocean_chars (move.c, where it is static) and
#                trades_to_check (buy.c), which print() lists
#   io_files.c   io.c up to delete_deadchars(), just past save_system(),
#                and delete_deadchars(), cleanup_posts() and remove_dir()
#                from further down
#   buy_trade.c  buy.c's trade-route block, from TRADE_ROUTE_DEFAULT_SECRET
#                through save_trade_routes()
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

/* Called only by display_kind() and just_name_qty(), not reached. */
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

/* Called by add_item() only for an item with lore, not reached. */
int test_known(int who, int i) STUB("test_known")
void queue_lore(int who, int num, int anyway) STUB("queue_lore")
int pledge_backlinks = FALSE;
""")

# Each file's helpers, in the order they appear in it.
helpers = {
    "sout.c": ["sout"],
    "u.c": ["loc_depth", "clear_temps", "min", "cap", "test_bit", "set_bit",
            "lookup", "stage", "add_item", "create_unique_item_alloc"],
    "seed.c": ["int_comp"],
    "loc.c": ["add_char_here", "all_stack", "region", "subloc", "viewloc",
              "add_to_here_list", "remove_from_here_list", "in_here_list",
              "set_where", "first_character"],
    "stack.c": ["stack_parent", "stack_leader"],
    "swear.c": ["lord", "player", "set_loyal", "set_lord", "is_unit"],
    "use.c": ["skill_school"],
    "move.c": ["restore_stack_actions", "init_ocean_chars"],
    "day.c": ["touch_loc_pl", "touch_loc"],
    "garr.c": ["garrison_here"],
    "order.c": ["p_order_head", "rp_order_head", "top_order", "pop_order",
                "queue_order", "save_player_orders", "load_player_orders",
                "load_orders", "save_orders"],
    "input.c": ["parse_line", "remove_comment", "remove_ctrl_chars", "parse_arg",
                "find_command", "oly_parse_cmd", "oly_parse", "get_command",
                "set_state", "load_command", "init_load_sup"],
    "dir.c": ["region_row", "region_col", "determine_map_edges"],
    "add.c": ["rename_act_join_files"],
    "buy.c": ["investigate_possible_trade"],
    "quest.c": ["create_a_relic", "create_relics"],
}
with open(f"{out}/helpers.c", "w") as f:
    f.write('#include <ctype.h>\n#include <dirent.h>\n#include <stdarg.h>\n'
            '#include <stdio.h>\n#include <stdlib.h>\n#include <string.h>\n'
            '#include <sys/stat.h>\n#include <time.h>\n#include <unistd.h>\n'
            '#include "z.h"\n#include "oly.h"\n\n'
            "#define MAX_BUFFER 50\n#define BUF_LEN LEN\n#define MAX_PRI 5\n"
            "ilist load_q[MAX_PRI];\nilist run_q;\nint fuzzy_find;\n"
            "int max_map_row = 0;\nint max_map_col = 0;\nint max_map_init = FALSE;\n"
            "ilist new_players = NULL;\n"
            "int faery_region = 0;\nint faery_player = 0;\n"
            "int hades_region = 0;\nint hades_pit = 0;\nint hades_player = 0;\n"
            "int nowhere_region = 0;\nint nowhere_loc = 0;\n"
            "int tunnel_region = 0;\nint under_region = 0;\nint cloud_region = 0;\n"
            "ilist ocean_chars = NULL;\nilist trades_to_check = NULL;\n"
            "extern int time_self;\n"
            "void reswear_all_sworn(int who);\n"
            "int test_known(int who, int i);\n"
            "void queue_lore(int who, int num, int anyway);\n\n")
    for path, names in helpers.items():
        for _, text in sorted(func(path, n) for n in names):
            f.write(text)

io = open(f"{src}/io.c").read()
end = io.find("static void\ndelete_deadchars(")
if end < 0 or io.find("\nsave_system(void)\n") > end:
    sys.exit("delete_deadchars() not found after save_system() in io.c")
m = re.search(r"^/\*\n(?: \*[^\n]*\n)* \*/\nvoid\nremove_dir\(char \*dir\)\n\{.*?^\}\n", io, re.M | re.S)
if not m:
    sys.exit("remove_dir() not found in io.c")
with open(f"{out}/io_files.c", "w") as f:
    # so that assert() names io.c and its lines, as the engine's build does
    f.write('#line 1 "olympia/io.c"\n')
    f.write(io[:end])
    # delete_deadchars(), which load_db() runs last, and cleanup_posts(),
    # which save_db() runs, from further down io.c
    f.write("\n" + func("io.c", "delete_deadchars")[1])
    f.write("\n" + func("io.c", "cleanup_posts")[1])
    # remove_dir(), which write_all_boxes() and save_orders() call
    f.write(f'\n#line {io[:m.start()].count(chr(10)) + 1} "olympia/io.c"\n')
    f.write(m.group(0))

buy = open(f"{src}/buy.c").read()
m = re.search(r"^#define\tTRADE_ROUTE_DEFAULT_SECRET.*?^save_trade_routes\(void\)\n\{.*?^\}\n", buy, re.M | re.S)
if not m:
    sys.exit("the trade-route block not found in buy.c")
with open(f"{out}/buy_trade.c", "w") as f:
    f.write(f'#line {buy[:m.start()].count(chr(10)) + 1} "olympia/buy.c"\n')
    f.write(m.group(0))
EOF

cat >"${WORK}/driver.c" <<'EOF'
/*
 *  Differential driver for tools/loadsave-compat/run.sh.
 *
 *  usage: driver CASES SECTION   runs one section of loadsave-cases.txt
 *                                (SECTION counts from 0); see that file
 *         driver -g LIBDIR       what load_db() runs of the routines
 *                                checked, print, and what save_db() runs
 *                                of them, on LIBDIR
 */
#include <dirent.h>
#include <limits.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

/* io.c's top includes z.h and oly.h, which have no include guards */
#include "io_files.c"
#include "buy_trade.c"

void glob_init(void);
void init_lower(void);
void check_db(void);
void load_orders(void);
void save_orders(void);
void determine_map_edges(void);
void rename_act_join_files(void);
void create_relics(void);
void init_ocean_chars(void);

extern int max_map_row, max_map_col, max_map_init;
extern ilist ocean_chars, trades_to_check;

int pretty_data_files = FALSE;
int win_flag = FALSE;
int time_self = FALSE;

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

static char *pos;

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

/* the box-file passes, as load_db() runs them */
static void
boxes(void)
{
	if (!fast_scan())	/* pass 1: call alloc_box for each entity */
		scan_all_boxes();

	read_all_boxes();	/* pass 2: read the entity attributes */
}

/* what load_db() runs of the routines checked, in its order */
static void
load(void)
{
	init_trade_routes();
	load_system();
	boxes();
	load_orders();
	check_db();
	determine_map_edges();
	create_relics();
	init_ocean_chars();
	delete_deadchars();
}

/* what save_db() runs of them, in its order */
static void
save(void)
{
	cleanup_posts();
	save_trade_routes();
	save_system();
	write_all_boxes();
	write_master();
	save_orders();
	rename_act_join_files();
}

/* the globals the routines set, and every player's order queues */
static void
print(void)
{
	int pl, i, j;
	struct entity_player *p;

	printf("sysclock %d %d %d\n", sysclock.turn, sysclock.day,
			sysclock.days_since_epoch);
	printf("from_host=%s\n", from_host);
	printf("reply_host=%s\n", reply_host);
	printf("gm_address=%s\n", gm_address);
	printf("game_title=%s\n", game_title);
	printf("game_url=%s\n", game_url);
	printf("rules_url=%s\n", rules_url);
	printf("times_url=%s\n", times_url);
	printf("htpasswd_loc=%s\n", htpasswd_loc);
	printf("garrison_pay=%d auto_quit=%d army_slow_factor=%d\n",
			garrison_pay, auto_quit_turns, army_slow_factor);
	printf("indep=%d gm=%d skill=%d post=%d init=%d\n",
			indep_player, gm_player, skill_player,
			post_has_been_run, seed_has_been_run);
	printf("fr=%d tr=%d ur=%d fp=%d hr=%d hp=%d hl=%d\n",
			faery_region, tunnel_region, under_region,
			faery_player, hades_region, hades_pit, hades_player);
	printf("nr=%d nl=%d np=%d cr=%d cp=%d mo=%d\n",
			nowhere_region, nowhere_loc, npc_pl, cloud_region,
			combat_pl, mount_olympus);
	printf("max_map %d %d %d\n", max_map_row, max_map_col, max_map_init);
	printf("trade_route %d %llu %u\n", trade_route_have_seed,
			(unsigned long long) trade_route_seed,
			(unsigned) trade_route_secret);
	printf("ocean_chars:");
	for (i = 0; i < ilist_len(ocean_chars); i++)
		printf(" %d", ocean_chars[i]);
	printf("\n");
	printf("trades_to_check:");
	for (i = 0; i < ilist_len(trades_to_check); i++)
		printf(" %d", trades_to_check[i]);
	printf("\n");

	loop_player(pl)
	{
		p = rp_player(pl);
		if (p == NULL)
			continue;
		for (i = 0; i < orders_len(p->orders); i++)
		{
			printf("orders %d %d:", pl, p->orders[i]->unit);
			for (j = 0; j < cstrings_len(p->orders[i]->l); j++)
				printf(" [%s]", p->orders[i]->l[j]);
			printf("\n");
		}
	}
	next_player;
}

int
main(int argc, char **argv)
{
	static char buf[65536];
	FILE *fp;
	int want, section = -1;

	/* the part of call_init_routines() in main.c that these need */
	init_lower();
	glob_init();

	if (argc == 3 && strcmp(argv[1], "-g") == 0)
	{
		libdir = argv[2];
		load();
		print();
		save();
		return 0;
	}

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
		int nread = 0, len;
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

		/* the database's lines are not echoed */
		op[0] = '\0';
		sscanf(buf, "%31s%n", op, &nread);
		pos = buf + nread;

		if (strcmp(op, "put") == 0 || strcmp(op, "putl") == 0)
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
			continue;
		}

		printf("> %s\n", buf);

		if (strcmp(op, "mkdir") == 0)
			mkdir(text(NULL), 0755);
		else if (strcmp(op, "rm") == 0)
			unlink(text(NULL));
		else if (strcmp(op, "chmod") == 0)
		{
			int mode = (int) strtol(word(), NULL, 8);
			chmod(text(NULL), mode);
		}
		else if (strcmp(op, "libdir") == 0)
			libdir = str_save(text(NULL));
		else if (strcmp(op, "seed") == 0)
		{
			trade_route_seed_pending = TRUE;
			trade_route_seed_value = strtoull(text(NULL), NULL, 0);
		}
		else if (strcmp(op, "newplayer") == 0)
			ilist_append(&new_players, atoi(text(NULL)));
		else if (strcmp(op, "trade") == 0)
			init_trade_routes();
		else if (strcmp(op, "system") == 0)
			load_system();
		else if (strcmp(op, "boxes") == 0)
			boxes();
		else if (strcmp(op, "orders") == 0)
			load_orders();
		else if (strcmp(op, "check") == 0)
			check_db();
		else if (strcmp(op, "edges") == 0)
			determine_map_edges();
		else if (strcmp(op, "relics") == 0)
			create_relics();
		else if (strcmp(op, "ocean") == 0)
			init_ocean_chars();
		else if (strcmp(op, "deadchars") == 0)
			delete_deadchars();
		else if (strcmp(op, "load") == 0)
			load();
		else if (strcmp(op, "posts") == 0)
			cleanup_posts();
		else if (strcmp(op, "savetrade") == 0)
			save_trade_routes();
		else if (strcmp(op, "savesystem") == 0)
			save_system();
		else if (strcmp(op, "saveboxes") == 0)
		{
			write_all_boxes();
			write_master();
		}
		else if (strcmp(op, "saveorders") == 0)
			save_orders();
		else if (strcmp(op, "rename") == 0)
			rename_act_join_files();
		else if (strcmp(op, "save") == 0)
			save();
		else if (strcmp(op, "print") == 0)
			print();
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
	"${WORK}/driver.c" "${WORK}/gen/helpers.c" olympia/check.c olympia/code.c olympia/glob.c \
	olympia/z.c "${WORK}/gen/stubs.c" lib/checked_alloc.c lib/ilist.c lib/rng.c lib/rnd.c \
	lib/admits.c lib/cstrings.c lib/item_ents.c lib/orders.c lib/req_ents.c lib/skill_ents.c \
	lib/trades.c \
	-o "${WORK}/loadsave-driver") >"${WORK}/build.log" 2>&1 ||
	die "C driver build failed; see ${WORK}/build.log"

# Check 1: the recorded cases.
#
# section.py runs every section in a fresh directory with umask 022 and
# writes, for each, the driver's output (its stderr goes to its stdout),
# how it ended and the directory's tree. A section that runs for more than
# RUN_TIMEOUT seconds is stopped, and the check fails. The tree lists each path
# under the directory, sorted: "d MODE PATH" for a directory ("unreadable"
# if it cannot be listed), "f MODE SIZE PATH" and the contents for a file,
# as "| LINE", escaped as the driver escapes, and "\ no newline at end"
# when the last line has none.
cat >"${WORK}/section.py" <<'EOF'
import os, signal, subprocess, sys

driver, cases, sections, outdir, runner = sys.argv[1], sys.argv[2], int(sys.argv[3]), sys.argv[4], sys.argv[5]

def esc(b):
    return "".join(chr(c) if 32 <= c <= 126 and c != 0x5c else "\\x%02x" % c for c in b)

def tree(root, out):
    def walk(rel):
        path = os.path.join(root, rel) if rel else root
        try:
            names = sorted(os.listdir(path))
        except OSError:
            out.append(f"d {os.lstat(path).st_mode & 0o777:03o} {rel} unreadable")
            return
        if rel:
            out.append(f"d {os.lstat(path).st_mode & 0o777:03o} {rel}")
        for name in names:
            r = f"{rel}/{name}" if rel else name
            p = os.path.join(root, r)
            st = os.lstat(p)
            if os.path.isdir(p) and not os.path.islink(p):
                walk(r)
            elif os.path.isfile(p) and not os.path.islink(p):
                out.append(f"f {st.st_mode & 0o777:03o} {st.st_size} {r}")
                try:
                    data = open(p, "rb").read()
                except OSError:
                    out.append("unreadable")
                    continue
                lines = data.split(b"\n")
                last = lines.pop()
                for l in lines:
                    out.append("| " + esc(l))
                if last:
                    out.append("| " + esc(last))
                    out.append("\\ no newline at end")
            else:
                out.append(f"? {r}")
    walk("")

def run(i):
    d = os.path.join(outdir, str(i))
    os.mkdir(d)
    p = subprocess.Popen(["bash", runner, driver, cases, str(i)], cwd=d, stdout=subprocess.PIPE,
                         stderr=subprocess.DEVNULL)
    data, _ = p.communicate()
    if p.returncode == 124:
        sys.exit(f"section {i} timed out")
    lines = data.split(b"\n")
    data = lines.pop()
    out = [l.decode("latin-1") for l in lines]
    if data:
        out.append(data.decode("latin-1") + " (no newline)")
    rc = p.returncode
    out.append(f"exit {128 - rc if rc < 0 else rc}")
    tree(d, out)
    for dp, dns, _ in os.walk(d):
        for n in dns:
            os.chmod(os.path.join(dp, n), 0o755)
    return "".join(l + "\n" for l in out)

os.umask(0o022)
signal.signal(signal.SIGPIPE, signal.SIG_DFL)
with open(os.path.join(outdir, "out.txt"), "w", encoding="latin-1") as f:
    for i in range(sections):
        f.write(run(i))
EOF

# Also used by go test -exec, so every driver/test process has the same
# supervisor and signal statuses, without shell crash diagnostics.
cat >"${WORK}/limited.sh" <<EOF
. "${HERE}/../lib/timeout.sh"
run_limited "${RUN_TIMEOUT}" "\$1" "\$@"
EOF

sections=$(grep -c '^section ' "${CASES}")
for run in run1 run2; do
	mkdir -p "${WORK}/${run}"
	python3 -I "${WORK}/section.py" "${WORK}/loadsave-driver" "${CASES}" "${sections}" "${WORK}/${run}" "${WORK}/limited.sh"
done
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
(cd "${ROOT}" && go test -exec "bash ${WORK}/limited.sh" -count=1 -run 'TestLoadSaveMatchesC' ./internal/olympia)

# Check 2: the saved games.

if [ "${NGAMES}" -eq 0 ]; then
	echo "saved games: not checked (--full checks every game, --games N a sample)"
	echo "PASS"
	exit 0
fi

mkdir -p "${WORK}/world" "${WORK}/games"
tar -xzf "${WORK}/g3/tests/olympia/fixtures/world.tgz" -C "${WORK}/world"
libs=("${WORK}/world/lib")
if [ -d "${GAMES}" ]; then
	for d in "${GAMES}"/*/input/lib; do
		[ -d "${d}" ] && libs+=("${d}")
	done
fi
# --games N: N of them, evenly spaced, from the first (world.tgz) to the last
if [ "${NGAMES}" -gt 0 ] && [ "${NGAMES}" -lt ${#libs[@]} ]; then
	sample=()
	for ((i = 0; i < NGAMES; i++)); do
		if [ "${NGAMES}" -eq 1 ]; then j=0; else j=$((i * (${#libs[@]} - 1) / (NGAMES - 1))); fi
		sample+=("${libs[j]}")
	done
	libs=("${sample[@]}")
fi

# The Go side runs from a test binary built once, one game at a time, so that
# each game's copies can be compared and deleted before the next are made.
echo "building the Go test binary"
(cd "${ROOT}" && go test -c -o "${WORK}/olympia.test" ./internal/olympia)

# copy LIB DIR: DIR/lib, a copy of LIB (a clone where the file system can)
copy() {
	mkdir -p "$2"
	cp -cR "$1" "$2/lib" 2>/dev/null || { rm -rf "$2/lib"; cp -R "$1" "$2/lib"; }
}

# listing DIR: every path under DIR with its mode, sorted
listing() {
	(cd "$1" && find . -exec stat -f '%Sp %N' {} + | LC_ALL=C sort)
}

# What check_db() reported, over every game: each message from C's stderr
# after its "check_db()" line, with its box codes and numbers replaced by N,
# and the number of games that gave it. A summary, for the README.
tally="${WORK}/games/messages.txt"
: >"${tally}"

umask 022
fail=0 boxes=0 total=0
dir="${WORK}/games"
k=0
for lib in "${libs[@]}"; do
	copy "${lib}" "${dir}/c/${k}"
	copy "${lib}" "${dir}/go/${k}"
	status=0
	(cd "${dir}/c/${k}" && run_limited "${RUN_TIMEOUT}" loadsave-driver "${WORK}/loadsave-driver" -g lib) \
		>"${dir}/c/${k}.out" 2>"${dir}/c/${k}.err" || status=$?
	[ "${status}" -ne "${TIMEOUT_STATUS}" ] || die "C driver timed out after ${RUN_TIMEOUT}s on ${lib}"
	echo "exit ${status}" >>"${dir}/c/${k}.out"
	list="${dir}/go/${k}.list"
	printf '%s\t%s\n' "${dir}/go/${k}" "${dir}/go/${k}" >"${list}"
	(cd "${ROOT}/internal/olympia" && G3G0_FULL=1 LOADSAVE_COMPAT_GAMES="${list}" \
		run_limited "${RUN_TIMEOUT}" olympia.test "${WORK}/olympia.test" -test.count=1 -test.run '^TestLoadSaveWorlds$' >"${dir}/go/${k}.log" 2>&1) ||
		{ cat "${dir}/go/${k}.log"; die "TestLoadSaveWorlds failed on ${lib}"; }

	ok=1
	for s in out err; do
		if ! cmp -s "${dir}/c/${k}.${s}" "${dir}/go/${k}.${s}"; then
			echo "FAIL: ${lib}: std${s} differs"
			diff -u "${dir}/c/${k}.${s}" "${dir}/go/${k}.${s}" | head -20
			ok=0
		fi
	done
	if ! diff -r "${dir}/c/${k}/lib" "${dir}/go/${k}/lib" >"${dir}/${k}.diff.txt" 2>&1; then
		echo "FAIL: ${lib}: lib differs"
		head -20 "${dir}/${k}.diff.txt"
		ok=0
	fi
	if ! diff <(listing "${dir}/c/${k}/lib") <(listing "${dir}/go/${k}/lib") >"${dir}/${k}.modes.txt"; then
		echo "FAIL: ${lib}: lib's files or modes differ"
		head -20 "${dir}/${k}.modes.txt"
		ok=0
	fi
	[ "${status}" -eq 0 ] || echo "note: ${lib}: C exits ${status}"
	n=$(grep -c '^[0-9][0-9]*	' "${dir}/c/${k}/lib/master" 2>/dev/null || true)
	[ "${k}" -eq 0 ] && boxes=${n}
	total=$((total + n))
	sed -n '/^check_db()$/,$p' "${dir}/c/${k}.err" | sed 1d |
		sed -E 's/\[[0-9a-z]+\]/[N]/g; s/[0-9]+/N/g' | LC_ALL=C sort -u >>"${tally}"
	if [ "${ok}" -eq 1 ]; then
		rm -rf "${dir}/c/${k}" "${dir}/go/${k}" "${dir}/c/${k}".* "${dir}/go/${k}".* "${dir}/${k}".*.txt
	else
		echo "kept: ${dir}/{c,go}/${k}"
		fail=1
	fi
	k=$((k + 1))
done
echo "saved games: ${#libs[@]} (world.tgz: ${boxes} boxes; all checked: ${total} boxes)"
echo "check_db() messages (games that gave each, box codes and numbers as N): $(LC_ALL=C sort -u "${tally}" | wc -l | tr -d ' ')"
LC_ALL=C sort "${tally}" | uniq -c | sed 's/^/  /'
[ "${fail}" -eq 0 ] || exit 1
echo "saved games: C and Go lib trees, stdout and stderr identical"
echo "PASS"
