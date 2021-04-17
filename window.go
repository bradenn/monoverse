package main

import "github.com/bradenn/monoverse/graphics"

type View struct {
	location graphics.F3
	size     graphics.F3
	camera   *Camera
}

func (w *View) NewView(name string, location graphics.F3, size graphics.F3) {

}

func (w *View) Draw() {

}
