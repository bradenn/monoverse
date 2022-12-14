package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	_ "github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

// Verse
// the verse object
// Contains a shit load of hex quadrants on a grid
//
//
type Verse struct {
	name        string
	octree      *Octree
	stats       *List
	grid        [][]bool
	clocks      F2
	location    F2
	perspective F3
	rotation    F3
	size        F2
	world       *World
	physics     *Physics
	cursor      *Cursor
	timer       *Timer
}

type Cursor struct {
	location F3
}

func (v *Cursor) Select(f F2) {
	v.location = F3{f.X, f.Y, 0}
}
func (v *Cursor) Draw(g *Graphics) {
	g.Color(0.5, 0.8, 1, 0.8)
	g.Tet(v.location, F3{64, 64, 64})
}
func (v *Verse) GetName() string {
	return v.name
}

func NewVerse(name string, location F2, size F2) *Verse {
	verse := &Verse{
		name:     name,
		location: location,
		stats:    nil,
		size:     size,
	}
	return verse
}

func (v *Verse) HandleEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.MouseButtonEvent:
		if t.Button == sdl.BUTTON_LEFT {
			v.cursor.Select(F2{float64(t.X), float64(t.Y)})
		}
		break
	case *sdl.MouseMotionEvent:
		if t.State == 1 {
			v.perspective.X += float64(t.XRel) * math.E / 2
			v.perspective.Y += float64(t.YRel) * math.E / 2

		} else if t.State == 4 {
			v.rotation.X -= float64(t.YRel) / math.E
			v.rotation.Y += float64(t.XRel) / math.E
		}
		break
	case *sdl.MouseWheelEvent:
		v.perspective.Z += float64(t.Y) * 0.01
		if v.perspective.Z <= 0.00001 {
			v.perspective.Z = 0.00001
		}
		break
	}
}

func (v *Verse) GetLocation() F2 {
	return v.location
}

func (v *Verse) GetSize() F2 {
	return v.size
}

func (v *Verse) Configure() {

	v.cursor = new(Cursor)
	v.cursor.location = F3{}
	v.perspective.Z = 1280 / 1280

	v.world = &World{
		cells: nil,
		size:  256,
	}
	v.world.Populate()
}

func (v *Verse) Draw(g *Graphics) {

	v.stats.elements = []*Item{
		NewItem("Camera", ""),
		NewItem("Perspective", fmt.Sprintf("X: %.2f, Y: %.2f, Z: %.2f", v.perspective.X, v.perspective.Y,
			v.perspective.Z)),
		NewItem("Rotation", fmt.Sprintf("X: %.2f, Y: %.2f, Z: %.2f", v.rotation.X, v.rotation.Y,
			v.rotation.Z)),
	}

	gl.PushMatrix()

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.PushMatrix()

	gl.Translatef(2, 28, 0)
	gl.Color4f(0, 0, 0, 0.1)
	gl.Scalef(1, 1, 1)
	gl.Begin(gl.QUADS)
	gl.Vertex3f(0, 0, -2040)
	gl.Vertex3f(0, float32(v.size.Y), -2040)
	gl.Vertex3f(float32(v.size.X), float32(v.size.Y), -2040)
	gl.Vertex3f(float32(v.size.X), 0, -2040)
	gl.End()

	gl.PopMatrix()

	gl.PushMatrix()

	gl.Translatef(float32(v.size.X-64), 64, 64)
	gl.Scalef(20, 20, 20)
	gl.Rotatef(float32(v.rotation.X), 1, 0, 0)
	gl.Rotatef(float32(v.rotation.Y), 0, 1, 0)
	gl.Rotatef(float32(v.rotation.Z), 0, 0, 1)

	gl.Begin(gl.LINES)

	gl.Color4f(1, 0, 0, 1)
	gl.Vertex3f(-1, 0, 0)
	gl.Vertex3f(1, 0, 0)

	gl.Color4f(0, 1, 0, 1)
	gl.Vertex3f(0, -1, 0)
	gl.Vertex3f(0, 1, 0)

	gl.Color4f(0, 0, 1, 1)
	gl.Vertex3f(0, 0, -1)
	gl.Vertex3f(0, 0, 1)

	gl.End()

	gl.PopMatrix()

	gl.PushMatrix()

	gl.Translatef(float32(v.size.X/2+v.perspective.X), float32(v.size.Y/2+v.perspective.Y), 0)
	gl.Scalef(float32(v.perspective.Z), float32(v.perspective.Z), float32(v.perspective.Z))
	gl.Rotatef(float32(v.rotation.X), 1, 0, 0)
	gl.Rotatef(float32(v.rotation.Y), 0, 1, 0)
	gl.Rotatef(float32(v.rotation.Z), 0, 0, 1)

	gl.Enable(gl.BLEND)

	gl.Color4f(1, 1, 1, 0.2)

	lightPosition := []float32{0, 0, 2, 0}
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])

	gl.Lightfv(gl.LIGHT0, gl.SPOT_EXPONENT, &[]float32{2}[0])
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &[]float32{0, 0, 0.2, 0.2}[0])

	lightAmbient := []float32{0.5, 0.5, 0.5, 0.6}
	gl.Materialfv(gl.FRONT_AND_BACK, gl.AMBIENT, &lightAmbient[0])

	gl.Translatef(float32(-v.world.size/4), float32(-v.world.size/4), 0)
	v.world.Draw(g)

	gl.Disable(gl.BLEND)

	gl.PopMatrix()

	gl.PopMatrix()

	return
}

func (v *Verse) Update() {

	// v.physics.Tick()

}
