package rtdb_api

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestRtdbGetApiVersion(t *testing.T) {
	major, minor, beta, err := RtdbGetApiVersionWarp()
	fmt.Println(major, minor, beta, err)
}

func TestRtdbConnectWarp(t *testing.T) {
	handle, err := RtdbConnectWarp("localhost", 6327)
	if !err.IsOk() {
		fmt.Println("创建连接失败", err.Error())
		return
	}
	fmt.Println(handle)
}

func TestRtdbConnectionCountWarp(t *testing.T) {
	handle, err := RtdbConnectWarp("localhost", 6327)
	if !err.IsOk() {
		t.Error("创建连接失败", err.Error())
		return
	}
	count, err := RtdbConnectionCountWarp(handle, 0)
	if !err.IsOk() {
		t.Error("获取Count失败", err.Error())
		return
	}
	fmt.Println("Connect Count: ", count)
}

func CFunc(code string) string {
	// 1. 提取注释块
	commentRe := regexp.MustCompile(`(?s)/\*\*.*?\*/`)
	comment := commentRe.FindString(code)

	// 2. 去掉注释，只保留函数声明部分
	codeNoComment := commentRe.ReplaceAllString(code, "")

	// 3. 把多行压成一行，便于解析
	normalized := strings.Join(strings.Fields(codeNoComment), " ")

	// 示例结果：
	// RTDBAPI rtdb_error RTDBAPI_CALLRULE rtdb_get_own_connection( rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32* socket );

	// 4. 解析函数声明
	funcRe := regexp.MustCompile(
		`RTDBAPI\s+(\w+)\s+RTDBAPI_CALLRULE\s+(\w+)\s*\((.*?)\)\s*;`,
	)
	matches := funcRe.FindStringSubmatch(normalized)
	if len(matches) != 4 {
		return "// parse error: invalid function signature"
	}

	retType := matches[1]
	funcName := matches[2]
	params := strings.TrimSpace(matches[3])

	// 5. wrapper 名称
	wrapName := funcName + "_warp"
	typedefName := funcName + "_fn"

	// 6. 生成 wrapper
	var sb strings.Builder

	if comment != "" {
		sb.WriteString(comment)
		sb.WriteString("\n")
	}

	sb.WriteString(retType)
	sb.WriteString(" RTDBAPI_CALLRULE ")
	sb.WriteString(wrapName)
	sb.WriteString("(")
	sb.WriteString(params)
	sb.WriteString(")\n{\n")

	sb.WriteString("    typedef ")
	sb.WriteString(retType)
	sb.WriteString(" (RTDBAPI_CALLRULE *")
	sb.WriteString(typedefName)
	sb.WriteString(")(")
	sb.WriteString(params)
	sb.WriteString(");\n")

	sb.WriteString("    ")
	sb.WriteString(typedefName)
	sb.WriteString(" fn = (")
	sb.WriteString(typedefName)
	sb.WriteString(")get_function(\"")
	sb.WriteString(funcName)
	sb.WriteString("\");\n")

	sb.WriteString("    return fn(")
	sb.WriteString(extractParamNames(params))
	sb.WriteString(");\n")

	sb.WriteString("}\n")

	return sb.String()
}

func extractParamNames(params string) string {
	if strings.TrimSpace(params) == "" {
		return ""
	}

	parts := strings.Split(params, ",")
	names := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "void" {
			continue
		}
		fields := strings.Fields(p)
		if len(fields) == 0 {
			continue
		}
		// 参数名永远在最后（支持 rtdb_int32* socket 这种）
		names = append(names, fields[len(fields)-1])
	}

	return strings.Join(names, ", ")
}

func TestFunc(t *testing.T) {
	code := `
/**
*
* \brief 批量读取开关量、模拟量快照数值
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表
* \param datetimes 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param values    双精度浮点型数组，输出，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的快照值；否则为 0
* \param states    64 位整型数组，输出，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的快照值；否则为 0
* \param qualities 短整型数组，输出，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量写入开关量、模拟量快照数值
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功写入实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
*                    但它们的时间戳必需是递增的。
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param values    双精度浮点型数组，输入，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
* \param states    64 位整型数组，输入，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_float64* values,
    const rtdb_int64* states,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量写入开关量、模拟量快照数值
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功写入实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
*                    但它们的时间戳必需是递增的。
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param values    双精度浮点型数组，输入，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
* \param states    64 位整型数组，输入，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
*        仅当输入时间戳与当前快照时间戳完全相等时，会替换当前快照的值和质量；
*        其余情况下会调用 rtdbs_put_snapshots()
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_fix_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_float64* values,
    const rtdb_int64* states,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量回溯快照
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功写入实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
*
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param values    双精度浮点型数组，输入，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
* \param states    64 位整型数组，输入，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
* 功能说明：
*       批量将标签点的快照值vtmq改成传入的vtmq，如果传入的时间戳早于当前快照，会删除传入时间戳到当前快照的历史存储值。
*       如果传入的时间戳等于或者晚于当前快照，什么也不做。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_back_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_float64* values,
    const rtdb_int64* states,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量读取坐标实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表
* \param datetimes 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param x         单精度浮点型数组，输出，实时浮点型横坐标数值列表
* \param y         单精度浮点型数组，输出，实时浮点型纵坐标数值列表
* \param qualities 短整型数组，输出，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_coor_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float32* x,
    rtdb_float32* y,
    rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量写入坐标实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param x         单精度浮点型数组，输入，实时浮点型横坐标数值列表
* \param y         单精度浮点型数组，输入，实时浮点型纵坐标数值列表
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时坐标数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_coor_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_float32* x,
    const rtdb_float32* y,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量写入坐标实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param x         单精度浮点型数组，输入，实时浮点型横坐标数值列表
* \param y         单精度浮点型数组，输入，实时浮点型纵坐标数值列表
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时坐标数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*        仅当输入时间戳与当前快照时间戳完全相等时，会替换当前快照的值和质量；
*        其余情况下会调用 rtdbs_put_coor_snapshots()
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_fix_coor_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_float32* x,
    const rtdb_float32* y,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 读取二进制/字符串实时数据
*
* \param handle    连接句柄
* \param id        整型，输入，标签点标识
* \param datetime  整型，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param blob      字节型数组，输出，实时二进制/字符串数值
* \param len       短整型，输出，二进制/字符串数值长度
* \param quality   短整型，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_blob_snapshot64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    rtdb_byte* blob,
    rtdb_length_type* len,
    rtdb_int16* quality
);

/**
*
* \brief 批量读取二进制/字符串实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、blobs、lens、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识
* \param datetimes 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param blobs     字节型指针数组，输出，实时二进制/字符串数值
* \param lens      短整型数组，输入/输出，二进制/字符串数值长度，
*                    输入时表示对应的 blobs 指针指向的缓冲区长度，
*                    输出时表示实际得到的 blob 长度，如果 blob 的长度大于缓冲区长度，会被截断。
* \param qualities 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_blob_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_byte* const* blobs,
    rtdb_length_type* lens,
    rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 写入二进制/字符串实时数据
*
* \param handle    连接句柄
* \param id        整型，输入，标签点标识
* \param datetime  整型，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param blob      字节型数组，输入，实时二进制/字符串数值
* \param len       短整型，输入，二进制/字符串数值长度，超过一个页大小数据将被截断。
* \param quality   短整型，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_blob_snapshot64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    const rtdb_byte* blob,
    rtdb_length_type len,
    rtdb_int16 quality
);

/**
*
* \brief 批量写入二进制/字符串实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、blobs、lens、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param blobs     字节型指针数组，输入，实时二进制/字符串数值
* \param lens      短整型数组，输入，二进制/字符串数值长度，
*                    表示对应的 blobs 指针指向的缓冲区长度，超过一个页大小数据将被截断。
* \param qualities 短整型数组，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_blob_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_byte* const* blobs,
    const rtdb_length_type* lens,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量读取datetime类型标签点实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、dtvalues、dtlens、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识
* \param datetimes 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param dtvalues  字节型指针数组，输出，实时datetime数值
* \param dtlens    短整型数组，输入/输出，datetime数值长度，
*                    输入时表示对应的 dtvalues 指针指向的缓冲区长度，
* \param qualities 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \param type      短整型，输入，所有标签点的显示类型，如“yyyy-mm-dd hh:mm:ss.000”的type为1，默认类型1，
*                    “yyyy/mm/dd hh:mm:ss.000”的type为2
*                    如果不传type，则按照标签点属性显示，否则按照type类型显示
* \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_datetime_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_byte* const* dtvalues,
    rtdb_length_type* dtlens,
    rtdb_int16* qualities,
    rtdb_error* errors,
    rtdb_int16 type
);

/**
*
* \brief 批量插入datetime类型标签点数据
*
* \param handle      连接句柄
* \param count       整型，输入/输出，标签点个数，
*                      输入时表示 ids、datetimes、ms、dtvalues、dtlens、qualities、errors的长度，
*                      输出时表示成功写入的标签点个数
* \param ids         整型数组，输入，标签点标识
* \param datetimes   整型数组，输入，实时值时间列表
*                      表示距离1970年1月1日08:00:00的秒数
* \param ms          短整型数组，输入，实时数值时间列表，
*                      对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param dtvalues    字节型指针数组，输入，datetime标签点的值
* \param dtlens      短整型数组，输入，数值长度
* \param qualities   短整型数组，输入，实时数值品质，，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors      无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 被接口只对数据类型 RTDB_DATETIME 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_datetime_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_byte* const* dtvalues,
    const rtdb_length_type* dtlens,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量标签点快照改变的通知订阅
*
* \param handle         连接句柄
* \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
*                           输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
* \param ids            整型数组，输入，标签点标识列表。
* \param options        订阅选项
*                           RTDB_O_AUTOCONN 自动重连
* \param param          用户自定义参数
* \param callback       rtdbs_snaps_event_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
*                       当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                       作为event_type取值调用回掉函数后退出订阅。
*                       当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                       作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
*                       网络中断期间回掉函数调用频率为最少3秒
*                       event_type参数值含义如下：
*                         RTDB_E_DATA        标签点快照改变
*                         RTDB_E_DISCONNECT  订阅客户端与数据库网络断开
*                         RTDB_E_RECOVERY    订阅客户端与数据库网络及订阅恢复
*                         RTDB_E_CHANGED     客户端修改订阅标签点信息，即通过rtdbs_change_subscribe_snapshots修改订阅信息的结果
*                       handle 产生订阅回掉的连接句柄，调用rtdbs_subscribe_snapshots_ex时的handle参数
*                       param  用户自定义参数，调用rtdbs_subscribe_snapshots_ex时的param参数
*                       count  event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示ids，datetimes,values等的数量
*                              event_type为其它值时，count值为0
*                       ids    event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示快照改变的标签点ID，数量由count指定
*                              event_type为其它值时，ids值为NULL
*                       datetimes event_type为RTDB_E_DATA时表示快照时间，数量由count指定
*                                 event_type为RTDB_E_CHANGED时表示changed_types，即通过rtdbs_change_subscribe_snapshots传入的changed_types
*                                 event_type为其它值时，datetimes值为NULL
*                       ms     event_type为RTDB_E_DATA时表示快照的毫秒，数量由count指定
*                              event_type为其它值时，ms值为NULL
*                       values event_type为RTDB_E_DATA时表示浮点数据类型快照值，数量由count指定
*                              event_type为RTDB_E_CHANGED时，表示delta_values，即通过rtdbs_change_subscribe_snapshots传入的delta_values
*                              event_type为其它值时，values值为NULL
*                       states event_type为RTDB_E_DATA时表示整形数据类型快照值，数量由count指定
*                              event_type为RTDB_E_CHANGED时，表示delta_states，通过rtdbs_change_subscribe_snapshots传入的delta_states
*                              event_type为其它值时，states值为NULL
*                       qualities event_type为RTDB_E_DATA时表示快照质量码，数量由count指定
*                              event_type为其它值时，qualities值为NULL
*                       errors event_type为RTDB_E_DATA时表示快照错误码，数量由count指定
*                              event_type为RTDB_E_CHANGED时，表示修改结果对应的错误码，数量由count指定
*                              event_type为其它值时，errors值为NULL
* \param errors         无符号整型数组，输出，
*                           写入实时数据的返回值列表，参考rtdb_error.h
* \remark   用户须保证 ids、errors 的长度与 count 一致。
*        用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
*        否则返回 RtE_OTHER_SDK_DOING 错误。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_subscribe_snapshots_ex64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_uint32 options,
    void* param,
    rtdbs_snaps_event_ex64 callback,
    rtdb_error* errors
);

/**
*
* \brief 批量标签点快照改变的通知订阅
*
* \param handle         连接句柄
* \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
*                           输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
* \param ids            整型数组，输入，标签点标识列表。
* \param delta_values   double型数组，输入，订阅浮点类型标签点的容差值，变化超过设置的容差值才会推送
* \param delta_values   整型数组，输入，订阅整型标签点的容差值，变化超过设置的容差值才会推送
* \param options        订阅选项
*                           RTDB_O_AUTOCONN 自动重连
* \param param          用户自定义参数
* \param callback       rtdbs_snaps_event_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
*                         当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                         作为event_type取值调用回掉函数后退出订阅。
*                         当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                         作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
*                         网络中断期间回掉函数调用频率为最少3秒
*                         event_type参数值含义如下：
*                           RTDB_E_DATA        标签点快照改变
*                           RTDB_E_DISCONNECT  订阅客户端与数据库网络断开
*                           RTDB_E_RECOVERY    订阅客户端与数据库网络及订阅恢复
*                           RTDB_E_CHANGED     客户端修改订阅标签点信息，即通过rtdbs_change_subscribe_snapshots修改订阅信息的结果
*                         handle 产生订阅回掉的连接句柄，调用rtdbs_subscribe_snapshots_ex时的handle参数
*                         param  用户自定义参数，调用rtdbs_subscribe_snapshots_ex时的param参数
*                         count  event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示ids，datetimes,values等的数量
*                                event_type为其它值时，count值为0
*                         ids    event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示快照改变的标签点ID，数量由count指定
*                                event_type为其它值时，ids值为NULL
*                         datetimes event_type为RTDB_E_DATA时表示快照时间，数量由count指定
*                                   event_type为RTDB_E_CHANGED时表示changed_types，即通过rtdbs_change_subscribe_snapshots传入的changed_types
*                                   event_type为其它值时，datetimes值为NULL
*                         ms     event_type为RTDB_E_DATA时表示快照的毫秒，数量由count指定
*                                event_type为其它值时，ms值为NULL
*                         values event_type为RTDB_E_DATA时表示浮点数据类型快照值，数量由count指定
*                                event_type为RTDB_E_CHANGED时，表示delta_values，即通过rtdbs_change_subscribe_snapshots传入的delta_values
*                                event_type为其它值时，values值为NULL
*                         states event_type为RTDB_E_DATA时表示整形数据类型快照值，数量由count指定
*                                event_type为RTDB_E_CHANGED时，表示delta_states，通过rtdbs_change_subscribe_snapshots传入的delta_states
*                                event_type为其它值时，states值为NULL
*                         qualities event_type为RTDB_E_DATA时表示快照质量码，数量由count指定
*                                event_type为其它值时，qualities值为NULL
*                         errors event_type为RTDB_E_DATA时表示快照错误码，数量由count指定
*                                event_type为RTDB_E_CHANGED时，表示修改结果对应的错误码，数量由count指定
*                                event_type为其它值时，errors值为NULL
* \param errors         无符号整型数组，输出，
*                           写入实时数据的返回值列表，参考rtdb_error.h
* \remark delta_values和delta_states可以为空指针，表示不设置容差值。 只有两个参数都不为空时，设置容差值才会生效。
            用户须保证 ids、delta_values、delta_states、errors 的长度与 count 一致
*           用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
*           否则返回 RtE_OTHER_SDK_DOING 错误。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_subscribe_delta_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_float64* delta_values,
    const rtdb_int64* delta_states,
    rtdb_uint32 options,
    void* param,
    rtdbs_snaps_event_ex64 callback,
    rtdb_error* errors
);

/**
*
* \brief 批量修改订阅标签点信息
*
* \param handle         连接句柄
* \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
*                           输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
* \param ids            整型数组，输入，标签点标识列表。
* \param delta_values   double型数组，输入，订阅浮点类型标签点的容差值，变化超过设置的容差值才会推送
* \param delta_values   整型数组，输入，订阅整型标签点的容差值，变化超过设置的容差值才会推送
* \param changed_types  整型数组，输入，修改类型，参考RTDB_SUBSCRIBE_CHANGE_TYPE
* \param errors         异步调用，保留参数，暂时不启用
* \remark   用户须保证 ids、delta_values、delta_states、errors 的长度与 count 一致。
*               可以同时添加、修改、删除订阅的标签点信息，
*               delta_values和delta_states，可以为空指针，为空，则表示不设置容差值，即写入新数据即推送
*               只有delta_values和delta_states都不为空时，设置的容差值才有效。
*               用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
*               否则返回 RtE_OTHER_SDK_DOING 错误。
*               此方法是异步方法，当网络中断等异常情况时，会通过方法的返回值返回错误，参考rtdb_error.h。
*               当方法返回值为RtE_OK时，表示已经成功发送给数据库，但是并没有等待修改结果。
*               数据库的修改结果，会异步通知给api的回调函数，通过rtdbs_snaps_event_ex的RTDB_E_CHANGED事件通知修改结果
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_change_subscribe_snapshots(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_float64* delta_values,
    const rtdb_int64* delta_states,
    const rtdb_int32* changed_types,
    rtdb_error* errors
);

/**
*
* \brief 取消标签点快照更改通知订阅
*
* \param handle    连接句柄
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_cancel_subscribe_snapshots(
  rtdb_int32 handle
  );

/**
* 命名：rtdbs_get_named_type_snapshot32
* 功能：获取自定义类型测点的单个快照
* 参数：
*        [handle]    连接句柄
*        [id]        整型，输入，标签点标识
*        [datetime]  整型，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
*        [object]    字节型数组，输出，实时自定义类型标签点的数值
*        [length]    短整型，输入/输出，自定义类型标签点的数值长度
*        [quality]   短整型，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_named_type_snapshot64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    void* object,
    rtdb_length_type* length,
    rtdb_int16* quality
);

/**
* 命名：rtdbs_get_named_type_snapshots32
* 功能：批量获取自定义类型测点的快照
* 参数：
*        [handle]    连接句柄
*        [count]     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、objects、lengths、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
*        [ids]       整型数组，输入，标签点标识
*        [datetimes] 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
*        [objects]   指针数组，输出，自定义类型标签点数值
*        [lengths]   短整型数组，输入/输出，自定义类型标签点数值长度，
*                    输入时表示对应的 objects 指针指向的缓冲区长度，
*                    输出时表示实际得到的 objects 长度，如果 objects 的长度大于缓冲区长度，会被截断。
*        [qualities] 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [errors]    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_named_type_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    void* const* objects,
    rtdb_length_type* lengths,
    rtdb_int16* qualities,
    rtdb_error* errors
);

/**
* 命名：rtdbs_put_named_type_snapshot32
* 功能：写入单个自定义类型标签点的快照
* 参数：
*        [handle]    连接句柄
*        [id]        整型，输入，标签点标识
*        [datetime]  整型，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
*        [object]    void类型数组，输入，自定义类型标签点数值
*        [length]    短整型，输入，自定义类型标签点数值长度，超过一个页大小数据将被截断。
*        [quality]   短整型，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_named_type_snapshot64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    const void* object,
    rtdb_length_type length,
    rtdb_int16 quality
);

/**
* 命名：rtdbs_put_named_type_snapshots32
* 功能：批量写入自定义类型标签点的快照
* 参数：
*        [handle]    连接句柄
*        [count]     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、objects、lengths、qualities、errors 的长度，
*                    输出时表示成功写入实时值的标签点个数
*        [ids]       整型数组，输入，标签点标识
*        [datetimes] 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
*        [objects]   void类型指针数组，输入，自定义类型标签点数值
*        [lengths]   短整型数组，输入，自定义类型标签点数值长度，
*                    表示对应的 objects 指针指向的缓冲区长度，超过一个页大小数据将被截断。
*        [qualities] 短整型数组，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [errors]    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_named_type_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const void* const* objects,
    const rtdb_length_type* lengths,
    const rtdb_int16* qualities,
    rtdb_error* errors
);
`
	sp := strings.Split(code, "/**")
	for _, fn := range sp {
		ff := "/**\n" + fn
		fmt.Println(CFunc(ff))
	}
}

func TestHello(t *testing.T) {
	fmt.Println("323")
}
