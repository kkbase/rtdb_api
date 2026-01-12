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
