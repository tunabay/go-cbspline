// Copyright (c) 2022 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package cbspline

import (
	"errors"
	"fmt"

	"gonum.org/v1/gonum/mat"
)

var errSolveFailure = errors.New("failed to solve equations")

// solveX solve the systems of equations [M] * [x] = [r] for x.
func solveX(m, r []float64) ([]float64, error) {
	n := len(r)
	switch {
	case n < 2:
		return nil, fmt.Errorf("%w: short vector: n=%d < 2", errSolveFailure, n)
	case len(m) != n*n:
		return nil, fmt.Errorf("%w: m must be %dx%d matrix", errSolveFailure, n, n)
	}

	lu := &mat.LU{}
	lu.Factorize(mat.NewDense(n, n, m))
	vecR := mat.NewVecDense(n, r)
	dst := &mat.VecDense{}
	if err := lu.SolveVecTo(dst, false, vecR); err != nil {
		return nil, fmt.Errorf("%w: %v", errSolveFailure, err)
	}
	if l := dst.Len(); l != n {
		return nil, fmt.Errorf("%w: unexpected result: %d != %d", errSolveFailure, l, n)
	}

	ret := make([]float64, n)
	for i := range ret {
		ret[i] = dst.At(i, 0)
	}

	return ret, nil
}
