package main

type Particle interface {
	GetMass() float64
	GetCharge() float64
	GetSpin() float64
}

type Quark struct {
	mass   float64
	charge float64
	spin   float64
}

type Lepton struct {
	mass   float64
	charge float64
	spin   float64
}

type Boson struct {
	mass   float64
	charge float64
	spin   float64
}
