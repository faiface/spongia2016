package main

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/gogame"
)

const (
	dangerPhase = iota
	explosionPhase
	colorfulPhase
)

type game struct {
	sync.Mutex
	width, height float64
	phase         int
	things        []*thing
	hits          []*hit
	pointer       gogame.Vec
	background    gogame.Color
}

func newGame(width, height float64) *game {
	g := &game{
		width:  width,
		height: height,
		phase:  dangerPhase,
	}
	go play(g)
	return g
}

func (g *game) update(dt float64) {
	g.Lock()
	if g.phase == explosionPhase {
		for _, h := range g.hits {
			h.update(dt)
		}
	}
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

func (g *game) removeThing(t *thing) *thing {
	g.Lock()
	defer g.Unlock()
	for i, x := range g.things {
		if x == t {
			g.things = append(g.things[:i], g.things[i+1:]...)
			return t
		}
	}
	return nil
}

func (g *game) addHit(h *hit) *hit {
	g.Lock()
	defer g.Unlock()
	g.hits = append(g.hits, h)
	return h
}

func play(g *game) {
	depth := 1.5 * math.Max(g.width, g.height)

	g.setBackground(gogame.Color{
		R: 0.1,
		G: 0.1,
		B: 0.1,
		A: 1,
	})

	// remaining things
	for i := 0; i < 12; i++ {
		thing := g.addThing(newThing(
			randomThingColor(),
			gogame.Vec{X: g.width * rand.Float64(), Y: g.height * 1.1},
			depth,
		))

		for {
			ok := func() bool {
				g.Lock()
				defer g.Unlock()

				if g.pointer.S(thing.position).Len() <= 77 {
					return false
				}

				if i == 0 { // first one
					currentHeight := thing.position.Y
					targetHeight := g.height * 0.8
					thing.velocity.Y = 1.1 * (targetHeight - currentHeight)
					if math.Abs(thing.velocity.Y) > 50 {
						thing.velocity.Y *= 50 / math.Abs(thing.velocity.Y)
					}
				} else { // remaining ones
					if thing.position.Y <= 0 {
						thing.velocity.Y = 0
						return false
					}
					thing.velocity.Y = -100
				}

				return true
			}()
			if !ok {
				break
			}
		}

		g.addHit(newHit(thing))

		for {
			ok := func() bool {
				g.Lock()
				defer g.Unlock()

				if thing.position.Y >= g.height*1.1 {
					thing.acceleration.Y = 0
					thing.velocity.Y = 0
					return false
				}

				thing.acceleration.Y = 1000

				return true
			}()
			if !ok {
				break
			}
		}

		g.removeThing(thing)
	}

	time.Sleep(5 * time.Second)

	g.phase = explosionPhase
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
