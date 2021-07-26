// Copyright 2021 The radialcells Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Written by Ilya Lisunov <elijahlisunov@gmail.com>

package radialcells

type cellsheap []cell

func (h *cellsheap) reset() {
	*h = (*h)[:0]
}

func (h cellsheap) less(i, j int) bool {
	if h[i].r == h[j].r {
		return h[i].c < h[j].c
	}
	return h[i].r < h[j].r
}
func (h cellsheap) equal(i, j int) bool {
	return h[i].r == h[j].r && h[i].c == h[j].c
}
func (h cellsheap) swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h cellsheap) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h cellsheap) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(j, i) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}

func (h *cellsheap) push(x cell) {
	*h = append(*h, x)
	h.up(len(*h) - 1)
}

func (h *cellsheap) pop() cell {
	h_ := *h
	n := len(h_) - 1
	h_.swap(0, n)
	h_.down(0, n)

	x := h_[n]
	*h = h_[:n]
	return x
}

func (h cellsheap) condswap(i, j int) {
	if !h.less(i, j) {
		h.swap(i, j)
	}
}
func (h cellsheap) sort4() {
	if len(h) > 4 {
		panic("len(h) > 4")
	}
	h.condswap(0, 1)
	h.condswap(2, 3)
	h.condswap(0, 2)
	h.condswap(1, 3)
	h.condswap(1, 2)
}

func (h *cellsheap) deduplicate() {
	j := 0
	h_ := *h
	for i := 1; i < len(h_); i++ {
		if h.equal(i, j) {
			continue
		}
		j++
		h_[j] = h_[i]
	}
	(*h) = (*h)[:j+1]
}

func (h_ *cellsheap) sortDeduplicateInplace() {
	h := *h_
	ln := len(h)
	// move smallest to end
	// in descending order
	for n := ln - 1; n >= 0; n-- {
		h.swap(0, n)
		h.down(0, n)
	}
	// reverse slice
	for i := 0; i < ln/2; i++ {
		j := ln - 1 - i
		h[i], h[j] = h[j], h[i]
	}
	h_.deduplicate()
}
