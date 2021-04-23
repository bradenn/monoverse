package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
	_ "github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"math/rand"
	"time"
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
	window      *Window
	location    F2
	perspective F3
	rotation    F3
	size        F2
}

func NewVerse(name string, location F2, size F2) *Verse {
	verse := &Verse{
		name:     name,
		ticks:    0,
		radius:   0,
		matter:   nil,
		location: location,
		size:     size,
	}
	return verse
}

func (v *Verse) HandleEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.MouseButtonEvent:
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
		v.perspective.Z += float64(t.Y)
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
	for i := 0.0; i < 3000; i++ {
		matter := new(Matter)
		matter.position = F3{-diam/2 + rand.Float64()*diam, -diam/2 + rand.Float64()*diam,
			-diam/2 + rand.Float64()*diam}
		v.matter = append(v.matter, matter)
	}
	// w.updateLastTick = time.Now()
	// w.renderLastTick = time.Now()
}

func (v *Verse) Draw(g *Graphics) {
	g.DrawFrame("Verse", v)
	gl.Scissor(8, int32(v.location.Y*2)+28+8,
		int32((v.GetSize().X)*2)-14,
		int32((v.GetSize().Y-28)*2)-2)
	gl.Enable(gl.SCISSOR_TEST)

	v.perspective.X += (-v.perspective.X) / 60
	v.perspective.Y += (-v.perspective.Y) / 60
	v.perspective.Z = math.Min(v.perspective.Z, 40)
	v.perspective.Z = math.Max(v.perspective.Z, 0) // No below ZEWOWOWS !!!
	v.ticks++

	gl.PushMatrix()

	gl.MatrixMode(gl.PROJECTION)

	// gl.Frustum(0, v.size.X, v.size.Y,0, 0, 100)
	thing := mgl32.LookAt(58, 36, 15, 0, 0, 0, 15, 15, 15)
	id, _ := g.window.GetID()
	gl.Uniform4fv(gl.GetUniformLocation(id, nil), 4, (*float32)(gl.Ptr(&thing[0])))
	// gl.matrixmu
	divisions := 10.0
	gl.MatrixMode(gl.MODELVIEW)
	gl.Translatef(float32(v.size.X/2+v.perspective.X), float32(v.size.Y/2+v.perspective.Y), 0)
	gl.Rotatef(float32(v.rotation.X), 1, 0, 0)
	gl.Rotatef(float32(v.rotation.Y), 0, 1, 0)
	gl.Rotatef(float32(v.rotation.Z), 0, 0, 1)
	// gl.Rotatef(float32(v.ticks/2), 0, math.E/10, 0)
	gl.Scalef(float32(v.perspective.Z+1), float32(v.perspective.Z+1), float32(v.perspective.Z+1))
	gl.Begin(gl.LINES)
	for i := -1.0; i < 1.0; i += 1 / divisions {
		gl.Vertex3f(float32(i*512), -512, 0)
		gl.Vertex3f(float32(i*512), 512, 0)

		gl.Vertex3f(-512, float32(i*512), 0)
		gl.Vertex3f(512, float32(i*512), 0)

		gl.Vertex3f(float32(i*512), 0, -512)
		gl.Vertex3f(float32(i*512), 0, 512)

		gl.Vertex3f(-512, 0, float32(i*512))
		gl.Vertex3f(512, 0, float32(i*512))

		gl.Vertex3f(0, float32(i*512), -512)
		gl.Vertex3f(0, float32(i*512), 512)

		gl.Vertex3f(0, -512, float32(i*512))
		gl.Vertex3f(0, 512, float32(i*512))

	}
	gl.Vertex3f(0, 0, 0)
	gl.Vertex3f(512, 0, 0)
	gl.Vertex3f(0, 0, 0)
	gl.Vertex3f(0, 512, 0)
	gl.Vertex3f(0, 0, 0)
	gl.Vertex3f(0, 0, 512)
	gl.End()
	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.COLOR_MATERIAL)

	gold := [13]float32{0.24725, 0.1995, 0.0745, 1.0, /* ambient */
		0.75164, 0.60648, 0.22648, 1.0, /* diffuse */
		0.628281, 0.555802, 0.366065, 1.0, /* specular */
		50.0, /* shininess */
	}

	gl.Lightfv(gl.LIGHT0, gl.POSITION, (*float32)(gl.Ptr([]float32{256, 256, 0, 1})))
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, (*float32)(gl.Ptr([]float32{1, 1, 1, 1})))
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, (*float32)(gl.Ptr([]float32{0.2, 0.2, 0.2, 1})))
	gl.Lightfv(gl.LIGHT0, gl.SPECULAR, (*float32)(gl.Ptr([]float32{1, 1, 1, 1})))
	gl.Materialfv(gl.FRONT_AND_BACK, gl.EMISSION, &gold[0])
	gl.Materialfv(gl.FRONT_AND_BACK, gl.DIFFUSE, &gold[4])
	gl.Materialfv(gl.FRONT_AND_BACK, gl.SPECULAR, &gold[8])
	gl.Materialfv(gl.FRONT_AND_BACK, gl.SHININESS, &gold[12])
	gl.Enable(gl.LIGHT0)
	for _, mass := range v.matter {

		g.Cube(mass.position, F3{16, 16, 16})

	}
	gl.Disable(gl.LIGHT0)
	gl.Disable(gl.LIGHTING)
	gl.PopMatrix()
	gl.Disable(gl.SCISSOR_TEST)

	return
}

func (v *Verse) Update() {
	// v.updateDelta = time.Duration(time.Since(v.updateLastTick).Nanoseconds())
	// v.updateLastTick = time.Now()
	//

}

func (v *Verse) computeIteration() (err error) {

	return

}

type Entropy struct {
	ticks    F1
	complete B
	step     F1
	relative F1
	start    time.Time
}

func NewEntropy() *Entropy {
	entropy := &Entropy{ticks: 0.0, step: 1.0, start: time.Now()}

	return entropy
}

func (e *Entropy) Print() {
	fmt.Printf("%f ticks, %f seconds", e.ticks, e.relative)
}

func (e *Entropy) Available() B {
	if e.complete {
		e.complete = false
		return true
	} else {
		fmt.Errorf("%s", "Tick Courrupted")
	}
	return true
}

func (e *Entropy) Ticks() F1 {
	return e.ticks
}

func (e *Entropy) Tick() {
	e.ticks += e.step
	e.complete = true
	e.relative = F1(time.Since(e.start).Seconds())
}

type Matter struct {
	position, velocity, force F3
	mass, density             F1
}

func sdds() {
	entropy := NewEntropy()

	for entropy.Available() {

		entropy.Tick()
	}
	entropy.Print()

}
