// SPDX-License-Identifier: Unlicense OR MIT

// Package material implements the Material design.
package material

import (
	"image"
	"image/color"
	"image/draw"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"golang.org/x/exp/shiny/iconvg"
)

type Button struct {
	Text string
	// Color is the text color.
	Color      color.RGBA
	Font       text.Font
	Background color.RGBA

	shaper *text.Shaper
}

type IconButton struct {
	Background color.RGBA
	Icon       *Icon
	Size       unit.Value
	Padding    unit.Value
}

type Icon struct {
	src  []byte
	size unit.Value

	// Cached values.
	op      paint.ImageOp
	imgSize int
}

func (t *Theme) Button(txt string) Button {
	return Button{
		Text:       txt,
		Color:      rgb(0xffffff),
		Background: t.Color.Primary,
		Font: text.Font{
			Size: t.TextSize.Scale(14.0 / 16.0),
		},
		shaper: t.Shaper,
	}
}

// NewIcon returns a new Icon from IconVG data.
func NewIcon(data []byte) (*Icon, error) {
	_, err := iconvg.DecodeMetadata(data)
	if err != nil {
		return nil, err
	}
	return &Icon{src: data}, nil
}

func (t *Theme) IconButton(icon *Icon) IconButton {
	return IconButton{
		Background: t.Color.Primary,
		Icon:       icon,
		Size:       unit.Dp(56),
		Padding:    unit.Dp(16),
	}
}

func (b Button) Layout(gtx *layout.Context, button *widget.Button) {
	col := b.Color
	bgcol := b.Background
	if !button.Active() {
		col = rgb(0x888888)
		bgcol = rgb(0xcccccc)
	}
	st := layout.Stack{Alignment: layout.Center}
	hmin := gtx.Constraints.Width.Min
	vmin := gtx.Constraints.Height.Min
	lbl := st.Rigid(gtx, func() {
		gtx.Constraints.Width.Min = hmin
		gtx.Constraints.Height.Min = vmin
		layout.Align(layout.Center).Layout(gtx, func() {
			layout.UniformInset(unit.Dp(16)).Layout(gtx, func() {
				paint.ColorOp{Color: col}.Add(gtx.Ops)
				widget.Label{}.Layout(gtx, b.shaper, b.Font, b.Text)
			})
		})
		pointer.RectAreaOp{Rect: image.Rectangle{Max: gtx.Dimensions.Size}}.Add(gtx.Ops)
		button.Layout(gtx)
	})
	bg := st.Expand(gtx, func() {
		rr := float32(gtx.Px(unit.Dp(4)))
		rrect(gtx.Ops,
			float32(gtx.Constraints.Width.Min),
			float32(gtx.Constraints.Height.Min),
			rr, rr, rr, rr,
		)
		fill(gtx, bgcol)
		for _, c := range button.History() {
			drawInk(gtx, c)
		}
	})
	st.Layout(gtx, bg, lbl)
}

func (b IconButton) Layout(gtx *layout.Context, button *widget.Button) {
	st := layout.Stack{}
	ico := st.Rigid(gtx, func() {
		layout.UniformInset(b.Padding).Layout(gtx, func() {
			size := gtx.Px(b.Size) - 2*gtx.Px(b.Padding)
			if b.Icon != nil {
				b.Icon.Layout(gtx, unit.Px(float32(size)))
			}
			gtx.Dimensions = layout.Dimensions{
				Size: image.Point{X: size, Y: size},
			}
		})
		pointer.EllipseAreaOp{Rect: image.Rectangle{Max: gtx.Dimensions.Size}}.Add(gtx.Ops)
		button.Layout(gtx)
	})
	bgcol := b.Background
	if !button.Active() {
		bgcol = rgb(0xcccccc)
	}
	bg := st.Expand(gtx, func() {
		size := float32(gtx.Constraints.Width.Min)
		rr := float32(size) * .5
		rrect(gtx.Ops,
			size,
			size,
			rr, rr, rr, rr,
		)
		fill(gtx, bgcol)
		for _, c := range button.History() {
			drawInk(gtx, c)
		}
	})
	st.Layout(gtx, bg, ico)
}

func (ic *Icon) Layout(gtx *layout.Context, sz unit.Value) {
	ico := ic.image(gtx.Px(sz))
	ico.Add(gtx.Ops)
	paint.PaintOp{
		Rect: f32.Rectangle{
			Max: toPointF(ico.Size()),
		},
	}.Add(gtx.Ops)
}

func (ic *Icon) image(sz int) paint.ImageOp {
	if sz == ic.imgSize {
		return ic.op
	}
	m, _ := iconvg.DecodeMetadata(ic.src)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	// Use white for icons.
	m.Palette[0] = color.RGBA{A: 0xff, R: 0xff, G: 0xff, B: 0xff}
	iconvg.Decode(&ico, ic.src, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	ic.op = paint.NewImageOp(img)
	ic.imgSize = sz
	return ic.op
}

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func toRectF(r image.Rectangle) f32.Rectangle {
	return f32.Rectangle{
		Min: toPointF(r.Min),
		Max: toPointF(r.Max),
	}
}

func drawInk(gtx *layout.Context, c widget.Click) {
	d := gtx.Now().Sub(c.Time)
	t := float32(d.Seconds())
	const duration = 0.5
	if t > duration {
		return
	}
	t = t / duration
	var stack op.StackOp
	stack.Push(gtx.Ops)
	size := float32(gtx.Px(unit.Dp(700))) * t
	rr := size * .5
	col := byte(0xaa * (1 - t*t))
	ink := paint.ColorOp{Color: color.RGBA{A: col, R: col, G: col, B: col}}
	ink.Add(gtx.Ops)
	op.TransformOp{}.Offset(c.Position).Offset(f32.Point{
		X: -rr,
		Y: -rr,
	}).Add(gtx.Ops)
	rrect(gtx.Ops, float32(size), float32(size), rr, rr, rr, rr)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(size), Y: float32(size)}}}.Add(gtx.Ops)
	stack.Pop()
	op.InvalidateOp{}.Add(gtx.Ops)
}
