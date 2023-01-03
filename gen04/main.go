// Package main creates an image, gif or video.
package main

import (
	"math"

	"github.com/bit101/bitlib/blmath"
	"github.com/bit101/bitlib/geom"
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
		render.MakeGIF("ffmpeg", "frames", "out.gif", 30)
		render.ViewImage("out.gif")
		break

	case target.Video:
		render.Frames(600, 600, 240, "frames", renderFrame)
		render.ConvertToVideo("frames", "out.mp4", 600, 600, 30)
		render.MPV("out.mp4", true)
		break
	}
}

func renderFrame(context *cairo.Context, width, height, percent float64) {
	context.BlackOnWhite()
	context.SetLineWidth(0.5)
	context.DisableDash()
	context.FillText("genuary4 2023", 5, height-5)

	a := percent * blmath.Tau

	c0 := geom.NewCircle(150, 450, 100)
	c1 := geom.NewCircle(170+math.Cos(a)*80, 420+math.Sin(a)*80, 60)
	p0, p1, _ := c0.Intersection(c1)

	context.StrokeCircleObject(c0)
	context.StrokeCircleObject(c1)

	c2 := geom.NewCircle(450, 150, 100)
	c3 := geom.NewCircle(460+math.Cos(a+1.3)*80, 140+math.Sin(a+1.3)*80, 40)
	p2, p3, _ := c2.Intersection(c3)

	context.StrokeCircleObject(c2)
	context.StrokeCircleObject(c3)

	context.SimpleDash(5, 5)
	context.LineThrough(p0.X, p0.Y, p1.X, p1.Y, 800)
	context.LineThrough(p2.X, p2.Y, p3.X, p3.Y, 800)

	x, y, _ := geom.LineOnLine(p0.X, p0.Y, p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y)
	context.SetSourceRGB(1, 0, 0)
	context.FillCircle(x, y, 10)
}
