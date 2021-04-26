package main

type Voxel struct {
	location, size F3
}

func (v *Voxel) Draw(g *Graphics) {
	g.Cube(v.location, v.size)
}

func (v *Voxel) SubVoxel(n int) Voxel {

	p := v.size.X / 4
	q := v.size.Y / 4
	r := v.size.Z / 4

	vertices := []F3{
		F3{p, q, r},
		F3{-p, q, r},
		F3{p, q, -r},
		F3{-p, q, -r},
		F3{p, -q, r},
		F3{-p, -q, r},
		F3{p, -q, -r},
		F3{-p, -q, -r},
	}

	return Voxel{
		location: AddF3(v.location, vertices[n]),
		size:     DivF3(v.size, F3{2, 2, 2}),
	}

}

func (v *Voxel) Contains(n F3) bool {

	xBound := n.X >= (v.location.X-v.size.X/2) && n.X <= (v.location.X+v.size.X/2)
	yBound := n.Y >= (v.location.Y-v.size.Y/2) && n.Y <= (v.location.Y+v.size.Y/2)
	zBound := n.Z >= (v.location.Z-v.size.Z/2) && n.Z <= (v.location.Z+v.size.Z/2)

	if xBound && yBound && zBound {
		return true
	} else {
		return false
	}
}

// The Octree
// Each item node in the tree represents
type Octree struct {
	node     Object
	voxel    Voxel
	children []*Octree
}

func (o *Octree) ApplyForces(object *Object) {
	physics := Physics{}
	if o.node != nil {
		if o.children != nil {
			if o.node != *object {
				physics.AddForce(*object, o.node)
			}
		} else {
			theta := o.voxel.size.X / physics.Distance(*object, o.node)
			if theta < 0.9 {
				physics.AddForce(*object, o.node)
			} else {
				for _, child := range o.children {
					child.ApplyForces(object)
				}
			}
		}
	}
}

func (o *Octree) Insert(object Object) {
	// 1. If the node is empty, fill it and birth virgins.
	if o.node == nil {
		o.node = object
		o.Subdivide()
		// 2. If its not empty, mail it to the right kids.
	} else {
		for _, node := range o.children {
			if node.voxel.Contains(object.GetLocation()) {
				node.Insert(object)
			}
		}
	}
	// 3. Make a copy and add it to your own mass and location before mailing it.
	physics := Physics{}
	physics.AddBody(o.node, object)

	// 4. Take a hit!
}

func (o *Octree) Subdivide() {
	if len(o.children) < 8 {
		for i := 0; i < 8; i++ {
			o.voxel.SubVoxel(i)
		}
	}
}

func (o *Octree) ComputeForce(matter []*Object, tick float64) {

}
