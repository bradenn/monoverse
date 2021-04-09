package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	fmt.Printf("BIG BANG! %s\n", time.Since(start))
	sim := Simulation{
		LinearPrecision:   1,
		EnergyPrecision:   1,
		TemporalPrecision: 1,
		UnitLength:        1,
		Radius:            10,
		Iteration:         1,
		Delta:             1.0/26000.0,
	}
	particle := &Particle{
		x:    0,
		y:    0,
		mass: 10,
	}
	sim.particles = append(sim.particles, particle)
	for i := 0; i < 10; i++ {
		sim.Iterate()
	}
	sim.running = true
	sim.Print()
	sim.Run()
	fmt.Printf("Done %s\n", time.Since(start))


}
