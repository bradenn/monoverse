package main

import "github.com/go-gl/gl/v2.1/gl"

type Matter struct {
	location, velocity, force     F3
	mass, density, volume, charge float64
}

func (m *Matter) GetCharge() float64 {
	return m.charge
}

func (m *Matter) SetCharge(f float64) {
	m.charge = f
}

func (m *Matter) GetMass() float64 {
	return m.mass
}

func (m *Matter) Draw(g *Graphics) {
	gl.PushMatrix()

	g.Color(
		MapF2(m.GetCharge(), F2{-1e-4, 1e-4}, F2{0.8, 0}),
		MapF2(m.GetCharge(), F2{-20000, 20000}, F2{0, 1}),
		MapF2(m.GetCharge(), F2{-20000, 20000}, F2{0, 1}),
		1)
	sz := MapF2(m.GetMass(), F2{0, 6e14}, F2{1, 16})
	gl.Translatef(float32(m.location.X), float32(m.location.Y), float32(m.location.Z))
	g.Cube(F3{0, 0, 0}, F3{sz, sz, sz})

	// gl.Begin(gl.LINES)
	// gl.Vertex3f(0, 0, 0)
	// gl.Vertex3f(float32(MapF2(m.force.X, F2{0, 1e-6}, F2{0, 1})), float32(MapF2(m.force.Y, F2{0, 1e-6}, F2{0, 1})),
	// 	float32(MapF2(m.force.Z, F2{0, 1e-6}, F2{0, 1})))
	// gl.End()
	gl.PopMatrix()
}

func (m *Matter) SetMass(f float64) {
	m.mass = f
}

func (m *Matter) GetDensity() float64 {
	return m.density
}

func (m *Matter) SetDensity(f float64) {
	m.density = f
}

func (m *Matter) GetVolume(f float64) {
	m.volume = f
}

func (m *Matter) SetVolume() float64 {
	return m.volume
}

func (m *Matter) GetLocation() F3 {
	return m.location
}

func (m *Matter) GetVelocity() F3 {
	return m.velocity
}

func (m *Matter) GetForce() F3 {
	return m.force
}

func (m *Matter) SetLocation(f F3) {
	m.location = f
}

func (m *Matter) SetVelocity(f F3) {
	m.velocity = f
}

func (m *Matter) SetForce(f F3) {
	m.force = f
}
