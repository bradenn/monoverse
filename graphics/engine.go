package graphics

import (
	"github.com/veandco/go-sdl2/sdl"
)

type B bool
type F1 float64

type F2 struct {
	X float64
	Y float64
}

type F3 struct {
	X float64
	Y float64
	Z float64
}
type Drawable interface {
	Draw() (x float64, y float64, z float64, m float64)
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

	sdl.SetHint("HINT_FRAMEBUFFER_ACCELERATION", "3")
	sdl.SetHint("HINT_RENDER_DRIVER", "metal")
	sdl.SetHint("HINT_RENDER_VSYNC", "0")
	sdl.SetHint("HINT_VIDEO_DOUBLE_BUFFER", "1")
	sdl.SetHint("HINT_OVERRIDE", "1")

	e.Window, err = sdl.CreateWindow("Monoverse v0.1.22 beta", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		width, height, sdl.WINDOW_ALLOW_HIGHDPI)
	e.Renderer, err = sdl.CreateRenderer(e.Window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)

	// High DPI Scale
	_ = e.Renderer.SetLogicalSize(1920, 1080)
	_ = e.Renderer.SetScale(2, 2)
	_ = e.Renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	HandleError(err)
	e.Running = true
	return e, err
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func (e *Engine) Clear() {
	_ = e.Renderer.Clear()

}

func (e *Engine) Render() {
	_ = e.Renderer.SetDrawColor(0, 0, 0, 255)
	e.Renderer.Present()
}

func (e *Engine) KeyInput(key *sdl.KeyboardEvent) {
	switch key.Keysym.Sym {
	case sdl.K_p:
		e.Window.SetSize(e.width+1, e.height-1)
	}
}

func (e *Engine) Cleanup() {
	err := e.Window.Destroy()
	if err != nil {
		panic(err)
	}
	err = e.Renderer.Destroy()
	if err != nil {
		panic(err)
	}
}
