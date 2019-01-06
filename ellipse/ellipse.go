package ellipse

import (
	"fmt"
	"math"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/plotter"
)

// Ellipse is 2D ellipse
//
// For more information see: https://en.wikipedia.org/wiki/Ellipse
type Ellipse struct {
	x     float64
	y     float64
	a     float64
	b     float64
	angle float64
}

// New creates new Ellipse with origin [x,y], length of major/minor axis (mx,my) and rotation angle radians.
// It returns error if either of the axis (a or b).
// Note: the lengths of major/minor axis are defined as a distance between the ellipse vertices on major/minor axis.
// however, this function accepts semi-major/minor lengths from the Ellipse origin.
//
// For more information please see: https://en.wikipedia.org/wiki/Semi-major_and_semi-minor_axes
func New(x, y, a, b, angle float64) (*Ellipse, error) {
	if a <= 0 || b <= 0 {
		return nil, fmt.Errorf("Invald ellipse axis: (a: %.2f, b: %.2f)", a, b)
	}

	return &Ellipse{a: a, b: b, angle: angle, x: x, y: y}, nil
}

// NewWithDataConfidence creates new Ellipse from data with origin being data mean and confidence probability.
// The data is assumed to be of the Normal (a.k.a. Gaussian) distribution.
// It panics if either of the folllowing happens:
// * supplied data matrix is nil
// * principal components could not be calculated from the supplied data
// It returns error if confidence is not in (0,1> interval.
func NewWithDataConfidence(data mat.Matrix, confidence float64) (*Ellipse, error) {
	if confidence <= 0 || confidence > 1 {
		return nil, fmt.Errorf("Invalid confidence level: %.2f", confidence)
	}

	// calculate x and y mean values
	rows, _ := data.Dims()
	vals := make([]float64, rows)
	xmean := stat.Mean(mat.Col(vals, 0, data), nil)
	ymean := stat.Mean(mat.Col(vals, 1, data), nil)

	// calculate data eigenvectors and eigenvalues
	var pc stat.PC
	ok := pc.PrincipalComponents(data, nil)
	if !ok {
		panic("Could not determine Principal Components")
	}
	eigVals := pc.VarsTo(nil)
	eigVecs := pc.VectorsTo(nil)

	// Calculate Ellipse rotation angle from the largest eigenvector
	// pc.VectorsTo returns eigenvalues/vectors in descending order
	angle := math.Atan2(eigVecs.At(0, 1), eigVecs.At(0, 0))
	if angle < 0 {
		// Shift the angle to the <0, 2*pi> interval instead of <-pi, pi>
		angle = angle + 2*math.Pi
	}

	// The sum of square Gaussian is distributed according to Chi-squared distribution:
	// https://en.wikipedia.org/wiki/Chi-squared_distribution
	src := rand.New(rand.NewSource(1))
	chi2 := distuv.ChiSquared{K: 2, Src: src}

	// pc.VarsTo returns eigenvalues in descending order
	a := math.Sqrt(chi2.Quantile(confidence) * eigVals[0])
	b := math.Sqrt(chi2.Quantile(confidence) * eigVals[1])

	return &Ellipse{x: xmean, y: ymean, a: a, b: b, angle: angle}, nil
}

// LinePoints returns both plotter.Line and plotter.Scatter which can be used to plot Ellipse.
// It returns error if at least one of the ellipse data points contains a NaN or Infinity.
func (e *Ellipse) LinePoints(size int) (*plotter.Line, *plotter.Scatter, error) {
	// generate size number of ellipse points
	points := floats.Span(make([]float64, size), 0, 2*math.Pi)
	x := make([]float64, len(points))
	y := make([]float64, len(points))

	// Parametric representation of ellipse can be obtained as:
	// (a*cos(angle), b*sin(angl)),  where angle is <0, 2*pi>
	for i, point := range points {
		x[i] = e.a * math.Cos(point)
		y[i] = e.b * math.Sin(point)
	}

	// ellipse data matrix: it will be rotated in the next step
	ellipseMx := mat.NewDense(len(points), 2, nil)
	ellipseMx.SetCol(0, x)
	ellipseMx.SetCol(1, y)

	// We need to rotate the data around X axis by angle radians
	rotateData := []float64{
		math.Cos(e.angle), math.Sin(e.angle),
		-math.Sin(e.angle), math.Cos(e.angle),
	}
	rotateMx := mat.NewDense(2, 2, rotateData)

	// Perform data rotation
	ellipseMx.Mul(ellipseMx, rotateMx)

	// get Ellipse data points
	ellipseXYs := XYFromDense(ellipseMx)

	// we need to shift ellipse points by data mean values
	for i := range ellipseXYs {
		ellipseXYs[i].X = ellipseXYs[i].X + e.x
		ellipseXYs[i].Y = ellipseXYs[i].Y + e.y
	}

	return plotter.NewLinePoints(ellipseXYs)
}

// Eccentricity returns eccentricity of the ellipse
func (e *Ellipse) Eccentricity() float64 {
	return math.Sqrt(1 - (e.a*e.a)/(e.b*e.b))
}

// String implements fmt.Stringer interface
func (e *Ellipse) String() string {
	return fmt.Sprintf("Ellipse{x: %.2f, y: %.2f, a: %.2f, b: %.2f, angle: %.2f}", e.x, e.y, e.a, e.b, e.angle)
}

// XYFromDense returns plotter.XYs from m, which stores X and Y coordinates in its 1st and 2nd column.
// It panics if either m is nil or if m doesn't have at least 2 columns.
func XYFromDense(m *mat.Dense) plotter.XYs {
	r, _ := m.Dims()
	pts := make(plotter.XYs, r)
	for i := 0; i < r; i++ {
		pts[i].X = m.At(i, 0)
		pts[i].Y = m.At(i, 1)
	}

	return pts
}
