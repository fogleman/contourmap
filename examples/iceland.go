package main

import (
	"image"
	"log"
	"os"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/fogleman/colormap"
	"github.com/fogleman/contourmap"
	"github.com/fogleman/gg"
)

const (
	N          = 12
	Scale      = 1
	Background = "77C4D3"
	Palette    = "70a80075ab007bb00080b30087b8008ebd0093bf009ac400a1c900a7cc00aed100b6d600bcd900c4de00cce300d2e600dbeb00e1ed00eaf200f3f700fafa00ffff05ffff12ffff1cffff29ffff36ffff42ffff4fffff5cffff66ffff73ffff80ffff8cffff99ffffa3ffffb0ffffbdffffc9ffffd6ffffe3ffffedfffffafcfcfcf7f7f7f5f5f5f0f0f0edededebebebe6e6e6e3e3e3dedededbdbdbd6d6d6d4d4d4cfcfcfccccccc7c7c7c4c4c4c2c2c2bdbdbdbababab5b5b5b3b3b3b3b3b3"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run iceland.go iceland.jpg")
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
	dc.SetColor(colormap.ParseColor(Background))
	dc.Clear()
	dc.Scale(Scale, Scale)

	options := vidio.Options{FPS: 1, Loop: 0, Delay: 50}
	video, err := vidio.NewVideoWriter("output.gif", w, h, &options)
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]uint8, w*h*3)
	img, _ := dc.Image().(*image.RGBA)

	pal := colormap.New(colormap.ParseColors(Palette))
	for i := 0; i < N; i++ {
		t := float64(i) / (N - 1)
		z := z0 + (z1-z0)*t
		contours := m.Contours(z + 1e-9)
		for _, c := range contours {
			dc.NewSubPath()
			for _, p := range c {
				dc.LineTo(p.X, p.Y)
			}
		}
		dc.SetColor(pal.At(t))
		dc.FillPreserve()
		dc.SetRGB(0, 0, 0)
		dc.SetLineWidth(1)
		dc.Stroke()

		if err == nil {
			fillBuffer(img, buffer)
			video.Write(buffer)
		}
	}

	dc.SavePNG("out.png")
}

func fillBuffer(im *image.RGBA, buffer []uint8) {
	index := 0
	for y := 0; y < im.Bounds().Dy(); y++ {
		for x := 0; x < im.Bounds().Dx(); x++ {
			r, g, b, _ := im.At(x, y).RGBA()
			buffer[index+0] = uint8(r >> 8)
			buffer[index+1] = uint8(g >> 8)
			buffer[index+2] = uint8(b >> 8)
			index += 3
		}
	}
}
