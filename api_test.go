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
* \brief 获得指定服务器端文件的大小
*
* \param handle     连接句柄
* \param file       字符串，输入，文件名
* \param size       64 位整数，输出，文件大小
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_file_size(
  rtdb_int32 handle,
  const char *file,
  rtdb_int64 *size
  );

/**
*
* \brief 读取服务器端指定文件的内容
*
* \param handle       连接句柄
* \param file         字符串，输入，要读取内容的文件名
* \param content      字符数组，输出，文件内容
* \param pos          64 位整型，输入，读取文件的起始位置
* \param size         整型，输入/输出，
*                     输入时表示要读取文件内容的字节大小；
*                     输出时表示实际读取的字节数
* \remark 用户须保证分配给 content 的空间与 size 相符。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_read_file(
  rtdb_int32 handle,
  const char *file,
  char *content,
  rtdb_int64 pos,
  rtdb_int64 *size
  );


/**
*
* \brief 取得数据库允许的blob与str类型测点的最大长度
*
* \param handle       连接句柄
* \param len          整形，输出参数，代表数据库允许的blob、str类型测点的最大长度
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_max_blob_len(
  rtdb_int32 handle,
  rtdb_int32 *len
  );

/**
*
* \brief 取得质量码对应的定义
*
* \param handle       连接句柄
* \param count        质量码个数，输入参数，
* \param qualities    质量码，输入参数
* \param definitions  输出参数，0~255为RTDB质量码（参加rtdb.h文件），256~511为OPC质量码，大于511为用户自定义质量码
* \param lens         输出参数，每个定义对应的长度
* \remark OPC质量码把8位分3部分定义：XX XXXX XX，对应着：品质位域 分状态位域 限定位域
* 品质位域：00（Bad）01（Uncertain）10（N/A）11（Good）
* 分状态位域：不同品质位域对应各自的分状态位域
* 限定位域：00（Not limited）01（Low limited）10（high limited）11（Constant）
* 三个域之间用逗号隔开，输出到definitions参数中，前面有有RTDB，OPC或者USER标识，说明标签点类别
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_format_quality(
  rtdb_int32 handle,
  rtdb_int32 *count,
  rtdb_int16 *qualities,
  rtdb_byte **definitions,
  rtdb_int32 *lens
  );

/*
* 命名：rtdb_write_named_type_field_by_name32
* 功能：按名称填充自定义类型数值中字段的内容
* 参数：
*      [handle]       连接句柄
*      [type_name]    自定义类型的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数，
*      [field_name]   自定义类型中需要填充的字段的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数
*      [field_type]   field_name字段的类型，RTDB_TYPE所支持的基础类型，输入参数
*      [object]       自定义类型数值的缓冲区,输入/输出参数
*     [object_len]   object缓冲区的长度,输入参数
*     [field]      需要填充的字段数值的缓冲区,输入参数
*     [field_len]    自定义类型中字段数值的缓冲区中数据的长度,输入参数
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_write_named_type_field_by_name32(
  rtdb_int32 handle,
  const char* type_name,
  const char* field_name,
  rtdb_int32 field_type,
  void* object,
  rtdb_length_type object_len,
  const void* field,
  rtdb_length_type field_len
  );

/*
* 命名：rtdb_write_named_type_field_by_pos32
* 功能：按位置填充自定义类型数值中字段的内容
* 参数：
*      [handle]       连接句柄
*      [type_name]    自定义类型的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数，
*      [field_pos]    自定义类型中需要填充的字段的位置，指字段在所有字段中的位置，从0开始，输入参数
*      [field_type]   field_pos位置所在字段的类型，RTDB_TYPE所支持的基础类型，输入参数
*      [object]       自定义类型数值的缓冲区,输入/输出参数
*     [object_len]   object缓冲区的长度,输入参数
*     [field]      需要填充的字段数值的缓冲区,输入参数
*     [field_len]    自定义类型中字段数值的缓冲区中数据的长度,输入参数
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_write_named_type_field_by_pos32(
  rtdb_int32 handle,
  const char* type_name,
  rtdb_int32 field_pos,
  rtdb_int32 field_type,
  void* object,
  rtdb_length_type object_len,
  const void* field,
  rtdb_length_type field_len
  );

/*
* 命名：rtdb_read_named_type_field_by_name32
* 功能：按名称提取自定义类型数值中字段的内容
* 参数：
*      [handle]       连接句柄
*      [type_name]    自定义类型的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数，
*      [field_name]   自定义类型中需要提取的字段的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数
*      [field_type]   field_name字段的类型，RTDB_TYPE所支持的基础类型，输入参数
*      [object]       自定义类型数值的缓冲区,输入参数
*     [object_len]   object缓冲区的长度,输入参数
*     [field]      被读取的字段的数值的缓冲区,输入/输出参数
*     [field_len]    field字段数值缓冲区的长度,输入参数
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_read_named_type_field_by_name32(
  rtdb_int32 handle,
  const char* type_name,
  const char* field_name,
  rtdb_int32 field_type,
  const void* object,
  rtdb_length_type object_len,
  void* field,
  rtdb_length_type field_len
  );

/*
* 命名：rtdb_read_named_type_field_by_pos32
* 功能：按位置提取自定义类型数值中字段的内容
* 参数：
*      [handle]       连接句柄
*      [type_name]    自定义类型的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数，
*      [field_pos]    自定义类型中需要提取的字段的位置，指字段在所有字段中的位置，从0开始，输入参数
*      [field_type]   field_pos位置所在字段的类型，RTDB_TYPE所支持的基础类型，输入参数
*      [object]       自定义类型数值的缓冲区,输入参数
*     [object_len]   object缓冲区的长度,输入参数
*     [field]      被读取的字段的数值的缓冲区,输入/输出参数
*     [field_len]    field字段数值缓冲区的长度,输入参数
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_read_named_type_field_by_pos32(
  rtdb_int32 handle,
  const char* type_name,
  rtdb_int32 field_pos,
  rtdb_int32 field_type,
  const void* object,
  rtdb_length_type object_len,
  void* field,
  rtdb_length_type field_len
  );

/*
* 命名：rtdb_named_type_name_field_check
* 功能：检查自定义类型名称及字段命名是否符合规则；
* 规则：1. 只允许使用26个英文字母,数字0-9，下划线；
*       2. 必须以字母作为首字母；
*       3. 大小写不敏感。
* 参数：
*      [check_name]   需要检查的名称
*      [flag]         标志0--类型名称，其他 -- 字段名称，暂不启用
*/

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_named_type_name_field_check(
  const char* check_name,
  rtdb_byte flag GAPI_DEFAULT_VALUE(0)
  );

/**
*
* \brief 判断连接是否可用
* \param handle   连接句柄
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_judge_connect_status(rtdb_int32 handle,
    rtdb_int8* change_connection GAPI_DEFAULT_VALUE(0),
    char* current_ip_addr GAPI_DEFAULT_VALUE(0),
    rtdb_int32 size GAPI_DEFAULT_VALUE(0));

/**
* 命名：rtdb_format_ipaddr
* 功能：将整形IP转换为字符串形式的IP
* 参数：
*      [ip]        无符号整型，输入，整形的IP地址
*      [ip_addr]      字符串，输出，字符串IP地址缓冲区
*      [size]         整型，输入，ip_addr 参数的字节长度
* 备注：用户须保证分配给 ip_addr 的空间与 size 相符
*/
RTDBAPI
void
RTDBAPI_CALLRULE
rtdb_format_ipaddr(rtdb_uint32 ip, char* ip_addr, rtdb_int32 size);
`
	sp := strings.Split(code, "/**")
	for _, fn := range sp {
		ff := "/**\n" + fn
		fmt.Println(CFunc(ff))
	}
}
