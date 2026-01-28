package rtdb_api

// #include<stdlib.h>
// #include<string.h>
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

	cStr := C.CString(s)
	defer C.free(unsafe.Pointer(cStr))

	C.strncpy(p, cStr, C.size_t(n-1))
	p[n-1] = 0 // 确保以'\0'结尾
}

/*
func GoStringToCCharArray(s string, p *C.char, n int) {
	if p == nil || n <= 0 {
		return
	}

	// 直接使用C.CString的理念，但限制长度
	cStr := C.CString(s)
	defer C.free(unsafe.Pointer(cStr))

	// 安全拷贝
	dst := unsafe.Slice(unsafe.Pointer(p), n)
	src := unsafe.Slice(unsafe.Pointer(cStr), len(s)+1) // 包含'\0'

	copied := copy(dst, src)
	if copied == n {
		dst[n-1] = 0 // 确保终止符
	}
}

*/

func RtdbErrorListToErrorList(errs []RtdbError) []error {
	rtn := make([]error, 0)
	for _, err := range errs {
		rtn = append(rtn, err.GoError())
	}
	return rtn
}

// SafeSlice 安全获取切片子集，自动处理越界问题
func SafeSlice[T any](slice []T, start, count int32) []T {
	// 处理空切片或无效参数
	if slice == nil || start < 0 || count <= 0 {
		return []T{}
	}

	// 转换为 int 类型便于操作
	s := int(start)
	c := int(count)
	length := len(slice)

	// 检查起始位置是否超出范围
	if s >= length {
		return []T{}
	}

	// 计算实际结束位置
	end := s + c
	if end > length {
		end = length
	}

	// 返回有效的子切片
	return slice[s:end]
}

// BoolToInt64 bool转换为Int
func BoolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
