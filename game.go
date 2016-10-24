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
	explosions    []*explosion
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
		for _, e := range g.explosions {
			e.update(dt)
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

	for _, e := range g.explosions {
		e.draw(out)
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
	g.setBackground(gogame.Color{
		R: 0.1,
		G: 0.1,
		B: 0.1,
		A: 1,
	})

	time.Sleep(5 * time.Second)

	for first := true; ; first = false {
		firstPart(g, first)
		switch g.phase {
		case explosionPhase:
			secondPart(g)
		case colorfulPhase:
			thirdPart(g)
		}
	}
}

func firstPart(g *game, first bool) {
	depth := 1.5 * math.Max(g.width, g.height)

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

				if first && i == 0 { // first one
					currentHeight := thing.position.Y
					targetHeight := g.height * 0.8
					thing.velocity.Y = 1.1 * (targetHeight - currentHeight)
					if math.Abs(thing.velocity.Y) > 50 {
						thing.velocity.Y *= 50 / math.Abs(thing.velocity.Y)
					}
				} else { // remaining ones
					if thing.position.Y <= -77 {
						g.phase = colorfulPhase
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

		if g.phase == colorfulPhase {
			return
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

func secondPart(g *game) {
	time.Sleep(time.Duration(hitLightUpTime * float64(time.Second)))

	g.Lock()
	for _, h := range g.hits {
		g.explosions = append(g.explosions, newExplosion(h))
	}
	g.hits = nil
	g.Unlock()

	time.Sleep(5 * time.Second)

	g.Lock()
	g.explosions = nil
	g.phase = dangerPhase
	g.Unlock()
}

func thirdPart(g *game) {
	depth := 1.5 * math.Max(g.width, g.height)

	g.Lock()
	g.explosions = nil
	g.Unlock()

	for {
		g.Lock()
		thing := g.things[len(g.things)-1]
		g.Unlock()

		for {
			ok := func() bool {
				g.Lock()
				defer g.Unlock()
				if thing.position.Y <= -thing.depth {
					return false
				}
				return true
			}()
			if !ok {
				break
			}
		}

		g.setBackground(thing.color)
		g.removeThing(thing)
		g.Lock()
		g.hits = nil
		g.Unlock()

		newThing := g.addThing(newThing(
			randomThingColor(),
			gogame.Vec{X: g.width * rand.Float64(), Y: g.height * 1.1},
			depth,
		))
		newThing.velocity.Y = -1000
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
