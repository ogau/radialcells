// Copyright 2021 The radialcells Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Ilya Lisunov <elijahlisunov@gmail.com>

package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

var (
	white   = color.RGBA{255, 255, 255, 255}
	gray    = color.RGBA{128, 128, 128, 255}
	blue    = color.RGBA{0, 0, 255, 255}
	bgblue  = color.RGBA{10, 10, 70, 255}
	green   = color.RGBA{0, 255, 0, 255}
	bggreen = color.RGBA{10, 70, 10, 255}
	red     = color.RGBA{255, 0, 0, 255}
	bgred   = color.RGBA{70, 10, 10, 255}
)

type canvas struct {
	img        *image.RGBA
	width      int
	height     int
	gridstep   int
	rows, cols int
}

func newCanvas(width, height, gridstep int) *canvas {
	cv := new(canvas)
	cv.img = image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cv.img.SetRGBA(x, y, color.RGBA{128, 128, 128, 255})
		}
	}
	cv.gridstep = gridstep
	cv.width, cv.height = width, height
	cv.rows = int(math.Ceil(float64(float32(height) / float32(gridstep))))
	cv.cols = int(math.Ceil(float64(float32(width) / float32(gridstep))))
	return cv
}

func (cv *canvas) drawgrid() {
	for y := 0; y < cv.rows; y++ {
		for x := 0; x < cv.cols; x++ {
			c := 128 - uint8((y+x)%2*45)
			cv.drawcell(y, x, color.RGBA{c, c, c, 255})
		}
	}
}

func (cv *canvas) drawpoints(points []Point) {
	for _, pt := range points {
		r, c := pt.Row, pt.Col
		cv.drawpt(int(r), int(c), blue)
	}
}

func (cv *canvas) drawcell(r, c int, val color.RGBA) {
	by, bx := r*cv.gridstep, c*cv.gridstep
	sy, sx := (r+1)*cv.gridstep, (c+1)*cv.gridstep
	for y := by; y < sy; y++ {
		for x := bx; x < sx; x++ {
			cv.img.Set(x, y, val)
		}
	}
}

func (cv *canvas) drawpt(r, c int, val color.RGBA) {
	cv.img.Set(c, r, val)
}

func saveImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
