package highs

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BuildExampleMipHighs(t *testing.T) *Highs {
	cols := []float64{2.0, 3.0}
	bnds := [][2]float64{{0.0, 3.0}, {1.0, math.Inf(1)}}
	rows := [][]float64{
		{math.Inf(-1), 0.0, 1.0, 6.0},
		{10.0, 1.0, 2.0, 14.0},
		{8.0, 2.0, 1.0, math.Inf(1)}}
	intg := []int{1, 1}

	h, err := New(cols, bnds, rows, intg)
	assert.NoError(t, err)

	return h
}

func BuildExampleHighs(t *testing.T) *Highs {
	cols := []float64{2.0, 3.0}
	bnds := [][2]float64{{0.0, 3.0}, {1.0, math.Inf(1)}}

	rows := [][]float64{
		{math.Inf(-1), 0.0, 1.0, 6.0},
		{10.0, 1.0, 2.0, 14.0},
		{8.0, 2.0, 1.0, math.Inf(1)}}

	h, err := New(cols, bnds, rows, []int{})
	assert.NoError(t, err)

	return h
}

func TestCreateHighs(t *testing.T) {
	_, err := New([]float64{}, [][2]float64{}, [][]float64{}, []int{})
	assert.Nil(t, err, "Error returned when allocating new highs object.")
}

func TestSetCols(t *testing.T) {
	h := BuildExampleHighs(t)
	cols := []float64{2.0, 3.0}
	assert.Equal(t, cols, h.cols, "Error returned when adding columns to highs object")
}

func TestSeparateBounds(t *testing.T) {
	rows := [][]float64{
		{-1, 0.0, 1.0, 2.0},
		{10.0, 1.0, 2.0, 14.0},
		{8.0, 2.0, 1.0, 20}}

	rows, lbs, ubs := separateBounds(rows)

	assert.Equal(t, []float64{-1, 10, 8}, lbs)
	assert.Equal(t, []float64{2, 14, 20}, ubs)
	assert.Equal(t, [][]float64{{0, 1}, {1, 2}, {2, 1}}, rows)
}

func TestPackRows(t *testing.T) {
	matrix := [][]float64{{0.0, 1.0}, {1.0, 2.0}, {2.0, 1.0}}
	pm := packMatrix(matrix)

	assert.Equal(t, []int{0, 1, 3}, pm.arStart, "malformed start index slice")
	assert.Equal(t, []int{1, 0, 1, 0, 1}, pm.arIndex, "malformed decision variable index slice")
	assert.Equal(t, []float64{1.0, 1.0, 2.0, 2.0, 1.0}, pm.arValue, "malformed decision variable coefficient slice")
}

func TestSetRows(t *testing.T) {
	h := BuildExampleHighs(t)
	rows := [][]float64{
		{math.Inf(-1), 0.0, 1.0, 6.0},
		{10.0, 1.0, 2.0, 14.0},
		{8.0, 2.0, 1.0, math.Inf(1)}}

	assert.Equal(t, rows, h.rows)
}

func TestAllocateColumns(t *testing.T) {
	h := BuildExampleHighs(t)

	cols := []float64{2.0, 3.0}

	h.allocateColumns()
	assert.Equal(t, cols, (copyDoubles(h.allocs[pCols], len(h.cols))))
	assert.Equal(t, h.GetLowerBounds(), (copyDoubles(h.allocs[pColLbs], len(h.cols))))
	assert.Equal(t, h.GetUpperBounds(), (copyDoubles(h.allocs[pColUbs], len(h.cols))))
}

func TestAllocateRows(t *testing.T) {
	h := BuildExampleHighs(t)

	h.allocateColumns()

	bounded_rows := [][]float64{
		{math.Inf(-1), 0.0, 1.0, 6.0},
		{10.0, 1.0, 2.0, 14.0},
		{8.0, 2.0, 1.0, math.Inf(1)}}

	err := h.allocateRows()
	assert.NoError(t, err)

	r, l, u := separateBounds(bounded_rows)
	pm := packMatrix(r)

	assert.Equal(t, pm.arStart, (copyInts(h.allocs[pArStart], h.dims.ArStartSize)), "arStart malformed")
	assert.Equal(t, pm.arIndex, (copyInts(h.allocs[pArIndex], h.dims.ArIndexSize)), "arIndex malformed")
	assert.Equal(t, pm.arValue, (copyDoubles(h.allocs[pArValue], h.dims.ArIndexSize)), "arValue malformed")
	assert.Equal(t, l, (copyDoubles(h.allocs[pRowLbs], len(h.rows))))
	assert.Equal(t, u, (copyDoubles(h.allocs[pRowUbs], len(h.rows))))
}

func TestRunLpSolver(t *testing.T) {
	h := BuildExampleHighs(t)

	_, err := h.RunSolver()
	assert.NoError(t, err)
}

func TestRunMipSolver(t *testing.T) {
	h := BuildExampleMipHighs(t)

	_, err := h.RunSolver()
	assert.NoError(t, err)
}

// SET/GET option

func TestSetBoolOptionValue(t *testing.T) {
	h := BuildExampleHighs(t)
	h.SetBoolOptionValue("output_flag", true)

	r := h.GetBoolOptionValue("output_flag")

	assert.Equal(t, true, r)
}

func TestSetIntOptionValue(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestSetDoubleOptionValue(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestSetStringOptionValue(t *testing.T) {
	h := BuildExampleHighs(t)
	opt := "solver"
	val := "ipm"
	h.SetStringOptionValue(opt, val)

	r := h.GetStringOptionValue(opt)

	assert.Equal(t, val, r, "value written to option was not returned")

}

// Objective Sense

func TestChangeObjectiveSense(t *testing.T) {
	h := BuildExampleHighs(t)

	h.SetObjectiveSense(Minimize)
	s := h.GetObjectiveSense()
	assert.Equal(t, Minimize, s, "sense is not Minimize")

	h.SetObjectiveSense(Maximize)
	s = h.GetObjectiveSense()
	assert.Equal(t, Maximize, s, "sense is not Maximize")

}

// Integrality

func TestChangeColIntegrality(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestChangeColIntegralityByMask(t *testing.T) {
	assert.Fail(t, "unimplemented")
}

func TestGetModelStatus(t *testing.T) {
	assert.Fail(t, "unimplemented")
}
