package simulation

import (
	"math/rand"
)

type Body struct {
	X, Y, Z    float64
	M          float64
	Fx, Fy, Fz float64
	Vx, Vy, Vz float64
}

func (Body) Random() *Body {
	diam := 1024.0
	body := &Body{
		X:  diam/2 - float64(rand.Float32())*diam,
		Y:  diam/2 - float64(rand.Float32())*diam,
		Z:  0,
		M:  float64(rand.Float32()) * 1e10,
		Fx: 0,
		Fy: 0,
		Fz: 0,
		Vx: 0,
		Vy: 0,
		Vz: 0,
	}
	return body
}

func (s *Body) Draw() (x float64, y float64, z float64, m float64) {
	return s.X, s.Y, s.Z, s.M
}

func (s *Body) ClearForces() {
	s.Fx = 0
	s.Fy = 0
	s.Fz = 0
}

func (s *Body) AddForce(t Body) {
	// // ds := physics.DistanceBetween(*s, t)
	// df := (physics.G * s.M * t.M) / (math.Pow(ds, 2))
	// t.Fx += df*((s.X-t.X)/ds)
	// t.Fy += df*((s.Y-t.Y)/ds)
	// t.Fz += df*((s.Z-t.Z)/ds)
}

func (s *Body) UpdatePosition(dt float64) {
	s.SetVelocity(s.Vx+dt*s.Fx/s.M, s.Vy+dt*s.Fy/s.M, s.Vz+dt*s.Fz/s.M)
	s.SetPosition(s.X+dt*s.Vx, s.Y+dt*s.Vy, s.Z+dt*s.Vz)
}

func (s *Body) AddForces(x float64, y float64, z float64) {
	s.Fx += x
	s.Fy += y
	s.Fz += z
}

func (s *Body) SetPosition(x float64, y float64, z float64) {
	s.X = x
	s.Y = y
	s.Z = z
}

func (s *Body) SetVelocity(x float64, y float64, z float64) {
	s.Vx = x
	s.Vy = y
	s.Vz = z
}

func (s *Body) print() {

}
