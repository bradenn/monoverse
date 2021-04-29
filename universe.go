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
	GetName() string
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
	updateHistory  []float64

	renderLastTick time.Time
	renderDelta    time.Duration
	renderHistory  []float64

	rssPrev    float64
	rssHistory []float64
}

func (f *fpsWidget) GetName() string {
	return "Performance"
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

func MapF2(value float64, input F2, output F2) float64 {
	return output.X + (output.Y-output.X)*((value-input.X)/(input.Y-input.X))
}

func (f *fpsWidget) Update() {
	f.updateDelta = time.Duration(time.Since(f.updateLastTick).Nanoseconds())
	f.updateHistory = append(f.updateHistory, f.updateDelta.Seconds())
	if len(f.updateHistory) > 240*2 {
		f.updateHistory = f.updateHistory[1:]
	}
	f.updateLastTick = time.Now()
}

func (f *fpsWidget) Draw(g *Graphics) {

	f.renderDelta = time.Duration(time.Since(f.renderLastTick).Nanoseconds())
	f.renderHistory = append(f.renderHistory, f.renderDelta.Seconds())
	if len(f.renderHistory) > 120*2 {
		f.renderHistory = f.renderHistory[1:]
	}
	f.renderLastTick = time.Now()

	list := NewList("Performance", f.location, f.size)

	renAvg := 0.0
	for _, it := range f.renderHistory {
		renAvg += it * 1000
	}
	renAvg /= float64(len(f.renderHistory))
	item := NewItem("Render", fmt.Sprintf("%.2f FPS", 1000.0/renAvg))

	list.AddItem(item)

	upsAvg := 0.0
	for _, it := range f.updateHistory {
		upsAvg += it * 1000
	}
	upsAvg /= float64(len(f.updateHistory))
	item2 := NewItem("Update", fmt.Sprintf("%.2f UPS", 1000.0/upsAvg))

	list.AddItem(item2)
	rusage := new(syscall.Rusage)

	syscall.Getrusage(0, rusage)
	f.rssHistory = append(f.rssHistory, float64(rusage.Maxrss)/1024/1024)
	if len(f.rssHistory) > 360 {
		f.rssHistory = f.rssHistory[1:]
	}
	f.rssPrev = float64(rusage.Maxrss) / 1024 / 1024

	item3 := NewItem("MaxRSS", fmt.Sprintf("%.2f MB ", float64(rusage.Maxrss)/1024/1024))
	list.AddItem(item3)
	list.AddItem(NewItem("Frame Duration", fmt.Sprintf("%.2f MS ", renAvg)))
	list.AddItem(NewItem("Update Duration", fmt.Sprintf("%.2f MS ", upsAvg)))
	list.Draw(g)
}

func (f *fpsWidget) GetBounds() (F2, F2) {
	return f.location, f.size
}

type Clock interface {
	GetTicks() float64
	GetDelta() float64
	GetRate() float64
	Tick()
}

type PerformanceClock struct {
	ticks    float64
	delta    float64
	rate     float64
	previous float64
}

func (p *PerformanceClock) Tick() {
	p.delta = float64(sdl.GetTicks()) - p.previous
	p.previous = float64(sdl.GetTicks())
	p.rate = 1000 / p.delta
	p.ticks++
}

func (p *PerformanceClock) GetTicks() float64 {
	return p.ticks
}

func (p *PerformanceClock) GetDelta() float64 {
	return p.delta
}

func (p *PerformanceClock) GetRate() float64 {
	return p.rate
}

type Performance struct {
	location, size F2
	clocks         []Clock
}

func (f *Performance) GetName() string {
	return "Performance"
}

func (f *Performance) Configure() {
	f.clocks = append(f.clocks, &PerformanceClock{}, &PerformanceClock{})

}

func (f *Performance) HandleEvent(event sdl.Event) {

}

func (f *Performance) GetLocation() F2 {
	return f.location
}

func (f *Performance) GetSize() F2 {
	return f.size
}

func (f *Performance) Update() {

	f.clocks[0].Tick()
}

func (f *Performance) Draw(g *Graphics) {
	f.clocks[1].Tick()
	list := NewList("Performance", f.location, f.size)
	render := NewItem("Render", fmt.Sprintf("%.2f FPS", f.clocks[1].GetRate()))
	list.AddItem(render)
	update := NewItem("Update", fmt.Sprintf("%.2f UPS", 1000/f.clocks[0].GetRate()))
	list.AddItem(update)
	list.Draw(g)
}
