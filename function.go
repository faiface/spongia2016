package main

import (
	"math"

	"github.com/faiface/gogame"
)

type function struct {
	times  []float64
	points []gogame.Vec
	cyclic bool
}

func (f *function) add(time float64, point gogame.Vec) {
	f.times = append(f.times, time)
	f.points = append(f.points, point)
}

func (f *function) at(time float64) gogame.Vec {
	if f.cyclic {
		duration := f.times[len(f.times)-1]
		time = math.Mod(math.Mod(time, duration)+duration, duration)
	}

	if time <= 0 {
		return f.points[0]
	}
	for i, t := range f.times[1:] {
		i = i + 1
		a, b := f.points[i-1], f.points[i]
		if time <= t {
			fraction := (time - f.times[i-1]) / (f.times[i] - f.times[i-1])
			return a.M(1 - fraction).A(b.M(fraction))
		}
	}
	return f.points[len(f.points)-1]
}
