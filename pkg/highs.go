package highs

// #cgo pkg-config: highs
// #include <stdlib.h>
// #include <stdio.h>
// #include "interfaces/highs_c_api.h"
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

type ModelStatus int

const (
	ModelNotset ModelStatus = iota
	ModelLoadError
	ModelError
	ModelPresolveError
	ModelSolveError
	ModelPostsolveError
	ModelEmpty
	ModelOptimal
	ModelInfeasible
	ModelUnboundedOrInfeasible
	ModelUnbounded
	ModelObjectiveBound
	ModelObjectiveTarget
	ModelTimeLimit
	ModelIterationLimit
	ModelUnknown
)

func (d ModelStatus) String() string {
	return [...]string{
		"Model Not Set",
		"Model Load Error",
		"Model Error",
		"Model Presolve Error",
		"Model Solve Error",
		"Model Postsolve Error",
		"Model Empty",
		"Model Optimal",
		"Model Infeasible",
		"Model Unbounded or Infeasible",
		"Model Unbounded",
		"Model Objetive Bound",
		"Model Objective Target",
		"Model TimeL imit",
		"Model Iteration Limit",
		"Model Unknown"}[d]
}

type SolutionStatus int

const (
	SolutionNone SolutionStatus = iota
	SolutionInfeasible
	SolutionFeasible
)

type Sense int

const (
	Maximize Sense = -1
	Minimize Sense = 1
)

type Integrality int

const (
	Continious Integrality = iota
	Integer
	ImplicitInteger
)

type Highs struct {
	obj         unsafe.Pointer
	allocs      []unsafe.Pointer
	cols        []float64
	bounds      [][2]float64
	rows        [][]float64
	integrality []int
	ptrs        highsPtrs
}

type highsPtrs struct {
	pCols       unsafe.Pointer
	pColLbs     unsafe.Pointer
	pColUbs     unsafe.Pointer
	pArStart    unsafe.Pointer
	pArIndex    unsafe.Pointer
	pArValue    unsafe.Pointer
	pRowLbs     unsafe.Pointer
	pRowUbs     unsafe.Pointer
	pIntg       unsafe.Pointer
	ArStartSize int
	ArIndexSize int
}

// New returns an allocated Highs object
func New() (*Highs, error) {
	h := &Highs{
		obj:         C.Highs_create(),
		allocs:      []unsafe.Pointer{},
		cols:        []float64{},
		bounds:      [][2]float64{},
		rows:        [][]float64{},
		integrality: []int{},
	}

	runtime.SetFinalizer(h, func(h *Highs) {
		h.destroy()
	})

	// TODO: how to check for a bad malloc?
	return h, nil
}

func (h *Highs) destroy() {
	if h.obj != nil {
		C.Highs_destroy(h.obj)
		h.obj = nil

		for _, p := range h.allocs {
			cFree(p)
		}
		h.allocs = nil
		h.ptrs = highsPtrs{}
	}
}

func (h *Highs) SetColumns(c []float64) {
	h.cols = c
}

func (h *Highs) SetRows(r [][]float64) {
	h.rows = r
}

func (h *Highs) SetBounds(b [][2]float64) {
	h.bounds = b
}

func (h *Highs) SetIntegrality(i []int) {
	h.integrality = i
}

func (h *Highs) GetLowerBounds() []float64 {
	lbs := make([]float64, len(h.cols))
	for i, lb := range h.bounds {
		lbs[i] = lb[0]
	}

	return lbs
}

func (h *Highs) GetUpperBounds() []float64 {
	ubs := make([]float64, len(h.cols))
	for i, ub := range h.bounds {
		ubs[i] = ub[1]
	}

	return ubs
}

func (h *Highs) allocate() {
	h.allocateColumns()
	h.allocateRows()
	h.allocateIntegrality()
}

func (h *Highs) allocateColumns() error {
	n := len(h.cols)
	if n < len(h.bounds) {
		return errors.New("columns are under bounded")
	}

	if n > len(h.bounds) {
		return errors.New("columns are over bounded")
	}

	pCols := cMalloc(n, C.double(0))
	h.allocs = append(h.allocs, pCols)
	cSetArrayDoubles(pCols, h.cols)
	h.ptrs.pCols = pCols

	pLb := cMalloc(n, C.double(0))
	h.allocs = append(h.allocs, pLb)
	cSetArrayDoubles(pLb, h.GetLowerBounds())
	h.ptrs.pColLbs = pLb

	pUb := cMalloc(n, C.double(0))
	h.allocs = append(h.allocs, pUb)
	cSetArrayDoubles(pUb, h.GetUpperBounds())
	h.ptrs.pColUbs = pUb

	return nil
}

func (h *Highs) allocateRows() error {
	if len(h.rows[0])-2 != len(h.cols) {
		return errors.New("row size mismatch len(row[i]) != len(col)")
	}

	rows, lbs, ubs := separateBounds(h.rows)

	pm := packMatrix(rows) // [lb, row constraints..., ub]

	h.ptrs.ArStartSize = len(pm.arStart)
	h.ptrs.ArIndexSize = len(pm.arIndex)

	pArStart := cMalloc(len(pm.arStart), C.int(0))
	h.allocs = append(h.allocs, pArStart)
	cSetArrayInts(pArStart, pm.arStart)
	h.ptrs.pArStart = pArStart

	pArIndex := cMalloc(len(pm.arIndex), C.int(0))
	h.allocs = append(h.allocs, pArIndex)
	cSetArrayInts(pArIndex, pm.arIndex)
	h.ptrs.pArIndex = pArIndex

	pArValue := cMalloc(len(pm.arValue), C.double(0))
	h.allocs = append(h.allocs, pArValue)
	cSetArrayDoubles(pArValue, pm.arValue)
	h.ptrs.pArValue = pArValue

	m := len(h.rows)
	pLb := cMalloc(m, C.double(0))
	h.allocs = append(h.allocs, pLb)
	cSetArrayDoubles(pLb, lbs)
	h.ptrs.pRowLbs = pLb

	pUb := cMalloc(m, C.double(0))
	h.allocs = append(h.allocs, pUb)
	cSetArrayDoubles(pUb, ubs)
	h.ptrs.pRowUbs = pUb

	return nil
}

func (h *Highs) allocateIntegrality() error {
	n := len(h.cols)
	if n != len(h.integrality) {
		return errors.New("integrality len does not match column len")
	}

	pIntg := cMalloc(n, C.int(0))
	h.allocs = append(h.allocs, pIntg)
	cSetArrayInts(pIntg, h.integrality)
	h.ptrs.pIntg = pIntg

	return nil
}

func (h *Highs) RunSolver() (Solution, error) {

	h.allocate()
	if len(h.integrality) > 0 {
		return h.runMipsSolver()
	}
	return h.runLpSolver()
}

func (h *Highs) runMipsSolver() (Solution, error) {
	h.PassMip()
	h.Run()
	status := h.GetModelStatus()
	if status == ModelOptimal {
		return h.GetSolution(), nil
	}

	return Solution{}, fmt.Errorf("solver error: %s", status)
}

func (h *Highs) PassMip() SolutionStatus {

	s := C.Highs_passMip(
		h.obj,
		C.int(len(h.cols)),
		C.int(len(h.rows)),
		C.int(h.ptrs.ArIndexSize),
		C.int(1),
		(*C.double)(h.ptrs.pCols),
		(*C.double)(h.ptrs.pColLbs),
		(*C.double)(h.ptrs.pColUbs),
		(*C.double)(h.ptrs.pRowLbs),
		(*C.double)(h.ptrs.pRowUbs),
		(*C.int)(h.ptrs.pArStart),
		(*C.int)(h.ptrs.pArIndex),
		(*C.double)(h.ptrs.pArValue),
		(*C.int)(h.ptrs.pIntg))

	return (SolutionStatus)(s)
}

func (h *Highs) runLpSolver() (Solution, error) {
	h.PassLp()
	h.Run()
	status := h.GetModelStatus()
	if status == ModelOptimal {
		return h.GetSolution(), nil
	}

	return Solution{}, fmt.Errorf("solver error: %s", status)
}

func (h *Highs) PassLp() SolutionStatus {
	s := C.Highs_passLp(
		h.obj,
		C.int(len(h.cols)),
		C.int(len(h.rows)),
		C.int(h.ptrs.ArIndexSize),
		C.int(1),
		(*C.double)(h.ptrs.pCols),
		(*C.double)(h.ptrs.pColLbs),
		(*C.double)(h.ptrs.pColUbs),
		(*C.double)(h.ptrs.pRowLbs),
		(*C.double)(h.ptrs.pRowUbs),
		(*C.int)(h.ptrs.pArStart),
		(*C.int)(h.ptrs.pArIndex),
		(*C.double)(h.ptrs.pArValue))

	return (SolutionStatus)(s)
}

func (h *Highs) callLpSolver() (Solution, error) {
	return Solution{}, nil
}

func (h *Highs) SetObjectiveSense(s Sense) {
	C.Highs_changeObjectiveSense(h.obj, C.int(s))
}

func (h *Highs) GetObjectiveSense() Sense {
	pS := cMalloc(1, C.int(0))
	h.allocs = append(h.allocs, pS)
	C.Highs_getObjectiveSense(h.obj, (*C.int)(pS))

	return (Sense)(int(*(*C.int)(pS)))
}

func (h *Highs) SetStringOptionValue(opt string, val string) {

	pOpt := C.CString(opt)
	defer cFree(unsafe.Pointer(pOpt))

	pVal := C.CString(val)
	defer cFree(unsafe.Pointer(pVal))

	C.Highs_setStringOptionValue(h.obj, pOpt, pVal)
}

func (h *Highs) GetStringOptionValue(opt string) string {
	pOpt := C.CString(opt)
	defer cFree(unsafe.Pointer(pOpt))

	pVal := cMalloc(1024, C.char('A'))
	defer cFree(pVal)

	C.Highs_getStringOptionValue(h.obj, pOpt, (*C.char)(pVal))

	return C.GoString((*C.char)(pVal))
}

func (h *Highs) Run() SolutionStatus {
	s := C.Highs_run(h.obj)
	return SolutionStatus(s)
}

type Solution struct {
	colValue []float64
	colDual  []float64
	rowValue []float64
	rowDual  []float64
}

func NewSolution() Solution {
	return Solution{[]float64{}, []float64{}, []float64{}, []float64{}}
}

func (h *Highs) GetSolution() Solution {
	// this is dangerous, need to be sure the column size isn't changed after
	// allocation and pass to HiGHs
	n := len(h.cols)
	m := len(h.rows)

	pColValue := cMalloc(n, C.double(0))
	h.allocs = append(h.allocs, pColValue)

	pColDual := cMalloc(n, C.double(0))
	h.allocs = append(h.allocs, pColDual)

	pRowValue := cMalloc(m, C.double(0))
	h.allocs = append(h.allocs, pRowValue)

	pRowDual := cMalloc(m, C.double(0))
	h.allocs = append(h.allocs, pRowDual)

	C.Highs_getSolution(h.obj, (*C.double)(pColValue), (*C.double)(pColDual), (*C.double)(pRowValue), (*C.double)(pRowDual))

	s := NewSolution()
	s.colValue = copyDoubles(pColValue, n)
	cFree(pColValue)

	s.colDual = copyDoubles(pColDual, n)
	cFree(pColDual)

	s.rowValue = copyDoubles(pRowValue, m)
	cFree(pRowValue)

	s.rowDual = copyDoubles(pRowDual, m)
	cFree(pRowDual)

	return s
}

func (h *Highs) GetModelStatus() ModelStatus {
	s := C.Highs_getModelStatus(h.obj)
	return ModelStatus(s)
}

type PackedMatrix struct {
	arStart []int
	arIndex []int
	arValue []float64
}

// packRows returns the rows in a flat packed form use in the HiGHS c interface
func packMatrix(matrix [][]float64) PackedMatrix {
	arStart := []int{}
	arIndex := []int{}
	arValue := []float64{}

	idx := 0
	for _, cols := range matrix {
		arStart = append(arStart, idx)
		for i, v := range cols {
			if v != 0 {
				arIndex = append(arIndex, i)
				arValue = append(arValue, v)
			}
		}
		idx = len(arIndex)
	}

	return PackedMatrix{arStart, arIndex, arValue}
}

// separateBounds splits a bounded rows [lb, row..., ub] into its components lb, rows, ub
func separateBounds(bound_rows [][]float64) ([][]float64, []float64, []float64) {
	col_size := len(bound_rows[0])

	rows := [][]float64{}
	lbs := []float64{}
	ubs := []float64{}
	for _, row := range bound_rows {
		rows = append(rows, row[1:col_size-1])
		lbs = append(lbs, row[0])
		ubs = append(ubs, row[col_size-1])
	}

	return rows, lbs, ubs
}

// cMalloc asks C to allocate memory.  For convenience to Go, the arguments
// are like calloc's except that the size argument is a value, which cMalloc
// will take the size of.  cMalloc panics on error (typically, out of memory).
func cMalloc(nmemb int, sizeVal interface{}) unsafe.Pointer {
	size := reflect.TypeOf(sizeVal).Size()
	mem := C.malloc(C.size_t(uintptr(nmemb) * size))
	if mem == nil {
		panic("HiGHS: malloc failed")
	}
	return mem
}

// cFree asks C to free memory.
func cFree(mem unsafe.Pointer) {
	C.free(mem)
}

func cSetArrayInts(a unsafe.Pointer, vs []int) {
	for i, v := range vs {
		cSetArrayInt(a, i, v)
	}
}

// cSetArrayInt assigns a[i] = v where a is a C.int array allocated by
// cMalloc and i and v are Go ints.
func cSetArrayInt(a unsafe.Pointer, i, v int) {
	eSize := unsafe.Sizeof(C.int(0))
	ptr := unsafe.Pointer(uintptr(a) + uintptr(i)*eSize)
	*(*C.int)(ptr) = C.int(v)
}

func copyInts(a unsafe.Pointer, size int) []int {
	gs := make([]int, size)
	for i := range gs {
		gs[i] = cGetArrayInt(a, i)
	}

	return gs
}

func cGetArrayInt(a unsafe.Pointer, i int) int {
	eSize := unsafe.Sizeof(C.int(0))
	ptr := unsafe.Pointer(uintptr(a) + uintptr(i)*eSize)
	return int(*(*C.int)(ptr))
}

func cSetArrayDoubles(a unsafe.Pointer, vs []float64) {
	for i, v := range vs {
		cSetArrayDouble(a, i, v)
	}
}

// cSetArrayDouble assigns a[i] = v where a is a C.double array allocated by
// cMalloc, i is an int, and v is a Go float64.
func cSetArrayDouble(a unsafe.Pointer, i int, v float64) {
	eSize := unsafe.Sizeof(C.double(0.0))
	ptr := unsafe.Pointer(uintptr(a) + uintptr(i)*eSize)
	*(*C.double)(ptr) = C.double(v)
}

func copyDoubles(a unsafe.Pointer, size int) []float64 {
	gs := make([]float64, size)
	for i := range gs {
		gs[i] = cGetArrayDouble(a, i)
	}

	return gs
}

// cGetArrayDouble returns a[i] as a Go float64 where a is a C.double array
// allocated by cMalloc and i is an int.
func cGetArrayDouble(a unsafe.Pointer, i int) float64 {
	eSize := unsafe.Sizeof(C.double(0.0))
	ptr := unsafe.Pointer(uintptr(a) + uintptr(i)*eSize)
	return float64(*(*C.double)(ptr))
}
