package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/faiface/gogame"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	err := gogame.Init()
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

	g := newGame(float64(cfg.Width), float64(cfg.Height))

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	err = gogame.Loop(cfg, func(ctx gogame.Context) {
		g.setPointer(ctx.MousePosition())

		g.update(ctx.Dt)
		g.draw(ctx)

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
