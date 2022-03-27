// Copyright (c) 2022 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package cbspline

import (
	"fmt"
	"math"
	"sort"
)

type spline struct {
	seg []*segment
}

// V(t)   = c[0]*(t-t0)^3 + c[1]*(t-t0)^2 + c[2]*(t-t0) + c[3]
// V'(t)  = 3*c[0]*(t-t0)^2 + 2*c[1]*(t-t0) + c[2]
// V''(t) = 6*c[0]*(t-t0) + 2*c[1]
// t = [t0, t1]
type segment struct {
	t0, t1, h float64
	v0, v1, d float64
	c         [4]float64
}

func (sp *spline) segAt(t float64) *segment {
	idx := sort.Search(len(sp.seg), func(i int) bool { return t < sp.seg[i].t0 })
	if idx != 0 {
		idx--
	}
	return sp.seg[idx]
}

func (sp *spline) at(t float64) float64 {
	switch {
	case t < sp.seg[0].t0:
		t = sp.seg[0].t0
	case sp.seg[len(sp.seg)-1].t1 < t:
		t = sp.seg[len(sp.seg)-1].t1
	}
	seg := sp.segAt(t)
	x := t - seg.t0

	return ((seg.c[0]*x+seg.c[1])*x+seg.c[2])*x + seg.c[3]
}

func newSpline(t []float64, v []float64, b Boundary) (*spline, error) {
	nKnots := len(t)
	if nKnots < 2 {
		return nil, fmt.Errorf("%w: require at least two knots", ErrInvalidData)
	}

	switch {
	case !sort.Float64sAreSorted(t):
		return nil, fmt.Errorf("%w: t is not sorted in increasing order", ErrInvalidData)
	case math.IsInf(t[0], 0), math.IsInf(t[nKnots-1], 0):
		return nil, fmt.Errorf("%w: t can not be Inf", ErrInvalidData)
	case math.IsNaN(t[0]):
		return nil, fmt.Errorf("%w: t can not be NaN", ErrInvalidData)
	}

	if len(v) != nKnots {
		return nil, fmt.Errorf("%w: v must have exactly %d values", ErrInvalidData, nKnots)
	}

	nSegs := nKnots - 1
	seg := make([]*segment, nSegs)
	for i := 0; i < nSegs; i++ {
		s := &segment{
			t0: t[i],
			t1: t[i+1],
			v0: v[i],
			v1: v[i+1],
		}
		s.h = s.t1 - s.t0
		s.d = (s.v1 - s.v0)
		if s.d != 0 {
			if s.h == 0 {
				return nil, fmt.Errorf("%w: #%d infinite dv/dt at t=%v", ErrInvalidData, i, s.t0)
			}
			s.d /= s.h
		}
		s.c[3] = s.v0

		seg[i] = s
	}

	// V(t)   = c[0]*(t-t0)^3 + c[1]*(t-t0)^2 + c[2]*(t-t0) + c[3]
	// V'(t)  = 3*c[0]*(t-t0)^2 + 2*c[1]*(t-t0) + c[2]
	// V''(t) = 6*c[0]*(t-t0) + 2*c[1]
	//
	// (d): c[i][0] = (c[i+1][1] - c[i][1]) / 3 * h[i]
	// (c): c[i][3] = (c[i+1][3] - c[i][3]) / h[i] - c[i][1]*h[i] - c[i][0]*(h[i]^2)
	// (b): c[i][2] = (c[i+1][3] - c[i][3]) / h[i] - c[i][1]*h[i] - c[i][0]*(h[i]^2)
	//              = (c[i+1][3] - c[i][3]) / h[i] - h[i]*(c[i+1][1] + 2*c[i][1]) / 3
	// (a): c[i][3] = v[i]

	m := make([]float64, nKnots*nKnots)
	r := make([]float64, nKnots)
	mSet := func(x, y int, v float64) { m[nKnots*y+x] = v }

	for i := 0; i < nKnots-2; i++ {
		mSet(i+0, i+1, seg[i].h)
		mSet(i+1, i+1, (seg[i].h+seg[i+1].h)*2.0)
		mSet(i+2, i+1, seg[i+1].h)
		r[i+1] = (seg[i+1].d - seg[i].d) * 6.0
	}

	switch b { // Boundary conditions
	case Natural:
		mSet(0, 0, 1)
		mSet(nKnots-1, nKnots-1, 1)
	case NotAKnot:
		mSet(0, 0, seg[1].h)
		mSet(1, 0, -seg[0].h-seg[1].h)
		mSet(2, 0, seg[0].h)
		mSet(nKnots-3, nKnots-1, seg[nSegs-1].h)
		mSet(nKnots-2, nKnots-1, -seg[nSegs-1].h-seg[nSegs-2].h)
		mSet(nKnots-1, nKnots-1, seg[nSegs-2].h)
	case Clamped:
		mSet(0, 0, seg[0].h*2.0)
		mSet(1, 0, seg[0].h)
		r[0] = (seg[0].d - 0) * 6.0
		mSet(nKnots-2, nKnots-1, seg[nSegs-1].h)
		mSet(nKnots-1, nKnots-1, seg[nSegs-1].h*2.0)
		r[nKnots-1] = (0 - seg[nSegs-1].d) * 6.0
	case Cyclic:
		mSet(nKnots-1, 0, seg[nSegs-1].h)
		mSet(0, 0, (seg[nSegs-1].h+seg[0].h)*2.0)
		mSet(1, 0, seg[0].h)
		r[0] = (seg[0].d - seg[nSegs-1].d) * 6.0
		mSet(nKnots-2, nKnots-1, seg[nSegs-1].h)
		mSet(nKnots-1, nKnots-1, (seg[nSegs-1].h+seg[0].h)*2.0)
		mSet(0, nKnots-1, seg[0].h)
		r[nKnots-1] = r[0]
	}

	x, err := solveX(m, r)
	if err != nil {
		return nil, fmt.Errorf("%w: unsolvable: %v", ErrInvalidData, err)
	}
	// fmt.Printf("x: %+v\n", x)

	for i, s := range seg {
		s.c[0] = (x[i+1] - x[i]) / 6.0 / s.h // this could be +/-Inf when s.h == 0
		s.c[1] = x[i] / 2.0
		s.c[2] = s.d - s.h*(2.0*x[i]+x[i+1])/6.0
	}

	// for i, s := range seg {
	// 	fmt.Printf("S(%d): %+v\n", i, s.c)
	// }
	sp := &spline{
		seg: seg,
	}

	return sp, nil
}
