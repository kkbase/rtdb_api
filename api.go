package rtdb_api

// #cgo CFLAGS: -DPNG_DEBUG=1 -I./cinclude
// #cgo CXXFLAGS: -std=c++11
// #include <stdlib.h>
// #include "api.h"
import "C"
import (
	_ "embed"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"
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
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.load_library(cPath)
}

// RtdbError 数据库错误
type RtdbError uint32

func (re RtdbError) IsOk() bool {
	return errors.Is(re, RteOk)
}

func (re RtdbError) Error() string {
	desc := ""
	switch re {
	case RteUnknownError:
		desc = "未知错误"
	case RteOk:
		desc = "操作成功完成"
	case RteWindowsError:
		desc = "Windows操作系统错误的起始值"
	case RteWindowsErrorMax:
		desc = "Windows操作系统错误的结束值"
	case RteInvalidOpenmode:
		desc = "无效的文件打开方式"
	case RteOpenfileFailed:
		desc = "打开文件失败"
	case RteMovetoendFailed:
		desc = "移动文件指针到文件尾失败"
	case RteDifferReadbytes:
		desc = "读取内容长度与要求不符"
	case RteGetfileposFailed:
		desc = "获取当前文件指针失败"
	case RteFlushfileFailed:
		desc = "刷新文件缓冲区失败"
	case RteSetsizeFailed:
		desc = "设置文件大小失败"
	case RteFileNotClosed:
		desc = "试图用未关闭的文件对象创建或打开文件"
	case RteFileUnknown:
		desc = "创建或打开文件时必须指定文件名"
	case RteInvalidHeader:
		desc = "数据文件头信息错误"
	case RteDisabledFile:
		desc = "数据文件无效，试图访问无效数据文件"
	case RteFileNotOpened:
		desc = "试图访问尚未打开的数据文件"
	case RtePointNotFound:
		desc = "要求访问的标签点不存在或无效"
	case RteReadyblockNotFound:
		desc = "数据文件中找不到从指定数据块以后的可用的空块"
	case RteFileIsIncult:
		desc = "文件未被使用过"
	case RteFileIsFull:
		desc = "数据文件已满"
	case RteFileexIsFull:
		desc = "数据文件扩展区已满，无法继续装载数据"
	case RteInvalidDataType:
		desc = "无效的数据类型"
	case RteDatablockNotFound:
		desc = "找不到符合时间条件的数据块"
	case RteDataBetweenBlock:
		desc = "数据时间位于找到的块和下一个数据块之间"
	case RteCantModifyExistValue:
		desc = "不允许修改已存在的数据"
	case RteWrongdataInBlock:
		desc = "块中有错误数据导致数据块头信息不符"
	case RteDatatimeNotIn:
		desc = "数据文件中没有该标签点指定时间的数据"
	case RteNullArchivePath:
		desc = "操作的数据文件路径为空"
	case RteRegArchivePath:
		desc = "数据文件已被注册"
	case RteUnregArchivePath:
		desc = "未注册的数据文件"
	case RteFileInexistence:
		desc = "指定的文件不存在"
	case RteDataTypeNotMatch:
		desc = "数据类型不匹配"
	case RteFileIsReadonly:
		desc = "不允许修改只读数据文件中的数据"
	case RteTomanyArchiveFile:
		desc = "过多的数据文件"
	case RteNoPointsList:
		desc = "缺少标签点列表"
	case RteNoActivedArchive:
		desc = "缺少活动文档"
	case RteNoArchiveFile:
		desc = "缺少数据文档"
	case RteNeedActivedArchive:
		desc = "只能在活动文档上执行该操作"
	case RteInvalidTimestamp:
		desc = "无效的时间戳"
	case RteNeedMoreWritable:
		desc = "非只读文档个数太少"
	case RteNoArchiveForPut:
		desc = "找不到合适的追加历史数据的文档"
	case RteInvalidValueMode:
		desc = "无效的取值模式"
	case RteDataNotFound:
		desc = "找不到需要的数据"
	case RteInvalidParameter:
		desc = "无效的参数"
	case RteReduplicateTag:
		desc = "重复的标签点名"
	case RteReduplicateTabname:
		desc = "重复的表名称"
	case RteReduplicateTabid:
		desc = "重复的表ID"
	case RteTableNotFound:
		desc = "指定的表不存在"
	case RteUnsupportedClassof:
		desc = "不支持的标签点类别"
	case RteWrongOrDuplicTag:
		desc = "错误的或重复的标签点名"
	case RteReduplicatePt:
		desc = "重复的标签点标识"
	case RtePointLicenseFull:
		desc = "标签点数超出了许可证规定的最大数目"
	case RteTableLicenseFull:
		desc = "标签点表个数超出了许可证规定的最大数目"
	case RteWrongOrDuplicTabname:
		desc = "错误的或重复的表名称"
	case RteInvalidFileFormat:
		desc = "无效的数据文件格式"
	case RteWrongTabname:
		desc = "错误的表名称"
	case RteWrongTag:
		desc = "错误的标签点名"
	case RteNotInScope:
		desc = "数值超出了应属的范围"
	case RteCantLoadBase:
		desc = "不能同标签点信息服务取得联系"
	case RteCantLoadSnapshot:
		desc = "不能同快照数据服务取得联系"
	case RteCantLoadHistory:
		desc = "不能同历史数据服务取得联系"
	case RteCantLoadEquation:
		desc = "不能同实施方程式服务取得联系"
	case RteArraySizeNotMatch:
		desc = "数组尺寸不匹配"
	case RteInvalidHostAddress:
		desc = "无效的主机地址"
	case RteConnectFalse:
		desc = "连接已断开"
	case RteToomanyBytesRecved:
		desc = "接收到的数据长度超出了指定字节长度"
	case RteReqidRespidNotMatch:
		desc = "应答与请求的ID不一致"
	case RteLessBytesRecved:
		desc = "接收到的数据长度小于指定字节长度"
	case RteUnsupportedCalcMode:
		desc = "不支持的计算模式"
	case RteUnsupportedDataType:
		desc = "不支持的标签点类型"
	case RteInvalidExpression:
		desc = "无效的表达式"
	case RteIncondDataNotFound:
		desc = "找不到符合条件的数据"
	case RteValidDataNotFound:
		desc = "找不到需要的有效数据"
	case RteValueOrStateIsNan:
		desc = "数据或状态不正常，为NAN"
	case RteCreateMutexFailed:
		desc = "创建互斥对象失败"
	case RteTlsallocfail:
		desc = "处理TLS时调用系统函数LocalAlloc()失败，可能因为内存不足导致"
	case RteToManyPoints:
		desc = "正在调用的API函数不支持过多的标签点数量，请参考函数声明和开发手册"
	case RteLicInfoError:
		desc = "获取授权许可协议信息时发生错误"
	case RteArchiveBufferFull:
		desc = "标签点的历史数据补写缓冲区已满，请稍后再补"
	case RteUserNotExist:
		desc = "用户不存在"
	case RteUserIsLocked:
		desc = "帐户被锁定,需要管理员为您解锁"
	case RteWrongPassword:
		desc = "错误的口令"
	case RteAccessIsDenied:
		desc = "访问被拒绝，请确定是否具有足够的权限"
	case RteHaveNotLogin:
		desc = "您尚未登录，请先登录"
	case RteUserIsDeleted:
		desc = "帐户已被删除"
	case RteUserAlreadyExist:
		desc = "帐户已存在"
	case RteWrongCreateTabname:
		desc = "创建删除表失败"
	case RteWrongFieldValue:
		desc = "标签点属性值有错误"
	case RteInvalidTagId:
		desc = "无效的标签点ID"
	case RteCheckNamedTypeNameError:
		desc = "无效的自定义类型名称或字段名称"
	case RteCantLoadDispatch:
		desc = "不能同转发服务器取得联系"
	case RteConnectTimeOut:
		desc = "连接已超时，需要重新登录"
	case RteWrongLogin4:
		desc = "账户信息验证失败，还有4次尝试机会"
	case RteWrongLogin3:
		desc = "账户信息验证失败，还有3次尝试机会"
	case RteWrongLogin2:
		desc = "账户信息验证失败，还有2次尝试机会"
	case RteWrongLogin1:
		desc = "账户信息验证失败，还有1次尝试机会"
	case RteWrongDesc:
		desc = "错误的表描述"
	case RteWrongUnit:
		desc = "错误的工程单位"
	case RteWrongChanger:
		desc = "错误的最后一次被修改的用户名"
	case RteWrongCreator:
		desc = "错误的标签点创建者用户名"
	case RteWrongFull:
		desc = "错误的标签点全名"
	case RteWrongSource:
		desc = "错误的数据源"
	case RteWrongInstrument:
		desc = "错误的设备标签"
	case RteWrongUser:
		desc = "错误的创建者"
	case RteWrongEquation:
		desc = "错误的实时方程式"
	case RteWrongTypeName:
		desc = "错误的自定义类型名称"
	case RteWrongEncode:
		desc = "编码转换时出错"
	case RteWrongOthermask:
		desc = "错误的搜索类型转换mask值"
	case RteWrongType:
		desc = "错误的搜索类型"
	case RtePointHardwareLimited:
		desc = "由于硬件资源限制，创建或恢复标签点失败"
	case RteWaitingRecoverData:
		desc = "正在等待恢复数据完成，请稍后尝试连接"
	case RteReplicationLicMismatch:
		desc = "双活数据库授权不一致"
	case RteReadConfigFailed:
		desc = "读取配置文件失败"
	case RteUpdateConfigFailed:
		desc = "更新配置文件失败"
	case RteFilterTooLong:
		desc = "filter超过最大长度"
	case RteGetArchiveNameFail:
		desc = "获取存档文件名失败"
	case RteAutoMoveFailed:
		desc = "自动移动存档文件失败"
	case RteTimeGreaterThanHotTailArc:
		desc = "创建/入列非闪盘存档文件的时间大于闪盘最早的存档文件"
	case RteTimeLessThanColdBeginArc:
		desc = "创建/入列闪盘的存档文件时间小于非闪盘最新的存档文件"
	case RteRemoveEarliestArcFailed:
		desc = "删除最早的存档文件失败（存档文件列表为空）"
	case RteNoFreeTableId:
		desc = "没有空闲的表ID可用"
	case RteNoFreeTagPosition:
		desc = "没有空闲的标签点位址可用"
	case RteNoFreeScanTagPosition:
		desc = "没有空闲的采集标签点位址可用"
	case RteNoFreeCalcTagPosition:
		desc = "没有空闲的计算标签点位址可用"
	case RteInvalidIpcPosition:
		desc = "无效的位址被用于进程间内存共享"
	case RteWrongIpcPosition:
		desc = "错误的位址被用于进程间内存共享"
	case RteIpcAccessException:
		desc = "共享内存访问异常"
	case RteArvPageNotReady:
		desc = "没有空闲的历史数据缓存页"
	case RteArvexPageNotReady:
		desc = "没有空闲的补历史数据缓存页"
	case RteInvalidPositionFromId:
		desc = "依据标签点ID获得的位址无效"
	case RteNoActivePageAllocator:
		desc = "新的活动存档无法加载页分配器"
	case RteMapIsNotReady:
		desc = "内存映射尚未就绪"
	case RteFileMapFailed:
		desc = "文件映射到内存失败"
	case RteTimeRangeNotAllowed:
		desc = "不允许使用的时间区间"
	case RteNoDataForSummary:
		desc = "找不到用于统计的源数据"
	case RteCantOperateOnActived:
		desc = "不允许操作活动存档文件"
	case RteScanPointLicenseFull:
		desc = "采集标签点数超出了许可证规定的最大数目"
	case RteCalcPointLicenseFull:
		desc = "计算标签点数超出了许可证规定的最大数目"
	case RteHistorianIsShuttingdown:
		desc = "历史数据服务正在停止"
	case RteSnapshotIsShuttingdown:
		desc = "实时数据服务正在停止"
	case RteEquationIsShuttingdown:
		desc = "实时方程式服务正在停止"
	case RteBaseIsShuttingdown:
		desc = "标签点信息服务正在停止"
	case RteServerIsShuttingdown:
		desc = "网络通信服务正在停止"
	case RteOutOfMemory:
		desc = "内存不足"
	case RteInvalidPage:
		desc = "无效的数据页，有可能是未加载"
	case RtePageIsEmpty:
		desc = "遇到空的数据页"
	case RteStrOrBlobTooLong:
		desc = "字符串或BLOB数据长度超出限制"
	case RteCreatedOrOverdue:
		desc = "尚未产生任何快照或快照已过期"
	case RteArchiveInfoNotMatching:
		desc = "历史存档文件头部信息与实际不符"
	case RteTimeRangeOverlapping:
		desc = "指定的时间范围与已有存档文件重叠"
	case RteCannotShiftToActived:
		desc = "找不到合适的存档文件用于切换成活动存档"
	case RteIndexNotReady:
		desc = "历史存档文件对应的索引尚未就绪"
	case RteIndexNodeNotMatch:
		desc = "索引节点与指向的内容不符"
	case RteCanNotCreateIndex:
		desc = "无法创建索引节点"
	case RteCanNotRemoveIndex:
		desc = "无法删除索引节点"
	case RteInvalidFilterExpress:
		desc = "无效的过滤器表达式"
	case RteMoreVarInFilterExp:
		desc = "过滤器表达式中的包含了过多的变量"
	case RteInvalidArvPageAllocate:
		desc = "刚分配的历史数据缓存页ID与标签点事件对象ID不匹配"
	case RteInvalidArvexPageAllocate:
		desc = "刚分配的补历史数据缓存页ID与标签点事件对象ID不匹配"
	case RteBigJobIsNotDone:
		desc = "正在执行重要的任务，请稍后再试"
	case RteDatabaseNeedRestart:
		desc = "数据库需要重新启动以便应用新的参数"
	case RteInvalidTimeFormat:
		desc = "无效的时间格式字符串"
	case RteDataPlaybackDone:
		desc = "历史数据回放过程已结束"
	case RteBadEquation:
		desc = "错误的方程式"
	case RteNotEnoughSapce:
		desc = "剩余磁盘空间不足"
	case RteActivedArchiveExist:
		desc = "已存在活动存档"
	case RteArchiveHaveExFiles:
		desc = "指定的存档文件具有附属文件"
	case RteArchiveIsNotLatest:
		desc = "指定的存档文件不是最晚的"
	case RteDbSystemNotRunning:
		desc = "数据库管理系统尚未完全启动"
	case RteArchiveIsAltered:
		desc = "存档文件内容发生变更"
	case RteArchiveIsTooSmall:
		desc = "不允许创建太小的存档文件和附属文件"
	case RteInvalidIndexNode:
		desc = "遇到无效的索引节点"
	case RteModifySnapshotNotAllowed:
		desc = "不允许删除或修改快照事件"
	case RteSearchInterrupted:
		desc = "因目标正被创建、删除或恢复，搜索被迫中断，请稍后再试"
	case RteRecycleShutdown:
		desc = "回收站已失效，相关操作无法完成"
	case RteNeedToReindex:
		desc = "索引文件缺失，或部分索引节点被损坏，需要重建索引"
	case RteInvalidQuality:
		desc = "无效的质量码"
	case RteEquationNotReady:
		desc = "实时方程式服务正在解析，请稍后再试"
	case RteArchivesLicenseFull:
		desc = "存档文件数已达到许可证规定的最大数目"
	case RteRecycledLicenseFull:
		desc = "标签点回收站容量超出了许可证规定的最大数目"
	case RteStrBlobLicenseFull:
		desc = "字符串或BLOB类型标签点数量超出了许可证规定的最大数目"
	case RteNotSupportWhenDebug:
		desc = "此功能被某个调试选项禁用"
	case RteMappingAlreadyLoaded:
		desc = "映射已经被加载，不允许重复加载"
	case RteArchiveIsModified:
		desc = "存档文件被修改，动作被中断"
	case RteActiveArchiveFull:
		desc = "活动文档已满"
	case RteSplitNoData:
		desc = "拆分数据页后所给时间区间内没有数据"
	case RteInvalidDirectory:
		desc = "指定的路径不存在或无效"
	case RteArchiveLackExFiles:
		desc = "指定存档文件的部分附属文件缺失"
	case RteBigJobIsCanceled:
		desc = "后台任务被取消"
	case RteArvexBlobPageNotReady:
		desc = "没有空闲的blob补历史数据缓存页"
	case RteInvalidArvexBlobPageAllocate:
		desc = "刚分配的blob补历史数据缓存页ID与标签点事件对象ID不匹配"
	case RteTimestampEqualtoSnapshot:
		desc = "写入的时间与快照时间相同"
	case RteTimestampEarlierThanSnapshot:
		desc = "写入的时间比当前快照时间较早"
	case RteTimestampGreaterThanAllow:
		desc = "写入的时间超过了允许的时间"
	case RteTimestampBegintimeGreagerThanEndtime:
		desc = "开始时间大于结束时间"
	case RteTimestampBegintimeEqualtoEndtime:
		desc = "开始时间等于结束时间"
	case RteInvalidCount:
		desc = "无效的count"
	case RteInvalidCapacity:
		desc = "无效的capacity"
	case RteInvalidPath:
		desc = "无效的路径"
	case RteInvalidPosition:
		desc = "无效的position"
	case RteInvalidArvPage:
		desc = "无效的rtdb_arv_page<RTDB_T>,未加载，或者size小于等于0"
	case RteInvalidHisinfoItemState:
		desc = "无效的历史信息条目"
	case RteInvalidInterval:
		desc = "无效的间隔"
	case RteInvalidLength:
		desc = "无效的字符串长度"
	case RteInvalidSerachMode:
		desc = "无效的search mode"
	case RteInvalidFileId:
		desc = "无效的存档文件ID"
	case RteInvalidMillisecond:
		desc = "无效的毫秒值/纳秒值"
	case RteInvalidDeadline:
		desc = "无效的截止时间"
	case RteInvalidJobname:
		desc = "无效的Job名称"
	case RteInvalidJobstate:
		desc = "无效的Job状态"
	case RteInvalidProcessRate:
		desc = "无效的Process速率"
	case RteInvalidTableId:
		desc = "无效的表ID"
	case RteInvalidDataSource:
		desc = "无效的数据源格式"
	case RteInvalidTriggerMethod:
		desc = "无效的触发方式"
	case RteInvalidCalcTimeRes:
		desc = "无效的计算结果时间戳参考方式"
	case RteInvalidTriggerTimer:
		desc = "无效的定时触发触发周期,不能小于1秒"
	case RteInvalidLimit:
		desc = "工程上限不得低于工程下限"
	case RteInvalidCompTime:
		desc = "无效的压缩间隔，最长压缩间隔不得小于最短压缩间隔"
	case RteInvalidExtTime:
		desc = "无效的例外间隔，最长例外间隔不得小于最短例外间隔"
	case RteInvalidDigits:
		desc = "无效的数值位数，数值位数超出了范围,-20~10"
	case RteInvalidFullTagName:
		desc = "标签点全名有误，找不到表名与点名的分隔符“.”"
	case RteInvalidTableDesc:
		desc = "表描述信息过有误"
	case RteInvalidUserCount:
		desc = "非法的用户个数，小于0"
	case RteInvalidBlacklistCount:
		desc = "非法的黑名单个数，小于0"
	case RteInvalidAuthorizationCount:
		desc = "非法的信任连接个数，小于0"
	case RteInvalidBigJobType:
		desc = "非法的大任务类型"
	case RteInvalidSysParam:
		desc = "无效的系统参数，调用db_set_db_info2时，参数有误"
	case RteInvalidFileParam:
		desc = "无效的文件路径参数，调用db_set_db_info1时，参数有误"
	case RteInvalidFileSize:
		desc = "文件长度有误  < 1 baserecycle.dat、scanrecycle.dat、calcrecycle.dat、snaprecycle.dat"
	case RteInvalidTagType:
		desc = "标签点类型有误，合法（ rtdb_bool ~ rtdb_blob)，但是不属于相应函数的处理范围"
	case RteInvalidRecyStructPos:
		desc = "回收站对象最后一个结构体位置非法"
	case RteInvalidRecycleFile:
		desc = "scanrecycle.dat、baserecycle.dat  snaprecycle.dat文件不存在或失效"
	case RteInvalidSuffixName:
		desc = "无效的文件后缀名"
	case RteInsertStringFalse:
		desc = "向数据页中插入字符串数据失败"
	case RteBlobPageFull:
		desc = "blob数据页已满"
	case RteInvalidStringIteratorPointer:
		desc = "无效的str/blob迭代器指针"
	case RteNotEqualTagid:
		desc = "目标页标签点ID 与 当前ID不一致"
	case RtePathsOfArchiveAndAutobackAreSame:
		desc = "存档文件路径与自动备份路径相同"
	case RteXmlParseFail:
		desc = "xml文件解析失败"
	case RteXmlElementsAbsent:
		desc = "xml清单文件文件内容缺失"
	case RteXmlMismatchOnName:
		desc = "xml清单文件与本产品不匹配"
	case RteXmlMismatchOnVersion:
		desc = "xml清单文件版本不匹配"
	case RteXmlMismatchOnDatasize:
		desc = "xml清单文件数据尺寸不匹配"
	case RteXmlMismatchOnFileinfo:
		desc = "xml清单文件中数据文件信息不匹配"
	case RteXmlMismatchOnWindow:
		desc = "xml清单文件中所有数据文件的窗口大小必须一致"
	case RteXmlMismatchOnTypecount:
		desc = "xml清单文件自定义数据类型的数量不匹配"
	case RteXmlMismatchOnFieldcount:
		desc = "xml清单文件自定义数据类型的field不匹配"
	case RteXmlFieldMustInType:
		desc = "xml清单文件中field标签必须嵌套在type标签中"
	case RteInvalidNamedTypeFieldCount:
		desc = "无效的FIELD数量"
	case RteReduplicateFieldName:
		desc = "重复的FIELD名字"
	case RteInvalidNamedTypeName:
		desc = "无效的自定义数据类型的名字"
	case RteReduplicateNamedType:
		desc = "已经存在的自定义数据类型"
	case RteNotExistNamedType:
		desc = "不存在的自定义数据类型"
	case RteUpdateXmlFailed:
		desc = "更新XML清单文件失败"
	case RteNamedTypeUsedWithPoint:
		desc = "有些标签点正在使用此自定义数据类型，不允许删除"
	case RteNamedTypeUnsupportCalcPoint:
		desc = "自定义数据类型不支持计算点"
	case RteXmlMismatchOnMaxId:
		desc = "自定义数据类型的最大ID与实际的自定义数据类型数量不一致"
	case RteNamedTypeLicenseFull:
		desc = "自定义数据类型的数量超出了许可证规定的最大数目"
	case RteNoFreeNamedTypeId:
		desc = "没有空闲的自定义数据类型的ID"
	case RteInvalidNamedTypeId:
		desc = "无效的自定义数据类型ID"
	case RteInvalidNamedTypeFieldName:
		desc = "无效的自定义数据类型的字段名字"
	case RteNamedTypeUsedWithRecyclePoint:
		desc = "有些回收站中的标签点正在使用此自定义数据类型，不允许删除"
	case RteNamedTypeNameTooLong:
		desc = "自定义类型的名字超过了允许的最大长度"
	case RteNamedTypeFieldNameTooLong:
		desc = "自定义类型的field 名字超过了允许的最大长度"
	case RteInvalidNamedTypeFieldLength:
		desc = "无效的自定义数据类型的字段长度"
	case RteInvalidSearchMask:
		desc = "无效的高级搜索的标签点属性mask"
	case RteRecycledSpaceNotEnough:
		desc = "标签点回收站空闲空间不足"
	case RteDynamicLoadedMemoryNotInit:
		desc = "动态加载的内存未初始化"
	case RteForbidDynamicAllocType:
		desc = "内存库禁止动态分配类型"
	case RteMemorydbIndexCreateFailed:
		desc = "内存库索引创建失败"
	case RteWgMakeQueryReturnNull:
		desc = "whitedb make_query_rc返回null"
	case RteThtreadPoolCreatedFailed:
		desc = "内存库创建线程池失败"
	case RteMemorydbRemoveRecordFailed:
		desc = "内存库删除记录失败"
	case RteMemorydbConfigLoadFailed:
		desc = "内存库配置文件加载失败"
	case RteMemorydbProhibitDynamicAlloType:
		desc = "内存库禁止动态分配类型"
	case RteMemorydbDynamicAllocTypeFailed:
		desc = "内存库动态分配类型失败"
	case RteMemorydbStorageFileNameParseFailed:
		desc = "内存库优先级文件名解析失败"
	case RteMemorydbTtreeIndexDamage:
		desc = "内存库T树索引损坏"
	case RteMemorydbConfigFailed:
		desc = "内存库配置文件错误"
	case RteMemorydbValueCountNotMatch:
		desc = "内存库记录的值个数不匹配。"
	case RteMemorydbFieldTypeNotMatch:
		desc = "内存库的字段类型不匹配"
	case RteMemorydbMemoryAllocFailed:
		desc = "内存库内存分配失败"
	case RteMemorydbMethodParamErr:
		desc = "内存库方法参数错误"
	case RteMemorydbQueryResultAllocFailed:
		desc = "内存库查询结果缓存分配失败"
	case RteFilePathLength:
		desc = "指定的文件路径长度错误"
	case RteMemorydbFileVersionMatch:
		desc = "内存库文件版本不匹配"
	case RteMemorydbFileCrcError:
		desc = "内存库文件CRC错误"
	case RteMemorydbFileFlagMatch:
		desc = "内存库文件标志错误"
	case RteMemorydbInexistence:
		desc = "存储库不存在"
	case RteMemorydbLoadFailed:
		desc = "存储库加载失败"
	case RteNoDataInInterval:
		desc = "指定的查询区间内没有数据。"
	case RteCantLoadMemorydb:
		desc = "不能与内存服务取得联系"
	case RteQueryInWhitedb:
		desc = "查询内存库过程中出现了错误，这是whitedb内部错误"
	case RteNoDatabaseMemorydb:
		desc = "没有找到指定数据类型所对应的分库"
	case RteRecordNotGet:
		desc = "从whitedb中获取记录失败"
	case RteMemoryAllocErr:
		desc = "内存库用于接收快照的缓冲区分配失败"
	case RteEventCreateFailed:
		desc = "用于内存库接收缓冲区的事件创建失败"
	case RteGetPointFailed:
		desc = "获取标签点失败"
	case RteMemoryInitFailed:
		desc = "内存库初始化失败"
	case RteDatatypeNotMatch:
		desc = "数据类型不匹配"
	case RteGetFieldErr:
		desc = "在whitedb获取记录的字段时出现了错误"
	case RteMemorydbInternalErr:
		desc = "whitedb内部未知错误"
	case RteMemorydbRecordCreatedFailed:
		desc = "内存库创建记录失败"
	case RteParseNormalTypeSnapshotErr:
		desc = "解析普通数据类型的快照失败"
	case RteParseNamedTypeSnapshotErr:
		desc = "解析自定义数据类型的快照失败"
	case RteStringBlobTypeUnsupportCalcPoint:
		desc = "string、blob类型不支持计算点"
	case RteCoorTypeUnsupportCalcPoint:
		desc = "坐标类型不支持计算点"
	case RteIncludeHisData:
		desc = "记录是历史数据，可能是无效过期的脏数据"
	case RteThreadCreateErr:
		desc = "线程创建失败"
	case RteXmlCrcError:
		desc = "xml文件crc校验失败"
	case RteOversizeIntervals:
		desc = "intervals >"
	case RteDatetimesMustAscendingOrder:
		desc = "时间必须按升序排序"
	case RteCantLoadPerf:
		desc = "不能同性能计数服务取得联系"
	case RtePerfTagNotFound:
		desc = "性能计数点不存在"
	case RteWaitDataEmpty:
		desc = "数据为空"
	case RteWaitDataFull:
		desc = "数据满了"
	case RteDataTypeCountLess:
		desc = "数据类型数量最小值"
	case RteMemorydbCreateFailed:
		desc = "内存库创建失败"
	case RteMemorydbFieldEncodeFailed:
		desc = "内存库字段编码失败"
	case RteRecordCreateFailed:
		desc = "内存库记录创建失败"
	case RteRemoveRecordErr:
		desc = "内存库记录删除失败"
	case RteMemorydbFileOpenField:
		desc = "内存库打开文件失败"
	case RteMemorydbFileWriteFailed:
		desc = "内存库写入文件失败"
	case RteFilterWtihFloatAndEqual:
		desc = "含有浮点数不等式中不能有"
	case RteDispatchPluginNotExsit:
		desc = "转发服务器插件不存在"
	case RteDispatchPluginFileNotExsit:
		desc = "转发服务器插件DLL文件不存在"
	case RteDispatchPluginAlreadyExsit:
		desc = "转发服务器插件已存在"
	case RteDispatchRegisterPluginFailure:
		desc = "插件注册失败"
	case RteDispatchStartPluginFailure:
		desc = "启动插件失败"
	case RteDispatchStopPluginFailure:
		desc = "停止插件失败"
	case RteDispatchSetPluginEnableStatusFailure:
		desc = "设置插件状态失败"
	case RteDispatchGetPluginCountFailure:
		desc = "获取插件个数信息失败"
	case RteDispatchConfigfileNotExist:
		desc = "转发服务配置文件不存在"
	case RteDispatchConfigDataParseErr:
		desc = "转发服务配置数据解析错误"
	case RteDispatchPluginAlreadyRunning:
		desc = "转发服务器插件已经运行"
	case RteDispatchPluginCannotRun:
		desc = "转发服务器插件禁止运行"
	case RteDispatchPluginContainerUnrun:
		desc = "转发服务器插件容器未运行"
	case RteDispatchPluginInterfaceErr:
		desc = "转发服务器插件接口未实现"
	case RteDispatchPluginSaveConfigErr:
		desc = "转发服务器保存配置文件出错"
	case RteDispatchPluginStartErr:
		desc = "转发服务器插件启动时失败"
	case RteDispatchPluginStopErr:
		desc = "转发服务器插件停止时失败"
	case RteDispatchParseDataPageErr:
		desc = "不支持的数据页类型"
	case RteDispatchNotRun:
		desc = "转发服务未启用"
	case RteBigJobIsCanceledBecauseArcRoll:
		desc = "因存档文件滚动，后台任务被取消"
	case RtePerfForbiddenOperation:
		desc = "禁止对性能表的操作"
	case RteReduplicateTagInDestTable:
		desc = "目标表中存在同名的标签点（用于标签点移动）"
	case RteProtocolnotimpl:
		desc = "用户请求的报文未实现"
	case RteCrcerror:
		desc = "报文CRC校验错误"
	case RteWrongUserpw:
		desc = "验证用户名密码失败"
	case RteChangeUserpw:
		desc = "修改用户名密码失败"
	case RteInvalidHandle:
		desc = "无效的句柄"
	case RteInvalidSocketHandle:
		desc = "无效的套接字句柄"
	case RteFalse:
		desc = "操作未成功完成，具体原因查看小错误码。"
	case RteScanPointNotFound:
		desc = "要求访问的采集标签点不存在或无效"
	case RteCalcPointNotFound:
		desc = "要求访问的计算标签点不存在或无效"
	case RteReduplicateId:
		desc = "重复的标签点标识"
	case RteHandleSubscribed:
		desc = "句柄已经被订阅"
	case RteOtherSdkDoing:
		desc = "另一个API正在执行"
	case RteBatchEnd:
		desc = "分段数据返回结束"
	case RteAuthNotFound:
		desc = "信任连接段不存在"
	case RteAuthExist:
		desc = "连接地址段已经位于信任列表中"
	case RteAuthFull:
		desc = "信任连接段已满"
	case RteUserFull:
		desc = "用户已满"
	case RteVersionUnmatch:
		desc = "报文或数据版本不匹配"
	case RteInvalidPriv:
		desc = "无效的权限"
	case RteInvalidMask:
		desc = "无效的子网掩码"
	case RteInvalidUsername:
		desc = "无效的用户名"
	case RteInvalidMark:
		desc = "无法识别的报文头标记"
	case RteUnexpectedMethod:
		desc = "意外的消息 ID"
	case RteInvalidParamIndex:
		desc = "无效的系统参数索引值"
	case RteDecodePacketError:
		desc = "解包错误"
	case RteEncodePacketError:
		desc = "编包错误"
	case RteBlacklistFull:
		desc = "阻止连接段已满"
	case RteBlacklistExist:
		desc = "连接地址段已经位于黑名单中"
	case RteBlacklistNotFound:
		desc = "阻止连接段不存在"
	case RteInBlacklist:
		desc = "连接地址位于黑名单中，被主动拒绝"
	case RteIncreaseFileFailed:
		desc = "试图增大文件失败"
	case RteRpcInterfaceFailed:
		desc = "远程过程接口调用失败"
	case RteConnectionFull:
		desc = "连接已满"
	case RteOneClientConnectionFull:
		desc = "连接已达到单个客户端允许连接数的最大值"
	case RteServerClutterPoolNotEnough:
		desc = "网络数据交换空间不足"
	case RteEquationClutterPoolNotEnough:
		desc = "实时方程式交换空间不足"
	case RteNamedTypeNameLenError:
		desc = "自定义类型的名称过长"
	case RteNamedTypeLengthNotMatch:
		desc = "数值长度与自定义类型的定义不符"
	case RteCanNotUpdateSummary:
		desc = "无法更新卫星数据"
	case RteTooManyArvexFile:
		desc = "附属文件太多，无法继续创建附属文件"
	case RteNotSupportedFeature:
		desc = "测试版本，暂时不支持此功能"
	case RteEnsureError:
		desc = "验证信息失败，详细信息请查看数据库日志"
	case RteOperatorIsCancel:
		desc = "操作被取消"
	case RteMsgbodyRevError:
		desc = "报文体接收失败"
	case RteUncompressFailed:
		desc = "解压缩失败"
	case RteCompressFailed:
		desc = "压缩失败"
	case RteSubscribeError:
		desc = "订阅失败，前一个订阅线程尚未退出"
	case RteSubscribeCancelError:
		desc = "取消订阅失败"
	case RteSubscribeCallbackFailed:
		desc = "订阅回掉函数中不能调用取消订阅、断开连接"
	case RteSubscribeGreaterMaxCount:
		desc = "超过单连接可订阅标签点数量"
	case RteKillConnectionFailed:
		desc = "断开连接失败，无法断开自身连接"
	case RteSubscribeNotMatch:
		desc = "请求的方法与当前的订阅不匹配"
	case RteNoSubscribe:
		desc = "连接还未发起订阅，或者标签点还未订阅"
	case RteAlreadySubscribe:
		desc = "标签点已经被订阅"
	case RteCalcPointUnsupportedWriteData:
		desc = "计算点不支持写入数据"
	case RteFeatureDeprecated:
		desc = "不再支持此功能"
	case RteInvalidValue:
		desc = "无效的数据"
	case RteVerifyVercodeFailed:
		desc = "验证授权码失败"
	case RteInvalidPageSize:
		desc = "无效的数据页的大小"
	case RteInvalidPrecision:
		desc = "无效的时间戳精度"
	case RteInvalidPageVersion:
		desc = "无效的数据页版本"
	case RtePageIsFull:
		desc = "数据页已满"
	case RtePageNotLoaded:
		desc = "还未加载数据页"
	case RtePageAlreadyLoaded:
		desc = "已经加载了数据页"
	case RtePageTooSmall:
		desc = "数据页太小，有效空间小于数据长度"
	case RtePageNoEnoughData:
		desc = "数据页中没有足够的数据"
	case RtePageInsertFailed:
		desc = "数据页插入数据失败"
	case RtePageNoEnoughSpace:
		desc = "数据页没有足够的空间"
	case RteModifingMetaData:
		desc = "正在修改元数据，请稍后再试"
	case RtePageSizeNotMatch:
		desc = "数据页大小不匹配"
	case RteSyncBegin:
		desc = "元数据同步错误码起始值"
	case RteSyncInvalidConfig:
		desc = "元数据同步-无效的配置"
	case RteSyncInvalidVersion:
		desc = "元数据同步-无效的版本号"
	case RteSyncConfirmExpired:
		desc = "元数据同步-等待确认信息过期"
	case RteSyncTooManyFwdinfo:
		desc = "元数据同步-转发信息过多"
	case RteSyncNotMaster:
		desc = "元数据同步-不是主库"
	case RteSyncSyncing:
		desc = "元数据同步-正在同步"
	case RteSyncUnsynced:
		desc = "元数据同步-未同步"
	case RteSyncTablePosConflict:
		desc = "元数据同步-表位置冲突"
	case RteSyncInvalidPointId:
		desc = "元数据同步-无效的标签点ID"
	case RteSyncInvalidTableId:
		desc = "元数据同步-无效的表ID"
	case RteSyncInvalidNamedTypeId:
		desc = "元数据同步-无效的自定义类型ID"
	case RteSyncRestoring:
		desc = "元数据同步-正在重建元数据"
	case RteSyncServerIsNotRunning:
		desc = "元数据同步-网络服务不是运行状态"
	case RteSyncWriteWalFailed:
		desc = "元数据同步-写WAL失败"
	case RteSyncEnd:
		desc = "元数据同步错误码结束值"
	case RteNetError:
		desc = "网络错误的起始值"
	case RteSockWsaeintr:
		desc = "（阻塞）调用被 WSACancelBlockingCall() 函数取消"
	case RteSockWsaeacces:
		desc = "请求地址是广播地址，但是相应的 flags 没设置"
	case RteSockWsaefault:
		desc = "非法内存访问"
	case RteSockWsaemfile:
		desc = "无多余的描述符可用"
	case RteSockWsaewouldblock:
		desc = "套接字被标识为非阻塞，但操作将被阻塞"
	case RteSockWsaeinprogress:
		desc = "一个阻塞的 Windows Sockets 操作正在进行"
	case RteSockWsaealready:
		desc = "一个非阻塞的 connect() 调用已经在指定的套接字上进行"
	case RteSockWsaenotsock:
		desc = "描述符不是套接字描述符"
	case RteSockWsaedestaddrreq:
		desc = "要求（未指定）目的地址"
	case RteSockWsaemsgsize:
		desc = "套接字为基于消息的，消息太大（大于底层传输支持的最大值）"
	case RteSockWsaeprototype:
		desc = "对此套接字来说，指定协议是错误的类型"
	case RteSockWsaeprotonosupport:
		desc = "不支持指定协议"
	case RteSockWsaesocktnosupport:
		desc = "在此地址族中不支持指定套接字类型"
	case RteSockWsaeopnotsupp:
		desc = "MSG_OOB 被指定，但是套接字不是流风格的"
	case RteSockWsaeafnosupport:
		desc = "不支持指定的地址族"
	case RteSockWsaeaddrinuse:
		desc = "套接字的本地地址已被使用"
	case RteSockWsaeaddrnotavail:
		desc = "远程地址非法"
	case RteSockWsaenetdown:
		desc = "Windows Sockets 检测到网络系统已经失效"
	case RteSockWsaenetunreach:
		desc = "网络无法到达主机"
	case RteSockWsaenetreset:
		desc = "在操作进行时 keep-alive 活动检测到一个失败，连接被中断"
	case RteSockWsaeconnaborted:
		desc = "连接因超时或其他失败而中断"
	case RteSockWsaeconnreset:
		desc = "连接被复位"
	case RteSockWsaenobufs:
		desc = "无缓冲区空间可用"
	case RteSockWsaeisconn:
		desc = "连接已建立"
	case RteSockWsaenotconn:
		desc = "套接字未建立连接"
	case RteSockWsaeshutdown:
		desc = "套接字已 shutdown，连接已断开"
	case RteSockWsaetimedout:
		desc = "连接请求超时，未能建立连接"
	case RteSockWsaeconnrefused:
		desc = "连接被拒绝"
	case RteSockWsaeclose:
		desc = "连接被关闭"
	case RteSockWsanotinitialised:
		desc = "Windows Sockets DLL 未初始化"
	case RteCErrnoError:
		desc = "C语言errno错误的起始值"
	case RteCErrnoEperm:
		desc = "Operation not permitted"
	case RteCErrnoEnoent:
		desc = "No such file or directory"
	case RteCErrnoEsrch:
		desc = "No such process"
	case RteCErrnoEintr:
		desc = "Interrupted system call"
	case RteCErrnoEio:
		desc = "I/O error"
	case RteCErrnoEnxio:
		desc = "No such device or address"
	case RteCErrnoE2big:
		desc = "Argument list too long"
	case RteCErrnoEnoexec:
		desc = "Exec format error"
	case RteCErrnoEbadf:
		desc = "Bad file number"
	case RteCErrnoEchild:
		desc = "No child processes"
	case RteCErrnoEagain:
		desc = "Try again"
	case RteCErrnoEnomem:
		desc = "Out of memory"
	case RteCErrnoEacces:
		desc = "Permission denied"
	case RteCErrnoEfault:
		desc = "Bad address"
	case RteCErrnoEnotblk:
		desc = "Block device required"
	case RteCErrnoEbusy:
		desc = "Device or resource busy"
	case RteCErrnoEexist:
		desc = "File exists"
	case RteCErrnoExdev:
		desc = "Cross-device link"
	case RteCErrnoEnodev:
		desc = "No such device"
	case RteCErrnoEnotdir:
		desc = "Not a directory"
	case RteCErrnoEisdir:
		desc = "Is a directory"
	case RteCErrnoEinval:
		desc = "Invalid argument"
	case RteCErrnoEnfile:
		desc = "File table overflow"
	case RteCErrnoEmfile:
		desc = "Too many open files"
	case RteCErrnoEnotty:
		desc = "Not a typewriter"
	case RteCErrnoEtxtbsy:
		desc = "Text file busy"
	case RteCErrnoEfbig:
		desc = "File too large"
	case RteCErrnoEnospc:
		desc = "No space left on device"
	case RteCErrnoEspipe:
		desc = "Illegal seek"
	case RteCErrnoErofs:
		desc = "Read-only file system"
	case RteCErrnoEmlink:
		desc = "Too many links"
	case RteCErrnoEpipe:
		desc = "Broken pipe"
	case RteCErrnoEdom:
		desc = "Math argument out of domain of func"
	case RteCErrnoErange:
		desc = "Math result not representable"
	case RteCErrnoEdeadlk:
		desc = "Resource deadlock would occur"
	case RteCErrnoEnametoolong:
		desc = "File name too long"
	case RteCErrnoEnolck:
		desc = "No record locks available"
	case RteCErrnoEnosys:
		desc = "Function not implemented"
	case RteCErrnoEnotempty:
		desc = "Directory not empty"
	case RteCErrnoEloop:
		desc = "Too many symbolic links encountered"
	case RteCErrnoEnomsg:
		desc = "No message of desired type"
	case RteCErrnoEidrm:
		desc = "Identifier removed"
	case RteCErrnoEchrng:
		desc = "Channel number out of range"
	case RteCErrnoEl2nsync:
		desc = "Level 2 not synchronized"
	case RteCErrnoEl3hlt:
		desc = "Level 3 halted"
	case RteCErrnoEl3rst:
		desc = "Level 3 reset"
	case RteCErrnoElnrng:
		desc = "Link number out of range"
	case RteCErrnoEunatch:
		desc = "Protocol driver not attached"
	case RteCErrnoEnocsi:
		desc = "No CSI structure available"
	case RteCErrnoEl2hlt:
		desc = "Level 2 halted"
	case RteCErrnoEbade:
		desc = "Invalid exchange"
	case RteCErrnoEbadr:
		desc = "Invalid request descriptor"
	case RteCErrnoExfull:
		desc = "Exchange full"
	case RteCErrnoEnoano:
		desc = "No anode"
	case RteCErrnoEbadrqc:
		desc = "Invalid request code"
	case RteCErrnoEbadslt:
		desc = "Invalid slot"
	case RteCErrnoEbfont:
		desc = "Bad font file format"
	case RteCErrnoEnostr:
		desc = "Device not a stream"
	case RteCErrnoEnodata:
		desc = "No data available"
	case RteCErrnoEtime:
		desc = "Timer expired"
	case RteCErrnoEnosr:
		desc = "Out of streams resources"
	case RteCErrnoEnonet:
		desc = "Machine is not on the network"
	case RteCErrnoEnopkg:
		desc = "Package not installed"
	case RteCErrnoEremote:
		desc = "Object is remote"
	case RteCErrnoEnolink:
		desc = "Link has been severed"
	case RteCErrnoEadv:
		desc = "Advertise error"
	case RteCErrnoEsrmnt:
		desc = "Srmount error"
	case RteCErrnoEcomm:
		desc = "Communication error on send"
	case RteCErrnoEproto:
		desc = "Protocol error"
	case RteCErrnoEmultihop:
		desc = "Multihop attempted"
	case RteCErrnoEdotdot:
		desc = "RFS specific error"
	case RteCErrnoEbadmsg:
		desc = "Not a data message"
	case RteCErrnoEoverflow:
		desc = "Value too large for defined data type"
	case RteCErrnoEnotuniq:
		desc = "Name not unique on network"
	case RteCErrnoEbadfd:
		desc = "File descriptor in bad state"
	case RteCErrnoEremchg:
		desc = "Remote address changed"
	case RteCErrnoElibacc:
		desc = "Can not access a needed shared library"
	case RteCErrnoElibbad:
		desc = "Accessing a corrupted shared library"
	case RteCErrnoElibscn:
		desc = ".lib section in a.out corrupted"
	case RteCErrnoElibmax:
		desc = "Attempting to link in too many shared libraries"
	case RteCErrnoElibexec:
		desc = "Cannot exec a shared library directly"
	case RteCErrnoEilseq:
		desc = "Illegal byte sequence"
	case RteCErrnoErestart:
		desc = "Interrupted system call should be restarted"
	case RteCErrnoEstrpipe:
		desc = "Streams pipe error"
	case RteCErrnoEusers:
		desc = "Too many users"
	case RteCErrnoEnotsock:
		desc = "Socket operation on non-socket"
	case RteCErrnoEdestaddrreq:
		desc = "Destination address required"
	case RteCErrnoEmsgsize:
		desc = "Message too long"
	case RteCErrnoEprototype:
		desc = "Protocol wrong type for socket"
	case RteCErrnoEnoprotoopt:
		desc = "Protocol not available"
	case RteCErrnoEprotonosupport:
		desc = "Protocol not supported"
	case RteCErrnoEsocktnosupport:
		desc = "Socket type not supported"
	case RteCErrnoEopnotsupp:
		desc = "Operation not supported on transport endpoint"
	case RteCErrnoEpfnosupport:
		desc = "Protocol family not supported"
	case RteCErrnoEafnosupport:
		desc = "Address family not supported by protocol"
	case RteCErrnoEaddrinuse:
		desc = "Address already in use"
	case RteCErrnoEaddrnotavail:
		desc = "Cannot assign requested address"
	case RteCErrnoEnetdown:
		desc = "Network is down"
	case RteCErrnoEnetunreach:
		desc = "Network is unreachable"
	case RteCErrnoEnetreset:
		desc = "Network dropped connection because of reset"
	case RteCErrnoEconnaborted:
		desc = "Software caused connection abort"
	case RteCErrnoEconnreset:
		desc = "Connection reset by peer"
	case RteCErrnoEnobufs:
		desc = "No buffer space available"
	case RteCErrnoEisconn:
		desc = "Transport endpoint is already connected"
	case RteCErrnoEnotconn:
		desc = "Transport endpoint is not connected"
	case RteCErrnoEshutdown:
		desc = "Cannot send after transport endpoint shutdown"
	case RteCErrnoEtoomanyrefs:
		desc = "Too many references: cannot splice"
	case RteCErrnoEtimedout:
		desc = "Connection timed out"
	case RteCErrnoEconnrefused:
		desc = "Connection refused"
	case RteCErrnoEhostdown:
		desc = "Host is down"
	case RteCErrnoEhostunreach:
		desc = "No route to host"
	case RteCErrnoEalready:
		desc = "Operation already in progress"
	case RteCErrnoEinprogress:
		desc = "Operation now in progress"
	case RteCErrnoEstale:
		desc = "Stale file handle"
	case RteCErrnoEuclean:
		desc = "Structure needs cleaning"
	case RteCErrnoEnotnam:
		desc = "Not a XENIX named type file"
	case RteCErrnoEnavail:
		desc = "No XENIX semaphores available"
	case RteCErrnoEisnam:
		desc = "Is a named type file"
	case RteCErrnoEremoteio:
		desc = "Remote I/O error"
	case RteCErrnoEdquot:
		desc = "Quota exceeded"
	case RteCErrnoEnomedium:
		desc = "No medium found"
	case RteCErrnoEmediumtype:
		desc = "Wrong medium type"
	case RteCErrnoEcanceled:
		desc = "Operation Canceled"
	case RteCErrnoEnokey:
		desc = "Required key not available"
	case RteCErrnoEkeyexpired:
		desc = "Key has expired"
	case RteCErrnoEkeyrevoked:
		desc = "Key has been revoked"
	case RteCErrnoEkeyrejected:
		desc = "Key was rejected by service"
	case RteCErrnoEownerdead:
		desc = "Owner died"
	case RteCErrnoEnotrecoverable:
		desc = "State not recoverable"
	case RteCErrnoErfkill:
		desc = "Operation not possible due to RF-kill"
	case RteCErrnoEhwpoison:
		desc = "Memory page has hardware error"
	case RteIpcError:
		desc = "ipc error begin"
	case RteIpcErrorEnd:
		desc = "ipc error end"
	default:
		desc = "未知错误"
	}
	return desc
}

func (re RtdbError) GoError() error {
	if re.IsOk() {
		return nil
	} else {
		return re
	}
}

const (
	// RteUnknownError  未知错误
	RteUnknownError = RtdbError(C.RtE_UNKNOWN_ERROR)

	// RteOk  操作成功完成
	RteOk = RtdbError(C.RtE_OK)

	// RteWindowsError  Windows操作系统错误的起始值
	RteWindowsError = RtdbError(C.RtE_WINDOWS_ERROR)

	// RteWindowsErrorMax  Windows操作系统错误的结束值
	RteWindowsErrorMax = RtdbError(C.RtE_WINDOWS_ERROR_MAX)

	// RteInvalidOpenmode  无效的文件打开方式
	RteInvalidOpenmode = RtdbError(C.RtE_INVALID_OPENMODE)

	// RteOpenfileFailed  打开文件失败
	RteOpenfileFailed = RtdbError(C.RtE_OPENFILE_FAILED)

	// RteMovetoendFailed  移动文件指针到文件尾失败
	RteMovetoendFailed = RtdbError(C.RtE_MOVETOEND_FAILED)

	// RteDifferReadbytes  读取内容长度与要求不符
	RteDifferReadbytes = RtdbError(C.RtE_DIFFER_READBYTES)

	// RteGetfileposFailed  获取当前文件指针失败
	RteGetfileposFailed = RtdbError(C.RtE_GETFILEPOS_FAILED)

	// RteFlushfileFailed  刷新文件缓冲区失败
	RteFlushfileFailed = RtdbError(C.RtE_FLUSHFILE_FAILED)

	// RteSetsizeFailed  设置文件大小失败
	RteSetsizeFailed = RtdbError(C.RtE_SETSIZE_FAILED)

	// RteFileNotClosed  试图用未关闭的文件对象创建或打开文件
	RteFileNotClosed = RtdbError(C.RtE_FILE_NOT_CLOSED)

	// RteFileUnknown  创建或打开文件时必须指定文件名
	RteFileUnknown = RtdbError(C.RtE_FILE_UNKNOWN)

	// RteInvalidHeader  数据文件头信息错误
	RteInvalidHeader = RtdbError(C.RtE_INVALID_HEADER)

	// RteDisabledFile  数据文件无效，试图访问无效数据文件
	RteDisabledFile = RtdbError(C.RtE_DISABLED_FILE)

	// RteFileNotOpened  试图访问尚未打开的数据文件
	RteFileNotOpened = RtdbError(C.RtE_FILE_NOT_OPENED)

	// RtePointNotFound  要求访问的标签点不存在或无效
	RtePointNotFound = RtdbError(C.RtE_POINT_NOT_FOUND)

	// RteReadyblockNotFound  数据文件中找不到从指定数据块以后的可用的空块
	RteReadyblockNotFound = RtdbError(C.RtE_READYBLOCK_NOT_FOUND)

	// RteFileIsIncult  文件未被使用过
	RteFileIsIncult = RtdbError(C.RtE_FILE_IS_INCULT)

	// RteFileIsFull  数据文件已满
	RteFileIsFull = RtdbError(C.RtE_FILE_IS_FULL)

	// RteFileexIsFull  数据文件扩展区已满，无法继续装载数据
	RteFileexIsFull = RtdbError(C.RtE_FILEEX_IS_FULL)

	// RteInvalidDataType  无效的数据类型
	RteInvalidDataType = RtdbError(C.RtE_INVALID_DATA_TYPE)

	// RteDatablockNotFound  找不到符合时间条件的数据块
	RteDatablockNotFound = RtdbError(C.RtE_DATABLOCK_NOT_FOUND)

	// RteDataBetweenBlock  数据时间位于找到的块和下一个数据块之间
	RteDataBetweenBlock = RtdbError(C.RtE_DATA_BETWEEN_BLOCK)

	// RteCantModifyExistValue  不允许修改已存在的数据
	RteCantModifyExistValue = RtdbError(C.RtE_CANT_MODIFY_EXIST_VALUE)

	// RteWrongdataInBlock  块中有错误数据导致数据块头信息不符
	RteWrongdataInBlock = RtdbError(C.RtE_WRONGDATA_IN_BLOCK)

	// RteDatatimeNotIn  数据文件中没有该标签点指定时间的数据
	RteDatatimeNotIn = RtdbError(C.RtE_DATATIME_NOT_IN)

	// RteNullArchivePath  操作的数据文件路径为空
	RteNullArchivePath = RtdbError(C.RtE_NULL_ARCHIVE_PATH)

	// RteRegArchivePath  数据文件已被注册
	RteRegArchivePath = RtdbError(C.RtE_REG_ARCHIVE_PATH)

	// RteUnregArchivePath  未注册的数据文件
	RteUnregArchivePath = RtdbError(C.RtE_UNREG_ARCHIVE_PATH)

	// RteFileInexistence  指定的文件不存在
	RteFileInexistence = RtdbError(C.RtE_FILE_INEXISTENCE)

	// RteDataTypeNotMatch  数据类型不匹配
	RteDataTypeNotMatch = RtdbError(C.RtE_DATA_TYPE_NOT_MATCH)

	// RteFileIsReadonly  不允许修改只读数据文件中的数据
	RteFileIsReadonly = RtdbError(C.RtE_FILE_IS_READONLY)

	// RteTomanyArchiveFile  过多的数据文件
	RteTomanyArchiveFile = RtdbError(C.RtE_TOMANY_ARCHIVE_FILE)

	// RteNoPointsList  缺少标签点列表
	RteNoPointsList = RtdbError(C.RtE_NO_POINTS_LIST)

	// RteNoActivedArchive  缺少活动文档
	RteNoActivedArchive = RtdbError(C.RtE_NO_ACTIVED_ARCHIVE)

	// RteNoArchiveFile  缺少数据文档
	RteNoArchiveFile = RtdbError(C.RtE_NO_ARCHIVE_FILE)

	// RteNeedActivedArchive  只能在活动文档上执行该操作
	RteNeedActivedArchive = RtdbError(C.RtE_NEED_ACTIVED_ARCHIVE)

	// RteInvalidTimestamp  无效的时间戳
	RteInvalidTimestamp = RtdbError(C.RtE_INVALID_TIMESTAMP)

	// RteNeedMoreWritable  非只读文档个数太少
	RteNeedMoreWritable = RtdbError(C.RtE_NEED_MORE_WRITABLE)

	// RteNoArchiveForPut  找不到合适的追加历史数据的文档
	RteNoArchiveForPut = RtdbError(C.RtE_NO_ARCHIVE_FOR_PUT)

	// RteInvalidValueMode  无效的取值模式
	RteInvalidValueMode = RtdbError(C.RtE_INVALID_VALUE_MODE)

	// RteDataNotFound  找不到需要的数据
	RteDataNotFound = RtdbError(C.RtE_DATA_NOT_FOUND)

	// RteInvalidParameter  无效的参数
	RteInvalidParameter = RtdbError(C.RtE_INVALID_PARAMETER)

	// RteReduplicateTag  重复的标签点名
	RteReduplicateTag = RtdbError(C.RtE_REDUPLICATE_TAG)

	// RteReduplicateTabname  重复的表名称
	RteReduplicateTabname = RtdbError(C.RtE_REDUPLICATE_TABNAME)

	// RteReduplicateTabid  重复的表ID
	RteReduplicateTabid = RtdbError(C.RtE_REDUPLICATE_TABID)

	// RteTableNotFound  指定的表不存在
	RteTableNotFound = RtdbError(C.RtE_TABLE_NOT_FOUND)

	// RteUnsupportedClassof  不支持的标签点类别
	RteUnsupportedClassof = RtdbError(C.RtE_UNSUPPORTED_CLASSOF)

	// RteWrongOrDuplicTag  错误的或重复的标签点名
	RteWrongOrDuplicTag = RtdbError(C.RtE_WRONG_OR_DUPLIC_TAG)

	// RteReduplicatePt  重复的标签点标识
	RteReduplicatePt = RtdbError(C.RtE_REDUPLICATE_PT)

	// RtePointLicenseFull  标签点数超出了许可证规定的最大数目
	RtePointLicenseFull = RtdbError(C.RtE_POINT_LICENSE_FULL)

	// RteTableLicenseFull  标签点表个数超出了许可证规定的最大数目
	RteTableLicenseFull = RtdbError(C.RtE_TABLE_LICENSE_FULL)

	// RteWrongOrDuplicTabname  错误的或重复的表名称
	RteWrongOrDuplicTabname = RtdbError(C.RtE_WRONG_OR_DUPLIC_TABNAME)

	// RteInvalidFileFormat  无效的数据文件格式
	RteInvalidFileFormat = RtdbError(C.RtE_INVALID_FILE_FORMAT)

	// RteWrongTabname  错误的表名称
	RteWrongTabname = RtdbError(C.RtE_WRONG_TABNAME)

	// RteWrongTag  错误的标签点名
	RteWrongTag = RtdbError(C.RtE_WRONG_TAG)

	// RteNotInScope  数值超出了应属的范围
	RteNotInScope = RtdbError(C.RtE_NOT_IN_SCOPE)

	// RteCantLoadBase  不能同标签点信息服务取得联系
	RteCantLoadBase = RtdbError(C.RtE_CANT_LOAD_BASE)

	// RteCantLoadSnapshot  不能同快照数据服务取得联系
	RteCantLoadSnapshot = RtdbError(C.RtE_CANT_LOAD_SNAPSHOT)

	// RteCantLoadHistory  不能同历史数据服务取得联系
	RteCantLoadHistory = RtdbError(C.RtE_CANT_LOAD_HISTORY)

	// RteCantLoadEquation  不能同实施方程式服务取得联系
	RteCantLoadEquation = RtdbError(C.RtE_CANT_LOAD_EQUATION)

	// RteArraySizeNotMatch  数组尺寸不匹配
	RteArraySizeNotMatch = RtdbError(C.RtE_ARRAY_SIZE_NOT_MATCH)

	// RteInvalidHostAddress  无效的主机地址
	RteInvalidHostAddress = RtdbError(C.RtE_INVALID_HOST_ADDRESS)

	// RteConnectFalse  连接已断开
	RteConnectFalse = RtdbError(C.RtE_CONNECT_FALSE)

	// RteToomanyBytesRecved  接收到的数据长度超出了指定字节长度
	RteToomanyBytesRecved = RtdbError(C.RtE_TOOMANY_BYTES_RECVED)

	// RteReqidRespidNotMatch  应答与请求的ID不一致
	RteReqidRespidNotMatch = RtdbError(C.RtE_REQID_RESPID_NOT_MATCH)

	// RteLessBytesRecved  接收到的数据长度小于指定字节长度
	RteLessBytesRecved = RtdbError(C.RtE_LESS_BYTES_RECVED)

	// RteUnsupportedCalcMode  不支持的计算模式
	RteUnsupportedCalcMode = RtdbError(C.RtE_UNSUPPORTED_CALC_MODE)

	// RteUnsupportedDataType  不支持的标签点类型
	RteUnsupportedDataType = RtdbError(C.RtE_UNSUPPORTED_DATA_TYPE)

	// RteInvalidExpression  无效的表达式
	RteInvalidExpression = RtdbError(C.RtE_INVALID_EXPRESSION)

	// RteIncondDataNotFound  找不到符合条件的数据
	RteIncondDataNotFound = RtdbError(C.RtE_INCOND_DATA_NOT_FOUND)

	// RteValidDataNotFound  找不到需要的有效数据
	RteValidDataNotFound = RtdbError(C.RtE_VALID_DATA_NOT_FOUND)

	// RteValueOrStateIsNan  数据或状态不正常，为NAN
	RteValueOrStateIsNan = RtdbError(C.RtE_VALUE_OR_STATE_IS_NAN)

	// RteCreateMutexFailed  创建互斥对象失败
	RteCreateMutexFailed = RtdbError(C.RtE_CREATE_MUTEX_FAILED)

	// RteTlsallocfail  处理TLS时调用系统函数LocalAlloc()失败，可能因为内存不足导致
	RteTlsallocfail = RtdbError(C.RtE_TLSALLOCFAIL)

	// RteToManyPoints  正在调用的API函数不支持过多的标签点数量，请参考函数声明和开发手册
	RteToManyPoints = RtdbError(C.RtE_TO_MANY_POINTS)

	// RteLicInfoError  获取授权许可协议信息时发生错误
	RteLicInfoError = RtdbError(C.RtE_LIC_INFO_ERROR)

	// RteArchiveBufferFull  标签点的历史数据补写缓冲区已满，请稍后再补
	RteArchiveBufferFull = RtdbError(C.RtE_ARCHIVE_BUFFER_FULL)

	// RteUserNotExist  用户不存在
	RteUserNotExist = RtdbError(C.RtE_USER_NOT_EXIST)

	// RteUserIsLocked  帐户被锁定,需要管理员为您解锁
	RteUserIsLocked = RtdbError(C.RtE_USER_IS_LOCKED)

	// RteWrongPassword  错误的口令
	RteWrongPassword = RtdbError(C.RtE_WRONG_PASSWORD)

	// RteAccessIsDenied  访问被拒绝，请确定是否具有足够的权限
	RteAccessIsDenied = RtdbError(C.RtE_ACCESS_IS_DENIED)

	// RteHaveNotLogin  您尚未登录，请先登录
	RteHaveNotLogin = RtdbError(C.RtE_HAVE_NOT_LOGIN)

	// RteUserIsDeleted  帐户已被删除
	RteUserIsDeleted = RtdbError(C.RtE_USER_IS_DELETED)

	// RteUserAlreadyExist  帐户已存在
	RteUserAlreadyExist = RtdbError(C.RtE_USER_ALREADY_EXIST)

	// RteWrongCreateTabname  创建删除表失败
	RteWrongCreateTabname = RtdbError(C.RtE_WRONG_CREATE_TABNAME)

	// RteWrongFieldValue  标签点属性值有错误
	RteWrongFieldValue = RtdbError(C.RtE_WRONG_FIELD_VALUE)

	// RteInvalidTagId  无效的标签点ID
	RteInvalidTagId = RtdbError(C.RtE_INVALID_TAG_ID)

	// RteCheckNamedTypeNameError  无效的自定义类型名称或字段名称
	RteCheckNamedTypeNameError = RtdbError(C.RtE_CHECK_NAMED_TYPE_NAME_ERROR)

	// RteCantLoadDispatch  不能同转发服务器取得联系
	RteCantLoadDispatch = RtdbError(C.RtE_CANT_LOAD_DISPATCH)

	// RteConnectTimeOut  连接已超时，需要重新登录
	RteConnectTimeOut = RtdbError(C.RtE_CONNECT_TIME_OUT)

	// RteWrongLogin4  账户信息验证失败，还有4次尝试机会
	RteWrongLogin4 = RtdbError(C.RtE_WRONG_LOGIN_4)

	// RteWrongLogin3  账户信息验证失败，还有3次尝试机会
	RteWrongLogin3 = RtdbError(C.RtE_WRONG_LOGIN_3)

	// RteWrongLogin2  账户信息验证失败，还有2次尝试机会
	RteWrongLogin2 = RtdbError(C.RtE_WRONG_LOGIN_2)

	// RteWrongLogin1  账户信息验证失败，还有1次尝试机会
	RteWrongLogin1 = RtdbError(C.RtE_WRONG_LOGIN_1)

	// RteWrongDesc  错误的表描述
	RteWrongDesc = RtdbError(C.RtE_WRONG_DESC)

	// RteWrongUnit  错误的工程单位
	RteWrongUnit = RtdbError(C.RtE_WRONG_UNIT)

	// RteWrongChanger  错误的最后一次被修改的用户名
	RteWrongChanger = RtdbError(C.RtE_WRONG_CHANGER)

	// RteWrongCreator  错误的标签点创建者用户名
	RteWrongCreator = RtdbError(C.RtE_WRONG_CREATOR)

	// RteWrongFull  错误的标签点全名
	RteWrongFull = RtdbError(C.RtE_WRONG_FULL)

	// RteWrongSource  错误的数据源
	RteWrongSource = RtdbError(C.RtE_WRONG_SOURCE)

	// RteWrongInstrument  错误的设备标签
	RteWrongInstrument = RtdbError(C.RtE_WRONG_INSTRUMENT)

	// RteWrongUser  错误的创建者
	RteWrongUser = RtdbError(C.RtE_WRONG_USER)

	// RteWrongEquation  错误的实时方程式
	RteWrongEquation = RtdbError(C.RtE_WRONG_EQUATION)

	// RteWrongTypeName  错误的自定义类型名称
	RteWrongTypeName = RtdbError(C.RtE_WRONG_TYPE_NAME)

	// RteWrongEncode  编码转换时出错
	RteWrongEncode = RtdbError(C.RtE_WRONG_ENCODE)

	// RteWrongOthermask  错误的搜索类型转换mask值
	RteWrongOthermask = RtdbError(C.RtE_WRONG_OTHERMASK)

	// RteWrongType  错误的搜索类型
	RteWrongType = RtdbError(C.RtE_WRONG_TYPE)

	// RtePointHardwareLimited  由于硬件资源限制，创建或恢复标签点失败
	RtePointHardwareLimited = RtdbError(C.RtE_POINT_HARDWARE_LIMITED)

	// RteWaitingRecoverData  正在等待恢复数据完成，请稍后尝试连接
	RteWaitingRecoverData = RtdbError(C.RtE_WAITING_RECOVER_DATA)

	// RteReplicationLicMismatch  双活数据库授权不一致
	RteReplicationLicMismatch = RtdbError(C.RtE_REPLICATION_LIC_MISMATCH)

	// RteReadConfigFailed  读取配置文件失败
	RteReadConfigFailed = RtdbError(C.RtE_READ_CONFIG_FAILED)

	// RteUpdateConfigFailed  更新配置文件失败
	RteUpdateConfigFailed = RtdbError(C.RtE_UPDATE_CONFIG_FAILED)

	// RteFilterTooLong  filter超过最大长度
	RteFilterTooLong = RtdbError(C.RtE_FILTER_TOO_LONG)

	// RteGetArchiveNameFail  获取存档文件名失败
	RteGetArchiveNameFail = RtdbError(C.RtE_GET_ARCHIVE_NAME_FAIL)

	// RteAutoMoveFailed  自动移动存档文件失败
	RteAutoMoveFailed = RtdbError(C.RtE_AUTO_MOVE_FAILED)

	// RteTimeGreaterThanHotTailArc  创建/入列非闪盘存档文件的时间大于闪盘最早的存档文件
	RteTimeGreaterThanHotTailArc = RtdbError(C.RtE_TIME_GREATER_THAN_HOT_TAIL_ARC)

	// RteTimeLessThanColdBeginArc  创建/入列闪盘的存档文件时间小于非闪盘最新的存档文件
	RteTimeLessThanColdBeginArc = RtdbError(C.RtE_TIME_LESS_THAN_COLD_BEGIN_ARC)

	// RteRemoveEarliestArcFailed  删除最早的存档文件失败（存档文件列表为空）
	RteRemoveEarliestArcFailed = RtdbError(C.RtE_REMOVE_EARLIEST_ARC_FAILED)

	// RteNoFreeTableId  没有空闲的表ID可用
	RteNoFreeTableId = RtdbError(C.RtE_NO_FREE_TABLE_ID)

	// RteNoFreeTagPosition  没有空闲的标签点位址可用
	RteNoFreeTagPosition = RtdbError(C.RtE_NO_FREE_TAG_POSITION)

	// RteNoFreeScanTagPosition  没有空闲的采集标签点位址可用
	RteNoFreeScanTagPosition = RtdbError(C.RtE_NO_FREE_SCAN_TAG_POSITION)

	// RteNoFreeCalcTagPosition  没有空闲的计算标签点位址可用
	RteNoFreeCalcTagPosition = RtdbError(C.RtE_NO_FREE_CALC_TAG_POSITION)

	// RteInvalidIpcPosition  无效的位址被用于进程间内存共享
	RteInvalidIpcPosition = RtdbError(C.RtE_INVALID_IPC_POSITION)

	// RteWrongIpcPosition  错误的位址被用于进程间内存共享
	RteWrongIpcPosition = RtdbError(C.RtE_WRONG_IPC_POSITION)

	// RteIpcAccessException  共享内存访问异常
	RteIpcAccessException = RtdbError(C.RtE_IPC_ACCESS_EXCEPTION)

	// RteArvPageNotReady  没有空闲的历史数据缓存页
	RteArvPageNotReady = RtdbError(C.RtE_ARV_PAGE_NOT_READY)

	// RteArvexPageNotReady  没有空闲的补历史数据缓存页
	RteArvexPageNotReady = RtdbError(C.RtE_ARVEX_PAGE_NOT_READY)

	// RteInvalidPositionFromId  依据标签点ID获得的位址无效
	RteInvalidPositionFromId = RtdbError(C.RtE_INVALID_POSITION_FROM_ID)

	// RteNoActivePageAllocator  新的活动存档无法加载页分配器
	RteNoActivePageAllocator = RtdbError(C.RtE_NO_ACTIVE_PAGE_ALLOCATOR)

	// RteMapIsNotReady  内存映射尚未就绪
	RteMapIsNotReady = RtdbError(C.RtE_MAP_IS_NOT_READY)

	// RteFileMapFailed  文件映射到内存失败
	RteFileMapFailed = RtdbError(C.RtE_FILE_MAP_FAILED)

	// RteTimeRangeNotAllowed  不允许使用的时间区间
	RteTimeRangeNotAllowed = RtdbError(C.RtE_TIME_RANGE_NOT_ALLOWED)

	// RteNoDataForSummary  找不到用于统计的源数据
	RteNoDataForSummary = RtdbError(C.RtE_NO_DATA_FOR_SUMMARY)

	// RteCantOperateOnActived  不允许操作活动存档文件
	RteCantOperateOnActived = RtdbError(C.RtE_CANT_OPERATE_ON_ACTIVED)

	// RteScanPointLicenseFull  采集标签点数超出了许可证规定的最大数目
	RteScanPointLicenseFull = RtdbError(C.RtE_SCAN_POINT_LICENSE_FULL)

	// RteCalcPointLicenseFull  计算标签点数超出了许可证规定的最大数目
	RteCalcPointLicenseFull = RtdbError(C.RtE_CALC_POINT_LICENSE_FULL)

	// RteHistorianIsShuttingdown  历史数据服务正在停止
	RteHistorianIsShuttingdown = RtdbError(C.RtE_HISTORIAN_IS_SHUTTINGDOWN)

	// RteSnapshotIsShuttingdown  实时数据服务正在停止
	RteSnapshotIsShuttingdown = RtdbError(C.RtE_SNAPSHOT_IS_SHUTTINGDOWN)

	// RteEquationIsShuttingdown  实时方程式服务正在停止
	RteEquationIsShuttingdown = RtdbError(C.RtE_EQUATION_IS_SHUTTINGDOWN)

	// RteBaseIsShuttingdown  标签点信息服务正在停止
	RteBaseIsShuttingdown = RtdbError(C.RtE_BASE_IS_SHUTTINGDOWN)

	// RteServerIsShuttingdown  网络通信服务正在停止
	RteServerIsShuttingdown = RtdbError(C.RtE_SERVER_IS_SHUTTINGDOWN)

	// RteOutOfMemory  内存不足
	RteOutOfMemory = RtdbError(C.RtE_OUT_OF_MEMORY)

	// RteInvalidPage  无效的数据页，有可能是未加载
	RteInvalidPage = RtdbError(C.RtE_INVALID_PAGE)

	// RtePageIsEmpty  遇到空的数据页
	RtePageIsEmpty = RtdbError(C.RtE_PAGE_IS_EMPTY)

	// RteStrOrBlobTooLong  字符串或BLOB数据长度超出限制
	RteStrOrBlobTooLong = RtdbError(C.RtE_STR_OR_BLOB_TOO_LONG)

	// RteCreatedOrOverdue  尚未产生任何快照或快照已过期
	RteCreatedOrOverdue = RtdbError(C.RtE_CREATED_OR_OVERDUE)

	// RteArchiveInfoNotMatching  历史存档文件头部信息与实际不符
	RteArchiveInfoNotMatching = RtdbError(C.RtE_ARCHIVE_INFO_NOT_MATCHING)

	// RteTimeRangeOverlapping  指定的时间范围与已有存档文件重叠
	RteTimeRangeOverlapping = RtdbError(C.RtE_TIME_RANGE_OVERLAPPING)

	// RteCannotShiftToActived  找不到合适的存档文件用于切换成活动存档
	RteCannotShiftToActived = RtdbError(C.RtE_CANNOT_SHIFT_TO_ACTIVED)

	// RteIndexNotReady  历史存档文件对应的索引尚未就绪
	RteIndexNotReady = RtdbError(C.RtE_INDEX_NOT_READY)

	// RteIndexNodeNotMatch  索引节点与指向的内容不符
	RteIndexNodeNotMatch = RtdbError(C.RtE_INDEX_NODE_NOT_MATCH)

	// RteCanNotCreateIndex  无法创建索引节点
	RteCanNotCreateIndex = RtdbError(C.RtE_CAN_NOT_CREATE_INDEX)

	// RteCanNotRemoveIndex  无法删除索引节点
	RteCanNotRemoveIndex = RtdbError(C.RtE_CAN_NOT_REMOVE_INDEX)

	// RteInvalidFilterExpress  无效的过滤器表达式
	RteInvalidFilterExpress = RtdbError(C.RtE_INVALID_FILTER_EXPRESS)

	// RteMoreVarInFilterExp  过滤器表达式中的包含了过多的变量
	RteMoreVarInFilterExp = RtdbError(C.RtE_MORE_VAR_IN_FILTER_EXP)

	// RteInvalidArvPageAllocate  刚分配的历史数据缓存页ID与标签点事件对象ID不匹配
	RteInvalidArvPageAllocate = RtdbError(C.RtE_INVALID_ARV_PAGE_ALLOCATE)

	// RteInvalidArvexPageAllocate  刚分配的补历史数据缓存页ID与标签点事件对象ID不匹配
	RteInvalidArvexPageAllocate = RtdbError(C.RtE_INVALID_ARVEX_PAGE_ALLOCATE)

	// RteBigJobIsNotDone  正在执行重要的任务，请稍后再试
	RteBigJobIsNotDone = RtdbError(C.RtE_BIG_JOB_IS_NOT_DONE)

	// RteDatabaseNeedRestart  数据库需要重新启动以便应用新的参数
	RteDatabaseNeedRestart = RtdbError(C.RtE_DATABASE_NEED_RESTART)

	// RteInvalidTimeFormat  无效的时间格式字符串
	RteInvalidTimeFormat = RtdbError(C.RtE_INVALID_TIME_FORMAT)

	// RteDataPlaybackDone  历史数据回放过程已结束
	RteDataPlaybackDone = RtdbError(C.RtE_DATA_PLAYBACK_DONE)

	// RteBadEquation  错误的方程式
	RteBadEquation = RtdbError(C.RtE_BAD_EQUATION)

	// RteNotEnoughSapce  剩余磁盘空间不足
	RteNotEnoughSapce = RtdbError(C.RtE_NOT_ENOUGH_SAPCE)

	// RteActivedArchiveExist  已存在活动存档
	RteActivedArchiveExist = RtdbError(C.RtE_ACTIVED_ARCHIVE_EXIST)

	// RteArchiveHaveExFiles  指定的存档文件具有附属文件
	RteArchiveHaveExFiles = RtdbError(C.RtE_ARCHIVE_HAVE_EX_FILES)

	// RteArchiveIsNotLatest  指定的存档文件不是最晚的
	RteArchiveIsNotLatest = RtdbError(C.RtE_ARCHIVE_IS_NOT_LATEST)

	// RteDbSystemNotRunning  数据库管理系统尚未完全启动
	RteDbSystemNotRunning = RtdbError(C.RtE_DB_SYSTEM_NOT_RUNNING)

	// RteArchiveIsAltered  存档文件内容发生变更
	RteArchiveIsAltered = RtdbError(C.RtE_ARCHIVE_IS_ALTERED)

	// RteArchiveIsTooSmall  不允许创建太小的存档文件和附属文件
	RteArchiveIsTooSmall = RtdbError(C.RtE_ARCHIVE_IS_TOO_SMALL)

	// RteInvalidIndexNode  遇到无效的索引节点
	RteInvalidIndexNode = RtdbError(C.RtE_INVALID_INDEX_NODE)

	// RteModifySnapshotNotAllowed  不允许删除或修改快照事件
	RteModifySnapshotNotAllowed = RtdbError(C.RtE_MODIFY_SNAPSHOT_NOT_ALLOWED)

	// RteSearchInterrupted  因目标正被创建、删除或恢复，搜索被迫中断，请稍后再试
	RteSearchInterrupted = RtdbError(C.RtE_SEARCH_INTERRUPTED)

	// RteRecycleShutdown  回收站已失效，相关操作无法完成
	RteRecycleShutdown = RtdbError(C.RtE_RECYCLE_SHUTDOWN)

	// RteNeedToReindex  索引文件缺失，或部分索引节点被损坏，需要重建索引
	RteNeedToReindex = RtdbError(C.RtE_NEED_TO_REINDEX)

	// RteInvalidQuality  无效的质量码
	RteInvalidQuality = RtdbError(C.RtE_INVALID_QUALITY)

	// RteEquationNotReady  实时方程式服务正在解析，请稍后再试
	RteEquationNotReady = RtdbError(C.RtE_EQUATION_NOT_READY)

	// RteArchivesLicenseFull  存档文件数已达到许可证规定的最大数目
	RteArchivesLicenseFull = RtdbError(C.RtE_ARCHIVES_LICENSE_FULL)

	// RteRecycledLicenseFull  标签点回收站容量超出了许可证规定的最大数目
	RteRecycledLicenseFull = RtdbError(C.RtE_RECYCLED_LICENSE_FULL)

	// RteStrBlobLicenseFull  字符串或BLOB类型标签点数量超出了许可证规定的最大数目
	RteStrBlobLicenseFull = RtdbError(C.RtE_STR_BLOB_LICENSE_FULL)

	// RteNotSupportWhenDebug  此功能被某个调试选项禁用
	RteNotSupportWhenDebug = RtdbError(C.RtE_NOT_SUPPORT_WHEN_DEBUG)

	// RteMappingAlreadyLoaded  映射已经被加载，不允许重复加载
	RteMappingAlreadyLoaded = RtdbError(C.RtE_MAPPING_ALREADY_LOADED)

	// RteArchiveIsModified  存档文件被修改，动作被中断
	RteArchiveIsModified = RtdbError(C.RtE_ARCHIVE_IS_MODIFIED)

	// RteActiveArchiveFull  活动文档已满
	RteActiveArchiveFull = RtdbError(C.RtE_ACTIVE_ARCHIVE_FULL)

	// RteSplitNoData  拆分数据页后所给时间区间内没有数据
	RteSplitNoData = RtdbError(C.RtE_SPLIT_NO_DATA)

	// RteInvalidDirectory  指定的路径不存在或无效
	RteInvalidDirectory = RtdbError(C.RtE_INVALID_DIRECTORY)

	// RteArchiveLackExFiles  指定存档文件的部分附属文件缺失
	RteArchiveLackExFiles = RtdbError(C.RtE_ARCHIVE_LACK_EX_FILES)

	// RteBigJobIsCanceled  后台任务被取消
	RteBigJobIsCanceled = RtdbError(C.RtE_BIG_JOB_IS_CANCELED)

	// RteArvexBlobPageNotReady  没有空闲的blob补历史数据缓存页
	RteArvexBlobPageNotReady = RtdbError(C.RtE_ARVEX_BLOB_PAGE_NOT_READY)

	// RteInvalidArvexBlobPageAllocate  刚分配的blob补历史数据缓存页ID与标签点事件对象ID不匹配
	RteInvalidArvexBlobPageAllocate = RtdbError(C.RtE_INVALID_ARVEX_BLOB_PAGE_ALLOCATE)

	// RteTimestampEqualtoSnapshot  写入的时间与快照时间相同
	RteTimestampEqualtoSnapshot = RtdbError(C.RtE_TIMESTAMP_EQUALTO_SNAPSHOT)

	// RteTimestampEarlierThanSnapshot  写入的时间比当前快照时间较早
	RteTimestampEarlierThanSnapshot = RtdbError(C.RtE_TIMESTAMP_EARLIER_THAN_SNAPSHOT)

	// RteTimestampGreaterThanAllow  写入的时间超过了允许的时间
	RteTimestampGreaterThanAllow = RtdbError(C.RtE_TIMESTAMP_GREATER_THAN_ALLOW)

	// RteTimestampBegintimeGreagerThanEndtime  开始时间大于结束时间
	RteTimestampBegintimeGreagerThanEndtime = RtdbError(C.RtE_TIMESTAMP_BEGINTIME_GREAGER_THAN_ENDTIME)

	// RteTimestampBegintimeEqualtoEndtime  开始时间等于结束时间
	RteTimestampBegintimeEqualtoEndtime = RtdbError(C.RtE_TIMESTAMP_BEGINTIME_EQUALTO_ENDTIME)

	// RteInvalidCount  无效的count
	RteInvalidCount = RtdbError(C.RtE_INVALID_COUNT)

	// RteInvalidCapacity  无效的capacity
	RteInvalidCapacity = RtdbError(C.RtE_INVALID_CAPACITY)

	// RteInvalidPath  无效的路径
	RteInvalidPath = RtdbError(C.RtE_INVALID_PATH)

	// RteInvalidPosition  无效的position
	RteInvalidPosition = RtdbError(C.RtE_INVALID_POSITION)

	// RteInvalidArvPage  无效的rtdb_arv_page<RTDB_T>,未加载，或者size小于等于0
	RteInvalidArvPage = RtdbError(C.RtE_INVALID_ARV_PAGE)

	// RteInvalidHisinfoItemState  无效的历史信息条目
	RteInvalidHisinfoItemState = RtdbError(C.RtE_INVALID_HISINFO_ITEM_STATE)

	// RteInvalidInterval  无效的间隔
	RteInvalidInterval = RtdbError(C.RtE_INVALID_INTERVAL)

	// RteInvalidLength  无效的字符串长度
	RteInvalidLength = RtdbError(C.RtE_INVALID_LENGTH)

	// RteInvalidSerachMode  无效的search mode
	RteInvalidSerachMode = RtdbError(C.RtE_INVALID_SERACH_MODE)

	// RteInvalidFileId  无效的存档文件ID
	RteInvalidFileId = RtdbError(C.RtE_INVALID_FILE_ID)

	// RteInvalidMillisecond  无效的毫秒值/纳秒值
	RteInvalidMillisecond = RtdbError(C.RtE_INVALID_MILLISECOND)

	// RteInvalidDeadline  无效的截止时间
	RteInvalidDeadline = RtdbError(C.RtE_INVALID_DEADLINE)

	// RteInvalidJobname  无效的Job名称
	RteInvalidJobname = RtdbError(C.RtE_INVALID_JOBNAME)

	// RteInvalidJobstate  无效的Job状态
	RteInvalidJobstate = RtdbError(C.RtE_INVALID_JOBSTATE)

	// RteInvalidProcessRate  无效的Process速率
	RteInvalidProcessRate = RtdbError(C.RtE_INVALID_PROCESS_RATE)

	// RteInvalidTableId  无效的表ID
	RteInvalidTableId = RtdbError(C.RtE_INVALID_TABLE_ID)

	// RteInvalidDataSource  无效的数据源格式
	RteInvalidDataSource = RtdbError(C.RtE_INVALID_DATA_SOURCE)

	// RteInvalidTriggerMethod  无效的触发方式
	RteInvalidTriggerMethod = RtdbError(C.RtE_INVALID_TRIGGER_METHOD)

	// RteInvalidCalcTimeRes  无效的计算结果时间戳参考方式
	RteInvalidCalcTimeRes = RtdbError(C.RtE_INVALID_CALC_TIME_RES)

	// RteInvalidTriggerTimer  无效的定时触发触发周期,不能小于1秒
	RteInvalidTriggerTimer = RtdbError(C.RtE_INVALID_TRIGGER_TIMER)

	// RteInvalidLimit  工程上限不得低于工程下限
	RteInvalidLimit = RtdbError(C.RtE_INVALID_LIMIT)

	// RteInvalidCompTime  无效的压缩间隔，最长压缩间隔不得小于最短压缩间隔
	RteInvalidCompTime = RtdbError(C.RtE_INVALID_COMP_TIME)

	// RteInvalidExtTime  无效的例外间隔，最长例外间隔不得小于最短例外间隔
	RteInvalidExtTime = RtdbError(C.RtE_INVALID_EXT_TIME)

	// RteInvalidDigits  无效的数值位数，数值位数超出了范围,-20~10
	RteInvalidDigits = RtdbError(C.RtE_INVALID_DIGITS)

	// RteInvalidFullTagName  标签点全名有误，找不到表名与点名的分隔符“.”
	RteInvalidFullTagName = RtdbError(C.RtE_INVALID_FULL_TAG_NAME)

	// RteInvalidTableDesc  表描述信息过有误
	RteInvalidTableDesc = RtdbError(C.RtE_INVALID_TABLE_DESC)

	// RteInvalidUserCount  非法的用户个数，小于0
	RteInvalidUserCount = RtdbError(C.RtE_INVALID_USER_COUNT)

	// RteInvalidBlacklistCount  非法的黑名单个数，小于0
	RteInvalidBlacklistCount = RtdbError(C.RtE_INVALID_BLACKLIST_COUNT)

	// RteInvalidAuthorizationCount  非法的信任连接个数，小于0
	RteInvalidAuthorizationCount = RtdbError(C.RtE_INVALID_AUTHORIZATION_COUNT)

	// RteInvalidBigJobType  非法的大任务类型
	RteInvalidBigJobType = RtdbError(C.RtE_INVALID_BIG_JOB_TYPE)

	// RteInvalidSysParam  无效的系统参数，调用db_set_db_info2时，参数有误
	RteInvalidSysParam = RtdbError(C.RtE_INVALID_SYS_PARAM)

	// RteInvalidFileParam  无效的文件路径参数，调用db_set_db_info1时，参数有误
	RteInvalidFileParam = RtdbError(C.RtE_INVALID_FILE_PARAM)

	// RteInvalidFileSize  文件长度有误  < 1 baserecycle.dat、scanrecycle.dat、calcrecycle.dat、snaprecycle.dat
	RteInvalidFileSize = RtdbError(C.RtE_INVALID_FILE_SIZE)

	// RteInvalidTagType  标签点类型有误，合法（ rtdb_bool ~ rtdb_blob)，但是不属于相应函数的处理范围
	RteInvalidTagType = RtdbError(C.RtE_INVALID_TAG_TYPE)

	// RteInvalidRecyStructPos  回收站对象最后一个结构体位置非法
	RteInvalidRecyStructPos = RtdbError(C.RtE_INVALID_RECY_STRUCT_POS)

	// RteInvalidRecycleFile  scanrecycle.dat、baserecycle.dat  snaprecycle.dat文件不存在或失效
	RteInvalidRecycleFile = RtdbError(C.RtE_INVALID_RECYCLE_FILE)

	// RteInvalidSuffixName  无效的文件后缀名
	RteInvalidSuffixName = RtdbError(C.RtE_INVALID_SUFFIX_NAME)

	// RteInsertStringFalse  向数据页中插入字符串数据失败
	RteInsertStringFalse = RtdbError(C.RtE_INSERT_STRING_FALSE)

	// RteBlobPageFull  blob数据页已满
	RteBlobPageFull = RtdbError(C.RtE_BLOB_PAGE_FULL)

	// RteInvalidStringIteratorPointer  无效的str/blob迭代器指针
	RteInvalidStringIteratorPointer = RtdbError(C.RtE_INVALID_STRING_ITERATOR_POINTER)

	// RteNotEqualTagid  目标页标签点ID 与 当前ID不一致
	RteNotEqualTagid = RtdbError(C.RtE_NOT_EQUAL_TAGID)

	// RtePathsOfArchiveAndAutobackAreSame  存档文件路径与自动备份路径相同
	RtePathsOfArchiveAndAutobackAreSame = RtdbError(C.RtE_PATHS_OF_ARCHIVE_AND_AUTOBACK_ARE_SAME)

	// RteXmlParseFail  xml文件解析失败
	RteXmlParseFail = RtdbError(C.RtE_XML_PARSE_FAIL)

	// RteXmlElementsAbsent  xml清单文件文件内容缺失
	RteXmlElementsAbsent = RtdbError(C.RtE_XML_ELEMENTS_ABSENT)

	// RteXmlMismatchOnName  xml清单文件与本产品不匹配
	RteXmlMismatchOnName = RtdbError(C.RtE_XML_MISMATCH_ON_NAME)

	// RteXmlMismatchOnVersion  xml清单文件版本不匹配
	RteXmlMismatchOnVersion = RtdbError(C.RtE_XML_MISMATCH_ON_VERSION)

	// RteXmlMismatchOnDatasize  xml清单文件数据尺寸不匹配
	RteXmlMismatchOnDatasize = RtdbError(C.RtE_XML_MISMATCH_ON_DATASIZE)

	// RteXmlMismatchOnFileinfo  xml清单文件中数据文件信息不匹配
	RteXmlMismatchOnFileinfo = RtdbError(C.RtE_XML_MISMATCH_ON_FILEINFO)

	// RteXmlMismatchOnWindow  xml清单文件中所有数据文件的窗口大小必须一致
	RteXmlMismatchOnWindow = RtdbError(C.RtE_XML_MISMATCH_ON_WINDOW)

	// RteXmlMismatchOnTypecount  xml清单文件自定义数据类型的数量不匹配
	RteXmlMismatchOnTypecount = RtdbError(C.RtE_XML_MISMATCH_ON_TYPECOUNT)

	// RteXmlMismatchOnFieldcount  xml清单文件自定义数据类型的field不匹配
	RteXmlMismatchOnFieldcount = RtdbError(C.RtE_XML_MISMATCH_ON_FIELDCOUNT)

	// RteXmlFieldMustInType  xml清单文件中field标签必须嵌套在type标签中
	RteXmlFieldMustInType = RtdbError(C.RtE_XML_FIELD_MUST_IN_TYPE)

	// RteInvalidNamedTypeFieldCount  无效的FIELD数量
	RteInvalidNamedTypeFieldCount = RtdbError(C.RtE_INVALID_NAMED_TYPE_FIELD_COUNT)

	// RteReduplicateFieldName  重复的FIELD名字
	RteReduplicateFieldName = RtdbError(C.RtE_REDUPLICATE_FIELD_NAME)

	// RteInvalidNamedTypeName  无效的自定义数据类型的名字
	RteInvalidNamedTypeName = RtdbError(C.RtE_INVALID_NAMED_TYPE_NAME)

	// RteReduplicateNamedType  已经存在的自定义数据类型
	RteReduplicateNamedType = RtdbError(C.RtE_REDUPLICATE_NAMED_TYPE)

	// RteNotExistNamedType  不存在的自定义数据类型
	RteNotExistNamedType = RtdbError(C.RtE_NOT_EXIST_NAMED_TYPE)

	// RteUpdateXmlFailed  更新XML清单文件失败
	RteUpdateXmlFailed = RtdbError(C.RtE_UPDATE_XML_FAILED)

	// RteNamedTypeUsedWithPoint  有些标签点正在使用此自定义数据类型，不允许删除
	RteNamedTypeUsedWithPoint = RtdbError(C.RtE_NAMED_TYPE_USED_WITH_POINT)

	// RteNamedTypeUnsupportCalcPoint  自定义数据类型不支持计算点
	RteNamedTypeUnsupportCalcPoint = RtdbError(C.RtE_NAMED_TYPE_UNSUPPORT_CALC_POINT)

	// RteXmlMismatchOnMaxId  自定义数据类型的最大ID与实际的自定义数据类型数量不一致
	RteXmlMismatchOnMaxId = RtdbError(C.RtE_XML_MISMATCH_ON_MAX_ID)

	// RteNamedTypeLicenseFull  自定义数据类型的数量超出了许可证规定的最大数目
	RteNamedTypeLicenseFull = RtdbError(C.RtE_NAMED_TYPE_LICENSE_FULL)

	// RteNoFreeNamedTypeId  没有空闲的自定义数据类型的ID
	RteNoFreeNamedTypeId = RtdbError(C.RtE_NO_FREE_NAMED_TYPE_ID)

	// RteInvalidNamedTypeId  无效的自定义数据类型ID
	RteInvalidNamedTypeId = RtdbError(C.RtE_INVALID_NAMED_TYPE_ID)

	// RteInvalidNamedTypeFieldName  无效的自定义数据类型的字段名字
	RteInvalidNamedTypeFieldName = RtdbError(C.RtE_INVALID_NAMED_TYPE_FIELD_NAME)

	// RteNamedTypeUsedWithRecyclePoint  有些回收站中的标签点正在使用此自定义数据类型，不允许删除
	RteNamedTypeUsedWithRecyclePoint = RtdbError(C.RtE_NAMED_TYPE_USED_WITH_RECYCLE_POINT)

	// RteNamedTypeNameTooLong  自定义类型的名字超过了允许的最大长度
	RteNamedTypeNameTooLong = RtdbError(C.RtE_NAMED_TYPE_NAME_TOO_LONG)

	// RteNamedTypeFieldNameTooLong  自定义类型的field 名字超过了允许的最大长度
	RteNamedTypeFieldNameTooLong = RtdbError(C.RtE_NAMED_TYPE_FIELD_NAME_TOO_LONG)

	// RteInvalidNamedTypeFieldLength  无效的自定义数据类型的字段长度
	RteInvalidNamedTypeFieldLength = RtdbError(C.RtE_INVALID_NAMED_TYPE_FIELD_LENGTH)

	// RteInvalidSearchMask  无效的高级搜索的标签点属性mask
	RteInvalidSearchMask = RtdbError(C.RtE_INVALID_SEARCH_MASK)

	// RteRecycledSpaceNotEnough  标签点回收站空闲空间不足
	RteRecycledSpaceNotEnough = RtdbError(C.RtE_RECYCLED_SPACE_NOT_ENOUGH)

	// RteDynamicLoadedMemoryNotInit  动态加载的内存未初始化
	RteDynamicLoadedMemoryNotInit = RtdbError(C.RtE_DYNAMIC_LOADED_MEMORY_NOT_INIT)

	// RteForbidDynamicAllocType  内存库禁止动态分配类型
	RteForbidDynamicAllocType = RtdbError(C.RtE_FORBID_DYNAMIC_ALLOC_TYPE)

	// RteMemorydbIndexCreateFailed  内存库索引创建失败
	RteMemorydbIndexCreateFailed = RtdbError(C.RtE_MEMORYDB_INDEX_CREATE_FAILED)

	// RteWgMakeQueryReturnNull  whitedb make_query_rc返回null
	RteWgMakeQueryReturnNull = RtdbError(C.RtE_WG_MAKE_QUERY_RETURN_NULL)

	// RteThtreadPoolCreatedFailed  内存库创建线程池失败
	RteThtreadPoolCreatedFailed = RtdbError(C.RtE_THTREAD_POOL_CREATED_FAILED)

	// RteMemorydbRemoveRecordFailed  内存库删除记录失败
	RteMemorydbRemoveRecordFailed = RtdbError(C.RtE_MEMORYDB_REMOVE_RECORD_FAILED)

	// RteMemorydbConfigLoadFailed  内存库配置文件加载失败
	RteMemorydbConfigLoadFailed = RtdbError(C.RtE_MEMORYDB_CONFIG_LOAD_FAILED)

	// RteMemorydbProhibitDynamicAlloType  内存库禁止动态分配类型
	RteMemorydbProhibitDynamicAlloType = RtdbError(C.RtE_MEMORYDB_PROHIBIT_DYNAMIC_ALLO_TYPE)

	// RteMemorydbDynamicAllocTypeFailed  内存库动态分配类型失败
	RteMemorydbDynamicAllocTypeFailed = RtdbError(C.RtE_MEMORYDB_DYNAMIC_ALLOC_TYPE_FAILED)

	// RteMemorydbStorageFileNameParseFailed  内存库优先级文件名解析失败
	RteMemorydbStorageFileNameParseFailed = RtdbError(C.RtE_MEMORYDB_STORAGE_FILE_NAME_PARSE_FAILED)

	// RteMemorydbTtreeIndexDamage  内存库T树索引损坏
	RteMemorydbTtreeIndexDamage = RtdbError(C.RtE_MEMORYDB_TTREE_INDEX_DAMAGE)

	// RteMemorydbConfigFailed  内存库配置文件错误
	RteMemorydbConfigFailed = RtdbError(C.RtE_MEMORYDB_CONFIG_FAILED)

	// RteMemorydbValueCountNotMatch  内存库记录的值个数不匹配。
	RteMemorydbValueCountNotMatch = RtdbError(C.RtE_MEMORYDB_VALUE_COUNT_NOT_MATCH)

	// RteMemorydbFieldTypeNotMatch  内存库的字段类型不匹配
	RteMemorydbFieldTypeNotMatch = RtdbError(C.RtE_MEMORYDB_FIELD_TYPE_NOT_MATCH)

	// RteMemorydbMemoryAllocFailed  内存库内存分配失败
	RteMemorydbMemoryAllocFailed = RtdbError(C.RtE_MEMORYDB_MEMORY_ALLOC_FAILED)

	// RteMemorydbMethodParamErr  内存库方法参数错误
	RteMemorydbMethodParamErr = RtdbError(C.RtE_MEMORYDB_METHOD_PARAM_ERR)

	// RteMemorydbQueryResultAllocFailed  内存库查询结果缓存分配失败
	RteMemorydbQueryResultAllocFailed = RtdbError(C.RtE_MEMORYDB_QUERY_RESULT_ALLOC_FAILED)

	// RteFilePathLength  指定的文件路径长度错误
	RteFilePathLength = RtdbError(C.RtE_FILE_PATH_LENGTH)

	// RteMemorydbFileVersionMatch  内存库文件版本不匹配
	RteMemorydbFileVersionMatch = RtdbError(C.RtE_MEMORYDB_FILE_VERSION_MATCH)

	// RteMemorydbFileCrcError  内存库文件CRC错误
	RteMemorydbFileCrcError = RtdbError(C.RtE_MEMORYDB_FILE_CRC_ERROR)

	// RteMemorydbFileFlagMatch  内存库文件标志错误
	RteMemorydbFileFlagMatch = RtdbError(C.RtE_MEMORYDB_FILE_FLAG_MATCH)

	// RteMemorydbInexistence  存储库不存在
	RteMemorydbInexistence = RtdbError(C.RtE_MEMORYDB_INEXISTENCE)

	// RteMemorydbLoadFailed  存储库加载失败
	RteMemorydbLoadFailed = RtdbError(C.RtE_MEMORYDB_LOAD_FAILED)

	// RteNoDataInInterval  指定的查询区间内没有数据。
	RteNoDataInInterval = RtdbError(C.RtE_NO_DATA_IN_INTERVAL)

	// RteCantLoadMemorydb  不能与内存服务取得联系
	RteCantLoadMemorydb = RtdbError(C.RtE_CANT_LOAD_MEMORYDB)

	// RteQueryInWhitedb  查询内存库过程中出现了错误，这是whitedb内部错误
	RteQueryInWhitedb = RtdbError(C.RtE_QUERY_IN_WHITEDB)

	// RteNoDatabaseMemorydb  没有找到指定数据类型所对应的分库
	RteNoDatabaseMemorydb = RtdbError(C.RtE_NO_DATABASE_MEMORYDB)

	// RteRecordNotGet  从whitedb中获取记录失败
	RteRecordNotGet = RtdbError(C.RtE_RECORD_NOT_GET)

	// RteMemoryAllocErr  内存库用于接收快照的缓冲区分配失败
	RteMemoryAllocErr = RtdbError(C.RtE_MEMORY_ALLOC_ERR)

	// RteEventCreateFailed  用于内存库接收缓冲区的事件创建失败
	RteEventCreateFailed = RtdbError(C.RtE_EVENT_CREATE_FAILED)

	// RteGetPointFailed  获取标签点失败
	RteGetPointFailed = RtdbError(C.RtE_GET_POINT_FAILED)

	// RteMemoryInitFailed  内存库初始化失败
	RteMemoryInitFailed = RtdbError(C.RtE_MEMORY_INIT_FAILED)

	// RteDatatypeNotMatch  数据类型不匹配
	RteDatatypeNotMatch = RtdbError(C.RtE_DATATYPE_NOT_MATCH)

	// RteGetFieldErr  在whitedb获取记录的字段时出现了错误
	RteGetFieldErr = RtdbError(C.RtE_GET_FIELD_ERR)

	// RteMemorydbInternalErr  whitedb内部未知错误
	RteMemorydbInternalErr = RtdbError(C.RtE_MEMORYDB_INTERNAL_ERR)

	// RteMemorydbRecordCreatedFailed  内存库创建记录失败
	RteMemorydbRecordCreatedFailed = RtdbError(C.RtE_MEMORYDB_RECORD_CREATED_FAILED)

	// RteParseNormalTypeSnapshotErr  解析普通数据类型的快照失败
	RteParseNormalTypeSnapshotErr = RtdbError(C.RtE_PARSE_NORMAL_TYPE_SNAPSHOT_ERR)

	// RteParseNamedTypeSnapshotErr  解析自定义数据类型的快照失败
	RteParseNamedTypeSnapshotErr = RtdbError(C.RtE_PARSE_NAMED_TYPE_SNAPSHOT_ERR)

	// RteStringBlobTypeUnsupportCalcPoint  string、blob类型不支持计算点
	RteStringBlobTypeUnsupportCalcPoint = RtdbError(C.RtE_STRING_BLOB_TYPE_UNSUPPORT_CALC_POINT)

	// RteCoorTypeUnsupportCalcPoint  坐标类型不支持计算点
	RteCoorTypeUnsupportCalcPoint = RtdbError(C.RtE_COOR_TYPE_UNSUPPORT_CALC_POINT)

	// RteIncludeHisData  记录是历史数据，可能是无效过期的脏数据
	RteIncludeHisData = RtdbError(C.RtE_INCLUDE_HIS_DATA)

	// RteThreadCreateErr  线程创建失败
	RteThreadCreateErr = RtdbError(C.RtE_THREAD_CREATE_ERR)

	// RteXmlCrcError  xml文件crc校验失败
	RteXmlCrcError = RtdbError(C.RtE_XML_CRC_ERROR)

	// RteOversizeIntervals  intervals >
	RteOversizeIntervals = RtdbError(C.RtE_OVERSIZE_INTERVALS)

	// RteDatetimesMustAscendingOrder  时间必须按升序排序
	RteDatetimesMustAscendingOrder = RtdbError(C.RtE_DATETIMES_MUST_ASCENDING_ORDER)

	// RteCantLoadPerf  不能同性能计数服务取得联系
	RteCantLoadPerf = RtdbError(C.RtE_CANT_LOAD_PERF)

	// RtePerfTagNotFound  性能计数点不存在
	RtePerfTagNotFound = RtdbError(C.RtE_PERF_TAG_NOT_FOUND)

	// RteWaitDataEmpty  数据为空
	RteWaitDataEmpty = RtdbError(C.RtE_WAIT_DATA_EMPTY)

	// RteWaitDataFull  数据满了
	RteWaitDataFull = RtdbError(C.RtE_WAIT_DATA_FULL)

	// RteDataTypeCountLess  数据类型数量最小值
	RteDataTypeCountLess = RtdbError(C.RtE_DATA_TYPE_COUNT_LESS)

	// RteMemorydbCreateFailed  内存库创建失败
	RteMemorydbCreateFailed = RtdbError(C.RtE_MEMORYDB_CREATE_FAILED)

	// RteMemorydbFieldEncodeFailed  内存库字段编码失败
	RteMemorydbFieldEncodeFailed = RtdbError(C.RtE_MEMORYDB_FIELD_ENCODE_FAILED)

	// RteRecordCreateFailed  内存库记录创建失败
	RteRecordCreateFailed = RtdbError(C.RtE_RECORD_CREATE_FAILED)

	// RteRemoveRecordErr  内存库记录删除失败
	RteRemoveRecordErr = RtdbError(C.RtE_REMOVE_RECORD_ERR)

	// RteMemorydbFileOpenField  内存库打开文件失败
	RteMemorydbFileOpenField = RtdbError(C.RtE_MEMORYDB_FILE_OPEN_FIELD)

	// RteMemorydbFileWriteFailed  内存库写入文件失败
	RteMemorydbFileWriteFailed = RtdbError(C.RtE_MEMORYDB_FILE_WRITE_FAILED)

	// RteFilterWtihFloatAndEqual  含有浮点数不等式中不能有
	RteFilterWtihFloatAndEqual = RtdbError(C.RtE_FILTER_WTIH_FLOAT_AND_EQUAL)

	// RteDispatchPluginNotExsit  转发服务器插件不存在
	RteDispatchPluginNotExsit = RtdbError(C.RtE_DISPATCH_PLUGIN_NOT_EXSIT)

	// RteDispatchPluginFileNotExsit  转发服务器插件DLL文件不存在
	RteDispatchPluginFileNotExsit = RtdbError(C.RtE_DISPATCH_PLUGIN_FILE_NOT_EXSIT)

	// RteDispatchPluginAlreadyExsit  转发服务器插件已存在
	RteDispatchPluginAlreadyExsit = RtdbError(C.RtE_DISPATCH_PLUGIN_ALREADY_EXSIT)

	// RteDispatchRegisterPluginFailure  插件注册失败
	RteDispatchRegisterPluginFailure = RtdbError(C.RtE_DISPATCH_REGISTER_PLUGIN_FAILURE)

	// RteDispatchStartPluginFailure  启动插件失败
	RteDispatchStartPluginFailure = RtdbError(C.RtE_DISPATCH_START_PLUGIN_FAILURE)

	// RteDispatchStopPluginFailure  停止插件失败
	RteDispatchStopPluginFailure = RtdbError(C.RtE_DISPATCH_STOP_PLUGIN_FAILURE)

	// RteDispatchSetPluginEnableStatusFailure  设置插件状态失败
	RteDispatchSetPluginEnableStatusFailure = RtdbError(C.RtE_DISPATCH_SET_PLUGIN_ENABLE_STATUS_FAILURE)

	// RteDispatchGetPluginCountFailure  获取插件个数信息失败
	RteDispatchGetPluginCountFailure = RtdbError(C.RtE_DISPATCH_GET_PLUGIN_COUNT_FAILURE)

	// RteDispatchConfigfileNotExist  转发服务配置文件不存在
	RteDispatchConfigfileNotExist = RtdbError(C.RtE_DISPATCH_CONFIGFILE_NOT_EXIST)

	// RteDispatchConfigDataParseErr  转发服务配置数据解析错误
	RteDispatchConfigDataParseErr = RtdbError(C.RtE_DISPATCH_CONFIG_DATA_PARSE_ERR)

	// RteDispatchPluginAlreadyRunning  转发服务器插件已经运行
	RteDispatchPluginAlreadyRunning = RtdbError(C.RtE_DISPATCH_PLUGIN_ALREADY_RUNNING)

	// RteDispatchPluginCannotRun  转发服务器插件禁止运行
	RteDispatchPluginCannotRun = RtdbError(C.RtE_DISPATCH_PLUGIN_CANNOT_RUN)

	// RteDispatchPluginContainerUnrun  转发服务器插件容器未运行
	RteDispatchPluginContainerUnrun = RtdbError(C.RtE_DISPATCH_PLUGIN_CONTAINER_UNRUN)

	// RteDispatchPluginInterfaceErr  转发服务器插件接口未实现
	RteDispatchPluginInterfaceErr = RtdbError(C.RtE_DISPATCH_PLUGIN_INTERFACE_ERR)

	// RteDispatchPluginSaveConfigErr  转发服务器保存配置文件出错
	RteDispatchPluginSaveConfigErr = RtdbError(C.RtE_DISPATCH_PLUGIN_SAVE_CONFIG_ERR)

	// RteDispatchPluginStartErr  转发服务器插件启动时失败
	RteDispatchPluginStartErr = RtdbError(C.RtE_DISPATCH_PLUGIN_START_ERR)

	// RteDispatchPluginStopErr  转发服务器插件停止时失败
	RteDispatchPluginStopErr = RtdbError(C.RtE_DISPATCH_PLUGIN_STOP_ERR)

	// RteDispatchParseDataPageErr  不支持的数据页类型
	RteDispatchParseDataPageErr = RtdbError(C.RtE_DISPATCH_PARSE_DATA_PAGE_ERR)

	// RteDispatchNotRun  转发服务未启用
	RteDispatchNotRun = RtdbError(C.RtE_DISPATCH_NOT_RUN)

	// RteBigJobIsCanceledBecauseArcRoll  因存档文件滚动，后台任务被取消
	RteBigJobIsCanceledBecauseArcRoll = RtdbError(C.RtE_BIG_JOB_IS_CANCELED_BECAUSE_ARC_ROLL)

	// RtePerfForbiddenOperation  禁止对性能表的操作
	RtePerfForbiddenOperation = RtdbError(C.RtE_PERF_FORBIDDEN_OPERATION)

	// RteReduplicateTagInDestTable  目标表中存在同名的标签点（用于标签点移动）
	RteReduplicateTagInDestTable = RtdbError(C.RtE_REDUPLICATE_TAG_IN_DEST_TABLE)

	// RteProtocolnotimpl  用户请求的报文未实现
	RteProtocolnotimpl = RtdbError(C.RtE_PROTOCOLNOTIMPL)

	// RteCrcerror  报文CRC校验错误
	RteCrcerror = RtdbError(C.RtE_CRCERROR)

	// RteWrongUserpw  验证用户名密码失败
	RteWrongUserpw = RtdbError(C.RtE_WRONG_USERPW)

	// RteChangeUserpw  修改用户名密码失败
	RteChangeUserpw = RtdbError(C.RtE_CHANGE_USERPW)

	// RteInvalidHandle  无效的句柄
	RteInvalidHandle = RtdbError(C.RtE_INVALID_HANDLE)

	// RteInvalidSocketHandle  无效的套接字句柄
	RteInvalidSocketHandle = RtdbError(C.RtE_INVALID_SOCKET_HANDLE)

	// RteFalse  操作未成功完成，具体原因查看小错误码。
	RteFalse = RtdbError(C.RtE_FALSE)

	// RteScanPointNotFound  要求访问的采集标签点不存在或无效
	RteScanPointNotFound = RtdbError(C.RtE_SCAN_POINT_NOT_FOUND)

	// RteCalcPointNotFound  要求访问的计算标签点不存在或无效
	RteCalcPointNotFound = RtdbError(C.RtE_CALC_POINT_NOT_FOUND)

	// RteReduplicateId  重复的标签点标识
	RteReduplicateId = RtdbError(C.RtE_REDUPLICATE_ID)

	// RteHandleSubscribed  句柄已经被订阅
	RteHandleSubscribed = RtdbError(C.RtE_HANDLE_SUBSCRIBED)

	// RteOtherSdkDoing  另一个API正在执行
	RteOtherSdkDoing = RtdbError(C.RtE_OTHER_SDK_DOING)

	// RteBatchEnd  分段数据返回结束
	RteBatchEnd = RtdbError(C.RtE_BATCH_END)

	// RteAuthNotFound  信任连接段不存在
	RteAuthNotFound = RtdbError(C.RtE_AUTH_NOT_FOUND)

	// RteAuthExist  连接地址段已经位于信任列表中
	RteAuthExist = RtdbError(C.RtE_AUTH_EXIST)

	// RteAuthFull  信任连接段已满
	RteAuthFull = RtdbError(C.RtE_AUTH_FULL)

	// RteUserFull  用户已满
	RteUserFull = RtdbError(C.RtE_USER_FULL)

	// RteVersionUnmatch  报文或数据版本不匹配
	RteVersionUnmatch = RtdbError(C.RtE_VERSION_UNMATCH)

	// RteInvalidPriv  无效的权限
	RteInvalidPriv = RtdbError(C.RtE_INVALID_PRIV)

	// RteInvalidMask  无效的子网掩码
	RteInvalidMask = RtdbError(C.RtE_INVALID_MASK)

	// RteInvalidUsername  无效的用户名
	RteInvalidUsername = RtdbError(C.RtE_INVALID_USERNAME)

	// RteInvalidMark  无法识别的报文头标记
	RteInvalidMark = RtdbError(C.RtE_INVALID_MARK)

	// RteUnexpectedMethod  意外的消息 ID
	RteUnexpectedMethod = RtdbError(C.RtE_UNEXPECTED_METHOD)

	// RteInvalidParamIndex  无效的系统参数索引值
	RteInvalidParamIndex = RtdbError(C.RtE_INVALID_PARAM_INDEX)

	// RteDecodePacketError  解包错误
	RteDecodePacketError = RtdbError(C.RtE_DECODE_PACKET_ERROR)

	// RteEncodePacketError  编包错误
	RteEncodePacketError = RtdbError(C.RtE_ENCODE_PACKET_ERROR)

	// RteBlacklistFull  阻止连接段已满
	RteBlacklistFull = RtdbError(C.RtE_BLACKLIST_FULL)

	// RteBlacklistExist  连接地址段已经位于黑名单中
	RteBlacklistExist = RtdbError(C.RtE_BLACKLIST_EXIST)

	// RteBlacklistNotFound  阻止连接段不存在
	RteBlacklistNotFound = RtdbError(C.RtE_BLACKLIST_NOT_FOUND)

	// RteInBlacklist  连接地址位于黑名单中，被主动拒绝
	RteInBlacklist = RtdbError(C.RtE_IN_BLACKLIST)

	// RteIncreaseFileFailed  试图增大文件失败
	RteIncreaseFileFailed = RtdbError(C.RtE_INCREASE_FILE_FAILED)

	// RteRpcInterfaceFailed  远程过程接口调用失败
	RteRpcInterfaceFailed = RtdbError(C.RtE_RPC_INTERFACE_FAILED)

	// RteConnectionFull  连接已满
	RteConnectionFull = RtdbError(C.RtE_CONNECTION_FULL)

	// RteOneClientConnectionFull  连接已达到单个客户端允许连接数的最大值
	RteOneClientConnectionFull = RtdbError(C.RtE_ONE_CLIENT_CONNECTION_FULL)

	// RteServerClutterPoolNotEnough  网络数据交换空间不足
	RteServerClutterPoolNotEnough = RtdbError(C.RtE_SERVER_CLUTTER_POOL_NOT_ENOUGH)

	// RteEquationClutterPoolNotEnough  实时方程式交换空间不足
	RteEquationClutterPoolNotEnough = RtdbError(C.RtE_EQUATION_CLUTTER_POOL_NOT_ENOUGH)

	// RteNamedTypeNameLenError  自定义类型的名称过长
	RteNamedTypeNameLenError = RtdbError(C.RtE_NAMED_TYPE_NAME_LEN_ERROR)

	// RteNamedTypeLengthNotMatch  数值长度与自定义类型的定义不符
	RteNamedTypeLengthNotMatch = RtdbError(C.RtE_NAMED_TYPE_LENGTH_NOT_MATCH)

	// RteCanNotUpdateSummary  无法更新卫星数据
	RteCanNotUpdateSummary = RtdbError(C.RtE_CAN_NOT_UPDATE_SUMMARY)

	// RteTooManyArvexFile  附属文件太多，无法继续创建附属文件
	RteTooManyArvexFile = RtdbError(C.RtE_TOO_MANY_ARVEX_FILE)

	// RteNotSupportedFeature  测试版本，暂时不支持此功能
	RteNotSupportedFeature = RtdbError(C.RtE_NOT_SUPPORTED_FEATURE)

	// RteEnsureError  验证信息失败，详细信息请查看数据库日志
	RteEnsureError = RtdbError(C.RtE_ENSURE_ERROR)

	// RteOperatorIsCancel  操作被取消
	RteOperatorIsCancel = RtdbError(C.RtE_OPERATOR_IS_CANCEL)

	// RteMsgbodyRevError  报文体接收失败
	RteMsgbodyRevError = RtdbError(C.RtE_MSGBODY_REV_ERROR)

	// RteUncompressFailed  解压缩失败
	RteUncompressFailed = RtdbError(C.RtE_UNCOMPRESS_FAILED)

	// RteCompressFailed  压缩失败
	RteCompressFailed = RtdbError(C.RtE_COMPRESS_FAILED)

	// RteSubscribeError  订阅失败，前一个订阅线程尚未退出
	RteSubscribeError = RtdbError(C.RtE_SUBSCRIBE_ERROR)

	// RteSubscribeCancelError  取消订阅失败
	RteSubscribeCancelError = RtdbError(C.RtE_SUBSCRIBE_CANCEL_ERROR)

	// RteSubscribeCallbackFailed  订阅回掉函数中不能调用取消订阅、断开连接
	RteSubscribeCallbackFailed = RtdbError(C.RtE_SUBSCRIBE_CALLBACK_FAILED)

	// RteSubscribeGreaterMaxCount  超过单连接可订阅标签点数量
	RteSubscribeGreaterMaxCount = RtdbError(C.RtE_SUBSCRIBE_GREATER_MAX_COUNT)

	// RteKillConnectionFailed  断开连接失败，无法断开自身连接
	RteKillConnectionFailed = RtdbError(C.RtE_KILL_CONNECTION_FAILED)

	// RteSubscribeNotMatch  请求的方法与当前的订阅不匹配
	RteSubscribeNotMatch = RtdbError(C.RtE_SUBSCRIBE_NOT_MATCH)

	// RteNoSubscribe  连接还未发起订阅，或者标签点还未订阅
	RteNoSubscribe = RtdbError(C.RtE_NO_SUBSCRIBE)

	// RteAlreadySubscribe  标签点已经被订阅
	RteAlreadySubscribe = RtdbError(C.RtE_ALREADY_SUBSCRIBE)

	// RteCalcPointUnsupportedWriteData  计算点不支持写入数据
	RteCalcPointUnsupportedWriteData = RtdbError(C.RtE_CALC_POINT_UNSUPPORTED_WRITE_DATA)

	// RteFeatureDeprecated  不再支持此功能
	RteFeatureDeprecated = RtdbError(C.RtE_FEATURE_DEPRECATED)

	// RteInvalidValue  无效的数据
	RteInvalidValue = RtdbError(C.RtE_INVALID_VALUE)

	// RteVerifyVercodeFailed  验证授权码失败
	RteVerifyVercodeFailed = RtdbError(C.RtE_VERIFY_VERCODE_FAILED)

	// RteInvalidPageSize  无效的数据页的大小
	RteInvalidPageSize = RtdbError(C.RtE_INVALID_PAGE_SIZE)

	// RteInvalidPrecision  无效的时间戳精度
	RteInvalidPrecision = RtdbError(C.RtE_INVALID_PRECISION)

	// RteInvalidPageVersion  无效的数据页版本
	RteInvalidPageVersion = RtdbError(C.RtE_INVALID_PAGE_VERSION)

	// RtePageIsFull  数据页已满
	RtePageIsFull = RtdbError(C.RtE_PAGE_IS_FULL)

	// RtePageNotLoaded  还未加载数据页
	RtePageNotLoaded = RtdbError(C.RtE_PAGE_NOT_LOADED)

	// RtePageAlreadyLoaded  已经加载了数据页
	RtePageAlreadyLoaded = RtdbError(C.RtE_PAGE_ALREADY_LOADED)

	// RtePageTooSmall  数据页太小，有效空间小于数据长度
	RtePageTooSmall = RtdbError(C.RtE_PAGE_TOO_SMALL)

	// RtePageNoEnoughData  数据页中没有足够的数据
	RtePageNoEnoughData = RtdbError(C.RtE_PAGE_NO_ENOUGH_DATA)

	// RtePageInsertFailed  数据页插入数据失败
	RtePageInsertFailed = RtdbError(C.RtE_PAGE_INSERT_FAILED)

	// RtePageNoEnoughSpace  数据页没有足够的空间
	RtePageNoEnoughSpace = RtdbError(C.RtE_PAGE_NO_ENOUGH_SPACE)

	// RteModifingMetaData  正在修改元数据，请稍后再试
	RteModifingMetaData = RtdbError(C.RtE_MODIFING_META_DATA)

	// RtePageSizeNotMatch  数据页大小不匹配
	RtePageSizeNotMatch = RtdbError(C.RtE_PAGE_SIZE_NOT_MATCH)

	// RteSyncBegin  元数据同步错误码起始值
	RteSyncBegin = RtdbError(C.RtE_SYNC_BEGIN)

	// RteSyncInvalidConfig  元数据同步-无效的配置
	RteSyncInvalidConfig = RtdbError(C.RtE_SYNC_INVALID_CONFIG)

	// RteSyncInvalidVersion  元数据同步-无效的版本号
	RteSyncInvalidVersion = RtdbError(C.RtE_SYNC_INVALID_VERSION)

	// RteSyncConfirmExpired  元数据同步-等待确认信息过期
	RteSyncConfirmExpired = RtdbError(C.RtE_SYNC_CONFIRM_EXPIRED)

	// RteSyncTooManyFwdinfo  元数据同步-转发信息过多
	RteSyncTooManyFwdinfo = RtdbError(C.RtE_SYNC_TOO_MANY_FWDINFO)

	// RteSyncNotMaster  元数据同步-不是主库
	RteSyncNotMaster = RtdbError(C.RtE_SYNC_NOT_MASTER)

	// RteSyncSyncing  元数据同步-正在同步
	RteSyncSyncing = RtdbError(C.RtE_SYNC_SYNCING)

	// RteSyncUnsynced  元数据同步-未同步
	RteSyncUnsynced = RtdbError(C.RtE_SYNC_UNSYNCED)

	// RteSyncTablePosConflict  元数据同步-表位置冲突
	RteSyncTablePosConflict = RtdbError(C.RtE_SYNC_TABLE_POS_CONFLICT)

	// RteSyncInvalidPointId  元数据同步-无效的标签点ID
	RteSyncInvalidPointId = RtdbError(C.RtE_SYNC_INVALID_POINT_ID)

	// RteSyncInvalidTableId  元数据同步-无效的表ID
	RteSyncInvalidTableId = RtdbError(C.RtE_SYNC_INVALID_TABLE_ID)

	// RteSyncInvalidNamedTypeId  元数据同步-无效的自定义类型ID
	RteSyncInvalidNamedTypeId = RtdbError(C.RtE_SYNC_INVALID_NAMED_TYPE_ID)

	// RteSyncRestoring  元数据同步-正在重建元数据
	RteSyncRestoring = RtdbError(C.RtE_SYNC_RESTORING)

	// RteSyncServerIsNotRunning  元数据同步-网络服务不是运行状态
	RteSyncServerIsNotRunning = RtdbError(C.RtE_SYNC_SERVER_IS_NOT_RUNNING)

	// RteSyncWriteWalFailed  元数据同步-写WAL失败
	RteSyncWriteWalFailed = RtdbError(C.RtE_SYNC_WRITE_WAL_FAILED)

	// RteSyncEnd  元数据同步错误码结束值
	RteSyncEnd = RtdbError(C.RtE_SYNC_END)

	// RteNetError  网络错误的起始值
	RteNetError = RtdbError(C.RtE_NET_ERROR)

	// RteSockWsaeintr  （阻塞）调用被 WSACancelBlockingCall() 函数取消
	RteSockWsaeintr = RtdbError(C.RtE_SOCK_WSAEINTR)

	// RteSockWsaeacces  请求地址是广播地址，但是相应的 flags 没设置
	RteSockWsaeacces = RtdbError(C.RtE_SOCK_WSAEACCES)

	// RteSockWsaefault  非法内存访问
	RteSockWsaefault = RtdbError(C.RtE_SOCK_WSAEFAULT)

	// RteSockWsaemfile  无多余的描述符可用
	RteSockWsaemfile = RtdbError(C.RtE_SOCK_WSAEMFILE)

	// RteSockWsaewouldblock  套接字被标识为非阻塞，但操作将被阻塞
	RteSockWsaewouldblock = RtdbError(C.RtE_SOCK_WSAEWOULDBLOCK)

	// RteSockWsaeinprogress  一个阻塞的 Windows Sockets 操作正在进行
	RteSockWsaeinprogress = RtdbError(C.RtE_SOCK_WSAEINPROGRESS)

	// RteSockWsaealready  一个非阻塞的 connect() 调用已经在指定的套接字上进行
	RteSockWsaealready = RtdbError(C.RtE_SOCK_WSAEALREADY)

	// RteSockWsaenotsock  描述符不是套接字描述符
	RteSockWsaenotsock = RtdbError(C.RtE_SOCK_WSAENOTSOCK)

	// RteSockWsaedestaddrreq  要求（未指定）目的地址
	RteSockWsaedestaddrreq = RtdbError(C.RtE_SOCK_WSAEDESTADDRREQ)

	// RteSockWsaemsgsize  套接字为基于消息的，消息太大（大于底层传输支持的最大值）
	RteSockWsaemsgsize = RtdbError(C.RtE_SOCK_WSAEMSGSIZE)

	// RteSockWsaeprototype  对此套接字来说，指定协议是错误的类型
	RteSockWsaeprototype = RtdbError(C.RtE_SOCK_WSAEPROTOTYPE)

	// RteSockWsaeprotonosupport  不支持指定协议
	RteSockWsaeprotonosupport = RtdbError(C.RtE_SOCK_WSAEPROTONOSUPPORT)

	// RteSockWsaesocktnosupport  在此地址族中不支持指定套接字类型
	RteSockWsaesocktnosupport = RtdbError(C.RtE_SOCK_WSAESOCKTNOSUPPORT)

	// RteSockWsaeopnotsupp  MSG_OOB 被指定，但是套接字不是流风格的
	RteSockWsaeopnotsupp = RtdbError(C.RtE_SOCK_WSAEOPNOTSUPP)

	// RteSockWsaeafnosupport  不支持指定的地址族
	RteSockWsaeafnosupport = RtdbError(C.RtE_SOCK_WSAEAFNOSUPPORT)

	// RteSockWsaeaddrinuse  套接字的本地地址已被使用
	RteSockWsaeaddrinuse = RtdbError(C.RtE_SOCK_WSAEADDRINUSE)

	// RteSockWsaeaddrnotavail  远程地址非法
	RteSockWsaeaddrnotavail = RtdbError(C.RtE_SOCK_WSAEADDRNOTAVAIL)

	// RteSockWsaenetdown  Windows Sockets 检测到网络系统已经失效
	RteSockWsaenetdown = RtdbError(C.RtE_SOCK_WSAENETDOWN)

	// RteSockWsaenetunreach  网络无法到达主机
	RteSockWsaenetunreach = RtdbError(C.RtE_SOCK_WSAENETUNREACH)

	// RteSockWsaenetreset  在操作进行时 keep-alive 活动检测到一个失败，连接被中断
	RteSockWsaenetreset = RtdbError(C.RtE_SOCK_WSAENETRESET)

	// RteSockWsaeconnaborted  连接因超时或其他失败而中断
	RteSockWsaeconnaborted = RtdbError(C.RtE_SOCK_WSAECONNABORTED)

	// RteSockWsaeconnreset  连接被复位
	RteSockWsaeconnreset = RtdbError(C.RtE_SOCK_WSAECONNRESET)

	// RteSockWsaenobufs  无缓冲区空间可用
	RteSockWsaenobufs = RtdbError(C.RtE_SOCK_WSAENOBUFS)

	// RteSockWsaeisconn  连接已建立
	RteSockWsaeisconn = RtdbError(C.RtE_SOCK_WSAEISCONN)

	// RteSockWsaenotconn  套接字未建立连接
	RteSockWsaenotconn = RtdbError(C.RtE_SOCK_WSAENOTCONN)

	// RteSockWsaeshutdown  套接字已 shutdown，连接已断开
	RteSockWsaeshutdown = RtdbError(C.RtE_SOCK_WSAESHUTDOWN)

	// RteSockWsaetimedout  连接请求超时，未能建立连接
	RteSockWsaetimedout = RtdbError(C.RtE_SOCK_WSAETIMEDOUT)

	// RteSockWsaeconnrefused  连接被拒绝
	RteSockWsaeconnrefused = RtdbError(C.RtE_SOCK_WSAECONNREFUSED)

	// RteSockWsaeclose  连接被关闭
	RteSockWsaeclose = RtdbError(C.RtE_SOCK_WSAECLOSE)

	// RteSockWsanotinitialised  Windows Sockets DLL 未初始化
	RteSockWsanotinitialised = RtdbError(C.RtE_SOCK_WSANOTINITIALISED)

	// RteCErrnoError  C语言errno错误的起始值
	RteCErrnoError = RtdbError(C.RtE_C_ERRNO_ERROR)

	// RteCErrnoEperm  Operation not permitted
	RteCErrnoEperm = RtdbError(C.RtE_C_ERRNO_EPERM)

	// RteCErrnoEnoent  No such file or directory
	RteCErrnoEnoent = RtdbError(C.RtE_C_ERRNO_ENOENT)

	// RteCErrnoEsrch  No such process
	RteCErrnoEsrch = RtdbError(C.RtE_C_ERRNO_ESRCH)

	// RteCErrnoEintr  Interrupted system call
	RteCErrnoEintr = RtdbError(C.RtE_C_ERRNO_EINTR)

	// RteCErrnoEio  I/O error
	RteCErrnoEio = RtdbError(C.RtE_C_ERRNO_EIO)

	// RteCErrnoEnxio  No such device or address
	RteCErrnoEnxio = RtdbError(C.RtE_C_ERRNO_ENXIO)

	// RteCErrnoE2big  Argument list too long
	RteCErrnoE2big = RtdbError(C.RtE_C_ERRNO_E2BIG)

	// RteCErrnoEnoexec  Exec format error
	RteCErrnoEnoexec = RtdbError(C.RtE_C_ERRNO_ENOEXEC)

	// RteCErrnoEbadf  Bad file number
	RteCErrnoEbadf = RtdbError(C.RtE_C_ERRNO_EBADF)

	// RteCErrnoEchild  No child processes
	RteCErrnoEchild = RtdbError(C.RtE_C_ERRNO_ECHILD)

	// RteCErrnoEagain  Try again
	RteCErrnoEagain = RtdbError(C.RtE_C_ERRNO_EAGAIN)

	// RteCErrnoEnomem  Out of memory
	RteCErrnoEnomem = RtdbError(C.RtE_C_ERRNO_ENOMEM)

	// RteCErrnoEacces  Permission denied
	RteCErrnoEacces = RtdbError(C.RtE_C_ERRNO_EACCES)

	// RteCErrnoEfault  Bad address
	RteCErrnoEfault = RtdbError(C.RtE_C_ERRNO_EFAULT)

	// RteCErrnoEnotblk  Block device required
	RteCErrnoEnotblk = RtdbError(C.RtE_C_ERRNO_ENOTBLK)

	// RteCErrnoEbusy  Device or resource busy
	RteCErrnoEbusy = RtdbError(C.RtE_C_ERRNO_EBUSY)

	// RteCErrnoEexist  File exists
	RteCErrnoEexist = RtdbError(C.RtE_C_ERRNO_EEXIST)

	// RteCErrnoExdev  Cross-device link
	RteCErrnoExdev = RtdbError(C.RtE_C_ERRNO_EXDEV)

	// RteCErrnoEnodev  No such device
	RteCErrnoEnodev = RtdbError(C.RtE_C_ERRNO_ENODEV)

	// RteCErrnoEnotdir  Not a directory
	RteCErrnoEnotdir = RtdbError(C.RtE_C_ERRNO_ENOTDIR)

	// RteCErrnoEisdir  Is a directory
	RteCErrnoEisdir = RtdbError(C.RtE_C_ERRNO_EISDIR)

	// RteCErrnoEinval  Invalid argument
	RteCErrnoEinval = RtdbError(C.RtE_C_ERRNO_EINVAL)

	// RteCErrnoEnfile  File table overflow
	RteCErrnoEnfile = RtdbError(C.RtE_C_ERRNO_ENFILE)

	// RteCErrnoEmfile  Too many open files
	RteCErrnoEmfile = RtdbError(C.RtE_C_ERRNO_EMFILE)

	// RteCErrnoEnotty  Not a typewriter
	RteCErrnoEnotty = RtdbError(C.RtE_C_ERRNO_ENOTTY)

	// RteCErrnoEtxtbsy  Text file busy
	RteCErrnoEtxtbsy = RtdbError(C.RtE_C_ERRNO_ETXTBSY)

	// RteCErrnoEfbig  File too large
	RteCErrnoEfbig = RtdbError(C.RtE_C_ERRNO_EFBIG)

	// RteCErrnoEnospc  No space left on device
	RteCErrnoEnospc = RtdbError(C.RtE_C_ERRNO_ENOSPC)

	// RteCErrnoEspipe  Illegal seek
	RteCErrnoEspipe = RtdbError(C.RtE_C_ERRNO_ESPIPE)

	// RteCErrnoErofs  Read-only file system
	RteCErrnoErofs = RtdbError(C.RtE_C_ERRNO_EROFS)

	// RteCErrnoEmlink  Too many links
	RteCErrnoEmlink = RtdbError(C.RtE_C_ERRNO_EMLINK)

	// RteCErrnoEpipe  Broken pipe
	RteCErrnoEpipe = RtdbError(C.RtE_C_ERRNO_EPIPE)

	// RteCErrnoEdom  Math argument out of domain of func
	RteCErrnoEdom = RtdbError(C.RtE_C_ERRNO_EDOM)

	// RteCErrnoErange  Math result not representable
	RteCErrnoErange = RtdbError(C.RtE_C_ERRNO_ERANGE)

	// RteCErrnoEdeadlk  Resource deadlock would occur
	RteCErrnoEdeadlk = RtdbError(C.RtE_C_ERRNO_EDEADLK)

	// RteCErrnoEnametoolong  File name too long
	RteCErrnoEnametoolong = RtdbError(C.RtE_C_ERRNO_ENAMETOOLONG)

	// RteCErrnoEnolck  No record locks available
	RteCErrnoEnolck = RtdbError(C.RtE_C_ERRNO_ENOLCK)

	// RteCErrnoEnosys  Function not implemented
	RteCErrnoEnosys = RtdbError(C.RtE_C_ERRNO_ENOSYS)

	// RteCErrnoEnotempty  Directory not empty
	RteCErrnoEnotempty = RtdbError(C.RtE_C_ERRNO_ENOTEMPTY)

	// RteCErrnoEloop  Too many symbolic links encountered
	RteCErrnoEloop = RtdbError(C.RtE_C_ERRNO_ELOOP)

	// RteCErrnoEnomsg  No message of desired type
	RteCErrnoEnomsg = RtdbError(C.RtE_C_ERRNO_ENOMSG)

	// RteCErrnoEidrm  Identifier removed
	RteCErrnoEidrm = RtdbError(C.RtE_C_ERRNO_EIDRM)

	// RteCErrnoEchrng  Channel number out of range
	RteCErrnoEchrng = RtdbError(C.RtE_C_ERRNO_ECHRNG)

	// RteCErrnoEl2nsync  Level 2 not synchronized
	RteCErrnoEl2nsync = RtdbError(C.RtE_C_ERRNO_EL2NSYNC)

	// RteCErrnoEl3hlt  Level 3 halted
	RteCErrnoEl3hlt = RtdbError(C.RtE_C_ERRNO_EL3HLT)

	// RteCErrnoEl3rst  Level 3 reset
	RteCErrnoEl3rst = RtdbError(C.RtE_C_ERRNO_EL3RST)

	// RteCErrnoElnrng  Link number out of range
	RteCErrnoElnrng = RtdbError(C.RtE_C_ERRNO_ELNRNG)

	// RteCErrnoEunatch  Protocol driver not attached
	RteCErrnoEunatch = RtdbError(C.RtE_C_ERRNO_EUNATCH)

	// RteCErrnoEnocsi  No CSI structure available
	RteCErrnoEnocsi = RtdbError(C.RtE_C_ERRNO_ENOCSI)

	// RteCErrnoEl2hlt  Level 2 halted
	RteCErrnoEl2hlt = RtdbError(C.RtE_C_ERRNO_EL2HLT)

	// RteCErrnoEbade  Invalid exchange
	RteCErrnoEbade = RtdbError(C.RtE_C_ERRNO_EBADE)

	// RteCErrnoEbadr  Invalid request descriptor
	RteCErrnoEbadr = RtdbError(C.RtE_C_ERRNO_EBADR)

	// RteCErrnoExfull  Exchange full
	RteCErrnoExfull = RtdbError(C.RtE_C_ERRNO_EXFULL)

	// RteCErrnoEnoano  No anode
	RteCErrnoEnoano = RtdbError(C.RtE_C_ERRNO_ENOANO)

	// RteCErrnoEbadrqc  Invalid request code
	RteCErrnoEbadrqc = RtdbError(C.RtE_C_ERRNO_EBADRQC)

	// RteCErrnoEbadslt  Invalid slot
	RteCErrnoEbadslt = RtdbError(C.RtE_C_ERRNO_EBADSLT)

	// RteCErrnoEbfont  Bad font file format
	RteCErrnoEbfont = RtdbError(C.RtE_C_ERRNO_EBFONT)

	// RteCErrnoEnostr  Device not a stream
	RteCErrnoEnostr = RtdbError(C.RtE_C_ERRNO_ENOSTR)

	// RteCErrnoEnodata  No data available
	RteCErrnoEnodata = RtdbError(C.RtE_C_ERRNO_ENODATA)

	// RteCErrnoEtime  Timer expired
	RteCErrnoEtime = RtdbError(C.RtE_C_ERRNO_ETIME)

	// RteCErrnoEnosr  Out of streams resources
	RteCErrnoEnosr = RtdbError(C.RtE_C_ERRNO_ENOSR)

	// RteCErrnoEnonet  Machine is not on the network
	RteCErrnoEnonet = RtdbError(C.RtE_C_ERRNO_ENONET)

	// RteCErrnoEnopkg  Package not installed
	RteCErrnoEnopkg = RtdbError(C.RtE_C_ERRNO_ENOPKG)

	// RteCErrnoEremote  Object is remote
	RteCErrnoEremote = RtdbError(C.RtE_C_ERRNO_EREMOTE)

	// RteCErrnoEnolink  Link has been severed
	RteCErrnoEnolink = RtdbError(C.RtE_C_ERRNO_ENOLINK)

	// RteCErrnoEadv  Advertise error
	RteCErrnoEadv = RtdbError(C.RtE_C_ERRNO_EADV)

	// RteCErrnoEsrmnt  Srmount error
	RteCErrnoEsrmnt = RtdbError(C.RtE_C_ERRNO_ESRMNT)

	// RteCErrnoEcomm  Communication error on send
	RteCErrnoEcomm = RtdbError(C.RtE_C_ERRNO_ECOMM)

	// RteCErrnoEproto  Protocol error
	RteCErrnoEproto = RtdbError(C.RtE_C_ERRNO_EPROTO)

	// RteCErrnoEmultihop  Multihop attempted
	RteCErrnoEmultihop = RtdbError(C.RtE_C_ERRNO_EMULTIHOP)

	// RteCErrnoEdotdot  RFS specific error
	RteCErrnoEdotdot = RtdbError(C.RtE_C_ERRNO_EDOTDOT)

	// RteCErrnoEbadmsg  Not a data message
	RteCErrnoEbadmsg = RtdbError(C.RtE_C_ERRNO_EBADMSG)

	// RteCErrnoEoverflow  Value too large for defined data type
	RteCErrnoEoverflow = RtdbError(C.RtE_C_ERRNO_EOVERFLOW)

	// RteCErrnoEnotuniq  Name not unique on network
	RteCErrnoEnotuniq = RtdbError(C.RtE_C_ERRNO_ENOTUNIQ)

	// RteCErrnoEbadfd  File descriptor in bad state
	RteCErrnoEbadfd = RtdbError(C.RtE_C_ERRNO_EBADFD)

	// RteCErrnoEremchg  Remote address changed
	RteCErrnoEremchg = RtdbError(C.RtE_C_ERRNO_EREMCHG)

	// RteCErrnoElibacc  Can not access a needed shared library
	RteCErrnoElibacc = RtdbError(C.RtE_C_ERRNO_ELIBACC)

	// RteCErrnoElibbad  Accessing a corrupted shared library
	RteCErrnoElibbad = RtdbError(C.RtE_C_ERRNO_ELIBBAD)

	// RteCErrnoElibscn  .lib section in a.out corrupted
	RteCErrnoElibscn = RtdbError(C.RtE_C_ERRNO_ELIBSCN)

	// RteCErrnoElibmax  Attempting to link in too many shared libraries
	RteCErrnoElibmax = RtdbError(C.RtE_C_ERRNO_ELIBMAX)

	// RteCErrnoElibexec  Cannot exec a shared library directly
	RteCErrnoElibexec = RtdbError(C.RtE_C_ERRNO_ELIBEXEC)

	// RteCErrnoEilseq  Illegal byte sequence
	RteCErrnoEilseq = RtdbError(C.RtE_C_ERRNO_EILSEQ)

	// RteCErrnoErestart  Interrupted system call should be restarted
	RteCErrnoErestart = RtdbError(C.RtE_C_ERRNO_ERESTART)

	// RteCErrnoEstrpipe  Streams pipe error
	RteCErrnoEstrpipe = RtdbError(C.RtE_C_ERRNO_ESTRPIPE)

	// RteCErrnoEusers  Too many users
	RteCErrnoEusers = RtdbError(C.RtE_C_ERRNO_EUSERS)

	// RteCErrnoEnotsock  Socket operation on non-socket
	RteCErrnoEnotsock = RtdbError(C.RtE_C_ERRNO_ENOTSOCK)

	// RteCErrnoEdestaddrreq  Destination address required
	RteCErrnoEdestaddrreq = RtdbError(C.RtE_C_ERRNO_EDESTADDRREQ)

	// RteCErrnoEmsgsize  Message too long
	RteCErrnoEmsgsize = RtdbError(C.RtE_C_ERRNO_EMSGSIZE)

	// RteCErrnoEprototype  Protocol wrong type for socket
	RteCErrnoEprototype = RtdbError(C.RtE_C_ERRNO_EPROTOTYPE)

	// RteCErrnoEnoprotoopt  Protocol not available
	RteCErrnoEnoprotoopt = RtdbError(C.RtE_C_ERRNO_ENOPROTOOPT)

	// RteCErrnoEprotonosupport  Protocol not supported
	RteCErrnoEprotonosupport = RtdbError(C.RtE_C_ERRNO_EPROTONOSUPPORT)

	// RteCErrnoEsocktnosupport  Socket type not supported
	RteCErrnoEsocktnosupport = RtdbError(C.RtE_C_ERRNO_ESOCKTNOSUPPORT)

	// RteCErrnoEopnotsupp  Operation not supported on transport endpoint
	RteCErrnoEopnotsupp = RtdbError(C.RtE_C_ERRNO_EOPNOTSUPP)

	// RteCErrnoEpfnosupport  Protocol family not supported
	RteCErrnoEpfnosupport = RtdbError(C.RtE_C_ERRNO_EPFNOSUPPORT)

	// RteCErrnoEafnosupport  Address family not supported by protocol
	RteCErrnoEafnosupport = RtdbError(C.RtE_C_ERRNO_EAFNOSUPPORT)

	// RteCErrnoEaddrinuse  Address already in use
	RteCErrnoEaddrinuse = RtdbError(C.RtE_C_ERRNO_EADDRINUSE)

	// RteCErrnoEaddrnotavail  Cannot assign requested address
	RteCErrnoEaddrnotavail = RtdbError(C.RtE_C_ERRNO_EADDRNOTAVAIL)

	// RteCErrnoEnetdown  Network is down
	RteCErrnoEnetdown = RtdbError(C.RtE_C_ERRNO_ENETDOWN)

	// RteCErrnoEnetunreach  Network is unreachable
	RteCErrnoEnetunreach = RtdbError(C.RtE_C_ERRNO_ENETUNREACH)

	// RteCErrnoEnetreset  Network dropped connection because of reset
	RteCErrnoEnetreset = RtdbError(C.RtE_C_ERRNO_ENETRESET)

	// RteCErrnoEconnaborted  Software caused connection abort
	RteCErrnoEconnaborted = RtdbError(C.RtE_C_ERRNO_ECONNABORTED)

	// RteCErrnoEconnreset  Connection reset by peer
	RteCErrnoEconnreset = RtdbError(C.RtE_C_ERRNO_ECONNRESET)

	// RteCErrnoEnobufs  No buffer space available
	RteCErrnoEnobufs = RtdbError(C.RtE_C_ERRNO_ENOBUFS)

	// RteCErrnoEisconn  Transport endpoint is already connected
	RteCErrnoEisconn = RtdbError(C.RtE_C_ERRNO_EISCONN)

	// RteCErrnoEnotconn  Transport endpoint is not connected
	RteCErrnoEnotconn = RtdbError(C.RtE_C_ERRNO_ENOTCONN)

	// RteCErrnoEshutdown  Cannot send after transport endpoint shutdown
	RteCErrnoEshutdown = RtdbError(C.RtE_C_ERRNO_ESHUTDOWN)

	// RteCErrnoEtoomanyrefs  Too many references: cannot splice
	RteCErrnoEtoomanyrefs = RtdbError(C.RtE_C_ERRNO_ETOOMANYREFS)

	// RteCErrnoEtimedout  Connection timed out
	RteCErrnoEtimedout = RtdbError(C.RtE_C_ERRNO_ETIMEDOUT)

	// RteCErrnoEconnrefused  Connection refused
	RteCErrnoEconnrefused = RtdbError(C.RtE_C_ERRNO_ECONNREFUSED)

	// RteCErrnoEhostdown  Host is down
	RteCErrnoEhostdown = RtdbError(C.RtE_C_ERRNO_EHOSTDOWN)

	// RteCErrnoEhostunreach  No route to host
	RteCErrnoEhostunreach = RtdbError(C.RtE_C_ERRNO_EHOSTUNREACH)

	// RteCErrnoEalready  Operation already in progress
	RteCErrnoEalready = RtdbError(C.RtE_C_ERRNO_EALREADY)

	// RteCErrnoEinprogress  Operation now in progress
	RteCErrnoEinprogress = RtdbError(C.RtE_C_ERRNO_EINPROGRESS)

	// RteCErrnoEstale  Stale file handle
	RteCErrnoEstale = RtdbError(C.RtE_C_ERRNO_ESTALE)

	// RteCErrnoEuclean  Structure needs cleaning
	RteCErrnoEuclean = RtdbError(C.RtE_C_ERRNO_EUCLEAN)

	// RteCErrnoEnotnam  Not a XENIX named type file
	RteCErrnoEnotnam = RtdbError(C.RtE_C_ERRNO_ENOTNAM)

	// RteCErrnoEnavail  No XENIX semaphores available
	RteCErrnoEnavail = RtdbError(C.RtE_C_ERRNO_ENAVAIL)

	// RteCErrnoEisnam  Is a named type file
	RteCErrnoEisnam = RtdbError(C.RtE_C_ERRNO_EISNAM)

	// RteCErrnoEremoteio  Remote I/O error
	RteCErrnoEremoteio = RtdbError(C.RtE_C_ERRNO_EREMOTEIO)

	// RteCErrnoEdquot  Quota exceeded
	RteCErrnoEdquot = RtdbError(C.RtE_C_ERRNO_EDQUOT)

	// RteCErrnoEnomedium  No medium found
	RteCErrnoEnomedium = RtdbError(C.RtE_C_ERRNO_ENOMEDIUM)

	// RteCErrnoEmediumtype  Wrong medium type
	RteCErrnoEmediumtype = RtdbError(C.RtE_C_ERRNO_EMEDIUMTYPE)

	// RteCErrnoEcanceled  Operation Canceled
	RteCErrnoEcanceled = RtdbError(C.RtE_C_ERRNO_ECANCELED)

	// RteCErrnoEnokey  Required key not available
	RteCErrnoEnokey = RtdbError(C.RtE_C_ERRNO_ENOKEY)

	// RteCErrnoEkeyexpired  Key has expired
	RteCErrnoEkeyexpired = RtdbError(C.RtE_C_ERRNO_EKEYEXPIRED)

	// RteCErrnoEkeyrevoked  Key has been revoked
	RteCErrnoEkeyrevoked = RtdbError(C.RtE_C_ERRNO_EKEYREVOKED)

	// RteCErrnoEkeyrejected  Key was rejected by service
	RteCErrnoEkeyrejected = RtdbError(C.RtE_C_ERRNO_EKEYREJECTED)

	// RteCErrnoEownerdead  Owner died
	RteCErrnoEownerdead = RtdbError(C.RtE_C_ERRNO_EOWNERDEAD)

	// RteCErrnoEnotrecoverable  State not recoverable
	RteCErrnoEnotrecoverable = RtdbError(C.RtE_C_ERRNO_ENOTRECOVERABLE)

	// RteCErrnoErfkill  Operation not possible due to RF-kill
	RteCErrnoErfkill = RtdbError(C.RtE_C_ERRNO_ERFKILL)

	// RteCErrnoEhwpoison  Memory page has hardware error
	RteCErrnoEhwpoison = RtdbError(C.RtE_C_ERRNO_EHWPOISON)

	// RteIpcError  ipc error begin
	RteIpcError = RtdbError(C.RtE_IPC_ERROR)

	// RteIpcErrorEnd  ipc error end
	RteIpcErrorEnd = RtdbError(C.RtE_IPC_ERROR_END)
)

// RtdbApiOption 用于设置API的工作模式的参数选项
// \see RtdbSetOptionWarp
type RtdbApiOption uint32

const (
	// RtdbApiAutoReconn api 在连接中断后是否自动重连, 0 不重连；1 重连。默认为 0 不重连
	RtdbApiAutoReconn = RtdbApiOption(C.RTDB_API_AUTO_RECONN)

	// RtdbApiConnTimeout api 连接超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为1000
	RtdbApiConnTimeout = RtdbApiOption(C.RTDB_API_CONN_TIMEOUT)

	// RtdbApiSendTimeout api 发送超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为1000
	RtdbApiSendTimeout = RtdbApiOption(C.RTDB_API_SEND_TIMEOUT)

	// RtdbApiRecvTimeout api 接收超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为60000
	RtdbApiRecvTimeout = RtdbApiOption(C.RTDB_API_RECV_TIMEOUT)

	// RtdbApiUserTimeout api TCP_USER_TIMEOUT超时值设置（单位：毫秒），默认为10000，Linux内核2.6.37以上有效
	RtdbApiUserTimeout = RtdbApiOption(C.RTDB_API_USER_TIMEOUT)

	// RtdbApiDefaultPrecision api 默认的时间戳精度，当使用旧版相关的api，以及新版api中未设置时间戳精度时，则使用此默认时间戳精度。 默认为毫秒精度
	RtdbApiDefaultPrecision = RtdbApiOption(C.RTDB_API_DEFAULT_PRECISION)

	// RtdbApiServerPrecision api 连接3.0数据库时，设置3.0数据库的时间戳精度，0表示毫秒精度，非0表示纳秒精度，默认为毫秒精度
	RtdbApiServerPrecision = RtdbApiOption(C.RTDB_API_SERVER_PRECISION)
)

// DatagramHandle 流句柄, 用于数据流订阅
type DatagramHandle struct {
	handle C.rtdb_datagram_handle
}

// RawRtdbGetApiVersionWarp 取得 rtdbapi 库的版本号
// \param [out]  major   主版本号
// \param [out]  minor   次版本号
// \param [out]  beta    发布版本号
// \return rtdb_error
// \remark 如果返回的版本号与 rtdb.h 中定义的不匹配(RTDB_API_XXX_VERSION)，则应用程序使用了错误的库。
//
//	应输出一条错误信息并退出，否则可能在调用某些 api 时会导致崩溃
func RawRtdbGetApiVersionWarp() (int32, int32, int32, error) {
	major, minor, beta := C.rtdb_int32(0), C.rtdb_int32(0), C.rtdb_int32(0)
	err := C.rtdb_get_api_version_warp(&major, &minor, &beta)
	return int32(major), int32(minor), int32(beta), RtdbError(err).GoError()
}

// RawRtdbSetOptionWarp 配置 api 行为参数，参见枚举 \ref RTDB_API_OPTION
// \param [in] type  选项类别
// \param [in] value 选项值
// \return rtdb_error
// \remark 选项设置后在下一次调用 api 时才生效
func RawRtdbSetOptionWarp(optionType RtdbApiOption, value int32) error {
	err := C.rtdb_set_option_warp(C.rtdb_int32(optionType), C.rtdb_int32(value))
	return RtdbError(err).GoError()
}

// RawRtdbCreateDatagramHandleWarp 创建数据流
// * \param [in] in 端口
// * \param [out] remotehost 对端地址
// * \param [out] handle 数据流句柄
// * \return rtdb_error
// * \remark 创建数据流 (备注：C代码中没文档，Go这边补的)
func RawRtdbCreateDatagramHandleWarp(port int32, remoteHost string) (DatagramHandle, error) {
	var handle C.rtdb_datagram_handle
	cRemoteHost := C.CString(remoteHost)
	defer C.free(unsafe.Pointer(cRemoteHost))
	err := C.rtdb_create_datagram_handle_warp(C.rtdb_int32(port), cRemoteHost, &handle)
	return DatagramHandle{handle: handle}, RtdbError(err).GoError()
}

// RawRtdbRemoveDatagramHandleWarp 删除数据流
// * \param [in] handle 数据流句柄
// * \return rtdb_error
// * \remark 删除数据流 (备注：C代码中没文档，Go这边补的)
func RawRtdbRemoveDatagramHandleWarp(handle DatagramHandle) error {
	err := C.rtdb_remove_datagram_handle_warp(handle.handle)
	return RtdbError(err).GoError()
}

// RawRtdbRecvDatagramWarp 接收数据流
// * \param [in] cacheLen 缓存大小(会分配一个[]byte)
// * \param [in] handle 数据流句柄
// * \param [in] remote_addr 对端地址
// * \param [in] timeout 超时时间
// * \return rtdb_error
// * \remark 接收数据流 (备注：C代码中没文档，Go这边补的)
func RawRtdbRecvDatagramWarp(cacheLen int32, handle DatagramHandle, remoteAddr string, timeout int32) ([]byte, error) {
	message := make([]byte, cacheLen)
	messageLen := C.rtdb_int32(cacheLen)
	cRemoteAddr := C.CString(remoteAddr)
	defer C.free(unsafe.Pointer(cRemoteAddr))
	err := C.rtdb_recv_datagram_warp((*C.char)(unsafe.Pointer(&message[0])), &messageLen, handle.handle, cRemoteAddr, C.rtdb_int32(timeout))
	return message[0:messageLen], RtdbError(err).GoError()
}

type ConnectHandle int32

// RawRtdbConnectWarp 建立同 RTDB 数据库的网络连接
// * \param [in] hostname     RTDB 数据平台服务器的网络地址或机器名
// * \param [in] port         连接断开，缺省值 6327
// * \param [out]  handle  连接句柄
// * \return rtdb_error
// * \remark 在调用所有的接口函数之前，必须先调用本函数建立同Rtdb服务器的连接
func RawRtdbConnectWarp(hostname string, port int32) (ConnectHandle, error) {
	cHostname := C.CString(hostname)
	defer C.free(unsafe.Pointer(cHostname))
	cPort := C.rtdb_int32(port)
	cHandle := C.rtdb_int32(0)
	err := C.rtdb_connect_warp(cHostname, cPort, &cHandle)
	return ConnectHandle(cHandle), RtdbError(err).GoError()
}

// RawRtdbConnectionCountWarp 获取 RTDB 服务器当前连接个数
// * \param [in] handle   连接句柄 参见 \ref rtdb_connect
// * \param [in] node_number   双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
// * \param [out]  count 返回当前主机的连接个数
// * \return rtdb_error
func RawRtdbConnectionCountWarp(handle ConnectHandle, nodeNumber int32) (int32, error) {
	count := C.rtdb_int32(0)
	err := C.rtdb_connection_count_warp(C.rtdb_int32(handle), C.rtdb_int32(nodeNumber), &count)
	return int32(count), RtdbError(err).GoError()
}

// RawRtdbGetConnectionsWarp 列出 RTDB 服务器的所有连接句柄
// * \param [in] handle       连接句柄
// * \param [in] node_number   双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
// * \param [out] sockets    整形数组，所有连接的套接字句柄
// * \param [in,out]  count   输入时表示sockets的长度，输出时表示返回的连接个数
// * \return rtdb_error
// * \remark 用户须保证分配给 sockets 的空间与 count 相符。如果输入的 count 小于输出的 count，则只返回部分连接
// rtdb_error RTDBAPI_CALLRULE rtdb_get_connections_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32 *sockets, rtdb_int32 *count)
func RawRtdbGetConnectionsWarp() {}

// RawRtdbGetOwnConnectionWarp 获取当前连接的socket句柄
// * \param [in] handle       连接句柄
// * \param [in] node_number   双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
// * \param [out] sockets    整形数组，所有连接的套接字句柄
// rtdb_error RTDBAPI_CALLRULE rtdb_get_own_connection_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32* socket)
func RawRtdbGetOwnConnectionWarp() {}

// RawRtdbGetConnectionInfoWarp 获取 RTDB 服务器指定连接的信息
// * \param [in] handle          连接句柄，参见 \ref rtdb_connect
// * \param [in] node_number   双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
// * \param [in] socket          指定的连接
// * \param [out] info          与连接相关的信息，参见 \ref RTDB_HOST_CONNECT_INFO
// * \return rtdb_error
// rtdb_error RTDBAPI_CALLRULE rtdb_get_connection_info_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32 socket, RTDB_HOST_CONNECT_INFO *info)
func RawRtdbGetConnectionInfoWarp() {}

// RawRtdbGetConnectionInfoIpv6Warp 获取 RTDB 服务器指定连接的ipv6版本
// * \param [in] handle          连接句柄，参见 \ref rtdb_connect
// * \param [in] node_number     双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP，双活模式仅支持ipv4
// * \param [in] socket          指定的连接
// * \param [out] info           与连接相关的信息，参见 \ref RTDB_HOST_CONNECT_INFO_IPV6
// * \return rtdb_error
// rtdb_error RTDBAPI_CALLRULE rtdb_get_connection_info_ipv6_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32 socket, RTDB_HOST_CONNECT_INFO_IPV6* info)
func RawRtdbGetConnectionInfoIpv6Warp() {}

// RawRtdbDisconnectWarp 断开同 RTDB 数据平台的连接
// * \param handle  连接句柄
// * \return rtdb_error
// * \remark 完成对 RTDB 的访问后调用本函数断开连接。连接一旦断开，则需要重新连接后才能调用其他的接口函数。
// rtdb_error RTDBAPI_CALLRULE rtdb_disconnect_warp(rtdb_int32 handle)
func RawRtdbDisconnectWarp() {}

// RawRtdbLoginWarp 以有效帐户登录
// * \param handle          连接句柄
// * \param user            登录帐户
// * \param password        帐户口令
// * \param [out] priv     账户权限， 枚举 \ref RTDB_PRIV_GROUP 的值之一
// * \return rtdb_error
// rtdb_error RTDBAPI_CALLRULE rtdb_login_warp(rtdb_int32 handle, const char *user, const char *password, rtdb_int32 *priv)
func RawRtdbLoginWarp() {}

// RawRTDBOSINVALID 获取连接句柄所连接的服务器操作系统类型
// * \param     handle          连接句柄
// * \param     ostype   操作系统类型 枚举 \ref RTDB_OS_TYPE 的值之一
// * \return    rtdb_error
// * \remark 如句柄未链接任何服务器，返回RTDB_OS_INVALID(当前支持操作系统类型：windows、linux)。
// rtdb_error RTDBAPI_CALLRULE rtdb_get_linked_ostype_warp(rtdb_int32 handle, RTDB_OS_TYPE* ostype)
func RawRTDBOSINVALID() {}

// RawRtdbChangePasswordWarp 修改用户帐户口令
// * \param handle    连接句柄
// * \param user      已有帐户
// * \param password  帐户新口令
// * \return rtdb_error
// * \remark 只有系统管理员可以修改其它用户的密码
// rtdb_error RTDBAPI_CALLRULE rtdb_change_password_warp(rtdb_int32 handle, const char *user, const char *password)
func RawRtdbChangePasswordWarp() {}

// RawRtdbChangeMyPasswordWarp 用户修改自己帐户口令
// * \param handle  连接句柄
// * \param old_pwd 帐户原口令
// * \param new_pwd 帐户新口令
// * \return rtdb_error
// rtdb_error RTDBAPI_CALLRULE rtdb_change_my_password_warp(rtdb_int32 handle, const char *old_pwd, const char *new_pwd)
func RawRtdbChangeMyPasswordWarp() {}

// RawRtdbGetPrivWarp 获取连接权限
// * \param handle          连接句柄
// * \param [out] priv  帐户权限， 枚举 \ref RTDB_PRIV_GROUP 的值之一
// * \return rtdb_error
// * \remark 如果还未登陆或不在服务器信任连接中，对应权限为-1，表示没有任何权限
// rtdb_error RTDBAPI_CALLRULE rtdb_get_priv_warp(rtdb_int32 handle, rtdb_int32 *priv)
func RawRtdbGetPrivWarp() {}

// RawRtdbChangePrivWarp 修改用户帐户权限
// * \param handle  连接句柄
// * \param user    已有帐户
// * \param priv    帐户权限， 枚举 \ref RTDB_PRIV_GROUP 的值之一
// * \return rtdb_error
// * \remark 只有管理员有修改权限
// rtdb_error RTDBAPI_CALLRULE rtdb_change_priv_warp(rtdb_int32 handle, const char *user, rtdb_int32 priv)
func RawRtdbChangePrivWarp() {}

// RawRtdbAddUserWarp 添加用户帐户
// * \param handle    连接句柄
// * \param user      帐户
// * \param password  帐户初始口令
// * \param priv      帐户权限， 枚举 \ref RTDB_PRIV_GROUP 的值之一
// * \return rtdb_error
// * \remark 只有管理员有添加用户权限
// rtdb_error RTDBAPI_CALLRULE rtdb_add_user_warp(rtdb_int32 handle, const char *user, const char *password, rtdb_int32 priv)
func RawRtdbAddUserWarp() {}

// RawRtdbRemoveUserWarp 删除用户帐户
// * \param handle  连接句柄
// * \param user    帐户
// * \return rtdb_error
// * \remark 只有管理员有删除用户权限
// rtdb_error RTDBAPI_CALLRULE rtdb_remove_user_warp(rtdb_int32 handle, const char *user)
func RawRtdbRemoveUserWarp() {}

// RawRtdbLockUserWarp 启用或禁用用户
// * \param     handle    连接句柄
// * \param     user      字符串，输入，帐户名
// * \param     lock      布尔，输入，是否禁用
// * \return    rtdb_error
// * \remark 只有管理员有启用禁用权限
// rtdb_error RTDBAPI_CALLRULE rtdb_lock_user_warp(rtdb_int32 handle, const char *user, rtdb_int8 lock)
func RawRtdbLockUserWarp() {}

// RawRtdbGetUsersWarp 获得所有用户
// * \param handle          连接句柄
// * \param [in,out]  count 输入时表示 users、privs 的长度，即用户个数；输出时表示成功返回的用户信息个数
// * \param [out] users     字符串指针数组，用户名称
// * \param [out] privs    整型数组，用户权限，枚举 \ref RTDB_PRIV_GROUP 的值之一
// * \return rtdb_error
// * \remark 用户须保证分配给 users, privs 的空间与 count 相符，如果输入的 count 小于总的用户数，则只返回部分用户信息。且每个指针指向的字符串缓冲区尺寸不小于 \ref RTDB_USER_SIZE。
// rtdb_error RTDBAPI_CALLRULE rtdb_get_users_warp(rtdb_int32 handle, rtdb_int32 *count, RTDB_USER_INFO *infos)
func RawRtdbGetUsersWarp() {}

// RawRtdbAddBlacklistWarp 添加连接黑名单项
// * \param handle  连接句柄
// * \param [in] addr    阻止连接段地址
// * \param [in] mask    阻止连接段子网掩码
// * \param [in] desc    阻止连接段的说明，超过 511 字符将被截断
// * \return rtdb_error
// * \remark addr 和 mask 进行与运算形成一个子网，
// * 来自该子网范围内的连接都将被阻止，黑名单的优先级高于信任连接，
// * 如果有连接既位于黑名单中，也位于信任连接中，则它将先被阻止。
// * 有效的子网掩码的所有 1 位于 0 左侧，例如："255.255.254.0"。
// * 当全部为 1 时，表示该子网中只有 addr 一个地址；但不能全部为 0。
// rtdb_error RTDBAPI_CALLRULE rtdb_add_blacklist_warp(rtdb_int32 handle, const char *addr, const char *mask, const char *desc)
func RawRtdbAddBlacklistWarp() {}

// RawRtdbUpdateBlacklistWarp 更新连接连接黑名单项
// * \param handle    连接句柄
// * \param addr      原阻止连接段地址
// * \param mask      原阻止连接段子网掩码
// * \param addr_new  新的阻止连接段地址
// * \param mask_new  新的阻止连接段子网掩码
// * \param desc      新的阻止连接段的说明，超过 511 字符将被截断
// rtdb_error RTDBAPI_CALLRULE rtdb_update_blacklist_warp(rtdb_int32 handle, const char *addr, const char *mask, const char *addr_new, const char *mask_new, const char *desc)
func RawRtdbUpdateBlacklistWarp() {}

// RawRtdbRemoveBlacklistWarp 删除连接黑名单项
// * \param handle  连接句柄
// * \param addr    阻止连接段地址
// * \param mask    阻止连接段子网掩码
// * \remark 只有 addr 与 mask 完全相同才视为同一个阻止连接段
// rtdb_error RTDBAPI_CALLRULE rtdb_remove_blacklist_warp(rtdb_int32 handle, const char *addr, const char *mask)
func RawRtdbRemoveBlacklistWarp() {}

// RawRtdbGetBlacklistWarp 获得连接黑名单
// * \param handle          连接句柄
// * \param addrs           字符串指针数组，输出，阻止连接段地址列表
// * \param masks           字符串指针数组，输出，阻止连接段子网掩码列表
// * \param descs           字符串指针数组，输出，阻止连接段的说明。
// * \param [in,out]  count 整型，输入/输出，用户个数
// * \remark 用户须保证分配给 addrs, masks, descs 的空间与 count 相符，
// * 如果输入的 count 小于输出的 count，则只返回部分阻止连接段，
// * addrs, masks 中每个字符串指针所指缓冲区尺寸不得小于 32 字节，
// * descs 中每个字符串指针所指缓冲区尺寸不得小于 512 字节。
// rtdb_error RTDBAPI_CALLRULE rtdb_get_blacklist_warp(rtdb_int32 handle, char* const* addrs, char* const* masks, char* const* descs, rtdb_int32 *count)
func RawRtdbGetBlacklistWarp() {}

// RawRtdbAddAuthorizationWarp 添加信任连接段
// * \param handle  连接句柄
// * \param addr    字符串，输入，信任连接段地址
// * \param mask    字符串，输入，信任连接段子网掩码。
// * \param priv    整数，输入，信任连接段拥有的用户权限。
// * \param desc    字符串，输入，信任连接段的说明，超过 511 字符将被截断。
// * \remark addr 和 mask 进行与运算形成一个子网，
// *        来自该子网范围内的连接都被视为可信任的，
// *        可以不用登录 (rtdb_login)，就直接调用 API，
// *        它所拥有的权限在 priv 中指定。
// *        有效的子网掩码的所有 1 位于 0 左侧，
// *        例如："255.255.254.0"。当全部为 1 时，
// *        表示该子网中只有 addr 一个地址；
// *        但不能全部为 0。
// rtdb_error RTDBAPI_CALLRULE rtdb_add_authorization_warp(rtdb_int32 handle, const char *addr, const char *mask, rtdb_int32 priv, const char *desc)
func RawRtdbAddAuthorizationWarp() {}

// RawRtdbUpdateAuthorizationWarp 更新信任连接段
// * \param handle    连接句柄
// * \param addr      字符串，输入，原信任连接段地址。
// * \param mask      字符串，输入，原信任连接段子网掩码。
// * \param addr_new  字符串，输入，新的信任连接段地址。
// * \param mask_new  字符串，输入，新的信任连接段子网掩码。
// * \param priv      整数，输入，新的信任连接段拥有的用户权限。
// * \param desc      字符串，输入，新的信任连接段的说明，超过 511 字符将被截断。
// rtdb_error RTDBAPI_CALLRULE rtdb_update_authorization_warp(rtdb_int32 handle, const char *addr, const char *mask, const char *addr_new, const char *mask_new, rtdb_int32 priv, const char *desc)
func RawRtdbUpdateAuthorizationWarp() {}

// RawRtdbRemoveAuthorizationWarp 删除信任连接段
// * \param handle  连接句柄
// * \param addr    字符串，输入，信任连接段地址。
// * \param mask    字符串，输入，信任连接段子网掩码。
// * \remark 只有 addr 与 mask 完全相同才视为同一个信任连接段
// rtdb_error RTDBAPI_CALLRULE rtdb_remove_authorization_warp(rtdb_int32 handle, const char *addr, const char *mask)
func RawRtdbRemoveAuthorizationWarp() {}

// RawRtdbGetAuthorizationsWarp 获得所有信任连接段
// * \param handle          连接句柄
// * \param addrs           字符串指针数组，输出，信任连接段地址列表
// * \param masks           字符串指针数组，输出，信任连接段子网掩码列表
// * \param [in,out]  privs 整型数组，输出，信任连接段权限列表
// * \param descs           字符串指针数组，输出，信任连接段的说明。
// * \param [in,out]  count 整型，输入/输出，用户个数
// * \remark 用户须保证分配给 addrs, masks, privs, descs 的空间与 count 相符，
// * 如果输入的 count 小于输出的 count，则只返回部分信任连接段，
// * addrs, masks 中每个字符串指针所指缓冲区尺寸不得小于 32 字节，
// * descs 中每个字符串指针所指缓冲区尺寸不得小于 512 字节。
// rtdb_error RTDBAPI_CALLRULE rtdb_get_authorizations_warp(rtdb_int32 handle, char* const* addrs, char* const* masks, rtdb_int32 *privs, char* const* descs, rtdb_int32 *count)
func RawRtdbGetAuthorizationsWarp() {}

// RawRtdbHostTimeWarp 获取 RTDB 服务器当前UTC时间
// * \param handle       连接句柄
// * \param hosttime     整型，输出，Rtdb服务器的当前UTC时间，
// * 表示距离1970年1月1日08:00:00的秒数。
// rtdb_error RTDBAPI_CALLRULE rtdb_host_time_warp(rtdb_int32 handle, rtdb_int32 *hosttime)
func RawRtdbHostTimeWarp() {}

// RawRtdbHostTime64Warp 获取 RTDB 服务器当前UTC时间
// * \param handle       连接句柄
// * \param hosttime     整型，输出，Rtdb服务器的当前UTC时间，
// * 表示距离1970年1月1日08:00:00的秒数。
// rtdb_error RTDBAPI_CALLRULE rtdb_host_time64_warp(rtdb_int32 handle, rtdb_timestamp_type* hosttime)
func RawRtdbHostTime64Warp() {}

// RawRtdbFormatTimespanWarp 根据时间跨度值生成时间格式字符串
// * \param str          字符串，输出，时间格式字符串，形如:
// * "1d" 表示时间跨度为24小时。
// * 具体含义参见 rtdb_parse_timespan 注释。
// * \param timespan     整型，输入，要处理的时间跨度秒数。
// * \remark 字符串缓冲区大小不应小于 32 字节。
// rtdb_error RTDBAPI_CALLRULE rtdb_format_timespan_warp(char *str, rtdb_int32 timespan)
func RawRtdbFormatTimespanWarp() {}

// RawRtdbParseTimespanWarp 根据时间格式字符串解析时间跨度值
// * \param str          字符串，输入，时间格式字符串，规则如下：
// *                     [time_span]
// *                     时间跨度部分可以出现多次，
// *                     可用的时间跨度代码及含义如下：
// *                     ?y            ?年, 1年 = 365日
// *                     ?m            ?月, 1月 = 30 日
// *                     ?d            ?日
// *                     ?h            ?小时
// *                     ?n            ?分钟
// *                     ?s            ?秒
// *                     例如："1d" 表示时间跨度为24小时。
// * \param timespan     整型，输出，返回解析得到的时间跨度秒数。
// rtdb_error RTDBAPI_CALLRULE rtdb_parse_timespan_warp(const char *str, rtdb_int32 *timespan)
func RawRtdbParseTimespanWarp() {}

// RawRtdbParseTimeWarp 根据时间格式字符串解析时间值
// * \param str          字符串，输入，时间格式字符串，规则如下：
// *                     base_time [+|- offset_time]
// *
// *                     其中 base_time 表示基本时间，有三种形式：
// *                     1. 时间字符串，如 "2010-1-1" 及 "2010-1-1 8:00:00"；
// *                     2. 时间代码，表示客户端相对时间；
// *                     可用的时间代码及含义如下：
// *                     td             当天零点
// *                     yd             昨天零点
// *                     tm             明天零点
// *                     mon            本周一零点
// *                     tue            本周二零点
// *                     wed            本周三零点
// *                     thu            本周四零点
// *                     fri            本周五零点
// *                     sat            本周六零点
// *                     sun            本周日零点
// *                     3. 星号('*')，表示客户端当前时间。
// *                     offset_time 是可选的，可以出现多次，
// *                     可用的时间偏移代码及含义如下：
// *                     [+|-] ?y            偏移?年, 1年 = 365日
// *                     [+|-] ?m            偏移?月, 1月 = 30 日
// *                     [+|-] ?d            偏移?日
// *                     [+|-] ?h            偏移?小时
// *                     [+|-] ?n            偏移?分钟
// *                     [+|-] ?s            偏移?秒
// *                     [+|-] ?ms           偏移?毫秒
// *                     例如："*-1d" 表示当前时刻减去24小时。
// * \param datetime     整型，输出，返回解析得到的时间值。
// * \param ms           短整型，输出，返回解析得到的时间毫秒值。
// *  备注：ms 可以为空指针，相应的毫秒信息将不再返回。
// rtdb_error RTDBAPI_CALLRULE rtdb_parse_time_warp(const char *str, rtdb_int64 *datetime, rtdb_int16 *ms)
func RawRtdbParseTimeWarp() {}

// RawRtdbFormatMessageWarp 获取 Rtdb API 调用返回值的简短描述
// * \param ecode        无符号整型，输入，Rtdb API调用后的返回值，详见rtdb_error.h头文件
// * \param message      字符串，输出，返回错误码简短描述
// * \param name         字符串，输出，返回错误码宏名称
// * \param size         整型，输入，message 参数的字节长度
// * \remark 用户须保证分配给 message， name 的空间与 size 相符,
// * name 或 message 可以为空指针，对应的信息将不再返回。
// void RTDBAPI_CALLRULE rtdb_format_message_warp(rtdb_error ecode, char *message, char *name, rtdb_int32 size)
func RawRtdbFormatMessageWarp() {}

// RawRtdbJobMessageWarp 获取任务的简短描述
// * \param job_id       整型，输入，RTDB_HOST_CONNECT_INFO::job 字段所表示的最近任务的描述
// * \param desc         字符串，输出，返回任务描述
// * \param name         字符串，输出，返回任务名称
// * \param size         整型，输入，desc、name 参数的字节长度
// * \remark 用户须保证分配给 desc、name 的空间与 size 相符，
// * name 或 message 可以为空指针，对应的信息将不再返回。
// void RTDBAPI_CALLRULE rtdb_job_message_warp(rtdb_int32 job_id, char *desc, char *name, rtdb_int32 size)
func RawRtdbJobMessageWarp() {}

// RawRtdbSetTimeoutWarp 设置连接超时时间
// * \param handle   连接句柄
// * \param socket   整型，输入，要设置超时时间的连接
// * \param timeout  整型，输入，超时时间，单位为秒，0 表示始终保持
// rtdb_error RTDBAPI_CALLRULE rtdb_set_timeout_warp(rtdb_int32 handle, rtdb_int32 socket, rtdb_int32 timeout)
func RawRtdbSetTimeoutWarp() {}

// RawRtdbGetTimeoutWarp 获得连接超时时间
// * \param handle   连接句柄
// * \param socket   整型，输入，要获取超时时间的连接
// * \param timeout  整型，输出，超时时间，单位为秒，0 表示始终保持
// rtdb_error RTDBAPI_CALLRULE rtdb_get_timeout_warp(rtdb_int32 handle, rtdb_int32 socket, rtdb_int32 *timeout)
func RawRtdbGetTimeoutWarp() {}

// RawRtdbKillConnectionWarp 断开已知连接
// * \param handle    连接句柄
// * \param socket    整型，输入，要断开的连接
// rtdb_error RTDBAPI_CALLRULE rtdb_kill_connection_warp(rtdb_int32 handle, rtdb_int32 socket)
func RawRtdbKillConnectionWarp() {}

// RawRtdbGetDbInfo1Warp 获得字符串型数据库系统参数
// * \param handle    连接句柄
// * \param index     整型，输入，要取得的参数索引，参见枚举 RTDB_DB_PARAM_INDEX。
// * \param str       字符串型，输出，存放取得的字符串参数值。
// * \param size      整型，输入，字符串缓冲区尺寸。
// * \remark 本接口只接受 [RTDB_PARAM_STR_FIRST, RTDB_PARAM_STR_LAST) 范围之内参数索引。
// rtdb_error RTDBAPI_CALLRULE rtdb_get_db_info1_warp(rtdb_int32 handle, rtdb_int32 index, char *str, rtdb_int32 size)
func RawRtdbGetDbInfo1Warp() {}

// RawRtdbGetDbInfo2Warp 获得整型数据库系统参数
// * \param handle    连接句柄
// * \param index     整型，输入，要取得的参数索引，参见枚举 RTDB_DB_PARAM_INDEX。
// * \param value     无符号整型，输出，存放取得的整型参数值。
// * \remark 本接口只接受 [RTDB_PARAM_INT_FIRST, RTDB_PARAM_INT_LAST) 范围之内参数索引。
// rtdb_error RTDBAPI_CALLRULE rtdb_get_db_info2_warp(rtdb_int32 handle, rtdb_int32 index, rtdb_uint32 *value)
func RawRtdbGetDbInfo2Warp() {}

// RawRtdbSetDbInfo1Warp 设置字符串型数据库系统参数
// * \param handle    连接句柄
// * \param index     整型，输入，要设置的参数索引，参见枚举 RTDB_DB_PARAM_INDEX。
// * 其中，仅以下列出的枚举值可用：
// * RTDB_PARAM_AUTO_BACKUP_PATH,
// * RTDB_PARAM_SERVER_SENDER_IP,
// * \param str       字符串型，输入，新的参数值。
// * \remark 如果修改了启动参数，将返回 RtE_DATABASE_NEED_RESTART 提示码。
// rtdb_error RTDBAPI_CALLRULE rtdb_set_db_info1_warp(rtdb_int32 handle, rtdb_int32 index, const char *str)
func RawRtdbSetDbInfo1Warp() {}

// RawRtdbSetDbInfo2Warp 设置整型数据库系统参数
// * \param handle    连接句柄
// * \param index     整型，输入，要取得的参数索引，参见枚举 RTDB_DB_PARAM_INDEX。
// * 其中，仅以下列出的枚举值可用：
// * RTDB_PARAM_SERVER_IPC_SIZE,
// * RTDB_PARAM_EQUATION_IPC_SIZE,
// * RTDB_PARAM_HASH_TABLE_SIZE,
// * RTDB_PARAM_TAG_DELETE_TIMES,
// * RTDB_PARAM_SERVER_PORT,
// * RTDB_PARAM_SERVER_SENDER_PORT,
// * RTDB_PARAM_SERVER_RECEIVER_PORT,
// * RTDB_PARAM_SERVER_MODE,
// * RTDB_PARAM_ARV_PAGES_NUMBER,
// * RTDB_PARAM_ARVEX_PAGES_NUMBER,
// * RTDB_PARAM_EXCEPTION_AT_SERVER,
// * RTDB_PARAM_EX_ARCHIVE_SIZE,
// * RTDB_PARAM_ARCHIVE_BATCH_SIZE,
// * RTDB_PARAM_ARV_ASYNC_QUEUE_SLOWER_DOOR,
// * RTDB_PARAM_ARV_ASYNC_QUEUE_NORMAL_DOOR,
// * RTDB_PARAM_INDEX_ALWAYS_IN_MEMORY,
// * RTDB_PARAM_DISK_MIN_REST_SIZE,
// * RTDB_PARAM_DELAY_OF_AUTO_MERGE_OR_ARRANGE,
// * RTDB_PARAM_START_OF_AUTO_MERGE_OR_ARRANGE,
// * RTDB_PARAM_STOP_OF_AUTO_MERGE_OR_ARRANGE,
// * RTDB_PARAM_START_OF_AUTO_BACKUP,
// * RTDB_PARAM_STOP_OF_AUTO_BACKUP,
// * RTDB_PARAM_MAX_LATENCY_OF_SNAPSHOT,
// * RTDB_PARAM_PAGE_ALLOCATOR_RESERVE_SIZE,
// * \param value     无符号整型，输入，新的参数值。
// * \remark 如果修改了启动参数，将返回 RtE_DATABASE_NEED_RESTART 提示码。
// rtdb_error RTDBAPI_CALLRULE rtdb_set_db_info2_warp(rtdb_int32 handle, rtdb_int32 index, rtdb_uint32 value)
func RawRtdbSetDbInfo2Warp() {}

// RawRtdbGetLogicalDriversWarp 获得逻辑盘符
// * \param handle     连接句柄
// * \param drivers    字符数组，输出，
// * 返回逻辑盘符组成的字符串，每个盘符占一个字符。
// * \remark drivers 的内存空间由用户负责维护，长度应不小于 32。
// rtdb_error RTDBAPI_CALLRULE rtdb_get_logical_drivers_warp(rtdb_int32 handle, char *drivers)
func RawRtdbGetLogicalDriversWarp() {}

// RawRtdbOpenPathWarp 打开目录以便遍历其中的文件和子目录。
// * \param handle       连接句柄
// * \param dir          字符串，输入，要打开的目录
// rtdb_error RTDBAPI_CALLRULE rtdb_open_path_warp(rtdb_int32 handle, const char *dir)
func RawRtdbOpenPathWarp() {}

// RawRtdbReadPathWarp 读取目录中的文件或子目录
// * \param handle      连接句柄
// * \param path        字符数组，输出，返回的文件、子目录全路径
// * \param is_dir      短整数，输出，返回 1 为目录，0 为文件
// * \param atime       整数，输出，为文件时，返回访问时间
// * \param ctime       整数，输出，为文件时，返回建立时间
// * \param mtime       整数，输出，为文件时，返回修改时间
// * \param size        64 位整数，输出，为文件时，返回文件大小
// * \remark path 的内存空间由用户负责维护，尺寸应不小于 RTDB_MAX_PATH。
// * 当返回值为 RtE_BATCH_END 时表示目录下所有子目录和文件已经遍历完毕。
// rtdb_error RTDBAPI_CALLRULE rtdb_read_path_warp(rtdb_int32 handle, char *path, rtdb_int16 *is_dir, rtdb_int32 *atime, rtdb_int32 *ctime, rtdb_int32 *mtime, rtdb_int64 *size)
func RawRtdbReadPathWarp() {}

// RawRtdbReadPath64Warp 读取目录中的文件或子目录
// * \param handle      连接句柄
// * \param path        字符数组，输出，返回的文件、子目录全路径
// * \param is_dir      短整数，输出，返回 1 为目录，0 为文件
// * \param atime       整数，输出，为文件时，返回访问时间
// * \param ctime       整数，输出，为文件时，返回建立时间
// * \param mtime       整数，输出，为文件时，返回修改时间
// * \param size        64 位整数，输出，为文件时，返回文件大小
// * \remark path 的内存空间由用户负责维护，尺寸应不小于 RTDB_MAX_PATH。
// * 当返回值为 RtE_BATCH_END 时表示目录下所有子目录和文件已经遍历完毕。
// rtdb_error RTDBAPI_CALLRULE rtdb_read_path64_warp(rtdb_int32 handle, char* path, rtdb_int16* is_dir, rtdb_timestamp_type* atime, rtdb_timestamp_type* ctime, rtdb_timestamp_type* mtime, rtdb_int64* size)
func RawRtdbReadPath64Warp() {}

// RawRtdbClosePathWarp 关闭当前遍历的目录
// * \param handle      连接句柄
// rtdb_error RTDBAPI_CALLRULE rtdb_close_path_warp(rtdb_int32 handle)
func RawRtdbClosePathWarp() {}

// RawRtdbMkdirWarp 建立目录
// * \param handle       连接句柄
// * \param dir          字符串，输入，新建目录的全路径
// rtdb_error RTDBAPI_CALLRULE rtdb_mkdir_warp(rtdb_int32 handle, const char *dir)
func RawRtdbMkdirWarp() {}

// RawRtdbGetFileSizeWarp 获得指定服务器端文件的大小
// * \param handle     连接句柄
// * \param file       字符串，输入，文件名
// * \param size       64 位整数，输出，文件大小
// rtdb_error RTDBAPI_CALLRULE rtdb_get_file_size_warp(rtdb_int32 handle, const char *file, rtdb_int64 *size)
func RawRtdbGetFileSizeWarp() {}

// RawRtdbReadFileWarp 读取服务器端指定文件的内容
// * \param handle       连接句柄
// * \param file         字符串，输入，要读取内容的文件名
// * \param content      字符数组，输出，文件内容
// * \param pos          64 位整型，输入，读取文件的起始位置
// * \param size         整型，输入/输出，
// *                     输入时表示要读取文件内容的字节大小；
// *                     输出时表示实际读取的字节数
// * \remark 用户须保证分配给 content 的空间与 size 相符。
// rtdb_error RTDBAPI_CALLRULE rtdb_read_file_warp(rtdb_int32 handle, const char *file, char *content, rtdb_int64 pos, rtdb_int64 *size)
func RawRtdbReadFileWarp() {}

// RawRtdbGetMaxBlobLenWarp 取得数据库允许的blob与str类型测点的最大长度
// * \param handle       连接句柄
// * \param len          整形，输出参数，代表数据库允许的blob、str类型测点的最大长度
// rtdb_error RTDBAPI_CALLRULE rtdb_get_max_blob_len_warp(rtdb_int32 handle, rtdb_int32 *len)
func RawRtdbGetMaxBlobLenWarp() {}

// RawRtdbFormatQualityWarp 取得质量码对应的定义
// * \param handle       连接句柄
// * \param count        质量码个数，输入参数，
// * \param qualities    质量码，输入参数
// * \param definitions  输出参数，0~255为RTDB质量码（参加rtdb.h文件），256~511为OPC质量码，大于511为用户自定义质量码
// * \param lens         输出参数，每个定义对应的长度
// * \remark OPC质量码把8位分3部分定义：XX XXXX XX，对应着：品质位域 分状态位域 限定位域
// * 品质位域：00（Bad）01（Uncertain）10（N/A）11（Good）
// * 分状态位域：不同品质位域对应各自的分状态位域
// * 限定位域：00（Not limited）01（Low limited）10（high limited）11（Constant）
// * 三个域之间用逗号隔开，输出到definitions参数中，前面有有RTDB，OPC或者USER标识，说明标签点类别
// rtdb_error RTDBAPI_CALLRULE rtdb_format_quality_warp(rtdb_int32 handle, rtdb_int32 *count, rtdb_int16 *qualities, rtdb_byte **definitions, rtdb_int32 *lens)
func RawRtdbFormatQualityWarp() {}

// RawRtdbJudgeConnectStatusWarp 判断连接是否可用
// * \param handle   连接句柄
// rtdb_error RTDBAPI_CALLRULE rtdb_judge_connect_status_warp(rtdb_int32 handle, rtdb_int8* change_connection GAPI_DEFAULT_VALUE(0), char* current_ip_addr GAPI_DEFAULT_VALUE(0), rtdb_int32 size GAPI_DEFAULT_VALUE(0))
func RawRtdbJudgeConnectStatusWarp() {}

// RawRtdbFormatIpaddrWarp 将整形IP转换为字符串形式的IP
// * [ip]        无符号整型，输入，整形的IP地址
// * [ip_addr]      字符串，输出，字符串IP地址缓冲区
// * [size]         整型，输入，ip_addr 参数的字节长度
// * 备注：用户须保证分配给 ip_addr 的空间与 size 相符
// void RTDBAPI_CALLRULE rtdb_format_ipaddr_warp(rtdb_uint32 ip, char* ip_addr, rtdb_int32 size)
func RawRtdbFormatIpaddrWarp() {}

// RawRtdbbGetEquationByFileNameWarp 根据文件名获取方程式
// *      [handle]   连接句柄
// *      [file_name] 输入，字符串，方程式路径
// *      [equation]  输出，返回的方程式长度最长为RTDB_MAX_EQUATION_SIZE-1
// *备注：用户调用时为equation分配的空间不得小于RTDB_MAX_EQUATION_SIZE
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_equation_by_file_name_warp(rtdb_int32 handle, const char* file_name, char equation[RTDB_MAX_EQUATION_SIZE])
func RawRtdbbGetEquationByFileNameWarp() {}

// /*
// *
//   - 命名：rtdbb_get_equation_by_id
//   - 功能：根ID径获取方程式
//   - 参数：
//   - [handle]   连接句柄
//   - [id]				输入，整型，方程式ID
//   - [equation] 输出，返回的方程式长度最长为RTDB_MAX_EQUATION_SIZE-1
//     *
//     *备注：用户调用时为equation分配的空间不得小于RTDB_MAX_EQUATION_SIZE
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_equation_by_id_warp(rtdb_int32 handle, rtdb_int32 id, char equation[RTDB_MAX_EQUATION_SIZE])
// */
func RawRtdbbGetEquationByIdWarp() {}

// /*
// *
//
//	*
//	* \brief 添加新表
//	*
//	* \param handle   连接句柄
//	* \param field    RTDB_TABLE 结构，输入/输出，表信息。
//	*                 在输入时，type、name、desc 字段有效；
//	*                 输出时，id 字段由系统自动分配并返回给用户。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_append_table_warp(rtdb_int32 handle, RTDB_TABLE *field)
// */
func RawRtdbbAppendTableWarp() {}

// /*
// *
//
//	*
//	* \brief 取得标签点表总数
//	*
//	* \param handle   连接句柄
//	* \param count    整型，输出，标签点表总数
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_tables_count_warp(rtdb_int32 handle, rtdb_int32 *count)
// */
func RawRtdbbTablesCountWarp() {}

// /*
// *
//
//	*
//	* \brief 取得所有标签点表的ID
//	*
//	* \param handle   连接句柄
//	* \param ids      整型数组，输出，标签点表的id
//	* \param count    整型，输入/输出，
//	*                 输入表示 ids 的长度，输出表示标签点表个数
//	* \remark 用户须保证分配给 ids 的空间与 count 相符
//	*      如果输入的 count 小于输出的 count，则只返回部分表id
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_tables_warp(rtdb_int32 handle, rtdb_int32 *ids, rtdb_int32 *count)
// */
func RawRtdbbGetTablesWarp() {}

// /*
// *
//
//	*
//	* \brief 根据表 id 获取表中包含的标签点数量
//	*
//	* \param handle   连接句柄
//	* \param id       整型，输入，表ID
//	* \param size     整型，输出，表中标签点数量
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_size_by_id_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 *size)
// */
func RawRtdbbGetTableSizeByIdWarp() {}

// /*
// *
// *
// * \brief 根据表名称获取表中包含的标签点数量
// *
// * \param handle   连接句柄
// * \param name     字符串，输入，表名称
// * \param size     整型，输出，表中标签点数量
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_size_by_name_warp(rtdb_int32 handle, const char *name, rtdb_int32 *size)
// */
func RawRtdbbGetTableSizeByNameWarp() {}

// /*
// *
//
//	*
//	* \brief 根据表 id 获取表中实际包含的标签点数量
//	*
//	* \param handle   连接句柄
//	* \param id       整型，输入，表ID
//	* \param size     整型，输出，表中标签点数量
//	* 注意：通过此API获取标签点数量，然后搜索此表中的标签点得到的数量可能会不一致，这是由于服务内部批量建点采取了异步的方式。
//	*       一般情况下请使用rtdbb_get_table_size_by_id来获取表中的标签点数量。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_real_size_by_id_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 *size)
// */
func RawRtdbbGetTableRealSizeByIdWarp() {}

// /*
// *
//
//	*
//	* \brief 根据标签点表 id 获取表属性
//	*
//	* \param handle 连接句柄
//	* \param field  RTDB_TABLE 结构，输入/输出，标签点表属性，
//	*               输入时指定 id 字段，输出时返回 type、name、desc 字段。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_property_by_id_warp(rtdb_int32 handle, RTDB_TABLE *field)
// */
func RawRtdbbGetTablePropertyByIdWarp() {}

// /*
// *
//
//	*
//	* \brief 根据表名获取标签点表属性
//	*
//	* \param handle 连接句柄
//	* \param field  RTDB_TABLE 结构，输入/输出，标签点表属性
//	*               输入时指定 name 字段，输出时返回 id、type、desc 字段。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_property_by_name_warp(rtdb_int32 handle, RTDB_TABLE *field)
// */
func RawRtdbbGetTablePropertyByNameWarp() {}

// /*
// *
//
//	*
//	* \brief 使用完整的属性集来创建单个标签点
//	*
//	* \param handle 连接句柄
//	* \param base RTDB_POINT 结构，输入/输出，
//	*      输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
//	* \param scan RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
//	* \param calc RTDB_CALC_POINT 结构，输入，计算标签点扩展属性集。
//	* \remark 如果新建的标签点没有对应的扩展属性集，可置为空指针。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_insert_point_warp(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_CALC_POINT *calc)
// */
func RawRtdbbInsertPointWarp() {}

// /*
// *
//   - 命名：rtdbb_insert_max_point
//   - 功能：使用最大长度的完整属性集来创建单个标签点
//   - 参数：
//   - [handle] 连接句柄
//   - [base] RTDB_POINT 结构，输入/输出，
//   - 输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
//   - [scan] RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
//   - [calc] RTDB_MAX_CALC_POINT 结构，输入，计算标签点扩展属性集。
//   - 备注：如果新建的标签点没有对应的扩展属性集，可置为空指针。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_insert_max_point_warp(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_MAX_CALC_POINT *calc)
// */
func RawRtdbbInsertMaxPointWarp() {}

// /*
// *
//   - 命名：rtdbb_insert_max_points
//   - 功能：使用最大长度的完整属性集来批量创建标签点
//   - 参数：
//   - [handle] 连接句柄
//   - [count] count, 输入，base/scan/calc数组的长度；输出，成功的个数
//   - [bases] RTDB_POINT 数组，输入/输出，
//   - 输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
//   - [scans] RTDB_SCAN_POINT 数组，输入，采集标签点扩展属性集。
//   - [calcs] RTDB_MAX_CALC_POINT 数组，输入，计算标签点扩展属性集。
//   - [errors] rtdb_error数组，输出，对应每个标签点的结果
//   - 备注：如果新建的标签点没有对应的扩展属性集，可置为空指针。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_insert_max_points_warp(rtdb_int32 handle, rtdb_int32* count, RTDB_POINT* bases, RTDB_SCAN_POINT* scans, RTDB_MAX_CALC_POINT* calcs, rtdb_error* errors)
// */
func RawRtdbbInsertMaxPointsWarp() {}

// /*
// *
// *
// * 功能  使用最小的属性集来创建单个标签点
// *
// * \param handle     连接句柄
// * \param tag        字符串，输入，标签点名称
// * \param type       整型，输入，标签点数据类型，取值 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、
// *                   RTDB_CHAR、RTDB_UINT16、RTDB_UINT32、RTDB_INT32、RTDB_INT64、
// *                   RTDB_REAL16、RTDB_REAL32、RTDB_REAL64、RTDB_COOR、RTDB_STRING、RTDB_BLOB 之一。
// * \param table_id   整型，输入，标签点所属表 id
// * \param use_ms     短整型，输入，标签点时间戳精度，0 为秒；1 为纳秒。
// * \param point_id   整型，输出，标签点 id
// * \remark 标签点的其余属性将取默认值。
// rtdb_error RTDBAPI_CALLRULE rtdbb_insert_base_point_warp(rtdb_int32 handle, const char *tag, rtdb_int32 type, rtdb_int32 table_id, rtdb_int16 use_ms, rtdb_int32 *point_id)
// */
func RawRtdbbInsertBasePointWarp() {}

// /*
// *
//   - 命名：rtdbb_insert_named_type_point
//   - 功能：使用完整的属性集来创建单个自定义数据类型标签点
//   - 参数：
//   - [handle] 连接句柄
//   - [base] RTDB_POINT 结构，输入/输出，
//   - 输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
//   - [scan] RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
//   - [name] 字符串，输入，自定义数据类型的名字。
//   - 备注：如果新建的标签点没有对应的扩展属性集，可置为空指针。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_insert_named_type_point_warp(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, const char* name)
// */
func RawRtdbbInsertNamedTypePointWarp() {}

// /*
// *
//
//	*
//	* \brief 根据 id 删除单个标签点
//	*
//	* \param handle 连接句柄
//	* \param id     整型，输入，标签点标识
//	* \remark 通过本接口删除的标签点为可回收标签点，
//	*        可以通过 rtdbb_recover_point 接口恢复。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_remove_point_by_id_warp(rtdb_int32 handle, rtdb_int32 id)
// */
func RawRtdbbRemovePointByIdWarp() {}

// /*
// *
// *
// * \brief 根据标签点全名删除单个标签点
// * \param handle        连接句柄
// * \param table_dot_tag  字符串，输入，标签点全名称："表名.标签点名"
// * \remark 通过本接口删除的标签点为可回收标签点，
// *        可以通过 rtdbb_recover_point 接口恢复。
// rtdb_error RTDBAPI_CALLRULE rtdbb_remove_point_by_name_warp(rtdb_int32 handle, const char *table_dot_tag)
// */
func RawRtdbbRemovePointByNameWarp() {}

// /*
// *
//   - 命名：rtdbb_move_point_by_id
//   - 功能：根据 id 移动单个标签点到其他表
//   - 参数：
//   - [handle] 连接句柄
//   - [id]     整型，输入，标签点标识
//   - [dest_table_name] 字符串，输入，移动的目标表名称
//   - 备注：通过本接口移动标签点后不改变标签点的id，且快照
//   - 和历史数据都不受影响
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_move_point_by_id_warp(rtdb_int32 handle, rtdb_int32 id, const char* dest_table_name)
// */
func RawRtdbbMovePointByIdWarp() {}

// /*
// *
//
//	*
//	* \brief 批量获取标签点属性
//	*
//	* \param handle 连接句柄
//	* \param count  整数，输入，表示标签点个数。
//	* \param base   RTDB_POINT 结构数组，输入/输出，标签点基本属性列表，
//	*                 输入时，id 字段指定需要得到属性的标签点，输出时，其它字段返回标签点属性值。
//	* \param scan   RTDB_SCAN_POINT 结构数组，输出，采集标签点扩展属性列表
//	* \param calc   RTDB_CALC_POINT 结构数组，输出，计算标签点扩展属性列表
//	* \param errors 无符号整型数组，输出，获取标签属性的返回值列表，参考rtdb_error.h
//	* \remark 用户须保证分配给 base、scan、calc、errors 的空间与 count 相符，
//	*        扩展属性集 scan、calc 可为空指针，此时将不返回对应的扩展属性集。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_points_property_warp(rtdb_int32 handle, rtdb_int32 count, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_CALC_POINT *calc, rtdb_error *errors)
// */
func RawRtdbbGetPointsPropertyWarp() {}

// /*
// *
//   - 命名：rtdbb_get_max_points_property
//   - 功能：按最大长度批量获取标签点属性
//   - 参数：
//   - [handle] 连接句柄
//   - [count]  整数，输入，表示标签点个数。
//   - [base]   RTDB_POINT 结构数组，输入/输出，标签点基本属性列表，
//   - 输入时，id 字段指定需要得到属性的标签点，输出时，其它字段返回标签点属性值。
//   - [scan]   RTDB_SCAN_POINT 结构数组，输出，采集标签点扩展属性列表
//   - [calc]   RTDB_MAX_CALC_POINT 结构数组，输出，计算标签点扩展属性列表
//   - [errors] 无符号整型数组，输出，获取标签属性的返回值列表，参考rtdb_error.h
//   - 备注：用户须保证分配给 base、scan、calc、errors 的空间与 count 相符，
//   - 扩展属性集 scan、calc 可为空指针，此时将不返回对应的扩展属性集。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_max_points_property_warp(rtdb_int32 handle, rtdb_int32 count, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_MAX_CALC_POINT *calc, rtdb_error *errors)
// */
func RawRtdbbGetMaxPointsPropertyWarp() {}

// /*
// *
// *
// * \brief 搜索符合条件的标签点，使用标签点名时支持通配符
// *
// * \param handle        连接句柄
// * \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
// * \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
// * \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
// *                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
// * \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
// *                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
// * \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
// *                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
// * \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
// * \param mode          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
// *                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
// * \param ids           整型数组，输出，返回搜索到的标签点标识列表
// * \param count         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
// * \remark 用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
// *        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
// *        如果 tagmask、tablemask 为空指针，则表示使用缺省设置"*",
// *        多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
// rtdb_error RTDBAPI_CALLRULE rtdbb_search_warp(rtdb_int32 handle, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
// */
func RawRtdbbSearchWarp() {}

// /*
// *
//
//	*
//	* \brief 分批继续搜索符合条件的标签点，使用标签点名时支持通配符
//	*
//	* \param handle        连接句柄
//	* \param start         整型，输入，搜索起始位置。
//	* \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//	* \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//	* \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
//	*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
//	* \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
//	*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
//	* \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
//	*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
//	* \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
//	* \param mode          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
//	*                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
//	* \param ids           整型数组，输出，返回搜索到的标签点标识列表
//	* \param count         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
//	* \remark 用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
//	*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
//	*        如果 tagmask、tablemask 为空指针，则表示使用缺省设置"*"。
//	*        当搜索到的标签点数比提供的要小时，表示这是最后一批符合条件的标签点 (即全部搜索完毕),
//	*        多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_search_in_batches_warp(rtdb_int32 handle, rtdb_int32 start, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
// */
func RawRtdbbSearchInBatchesWarp() {}

// /*
// *
// *
// * \brief 搜索符合条件的标签点，使用标签点名时支持通配符
// *
// * \param handle        连接句柄
// * \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
// * \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
// * \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
// *                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
// * \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
// *                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
// * \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
// *                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
// * \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
// * \param typemask      字符串，输入参数，标签点类型名称。缺省设置为空，长度不得超过 RTDB_TYPE_NAME_SIZE,
// *                        内置的普通数据类型可以使用 bool、uint8、datetime等字符串表示，不区分大小写，支持模糊搜索。
// * \param classofmask   整型，输入参数，标签点的类别，缺省设置为-1，表示可以是任意类型的标签点，
// *                        当使用标签点类型作为搜索条件时，必须是RTDB_CLASS枚举中的一项或者多项的组合。
// * \param timeunitmask  整型，输入参数，标签点的时间戳精度，缺省设置为-1，表示可以是任意时间戳精度，
// *                        当使用此时间戳精度作为搜索条件时，timeunitmask的值可以为0或1，0表示时间戳精度为秒，1表示纳秒
// * \param othertypemask 整型，输入参数，使用其他标签点属性作为搜索条件，缺省设置为0，表示不作为搜索条件，
// *                        当使用此参数作为搜索条件时，othertypemaskvalue作为对应的搜索值，
// *                        此参数的取值可以参考rtdb.h文件中的RTDB_SEARCH_MASK。
// * \param othertypemaskvalue
// *                        字符串，输入参数，当使用其他标签点属性作为搜索条件时，此参数作为对应的搜索值，缺省设置为0，表示不作为搜索条件，
// *                        如果othertypemask的值为0，或者RTDB_SEARCH_NULL，则此参数被忽略,
// *                        当othertypemask对应的标签点属性为数值类型时，此搜索值只支持相等判断，
// *                        当othertypemask对应的标签点属性为字符串类型时，此搜索值支持模糊搜索。
// * \param mode          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
// *                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
// * \param ids           整型数组，输出，返回搜索到的标签点标识列表
// * \param count         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
// * \remark 用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
// *        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
// *        如果 tagmask、tablemask 为空指针，则表示使用缺省设置"*",
// *        多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
// rtdb_error RTDBAPI_CALLRULE rtdbb_search_ex_warp(rtdb_int32 handle, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, const char *typemask, rtdb_int32 classofmask, rtdb_int32 timeunitmask, rtdb_int32 othertypemask, const char *othertypemaskvalue, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
// */
func RawRtdbbSearchExWarp() {}

// /*
// *
// *
// * \brief 搜索符合条件的标签点，使用标签点名时支持通配符
// *
// * \param handle        连接句柄
// * \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
// * \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
// * \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
// *                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
// * \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
// *                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
// * \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
// *                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
// * \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
// * \param typemask      字符串，输入参数，标签点类型名称。缺省设置为空，长度不得超过 RTDB_TYPE_NAME_SIZE,
// *                        内置的普通数据类型可以使用 bool、uint8、datetime等字符串表示，不区分大小写，支持模糊搜索。
// * \param classofmask   整型，输入参数，标签点的类别，缺省设置为-1，表示可以是任意类型的标签点，
// *                        当使用标签点类型作为搜索条件时，必须是RTDB_CLASS枚举中的一项或者多项的组合。
// * \param timeunitmask  整型，输入参数，标签点的时间戳精度，缺省设置为-1，表示可以是任意时间戳精度，
// *                        当使用此时间戳精度作为搜索条件时，timeunitmask的值可以为0或1，0表示时间戳精度为秒，1表示纳秒
// * \param othertypemask 整型，输入参数，使用其他标签点属性作为搜索条件，缺省设置为0，表示不作为搜索条件，
// *                        当使用此参数作为搜索条件时，othertypemaskvalue作为对应的搜索值，
// *                        此参数的取值可以参考rtdb.h文件中的RTDB_SEARCH_MASK。
// * \param othertypemaskvalue
// *                        字符串，输入参数，当使用其他标签点属性作为搜索条件时，此参数作为对应的搜索值，缺省设置为0，表示不作为搜索条件，
// *                        如果othertypemask的值为0，或者RTDB_SEARCH_NULL，则此参数被忽略,
// *                        当othertypemask对应的标签点属性为数值类型时，此搜索值只支持相等判断，
// *                        当othertypemask对应的标签点属性为字符串类型时，此搜索值支持模糊搜索。
// * \param count         整型，输出，表示搜索到的标签点个数
// * \remark  各参数中包含的搜索条件之间的关系为"与"的关系，
// *        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
// *        如果 tagmask、tablemask 为空指针，则表示使用缺省设置"*",
// *        多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
// rtdb_error RTDBAPI_CALLRULE rtdbb_search_points_count_warp(rtdb_int32 handle, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, const char *typemask, rtdb_int32 classofmask, rtdb_int32 timeunitmask, rtdb_int32 othertypemask, const char *othertypemaskvalue, rtdb_int32 *count)
// */
func RawRtdbbSearchPointsCountWarp() {}

// /*
// *
//   - 命名：rtdbb_remove_table_by_id
//   - \brief 根据表 id 删除表及表中标签点
//     *
//   - \param handle        连接句柄
//   - \param id            整型，输入，表 id
//   - \remark 删除的表不可恢复，删除的标签点可以通过 rtdbb_recover_point 接口恢复。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_remove_table_by_id_warp(rtdb_int32 handle, rtdb_int32 id)
// */
func RawRtdbbRemoveTableByIdWarp() {}

// /*
// *
// *
// * \brief 根据表名删除表及表中标签点
// *
// * \param handle        连接句柄
// * \param name          字符串，输入，表名称
// * \remark 删除的表不可恢复，删除的标签点可以通过 rtdbb_recover_point 接口恢复。
// rtdb_error RTDBAPI_CALLRULE rtdbb_remove_table_by_name_warp(rtdb_int32 handle, const char *name)
// */
func RawRtdbbRemoveTableByNameWarp() {}

// /*
// *
// *
// * \brief 更新单个标签点属性
// *
// * \param handle        连接句柄
// * \param base RTDB_POINT 结构，输入，基本标签点属性集。
// * \param scan RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
// * \param calc RTDB_CALC_POINT 结构，输入，计算标签点扩展属性集。
// * \remark 标签点由 base 参数的 id 字段指定，其中 id、table、type、millisecond 字段不能修改，
// *      changedate、changer、createdate、creator 字段由系统维护，其余字段均可修改，
// *      包括 classof 字段。输入参数中 scan、calc 可为空指针，对应的扩展属性将保持不变。
// rtdb_error RTDBAPI_CALLRULE rtdbb_update_point_property_warp(rtdb_int32 handle, const RTDB_POINT *base, const RTDB_SCAN_POINT *scan, const RTDB_CALC_POINT *calc)
// */
func RawRtdbbUpdatePointPropertyWarp() {}

// /*
// *
// * 命名：rtdbb_update_max_point_property
// * 功能：按最大长度更新单个标签点属性
// * 参数：
// *        [handle]        连接句柄
// *        [base] RTDB_POINT 结构，输入，基本标签点属性集。
// *        [scan] RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
// *        [calc] RTDB_MAX_CALC_POINT 结构，输入，计算标签点扩展属性集。
// * 备注：标签点由 base 参数的 id 字段指定，其中 id、table、type、millisecond 字段不能修改，
// *      changedate、changer、createdate、creator 字段由系统维护，其余字段均可修改，
// *      包括 classof 字段。输入参数中 scan、calc 可为空指针，对应的扩展属性将保持不变。
// rtdb_error RTDBAPI_CALLRULE rtdbb_update_max_point_property_warp(rtdb_int32 handle, const RTDB_POINT *base, const RTDB_SCAN_POINT *scan, const RTDB_MAX_CALC_POINT *calc)
// */
func RawRtdbbUpdateMaxPointPropertyWarp() {}

// /*
// *
//
//	*
//	* \brief 根据 "表名.标签点名" 格式批量获取标签点标识
//	*
//	* \param handle           连接句柄
//	* \param count            整数，输入/输出，输入时表示标签点个数
//	*                           (即table_dot_tags、ids、types、classof、use_ms 的长度)，
//	*                           输出时表示找到的标签点个数
//	* \param table_dot_tags   字符串指针数组，输入，"表名.标签点名" 列表
//	* \param ids              整型数组，输出，标签点标识列表, 返回 0 表示未找到
//	* \param types            整型数组，输出，标签点数据类型
//	* \param classof          整型数组，输出，标签点类别
//	* \param use_ms           短整型数组，输出，时间戳精度，
//	*                           返回 1 表示时间戳精度为纳秒， 为 0 表示为秒。
//	* \remark 用户须保证分配给 table_dot_tags、ids、types、classof、use_ms 的空间与count相符，
//	*        其中 types、classof、use_ms 可为空指针，对应的字段将不再返回。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_find_points_warp(rtdb_int32 handle, rtdb_int32 *count, const char* const* table_dot_tags, rtdb_int32 *ids, rtdb_int32 *types, rtdb_int32 *classof, rtdb_int16 *use_ms)
// */
func RawRtdbbFindPointsWarp() {}

// /*
// *
//
//	*
//	* \brief 根据 "表名.标签点名" 格式批量获取标签点标识
//	*
//	* \param handle           连接句柄
//	* \param count            整数，输入/输出，输入时表示标签点个数
//	*                           (即table_dot_tags、ids、types、classof、use_ms 的长度)，
//	*                           输出时表示找到的标签点个数
//	* \param table_dot_tags   字符串指针数组，输入，"表名.标签点名" 列表
//	* \param ids              整型数组，输出，标签点标识列表, 返回 0 表示未找到
//	* \param types            整型数组，输出，标签点数据类型
//	* \param classof          整型数组，输出，标签点类别
//	* \param precisions       数组，输出，时间戳精度，
//	*                           0表示秒，1表示毫秒，2表示微秒，3纳秒。
//	* \param errors           无符号整型数组，输出，表示每个标签点的查询结果的错误码
//	* \remark 用户须保证分配给 table_dot_tags、ids、types、classof、precisions、errors 的空间与count相符，
//	*        其中 types、classof、precisions、errors 可为空指针，对应的字段将不再返回。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_find_points_ex_warp(rtdb_int32 handle, rtdb_int32* count, const char* const* table_dot_tags, rtdb_int32* ids, rtdb_int32* types, rtdb_int32* classof, rtdb_precision_type* precisions, rtdb_error* errors)
// */
func RawRtdbbFindPointsExWarp() {}

// /*
// *
//
//	*
//	* \brief 根据标签属性字段对标签点标识进行排序
//	*
//	* \param handle           连接句柄
//	* \param count            整数，输入，表示标签点个数, 即 ids 的长度
//	* \param ids              整型数组，输入，标签点标识列表
//	* \param index            整型，输入，属性字段枚举，参见 RTDB_TAG_FIELD_INDEX，
//	*                           将根据该字段对 ID 进行排序。
//	* \param flag             整型，输入，标志位组合，参见 RTDB_TAG_SORT_FLAG 枚举，其中
//	*                           RTDB_SORT_FLAG_DESCEND             表示降序排序，不设置表示升序排列；
//	*                           RTDB_SORT_FLAG_CASE_SENSITIVE      表示进行字符串类型字段比较时大小写敏感，不设置表示不区分大小写；
//	*                           RTDB_SORT_FLAG_RECYCLED            表示对可回收标签进行排序，不设置表示对正常标签排序，
//	*                           不同的标志位可通过"或"运算连接在一起，
//	*                           当对可回收标签排序时，以下字段索引不可使用：
//	*                               RTDB_TAG_INDEX_TIMESTAMP
//	*                               RTDB_TAG_INDEX_VALUE
//	*                               RTDB_TAG_INDEX_QUALITY
//	* \remark 用户须保证分配给 ids 的空间与 count 相符, 如果 ID 指定的标签并不存在，
//	*        或标签不具备要求排序的字段 (如对非计算点进行方程式排序)，它们将被放置在数组的尾部。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_sort_points_warp(rtdb_int32 handle, rtdb_int32 count, rtdb_int32 *ids, rtdb_int32 index, rtdb_int32 flag)
// */
func RawRtdbbSortPointsWarp() {}

// /*
// *
//
//	*
//	* \brief 根据表 ID 更新表名称。
//	*
//	* \param handle    连接句柄
//	* \param tab_id    整型，输入，要修改表的标识
//	* \param name      字符串，输入，新的标签点表名称。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_update_table_name_warp(rtdb_int32 handle, rtdb_int32 tab_id, const char *name)
// */
func RawRtdbbUpdateTableNameWarp() {}

// /*
// *
//
//	*
//	* \brief 根据表 ID 更新表描述。
//	*
//	* \param handle    连接句柄
//	* \param tab_id    整型，输入，要修改表的标识
//	* \param desc      字符串，输入，新的表描述。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_update_table_desc_by_id_warp(rtdb_int32 handle, rtdb_int32 tab_id, const char *desc)
// */
func RawRtdbbUpdateTableDescByIdWarp() {}

// /*
// *
// *
// * \brief 根据表名称更新表描述。
// *
// * \param handle    连接句柄
// * \param name      字符串，输入，要修改表的名称。
// * \param desc      字符串，输入，新的表描述。
// rtdb_error RTDBAPI_CALLRULE rtdbb_update_table_desc_by_name_warp(rtdb_int32 handle, const char *name, const char *desc)
// */
func RawRtdbbUpdateTableDescByNameWarp() {}

// /*
// *
//
//	*
//	* \brief 恢复已删除标签点
//	*
//	* \param handle    连接句柄
//	* \param table_id  整型，输入，要将标签点恢复到的表标识
//	* \param point_id  整型，输入，待恢复的标签点标识
//	* 备注: 本接口只对可回收标签点(通过接口rtdbb_remove_point_by_id/rtdbb_remove_point_by_tag)有效，
//	*        对正常的标签点没有作用。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_recover_point_warp(rtdb_int32 handle, rtdb_int32 table_id, rtdb_int32 point_id)
// */
func RawRtdbbRecoverPointWarp() {}

// /*
// *
//
//	*
//	* \brief 清除标签点
//	*
//	* \param handle    连接句柄
//	* \param id        整数，输入，要清除的标签点标识
//	* 备注: 本接口仅对可回收标签点(通过接口rtdbb_remove_point_by_id/rtdbb_remove_point_by_name)有效，
//	*      对正常的标签点没有作用。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_purge_point_warp(rtdb_int32 handle, rtdb_int32 id)
// */
func RawRtdbbPurgePointWarp() {}

// /*
// *
//
//	*
//	* \brief 获取可回收标签点数量
//	*
//	* \param handle    连接句柄
//	* \param count     整型，输出，可回收标签点的数量
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_points_count_warp(rtdb_int32 handle, rtdb_int32 *count)
// */
func RawRtdbbGetRecycledPointsCountWarp() {}

// /*
// *
//
//	*
//	* \brief 获取可回收标签点 id 列表
//	*
//	* \param handle    连接句柄
//	* \param ids       整型数组，输出，可回收标签点 id
//	* \param count     整型，输入/输出，标签点个数，
//	*                    输入时表示 ids 的长度，
//	*                    输出时表示成功获取标签点的个数。
//	* \remark 用户须保证 ids 的长度与 count 一致
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_points_warp(rtdb_int32 handle, rtdb_int32 *ids, rtdb_int32 *count)
// */
func RawRtdbbGetRecycledPointsWarp() {}

// /*
// *
// * 命名：rtdbb_search_recycled_points
// * 功能：搜索符合条件的可回收标签点，使用标签点名时支持通配符
// * 参数：
// *        [handle]        连接句柄
// *        [tagmask]       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
// *        [tablemask]     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
// *        [source]        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
// *                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
// *        [unit]          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
// *                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
// *        [desc]          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
// *                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
// *        [instrument]    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
// *        [mode]          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
// *                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
// *        [ids]           整型数组，输出，返回搜索到的标签点标识列表
// *        [count]         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
// * 备注：用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
// *        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
// *        如果 tagmask、fullmask 为空指针，则表示使用缺省设置"*"
// rtdb_error RTDBAPI_CALLRULE rtdbb_search_recycled_points_warp(rtdb_int32 handle, const char *tagmask, const char *fullmask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
// */
func RawRtdbbSearchRecycledPointsWarp() {}

// /*
// *
//
//	*
//	* \brief 分批搜索符合条件的可回收标签点，使用标签点名时支持通配符
//	*
//	* \param handle        连接句柄
//	* \param start         整型，输入，搜索的起始位置。
//	* \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
//	* \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
//	* \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
//	*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
//	* \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
//	*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
//	* \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
//	*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
//	* \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
//	* \param mode          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
//	*                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
//	* \param ids           整型数组，输出，返回搜索到的标签点标识列表
//	* \param count         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
//	* \remark 用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
//	*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
//	*        如果 tagmask、fullmask 为空指针，则表示使用缺省设置"*"
//	*        当搜索到的标签点数比提供的要小时，表示这是最后一批符合条件的标签点 (即全部搜索完毕)。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_search_recycled_points_in_batches_warp(rtdb_int32 handle, rtdb_int32 start, const char *tagmask, const char *fullmask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
// */
func RawRtdbbSearchRecycledPointsInBatchesWarp() {}

// /*
// *
//
//	*
//	* \brief 获取可回收标签点的属性
//	*
//	* \param handle   连接句柄
//	* \param base     RTDB_POINT 结构，输入/输出，标签点基本属性。
//	输入时，由 id 字段指定要取得的可回收标签点。
//	* \param scan     RTDB_SCAN_POINT 结构，输出，标签点采集扩展属性
//	* \param calc     RTDB_CALC_POINT 结构，输出，标签点计算扩展属性
//	* \remark scan、calc 可为空指针，对应的扩展信息将不返回。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_point_property_warp(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_CALC_POINT *calc)
// */
func RawRtdbbGetRecycledPointPropertyWarp() {}

// /*
// *
//   - 命名：rtdbb_get_recycled_max_point_property
//   - 功能：按最大长度获取可回收标签点的属性
//   - 参数：
//   - [handle]   连接句柄
//   - [base]     RTDB_POINT 结构，输入/输出，标签点基本属性。
//   - 输入时，由 id 字段指定要取得的可回收标签点。
//   - [scan]     RTDB_SCAN_POINT 结构，输出，标签点采集扩展属性
//   - [calc]     RTDB_MAX_CALC_POINT 结构，输出，标签点计算扩展属性
//   - 备注：scan、calc 可为空指针，对应的扩展信息将不返回。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_max_point_property_warp(rtdb_int32 handle, RTDB_POINT* base, RTDB_SCAN_POINT* scan, RTDB_MAX_CALC_POINT* calc)
// */
func RawRtdbbGetRecycledMaxPointPropertyWarp() {}

// /*
// *
//
//	*
//	* \brief 清空标签点回收站
//	*
//	* \param handle   连接句柄
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_clear_recycler_warp(rtdb_int32 handle)
// */
func RawRtdbbClearRecyclerWarp() {}

// /*
// *
//   - 命名：rtdbb_subscribe_tags_ex
//   - 功能：标签点属性更改通知订阅
//   - 参数：
//   - [handle]    连接句柄
//   - [options]   整型，输入，订阅选项，参见枚举RTDB_OPTION
//   - RTDB_O_AUTOCONN 订阅客户端与数据库服务器网络中断后自动重连并订阅
//   - [param]     输入，用户参数，
//   - 作为rtdbb_tags_change_ex的param参数
//   - [callback]  rtdbb_tags_change_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
//   - 当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
//   - 作为event_type取值调用回掉函数后退出订阅。
//   - 当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
//   - 作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
//   - 网络中断期间回掉函数调用频率为最少3秒
//   - event_type参数值含义如下：
//   - RTDB_E_DATA        标签点属性发生更改
//   - RTDB_E_DISCONNECT  订阅客户端与数据库网络断开
//   - RTDB_E_RECOVERY    订阅客户端与数据库网络及订阅恢复
//   - handle 产生订阅回掉的连接句柄，调用rtdbb_subscribe_tags_ex时的handle参数
//   - param  用户自定义参数，调用rtdbb_subscribe_tags_ex时的param参数
//   - count  event_type为RTDB_E_DATA时表示ids的数量
//   - event_type为其它值时，count值为0
//   - ids    event_type为RTDB_E_DATA时表示属性更改的标签点ID，数量由count指定
//   - event_type为其它值时，ids值为NULL
//   - what   event_type为RTDB_E_DATA时表示属性变更原因，参考RTDB_TAG_CHANGE_REASON
//   - event_type为其它值时，what时值为0
//   - 备注：用于订阅测点的连接句柄必需是独立的，不能再用来调用其它 api，
//   - 否则返回 RtE_OTHER_SDK_DOING 错误。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_subscribe_tags_ex_warp(rtdb_int32 handle, rtdb_uint32 options, void* param, rtdbb_tags_change_event_ex callback)
// */
func RawRtdbbSubscribeTagsExWarp() {}

// /*
// *
//
//	*
//	* \brief 取消标签点属性更改通知订阅
//	*
//	* \param handle    连接句柄
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_cancel_subscribe_tags_warp(rtdb_int32 handle)
// */
func RawRtdbbCancelSubscribeTagsWarp() {}

// /*
// *
// * 命名：rtdbb_create_named_type
// * 功能：创建自定义类型
// * 参数：
// *        [handle]      连接句柄，输入参数
// *        [name]        自定义类型的名称，类型的唯一标示,不能重复，长度不能超过RTDB_TYPE_NAME_SIZE，输入参数
// *        [field_count]    自定义类型中包含的字段的个数,输入参数
// *        [fields]      自定义类型中包含的字段的属性，RTDB_DATA_TYPE_FIELD结构的数组，个数与field_count相等，输入参数
// *              RTDB_DATA_TYPE_FIELD中的length只对type为str或blob类型的数据有效。其他类型忽略
// * 备注：自定义类型的大小必须要小于数据页大小(小于数据页大小的2/3，即需要合理定义字段的个数及每个字段的长度)。
// rtdb_error RTDBAPI_CALLRULE rtdbb_create_named_type_warp(rtdb_int32 handle, const char* name, rtdb_int32 field_count, const RTDB_DATA_TYPE_FIELD* fields, char desc[RTDB_DESC_SIZE])
// */
func RawRtdbbCreateNamedTypeWarp() {}

// /*
// *
//   - 命名：rtdbb_get_named_types_count
//   - 功能：获取所有的自定义类型的总数
//   - 参数：
//   - [handle]      连接句柄，输入参数
//   - [count]      返回所有的自定义类型的总数，输入/输出参数
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_named_types_count_warp(rtdb_int32 handle, rtdb_int32* count)
// */
func RawRtdbbGetNamedTypesCountWarp() {}

// /*
// *
//   - 命名：rtdbb_get_all_named_types
//   - 功能：获取所有的自定义类型
//   - 参数：
//   - [handle]      连接句柄，输入参数
//   - [count]      返回所有的自定义类型的总数，输入/输出参数，输入:为name,field_counts数组的长度，输出:获取的实际自定义类型的个数
//   - [name]        返回所有的自定义类型的名称的数组，每个自定义类型的名称的长度不超过RTDB_TYPE_NAME_SIZE，输入/输出参数
//   - 输入：name数组长度要等于count.输出：实际获取的自定义类型名称的数组
//   - [field_counts]    返回所有的自定义类型所包含字段个数的数组，输入/输出参数
//   - 输入：field_counts数组长度要等于count。输出:实际每个自定义类型所包含的字段的个数的数组
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_all_named_types_warp(rtdb_int32 handle, rtdb_int32* count, char* name[RTDB_TYPE_NAME_SIZE], rtdb_int32* field_counts)
// */
func RawRtdbbGetAllNamedTypesWarp() {}

// /*
// *
// * 命名：rtdbb_get_named_type
// * 功能：获取自定义类型的所有字段
// * 参数：
// *        [handle]         连接句柄，输入参数
// *        [name]           自定义类型的名称，输入参数
// *        [field_count]    返回name指定的自定义类型的字段个数，输入/输出参数
// *                         输入：指定fields数组长度.输出：实际的name自定义类型的字段的个数
// *        [fields]         返回由name所指定的自定义类型所包含字段RTDB_DATA_TYPE_FIELD结构的数组，输入/输出参数
// *                         输入：fields数组长度要等于count。输出:RTDB_DATA_TYPE_FIELD结构的数组
// *        [type_size]      所有自定义类型fields结构中长度字段的累加和，输出参数
// *        [desc]           自定义类型的描述，输出参数
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_named_type_warp(rtdb_int32 handle, const char* name, rtdb_int32* field_count, RTDB_DATA_TYPE_FIELD* fields, rtdb_int32* type_size, char desc[RTDB_DESC_SIZE])
// */
func RawRtdbbGetNamedTypeWarp() {}

// /*
// *
// * 命名：rtdbb_remove_named_type
// * 功能：删除自定义类型
// * 参数：
// *        [handle]      连接句柄，输入参数
// *        [name]        自定义类型的名称，输入参数
// *        [reserved]      保留字段,暂时不用
// rtdb_error RTDBAPI_CALLRULE rtdbb_remove_named_type_warp(rtdb_int32 handle, const char* name, rtdb_int32 reserved GAPI_DEFAULT_VALUE(0))
// */
func RawRtdbbRemoveNamedTypeWarp() {}

// /*
// *
//   - 命名：rtdbb_get_named_type_names_property
//   - 功能：根据标签点id查询标签点所对应的自定义类型的名字和字段总数
//   - 参数：
//   - [handle]           连接句柄
//   - [count]            输入/输出，标签点个数，
//   - 输入时表示 ids、named_type_names、field_counts、errors 的长度，
//   - 输出时表示成功获取自定义类型名字的标签点个数
//   - [ids]              整型数组，输入，标签点标识列表
//   - [named_type_names] 字符串数组，输出，标签点自定义类型的名字
//   - [field_counts]     整型数组，输出，标签点自定义类型的字段个数
//   - [errors]           无符号整型数组，输出，获取自定义类型名字的返回值列表，参考rtdb_error.h
//   - 备注：用户须保证 ids、named_type_names、field_counts、errors 的长度与 count 一致。
//   - 本接口只对数据类型为 RTDB_NAMED_T 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_named_type_names_property_warp(rtdb_int32 handle, rtdb_int32 *count, rtdb_int32 *ids, char* const *named_type_names, rtdb_int32 *field_counts, rtdb_error *errors)
// */
func RawRtdbbGetNamedTypeNamesPropertyWarp() {}

// /*
// *
//   - 命名：rtdbb_get_recycled_named_type_names_property
//   - 功能：根据回收站标签点id查询标签点所对应的自定义类型的名字和字段总数
//   - 参数：
//   - [handle]           连接句柄
//   - [count]            输入/输出，标签点个数，
//   - 输入时表示 ids、named_type_names、field_counts、errors 的长度，
//   - 输出时表示成功获取自定义类型名字的标签点个数
//   - [ids]              整型数组，输入，回收站标签点标识列表
//   - [named_type_names] 字符串数组，输出，标签点自定义类型的名字
//   - [field_counts]     整型数组，输出，标签点自定义类型的字段个数
//   - [errors]           无符号整型数组，输出，获取自定义类型名字的返回值列表，参考rtdb_error.h
//   - 备注：用户须保证 ids、named_type_names、field_counts、errors 的长度与 count 一致。
//   - 本接口只对数据类型为 RTDB_NAMED_T 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_named_type_names_property_warp(rtdb_int32 handle, rtdb_int32 *count, rtdb_int32 *ids, char* const *named_type_names, rtdb_int32 *field_counts, rtdb_error *errors)
// */
func RawRtdbbGetRecycledNamedTypeNamesPropertyWarp() {}

// /*
// *
// * 命名：rtdbb_get_named_type_points_count
// * 功能：获取该自定义类型的所有标签点个数
// * 参数：
// *        [handle]           连接句柄，输入参数
// *        [name]             自定义类型的名称，输入参数
// *        [points_count]     返回name指定的自定义类型的标签点个数，输入参数
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_named_type_points_count_warp(rtdb_int32 handle, const char* name, rtdb_int32 *points_count)
// */
func RawRtdbbGetNamedTypePointsCountWarp() {}

// /*
// *
// *
// * \brief 获取该内置的基本类型的所有标签点个数
// *
// * \param handle           整型，输入参数，连接句equation[RTDB_MAX_EQUATION_SIZE]柄
// * \param type             整型，输入参数，内置的基本类型，参数的值可以是除RTDB_NAME_T以外的所有RTDB_TYPE枚举值
// * \param points_count     整型，输入参数，返回type指定的内置基本类型的标签点个数
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_base_type_points_count_warp(rtdb_int32 handle, rtdb_int32 type, rtdb_int32 *points_count)
// */
func RawRtdbbGetBaseTypePointsCountWarp() {}

// /*
// *
// * 命名：rtdbb_modify_named_type
// * 功能：修改自定义类型名称,描述,字段名称,字段描述
// * 参数：
// *        [handle]             连接句柄，输入参数
// *        [name]               自定义类型的名称，输入参数
// *        [modify_name]        要修改的自定义类型名称，输入参数
// *        [modify_desc]        要修改的自定义类型的描述，输入参数
// *        [modify_field_name]  要修改的自定义类型字段的名称，输入参数
// *        [modify_field_desc]  要修改的自定义类型字段的描述，输入参数
// *        [field_count]        自定义类型字段的个数，输入参数
// rtdb_error RTDBAPI_CALLRULE rtdbb_modify_named_type_warp(rtdb_int32 handle, const char* name, const char* modify_name, const char* modify_desc, const char* modify_field_name[RTDB_TYPE_NAME_SIZE], const char* modify_field_desc[RTDB_DESC_SIZE], rtdb_int32 field_count)
// */
func RawRtdbbModifyNamedTypeWarp() {}

// /*
// *
//
//	*
//	* \brief 获取元数据同步信息
//	*
//	* \param handle           整型，输入参数，连接句柄
//	* \param node_number      整型，输入参数，双活节点id，1表示第一个节点，2表示第二个节点。0表示所有节点
//	* \param count            整型，输入参数，sync_infos参数的数量
//	*                              输出参数，输出实际获取到的sync_infos的个数
//	* \param sync_infos       RTDB_SYNC_INFO数组，输出参数，输出实际获取到的同步信息
//	* \param errors           rtdb_error数组，输出参数，输出对应节点的错误信息
//
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_meta_sync_info_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32* count, RTDB_SYNC_INFO* sync_infos, rtdb_error* errors)
// */
func RawRtdbbGetMetaSyncInfoWarp() {}

// /*
// *
//
//	*
//	* \brief 批量读取开关量、模拟量快照数值
//	*
//	* \param handle    连接句柄
//	* \param count     整型，输入/输出，标签点个数，
//	*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
//	*                    输出时表示成功获取实时值的标签点个数
//	* \param ids       整型数组，输入，标签点标识列表
//	* \param datetimes 整型数组，输出，实时数值时间列表,
//	*                    表示距离1970年1月1日08:00:00的秒数
//	* \param ms        短整型数组，输出，实时数值时间列表，
//	*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
//	* \param values    双精度浮点型数组，输出，实时浮点型数值列表，
//	*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的快照值；否则为 0
//	* \param states    64 位整型数组，输出，实时整型数值列表，
//	*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的快照值；否则为 0
//	* \param qualities 短整型数组，输出，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
//	* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_get_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbsGetSnapshots64Warp() {}

// /*
//
//	*
//	* \brief 批量写入开关量、模拟量快照数值
//	*
//	* \param handle    连接句柄
//	* \param count     整型，输入/输出，标签点个数，
//	*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
//	*                    输出时表示成功写入实时值的标签点个数
//	* \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
//	*                    但它们的时间戳必需是递增的。
//	* \param datetimes 整型数组，输入，实时数值时间列表,
//	*                    表示距离1970年1月1日08:00:00的秒数
//	* \param ms        短整型数组，输入，实时数值时间列表，
//	*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
//	* \param values    双精度浮点型数组，输入，实时浮点型数值列表，
//	*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
//	* \param states    64 位整型数组，输入，实时整型数值列表，
//	*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
//	* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
//	* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_put_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbsPutSnapshots64Warp() {}

// /*
//
//	*
//	* \brief 批量写入开关量、模拟量快照数值
//	*
//	* \param handle    连接句柄
//	* \param count     整型，输入/输出，标签点个数，
//	*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
//	*                    输出时表示成功写入实时值的标签点个数
//	* \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
//	*                    但它们的时间戳必需是递增的。
//	* \param datetimes 整型数组，输入，实时数值时间列表,
//	*                    表示距离1970年1月1日08:00:00的秒数
//	* \param ms        短整型数组，输入，实时数值时间列表，
//	*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
//	* \param values    双精度浮点型数组，输入，实时浮点型数值列表，
//	*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
//	* \param states    64 位整型数组，输入，实时整型数值列表，
//	*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
//	* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
//	* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
//	*        仅当输入时间戳与当前快照时间戳完全相等时，会替换当前快照的值和质量；
//	*        其余情况下会调用 rtdbs_put_snapshots()
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_fix_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbsPutSnapshots() {}

// /*
//
//	*
//	* \brief 批量回溯快照
//	*
//	* \param handle    连接句柄
//	* \param count     整型，输入/输出，标签点个数，
//	*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
//	*                    输出时表示成功写入实时值的标签点个数
//	* \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
//	*
//	* \param datetimes 整型数组，输入，实时数值时间列表,
//	*                    表示距离1970年1月1日08:00:00的秒数
//	* \param ms        短整型数组，输入，实时数值时间列表，
//	*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
//	* \param values    双精度浮点型数组，输入，实时浮点型数值列表，
//	*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
//	* \param states    64 位整型数组，输入，实时整型数值列表，
//	*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
//	* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
//	* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
//	* 功能说明：
//	*       批量将标签点的快照值vtmq改成传入的vtmq，如果传入的时间戳早于当前快照，会删除传入时间戳到当前快照的历史存储值。
//	*       如果传入的时间戳等于或者晚于当前快照，什么也不做。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_back_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbsBackSnapshots64Warp() {}

// /*
//
//	*
//	* \brief 批量读取坐标实时数据
//	*
//	* \param handle    连接句柄
//	* \param count     整型，输入/输出，标签点个数，
//	*                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
//	*                    输出时表示成功获取实时值的标签点个数
//	* \param ids       整型数组，输入，标签点标识列表
//	* \param datetimes 整型数组，输出，实时数值时间列表,
//	*                    表示距离1970年1月1日08:00:00的秒数
//	* \param ms        短整型数组，输出，实时数值时间列表，
//	*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
//	* \param x         单精度浮点型数组，输出，实时浮点型横坐标数值列表
//	* \param y         单精度浮点型数组，输出，实时浮点型纵坐标数值列表
//	* \param qualities 短整型数组，输出，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
//	* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
//	*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_get_coor_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float32* x, rtdb_float32* y, rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbsGetCoorSnapshots64Warp() {}

// /*
//
//	*
//	* \brief 批量写入坐标实时数据
//	*
//	* \param handle    连接句柄
//	* \param count     整型，输入/输出，标签点个数，
//	*                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
//	*                    输出时表示成功获取实时值的标签点个数
//	* \param ids       整型数组，输入，标签点标识列表
//	* \param datetimes 整型数组，输入，实时数值时间列表,
//	*                    表示距离1970年1月1日08:00:00的秒数
//	* \param ms        短整型数组，输入，实时数值时间列表，
//	*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
//	* \param x         单精度浮点型数组，输入，实时浮点型横坐标数值列表
//	* \param y         单精度浮点型数组，输入，实时浮点型纵坐标数值列表
//	* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \param errors    无符号整型数组，输出，写入实时坐标数据的返回值列表，参考rtdb_error.h
//	* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
//	*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_put_coor_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float32* x, const rtdb_float32* y, const rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbsPutCoorSnapshots64Warp() {}

// /*
//
//	*
//	* \brief 批量写入坐标实时数据
//	*
//	* \param handle    连接句柄
//	* \param count     整型，输入/输出，标签点个数，
//	*                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
//	*                    输出时表示成功获取实时值的标签点个数
//	* \param ids       整型数组，输入，标签点标识列表
//	* \param datetimes 整型数组，输入，实时数值时间列表,
//	*                    表示距离1970年1月1日08:00:00的秒数
//	* \param ms        短整型数组，输入，实时数值时间列表，
//	*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
//	* \param x         单精度浮点型数组，输入，实时浮点型横坐标数值列表
//	* \param y         单精度浮点型数组，输入，实时浮点型纵坐标数值列表
//	* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \param errors    无符号整型数组，输出，写入实时坐标数据的返回值列表，参考rtdb_error.h
//	* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
//	*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
//	*        仅当输入时间戳与当前快照时间戳完全相等时，会替换当前快照的值和质量；
//	*        其余情况下会调用 rtdbs_put_coor_snapshots()
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_fix_coor_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float32* x, const rtdb_float32* y, const rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbsPutCoorSnapshots() {}

// /*
//
//	*
//	* \brief 读取二进制/字符串实时数据
//	*
//	* \param handle    连接句柄
//	* \param id        整型，输入，标签点标识
//	* \param datetime  整型，输出，实时数值时间列表,
//	*                    表示距离1970年1月1日08:00:00的秒数
//	* \param ms        短整型，输出，实时数值时间列表，
//	*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
//	* \param blob      字节型数组，输出，实时二进制/字符串数值
//	* \param len       短整型，输出，二进制/字符串数值长度
//	* \param quality   短整型，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_get_blob_snapshot64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type* datetime, rtdb_subtime_type* subtime, rtdb_byte* blob, rtdb_length_type* len, rtdb_int16* quality)
// */
func RawRtdbsGetBlobSnapshot64Warp() {}

// /*
//
//	*
//	* \brief 批量读取二进制/字符串实时数据
//	*
//	* \param handle    连接句柄
//	* \param count     整型，输入/输出，标签点个数，
//	*                    输入时表示 ids、datetimes、ms、blobs、lens、qualities、errors 的长度，
//	*                    输出时表示成功获取实时值的标签点个数
//	* \param ids       整型数组，输入，标签点标识
//	* \param datetimes 整型数组，输出，实时数值时间列表,
//	*                    表示距离1970年1月1日08:00:00的秒数
//	* \param ms        短整型数组，输出，实时数值时间列表，
//	*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
//	* \param blobs     字节型指针数组，输出，实时二进制/字符串数值
//	* \param lens      短整型数组，输入/输出，二进制/字符串数值长度，
//	*                    输入时表示对应的 blobs 指针指向的缓冲区长度，
//	*                    输出时表示实际得到的 blob 长度，如果 blob 的长度大于缓冲区长度，会被截断。
//	* \param qualities 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
//	* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_get_blob_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_byte* const* blobs, rtdb_length_type* lens, rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbsGetBlobSnapshots64Warp() {}

//   - \brief 写入二进制/字符串实时数据
//     *
//   - \param handle    连接句柄
//   - \param id        整型，输入，标签点标识
//   - \param datetime  整型，输入，实时数值时间列表,
//   - 表示距离1970年1月1日08:00:00的秒数
//   - \param ms        短整型，输入，实时数值时间列表，
//   - 对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
//   - \param blob      字节型数组，输入，实时二进制/字符串数值
//   - \param len       短整型，输入，二进制/字符串数值长度，超过一个页大小数据将被截断。
//   - \param quality   短整型，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_put_blob_snapshot64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, const rtdb_byte* blob, rtdb_length_type len, rtdb_int16 quality)
// */
func RawRtdbsPutBlobSnapshot64Warp() {}

// /*
//
//	*
//	* \brief 批量写入二进制/字符串实时数据
//	*
//	* \param handle    连接句柄
//	* \param count     整型，输入/输出，标签点个数，
//	*                    输入时表示 ids、datetimes、ms、blobs、lens、qualities、errors 的长度，
//	*                    输出时表示成功获取实时值的标签点个数
//	* \param ids       整型数组，输入，标签点标识
//	* \param datetimes 整型数组，输入，实时数值时间列表,
//	*                    表示距离1970年1月1日08:00:00的秒数
//	* \param ms        短整型数组，输入，实时数值时间列表，
//	*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
//	* \param blobs     字节型指针数组，输入，实时二进制/字符串数值
//	* \param lens      短整型数组，输入，二进制/字符串数值长度，
//	*                    表示对应的 blobs 指针指向的缓冲区长度，超过一个页大小数据将被截断。
//	* \param qualities 短整型数组，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
//	* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_put_blob_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_byte* const* blobs, const rtdb_length_type* lens, const rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbsPutBlobSnapshots64Warp() {}

// /*
//
//	*
//	* \brief 批量读取datetime类型标签点实时数据
//	*
//	* \param handle    连接句柄
//	* \param count     整型，输入/输出，标签点个数，
//	*                    输入时表示 ids、datetimes、ms、dtvalues、dtlens、qualities、errors 的长度，
//	*                    输出时表示成功获取实时值的标签点个数
//	* \param ids       整型数组，输入，标签点标识
//	* \param datetimes 整型数组，输出，实时数值时间列表,
//	*                    表示距离1970年1月1日08:00:00的秒数
//	* \param ms        短整型数组，输出，实时数值时间列表，
//	*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
//	* \param dtvalues  字节型指针数组，输出，实时datetime数值
//	* \param dtlens    短整型数组，输入/输出，datetime数值长度，
//	*                    输入时表示对应的 dtvalues 指针指向的缓冲区长度，
//	* \param qualities 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
//	* \param type      短整型，输入，所有标签点的显示类型，如“yyyy-mm-dd hh:mm:ss.000”的type为1，默认类型1，
//	*                    “yyyy/mm/dd hh:mm:ss.000”的type为2
//	*                    如果不传type，则按照标签点属性显示，否则按照type类型显示
//	* \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_get_datetime_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_byte* const* dtvalues, rtdb_length_type* dtlens, rtdb_int16* qualities, rtdb_error* errors, rtdb_int16 type)
// */
func RawRtdbsGetDatetimeSnapshots64Warp() {}

// /*
//
//	*
//	* \brief 批量插入datetime类型标签点数据
//	*
//	* \param handle      连接句柄
//	* \param count       整型，输入/输出，标签点个数，
//	*                      输入时表示 ids、datetimes、ms、dtvalues、dtlens、qualities、errors的长度，
//	*                      输出时表示成功写入的标签点个数
//	* \param ids         整型数组，输入，标签点标识
//	* \param datetimes   整型数组，输入，实时值时间列表
//	*                      表示距离1970年1月1日08:00:00的秒数
//	* \param ms          短整型数组，输入，实时数值时间列表，
//	*                      对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
//	* \param dtvalues    字节型指针数组，输入，datetime标签点的值
//	* \param dtlens      短整型数组，输入，数值长度
//	* \param qualities   短整型数组，输入，实时数值品质，，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \param errors      无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
//	* \remark 被接口只对数据类型 RTDB_DATETIME 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_put_datetime_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_byte* const* dtvalues, const rtdb_length_type* dtlens, const rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbsPutDatetimeSnapshots64Warp() {}

// /*
//
//	*
//	* \brief 批量标签点快照改变的通知订阅
//	*
//	* \param handle         连接句柄
//	* \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
//	*                           输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
//	* \param ids            整型数组，输入，标签点标识列表。
//	* \param options        订阅选项
//	*                           RTDB_O_AUTOCONN 自动重连
//	* \param param          用户自定义参数
//	* \param callback       rtdbs_snaps_event_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
//	*                       当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
//	*                       作为event_type取值调用回掉函数后退出订阅。
//	*                       当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
//	*                       作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
//	*                       网络中断期间回掉函数调用频率为最少3秒
//	*                       event_type参数值含义如下：
//	*                         RTDB_E_DATA        标签点快照改变
//	*                         RTDB_E_DISCONNECT  订阅客户端与数据库网络断开
//	*                         RTDB_E_RECOVERY    订阅客户端与数据库网络及订阅恢复
//	*                         RTDB_E_CHANGED     客户端修改订阅标签点信息，即通过rtdbs_change_subscribe_snapshots修改订阅信息的结果
//	*                       handle 产生订阅回掉的连接句柄，调用rtdbs_subscribe_snapshots_ex时的handle参数
//	*                       param  用户自定义参数，调用rtdbs_subscribe_snapshots_ex时的param参数
//	*                       count  event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示ids，datetimes,values等的数量
//	*                              event_type为其它值时，count值为0
//	*                       ids    event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示快照改变的标签点ID，数量由count指定
//	*                              event_type为其它值时，ids值为NULL
//	*                       datetimes event_type为RTDB_E_DATA时表示快照时间，数量由count指定
//	*                                 event_type为RTDB_E_CHANGED时表示changed_types，即通过rtdbs_change_subscribe_snapshots传入的changed_types
//	*                                 event_type为其它值时，datetimes值为NULL
//	*                       ms     event_type为RTDB_E_DATA时表示快照的毫秒，数量由count指定
//	*                              event_type为其它值时，ms值为NULL
//	*                       values event_type为RTDB_E_DATA时表示浮点数据类型快照值，数量由count指定
//	*                              event_type为RTDB_E_CHANGED时，表示delta_values，即通过rtdbs_change_subscribe_snapshots传入的delta_values
//	*                              event_type为其它值时，values值为NULL
//	*                       states event_type为RTDB_E_DATA时表示整形数据类型快照值，数量由count指定
//	*                              event_type为RTDB_E_CHANGED时，表示delta_states，通过rtdbs_change_subscribe_snapshots传入的delta_states
//	*                              event_type为其它值时，states值为NULL
//	*                       qualities event_type为RTDB_E_DATA时表示快照质量码，数量由count指定
//	*                              event_type为其它值时，qualities值为NULL
//	*                       errors event_type为RTDB_E_DATA时表示快照错误码，数量由count指定
//	*                              event_type为RTDB_E_CHANGED时，表示修改结果对应的错误码，数量由count指定
//	*                              event_type为其它值时，errors值为NULL
//	* \param errors         无符号整型数组，输出，
//	*                           写入实时数据的返回值列表，参考rtdb_error.h
//	* \remark   用户须保证 ids、errors 的长度与 count 一致。
//	*        用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
//	*        否则返回 RtE_OTHER_SDK_DOING 错误。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_subscribe_snapshots_ex64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_uint32 options, void* param, rtdbs_snaps_event_ex64 callback, rtdb_error* errors)
// */
func RawRtdbsSubscribeSnapshotsEx64Warp() {}

// /*
//
//	*
//	* \brief 批量标签点快照改变的通知订阅
//	*
//	* \param handle         连接句柄
//	* \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
//	*                           输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
//	* \param ids            整型数组，输入，标签点标识列表。
//	* \param delta_values   double型数组，输入，订阅浮点类型标签点的容差值，变化超过设置的容差值才会推送
//	* \param delta_values   整型数组，输入，订阅整型标签点的容差值，变化超过设置的容差值才会推送
//	* \param options        订阅选项
//	*                           RTDB_O_AUTOCONN 自动重连
//	* \param param          用户自定义参数
//	* \param callback       rtdbs_snaps_event_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
//	*                         当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
//	*                         作为event_type取值调用回掉函数后退出订阅。
//	*                         当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
//	*                         作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
//	*                         网络中断期间回掉函数调用频率为最少3秒
//	*                         event_type参数值含义如下：
//	*                           RTDB_E_DATA        标签点快照改变
//	*                           RTDB_E_DISCONNECT  订阅客户端与数据库网络断开
//	*                           RTDB_E_RECOVERY    订阅客户端与数据库网络及订阅恢复
//	*                           RTDB_E_CHANGED     客户端修改订阅标签点信息，即通过rtdbs_change_subscribe_snapshots修改订阅信息的结果
//	*                         handle 产生订阅回掉的连接句柄，调用rtdbs_subscribe_snapshots_ex时的handle参数
//	*                         param  用户自定义参数，调用rtdbs_subscribe_snapshots_ex时的param参数
//	*                         count  event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示ids，datetimes,values等的数量
//	*                                event_type为其它值时，count值为0
//	*                         ids    event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示快照改变的标签点ID，数量由count指定
//	*                                event_type为其它值时，ids值为NULL
//	*                         datetimes event_type为RTDB_E_DATA时表示快照时间，数量由count指定
//	*                                   event_type为RTDB_E_CHANGED时表示changed_types，即通过rtdbs_change_subscribe_snapshots传入的changed_types
//	*                                   event_type为其它值时，datetimes值为NULL
//	*                         ms     event_type为RTDB_E_DATA时表示快照的毫秒，数量由count指定
//	*                                event_type为其它值时，ms值为NULL
//	*                         values event_type为RTDB_E_DATA时表示浮点数据类型快照值，数量由count指定
//	*                                event_type为RTDB_E_CHANGED时，表示delta_values，即通过rtdbs_change_subscribe_snapshots传入的delta_values
//	*                                event_type为其它值时，values值为NULL
//	*                         states event_type为RTDB_E_DATA时表示整形数据类型快照值，数量由count指定
//	*                                event_type为RTDB_E_CHANGED时，表示delta_states，通过rtdbs_change_subscribe_snapshots传入的delta_states
//	*                                event_type为其它值时，states值为NULL
//	*                         qualities event_type为RTDB_E_DATA时表示快照质量码，数量由count指定
//	*                                event_type为其它值时，qualities值为NULL
//	*                         errors event_type为RTDB_E_DATA时表示快照错误码，数量由count指定
//	*                                event_type为RTDB_E_CHANGED时，表示修改结果对应的错误码，数量由count指定
//	*                                event_type为其它值时，errors值为NULL
//	* \param errors         无符号整型数组，输出，
//	*                           写入实时数据的返回值列表，参考rtdb_error.h
//	* \remark delta_values和delta_states可以为空指针，表示不设置容差值。 只有两个参数都不为空时，设置容差值才会生效。
//	*           用户须保证 ids、delta_values、delta_states、errors 的长度与 count 一致
//	*           用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
//	*           否则返回 RtE_OTHER_SDK_DOING 错误。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_subscribe_delta_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_float64* delta_values, const rtdb_int64* delta_states, rtdb_uint32 options, void* param, rtdbs_snaps_event_ex64 callback, rtdb_error* errors)
// */
func RawRtdbsSubscribeDeltaSnapshots64Warp() {}

// /*
//
//	*
//	* \brief 批量修改订阅标签点信息
//	*
//	* \param handle         连接句柄
//	* \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
//	*                           输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
//	* \param ids            整型数组，输入，标签点标识列表。
//	* \param delta_values   double型数组，输入，订阅浮点类型标签点的容差值，变化超过设置的容差值才会推送
//	* \param delta_values   整型数组，输入，订阅整型标签点的容差值，变化超过设置的容差值才会推送
//	* \param changed_types  整型数组，输入，修改类型，参考RTDB_SUBSCRIBE_CHANGE_TYPE
//	* \param errors         异步调用，保留参数，暂时不启用
//	* \remark   用户须保证 ids、delta_values、delta_states、errors 的长度与 count 一致。
//	*               可以同时添加、修改、删除订阅的标签点信息，
//	*               delta_values和delta_states，可以为空指针，为空，则表示不设置容差值，即写入新数据即推送
//	*               只有delta_values和delta_states都不为空时，设置的容差值才有效。
//	*               用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
//	*               否则返回 RtE_OTHER_SDK_DOING 错误。
//	*               此方法是异步方法，当网络中断等异常情况时，会通过方法的返回值返回错误，参考rtdb_error.h。
//	*               当方法返回值为RtE_OK时，表示已经成功发送给数据库，但是并没有等待修改结果。
//	*               数据库的修改结果，会异步通知给api的回调函数，通过rtdbs_snaps_event_ex的RTDB_E_CHANGED事件通知修改结果
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_change_subscribe_snapshots_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_float64* delta_values, const rtdb_int64* delta_states, const rtdb_int32* changed_types, rtdb_error* errors)
// */
func RawRtdbsChangeSubscribeSnapshotsWarp() {}

// /*
//
//	*
//	* \brief 取消标签点快照更改通知订阅
//	*
//	* \param handle    连接句柄
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_cancel_subscribe_snapshots_warp(rtdb_int32 handle)
// */
func RawRtdbsCancelSubscribeSnapshotsWarp() {}

// /*
//   - 命名：rtdbs_get_named_type_snapshot32
//   - 功能：获取自定义类型测点的单个快照
//   - 参数：
//   - [handle]    连接句柄
//   - [id]        整型，输入，标签点标识
//   - [datetime]  整型，输出，实时数值时间列表,
//   - 表示距离1970年1月1日08:00:00的秒数
//   - [ms]        短整型，输出，实时数值时间列表，
//   - 对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
//   - [object]    字节型数组，输出，实时自定义类型标签点的数值
//   - [length]    短整型，输入/输出，自定义类型标签点的数值长度
//   - [quality]   短整型，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_get_named_type_snapshot64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type* datetime, rtdb_subtime_type* subtime, void* object, rtdb_length_type* length, rtdb_int16* quality)
// */
func RawRtdbsGetNamedTypeSnapshot64Warp() {}

// /*
//   - 命名：rtdbs_get_named_type_snapshots32
//   - 功能：批量获取自定义类型测点的快照
//   - 参数：
//   - [handle]    连接句柄
//   - [count]     整型，输入/输出，标签点个数，
//   - 输入时表示 ids、datetimes、ms、objects、lengths、qualities、errors 的长度，
//   - 输出时表示成功获取实时值的标签点个数
//   - [ids]       整型数组，输入，标签点标识
//   - [datetimes] 整型数组，输出，实时数值时间列表,
//   - 表示距离1970年1月1日08:00:00的秒数
//   - [ms]        短整型数组，输出，实时数值时间列表，
//   - 对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
//   - [objects]   指针数组，输出，自定义类型标签点数值
//   - [lengths]   短整型数组，输入/输出，自定义类型标签点数值长度，
//   - 输入时表示对应的 objects 指针指向的缓冲区长度，
//   - 输出时表示实际得到的 objects 长度，如果 objects 的长度大于缓冲区长度，会被截断。
//   - [qualities] 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - [errors]    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_get_named_type_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, void* const* objects, rtdb_length_type* lengths, rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbsGetNamedTypeSnapshots64Warp() {}

// /*
//   - 命名：rtdbs_put_named_type_snapshot32
//   - 功能：写入单个自定义类型标签点的快照
//   - 参数：
//   - [handle]    连接句柄
//   - [id]        整型，输入，标签点标识
//   - [datetime]  整型，输入，实时数值时间列表,
//   - 表示距离1970年1月1日08:00:00的秒数
//   - [ms]        短整型，输入，实时数值时间列表，
//   - 对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
//   - [object]    void类型数组，输入，自定义类型标签点数值
//   - [length]    短整型，输入，自定义类型标签点数值长度，超过一个页大小数据将被截断。
//   - [quality]   短整型，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_put_named_type_snapshot64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, const void* object, rtdb_length_type length, rtdb_int16 quality)
// */
func RawRtdbsPutNamedTypeSnapshot64Warp() {}

// /*
// *
//   - 命名：rtdbs_put_named_type_snapshots32
//   - 功能：批量写入自定义类型标签点的快照
//   - 参数：
//   - [handle]    连接句柄
//   - [count]     整型，输入/输出，标签点个数，
//   - 输入时表示 ids、datetimes、ms、objects、lengths、qualities、errors 的长度，
//   - 输出时表示成功写入实时值的标签点个数
//   - [ids]       整型数组，输入，标签点标识
//   - [datetimes] 整型数组，输入，实时数值时间列表,
//   - 表示距离1970年1月1日08:00:00的秒数
//   - [ms]        短整型数组，输入，实时数值时间列表，
//   - 对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
//   - [objects]   void类型指针数组，输入，自定义类型标签点数值
//   - [lengths]   短整型数组，输入，自定义类型标签点数值长度，
//   - 表示对应的 objects 指针指向的缓冲区长度，超过一个页大小数据将被截断。
//   - [qualities] 短整型数组，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - [errors]    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_put_named_type_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const void* const* objects, const rtdb_length_type* lengths, const rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbsPutNamedTypeSnapshots64Warp() {}

// /*
//
//	*
//	* \brief 获取存档文件数量
//	*
//	* \param handle    连接句柄
//	* \param count     整型，输出，存档文件数量
//
// rtdb_error RTDBAPI_CALLRULE rtdba_get_archives_count_warp(rtdb_int32 handle, rtdb_int32 *count)
// */
func RawRtdbaGetArchivesCountWarp() {}

// /*
// * \brief 新建指定时间范围的历史存档文件并插入到历史数据库
// *
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名。
// * \param begin      整数，输入，起始时间，距离1970年1月1日08:00:00的秒数
// * \param end        整数，输入，终止时间，距离1970年1月1日08:00:00的秒数
// * \param mb_size    整型，输入，文件兆字节大小，单位为 MB。
// rtdb_error RTDBAPI_CALLRULE rtdba_create_ranged_archive64_warp(rtdb_int32 handle, const char* path, const char* file, rtdb_timestamp_type begin, rtdb_timestamp_type end, rtdb_int32 mb_size)
// */
func RawRtdbaCreateRangedArchive64Warp() {}

// /*
// * \brief 追加磁盘上的历史存档文件到历史数据库。
// *
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名，后缀名应为.rdf。
// * \param state      整型，输入，取值 RTDB_ACTIVED_ARCHIVE、RTDB_NORMAL_ARCHIVE、
// *                     RTDB_READONLY_ARCHIVE 之一，表示文件状态
// rtdb_error RTDBAPI_CALLRULE rtdba_append_archive_warp(rtdb_int32 handle, const char *path, const char *file, rtdb_int32 state)
// */
func RawRtdbaAppendArchiveWarp() {}

// /*
// * \brief 从历史数据库中移出历史存档文件。
// *
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名。
// rtdb_error RTDBAPI_CALLRULE rtdba_remove_archive_warp(rtdb_int32 handle, const char *path, const char *file)
// */
func RawRtdbaRemoveArchiveWarp() {}

// /*
//
//	*
//	* \brief 切换活动文件
//	*
//	* \param handle     连接句柄
//	* \remark 当前活动文件被写满时该事务被启动，
//	*        改变当前活动文件的状态为普通状态，
//	*        在所有历史数据存档文件中寻找未被使用过的
//	*        插入到前活动文件的右侧并改为活动状态，
//	*        若找不到则将前活动文件右侧的文件改为活动状态，
//	*        并将active_archive_指向该文件。该事务进行过程中，
//	*        用锁保证所有读写操作都暂停等待该事务完成。
//
// rtdb_error RTDBAPI_CALLRULE rtdba_shift_actived_warp(rtdb_int32 handle)
// */
func RawRtdbaShiftActivedWarp() {}

// /*
// *
//   - 命名：rtdba_get_archives
//   - 功能：获取存档文件的路径、名称、状态和最早允许写入时间。
//   - 参数：
//   - [handle]          连接句柄
//   - [paths]            字符串数组，输出，存档文件的目录路径，长度至少为 RTDB_PATH_SIZE。
//   - [files]            字符串数组，输出，存档文件的名称，长度至少为 RTDB_FILE_NAME_SIZE。
//   - [states]           整型数组，输出，取值 RTDB_INVALID_ARCHIVE、RTDB_ACTIVED_ARCHIVE、
//   - RTDB_NORMAL_ARCHIVE、RTDB_READONLY_ARCHIVE 之一，表示文件状态
//
// rtdb_error RTDBAPI_CALLRULE rtdba_get_archives_warp(rtdb_int32 handle, rtdb_int32* count, rtdb_path_string* paths, rtdb_filename_string* files, rtdb_int32 *states)
// */
func RawRtdbaGetArchivesWarp() {}

// /*
// *
//   - 功能：获取存档信息
//   - 参数：
//   - [handle]: in, 句柄
//   - [count]: out, 数量
//   - [paths]: out, 路径
//   - [files]: out, 文件
//   - [infos]: out, 存档信息
//   - [errors]: out, 错误
//
// rtdb_error RTDBAPI_CALLRULE rtdba_get_archives_info_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_path_string* const paths, const rtdb_filename_string* const files, RTDB_HEADER_PAGE *infos, rtdb_error* errors)
// */
func RawRtdbaGetArchivesInfoWarp() {}

// /*
// *
//   - 功能：获取存档的实时信息
//   - 参数：
//   - [handle]: in, 句柄
//   - [count]: out, 数量
//   - [paths]: out, 路径
//   - [files]: out, 文件
//   - [real_time_datas]: out, 实时数据
//   - [total_datas]: 总数
//   - [errors]: 错误
//
// rtdb_error RTDBAPI_CALLRULE rtdba_get_archives_perf_data_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_path_string* const paths, const rtdb_filename_string* const files, RTDB_ARCHIVE_PERF_DATA* real_time_datas, RTDB_ARCHIVE_PERF_DATA* total_datas, rtdb_error* errors)
// */
func RawRtdbaGetArchivesPerfDataWarp() {}

// /*
// *
//   - 功能：获取存档状态
//   - 参数：
//   - [handle]: in, 句柄
//   - [status]: out, 存档状态
//
// rtdb_error RTDBAPI_CALLRULE rtdba_get_archives_status_warp(rtdb_int32 handle, rtdb_error* status)
// */
func RawRtdbaGetArchivesStatusWarp() {}

// /*
// * \brief 获取存档文件及其附属文件的详细信息。
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名。
// * \param file_id    整型，输入，附属文件标识，0 表示获取主文件信息。
// * \param info       RTDB_HEADER_PAGE 结构，输出，存档文件信息
// rtdb_error RTDBAPI_CALLRULE rtdba_get_archive_info_warp(rtdb_int32 handle, const char *path, const char *file, rtdb_int32 file_id, RTDB_HEADER_PAGE *info)
// */
func RawRtdbaGetArchiveInfoWarp() {}

// /*
// * \brief 修改存档文件的可配置项。
// *
// * \param handle         连接句柄
// * \param path           字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file           字符串，输入，文件名。
// * \param rated_capacity 整型，输入，文件额定大小，单位为 MB。
// * \param ex_capacity    整型，输入，附属文件大小，单位为 MB。
// * \param auto_merge     短整型，输入，是否自动合并附属文件。
// * \param auto_arrange   短整型，输入，是否自动整理存档文件。
// * 备注: rated_capacity 与 ex_capacity 参数可为 0，表示不修改对应的配置项。
// rtdb_error RTDBAPI_CALLRULE rtdba_update_archive_warp(rtdb_int32 handle, const char *path, const char *file, rtdb_int32 rated_capacity, rtdb_int32 ex_capacity, rtdb_int16 auto_merge, rtdb_int16 auto_arrange)
// */
func RawRtdbaUpdateArchiveWarp() {}

// /*
// * \brief 整理存档文件，将同一标签点的数据块存放在一起以提高查询效率。
// *
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名。
// rtdb_error RTDBAPI_CALLRULE rtdba_arrange_archive_warp(rtdb_int32 handle, const char *path, const char *file)
// */
func RawRtdbaArrangeArchiveWarp() {}

// /*
// * \brief 为存档文件重新生成索引，用于恢复数据。
// *
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名。
// rtdb_error RTDBAPI_CALLRULE rtdba_reindex_archive_warp(rtdb_int32 handle, const char *path, const char *file)
// */
func RawRtdbaReindexArchiveWarp() {}

// /*
// * \brief 备份主存档文件及其附属文件到指定路径
// *
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名。
// * \param dest       字符串，输入，备份目录路径，必须以"\"或"/"结尾。
// rtdb_error RTDBAPI_CALLRULE rtdba_backup_archive_warp(rtdb_int32 handle, const char *path, const char *file, const char *dest)
// */
func RawRtdbaBackupArchiveWarp() {}

// /*
// *
// * 命名：rtdba_move_archive
// * 功能：将存档文件移动到指定目录
// * 参数：
// *        [handle]     连接句柄
// *        [path]       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// *        [file]       字符串，输入，文件名。
// *        [dest]       字符串，输入，移动目录路径，必须以"\"或"/"结尾。
// rtdb_error RTDBAPI_CALLRULE rtdba_move_archive_warp(rtdb_int32 handle, const char *path, const char *file, const char *dest)
// */
func RawRtdbaMoveArchiveWarp() {}

// /*
// *
// * 命名：rtdba_reindex_archive
// * 功能：为存档文件转换索引格式。
// * 参数：
// *        [handle]     连接句柄
// *        [path]       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// *        [file]       字符串，输入，文件名。
// rtdb_error RTDBAPI_CALLRULE rtdba_convert_index_warp(rtdb_int32 handle, const char *path, const char *file)
// */
func RawRtdbaConvertIndexWarp() {}

// /*
// *
//   - 命名：rtdba_query_big_job
//   - \brief 查询进程正在执行的后台任务类型、状态和进度
//     *
//   - \param handle     连接句柄
//   - \param process    所查询的进程代号，进程的标识参见枚举 RTDB_PROCESS_NAME,
//   - RTDB_PROCESS_HISTORIAN: 历史服务进程，具有以下任务类型：
//   - RTDB_MERGE: 合并附属文件到主文件;
//   - RTDB_ARRANGE: 整理存档文件;
//   - RTDB_REINDEX: 重建索引;
//   - RTDB_BACKUP: 备份;
//   - RTDB_REACTIVE: 激活为活动存档;
//   - RTDB_PROCESS_EQUATION: 方程式服务进程，具有以下任务类型：
//   - RTDB_COMPUTE: 历史计算;
//   - RTDB_PROCESS_BASE: 标签信息服务进程，具有以下任务类型：
//   - RTDB_UPDATE_TABLE: 修改表名称;
//   - RTDB_REMOVE_TABLE: 删除表;
//   - \param path       字符串，输出，长度至少为 RTDB_PATH_SIZE，
//   - 对以下任务，这个字段表示存档文件所在目录路径：
//   - RTDB_MERGE
//   - RTDB_ARRANGE
//   - RTDB_REINDEX
//   - RTDB_BACKUP
//   - RTDB_REACTIVE
//   - 对于以下任务，这个字段表示原来的表名：
//   - RTDB_UPDATE_TABLE
//   - RTDB_REMOVE_TABLE
//   - 对于其它任务不可用。
//   - \param file       字符串，输出，长度至少为 RTDB_FILE_NAME_SIZE，
//   - 对以下任务，这个字段表示存档文件名：
//   - RTDB_MERGE
//   - RTDB_ARRANGE
//   - RTDB_REINDEX
//   - RTDB_BACKUP
//   - RTDB_REACTIVE
//   - 对于以下任务，这个字段表示修改后的表名：
//   - RTDB_UPDATE_TABLE
//   - 对于其它任务不可用。
//   - \param job        短整型，输出，任务的标识参见枚举 RTDB_BIG_JOB_NAME。
//   - \param state      整型，输出，任务的执行状态，参考 rtdb_error.h
//   - \param end_time   整型，输出，任务的完成时间。
//   - \param progress   单精度浮点型，输出，任务的进度百分比。
//   - \remark path 及 file 参数可传空指针，对应的信息将不再返回。
//
// rtdb_error RTDBAPI_CALLRULE rtdba_query_big_job64_warp(rtdb_int32 handle, rtdb_int32 process, char* path, char* file, rtdb_int16* job, rtdb_int32* state, rtdb_timestamp_type* end_time, rtdb_float32* progress)
// */
func RawRtdbaQueryBigJob64Warp() {}

// /*
// *
//   - 命名：rtdba_cancel_big_job
//   - 功能：取消进程正在执行的后台任务
//   - 参数：
//   - [handle]     连接句柄
//   - [process]    所查询的进程代号，进程的标识参见枚举 RTDB_PROCESS_NAME,
//   - RTDB_PROCESS_HISTORIAN: 历史服务进程，具有以下任务类型：
//   - RTDB_MERGE: 合并附属文件到主文件;
//   - RTDB_ARRANGE: 整理存档文件;
//   - RTDB_REINDEX: 重建索引;
//   - RTDB_BACKUP: 备份;
//   - RTDB_REACTIVE: 激活为活动存档;
//   - RTDB_PROCESS_EQUATION: 方程式服务进程，具有以下任务类型：
//   - RTDB_COMPUTE: 历史计算;
//   - RTDB_PROCESS_BASE: 标签信息服务进程，具有以下任务类型：
//   - RTDB_UPDATE_TABLE: 修改表名称;
//   - RTDB_REMOVE_TABLE: 删除表;
//   - 备注：path 及 file 参数可传空指针，对应的信息将不再返回。
//
// rtdb_error RTDBAPI_CALLRULE rtdba_cancel_big_job_warp(rtdb_int32 handle, rtdb_int32 process)
// */
func RawRtdbaCancelBigJobWarp() {}

// /*
//
//	*
//	* \brief 获取单个标签点在一段时间范围内的存储值数量.
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
//	* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
//	* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
//	* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
//	* \param count         整型，输出，返回上述时间范围内的存储值数量
//	* \remark 由 datetime1、ms1 形成的时间可以大于 datetime2、ms2 表示的时间，
//	*        此时前者表示结束时间，后者表示起始时间。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_archived_values_count64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, rtdb_int32* count)
// */
func RawRtdbhArchivedValuesCount64Warp() {}

// /*
//
//	*
//	* \brief 获取单个标签点在一段时间范围内的真实的存储值数量.
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
//	* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
//	* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
//	* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
//	* \param count         整型，输出，返回上述时间范围内的存储值数量
//	* \remark 由 datetime1、ms1 形成的时间可以大于 datetime2、ms2 表示的时间，
//	*        此时前者表示结束时间，后者表示起始时间。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_archived_values_real_count64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, rtdb_int32* count)
// */
func RawRtdbhArchivedValuesRealCount64Warp() {}

// /*
//
//	*
//	* \brief 读取单个标签点一段时间内的储存数据
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param count         整型，输入/输出，
//	*                        输入时表示 datetimes、ms、values、states、qualities 的长度；
//	*                        输出时返回实际得到的数值个数
//	* \param datetimes     整型数组，输入/输出，
//	*                        输入时第一个元素表示起始时间秒数，
//	*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
//	*                        输出时表示对应的历史数值时间秒数。
//	* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	*                        则输入时第一个元素表示起始时间纳秒，
//	*                        最后一个元素表示结束时间纳秒；
//	*                        输出时表示对应的历史数值时间纳秒。
//	*                        否则忽略输入，输出时为 0。
//	* \param values        双精度浮点数数组，输出，历史浮点型数值列表
//	*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史存储值；否则为 0
//	* \param states        64 位整数数组，输出，历史整型数值列表，
//	*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史存储值；否则为 0
//	* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
//	*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
//	*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
//	*        最后一个元素表示开始时间。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
// */
func RawRtdbhGetArchivedValues64Warp() {}

// /*
//
//	*
//	* \brief 逆向读取单个标签点一段时间内的储存数据
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param count         整型，输入/输出，
//	*                        输入时表示 datetimes、ms、values、states、qualities 的长度；
//	*                        输出时返回实际得到的数值个数
//	* \param datetimes     整型数组，输入/输出，
//	*                        输入时第一个元素表示起始时间秒数，
//	*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
//	*                        输出时表示对应的历史数值时间秒数。
//	* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	*                        则输入时第一个元素表示起始时间纳秒，
//	*                        最后一个元素表示结束时间纳秒；
//	*                        输出时表示对应的历史数值时间纳秒。
//	*                        否则忽略输入，输出时为 0。
//	* \param values        双精度浮点数数组，输出，历史浮点型数值列表
//	*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史存储值；否则为 0
//	* \param states        64 位整数数组，输出，历史整型数值列表，
//	*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史存储值；否则为 0
//	* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
//	*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
//	*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
//	*        最后一个元素表示开始时间。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_values_backward64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
// */
func RawRtdbhGetArchivedValuesBackward64Warp() {}

// /*
//
//	*
//	* \brief 读取单个标签点一段时间内的坐标型储存数据
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param count         整型，输入/输出，
//	*                        输入时表示 datetimes、ms、x、y、qualities 的长度；
//	*                        输出时返回实际得到的数值个数
//	* \param datetimes     整型数组，输入/输出，
//	*                        输入时第一个元素表示起始时间秒数，
//	*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
//	*                        输出时表示对应的历史数值时间秒数。
//	* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	*                        则输入时第一个元素表示起始时间纳秒，
//	*                        最后一个元素表示结束时间纳秒；
//	*                        输出时表示对应的历史数值时间纳秒。
//	*                        否则忽略输入，输出时为 0。
//	* \param x             单精度浮点型数组，输出，浮点型横坐标历史数值列表
//	* \param y             单精度浮点型数组，输出，浮点型纵坐标历史数值列表
//	* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 用户须保证 datetimes、ms、x、y、qualities 的长度与 count 一致，
//	*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
//	*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
//	*        最后一个元素表示开始时间。
//	*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_coor_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float32* x, rtdb_float32* y, rtdb_int16* qualities)
// */
func RawRtdbhGetArchivedCoorValues64Warp() {}

// /*
//
//	*
//	* \brief 逆向读取单个标签点一段时间内的坐标型储存数据
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param count         整型，输入/输出，
//	*                        输入时表示 datetimes、ms、x、y、qualities 的长度；
//	*                        输出时返回实际得到的数值个数
//	* \param datetimes     整型数组，输入/输出，
//	*                        输入时第一个元素表示起始时间秒数，
//	*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
//	*                        输出时表示对应的历史数值时间秒数。
//	* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	*                        则输入时第一个元素表示起始时间纳秒，
//	*                        最后一个元素表示结束时间纳秒；
//	*                        输出时表示对应的历史数值时间纳秒。
//	*                        否则忽略输入，输出时为 0。
//	* \param x             单精度浮点型数组，输出，浮点型横坐标历史数值列表
//	* \param y             单精度浮点型数组，输出，浮点型纵坐标历史数值列表
//	* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 用户须保证 datetimes、ms、x、y、qualities 的长度与 count 一致，
//	*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
//	*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
//	*        最后一个元素表示开始时间。
//	*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_coor_values_backward64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float32* x, rtdb_float32* y, rtdb_int16* qualities)
// */
func RawRtdbhGetArchivedCoorValuesBackward64Warp() {}

// /*
//
//	*
//	* \brief 开始以分段返回方式读取一段时间内的储存数据
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
//	* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
//	* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
//	* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
//	* \param count         整型，输出，返回上述时间范围内的存储值数量
//	* \param batch_count   整型，输出，每次分段返回的长度，用于继续调用 rtdbh_get_next_archived_values 接口
//	* \remark 由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
//	*        此时前者表示结束时间，后者表示起始时间。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_values_in_batches64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, rtdb_int32* count, rtdb_int32* batch_count)
// */
func RawRtdbhGetArchivedValuesInBatches64Warp() {}

// /*
//
//	*
//	* \brief 分段读取一段时间内的储存数据
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param count         整形，输入/输出，
//	*                        输入时表示 datetimes、ms、values、states、qualities 的长度；
//	*                        输出时表示实际得到的存储值个数。
//	* \param datetimes     整型数组，输出，历史数值时间列表,
//	*                        表示距离1970年1月1日08:00:00的秒数
//	* \param ms            短整型数组，输出，历史数值时间列表，
//	*                        对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
//	* \param values        双精度浮点型数组，输出，历史浮点型数值列表，
//	*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的历史存储值；否则为 0
//	* \param states        64 位整型数组，输出，历史整型数值列表，
//	*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的历史存储值；否则为 0
//	* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 相符，
//	*        且 count 不能小于 rtdbh_get_archived_values_in_batches 接口中返回的 batch_count 的值，
//	*        当返回 RtE_BATCH_END 表示全部数据获取完毕。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_next_archived_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
// */
func RawRtdbhGetNextArchivedValues64Warp() {}

// /*
//
//	*
//	* \brief 获取单个标签点的单调递增时间序列历史插值。
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param count         整型，输入，表示 datetimes、ms、values、states、qualities 的长度。
//	* \param datetimes     整型数组，输入，表示需要的单调递增时间列表，
//	*                        为距离1970年1月1日08:00:00的秒数
//	* \param ms            短整型数组，输入，对于时间精度为纳秒的标签点，
//	*                        表示需要的单调递增时间对应的纳秒值；否则忽略。
//	* \param values        双精度浮点型数组，输出，历史浮点型数值列表，
//	*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的历史插值；否则为 0
//	* \param states        64 位整型数组，输出，历史整型数值列表，
//	*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的历史插值；否则为 0
//	* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 相符，
//	*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_timed_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 count, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
// */
func RawRtdbhGetTimedValues64Warp() {}

// /*
//
//	*
//	* \brief 获取单个坐标标签点的单调递增时间序列历史插值。
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param count         整型，输入，表示 datetimes、ms、x、y、qualities 的长度。
//	* \param datetimes     整型数组，输入，表示需要的单调递增时间列表，
//	*                        为距离1970年1月1日08:00:00的秒数
//	* \param ms            短整型数组，输入，对于时间精度为纳秒的标签点，
//	*                        表示需要的单调递增时间对应的纳秒值；否则忽略。
//	* \param x             单精度浮点型数组，输出，浮点型横坐标历史插值数值列表
//	* \param y             单精度浮点型数组，输出，浮点型纵坐标历史插值数值列表
//	* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 用户须保证 datetimes、ms、x、y、qualities 的长度与 count 相符，
//	*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_timed_coor_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 count, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, rtdb_float32* x, rtdb_float32* y, rtdb_int16* qualities)
// */
func RawRtdbhGetTimedCoorValues64Warp() {}

// /*
//
//	*
//	* \brief 获取单个标签点一段时间内等间隔历史插值
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param count         整型，输入/输出，
//	*                        输入时表示 datetimes、ms、values、states、qualities 的长度，
//	*                        即需要的插值个数；输出时返回实际得到的插值个数
//	* \param datetimes     整型数组，输入/输出，
//	*                        输入时第一个元素表示起始时间秒数，
//	*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
//	*                        输出时表示对应的历史数值时间秒数。
//	* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	*                        则输入时第一个元素表示起始时间纳秒，
//	*                        最后一个元素表示结束时间纳秒；
//	*                        输出时表示对应的历史数值时间纳秒。
//	*                        否则忽略输入，输出时为 0。
//	* \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
//	*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
//	* \param states        64 位整数数组，输出，整型历史插值数值列表，
//	*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
//	* \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
//	*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
//	*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
//	*        最后一个元素表示开始时间。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_interpo_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
// */
func RawRtdbhGetInterpoValues64Warp() {}

// /*
//
//	*
//	* \brief 读取单个标签点某个时刻之后一定数量的等间隔内插值替换的历史数值
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param interval      整型，输入，插值时间间隔，单位为纳秒
//	* \param count         整型，输入，表示 datetimes、ms、values、states、qualities 的长度，
//	*                        即需要的插值个数。
//	* \param datetimes     整型数组，输入/输出，
//	*                        输入时第一个元素表示起始时间秒数；
//	*                        输出时表示对应的历史数值时间秒数。
//	* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	*                        则输入时第一个元素表示起始时间纳秒；
//	*                        输出时表示对应的历史数值时间纳秒。
//	*                        否则忽略输入，输出时为 0。
//	* \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
//	*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
//	* \param states        64 位整数数组，输出，整型历史插值数值列表，
//	*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
//	* \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
//	*        在输入时，datetimes、ms 中至少应有一个元素用于存放起始时间。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_interval_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int64 interval, rtdb_int32 count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
// */
func RawRtdbhGetIntervalValues64Warp() {}

// /*
//
//	*
//	* \brief 读取单个标签点某个时间的历史数据
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT、RTDB_INTER 之一：
//	*                        RTDB_NEXT 寻找下一个最近的数据；
//	*                        RTDB_PREVIOUS 寻找上一个最近的数据；
//	*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
//	*                        RTDB_INTER 取指定时间的内插值数据。
//	* \param datetime      整型，输入/输出，输入时表示时间秒数；
//	*                        输出时表示实际取得的历史数值对应的时间秒数。
//	* \param ms            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
//	*                        否则忽略输入，输出时为 0。
//	* \param value         双精度浮点数，输出，浮点型历史数值
//	*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史值；否则为 0
//	* \param state         64 位整数，输出，整型历史数值，
//	*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史值；否则为 0
//	* \param quality       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_single_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 mode, rtdb_timestamp_type* datetime, rtdb_subtime_type* subtime, rtdb_float64* value, rtdb_int64* state, rtdb_int16* quality)
// */
func RawRtdbhGetSingleValue64Warp() {}

// /*
//
//	*
//	* \brief 读取单个标签点某个时间的坐标型历史数据
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT、RTDB_INTER 之一：
//	*                        RTDB_NEXT 寻找下一个最近的数据；
//	*                        RTDB_PREVIOUS 寻找上一个最近的数据；
//	*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
//	*                        RTDB_INTER 取指定时间的内插值数据。
//	* \param datetime      整型，输入/输出，输入时表示时间秒数；
//	*                        输出时表示实际取得的历史数值对应的时间秒数。
//	* \param ms            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
//	*                        否则忽略输入，输出时为 0。
//	* \param x             单精度浮点型，输出，横坐标历史数值
//	* \param y             单精度浮点型，输出，纵坐标历史数值
//	* \param quality       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 本接口只对数据类型为 RTDB_COOR 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_single_coor_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 mode, rtdb_timestamp_type* datetime, rtdb_subtime_type* subtime, rtdb_float32* x, rtdb_float32* y, rtdb_int16* quality)
// */
func RawRtdbhGetSingleCoorValue64Warp() {}

// /*
//
//	*
//	* \brief 读取单个标签点某个时间的二进制/字符串型历史数据
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT 之一：
//	*                        RTDB_NEXT 寻找下一个最近的数据；
//	*                        RTDB_PREVIOUS 寻找上一个最近的数据；
//	*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
//	* \param datetime      整型，输入/输出，输入时表示时间秒数；
//	*                        输出时表示实际取得的历史数值对应的时间秒数。
//	* \param ms            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
//	*                        否则忽略输入，输出时为 0。
//	* \param blob          字节型数组，输出，二进制/字符串历史值
//	* \param len           短整型，输入/输出，输入时表示 blob 的长度，
//	*                        输出时表示实际获取的二进制/字符串数据长度。
//	* \param quality       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_single_blob_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 mode, rtdb_timestamp_type* datetime, rtdb_subtime_type* subtime, rtdb_byte* blob, rtdb_length_type* len, rtdb_int16* quality)
// */
func RawRtdbhGetSingleBlobValue64Warp() {}

// /*
//   - \brief 读取单个标签点一段时间的二进制/字符串型历史数据
//     *
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
//   - \param count         整型，输入/输出，输入表示想要查询多少数据
//   - 输出表示实际查到多少数据
//   - \param datetime1     整型，输入，表示开始时间秒数；
//   - \param ms1           短整型，输入，指定的标签点时间精度为纳秒，
//   - 表示时间纳秒数；
//   - \param datetime2     整型，输入,表示结束时间秒数；
//   - \param ms2           短整型，输入，指定的标签点时间精度为纳秒，
//   - 表示时间纳秒数；
//   - \param datetimes     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
//   - \param ms            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
//   - 表示实际取得的历史数值时间纳秒数。
//   - \param lens          短整型数组，输入/输出，输入时表示 blob 的长度，
//   - 输出时表示实际获取的二进制/字符串数据长度。
//   - 当blobs为空指针时，表示只获取每条数据的长度，此时会忽略输入的lens
//   - \param blobs         字节型数组，输出，二进制/字符串历史值。可以设置为空指针，表示只获取每条数据的长度
//   - \param qualities     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_blob_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_length_type* lens, rtdb_byte* const* blobs, rtdb_int16* qualities)
// */
func RawRtdbhGetArchivedBlobValues64Warp() {}

// /*
//
//	*
//	* \brief 读取并模糊搜索单个标签点一段时间的二进制/字符串型历史数据
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
//	* \param count         整型，输入/输出，输入表示想要查询多少数据
//	*                        输出表示实际查到多少数据
//	* \param datetime1     整型，输入，表示开始时间秒数；
//	* \param ms1           短整型，输入，指定的标签点时间精度为纳秒，
//	*                        表示时间纳秒数；
//	* \param datetime2     整型，输入,表示结束时间秒数；
//	* \param ms2           短整型，输入，指定的标签点时间精度为纳秒，
//	*                        表示时间纳秒数；
//	* \param filter        字符串，输入，支持通配符的模糊搜索字符串，多个模糊搜索的条件通过空格分隔，只针对string类型有效
//	*                        当filter为空指针时，表示不进行过滤,
//	*                        限制最大长度为RTDB_EQUATION_SIZE-1，超过此长度会返回错误
//	* \param datetimes     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
//	* \param ms            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
//	*                        表示实际取得的历史数值时间纳秒数。
//	* \param lens          短整型数组，输入/输出，输入时表示 blob 的长度，
//	*                        输出时表示实际获取的二进制/字符串数据长度。
//	*                        当blobs为空指针时，表示只获取每条数据的长度，此时会忽略输入的lens
//	* \param blobs         字节型数组，输出，二进制/字符串历史值。可以设置为空指针，表示只获取每条数据的长度
//	* \param qualities     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_blob_values_filt64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, const char* filter, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_length_type* lens, rtdb_byte* const* blobs, rtdb_int16* qualities)
// */
func RawRtdbhGetArchivedBlobValuesFilt64Warp() {}

// /*
//   - \brief 读取单个标签点某个时间的datetime历史数据
//     *
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT 之一：
//   - RTDB_NEXT 寻找下一个最近的数据；
//   - RTDB_PREVIOUS 寻找上一个最近的数据；
//   - RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
//   - \param datetime      整型，输入/输出，输入时表示时间秒数；
//   - 输出时表示实际取得的历史数值对应的时间秒数。
//   - \param ms            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//   - 则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
//   - 否则忽略输入，输出时为 0。
//   - \param dtblob          字节型数组，输出，datetime历史值
//   - \param dtlen           短整型，输入/输出，输入时表示 blob 的长度，
//   - 输出时表示实际获取的datetime数据长度。
//   - \param quality       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \param type           短整型 datetime字符串的格式类型，默认为-1
//   - \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_single_datetime_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 mode, rtdb_timestamp_type* datetime, rtdb_subtime_type* subtime, rtdb_byte* dtblob, rtdb_length_type* dtlen, rtdb_int16* quality, rtdb_int16 type)
// */
func RawRtdbhGetSingleDatetimeValue64Warp() {}

// /*
//   - \brief 读取单个标签点一段时间的时间类型历史数据
//     *
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
//   - \param count         整型，输入/输出，输入表示想要查询多少数据
//   - 输出表示实际查到多少数据
//   - \param datetime1     整型，输入，表示开始时间秒数；
//   - \param ms1           短整型，输入，指定的标签点时间精度为纳秒，
//   - 表示时间纳秒数；
//   - \param datetime2     整型，输入,表示结束时间秒数；
//   - \param ms2           短整型，输入，指定的标签点时间精度为纳秒，
//   - 表示时间纳秒数；
//   - \param datetimes     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
//   - \param ms            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
//   - 表示实际取得的历史数值时间纳秒数。
//   - \param dtlens          短整型数组，输入/输出，输入时表示 blob 的长度，
//   - 输出时表示实际获取的时间数据长度。
//   - \param dtvalues         字节型数组，输出，时间历史值
//   - \param qualities     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \param type          短整型，输入，“yyyy-mm-dd hh:mm:ss.000”的type为1， 同样默认输入格式也为 “yyyy-mm-dd hh:mm:ss.000”
//   - “yyyy/mm/dd hh:mm:ss.000”的type为2
//   - \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_datetime_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_length_type* lens, rtdb_byte* const* blobs, rtdb_int16* qualities, rtdb_int16 type)
// */
func RawRtdbhGetArchivedDatetimeValues64Warp() {}

// /*
//   - \brief 写入批量标签点批量时间型历史存储数据
//     *
//   - \param handle        连接句柄
//   - \param count         整型，输入/输出，
//   - 输入时表示 ids、datetimes、ms、dtlens、dtvalues、qualities、errors 的长度，
//   - 即历史值个数；输出时返回实际写入的数值个数
//   - \param ids           整型数组，输入，标签点标识
//   - \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
//   - \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
//   - 表示对应的历史数值时间纳秒；否则忽略。
//   - \param dtvalues      字节型指针数组，输入，实时时间数值
//   - \param dtlens        短整型数组，输入，时间数值长度，
//   - 表示对应的 dtvalues 指针指向的缓冲区长度，超过一个页大小数据将被截断。
//   - \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
//   - \remark 用户须保证 ids、datetimes、ms、dtlens、dtvalues、qualities、errors 的长度与 count 一致，
//   - 本接口仅对数据类型为 RTDB_DATETIME 的标签点有效。
//   - 如果 datetimes、ms 标识的数据已经存在，其值将被替换。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_archived_datetime_values64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_byte* const* dtvalues, const rtdb_length_type* dtlens, const rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbhPutArchivedDatetimeValues64Warp() {}

// /*
//   - \brief 获取单个标签点一段时间内的统计值。
//     *
//   - \param handle            连接句柄
//   - \param id                整型，输入，标签点标识
//   - \param datetime1         整型，输入/输出，输入时表示起始时间秒数。
//   - 如果为 0，表示从存档中最早时间的数据开始进行统计。
//   - 输出时返回最大值的时间秒数。
//   - \param ms1               短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//   - 表示起始时间对应的纳秒，输出时表示最大值的时间纳秒数；否则忽略，返回值为 0
//   - \param datetime2         整型，输入/输出，输入时表示结束时间秒数。
//   - 如果为 0，表示统计到存档中最近时间的数据为止。
//   - 输出时返回最小值的时间秒数。
//   - \param ms2               短整型，如果 id 指定的标签点时间精度为纳秒，
//   - 表示结束时间对应的纳秒，输出时表示最小值的时间纳秒数；否则忽略，返回值为 0
//   - \param max_value         双精度浮点型，输出，表示统计时间段内的最大数值。
//   - \param min_value         双精度浮点型，输出，表示统计时间段内的最小数值。
//   - \param total_value       双精度浮点型，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
//   - \param calc_avg          双精度浮点型，输出，表示统计时间段内的算术平均值。
//   - \param power_avg         双精度浮点型，输出，表示统计时间段内的加权平均值。
//   - \remark 由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
//   - 此时前者表示结束时间，后者表示起始时间。
//   - 如果输出的最大值或最小值的时间戳秒值为 0，
//   - 则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
//   - 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_summary_data_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, RTDB_SUMMARY_DATA* summary_data)
// */
func RawRtdbhSummaryDataWarp() {}

// /*
// *
//   - 命名：rtdbh_summary_in_batches
//   - \brief 分批获取单一标签点一段时间内的统计值
//     *
//   - \param handle            连接句柄
//   - \param id                整型，输入，标签点标识
//   - \param count             整形，输入/输出，输入时表示 datatimes1、ms1、datatimes2、ms2、
//   - max_values、min_values、total_values、calc_avgs、power_avgs、errors 的长度，
//   - 即分段的个数；输出时表示成功取得统计值的分段个数。
//   - \param interval          64 位整型，输入，分段时间间隔，单位为纳秒。
//   - 如果为纳秒点，输入时间必须大于1纳秒，如果为秒级点，则必须大于1000000000纳秒。
//   - \param datetimes1        整型数组，输入/输出，输入时第一个元素表示起始时间秒数。
//   - 如果为 0，表示从存档中最早时间的数据开始进行统计。
//   - 输出时返回各个分段对应的最大值的时间秒数。
//   - \param ms1               短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//   - 第一个元素表示起始时间对应的纳秒，
//   - 输出时返回各个分段对应的最大值的时间纳秒数；否则忽略，返回值为 0
//   - \param datetimes2        整型数组，输入/输出，输入时第一个元素表示结束时间秒数。
//   - 如果为 0，表示统计到存档中最近时间的数据为止。
//   - 输出时返回各个分段对应的最小值的时间秒数。
//   - \param ms2               短整型数组，如果 id 指定的标签点时间精度为纳秒，
//   - 第一个元素表示结束时间对应的纳秒，
//   - 输出时返回各个分段对应的最小值的时间纳秒数；否则忽略，返回值为 0
//   - \param max_values        双精度浮点型数组，输出，表示统计时间段内的最大数值。
//   - \param min_values        双精度浮点型数组，输出，表示统计时间段内的最小数值。
//   - \param total_values      双精度浮点型数组，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
//   - \param calc_avgs         双精度浮点型数组，输出，表示统计时间段内的算术平均值。
//   - \param power_avgs        双精度浮点型数组，输出，表示统计时间段内的加权平均值。
//   - \param errors            无符号整型数组，输出，表示各个分段取得统计值的返回值。
//   - \remark 由 datetimes1[0]、ms1[0] 表示的时间可以大于 datetimes2[0]、ms2[0] 表示的时间，
//   - 此时前者表示结束时间，后者表示起始时间。
//   - 如果输出的最大值或最小值的时间戳秒值为 0，
//   - 则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
//   - 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_summary_data_in_batches_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_int64 interval, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, RTDB_SUMMARY_DATA* summary_datas, rtdb_error* errors)
// */
func RawRtdbhSummaryDataInBatchesWarp() {}

// /*
//
//	*
//	* \brief 获取单个标签点一段时间内用于绘图的历史数据
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param interval      整型，输入，时间区间数量，单位为个，
//	*                        一般会使用绘图的横轴(时间轴)所用屏幕像素数，
//	*                        该功能将起始至结束时间等分为 interval 个区间，
//	*                        并返回每个区间的第一个和最后一个数值、最大和最小数值、一条异常数值；
//	*                        故参数 count 有可能输出五倍于 interval 的历史值个数，
//	*                        所以推荐输入的 count 至少是 interval 的五倍。
//	* \param count         整型，输入/输出，输入时表示 datetimes、ms、values、states、qualities 的长度，
//	*                        即需要获取的最大历史值个数，输出时返回实际得到的历史值个数。
//	* \param datetimes     整型数组，输入/输出，
//	*                        输入时第一个元素表示起始时间秒数，
//	*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
//	*                        输出时表示对应的历史数值时间秒数。
//	* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	*                        则输入时第一个元素表示起始时间纳秒，
//	*                        最后一个元素表示结束时间纳秒；
//	*                        输出时表示对应的历史数值时间纳秒。
//	*                        否则忽略输入，输出时为 0。
//	* \param values        双精度浮点数数组，输出，浮点型历史值数值列表
//	*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史值；否则为 0
//	* \param states        64 位整数数组，输出，整型历史值数值列表，
//	*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史值；否则为 0
//	* \param qualities     短整型数组，输出，历史值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
//	*        在输入时，datetimes、ms 中至少应有一个元素，用以存放起始及结束时间。
//	*        第一个元素形成的时间可以大于最后一个元素形成的时间，
//	*        此时第一个元素表示结束时间，最后一个元素表示开始时间。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_plot_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 interval, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
// */
func RawRtdbhGetPlotValues64Warp() {}

// /*
// * \brief 获取批量标签点在某一时间的历史断面数据
// *
// * \param handle        连接句柄
// * \param ids           整型数组，输入，标签点标识列表
// * \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT、RTDB_INTER 之一：
// *                        RTDB_NEXT 寻找下一个最近的数据；
// *                        RTDB_PREVIOUS 寻找上一个最近的数据；
// *                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
// *                        RTDB_INTER 取指定时间的内插值数据。
// * \param count         整型，输入，表示 ids、datetimes、ms、values、states、qualities 的长度，即标签点个数。
// * \param datetimes     整型数组，输入/输出，输入时表示对应标签点的历史数值时间秒数，
// *                        输出时表示根据 mode 实际寻找到的数值时间秒数。
// * \param ms            短整型数组，输入/输出，对于时间精度为纳秒的标签点，
// *                        输入时表示历史数值时间纳秒数，存放相应的纳秒值，
// *                        输出时表示根据 mode 实际寻找到的数值时间纳秒数；否则忽略输入，输出时为 0。
// * \param values        双精度浮点数数组，输出，浮点型历史值数值列表
// *                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史值；否则为 0
// * \param states        64 位整数数组，输出，整型历史值数值列表，
// *                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
// *                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史值；否则为 0
// * \param qualities     短整型数组，输出，历史值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \param errors        无符号整型数组，输出，读取历史数据的返回值列表，参考rtdb_error.h
// * \remark 用户须保证 ids、datetimes、ms、values、states、qualities 的长度与 count 一致，
// *        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_cross_section_values64_warp(rtdb_int32 handle, const rtdb_int32* ids, rtdb_int32 mode, rtdb_int32 count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbhGetCrossSectionValues64Warp() {}

// /*
//   - 命名：rtdbh_get_archived_values_filt
//   - 功能：读取单个标签点在一段时间内经复杂条件筛选后的历史储存值
//     *
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param filter        字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
//   - 长度不得超过 RTDB_EQUATION_SIZE，为 0 则不进行条件筛选。
//   - \param count         整型，输入/输出，
//   - 输入时表示 datetimes、ms、values、states、qualities 的长度，
//   - 即需要的数值个数；输出时返回实际得到的数值个数。
//   - \param datetimes     整型数组，输入/输出，
//   - 输入时第一个元素表示起始时间秒数，
//   - 最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
//   - 输出时表示对应的历史数值时间秒数。
//   - \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//   - 则输入时第一个元素表示起始时间纳秒，
//   - 最后一个元素表示结束时间纳秒；
//   - 输出时表示对应的历史数值时间纳秒。
//   - 否则忽略输入，输出时为 0。
//   - \param values        双精度浮点数数组，输出，浮点型历史数值列表
//   - 对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史存储值；否则为 0
//   - \param states        64 位整数数组，输出，整型历史数值列表，
//   - 对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//   - RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史存储值；否则为 0
//   - \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
//   - 在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
//   - 大于最后一个元素形成的时间，此时第一个元素表示结束时间，
//   - 最后一个元素表示开始时间。
//   - 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_values_filt64_warp(rtdb_int32 handle, rtdb_int32 id, const char* filter, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
// */
func RawRtdbhGetArchivedValuesFilt64Warp() {}

// /*
//
//	*
//	* \brief 读取单个标签点某个时刻之后经复杂条件筛选后一定数量的等间隔内插值替换的历史数值
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param filter        字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
//	*                        长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
//	* \param interval      整型，输入，插值时间间隔，单位为纳秒
//	* \param count         整型，输入，表示 datetimes、ms、values、states、qualities 的长度，
//	*                        即需要的插值个数。
//	* \param datetimes     整型数组，输入/输出，
//	*                        输入时第一个元素表示起始时间秒数；
//	*                        输出时表示对应的历史数值时间秒数。
//	* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	*                        则输入时第一个元素表示起始时间纳秒；
//	*                        输出时表示对应的历史数值时间纳秒。
//	*                        否则忽略输入，输出时为 0。
//	* \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
//	*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
//	* \param states        64 位整数数组，输出，整型历史插值数值列表，
//	*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
//	* \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
//	*        在输入时，datetimes、ms 中至少应有一个元素用于表示起始时间。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_interval_values_filt64_warp(rtdb_int32 handle, rtdb_int32 id, const char* filter, rtdb_int64 interval, rtdb_int32 count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
// */
func RawRtdbhGetIntervalValuesFilt64Warp() {}

// /*
//
//	*
//	* \brief 获取单个标签点一段时间内经复杂条件筛选后的等间隔插值
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param filter        字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
//	*                        长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
//	* \param count         整型，输入/输出，
//	*                        输入时表示 datetimes、ms、values、states、qualities 的长度，
//	*                        即需要的插值个数；输出时返回实际得到的插值个数
//	* \param datetimes     整型数组，输入/输出，
//	*                        输入时第一个元素表示起始时间秒数，
//	*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
//	*                        输出时表示对应的历史数值时间秒数。
//	* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	*                        则输入时第一个元素表示起始时间纳秒，
//	*                        最后一个元素表示结束时间纳秒；
//	*                        输出时表示对应的历史数值时间纳秒。
//	*                        否则忽略输入，输出时为 0。
//	* \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
//	*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
//	* \param states        64 位整数数组，输出，整型历史插值数值列表，
//	*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
//	* \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
//	*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
//	*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
//	*        最后一个元素表示开始时间。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_interpo_values_filt64_warp(rtdb_int32 handle, rtdb_int32 id, const char* filter, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
// */
func RawRtdbhGetInterpoValuesFilt64Warp() {}

// /*
//
//	*
//	* \brief 获取单个标签点一段时间内经复杂条件筛选后的统计值
//	*
//	* \param handle            连接句柄
//	* \param id                整型，输入，标签点标识
//	* \param filter            字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
//	*                            长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
//	* \param datetime1         整型，输入/输出，输入时表示起始时间秒数。
//	*                            如果为 0，表示从存档中最早时间的数据开始进行统计。
//	*                            输出时返回最大值的时间秒数。
//	* \param ms1               短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	*                            表示起始时间对应的纳秒，输出时表示最大值的时间纳秒数；否则忽略，返回值为 0
//	* \param datetime2         整型，输入/输出，输入时表示结束时间秒数。
//	*                            如果为 0，表示统计到存档中最近时间的数据为止。
//	*                            输出时返回最小值的时间秒数。
//	* \param ms2               短整型，如果 id 指定的标签点时间精度为纳秒，
//	*                            表示结束时间对应的纳秒，输出时表示最小值的时间纳秒数；否则忽略，返回值为 0
//	* \param max_value         双精度浮点型，输出，表示统计时间段内的最大数值。
//	* \param min_value         双精度浮点型，输出，表示统计时间段内的最小数值。
//	* \param total_value       双精度浮点型，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
//	* \param calc_avg          双精度浮点型，输出，表示统计时间段内的算术平均值。
//	* \param power_avg         双精度浮点型，输出，表示统计时间段内的加权平均值。
//	* \remark 由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
//	*        此时前者表示结束时间，后者表示起始时间。
//	*        如果输出的最大值或最小值的时间戳秒值为 0，
//	*        则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
//	*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_summary_data_filt_warp(rtdb_int32 handle, rtdb_int32 id, const char* filter, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, RTDB_SUMMARY_DATA* summary_data)
// */
func RawRtdbhSummaryDataFiltWarp() {}

// /*
// *
//   - 命名：rtdbh_summary_filt_in_batches
//   - 功能：分批获取单一标签点一段时间内经复杂条件筛选后的统计值
//     *
//   - \param handle            连接句柄
//   - \param id                整型，输入，标签点标识
//   - \param filter            字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
//   - 长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
//   - \param count             整形，输入/输出，输入时表示 datatimes1、ms1、datatimes2、ms2、
//   - max_values、min_values、total_values、calc_avgs、power_avgs、errors 的长度，
//   - 即分段的个数；输出时表示成功取得统计值的分段个数。
//   - \param interval          64 位整型，输入，分段时间间隔，单位为纳秒。
//   - \param datetimes1        整型数组，输入/输出，输入时第一个元素表示起始时间秒数。
//   - 如果为 0，表示从存档中最早时间的数据开始进行统计。
//   - 输出时返回各个分段对应的最大值的时间秒数。
//   - \param ms1               短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//   - 第一个元素表示起始时间对应的纳秒，
//   - 输出时返回各个分段对应的最大值的时间纳秒数；否则忽略，返回值为 0
//   - \param datetimes2        整型数组，输入/输出，输入时第一个元素表示结束时间秒数。
//   - 如果为 0，表示统计到存档中最近时间的数据为止。
//   - 输出时返回各个分段对应的最小值的时间秒数。
//   - \param ms2               短整型数组，如果 id 指定的标签点时间精度为纳秒，
//   - 第一个元素表示结束时间对应的纳秒，
//   - 输出时返回各个分段对应的最小值的时间纳秒数；否则忽略，返回值为 0
//   - \param max_values        双精度浮点型数组，输出，表示统计时间段内的最大数值。
//   - \param min_values        双精度浮点型数组，输出，表示统计时间段内的最小数值。
//   - \param total_values      双精度浮点型数组，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
//   - \param calc_avgs         双精度浮点型数组，输出，表示统计时间段内的算术平均值。
//   - \param power_avgs        双精度浮点型数组，输出，表示统计时间段内的加权平均值。
//   - \param errors            无符号整型数组，输出，表示各个分段取得统计值的返回值。
//   - \remark 由 datetimes1[0]、ms1[0] 表示的时间可以大于 datetimes2[0]、ms2[0] 表示的时间，
//   - 此时前者表示结束时间，后者表示起始时间。
//   - 如果输出的最大值或最小值的时间戳秒值为 0，
//   - 则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
//   - 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_summary_data_filt_in_batches_warp(rtdb_int32 handle, rtdb_int32 id, const char* filter, rtdb_int32* count, rtdb_int64 interval, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, RTDB_SUMMARY_DATA* summary_datas, rtdb_error* errors)
// */
func RawRtdbhSummaryDataFiltInBatchesWarp() {}

// /*
//
//	*
//	* \brief 修改单个标签点某一时间的历史存储值.
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param datetime      整型，输入，时间秒数
//	* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
//	*                        表示时间纳秒数；否则忽略。
//	* \param value         双精度浮点数，输入，浮点型历史数值
//	*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放新的历史值；否则忽略
//	* \param state         64 位整数，输入，整型历史数值，
//	*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放新的历史值；否则忽略
//	* \param quality       短整型，输入，新的历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_update_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, rtdb_float64 value, rtdb_int64 state, rtdb_int16 quality)
// */
func RawRtdbhUpdateValue64Warp() {}

// /*
//   - \brief 修改单个标签点某一时间的历史存储值.
//     *
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param datetime      整型，输入，时间秒数
//   - \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
//   - 表示时间纳秒数；否则忽略。
//   - \param x             单精度浮点型，输入，新的横坐标历史数值
//   - \param y             单精度浮点型，输入，新的纵坐标历史数值
//   - \param quality       短整型，输入，新的历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 本接口仅对数据类型为 RTDB_COOR 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_update_coor_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, rtdb_float32 x, rtdb_float32 y, rtdb_int16 quality)
// */
func RawRtdbhUpdateCoorValue64Warp() {}

// /*
//
//	*
//	* \brief 删除单个标签点某个时间的历史存储值
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param datetime      整型，输入，时间秒数
//	* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
//	*                        表示时间纳秒数；否则忽略。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_remove_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime)
// */
func RawRtdbhRemoveValue64Warp() {}

// /*
//   - \brief 删除单个标签点一段时间内的历史存储值
//     *
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
//   - \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
//   - \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
//   - \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
//   - \param count         整形，输出，表示删除的历史值个数
//   - \remark 由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
//   - 此时前者表示结束时间，后者表示起始时间。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_remove_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, rtdb_int32* count)
// */
func RawRtdbhRemoveValues64Warp() {}

// /*
//
//	*
//	* \brief 写入单个标签点在某一时间的历史数据。
//	*
//	* \param handle        连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param datetime      整型，输入，时间秒数
//	* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
//	*                        表示时间纳秒数；否则忽略。
//	* \param value         双精度浮点数，输入，浮点型历史数值
//	*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放历史值；否则忽略
//	* \param state         64 位整数，输入，整型历史数值，
//	*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放历史值；否则忽略
//	* \param quality       短整型，输入，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//	*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_single_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, rtdb_float64 value, rtdb_int64 state, rtdb_int16 quality)
// */
func RawRtdbhPutSingleValue64Warp() {}

// /*
//
//	*
//	* \brief 写入单个标签点在某一时间的坐标型历史数据。
//	*
//	* \param handle              连接句柄
//	* \param id            整型，输入，标签点标识
//	* \param datetime      整型，输入，时间秒数
//	* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
//	*                        表示时间纳秒数；否则忽略。
//	* \param x             单精度浮点型，输入，横坐标历史数值
//	* \param y             单精度浮点型，输入，纵坐标历史数值
//	* \param quality       短整型，输入，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//	*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_single_coor_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, rtdb_float32 x, rtdb_float32 y, rtdb_int16 quality)
// */
func RawRtdbhPutSingleCoorValue64Warp() {}

// /*
//
//	*
//	* \brief 写入单个二进制/字符串标签点在某一时间的历史数据
//	*
//	* \param handle    连接句柄
//	* \param id        整型，输入，标签点标识
//	* \param datetime  整型，输入，数值时间列表,
//	*                    表示距离1970年1月1日08:00:00的秒数
//	* \param ms        短整型，输入，历史数值时间，
//	*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
//	* \param blob      字节型数组，输入，历史二进制/字符串数值
//	* \param len       短整型，输入，二进制/字符串数值长度，超过一个页大小数据将被截断。
//	* \param quality   短整型，输入，历史数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//	* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_single_blob_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, const rtdb_byte* blob, rtdb_length_type len, rtdb_int16 quality)
// */
func RawRtdbhPutSingleBlobValue64Warp() {}

// /*
//   - \brief 写入批量标签点批量历史存储数据
//     *
//   - \param handle        连接句柄
//   - \param count         整型，输入/输出，
//   - 输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
//   - 即历史值个数；输出时返回实际写入的数值个数
//   - \param ids           整型数组，输入，标签点标识，同一个标签点标识可以出现多次，
//   - 但它们的时间戳必需是递增的。
//   - \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
//   - \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
//   - 表示对应的历史数值时间纳秒；否则忽略。
//   - \param values        双精度浮点数数组，输入，浮点型历史数值列表
//   - 对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，表示相应的历史存储值；否则忽略
//   - \param states        64 位整数数组，输入，整型历史数值列表，
//   - 对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//   - RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，表示相应的历史存储值；否则忽略
//   - \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
//   - \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致，
//   - 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//   - 如果 datetimes、ms 标识的数据已经存在，其值将被替换。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_archived_values64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbhPutArchivedValues64Warp() {}

// /*
//   - \brief 写入批量标签点批量坐标型历史存储数据
//     *
//   - \param handle        连接句柄
//   - \param count         整型，输入/输出，
//   - 输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
//   - 即历史值个数；输出时返回实际写入的数值个数
//   - \param ids           整型数组，输入，标签点标识
//   - \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
//   - \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
//   - 表示对应的历史数值时间纳秒；否则忽略。
//   - \param x             单精度浮点型数组，输入，浮点型横坐标历史数值列表
//   - \param y             单精度浮点型数组，输入，浮点型纵坐标历史数值列表
//   - \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
//   - \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致，
//   - 本接口仅对数据类型为 RTDB_COOR 的标签点有效。
//   - 如果 datetimes、ms 标识的数据已经存在，其值将被替换。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_archived_coor_values64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float32* x, const rtdb_float32* y, const rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbhPutArchivedCoorValues64Warp() {}

// /*
//   - \brief 写入单个datetime标签点在某一时间的历史数据
//     *
//   - \param handle    连接句柄
//   - \param id        整型，输入，标签点标识
//   - \param datetime  整型，输入，数值时间列表,
//   - 表示距离1970年1月1日08:00:00的秒数
//   - \param ms        短整型，输入，历史数值时间，
//   - 对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
//   - \param blob      字节型数组，输入，历史datetime数值
//   - \param len       短整型，输入，datetime数值长度，超过一个页大小数据将被截断。
//   - \param quality   短整型，输入，历史数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_single_datetime_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, const rtdb_byte* blob, rtdb_length_type len, rtdb_int16 quality)
// */
func RawRtdbhPutSingleDatetimeValue64Warp() {}

// /*
//   - \brief 写入批量标签点批量字符串型历史存储数据
//     *
//   - \param handle        连接句柄
//   - \param count         整型，输入/输出，
//   - 输入时表示 ids、datetimes、ms、lens、blobs、qualities、errors 的长度，
//   - 即历史值个数；输出时返回实际写入的数值个数
//   - \param ids           整型数组，输入，标签点标识
//   - \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
//   - \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
//   - 表示对应的历史数值时间纳秒；否则忽略。
//   - \param blobs         字节型指针数组，输入，实时二进制/字符串数值
//   - \param lens          短整型数组，输入，二进制/字符串数值长度，
//   - 表示对应的 blobs 指针指向的缓冲区长度，超过一个页大小数据将被截断。
//   - \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
//   - \remark 用户须保证 ids、datetimes、ms、lens、blobs、qualities、errors 的长度与 count 一致，
//   - 本接口仅对数据类型为 RTDB_STRING、RTDB_BLOB 的标签点有效。
//   - 如果 datetimes、ms 标识的数据已经存在，其值将被替换。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_archived_blob_values64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_byte* const* blobs, const rtdb_length_type* lens, const rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbhPutArchivedBlobValues64Warp() {}

// /*
//   - \brief 将标签点未写满的补历史缓存页写入存档文件中。
//     *
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识。
//   - \param count         整型，输出，缓存页中数据个数。
//   - \remark 补历史缓存页写满后会自动写入存档文件中，不满的历史缓存页也会写入文件，
//   - 但会有一个时间延迟，在此期间此段数据可能查询不到，为了及时看到补历史的结果，
//   - 应在结束补历史后调用本接口。
//   - count 参数可为空指针，对应的信息将不再返回。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_flush_archived_values_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 *count)
// */
func RawRtdbhFlushArchivedValuesWarp() {}

// /*
//   - 命名：rtdbh_get_single_named_type_value32
//   - 功能：读取单个自定义类型标签点某个时间的历史数据
//   - 参数：
//   - [handle]        连接句柄
//   - [id]            整型，输入，标签点标识
//   - [mode]          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT 之一：
//   - RTDB_NEXT 寻找下一个最近的数据；
//   - RTDB_PREVIOUS 寻找上一个最近的数据；
//   - RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
//   - [datetime]      整型，输入/输出，输入时表示时间秒数；
//   - 输出时表示实际取得的历史数值对应的时间秒数。
//   - [ms]            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//   - 则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
//   - 否则忽略输入，输出时为 0。
//   - [object]        void数组，输出，自定义类型标签点历史值
//   - [length]        短整型，输入/输出，输入时表示 object 的长度，
//   - 输出时表示实际获取的自定义类型标签点数据长度。
//   - [quality]       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_single_named_type_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 mode, rtdb_timestamp_type* datetime, rtdb_subtime_type* subtime, void* object, rtdb_length_type* length, rtdb_int16* quality)
// */
func RawRtdbhGetSingleNamedTypeValue64Warp() {}

// /*
// *
//   - 命名：rtdbh_get_archived_named_type_values32
//   - 功能：连续读取自定义类型标签点的历史数据
//   - 参数：
//   - [handle]        连接句柄
//   - [id]            整型，输入，标签点标识
//   - RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
//   - [datetime1]     整型，输入，表示开始时间秒数；
//   - [ms1]           短整型，输入，指定的标签点时间精度为纳秒，
//   - 表示时间纳秒数；
//   - [datetime2]     整型，输入,表示结束时间秒数；
//   - [ms2]           短整型，输入，指定的标签点时间精度为纳秒，
//   - 表示时间纳秒数；
//   - [length]        短整型数组，输入，输入时表示 objects 的长度，
//   - [count]         整型，输入/输出，输入表示想要查询多少数据
//   - 输出表示实际查到多少数据
//   - [datetimes]     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
//   - [ms]            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
//   - 表示实际取得的历史数值时间纳秒数。
//   - [objects]       void类型数组，输出，自定义类型标签点历史值
//   - [qualities]     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_named_type_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, rtdb_length_type length, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, void* const* objects, rtdb_int16* qualities)
// */
func RawRtdbhGetArchivedNamedTypeValues64Warp() {}

// /*
// *
//   - 命名：rtdbh_put_single_named_type_value32
//   - 功能：写入自定义类型标签点的单个历史事件
//   - 参数：
//   - [handle]    连接句柄
//   - [id]        整型，输入，标签点标识
//   - [datetime]  整型，输入，数值时间列表,
//   - 表示距离1970年1月1日08:00:00的秒数
//   - [ms]        短整型，输入，历史数值时间，
//   - 对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
//   - [object]    void数组，输入，历史自定义类型标签点数值
//   - [length]    短整型，输入，自定义类型标签点数值长度，超过一个页大小数据将被截断。
//   - [quality]   短整型，输入，历史数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_single_named_type_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, const void* object, rtdb_length_type length, rtdb_int16 quality)
// */
func RawRtdbhPutSingleNamedTypeValue64Warp() {}

// /*
// *
//   - 命名：rtdbh_put_archived_named_type_values32
//   - 功能：批量补写自定义类型标签点的历史事件
//   - 参数：
//   - [handle]        连接句柄
//   - [count]         整型，输入/输出，
//   - 输入时表示 ids、datetimes、ms、lens、blobs、qualities、errors 的长度，
//   - 即历史值个数；输出时返回实际写入的数值个数
//   - [ids]           整型数组，输入，标签点标识
//   - [datetimes]     整型数组，输入，表示对应的历史数值时间秒数。
//   - [ms]            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
//   - 表示对应的历史数值时间纳秒；否则忽略。
//   - [objects]       void类型指针数组，输入，自定义类型标签点数值
//   - [lengths]       短整型数组，输入，自定义类型标签点数值长度，
//   - 表示对应的 objects 指针指向的缓冲区长度，超过一个页大小数据将被截断。
//   - [qualities]     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - [errors]        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
//   - 备注：用户须保证 ids、datetimes、ms、lens、objects、qualities、errors 的长度与 count 一致，
//   - 如果 datetimes、ms 标识的数据已经存在，其值将被替换。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_archived_named_type_values64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const void* const* objects, const rtdb_length_type* lengths, const rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbhPutArchivedNamedTypeValues64Warp() {}

// /*
//
//	*
//	* \brief 重算或补算批量计算标签点历史数据
//	*
//	* \param handle        连接句柄
//	* \param count         整型，输入/输出，
//	*                        输入时表示 ids、errors 的长度，
//	*                        即标签点个数；输出时返回成功开始计算的标签点个数
//	* \param flag          短整型，输入，不为 0 表示进行重算，删除时间范围内已经存在历史数据；
//	*                        为 0 表示补算，保留时间范围内已经存在历史数据，覆盖同时刻的计算值。
//	* \param datetime1     整型，输入，表示起始时间秒数。
//	* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
//	* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示计算直至存档中数据的最后时间
//	* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
//	* \param ids           整型数组，输入，标签点标识
//	* \param errors        无符号整型数组，输出，计算历史数据的返回值列表，参考rtdb_error.h
//	* \remark 用户须保证 ids、errors 的长度与 count 一致，本接口仅对带有计算扩展属性的标签点有效。
//	*        由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
//	*        此时前者表示结束时间，后者表示起始时间。
//
// rtdb_error RTDBAPI_CALLRULE rtdbe_compute_history64_warp(rtdb_int32 handle, rtdb_int32* count, rtdb_int16 flag, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, const rtdb_int32* ids, rtdb_error* errors)
// */
func RawRtdbeComputeHistory64Warp() {}

// /*
//   - 命名：rtdbe_get_equation_graph_count
//   - 功能：根据标签点 id 获取相关联方程式键值对数量
//   - 参数：
//   - [handle]   连接句柄
//   - [id]       整型，输入，标签点标识
//   - [flag]     枚举，输入，获取的拓扑图的关系
//   - [count]    整型，输入，拓扑图键值对数量
//   - 备注：键值对为数据结构，存储方程式涉及到的各标签点ID、及其父ID等
//   - 具体参考rtdbe_get_equation_graph_datas
//
// rtdb_error RTDBAPI_CALLRULE rtdbe_get_equation_graph_count_warp(rtdb_int32 handle, rtdb_int32 id, RTDB_GRAPH_FLAG flag, rtdb_int32 *count)
// */
func RawRtdbeGetEquationGraphCountWarp() {}

// /*
//   - 命名：rtdbe_get_equation_graph_datas
//   - 功能：根据标签点 id 获取相关联方程式键值对数据
//   - 参数：
//   - [handle]   连接句柄
//   - [id]       整型，输入，标签点标识
//   - [flag]     枚举，输入，获取的拓扑图的关系
//   - [count]    整型，输出
//   - 输入时，表示拓扑图键值对数量
//   - 输出时，表示实际获取到的拓扑图键值对数量
//   - [graph]    输出，GOLDE_GRAPH数据结构，拓扑图键值对信息
//   - 备注：键值对为数据结构，存储方程式涉及到的各标签点ID、及其父ID等
//
// rtdb_error RTDBAPI_CALLRULE rtdbe_get_equation_graph_datas_warp(rtdb_int32 handle, rtdb_int32 id, RTDB_GRAPH_FLAG flag, rtdb_int32 *count, RTDB_GRAPH *graph)
// */
func RawRtdbeGetEquationGraphDatasWarp() {}

// /*
// *
//   - 命名：rtdbp_get_perf_tags_count
//   - 功能：获取Perf服务中支持的性能计数点的数量
//   - 参数：
//   - [handle]   连接句柄
//   - [count]    整型，输出，表示实际获取到的Perf服务中支持的性能计数点的数量
//
// rtdb_error RTDBAPI_CALLRULE rtdbp_get_perf_tags_count_warp(rtdb_int32 handle, int* count)
// */
func RawRtdbpGetPerfTagsCountWarp() {}

// /*
// *
//   - 命名：rtdbp_get_perf_tags_info
//   - 功能：根据性能计数点ID获取相关的性能计数点信息
//   - 参数：
//   - [handle]   连接句柄
//   - [count]    整型，输入，输出
//   - 输入时，表示想要获取的性能计数点信息的数量，也表示tags_info，errors等的长度
//   - 输出时，表示实际获取到的性能计数点信息的数量
//   - [errors] 无符号整型数组，输出，获取性能计数点信息的返回值列表，参考rtdb_error.h
//   - 备注：用户须保证分配给 tags_info，errors 的空间与 count 相符
//
// rtdb_error RTDBAPI_CALLRULE rtdbp_get_perf_tags_info_warp(rtdb_int32 handle, rtdb_int32* count, RTDB_PERF_TAG_INFO* tags_info, rtdb_error* errors)
// */
func RawRtdbpGetPerfTagsInfoWarp() {}

// /*
// *
//   - 命名：rtdbp_get_perf_values
//   - 功能：批量读取性能计数点的当前快照数值
//   - 参数：
//   - [handle]    连接句柄
//   - [count]     整型，输入/输出，性能点个数，
//   - 输入时表示 perf_ids、datetimes、ms、values、states、qualities、errors 的长度，
//   - 输出时表示成功获取实时值的性能计数点个数
//   - [perf_ids]  整型数组，输入，性能计数点标识列表，参考RTDB_PERF_TAG_ID
//   - [datetimes] 整型数组，输出，实时数值时间列表,
//   - 表示距离1970年1月1日08:00:00的秒数
//   - [ms]        短整型数组，输出，实时数值时间列表，
//   - 对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
//   - [values]    双精度浮点型数组，输出，实时浮点型数值列表，
//   - 对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的快照值；否则为 0
//   - [states]    64 位整型数组，输出，实时整型数值列表，
//   - 对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//   - RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的快照值；否则为 0
//   - [qualities] 短整型数组，输出，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - [errors]    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
//   - 备注：用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
//
// rtdb_error RTDBAPI_CALLRULE rtdbp_get_perf_values64_warp(rtdb_int32 handle, rtdb_int32* count, int* perf_ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities, rtdb_error* errors)
// */
func RawRtdbpGetPerfValues64Warp() {}
