package contourmap

// Function returns a height Z for the specified X, Y point in a 2D grid.
type Function func(x, y int) float64
