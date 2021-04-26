package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type Window struct {
	title       string
	location    F2
	perspective F2
	size        F2

	updateLastTick time.Time
	updateDelta    time.Duration

	renderLastTick time.Time
	renderDelta    time.Duration
}

func (w *Window) GetLocation() F2 {
	return w.location
}

func (w *Window) GetSize() F2 {
	return w.size
}

func (w *Window) Configure() {
	w.updateLastTick = time.Now()
	w.renderLastTick = time.Now()
}

func (w *Window) HandleEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.MouseButtonEvent:
		break
	case *sdl.MouseMotionEvent:
		w.perspective.X += float64(t.XRel)
		w.perspective.Y += float64(t.YRel)
		break
	case *sdl.MouseWheelEvent:
		break
	}
}

func (w *Window) DrawFrame(g *Graphics) {
	g.ColorSecondary()
	g.FillRect(F2{1, 1}, F2{w.size.X - 2, 28})
	g.ColorText()
	g.Text(w.title, F2{8, 6})
	g.ColorPrimary()
	g.Rect(F2{1, 1}, F2{w.size.X - 2, w.size.Y - 2})
}
func (w *Window) Update() {
	w.updateDelta = time.Duration(time.Since(w.updateLastTick).Nanoseconds())
	w.updateLastTick = time.Now()
}

func (w *Window) Draw(g *Graphics) {
	w.renderDelta = time.Duration(time.Since(w.renderLastTick).Nanoseconds())
	w.renderLastTick = time.Now()
}

func (w *Window) GetBounds() (F2, F2) {
	return w.location, w.size
}
