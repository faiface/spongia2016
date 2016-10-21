package main

import (
	"math"
	"math/rand"

	"github.com/faiface/gogame"
)

var squares []*gogame.Picture

func init() {
	for size := 40; size <= 140; size += 5 {
		canvas := gogame.NewCanvas(size, size)
		canvas.Clear(gogame.Colors["white"])
		squares = append(squares, canvas.Picture().Copy())
	}
}

const (
	waveSpeed   = 40
	waveMinFreq = 0.5
	waveMaxFreq = 2.5

	waveMinSpawnTime = 0.1
	waveMaxSpawnTime = 0.4

	hitLightUpTime = 5.0

	explWaveCount    = 16
	explWaveMinSpeed = 1000
	explWaveMaxSpeed = 5000
	explWaveMinFreq  = 4
	explWaveMaxFreq  = 16
)

type thing struct {
	color        gogame.Color
	position     gogame.Vec
	velocity     gogame.Vec
	acceleration gogame.Vec
	depth        float64
	waves        []wave

	time float64
}

func newThing(color gogame.Color, position gogame.Vec, depth float64) *thing {
	tg := &thing{
		color:    color,
		position: position,
		depth:    depth,
	}
	for t := 0.0; t < depth/waveSpeed; t += 1.0 / 64 {
		tg.update(1.0 / 64)
	}
	return tg
}

func (t *thing) draw(out gogame.VideoOutput) {
	bottom := out.OutputRect().Y + out.OutputRect().H
	left, right := t.position.X-(t.position.Y-bottom), t.position.X+(t.position.Y-bottom)

	out.SetMask(t.color)
	out.DrawPolygon([]gogame.Vec{t.position, {X: left, Y: bottom}, {X: right, Y: bottom}}, 0, gogame.Colors["white"])

	for i := range t.waves {
		t.waves[i].draw(out)
	}
}

func (t *thing) update(dt float64) {
	t.time -= dt

	t.velocity = t.velocity.A(t.acceleration.M(dt))
	t.position = t.position.A(t.velocity.M(dt))

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
		dir := gogame.Vec{X: 0, Y: +1}
		if rand.Intn(2) == 0 {
			dir.X = -1
		} else {
			dir.X = +1
		}

		freq := rand.Float64()*(waveMaxFreq-waveMinFreq) + waveMinFreq
		if dir.X < 0 {
			freq *= -1
		}

		t.waves = append(t.waves, wave{
			square:    squares[rand.Intn(len(squares))],
			start:     t.position,
			dir:       dir.M(waveSpeed),
			frequency: freq,
			time:      -1,
		})
		t.time = rand.Float64()*(waveMaxSpawnTime-waveMinSpawnTime) + waveMinSpawnTime
	}
}

type wave struct {
	square     *gogame.Picture
	start, dir gogame.Vec
	frequency  float64
	time       float64
}

func (w *wave) position() gogame.Vec {
	if w.time < 0 {
		return w.start.A(gogame.Vec{X: 0, Y: 1}.M(w.dir.Len()).M(-w.time))
	}
	return w.start.A(w.dir.M(w.time))
}

func (w *wave) angle() float64 {
	return w.time * w.frequency
}

func (w *wave) draw(out gogame.VideoOutput) {
	position := w.position()
	angle := w.angle()

	sizeX, sizeY := w.square.Size()
	rect := gogame.Rect{
		X: position.X - float64(sizeX)/2,
		Y: position.Y - float64(sizeY)/2,
		W: float64(sizeX),
		H: float64(sizeY),
	}

	out.DrawPicture(rect, w.square.Rotated(angle))
}

func (w *wave) update(dt float64) {
	w.time += dt
}

type hit struct {
	color    gogame.Color
	square   *gogame.Picture
	position gogame.Vec
	angle    float64
	time     float64
}

func newHit(t *thing) *hit {
	bestI, bestY := 0, t.waves[0].position().Y
	for i, w := range t.waves {
		y := w.position().Y
		if y < bestY {
			bestI, bestY = i, y
		}
	}

	return &hit{
		color:    t.color,
		square:   t.waves[bestI].square,
		position: t.waves[bestI].position(),
		angle:    t.waves[bestI].angle(),
	}
}

func (h *hit) update(dt float64) {
	h.time += dt
}

func (h *hit) draw(out gogame.VideoOutput) {
	mul := math.Min(1, h.time/hitLightUpTime)
	mul *= mul
	mul = 0.3 + 0.7*mul
	mulColor := gogame.Color{
		R: mul,
		G: mul,
		B: mul,
		A: 1,
	}
	out.SetMask(h.color.Mul(mulColor))

	sizeX, sizeY := h.square.Size()
	rect := gogame.Rect{
		X: h.position.X - float64(sizeX)/2,
		Y: h.position.Y - float64(sizeY)/2,
		W: float64(sizeX),
		H: float64(sizeY),
	}

	out.DrawPicture(rect, h.square.Rotated(h.angle))
}

type explosion struct {
	waves []wave
}

func newExplosion(h *hit) *explosion {
	e := new(explosion)
	e.waves = make([]wave, explWaveCount)
	for i := range e.waves {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*(explWaveMaxSpeed-explWaveMinSpeed) + explWaveMinSpeed
		dir := gogame.Vec{
			X: math.Cos(angle) * speed,
			Y: math.Sin(angle) * speed,
		}
		freq := rand.Float64()*(explWaveMaxFreq-explWaveMinFreq) + explWaveMinFreq
		e.waves[i] = wave{
			square:    squares[rand.Intn(len(squares))],
			start:     h.position,
			dir:       dir,
			frequency: freq,
			time:      0,
		}
	}
	return e
}

func (e *explosion) update(dt float64) {
	for i := range e.waves {
		e.waves[i].update(dt)
	}
}

func (e *explosion) draw(out gogame.VideoOutput) {
	out.SetMask(gogame.Colors["white"])
	for i := range e.waves {
		e.waves[i].draw(out)
	}
}
