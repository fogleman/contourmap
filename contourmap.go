package contourmap

import (
	"image"
	"math"
)

type ContourMap struct {
	W, H int
	Min  float64
	Max  float64
	grid []float64
}

func FromFloat64s(w, h int, grid []float64) *ContourMap {
	min := math.Inf(1)
	max := math.Inf(-1)
	for _, x := range grid {
		if x == closed {
			continue
		}
		min = math.Min(min, x)
		max = math.Max(max, x)
	}
	return &ContourMap{w, h, min, max, grid}
}

func FromFunction(w, h int, f Function) *ContourMap {
	grid := make([]float64, w*h)
	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			grid[i] = f(x, y)
			i++
		}
	}
	return FromFloat64s(w, h, grid)
}

func FromImage(im image.Image) *ContourMap {
	gray := imageToGray16(im)
	w := gray.Bounds().Size().X
	h := gray.Bounds().Size().Y
	grid := make([]float64, w*h)
	j := 0
	for i := range grid {
		x := int(gray.Pix[j])<<8 | int(gray.Pix[j+1])
		grid[i] = float64(x) / 0xffff
		j += 2
	}
	return FromFloat64s(w, h, grid)
}

func (m *ContourMap) at(x, y int) float64 {
	return m.grid[y*m.W+x]
}

func (m *ContourMap) Contours(z float64) []Contour {
	return marchingSquares(m, m.W, m.H, z)
}

func (m *ContourMap) Closed() *ContourMap {
	w := m.W + 2
	h := m.H + 2
	grid := make([]float64, w*h)
	for i := range grid {
		grid[i] = closed
	}
	for y := 0; y < m.H; y++ {
		i := (y+1)*w + 1
		j := y * m.W
		copy(grid[i:], m.grid[j:j+m.W])
	}
	return FromFloat64s(w, h, grid)
}
