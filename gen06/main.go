// Package main creates an image, gif or video.
package main

import (
	"math"

	"github.com/bit101/bitlib/blmath"
	"github.com/bit101/bitlib/random"
	cairo "github.com/bit101/blcairo"
	"github.com/bit101/blcairo/render"
	"github.com/bit101/blcairo/target"
)

func main() {
	renderTarget := target.Video

	switch renderTarget {
	case target.Image:
		render.Image(600, 600, "out.png", renderFrame, 0.5)
		render.ViewImage("out.png")
		break

	case target.Gif:
		render.Frames(400, 400, 60, "frames", renderFrame)
		render.MakeGIF("ffmpeg", "frames", "out.gif", 400, 400, 30, 2)
		render.ViewImage("out.gif")
		break

	case target.Video:
		render.Frames(600, 600, 600, "frames", renderFrame)
		render.ConvertToVideo("frames", "out.mp4", 600, 600, 30, 20)
		render.PlayVideo("out.mp4")
		break
	}
}

type slice struct {
	angles  []float64
	lengths []float64
}

func newSlice() *slice {
	s := slice{}
	length := 120.0
	for i := 0; i < 6; i++ {
		s.angles = append(s.angles, random.FloatRange(0, -math.Pi/2))
		s.lengths = append(s.lengths, length)
		length *= 0.8
	}
	return &s
}

func (s *slice) tweak(amt float64) {
	for i := 0; i < len(s.angles); i++ {
		s.angles[i] += random.FloatRange(-amt, amt)
	}
}

func (s *slice) render(context *cairo.Context) {
	x := 50.0
	y := 0.0
	r := 5.0
	context.Circle(x, y, r)
	for i := 0; i < len(s.angles); i++ {
		context.MoveTo(x, y)
		angle := s.angles[i]
		length := s.lengths[i]
		r *= 0.9
		x += math.Cos(angle) * length
		y += math.Sin(angle) * length
		context.LineTo(x, y)
		context.Circle(x, y, r)
	}
}

var slices []*slice

func init() {
	for i := 0; i < 20; i++ {
		slices = append(slices, newSlice())
	}
}

func renderFrame(context *cairo.Context, width, height, percent float64) {
	context.BlackOnWhite()
	context.Save()
	context.Translate(width/2, height*0.8)
	context.SetLineWidth(0.5)
	for i := 0; i < len(slices); i++ {
		context.Save()
		s := slices[i]
		t := float64(i)
		l := float64(len(slices))
		a := t/l*blmath.Tau + percent*blmath.Tau
		scale := math.Sin(a)
		context.Scale(scale+0.0001, 1)
		s.render(context)
		context.Restore()
		context.Stroke()
		s.tweak(0.025)
	}
	context.Restore()
	context.SetSourceGray(0.4)
	context.FillText("genuary6 2023", 5, height-20)
	context.FillText("stolen from Jared Tarbell and Josh Davis", 5, height-5)
}
