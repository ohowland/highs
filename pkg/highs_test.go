package highs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BuildExampleMipHighs(t *testing.T) *Highs {
	cols := []float64{2.0, 3.0}
	bnds := [][2]float64{{0.0, 3.0}, {1.0, 1e30}}
	rows := [][]float64{
		{-1e30, 0.0, 1.0, 6.0},
		{10.0, 1.0, 2.0, 14.0},
		{8.0, 2.0, 1.0, 1e30}}
	intg := []int{1, 1}

	h, err := New(cols, bnds, rows, intg)
	assert.NoError(t, err)

	return h
}

func BuildExampleHighs(t *testing.T) *Highs {
	h, err := New()
	assert.NoError(t, err)

	cols := []float64{2.0, 3.0}
	h.SetColumns(cols)
	bnds := [][2]float64{{0.0, 3.0}, {1.0, 1e30}}
	h.SetBounds(bnds)

	rows := [][]float64{
		{-1e30, 0.0, 1.0, 6.0},
		{10.0, 1.0, 2.0, 14.0},
		{8.0, 2.0, 1.0, 1e30}}

	h.SetRows(rows)
	assert.NoError(t, err)

	return h
}

func TestCreateHighs(t *testing.T) {
	_, err := New()
	assert.Nil(t, err, "Error returned when allocating new highs object.")
}

func TestSetCols(t *testing.T) {
	h, _ := New()
	cols := []float64{2.0, 3.0}
	h.SetColumns(cols)
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
	h, _ := New()
	rows := [][]float64{
		{-10e30, 0.0, 1.0, 6.0},
		{10.0, 1.0, 2.0, 14.0},
		{8.0, 2.0, 1.0, 1.0e30}}
	h.SetRows(rows)

	assert.Equal(t, rows, h.rows)
}

func TestAllocateColumns(t *testing.T) {
	h, _ := New()

	cols := []float64{2.0, 3.0}
	h.SetColumns(cols)
	bnds := [][2]float64{{0.0, 3.0}, {1.0, 1e30}}
	h.SetBounds(bnds)

	h.allocateColumns()
	assert.Equal(t, cols, (copyDoubles(h.ptrs.pCols, len(h.cols))))
	assert.Equal(t, h.GetLowerBounds(), (copyDoubles(h.ptrs.pColLbs, len(h.cols))))
	assert.Equal(t, h.GetUpperBounds(), (copyDoubles(h.ptrs.pColUbs, len(h.cols))))
}

func TestAllocateRows(t *testing.T) {
	h, _ := New()

	cols := []float64{2.0, 3.0}
	h.SetColumns(cols)
	bnds := [][2]float64{{0.0, 3.0}, {1.0, 1e30}}
	h.SetBounds(bnds)
	h.allocateColumns()

	bounded_rows := [][]float64{
		{-1e30, 0.0, 1.0, 6.0},
		{10.0, 1.0, 2.0, 14.0},
		{8.0, 2.0, 1.0, 1e30}}

	h.SetRows(bounded_rows)
	err := h.allocateRows()
	assert.NoError(t, err)

	r, l, u := separateBounds(bounded_rows)
	pm := packMatrix(r)

	assert.Equal(t, pm.arStart, (copyInts(h.ptrs.pArStart, h.ptrs.ArStartSize)), "arStart malformed")
	assert.Equal(t, pm.arIndex, (copyInts(h.ptrs.pArIndex, h.ptrs.ArIndexSize)), "arIndex malformed")
	assert.Equal(t, pm.arValue, (copyDoubles(h.ptrs.pArValue, h.ptrs.ArIndexSize)), "arValue malformed")
	assert.Equal(t, l, (copyDoubles(h.ptrs.pRowLbs, len(h.rows))))
	assert.Equal(t, u, (copyDoubles(h.ptrs.pRowUbs, len(h.rows))))
}

func TestRunLpSolver(t *testing.T) {
	h := BuildExampleHighs(t)

	s, err := h.RunSolver()
	fmt.Println(s.colValue)
	assert.NoError(t, err)
}

func TestRunMipSolver(t *testing.T) {
	h := BuildExampleMipHighs(t)

	s, err := h.RunSolver()
	fmt.Println(s.colValue)
	assert.NoError(t, err)
}

// SET/GET option

func TestSetBoolOptionValue(t *testing.T) {
	h, _ := New()
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
	h, _ := New()
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
