// Package main creates an image, gif or video.
package main

import (
	"fmt"
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
		render.Image(800, 800, "out.png", renderFrame, 0.5)
		render.ViewImage("out.png")
		break

	case target.Gif:
		render.Frames(400, 400, 60, "frames", renderFrame)
		render.MakeGIF("ffmpeg", "frames", "out.gif", 30)
		render.ViewImage("out.gif")
		break

	case target.Video:
		render.Frames(600, 600, 300, "frames", renderFrame)
		render.ConvertToVideo("frames", "out.mp4", 600, 600, 30)
		render.MPV("out.mp4", true)
		break
	}
}

var points []*geom.Point
var offsets []float64

func init() {
	for i := 0.0; i < 12; i++ {
		a := blmath.Tau * i / 12
		points = append(points, geom.NewPoint(math.Cos(a)*250, math.Sin(a)*250))
		offsets = append(offsets, random.FloatRange(0, blmath.Tau))
	}
}

func renderFrame(context *cairo.Context, width, height, percent float64) {
	context.ClearBlack()
	context.SetLineWidth(0.5)
	context.Save()
	context.TranslateCenter()

	a := blmath.Tau * percent * 2
	ops := []*geom.Point{}
	radius := 50.0

	for i := 0; i < len(points); i++ {
		point := points[i]
		offset := offsets[i]
		ops = append(ops, geom.NewPoint(point.X+math.Cos(a+offset)*radius, point.Y+math.Sin(a+offset)*radius))
	}

	pattern := cairo.CreateRadialGradient(60+math.Cos(a+math.Pi)*40, -100+math.Sin(a+math.Pi)*40, 0, 0, 0, 250)
	pattern.AddColorStopRGB(0, 1, 1, 1)
	pattern.AddColorStopRGB(1, 1, 0.5, 0)
	context.SetSource(pattern)
	context.FillMultiLoop(ops)
	context.FillText("genuary5 2023", -width/2+5, height/2-5)

	context.Rectangle(-width/2+blmath.LerpSin(percent, 0, 1)*width, -height/2, width, height)
	context.Clip()

	context.Blueprint()
	context.GridFull(20, 0.1)
	context.LabelPoints(points, true)
	context.Points(points, 2)
	context.Points(ops, 2)
	context.LabelPoints(ops, false)
	context.StrokeMultiLoop(ops)
	for i := 0; i < len(points); i++ {
		context.DisableDash()
		c := points[i]
		context.StrokeCircle(c.X, c.Y, radius)
		p := ops[i]
		context.StrokeLine(p.X, p.Y, c.X, c.Y)
		context.FillText(fmt.Sprintf("%.3f", math.Mod(a+offsets[i], blmath.Tau)), c.X-10, c.Y+20)
	}
	context.FillText("genuary5 2023", -width/2+5, height/2-5)
	context.SimpleDash(5, 5)
	context.StrokePath(ops, true)

	context.Restore()
}
