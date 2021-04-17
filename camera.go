package main

import (
	"math"
)

type Camera struct {
	Location F3
	Size     F3
	Scale    float64
	Rotation F3
	Fov      float64
}

func (c *Camera) matMul(a [][]float64, b []float64) (f []float64) {
	colsA := len(a[0])
	rowsA := len(a)
	f = make([]float64, rowsA)

	colsB := len(b)

	for i := 0; i < rowsA; i++ {
		for j := 0; j < colsB; j++ {
			sum := 0.0
			for k := 0; k < colsA; k++ {
				sum += a[i][k] * b[k]
			}
			f[i] = sum
		}
	}
	return f
}

func (c *Camera) ProjectPoint(point F3) F2 {
	// X Rotation
	x := c.Rotation.X
	rotX := c.matMul([][]float64{
		{math.Cos(x), 0, math.Sin(x)},
		{0, 1, 0},
		{-math.Sin(x), 0, math.Cos(x)},
	}, []float64{point.X, point.Y, point.Z})
	// Y Rotation
	y := c.Rotation.Y
	rotY := c.matMul([][]float64{
		{1, 0, 0},
		{0, math.Cos(y), math.Sin(y)},
		{0, math.Sin(y), math.Cos(y)},
	}, rotX)
	// Z Rotation
	z := c.Rotation.Z
	rotZ := c.matMul([][]float64{
		{math.Cos(z), -math.Sin(z), 0},
		{math.Sin(z), math.Cos(z), 0},
		{0, 0, 0},
	}, rotY)
	// X, Y, Z Translation
	translation := c.matMul([][]float64{
		{1, 0, 0, c.Location.X + c.Size.X/2},
		{0, 1, 0, c.Location.Y + c.Size.Y/2},
		{0, 0, 1, c.Location.Z},
		{0, 0, 0, 1},
	}, []float64{rotZ[0], rotZ[1], rotZ[2], 1})
	// Scale
	scale := c.matMul([][]float64{
		{c.Scale, 0, 0, 0},
		{0, c.Scale, 0, 0},
		{0, 0, c.Scale, 0},
		{0, 0, 0, 1},
	}, []float64{translation[0], translation[1], translation[2], 1})
	// Convert to 2D
	points := c.matMul([][]float64{
		{1, 0, 0},
		{0, 1, 0},
	}, []float64{scale[0], scale[1], scale[2]})
	return F2{points[0], points[1]}

}
