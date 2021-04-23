package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
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

type Monoverse struct {
	graphics *Graphics
	grid     Grid
	views    []View
	running  bool
}

func (m *Monoverse) AddView(v View) {
	if v != nil {
		v.Configure()
		m.views = append(m.views, v)
	}
}

func (m *Monoverse) Configure() {
	w := 1920.0
	h := 1080.0
	m.grid = Grid{
		size: F2{w, h},
		cell: F2{1920 / 12, 28},
	}
	m.graphics, _ = NewGraphics(F3{w, h, 0})
}

type Timer interface {
	Configure()
	TickCommence()
	TickConclude()
}

type UpdateClock struct {
	inner      float64
	innerDelta float64

	outer      float64
	outerDelta float64

	ticks   float64
	history []time.Duration
	ready   bool
}

func (u *UpdateClock) TickCommence(now float64) {
	if u.ready {
		u.ready = false
	} else {
		u.TickConclude(now)
	}
	u.outerDelta = now - u.outer
	u.inner = now
}

func (u *UpdateClock) TickConclude(now float64) {
	u.innerDelta = now - u.inner
	u.outer = now
	u.ticks += 1.0
	u.ready = true
}

func (u *UpdateClock) GetUpsPerSecond() float64 {
	return u.outerDelta
}

func (u *UpdateClock) Configure() {

}

func (m *Monoverse) Run() {
	uc := new(UpdateClock)
	uc.Configure()

	rc := new(UpdateClock)
	rc.Configure()

	m.Configure()
	verse := NewVerse("Verse", m.grid.GetLocation(F2{0, 0}), m.grid.GetSize(F2{10, 38}))
	m.AddView(verse)

	list := NewList("Scene Info", m.grid.GetLocation(F2{10, 4}), m.grid.GetSize(F2{2, 10}))
	list.AddItem(NewItem("MaxRSS", "Anal Leisure"))
	m.AddView(&list)

	widget := &fpsWidget{
		location:    m.grid.GetLocation(F2{10, 0}),
		size:        m.grid.GetSize(F2{2, 5}),
		updateDelta: 0,
		renderDelta: 0,
	}
	m.AddView(widget)

	m.running = true

	logicInterval := 1000.0 / 144.0 // ms
	updateClock := 0.0
	renderClock := 0.0
	delta := 0.0

	prevTime := 0.0

	for m.running {
		uc.TickCommence(float64(sdl.GetTicks()))
		delta = float64(sdl.GetTicks()) - prevTime
		prevTime = float64(sdl.GetTicks())

		updateClock += delta
		for updateClock > logicInterval {

			for _, view := range m.views {
				view.Update()
			}
			updateClock -= logicInterval

		}

		renderClock += delta
		if renderClock > (16.666) {
			rc.TickCommence(float64(sdl.GetTicks()))
			m.HandleEvents()
			m.graphics.Clear()
			m.grid.Draw(m.graphics)
			for _, view := range m.views {
				m.graphics.RenderView(view)
			}
			m.graphics.Render()
			rc.TickConclude(float64(sdl.GetTicks()))
			renderClock = 0
		}
		uc.TickConclude(float64(sdl.GetTicks()))
	}
}

func (m *Monoverse) HandleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		for _, view := range m.views {
			view.HandleEvent(event)
		}
		switch event.(type) {
		case *sdl.QuitEvent:
			m.running = false
			break
		case *sdl.MouseButtonEvent:

		case *sdl.MouseMotionEvent:

			break
		case *sdl.MouseWheelEvent:

			break
		}
	}

}
