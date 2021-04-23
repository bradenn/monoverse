package main

import "github.com/veandco/go-sdl2/sdl"

const (
	menuX = 16
	menuY = 16
)

type Grid struct {
	size F2
	cell F2
}

func (g *Grid) Draw(graphics *Graphics) {
	divX := g.size.X / g.cell.X
	for i := 0.0; i < divX; i++ {
		graphics.Color(0.1, 0.08, 0.08, 1)
		graphics.Line(F3{(g.cell.X) * i, 0, 0}, F3{g.cell.X * i, g.size.Y, -1})
	}
	divY := g.size.Y / g.cell.Y
	for i := 0.0; i < divY; i++ {
		graphics.Color(0.1, 0.08, 0.08, 1)
		graphics.Line(F3{0, (g.cell.Y) * i, 0}, F3{g.size.X, g.cell.Y * i, 2})
	}
}

func (g *Grid) GetLocation(point F2) F2 {
	return F2{point.X * g.cell.X, point.Y * g.cell.Y}
}

func (g *Grid) GetSize(size F2) F2 {
	return F2{size.X * g.cell.X, size.Y * g.cell.Y}
}

type Color struct {
	r, g, b, a float64
}

func From32ARGB(r, g, b, a float64) Color {
	c := Color{
		r: (1 / 255) * r,
		g: (1 / 255) * g,
		b: (1 / 255) * b,
		a: (1 / 255) * a,
	}
	return c
}

type ListItem interface {
	Draw(g *Graphics, size F2)
}

type List struct {
	label    string
	elements []*Item
	location F2
	size     F2
}

func NewList(label string, location F2, size F2) List {
	return List{
		label:    label,
		elements: nil,
		location: location,
		size:     size,
	}
}

func (l *List) GetLocation() F2 {
	return l.location
}

func (l *List) GetSize() F2 {
	return l.size
}

func (l *List) Configure() {
	l.elements = make([]*Item, 0)
}
func (l *List) HandleEvent(event sdl.Event) {

}
func (l *List) Update() {

}

func (l *List) AddItem(item *Item) {
	l.elements = append(l.elements, item)
}

func (l *List) Draw(g *Graphics) {
	g.DrawFrame(l.label, l)
	for _, item := range l.elements {
		g.Translate(F3{0, 28, 0})
		item.Draw(g, F2{l.size.X, float64(28)})
	}
}

type Item struct {
	label  string
	value  string
	color  Color
	indent int
}

func DrawWindow(g *Graphics, name string, location F2, size F2) {
	g.Color((1/255)*43, (1/255)*45, (1/255)*45, 1)
	g.Rect(
		F2{
			X: 1,
			Y: 1,
		},
		F2{
			X: size.X - 2,
			Y: size.Y - 2,
		})
	g.Color((1/255)*43, (1/255)*45, (1/255)*45, 1)
	g.FillRect(
		F2{
			X: 1,
			Y: 1,
		},
		F2{
			X: size.X - 2,
			Y: 28,
		})
	g.Text(name, location)
}

func NewItem(label string, value string) *Item {
	return &Item{
		label:  label,
		value:  value,
		indent: 0,
	}
}

func (u *Item) Draw(g *Graphics, size F2) {
	color := 1 / 255.0
	g.Color(color*66, color*66, color*66, 1)
	g.Line(F3{1, size.Y, 0}, F3{size.X - 2, size.Y, 0})
	g.Color(color*255, color*255, color*255, 1)
	g.Text(u.label, F2{8, 6})
	length, _, _ := g.font.SizeUTF8(u.value)
	g.Text(u.value, F2{size.X - float64(length/2) - 8, 6})
}

func (u *Item) Indent() {
	u.indent++
}
