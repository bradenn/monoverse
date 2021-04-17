package simulation

import (
	"fmt"
	"github.com/bradenn/monoverse"
	"github.com/bradenn/monoverse/graphics"
	"github.com/bradenn/monoverse/physics"
	"github.com/veandco/go-sdl2/sdl"
	_ "github.com/veandco/go-sdl2/sdl"
	"math"
	"math/rand"
	"syscall"
)

// Simulation [Renderer, Octree, State]

// Smallest Traversable Distance
// 1 unit length / precision
// 1 unit / 100 = smallest traversable distance = 0.01 units

type B bool
type F1 float64

type F2 struct {
	X float64
	Y float64
}

type F3 struct {
	X float64
	Y float64
	Z float64
}
type Simulation struct {
	Engine *graphics.Engine
	frame  *main.Frame
	width  float64

	delta float64
	clock float64

	running bool
	octree  *Octree
	bodies  []*physics.GBody
}

func NewSimulation() (s *Simulation) {
	s = new(Simulation)
	s.width = 1024
	s.Engine, _ = graphics.NewEngine("Monoverse", int32(1920), int32(1080))
	s.delta = 1
	s.running = true
	s.frame = main.NewFrame(s.Engine, 1920, 1080, 16, 9)
	for i := 1; i < 256; i++ {
		s.bodies = append(s.bodies, physics.NewBody(s.width/2-rand.Float64()*s.width, s.width/2-rand.Float64()*s.width, s.width/2-rand.Float64()*s.width,
			rand.Float64()*100+100))
	}
	return
}

func (s *Simulation) Unload() {
	s.Engine.Cleanup()
}

func (s *Simulation) Run() {

	fmt.Println("== Beginning Simulation ==")
	logicTicks := 0.0
	clock := 0.0
	logicInterval := 1000.0 / 60.0 // 15.625 ms
	delta := 0.0

	currTime := sdl.GetTicks()
	prevTime := sdl.GetTicks()

	pane := graphics.NewPane(s.Engine, "Simulation", graphics.F3{0, 0, 0}, graphics.F2{1440, 28})
	pane2 := graphics.NewPane(s.Engine, "Monitor", graphics.F3{1440, 0, 0}, graphics.F2{480, 28})
	pane4 := graphics.NewPane(s.Engine, "Interface", graphics.F3{Y: 28}, graphics.F2{1440, 1052})
	pane3 := graphics.NewPane(s.Engine, "Statistics", graphics.F3{1440, 28, 0}, graphics.F2{480, 1052})

	s.frame.Panes = append(s.frame.Panes, pane)
	s.frame.Panes = append(s.frame.Panes, pane2)
	s.frame.Panes = append(s.frame.Panes, pane3)

	s.frame.Panes = append(s.frame.Panes, pane4)
	renderClock := 0.0
	for s.running {

		prevTime = currTime
		currTime = sdl.GetTicks()

		rusage := new(syscall.Rusage)
		syscall.Getrusage(0, rusage)

		delta = float64(currTime - prevTime)
		clock += delta
		if clock > logicInterval {
			s.HandleEvents()
			s.frame.Clear()

			pane2.Renderer.UIText(fmt.Sprintf("%.2f FPS %.2f MB", 1000.0/delta, float64(rusage.Maxrss)/1024/1024), graphics.F3{
				X: pane2.Renderer.Camera.Size.X - 144,
				Y: pane2.Renderer.Camera.Size.Y/2 - 8,
				Z: 0,
			})
			pane4.Camera.Rotation.X += math.Pi / 2048
			pane4.Camera.Rotation.Y += math.Pi / 2048

			pane4.Camera.Scale = 1
			outer := 512.0

			pane4.Renderer.SetColor(255, 131, 0, 64)
			for n := -1.0; n < 1; n += 0.1 {
				if n == 0 {
					pane4.Renderer.SetColor(255, 131, 0, 128)
				} else {
					pane4.Renderer.SetColor(255, 131, 0, 64)
				}

				pane4.Renderer.DrawLine(graphics.F3{-outer, n * outer, 0}, graphics.F3{outer, n * outer, 0})
				pane4.Renderer.DrawLine(graphics.F3{n * outer, -outer, 0}, graphics.F3{n * outer, outer, 0})

				pane4.Renderer.DrawLine(graphics.F3{-outer, 0, n * outer}, graphics.F3{outer, 0, n * outer})
				pane4.Renderer.DrawLine(graphics.F3{n * outer, 0, -outer}, graphics.F3{n * outer, 0, outer})

				pane4.Renderer.DrawLine(graphics.F3{0, -outer, n * outer}, graphics.F3{0, outer, n * outer})
				pane4.Renderer.DrawLine(graphics.F3{0, n * outer, -outer}, graphics.F3{0, n * outer, outer})
			}
			pane4.Renderer.SetColor(255, 131, 0, 255)
			for _, body := range s.bodies {

				pane4.Renderer.DrawPoint(graphics.F3(body.Location))
			}
			s.frame.Render()
			// s.Engine.Renderer.SetDrawColor(255, 128, 128, 255)
			// pane4.Renderer.DrawLine(graphics.F3{0, 0, 0}, graphics.F3{64, 0, 0})
			// s.Engine.Renderer.SetDrawColor(128, 255, 128, 255)
			// pane4.Renderer.DrawLine(graphics.F3{0, 0, 0}, graphics.F3{0, 64, 0})
			// s.Engine.Renderer.SetDrawColor(128, 128, 255, 255)
			// pane4.Renderer.DrawLine(graphics.F3{0, 0, 0}, graphics.F3{0, 0, 64})

			// for n := 0.0; n < 10; n++ {
			// 	for m := 0.0; m < 10; m++ {
			// 		rad := 50.0 / 2.0
			// 		col := math.Sqrt(3.0) * rad
			// 		row := rad * 2.0
			// 		b := math.Mod(m, 2) == 0
			// 		if b {
			// 			pane4.Renderer.DrawHex(graphics.F3(graphics.Vertex3D{
			// 				X: 200 + col * n,
			// 				Y: 200 + row * 3 / 4 * m,
			// 				Z: 0}), rad)
			// 		} else {
			// 			pane4.Renderer.DrawHex(graphics.F3(graphics.Vertex3D{
			// 				X: 200 + col*n + (rad) - 2,
			// 				Y: 200 + row * 3 / 4 * m,
			// 				Z: 0}), rad)
			// 		}
			//
			// 	}
			//
			// }

			logicTicks += 1.0
			clock -= logicInterval
		}

		renderClock += delta
		if clock > 1000/24 {

			renderClock = 0
		}

	}

	// for i := range s.bodies {
	// 	s.bodies[i].ResetForces()
	// }

	// wg := new(sync.WaitGroup)
	// wg.Add(len(s.bodies) * len(s.bodies))
	// for _, body := range s.bodies {
	// 	for _, target := range s.bodies {
	// 		go func(so *physics.GBody, ta *physics.GBody) {
	// 			if *so != *ta {
	// 				so.AddForce(*ta)
	// 			}
	// 			wg.Done()
	// 		}(body, target)
	// 	}
	// }
	// wg.Wait()

	// for _, body := range s.bodies {
	// 	body.UpdatePosition(logicTicks)
	// }

	//
	//
	//
	// step := 1000.0 / 24.0
	// frame := 1000.0 / 60.0
	//
	// updateClock := 0.0
	// renderClock := 0.0
	//
	//
	//
	// updateDelta := 0.0
	// frameDelta := 0.0
	//
	// for s.running && s.Engine.Running {
	// 	prevTime = newTime
	// 	newTime = time.Now()
	// 	s.clock += s.delta
	// 	frameDelta = newTime.Sub(prevTime).Seconds() * 1000
	//
	//
	// 	updateClock += updateDelta
	// 	for updateClock > step {
	// 		s.HandleEvents()
	//
	// 		s.octree = Init(s.width)
	// 		for i := range s.bodies {
	// 			s.bodies[i].ResetForces()
	// 			s.octree.insert(s.bodies[i])
	// 		}
	//
	// 		for _, body := range s.bodies {
	// 			s.octree.ApplyForces(body)
	// 			// for _, target := range s.bodies {
	// 			// 	if body != target {
	// 			// 		body.AddForce(*target)
	// 			// 	}
	// 			// }
	// 			// if i + 1 < len(s.bodies) {
	// 			// 	body.AddForce(*s.bodies[i+1])
	// 			// }
	// 		}
	//
	// 		for _, body := range s.bodies {
	// 			body.UpdatePosition(s.clock)
	// 		}
	//
	// 		updateClock -= step
	// 	}
	// 	updateDelta = newTime.Sub(prevTime).Seconds() * 1000
	// 	renderClock += frameDelta
	// 	if renderClock > frame {
	// 		rusage := new(syscall.Rusage)
	// 		syscall.Getrusage(0, rusage)
	//
	// 		// s.frame.Panes[0].RotateXBy(math.Pi / 480)
	// 		// s.frame.Panes[0].RotateYBy(math.Pi / 480)
	// 		s.Engine.View.SetTitle(fmt.Sprintf("%.2f FPS %.2f MB",
	// 			1000.0/frameDelta, float64(rusage.Maxrss) / 1024 / 1024))
	// 		s.Engine.Clear()
	//
	// 		if s.octree != nil {
	// 			s.octree.draw(s.frame, 0)
	// 		}
	// 		s.frame.Panes[0].SetScale((1024/s.width)/2)
	// 		for _, body := range s.bodies {
	// 			s.Engine.Renderer.SetDrawColor(255, 255, 255, 255)
	// 			s.frame.Panes[0].DrawPoint(graphics.Vertex3D(body.Location))
	// 		}
	//
	//
	// 		s.Render()
	//
	// 		renderClock = 0
	// 	}
	//
	//
	//
	// }
}

func (s *Simulation) HandleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			s.running = false
			break
		case *sdl.KeyboardEvent:
			break
		case *sdl.MouseButtonEvent:
			s.frame.HandleClick(graphics.F2{float64(t.X), float64(t.Y)})
		}
	}
}

func (s *Simulation) GetOuterBound() float64 {
	max := s.width
	for _, body := range s.bodies {
		bd := math.Abs(body.DistanceTo(physics.Vector3D{}))
		if bd*2 > max {
			max = bd * 2
		}
	}
	s.width = max
	return max
}

func (s *Simulation) Update() {
	// s.GetOuterBound()

}

func (s *Simulation) Render() {

	s.Engine.Render()
}
