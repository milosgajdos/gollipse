package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/milosgajdos/gollipse/ellipse"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	// Create new plot
	p, err := plot.New()
	if err != nil {
		log.Fatalf("Failed to create mew plot: %v", err)
	}
	p.Title.Text = "Ellipse Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// generate random data
	size := 200
	data := mat.NewDense(size, 2, nil)
	for i := 0; i < size; i++ {
		data.Set(i, 0, 5*rand.NormFloat64())
		data.Set(i, 1, 5*rand.NormFloat64())
	}

	// plot data points
	dataXY := ellipse.XYFromDense(data)
	scatterData, err := plotter.NewScatter(dataXY)
	if err != nil {
		log.Fatalf("Failed to create data scatter: %v", err)
	}
	scatterData.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	scatterData.GlyphStyle.Radius = vg.Points(1)

	p.Add(scatterData)
	p.Legend.Add("Data", scatterData)

	// Calculate mean values for X an Y coordinates
	vals := make([]float64, size)

	xMean := stat.Mean(mat.Col(vals, 0, data), nil)
	yMean := stat.Mean(mat.Col(vals, 1, data), nil)

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

	// create new ellipse with 95% data confidence
	confidence := 0.95
	ell, err := ellipse.NewWithDataConfidence(data, confidence)
	if err != nil {
		log.Fatalf("Failed to create new ellipse: %v", err)
	}

	line, _, err := ell.LinePoints(size)
	if err != nil {
		log.Fatalf("Failed to compute ellipse points: %v", err)
	}

	// ellipse line plot
	line.Color = color.RGBA{B: 255, A: 255}
	p.Add(line)
	p.Legend.Add(fmt.Sprintf("%.2f%%", 100*confidence), line)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "confidence.png"); err != nil {
		log.Fatalf("Failed to plot confidence ellipse: %v", err)
	}
}
