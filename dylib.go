package rtdb_api

import (
	"fmt"
)

// #cgo CFLAGS: -DPNG_DEBUG=1 -I./cinclude
// #cgo CXXFLAGS: -std=c++11
// #include "dylib.h"
import "C"

type FuncPtr uintptr

func LoadLibrary(path string) (C.DylibHandle, error) {
	h := C.load_library(C.CString(path))
	if h.handle == nil {
		return C.DylibHandle{}, fmt.Errorf("LoadLibrary failed")
	}
	return h, nil
}

func GetProcAddress(lib C.DylibHandle, name string) (FuncPtr, error) {
	p := C.get_function(lib, C.CString(name))
	if p == nil {
		return 0, fmt.Errorf("GetProcAddress failed for %s", name)
	}
	return FuncPtr(p), nil
}

func FreeLibrary(lib C.DylibHandle) {
	C.free_library(lib)
}
