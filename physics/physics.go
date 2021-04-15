package physics

import (
	"math"
)

const (
	G = 1
	D = 2
)

type VectorBody interface {
}

type Vector3D struct {
	X float64
	Y float64
	Z float64
}

type GBody struct {
	Location     Vector3D
	Velocity     Vector3D
	Acceleration Vector3D
	Mass         float64
}

func NewBody(x float64, y float64, z float64, m float64) *GBody {
	g := &GBody{
		Location:     Vector3D{x, y, z},
		Velocity:     Vector3D{0, 0, 0},
		Acceleration: Vector3D{0, 0, 0},
		Mass:         m,
	}
	return g
}

func (g *GBody) Draw() (x float64, y float64, z float64, m float64) {
	return g.Location.X, g.Location.Y, g.Location.Z, g.Mass
}

func (g *GBody) ResetForces() {
	g.Acceleration.X = 0
	g.Acceleration.Y = 0
	g.Acceleration.Z = 0
}

func (g *GBody) UpdatePosition(dt float64) {

	g.Velocity.X += dt * g.Acceleration.X / g.Mass
	g.Velocity.Y += dt * g.Acceleration.Y / g.Mass
	g.Velocity.Z += dt * g.Acceleration.Z / g.Mass

	g.Location.X += dt * g.Velocity.X
	g.Location.Y += dt * g.Velocity.Y
	g.Location.Z += dt * g.Velocity.Z

}

func (g *GBody) DistanceTo(d Vector3D) float64 {
	dist := math.Sqrt(
		math.Pow(d.X-g.Location.X, 2) +
			math.Pow(d.Y-g.Location.Y, 2) +
			math.Pow(d.Z-g.Location.Z, 2))
	if dist != 0 {
		return dist
	} else {
		return 0
	}
}

func (g *GBody) AddForce(d GBody) {
	if *g != d {
		ds := d.DistanceTo(g.Location)
		df := (G * g.Mass * d.Mass) / (math.Pow(ds, 2) + math.Pow(D, 2))
		g.Acceleration.X += df * ((d.Location.X - g.Location.X) / ds)
		g.Acceleration.Y += df * ((d.Location.Y - g.Location.Y) / ds)
		g.Acceleration.Z += df * ((d.Location.Z - g.Location.Z) / ds)
	}
}

func (g *GBody) AddBody(d GBody) {
	ms := g.Mass + d.Mass
	g.Location.X = (g.Location.X*g.Mass + d.Location.X*d.Mass) / ms
	g.Location.Y = (g.Location.Y*g.Mass + d.Location.Y*d.Mass) / ms
	g.Location.Z = (g.Location.Z*g.Mass + d.Location.Z*d.Mass) / ms
	g.Mass = ms
}
