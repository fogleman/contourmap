package contourmap

import (
	"image"
	"math"
	"sort"
)

type ContourMap struct {
	W    int     // width of the contour map in pixels
	H    int     // height of the contour map in pixels
	Min  float64 // minimum value contained in this contour map
	Max  float64 // maximum value contained in this contour map
	grid []float64
}

// FromFloat64s returns a new ContourMap for the provided 2D grid of values.
// len(grid) must equal w * h.
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

// FromFloat64s returns a new ContourMap for the provided function.
// The function will be called for all points x = [0, w) and y = [0, h) to
// determine the Z value at each point.
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

// FromImage returns a new ContourMap for the provided image. The image is
// converted to 16-bit grayscale and will have Z values mapped from
// [0, 65535] to [0, 1].
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

func (m *ContourMap) HistogramZs(numLevels int) []float64 {
	// compute histogram
	hist := make(map[float64]int)
	for _, v := range m.grid {
		hist[v]++
	}

	// sort histogram keys
	keys := make([]float64, 0, len(hist))
	for key := range hist {
		keys = append(keys, key)
	}
	sort.Float64s(keys)

	result := make([]float64, numLevels)
	numPixels := len(m.grid)
	for i := 0; i < numLevels; i++ {
		// compute number of pixels for this level
		t := (float64(i) + 0.5) / float64(numLevels)
		pixelCount := int(t * float64(numPixels))
		// find z
		var total int
		for _, k := range keys {
			total += hist[k]
			if total >= pixelCount {
				result[i] = k
				break
			}
		}
	}
	return result
}

// Contours returns a list of contours that represent isolines at the specified
// Z value.
func (m *ContourMap) Contours(z float64) []Contour {
	return marchingSquares(m, m.W, m.H, z)
}

// Closed returns a new ContourMap that will ensure all Contours are closed
// paths by following the border when they would normally stop at the edge
// of the grid.
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
