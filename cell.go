package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"sync"
)

type I2 struct {
	X, Y int
}
type World struct {
	cells    map[I2]*Cell
	name     string
	location F2
	frame    F2
	size     int
}

func (w *World) GetName() string {
	return w.name
}

func (w *World) GetLocation() F2 {
	return w.location
}

func (w *World) HandleEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.MouseMotionEvent:
		if t.State == 1 {
			fmt.Println(t.X)
		}
		break
	}
}

func (w *World) Configure() {
	w.Populate()
}

func NewWorld(name string, size int, location F2, frame F2) *World {
	return &World{
		cells:    make(map[I2]*Cell),
		size:     size,
		name:     name,
		location: location,
		frame:    frame,
	}
}

func (w *World) Populate() {
	for y := 0; y < w.size; y++ {
		for x := 0; x < w.size; x++ {
			w.cells[I2{x, y}] = &Cell{false, 0}
		}
	}
	for y := 0; y < 500; y++ {
		w.cells[I2{rand.Int() % w.size, rand.Int() % w.size}] = &Cell{true, 0}
	}

}

func (w *World) GetSize() F2 {
	return w.frame
}

func (w *World) GetNeighbours(n I2) int {
	tally := 0
	points := [][]int{{-1, 0}, {1, 0}, {1, -1}, {0, 1}, {-1, 1}, {0, -1}}
	for _, p := range points {
		cell := w.cells[I2{n.X + p[0], n.Y + p[1]}]
		if cell != nil {
			if cell.alive {
				tally++
			}
		}

	}

	return tally
}

func (w *World) Draw(g *Graphics) {
	g.Line(F3{0, w.frame.Y / 2, 0}, F3{w.frame.X, w.frame.Y / 2, 0})
	g.Line(F3{w.frame.X / 2, 0, 0}, F3{w.frame.X / 2, w.frame.Y, 0})

	gl.Translatef(float32(w.frame.X/2)-float32(w.size)*0.8660254038,
		float32(w.frame.Y/2)-(float32(w.size)*1.5)/2, 0)
	for u, cell := range w.cells {
		offset := 0.0
		if u.Y%2 == 0 {
			offset = 0.5
		}
		if cell.alive {
			cx := (float64(u.X) + offset) * 2 * 0.8660254038
			cy := float64(u.Y) * 1.5
			g.Color(0.5, 0.5, 1, 1)
			g.Hexagon(F2{cx, cy})
		}
	}

}

func (w *World) Update() {

	d := w.cells
	wg := new(sync.WaitGroup)
	wg.Add(len(w.cells))
	for loc, current := range w.cells {
		go func(cell *Cell, delta map[I2]*Cell, u I2) {
			n := w.GetNeighbours(u)
			delta[u].n = n
			if cell.alive {
				if n <= 1 {
					delta[u].alive = false
				} else if n == 2 {
					delta[u].alive = true
				} else if n >= 3 {
					delta[u].alive = false
				}
			} else if n == 2 {
				delta[u].alive = true
			} else {
				delta[u].alive = false
			}
			wg.Done()
		}(current, d, loc)

	}
	wg.Wait()
	w.cells = d
}

type Cell struct {
	alive bool
	n     int
}
