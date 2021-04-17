package main

import (
	"github.com/bradenn/monoverse/graphics"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	bx = 2
	by = 2
)

type Frame struct {
	Renderer      *sdl.Renderer
	Views         []*View
	cols, rows    int
	width, height float64
}

func NewFrame(eng *graphics.Engine, width float64, height float64, cols int, rows int) (f *Frame) {

	return &Frame{
		cols:     cols,
		rows:     rows,
		width:    width,
		height:   height,
		Renderer: eng.Renderer,
		Views:    []*Views{},
	}

}

func (f *Frame) Clear() {
	err := f.Renderer.Clear()
	if err != nil {
		return
	}
}

func (f *Frame) HandleClick(location graphics.F2) {

	for _, pane := range f.Panes {
		if pane != nil {
			pane.Contains(location)
		}
	}
}
func (f *Frame) Render() {

	for _, pane := range f.Panes {
		if pane != nil {
			pane.Render()
		}
	}
	_ = f.Renderer.SetDrawColor(0, 0, 0, 255)
	f.Renderer.Present()
}
