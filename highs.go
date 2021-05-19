package main

// #cgo pkg-config: highs
// #include "highs-interface.h"
import "C"

func main() {
	C.full_api()
}
