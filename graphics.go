package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	_ "github.com/nullboundary/glfont"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"math"
)

type Graphics struct {
	window   *sdl.Window
	program  uint32
	renderer *sdl.Renderer
	font     *ttf.Font
	ctx      sdl.GLContext
	size     F3
}

func NewGraphics(size F3) (g *Graphics, err error) {
	g = &Graphics{size: size}
	err = g.configure()

	return g, err
}

func (g *Graphics) drawSphere(r float64, lats float64, longs float64) {
	i := 0.0
	j := 0.0
	for i = 0.0; i <= lats; i++ {
		lat0 := math.Pi * (-0.5 + (i-1)/lats)
		z0 := math.Sin(lat0)
		zr0 := math.Cos(lat0)

		lat1 := math.Pi * (-0.5 + i/lats)
		z1 := math.Sin(lat1)
		zr1 := math.Cos(lat1)

		gl.Begin(gl.QUAD_STRIP)
		for j = 0; j <= longs; j++ {
			lng := 2 * math.Pi * (j - 1) / longs
			x := math.Cos(lng)
			y := math.Sin(lng)

			gl.Normal3f(float32(x*zr0), float32(y*zr0), float32(z0))
			gl.Vertex3f(float32(r*x*zr0), float32(r*y*zr0), float32(r*z0))
			gl.Normal3f(float32(x*zr1), float32(y*zr1), float32(z1))
			gl.Vertex3f(float32(r*x*zr1), float32(r*y*zr1), float32(r*z1))
		}
		gl.End()
	}
}

func (g *Graphics) configure() (err error) {

	err = sdl.Init(sdl.INIT_EVERYTHING)
	err = ttf.Init()

	sdl.SetHint(sdl.HINT_RENDER_DRIVER, sdl.HINT_RENDER_OPENGL_SHADERS)

	err = sdl.GLSetAttribute(sdl.GL_ACCELERATED_VISUAL, 1)
	err = sdl.GLSetAttribute(sdl.GL_BUFFER_SIZE, 8)
	err = sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	err = sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 2048)
	err = sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 8)

	g.window, _ = sdl.CreateWindow("Monoverse - v0.12 beta", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(g.size.X), int32(g.size.Y), sdl.WINDOW_OPENGL|sdl.WINDOW_ALLOW_HIGHDPI)
	g.window.SetResizable(true)
	g.ctx, _ = g.window.GLCreateContext()

	err = gl.Init()

	g.font = g.loadFont("JetBrainsMono-Medium.ttf", 24)

	gl.Viewport(0, 0, int32(g.size.X*2), int32(g.size.Y*2))

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	gl.Ortho(0, g.size.X, g.size.Y, 0, -4096, 4096)

	gl.MatrixMode(gl.MODELVIEW)

	gl.Enable(gl.DOUBLEBUFFER)

	program, err := newProgramShader("/vertex.glsl", "/fragment.glsl")
	if err != nil {
		return err
	}
	gl.UseProgram(program)

	return nil
}

func (g *Graphics) Circle(n F3, m F3) {
	gl.PushMatrix()
	gl.Scalef(float32(m.X), float32(m.Y), float32(m.Z))
	gl.Begin(gl.LINE_LOOP)
	for r := 0.0; r < 360; r += 1 {
		gl.Vertex3f(float32(n.X+math.Cos(r)), float32(n.Y+math.Sin(r)), float32(n.Z))
	}
	gl.End()
	gl.Begin(gl.LINE_LOOP)
	for r := 0.0; r < 360; r += 0.2 {
		gl.Vertex3f(float32(n.X), float32(n.Y+math.Cos(r)), float32(n.Z+math.Sin(r)))
	}
	gl.End()
	gl.PopMatrix()
}

func (g *Graphics) RenderView(view View) {
	gl.Enable(gl.BLEND)

	gl.Enable(gl.DEPTH_TEST)

	gl.ClearColor(0, 0, 0, 1)
	gl.ClearDepth(1)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.PushMatrix()
	gl.Translatef(float32(view.GetLocation().X), float32(view.GetLocation().Y), 0)
	location := view.GetLocation()
	size := view.GetSize()
	gl.Scissor(int32(location.X+2)*2, int32((g.size.Y-(location.Y))-size.Y-25)*2, int32(size.X)*2,
		int32(size.Y)*2)
	g.DrawFrame(view.GetName(), view)
	gl.Translatef(0, 28, 0)
	gl.Enable(gl.SCISSOR_TEST)

	view.Draw(g)

	gl.Disable(gl.SCISSOR_TEST)
	gl.PopMatrix()

	gl.Disable(gl.DEPTH_TEST)

	gl.Disable(gl.BLEND)

}

func (g *Graphics) DrawFrame(name string, w View) {
	g.ColorSecondary()
	g.FillRect(F2{1, 1}, F2{w.GetSize().X - 2, 28 - 2})
	g.ColorText()
	g.Text(name, F2{8, 6})
	g.ColorPrimary()
	g.Rect(F2{1, 1}, F2{w.GetSize().X - 2, w.GetSize().Y - 2 + 28})
	g.ColorPrimary()
	g.Rect(F2{2, 2}, F2{w.GetSize().X - 4, w.GetSize().Y - 4 + 28})
}

func (g *Graphics) renderTextCenter(name string, location F3) {

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.TEXTURE_2D)
	surface, err := g.font.RenderUTF8Blended(name, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	defer func() { surface.Free() }()
	if err != nil {
		return
	}

	var textId uint32
	gl.GenTextures(1, &textId)
	gl.BindTexture(gl.TEXTURE_2D, textId)

	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, surface.W,
		surface.H, 0, gl.BGRA, gl.UNSIGNED_BYTE, surface.Data())

	gl.PushMatrix()

	gl.Translatef(float32(location.X-float64(surface.W/4)), float32(location.Y-float64(surface.H/4)),
		float32(location.Z))

	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(0, 0)

	gl.TexCoord2f(1, 0)
	gl.Vertex2f(float32(surface.W)/2, 0)

	gl.TexCoord2f(1, 1)
	gl.Vertex2f(float32(surface.W)/2, float32(surface.H)/2)

	gl.TexCoord2f(0, 1)
	gl.Vertex2f(0, float32(surface.H)/2)
	gl.End()
	gl.DeleteTextures(1, &textId) // Hehe, don't forget this
	gl.PopMatrix()

	gl.Disable(gl.TEXTURE_2D)
	gl.Disable(gl.BLEND)

	//

}

func (g *Graphics) renderText(name string, location F3) {

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.TEXTURE_2D)
	surface, err := g.font.RenderUTF8Blended(name, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	defer func() { surface.Free() }()
	if err != nil {
		return
	}

	var textId uint32
	gl.GenTextures(1, &textId)
	gl.BindTexture(gl.TEXTURE_2D, textId)

	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, surface.W,
		surface.H, 0, gl.BGRA, gl.UNSIGNED_BYTE, surface.Data())

	gl.PushMatrix()

	gl.Translatef(float32(location.X), float32(location.Y), float32(location.Z))

	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(0, 0)

	gl.TexCoord2f(1, 0)
	gl.Vertex2f(float32(surface.W)/2, 0)

	gl.TexCoord2f(1, 1)
	gl.Vertex2f(float32(surface.W)/2, float32(surface.H)/2)

	gl.TexCoord2f(0, 1)
	gl.Vertex2f(0, float32(surface.H)/2)
	gl.End()
	gl.DeleteTextures(1, &textId) // Hehe, don't forget this
	gl.PopMatrix()

	gl.Disable(gl.TEXTURE_2D)
	gl.Disable(gl.BLEND)

	//

}

func (g *Graphics) loadFont(name string, scale float64) *ttf.Font {
	font, err := ttf.OpenFont(name, int(scale))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return font
}

func (g *Graphics) Clear() {
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)
	gl.LineWidth(1)
	gl.Enable(gl.BLEND)

}

func (g *Graphics) Render() {
	gl.Disable(gl.BLEND)
	g.Color(1, 1, 1, 1)
	g.window.GLSwap()
}

func (g *Graphics) Text(m string, location F2) {
	gl.PushMatrix()
	g.renderText(m, F3{location.X, location.Y, 1})

	gl.PopMatrix()
}

func (g *Graphics) TextCenter(m string, location F2) {
	gl.PushMatrix()
	g.renderTextCenter(m, F3{location.X, location.Y, 1})

	gl.PopMatrix()
}

func (g *Graphics) TextCenter3D(m string, location F3) {
	gl.PushMatrix()
	g.renderTextCenter(m, location)

	gl.PopMatrix()
}

func (g *Graphics) Color(a, b, c, d float64) {
	gl.Color4f(float32(a), float32(b), float32(c), float32(d))
	gl.Materialfv(gl.FRONT_AND_BACK, gl.DIFFUSE, &[]float32{float32(a), float32(b), float32(c), float32(d)}[0])
}
func (g *Graphics) Rect(n F2, m F2) {

	gl.Begin(gl.LINE_LOOP)
	gl.Vertex3f(float32(n.X), float32(n.Y), 0)
	// gl.Vertex3f(float32(n.X+m.X), float32(n.Y), 0)
	gl.Vertex3f(float32(n.X+m.X), float32(n.Y), 0)
	// gl.Vertex3f(float32(n.X+m.X), float32(n.Y+m.Y), 0)
	gl.Vertex3f(float32(n.X+m.X), float32(n.Y+m.Y), 0)
	// gl.Vertex3f(float32(n.X), float32(n.Y+m.Y), 0)
	gl.Vertex3f(float32(n.X), float32(n.Y+m.Y), 0)
	// gl.Vertex3f(float32(n.X), float32(n.Y), 0)
	gl.End()
}

func (g *Graphics) Hexagon(n F2) {
	gl.PushMatrix()
	gl.Translatef(float32(n.X), float32(n.Y), 0)
	gl.Begin(gl.POLYGON)
	for i := 0.0; i < 6; i++ {
		gl.Vertex3f(float32(math.Sin(i/6.0*2*math.Pi)), float32(math.Cos(i/6.0*2*math.Pi)), 0)
	}
	gl.End()
	gl.PopMatrix()
}

func (g *Graphics) FillRect(n F2, m F2) {
	gl.Rectf(float32(n.X), float32(n.Y), float32(n.X+m.X), float32(n.Y+m.Y))
}

func (g *Graphics) Line(n F3, m F3) {
	gl.Begin(gl.LINES)
	gl.Vertex3f(float32(n.X), float32(n.Y), float32(n.Z))
	gl.Vertex3f(float32(m.X), float32(m.Y), float32(m.Z))
	gl.End()
}

func (g *Graphics) WireCube(n F3, m F3) {
	gl.PushMatrix()
	gl.Enable(gl.BLEND)
	gl.Translatef(float32(n.X), float32(n.Y), float32(n.Z))
	gl.Scalef(float32(m.X), float32(m.Y), float32(m.Z))
	cube := []float32{
		-1.0, -1.0, -1.0, // triangle 1 : begin
		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0, // triangle 1 : end
		1.0, 1.0, -1.0, // triangle 2 : begin
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0, // triangle 2 : end
		1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
	}

	gl.Begin(gl.LINE_LOOP)
	for i := 0; i < len(cube)/3; i++ {
		gl.Vertex3f(cube[i*3+0], cube[i*3+1], cube[i*3+2])
	}
	gl.End()

	gl.Disable(gl.BLEND)
	gl.PopMatrix()
}

func (g *Graphics) Tet(n F3, m F3) {
	gl.PushMatrix()
	gl.Enable(gl.BLEND)
	gl.Translatef(float32(n.X), float32(n.Y), float32(n.Z))
	gl.Scalef(float32(m.X/2), float32(m.Y/2), float32(m.Z/2))
	vertices := [][]float64{
		{1, 1, 1}, {1, 1, -1},
		{1, 1, 1}, {-1, 1, 1},

		{1, 1, 1}, {1, -1, 1},
		{-1, -1, -1}, {1, -1, -1},

		{-1, -1, -1}, {-1, -1, 1},
		{-1, -1, -1}, {-1, 1, -1},

		{-1, 1, -1}, {1, 1, -1},
		{-1, 1, -1}, {-1, 1, 1},

		{-1, 1, -1}, {-1, 1, 1},
		{-1, -1, 1}, {-1, 1, 1},

		{-1, -1, 1}, {-1, 1, 1},
		{-1, -1, 1}, {1, -1, 1},

		{1, -1, -1}, {1, -1, 1},
		{1, -1, -1}, {1, 1, -1},

		{-1, -1, 1}, {1, -1, 1},
		{-1, 1, 1}, {1, 1, 1},
	}
	gl.Begin(gl.LINES)
	for i := 0; i < len(vertices); i += 1 {

		gl.Vertex3f(float32(vertices[i][0]), float32(vertices[i][1]), float32(vertices[i][2]))
	}
	gl.End()

	gl.Disable(gl.BLEND)
	gl.PopMatrix()
}

func (g *Graphics) Sphere(n F3, m F3) {

}

func (g *Graphics) Cube(n F3, m F3) {
	gl.PushMatrix()
	gl.Enable(gl.BLEND)
	gl.Enable(gl.DOUBLEBUFFER)
	gl.Translatef(float32(n.X), float32(n.Y), float32(n.Z))
	gl.Scalef(float32(m.X), float32(m.Y), float32(m.Z))

	gl.Begin(gl.QUADS)
	gl.Vertex3f(1.0, 1.0, -1.0)
	gl.Vertex3f(-1.0, 1.0, -1.0)
	gl.Vertex3f(-1.0, 1.0, 1.0)
	gl.Vertex3f(1.0, 1.0, 1.0)

	// Bottom face (y = -1.0)
	// gl.Color3f(1.0, 0.5, 0.0) // Orange
	gl.Vertex3f(1.0, -1.0, 1.0)
	gl.Vertex3f(-1.0, -1.0, 1.0)
	gl.Vertex3f(-1.0, -1.0, -1.0)
	gl.Vertex3f(1.0, -1.0, -1.0)

	// Front face  (z = 1.0)
	// gl.Color3f(1.0, 0.0, 0.0) // Red
	gl.Vertex3f(1.0, 1.0, 1.0)
	gl.Vertex3f(-1.0, 1.0, 1.0)
	gl.Vertex3f(-1.0, -1.0, 1.0)
	gl.Vertex3f(1.0, -1.0, 1.0)

	// Back face (z = -1.0)
	// gl.Color3f(1.0, 1.0, 0.0) // Yellow
	gl.Vertex3f(1.0, -1.0, -1.0)
	gl.Vertex3f(-1.0, -1.0, -1.0)
	gl.Vertex3f(-1.0, 1.0, -1.0)
	gl.Vertex3f(1.0, 1.0, -1.0)

	// Left face (x = -1.0)
	// gl.Color3f(0.0, 0.0, 1.0) // Blue
	gl.Vertex3f(-1.0, 1.0, 1.0)
	gl.Vertex3f(-1.0, 1.0, -1.0)
	gl.Vertex3f(-1.0, -1.0, -1.0)
	gl.Vertex3f(-1.0, -1.0, 1.0)

	// Right face (x = 1.0)
	// gl.Color3f(1.0, 0.0, 1.0) // Magenta
	gl.Vertex3f(1.0, 1.0, -1.0)
	gl.Vertex3f(1.0, 1.0, 1.0)
	gl.Vertex3f(1.0, -1.0, 1.0)
	gl.Vertex3f(1.0, -1.0, -1.0)
	gl.End()

	gl.Disable(gl.BLEND)
	gl.Disable(gl.DOUBLEBUFFER)
	gl.PopMatrix()
}

func (g *Graphics) ColorPrimary() {
	color := 1 / 255.0
	gl.Color4f(float32(color*66), float32(color*66), float32(color*66), 1)
}
func (g *Graphics) ColorSecondary() {
	color := 1 / 255.0
	gl.Color4f(float32(color*45), float32(color*45), float32(color*45), 1)
}
func (g *Graphics) ColorText() {
	gl.Color4f(0.8, 0.8, 0.8, 1)
}

func (g *Graphics) HandleClick(location F2, button bool) {

}

func (g *Graphics) MakeViewport(location F2) {
	gl.Viewport(int32(location.X), int32(location.Y), 400, 200)
	gl.Ortho(0, 400, 200, 0, 12, 128)
	gl.Rectf(0, 0, 400, 200)
}

func (g *Graphics) Cleanup() {
	err := g.window.Destroy()
	if err != nil {
		panic(err)
	}
	err = g.renderer.Destroy()
	sdl.GLDeleteContext(g.ctx)
	if err != nil {
		panic(err)
	}
}

func (g *Graphics) Push() {
	gl.PushMatrix()
}

func (g *Graphics) Pop() {
	gl.PopMatrix()
}

func (g *Graphics) Translate(f3 F3) {
	gl.Translatef(float32(f3.X), float32(f3.Y), float32(f3.Z))
}
