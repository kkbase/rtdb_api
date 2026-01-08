package rtdb_api

// #cgo CFLAGS: -DPNG_DEBUG=1 -I./cinclude
// #cgo CXXFLAGS: -std=c++11
// #include "api.h"
import "C"
import (
	_ "embed"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed clibrary/linux_amd64/librtdbapi.so
var LinuxAmd64RtdbSo []byte

// TODO, Windows, Linux

func init() {
	// 跨平台加载SO路径
	data := make([]byte, 0)
	name := ""
	if runtime.GOOS == "linux" {
		if runtime.GOARCH == "amd64" {
			data = LinuxAmd64RtdbSo
			name = "librtdb.so"
		} else if runtime.GOARCH == "arm64" {

		} else {
			panic("不支持的平台，分支不可达")
		}
	} else if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
		} else {
			panic("不支持的平台，分支不可达")
		}
	} else {
		panic("不支持的平台，分支不可达")
	}

	// 将动态库写入到临时文件
	path := filepath.Join(os.TempDir(), name)
	if err := os.WriteFile(path, data, 0755); err != nil {
		panic(err)
	}

	// 加载动态库
	C.load_library(C.CString(path))
}

type RtdbError int32

func RtdbGetApiVersion() (int32, int32, int32, RtdbError) {
	major, minor, beta := C.rtdb_int32(0), C.rtdb_int32(0), C.rtdb_int32(0)
	err := C.rtdb_get_api_version_warp(&major, &minor, &beta)
	return int32(major), int32(minor), int32(beta), RtdbError(err)
}
