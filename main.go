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
		R: rand.Float64(),
		G: rand.Float64(),
		B: rand.Float64(),
		A: 1,
	}
	thing := newThing(color, gogame.Vec{}, 1.5*math.Max(float64(cfg.Width), float64(cfg.Height)))

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	err = gogame.Loop(cfg, func(ctx gogame.Context) {
		thing.position = ctx.MousePosition()
		thing.update(ctx.Dt)

		ctx.Clear(gogame.Colors["black"])
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
