package main

import "github.com/faiface/gogame"

type thing struct {
	color    gogame.Color
	position gogame.Vec
	waves    []wave
}

type wave struct {
	color           gogame.Color
	start, dir      gogame.Vec
	size, frequency float64
	time            float64
}

func newThing(color gogame.Color, position gogame.Vec) *thing {
	return &thing{
		color:    color,
		position: position,
	}
}

func (t *thing) draw(out gogame.VideoOutput) {
	bottom := out.OutputRect().Y + out.OutputRect().H
	left, right := t.position.X-(t.position.Y-bottom), t.position.X+(t.position.Y-bottom)

	out.DrawPolygon([]gogame.Vec{t.position, {X: left, Y: bottom}, {X: right, Y: bottom}}, 0, t.color)
}
