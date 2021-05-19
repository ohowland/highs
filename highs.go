package highs

// #cgo pkg-config: highs
// #include <stdlib.h>
// #include "highs-interface.h"
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

type Highs struct {
	matrix *C.highs_obj
	allocs []unsafe.Pointer
}

// NewHighsObj returns an allocated Highs object
func New() (*Highs, error) {
	highs := &Highs{
		matrix: C.highsiface_create(),
		allocs: make([]unsafe.Pointer, 0, 64),
	}

	// TODO: how to check for a bad malloc?
	return highs, nil
}

func (h *Highs) Free() {
	if h.matrix != nil {
		C.highsiface_free(h.matrix)
		h.matrix = nil

		for _, p := range h.allocs {
			cFree(p)
		}
		h.allocs = nil
	}
}

func (h *Highs) AddColumns(cols []float64, lb []float64, ub []float64) error {
	n := len(cols)
	if n != len(lb) || n != len(ub) {
		return errors.New("all slice parameters must be equal length.")
	}

	pCols := cMalloc(n, C.double(0))
	h.allocs = append(h.allocs, pCols)
	cSetArrayDoubles(pCols, cols)

	pLb := cMalloc(n, C.double(0))
	h.allocs = append(h.allocs, pLb)
	cSetArrayDoubles(pLb, lb)

	pUb := cMalloc(n, C.double(0))
	h.allocs = append(h.allocs, pUb)
	cSetArrayDoubles(pUb, ub)

	err := C.highsiface_add_cols(
		h.matrix,
		C.int(n),
		(*C.double)(pCols),
		(*C.double)(pLb),
		(*C.double)(pUb))

	if err == 0 {
		return errors.New(fmt.Sprintf("unable to add columns; returned error: %d", err))
	}

	return nil
}

func (h *Highs) AddRows(rows [][]float64, lb []float64, ub []float64) error {
	n := len(rows)
	if n != len(lb) || n != len(ub) {
		return errors.New("an upper and lower bound must be specified for all rows, len(row[0]) != len(ub) or len(lb)")
	}

	arStart, arIndex, arValue := packRows(rows)

	pArStart := cMalloc(len(arStart), C.int(0))
	h.allocs = append(h.allocs, pArStart)
	cSetArrayInts(pArStart, arStart)

	pArIndex := cMalloc(len(arIndex), C.int(0))
	h.allocs = append(h.allocs, pArIndex)
	cSetArrayInts(pArIndex, arIndex)

	pArValue := cMalloc(len(arValue), C.double(0))
	h.allocs = append(h.allocs, pArValue)
	cSetArrayDoubles(pArValue, arValue)

	pLb := cMalloc(n, C.double(0))
	h.allocs = append(h.allocs, pLb)
	cSetArrayDoubles(pLb, lb)

	pUb := cMalloc(n, C.double(0))
	h.allocs = append(h.allocs, pUb)
	cSetArrayDoubles(pUb, ub)

	err := C.highsiface_add_rows(
		h.matrix,
		C.int(n),
		(*C.double)(pLb),
		(*C.double)(pUb),
		C.int(len(arIndex)),
		(*C.int)(pArStart),
		(*C.int)(pArIndex),
		(*C.double)(pArValue))

	if err == 0 {
		return errors.New(fmt.Sprintf("unable to add rows; returned error: %d", err))
	}

	return nil
}

func (h *Highs) Run() error {
	C.highsiface_run(h.matrix)

	return nil
}

// packRows returns the rows in a flat packed form use in the HiGHS c interface
func packRows(rows [][]float64) ([]int, []int, []float64) {
	arStart := []int{}
	arIndex := []int{}
	arValue := []float64{}

	idx := 0
	for _, cols := range rows {
		arStart = append(arStart, idx)
		for i, v := range cols {
			if v != 0 {
				arIndex = append(arIndex, i)
				arValue = append(arValue, v)
			}
		}
		idx = len(arIndex)
	}

	return arStart, arIndex, arValue
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

*/

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

/*

// c__GetArrayInt returns a[i] as a Go int where a is a C.int array allocated
// by cMalloc and i is a Go ints.
func cGetArrayInt(a unsafe.Pointer, i int) int {
	eSize := unsafe.Sizeof(C.int(0))
	ptr := unsafe.Pointer(uintptr(a) + uintptr(i)*eSize)
	return int(*(*C.int)(ptr))
}

*/

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

/*

// cGetArrayDouble returns a[i] as a Go float64 where a is a C.double array
// allocated by cMalloc and i is an int.
func cGetArrayDouble(a unsafe.Pointer, i int) float64 {
	eSize := unsafe.Sizeof(C.double(0.0))
	ptr := unsafe.Pointer(uintptr(a) + uintptr(i)*eSize)
	return float64(*(*C.double)(ptr))
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
