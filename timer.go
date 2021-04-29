package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"time"
)

// A SlipStack is a stack of pairs.
//  | |
//  | |
type SlipStack struct {
	begin time.Time
	flags []Flag
}

type Bounds struct {
	location F2
	size     F2
}

type Flag struct {
	name  string
	time  time.Time
	delta time.Duration
}

func (f *Flag) GetName() string {
	return f.name
}

type Timer struct {
	bounds  Bounds
	flags   []Flag
	history [][]Flag
	ticks   float64
	listen  bool
}

func (t *Timer) GetName() string {
	return "Timer"
}

func (t *Timer) GetLocation() F2 {
	return t.bounds.location
}

func (t *Timer) InnerDelta() time.Duration {
	if len(t.flags) >= 2 {
		return t.flags[len(t.flags)-1].time.Sub(t.flags[0].time)
	}
	return time.Duration(0)
}

func (t *Timer) GetSize() F2 {
	return t.bounds.size
}

func (t *Timer) GetBounds() Bounds {
	return t.bounds
}

func (t *Timer) HandleEvent(event sdl.Event) {

}

func (t *Timer) Configure() {
	t.flags = []Flag{}
}

func (t *Timer) Update() {

}

func (t *Timer) Begin(event string) {
	t.flags = []Flag{}
	t.listen = true
	t.Flag(event)
}
func (t *Timer) End() {
	t.listen = false
	t.UpdateDelta()
	t.flags = append(t.flags, Flag{
		name:  "EXIT",
		time:  time.Now(),
		delta: 0,
	})
}
func (t *Timer) UpdateDelta() {
	if len(t.flags) > 0 {
		item := t.flags[len(t.flags)-1]
		t.flags[len(t.flags)-1].delta = time.Now().Sub(item.time)
	}
}
func (t *Timer) Flag(event string) {
	t.UpdateDelta()
	t.flags = append(t.flags, Flag{
		name:  event,
		time:  time.Now(),
		delta: 0,
	})
}

func (t *Timer) Draw(g *Graphics) {

	if len(t.flags) < 1 {
		return
	}
	t.history = append(t.history, t.flags)
	gl.Enable(gl.BLEND)
	g.Color(1, 1, 1, 1)
	g.Text(fmt.Sprintf("Ticks: %.2f", t.ticks), F2{})
	g.Line(F3{0, t.bounds.size.Y / 2, 0}, F3{t.bounds.size.X, t.bounds.size.Y / 2, 0})
	// begin := t.flags[0].time
	// end := t.flags[len(t.flags)-1].time
	gl.LineWidth(2)
	g.Color(1, 0.5, 0.5, 1)

	track := MapF2(float64(time.Now().Nanosecond()),
		F2{0, 5},
		F2{t.bounds.size.X / 16, t.bounds.size.X - t.bounds.size.X/16})
	g.Line(F3{track, 0, 0}, F3{track, t.bounds.size.Y, 0})

	ht := t.bounds.size.Y / 8
	// aap := float64(t.flags[len(t.flags)-1].time.Sub(t.flags[0].time).Microseconds())
	// float64(t.flags[0].time.Nanosecond()), float64(t.flags[0].time.Nanosecond())
	aap := math.Round(float64(t.InnerDelta().Microseconds())/5000) * 5000
	for c := 0.0; c <= aap; c++ {
		k := MapF2(c,
			F2{0, aap},
			F2{16, t.bounds.size.X - 32})
		gl.LineWidth(1)
		g.Color(0.4, 0.4, 0.4, 0.1)
		if int(c)%1000 == 0 {
			hm := ht * 5
			g.Line(F3{k, t.bounds.size.Y, 0}, F3{k, t.bounds.size.Y - hm, 0})
			g.Color(1, 1, 1, 1)
			g.TextCenter(fmt.Sprintf("%s", time.Duration(c*1000).String()),
				F2{k, (t.bounds.size.Y / 2) - hm/2})
		} else if int(c)%100 == 0 {
			hm := ht * 3
			g.Line(F3{k, t.bounds.size.Y, 0}, F3{k, t.bounds.size.Y - hm, 0})
		} else if int(c)%50 == 0 {
			hm := ht * 2
			g.Line(F3{k, t.bounds.size.Y, 0}, F3{k, t.bounds.size.Y - hm, 0})
		} else if int(c)%10 == 0 {
			hm := ht / 2
			g.Line(F3{k, t.bounds.size.Y, 0}, F3{k, t.bounds.size.Y - hm, 0})
		}
	}

	gl.Enable(gl.BLEND)
	//
	// for i := 0.0; i < float64(t.flags[len(t.flags)-1].time.Sub(t.flags[0].time).Microseconds())/8; i+=8 {
	// 	// g.Line(F3{s+i*sp, (t.bounds.size.Y / 2) - 2, 0}, F3{s+i*sp, (t.bounds.size.Y / 2) - 16, 0})
	// 	for j := 0; j < 8; j++ {
	// 	g.Rect(F2{i*8, float64(j * 8)}, F2{8, 8})
	// 	}
	//
	// }
	g.TextCenter3D(fmt.Sprintf("%s", t.InnerDelta().Microseconds()), F3{20, 20, 0})

	py := t.bounds.size.X / 16
	for _, flag := range t.flags {

		y := MapF2(float64(flag.time.Sub(t.flags[0].time).Microseconds()),
			F2{0, float64(t.InnerDelta().Microseconds())},
			F2{t.bounds.size.X / 16, t.bounds.size.X - t.bounds.size.X/16})

		g.Color(1, 0, 1, 1)
		g.FillRect(F2{py, (t.bounds.size.Y / 2) - 16}, F2{y, 32})
		g.Color(1, 1, 0, 1)
		g.Line(F3{py, (t.bounds.size.Y / 2) + 16, 0}, F3{py, (t.bounds.size.Y / 2) - 16, 0})
		g.Text(fmt.Sprintf("%s", flag.name), F2{py - 8, t.bounds.size.Y / 2})
		py += y
	}
	gl.Disable(gl.BLEND)

	t.flags = []Flag{}

}
