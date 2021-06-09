package highs

// #cgo pkg-config: highs
// #include <stdlib.h>
// #include <stdio.h>
// #include "interfaces/highs_c_api.h"
import "C"
import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

// Highs ecapsulates the HiGHS C API
type Highs struct {
	mutex       *sync.Mutex
	obj         unsafe.Pointer
	allocs      map[highsPtr]unsafe.Pointer
	dims        dims
	cols        []float64
	bounds      [][2]float64
	rows        [][]float64
	integrality []int
}

type dims struct {
	ArStartSize int
	ArIndexSize int
}

type highsPtr int

const (
	pCols highsPtr = iota
	pColLbs
	pColUbs
	pArStart
	pArIndex
	pArValue
	pRowLbs
	pRowUbs
	pIntg
)

// New returns an allocated Highs object
func New(cost_coefficients []float64, bounds [][2]float64, constraints [][]float64, integrality []int) (*Highs, error) {
	h := &Highs{
		obj:         C.Highs_create(),
		allocs:      make(map[highsPtr]unsafe.Pointer),
		cols:        cost_coefficients,
		bounds:      bounds,
		rows:        constraints,
		integrality: integrality,
	}

	runtime.SetFinalizer(h, func(h *Highs) {
		h.destroy()
	})

	// TODO: how to check for a bad malloc?
	return h, nil
}

// destroy deallocates the Highs object and all heap allocations made by CGO.
func (h *Highs) destroy() {
	if h.obj != nil {
		C.Highs_destroy(h.obj)
		h.obj = nil
	}

	for i, ptr := range h.allocs {
		cFree(ptr)
		delete(h.allocs, i)
	}
}

// allocate take the linear program defined by go slices and converts it to C array heap allocations.
func (h *Highs) allocate() {
	h.allocateColumns()
	h.allocateRows()
	h.allocateIntegrality()
}

func (h *Highs) allocateColumns() error {
	n := len(h.cols)
	if n < len(h.bounds) {
		return errors.New("columns are under-bounded")
	}

	if n > len(h.bounds) {
		return errors.New("columns are over-bounded")
	}

	h.validate(pCols, cMalloc(n, C.double(0)))
	cSetArrayDoubles(h.allocs[pCols], h.cols)

	h.validate(pColLbs, cMalloc(n, C.double(0)))
	cSetArrayDoubles(h.allocs[pColLbs], h.GetLowerBounds())

	h.validate(pColUbs, cMalloc(n, C.double(0)))
	cSetArrayDoubles(h.allocs[pColUbs], h.GetUpperBounds())

	return nil
}

func (h *Highs) validate(n highsPtr, p unsafe.Pointer) {
	old_p, found := h.allocs[n]
	if found {
		log.Println("reallocation of mem:", p)
		cFree(old_p)
	}

	h.allocs[n] = p
}

func (h *Highs) allocateRows() error {
	if len(h.rows[0])-2 != len(h.cols) {
		return errors.New("row size mismatch len(row[i]) != len(col)")
	}

	rows, lbs, ubs := separateBounds(h.rows) // [lb, row constraints..., ub]
	pm := packMatrix(rows)

	h.dims.ArStartSize = len(pm.arStart)
	h.dims.ArIndexSize = len(pm.arIndex)

	h.validate(pArStart, cMalloc(len(pm.arStart), C.int(0)))
	cSetArrayInts(h.allocs[pArStart], pm.arStart)

	h.validate(pArIndex, cMalloc(len(pm.arIndex), C.int(0)))
	cSetArrayInts(h.allocs[pArIndex], pm.arIndex)

	h.validate(pArValue, cMalloc(len(pm.arValue), C.double(0)))
	cSetArrayDoubles(h.allocs[pArValue], pm.arValue)

	m := len(h.rows)
	h.validate(pRowLbs, cMalloc(m, C.double(0)))
	cSetArrayDoubles(h.allocs[pRowLbs], lbs)

	h.validate(pRowUbs, cMalloc(m, C.double(0)))
	cSetArrayDoubles(h.allocs[pRowUbs], ubs)

	return nil
}

func (h *Highs) allocateIntegrality() error {
	n := len(h.cols)
	if n != len(h.integrality) {
		return errors.New("integrality len does not match column len")
	}

	h.validate(pIntg, cMalloc(n, C.int(0)))
	cSetArrayInts(h.allocs[pIntg], h.integrality)

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
		C.int(h.dims.ArIndexSize),
		C.int(1),
		(*C.double)(h.allocs[pCols]),
		(*C.double)(h.allocs[pColLbs]),
		(*C.double)(h.allocs[pColUbs]),
		(*C.double)(h.allocs[pRowLbs]),
		(*C.double)(h.allocs[pRowUbs]),
		(*C.int)(h.allocs[pArStart]),
		(*C.int)(h.allocs[pArIndex]),
		(*C.double)(h.allocs[pArValue]),
		(*C.int)(h.allocs[pIntg]))

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
		C.int(h.dims.ArIndexSize),
		C.int(1),
		(*C.double)(h.allocs[pCols]),
		(*C.double)(h.allocs[pColLbs]),
		(*C.double)(h.allocs[pColUbs]),
		(*C.double)(h.allocs[pRowLbs]),
		(*C.double)(h.allocs[pRowUbs]),
		(*C.int)(h.allocs[pArStart]),
		(*C.int)(h.allocs[pArIndex]),
		(*C.double)(h.allocs[pArValue]))

	return (SolutionStatus)(s)
}

func (h *Highs) callLpSolver() (Solution, error) {
	return Solution{}, nil
}

func (h *Highs) Run() SolutionStatus {
	s := C.Highs_run(h.obj)
	return SolutionStatus(s)
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
	defer cFree(pColValue)

	pColDual := cMalloc(n, C.double(0))
	defer cFree(pColDual)

	pRowValue := cMalloc(m, C.double(0))
	defer cFree(pRowValue)

	pRowDual := cMalloc(m, C.double(0))
	defer cFree(pRowDual)

	C.Highs_getSolution(h.obj, (*C.double)(pColValue), (*C.double)(pColDual), (*C.double)(pRowValue), (*C.double)(pRowDual))

	s := NewSolution()
	s.colValue = copyDoubles(pColValue, n)
	s.colDual = copyDoubles(pColDual, n)
	s.rowValue = copyDoubles(pRowValue, m)
	s.rowDual = copyDoubles(pRowDual, m)

	return s
}

func (h *Highs) PrimalColumnSolution() []float64 {
	s := h.GetSolution()
	return s.colValue
}

func (h *Highs) GetModelStatus() ModelStatus {
	s := C.Highs_getModelStatus(h.obj)
	return ModelStatus(s)
}

func (h *Highs) SetObjectiveSense(s Sense) {
	C.Highs_changeObjectiveSense(h.obj, C.int(s))
}

func (h *Highs) GetObjectiveSense() Sense {
	pS := cMalloc(1, C.int(0))
	defer cFree(pS)
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

func (h *Highs) SetBoolOptionValue(opt string, val bool) {

	pOpt := C.CString(opt)
	defer cFree(unsafe.Pointer(pOpt))

	var v C.int
	if val {
		v = C.int(1)
	} else {
		v = C.int(0)
	}

	C.Highs_setBoolOptionValue(h.obj, pOpt, v)
}

func (h *Highs) GetBoolOptionValue(opt string) bool {
	pOpt := C.CString(opt)
	defer cFree(unsafe.Pointer(pOpt))

	pVal := cMalloc(1, C.int(0))
	defer cFree(pVal)

	C.Highs_getBoolOptionValue(h.obj, pOpt, (*C.int)(pVal))

	return int(*(*C.int)(pVal)) > 0
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
