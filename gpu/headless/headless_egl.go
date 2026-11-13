// SPDX-License-Identifier: Unlicense OR MIT

//go:build linux || freebsd || openbsd

package headless

import (
	"gioui.org/internal/egl"
)

func init() {
	newContextPrimary = func() (context, error) {
		return egl.NewContext(egl.PlatformUnknown, egl.EGL_DEFAULT_DISPLAY)
	}
}
