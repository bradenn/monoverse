package main

import (
	"fmt"
)

type Point struct {
	x, y, z float64
	r, m    float64
}

type Medium struct {
	x, y, z float64
	w       float64
}

//			0	       1			2		    3
// 0 - 7,  (1, 1, 1), (1, -1, 1), (-1, 1, 1), (-1, -1, 1)
func (m *Medium) GenerateSubtree(index int) (e *Octree) {

	a := m.x/4
	b := m.y/4
	c := m.z/4

	if index > 4 {
		a = -a
	}

	if index % 2 != 0 {
		b = -b
	}

	// We need alternating numbers now, like 00110011
	// We can use the decimal place of an repeating value...
	// So anything divided by 11 will be multiplied by nine and repeated in the decimal place.
	// So 5/11 = 0.454545, We need to get 1100 to repeat. So to find that, (1100/9)/1111
	// Math is weird... (1100/9)/1111 = 0.11001100

	swap := (1100.0/9.0)/1111.0
	swstr := fmt.Sprintf("%.8f", swap)
	if  swstr[2:][index] == 49 {
		c = -c
	}
	
	medium := &Medium{
		x: m.x + a,
		y: m.y + b,
		z: m.z + c,
		w: m.w/2,
	}
	
	return &Octree{
		medium: *medium,
	}
}

func (m *Medium) Contains(point *Point) bool {
	if point.x >= m.x && point.x <= m.x+m.w {
		if point.y >= m.y && point.y <= m.y+m.w {
			return true
		}
	}
	return false
}

type Octree struct {
	point    *Point
	medium   Medium
	children []*Octree
}

func (o *Octree) insert(point Point) {
	if o.point == nil {
		o.point = &point
	} else if len(o.children) > 0 {
		for _, child := range o.children {
			if child.medium.Contains(&point) {
				child.insert(point)
				break
			}
		}
	} else {
		o.children = make([]*Octree, 8)
		for i := range o.children {
			o.children[i] = o.medium.GenerateSubtree(i)
		}
	}
}

func (o *Octree) print(clk int) {
	for j := 0; j < clk; j++ {
		fmt.Print(" ")
	}
	fmt.Printf("%f\n", o.medium.w)
	clk++
	for _, child := range o.children {
		child.print(clk)
	}
}
