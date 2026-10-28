# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later

"""Local comparator controls, run by all three compat tools, even without --full."""

import contextlib
import io
import os
from pathlib import Path
import runpy
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import patch


SCRIPT = Path(__file__).with_name("tree-modes.py")
compare = runpy.run_path(str(SCRIPT))["main"]


class TreeModesTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.roots = [Path(self.temp.name) / name for name in ("c", "go")]
        for root in self.roots:
            root.mkdir(mode=0o755)
            (root / "dir").mkdir(mode=0o755)
            (root / "dir/file").write_bytes(b"same bytes\n")
            (root / "dir/file").chmod(0o644)
            (root / "dangling").symlink_to("absent")

    def result(self):
        out, err = io.StringIO(), io.StringIO()
        with contextlib.redirect_stdout(out), contextlib.redirect_stderr(err):
            status = compare(*(str(root) for root in self.roots))
        return status, out.getvalue(), err.getvalue()

    def test_equal(self):
        self.assertEqual(self.result(), (0, "", ""))

    def test_mode(self):
        for side in range(2):
            for path, mode in (("dir/file", 0o600), ("dir", 0o700), ("", 0o1777)):
                with self.subTest(side=side, path=path):
                    target = self.roots[side] / path
                    original = target.stat().st_mode & 0o7777
                    target.chmod(mode)
                    self.assertEqual(self.result()[0], 1)
                    target.chmod(original)

    def test_type(self):
        for side in range(2):
            target = self.roots[side] / "dir/file"
            target.unlink()
            target.mkdir(mode=0o644)  # same permission bits, different type
            self.assertEqual(self.result()[0], 1)
            target.rmdir()
            target.symlink_to("absent")
            self.assertEqual(self.result()[0], 1)
            target.unlink()
            target.write_bytes(b"same bytes\n")
            target.chmod(0o644)

    def test_path(self):
        for side in range(2):
            target = self.roots[side] / "extra\n path\t"
            target.touch()
            self.assertEqual(self.result()[0], 1)
            target.unlink()
            original = self.roots[side] / "dir/file"
            original.rename(target)
            self.assertEqual(self.result()[0], 1)
            target.rename(original)

    def test_listing_errors(self):
        # Inject actual traversal/stat failures, independent of root privileges.
        # Errors on either side or BOTH must never yield equality or a partial diff.
        for operation in ("scandir", "lstat"):
            original = getattr(os, operation)
            for sides in ((0,), (1,), (0, 1)):
                blocked = [str(self.roots[side] / "dir") for side in sides]

                def fail(path):
                    if path in blocked:
                        raise PermissionError(13, "control: denied", path)
                    return original(path)

                with self.subTest(operation=operation, sides=sides):
                    with patch.object(os, operation, side_effect=fail):
                        status, out, err = self.result()
                    self.assertEqual(status, 2)
                    self.assertEqual(out, "")
                    self.assertIn("listing failed", err)
                    self.assertIn("control: denied", err)

    def test_cli_errors(self):
        missing = str(Path(self.temp.name) / "missing")
        for args in ((missing, str(self.roots[1])), (str(self.roots[0]), missing),
                     (missing, missing)):
            result = subprocess.run([sys.executable, "-I", str(SCRIPT), *args],
                                    capture_output=True, text=True)
            self.assertEqual(result.returncode, 2)
            self.assertEqual(result.stdout, "")
            self.assertIn("listing failed", result.stderr)


if __name__ == "__main__":
    unittest.main()
