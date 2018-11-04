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

// Ellipse is a 2D ellipse
// For more information see: https://en.wikipedia.org/wiki/Ellipse
type Ellipse struct {
	a     float64
	b     float64
	angle float64
}

// New creates new Ellipse with given distance of major and minor vertices from the center rotated
// around the orgin by given angle in radians.
// It returns error if either the vertex distances are smaller than or equal to 0.
// Note, major/minor axis lengths are defined as distance of the ellipse vertices on major/minor axis.
// This function however accepts distance of major/minor vertices from the Ellipse center which lies
// in the middle of the distance between particular minor/major vertices.
// For more information please see: https://en.wikipedia.org/wiki/Ellipse#Ellipse_in_Cartesian_coordinates
func New(a, b, angle float64) (*Ellipse, error) {
	if a <= 0 || b <= 0 {
		return nil, fmt.Errorf("Invald ellipse axis: (a: %f, b: %f)", a, b)
	}

	return &Ellipse{a: a, b: b, angle: angle}, nil
}

// NewDataConfidence creates new Ellipse given the list of 2D data points coordinates and confidence.
// The suplied data is assumed to be from the Gaussian statistical distribution
// It panics if either of the folllowing happens:
// * either of the lengths of provided data coordinates is zero
// * the lengths of x and y points coordinates are not equal
// * principal components could not be calculated
// It returns error if the provided confidence value is not in (0,1> interval.
func NewDataConfidence(x, y []float64, confidence float64) (*Ellipse, error) {
	if len(x) == 0 || len(y) == 0 {
		panic("Empty coordinates supplied")
	}

	if len(x) != len(y) {
		panic("Coordinates dimension mismatch")
	}

	if confidence <= 0 || confidence > 1 {
		return nil, fmt.Errorf("Invalid confidence level: %f", confidence)
	}

	// Creates data matrix from provided data point coordinates: data points ares stored in rows
	// The matrix dimensions are: len(x) x 2; X/Y coordinates of each point are stored in matrix columns
	data := mat.NewDense(len(x), 2, nil)
	data.SetCol(0, x)
	data.SetCol(1, y)

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
	chi2 := distuv.ChiSquared{2, src}
	// pc.VarsTo returns eigenvalues in descending order
	a := math.Sqrt(chi2.Quantile(confidence)) * math.Sqrt(eigVals[0])
	b := math.Sqrt(chi2.Quantile(confidence)) * math.Sqrt(eigVals[1])

	return &Ellipse{a: a, b: b, angle: angle}, nil
}

// LinePoints returns both Ellipse line and size number of points which can be used to plot the Ellipse.
// It returns error if one of the ellipse data points contains a NaN or Infinity.
func (e *Ellipse) LinePoints(size int) (*plotter.Line, *plotter.Scatter, error) {
	// generate size points all around ellipse
	points := floats.Span(make([]float64, size), 0, 2*math.Pi)
	x := make([]float64, len(points))
	y := make([]float64, len(points))

	// Parametric representation of ellipse can be obtained as:
	// (a*cos(angle), b*sin(angl)),  where angle is <0, 2*pi>
	for i, point := range points {
		x[i] = e.a * math.Cos(point)
		y[i] = e.b * math.Sin(point)
	}

	// We need to rotate the data around X axis by angle radians
	rotateData := []float64{
		math.Cos(e.angle), math.Sin(e.angle),
		-math.Sin(e.angle), math.Cos(e.angle),
	}
	rotateMx := mat.NewDense(2, 2, rotateData)

	// Perform data rotation
	ellipsePoints := make([]float64, 2*len(points))
	ellipseMx := mat.NewDense(len(points), 2, ellipsePoints)
	ellipseMx.SetCol(0, x)
	ellipseMx.SetCol(1, y)
	ellipseMx.Mul(ellipseMx, rotateMx)

	ellipseXY := XYFromDense(ellipseMx)
	xyLine, xyPoints, err := plotter.NewLinePoints(ellipseXY)
	if err != nil {
		return nil, nil, err
	}

	return xyLine, xyPoints, nil
}

// Eccentricity returns eccentricity of the ellipse
func (e *Ellipse) Eccentricity() float64 {
	return math.Sqrt(1 - (e.a*e.a)/(e.b*e.b))
}

// XYFromDense returns plotter.XYs from matrix m.
// It panics if either m is nil or if it doesn't have at least 2 columns.
func XYFromDense(m *mat.Dense) plotter.XYs {
	r, _ := m.Dims()
	pts := make(plotter.XYs, r)
	for i := 0; i < r; i++ {
		pts[i].X = m.At(i, 0)
		pts[i].Y = m.At(i, 1)
	}

	return pts
}
