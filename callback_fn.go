package rtdb_api

// #cgo CFLAGS: -DPNG_DEBUG=1 -I./cinclude
// #cgo CXXFLAGS: -std=c++11
// #include <stdlib.h>
// #include "gofn.h"
import "C"
import "unsafe"

//export goSubscribeTagsEx
func goSubscribeTagsEx(
	eventType C.rtdb_uint32,
	handle C.rtdb_int32,
	param unsafe.Pointer,
	count C.rtdb_int32,
	ids *C.rtdb_int32,
	what C.rtdb_int32,
) C.rtdb_error {
	return C.rtdb_error(0)
}
