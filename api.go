package rtdb_api

// #cgo CFLAGS: -DPNG_DEBUG=1 -I./cinclude
// #cgo CXXFLAGS: -std=c++11
// #include "rtdb.h"
// #include "rtdbapi.h"
// #include "rtdb_error.h"
import "C"

func init() {

}

type RtdbError int32

func RtdbGetApiVersion() (int32, int32, int32, RtdbError) {
	major, minor, beta := C.rtdb_int32(0), C.rtdb_int32(0), C.rtdb_int32(0)
	err := C.rtdb_get_api_version(&major, &minor, &beta)
	return int32(major), int32(minor), int32(beta), RtdbError(err)
}
