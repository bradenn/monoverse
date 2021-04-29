package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	_ "github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"math/rand"
)

// Verse
// the verse object
// Contains a shit load of hex quadrants on a grid
//
//
type Verse struct {
	name        string
	ticks       F1
	radius      F1
	matter      []*Matter
	force       []Force
	window      *Window
	octree      *Octree
	stats       *List
	clocks      F2
	location    F2
	perspective F3
	rotation    F3
	size        F2
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
		ticks:    0,
		radius:   0,
		matter:   nil,
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
	diam := 1024.0
	for x := 0.0; x < 256; x += 1 {
		matter := new(Matter)
		matter.mass = 1000 * rand.Float64()
		matter.location = F3{diam/2 - rand.Float64()*(diam), diam/2 - rand.Float64()*(diam),
			diam/2 - rand.Float64()*(diam)}
		v.matter = append(v.matter, matter)
	}
	v.cursor = new(Cursor)
	v.cursor.location = F3{}
	v.perspective.Z = 1

}

func (v *Verse) Draw(g *Graphics) {
	gl.PopMatrix()

	v.stats.elements = []*Item{
		NewItem("Camera", ""),
		NewItem("Perspective", fmt.Sprintf("X: %.2f, Y: %.2f, Z: %.2f", v.perspective.X, v.perspective.Y,
			v.perspective.Z)),
		NewItem("Rotation", fmt.Sprintf("X: %.2f, Y: %.2f, Z: %.2f", v.rotation.X, v.rotation.Y,
			v.rotation.Z)),
		NewItem("", ""),
		NewItem("Simulation", ""),
		NewItem("Ticks", fmt.Sprintf("Tick: %.2f", v.ticks)),
	}

	v.perspective.X += (-v.perspective.X) / 60
	v.perspective.Y += (-v.perspective.Y) / 60

	gl.PushMatrix()

	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.Translatef(2, 28, 0)
	gl.Color4f(0.01, 0.02, 0.02, 0.1)
	gl.Scalef(1, 1, 1)
	gl.Begin(gl.QUADS)
	gl.Vertex3f(0, 0, -2040)
	gl.Vertex3f(0, float32(v.size.Y), -2040)
	gl.Vertex3f(float32(v.size.X), float32(v.size.Y), -2040)
	gl.Vertex3f(float32(v.size.X), 0, -2040)
	gl.End()
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
	gl.Color4f(1, 1, 1, 0.2)
	gl.Enable(gl.BLEND)
	gl.Begin(gl.LINES)
	// divisions := 512.0
	gl.LineWidth(1)
	diam := float32(math.Max(v.size.X, v.size.Y) * 10)
	// for i := -1.0; i < 1.0; i += 1 / divisions {
	// 	// gl.Vertex3f(float32(i*512), -512, 0)
	// 	// gl.Vertex3f(float32(i*512), 512, 0)
	// 	//
	// 	// gl.Vertex3f(-512, float32(i*512), 0)
	// 	// gl.Vertex3f(512, float32(i*512), 0)
	//
	// 	// gl.Vertex3f(float32(i)*diam, 0, -diam)
	// 	// gl.Vertex3f(float32(i)*diam, 0, diam)
	// 	//
	// 	// gl.Vertex3f(-diam, 0, float32(i)*diam)
	// 	// gl.Vertex3f(diam, 0, float32(i)*diam)
	//
	// 	// gl.Vertex3f(0, float32(i*512), -512)
	// 	// gl.Vertex3f(0, float32(i*512), 512)
	// 	//
	// 	// gl.Vertex3f(0, -512, float32(i*512))
	// 	// gl.Vertex3f(0, 512, float32(i*512))
	//
	// }
	gl.End()
	gl.Begin(gl.LINES)
	gl.LineWidth(4)
	gl.Color4f(1, 0, 0, 1)
	gl.Vertex3f(-diam, 0, 0)
	gl.Vertex3f(diam, 0, 0)

	gl.Color4f(0, 1, 0, 1)
	gl.Vertex3f(0, -diam, 0)
	gl.Vertex3f(0, diam, 0)

	gl.Color4f(0, 0, 1, 1)
	gl.Vertex3f(0, 0, -diam)
	gl.Vertex3f(0, 0, diam)

	gl.End()
	gl.Color4f(1, 1, 1, 1)

	for _, mass := range v.matter {
		mass.Draw(g)
	}
	if v.octree != nil {
		// v.octree.Draw(g)

	}
	//
	// gl.Color4f(0.9, 0.3, 0.7, 0.4)
	// g.Circle(F3{0, 0, 0}, F3{800, 800, 0})

	gl.Disable(gl.BLEND)
	gl.PopMatrix()

	gl.PopMatrix()
	v.cursor.Draw(g)
	return
}

func (v *Verse) Update() {
	v.timer.Begin("Update")

	v.timer.Flag("Physics")
	v.ticks += 0.1
	v.octree = &Octree{voxel: Voxel{
		location: F3{0, 0, 0},
		size:     F3{1024, 1024, 1024},
	}}

	v.timer.Flag("Matter")
	for _, mass := range v.matter {
		v.physics.ResetForces(mass)
		v.octree.Insert(mass)
	}
	v.timer.Flag("Gravity")
	for _, mass := range v.matter {
		for _, force := range v.physics.forces {
			v.octree.ApplyForces(force, mass)
		}
	}

	v.timer.Flag("Updating")
	for _, mass := range v.matter {
		loc, vel := v.physics.UpdatePosition(mass, float64(v.ticks))
		mass.SetLocation(loc)
		mass.SetVelocity(vel)
	}

	v.timer.End()
}
