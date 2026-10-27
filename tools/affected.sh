#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Lists the compat tools that a move of the olympia pin from OLD to NEW
# affects: those that read an upstream file that changed between the two
# revisions. See tools/README.md for the table this encodes; keep the two in
# step when a tool starts reading another upstream file.
#
# Usage: tools/affected.sh OLD NEW
#
# Prints one tool per line, then, on stderr, the changed files that no
# tool reads. Tools that build the whole C engine (olympia-cli-compat,
# olympia-roundtrip, olympia-spec) are listed whenever anything changed.
# genesis-compat, island-compat and mapgen-compat have pins of their own and
# are not considered.
#
# Environment:
#   G3_REPO  G3 clone that contains OLD and NEW
#            (default ~/Software/playbymail/olympia-g3)

set -euo pipefail

G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"

[ $# -eq 2 ] || { echo "usage: $0 OLD NEW" >&2; exit 2; }

changed=()
while IFS= read -r f; do
	changed+=("${f}")
done < <(git -C "${G3_REPO}" diff --name-only "$1" "$2")
[ ${#changed[@]} -gt 0 ] || exit 0

# reads TOOL PATH succeeds if TOOL reads the upstream PATH.
reads() {
	case "$1" in
	rng-compat)
		case "$2" in
		lib/rng.[ch] | lib/rnd.[ch] | lib/ilist.c | lib/lists.h) return 0 ;;
		olympia/buy.c | olympia/combat.c | olympia/npc.c | olympia/quest.c) return 0 ;;
		olympia/cloud.c | olympia/faery.c | olympia/hades.c | olympia/tunnel.c) return 0 ;;
		olympia/basic.c | olympia/c1.c | olympia/code.c | olympia/dir.c) return 0 ;;
		olympia/day.c | olympia/seed.c | olympia/storm.c | olympia/swear.c) return 0 ;;
		olympia/use.c | olympia/u.c) return 0 ;;
		esac ;;
	lists-compat)
		case "$2" in
		lib/*.[ch]) return 0 ;;
		esac ;;
	entity-compat)
		case "$2" in
		olympia/*.h | lib/*.h | olympia/glob.c | lib/checked_alloc.c) return 0 ;;
		esac ;;
	code-compat)
		case "$2" in
		olympia/*.h | lib/*.h) return 0 ;;
		olympia/code.c | olympia/glob.c | olympia/z.c | olympia/io.c | olympia/sout.c) return 0 ;;
		lib/checked_alloc.c | lib/ilist.c | lib/rng.c | lib/rnd.c) return 0 ;;
		esac ;;
	helper-compat)
		case "$2" in
		olympia/*.h | lib/*.h) return 0 ;;
		olympia/code.c | olympia/glob.c | olympia/z.c | olympia/sout.c | olympia/u.c) return 0 ;;
		olympia/seed.c | olympia/loc.c | olympia/stack.c | olympia/swear.c | olympia/use.c) return 0 ;;
		olympia/move.c | olympia/day.c | olympia/garr.c | olympia/order.c | olympia/input.c) return 0 ;;
		lib/checked_alloc.c | lib/ilist.c | lib/rng.c | lib/rnd.c) return 0 ;;
		lib/cstrings.c | lib/item_ents.c | lib/orders.c) return 0 ;;
		esac ;;
	check-compat)
		reads helper-compat "$2" && return 0
		case "$2" in
		olympia/check.c | olympia/io.c | tests/olympia/fixtures/world.tgz) return 0 ;;
		lib/admits.c | lib/req_ents.c | lib/skill_ents.c | lib/trades.c) return 0 ;;
		esac ;;
	io-compat | box-compat | boxfile-compat)
		case "$2" in
		olympia/*.h | lib/*.h) return 0 ;;
		olympia/code.c | olympia/glob.c | olympia/z.c | olympia/io.c) return 0 ;;
		olympia/sout.c | olympia/u.c | olympia/seed.c) return 0 ;;
		lib/checked_alloc.c | lib/ilist.c | lib/rng.c | lib/rnd.c) return 0 ;;
		lib/admits.c | lib/item_ents.c | lib/req_ents.c | lib/skill_ents.c | lib/trades.c) return 0 ;;
		esac
		[ "$1" = io-compat ] && return 1
		case "$2" in
		olympia/input.c | lib/cstrings.c | tests/olympia/fixtures/world.tgz) return 0 ;;
		esac
		[ "$1" = boxfile-compat ] && [ "$2" = olympia/swear.c ] && return 0
		;;
	olympia-cli-compat | olympia-roundtrip | olympia-spec)
		return 0 ;;
	esac
	return 1
}

tools=(rng-compat lists-compat entity-compat code-compat helper-compat io-compat box-compat boxfile-compat
	check-compat scanf-compat olympia-cli-compat olympia-roundtrip olympia-spec)

for t in "${tools[@]}"; do
	for f in "${changed[@]}"; do
		if reads "${t}" "${f}"; then
			echo "${t}"
			break
		fi
	done
done

# The files only the whole-engine tools read, for a reviewer's eye.
for f in "${changed[@]}"; do
	hit=0
	for t in rng-compat lists-compat entity-compat code-compat helper-compat io-compat box-compat boxfile-compat check-compat; do
		if reads "${t}" "${f}"; then hit=1; break; fi
	done
	[ "${hit}" -eq 1 ] || echo "only the whole-engine tools read: ${f}" >&2
done
