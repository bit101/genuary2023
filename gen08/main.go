// Package main creates an image, gif or video.
package main

import (
	"math"

	"github.com/bit101/bitlib/blmath"
	"github.com/bit101/bitlib/geom"
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
		render.Frames(600, 600, 60, "frames", renderFrame)
		render.MakeGIF("ffmpeg", "frames", "out.gif", 30)
		render.ViewImage("out.gif")
		break

	case target.Video:
		render.Frames(600, 600, 150, "frames", renderFrame)
		render.ConvertToVideo("frames", "out.mp4", 600, 600, 30)
		render.MPV("out.mp4", true)
		break
	}
}

var circles []*geom.Circle
var offsets []float64

func init() {
	random.Seed(5)
	for i := 0; i < 7; i++ {
		p := geom.NewCircle(random.FloatRange(0, 600), random.FloatRange(0, 600), 100)
		circles = append(circles, p)
		offsets = append(offsets, random.FloatRange(0, blmath.Tau))
	}
}

func renderFrame(context *cairo.Context, width, height, percent float64) {
	context.BlackOnWhite()
	a := percent * blmath.Tau
	context.ProcessPixels(func(context *cairo.Context, x, y float64) {
		s := 0.0
		for _, c := range circles {
			sdf := sdfCircle(c, x, y)
			s += math.Sin(sdf/30.0 - a)
		}
		v := blmath.Map(s, -7, 7, 0, 1)
		context.SetSourceRGB(v, v, v)

		context.FillRectangle(x, y, 1, 1)
	})

}

func sdfCircle(c *geom.Circle, x, y float64) float64 {
	d := math.Hypot(x-c.X, y-c.Y)
	return d - c.Radius
}
