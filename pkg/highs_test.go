package highs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Sense int

const (
	Minimize Sense = iota
	Maximize
)

type BooleanOption int

const 

func TestCreateHighs(t *testing.T) {
	h, err := New()
	defer h.Destroy()

	assert.Nil(t, err, "Error returned when allocating new highs object.")
}

func TestAddCols(t *testing.T) {
	h, _ := New()
	defer h.Destroy()

	cost := []float64{2.0, 3.0}
	lb := []float64{0.0, 1.0}
	ub := []float64{3.0, 1e30}

	err := h.AddColumns(cost, lb, ub)

	assert.Nil(t, err, "Error returned when adding columns to highs object")
}

func TestPackRows(t *testing.T) {
	matrix := [][]float64{{0.0, 1.0}, {1.0, 2.0}, {2.0, 1.0}}
	pm := packMatrix(matrix)

	assert.Equal(t, []int{0, 1, 3}, pm.arStart, "malformed start index slice")
	assert.Equal(t, []int{1, 0, 1, 0, 1}, pm.arIndex, "malformed decision variable index slice")
	assert.Equal(t, []float64{1.0, 1.0, 2.0, 2.0, 1.0}, pm.arValue, "malformed decision variable coefficient slice")
}

func TestAddRowsWithoutCols(t *testing.T) {
	h, _ := New()
	defer h.Destroy()

	rows := [][]float64{{0.0, 1.0}, {1.0, 2.0}, {2.0, 1.0}}
	lb := []float64{-10e30, 10.0, 8}
	ub := []float64{6.0, 14.0, 1.0e30}
	err := h.AddRows(rows, lb, ub)

	assert.Error(t, err, "No error returned when adding rows to a highs object that contains now columns")
}

func TestAddColsAndRows(t *testing.T) {
	h, _ := New()
	defer h.Destroy()

	cost := []float64{2.0, 3.0}
	lb := []float64{0.0, 1.0}
	ub := []float64{3.0, 1e30}

	err := h.AddColumns(cost, lb, ub)
	assert.Nil(t, err, "Error returned when adding columns to highs object")

	rows := [][]float64{{0.0, 1.0}, {1.0, 2.0}, {2.0, 1.0}}
	lb = []float64{-1e30, 10.0, 8}
	ub = []float64{6.0, 14.0, 1e30}

	err = h.AddRows(rows, lb, ub)
	assert.Nil(t, err, "Error returned when adding columns to highs object")
}

func TestRun(t *testing.T) {
	h, _ := New()
	defer h.Destroy()

	cost := []float64{2.0, 3.0}
	lb := []float64{0.0, 1.0}
	ub := []float64{3.0, 1e30}

	err := h.AddColumns(cost, lb, ub)
	assert.Nil(t, err, "Error returned when adding columns to highs object")

	rows := [][]float64{{0.0, 1.0}, {1.0, 2.0}, {2.0, 1.0}}
	lb = []float64{-1e30, 10.0, 8}
	ub = []float64{6.0, 14.0, 1e30}

	err = h.AddRows(rows, lb, ub)
	assert.Nil(t, err, "Error returned when adding columns to highs object")

	h.Run()
	s := h.GetSolution()

	assert.Equal(t, []float64{2, 4}, s.colValue, "unexpected primal solution")
}

// SET option

func TestSetBoolOptionValue(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestSetIntOptionValue(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestSetDoubleOptionValue(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestSetStringOptionValue(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestSetOptionValue(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

// GET option

func TestGetBoolOptionValue(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestGetIntOptionValue(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestGetDoubleOptionValue(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestGetStringOptionValue(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestGetOptionType(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

// Objective Sense

func TestChangeObjectiveSense(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestGetObjectiveSense(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

// Integrality

func TestChangeColIntegrality(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestChangeColIntegralityByMask(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestGetIntInfoValue(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestGetDoubleInfoValue(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestGetBasis(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestGetModelStatus(t *testing.T) {
	assert.Fail(t, "unimplemented")
}
