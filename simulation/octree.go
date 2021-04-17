package simulation

import (
	"fmt"
	"github.com/bradenn/monoverse"
	"github.com/bradenn/monoverse/physics"
)

type Octree struct {
	body     *physics.GBody
	medium   *Medium
	children []*Octree
}

func Init(width float64) (o *Octree) {
	o = new(Octree)
	o.children = make([]*Octree, 0)
	o.medium = &Medium{
		w: width,
	}
	o.body = nil
	return o
}

func (o *Octree) insert(body *physics.GBody) {
	if o.body == nil {
		o.body = body
		return
	} else if len(o.children) == 0 {
		children := make([]*Octree, 8)
		for i := range children {
			children[i] = o.medium.GenerateSubtree(i)

		}
		o.children = children
		o.pushBody(o.body)
	}
	o.body.AddBody(*body)
	o.pushBody(body)

}

func (o *Octree) pushBody(body *physics.GBody) {
	for _, child := range o.children {
		if child.medium.Contains(body) {
			child.insert(body)
			break
		}
	}
}

func (o *Octree) draw(frame *main.Frame, d int) {
	//
	//
	// for i := 0.0; i <= 6; i+=1{
	// 	x := 128 + 128 / 6 * i
	// 	y := (131 / 6) * i
	// 	z := 0
	// 	frame.Renderer.SetDrawColor(uint8(x), uint8(y), uint8(z),
	// 		255)
	// 	if i == 0 {
	//
	// 	}
	//
	//
	// }

	// frame.Renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	// frame.Renderer.SetDrawColor(255, 131, 0,
	// 	255)
	// frame.Panes[0].DrawPoint(graphics.Vertex3D{X: o.body.Location.X, Y: o.body.Location.Y, Z: o.body.Location.Z})
	// frame.Renderer.SetDrawColor(255, 131, 0,
	// 	64)
	// frame.Panes[0].DrawCube(graphics.Vertex3D{X: o.medium.X, Y: o.medium.Y, Z: o.medium.Z},
	// 	o.medium.w)

	// engine.Renderer.SetDrawColor(0, 131, 255, 12)
	// engine.DrawCube(graphics.Vertex3D{X: o.medium.X, Y: o.medium.Y, Z: o.medium.Z}, o.medium.w * 0.98)

	//
	// za := float32(o.medium.Z - (o.medium.w / 2))
	// zb := float32(o.medium.Z + (o.medium.w / 2))
	// engine.DrawLine()
	// engine.DrawPoint()
	// engine.DrawPoint(graphics.Vertex3D{X: o.medium.X, Y: o.medium.Y, Z: o.medium.Z})

	// renderer.DrawLineF(xa, ya, xb, ya)
	// renderer.DrawLineF(xa, yb, xb, yb)
	//
	// renderer.DrawLineF(xb, yb, xb, ya)
	// renderer.DrawLineF(xa, ya, xa, yb)

	// renderer.DrawLineF(xb, yb, xb, yb)
	// renderer.DrawLineF(xa, ya, xa, ya)

	for _, child := range o.children {
		if child != nil {
			if child.body != nil && d < 32 {
				child.draw(frame, d+1)
			}
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

func (o *Octree) ApplyForces(body *physics.GBody) {

	if o.body != nil {
		if o.children == nil {
			if body != o.body {
				body.AddForce(*o.body)
			}
		} else {
			theta := o.medium.w / body.DistanceTo(o.body.Location)
			if theta < 0.9 {
				body.AddForce(*o.body)
			} else {
				for _, child := range o.children {
					child.ApplyForces(body)
				}
			}
		}
	}
}
