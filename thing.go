package main

import (
	"math"
	"math/rand"

	"github.com/faiface/gogame"
)

type thing struct {
	color    gogame.Color
	position gogame.Vec
	depth    float64
	waves    []wave

	time float64
}

func newThing(color gogame.Color, position gogame.Vec, depth float64) *thing {
	return &thing{
		color:    color,
		position: position,
		depth:    depth,
	}
}

func (t *thing) draw(out gogame.VideoOutput) {
	bottom := out.OutputRect().Y + out.OutputRect().H
	left, right := t.position.X-(t.position.Y-bottom), t.position.X+(t.position.Y-bottom)

	out.DrawPolygon([]gogame.Vec{t.position, {X: left, Y: bottom}, {X: right, Y: bottom}}, 0, t.color)

	for i := range t.waves {
		t.waves[i].draw(out)
	}
}

func (t *thing) update(dt float64) {
	t.time -= dt

	var toDelete []int
	for i := range t.waves {
		t.waves[i].start = t.position
		t.waves[i].update(dt)

		if t.waves[i].dir.M(t.waves[i].time).Len() > t.depth {
			toDelete = append(toDelete, i)
		}
	}

	for j, i := range toDelete {
		j = len(t.waves) - j - 1
		t.waves[i], t.waves[j] = t.waves[j], t.waves[i]
	}
	t.waves = t.waves[:len(t.waves)-len(toDelete)]

	if t.time < 0 {
		var dir gogame.Vec
		if rand.Intn(2) == 0 {
			dir.X = -1
			dir.Y = +1
		} else {
			dir.X = +1
			dir.Y = +1
		}

		freq := rand.Float64()*2 + 0.5
		if dir.X < 0 {
			freq *= -1
		}

		t.waves = append(t.waves, wave{
			color:     t.color,
			start:     t.position,
			dir:       dir.M(40),
			size:      rand.Float64()*50 + 20,
			frequency: freq,
			time:      -1,
		})
		t.time = rand.Float64()*0.3 + 0.1
	}
}

type wave struct {
	color           gogame.Color
	start, dir      gogame.Vec
	size, frequency float64
	time            float64
}

func (w *wave) draw(out gogame.VideoOutput) {
	var position gogame.Vec
	if w.time < 0 {
		position = w.start.A(gogame.Vec{X: 0, Y: 1}.M(w.dir.Len()).M(-w.time))
	} else {
		position = w.start.A(w.dir.M(w.time))
	}
	angle := w.time * w.frequency

	base := gogame.Vec{X: 0, Y: w.size}
	a := base.Rotated(angle + 0*math.Pi/2).A(position)
	b := base.Rotated(angle + 1*math.Pi/2).A(position)
	c := base.Rotated(angle + 2*math.Pi/2).A(position)
	d := base.Rotated(angle + 3*math.Pi/2).A(position)

	out.DrawPolygon([]gogame.Vec{a, b, c, d}, 0, w.color)
}

func (w *wave) update(dt float64) {
	w.time += dt
}
