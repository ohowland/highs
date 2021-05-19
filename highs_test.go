package highs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateHighs(t *testing.T) {
	h, err := New()
	defer h.Free()

	assert.Nil(t, err, "Error returned when allocating new highs object.")
}

func TestAddCols(t *testing.T) {
	h, _ := New()
	defer h.Free()

	cost := []float64{2.0, 3.0}
	lb := []float64{0.0, 1.0}
	ub := []float64{3.0, 1e30}

	err := h.AddColumns(cost, lb, ub)

	assert.Nil(t, err, "Error returned when adding columns to highs object")
}

func TestPackRows(t *testing.T) {
	rows := [][]float64{{0.0, 1.0}, {1.0, 2.0}, {2.0, 1.0}}
	arStart, arIndex, arValue := packRows(rows)

	assert.Equal(t, []int{0, 1, 3}, arStart, "malformed start index slice")
	assert.Equal(t, []int{1, 0, 1, 0, 1}, arIndex, "malformed decision variable index slice")
	assert.Equal(t, []float64{1.0, 1.0, 2.0, 2.0, 1.0}, arValue, "malformed decision variable coefficient slice")
}

func TestAddRowsWithoutCols(t *testing.T) {
	h, _ := New()
	defer h.Free()

	rows := [][]float64{{0.0, 1.0}, {1.0, 2.0}, {2.0, 1.0}}
	lb := []float64{-10e30, 10.0, 8}
	ub := []float64{6.0, 14.0, 1.0e30}
	err := h.AddRows(rows, lb, ub)

	assert.Error(t, err, "No error returned when adding rows to a highs object that contains now columns")
}

func TestAddColsAndRows(t *testing.T) {
	h, _ := New()
	defer h.Free()

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
	defer h.Free()

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
}
