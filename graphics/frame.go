package graphics

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Frame struct {
	Renderer      *sdl.Renderer
	Panes         []*Pane
	cols, rows    int
	width, height float64
}

func NewFrame(eng *Engine, width float64, height float64, cols int, rows int) (f *Frame) {

	return &Frame{
		cols:     cols,
		rows:     rows,
		width:    width,
		height:   height,
		Renderer: eng.Renderer,
		Panes:    []*Pane{NewPane(eng, width/2, height/2)},
	}

}
