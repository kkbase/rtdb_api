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
* \brief 获取单个标签点在一段时间范围内的存储值数量.
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
* \param count         整型，输出，返回上述时间范围内的存储值数量
* \remark 由 datetime1、ms1 形成的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_archived_values_count64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_int32* count
);

/**
*
* \brief 获取单个标签点在一段时间范围内的真实的存储值数量.
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
* \param count         整型，输出，返回上述时间范围内的存储值数量
* \remark 由 datetime1、ms1 形成的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_archived_values_real_count64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_int32* count
);

/**
*
* \brief 读取单个标签点一段时间内的储存数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、values、states、qualities 的长度；
*                        输出时返回实际得到的数值个数
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，历史浮点型数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史存储值；否则为 0
* \param states        64 位整数数组，输出，历史整型数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史存储值；否则为 0
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);


/**
*
* \brief 逆向读取单个标签点一段时间内的储存数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、values、states、qualities 的长度；
*                        输出时返回实际得到的数值个数
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，历史浮点型数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史存储值；否则为 0
* \param states        64 位整数数组，输出，历史整型数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史存储值；否则为 0
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_values_backward64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 读取单个标签点一段时间内的坐标型储存数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、x、y、qualities 的长度；
*                        输出时返回实际得到的数值个数
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param x             单精度浮点型数组，输出，浮点型横坐标历史数值列表
* \param y             单精度浮点型数组，输出，浮点型纵坐标历史数值列表
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、x、y、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_coor_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float32* x,
    rtdb_float32* y,
    rtdb_int16* qualities
);

/**
*
* \brief 逆向读取单个标签点一段时间内的坐标型储存数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、x、y、qualities 的长度；
*                        输出时返回实际得到的数值个数
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param x             单精度浮点型数组，输出，浮点型横坐标历史数值列表
* \param y             单精度浮点型数组，输出，浮点型纵坐标历史数值列表
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、x、y、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_coor_values_backward64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float32* x,
    rtdb_float32* y,
    rtdb_int16* qualities
);


/**
*
* \brief 开始以分段返回方式读取一段时间内的储存数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
* \param count         整型，输出，返回上述时间范围内的存储值数量
* \param batch_count   整型，输出，每次分段返回的长度，用于继续调用 rtdbh_get_next_archived_values 接口
* \remark 由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_values_in_batches64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_int32* count,
    rtdb_int32* batch_count
);

/**
*
* \brief 分段读取一段时间内的储存数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整形，输入/输出，
*                        输入时表示 datetimes、ms、values、states、qualities 的长度；
*                        输出时表示实际得到的存储值个数。
* \param datetimes     整型数组，输出，历史数值时间列表,
*                        表示距离1970年1月1日08:00:00的秒数
* \param ms            短整型数组，输出，历史数值时间列表，
*                        对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param values        双精度浮点型数组，输出，历史浮点型数值列表，
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的历史存储值；否则为 0
* \param states        64 位整型数组，输出，历史整型数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的历史存储值；否则为 0
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 相符，
*        且 count 不能小于 rtdbh_get_archived_values_in_batches 接口中返回的 batch_count 的值，
*        当返回 RtE_BATCH_END 表示全部数据获取完毕。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_next_archived_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 获取单个标签点的单调递增时间序列历史插值。
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入，表示 datetimes、ms、values、states、qualities 的长度。
* \param datetimes     整型数组，输入，表示需要的单调递增时间列表，
*                        为距离1970年1月1日08:00:00的秒数
* \param ms            短整型数组，输入，对于时间精度为纳秒的标签点，
*                        表示需要的单调递增时间对应的纳秒值；否则忽略。
* \param values        双精度浮点型数组，输出，历史浮点型数值列表，
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的历史插值；否则为 0
* \param states        64 位整型数组，输出，历史整型数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的历史插值；否则为 0
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 相符，
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_timed_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 count,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 获取单个坐标标签点的单调递增时间序列历史插值。
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入，表示 datetimes、ms、x、y、qualities 的长度。
* \param datetimes     整型数组，输入，表示需要的单调递增时间列表，
*                        为距离1970年1月1日08:00:00的秒数
* \param ms            短整型数组，输入，对于时间精度为纳秒的标签点，
*                        表示需要的单调递增时间对应的纳秒值；否则忽略。
* \param x             单精度浮点型数组，输出，浮点型横坐标历史插值数值列表
* \param y             单精度浮点型数组，输出，浮点型纵坐标历史插值数值列表
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、x、y、qualities 的长度与 count 相符，
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_timed_coor_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 count,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    rtdb_float32* x,
    rtdb_float32* y,
    rtdb_int16* qualities
);

/**
*
* \brief 获取单个标签点一段时间内等间隔历史插值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、values、states、qualities 的长度，
*                        即需要的插值个数；输出时返回实际得到的插值个数
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
* \param states        64 位整数数组，输出，整型历史插值数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
* \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_interpo_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 读取单个标签点某个时刻之后一定数量的等间隔内插值替换的历史数值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param interval      整型，输入，插值时间间隔，单位为纳秒
* \param count         整型，输入，表示 datetimes、ms、values、states、qualities 的长度，
*                        即需要的插值个数。
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
* \param states        64 位整数数组，输出，整型历史插值数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
* \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素用于存放起始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_interval_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int64 interval,
    rtdb_int32 count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 读取单个标签点某个时间的历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT、RTDB_INTER 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
*                        RTDB_INTER 取指定时间的内插值数据。
* \param datetime      整型，输入/输出，输入时表示时间秒数；
*                        输出时表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
*                        否则忽略输入，输出时为 0。
* \param value         双精度浮点数，输出，浮点型历史数值
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史值；否则为 0
* \param state         64 位整数，输出，整型历史数值，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史值；否则为 0
* \param quality       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 mode,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    rtdb_float64* value,
    rtdb_int64* state,
    rtdb_int16* quality
);

/**
*
* \brief 读取单个标签点某个时间的坐标型历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT、RTDB_INTER 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
*                        RTDB_INTER 取指定时间的内插值数据。
* \param datetime      整型，输入/输出，输入时表示时间秒数；
*                        输出时表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
*                        否则忽略输入，输出时为 0。
* \param x             单精度浮点型，输出，横坐标历史数值
* \param y             单精度浮点型，输出，纵坐标历史数值
* \param quality       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_coor_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 mode,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    rtdb_float32* x,
    rtdb_float32* y,
    rtdb_int16* quality
);

/**
*
* \brief 读取单个标签点某个时间的二进制/字符串型历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
* \param datetime      整型，输入/输出，输入时表示时间秒数；
*                        输出时表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
*                        否则忽略输入，输出时为 0。
* \param blob          字节型数组，输出，二进制/字符串历史值
* \param len           短整型，输入/输出，输入时表示 blob 的长度，
*                        输出时表示实际获取的二进制/字符串数据长度。
* \param quality       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_blob_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 mode,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    rtdb_byte* blob,
    rtdb_length_type* len,
    rtdb_int16* quality
);

/**
*
* \brief 读取单个标签点一段时间的二进制/字符串型历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
* \param count         整型，输入/输出，输入表示想要查询多少数据
*                        输出表示实际查到多少数据
* \param datetime1     整型，输入，表示开始时间秒数；
* \param ms1           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param datetime2     整型，输入,表示结束时间秒数；
* \param ms2           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param datetimes     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
*                        表示实际取得的历史数值时间纳秒数。
* \param lens          短整型数组，输入/输出，输入时表示 blob 的长度，
*                        输出时表示实际获取的二进制/字符串数据长度。
*                        当blobs为空指针时，表示只获取每条数据的长度，此时会忽略输入的lens
* \param blobs         字节型数组，输出，二进制/字符串历史值。可以设置为空指针，表示只获取每条数据的长度
* \param qualities     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_blob_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_length_type* lens,
    rtdb_byte* const* blobs,
    rtdb_int16* qualities
);

/**
*
* \brief 读取并模糊搜索单个标签点一段时间的二进制/字符串型历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
* \param count         整型，输入/输出，输入表示想要查询多少数据
*                        输出表示实际查到多少数据
* \param datetime1     整型，输入，表示开始时间秒数；
* \param ms1           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param datetime2     整型，输入,表示结束时间秒数；
* \param ms2           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param filter        字符串，输入，支持通配符的模糊搜索字符串，多个模糊搜索的条件通过空格分隔，只针对string类型有效
*                        当filter为空指针时，表示不进行过滤,
*                        限制最大长度为RTDB_EQUATION_SIZE-1，超过此长度会返回错误
* \param datetimes     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
*                        表示实际取得的历史数值时间纳秒数。
* \param lens          短整型数组，输入/输出，输入时表示 blob 的长度，
*                        输出时表示实际获取的二进制/字符串数据长度。
*                        当blobs为空指针时，表示只获取每条数据的长度，此时会忽略输入的lens
* \param blobs         字节型数组，输出，二进制/字符串历史值。可以设置为空指针，表示只获取每条数据的长度
* \param qualities     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_blob_values_filt64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    const char* filter,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_length_type* lens,
    rtdb_byte* const* blobs,
    rtdb_int16* qualities
);

/**
*
* \brief 读取单个标签点某个时间的datetime历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
* \param datetime      整型，输入/输出，输入时表示时间秒数；
*                        输出时表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
*                        否则忽略输入，输出时为 0。
* \param dtblob          字节型数组，输出，datetime历史值
* \param dtlen           短整型，输入/输出，输入时表示 blob 的长度，
*                        输出时表示实际获取的datetime数据长度。
* \param quality       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param type           短整型 datetime字符串的格式类型，默认为-1
* \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_datetime_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 mode,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    rtdb_byte* dtblob,
    rtdb_length_type* dtlen,
    rtdb_int16* quality,
    rtdb_int16 type
);

/**
*
* \brief 读取单个标签点一段时间的时间类型历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
* \param count         整型，输入/输出，输入表示想要查询多少数据
*                        输出表示实际查到多少数据
* \param datetime1     整型，输入，表示开始时间秒数；
* \param ms1           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param datetime2     整型，输入,表示结束时间秒数；
* \param ms2           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param datetimes     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
*                        表示实际取得的历史数值时间纳秒数。
* \param dtlens          短整型数组，输入/输出，输入时表示 blob 的长度，
*                        输出时表示实际获取的时间数据长度。
* \param dtvalues         字节型数组，输出，时间历史值
* \param qualities     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param type          短整型，输入，“yyyy-mm-dd hh:mm:ss.000”的type为1， 同样默认输入格式也为 “yyyy-mm-dd hh:mm:ss.000”
*                       “yyyy/mm/dd hh:mm:ss.000”的type为2
* \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_datetime_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_length_type* lens,
    rtdb_byte* const* blobs,
    rtdb_int16* qualities,
    rtdb_int16 type
);


/**
*
* \brief 写入批量标签点批量时间型历史存储数据
*
* \param handle        连接句柄
* \param count         整型，输入/输出，
*                        输入时表示 ids、datetimes、ms、dtlens、dtvalues、qualities、errors 的长度，
*                        即历史值个数；输出时返回实际写入的数值个数
* \param ids           整型数组，输入，标签点标识
* \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示对应的历史数值时间纳秒；否则忽略。
* \param dtvalues      字节型指针数组，输入，实时时间数值
* \param dtlens        短整型数组，输入，时间数值长度，
*                        表示对应的 dtvalues 指针指向的缓冲区长度，超过一个页大小数据将被截断。
* \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、dtlens、dtvalues、qualities、errors 的长度与 count 一致，
*        本接口仅对数据类型为 RTDB_DATETIME 的标签点有效。
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_datetime_values64(
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
* \brief 获取单个标签点一段时间内的统计值。
*
* \param handle            连接句柄
* \param id                整型，输入，标签点标识
* \param datetime1         整型，输入/输出，输入时表示起始时间秒数。
*                            如果为 0，表示从存档中最早时间的数据开始进行统计。
*                            输出时返回最大值的时间秒数。
* \param ms1               短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                            表示起始时间对应的纳秒，输出时表示最大值的时间纳秒数；否则忽略，返回值为 0
* \param datetime2         整型，输入/输出，输入时表示结束时间秒数。
*                            如果为 0，表示统计到存档中最近时间的数据为止。
*                            输出时返回最小值的时间秒数。
* \param ms2               短整型，如果 id 指定的标签点时间精度为纳秒，
*                            表示结束时间对应的纳秒，输出时表示最小值的时间纳秒数；否则忽略，返回值为 0
* \param max_value         双精度浮点型，输出，表示统计时间段内的最大数值。
* \param min_value         双精度浮点型，输出，表示统计时间段内的最小数值。
* \param total_value       双精度浮点型，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
* \param calc_avg          双精度浮点型，输出，表示统计时间段内的算术平均值。
* \param power_avg         双精度浮点型，输出，表示统计时间段内的加权平均值。
* \remark 由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*        如果输出的最大值或最小值的时间戳秒值为 0，
*        则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_summary_data(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    RTDB_SUMMARY_DATA* summary_data
);

/**
* 命名：rtdbh_summary_in_batches
* \brief 分批获取单一标签点一段时间内的统计值
*
* \param handle            连接句柄
* \param id                整型，输入，标签点标识
* \param count             整形，输入/输出，输入时表示 datatimes1、ms1、datatimes2、ms2、
*                            max_values、min_values、total_values、calc_avgs、power_avgs、errors 的长度，
*                            即分段的个数；输出时表示成功取得统计值的分段个数。
* \param interval          64 位整型，输入，分段时间间隔，单位为纳秒。
*                            如果为纳秒点，输入时间必须大于1纳秒，如果为秒级点，则必须大于1000000000纳秒。
* \param datetimes1        整型数组，输入/输出，输入时第一个元素表示起始时间秒数。
*                            如果为 0，表示从存档中最早时间的数据开始进行统计。
*                            输出时返回各个分段对应的最大值的时间秒数。
* \param ms1               短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                            第一个元素表示起始时间对应的纳秒，
*                            输出时返回各个分段对应的最大值的时间纳秒数；否则忽略，返回值为 0
* \param datetimes2        整型数组，输入/输出，输入时第一个元素表示结束时间秒数。
*                            如果为 0，表示统计到存档中最近时间的数据为止。
*                            输出时返回各个分段对应的最小值的时间秒数。
* \param ms2               短整型数组，如果 id 指定的标签点时间精度为纳秒，
*                            第一个元素表示结束时间对应的纳秒，
*                            输出时返回各个分段对应的最小值的时间纳秒数；否则忽略，返回值为 0
* \param max_values        双精度浮点型数组，输出，表示统计时间段内的最大数值。
* \param min_values        双精度浮点型数组，输出，表示统计时间段内的最小数值。
* \param total_values      双精度浮点型数组，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
* \param calc_avgs         双精度浮点型数组，输出，表示统计时间段内的算术平均值。
* \param power_avgs        双精度浮点型数组，输出，表示统计时间段内的加权平均值。
* \param errors            无符号整型数组，输出，表示各个分段取得统计值的返回值。
* \remark 由 datetimes1[0]、ms1[0] 表示的时间可以大于 datetimes2[0]、ms2[0] 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*        如果输出的最大值或最小值的时间戳秒值为 0，
*        则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_summary_data_in_batches(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_int64 interval,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    RTDB_SUMMARY_DATA* summary_datas,
    rtdb_error* errors
);

/**
*
* \brief 获取单个标签点一段时间内用于绘图的历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param interval      整型，输入，时间区间数量，单位为个，
*                        一般会使用绘图的横轴(时间轴)所用屏幕像素数，
*                        该功能将起始至结束时间等分为 interval 个区间，
*                        并返回每个区间的第一个和最后一个数值、最大和最小数值、一条异常数值；
*                        故参数 count 有可能输出五倍于 interval 的历史值个数，
*                        所以推荐输入的 count 至少是 interval 的五倍。
* \param count         整型，输入/输出，输入时表示 datetimes、ms、values、states、qualities 的长度，
*                        即需要获取的最大历史值个数，输出时返回实际得到的历史值个数。
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史值数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史值；否则为 0
* \param states        64 位整数数组，输出，整型历史值数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史值；否则为 0
* \param qualities     短整型数组，输出，历史值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，用以存放起始及结束时间。
*        第一个元素形成的时间可以大于最后一个元素形成的时间，
*        此时第一个元素表示结束时间，最后一个元素表示开始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_plot_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 interval,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 获取批量标签点在某一时间的历史断面数据
*
* \param handle        连接句柄
* \param ids           整型数组，输入，标签点标识列表
* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT、RTDB_INTER 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
*                        RTDB_INTER 取指定时间的内插值数据。
* \param count         整型，输入，表示 ids、datetimes、ms、values、states、qualities 的长度，即标签点个数。
* \param datetimes     整型数组，输入/输出，输入时表示对应标签点的历史数值时间秒数，
*                        输出时表示根据 mode 实际寻找到的数值时间秒数。
* \param ms            短整型数组，输入/输出，对于时间精度为纳秒的标签点，
*                        输入时表示历史数值时间纳秒数，存放相应的纳秒值，
*                        输出时表示根据 mode 实际寻找到的数值时间纳秒数；否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史值数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史值；否则为 0
* \param states        64 位整数数组，输出，整型历史值数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史值；否则为 0
* \param qualities     短整型数组，输出，历史值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors        无符号整型数组，输出，读取历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities 的长度与 count 一致，
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_cross_section_values64(
    rtdb_int32 handle,
    const rtdb_int32* ids,
    rtdb_int32 mode,
    rtdb_int32 count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities,
    rtdb_error* errors
);

/**
* 命名：rtdbh_get_archived_values_filt
* 功能：读取单个标签点在一段时间内经复杂条件筛选后的历史储存值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param filter        字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
*                        长度不得超过 RTDB_EQUATION_SIZE，为 0 则不进行条件筛选。
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、values、states、qualities 的长度，
*                        即需要的数值个数；输出时返回实际得到的数值个数。
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史存储值；否则为 0
* \param states        64 位整数数组，输出，整型历史数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史存储值；否则为 0
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_values_filt64(
    rtdb_int32 handle,
    rtdb_int32 id,
    const char* filter,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 读取单个标签点某个时刻之后经复杂条件筛选后一定数量的等间隔内插值替换的历史数值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param filter        字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
*                        长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
* \param interval      整型，输入，插值时间间隔，单位为纳秒
* \param count         整型，输入，表示 datetimes、ms、values、states、qualities 的长度，
*                        即需要的插值个数。
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
* \param states        64 位整数数组，输出，整型历史插值数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
* \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素用于表示起始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_interval_values_filt64(
    rtdb_int32 handle,
    rtdb_int32 id,
    const char* filter,
    rtdb_int64 interval,
    rtdb_int32 count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 获取单个标签点一段时间内经复杂条件筛选后的等间隔插值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param filter        字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
*                        长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、values、states、qualities 的长度，
*                        即需要的插值个数；输出时返回实际得到的插值个数
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
* \param states        64 位整数数组，输出，整型历史插值数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
* \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_interpo_values_filt64(
    rtdb_int32 handle,
    rtdb_int32 id,
    const char* filter,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 获取单个标签点一段时间内经复杂条件筛选后的统计值
*
* \param handle            连接句柄
* \param id                整型，输入，标签点标识
* \param filter            字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
*                            长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
* \param datetime1         整型，输入/输出，输入时表示起始时间秒数。
*                            如果为 0，表示从存档中最早时间的数据开始进行统计。
*                            输出时返回最大值的时间秒数。
* \param ms1               短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                            表示起始时间对应的纳秒，输出时表示最大值的时间纳秒数；否则忽略，返回值为 0
* \param datetime2         整型，输入/输出，输入时表示结束时间秒数。
*                            如果为 0，表示统计到存档中最近时间的数据为止。
*                            输出时返回最小值的时间秒数。
* \param ms2               短整型，如果 id 指定的标签点时间精度为纳秒，
*                            表示结束时间对应的纳秒，输出时表示最小值的时间纳秒数；否则忽略，返回值为 0
* \param max_value         双精度浮点型，输出，表示统计时间段内的最大数值。
* \param min_value         双精度浮点型，输出，表示统计时间段内的最小数值。
* \param total_value       双精度浮点型，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
* \param calc_avg          双精度浮点型，输出，表示统计时间段内的算术平均值。
* \param power_avg         双精度浮点型，输出，表示统计时间段内的加权平均值。
* \remark 由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*        如果输出的最大值或最小值的时间戳秒值为 0，
*        则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_summary_data_filt(
    rtdb_int32 handle,
    rtdb_int32 id,
    const char* filter,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    RTDB_SUMMARY_DATA* summary_data
);

/**
* 命名：rtdbh_summary_filt_in_batches
* 功能：分批获取单一标签点一段时间内经复杂条件筛选后的统计值
*
* \param handle            连接句柄
* \param id                整型，输入，标签点标识
* \param filter            字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
*                            长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
* \param count             整形，输入/输出，输入时表示 datatimes1、ms1、datatimes2、ms2、
*                            max_values、min_values、total_values、calc_avgs、power_avgs、errors 的长度，
*                            即分段的个数；输出时表示成功取得统计值的分段个数。
* \param interval          64 位整型，输入，分段时间间隔，单位为纳秒。
* \param datetimes1        整型数组，输入/输出，输入时第一个元素表示起始时间秒数。
*                            如果为 0，表示从存档中最早时间的数据开始进行统计。
*                            输出时返回各个分段对应的最大值的时间秒数。
* \param ms1               短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                            第一个元素表示起始时间对应的纳秒，
*                            输出时返回各个分段对应的最大值的时间纳秒数；否则忽略，返回值为 0
* \param datetimes2        整型数组，输入/输出，输入时第一个元素表示结束时间秒数。
*                            如果为 0，表示统计到存档中最近时间的数据为止。
*                            输出时返回各个分段对应的最小值的时间秒数。
* \param ms2               短整型数组，如果 id 指定的标签点时间精度为纳秒，
*                            第一个元素表示结束时间对应的纳秒，
*                            输出时返回各个分段对应的最小值的时间纳秒数；否则忽略，返回值为 0
* \param max_values        双精度浮点型数组，输出，表示统计时间段内的最大数值。
* \param min_values        双精度浮点型数组，输出，表示统计时间段内的最小数值。
* \param total_values      双精度浮点型数组，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
* \param calc_avgs         双精度浮点型数组，输出，表示统计时间段内的算术平均值。
* \param power_avgs        双精度浮点型数组，输出，表示统计时间段内的加权平均值。
* \param errors            无符号整型数组，输出，表示各个分段取得统计值的返回值。
* \remark 由 datetimes1[0]、ms1[0] 表示的时间可以大于 datetimes2[0]、ms2[0] 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*        如果输出的最大值或最小值的时间戳秒值为 0，
*        则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_summary_data_filt_in_batches(
    rtdb_int32 handle,
    rtdb_int32 id,
    const char* filter,
    rtdb_int32* count,
    rtdb_int64 interval,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    RTDB_SUMMARY_DATA* summary_datas,
    rtdb_error* errors
);

/**
*
* \brief 修改单个标签点某一时间的历史存储值.
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime      整型，输入，时间秒数
* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；否则忽略。
* \param value         双精度浮点数，输入，浮点型历史数值
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放新的历史值；否则忽略
* \param state         64 位整数，输入，整型历史数值，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放新的历史值；否则忽略
* \param quality       短整型，输入，新的历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_update_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    rtdb_float64 value,
    rtdb_int64 state,
    rtdb_int16 quality
);

/**
*
* \brief 修改单个标签点某一时间的历史存储值.
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime      整型，输入，时间秒数
* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；否则忽略。
* \param x             单精度浮点型，输入，新的横坐标历史数值
* \param y             单精度浮点型，输入，新的纵坐标历史数值
* \param quality       短整型，输入，新的历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口仅对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_update_coor_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    rtdb_float32 x,
    rtdb_float32 y,
    rtdb_int16 quality
);


/**
*
* \brief 删除单个标签点某个时间的历史存储值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime      整型，输入，时间秒数
* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；否则忽略。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_remove_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime
);

/**
*
* \brief 删除单个标签点一段时间内的历史存储值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
* \param count         整形，输出，表示删除的历史值个数
* \remark 由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_remove_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_int32* count
);

/**
*
* \brief 写入单个标签点在某一时间的历史数据。
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime      整型，输入，时间秒数
* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；否则忽略。
* \param value         双精度浮点数，输入，浮点型历史数值
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放历史值；否则忽略
* \param state         64 位整数，输入，整型历史数值，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放历史值；否则忽略
* \param quality       短整型，输入，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    rtdb_float64 value,
    rtdb_int64 state,
    rtdb_int16 quality
);

/**
*
* \brief 写入单个标签点在某一时间的坐标型历史数据。
*
* \param handle              连接句柄
* \param id            整型，输入，标签点标识
* \param datetime      整型，输入，时间秒数
* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；否则忽略。
* \param x             单精度浮点型，输入，横坐标历史数值
* \param y             单精度浮点型，输入，纵坐标历史数值
* \param quality       短整型，输入，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_coor_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    rtdb_float32 x,
    rtdb_float32 y,
    rtdb_int16 quality
);

/**
*
* \brief 写入单个二进制/字符串标签点在某一时间的历史数据
*
* \param handle    连接句柄
* \param id        整型，输入，标签点标识
* \param datetime  整型，输入，数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型，输入，历史数值时间，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param blob      字节型数组，输入，历史二进制/字符串数值
* \param len       短整型，输入，二进制/字符串数值长度，超过一个页大小数据将被截断。
* \param quality   短整型，输入，历史数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_blob_value64(
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
* \brief 写入批量标签点批量历史存储数据
*
* \param handle        连接句柄
* \param count         整型，输入/输出，
*                        输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                        即历史值个数；输出时返回实际写入的数值个数
* \param ids           整型数组，输入，标签点标识，同一个标签点标识可以出现多次，
*                        但它们的时间戳必需是递增的。
* \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示对应的历史数值时间纳秒；否则忽略。
* \param values        双精度浮点数数组，输入，浮点型历史数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，表示相应的历史存储值；否则忽略
* \param states        64 位整数数组，输入，整型历史数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，表示相应的历史存储值；否则忽略
* \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致，
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_values64(
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
* \brief 写入批量标签点批量坐标型历史存储数据
*
* \param handle        连接句柄
* \param count         整型，输入/输出，
*                        输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
*                        即历史值个数；输出时返回实际写入的数值个数
* \param ids           整型数组，输入，标签点标识
* \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示对应的历史数值时间纳秒；否则忽略。
* \param x             单精度浮点型数组，输入，浮点型横坐标历史数值列表
* \param y             单精度浮点型数组，输入，浮点型纵坐标历史数值列表
* \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致，
*        本接口仅对数据类型为 RTDB_COOR 的标签点有效。
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_coor_values64(
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
* \brief 写入单个datetime标签点在某一时间的历史数据
*
* \param handle    连接句柄
* \param id        整型，输入，标签点标识
* \param datetime  整型，输入，数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型，输入，历史数值时间，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param blob      字节型数组，输入，历史datetime数值
* \param len       短整型，输入，datetime数值长度，超过一个页大小数据将被截断。
* \param quality   短整型，输入，历史数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_datetime_value64(
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
* \brief 写入批量标签点批量字符串型历史存储数据
*
* \param handle        连接句柄
* \param count         整型，输入/输出，
*                        输入时表示 ids、datetimes、ms、lens、blobs、qualities、errors 的长度，
*                        即历史值个数；输出时返回实际写入的数值个数
* \param ids           整型数组，输入，标签点标识
* \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示对应的历史数值时间纳秒；否则忽略。
* \param blobs         字节型指针数组，输入，实时二进制/字符串数值
* \param lens          短整型数组，输入，二进制/字符串数值长度，
*                        表示对应的 blobs 指针指向的缓冲区长度，超过一个页大小数据将被截断。
* \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、lens、blobs、qualities、errors 的长度与 count 一致，
*        本接口仅对数据类型为 RTDB_STRING、RTDB_BLOB 的标签点有效。
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_blob_values64(
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
* \brief 将标签点未写满的补历史缓存页写入存档文件中。
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识。
* \param count         整型，输出，缓存页中数据个数。
* \remark 补历史缓存页写满后会自动写入存档文件中，不满的历史缓存页也会写入文件，
*      但会有一个时间延迟，在此期间此段数据可能查询不到，为了及时看到补历史的结果，
*      应在结束补历史后调用本接口。
*      count 参数可为空指针，对应的信息将不再返回。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_flush_archived_values(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *count
  );

/**
* 命名：rtdbh_get_single_named_type_value32
* 功能：读取单个自定义类型标签点某个时间的历史数据
* 参数：
*        [handle]        连接句柄
*        [id]            整型，输入，标签点标识
*        [mode]          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
*        [datetime]      整型，输入/输出，输入时表示时间秒数；
*                        输出时表示实际取得的历史数值对应的时间秒数。
*        [ms]            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
*                        否则忽略输入，输出时为 0。
*        [object]        void数组，输出，自定义类型标签点历史值
*        [length]        短整型，输入/输出，输入时表示 object 的长度，
*                        输出时表示实际获取的自定义类型标签点数据长度。
*        [quality]       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_named_type_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 mode,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    void* object,
    rtdb_length_type* length,
    rtdb_int16* quality
);

/**
* 命名：rtdbh_get_archived_named_type_values32
* 功能：连续读取自定义类型标签点的历史数据
* 参数：
*        [handle]        连接句柄
*        [id]            整型，输入，标签点标识
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
*        [datetime1]     整型，输入，表示开始时间秒数；
*        [ms1]           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
*        [datetime2]     整型，输入,表示结束时间秒数；
*        [ms2]           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
*        [length]        短整型数组，输入，输入时表示 objects 的长度，
*        [count]         整型，输入/输出，输入表示想要查询多少数据
*                        输出表示实际查到多少数据
*        [datetimes]     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
*        [ms]            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
*                        表示实际取得的历史数值时间纳秒数。
*        [objects]       void类型数组，输出，自定义类型标签点历史值
*        [qualities]     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_named_type_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_length_type length,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    void* const* objects,
    rtdb_int16* qualities
);

/**
* 命名：rtdbh_put_single_named_type_value32
* 功能：写入自定义类型标签点的单个历史事件
* 参数：
*        [handle]    连接句柄
*        [id]        整型，输入，标签点标识
*        [datetime]  整型，输入，数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型，输入，历史数值时间，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
*        [object]    void数组，输入，历史自定义类型标签点数值
*        [length]    短整型，输入，自定义类型标签点数值长度，超过一个页大小数据将被截断。
*        [quality]   短整型，输入，历史数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_named_type_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    const void* object,
    rtdb_length_type length,
    rtdb_int16 quality
);

/**
* 命名：rtdbh_put_archived_named_type_values32
* 功能：批量补写自定义类型标签点的历史事件
* 参数：
*        [handle]        连接句柄
*        [count]         整型，输入/输出，
*                        输入时表示 ids、datetimes、ms、lens、blobs、qualities、errors 的长度，
*                        即历史值个数；输出时返回实际写入的数值个数
*        [ids]           整型数组，输入，标签点标识
*        [datetimes]     整型数组，输入，表示对应的历史数值时间秒数。
*        [ms]            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示对应的历史数值时间纳秒；否则忽略。
*        [objects]       void类型指针数组，输入，自定义类型标签点数值
*        [lengths]       短整型数组，输入，自定义类型标签点数值长度，
*                        表示对应的 objects 指针指向的缓冲区长度，超过一个页大小数据将被截断。
*        [qualities]     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [errors]        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
* 备注：用户须保证 ids、datetimes、ms、lens、objects、qualities、errors 的长度与 count 一致，
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_named_type_values64(
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
	fmt.Println("233")
}
