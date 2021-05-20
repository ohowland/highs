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

type Sense int

const (
	Minimize Sense = -1
	Maximize Sense = 1
)

type Integrality int

const (
	Continious Integrality = iota
	Discrete
)

type Dims struct {
	rows int
	cols int
}

type Highs struct {
	obj    unsafe.Pointer
	allocs []unsafe.Pointer
	dims   Dims
}

// New returns an allocated Highs object
func New() (*Highs, error) {
	h := &Highs{
		obj:    C.Highs_create(),
		allocs: make([]unsafe.Pointer, 0, 64),
		dims:   Dims{0, 0},
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
	}
}

func (h *Highs) AddColumns(cols []float64, lb []float64, ub []float64) error {
	h.dims.cols = len(cols)
	if h.dims.cols != len(lb) || h.dims.cols != len(ub) {
		return errors.New("all slice parameters must be equal length.")
	}

	pCols := cMalloc(h.dims.cols, C.double(0))
	h.allocs = append(h.allocs, pCols)
	cSetArrayDoubles(pCols, cols)

	pLb := cMalloc(h.dims.cols, C.double(0))
	h.allocs = append(h.allocs, pLb)
	cSetArrayDoubles(pLb, lb)

	pUb := cMalloc(h.dims.cols, C.double(0))
	h.allocs = append(h.allocs, pUb)
	cSetArrayDoubles(pUb, ub)

	err := C.Highs_addCols(
		h.obj,
		C.int(h.dims.cols),
		(*C.double)(pCols),
		(*C.double)(pLb),
		(*C.double)(pUb),
		C.int(0),
		nil,
		nil,
		nil)

	if err == 0 {
		return errors.New(fmt.Sprintf("unable to add columns; returned error: %d", err))
	}

	return nil
}

func (h *Highs) AddRows(rows [][]float64, lb []float64, ub []float64) error {
	h.dims.rows = len(rows)
	if h.dims.rows != len(lb) || h.dims.rows != len(ub) {
		return errors.New("an upper and lower bound must be specified for all rows, len(row[0]) != len(ub) or len(lb)")
	}

	pm := packMatrix(rows)

	pArStart := cMalloc(len(pm.arStart), C.int(0))
	h.allocs = append(h.allocs, pArStart)
	cSetArrayInts(pArStart, pm.arStart)

	pArIndex := cMalloc(len(pm.arIndex), C.int(0))
	h.allocs = append(h.allocs, pArIndex)
	cSetArrayInts(pArIndex, pm.arIndex)

	pArValue := cMalloc(len(pm.arValue), C.double(0))
	h.allocs = append(h.allocs, pArValue)
	cSetArrayDoubles(pArValue, pm.arValue)

	pLb := cMalloc(h.dims.rows, C.double(0))
	h.allocs = append(h.allocs, pLb)
	cSetArrayDoubles(pLb, lb)

	pUb := cMalloc(h.dims.rows, C.double(0))
	h.allocs = append(h.allocs, pUb)
	cSetArrayDoubles(pUb, ub)

	err := C.Highs_addRows(
		h.obj,
		C.int(h.dims.rows),
		(*C.double)(pLb),
		(*C.double)(pUb),
		C.int(len(pm.arIndex)),
		(*C.int)(pArStart),
		(*C.int)(pArIndex),
		(*C.double)(pArValue))

	if err == 0 {
		return fmt.Errorf("unable to add rows; returned error: %d", err)
	}

	return nil
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

func (h *Highs) SetIntegrality(col int, i Integrality) {
	_ = C.Highs_changeColIntegrality(h.obj, C.int(col), C.int(i))
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

func (h *Highs) Run() error {
	C.Highs_run(h.obj)
	return nil
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
	pColValue := cMalloc(h.dims.cols, C.double(0))
	h.allocs = append(h.allocs, pColValue)

	pColDual := cMalloc(h.dims.cols, C.double(0))
	h.allocs = append(h.allocs, pColDual)

	pRowValue := cMalloc(h.dims.rows, C.double(0))
	h.allocs = append(h.allocs, pRowValue)

	pRowDual := cMalloc(h.dims.rows, C.double(0))
	h.allocs = append(h.allocs, pRowDual)

	C.Highs_getSolution(h.obj, (*C.double)(pColValue), (*C.double)(pColDual), (*C.double)(pRowValue), (*C.double)(pRowDual))

	s := NewSolution()
	s.colValue = copyDoubles(pColValue, h.dims.cols)
	cFree(pColValue)

	s.colDual = copyDoubles(pColDual, h.dims.cols)
	cFree(pColDual)

	s.rowValue = copyDoubles(pRowValue, h.dims.rows)
	cFree(pRowValue)

	s.rowDual = copyDoubles(pRowDual, h.dims.rows)
	cFree(pRowDual)

	return s
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

/*
// A Nonzero represents an element in a sparse row or column.
type Nonzero struct {
	Index int     // Zero-based element offset
	Value float64 // Value at that offset
}

// A Matrix sparsely represents a set of linear expressions.  Each column
// represents a variable, each row represents an expression, and each cell
// containing a coefficient. Bounds on rows and columns are applied during
// model initialization.
type Matrix interface {
	AppendColumn(col []Nonzero) // Append a column given values for all of its nonzero elements
	Dims() (rows, cols int)     // Return the matrix's dimensions
}


// c__GetArrayInt returns a[i] as a Go int where a is a C.int array allocated
// by cMalloc and i is a Go ints.
func cGetArrayInt(a unsafe.Pointer, i int) int {
	eSize := unsafe.Sizeof(C.int(0))
	ptr := unsafe.Pointer(uintptr(a) + uintptr(i)*eSize)
	return int(*(*C.int)(ptr))
}

// copyIntsGoC copies a slice of Go ints to a slice of C ints.
func copyIntsGoC(cs []C.int, gs []int) {
	if len(gs) != len(cs) {
		panic("Slices of different sizes were passed to copyIntsGoC")
	}
	for i, g := range gs {
		cs[i] = C.int(g)
	}
}

// copyIntsCGo copies a slice of C ints to a slice of Go ints.
func copyIntsCGo(gs []int, cs []C.int) {
	if len(gs) != len(cs) {
		panic("Slices of different sizes were passed to copyIntsCGo")
	}
	for i, c := range cs {
		gs[i] = int(c)
	}
}

*/
