//go:build linux

// 一些跨平台函数

package rtdb_api

import "unsafe"

// UTF16PtrFromString 占位用的，此为Windows函数，Linux中不存在
func UTF16PtrFromString(_path string) (unsafe.Pointer, error) {
	return nil, nil
}
