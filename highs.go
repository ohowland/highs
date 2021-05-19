package main

// #cgo CFLAGS: -I/usr/lib/
// #cgo LDFLAGS: -L/usr/lib/ -lstdc++ -lhighs
// #include "call_highs.h"
import "C"

func main() {
	C.minimal_api()
}
