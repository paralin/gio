// SPDX-License-Identifier: Unlicense OR MIT

/*
Package key implements key and text events and operations.

The InputOp operations is used for declaring key input handlers. Use
an implementation of the Queue interface from package ui to receive
events.
*/
package key

import (
	"strings"

	"gioui.org/internal/opconst"
	"gioui.org/io/event"
	"gioui.org/op"
)

// InputOp declares a handler ready for key events.
// Key events are in general only delivered to the
// focused key handler. Set the Focus flag to request
// the focus.
type InputOp struct {
	Key   event.Key
	Focus bool
}

// HideInputOp request that any on screen text input
// be hidden.
type HideInputOp struct{}

// A FocusEvent is generated when a handler gains or loses
// focus.
type FocusEvent struct {
	Focus bool
}

// An Event is generated when a key is pressed. For text input
// use EditEvent.
type Event struct {
	// Name is the rune character that most closely
	// match the key. For letters, the upper case form
	// is used.
	Name rune
	// Modifiers is the set of active modifiers when
	// the key was pressed.
	Modifiers Modifiers
}

// An EditEvent is generated when text is input.
type EditEvent struct {
	Text string
}

// Modifiers
type Modifiers uint32

const (
	// ModCommand is the command modifier. On macOS
	// it is the Cmd key, on other platforms the Ctrl
	// key.
	ModCommand Modifiers = 1 << iota
	// THe shift key.
	ModShift
)

const (
	// Runes for special keys.
	NameLeftArrow      = '←'
	NameRightArrow     = '→'
	NameUpArrow        = '↑'
	NameDownArrow      = '↓'
	NameReturn         = '⏎'
	NameEnter          = '⌤'
	NameEscape         = '⎋'
	NameHome           = '⇱'
	NameEnd            = '⇲'
	NameDeleteBackward = '⌫'
	NameDeleteForward  = '⌦'
	NamePageUp         = '⇞'
	NamePageDown       = '⇟'
)

// Contain reports whether m contains all modifiers
// in m2.
func (m Modifiers) Contain(m2 Modifiers) bool {
	return m&m2 == m2
}

func (h InputOp) Add(o *op.Ops) {
	data := o.Write(opconst.TypeKeyInputLen, h.Key)
	data[0] = byte(opconst.TypeKeyInput)
	if h.Focus {
		data[1] = 1
	}
}

func (h HideInputOp) Add(o *op.Ops) {
	data := o.Write(opconst.TypeHideInputLen)
	data[0] = byte(opconst.TypeHideInput)
}

func (EditEvent) ImplementsEvent()  {}
func (Event) ImplementsEvent()      {}
func (FocusEvent) ImplementsEvent() {}

func (e Event) String() string {
	return "{" + string(e.Name) + " " + e.Modifiers.String() + "}"
}

func (m Modifiers) String() string {
	var strs []string
	if m.Contain(ModCommand) {
		strs = append(strs, "ModCommand")
	}
	if m.Contain(ModShift) {
		strs = append(strs, "ModShift")
	}
	return strings.Join(strs, "|")
}
