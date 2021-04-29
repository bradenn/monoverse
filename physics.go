package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type Force interface {
	Apply(n Object, m Object) float64
}

type Gravity struct{}

func (g *Gravity) GetConstant() float64 {
	return 6.674e-11
}

func (g *Gravity) Apply(n Object, m Object) float64 {
	distance := Distance(n, m)
	numerator := g.GetConstant() * n.GetMass() * m.GetMass()
	force := numerator / (math.Pow(distance, 2) + math.Pow(2, 2))
	return force
}

type Locatable interface {
	GetLocation() F3
	SetLocation(f F3)
}

// The Object
// An object is basically a non-point-sized mass.
// The Object, in theory should be able to represent most of the attributes
// of the Universe we know today.
//
type Object interface {
	GetLocation() F3
	SetLocation(f F3)

	GetVelocity() F3
	SetVelocity(f F3)

	GetForce() F3
	SetForce(f F3)

	GetMass() float64
	SetMass(f float64)

	GetCharge() float64
	SetCharge(f float64)

	GetDensity() float64
	SetDensity(f float64)

	GetVolume(f float64)
	SetVolume() float64
	Draw(g *Graphics)
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
	G = 6.67e-11
	D = 2e2
)

type Physics struct {
	location, size F2
	theta          float64
	forces         []Force
}

func (p *Physics) Tick(delta float64, matter []*Matter) {

}

func (p *Physics) GetName() string {
	return "Physics"
}

func (p *Physics) GetLocation() F2 {
	return p.location
}

func (p *Physics) GetSize() F2 {
	return p.size
}

func (p *Physics) HandleEvent(event sdl.Event) {
	// panic("implement me")
}

func (p *Physics) Configure() {
	p.forces = append(p.forces, &Gravity{})
	// panic("implement me")
}

func (p *Physics) Update() {
	// panic("implement me")
}

func (p *Physics) Draw(g *Graphics) {
	list := NewList("", F2{}, p.size)
	list.AddItem(NewItem("Forces Applied", "+"))
	list.AddItem(NewItem("    Gravity", "6.67E-11 N•m^2/kg^2"))
	list.AddItem(NewItem("    Electromagnetic", "OFF"))
	list.AddItem(NewItem("    Strong", "OFF"))
	list.AddItem(NewItem("    Weak", "OFF"))
	list.Draw(g)
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
//
// func (p *Physics) AddForce(n Object, m Object) {
// 	// distance := p.GetDistance(n, m)
// 	force := SubF3(m.GetForce(), n.GetForce())
//
// 	// divisor := math.Sqrt(math.Pow(force.X, 2) + math.Pow(force.Y, 2) + math.Pow(force.Z, 2))
// 	// numerator := 1.67E-11 * n.GetProperties().X * m.GetProperties().X
// 	// n.SetForce(AddF3(n.GetForce(), F3{numerator * (force.X / divisor), numerator * (force.Y / divisor),
// 	// 	numerator * (force.Z / divisor)}))
//
// }

func (p *Physics) ResetForces(n Object) {
	n.SetForce(F3{0, 0, 0})
}

func (p *Physics) UpdatePosition(n Object, dt float64) (location F3, velocity F3) {
	mass := n.GetMass()

	force := n.GetForce()
	newVelocity := AddF3(n.GetVelocity(), F3{dt * force.X / mass, dt * force.Y / mass, dt * force.Z / mass})

	velocity = newVelocity
	newLocation := AddF3(n.GetLocation(), F3{dt * velocity.X, dt * velocity.Y, dt * velocity.Z})
	return newLocation, newVelocity
}

func Distance(n Object, m Object) float64 {
	diffs := SubF3(m.GetLocation(), n.GetLocation())
	dist := math.Pow(diffs.X, 2) +
		math.Pow(diffs.Y, 2) +
		math.Pow(diffs.Z, 2)
	sqrt := math.Sqrt(dist)
	return sqrt
}

func (p *Physics) AddForce(f Force, n Object, m Object) F3 {
	net := f.Apply(n, m)
	ds := Distance(n, m)
	delta := SubF3(m.GetLocation(), n.GetLocation())
	return AddF3(n.GetForce(), F3{
		net * (delta.X / ds),
		net * (delta.Y / ds),
		net * (delta.Z / ds)})

}

func (p *Physics) AddBody(n Object, m Object) (b Object) {

	massLocationN := MulF3(n.GetLocation(), F3{n.GetMass(), n.GetMass(), n.GetMass()})
	massLocationM := MulF3(m.GetLocation(), F3{m.GetMass(), m.GetMass(), m.GetMass()})

	massSum := AddF3(massLocationN, massLocationM)

	totalMass := n.GetMass() + m.GetMass()

	newLocation := DivF3(massSum, F3{totalMass, totalMass, totalMass})

	b = &Matter{
		location: newLocation,
		velocity: F3{},
		force:    F3{},
		mass:     totalMass,
		density:  0,
		volume:   0,
	}

	return b

}
