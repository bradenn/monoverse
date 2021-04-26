package main

type Matter struct {
	location, velocity, force F3
	mass, density, volume     float64
}

func (m *Matter) GetMass() float64 {
	return m.mass
}
func (m *Matter) Draw(g *Graphics) {
	g.Color(1, 1, 1, 1)
	g.Tet(m.location, F3{2, 2, 2})
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
