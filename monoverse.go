package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"sync"
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

	verse := NewVerse("Verse", m.grid.GetLocation(F2{0, 0}), m.grid.GetSize(F2{10, 32}))
	m.AddView(verse)

	list := NewList("Scene Statistics", m.grid.GetLocation(F2{10, 30}), m.grid.GetSize(F2{2, 7}))
	verse.stats = &list
	m.AddView(&list)

	physics := &Physics{
		location: m.grid.GetLocation(F2{10, 6}),
		size:     m.grid.GetSize(F2{2, 24}),
	}
	verse.physics = physics
	m.AddView(physics)

	timer := &Timer{
		bounds: m.grid.GetBounds(F2{0, 33}, F2{10, 4}),
	}
	verse.timer = timer
	// m.AddView(timer)

	widget := &fpsWidget{
		location:    m.grid.GetLocation(F2{10, 0}),
		size:        m.grid.GetSize(F2{2, 6}),
		updateDelta: 0,
		renderDelta: 0,
	}
	m.AddView(widget)
}

// type Timer interface {
// 	Configure()
// 	TickCommence(float64)
// 	TickConclude(float64)
// }

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

	m.Configure()

	timer := &Timer{
		bounds: m.grid.GetBounds(F2{0, 33}, F2{10, 4}),
	}
	m.AddView(timer)

	m.running = true

	updateClock := 0.0

	renderInterval := 1000.0 / 60.0 // ms
	renderClock := 0.0

	delta := 0.0
	prevTime := 0.0

	for m.running {

		delta = float64(sdl.GetTicks()) - prevTime
		prevTime = float64(sdl.GetTicks())

		updateClock += delta
		renderClock += delta

		wg := new(sync.WaitGroup)
		// Since the duration of the update cycle can range from nanoseconds to hours or days,
		// so instead of a fixed frame rate, we run this for loop for as long as a rendered frame is not due.
		for renderClock < renderInterval {
			for i := 0; i < 1; i++ {
				wg.Add(len(m.views))
				for _, view := range m.views {
					go func(v View) {
						v.Update()
						wg.Done()
					}(view)
				}
				wg.Wait()
			}
			renderClock += float64(sdl.GetTicks()) - prevTime
		}
		if renderClock > renderInterval {
			m.HandleEvents()
			m.graphics.Clear()
			for _, view := range m.views {
				m.graphics.RenderView(view)
			}
			m.graphics.Render()
			renderClock = 0
		}
	}
}

func (m *Monoverse) HandleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		for _, view := range m.views {
			view.HandleEvent(event)
		}
		switch t := event.(type) {
		case *sdl.QuitEvent:
			m.running = false
			break
		case *sdl.MouseButtonEvent:

		case *sdl.MouseMotionEvent:

			break
		case *sdl.MouseWheelEvent:
		case *sdl.WindowEvent:
			if t.Type == sdl.WINDOWEVENT_RESIZED {
				// Some day
			}
			break
		}
	}

}
