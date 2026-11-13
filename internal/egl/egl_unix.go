// SPDX-License-Identifier: Unlicense OR MIT

//go:build linux || freebsd || openbsd

package egl

/*
#cgo linux,!android  pkg-config: egl
#cgo freebsd openbsd android LDFLAGS: -lEGL
#cgo freebsd CFLAGS: -I/usr/local/include
#cgo freebsd LDFLAGS: -L/usr/local/lib
#cgo openbsd CFLAGS: -I/usr/X11R6/include
#cgo openbsd LDFLAGS: -L/usr/X11R6/lib
#cgo CFLAGS: -DEGL_NO_X11

#include <EGL/egl.h>
#include <EGL/eglext.h>

// gio_eglGetPlatformDisplay calls eglGetPlatformDisplay when khr is set and
// eglGetPlatformDisplayEXT otherwise, loading either at run time since the
// link library may predate EGL 1.5. It sets *ok to 0 when the entry point
// is missing.
static EGLDisplay gio_eglGetPlatformDisplay(int khr, EGLenum platform, void *native, int *ok) {
	*ok = 1;
	if (khr) {
		PFNEGLGETPLATFORMDISPLAYPROC get = (PFNEGLGETPLATFORMDISPLAYPROC)eglGetProcAddress("eglGetPlatformDisplay");
		if (get != NULL) {
			return get(platform, native, NULL);
		}
	} else {
		PFNEGLGETPLATFORMDISPLAYEXTPROC get = (PFNEGLGETPLATFORMDISPLAYEXTPROC)eglGetProcAddress("eglGetPlatformDisplayEXT");
		if (get != NULL) {
			return get(platform, native, NULL);
		}
	}
	*ok = 0;
	return EGL_NO_DISPLAY;
}
*/
import "C"

import (
	"slices"
	"strings"
	"unsafe"
)

type (
	_EGLint           = C.EGLint
	_EGLDisplay       = C.EGLDisplay
	_EGLConfig        = C.EGLConfig
	_EGLContext       = C.EGLContext
	_EGLSurface       = C.EGLSurface
	NativeDisplayType = C.EGLNativeDisplayType
	NativeWindowType  = C.EGLNativeWindowType
)

func loadEGL() error {
	return nil
}

func eglChooseConfig(disp _EGLDisplay, attribs []_EGLint) (_EGLConfig, bool) {
	var cfg C.EGLConfig
	var ncfg C.EGLint
	if C.eglChooseConfig(disp, &attribs[0], &cfg, 1, &ncfg) != C.EGL_TRUE {
		return nilEGLConfig, false
	}
	return _EGLConfig(cfg), true
}

func eglCreateContext(disp _EGLDisplay, cfg _EGLConfig, shareCtx _EGLContext, attribs []_EGLint) _EGLContext {
	ctx := C.eglCreateContext(disp, cfg, shareCtx, &attribs[0])
	return _EGLContext(ctx)
}

func eglDestroySurface(disp _EGLDisplay, surf _EGLSurface) bool {
	return C.eglDestroySurface(disp, surf) == C.EGL_TRUE
}

func eglDestroyContext(disp _EGLDisplay, ctx _EGLContext) bool {
	return C.eglDestroyContext(disp, ctx) == C.EGL_TRUE
}

func eglGetConfigAttrib(disp _EGLDisplay, cfg _EGLConfig, attr _EGLint) (_EGLint, bool) {
	var val _EGLint
	ret := C.eglGetConfigAttrib(disp, cfg, attr, &val)
	return val, ret == C.EGL_TRUE
}

func eglGetError() _EGLint {
	return C.eglGetError()
}

func eglInitialize(disp _EGLDisplay) (_EGLint, _EGLint, bool) {
	var maj, min _EGLint
	ret := C.eglInitialize(disp, &maj, &min)
	return maj, min, ret == C.EGL_TRUE
}

func eglMakeCurrent(disp _EGLDisplay, draw, read _EGLSurface, ctx _EGLContext) bool {
	return C.eglMakeCurrent(disp, draw, read, ctx) == C.EGL_TRUE
}

func eglReleaseThread() bool {
	return C.eglReleaseThread() == C.EGL_TRUE
}

func eglSwapBuffers(disp _EGLDisplay, surf _EGLSurface) bool {
	return C.eglSwapBuffers(disp, surf) == C.EGL_TRUE
}

func eglSwapInterval(disp _EGLDisplay, interval _EGLint) bool {
	return C.eglSwapInterval(disp, interval) == C.EGL_TRUE
}

func eglTerminate(disp _EGLDisplay) bool {
	return C.eglTerminate(disp) == C.EGL_TRUE
}

func eglQueryString(disp _EGLDisplay, name _EGLint) string {
	return C.GoString(C.eglQueryString(disp, name))
}

func eglGetDisplay(disp NativeDisplayType) _EGLDisplay {
	return C.eglGetDisplay(disp)
}

// eglGetPlatformDisplay returns the display for disp on the named platform,
// with named false when the platform is unknown or EGL lacks the client
// extension for it.
func eglGetPlatformDisplay(platform Platform, disp NativeDisplayType) (_EGLDisplay, bool) {
	var name string
	switch platform {
	case PlatformX11:
		name = "platform_x11"
	case PlatformWayland:
		name = "platform_wayland"
	default:
		return nilEGLDisplay, false
	}
	// Client extensions are queried on EGL_NO_DISPLAY; without
	// EGL_EXT_client_extensions the query fails and returns NULL.
	exts := C.eglQueryString(nilEGLDisplay, C.EGL_EXTENSIONS)
	if exts == nil {
		return nilEGLDisplay, false
	}
	list := strings.Fields(C.GoString(exts))
	for _, khr := range []bool{true, false} {
		ext := "EGL_EXT_" + name
		if khr {
			ext = "EGL_KHR_" + name
		}
		if !slices.Contains(list, ext) {
			continue
		}
		var k, ok C.int
		if khr {
			k = 1
		}
		d := C.gio_eglGetPlatformDisplay(k, C.EGLenum(platform), unsafe.Pointer(disp), &ok)
		if ok != 0 {
			return d, true
		}
	}
	return nilEGLDisplay, false
}

func eglCreateWindowSurface(disp _EGLDisplay, conf _EGLConfig, win NativeWindowType, attribs []_EGLint) _EGLSurface {
	eglSurf := C.eglCreateWindowSurface(disp, conf, win, &attribs[0])
	return eglSurf
}

func eglWaitClient() bool {
	return C.eglWaitClient() == C.EGL_TRUE
}
