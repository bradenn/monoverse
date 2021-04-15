package simulation

import (
	"fmt"
	"github.com/bradenn/monoverse/physics"
)

type Medium struct {
	X, Y, Z float64
	w       float64
}

func (m *Medium) GenerateSubtree(index int) (e *Octree) {

	a := m.w / 4
	b := m.w / 4
	c := m.w / 4

	if index > 4 {
		a = -a
	}

	if index%2 != 0 {
		b = -b
	}

	// We need alternating numbers now, like 00110011
	// We can use the decimal place of an repeating value...
	// So anything divided by 11 will be multiplied by nine and repeated in the decimal place.
	// So 5/11 = 0.454545, We need to get 1100 to repeat. So to find that, (1100/9)/1111
	// Math is weird... (1100/9)/1111 = 0.11001100

	swap := (1100.0 / 9.0) / 1111.0
	swstr := fmt.Sprintf("%.10f", swap)
	if swstr[4:][index] == '0' {
		c = -c
	}

	medium := &Medium{
		X: m.X + a,
		Y: m.Y + b,
		Z: m.Z + c,
		w: m.w / 2,
	}

	return &Octree{
		medium: medium,
	}
}

func (m *Medium) Contains(body *physics.GBody) bool {
	xBound := body.Location.X >= (m.X-m.w/2) && body.Location.X <= (m.X+m.w/2)
	yBound := body.Location.Y >= (m.Y-m.w/2) && body.Location.Y <= (m.Y+m.w/2)
	zBound := body.Location.Z >= (m.Z-m.w/2) && body.Location.Z <= (m.Z+m.w/2)
	if xBound && yBound && zBound {
		return true
	} else {
		return false
	}
}
