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

// RtdbGetApiVersionWarp 取得 rtdbapi 库的版本号
// \param [out]  major   主版本号
// \param [out]  minor   次版本号
// \param [out]  beta    发布版本号
// \return rtdb_error
// \remark 如果返回的版本号与 rtdb.h 中定义的不匹配(RTDB_API_XXX_VERSION)，则应用程序使用了错误的库。
//
//	应输出一条错误信息并退出，否则可能在调用某些 api 时会导致崩溃
func RtdbGetApiVersionWarp() (int32, int32, int32, RtdbError) {
	major, minor, beta := C.rtdb_int32(0), C.rtdb_int32(0), C.rtdb_int32(0)
	err := C.rtdb_get_api_version_warp(&major, &minor, &beta)
	return int32(major), int32(minor), int32(beta), RtdbError(err)
}

// RtdbSetOptionWarp 配置 api 行为参数，参见枚举 \ref RtdbApiOption
// \param [in] type  选项类别
// \param [in] value 选项值
// \return rtdb_error
// \remark 选项设置后在下一次调用 api 时才生效
func RtdbSetOptionWarp(optionType RtdbApiOption, value int32) RtdbError {
	err := C.rtdb_set_option_warp(C.rtdb_int32(optionType), C.rtdb_int32(value))
	return RtdbError(err)
}
