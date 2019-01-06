package main

import (
	"image/color"
	"log"
	"math"

	"github.com/milosgajdos83/gollipse/ellipse"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

func main() {
	ell, err := ellipse.New(10, 10, 50, 10, math.Pi/2)
	if err != nil {
		log.Fatalf("Failed to create ellipse: %v", err)
	}

	// Create new plot
	p, err := plot.New()
	if err != nil {
		log.Fatalf("Failed to create mew plot: %v", err)
	}
	p.Title.Text = "Ellipse Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	legend, err := plot.NewLegend()
	if err != nil {
		log.Fatalf("Failed to create new plot legend")
	}
	legend.Top = true

	p.Legend = legend

	// generate ellipse curve: we request 100 points
	line, _, err := ell.LinePoints(100)
	if err != nil {
		log.Fatalf("Failed to compute ellipse curve points: %v", err)
	}

	// ellipse line plot
	line.Color = color.RGBA{B: 255, A: 255}
	p.Add(line)
	p.Legend.Add("a=10\nb=20\nangle=pi/2", line)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "simple.png"); err != nil {
		log.Fatalf("Failed to plot ellipse: %v", err)
	}
}
