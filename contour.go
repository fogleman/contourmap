package contourmap

// Contour is a list of Points which define an isoline.
// Contours may be open or closed. Closed contours have c[0] == c[len(c)-1].
type Contour []Point
