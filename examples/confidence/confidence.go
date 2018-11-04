package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/milosgajdos83/gollipse/ellipse"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	// generate random data
	size := 200
	x := make([]float64, size)
	y := make([]float64, size)

	for i := range x {
		x[i] = 5 * rand.NormFloat64()
		y[i] = 5 * rand.NormFloat64()
	}
	data := mat.NewDense(size, 2, nil)
	data.SetCol(0, x)
	data.SetCol(1, y)

	// Calculate mean values for X an Y coordinates
	xMean := stat.Mean(x, nil)
	yMean := stat.Mean(y, nil)

	confidence := 0.95
	ell, err := ellipse.NewDataConfidence(x, y, confidence)
	if err != nil {
		log.Fatalf("Failed to create new ellipse: %v", err)
	}

	line, points, err := ell.LinePoints(size)
	if err != nil {
		log.Fatalf("Failed to compute ellipse points: %v", err)
	}

	// we need to shift ellipse points by data mean values
	for i := range points.XYs {
		points.XYs[i].X = points.XYs[i].X + xMean
		points.XYs[i].Y = points.XYs[i].Y + yMean
	}

	// Create new plot
	p, err := plot.New()
	if err != nil {
		log.Fatalf("Failed to create mew plot: %v", err)
	}
	p.Title.Text = "Ellipse Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// original data points
	dataXY := ellipse.XYFromDense(data)
	scatterData, err := plotter.NewScatter(dataXY)
	if err != nil {
		log.Fatalf("Failed to create data scatter: %v", err)
	}
	scatterData.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	scatterData.GlyphStyle.Radius = vg.Points(1)
	p.Add(scatterData)
	p.Legend.Add("Data", scatterData)

	// ellipse line plot
	line.Color = color.RGBA{B: 255, A: 255}
	p.Add(line)
	p.Legend.Add(fmt.Sprintf("%.2f%%", 100*confidence), line)

	// plot mean values
	var meanXYs plotter.XYs
	meanVals := struct{ X, Y float64 }{xMean, yMean}
	meanXYs = append(meanXYs, meanVals)

	scatterMean, err := plotter.NewScatter(meanXYs)
	if err != nil {
		log.Fatalf("Failed to create mean values scatter: %v", meanXYs)
	}
	scatterMean.GlyphStyle.Color = color.RGBA{G: 255, A: 128}
	scatterMean.GlyphStyle.Shape = draw.PyramidGlyph{}
	scatterMean.GlyphStyle.Radius = vg.Points(3)

	p.Add(scatterMean)
	p.Legend.Add("Mean", scatterMean)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "confidence.png"); err != nil {
		log.Fatalf("Failed to plot confidence ellipse: %v", err)
	}
}
