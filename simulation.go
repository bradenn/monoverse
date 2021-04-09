package main

import (
	"math/rand"
)

// Smallest Traversable Distance
// 1 unit length / precision
// 1 unit / 100 = smallest traversable distance = 0.01 units

// Fastest Velocity Possible
// 1 unit length / 1 iteration
// 1 unit length /

type Particle struct {
	x    float64
	y    float64
	mass float64
}

type Simulation struct {
	LinearPrecision   float64
	EnergyPrecision   float64
	TemporalPrecision float64
	UnitLength        uint

	Radius    float64
	Iteration float64
	Delta     float64
	particles []*Particle
	running   bool
}

func (s *Simulation) Run() {
	for s.running {

	}
}

func (s *Simulation) Iterate() {
	oct := Octree{
		medium: Medium{0, 0, 0, 512},
	}

	for i := 0; i < 1000; i++ {
		oct.insert(Point{rand.Float64() * 512, rand.Float64() * 512, rand.Float64() * 512, 15, rand.Float64() * 1024})

	}
	oct.print(0)
}

func (s *Simulation) Print() {

}

/*
	LinearPrecision  : 	1E32
	TemporalPrecision:  1E4



*/
