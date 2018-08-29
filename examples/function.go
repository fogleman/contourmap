package main

import (
	"math"

	"github.com/fogleman/colormap"
	"github.com/fogleman/contourmap"
	"github.com/fogleman/gg"
)

const (
	W     = 1024
	H     = 1024
	N     = 48
	Scale = 170
)

func f(i, j int) float64 {
	x := (float64(i) - W/2) / Scale
	y := (float64(j) - H/2) / Scale
	return math.Sin(1.3*x)*math.Cos(0.9*y) +
		math.Cos(.8*x)*math.Sin(1.9*y) +
		math.Cos(y*.2*x)
}

func main() {
	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	m := contourmap.FromFunction(W, H, f).Closed()
	z0 := m.Min
	z1 := m.Max
	for i := 0; i < N; i++ {
		t := float64(i) / (N - 1)
		z := z0 + (z1-z0)*t
		contours := m.Contours(z)
		for _, c := range contours {
			dc.NewSubPath()
			for _, p := range c {
				dc.LineTo(p.X, p.Y)
			}
		}
		dc.SetColor(colormap.Viridis.At(t))
		dc.FillPreserve()
		dc.SetRGB(0, 0, 0)
		dc.Stroke()
	}

	dc.SavePNG("out.png")
}
