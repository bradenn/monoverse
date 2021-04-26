package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	_ "github.com/nullboundary/glfont"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Graphics struct {
	window   *sdl.Window
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

func (g *Graphics) configure() (err error) {

	err = sdl.Init(sdl.INIT_EVERYTHING)
	err = ttf.Init()

	sdl.SetHint(sdl.HINT_RENDER_DRIVER, sdl.HINT_RENDER_OPENGL_SHADERS)

	err = sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	err = sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 24)

	g.window, _ = sdl.CreateWindow("Monoverse - v0.12 beta", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(g.size.X), int32(g.size.Y), sdl.WINDOW_OPENGL|sdl.WINDOW_ALLOW_HIGHDPI)
	g.window.SetResizable(true)
	g.ctx, _ = g.window.GLCreateContext()

	err = gl.Init()

	g.font = g.loadFont("JetBrainsMono-Medium.ttf", 24)

	gl.Viewport(0, 0, int32(g.size.X*2), int32(g.size.Y*2))
	gl.DepthRange(-2048, 2048)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	gl.Ortho(0, g.size.X, g.size.Y, 0, -2048, 2048)

	gl.MatrixMode(gl.MODELVIEW)

	gl.Enable(gl.DOUBLEBUFFER)

	return nil
}

func (g *Graphics) RenderView(view View) {

	gl.PushMatrix()
	gl.MatrixMode(gl.MODELVIEW)

	gl.Translatef(float32(view.GetLocation().X), float32(view.GetLocation().Y), 0)

	g.DrawFrame(view.GetName(), view)

	// gl.Scissor(int32(view.GetLocation().X + 2)*2, int32(g.size.Y - view.GetLocation().Y)*2,
	// 	int32(view.GetSize().X - 4)*2,
	// 	int32(view.GetSize().Y - 28 - 1)*2)
	location := view.GetLocation()
	size := view.GetSize()
	gl.Scissor(int32(location.X+2)*2, int32(g.size.Y-(location.Y+size.Y))*2, int32(size.X-4)*2, int32(size.Y)*2)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.SCISSOR_TEST)

	view.Draw(g)

	gl.Disable(gl.SCISSOR_TEST)
	gl.Disable(gl.BLEND)

	gl.PopMatrix()

}

func (g *Graphics) DrawFrame(name string, w View) {
	g.ColorSecondary()
	g.FillRect(F2{1, 1}, F2{w.GetSize().X - 2, 28 - 2})
	g.ColorText()
	g.Text(name, F2{8, 6})
	g.ColorPrimary()
	g.Rect(F2{1, 1}, F2{w.GetSize().X - 2, w.GetSize().Y - 2})
	g.ColorPrimary()
	g.Rect(F2{2, 2}, F2{w.GetSize().X - 4, w.GetSize().Y - 4})
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

func (g *Graphics) Color(a, b, c, d float64) {
	gl.Color4f(float32(a), float32(b), float32(c), float32(d))
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
func (g *Graphics) FillRect(n F2, m F2) {
	gl.Rectf(float32(n.X), float32(n.Y), float32(n.X+m.X), float32(n.Y+m.Y))
}

func (g *Graphics) Line(n F3, m F3) {
	gl.Begin(gl.LINES)
	gl.Vertex3f(float32(n.X), float32(n.Y), float32(n.Z))
	gl.Vertex3f(float32(m.X), float32(m.Y), float32(m.Z))
	gl.End()
}

func (g *Graphics) Tet(n F3, m F3) {
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

	gl.Begin(gl.TRIANGLES)
	for i := 0; i < len(cube)/3; i++ {
		gl.Vertex3f(cube[i*3+0], cube[i*3+1], cube[i*3+2])
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

	gl.Color3f(1.0, 1.0, 1.0) // Green
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
