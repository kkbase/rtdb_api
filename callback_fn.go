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

//export goSnapsEventEx
func goSnapsEventEx(
	eventType C.rtdb_uint32,
	handle C.rtdb_int32,
	param unsafe.Pointer,
	count C.rtdb_int32,
	ids *C.rtdb_int32,
	datetimes *C.rtdb_timestamp_type,
	subtimes *C.rtdb_subtime_type,
	values *C.rtdb_float64,
	status *C.rtdb_int64,
	qualities *C.rtdb_int16,
	errors *C.rtdb_error,
) C.rtdb_error {
	return C.rtdb_error(0)
}
