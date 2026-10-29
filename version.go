// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package g3g0

import (
	"github.com/maloquacious/semver"
)

var (
	version = semver.Version{
		Major:      0,
		Minor:      115,
		Patch:      1,
		PreRelease: "alpha",
	}
)

func Version() semver.Version {
	return version
}
