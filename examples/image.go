package main

import (
	"log"
	"os"

	"github.com/fogleman/colormap"
	"github.com/fogleman/contourmap"
	"github.com/fogleman/gg"
)

const (
	N     = 20
	Scale = 2
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run image.go input.png")
	}

	im, err := gg.LoadImage(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	m := contourmap.FromImage(im).Closed()
	z0 := m.Min
	z1 := m.Max

	w := int(float64(m.W) * Scale)
	h := int(float64(m.H) * Scale)

	dc := gg.NewContext(w, h)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.Scale(Scale, Scale)

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
