package main

// #cgo CFLAGS: -I/usr/lib/
// #cgo LDFLAGS: -L/usr/lib/ -lstdc++ -lhighs
// #include "call_highs.hpp"
import "C"
import "fmt"

func main() {
	C.minimal_api()

	x := C.first(C.int(5))
	fmt.Println(x)
}
