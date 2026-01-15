//go:build linux

// 一些跨平台函数

package rtdb_api

import (
	"runtime"
	"unsafe"
)

// UTF16PtrFromString 占位用的，此为Windows函数，Linux中不存在
func UTF16PtrFromString(_path string) (unsafe.Pointer, error) {
	if runtime.GOOS == "linux" {
		panic("Linux中无法调用此函数，此函数为占位函数")
	}
	return nil, nil
}
