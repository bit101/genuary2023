// Package main creates an image, gif or video.
package main

import (
	"math"
	"sort"

	"github.com/bit101/bitlib/blcolor"
	"github.com/bit101/bitlib/blmath"
	"github.com/bit101/bitlib/random"
	cairo "github.com/bit101/blcairo"
	"github.com/bit101/blcairo/render"
	"github.com/bit101/blcairo/target"
)

const frames = 300.0

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
		render.Frames(600, 600, frames, "frames", renderFrame)
		render.ConvertToVideo("frames", "out.mp4", 600, 600, 30)
		render.MPV("out.mp4", true)
		break
	}
}

type particle struct {
	x, y, z, cz float64
	color       blcolor.Color
	radius      float64
	fl          float64
}

func newParticle(x, y, z, cz float64, color blcolor.Color) *particle {
	return &particle{x, y, z, cz, color, random.FloatRange(4, 16), 300}
}

func (p *particle) project() (float64, float64, float64) {
	scale := p.fl / (p.fl + p.z + p.cz)
	return p.x * scale, p.y * scale, scale
}

func (p *particle) rotate(t float64) {
	cos := math.Cos(t)
	sin := math.Sin(t)
	x := p.x*cos - p.z*sin
	z := p.z*cos + p.x*sin
	p.x = x
	p.z = z
}

type particleArray []*particle

func (a particleArray) Len() int {
	return len(a)
}

func (a particleArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a particleArray) Less(i, j int) bool {
	return a[i].z > a[j].z
}

var particles []*particle

func init() {
	colors := cairo.SampleColors("easter.png", 32)
	zc := 300.0
	for i := 0; i < 10000; i++ {
		c := random.IntRange(0, len(colors))
		r := random.FloatRange(100, 300)
		t := random.FloatRange(0, blmath.Tau)
		p := newParticle(math.Sin(t)*r, random.FloatRange(-500, 500), math.Cos(t)*r, zc, colors[c])
		particles = append(particles, p)
	}
}

func renderFrame(context *cairo.Context, width, height, percent float64) {
	context.WhiteOnBlack()
	context.Save()
	context.TranslateCenter()
	for _, p := range particles {
		p.rotate(1.0 / frames * blmath.Tau)
	}
	sort.Sort(particleArray(particles))

	for _, p := range particles {
		x, y, scale := p.project()
		context.SetSourceColor(p.color)
		context.FillCircle(x, y, p.radius*scale)
	}

	context.Restore()
	context.SetSourceWhite()
	context.FillText("genuary7 2023", 5, height-20)
	context.FillText("palette sampled from Easter by the Patti Smith Group", 5, height-5)
}
