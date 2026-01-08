package rtdb_api

// windows,amd64 LDFLAGS: -L./clibrary/windows_amd64 -lrtdbapi -lstdc++
// linux,amd64 LDFLAGS: -L./clibrary/linux_amd64 -lrtdbapi -lstdc++ -Wl,-rpath,'$ORIGIN/clibrary/linux_amd64'
// linux,arm64 LDFLAGS: -L./clibrary/linux_arm64 -lrtdbapi -lstdc++ -Wl,-rpath,'$ORIGIN/clibrary/linux_arm64'

// #cgo CFLAGS: -DPNG_DEBUG=1 -I./cinclude
// #cgo CXXFLAGS: -std=c++11
// #include "rtdb.h"
// #include "rtdbapi.h"
// #include "rtdb_error.h"
import "C"

type RtdbError int32

func RtdbGetApiVersion() (int32, int32, int32, RtdbError) {
	//	major, minor, beta := C.rtdb_int32(0), C.rtdb_int32(0), C.rtdb_int32(0)
	//	err := C.rtdb_get_api_version(&major, &minor, &beta)
	// 	return int32(major), int32(minor), int32(beta), RtdbError(err)
	return 0, 0, 0, 0
}
