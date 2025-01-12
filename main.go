package main

/*
#cgo LDFLAGS: ./lib/diffgen/target/release/libdiffgen.a -ldl
#include "./lib/diffgen.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
)

func RustDiff(before, after string) string {
	diffStr := C.diff(C.CString(string(before)), C.CString(string(after)))
	defer C.free(unsafe.Pointer(diffStr))
	return C.GoString(diffStr)
}

func GoDiff(before, after string) string {
	edits := myers.ComputeEdits("", before, after)
	if len(edits) == 0 {
		return ""
	}
	return fmt.Sprint(gotextdiff.ToUnified("before", "after", before, edits))
}

func main() {
	str1 := C.CString("world")
	str2 := C.CString("this is code from the static library")
	defer C.free(unsafe.Pointer(str1))
	defer C.free(unsafe.Pointer(str2))

	C.hello(str1)
}
