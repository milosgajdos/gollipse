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
	ell, err := ellipse.New(50, 10, math.Pi/2)
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

	// generate ellipse curve: we request 100 points
	line, _, err := ell.LinePoints(100)
	if err != nil {
		log.Fatalf("Failed to compute ellipse curve points: %v", err)
	}

	// ellipse line plot
	line.Color = color.RGBA{B: 255, A: 255}
	p.Add(line)
	p.Legend.Add("a=10,b=20,angle=pi/2", line)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "example.png"); err != nil {
		log.Fatalf("Failed to plot ellipse: %v", err)
	}
}
