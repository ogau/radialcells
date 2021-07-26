// Copyright 2021 The radialcells Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Ilya Lisunov <elijahlisunov@gmail.com>

package radialcells

import (
	"fmt"
	"math"
	"sort"
	"unsafe"
)

type (
	cell   struct{ r, c int }
	Point  struct{ Row, Col float32 }
	DIndex struct {
		Index    int
		Distance float32
	}
)

type grid struct {
	points  []Point
	indices []int
	offsets []int
	sizes   []int
	width   float32
	height  float32
	rows    int
	cols    int
	step    float32 // or cell size
}

func buildGrid(points []Point, width, height, gridstep float32) *grid {
	g := new(grid)

	g.width, g.height = width, height
	g.rows = int(math.Ceil(float64(height / gridstep)))
	g.cols = int(math.Ceil(float64(width / gridstep)))
	g.step = gridstep

	g.points = points
	g.indices = make([]int, len(points))
	g.offsets = make([]int, g.rows*g.cols)
	g.sizes = make([]int, g.rows*g.cols)

	g.sortpoints()
	return g
}

func (g *grid) Len() int { return len(g.indices) }
func (g *grid) Swap(i, j int) {
	g.points[i], g.points[j] = g.points[j], g.points[i]
	g.indices[i], g.indices[j] = g.indices[j], g.indices[i]
}
func (g *grid) Less(i, j int) bool { return g.indices[i] < g.indices[j] }

func (g *grid) sortpoints() {
	for i, pt := range g.points {
		if pt.Row < 0 || pt.Row >= g.height ||
			pt.Col < 0 || pt.Col >= g.width {
			panic(fmt.Sprintf("point %v out of bounds", pt))
		}
		r, c := g.point2gcoord(pt.Row, pt.Col)
		ind := g.at(r, c)
		if ind == 812 {
			fmt.Println(ind, r, c, pt, i)
		}
		g.indices[i] = ind
		g.sizes[ind]++
	}
	sort.Sort(g)
	var cum int
	for i, x := range g.sizes {
		g.offsets[i] = cum
		cum += x
	}
}
func (g *grid) anchor(val, d float32) float32 {
	return (floorf(val/g.step) + d) * g.step
}

// convert a value to a grid plane
func (g *grid) val2grid(val float32) int {
	return floor(val / g.step)
}

// convert point to grid coordinate plane
func (g *grid) point2gcoord(row, col float32) (int, int) {
	return g.val2grid(row), g.val2grid(col)
}

// convert point to cell in grid coordinate plane
func (g *grid) pointAsCell(row, col float32) cell {
	return cell{g.val2grid(row), g.val2grid(col)}
}

func (g *grid) notInBounds(c cell) bool {
	return c.r < 0 || c.c < 0 || c.r >= g.rows || c.c >= g.cols
}
func (g *grid) at(row, col int) int {
	return row*g.cols + col
}

// returns (start, end] indices for points array
func (g *grid) indCellAt(c cell) (int, int) {
	ind := c.r*g.cols + c.c
	off := g.offsets[ind]
	n := g.sizes[ind]
	return off, off + n
}

// returns (start, end] indices for points array
func (g *grid) indCellRange(start, end, row int) (int, int) {
	base := row * g.cols
	startcell := base + start
	endcell := base + end

	offstart := g.offsets[startcell]
	offend := g.offsets[endcell] + g.sizes[endcell]
	return offstart, offend
}

type RadialCells struct {
	*grid
	heapcache cellsheap
	result    []DIndex
}

func NewRadialCells(points []Point, width, height, step float32) *RadialCells {
	rc := new(RadialCells)

	rc.grid = buildGrid(points, width, height, step)
	rc.result = make([]DIndex, 0)
	rc.heapcache = make(cellsheap, 0)

	return rc
}

// func floor(x float32) int {
// 	return int(x+1<<15) - 1<<15
// }
func floor(x float32) int {
	return int(float64(x)+float64(1<<15)) - 1<<15
}
func floorf(x float32) float32 {
	return float32(floor(x))
}

func cmpsqd(dx, dy, r2 float32) bool {
	return dx*dx+dy*dy < r2
}

func intersectAxis(x, centerX, centerY, r2 float32) (float32, float32) {
	x -= centerX
	v := fastSqrt0(r2 - x*x)
	return -v + centerY, v + centerY
}

func sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}
func fastSqrt0(x float32) float32 {
	xhalf := float32(0.5) * x
	i := *(*int32)(unsafe.Pointer(&x))
	i = int32(0x5f3759df) - int32(i>>1)
	x = *(*float32)(unsafe.Pointer(&i))
	x = x * (1.5 - (xhalf * x * x))
	return 1 / x
}
func fastSqrt1(x float32) float32 {
	i := *(*int32)(unsafe.Pointer(&x))
	i = (1 << 29) + (i >> 1) - (1 << 22)
	ux := *(*float32)(unsafe.Pointer(&i))

	ux = ux + x/ux
	ux = 0.25*ux + x/ux
	return 1 / ux
}

// limiter
func clamp(x, min, max int) (int, bool) {
	if x < min {
		return min, true
	} else if x > max {
		return max, true
	}
	return x, false
}
