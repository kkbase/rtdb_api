//go:build windows

package rtdb_api

import (
	"syscall"
	"unsafe"
)

func UTF16PtrFromString(_path string) (unsafe.Pointer, error) {
	cPath, err := syscall.UTF16PtrFromString(path)
	return unsafe.Pointer(cPath), err
}
