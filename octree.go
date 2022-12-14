package main

type Voxel struct {
	location, size F3
}

func (v *Voxel) Draw(g *Graphics) {
	g.Color(0.5, 0.5, 0.5, 1)
	g.Tet(v.location, v.size)
}

func (v *Voxel) SubVoxel(n int) *Octree {

	p := v.size.X / 4
	q := v.size.Y / 4
	r := v.size.Z / 4

	vertices := make([]F3, 0)
	vertices = append(vertices,
		F3{p, q, r},
		F3{-p, q, r},
		F3{p, q, -r},
		F3{-p, q, -r},
		F3{p, -q, r},
		F3{-p, -q, r},
		F3{p, -q, -r},
		F3{-p, -q, -r})
	return &Octree{voxel: Voxel{
		location: AddF3(v.location, vertices[n]),
		size:     DivF3(v.size, F3{2, 2, 2}),
	}}

}

func (v *Voxel) Contains(b Locatable) bool {
	n := b.GetLocation()
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

func (o *Octree) ApplyForces(physics *Physics, object Object) {
	if o.node != nil {
		if o.children == nil { // External Node
			if o.node != object {
				physics.AddForce(object, o.node)
			}
		} else { // Internal Node
			theta := (o.voxel.size.X) / Distance(o.node, object)
			if theta < 1.2 {
				physics.AddForce(object, o.node)
			} else {
				for _, child := range o.children {
					child.ApplyForces(physics, object)
				}
			}
		}
	}
}

func (o *Octree) Draw(g *Graphics, depth int) {
	depth++
	if depth > 4 {
		return
	}

	if o.node != nil {
		g.Color(0.2, 0.2, 0.2, 0.2)
		o.voxel.Draw(g)
	}
	if o.children != nil {
		for _, child := range o.children {
			child.Draw(g, depth)
		}
	}
}

func (o *Octree) Push(object Object) {
	o.Subdivide()
	for _, node := range o.children {
		if node.voxel.Contains(object) {
			node.Insert(object)
			return
		}
	}
}

func (o *Octree) Insert(object Object) {
	if o.node == nil {
		o.node = object
		return
	} else if len(o.children) < 8 {
		// External Node
		o.Push(o.node)
	} else {
		// Internal Node

	}

	physics := Physics{}
	ss := physics.AddBody(o.node, object)
	o.node = ss
	o.Push(object)
}

func (o *Octree) Subdivide() {
	if o.children == nil {
		o.children = make([]*Octree, 8)
		for i := 0; i < 8; i++ {
			o.children[i] = o.voxel.SubVoxel(i)
		}
	}
}

func (o *Octree) ComputeForce(matter []*Object, tick float64) {

}
