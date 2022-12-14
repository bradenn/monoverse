package main

import (
	"math"
)

type Force interface {
	Apply(n Object, m Object) float64
	Draw(object Object, graphics *Graphics)
}

type Gravity struct{}

func (g *Gravity) Draw(object Object, graphics *Graphics) {
	force := g.Apply(object, &Matter{
		location: object.GetLocation(),
		velocity: F3{},
		force:    F3{},
		mass:     1,
		density:  0,
		volume:   0,
		charge:   1,
		radius:   0,
	})
	fAdj := MapF2(force, F2{0e-6, 0e-5}, F2{0, 200})

	graphics.Circle(object.GetLocation(), F3{fAdj, fAdj, fAdj})
}

func (g *Gravity) GetConstant() float64 {
	return 1 / 16
}

func (g *Gravity) Apply(n Object, m Object) float64 {
	distance := Distance(n, m)
	numerator := g.GetConstant() * n.GetCharge() * m.GetCharge()
	force := numerator / (math.Pow(distance, 2))
	return force
}

type Electromagnetic struct{}

func (g *Electromagnetic) Draw(object Object, graphics *Graphics) {

}

func (g *Electromagnetic) GetConstant() float64 {
	return 0.007299270073
}

func (g *Electromagnetic) Apply(n Object, m Object) float64 {
	distance := Distance(n, m)
	numerator := -g.GetConstant() * n.GetMass() * m.GetMass()
	force := numerator / (math.Pow(distance, 1))
	return force

}

type Strong struct{}

func (g *Strong) Draw(object Object, graphics *Graphics) {

}

func (g *Strong) GetConstant() float64 {
	return 1
}

func (g *Strong) Apply(n Object, m Object) float64 {
	distance := Distance(n, m)

	numerator := g.GetConstant() * n.GetMass() * m.GetMass()
	force := numerator / (math.Pow(distance, 1))
	return force

}
