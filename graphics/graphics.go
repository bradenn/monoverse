package graphics

type Point interface {
	Draw() (x float64, y float64, z float64, m float64)
}

type Vertex2D struct {
	X float64
	Y float64
}

type Vertex3D struct {
	X float64
	Y float64
	Z float64
}
