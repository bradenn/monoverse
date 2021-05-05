package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"math/rand"
	"sync"
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
		mass:     100,
		density:  0,
		volume:   0,
		charge:   1,
	})
	fAdj := MapF2(force, F2{0e-6, 0e-5}, F2{0, 2})
	fmt.Println(fAdj)
	graphics.Circle(object.GetLocation(), F3{fAdj, fAdj, fAdj})
}

func (g *Gravity) GetConstant() float64 {
	return 6.674e-11
}

func (g *Gravity) Apply(n Object, m Object) float64 {
	distance := Distance(n, m)
	numerator := g.GetConstant() * n.GetMass() * m.GetMass()
	force := numerator / (math.Pow(distance, 2) + math.Pow(D, 2))
	return force
}

type Electromagnetic struct{}

func (g *Electromagnetic) Draw(object Object, graphics *Graphics) {

}

func (g *Electromagnetic) GetConstant() float64 {
	return 8.98e9
}

func (g *Electromagnetic) Apply(n Object, m Object) float64 {
	distance := Distance(n, m)
	numerator := g.GetConstant() * n.GetCharge() * m.GetCharge()
	force := numerator / (math.Pow(distance, 2) + math.Pow(D, 2))
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
	D = 2e1
)

type Physics struct {
	location, size F2
	theta          float64
	forces         []Force
	stack          *SlipStack
	matter         []Object
	octree         *Octree
	ticks          float64
	radius         float64
	delta          float64
}

// The Tick
// This is the smallest unit of progression in a multiverse
func (p *Physics) Tick() {
	p.stack.Start()
	p.ticks += p.delta
	p.stack.Flag("Octree")
	p.getFurthestPoint()
	p.octree = &Octree{voxel: Voxel{
		location: F3{0, 0, 0},
		size:     F3{p.radius, p.radius, p.radius},
	}}

	for _, object := range p.matter {
		p.ResetForces(object)
		p.octree.Insert(object)
	}

	p.stack.Flag("Forces")
	wg := new(sync.WaitGroup)
	wg.Add(len(p.matter))
	for _, object := range p.matter {
		go func(obj Object) {
			p.octree.ApplyForces(p, obj)
			wg.Done()
		}(object)
	}
	wg.Wait()

	p.stack.Flag("Apply")
	for _, object := range p.matter {
		loc, vel := p.UpdatePosition(object, p.ticks)
		object.SetLocation(loc)
		object.SetVelocity(vel)
	}
	p.stack.Flag("End")
	p.stack.End()
}

func (p *Physics) AddForce(n Object, m Object) {

	ds := Distance(n, m)
	delta := SubF3(m.GetLocation(), n.GetLocation())
	for _, force := range p.forces {
		net := force.Apply(n, m)
		n.SetForce(AddF3(n.GetForce(), F3{
			net * (delta.X / ds),
			net * (delta.Y / ds),
			net * (delta.Z / ds)}))
	}

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
	// p.forces = append(p.forces, &Gravity{})
	p.forces = append(p.forces, &Electromagnetic{})
	p.stack = new(SlipStack)

	p.delta = 0.01
	p.radius = 4096
	diam := p.radius
	// 1 lightyear = 1 pixel
	for x := 0.0; x < 2048; x += 1 {
		matter := new(Matter)
		matter.mass = ((5.9e30) * rand.Float64()) * 5.291005291e-17
		matter.charge = 250 - 500*rand.Float64()
		matter.location = F3{diam/2 - rand.Float64()*(diam), diam/2 - rand.Float64()*(diam),
			diam/2 - rand.Float64()*(diam)}
		p.matter = append(p.matter, matter)
	}
	p.ticks = 0
	// panic("implement me")
}

func (p *Physics) Update() {
	// panic("implement me")
}

func (p *Physics) getFurthestPoint() {
	furthest := 0.0
	for _, object := range p.matter {
		if Distance(object, &Matter{location: F3{}}) > furthest {
			furthest = Distance(object, &Matter{location: F3{}})
		}
	}
	p.radius = furthest
}

func (p *Physics) DrawForces(object Object, g *Graphics) {
	for _, force := range p.forces {
		force.Draw(object, g)
	}
}

func (p *Physics) Draw(g *Graphics) {
	list := NewList("", F2{}, p.size)
	list.AddItem(NewItem("Verse", "+"))
	list.AddItem(NewItem("    Matter", fmt.Sprintf("%d OBJ", len(p.matter))))
	list.AddItem(NewItem("    Radius", fmt.Sprintf("%.2f LY", p.radius)))
	list.AddItem(NewItem("Emergent", "+"))
	list.AddItem(NewItem("    Ticks", fmt.Sprintf("%.4f", p.ticks)))
	list.AddItem(NewItem("    Delta", fmt.Sprintf("%.4f", p.delta)))
	list.AddItem(NewItem("Forces Applied", "+"))
	list.AddItem(NewItem("    Gravity", "6.67E-11"))
	list.AddItem(NewItem("    Electromagnetic", "8.98E9"))
	list.AddItem(NewItem("    Strong", "OFF"))
	list.AddItem(NewItem("    Weak", "OFF"))
	list.Draw(g)
}

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

func Distance(n Locatable, m Locatable) float64 {
	diffs := SubF3(m.GetLocation(), n.GetLocation())
	dist := math.Pow(diffs.X, 2) +
		math.Pow(diffs.Y, 2) +
		math.Pow(diffs.Z, 2)
	sqrt := math.Sqrt(dist)
	return sqrt
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

func (p *Physics) DrawMap(g *Graphics) {

}

func (p *Physics) DrawMatter(g *Graphics) {
	for _, object := range p.matter {
		object.Draw(g)
	}
}
