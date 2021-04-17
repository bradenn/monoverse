package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Graphics struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	font     *ttf.Font
	frame    *Frame
	size     F3
}

func NewGraphics(size F3) (g *Graphics, err error) {
	g = &Graphics{size: size}
	err = g.configure()

	return g, err
}

func (g *Graphics) configure() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)

	sdl.SetHint("HINT_FRAMEBUFFER_ACCELERATION", "3")
	sdl.SetHint("HINT_RENDER_DRIVER", "metal")
	sdl.SetHint("HINT_RENDER_VSYNC", "0")
	sdl.SetHint("HINT_VIDEO_DOUBLE_BUFFER", "1")
	sdl.SetHint("HINT_OVERRIDE", "1")

	g.window, err = sdl.CreateWindow("Monoverse - v0.12 beta", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(g.size.X), int32(g.size.X), sdl.WINDOW_ALLOW_HIGHDPI)
	g.renderer, err = sdl.CreateRenderer(g.window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)

	// High DPI Scale
	_ = g.renderer.SetLogicalSize(int32(g.size.X), int32(g.size.Y))
	_ = g.renderer.SetScale(2, 2)
	_ = g.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	return err
}

func (g *Graphics) Clear() {
	err := g.renderer.Clear()
	if err != nil {
		return
	}
}

func (g *Graphics) Render() {
	g.renderer.Present()
}

func (g *Graphics) Cleanup() {
	err := g.window.Destroy()
	if err != nil {
		panic(err)
	}
	err = g.renderer.Destroy()
	if err != nil {
		panic(err)
	}
}
