# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later

"""Compare exact paths, file types and permission bits; contents stay with diff -r.

Shared by boxfile/check/loadsave-compat. Exit 0 for equality, 1 for a
difference, 2 for a listing error. Never compare partial or failed listings.
"""

import difflib
import json
import os
import stat
import sys


def listing(root):
    rows = []

    def visit(path, relative):
        mode = os.lstat(path).st_mode
        # JSON quoting keeps even whitespace/newlines and undecodable names exact.
        rows.append(f"{mode:o} {json.dumps(relative)}\n")
        if stat.S_ISDIR(mode):
            with os.scandir(path) as entries:
                names = sorted(entry.name for entry in entries)
            for name in names:
                visit(os.path.join(path, name), relative + "/" + name)

    visit(root, ".")
    return rows


def main(left, right):
    try:
        before = listing(left)
        after = listing(right)
    except OSError as error:
        print(f"tree-modes: listing failed: {error}", file=sys.stderr)
        return 2
    if before == after:
        return 0
    sys.stdout.writelines(difflib.unified_diff(before, after, left, right))
    return 1


if __name__ == "__main__":
    sys.exit(main(*sys.argv[1:]))
