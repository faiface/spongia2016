package main

import (
	"log"
	"math"

	"github.com/faiface/gogame"
)

func main() {
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

	thing := newThing(gogame.Colors["red"], gogame.Vec{}, 1.5*math.Max(float64(cfg.Width), float64(cfg.Height)))

	err = gogame.Loop(cfg, func(ctx gogame.Context) {
		thing.position = ctx.MousePosition()
		thing.update(ctx.Dt)

		ctx.Clear(gogame.Colors["black"])
		thing.draw(ctx)
	})

	if err != nil {
		log.Fatal(err)
	}
}
