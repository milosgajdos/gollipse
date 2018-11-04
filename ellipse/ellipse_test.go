package ellipse

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

func TestNewEllipse(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		a     float64
		b     float64
		angle float64
		err   bool
	}{
		{0, 1.0, 4.5, true},
		{0, 0, 4.5, true},
		{1.0, -10.0, 20.3, true},
		{10.0, 10.0, -5.6, false},
	}

	for _, tc := range testCases {
		ell, err := New(tc.a, tc.b, tc.angle)
		if !tc.err {
			assert.NoError(err)
			assert.Equal(tc.a, ell.a)
			assert.Equal(tc.b, ell.b)
			continue
		}
		assert.Error(err)
		assert.Nil(ell)
	}
}

func TestNewDataConfidence(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		x          []float64
		y          []float64
		confidence float64
		err        bool
		pnc        bool
	}{
		{[]float64{}, []float64{1.0}, 0.05, false, true},
		{[]float64{1.0, 2.0}, []float64{1.0}, 0.05, false, true},
		{[]float64{1.0, 2.0}, []float64{1.0, 2.0}, 0, true, false},
		{[]float64{1.0, 2.0}, []float64{1.0, 2.0}, 2.0, true, false},
		{[]float64{1.0, 2.0}, []float64{1.0, 2.0}, 0.05, false, false},
	}

	for _, tc := range testCases {
		if tc.pnc {
			assert.Panics(func() { NewDataConfidence(tc.x, tc.y, tc.confidence) })
			continue
		}

		ell, err := NewDataConfidence(tc.x, tc.y, tc.confidence)
		if !tc.err {
			assert.NoError(err)
			assert.NotNil(ell)
			continue
		}
		assert.Error(err)
		assert.Nil(ell)
	}
}

func TestLinePoints(t *testing.T) {
	assert := assert.New(t)

	ell := Ellipse{a: 1.0, b: 3.0, angle: math.Pi}
	size := 10

	line, points, err := ell.LinePoints(size)
	assert.NoError(err)
	assert.NotNil(line)
	assert.Equal(size, points.Len())
}

func TestEccentricity(t *testing.T) {
	assert := assert.New(t)

	ell := Ellipse{a: 1.0, b: 3.0, angle: math.Pi}
	ecc := ell.Eccentricity()
	assert.NotZero(ecc)
}

func TestXYFromDense(t *testing.T) {
	assert := assert.New(t)

	testMx1 := mat.NewDense(5, 1, nil)
	testMx2 := mat.NewDense(5, 2, nil)

	testCases := []struct {
		m   *mat.Dense
		pnc bool
	}{
		{nil, true},
		{testMx1, true},
		{testMx2, false},
	}

	for _, tc := range testCases {
		if tc.pnc {
			assert.Panics(func() { XYFromDense(tc.m) })
			continue
		}
		xy := XYFromDense(tc.m)
		assert.NotNil(xy)
	}
}
