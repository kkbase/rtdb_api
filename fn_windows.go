//go:build windows

// 一些跨平台函数

package rtdb_api

import (
	"syscall"
	"unsafe"
)

// UTF16PtrFromString 将一个string转换成Windows中的宽字符
func UTF16PtrFromString(_path string) (unsafe.Pointer, error) {
	cPath, err := syscall.UTF16PtrFromString(path)
	return unsafe.Pointer(cPath), err
}
