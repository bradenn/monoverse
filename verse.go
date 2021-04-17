package main

import (
	"fmt"
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
	verses []*Verse
	gfx    *Graphics
}

// Verse
// the verse object
// Contains a shit load of hex quadrants on a grid
//
//
type Verse struct {
	ticks  F1
	radius F1
	matter *[]Matter
}

func (v *Verse) computeIteration() (err error) {

	return

}

type Entropy struct {
	ticks    F1
	complete B
	step     F1
	relative F1
	start    time.Time
}

func NewEntropy() *Entropy {
	entropy := &Entropy{ticks: 0.0, step: 1.0, start: time.Now()}

	return entropy
}

func (e *Entropy) Print() {
	fmt.Printf("%f ticks, %f seconds", e.ticks, e.relative)
}

func (e *Entropy) Available() B {
	if e.complete {
		e.complete = false
		return true
	} else {
		fmt.Errorf("%s", "Tick Courrupted")
	}
	return true
}

func (e *Entropy) Ticks() F1 {
	return e.ticks
}

func (e *Entropy) Tick() {
	e.ticks += e.step
	e.complete = true
	e.relative = F1(time.Since(e.start).Seconds())
}

type Matter struct {
	position, velocity, force F3
	mass, density             F1
}

func sdds() {
	entropy := NewEntropy()

	for entropy.Available() {

		entropy.Tick()
	}
	entropy.Print()

}
