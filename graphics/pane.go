package graphics

import (
	"fmt"
	"github.com/bradenn/monoverse"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Ratio struct {
	A float64
	B float64
}

type Pane struct {
	Aspect     Ratio
	name       string
	engine     *Engine
	texture    *sdl.Texture
	width      float64
	location   F3
	center     F3
	size       F3
	actualSize F2
	height     float64
	active     bool // Add an active object with stuff
	Camera     *main.Camera
	Renderer   *Renderer
}

func NewPane(eng *Engine, name string, location F3, size F3) *Pane {
	font, _ := ttf.OpenFont("./JetBrainsMonoNL-Regular.ttf", 24)
	pane := &Pane{
		Camera:     &main.Camera{location, size, F3{}, 1, 90},
		location:   location,
		engine:     eng,
		size:       size,
		name:       name,
		actualSize: size,
		Aspect: Ratio{
			A: 16,
			B: 9,
		},
		center: F3{
			X: location.X + size.X/2,
			Y: location.Y + size.Y/2,
			Z: location.Z,
		},
		Renderer: &Renderer{
			Ref:    eng.Renderer,
			Camera: &main.Camera{location, size, F3{}, 1, 90},
			Font:   font,
		},
	}
	pane.Renderer.Init()
	pane.Camera = pane.Renderer.Camera
	return pane
}

func (p *Pane) Contains(f F2) {
	xBound := p.location.X <= f.X && p.location.X+p.size.X >= f.X
	yBound := p.location.Y <= f.Y && p.location.Y+p.size.Y >= f.Y

	p.active = xBound && yBound
}

func (p *Pane) Render() {

	p.Renderer.SetColor(255, 131, 0, 255)
	p.Renderer.UIText(fmt.Sprintf("%s", p.name), F3{8, 6, 0})

	crossHair := 4.0
	if p.active {
		p.Renderer.DrawCube(F3{0, 0, 0}, 20)
		p.Renderer.DrawLine(
			F3{
				X: -crossHair,
				Y: 0,
				Z: 0,
			},
			F3{
				X: crossHair,
				Y: 0,
				Z: 0,
			})

		p.Renderer.DrawLine(
			F3{
				X: 0,
				Y: -crossHair,
				Z: 0,
			},
			F3{
				X: 0,
				Y: crossHair,
				Z: 0,
			})

		p.Renderer.DrawLine(
			F3{
				X: 0,
				Y: 0,
				Z: -crossHair,
			},
			F3{
				X: 0,
				Y: 0,
				Z: crossHair,
			})
	}
	// p.Renderer.UIRect(F3{0, 52, 0}, F2{X: p.size.X, Y: 2})

	p.Renderer.SetColor(47, 24, 0, 255)
	p.Renderer.UIRect(
		F3{
			X: p.location.X + 1,
			Y: p.location.Y + 1,
			Z: 0,
		},
		F2{
			X: p.size.X - 2,
			Y: p.size.Y - 2,
		})

	p.Renderer.UIRect(
		F3{
			X: p.location.X + 2,
			Y: p.location.Y + 2,
			Z: 0,
		},
		F2{
			X: p.size.X - 4,
			Y: p.size.Y - 4,
		})

}
