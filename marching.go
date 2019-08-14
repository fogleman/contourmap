package contourmap

import (
	"math"
)

const closed = -math.MaxFloat64

type edge struct {
	X0, Y0   int
	X1, Y1   int
	Boundary bool
}

func fraction(z0, z1, z float64) float64 {
	const eps = 1e-9
	var f float64
	if z0 == closed {
		f = 0
	} else if z1 == closed {
		f = 1
	} else if z0 != z1 {
		f = (z - z0) / (z1 - z0)
	}
	f = math.Max(f, eps)
	f = math.Min(f, 1-eps)
	return f
}

func marchingSquares(m *ContourMap, w, h int, z float64) []Contour {
	edgePoint := make(map[edge]Point)
	nextEdge := make(map[Point]edge)
	for y := 0; y < h-1; y++ {
		up := m.at(0, y)
		lp := m.at(0, y+1)
		for x := 0; x < w-1; x++ {
			ul := up
			ur := m.at(x+1, y)
			ll := lp
			lr := m.at(x+1, y+1)

			up = ur
			lp = lr

			var squareCase int
			if ul > z {
				squareCase |= 1
			}
			if ur > z {
				squareCase |= 2
			}
			if ll > z {
				squareCase |= 4
			}
			if lr > z {
				squareCase |= 8
			}

			if squareCase == 0 || squareCase == 15 {
				continue
			}

			fx := float64(x)
			fy := float64(y)

			t := Point{fx + fraction(ul, ur, z), fy}
			b := Point{fx + fraction(ll, lr, z), fy + 1}
			l := Point{fx, fy + fraction(ul, ll, z)}
			r := Point{fx + 1, fy + fraction(ur, lr, z)}

			te := edge{x, y, x + 1, y, y == 0}
			be := edge{x, y + 1, x + 1, y + 1, y+2 == h}
			le := edge{x, y, x, y + 1, x == 0}
			re := edge{x + 1, y, x + 1, y + 1, x+2 == w}

			const connectHigh = false
			switch squareCase {
			case 1:
				edgePoint[te] = t
				nextEdge[t] = le
			case 2:
				edgePoint[re] = r
				nextEdge[r] = te
			case 3:
				edgePoint[re] = r
				nextEdge[r] = le
			case 4:
				edgePoint[le] = l
				nextEdge[l] = be
			case 5:
				edgePoint[te] = t
				nextEdge[t] = be
			case 6:
				if connectHigh {
					edgePoint[le] = l
					nextEdge[l] = te
					edgePoint[re] = r
					nextEdge[r] = be
				} else {
					edgePoint[re] = r
					nextEdge[r] = te
					edgePoint[le] = l
					nextEdge[l] = be
				}
			case 7:
				edgePoint[re] = r
				nextEdge[r] = be
			case 8:
				edgePoint[be] = b
				nextEdge[b] = re
			case 9:
				if connectHigh {
					edgePoint[te] = t
					nextEdge[t] = re
					edgePoint[be] = b
					nextEdge[b] = le
				} else {
					edgePoint[te] = t
					nextEdge[t] = le
					edgePoint[be] = b
					nextEdge[b] = re
				}
			case 10:
				edgePoint[be] = b
				nextEdge[b] = te
			case 11:
				edgePoint[be] = b
				nextEdge[b] = le
			case 12:
				edgePoint[le] = l
				nextEdge[l] = re
			case 13:
				edgePoint[te] = t
				nextEdge[t] = re
			case 14:
				edgePoint[le] = l
				nextEdge[l] = te
			}
		}
	}

	// pick out all boundary edgePoints
	boundaryEdgePoint := make(map[edge]Point)
	for e, p := range edgePoint {
		if e.Boundary {
			boundaryEdgePoint[e] = p
		}
	}

	var contours []Contour
	for len(edgePoint) > 0 {
		var contour Contour

		// find an unused edge; prefer starting at a boundary
		var e edge
		if len(boundaryEdgePoint) > 0 {
			for e = range boundaryEdgePoint {
				break
			}
		} else {
			for e = range edgePoint {
				break
			}
		}
		e0 := e

		// add the first point
		// (this allows closed paths to start & end at the same point)
		p := edgePoint[e]
		contour = append(contour, p)
		e = nextEdge[p]

		// follow points until none remain
		for {
			p, ok := edgePoint[e]
			if !ok {
				break
			}
			contour = append(contour, p)
			delete(edgePoint, e)
			if e.Boundary {
				delete(boundaryEdgePoint, e)
			}
			e = nextEdge[p]
		}

		// make sure the first one gets deleted in case of open paths
		delete(edgePoint, e0)
		if e0.Boundary {
			delete(boundaryEdgePoint, e0)
		}

		// add the contour
		contours = append(contours, contour)
	}

	return contours
}
