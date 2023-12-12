// Package main creates an image, gif or video.
package main

import (
	"math"

	"github.com/bit101/bitlib/blmath"
	"github.com/bit101/bitlib/geom"
	"github.com/bit101/bitlib/noise"
	"github.com/bit101/bitlib/random"
	cairo "github.com/bit101/blcairo"
	"github.com/bit101/blcairo/render"
	"github.com/bit101/blcairo/target"
)

func main() {
	renderTarget := target.Video

	switch renderTarget {
	case target.Image:
		render.Image(800, 800, "out.png", renderFrame, 0.5)
		render.ViewImage("out.png")
		break

	case target.Gif:
		render.Frames(400, 400, 60, "frames", renderFrame)
		render.MakeGIF("ffmpeg", "frames", "out.gif", 400, 400, 30, 2)
		render.ViewImage("out.gif")
		break

	case target.Video:
		w := 600.0
		h := 600.0
		render.Frames(w, h, 240, "frames", renderFrame)
		render.ConvertToVideo("frames", "out.mp4", 400, 400, 30, 8)
		render.PlayVideo("out.mp4")
		break
	}
}

func renderFrame(context *cairo.Context, width, height, percent float64) {
	// random.Seed(0)
	context.BlackOnWhite()
	context.FillText("genuary3 2023", 5, height-5)
	context.Save()
	context.TranslateCenter()
	r := width / 4
	c := geom.NewCircle(0, 0, r)
	t := blmath.LerpSin(percent, 0, 1)
	fillCircle(context, c)
	context.SetOperator(cairo.OperatorDifference)
	context.SetSourceWhite()

	count := 10.0
	outer := geom.OuterCircles(c, int(count), blmath.Tau/count*percent*3)
	inner := geom.InnerCircles(c, int(count), blmath.Tau/count*percent*3)
	lerped := lerpCircles(t, inner, outer)
	fillCircles(context, lerped)

	count = 20.0
	outer = geom.OuterCircles(c, int(count), -blmath.Tau/count*percent*6)
	inner = geom.InnerCircles(c, int(count), -blmath.Tau/count*percent*6)
	lerped = lerpCircles(1-t, inner, outer)
	fillCircles(context, lerped)
	context.Restore()

	context.SetOperator(cairo.OperatorOver)
	glitch(context)
}

func lerpCircles(t float64, inner, outer []*geom.Circle) []*geom.Circle {
	circles := []*geom.Circle{}
	for i := 0; i < len(inner); i++ {
		in := inner[i]
		out := outer[i]
		c := geom.NewCircle(
			blmath.Lerp(t, in.X, out.X),
			blmath.Lerp(t, in.Y, out.Y),
			blmath.Lerp(t, in.Radius, out.Radius),
		)
		circles = append(circles, c)
	}
	return circles
}

func fillCircles(context *cairo.Context, circles []*geom.Circle) {
	for _, c := range circles {
		fillCircle(context, c)
	}

}

func fillCircle(context *cairo.Context, circle *geom.Circle) {
	s := 0.01
	o := 20.0
	for t := 0.0; t < blmath.Tau; t += 0.01 {
		x := circle.X + math.Cos(t)*circle.Radius
		y := circle.Y + math.Sin(t)*circle.Radius
		a := noise.Simplex2(x*s, y*s)
		x += math.Cos(a) * o
		y += math.Sin(a) * o
		context.LineTo(x, y)

	}
	context.Fill()
}

func glitch(context *cairo.Context) {
	for y := 0.0; y < context.Height; y++ {
		if random.WeightedBool(0.2) {
			context.SetSourceSurface(context.Surface, random.FloatRange(-20, 0), 0)
			context.FillRectangle(0, y, context.Width, 1)
		}
	}
}
