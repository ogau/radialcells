// Copyright 2021 The radialcells Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Ilya Lisunov <elijahlisunov@gmail.com>

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ogau/radialcells"
)

type Point = radialcells.Point

func randpoints(width, height, N int) []Point {
	points := make([]Point, N)
	seed := time.Now().Unix()
	fmt.Printf("[seed]: %v\n", seed)
	rnd := rand.New(rand.NewSource(seed))
	for i := 0; i < N; i++ {
		points[i].Row = float32(rnd.Intn(height))
		points[i].Col = float32(rnd.Intn(width))
	}
	return points
}

func main() {
	const width, height = 2000, 1400
	const gridstep = 78
	const npts = 300000

	points := randpoints(width, height, npts)

	rc := radialcells.NewRadialCells(points, width, height, gridstep)

	cv := newCanvas(width, height, gridstep)
	cv.drawgrid()
	cv.drawpoints(points)

	draw := func(cx, cy, radius float32) {
		result := rc.RadiusQuery(cx, cy, radius)
		for _, x := range result {
			pt := points[x.Index]
			cv.drawpt(int(pt.Row), int(pt.Col), white)
		}
		fmt.Printf("%p %v %v\n", result, cap(result), len(result))
	}

	draw(width/2, height/2, 222)
	draw(-777/2, 0, 777)
	draw(width+314/2, height-314/2, 314)
	draw(6*gridstep, 13*gridstep, gridstep/3)
	draw(6*gridstep+gridstep/2, 15*gridstep+gridstep/2, gridstep/2)

	if err := saveImage(cv.img, "test.png"); err != nil {
		panic(err)
	}
}
