## contourmap

Compute contour lines (isolines) for any 2D data in Go.

### Installation

    go get -u github.com/fogleman/contourmap
    
### Documentation

https://godoc.org/github.com/fogleman/contourmap

### Example Usage

#### Creating a ContourMap

A new `ContourMap` can be generated in many different ways, depending on what type of data you have.

Use `FromFloat64s` if you have an array of numbers. The length of the array must equal `width * height`. The two-dimensional data is stored in a flat array in row-major order.

```go
m := contourmap.FromFloat64s(width, height, data)
```

Use `FromImage` if you have an `image.Image`, such as a grayscale heightmap.

```go
m := contourmap.FromImage(im)
```

Use `FromFunction` to specify an arbitrary function that will provide a Z for any given X, Y coordinate.
The function will be called for all points `x = [0, w)` and `y = [0, h)` to determine the Z value at each point in the grid.

```go
var f func(x, y int) float64
...
m := contourmap.FromFunction(width, height, f)
```

#### Finding Contour Lines

Once your `ContourMap` is created, you can use the `Contours` function to find isolines at any given Z height. This function returns a list of contours where each contour is a list of X, Y points.
A `Contour` may be open or closed. Closed contours have `c[0] == c[len(c)-1]`.

```go
contours := m.Contours(z)
for _, contour := range contours {
    for _, point := range contour {
        // do something with points...
        fmt.Println(point.X, point.Y)
    }
}
```

#### Closing Contours at the Grid Perimeter

Contours may end at the edge of the grid data, forming open contours. If you want to force all contours to be closed by following the perimeter of the grid, you can use `ContourMap.Closed` which will generate a new ContourMap that can be used for this purpose:

```go
m = m.Closed() // now all contours will be closed paths
```

### Examples

Some examples are included to help you get started.

    $ cd go/src/github.com/fogleman/contourmap/examples
    $ go run iceland.go iceland.jpg

![Iceland Example](https://i.imgur.com/fd7fUnt.png)

    $ go run examples/function.go
    
![Function Example](https://i.imgur.com/lbGXPC9.png)

