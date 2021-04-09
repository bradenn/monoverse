package main

import "github.com/veandco/go-sdl2/sdl"

type Engine struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Canvases []*Canvas
	Width    int32
	Height   int32
	Running  bool
}

func NewWindow(width int32, height int32) {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	window, err := sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, width, height,
		sdl.WINDOW_ALLOW_HIGHDPI)
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)

	_ = renderer.SetScale(2, 2)
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
