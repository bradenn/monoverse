package graphics

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type Camera struct {
	Location Vector3D
	Rotation Vector3D
	Scale    float64
}

type Pane struct {
	engine  *Engine
	texture *sdl.Texture
	width   float64
	height  float64
	camera  Camera
}

func NewPane(eng *Engine, width float64, height float64) *Pane {
	pane := &Pane{
		width:  width,
		height: height,
		camera: Camera{},
		engine: eng,
	}
	return pane
}

func (p *Pane) SetScale(scale float64) {
	p.camera.Scale = scale
}

func (p *Pane) RotateXBy(rot float64) {
	p.camera.Rotation.X += rot
}

func (p *Pane) RotateYBy(rot float64) {
	p.camera.Rotation.Y += rot
}

func matrixMultiply(a [][]float64, b []float64) (f []float64) {
	colsA := len(a[0])
	rowsA := len(a)
	f = make([]float64, rowsA)

	colsB := len(b)

	for i := 0; i < rowsA; i++ {
		for j := 0; j < colsB; j++ {
			sum := 0.0
			for k := 0; k < colsA; k++ {
				sum += a[i][k] * b[k]
			}
			f[i] = sum
		}
	}
	return f
}

func (p *Pane) project3D(d Vertex3D) Vertex2D {
	x := p.camera.Rotation.X
	rotY := matrixMultiply([][]float64{
		{math.Cos(x), 0, math.Sin(x)},
		{0, 1, 0},
		{-math.Sin(x), 0, math.Cos(x)},
	}, []float64{d.X, d.Y, d.Z})
	y := p.camera.Rotation.Y
	rotP := matrixMultiply([][]float64{
		{1, 0, 0},
		{0, math.Cos(y), math.Sin(y)},
		{0, math.Sin(y), math.Cos(y)},
	}, []float64{rotY[0], rotY[1], rotY[2]})

	points := matrixMultiply([][]float64{
		{p.camera.Scale, 0, 0},
		{0, p.camera.Scale, 0},
	}, []float64{rotP[0], rotP[1], rotP[2]})

	return Vertex2D{512 + points[0], 512 + points[1]}
}

func (p *Pane) DrawRhombicDodecahedron(n Vertex3D) {

	// vertices := [][]float64{
	//
	// 	// 1 1 0
	// 	// -1 1 0
	// 	// -1 -1 0
	// 	// 1 -1 0
	// 	//
	//
	// 	{1, 1, -1},
	// 	{1, -1, -1},
	// 	{1, 1, 1},
	// 	{1, -1, 1},
	// 	{-1, 1, -1},
	// 	{-1, -1, -1},
	// 	{-1, 1, 1},
	// 	{-1, -1, 1},
	//
	// 	{2, 0, 0},
	//
	// 	{0, 2, 2},
	// 	{0, 2, -2},
	// 	{0, -2, 2},
	// 	{0, -2, -2},
	//
	// }
	p.DrawHex(n, 100)
	// last := Vertex3D{}
	// for _, v := range vertices {
	// 	curr := Vertex3D{n.X + (v[0])*64, n.Y + (v[1])*64, n.Z + (v[2])*64}
	// 	p.DrawLine(curr, last)
	// 	last = curr
	// }

}
func (p *Pane) DrawPoint(n Vertex3D) {
	points := p.project3D(n)
	p.engine.Renderer.DrawPointF(float32(points.X), float32(points.Y))
}

func (p *Pane) DrawLine(n Vertex3D, m Vertex3D) {
	start := p.project3D(n)
	end := p.project3D(m)
	err := p.engine.Renderer.DrawLineF(float32(start.X), float32(start.Y), float32(end.X), float32(end.Y))
	if err != nil {
		return
	}

}

func (p *Pane) DrawHex(m Vertex3D, w float64) {
	var last Vertex3D
	for j := 0.0; j < 7; j += 1.0 {
		lf := Vertex3D{
			X: m.X + math.Cos((60*j-30)*(math.Pi/180))*w,
			Y: m.Y + math.Sin((60*j-30)*(math.Pi/180))*w,
			Z: m.Z,
		}
		if last.X == 0 {
			last = lf
		}
		p.DrawLine(last, lf)
		last = lf
	}

}

func (p *Pane) DrawCube(m Vertex3D, w float64) {

	vertices := [][]float64{
		{1, 1, 1}, {1, 1, -1},
		{1, 1, 1}, {-1, 1, 1},

		{1, 1, 1}, {1, -1, 1},
		{-1, -1, -1}, {1, -1, -1},

		{-1, -1, -1}, {-1, -1, 1},
		{-1, -1, -1}, {-1, 1, -1},

		{-1, 1, -1}, {1, 1, -1},
		{-1, 1, -1}, {-1, 1, 1},

		{-1, -1, 1}, {-1, 1, 1},
		{-1, -1, 1}, {1, -1, 1},

		{1, -1, -1}, {1, -1, 1},
		{1, -1, -1}, {1, 1, -1},
	}

	for i := 0; i < len(vertices); i += 2 {
		start := Vertex3D{
			X: m.X + vertices[i][0]*(w/2),
			Y: m.Y + vertices[i][1]*(w/2),
			Z: m.Z + vertices[i][2]*(w/2),
		}
		end := Vertex3D{
			X: m.X + vertices[i+1][0]*(w/2),
			Y: m.Y + vertices[i+1][1]*(w/2),
			Z: m.Z + vertices[i+1][2]*(w/2),
		}

		p.DrawLine(start, end)

		// Project3D(Vertex3D{
		// 	X: w + vertices[i-1][0]*(w/2),
		// 	Y: w + vertices[i-1][1]*(w/2),
		// 	Z: w + vertices[i-1][2]*(w/2),
		// })
	}

}
