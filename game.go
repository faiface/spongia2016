package main

import (
	"math"
	"math/rand"
	"sync"

	"github.com/faiface/gogame"
)

type game struct {
	sync.Mutex
	width, height float64
	things        []*thing
	hits          []*hit
	pointer       gogame.Vec
	background    gogame.Color
}

func newGame(width, height float64) *game {
	g := &game{
		width:  width,
		height: height,
	}
	go play(g)
	return g
}

func (g *game) update(dt float64) {
	g.Lock()
	for _, t := range g.things {
		t.update(dt)
	}
	g.Unlock()
}

func (g *game) draw(out gogame.VideoOutput) {
	g.Lock()
	out.SetMask(gogame.Colors["white"])
	out.Clear(g.background)

	for _, h := range g.hits {
		h.draw(out)
	}

	for _, t := range g.things {
		t.draw(out)
	}

	g.Unlock()
}

func (g *game) setPointer(p gogame.Vec) {
	g.Lock()
	defer g.Unlock()
	g.pointer = p
}

func (g *game) setBackground(bg gogame.Color) {
	g.Lock()
	defer g.Unlock()
	g.background = bg
}

func (g *game) addThing(t *thing) *thing {
	g.Lock()
	defer g.Unlock()
	g.things = append(g.things, t)
	return t
}

func play(g *game) {
	depth := 1.5 * math.Max(g.width, g.height)

	g.setBackground(gogame.Color{
		R: 0.1,
		G: 0.1,
		B: 0.1,
		A: 1,
	})

	thing := g.addThing(newThing(
		randomThingColor(),
		gogame.Vec{X: g.width * rand.Float64(), Y: g.height * 1.1},
		depth,
	))

	for {
		ok := func() bool {
			g.Lock()
			defer g.Unlock()

			if thing.position.Y <= 0 {
				thing.velocity.Y = 0
				return false
			}

			thing.velocity.Y = -100

			return true
		}()
		if !ok {
			break
		}
	}
}

func randomThingColor() gogame.Color {
newColor:
	c := gogame.Color{
		R: rand.Float64(),
		G: rand.Float64(),
		B: rand.Float64(),
		A: 1,
	}
	sum := c.R + c.G + c.B
	if sum == 0 {
		goto newColor
	}
	k := 1.25 / sum
	c.R *= k
	c.G *= k
	c.B *= k
	return c
}
