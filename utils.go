package rtdb_api

import "C"
import (
	"fmt"
	"regexp"
	"strings"
)

// CFunc 把原始的C的函数转换成解析后的C函数
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

func sp(s string) []string {
	// 统一换行符（兼容 Windows / Unix）
	s = strings.ReplaceAll(s, "\r\n", "\n")

	// 正则：一个或多个“空白行”
	re := regexp.MustCompile(`\n\s*\n+`)

	parts := re.Split(strings.TrimSpace(s), -1)

	// 可选：过滤掉意外的空段
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			res = append(res, p)
		}
	}
	return res
}

// XXX 将 C 接口注释 + 函数定义，转换为 Go 的空函数定义
func XXX(s string) string {
	comment := extractComment(s)
	cFuncName := extractCFuncName(s)
	goFuncName := "Raw" + toGoFuncName(cFuncName)

	var b strings.Builder
	if comment != "" {
		b.WriteString(comment)
		b.WriteString("\n")
	}
	b.WriteString("func ")
	b.WriteString(goFuncName)
	b.WriteString("() {}\n")

	return b.String()
}

// 提取 /** ... */ 注释
func extractComment(s string) string {
	re := regexp.MustCompile(`(?s)/\*\*.*?\*/`)
	return re.FindString(s)
}

// 提取 C 函数名
func extractCFuncName(s string) string {
	re := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`)
	m := re.FindStringSubmatch(s)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

// rtdbh_update_value64_warp -> RawRtdbhUpdateValue64Warp
func toGoFuncName(cName string) string {
	parts := strings.Split(cName, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

func Param() {
	code := `
RTDB_PARAM_TABLE_FILE,                               // 标签点表文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_BASE_FILE,                                // 基本标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_SCAN_FILE,                                // 采集标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_CALC_FILE,                                // 计算标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_SNAP_FILE,                                // 标签点快照文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_LIC_FILE,                                 // 协议文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_HIS_FILE,                                 // 历史信息文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_LOG_DIR,                                  // 服务器端日志文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_USER_FILE,                                // 用户权限信息文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_SERVER_FILE,                              // 网络服务进程与其它进程交互所用的共享内存文件，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_EQAUTION_FILE,                            // 方程式服务进程与其它进程交互所用的共享内存文件，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_ARV_PAGES_FILE,                           // 历史数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_ARVEX_PAGES_FILE,                         // 补历史数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_ARVEX_PAGES_BLOB_FILE,                    // 补历史数据blob、str缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_AUTH_FILE,                                // 信任连接段信息文件路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_RECYCLED_BASE_FILE,                       // 可回收基本标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_RECYCLED_SCAN_FILE,                       // 可回收采集标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_RECYCLED_CALC_FILE,                       // 可回收计算标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_AUTO_BACKUP_PATH,                         // 自动备份目的地全路径，必须以“\”或“/”结束，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_SERVER_SENDER_IP,                         // 镜像发送地址，字符串最大长度为 RTDB_MAX_HOSTNAME_SIZE
RTDB_PARAM_BLACKLIST_FILE,                           // 连接黑名单文件路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_DB_VERSION,                               // 数据库版本
RTDB_PARAM_LIC_USER,                                 // 授权单位
RTDB_PARAM_LIC_TYPE,                                 // 授权方式
RTDB_PARAM_INDEX_DIR,                                // 索引文件存放目录，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_MIRROR_BUFFER_PATH,                       // 镜像缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_MIRROR_EX_BUFFER_PATH,                    // 补写镜像缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_EQAUTION_PATH_FILE,                       // 方程式长度超过规定长度时进行保存的文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_TAGS_FILE,                                // 标签点关键属性文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_RECYCLED_SNAP_FILE,                       // 可回收标签点快照事件文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_SWAP_PAGE_FILE,					         // 历史数据交换页文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_PAGE_ALLOCATOR_FILE,				         // 活动存档数据页分配器文件全路径，字符串最大长度为 RTDB_MAX_PATH, 该系统配置项2.1版数据库在使用，3.0数据库已去掉，但为了保证系统选项索引号, 与2.1一致，此处不能去掉。便于java sdk统一调用
RTDB_PARAM_NAMED_TYPE_FILE,					         // 自定义类型配置信息全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_STRBLOB_MIRROR_PATH,				         // BLOB/STRING镜像数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_STRBLOB_MIRROR_EX_PATH,			         // 补写BLOB/STRING镜像数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_BUFFER_DIR,						         // 临时数据缓存路径
RTDB_PARAM_POOL_CACHE_FLIE,					         // 曲线池索引文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_POOL_DATA_FILE_DIR,				         // 曲线池缓存文件目录，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_ARCHIVE_FILE_PATH,				         // 存档文件低速存储区路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_LIC_VERSION_TYPE,                         // 授权版本
RTDB_PARAM_AUTO_MOVE_PATH,                           // 自动移动目的地全路径，必须以“\”或“/”结束，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_REPLICATION_BUFFER_PATH,				     // 双活：数据同步缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_REPLICATION_EX_BUFFER_PATH,		   	     // 双活：数据同步补写数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_STRBLOB_REPLICATION_BUFFER_PATH,	         // 双活：数据同步BLOB/STRING数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_STRBLOB_REPLICATION_EX_BUFFER_PATH,	     // 双活：数据同步补写BLOB/STRING数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_REPLICATION_GROUP_IP,                     // 双活：同步组地址，字符串最大长度为 RTDB_MAX_HOSTNAME_SIZE
RTDB_PARAM_ARC_FILENAME_PREFIX_WHEN_USING_DATE,      // 是否归档文件使用日期作为文件名
RTDB_PARAM_HOT_ARCHIVE_FILE_PATH,                    // 存档文件高速存储区路径，字符串最大长度为 RTDB_MAX_PATH
RTDB_PARAM_LIC_TABLES_COUNT,                         // 协议中限定的标签点表数量
RTDB_PARAM_LIC_TAGS_COUNT,                           // 协议中限定的所有标签点数量
RTDB_PARAM_LIC_SCAN_COUNT,                           // 协议中限定的采集标签点数量
RTDB_PARAM_LIC_CALC_COUNT,                           // 协议中限定的计算标签点数量
RTDB_PARAM_LIC_ARCHICVE_COUNT,                       // 协议中限定的历史存档文件数量
RTDB_PARAM_SERVER_IPC_SIZE,                          // 网络服务进程与其它进程进行交互所使用的共享内存池的字节尺寸（单位：B）
RTDB_PARAM_EQUATION_IPC_SIZE,                        // 方程式服务进程与其它进程进行交互所使用的共享内存池的字节尺寸（单位：B）
RTDB_PARAM_HASH_TABLE_SIZE,                          // 标签点求余哈希表的尺寸
RTDB_PARAM_TAG_DELETE_TIMES,          		   	     // 可整库删除标签点的次数
RTDB_PARAM_SERVER_PORT,                              // 网络服务独立服务器端口
RTDB_PARAM_SERVER_SENDER_PORT,                       // 网络服务镜像发送端口
RTDB_PARAM_SERVER_RECEIVER_PORT,                     // 网络服务镜像接收端口
RTDB_PARAM_SERVER_MODE,                              // 网络服务启动模式
RTDB_PARAM_SERVER_CONNECTION_COUNT,                  // 协议中限定网络服务连接并发数量
RTDB_PARAM_ARV_PAGES_NUMBER,                         // 历史数据缓存中的页数量
RTDB_PARAM_ARVEX_PAGES_NUMBER,                       // 补历史数据缓存中的页数量
RTDB_PARAM_EXCEPTION_AT_SERVER,                      // 是否由服务器进行例外判定
RTDB_PARAM_ARV_PAGE_RECYCLE_DELAY,                   // 历史数据缓存页回收延时（毫秒）
RTDB_PARAM_EX_ARCHIVE_SIZE,                          // 历史数据存档文件文件自动增长大小（单位：MB）
RTDB_PARAM_ARCHIVE_BATCH_SIZE,                       // 历史存储值分段查询个数
RTDB_PARAM_DATAFILE_PAGESIZE,                        // 系统数据文件页大小
RTDB_PARAM_ARV_ASYNC_QUEUE_NORMAL_DOOR,              // 历史数据缓存队列中速归档区（单位：百分比）
RTDB_PARAM_INDEX_ALWAYS_IN_MEMORY,                   // 常驻内存的历史数据索引大小（单位：MB）
RTDB_PARAM_DISK_MIN_REST_SIZE,                       // 最低可用磁盘空间（单位：MB）
RTDB_PARAM_MIN_SIZE_OF_ARCHIVE,                      // 历史存档文件和附属文件的最小尺寸（单位：MB）
RTDB_PARAM_DELAY_OF_AUTO_MERGE_OR_ARRANGE,           // 自动合并/整理最小延迟（单位：小时）
RTDB_PARAM_START_OF_AUTO_MERGE_OR_ARRANGE,           // 自动合并/整理开始时间（单位：点钟）
RTDB_PARAM_STOP_OF_AUTO_MERGE_OR_ARRANGE,            // 自动合并/整理停止时间（单位：点钟）
RTDB_PARAM_START_OF_AUTO_BACKUP,                     // 自动备份开始时间（单位：点钟）
RTDB_PARAM_STOP_OF_AUTO_BACKUP,                      // 自动备份停止时间（单位：点钟）
RTDB_PARAM_MAX_LATENCY_OF_SNAPSHOT,                  // 允许服务器时间之后多少小时内的数据进入快照（单位：小时）
RTDB_PARAM_PAGE_ALLOCATOR_RESERVE_SIZE,              // 活动页分配器预留大小（单位：KB）， 0 表示使用操作系统视图大小
RTDB_PARAM_INCLUDE_SNAPSHOT_IN_QUERY,                // 决定取样本值和统计值时，快照是否应该出现在查询结果中
RTDB_PARAM_LIC_BLOB_COUNT,                           // 协议中限定的字符串或BLOB类型标签点数量
RTDB_PARAM_MIRROR_BUFFER_SIZE,                       // 镜像文件大小（单位：GB）
RTDB_PARAM_BLOB_ARVEX_PAGES_NUMBER,                  // blob、str补历史的默认缓存页数量
RTDB_PARAM_MIRROR_EVENT_QUEUE_CAPACITY,              // 镜像缓存队列容量
RTDB_PARAM_NOTIFY_NOT_ENOUGH_SPACE,                  // 提示磁盘空间不足，一旦启用，设置为ON，则通过API返回大错误码，否则只记录日志
RTDB_PARAM_ARCHIVE_FIXED_RANGE,                      // 历史数据存档文件的固定时间范围，默认为0表示不使用固定时间范围（单位：分钟）
RTDB_PARAM_ONE_CLINET_MAX_CONNECTION_COUNT,          // 单个客户端允许的最大连接数，默认为0表示不限制
RTDB_PARAM_ARV_PAGES_CAPACITY,                       // 历史数据缓存所占字节大小，单位：字节
RTDB_PARAM_ARVEX_PAGES_CAPACITY,                     // 历史数据补写缓存所占字节大小，单位：字节
RTDB_PARAM_BLOB_ARVEX_PAGES_CAPACITY,                // blob、string类型标签点历史数据补写缓存所占字节大小，单位：字节
RTDB_PARAM_LOCKED_PAGES_MEM,                         // 指定分配给数据库用的内存大小，单位：MB
RTDB_PARAM_LIC_RECYCLE_COUNT,                        // 协议中回收站的容量
RTDB_PARAM_ARCHIVED_POLICY,                          // 快照数据和补写数据的归档策略
RTDB_PARAM_NETWORK_ISOLATION_ACK_BYTE,               // 网络隔离装置ACK字节
RTDB_PARAM_ENABLE_LOGGER,                            // 启用日志输出，0为不启用
RTDB_PARAM_LOG_ENCODE,                               // 启用日志加密，0为不启用
RTDB_PARAM_LOGIN_TRY,                                // 启用登录失败次数验证，0为不启用
RTDB_PARAM_USER_LOG,                                 // 启用用户详细日志，0为不启用
RTDB_PARAM_COVER_WRITE_LOG,                          // 启用日志覆盖写功能，0为不启用
RTDB_PARAM_LIC_NAMED_TYPE_COUNT,				     // 协议中限定的自定义类型标签点数量
RTDB_PARAM_MIRROR_RECEIVER_THREADPOOL_SIZE,		     // 镜像接收线程数量
RTDB_PARAM_SNAPSHOT_USE_ARCHIVE_INTERFACE,		     // 按照补历史流程归档快照数据页
RTDB_PARAM_NO_ARCDATA_WRITE_LOG,				     // 归档无对应存档文件的数据时记录日志
RTDB_PARAM_PUT_ARCHIVE_THREAD_NUM,				     // 补历史归档线程数
RTDB_PARAM_ARVEX_DATA_ARCHIVED_THRESHOLD,		     // 单次补写数据归档阈值
RTDB_PARAM_SNAPSHOT_FLUSH_BUFFER_DELAY,		   	     // 快照服务的共享缓存刷新到磁盘的周期
RTDB_PARAM_DATA_SPEED,						   	     // 查询时使用加速统计
RTDB_PARAM_USE_NEW_PLOT_ALGO,				   	     // 启用新的曲线算法
RTDB_PARAM_QUERY_THREAD_POOL_SIZE,				     // 曲线查询线程池中线程数量
RTDB_PARAM_ARCHIVED_VALUES,                          // 使用查询线程池查询历史数据
RTDB_PARAM_ARCHIVED_VALUES_COUNT,                    // 使用查询线程池查询历史数据的条数
RTDB_PARAM_POOL_USE_FLAG,					   	     // 启用曲线池
RTDB_PARAM_POOL_OUT_LOG_FLAG,				   	     // 输出曲线池日志
RTDB_PARAM_POOL_TIME_USE_POOL_FLAG,				     // 使用曲线池缓存计算插值
RTDB_PARAM_POOL_MAX_POINT_COUNT,				     // 曲线池的标签点容量
RTDB_PARAM_POOL_ONE_FILE_SAVE_POINT_COUNT,		     // 曲线池每个数据文件的标签点容量
RTDB_PARAM_POOL_SAVE_MEMORY_SIZE,				     // 曲线缓存退出时临时缓冲区大小
RTDB_PARAM_POOL_MIN_TIME_UNIT_SECONDS,		   	     // 曲线池缓存数据当前时间单位
RTDB_PARAM_POOL_TIME_UNIT_VIEW_RATE,			     // 曲线池查询数据最小时间单位显示系数
RTDB_PARAM_POOL_TIMER_INTERVAL_SECONDS,		   	     // 曲线池定时器刷新周期
RTDB_PARAM_POOL_PERF_TIMER_INTERVAL_SECONDS,	     // 曲线池性能计算点刷新周期
RTDB_PARAM_ARCHIVE_INIT_FILE_SIZE,				     // 存档文件初始大小
RTDB_PARAM_ARCHIVE_INCREASE_MODE,				     // 存档文件增长模式
RTDB_PARAM_ARCHIVE_INCREASE_SIZE,				     // 固定模式下文件增长大小
RTDB_PARAM_ARCHIVE_INCREASE_PERCENT,			     // 百分比模式下增长百分比
RTDB_PARAM_ALLOW_CONVERT_SKL_TO_RBT_INDEX,	         // 跳跃链表转换到红黑树
RTDB_PARAM_EARLY_DATA_TIME,					         // 冷数据时间
RTDB_PARAM_EARLY_INDEX_TIME,				   	     // 自动转换索引时间
RTDB_PARAM_ARRANGE_RBT_TIME,				   	     // 整理存档文件时决定索引格式的时间轴
RTDB_PARAM_ENABLE_BIG_DATA,					   	     // 将存档文件全部读取到内存中
RTDB_PARAM_AUTO_ARRANGE_PERCENT,				     // 自动整理存档文件时的实际使用率
RTDB_PARAM_EARLY_ARRANGE_TIME,					     // 自动整理存档文件的时间
RTDB_PARAM_MIN_AUTO_ARRANGE_ARCFILE_PERCENT,	     // 自动整理存档文件时的最小使用率
RTDB_PARAM_ARRANGE_ARC_WITH_MEMORY,			         // 在内存中整理存档文件
RTDB_PARAM_ARAANGE_ARC_MAX_MEM_PERCENT,		         // 整理存档文件最大内存使用率
RTDB_PARAM_MAX_DISK_SPACE_PERCENT,			         // 磁盘最大使用率
RTDB_PARAM_USE_DISPATH,						         // windows 用 linux 已禁用,是否启用转发服务
RTDB_PARAM_USE_SMART_PARAM,					         // windows 用 linux 已禁用,是否使用推荐参数
RTDB_PARAM_SUBSCRIBE_SNAPSHOT_COUNT,                 // 单连接快照事件订阅个数
RTDB_PARAM_SUBSCRIBE_QUEUE_SIZE,                     // 订阅事件队列大小
RTDB_PARAM_SUBSCRIBE_TIMEOUT,                        // 订阅事件超时时间
RTDB_PARAM_MIRROR_COMPRESS_ONOFF,				     // 镜像报文压缩是否打开
RTDB_PARAM_MIRROR_COMPRESS_TYPE,				     // 镜像报文压缩类型
RTDB_PARAM_MIRROR_COMPRESS_MIN,				    	 // 镜像报文压缩最小值
RTDB_PARAM_ARCHIVE_ROLL_TIME,				    	 // 存档文件滚动时间轴
RTDB_PARAM_HANDLE_TIME_OUT,					    	 // 连接超时断开，单位：秒
RTDB_PARAM_MOVE_ARV_TIME,					         // 移动存档文件时决定移动存档的时间轴
RTDB_PARAM_USE_NEW_INTERP_ALGO,				    	 // 启用新的插值算法
RTDB_PARAM_ENABLE_REPLICATION,                       // 启用双活，0为不启用，1为启用
RTDB_PARAM_REPLICATION_GROUP_PORT,                   // 双活：同步组端口
RTDB_PARAM_REPLICATION_THREAD_SIZE,				     // 双活：同步线程数
RTDB_PARAM_FORCE_ARCHIVE_INCOMPLETE_DATA_PAGE_DELAY, // 强制归档补历史缓存里面未满数据页的延迟时间
RTDB_PARAM_ARCHIVE_ROLL_DISK_PERCENTAGE,             // 存档文件滚动存储空间百分比
RTDB_PARAM_ENABLE_IPV6,                              // 启用ipv6设置
RTDB_PARAM_ENABLE_USE_ARCHIVED_VALUE,                // 按条件获取历史值时，是否直接获取条件中点的历史值，0:获取插值，1:获取历史值
RTDB_PARAM_TIMESTAMP_TYPE,                           // 获取服务器时间戳类型
RTDB_PARAM_ARC_FILENAME_USING_DATE,				     // 是否归档文件使用日期作为文件名
RTDB_PARAM_LOG_MAX_SPACE,				             // 日志文件占用的最大磁盘空间
RTDB_PARAM_LOG_FILE_SIZE,					   	     // 单个日志文件大小
RTDB_PARAM_IGNORE_TO_WRITE_NOARCBUFFER,              // 是否丢弃补历史数据
RTDB_PARAM_ARCHIVES_COUNT_FOR_CALC,                  // 统计存档文件平均大小的存档文件个数
RTDB_PARAM_MAX_BLOB_SIZE,                            // blob、str类型数据在数据库中允许的最大长度
`
	sp := strings.Split(code, "\n")
	for _, line := range sp {
		if strings.TrimSpace(line) == "" {
			continue
		}

		spp := strings.Split(line, ",")
		name := strings.TrimSpace(spp[0])
		spp2 := strings.Split(line, "//")
		desc := strings.TrimSpace(spp2[1])

		nsp := strings.Split(name, "_")
		nsp2 := make([]string, 0)
		for _, n := range nsp {
			nsp2 = append(nsp2, Capitalize(n))
		}
		name2 := strings.Join(nsp2, "")

		fmt.Printf("case %s:\n", name2)
		fmt.Printf("    return \"%s\"\n", desc)
		// fmt.Printf("// %s %s\n", name2, desc)
		// fmt.Printf("%s = RtdbParam(C.%s)\n\n", name2, name)
	}
}

// Capitalize returns s with the first letter uppercased
// and the remaining letters lowercased.
func Capitalize(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)

	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	for i := 1; i < len(runes); i++ {
		runes[i] = []rune(strings.ToLower(string(runes[i])))[0]
	}

	return string(runes)
}
