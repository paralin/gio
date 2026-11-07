// SPDX-License-Identifier: Unlicense OR MIT

package app

import (
	"testing"

	"gioui.org/app/internal/windows"
)

func TestPointerButtonChange(t *testing.T) {
	const (
		first   = windows.POINTER_FLAG_FIRSTBUTTON
		second  = windows.POINTER_FLAG_SECONDBUTTON
		inRange = windows.POINTER_FLAG_INRANGE
	)
	tests := []struct {
		name       string
		prev, next windows.POINTER_INFO_POINTER_FLAGS
		want       windows.POINTER_BUTTON_CHANGE_TYPE
	}{
		{"move", inRange, inRange, windows.POINTER_CHANGE_NONE},
		{"first down", inRange, inRange | first, windows.POINTER_CHANGE_FIRSTBUTTON_DOWN},
		{"first up", inRange | first, inRange, windows.POINTER_CHANGE_FIRSTBUTTON_UP},
		{"drag", inRange | first, inRange | first, windows.POINTER_CHANGE_NONE},
		{"second down while first held", first, first | second, windows.POINTER_CHANGE_SECONDBUTTON_DOWN},
		{"fifth up", windows.POINTER_FLAG_FIFTHBUTTON, 0, windows.POINTER_CHANGE_FIFTHBUTTON_UP},
	}
	for _, tt := range tests {
		if got := pointerButtonChange(tt.prev, tt.next); got != tt.want {
			t.Errorf("%s: pointerButtonChange(%#x, %#x) = %d, want %d", tt.name, tt.prev, tt.next, got, tt.want)
		}
	}
}

// TestPointerMessageFlags checks that the high word of a WM_POINTER wParam
// maps onto the POINTER_FLAG constants. The value is a WM_POINTERUPDATE that
// Wine sent for a mouse move under CrossOver.
func TestPointerMessageFlags(t *testing.T) {
	const wParam = 0x20020001
	flags := windows.POINTER_INFO_POINTER_FLAGS(uint16(wParam >> 16))
	if want := windows.POINTER_FLAG_PRIMARY | windows.POINTER_FLAG_INRANGE; flags != want {
		t.Errorf("flags = %#x, want %#x", flags, want)
	}
	if id := getPointerIDwParam(wParam); id != 1 {
		t.Errorf("pointer id = %d, want 1", id)
	}
}
