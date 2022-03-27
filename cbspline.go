// Copyright (c) 2022 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package cbspline

import "errors"

// CbSpline represents a cubic spline curve.
type CbSpline struct {
	*spline
}

// New creates and returns a cubic spline. t and v are the list of data points,
// and each element in them represents S(t[i]) = v[i]. t and v must be exactly
// the same length, and t must be sorted in ascending order.
//
// Elements of t with the same value can appear up to two times in a row, only
// if the corresponding v elements are also equivalent. If different v values
// are specified for the same t, or if the same t appears more than twice, it
// will probably return an error.
//
// BUG(tunabay): Currently b=Cyclic does not seem to work at all. Find out the
// correct algorithm.
func New(t, v []float64, b Boundary) (*CbSpline, error) {
	sp, err := newSpline(t, v, b)
	if err != nil {
		return nil, err
	}

	return &CbSpline{spline: sp}, nil
}

// At calculates and returns the S(t) value on the cubic spline. t must be in
// the range of minimum to maximum values in the t list passed to New(). The
// return value for t out of range is undefined.
func (sp *CbSpline) At(t float64) float64 { return sp.at(t) }

// ErrInvalidData is an error thrown when the given parameters are not
// appropriate and the cubic spline curve can not be determined.
var ErrInvalidData = errors.New("invalid data")
