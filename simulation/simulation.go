package simulation

import (
	"fmt"
	"github.com/bradenn/monoverse/graphics"
	"github.com/bradenn/monoverse/physics"
	_ "github.com/veandco/go-sdl2/sdl"
	"math"
	"math/rand"
	"syscall"
	"time"
)

// Simulation [Renderer, Octree, State]

// Smallest Traversable Distance
// 1 unit length / precision
// 1 unit / 100 = smallest traversable distance = 0.01 units

// Fastest Velocity Possible
// 1 unit length / 1 iteration
// 1 unit length /

type Simulation struct {
	Engine *graphics.Engine
	frame  *graphics.Frame
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
	s.Engine, _ = graphics.NewEngine("Monoverse", int32(s.width+1), int32(s.width+1))
	s.delta = 1
	s.running = true
	s.frame = graphics.NewFrame(s.Engine, s.width, s.width, 4, 3)
	for i := 1; i < 16; i++ {
		s.bodies = append(s.bodies, physics.NewBody(s.width/2-rand.Float64()*s.width, s.width/2-rand.Float64()*s.width, s.width/2-rand.Float64()*s.width,
			rand.Float64()*100+100))
	}
	return
}

func (s *Simulation) Unload() {
	s.Engine.Cleanup()
}

func (s *Simulation) Run() {

	// buffer := make([]time.Duration, 0)
	fmt.Println("== Beginning Simulation ==")
	logicTicks := 0.0
	clock := 0.0
	logicInterval := 1000.0 / 60.0 // 15.625 ms
	delta := 0.0

	currTime := time.Now()
	prevTime := time.Now()

	// start := time.Now()
	s.frame.Panes[0].SetScale(1)
	for s.Engine.Running {

		prevTime = currTime
		currTime = time.Now()

		delta = currTime.Sub(prevTime).Seconds() * 1000
		clock += delta
		for clock > logicInterval {
			s.clock += 1
			s.HandleEvents()
			s.Engine.Clear()

			rusage := new(syscall.Rusage)
			syscall.Getrusage(0, rusage)

			s.frame.Panes[0].RotateXBy(math.Pi / 1024)
			s.frame.Panes[0].RotateYBy(math.Pi / 1024)

			s.Engine.Window.SetTitle(fmt.Sprintf("%.2f FPS %.2f MB",
				1000/delta, float64(rusage.Maxrss)/1024/1024))

			s.Engine.Renderer.SetDrawColor(255, 255, 255, 255)
			for n := 0.0; n < 10; n++ {
				for m := 0.0; m < 10; m++ {
					rad := 50.0 / 2.0
					col := math.Sqrt(3.0) * rad
					row := rad * 2.0
					b := math.Mod(m, 2) == 0
					if b {
						s.frame.Panes[0].DrawHex(graphics.Vertex3D{
							X: col * n,
							Y: row * 3 / 4 * m,
							Z: n}, rad)
					} else {
						s.frame.Panes[0].DrawHex(graphics.Vertex3D{
							X: col*n + (rad) - 2,
							Y: row * 3 / 4 * m,
							Z: n}, rad)
					}

				}

			}

			logicTicks += 1.0
			clock -= logicInterval
			s.Engine.Render()
			// if logicTicks > 64.0 {
			// 	running = false
			// 	break
			// }
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
	// 		s.Engine.Window.SetTitle(fmt.Sprintf("%.2f FPS %.2f MB",
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
	s.Engine.HandleEvents()
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
