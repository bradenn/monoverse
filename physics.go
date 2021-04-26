package main

import (
	"math"
)

type Locatable interface {
	GetLocation() F3
	SetLocation(f F3)
}

// The Object
// An object is basically a non-point-sized mass.
// The Object, in theory should be able to represent most of the attributes
// of the Universe we know today.
type Object interface {
	GetLocation() F3
	SetLocation(f F3)

	GetVelocity() F3
	SetVelocity(f F3)

	GetForce() F3
	SetForce(f F3)

	GetMass() float64
	SetMass(f float64)

	GetDensity() float64
	SetDensity(f float64)

	GetVolume(f float64)
	SetVolume() float64
}

func ForEachF3(a F3, call func(f float64) float64) float64 {
	return call(a.X) + call(a.Y) + call(a.Z)

}

func AddF3(a F3, b F3) F3 {
	return F3{X: a.X + b.X, Y: a.Y + b.Y, Z: a.Z + b.Z}
}

func SubF3(a F3, b F3) F3 {
	return F3{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

func MulF3(a F3, b F3) F3 {
	return F3{X: a.X * b.X, Y: a.Y * b.Y, Z: a.Z * b.Z}
}

func DivF3(a F3, b F3) F3 {
	return F3{X: a.X / b.X, Y: a.Y / b.Y, Z: a.Z / b.Z}
}

const (
	G = 6.17e-11
	D = 2e2
)

type Physics struct {
}

func (p *Physics) GetDiameter(n Object) float64 {
	return 1
}

// func (p *Physics) Distance(n Locatable, m Locatable) float64 {
// 	difference := SubF3(n.GetLocation(), m.GetLocation())
//
// 	return math.Sqrt(ForEachF3(difference, func(f float64) float64 {
// 		return math.Pow(f, 2)
// 	}))
// }

// func (p *Physics) AddForce(n Object, m Object) {
// 	// // distance := p.GetDistance(n, m)
// 	// force := SubF3(m.GetForce(), n.GetForce())
// 	//
// 	// divisor := math.Sqrt(math.Pow(force.X, 2) + math.Pow(force.Y, 2) + math.Pow(force.Z, 2))
// 	// numerator := 1.67E-11 * n.GetProperties().X * m.GetProperties().X
// 	// n.SetForce(AddF3(n.GetForce(), F3{numerator * (force.X / divisor), numerator * (force.Y / divisor),
// 	// 	numerator * (force.Z / divisor)}))
//
// }

func (p *Physics) ResetForces(n Object) {
	n.SetForce(F3{})
}

func (p *Physics) UpdatePosition(n Object, dt float64) {
	mass := n.GetMass()

	force := n.GetForce()
	newVelocity := AddF3(n.GetVelocity(), F3{dt * force.X / mass, dt * force.Y / mass, dt * force.Z / mass})
	n.SetVelocity(newVelocity)

	velocity := n.GetVelocity()
	newLocation := AddF3(n.GetLocation(), F3{dt * velocity.X, dt * velocity.Y, dt * velocity.Z})
	n.SetLocation(newLocation)
}

func (p *Physics) Distance(n Object, m Object) float64 {
	diffs := SubF3(m.GetLocation(), n.GetLocation())
	dist := math.Sqrt(
		math.Pow(diffs.X, 2) +
			math.Pow(diffs.Y, 2) +
			math.Pow(diffs.Z, 2))
	if dist != 0 {
		return dist
	} else {
		return 0
	}
}

func (p *Physics) AddForce(n Object, m Object) {
	ds := p.Distance(n, m)
	df := (G * n.GetMass() * n.GetMass()) / (math.Pow(ds, 2) + math.Pow(D, 2))
	delta := SubF3(m.GetLocation(), n.GetLocation())
	n.SetForce(AddF3(n.GetForce(), F3{df * (delta.X / ds),
		df * (delta.Y / ds),
		df * (delta.Z / ds)}))
}

func (p *Physics) AddBody(n Object, m Object) {

	massLocationN := MulF3(n.GetLocation(), F3{n.GetMass(), n.GetMass(), n.GetMass()})
	massLocationM := MulF3(m.GetLocation(), F3{m.GetMass(), m.GetMass(), m.GetMass()})

	massSum := AddF3(massLocationN, massLocationM)

	totalMass := n.GetMass() + m.GetMass()
	n.SetMass(totalMass)

	newLocation := DivF3(massSum, F3{totalMass, totalMass, totalMass})
	n.SetLocation(newLocation)

}
