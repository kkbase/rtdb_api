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
* \brief 重算或补算批量计算标签点历史数据
*
* \param handle        连接句柄
* \param count         整型，输入/输出，
*                        输入时表示 ids、errors 的长度，
*                        即标签点个数；输出时返回成功开始计算的标签点个数
* \param flag          短整型，输入，不为 0 表示进行重算，删除时间范围内已经存在历史数据；
*                        为 0 表示补算，保留时间范围内已经存在历史数据，覆盖同时刻的计算值。
* \param datetime1     整型，输入，表示起始时间秒数。
* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示计算直至存档中数据的最后时间
* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
* \param ids           整型数组，输入，标签点标识
* \param errors        无符号整型数组，输出，计算历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、errors 的长度与 count 一致，本接口仅对带有计算扩展属性的标签点有效。
*        由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbe_compute_history64(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int16 flag,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    const rtdb_int32* ids,
    rtdb_error* errors
);

/**
* 命名：rtdbe_get_equation_graph_count
* 功能：根据标签点 id 获取相关联方程式键值对数量
* 参数：
*      [handle]   连接句柄
*      [id]       整型，输入，标签点标识
*      [flag]     枚举，输入，获取的拓扑图的关系
*      [count]    整型，输入，拓扑图键值对数量
* 备注：键值对为数据结构，存储方程式涉及到的各标签点ID、及其父ID等
*		具体参考rtdbe_get_equation_graph_datas
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbe_get_equation_graph_count(
  rtdb_int32 handle,
  rtdb_int32 id,
  RTDB_GRAPH_FLAG flag,
  rtdb_int32 *count
  );

/**
* 命名：rtdbe_get_equation_graph_datas
* 功能：根据标签点 id 获取相关联方程式键值对数据
* 参数：
*      [handle]   连接句柄
*      [id]       整型，输入，标签点标识
*      [flag]     枚举，输入，获取的拓扑图的关系
*      [count]    整型，输出
                    输入时，表示拓扑图键值对数量
                    输出时，表示实际获取到的拓扑图键值对数量
*      [graph]    输出，GOLDE_GRAPH数据结构，拓扑图键值对信息
* 备注：键值对为数据结构，存储方程式涉及到的各标签点ID、及其父ID等
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbe_get_equation_graph_datas(
  rtdb_int32 handle,
  rtdb_int32 id,
  RTDB_GRAPH_FLAG flag,
  rtdb_int32 *count,
  RTDB_GRAPH *graph
  );


/**
* 命名：rtdbp_get_perf_tags_count
* 功能：获取Perf服务中支持的性能计数点的数量
* 参数：
*      [handle]   连接句柄
*      [count]    整型，输出，表示实际获取到的Perf服务中支持的性能计数点的数量
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbp_get_perf_tags_count(
  rtdb_int32 handle,
  int* count);


/**
* 命名：rtdbp_get_perf_tags_info
* 功能：根据性能计数点ID获取相关的性能计数点信息
* 参数：
*      [handle]   连接句柄
*      [count]    整型，输入，输出
                    输入时，表示想要获取的性能计数点信息的数量，也表示tags_info，errors等的长度
                    输出时，表示实际获取到的性能计数点信息的数量
       [errors] 无符号整型数组，输出，获取性能计数点信息的返回值列表，参考rtdb_error.h
* 备注：用户须保证分配给 tags_info，errors 的空间与 count 相符
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbp_get_perf_tags_info(
  rtdb_int32 handle,
  rtdb_int32* count,
  RTDB_PERF_TAG_INFO* tags_info,
  rtdb_error* errors);


/**
* 命名：rtdbp_get_perf_values
* 功能：批量读取性能计数点的当前快照数值
* 参数：
*        [handle]    连接句柄
*        [count]     整型，输入/输出，性能点个数，
*                    输入时表示 perf_ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功获取实时值的性能计数点个数
*        [perf_ids]  整型数组，输入，性能计数点标识列表，参考RTDB_PERF_TAG_ID
*        [datetimes] 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
*        [values]    双精度浮点型数组，输出，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的快照值；否则为 0
*        [states]    64 位整型数组，输出，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的快照值；否则为 0
*        [qualities] 短整型数组，输出，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [errors]    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* 备注：用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbp_get_perf_values64(
    rtdb_int32 handle,
    rtdb_int32* count,
    int* perf_ids,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities,
    rtdb_error* errors);
`
	sp := strings.Split(code, "/**")
	for _, fn := range sp {
		ff := "/**\n" + fn
		fmt.Println(CFunc(ff))
	}
}

func TestHello(t *testing.T) {
	fmt.Println("243")
}
