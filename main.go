package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/gogame"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start().Stop()

	rand.Seed(time.Now().UnixNano())

	var err error

	err = gogame.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer gogame.Quit()

	cfg := gogame.Config{
		Title:       "Špongia 2016",
		Width:       1024,
		Height:      768,
		VSync:       true,
		QuitOnClose: true,
	}

	color := gogame.Color{
		R: 0.85,
		G: 0.15,
		B: 0.15,
		A: 1,
	}
	thing := newThing(color, gogame.Vec{}, 1.5*math.Max(float64(cfg.Width), float64(cfg.Height)))

	trajectory := new(function)
	trajectory.cyclic = true
	trajectory.add(0.0, gogame.Vec{X: 100, Y: 100})
	trajectory.add(3.0, gogame.Vec{X: 900, Y: 100})
	trajectory.add(6.0, gogame.Vec{X: 900, Y: 600})
	trajectory.add(9.0, gogame.Vec{X: 100, Y: 600})
	trajectory.add(12.0, gogame.Vec{X: 100, Y: 100})

	passed := 0.0

	for t := 0.0; t < 20.0; t += 1.0 / 60 {
		thing.update(1.0 / 60)
	}

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	err = gogame.Loop(cfg, func(ctx gogame.Context) {
		passed += ctx.Dt
		thing.position = trajectory.at(passed)
		thing.update(ctx.Dt)

		ctx.SetMask(gogame.Colors["white"])
		ctx.Clear(gogame.Color{R: 0.15, G: 0.15, B: 0.15, A: 1})
		thing.draw(ctx)

		frames++
		select {
		case <-second:
			ctx.WindowSetTitle(fmt.Sprintf("Špongia 2016 | FPS: %d", frames))
			frames = 0
		default:
		}
	})

	if err != nil {
		log.Fatal(err)
	}
}
