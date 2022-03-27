// Copyright (c) 2022 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package cbspline

import "fmt"

// Boundary specifies the boundary condition of a cubic spline.
type Boundary uint8

const (
	Natural Boundary = iota
	NotAKnot
	Clamped
	Cyclic // periodic
)

var boundaryStr = map[Boundary]string{
	Natural:  "natural",
	NotAKnot: "not-a-knot",
	Clamped:  "clamped",
	Cyclic:   "cyclic",
}

// String returns the string representation of Boundary.
func (b Boundary) String() string {
	if s, ok := boundaryStr[b]; ok {
		return s
	}
	return fmt.Sprintf("unknown(%d)", uint8(b))
}
