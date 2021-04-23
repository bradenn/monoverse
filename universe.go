package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"syscall"
	"time"
)

type mvObject interface {
	GetLocation() F3
	GetRotation() F3
	GetMotion() (F3, F3)
}

type mvSolid interface {
	GetLocation() F3
	GetRotation() F3
	GetMotion() (F3, F3)
	GetComposition() F3
}

type View interface {
	GetLocation() F2
	GetSize() F2
	HandleEvent(event sdl.Event)
	Configure()
	Update()
	Draw(g *Graphics)
}

type fpsWidget struct {
	location F2
	size     F2

	updateLastTick time.Time
	updateDelta    time.Duration
	updateHistory  []time.Duration

	renderLastTick time.Time
	renderDelta    time.Duration
	renderHistory  []time.Duration
}

func (f *fpsWidget) HandleEvent(event sdl.Event) {

}

func (f *fpsWidget) GetLocation() F2 {
	return f.location
}

func (f *fpsWidget) GetSize() F2 {
	return f.size
}

func (f *fpsWidget) Configure() {
	f.updateLastTick = time.Now()
	f.renderLastTick = time.Now()
}

func (f *fpsWidget) Update() {
	f.updateDelta = time.Duration(time.Since(f.updateLastTick).Nanoseconds())
	f.updateHistory = append(f.updateHistory, f.updateDelta)
	if len(f.updateHistory) > 120 {
		f.updateHistory = f.updateHistory[1:]
	}
	f.updateLastTick = time.Now()
}

func (f *fpsWidget) Draw(g *Graphics) {
	f.renderDelta = time.Duration(time.Since(f.renderLastTick).Nanoseconds())
	f.renderHistory = append(f.renderHistory, f.renderDelta)
	if len(f.renderHistory) > 60 {
		f.renderHistory = f.renderHistory[1:]
	}
	f.renderLastTick = time.Now()

	list := NewList("Performance", f.location, f.size)

	renAvg := 0.0
	for _, it := range f.renderHistory {
		renAvg += float64(it.Milliseconds())
	}
	renAvg /= float64(len(f.renderHistory))
	list.AddItem(NewItem("Render", fmt.Sprintf("%.2f FPS", 1000.0/renAvg)))

	upsAvg := 0.0
	for _, it := range f.updateHistory {
		upsAvg += float64(it.Milliseconds())
	}
	upsAvg /= float64(len(f.updateHistory))
	list.AddItem(NewItem("Update", fmt.Sprintf("%.2f UPS", 1000.0/upsAvg)))

	rusage := new(syscall.Rusage)
	syscall.Getrusage(0, rusage)

	list.AddItem(NewItem("MaxRSS", fmt.Sprintf("%.2f MB ", float64(rusage.Maxrss)/1024/1024)))

	list.Draw(g)
}

func (f *fpsWidget) GetBounds() (F2, F2) {
	return f.location, f.size
}
