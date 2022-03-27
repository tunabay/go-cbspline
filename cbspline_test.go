// Copyright (c) 2022 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package cbspline_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/tunabay/go-cbspline"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func TestCbSpline_0(t *testing.T) {
	ts := []float64{0, 1, 3, 4}
	vs := []float64{0, 0, 2, 2}
	testPlotSpline(t, ts, vs, -0.5, 4.5)
}

func TestCbSpline_1(t *testing.T) {
	ts := []float64{0, 1, 4, 5, 8}
	vs := []float64{0, 3, 4, 1, 2}
	testPlotSpline(t, ts, vs, -0.5, 8.5)
}

func TestCbSpline_2(t *testing.T) {
	ts := []float64{0, 1, 1.5, 2.5, 4, 5, 5.5, 7}
	vs := []float64{1, 5, 2.5, 2.5, 6, 4, 6, 1}
	testPlotSpline(t, ts, vs, -0.5, 7.5)
}

func testPlotSpline(t *testing.T, ts, vs []float64, xMin, xMax float64) {
	t.Helper()

	bs := []cbspline.Boundary{
		cbspline.Natural,
		cbspline.NotAKnot,
		cbspline.Clamped,
		cbspline.Cyclic,
	}

	const nSamples = 1024
	delta := (xMax - xMin) / float64(nSamples-1)

	p := plot.New()
	p.Title.Text = t.Name()
	p.X.Label.Text = "t"
	p.Y.Label.Text = "S(t)"
	p.Add(plotter.NewGrid())
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

	for i, b := range bs {
		sp, err := cbspline.New(ts, vs, b)
		if err != nil {
			t.Fatalf("cbspline.New: %s: %v", b, err)
		}
		curve := make(plotter.XYs, nSamples)
		for j := 0; j < nSamples; j++ {
			curve[j].X = xMin + delta*float64(j)
			curve[j].Y = sp.At(curve[j].X)
		}
		line, _ := plotter.NewLine(curve)
		line.Width = vg.Points(2)
		line.Color = plotutil.Color(i)
		p.Add(line)
		p.Legend.Add(b.String(), line)
	}

	ptxy := make(plotter.XYs, len(ts))
	for i, t := range ts {
		ptxy[i].X, ptxy[i].Y = t, vs[i]
	}
	pts, _ := plotter.NewScatter(ptxy)
	pts.GlyphStyle.Radius = vg.Points(4)
	p.Add(pts)

	outPath := fmt.Sprintf("%s.png", t.Name())
	if dir := os.Getenv("ARTIFACT_DIR"); dir != "" {
		outPath = filepath.Join(dir, outPath)
	}
	_ = p.Save(800, 450, outPath)
}
