package graphics

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type Drawable interface {
	Draw() (x float64, y float64, z float64, m float64)
}

type Window struct {
	Location Vector3D
	Rotation Vector3D
	Scale    float64
}

type Engine struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	width    int32
	height   int32
	Running  bool
}

func NewEngine(name string, width int32, height int32) (e *Engine, err error) {

	e = &Engine{width: width, height: height}
	err = sdl.Init(sdl.INIT_EVERYTHING)

	e.Window, err = sdl.CreateWindow(name, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		width, height, sdl.WINDOW_ALLOW_HIGHDPI)
	e.Renderer, err = sdl.CreateRenderer(e.Window, -1, sdl.RENDERER_ACCELERATED)

	// High DPI Scale
	_ = e.Renderer.SetScale(2, 2)

	HandleError(err)
	e.Running = true
	return e, err
}

func (e *Engine) Cleanup() {
	_ = e.Window.Destroy()
	_ = e.Renderer.Destroy()
}

func (e *Engine) Clear() {
	_ = e.Renderer.Clear()
}

func MatrixMultiply(a [][]float64, b []float64) (f []float64) {
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

func Project3D(d Vertex3D) Vertex2D {
	t := 0.0
	rotY := MatrixMultiply([][]float64{
		{math.Cos(t), 0, math.Sin(t)},
		{0, 1, 0},
		{-math.Sin(t), 0, math.Cos(t)},
	}, []float64{d.X, d.Y, d.Z})

	rotP := MatrixMultiply([][]float64{
		{1, 0, 0},
		{0, math.Cos(t), math.Sin(t)},
		{0, math.Sin(t), math.Cos(t)},
	}, []float64{rotY[0], rotY[1], rotY[2]})

	points := MatrixMultiply([][]float64{
		{0.75, 0, 0},
		{0, 0.75, 0},
	}, []float64{rotP[0], rotP[1], rotP[2]})

	return Vertex2D{512 + points[0], 512 + points[1]}
}

func (e *Engine) DrawCube(m Vertex3D, w float64) {

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

		e.DrawLine(start, end)

		// Project3D(Vertex3D{
		// 	X: w + vertices[i-1][0]*(w/2),
		// 	Y: w + vertices[i-1][1]*(w/2),
		// 	Z: w + vertices[i-1][2]*(w/2),
		// })
	}

}

// func (e *Engine) DrawPoint(p Vertex3D) {
// 	points := Project3D(p)
// 	e.Renderer.SetDrawColor(255, 255, 255, 255)
// 	e.Renderer.DrawPointF(float32(points.X), float32(points.Y))
// }

func (e *Engine) DrawLine(p Vertex3D, q Vertex3D) {
	start := Project3D(p)
	end := Project3D(q)
	err := e.Renderer.DrawLineF(float32(start.X), float32(start.Y), float32(end.X), float32(end.Y))
	if err != nil {
		return
	}
}

func (e *Engine) RenderDrawable(d Drawable) {
	x, y, z, _ := d.Draw()
	points := Project3D(Vertex3D{x, y, z})
	e.Renderer.DrawPointF(float32(points.X), float32(points.Y))
}

func (e *Engine) Render() {
	e.Renderer.SetDrawColor(0, 0, 0, 255)
	e.Renderer.Present()
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func (e *Engine) HandleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			e.Running = false
			break
		}
	}
}

func (e *Engine) Wait(ms int) {
	sdl.WaitEventTimeout(100)
}
