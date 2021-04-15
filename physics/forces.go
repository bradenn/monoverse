package physics

type Forcible interface {
	GetLocation() Vector3D
	GetAcceleration() Vector3D
	GetVelocity() Vector3D
	GetMass() Vector3D

	SetLocation(Vector3D)
	SetAcceleration(Vector3D)
	SetVelocity(Vector3D)
	SetMass(Vector3D)
}

func EmplaceForce(s *Forcible, t *Forcible) {

}
