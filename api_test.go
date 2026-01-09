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
* \brief 获取存档文件数量
*
* \param handle    连接句柄
* \param count     整型，输出，存档文件数量
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_archives_count(
  rtdb_int32 handle,
  rtdb_int32 *count
  );

/**
*
* \brief 新建指定时间范围的历史存档文件并插入到历史数据库
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名。
* \param begin      整数，输入，起始时间，距离1970年1月1日08:00:00的秒数
* \param end        整数，输入，终止时间，距离1970年1月1日08:00:00的秒数
* \param mb_size    整型，输入，文件兆字节大小，单位为 MB。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_create_ranged_archive64(
    rtdb_int32 handle,
    const char* path,
    const char* file,
    rtdb_timestamp_type begin,
    rtdb_timestamp_type end,
    rtdb_int32 mb_size
);


/**
*
* \brief 追加磁盘上的历史存档文件到历史数据库。
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名，后缀名应为.rdf。
* \param state      整型，输入，取值 RTDB_ACTIVED_ARCHIVE、RTDB_NORMAL_ARCHIVE、
*                     RTDB_READONLY_ARCHIVE 之一，表示文件状态
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_append_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file,
  rtdb_int32 state
  );

/**
*
* \brief 从历史数据库中移出历史存档文件。
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_remove_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file
  );


/**
*
* \brief 切换活动文件
*
* \param handle     连接句柄
* \remark 当前活动文件被写满时该事务被启动，
*        改变当前活动文件的状态为普通状态，
*        在所有历史数据存档文件中寻找未被使用过的
*        插入到前活动文件的右侧并改为活动状态，
*        若找不到则将前活动文件右侧的文件改为活动状态，
*        并将active_archive_指向该文件。该事务进行过程中，
*        用锁保证所有读写操作都暂停等待该事务完成。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_shift_actived(
  rtdb_int32 handle
  );


/**
* 命名：rtdba_get_archives
* 功能：获取存档文件的路径、名称、状态和最早允许写入时间。
* 参数：
*        [handle]          连接句柄
*        [paths]            字符串数组，输出，存档文件的目录路径，长度至少为 RTDB_PATH_SIZE。
*        [files]            字符串数组，输出，存档文件的名称，长度至少为 RTDB_FILE_NAME_SIZE。
*        [states]           整型数组，输出，取值 RTDB_INVALID_ARCHIVE、RTDB_ACTIVED_ARCHIVE、
*                          RTDB_NORMAL_ARCHIVE、RTDB_READONLY_ARCHIVE 之一，表示文件状态
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_archives(
  rtdb_int32 handle,
  rtdb_int32* count,
  rtdb_path_string* paths,
  rtdb_filename_string* files,
  rtdb_int32 *states
  );

/**
* 功能：获取存档信息
* 参数：
*    [handle]: in, 句柄
*    [count]: out, 数量
*    [paths]: out, 路径
*    [files]: out, 文件
*    [infos]: out, 存档信息
*    [errors]: out, 错误
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_archives_info(
  rtdb_int32 handle,
  rtdb_int32* count,
  const rtdb_path_string* const paths,
  const rtdb_filename_string* const files,
  RTDB_HEADER_PAGE *infos,
  rtdb_error* errors
  );

/**
* 功能：获取存档的实时信息
* 参数：
*    [handle]: in, 句柄
*    [count]: out, 数量
*    [paths]: out, 路径
*    [files]: out, 文件
*    [real_time_datas]: out, 实时数据
*    [total_datas]: 总数
*    [errors]: 错误
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_archives_perf_data(
  rtdb_int32 handle,
  rtdb_int32* count,
  const rtdb_path_string* const paths,
  const rtdb_filename_string* const files,
  RTDB_ARCHIVE_PERF_DATA* real_time_datas,
  RTDB_ARCHIVE_PERF_DATA* total_datas,
  rtdb_error* errors
  );

/**
* 功能：获取存档状态
* 参数：
*    [handle]: in, 句柄
*    [status]: out, 存档状态
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_archives_status(rtdb_int32 handle, rtdb_error* status);

/**
*
* \brief 获取存档文件及其附属文件的详细信息。
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名。
* \param file_id    整型，输入，附属文件标识，0 表示获取主文件信息。
* \param info       RTDB_HEADER_PAGE 结构，输出，存档文件信息
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_archive_info(
  rtdb_int32 handle,
  const char *path,
  const char *file,
  rtdb_int32 file_id,
  RTDB_HEADER_PAGE *info
  );


/**
*
* \brief 修改存档文件的可配置项。
*
* \param handle         连接句柄
* \param path           字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file           字符串，输入，文件名。
* \param rated_capacity 整型，输入，文件额定大小，单位为 MB。
* \param ex_capacity    整型，输入，附属文件大小，单位为 MB。
* \param auto_merge     短整型，输入，是否自动合并附属文件。
* \param auto_arrange   短整型，输入，是否自动整理存档文件。
* 备注: rated_capacity 与 ex_capacity 参数可为 0，表示不修改对应的配置项。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_update_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file,
  rtdb_int32 rated_capacity,
  rtdb_int32 ex_capacity,
  rtdb_int16 auto_merge,
  rtdb_int16 auto_arrange
  );

/**
*
* \brief 整理存档文件，将同一标签点的数据块存放在一起以提高查询效率。
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_arrange_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file
  );

/**
*
* \brief 为存档文件重新生成索引，用于恢复数据。
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_reindex_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file
  );

/**
*
* \brief 备份主存档文件及其附属文件到指定路径
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名。
* \param dest       字符串，输入，备份目录路径，必须以"\"或"/"结尾。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_backup_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file,
  const char *dest
  );

/**
* 命名：rtdba_move_archive
* 功能：将存档文件移动到指定目录
* 参数：
*        [handle]     连接句柄
*        [path]       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
*        [file]       字符串，输入，文件名。
*        [dest]       字符串，输入，移动目录路径，必须以"\"或"/"结尾。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_move_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file,
  const char *dest
  );

/**
* 命名：rtdba_reindex_archive
* 功能：为存档文件转换索引格式。
* 参数：
*        [handle]     连接句柄
*        [path]       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
*        [file]       字符串，输入，文件名。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_convert_index(
  rtdb_int32 handle,
  const char *path,
  const char *file
  );


/**
* 命名：rtdba_query_big_job
* \brief 查询进程正在执行的后台任务类型、状态和进度
*
* \param handle     连接句柄
* \param process    所查询的进程代号，进程的标识参见枚举 RTDB_PROCESS_NAME,
*                     RTDB_PROCESS_HISTORIAN: 历史服务进程，具有以下任务类型：
*                         RTDB_MERGE: 合并附属文件到主文件;
*                         RTDB_ARRANGE: 整理存档文件;
*                         RTDB_REINDEX: 重建索引;
*                         RTDB_BACKUP: 备份;
*                         RTDB_REACTIVE: 激活为活动存档;
*                     RTDB_PROCESS_EQUATION: 方程式服务进程，具有以下任务类型：
*                         RTDB_COMPUTE: 历史计算;
*                     RTDB_PROCESS_BASE: 标签信息服务进程，具有以下任务类型：
*                         RTDB_UPDATE_TABLE: 修改表名称;
*                         RTDB_REMOVE_TABLE: 删除表;
* \param path       字符串，输出，长度至少为 RTDB_PATH_SIZE，
*                     对以下任务，这个字段表示存档文件所在目录路径：
*                         RTDB_MERGE
*                         RTDB_ARRANGE
*                         RTDB_REINDEX
*                         RTDB_BACKUP
*                         RTDB_REACTIVE
*                     对于以下任务，这个字段表示原来的表名：
*                         RTDB_UPDATE_TABLE
*                         RTDB_REMOVE_TABLE
*                     对于其它任务不可用。
* \param file       字符串，输出，长度至少为 RTDB_FILE_NAME_SIZE，
*                     对以下任务，这个字段表示存档文件名：
*                         RTDB_MERGE
*                         RTDB_ARRANGE
*                         RTDB_REINDEX
*                         RTDB_BACKUP
*                         RTDB_REACTIVE
*                     对于以下任务，这个字段表示修改后的表名：
*                          RTDB_UPDATE_TABLE
*                     对于其它任务不可用。
* \param job        短整型，输出，任务的标识参见枚举 RTDB_BIG_JOB_NAME。
* \param state      整型，输出，任务的执行状态，参考 rtdb_error.h
* \param end_time   整型，输出，任务的完成时间。
* \param progress   单精度浮点型，输出，任务的进度百分比。
* \remark path 及 file 参数可传空指针，对应的信息将不再返回。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_query_big_job64(
    rtdb_int32 handle,
    rtdb_int32 process,
    char* path,
    char* file,
    rtdb_int16* job,
    rtdb_int32* state,
    rtdb_timestamp_type* end_time,
    rtdb_float32* progress
);

/**
* 命名：rtdba_cancel_big_job
* 功能：取消进程正在执行的后台任务
* 参数：
*        [handle]     连接句柄
*        [process]    所查询的进程代号，进程的标识参见枚举 RTDB_PROCESS_NAME,
*                     RTDB_PROCESS_HISTORIAN: 历史服务进程，具有以下任务类型：
*                         RTDB_MERGE: 合并附属文件到主文件;
*                         RTDB_ARRANGE: 整理存档文件;
*                         RTDB_REINDEX: 重建索引;
*                         RTDB_BACKUP: 备份;
*                         RTDB_REACTIVE: 激活为活动存档;
*                     RTDB_PROCESS_EQUATION: 方程式服务进程，具有以下任务类型：
*                         RTDB_COMPUTE: 历史计算;
*                     RTDB_PROCESS_BASE: 标签信息服务进程，具有以下任务类型：
*                         RTDB_UPDATE_TABLE: 修改表名称;
*                         RTDB_REMOVE_TABLE: 删除表;
* 备注：path 及 file 参数可传空指针，对应的信息将不再返回。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_cancel_big_job(rtdb_int32 handle, rtdb_int32 process);
`
	sp := strings.Split(code, "/**")
	for _, fn := range sp {
		ff := "/**\n" + fn
		fmt.Println(CFunc(ff))
	}
}

func TestHello(t *testing.T) {
	fmt.Println("123")
}
