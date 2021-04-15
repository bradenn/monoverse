package physics

// SpaceTime
// Locatable Object: x, y, z, vx, vy, vz... Force? no... no mass...
//
// Locatable Mass  : x, y, z, vx, vy, vz, fx, fy, fz, mass
//
// Force Aggregator:
// - Collects all forces acting on body.
// - Applies forces.
// - Where is it?
//
// Iterator:
// Go to each mass, aggregate each of their forces, apply all forces at the same time.
// Check if masses have collided or are going to collide. If so, combine or fragment.
//
type SpaceTime struct {
	spacial  uint8 // I imagine more than 255 spacial dimensions would be overkill...
	temporal uint8 // Ibid. Beyond my comprehension.
}

func (*SpaceTime) ComputeGravity(bodies []*GBody) {

}

func (*SpaceTime) ApplyForce(bodies []*GBody) {

}
