package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fogleman/contourmap"
	"github.com/fogleman/gg"
)

const (
	N     = 8
	Scale = 1
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
	// z0 := m.Min
	// z1 := m.Max

	w := int(float64(m.W) * Scale)
	h := int(float64(m.H) * Scale)

	dc := gg.NewContext(w, h)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.Scale(Scale, Scale)

	zs := []float64{0.5}
	// for i := 1; i < N; i++ {
	for _, z := range zs {
		// t := float64(i) / N
		// z := z0 + (z1-z0)*t
		// z = 0.333
		// fmt.Println(i, t, z)
		fmt.Println(z)
		contours := m.Contours(z)
		for _, c := range contours {
			dc.NewSubPath()
			for _, p := range c {
				dc.LineTo(p.X, p.Y)
			}
		}
		// dc.SetColor(colormap.Viridis.At(t))
		// dc.FillPreserve()
		dc.SetRGB(0, 0, 0)
		dc.Stroke()
		// break
	}

	dc.SavePNG("out.png")
}
