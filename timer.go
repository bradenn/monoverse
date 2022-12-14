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
	end   time.Time
	delta time.Duration
	flags []Flag
}

func (s *SlipStack) Start() {
	s.flags = []Flag{}
	s.begin = time.Now()
	s.Flag("Start")
}

func (s *SlipStack) Flag(event string) {
	s.flags = append(s.flags, Flag{
		name:  event,
		time:  time.Now(),
		delta: 0,
	})
}

func (s *SlipStack) Print() {

	for i := 0.0; i < float64(s.delta.Milliseconds()); i++ {

		fmt.Printf("Duration: %s\n", s.delta.String())
	}
}

func (s *SlipStack) End() {
	s.end = time.Now()
	s.delta = s.end.Sub(s.begin)
	s.Flag("End")
}

func (s *SlipStack) Draw(g *Graphics) {
	if len(s.flags) > 0 {
		for i, _ := range s.flags {
			xLocation := MapF2(float64(i), F2{0, float64(len(s.flags) - 1)}, F2{0, 1000})
			g.Line(F3{0, 10, 0}, F3{0, xLocation, 0})
		}
	}

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
	last    []Flag
	top     Flag
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

// item := t.flags[len(t.flags)-1]
// t.flags[len(t.flags)-1].delta = time.Now().Sub(item.time)
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
	t.last = []Flag{}
}

func (t *Timer) Update() {

}

func (t *Timer) Begin(event string) {
	t.flags = []Flag{}
	t.listen = true
	t.Flag(event)
}
func (t *Timer) End() {
	t.UpdateDelta()
	t.flags = append(t.flags, Flag{
		name:  "EXIT",
		time:  time.Now(),
		delta: 1,
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
	t.runningAverage()
}

func (t *Timer) GetTops(event string) {
	t.UpdateDelta()
	t.flags = append(t.flags, Flag{
		name:  event,
		time:  time.Now(),
		delta: 0,
	})
}

func (t *Timer) runningAverage() {
	for i, flag := range t.flags {
		if len(t.last)-1 < i {
			t.last = append(t.last, flag)
		} else {
			t.last[i].delta = (t.last[i].delta + flag.delta) / 2
			nt := new(time.Time)
			t.last[i].time = nt.Add(time.Duration((t.last[i].time.Nanosecond() + flag.time.Nanosecond()) / 2))
		}

	}

}

func (t *Timer) Draw(g *Graphics) {
	if len(t.flags) < 1 {
		return
	}
	t.history = append(t.history, t.flags)
	gl.Enable(gl.BLEND)
	g.Color(1, 1, 1, 1)
	gl.LineWidth(2)
	g.Color(1, 0.5, 0.5, 1)

	track := MapF2(float64(time.Now().Nanosecond()),
		F2{0, 5},
		F2{t.bounds.size.X / 16, t.bounds.size.X - t.bounds.size.X/16})
	g.Line(F3{track, 0, 0}, F3{track, t.bounds.size.Y, 0})

	ht := t.bounds.size.Y / 8
	// aap := float64(t.flags[len(t.flags)-1].time.Sub(t.flags[0].time).Microseconds())
	// float64(t.flags[0].time.Nanosecond()), float64(t.flags[0].time.Nanosecond())
	aap := math.Round(float64(t.InnerDelta().Microseconds())/10000) * 10000
	for c := 0.0; c <= aap; c++ {
		k := MapF2(c,
			F2{0, aap},
			F2{16, t.bounds.size.X - 64})
		gl.LineWidth(1)
		g.Color(0.4, 0.4, 0.4, 0.1)
		if int(c)%1000 == 0 || int(c) == 0 {
			hm := ht * 5
			g.Line(F3{k, t.bounds.size.Y, 0}, F3{k, t.bounds.size.Y - hm, 0})
			g.Color(1, 1, 1, 1)
			g.TextCenter(fmt.Sprintf("%s", time.Duration(c*1000).String()),
				F2{k, (t.bounds.size.Y / 2) - hm/2})
		} else if int(c)%100 == 0 {
			hm := ht * 3
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
	g.Text(fmt.Sprintf("%s", t.InnerDelta().String()), F2{8, 0})
	hold := 16.0
	for _, flag := range t.flags {

		deltaFs := math.Round(float64(flag.delta.Microseconds())/10) * 10
		duration := MapF2(deltaFs,
			F2{0, aap},
			F2{16, t.bounds.size.X - 32})

		g.Color(0.2, 0.2, 0.8, 1)
		g.FillRect(F2{hold, (t.bounds.size.Y) - 64}, F2{duration, 64})
		g.Color(1, 1, 1, 1)
		g.Rect(F2{hold - 1, (t.bounds.size.Y) - 65}, F2{duration + 2, 66})
		g.Line(F3{hold, t.bounds.size.Y, 0}, F3{hold, (t.bounds.size.Y) - 64, 0})
		if duration > 500 {
			g.TextCenter3D(fmt.Sprintf("%s", flag.name), F3{hold + duration/2, (t.bounds.size.Y / 2) + 16, 3})
		}
		hold += duration + 8

	}
	gl.Disable(gl.BLEND)

	t.flags = []Flag{}

}
