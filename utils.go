package rtdb_api

import "C"
import "unsafe"

func CCharArrayToString(p *C.char, n int) string {
	b := C.GoBytes(unsafe.Pointer(p), C.int(n))
	for i, v := range b {
		if v == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

func GoStringToCCharArray(s string, p *C.char, n int) {
	if p == nil || n <= 0 {
		return
	}

	b := []byte(s)

	// 最多拷贝 n-1 个字节，预留 '\0'
	max := n - 1
	if len(b) > max {
		b = b[:max]
	}

	base := uintptr(unsafe.Pointer(p))

	// 逐字节拷贝
	for i := 0; i < len(b); i++ {
		*(*byte)(unsafe.Pointer(base + uintptr(i))) = b[i]
	}

	// 补 '\0'
	*(*byte)(unsafe.Pointer(base + uintptr(len(b)))) = 0

	// （可选）清零剩余空间
	for i := len(b) + 1; i < n; i++ {
		*(*byte)(unsafe.Pointer(base + uintptr(i))) = 0
	}
}
