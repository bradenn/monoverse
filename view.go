package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"math"
)

type ViewB struct {
	name        string
	focused     bool
	location    F2
	perspective F3
	fov         float64
	rotation    F3
	size        F2
	scale       float64
}

func (v *ViewB) HandleClick(location F2, button bool) {
	v.focused = true
}

func (v *ViewB) RotateDelta(d F3) {
	v.rotation.X += d.X
	v.rotation.Y += d.Y
	v.rotation.Z += d.Z
}
func (v *ViewB) GetBounds() (F2, F2) {
	return F2{v.location.X, v.location.Y}, F2{v.size.X, v.size.Y}
}
func (v *ViewB) Configure() {

}
func (v *ViewB) Update() {

}

func NewView(name string, location F2, size F2) (w *ViewB) {
	w = &ViewB{
		name:        name,
		focused:     false,
		location:    location,
		perspective: F3{},
		fov:         0,
		rotation:    F3{},
		size:        size,
		scale:       1,
	}
	return w
}

func (v *ViewB) Contains(f F2) bool {
	xBound := v.location.X <= f.X && v.location.X+v.size.X >= f.X
	yBound := v.location.Y <= f.Y && v.location.Y+v.size.Y >= f.Y
	v.focused = xBound && yBound
	return xBound && yBound
}

func (v *ViewB) Draw(graphics *Graphics) {
	// _ = v.SetColor(255, 131, 0, 255)
	gl.Enable(gl.BLEND)
	gl.PushMatrix()
	// gl.Frustum(0, v.size.X, v.size.Y, 0, 0.001, 1024)
	gl.Ortho(0, v.size.X, v.size.Y, 0, 128, 1024)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(float32(v.location.X), float32(v.location.Y), 0)

	gl.Enable(gl.BLEND)

	if v.focused {
		graphics.Color(1, 1, 1, 0.8)
		graphics.Text(fmt.Sprintf("X: %.2f, Y: %.2f, Z: %.2f", v.perspective.X, v.perspective.Y, v.perspective.Z), F2{X: 8, Y: 36})
		graphics.Text(fmt.Sprintf("FOV: %.2f°, Focal: %.2fmm, Scale: %.2f", v.fov, 2*math.Atan(8/v.fov), v.scale), F2{X: 8, Y: 48})
	}
	// gl.Ortho(v.location.X, v.location.X+v.size.X, v.location.Y+v.size.Y, v.location.Y, 0.001, 1024)
	color := 1 / 255.0
	gl.LineWidth(2)
	graphics.Color(color*66, color*66, color*66, 1)
	graphics.Rect(
		F2{
			X: 1,
			Y: 1,
		},
		F2{
			X: v.size.X - 2,
			Y: v.size.Y - 2,
		})
	gl.LineWidth(1)
	graphics.Color(color*45, color*45, color*45, 1)
	graphics.FillRect(
		F2{
			X: 1,
			Y: 1,
		},
		F2{
			X: v.size.X - 2,
			Y: 28 - 2,
		})

	graphics.Rect(
		F2{
			X: 1,
			Y: 1,
		},
		F2{
			X: v.size.X - 2,
			Y: 28 - 2,
		})
	graphics.Color(1, 1, 1, 0.8)
	graphics.Text(fmt.Sprintf("%s", v.name), F2{8, 6})
	gl.Disable(gl.BLEND)
	// graphics.Rect(
	// 	F2{
	// 		X: 2,
	// 		Y: 2,
	// 	},
	// 	F2{
	// 		X: v.size.X - 4,
	// 		Y: v.size.Y - 4,
	// 	})
	gl.PopMatrix()
	if v.focused {
		gl.PushMatrix()
		// gl.Ortho(0, v.size.X, v.size.Y, 0, -1024, 1024)
		gl.MatrixMode(gl.MODELVIEW)
		// gl.perspec(v.location.X, v.location.X+v.size.X, v.location.Y+v.size.Y, v.location.Y, 0.001, 1024)
		gl.LoadIdentity()
		gl.Enable(gl.BLEND)
		graphics.Color(color*106, color*106, color*106, 0.21)
		outer := 256.0
		gl.Translatef(float32(v.size.X/2), float32(v.size.Y/2), 0)
		gl.Translatef(float32(v.location.X+v.perspective.X), float32(v.location.X+v.perspective.Y), float32(v.location.X+v.perspective.Z))
		gl.Rotatef(float32(v.rotation.X), 1, 0, 0)
		gl.Rotatef(float32(v.rotation.Y), 0, 1, 0)

		for n := -1.0; n < 1; n += 0.1 {
			graphics.Line(F3{-outer, n * outer, 0}, F3{outer, n * outer, 0})
			graphics.Line(F3{n * outer, -outer, 0}, F3{n * outer, outer, 0})

			graphics.Line(F3{-outer, 0, n * outer}, F3{outer, 0, n * outer})
			graphics.Line(F3{n * outer, 0, -outer}, F3{n * outer, 0, outer})

			graphics.Line(F3{0, -outer, n * outer}, F3{0, outer, n * outer})
			graphics.Line(F3{0, n * outer, -outer}, F3{0, n * outer, outer})
			graphics.Line(F3{0, 0, 0}, F3{500, 0, 0})
		}
		gl.Disable(gl.BLEND)
		gl.PopMatrix()
	}
	gl.Disable(gl.BLEND)
}
