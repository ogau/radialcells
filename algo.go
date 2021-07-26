// Copyright 2021 The radialcells Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Ilya Lisunov <elijahlisunov@gmail.com>

package radialcells

func (rc *RadialCells) radiusQueryGridIntersection(centerX, centerY, radius float32) cellsheap {
	const radiusEps = 1e-06
	radius -= radiusEps
	step := rc.step
	r2 := radius * radius
	visitadd := step / 2
	rc.heapcache.reset()

	vertstart := rc.anchor(centerX-radius, 1)
	vertend := rc.anchor(centerX+radius, 0)
	if vertstart < 0 {
		vertstart = 0 // start with gridstep
	}
	if vertend > rc.width {
		vertend = rc.width // without last vert
	}
	for ; vertstart <= vertend; vertstart += step {
		ytop, ybot := intersectAxis(vertstart, centerX, centerY, r2)
		vm := rc.val2grid(vertstart - visitadd)
		vp := rc.val2grid(vertstart + visitadd)
		if ytop >= 0 && ytop < rc.height {
			r := rc.val2grid(ytop)
			rc.heapcache.push(cell{r, vp})
			rc.heapcache.push(cell{r, vm})
		}
		if ybot >= 0 && ybot < rc.height {
			r := rc.val2grid(ybot)
			rc.heapcache.push(cell{r, vp})
			rc.heapcache.push(cell{r, vm})
		}
	}

	horzstart := rc.anchor(centerY-radius, 1)
	horzend := rc.anchor(centerY+radius, 0)
	if horzstart < 0 {
		horzstart = 0 // start with gridstep
	}
	if horzend > rc.height {
		horzend = rc.height // without last horz
	}
	for ; horzstart <= horzend; horzstart += step {
		xleft, xright := intersectAxis(horzstart, centerY, centerX, r2)
		vp := rc.val2grid(horzstart + visitadd)
		vm := rc.val2grid(horzstart - visitadd)
		if xleft >= 0 && xleft < rc.width {
			c := rc.val2grid(xleft)
			rc.heapcache.push(cell{vp, c})
			rc.heapcache.push(cell{vm, c})
		}
		if xright >= 0 && xright < rc.width {
			c := rc.val2grid(xright)
			rc.heapcache.push(cell{vp, c})
			rc.heapcache.push(cell{vm, c})
		}
	}
	return rc.heapcache
}

func (rc *RadialCells) drawedge(centerX, centerY, radius float32) cellsheap {
	// reset cache
	rc.heapcache.reset()

	const radiusEps = 1e-06
	radius -= radiusEps
	r2 := radius * radius
	gridstep := rc.step
	visitShift := gridstep / 2 // roll by axis

	var anchorX, anchorY, stoplineX, stoplineY float32

	// left to top cell
	anchorY, anchorX = rc.anchor(centerY, 0), rc.anchor(centerX-radius, 1)
	stoplineX = rc.anchor(centerX, 1) + visitShift // next anchor (top) with shift
	for {
		if cmpsqd(anchorX-centerX, anchorY-centerY, r2) {
			anchorY -= gridstep
		} else {
			anchorX += gridstep
		}
		if anchorX > stoplineX {
			break
		}
		rc.heapcache.push(rc.pointAsCell(anchorY+visitShift, anchorX-visitShift))
	}

	// top to right cell
	anchorY, anchorX = rc.anchor(centerY-radius, 1), rc.anchor(centerX, 1)
	stoplineY = rc.anchor(centerY, 1) + visitShift // next anchor (right) with shift
	for {
		if cmpsqd(anchorX-centerX, anchorY-centerY, r2) {
			anchorX += gridstep
		} else {
			anchorY += gridstep
		}
		if anchorY > stoplineY {
			break
		}
		rc.heapcache.push(rc.pointAsCell(anchorY-visitShift, anchorX-visitShift))
	}

	// right to bot cell
	anchorY, anchorX = rc.anchor(centerY, 1), rc.anchor(centerX+radius, 0)
	stoplineX = rc.anchor(centerX, 0) - visitShift // next anchor (bot) with shift
	for {
		if cmpsqd(anchorX-centerX, anchorY-centerY, r2) {
			anchorY += gridstep
		} else {
			anchorX -= gridstep
		}
		if anchorX < stoplineX {
			break
		}
		rc.heapcache.push(rc.pointAsCell(anchorY-visitShift, anchorX+visitShift))
	}

	// bot to left cell
	anchorY, anchorX = rc.anchor(centerY+radius, 0), rc.anchor(centerX, 0)
	stoplineY = rc.anchor(centerY, 0) - visitShift // next anchor (left) with shift
	for {
		if cmpsqd(anchorX-centerX, anchorY-centerY, r2) {
			anchorX -= gridstep
		} else {
			anchorY -= gridstep
		}
		if anchorY < stoplineY {
			break
		}
		rc.heapcache.push(rc.pointAsCell(anchorY+visitShift, anchorX+visitShift))
	}

	return rc.heapcache
}

func (rc *RadialCells) RadiusQuery(centerX, centerY, radius float32) []DIndex {
	if radius < 0 {
		panic("radius < 0")
	}
	return rc.radiusQuery(centerX, centerY, radius)
}

func (rc *RadialCells) radiusQuery(centerX, centerY, radius float32) []DIndex {
	// reset cache
	rc.result = rc.result[:0]

	// get edge cells in a circle
	edgecells := rc.drawedge(centerX, centerY, radius)

	// at end override var for reuse slice
	defer func() { rc.heapcache = edgecells }()

	// draw four neighboring cells for small radius
	if radius < rc.step {
		edgecells.push(rc.pointAsCell(centerY, centerX-radius))
		edgecells.push(rc.pointAsCell(centerY-radius, centerX))
		edgecells.push(rc.pointAsCell(centerY, centerX+radius))
		edgecells.push(rc.pointAsCell(centerY+radius, centerX))
	}

	// sorting and deduplication cells in place
	edgecells.sortDeduplicateInplace()

	r2 := radius * radius
	for _, cell := range edgecells {
		if rc.grid.notInBounds(cell) {
			continue
		}

		start, end := rc.indCellAt(cell)
		for i, pt := range rc.points[start:end] {
			dy := pt.Row - centerY
			dx := pt.Col - centerX
			d2 := dy*dy + dx*dx
			if d2 < r2 {
				rc.result = append(rc.result, DIndex{i + start, fastSqrt0(d2)})
			}
		}
	}
	// inline tool
	split := func(x cell) (int, int) {
		return x.r, x.c
	}
	prow, pcol := split(edgecells[0])
	for i := 1; i < len(edgecells); i++ {
		row, col := split(edgecells[i])
		if row < 0 {
			continue
		} else if row >= rc.rows {
			break
		} else if col-pcol > 1 && prow == row { // check if gap in one row
			startcol, beglim := clamp(pcol+1, 0, rc.cols-1) // start col
			endcol, endlim := clamp(col-1, 0, rc.cols-1)    // end col
			if beglim && endlim && startcol == endcol {
				continue
			}

			start, end := rc.indCellRange(startcol, endcol, row)
			for i, pt := range rc.points[start:end] {
				dy := pt.Row - centerY
				dx := pt.Col - centerX
				rc.result = append(rc.result, DIndex{i + start, fastSqrt0(dy*dy + dx*dx)})
			}
		}
		prow, pcol = row, col
	}

	return rc.result
}
