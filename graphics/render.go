package graphics

import (
	"fmt"
	"github.com/bradenn/monoverse"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"math"
)

type Renderer struct {
	Ref     *sdl.Renderer
	Font    *ttf.Font
	texture *sdl.Texture
	Camera  *main.Camera
}

func (r *Renderer) Init() {

}

func (r *Renderer) UIText(m string, v F3) {

	err := ttf.Init()
	if err != nil {
		fmt.Println(err)
	}
	if r.Font == nil {
		r.Font, _ = ttf.OpenFont("./JetBrainsMonoNL-Regular.ttf", 24)
	}
	solid, err := r.Font.RenderUTF8Blended(m, sdl.Color{
		R: 255,
		G: 131,
		B: 0,
		A: 255,
	})

	text, _ := r.Ref.CreateTextureFromSurface(solid)
	if err != nil {
		fmt.Println(err)
	}
	_ = r.Ref.SetClipRect(&sdl.Rect{
		X: int32(r.Camera.Location.X),
		Y: int32(r.Camera.Location.Y),
		W: int32(r.Camera.Size.X),
		H: int32(r.Camera.Size.Y),
	})

	err = r.Ref.CopyF(text, nil, &sdl.FRect{X: float32(r.Camera.Location.X + v.X), Y: float32(r.Camera.Location.Y + v.Y),
		W: float32(solid.W / 2),
		H: float32(solid.H / 2)})
	err = text.Destroy()
	solid.Free()
	if err != nil {
		return
	}
}

func (r *Renderer) SetPixel(f F2) error {

	err := r.Ref.DrawPointF(float32(f.X), float32(f.Y))
	if err != nil {
		return err
	}
	return nil
}

func (r *Renderer) UIRect(rect F3, size F2) {
	_ = r.Ref.DrawRectF(&sdl.FRect{X: float32(rect.X), Y: float32(rect.Y), W: float32(size.X), H: float32(size.Y)})
}

func (r *Renderer) SetColor(a float64, b float64, c float64, d float64) error {
	err := r.Ref.SetDrawColor(uint8(a), uint8(b), uint8(c), uint8(d))
	if err != nil {
		return err
	}
	return nil
}

func (r *Renderer) RenderPoint(Fp F3) error {
	err := r.SetPixel(r.project3D(Fp))
	if err != nil {
		return err
	}
	return nil
}

func matMul(a [][]float64, b []float64) (f []float64) {
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

func (r *Renderer) project3D(d F3) F2 {
	err := r.Ref.SetClipRect(&sdl.Rect{
		X: int32(r.Camera.Location.X),
		Y: int32(r.Camera.Location.Y),
		W: int32(r.Camera.Size.X),
		H: int32(r.Camera.Size.Y),
	})
	if err != nil {
		return F2{}
	}
	// X Rotation
	x := r.Camera.Rotation.X
	rotX := matMul([][]float64{
		{math.Cos(x), 0, math.Sin(x)},
		{0, 1, 0},
		{-math.Sin(x), 0, math.Cos(x)},
	}, []float64{d.X, d.Y, d.Z})
	// Y Rotation
	y := r.Camera.Rotation.Y
	rotY := matMul([][]float64{
		{1, 0, 0},
		{0, math.Cos(y), math.Sin(y)},
		{0, math.Sin(y), math.Cos(y)},
	}, rotX)
	// Z Rotation
	z := r.Camera.Rotation.Z
	rotZ := matMul([][]float64{
		{math.Cos(z), -math.Sin(z), 0},
		{math.Sin(z), math.Cos(z), 0},
		{0, 0, 0},
	}, rotY)
	// X, Y, Z Translation
	translation := matMul([][]float64{
		{1, 0, 0, r.Camera.Location.X + r.Camera.Size.X/2},
		{0, 1, 0, r.Camera.Location.Y + r.Camera.Size.Y/2},
		{0, 0, 1, r.Camera.Location.Z},
		{0, 0, 0, 1},
	}, []float64{rotZ[0], rotZ[1], rotZ[2], 1})
	scale := matMul([][]float64{
		{r.Camera.Scale, 0, 0, 0},
		{0, r.Camera.Scale, 0, 0},
		{0, 0, r.Camera.Scale, 0},
		{0, 0, 0, 1},
	}, []float64{translation[0], translation[1], translation[2], 1})

	// camera := matMul([][]float64{
	// 	{1, 0, 0, 0},
	// 	{0, 1, 0, 0},
	// 	{0, 0, 1, 0},
	// 	{0, 0, 0, 1},
	// }, translation)
	// far := 10
	// near := 1000
	// clip := matMul([][]float64{
	// 	{1 / r.Camera.Size.X, 0, 0, 0},
	// 	{0, 1 / r.Camera.Size.Y, 0, 0},
	// 	{0, 0, float64((far + near) / (far - near)), 1},
	// 	{0, 0, float64(-(2 * far * near) / (far - near)), 0},
	// }, []float64{scale[0], scale[1], scale[2], 1})
	//
	//
	points := matMul([][]float64{
		{1, 0, 0},
		{0, 1, 0},
	}, []float64{scale[0], scale[1], scale[2]})
	// if points[0] >= r.Camera.Location.X && points[0] <= r.Camera.Location.X + r.Camera.Size.X && points[1] >= r.Camera.Location.Y && points[1] <= r.Camera.Location.Y + r.Camera.Size.Y{
	//
	// }
	return F2{points[0], points[1]}

}

func (r *Renderer) DrawLine(a F3, b F3) {
	start := r.project3D(a)
	end := r.project3D(b)

	err := r.Ref.DrawLineF(float32(start.X), float32(start.Y), float32(end.X), float32(end.Y))
	if err != nil {
		return
	}

}

func (r *Renderer) DrawHex(v F3, w float64) {
	var last F3
	for j := 0.0; j < 7; j += 1.0 {
		lf := F3{
			X: v.X + math.Cos((60*j-30)*(math.Pi/180))*w,
			Y: v.Y + math.Sin((60*j-30)*(math.Pi/180))*w,
			Z: v.Z,
		}
		if last.X == 0 {
			last = lf
		}
		r.DrawLine(last, lf)
		last = lf
	}

}

func (r *Renderer) DrawCube(v F3, w float64) {

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
		start := F3{
			X: v.X + vertices[i][0]*(w/2),
			Y: v.Y + vertices[i][1]*(w/2),
			Z: v.Z + vertices[i][2]*(w/2),
		}
		end := F3{
			X: v.X + vertices[i+1][0]*(w/2),
			Y: v.Y + vertices[i+1][1]*(w/2),
			Z: v.Z + vertices[i+1][2]*(w/2),
		}
		r.DrawLine(start, end)
	}

}

func (r *Renderer) DrawRect(location F3, size F2) {
	point := r.project3D(location)
	_ = r.Ref.DrawRectF(&sdl.FRect{
		X: float32(point.X),
		Y: float32(point.Y),
		W: float32(size.X),
		H: float32(size.Y),
	})
}

func (r *Renderer) DrawPoint(location F3) {
	err := r.SetPixel(r.project3D(location))
	if err != nil {
		return
	}
}
func (r *Renderer) FillRect(location F3, size F3) {
	point := r.project3D(location)
	_ = r.Ref.FillRectF(&sdl.FRect{
		X: float32(point.X),
		Y: float32(point.Y),
		W: float32(size.X),
		H: float32(size.Y),
	})
}

func (r *Renderer) DrawRoundRect(location F3, size F2) {
	point := r.project3D(location)
	gfx.RoundedRectangleRGBA(r.Ref, int32(point.X), int32(point.Y), int32(size.X),
		int32(size.Y),
		8, 255,
		131, 0,
		128)
}

func (p *Pane) DrawText() {
	font, _ := ttf.OpenFont("IBMPlexSansCondensed-Regular.ttf", 3)

	surface, _ := font.RenderUTF8Solid("ANAL", sdl.Color{255, 0, 255, 0})
	tex, _ := p.engine.Renderer.CreateTextureFromSurface(surface)
	p.engine.Renderer.Copy(tex, nil, nil)
}
