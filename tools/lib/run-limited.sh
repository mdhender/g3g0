#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
# Adapter for go test -exec and upstream fixtures' engine overrides.
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${HERE}/timeout.sh"
run_limited "${RUN_TIMEOUT:-120}" "$1" "$@"
