# RadialCells
#### Fast radius nearest neighbors search on 2d plane
It is based on my method of exact rasterization of a circle
Usage example (also in ./example folder):
```go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ogau/radialcells"
)

type Point = radialcells.Point

func randpoints(width, height float32, N int) []Point {
	points := make([]Point, N)
	seed := time.Now().Unix()
	fmt.Printf("[seed]: %v\n", seed)
	rnd := rand.New(rand.NewSource(seed))
	for i := 0; i < N; i++ {
		points[i].Row = rnd.Float32() * height
		points[i].Col = rnd.Float32() * width
	}
	return points
}

func main() {
	const width, height = 2.0, 1.0
	const gridstep = 0.05
	const npts = 500

	points := randpoints(width, height, npts)

	rc := radialcells.NewRadialCells(points, width, height, gridstep)

	var cx, cy float32 = 1.0, 0.5
	var radius float32 = 0.15

	result := rc.RadiusQuery(cx, cy, radius)
	for _, x := range result {
		pt := points[x.Index]
		fmt.Println(pt, x.Distance)
	}

	fmt.Println(len(result))
}
```

##### Rasterization algorithm:
There is a circle of a given radius relative to the point with the center (cx, cy).
The method projects points from (cx, cy) on radius horizontally and vertically.
Next, from each point in the clockwise direction, we pave path to stopline (end of quadrant).
For each cell containing the projected point, a checkpoint is set in one of the corners, which is checked for entering the circle.
For example, for the upper-left quadrant (green dots), the test point is set to the upper-right corner of the cell, and if it enters the circle, then the next cell for the path is the upper one, otherwise the right one.