package rtdb_api

// #cgo CFLAGS: -DPNG_DEBUG=1 -I./cinclude
// #cgo CXXFLAGS: -std=c++11
// #include <stdlib.h>
// #include "api.h"
// #include "gofn.h"
import "C"
import (
	_ "embed"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"
)

//go:embed clibrary/linux_amd64/librtdbapi.so
var LinuxAmd64RtdbSo []byte

//go:embed clibrary/linux_arm64/librtdbapi.so
var LinuxArm64RtdbSo []byte

//go:embed clibrary/windows_amd32/rtdbapi.dll
var WindowsAmd32RtdbSo []byte

//go:embed clibrary/windows_amd64/rtdbapi.dll
var WindowsAmd64RtdbSo []byte

func init() {
	data := make([]byte, 0)
	name := ""
	if runtime.GOOS == "linux" {
		if runtime.GOARCH == "amd64" {
			data = LinuxAmd64RtdbSo
			name = "librtdb.so"
		} else if runtime.GOARCH == "arm64" {
			data = LinuxArm64RtdbSo
			name = "librtdb.so"
		} else {
			panic("不支持的平台，分支不可达")
		}
	} else if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
			data = WindowsAmd64RtdbSo
			name = "rtdbapi.dll"
		} else if runtime.GOARCH == "amd32" {
			data = WindowsAmd32RtdbSo
			name = "rtdbapi.dll"
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
	if runtime.GOOS == "linux" {
		cPath := C.CString(path)
		defer C.free(unsafe.Pointer(cPath))
		C.load_library_linux(cPath)
	} else {
		cPath, err := UTF16PtrFromString(path)
		if err != nil {
			panic("字符转换失败：" + err.Error())
		}
		C.load_library_windows((*C.wchar_t)(unsafe.Pointer(cPath)))
	}
}

// Switch 开关
type Switch byte

const (
	// OFF 关闭
	OFF = Switch(0)

	// ON 开启
	ON = Switch(1)
)

// RtdbError 数据库错误
type RtdbError uint32

// IsOk 判断当前错误是否为RteOk
func (re RtdbError) IsOk() bool {
	return errors.Is(re, RteOk)
}

// Error 实现 error 接口
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

// GoError 将RtdbError转换成 error
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
type RtdbApiOption uint32

const (
	// RtdbApiOptionAutoReconn api 在连接中断后是否自动重连, 0 不重连；1 重连。默认为 0 不重连
	RtdbApiOptionAutoReconn = RtdbApiOption(C.RTDB_API_AUTO_RECONN)

	// RtdbApiOptionConnTimeout api 连接超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为1000
	RtdbApiOptionConnTimeout = RtdbApiOption(C.RTDB_API_CONN_TIMEOUT)

	// RtdbApiOptionSendTimeout api 发送超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为1000
	RtdbApiOptionSendTimeout = RtdbApiOption(C.RTDB_API_SEND_TIMEOUT)

	// RtdbApiOptionRecvTimeout api 接收超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为60000
	RtdbApiOptionRecvTimeout = RtdbApiOption(C.RTDB_API_RECV_TIMEOUT)

	// RtdbApiOptionUserTimeout api TCP_USER_TIMEOUT超时值设置（单位：毫秒），默认为10000，Linux内核2.6.37以上有效
	RtdbApiOptionUserTimeout = RtdbApiOption(C.RTDB_API_USER_TIMEOUT)

	// RtdbApiOptionDefaultPrecision api 默认的时间戳精度，当使用旧版相关的api，以及新版api中未设置时间戳精度时，则使用此默认时间戳精度。 默认为毫秒精度
	RtdbApiOptionDefaultPrecision = RtdbApiOption(C.RTDB_API_DEFAULT_PRECISION)

	// RtdbApiOptionServerPrecision api 连接3.0数据库时，设置3.0数据库的时间戳精度，0表示毫秒精度，非0表示纳秒精度，默认为毫秒精度
	RtdbApiOptionServerPrecision = RtdbApiOption(C.RTDB_API_SERVER_PRECISION)
)

func (rao RtdbApiOption) Desc() string {
	switch rao {
	case RtdbApiOptionAutoReconn:
		return "api 在连接中断后是否自动重连, 0 不重连；1 重连。默认为 0 不重连"
	case RtdbApiOptionConnTimeout:
		return "api 连接超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为1000"
	case RtdbApiOptionSendTimeout:
		return "api 发送超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为1000"
	case RtdbApiOptionRecvTimeout:
		return "api 接收超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为60000"
	case RtdbApiOptionUserTimeout:
		return "api TCP_USER_TIMEOUT超时值设置（单位：毫秒），默认为10000，Linux内核2.6.37以上有效"
	case RtdbApiOptionDefaultPrecision:
		return "api 默认的时间戳精度，当使用旧版相关的api，以及新版api中未设置时间戳精度时，则使用此默认时间戳精度。 默认为毫秒精度"
	case RtdbApiOptionServerPrecision:
		return "api 连接3.0数据库时，设置3.0数据库的时间戳精度，0表示毫秒精度，非0表示纳秒精度，默认为毫秒精度"
	default:
		return "未知的RtdbApiOption"
	}
}

// DatagramHandle 流句柄, 用于数据流订阅
type DatagramHandle struct {
	handle C.rtdb_datagram_handle
}

// ApiVersion 版本号
type ApiVersion struct {
	Major int32 // 主版本号
	Minor int32 // 次版本号
	Beta  int32 // 发布版本号
}

// PrivGroup 用户权限
type PrivGroup uint32

const (
	// PrivGroupRtdbRO 只读
	PrivGroupRtdbRO = PrivGroup(C.RTDB_RO)

	// PrivGroupRtdbDW 数据记录
	PrivGroupRtdbDW = PrivGroup(C.RTDB_DW)

	// PrivGroupRtdbTA 标签点表管理员
	PrivGroupRtdbTA = PrivGroup(C.RTDB_TA)

	// PrivGroupRtdbSA 数据库管理员
	PrivGroupRtdbSA = PrivGroup(C.RTDB_SA)
)

func (pg PrivGroup) Desc() string {
	switch pg {
	case PrivGroupRtdbRO:
		return "只读"
	case PrivGroupRtdbDW:
		return "数据记录"
	case PrivGroupRtdbTA:
		return "标签点表管理员"
	case PrivGroupRtdbSA:
		return "数据库管理员"
	default:
		return "未知权限"
	}
}

// ConnectHandle 连接句柄, 用于描述一个 API库 与 数据库 之间的连接
type ConnectHandle int32

// SocketHandle socket连接句柄
type SocketHandle int32

// RtdbParam 查询系统参数时对应的索引
type RtdbParam int32

const (
	// RtdbParamTableFile 标签点表文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamTableFile = RtdbParam(C.RTDB_PARAM_TABLE_FILE)

	// RtdbParamBaseFile 基本标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamBaseFile = RtdbParam(C.RTDB_PARAM_BASE_FILE)

	// RtdbParamScanFile 采集标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamScanFile = RtdbParam(C.RTDB_PARAM_SCAN_FILE)

	// RtdbParamCalcFile 计算标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamCalcFile = RtdbParam(C.RTDB_PARAM_CALC_FILE)

	// RtdbParamSnapFile 标签点快照文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamSnapFile = RtdbParam(C.RTDB_PARAM_SNAP_FILE)

	// RtdbParamLicFile 协议文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamLicFile = RtdbParam(C.RTDB_PARAM_LIC_FILE)

	// RtdbParamHisFile 历史信息文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamHisFile = RtdbParam(C.RTDB_PARAM_HIS_FILE)

	// RtdbParamLogDir 服务器端日志文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamLogDir = RtdbParam(C.RTDB_PARAM_LOG_DIR)

	// RtdbParamUserFile 用户权限信息文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamUserFile = RtdbParam(C.RTDB_PARAM_USER_FILE)

	// RtdbParamServerFile 网络服务进程与其它进程交互所用的共享内存文件，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamServerFile = RtdbParam(C.RTDB_PARAM_SERVER_FILE)

	// RtdbParamEqautionFile 方程式服务进程与其它进程交互所用的共享内存文件，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamEqautionFile = RtdbParam(C.RTDB_PARAM_EQAUTION_FILE)

	// RtdbParamArvPagesFile 历史数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamArvPagesFile = RtdbParam(C.RTDB_PARAM_ARV_PAGES_FILE)

	// RtdbParamArvexPagesFile 补历史数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamArvexPagesFile = RtdbParam(C.RTDB_PARAM_ARVEX_PAGES_FILE)

	// RtdbParamArvexPagesBlobFile 补历史数据blob、str缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamArvexPagesBlobFile = RtdbParam(C.RTDB_PARAM_ARVEX_PAGES_BLOB_FILE)

	// RtdbParamAuthFile 信任连接段信息文件路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamAuthFile = RtdbParam(C.RTDB_PARAM_AUTH_FILE)

	// RtdbParamRecycledBaseFile 可回收基本标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamRecycledBaseFile = RtdbParam(C.RTDB_PARAM_RECYCLED_BASE_FILE)

	// RtdbParamRecycledScanFile 可回收采集标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamRecycledScanFile = RtdbParam(C.RTDB_PARAM_RECYCLED_SCAN_FILE)

	// RtdbParamRecycledCalcFile 可回收计算标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamRecycledCalcFile = RtdbParam(C.RTDB_PARAM_RECYCLED_CALC_FILE)

	// RtdbParamAutoBackupPath 自动备份目的地全路径，必须以“\”或“/”结束，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamAutoBackupPath = RtdbParam(C.RTDB_PARAM_AUTO_BACKUP_PATH)

	// RtdbParamServerSenderIp 镜像发送地址，字符串最大长度为 RTDB_MAX_HOSTNAME_SIZE
	RtdbParamServerSenderIp = RtdbParam(C.RTDB_PARAM_SERVER_SENDER_IP)

	// RtdbParamBlacklistFile 连接黑名单文件路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamBlacklistFile = RtdbParam(C.RTDB_PARAM_BLACKLIST_FILE)

	// RtdbParamDbVersion 数据库版本
	RtdbParamDbVersion = RtdbParam(C.RTDB_PARAM_DB_VERSION)

	// RtdbParamLicUser 授权单位
	RtdbParamLicUser = RtdbParam(C.RTDB_PARAM_LIC_USER)

	// RtdbParamLicType 授权方式
	RtdbParamLicType = RtdbParam(C.RTDB_PARAM_LIC_TYPE)

	// RtdbParamIndexDir 索引文件存放目录，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamIndexDir = RtdbParam(C.RTDB_PARAM_INDEX_DIR)

	// RtdbParamMirrorBufferPath 镜像缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamMirrorBufferPath = RtdbParam(C.RTDB_PARAM_MIRROR_BUFFER_PATH)

	// RtdbParamMirrorExBufferPath 补写镜像缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamMirrorExBufferPath = RtdbParam(C.RTDB_PARAM_MIRROR_EX_BUFFER_PATH)

	// RtdbParamEqautionPathFile 方程式长度超过规定长度时进行保存的文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamEqautionPathFile = RtdbParam(C.RTDB_PARAM_EQAUTION_PATH_FILE)

	// RtdbParamTagsFile 标签点关键属性文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamTagsFile = RtdbParam(C.RTDB_PARAM_TAGS_FILE)

	// RtdbParamRecycledSnapFile 可回收标签点快照事件文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamRecycledSnapFile = RtdbParam(C.RTDB_PARAM_RECYCLED_SNAP_FILE)

	// RtdbParamSwapPageFile 历史数据交换页文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamSwapPageFile = RtdbParam(C.RTDB_PARAM_SWAP_PAGE_FILE)

	// RtdbParamPageAllocatorFile 活动存档数据页分配器文件全路径，字符串最大长度为 RTDB_MAX_PATH, 该系统配置项2.1版数据库在使用，3.0数据库已去掉，但为了保证系统选项索引号, 与2.1一致，此处不能去掉。便于java sdk统一调用
	RtdbParamPageAllocatorFile = RtdbParam(C.RTDB_PARAM_PAGE_ALLOCATOR_FILE)

	// RtdbParamNamedTypeFile 自定义类型配置信息全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamNamedTypeFile = RtdbParam(C.RTDB_PARAM_NAMED_TYPE_FILE)

	// RtdbParamStrblobMirrorPath BLOB/STRING镜像数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamStrblobMirrorPath = RtdbParam(C.RTDB_PARAM_STRBLOB_MIRROR_PATH)

	// RtdbParamStrblobMirrorExPath 补写BLOB/STRING镜像数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamStrblobMirrorExPath = RtdbParam(C.RTDB_PARAM_STRBLOB_MIRROR_EX_PATH)

	// RtdbParamBufferDir 临时数据缓存路径
	RtdbParamBufferDir = RtdbParam(C.RTDB_PARAM_BUFFER_DIR)

	// RtdbParamPoolCacheFlie 曲线池索引文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamPoolCacheFlie = RtdbParam(C.RTDB_PARAM_POOL_CACHE_FLIE)

	// RtdbParamPoolDataFileDir 曲线池缓存文件目录，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamPoolDataFileDir = RtdbParam(C.RTDB_PARAM_POOL_DATA_FILE_DIR)

	// RtdbParamArchiveFilePath 存档文件低速存储区路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamArchiveFilePath = RtdbParam(C.RTDB_PARAM_ARCHIVE_FILE_PATH)

	// RtdbParamLicVersionType 授权版本
	RtdbParamLicVersionType = RtdbParam(C.RTDB_PARAM_LIC_VERSION_TYPE)

	// RtdbParamAutoMovePath 自动移动目的地全路径，必须以“\”或“/”结束，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamAutoMovePath = RtdbParam(C.RTDB_PARAM_AUTO_MOVE_PATH)

	// RtdbParamReplicationBufferPath 双活：数据同步缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamReplicationBufferPath = RtdbParam(C.RTDB_PARAM_REPLICATION_BUFFER_PATH)

	// RtdbParamReplicationExBufferPath 双活：数据同步补写数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamReplicationExBufferPath = RtdbParam(C.RTDB_PARAM_REPLICATION_EX_BUFFER_PATH)

	// RtdbParamStrblobReplicationBufferPath 双活：数据同步BLOB/STRING数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamStrblobReplicationBufferPath = RtdbParam(C.RTDB_PARAM_STRBLOB_REPLICATION_BUFFER_PATH)

	// RtdbParamStrblobReplicationExBufferPath 双活：数据同步补写BLOB/STRING数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamStrblobReplicationExBufferPath = RtdbParam(C.RTDB_PARAM_STRBLOB_REPLICATION_EX_BUFFER_PATH)

	// RtdbParamReplicationGroupIp 双活：同步组地址，字符串最大长度为 RTDB_MAX_HOSTNAME_SIZE
	RtdbParamReplicationGroupIp = RtdbParam(C.RTDB_PARAM_REPLICATION_GROUP_IP)

	// RtdbParamArcFilenamePrefixWhenUsingDate 是否归档文件使用日期作为文件名
	RtdbParamArcFilenamePrefixWhenUsingDate = RtdbParam(C.RTDB_PARAM_ARC_FILENAME_PREFIX_WHEN_USING_DATE)

	// RtdbParamHotArchiveFilePath 存档文件高速存储区路径，字符串最大长度为 RTDB_MAX_PATH
	RtdbParamHotArchiveFilePath = RtdbParam(C.RTDB_PARAM_HOT_ARCHIVE_FILE_PATH)

	// RtdbParamLicTablesCount 协议中限定的标签点表数量
	RtdbParamLicTablesCount = RtdbParam(C.RTDB_PARAM_LIC_TABLES_COUNT)

	// RtdbParamLicTagsCount 协议中限定的所有标签点数量
	RtdbParamLicTagsCount = RtdbParam(C.RTDB_PARAM_LIC_TAGS_COUNT)

	// RtdbParamLicScanCount 协议中限定的采集标签点数量
	RtdbParamLicScanCount = RtdbParam(C.RTDB_PARAM_LIC_SCAN_COUNT)

	// RtdbParamLicCalcCount 协议中限定的计算标签点数量
	RtdbParamLicCalcCount = RtdbParam(C.RTDB_PARAM_LIC_CALC_COUNT)

	// RtdbParamLicArchicveCount 协议中限定的历史存档文件数量
	RtdbParamLicArchicveCount = RtdbParam(C.RTDB_PARAM_LIC_ARCHICVE_COUNT)

	// RtdbParamServerIpcSize 网络服务进程与其它进程进行交互所使用的共享内存池的字节尺寸（单位：B）
	RtdbParamServerIpcSize = RtdbParam(C.RTDB_PARAM_SERVER_IPC_SIZE)

	// RtdbParamEquationIpcSize 方程式服务进程与其它进程进行交互所使用的共享内存池的字节尺寸（单位：B）
	RtdbParamEquationIpcSize = RtdbParam(C.RTDB_PARAM_EQUATION_IPC_SIZE)

	// RtdbParamHashTableSize 标签点求余哈希表的尺寸
	RtdbParamHashTableSize = RtdbParam(C.RTDB_PARAM_HASH_TABLE_SIZE)

	// RtdbParamTagDeleteTimes 可整库删除标签点的次数
	RtdbParamTagDeleteTimes = RtdbParam(C.RTDB_PARAM_TAG_DELETE_TIMES)

	// RtdbParamServerPort 网络服务独立服务器端口
	RtdbParamServerPort = RtdbParam(C.RTDB_PARAM_SERVER_PORT)

	// RtdbParamServerSenderPort 网络服务镜像发送端口
	RtdbParamServerSenderPort = RtdbParam(C.RTDB_PARAM_SERVER_SENDER_PORT)

	// RtdbParamServerReceiverPort 网络服务镜像接收端口
	RtdbParamServerReceiverPort = RtdbParam(C.RTDB_PARAM_SERVER_RECEIVER_PORT)

	// RtdbParamServerMode 网络服务启动模式
	RtdbParamServerMode = RtdbParam(C.RTDB_PARAM_SERVER_MODE)

	// RtdbParamServerConnectionCount 协议中限定网络服务连接并发数量
	RtdbParamServerConnectionCount = RtdbParam(C.RTDB_PARAM_SERVER_CONNECTION_COUNT)

	// RtdbParamArvPagesNumber 历史数据缓存中的页数量
	RtdbParamArvPagesNumber = RtdbParam(C.RTDB_PARAM_ARV_PAGES_NUMBER)

	// RtdbParamArvexPagesNumber 补历史数据缓存中的页数量
	RtdbParamArvexPagesNumber = RtdbParam(C.RTDB_PARAM_ARVEX_PAGES_NUMBER)

	// RtdbParamExceptionAtServer 是否由服务器进行例外判定
	RtdbParamExceptionAtServer = RtdbParam(C.RTDB_PARAM_EXCEPTION_AT_SERVER)

	// RtdbParamArvPageRecycleDelay 历史数据缓存页回收延时（毫秒）
	RtdbParamArvPageRecycleDelay = RtdbParam(C.RTDB_PARAM_ARV_PAGE_RECYCLE_DELAY)

	// RtdbParamExArchiveSize 历史数据存档文件文件自动增长大小（单位：MB）
	RtdbParamExArchiveSize = RtdbParam(C.RTDB_PARAM_EX_ARCHIVE_SIZE)

	// RtdbParamArchiveBatchSize 历史存储值分段查询个数
	RtdbParamArchiveBatchSize = RtdbParam(C.RTDB_PARAM_ARCHIVE_BATCH_SIZE)

	// RtdbParamDatafilePagesize 系统数据文件页大小
	RtdbParamDatafilePagesize = RtdbParam(C.RTDB_PARAM_DATAFILE_PAGESIZE)

	// RtdbParamArvAsyncQueueNormalDoor 历史数据缓存队列中速归档区（单位：百分比）
	RtdbParamArvAsyncQueueNormalDoor = RtdbParam(C.RTDB_PARAM_ARV_ASYNC_QUEUE_NORMAL_DOOR)

	// RtdbParamIndexAlwaysInMemory 常驻内存的历史数据索引大小（单位：MB）
	RtdbParamIndexAlwaysInMemory = RtdbParam(C.RTDB_PARAM_INDEX_ALWAYS_IN_MEMORY)

	// RtdbParamDiskMinRestSize 最低可用磁盘空间（单位：MB）
	RtdbParamDiskMinRestSize = RtdbParam(C.RTDB_PARAM_DISK_MIN_REST_SIZE)

	// RtdbParamMinSizeOfArchive 历史存档文件和附属文件的最小尺寸（单位：MB）
	RtdbParamMinSizeOfArchive = RtdbParam(C.RTDB_PARAM_MIN_SIZE_OF_ARCHIVE)

	// RtdbParamDelayOfAutoMergeOrArrange 自动合并/整理最小延迟（单位：小时）
	RtdbParamDelayOfAutoMergeOrArrange = RtdbParam(C.RTDB_PARAM_DELAY_OF_AUTO_MERGE_OR_ARRANGE)

	// RtdbParamStartOfAutoMergeOrArrange 自动合并/整理开始时间（单位：点钟）
	RtdbParamStartOfAutoMergeOrArrange = RtdbParam(C.RTDB_PARAM_START_OF_AUTO_MERGE_OR_ARRANGE)

	// RtdbParamStopOfAutoMergeOrArrange 自动合并/整理停止时间（单位：点钟）
	RtdbParamStopOfAutoMergeOrArrange = RtdbParam(C.RTDB_PARAM_STOP_OF_AUTO_MERGE_OR_ARRANGE)

	// RtdbParamStartOfAutoBackup 自动备份开始时间（单位：点钟）
	RtdbParamStartOfAutoBackup = RtdbParam(C.RTDB_PARAM_START_OF_AUTO_BACKUP)

	// RtdbParamStopOfAutoBackup 自动备份停止时间（单位：点钟）
	RtdbParamStopOfAutoBackup = RtdbParam(C.RTDB_PARAM_STOP_OF_AUTO_BACKUP)

	// RtdbParamMaxLatencyOfSnapshot 允许服务器时间之后多少小时内的数据进入快照（单位：小时）
	RtdbParamMaxLatencyOfSnapshot = RtdbParam(C.RTDB_PARAM_MAX_LATENCY_OF_SNAPSHOT)

	// RtdbParamPageAllocatorReserveSize 活动页分配器预留大小（单位：KB）， 0 表示使用操作系统视图大小
	RtdbParamPageAllocatorReserveSize = RtdbParam(C.RTDB_PARAM_PAGE_ALLOCATOR_RESERVE_SIZE)

	// RtdbParamIncludeSnapshotInQuery 决定取样本值和统计值时，快照是否应该出现在查询结果中
	RtdbParamIncludeSnapshotInQuery = RtdbParam(C.RTDB_PARAM_INCLUDE_SNAPSHOT_IN_QUERY)

	// RtdbParamLicBlobCount 协议中限定的字符串或BLOB类型标签点数量
	RtdbParamLicBlobCount = RtdbParam(C.RTDB_PARAM_LIC_BLOB_COUNT)

	// RtdbParamMirrorBufferSize 镜像文件大小（单位：GB）
	RtdbParamMirrorBufferSize = RtdbParam(C.RTDB_PARAM_MIRROR_BUFFER_SIZE)

	// RtdbParamBlobArvexPagesNumber blob、str补历史的默认缓存页数量
	RtdbParamBlobArvexPagesNumber = RtdbParam(C.RTDB_PARAM_BLOB_ARVEX_PAGES_NUMBER)

	// RtdbParamMirrorEventQueueCapacity 镜像缓存队列容量
	RtdbParamMirrorEventQueueCapacity = RtdbParam(C.RTDB_PARAM_MIRROR_EVENT_QUEUE_CAPACITY)

	// RtdbParamNotifyNotEnoughSpace 提示磁盘空间不足，一旦启用，设置为ON，则通过API返回大错误码，否则只记录日志
	RtdbParamNotifyNotEnoughSpace = RtdbParam(C.RTDB_PARAM_NOTIFY_NOT_ENOUGH_SPACE)

	// RtdbParamArchiveFixedRange 历史数据存档文件的固定时间范围，默认为0表示不使用固定时间范围（单位：分钟）
	RtdbParamArchiveFixedRange = RtdbParam(C.RTDB_PARAM_ARCHIVE_FIXED_RANGE)

	// RtdbParamOneClinetMaxConnectionCount 单个客户端允许的最大连接数，默认为0表示不限制
	RtdbParamOneClinetMaxConnectionCount = RtdbParam(C.RTDB_PARAM_ONE_CLINET_MAX_CONNECTION_COUNT)

	// RtdbParamArvPagesCapacity 历史数据缓存所占字节大小，单位：字节
	RtdbParamArvPagesCapacity = RtdbParam(C.RTDB_PARAM_ARV_PAGES_CAPACITY)

	// RtdbParamArvexPagesCapacity 历史数据补写缓存所占字节大小，单位：字节
	RtdbParamArvexPagesCapacity = RtdbParam(C.RTDB_PARAM_ARVEX_PAGES_CAPACITY)

	// RtdbParamBlobArvexPagesCapacity blob、string类型标签点历史数据补写缓存所占字节大小，单位：字节
	RtdbParamBlobArvexPagesCapacity = RtdbParam(C.RTDB_PARAM_BLOB_ARVEX_PAGES_CAPACITY)

	// RtdbParamLockedPagesMem 指定分配给数据库用的内存大小，单位：MB
	RtdbParamLockedPagesMem = RtdbParam(C.RTDB_PARAM_LOCKED_PAGES_MEM)

	// RtdbParamLicRecycleCount 协议中回收站的容量
	RtdbParamLicRecycleCount = RtdbParam(C.RTDB_PARAM_LIC_RECYCLE_COUNT)

	// RtdbParamArchivedPolicy 快照数据和补写数据的归档策略
	RtdbParamArchivedPolicy = RtdbParam(C.RTDB_PARAM_ARCHIVED_POLICY)

	// RtdbParamNetworkIsolationAckByte 网络隔离装置ACK字节
	RtdbParamNetworkIsolationAckByte = RtdbParam(C.RTDB_PARAM_NETWORK_ISOLATION_ACK_BYTE)

	// RtdbParamEnableLogger 启用日志输出，0为不启用
	RtdbParamEnableLogger = RtdbParam(C.RTDB_PARAM_ENABLE_LOGGER)

	// RtdbParamLogEncode 启用日志加密，0为不启用
	RtdbParamLogEncode = RtdbParam(C.RTDB_PARAM_LOG_ENCODE)

	// RtdbParamLoginTry 启用登录失败次数验证，0为不启用
	RtdbParamLoginTry = RtdbParam(C.RTDB_PARAM_LOGIN_TRY)

	// RtdbParamUserLog 启用用户详细日志，0为不启用
	RtdbParamUserLog = RtdbParam(C.RTDB_PARAM_USER_LOG)

	// RtdbParamCoverWriteLog 启用日志覆盖写功能，0为不启用
	RtdbParamCoverWriteLog = RtdbParam(C.RTDB_PARAM_COVER_WRITE_LOG)

	// RtdbParamLicNamedTypeCount 协议中限定的自定义类型标签点数量
	RtdbParamLicNamedTypeCount = RtdbParam(C.RTDB_PARAM_LIC_NAMED_TYPE_COUNT)

	// RtdbParamMirrorReceiverThreadpoolSize 镜像接收线程数量
	RtdbParamMirrorReceiverThreadpoolSize = RtdbParam(C.RTDB_PARAM_MIRROR_RECEIVER_THREADPOOL_SIZE)

	// RtdbParamSnapshotUseArchiveInterface 按照补历史流程归档快照数据页
	RtdbParamSnapshotUseArchiveInterface = RtdbParam(C.RTDB_PARAM_SNAPSHOT_USE_ARCHIVE_INTERFACE)

	// RtdbParamNoArcdataWriteLog 归档无对应存档文件的数据时记录日志
	RtdbParamNoArcdataWriteLog = RtdbParam(C.RTDB_PARAM_NO_ARCDATA_WRITE_LOG)

	// RtdbParamPutArchiveThreadNum 补历史归档线程数
	RtdbParamPutArchiveThreadNum = RtdbParam(C.RTDB_PARAM_PUT_ARCHIVE_THREAD_NUM)

	// RtdbParamArvexDataArchivedThreshold 单次补写数据归档阈值
	RtdbParamArvexDataArchivedThreshold = RtdbParam(C.RTDB_PARAM_ARVEX_DATA_ARCHIVED_THRESHOLD)

	// RtdbParamSnapshotFlushBufferDelay 快照服务的共享缓存刷新到磁盘的周期
	RtdbParamSnapshotFlushBufferDelay = RtdbParam(C.RTDB_PARAM_SNAPSHOT_FLUSH_BUFFER_DELAY)

	// RtdbParamDataSpeed 查询时使用加速统计
	RtdbParamDataSpeed = RtdbParam(C.RTDB_PARAM_DATA_SPEED)

	// RtdbParamUseNewPlotAlgo 启用新的曲线算法
	RtdbParamUseNewPlotAlgo = RtdbParam(C.RTDB_PARAM_USE_NEW_PLOT_ALGO)

	// RtdbParamQueryThreadPoolSize 曲线查询线程池中线程数量
	RtdbParamQueryThreadPoolSize = RtdbParam(C.RTDB_PARAM_QUERY_THREAD_POOL_SIZE)

	// RtdbParamArchivedValues 使用查询线程池查询历史数据
	RtdbParamArchivedValues = RtdbParam(C.RTDB_PARAM_ARCHIVED_VALUES)

	// RtdbParamArchivedValuesCount 使用查询线程池查询历史数据的条数
	RtdbParamArchivedValuesCount = RtdbParam(C.RTDB_PARAM_ARCHIVED_VALUES_COUNT)

	// RtdbParamPoolUseFlag 启用曲线池
	RtdbParamPoolUseFlag = RtdbParam(C.RTDB_PARAM_POOL_USE_FLAG)

	// RtdbParamPoolOutLogFlag 输出曲线池日志
	RtdbParamPoolOutLogFlag = RtdbParam(C.RTDB_PARAM_POOL_OUT_LOG_FLAG)

	// RtdbParamPoolTimeUsePoolFlag 使用曲线池缓存计算插值
	RtdbParamPoolTimeUsePoolFlag = RtdbParam(C.RTDB_PARAM_POOL_TIME_USE_POOL_FLAG)

	// RtdbParamPoolMaxPointCount 曲线池的标签点容量
	RtdbParamPoolMaxPointCount = RtdbParam(C.RTDB_PARAM_POOL_MAX_POINT_COUNT)

	// RtdbParamPoolOneFileSavePointCount 曲线池每个数据文件的标签点容量
	RtdbParamPoolOneFileSavePointCount = RtdbParam(C.RTDB_PARAM_POOL_ONE_FILE_SAVE_POINT_COUNT)

	// RtdbParamPoolSaveMemorySize 曲线缓存退出时临时缓冲区大小
	RtdbParamPoolSaveMemorySize = RtdbParam(C.RTDB_PARAM_POOL_SAVE_MEMORY_SIZE)

	// RtdbParamPoolMinTimeUnitSeconds 曲线池缓存数据当前时间单位
	RtdbParamPoolMinTimeUnitSeconds = RtdbParam(C.RTDB_PARAM_POOL_MIN_TIME_UNIT_SECONDS)

	// RtdbParamPoolTimeUnitViewRate 曲线池查询数据最小时间单位显示系数
	RtdbParamPoolTimeUnitViewRate = RtdbParam(C.RTDB_PARAM_POOL_TIME_UNIT_VIEW_RATE)

	// RtdbParamPoolTimerIntervalSeconds 曲线池定时器刷新周期
	RtdbParamPoolTimerIntervalSeconds = RtdbParam(C.RTDB_PARAM_POOL_TIMER_INTERVAL_SECONDS)

	// RtdbParamPoolPerfTimerIntervalSeconds 曲线池性能计算点刷新周期
	RtdbParamPoolPerfTimerIntervalSeconds = RtdbParam(C.RTDB_PARAM_POOL_PERF_TIMER_INTERVAL_SECONDS)

	// RtdbParamArchiveInitFileSize 存档文件初始大小
	RtdbParamArchiveInitFileSize = RtdbParam(C.RTDB_PARAM_ARCHIVE_INIT_FILE_SIZE)

	// RtdbParamArchiveIncreaseMode 存档文件增长模式
	RtdbParamArchiveIncreaseMode = RtdbParam(C.RTDB_PARAM_ARCHIVE_INCREASE_MODE)

	// RtdbParamArchiveIncreaseSize 固定模式下文件增长大小
	RtdbParamArchiveIncreaseSize = RtdbParam(C.RTDB_PARAM_ARCHIVE_INCREASE_SIZE)

	// RtdbParamArchiveIncreasePercent 百分比模式下增长百分比
	RtdbParamArchiveIncreasePercent = RtdbParam(C.RTDB_PARAM_ARCHIVE_INCREASE_PERCENT)

	// RtdbParamAllowConvertSklToRbtIndex 跳跃链表转换到红黑树
	RtdbParamAllowConvertSklToRbtIndex = RtdbParam(C.RTDB_PARAM_ALLOW_CONVERT_SKL_TO_RBT_INDEX)

	// RtdbParamEarlyDataTime 冷数据时间
	RtdbParamEarlyDataTime = RtdbParam(C.RTDB_PARAM_EARLY_DATA_TIME)

	// RtdbParamEarlyIndexTime 自动转换索引时间
	RtdbParamEarlyIndexTime = RtdbParam(C.RTDB_PARAM_EARLY_INDEX_TIME)

	// RtdbParamArrangeRbtTime 整理存档文件时决定索引格式的时间轴
	RtdbParamArrangeRbtTime = RtdbParam(C.RTDB_PARAM_ARRANGE_RBT_TIME)

	// RtdbParamEnableBigData 将存档文件全部读取到内存中
	RtdbParamEnableBigData = RtdbParam(C.RTDB_PARAM_ENABLE_BIG_DATA)

	// RtdbParamAutoArrangePercent 自动整理存档文件时的实际使用率
	RtdbParamAutoArrangePercent = RtdbParam(C.RTDB_PARAM_AUTO_ARRANGE_PERCENT)

	// RtdbParamEarlyArrangeTime 自动整理存档文件的时间
	RtdbParamEarlyArrangeTime = RtdbParam(C.RTDB_PARAM_EARLY_ARRANGE_TIME)

	// RtdbParamMinAutoArrangeArcfilePercent 自动整理存档文件时的最小使用率
	RtdbParamMinAutoArrangeArcfilePercent = RtdbParam(C.RTDB_PARAM_MIN_AUTO_ARRANGE_ARCFILE_PERCENT)

	// RtdbParamArrangeArcWithMemory 在内存中整理存档文件
	RtdbParamArrangeArcWithMemory = RtdbParam(C.RTDB_PARAM_ARRANGE_ARC_WITH_MEMORY)

	// RtdbParamAraangeArcMaxMemPercent 整理存档文件最大内存使用率
	RtdbParamAraangeArcMaxMemPercent = RtdbParam(C.RTDB_PARAM_ARAANGE_ARC_MAX_MEM_PERCENT)

	// RtdbParamMaxDiskSpacePercent 磁盘最大使用率
	RtdbParamMaxDiskSpacePercent = RtdbParam(C.RTDB_PARAM_MAX_DISK_SPACE_PERCENT)

	// RtdbParamUseDispath windows 用 linux 已禁用,是否启用转发服务
	RtdbParamUseDispath = RtdbParam(C.RTDB_PARAM_USE_DISPATH)

	// RtdbParamUseSmartParam windows 用 linux 已禁用,是否使用推荐参数
	RtdbParamUseSmartParam = RtdbParam(C.RTDB_PARAM_USE_SMART_PARAM)

	// RtdbParamSubscribeSnapshotCount 单连接快照事件订阅个数
	RtdbParamSubscribeSnapshotCount = RtdbParam(C.RTDB_PARAM_SUBSCRIBE_SNAPSHOT_COUNT)

	// RtdbParamSubscribeQueueSize 订阅事件队列大小
	RtdbParamSubscribeQueueSize = RtdbParam(C.RTDB_PARAM_SUBSCRIBE_QUEUE_SIZE)

	// RtdbParamSubscribeTimeout 订阅事件超时时间
	RtdbParamSubscribeTimeout = RtdbParam(C.RTDB_PARAM_SUBSCRIBE_TIMEOUT)

	// RtdbParamMirrorCompressOnoff 镜像报文压缩是否打开
	RtdbParamMirrorCompressOnoff = RtdbParam(C.RTDB_PARAM_MIRROR_COMPRESS_ONOFF)

	// RtdbParamMirrorCompressType 镜像报文压缩类型
	RtdbParamMirrorCompressType = RtdbParam(C.RTDB_PARAM_MIRROR_COMPRESS_TYPE)

	// RtdbParamMirrorCompressMin 镜像报文压缩最小值
	RtdbParamMirrorCompressMin = RtdbParam(C.RTDB_PARAM_MIRROR_COMPRESS_MIN)

	// RtdbParamArchiveRollTime 存档文件滚动时间轴
	RtdbParamArchiveRollTime = RtdbParam(C.RTDB_PARAM_ARCHIVE_ROLL_TIME)

	// RtdbParamHandleTimeOut 连接超时断开，单位：秒
	RtdbParamHandleTimeOut = RtdbParam(C.RTDB_PARAM_HANDLE_TIME_OUT)

	// RtdbParamMoveArvTime 移动存档文件时决定移动存档的时间轴
	RtdbParamMoveArvTime = RtdbParam(C.RTDB_PARAM_MOVE_ARV_TIME)

	// RtdbParamUseNewInterpAlgo 启用新的插值算法
	RtdbParamUseNewInterpAlgo = RtdbParam(C.RTDB_PARAM_USE_NEW_INTERP_ALGO)

	// RtdbParamEnableReplication 启用双活，0为不启用，1为启用
	RtdbParamEnableReplication = RtdbParam(C.RTDB_PARAM_ENABLE_REPLICATION)

	// RtdbParamReplicationGroupPort 双活：同步组端口
	RtdbParamReplicationGroupPort = RtdbParam(C.RTDB_PARAM_REPLICATION_GROUP_PORT)

	// RtdbParamReplicationThreadSize 双活：同步线程数
	RtdbParamReplicationThreadSize = RtdbParam(C.RTDB_PARAM_REPLICATION_THREAD_SIZE)

	// RtdbParamForceArchiveIncompleteDataPageDelay 强制归档补历史缓存里面未满数据页的延迟时间
	RtdbParamForceArchiveIncompleteDataPageDelay = RtdbParam(C.RTDB_PARAM_FORCE_ARCHIVE_INCOMPLETE_DATA_PAGE_DELAY)

	// RtdbParamArchiveRollDiskPercentage 存档文件滚动存储空间百分比
	RtdbParamArchiveRollDiskPercentage = RtdbParam(C.RTDB_PARAM_ARCHIVE_ROLL_DISK_PERCENTAGE)

	// RtdbParamEnableIpv6 启用ipv6设置
	RtdbParamEnableIpv6 = RtdbParam(C.RTDB_PARAM_ENABLE_IPV6)

	// RtdbParamEnableUseArchivedValue 按条件获取历史值时，是否直接获取条件中点的历史值，0:获取插值，1:获取历史值
	RtdbParamEnableUseArchivedValue = RtdbParam(C.RTDB_PARAM_ENABLE_USE_ARCHIVED_VALUE)

	// RtdbParamTimestampType 获取服务器时间戳类型
	RtdbParamTimestampType = RtdbParam(C.RTDB_PARAM_TIMESTAMP_TYPE)

	// RtdbParamArcFilenameUsingDate 是否归档文件使用日期作为文件名
	RtdbParamArcFilenameUsingDate = RtdbParam(C.RTDB_PARAM_ARC_FILENAME_USING_DATE)

	// RtdbParamLogMaxSpace 日志文件占用的最大磁盘空间
	RtdbParamLogMaxSpace = RtdbParam(C.RTDB_PARAM_LOG_MAX_SPACE)

	// RtdbParamLogFileSize 单个日志文件大小
	RtdbParamLogFileSize = RtdbParam(C.RTDB_PARAM_LOG_FILE_SIZE)

	// RtdbParamIgnoreToWriteNoarcbuffer 是否丢弃补历史数据
	RtdbParamIgnoreToWriteNoarcbuffer = RtdbParam(C.RTDB_PARAM_IGNORE_TO_WRITE_NOARCBUFFER)

	// RtdbParamArchivesCountForCalc 统计存档文件平均大小的存档文件个数
	RtdbParamArchivesCountForCalc = RtdbParam(C.RTDB_PARAM_ARCHIVES_COUNT_FOR_CALC)

	// RtdbParamMaxBlobSize blob、str类型数据在数据库中允许的最大长度
	RtdbParamMaxBlobSize = RtdbParam(C.RTDB_PARAM_MAX_BLOB_SIZE)
)

func (rp RtdbParam) Desc() string {
	switch rp {
	case RtdbParamTableFile:
		return "标签点表文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamBaseFile:
		return "基本标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamScanFile:
		return "采集标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamCalcFile:
		return "计算标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamSnapFile:
		return "标签点快照文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamLicFile:
		return "协议文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamHisFile:
		return "历史信息文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamLogDir:
		return "服务器端日志文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamUserFile:
		return "用户权限信息文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamServerFile:
		return "网络服务进程与其它进程交互所用的共享内存文件，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamEqautionFile:
		return "方程式服务进程与其它进程交互所用的共享内存文件，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamArvPagesFile:
		return "历史数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamArvexPagesFile:
		return "补历史数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamArvexPagesBlobFile:
		return "补历史数据blob、str缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamAuthFile:
		return "信任连接段信息文件路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamRecycledBaseFile:
		return "可回收基本标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamRecycledScanFile:
		return "可回收采集标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamRecycledCalcFile:
		return "可回收计算标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamAutoBackupPath:
		return "自动备份目的地全路径，必须以“\\”或“/”结束，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamServerSenderIp:
		return "镜像发送地址，字符串最大长度为 RTDB_MAX_HOSTNAME_SIZE"
	case RtdbParamBlacklistFile:
		return "连接黑名单文件路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamDbVersion:
		return "数据库版本"
	case RtdbParamLicUser:
		return "授权单位"
	case RtdbParamLicType:
		return "授权方式"
	case RtdbParamIndexDir:
		return "索引文件存放目录，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamMirrorBufferPath:
		return "镜像缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamMirrorExBufferPath:
		return "补写镜像缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamEqautionPathFile:
		return "方程式长度超过规定长度时进行保存的文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamTagsFile:
		return "标签点关键属性文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamRecycledSnapFile:
		return "可回收标签点快照事件文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamSwapPageFile:
		return "历史数据交换页文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamPageAllocatorFile:
		return "活动存档数据页分配器文件全路径，字符串最大长度为 RTDB_MAX_PATH, 该系统配置项2.1版数据库在使用，3.0数据库已去掉，但为了保证系统选项索引号, 与2.1一致，此处不能去掉。便于java sdk统一调用"
	case RtdbParamNamedTypeFile:
		return "自定义类型配置信息全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamStrblobMirrorPath:
		return "BLOB/STRING镜像数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamStrblobMirrorExPath:
		return "补写BLOB/STRING镜像数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamBufferDir:
		return "临时数据缓存路径"
	case RtdbParamPoolCacheFlie:
		return "曲线池索引文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamPoolDataFileDir:
		return "曲线池缓存文件目录，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamArchiveFilePath:
		return "存档文件低速存储区路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamLicVersionType:
		return "授权版本"
	case RtdbParamAutoMovePath:
		return "自动移动目的地全路径，必须以“\\”或“/”结束，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamReplicationBufferPath:
		return "双活：数据同步缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamReplicationExBufferPath:
		return "双活：数据同步补写数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamStrblobReplicationBufferPath:
		return "双活：数据同步BLOB/STRING数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamStrblobReplicationExBufferPath:
		return "双活：数据同步补写BLOB/STRING数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamReplicationGroupIp:
		return "双活：同步组地址，字符串最大长度为 RTDB_MAX_HOSTNAME_SIZE"
	case RtdbParamArcFilenamePrefixWhenUsingDate:
		return "是否归档文件使用日期作为文件名"
	case RtdbParamHotArchiveFilePath:
		return "存档文件高速存储区路径，字符串最大长度为 RTDB_MAX_PATH"
	case RtdbParamLicTablesCount:
		return "协议中限定的标签点表数量"
	case RtdbParamLicTagsCount:
		return "协议中限定的所有标签点数量"
	case RtdbParamLicScanCount:
		return "协议中限定的采集标签点数量"
	case RtdbParamLicCalcCount:
		return "协议中限定的计算标签点数量"
	case RtdbParamLicArchicveCount:
		return "协议中限定的历史存档文件数量"
	case RtdbParamServerIpcSize:
		return "网络服务进程与其它进程进行交互所使用的共享内存池的字节尺寸（单位：B）"
	case RtdbParamEquationIpcSize:
		return "方程式服务进程与其它进程进行交互所使用的共享内存池的字节尺寸（单位：B）"
	case RtdbParamHashTableSize:
		return "标签点求余哈希表的尺寸"
	case RtdbParamTagDeleteTimes:
		return "可整库删除标签点的次数"
	case RtdbParamServerPort:
		return "网络服务独立服务器端口"
	case RtdbParamServerSenderPort:
		return "网络服务镜像发送端口"
	case RtdbParamServerReceiverPort:
		return "网络服务镜像接收端口"
	case RtdbParamServerMode:
		return "网络服务启动模式"
	case RtdbParamServerConnectionCount:
		return "协议中限定网络服务连接并发数量"
	case RtdbParamArvPagesNumber:
		return "历史数据缓存中的页数量"
	case RtdbParamArvexPagesNumber:
		return "补历史数据缓存中的页数量"
	case RtdbParamExceptionAtServer:
		return "是否由服务器进行例外判定"
	case RtdbParamArvPageRecycleDelay:
		return "历史数据缓存页回收延时（毫秒）"
	case RtdbParamExArchiveSize:
		return "历史数据存档文件文件自动增长大小（单位：MB）"
	case RtdbParamArchiveBatchSize:
		return "历史存储值分段查询个数"
	case RtdbParamDatafilePagesize:
		return "系统数据文件页大小"
	case RtdbParamArvAsyncQueueNormalDoor:
		return "历史数据缓存队列中速归档区（单位：百分比）"
	case RtdbParamIndexAlwaysInMemory:
		return "常驻内存的历史数据索引大小（单位：MB）"
	case RtdbParamDiskMinRestSize:
		return "最低可用磁盘空间（单位：MB）"
	case RtdbParamMinSizeOfArchive:
		return "历史存档文件和附属文件的最小尺寸（单位：MB）"
	case RtdbParamDelayOfAutoMergeOrArrange:
		return "自动合并/整理最小延迟（单位：小时）"
	case RtdbParamStartOfAutoMergeOrArrange:
		return "自动合并/整理开始时间（单位：点钟）"
	case RtdbParamStopOfAutoMergeOrArrange:
		return "自动合并/整理停止时间（单位：点钟）"
	case RtdbParamStartOfAutoBackup:
		return "自动备份开始时间（单位：点钟）"
	case RtdbParamStopOfAutoBackup:
		return "自动备份停止时间（单位：点钟）"
	case RtdbParamMaxLatencyOfSnapshot:
		return "允许服务器时间之后多少小时内的数据进入快照（单位：小时）"
	case RtdbParamPageAllocatorReserveSize:
		return "活动页分配器预留大小（单位：KB）， 0 表示使用操作系统视图大小"
	case RtdbParamIncludeSnapshotInQuery:
		return "决定取样本值和统计值时，快照是否应该出现在查询结果中"
	case RtdbParamLicBlobCount:
		return "协议中限定的字符串或BLOB类型标签点数量"
	case RtdbParamMirrorBufferSize:
		return "镜像文件大小（单位：GB）"
	case RtdbParamBlobArvexPagesNumber:
		return "blob、str补历史的默认缓存页数量"
	case RtdbParamMirrorEventQueueCapacity:
		return "镜像缓存队列容量"
	case RtdbParamNotifyNotEnoughSpace:
		return "提示磁盘空间不足，一旦启用，设置为ON，则通过API返回大错误码，否则只记录日志"
	case RtdbParamArchiveFixedRange:
		return "历史数据存档文件的固定时间范围，默认为0表示不使用固定时间范围（单位：分钟）"
	case RtdbParamOneClinetMaxConnectionCount:
		return "单个客户端允许的最大连接数，默认为0表示不限制"
	case RtdbParamArvPagesCapacity:
		return "历史数据缓存所占字节大小，单位：字节"
	case RtdbParamArvexPagesCapacity:
		return "历史数据补写缓存所占字节大小，单位：字节"
	case RtdbParamBlobArvexPagesCapacity:
		return "blob、string类型标签点历史数据补写缓存所占字节大小，单位：字节"
	case RtdbParamLockedPagesMem:
		return "指定分配给数据库用的内存大小，单位：MB"
	case RtdbParamLicRecycleCount:
		return "协议中回收站的容量"
	case RtdbParamArchivedPolicy:
		return "快照数据和补写数据的归档策略"
	case RtdbParamNetworkIsolationAckByte:
		return "网络隔离装置ACK字节"
	case RtdbParamEnableLogger:
		return "启用日志输出，0为不启用"
	case RtdbParamLogEncode:
		return "启用日志加密，0为不启用"
	case RtdbParamLoginTry:
		return "启用登录失败次数验证，0为不启用"
	case RtdbParamUserLog:
		return "启用用户详细日志，0为不启用"
	case RtdbParamCoverWriteLog:
		return "启用日志覆盖写功能，0为不启用"
	case RtdbParamLicNamedTypeCount:
		return "协议中限定的自定义类型标签点数量"
	case RtdbParamMirrorReceiverThreadpoolSize:
		return "镜像接收线程数量"
	case RtdbParamSnapshotUseArchiveInterface:
		return "按照补历史流程归档快照数据页"
	case RtdbParamNoArcdataWriteLog:
		return "归档无对应存档文件的数据时记录日志"
	case RtdbParamPutArchiveThreadNum:
		return "补历史归档线程数"
	case RtdbParamArvexDataArchivedThreshold:
		return "单次补写数据归档阈值"
	case RtdbParamSnapshotFlushBufferDelay:
		return "快照服务的共享缓存刷新到磁盘的周期"
	case RtdbParamDataSpeed:
		return "查询时使用加速统计"
	case RtdbParamUseNewPlotAlgo:
		return "启用新的曲线算法"
	case RtdbParamQueryThreadPoolSize:
		return "曲线查询线程池中线程数量"
	case RtdbParamArchivedValues:
		return "使用查询线程池查询历史数据"
	case RtdbParamArchivedValuesCount:
		return "使用查询线程池查询历史数据的条数"
	case RtdbParamPoolUseFlag:
		return "启用曲线池"
	case RtdbParamPoolOutLogFlag:
		return "输出曲线池日志"
	case RtdbParamPoolTimeUsePoolFlag:
		return "使用曲线池缓存计算插值"
	case RtdbParamPoolMaxPointCount:
		return "曲线池的标签点容量"
	case RtdbParamPoolOneFileSavePointCount:
		return "曲线池每个数据文件的标签点容量"
	case RtdbParamPoolSaveMemorySize:
		return "曲线缓存退出时临时缓冲区大小"
	case RtdbParamPoolMinTimeUnitSeconds:
		return "曲线池缓存数据当前时间单位"
	case RtdbParamPoolTimeUnitViewRate:
		return "曲线池查询数据最小时间单位显示系数"
	case RtdbParamPoolTimerIntervalSeconds:
		return "曲线池定时器刷新周期"
	case RtdbParamPoolPerfTimerIntervalSeconds:
		return "曲线池性能计算点刷新周期"
	case RtdbParamArchiveInitFileSize:
		return "存档文件初始大小"
	case RtdbParamArchiveIncreaseMode:
		return "存档文件增长模式"
	case RtdbParamArchiveIncreaseSize:
		return "固定模式下文件增长大小"
	case RtdbParamArchiveIncreasePercent:
		return "百分比模式下增长百分比"
	case RtdbParamAllowConvertSklToRbtIndex:
		return "跳跃链表转换到红黑树"
	case RtdbParamEarlyDataTime:
		return "冷数据时间"
	case RtdbParamEarlyIndexTime:
		return "自动转换索引时间"
	case RtdbParamArrangeRbtTime:
		return "整理存档文件时决定索引格式的时间轴"
	case RtdbParamEnableBigData:
		return "将存档文件全部读取到内存中"
	case RtdbParamAutoArrangePercent:
		return "自动整理存档文件时的实际使用率"
	case RtdbParamEarlyArrangeTime:
		return "自动整理存档文件的时间"
	case RtdbParamMinAutoArrangeArcfilePercent:
		return "自动整理存档文件时的最小使用率"
	case RtdbParamArrangeArcWithMemory:
		return "在内存中整理存档文件"
	case RtdbParamAraangeArcMaxMemPercent:
		return "整理存档文件最大内存使用率"
	case RtdbParamMaxDiskSpacePercent:
		return "磁盘最大使用率"
	case RtdbParamUseDispath:
		return "windows 用 linux 已禁用,是否启用转发服务"
	case RtdbParamUseSmartParam:
		return "windows 用 linux 已禁用,是否使用推荐参数"
	case RtdbParamSubscribeSnapshotCount:
		return "单连接快照事件订阅个数"
	case RtdbParamSubscribeQueueSize:
		return "订阅事件队列大小"
	case RtdbParamSubscribeTimeout:
		return "订阅事件超时时间"
	case RtdbParamMirrorCompressOnoff:
		return "镜像报文压缩是否打开"
	case RtdbParamMirrorCompressType:
		return "镜像报文压缩类型"
	case RtdbParamMirrorCompressMin:
		return "镜像报文压缩最小值"
	case RtdbParamArchiveRollTime:
		return "存档文件滚动时间轴"
	case RtdbParamHandleTimeOut:
		return "连接超时断开，单位：秒"
	case RtdbParamMoveArvTime:
		return "移动存档文件时决定移动存档的时间轴"
	case RtdbParamUseNewInterpAlgo:
		return "启用新的插值算法"
	case RtdbParamEnableReplication:
		return "启用双活，0为不启用，1为启用"
	case RtdbParamReplicationGroupPort:
		return "双活：同步组端口"
	case RtdbParamReplicationThreadSize:
		return "双活：同步线程数"
	case RtdbParamForceArchiveIncompleteDataPageDelay:
		return "强制归档补历史缓存里面未满数据页的延迟时间"
	case RtdbParamArchiveRollDiskPercentage:
		return "存档文件滚动存储空间百分比"
	case RtdbParamEnableIpv6:
		return "启用ipv6设置"
	case RtdbParamEnableUseArchivedValue:
		return "按条件获取历史值时，是否直接获取条件中点的历史值，0:获取插值，1:获取历史值"
	case RtdbParamTimestampType:
		return "获取服务器时间戳类型"
	case RtdbParamArcFilenameUsingDate:
		return "是否归档文件使用日期作为文件名"
	case RtdbParamLogMaxSpace:
		return "日志文件占用的最大磁盘空间"
	case RtdbParamLogFileSize:
		return "单个日志文件大小"
	case RtdbParamIgnoreToWriteNoarcbuffer:
		return "是否丢弃补历史数据"
	case RtdbParamArchivesCountForCalc:
		return "统计存档文件平均大小的存档文件个数"
	case RtdbParamMaxBlobSize:
		return "blob、str类型数据在数据库中允许的最大长度"
	default:
		return "未知系统参数"
	}
}

// ParamString 字符串类型系统参数
type ParamString string

// ParamInt 数值类型系统参数
type ParamInt uint32

// RtdbConst 一些常量
type RtdbConst int32

const (
	// RtdbConstTagSize 标签点名称占用字节数
	RtdbConstTagSize = RtdbConst(C.RTDB_TAG_SIZE)

	// RtdbConstDescSize 标签点描述占用字节数
	RtdbConstDescSize = RtdbConst(C.RTDB_DESC_SIZE)

	// RtdbConstUnitSize 标签点单位占用字节数
	RtdbConstUnitSize = RtdbConst(C.RTDB_UNIT_SIZE)

	// RtdbConstUserSize 用户名占用字节数
	RtdbConstUserSize = RtdbConst(C.RTDB_USER_SIZE)

	// RtdbConstSourceSize 标签点数据源占用字节数
	RtdbConstSourceSize = RtdbConst(C.RTDB_SOURCE_SIZE)

	// RtdbConstInstrumentSize 标签点所属设备占用字节数
	RtdbConstInstrumentSize = RtdbConst(C.RTDB_INSTRUMENT_SIZE)

	// RtdbConstLocationsSize  采集标签点位址个数
	RtdbConstLocationsSize = RtdbConst(C.RTDB_LOCATIONS_SIZE)

	// RtdbConstUserintSize  采集标签点用户自定义整数个数
	RtdbConstUserintSize = RtdbConst(C.RTDB_USERINT_SIZE)

	// RtdbConstUserrealSize  采集标签点用户自定义浮点数个数
	RtdbConstUserrealSize = RtdbConst(C.RTDB_USERREAL_SIZE)

	// RtdbConstEquationSize 计算标签点方程式占用字节数
	RtdbConstEquationSize = RtdbConst(C.RTDB_EQUATION_SIZE)

	// RtdbConstTypeNameSize 自定义类型名称占用字节数
	RtdbConstTypeNameSize = RtdbConst(C.RTDB_TYPE_NAME_SIZE)

	// RtdbConstPackOfSnapshot   事件快照备用字节空间
	RtdbConstPackOfSnapshot = RtdbConst(C.RTDB_PACK_OF_SNAPSHOT)

	// RtdbConstPackOfPoint 标签点备用字节空间
	RtdbConstPackOfPoint = RtdbConst(C.RTDB_PACK_OF_POINT)

	// RtdbConstPackOfBasePoint 基本标签点备用字节空间
	RtdbConstPackOfBasePoint = RtdbConst(C.RTDB_PACK_OF_BASE_POINT)

	// RtdbConstPackOfScan 采集标签点备用字节空间
	RtdbConstPackOfScan = RtdbConst(C.RTDB_PACK_OF_SCAN)

	// RtdbConstPackOfCalc 计算标签点备用字节空间
	RtdbConstPackOfCalc = RtdbConst(C.RTDB_PACK_OF_CALC)

	// RtdbConstFileNameSize 文件名字符串字节长度
	RtdbConstFileNameSize = RtdbConst(C.RTDB_FILE_NAME_SIZE)

	// RtdbConstPathSize  路径字符串字节长度
	RtdbConstPathSize = RtdbConst(C.RTDB_PATH_SIZE)

	// RtdbConstMaxUserCount 最大用户个数
	RtdbConstMaxUserCount = RtdbConst(C.RTDB_MAX_USER_COUNT)

	// RtdbConstMaxAuthCount 最大信任连接段个数
	RtdbConstMaxAuthCount = RtdbConst(C.RTDB_MAX_AUTH_COUNT)

	// RtdbConstMaxBlacklistLen 连接黑名单最大长度
	RtdbConstMaxBlacklistLen = RtdbConst(C.RTDB_MAX_BLACKLIST_LEN)

	// RtdbConstMaxSubscribeSnapshots  单个连接最大订阅快照数量
	RtdbConstMaxSubscribeSnapshots = RtdbConst(C.RTDB_MAX_SUBSCRIBE_SNAPSHOTS)

	// RtdbConstApiServerDescriptionLen  API_SERVER中，decription描述字段的长度
	RtdbConstApiServerDescriptionLen = RtdbConst(C.RTDB_API_SERVER_DESCRIPTION_LEN)

	// RtdbConstMinEquationSize 缩减长度后的方程式占用字节数
	RtdbConstMinEquationSize = RtdbConst(C.RTDB_MIN_EQUATION_SIZE)

	// RtdbConstPackOfMinCalc 缩减长度后的计算标签点备用字节空间
	RtdbConstPackOfMinCalc = RtdbConst(C.RTDB_PACK_OF_MIN_CALC)

	// RtdbConstMaxEquationSize 最大长度的方程式占用字节数
	RtdbConstMaxEquationSize = RtdbConst(C.RTDB_MAX_EQUATION_SIZE)

	// RtdbConstPackOfMaxCalc 最大长度的计算标签点备用字节空间
	RtdbConstPackOfMaxCalc = RtdbConst(C.RTDB_PACK_OF_MAX_CALC)

	// RtdbConstMaxJsonSize  允许的json字符串的最大长度
	RtdbConstMaxJsonSize = RtdbConst(C.RTDB_MAX_JSON_SIZE)

	// RtdbConstIpv6AddrSize ipv6地址空间大小
	RtdbConstIpv6AddrSize = RtdbConst(C.RTDB_IPV6_ADDR_SIZE)

	// RtdbConstOutputPluginLibVersionLength 输入输出适配器插件库版本号长度64Bytes
	RtdbConstOutputPluginLibVersionLength = RtdbConst(C.RTDB_OUTPUT_PLUGIN_LIB_VERSION_LENGTH)

	// RtdbConstOutputPluginLibNameLength  输入输出适配器插件库名长度128Bytes
	RtdbConstOutputPluginLibNameLength = RtdbConst(C.RTDB_OUTPUT_PLUGIN_LIB_NAME_LENGTH)

	// RtdbConstOutputPluginActorNameLength  输入输出适配器执行实例长度128Bytes
	RtdbConstOutputPluginActorNameLength = RtdbConst(C.RTDB_OUTPUT_PLUGIN_ACTOR_NAME_LENGTH)

	// RtdbConstOutputPluginNameLength 输入输出适配器插件名长度255Bytes
	RtdbConstOutputPluginNameLength = RtdbConst(C.RTDB_OUTPUT_PLUGIN_NAME_LENGTH)

	// RtdbConstOutputPluginDirLength 输入输出适配器路径长度255Bytes
	RtdbConstOutputPluginDirLength = RtdbConst(C.RTDB_OUTPUT_PLUGIN_DIR_LENGTH)

	// RtdbConstPerOfUsefulMemSize 历史数据缓存/补历史数据缓存/blob补历史数据缓存/内存大小上限占可用内存百分比值占用的字节数
	RtdbConstPerOfUsefulMemSize = RtdbConst(C.RTDB_PER_OF_USEFUL_MEM_SIZE)
)

// DateTimeType 32位时间戳类型，秒级时间戳
type DateTimeType int32

// TimestampType 64位时间戳类型，秒级时间戳
type TimestampType int64

// SubtimeType 时间戳，小于秒的部分，根据设置的全局时间戳精度，表示毫秒、微秒、纳秒的部分
type SubtimeType int32

// RtdbHostConnectInfo 连接到RTDB数据库服务器的连接信息
// 备注， IPv6版本兼容此 RtdbHostConnectInfo ， 因此暂时注释掉
// type RtdbHostConnectInfo struct {
// 	IpAddr      int32        // 连接的客户端IP地址
// 	Port        uint16       // 连接端口
// 	Job         int32        // 连接最近处理的任务
// 	JobTime     DateTimeType // 最近处理任务的时间
// 	ConnectTime DateTimeType // 客户端连接时间
// 	Client      string       // 连接的客户端主机名称
// 	Process     string       // 连接的客户端程序名
// 	User        string       // 登录的用户
// 	Length      int32        // 记录用户名长度，用于加密传输
// }
// func cToRtdbHostConnectInfo(cInfo *C.RTDB_HOST_CONNECT_INFO) RtdbHostConnectInfo {
// 	goInfo := RtdbHostConnectInfo{
// 		IpAddr:      int32(cInfo.ipaddr),
// 		Port:        uint16(cInfo.port),
// 		Job:         int32(cInfo.job),
// 		JobTime:     DateTimeType(cInfo.job_time),
// 		ConnectTime: DateTimeType(cInfo.connect_time),
// 		Client:      CCharArrayToString(&cInfo.client[0], len(cInfo.client)),
// 		Process:     CCharArrayToString(&cInfo.process[0], len(cInfo.process)),
// 		User:        CCharArrayToString(&cInfo.user[0], len(cInfo.user)),
// 		Length:      int32(cInfo.length),
// 	}
// 	return goInfo
// }

type RtdbHostConnectInfoIpv6 struct {
	IpAddr      int32        // 连接的客户端IP地址
	IpAddr6     string       // ipv6地址
	Port        uint16       // 连接端口
	Job         int32        // 连接最近处理的任务
	JobTime     DateTimeType // 最近处理任务的时间
	ConnectTime DateTimeType // 客户端连接时间
	Client      string       // 连接的客户端主机名称
	Process     string       // 连接的客户端程序名
	User        string       // 登录的用户
	Length      int32        // 记录用户名长度，用于加密传输
}

func cToRtdbHostConnectInfoIpv6(cInfo *C.RTDB_HOST_CONNECT_INFO_IPV6) RtdbHostConnectInfoIpv6 {
	goInfo := RtdbHostConnectInfoIpv6{
		IpAddr:      int32(cInfo.ipaddr),
		IpAddr6:     CCharArrayToString(&cInfo.ipaddr6[0], len(cInfo.ipaddr6)),
		Port:        uint16(cInfo.port),
		Job:         int32(cInfo.job),
		JobTime:     DateTimeType(cInfo.job_time),
		ConnectTime: DateTimeType(cInfo.connect_time),
		Client:      CCharArrayToString(&cInfo.client[0], len(cInfo.client)),
		Process:     CCharArrayToString(&cInfo.process[0], len(cInfo.process)),
		User:        CCharArrayToString(&cInfo.user[0], len(cInfo.user)),
		Length:      int32(cInfo.length),
	}
	return goInfo
}

// RtdbOsType 操作系统类型
type RtdbOsType int8

const (
	RtdbOsWindows = RtdbOsType(C.RTDB_OS_WINDOWS)
	RtdbOsLinux   = RtdbOsType(C.RTDB_OS_LINUX)
	RtdbOsInvalid = RtdbOsType(C.RTDB_OS_INVALID)
)

func (ost RtdbOsType) Desc() string {
	switch ost {
	case RtdbOsWindows:
		return "windows"
	case RtdbOsLinux:
		return "linux"
	case RtdbOsInvalid:
		return "未知操作系统"
	default:
		return "无效OsType"
	}
}

// RtdbHandleInfo 连接信息
type RtdbHandleInfo struct {
	OsType RtdbOsType // 当前连接数据库的系统，参考 RTDB_OS_TYPE
	NewDB  Switch     // 当前连接数据库的版本，0(OFF)表示旧版本，1(ON)表示新版本
}

func cToRtdbHandleInfo(cOsType *C.RTDB_HANDLE_INFO) RtdbHandleInfo {
	goHandleInfo := RtdbHandleInfo{
		OsType: RtdbOsType(cOsType.os_type),
		NewDB:  Switch(cOsType.new_db),
	}
	return goHandleInfo
}

// RtdbUserInfo 用户信息
type RtdbUserInfo struct {
	User      string    // 用户名
	Privilege PrivGroup // 权限
	IsLocked  Switch    // 是否锁定
}

func cToRtdbUserInfo(cInfo *C.RTDB_USER_INFO) RtdbUserInfo {
	goInfo := RtdbUserInfo{
		User:      CCharArrayToString(&cInfo.user[0], len(cInfo.user)),
		Privilege: PrivGroup(cInfo.privilege),
		IsLocked:  Switch(cInfo.islocked),
	}
	return goInfo
}

// BlackList 黑名单
type BlackList struct {
	Addr string
	Mask string
	Desc string
}

// AuthorizationsList 白名单
type AuthorizationsList struct {
	Addr string
	Mask string
	Desc string
	Priv PrivGroup
}

// DirItem 目录项，表示目录下面的子条目，可能是子目录也可能是子文件
type DirItem struct {
	Path  string        // 文件、子目录全路径
	IsDir Switch        // ON(1)为目录，OFF(0)为文件
	ATime TimestampType // 访问时间
	CTime TimestampType // 建立时间
	MTime TimestampType // 修改时间
	Size  int64         // 文件大小
}

const (
	// RtdbMaxPath 系统支持的最大路径长度
	RtdbMaxPath = int32(2048)

	// RtdbMaxHostnameSize 系统支持的最大主机名长度
	RtdbMaxHostnameSize = int32(1024)
)

// Quality 质量码
type Quality int16

// TableID 表ID
type TableID int32

// RtdbTable 表
type RtdbTable struct {
	ID   TableID // 表ID
	Type int32   // 表类型，暂时无用
	Name string  // 表名
	Desc string  // 表描述
}

func cToRtdbTable(table *C.RTDB_TABLE) RtdbTable {
	rtn := RtdbTable{
		ID:   TableID(table.id),
		Type: int32(table._type),
		Name: CCharArrayToString(&table.name[0], int(C.RTDB_TAG_SIZE)),
		Desc: CCharArrayToString(&table.desc[0], int(C.RTDB_DESC_SIZE)),
	}
	return rtn
}

// PointID 点ID
type PointID int32

// RtdbType 数据类型
type RtdbType int32

const (
	// RtdbTypeAny 任意类型, 备注，只有标签点搜索的时候可以使用, 如：RawRtdbbSearchExWarp
	RtdbTypeAny = RtdbType(-1)

	// RtdbTypeBool 布尔类型，0值或1值
	RtdbTypeBool = RtdbType(C.RTDB_BOOL)

	// RtdbTypeUint8 无符号8位整数，占用1字节
	RtdbTypeUint8 = RtdbType(C.RTDB_UINT8)

	// RtdbTypeInt8 有符号8位整数，占用1字节
	RtdbTypeInt8 = RtdbType(C.RTDB_INT8)

	// RtdbTypeChar 单字节字符，占用1字节
	RtdbTypeChar = RtdbType(C.RTDB_CHAR)

	// RtdbTypeUint16 无符号16位整数，占用2字节
	RtdbTypeUint16 = RtdbType(C.RTDB_UINT16)

	// RtdbTypeInt16 有符号16位整数，占用2字节
	RtdbTypeInt16 = RtdbType(C.RTDB_INT16)

	// RtdbTypeUint32 无符号32位整数，占用4字节
	RtdbTypeUint32 = RtdbType(C.RTDB_UINT32)

	// RtdbTypeInt32 有符号32位整数，占用4字节
	RtdbTypeInt32 = RtdbType(C.RTDB_INT32)

	// RtdbTypeInt64 有符号64位整数，占用8字节
	RtdbTypeInt64 = RtdbType(C.RTDB_INT64)

	// RtdbTypeReal16 16位浮点数，占用2字节
	RtdbTypeReal16 = RtdbType(C.RTDB_REAL16)

	// RtdbTypeReal32 32位单精度浮点数，占用4字节
	RtdbTypeReal32 = RtdbType(C.RTDB_REAL32)

	// RtdbTypeReal64 64位双精度浮点数，占用8字节
	RtdbTypeReal64 = RtdbType(C.RTDB_REAL64)

	// RtdbTypeCoor 二维坐标，具有x、y两个维度的浮点数，占用8字节
	RtdbTypeCoor = RtdbType(C.RTDB_COOR)

	// RtdbTypeString 字符串，长度不超过存储页面大小
	RtdbTypeString = RtdbType(C.RTDB_STRING)

	// RtdbTypeBlob 二进制数据块，占用字节不超过存储页面大小
	RtdbTypeBlob = RtdbType(C.RTDB_BLOB)

	// RtdbTypeNamedT 自定义类型，由用户创建时确定字节长度
	RtdbTypeNamedT = RtdbType(C.RTDB_NAMED_T)

	// RtdbTypeDatetime 时间格式类型
	RtdbTypeDatetime = RtdbType(C.RTDB_DATETIME)

	// RtdbTypeFp16 定点数，占用2字节
	RtdbTypeFp16 = RtdbType(C.RTDB_FP16)

	// RtdbTypeFp32 定点数，占用4字节
	RtdbTypeFp32 = RtdbType(C.RTDB_FP32)

	// RtdbTypeFp64 定点数，占用8字节
	RtdbTypeFp64 = RtdbType(C.RTDB_FP64)
)

// RtdbMirror 镜像选项
type RtdbMirror int8

const (
	// RtdbMirrorPointOff 镜像关闭
	RtdbMirrorPointOff = RtdbMirror(C.RTDB_POINT_OFF)

	// RtdbPointSendRecv 镜像收发
	RtdbPointSendRecv = RtdbMirror(C.RTDB_POINT_SEND_RECV)

	// RtdbPointRecv 镜像接收
	RtdbPointRecv = RtdbMirror(C.RTDB_POINT_RECV)

	// RtdbPointSend 镜像发送
	RtdbPointSend = RtdbMirror(C.RTDB_POINT_SEND)
)

// RtdbPrecision 时间戳精度
type RtdbPrecision int8

const (
	// RtdbPrecisionAny 所有时间戳精度，只有标签点搜索的时候可用，例如: RawRtdbbSearchExWarp
	RtdbPrecisionAny = RtdbPrecision(-1)

	// RtdbPrecisionSecond 秒
	RtdbPrecisionSecond = RtdbPrecision(0)

	// RtdbPrecisionMilli 毫秒
	RtdbPrecisionMilli = RtdbPrecision(1)

	// RtdbPrecisionMicro 微秒
	RtdbPrecisionMicro = RtdbPrecision(2)

	// RtdbPrecisionNano 纳秒
	RtdbPrecisionNano = RtdbPrecision(3)
)

// RtdbClass 点类型
type RtdbClass uint32

const (
	// RtdbClassBase 基本标签点，所有类别标签点均在基本标签点的属性集上扩展自己的属性集。
	RtdbClassBase = RtdbClass(C.RTDB_BASE)

	// RtdbClassScan 采集标签点。
	RtdbClassScan = RtdbClass(C.RTDB_SCAN)

	// RtdbClassCalc 计算标签点
	RtdbClassCalc = RtdbClass(C.RTDB_CALC)

	// RtdbClassAlarm 报警标签点
	RtdbClassAlarm = RtdbClass(C.RTDB_ALARM)
)

// RtdbTrigger 计算标签点触发机制
type RtdbTrigger uint8

const (
	// RtdbNullTrigger 无触发
	RtdbNullTrigger = RtdbTrigger(C.RTDB_NULL_TRIGGER)

	// RtdbEventTrigger 事件触发
	RtdbEventTrigger = RtdbTrigger(C.RTDB_EVENT_TRIGGER)

	// RtdbTimerTrigger 周期触发
	RtdbTimerTrigger = RtdbTrigger(C.RTDB_TIMER_TRIGGER)

	// RtdbFixtimeTrigger 定时触发
	RtdbFixtimeTrigger = RtdbTrigger(C.RTDB_FIXTIME_TRIGGER)
)

// RtdbTimeCopy 计算结果时间戳参考
type RtdbTimeCopy uint8

const (
	// RtdbCalcTime 采用计算时间
	RtdbCalcTime = RtdbTimeCopy(C.RTDB_CALC_TIME)

	// RtdbLatestTime 采用最晚标签点时间
	RtdbLatestTime = RtdbTimeCopy(C.RTDB_LATEST_TIME)

	// RtdbEarliestTime 采用最早标签点时间
	RtdbEarliestTime = RtdbTimeCopy(C.RTDB_EARLIEST_TIME)
)

// RtdbPoint 基本标签点属性集
type RtdbPoint struct {
	// 标签点名称。
	// 用于在表中唯一标识一个标签点；
	// 该属性允许被修改；
	// 命名规则：
	// 1、第一个字符必须是26个字母之一或数字0-9之一或者"_"、"%"；
	// 2、不允许使用控制字符，比如换行符或制表符；
	// 3、不允许使用以下字符（'*'、'''、'?'、';'、'('、')'、'{'、'}'、'['、']'、'|'、'\'、'`'、'''、'"'、'.'）；
	// 4、字节长度不要超出 \b RTDB_TAG_SIZE，如果超出，系统会自动将后面的字符截断。
	Tag string

	// 全库唯一标识。
	// 只读属性，创建标签点时系统会自动为标签点分配的唯一标识，即标签点ID，标签点ID一经创建后永不更改。
	ID PointID

	// 标签点的数值类型。
	// 只读属性，在创建标签点时指定。
	Type RtdbType

	// 标签点所属表ID。
	Table TableID

	// 有关标签点的描述性文字。
	// 字节长度不要超出 \b RTDB_DESC_SIZE，多余的部分会被截断。
	Desc string

	// 工程单位。
	// 字节长度不要超出 \b RTDB_UNIT_SIZE，多余的部分会被截断。
	Unit string

	// 是否存档。
	// 缺省值：ON，1；
	// ON或1表示存档，OFF或0表示不存档。
	Archive Switch

	// 数值位数。
	// 缺省值：-5；
	// 范围：>=-20、<=10；
	// 用来控制数值的显示格式；
	// 如果为0或正数，表示数值的小数位数，如果为负数，表示数值的有效位数。
	Digits int16

	// 停机状态字（Shutdown）
	// 缺省值：0；
	// 定义该点在停机状态下是否补写停机状态值。
	// 1 表示补写；0 表示不补写。
	Shutdown Switch

	// 量程下限。
	// 缺省值：0；
	// 单位：标签点工程单位。
	LowLimit float32

	// 量程上限。
	// 缺省值：100；
	// 单位：标签点工程单位。
	HighLimit float32

	// 是否阶跃。
	// 缺省值：OFF，0；
	// 该属性决定了中间值的计算是用阶梯还是连续的内插值替换；
	// 缺省情况下该属性为OFF，即中间值的计算是用内插值替换；
	// 如果被设置为ON，则中间值的数值同前一个有记录的数值相同。
	// 在历史数据检索中，本设置可能被外部输入的阶跃开关覆盖。
	Step Switch

	// 典型值。
	// 缺省值：50；
	// 大于等于量程下限，小于等于量程上限。
	Typical float32

	// 是否压缩。
	// 缺省值：ON，1；
	// 如果该属性被关闭（OFF，0），任何到达数据存储服务器Server的数据都会被提交到历史数据库；否则（ON，1），只有满足压缩条件的数据才会被提交到历史数据库。
	// 需要手工录入的标签点应该将该属性设置为OFF，0。
	Compress Switch

	// 压缩偏差。
	// 单位：标签点工程单位；
	// 缺省值：1；
	// 当有新的数值被提交到数据存储服务器Server，如果从上一个被提交到历史数据库的数值开始有数值超出了压缩偏差外，则上一个被提交到数据存储服务器Server的数值被视为关键数值；
	// 该属性与[压缩偏差百分比（the percent of compression deviation）]属性含义相同，该属性与量程（high
	// limit - low limit）的百分比即[压缩偏差百分比（the percent of compression
	// deviation）]属性的值；
	// 对该属性的修改将导致对[压缩偏差百分比（the percent of compression
	// deviation）]的修改，同样，对[压缩偏差百分比（the percent of compression
	// deviation）]的修改也将修改该属性，如果两个同时被修改，[压缩偏差百分比（the percent of compression
	// deviation）]具有更高的优先权。
	CompDev float32

	// 压缩偏差百分比。
	// 单位：百分比；
	// \see compdev。
	CompDevPercent float32

	// 最大压缩间隔。
	// 单位：秒；
	// 缺省值：28800；
	// 如果某个数值与上一个被提交到历史数据库的数值的时间间隔大于或等于最大压缩间隔，无论是否满足压缩条件，该数值都应该被视为关键数值从而被提交到历史数据库的数据队列；
	// 数据库中两个标签点间的时间间隔有可能超出该属性值，因为数据存储服务器Server可能长时间接收不到提交的数据，而且任何系统绝对不会自己创造数据。
	CompTimeMax int32

	// 最短压缩间隔。
	// 单位：秒；
	// 缺省值：0；
	// 如果某个数值与上一个被提交到历史数据库的数值的时间间隔小于最短压缩间隔，该数值会被忽略；
	// 该属性有降噪（suppress noise）的作用。
	CompTimeMin int32

	// 例外偏差。
	// 单位：标签点工程单位；
	// 缺省值：0.5；
	// 如果某个数值与上一个被提交到数据存储服务器Server的数值的偏差大于该例外偏差（以数值的工程单位为准），则该数值被视为例外数值，应该被提交到数据存储服务器Server；
	// 建议例外偏差应该小于等于压缩偏差的一半；
	// 该属性与[例外偏差百分比（The Percent of Exception Deviation）]属性含义相同，该属性与量程（high
	// limit - low limit）的百分比即[例外偏差百分比（The Percent of Exception
	// Deviation）]属性的值；
	// 对该属性的修改将导致对[例外偏差百分比（The Percent of Exception
	// Deviation）]的修改，同样，对[例外偏差百分比（The Percent of Exception
	// Deviation）]的修改也将修改该属性，如果两个同时被修改，[例外偏差百分比（The Percent of Exception
	// Deviation）]具有更高的优先权。
	ExcDev float32

	// 例外偏差百分比。
	// 单位：百分比；
	// \see excdev。
	ExcDevPercent float32

	// 最大例外间隔。
	// 单位：秒；
	// 缺省值：600；
	// 如果某个数值与上一个被提交到数据存储服务器Server的数值的时间间隔大于或等于最大例外间隔，无论是否满足例外条件，该数值都应该被视为例外数值从而被提交到数据存储服务器Server；
	// 数据库中两个标签点间的时间间隔有可能超出该属性值，因为接口可能长时间采集不到数据，而且任何系统绝对不会自己创造数据；
	// 如果要关闭例外过滤，设置该属性为0即可。
	ExcTimeMax int32

	// 最短例外间隔。
	// 单位：秒；
	// 缺省值：0；
	// 如果某个数值与上一个被提交到数据存储服务器Server的数值的时间间隔小于最短例外间隔，无论是否满足例外条件，该数值都会被忽略；
	// 该属性有降噪（suppress noise）的作用。
	ExcTimeMin int32

	// 标签点类别。
	// RTDB_CLASS类型的组合，最多可以扩展至32种类型的组合；
	// 所有类别标签点均继承自"基本"类型标签点。
	// 不同类别的标签点具有不同的属性集，"采集"类别的标签点具有"设备标签"、"位号"、"自定义整数"和"自定义浮点数"等扩展属性，"计算"类别的标签点具有"扩展描述"、"触发机制"等扩展属性。
	Class RtdbClass

	// 标签点属性最后一次被修改的时间。
	ChangeDate DateTimeType

	// 标签点属性最后一次被修改的用户名。
	Changer string

	// 标签点被创建的时间。
	CreateDate DateTimeType

	// 标签点创建者的用户名。
	Creator string

	// 镜像收发控制。
	// 默认值：关闭，0
	// 开启镜像收发，1，既接收，又发送
	// 开启镜像接收，2，只接收，不发送
	// 开启镜像发送，3，只发送，不接收
	Mirror RtdbMirror

	// 时间戳精度。
	// 默认值：秒，0；
	// 用于设定标签点的历史值在存储中精确到"秒"（0）还是"毫秒/纳秒"（1）。
	// 标签点一经创建就不允许修改该属性。
	MilliSecond int8

	// 采集点扩展属性集存储地址索引。
	ScanIndex uint32

	// 计算点扩展属性集存储地址索引。
	CalcIndex uint32

	// 报警点扩展属性集存储地址索引。
	AlarmIndex uint32

	// 标签点全名，格式为“表名称.标签点名称”。
	TableDotTag string

	// 统计加速。
	// 默认值：关，0；
	// 用于设定是否生成标签点统计信息，从而加速历史数据统计过程。
	Summary Switch

	// 标签点对应自定义类型id，只用标签点类别为自定义类型时，才有意义。
	NamedTypeID uint16

	// 时间戳精度，0秒、1毫秒、2微秒、3纳秒
	Precision RtdbPrecision
}

// NewDefaultPoint 新建基本点结构，使用默认参数
func NewDefaultPoint(tag string, rtdbType RtdbType, tableID TableID, class RtdbClass, precision RtdbPrecision) *RtdbPoint {
	rtn := RtdbPoint{
		Tag:            tag,
		Type:           rtdbType,
		Table:          tableID,
		Desc:           "",
		Unit:           "",
		Archive:        ON,
		Digits:         -5,
		Shutdown:       OFF,
		LowLimit:       0,
		HighLimit:      100,
		Step:           OFF,
		Typical:        50,
		Compress:       ON,
		CompDev:        1,
		CompDevPercent: 0,
		CompTimeMax:    28800,
		CompTimeMin:    0,
		ExcDev:         0.5,
		ExcDevPercent:  0,
		ExcTimeMax:     600,
		ExcTimeMin:     0,
		Class:          class,
		Mirror:         RtdbMirrorPointOff,
		Summary:        OFF,
		NamedTypeID:    0,
		Precision:      precision,
	}
	switch rtn.Precision {
	case RtdbPrecisionSecond:
		rtn.MilliSecond = 0
	case RtdbPrecisionMilli, RtdbPrecisionMicro, RtdbPrecisionNano:
		rtn.MilliSecond = 1
	}
	return &rtn
}

// NewRtdbPoint 创建基本点结构
func NewRtdbPoint(
	tag string, rtdbType RtdbType, tableID TableID, desc, unit string, archive Switch, digits int16,
	shutdown Switch, lowLimit float32, highLimit float32, step Switch, typical float32, compress Switch,
	compDev float32, compDevPercent float32, compTimeMax int32, compTimeMin int32,
	excDev float32, excDevPercent float32, excTimeMax int32, excTimeMin int32, class RtdbClass,
	mirror RtdbMirror, summary Switch, precision RtdbPrecision,
) *RtdbPoint {
	rtn := RtdbPoint{
		Tag:            tag,
		Type:           rtdbType,
		Table:          tableID,
		Desc:           desc,
		Unit:           unit,
		Archive:        archive,
		Digits:         digits,
		Shutdown:       shutdown,
		LowLimit:       lowLimit,
		HighLimit:      highLimit,
		Step:           step,
		Typical:        typical,
		Compress:       compress,
		CompDev:        compDev,
		CompDevPercent: compDevPercent,
		CompTimeMax:    compTimeMax,
		CompTimeMin:    compTimeMin,
		ExcDev:         excDev,
		ExcDevPercent:  excDevPercent,
		ExcTimeMax:     excTimeMax,
		ExcTimeMin:     excTimeMin,
		Class:          class,
		Mirror:         mirror,
		Summary:        summary,
		Precision:      precision,
	}
	switch rtn.Precision {
	case RtdbPrecisionSecond:
		rtn.MilliSecond = 0
	case RtdbPrecisionMilli, RtdbPrecisionMicro, RtdbPrecisionNano:
		rtn.MilliSecond = 1
	}
	return &rtn
}

func goToCRtdbPoint(p *RtdbPoint) *C.RTDB_POINT {
	if p == nil {
		return nil
	}

	rtn := C.RTDB_POINT{}
	GoStringToCCharArray(p.Tag, &rtn.tag[0], int(C.RTDB_TAG_SIZE))
	rtn.id = C.int(p.ID)
	rtn._type = C.int(p.Type)
	rtn.table = C.int(p.Table)
	GoStringToCCharArray(p.Desc, &rtn.desc[0], int(C.RTDB_DESC_SIZE))
	GoStringToCCharArray(p.Unit, &rtn.unit[0], int(C.RTDB_UNIT_SIZE))
	rtn.archive = C.rtdb_byte(p.Archive)
	rtn.digits = C.short(p.Digits)
	rtn.shutdown = C.rtdb_byte(p.Shutdown)
	rtn.lowlimit = C.float(p.LowLimit)
	rtn.highlimit = C.float(p.HighLimit)
	rtn.step = C.rtdb_byte(p.Step)
	rtn.typical = C.float(p.Typical)
	rtn.compress = C.rtdb_byte(p.Compress)
	rtn.compdev = C.float(p.CompDev)
	rtn.compdevpercent = C.float(p.CompDevPercent)
	rtn.comptimemax = C.int(p.CompTimeMax)
	rtn.comptimemin = C.int(p.CompTimeMin)
	rtn.excdev = C.float(p.ExcDev)
	rtn.excdevpercent = C.float(p.ExcDevPercent)
	rtn.exctimemax = C.int(p.ExcTimeMax)
	rtn.exctimemin = C.int(p.ExcTimeMin)
	rtn.classof = C.uint(p.Class)
	rtn.changedate = C.rtdb_datetime_type(p.ChangeDate)
	GoStringToCCharArray(p.Changer, &rtn.changer[0], int(C.RTDB_USER_SIZE))
	rtn.createdate = C.rtdb_datetime_type(p.CreateDate)
	GoStringToCCharArray(p.Creator, &rtn.creator[0], int(C.RTDB_USER_SIZE))
	rtn.mirror = C.rtdb_byte(p.Mirror)
	rtn.millisecond = C.rtdb_byte(p.MilliSecond)
	rtn.scanindex = C.uint(p.ScanIndex)
	rtn.calcindex = C.uint(p.CalcIndex)
	rtn.alarmindex = C.uint(p.AlarmIndex)
	GoStringToCCharArray(p.TableDotTag, &rtn.table_dot_tag[0], int(C.RTDB_TAG_SIZE+C.RTDB_TAG_SIZE))
	rtn.summary = C.rtdb_byte(p.Summary)
	rtn.named_type_id = C.rtdb_uint16(p.NamedTypeID)
	rtn.precision = C.rtdb_precision_type(p.Precision)

	return &rtn
}

func cToRtdbPoint(p *C.RTDB_POINT) *RtdbPoint {
	if p == nil {
		return nil
	}

	rtn := RtdbPoint{
		Tag:            CCharArrayToString(&p.tag[0], int(C.RTDB_TAG_SIZE)),
		ID:             PointID(p.id),
		Type:           RtdbType(p._type),
		Table:          TableID(p.table),
		Desc:           CCharArrayToString(&p.desc[0], int(C.RTDB_DESC_SIZE)),
		Unit:           CCharArrayToString(&p.unit[0], int(C.RTDB_UNIT_SIZE)),
		Archive:        Switch(p.archive),
		Digits:         int16(p.digits),
		Shutdown:       Switch(p.shutdown),
		LowLimit:       float32(p.lowlimit),
		HighLimit:      float32(p.highlimit),
		Step:           Switch(p.step),
		Typical:        float32(p.typical),
		Compress:       Switch(p.compress),
		CompDev:        float32(p.compdev),
		CompDevPercent: float32(p.compdevpercent),
		CompTimeMax:    int32(p.comptimemax),
		CompTimeMin:    int32(p.comptimemin),
		ExcDev:         float32(p.excdev),
		ExcDevPercent:  float32(p.excdevpercent),
		ExcTimeMax:     int32(p.exctimemax),
		ExcTimeMin:     int32(p.exctimemin),
		Class:          RtdbClass(p.classof),
		ChangeDate:     DateTimeType(p.changedate),
		Changer:        CCharArrayToString(&p.changer[0], int(C.RTDB_USER_SIZE)),
		CreateDate:     DateTimeType(p.createdate),
		Creator:        CCharArrayToString(&p.creator[0], int(C.RTDB_USER_SIZE)),
		Mirror:         RtdbMirror(p.mirror),
		MilliSecond:    int8(p.millisecond),
		ScanIndex:      uint32(p.scanindex),
		CalcIndex:      uint32(p.calcindex),
		AlarmIndex:     uint32(p.alarmindex),
		TableDotTag:    CCharArrayToString(&p.table_dot_tag[0], int(C.RTDB_TAG_SIZE+C.RTDB_TAG_SIZE)),
		Summary:        Switch(p.summary),
		NamedTypeID:    uint16(p.named_type_id),
		Precision:      RtdbPrecision(p.precision),
	}

	return &rtn
}

// RtdbScan 采集标签点扩展属性集
type RtdbScan struct {
	// 全库唯一标识。0表示无效
	ID PointID

	// 数据源。
	// 缺省值：空（NULL）；
	// 将标签点同某些接口或某些模块相关联；
	// 每个数据源字符串只允许由26个字母（大小写敏感）和数字（0-9）组成，字节长度不要超出 \b RTDB_SOURCE_SIZE，多余的部分会被截断。
	Source string

	// 是否采集。
	// 缺省值：ON，1；
	// 该属性可能会被某些接口用到，如果该属性被关闭（OFF，0），从接口传来的数据可能不会被报告到数据库。
	Scan Switch

	// 设备标签。
	// 缺省值：空（NULL）；
	// 字节长度不要超出 \b RTDB_INSTRUMENT_SIZE，多余的部分会被截断。
	Instrument string

	// 共包含五个设备位址，缺省值全部为0。
	Locations [RtdbConstLocationsSize]int32

	// 共包含两个自定义整数，缺省值全部为0。
	UserInts [RtdbConstUserintSize]int32

	// 共包含两个自定义单精度浮点数，缺省值全部为0。
	UserReals [RtdbConstUserrealSize]float32
}

// NewRtdbScan 新建采集点结构
func NewRtdbScan(source string, scan Switch, instrument string) *RtdbScan {
	rtn := RtdbScan{
		Source:     source,
		Scan:       scan,
		Instrument: instrument,
	}
	return &rtn
}

func cToRtdbScan(p *C.RTDB_SCAN_POINT) *RtdbScan {
	if p == nil {
		return nil
	}

	rtn := RtdbScan{
		ID:         PointID(p.id),
		Source:     CCharArrayToString(&p.source[0], int(C.RTDB_SOURCE_SIZE)),
		Scan:       Switch(p.scan),
		Instrument: CCharArrayToString(&p.instrument[0], int(C.RTDB_INSTRUMENT_SIZE)),
	}
	for i := 0; i < int(RtdbConstLocationsSize); i++ {
		rtn.Locations[i] = int32(p.locations[i])
	}
	for i := 0; i < int(RtdbConstUserintSize); i++ {
		rtn.UserInts[i] = int32(p.userints[i])
	}
	for i := 0; i < int(RtdbConstUserrealSize); i++ {
		rtn.UserReals[i] = float32(p.userreals[i])
	}

	return &rtn
}

func goToCRtdbScan(p *RtdbScan) *C.RTDB_SCAN_POINT {
	if p == nil {
		return nil
	}

	rtn := C.RTDB_SCAN_POINT{}
	rtn.id = C.int(p.ID)
	GoStringToCCharArray(p.Source, &rtn.source[0], int(C.RTDB_SOURCE_SIZE))
	rtn.scan = C.rtdb_byte(p.Scan)
	GoStringToCCharArray(p.Instrument, &rtn.instrument[0], int(C.RTDB_INSTRUMENT_SIZE))
	for i := 0; i < int(RtdbConstLocationsSize); i++ {
		rtn.locations[i] = C.int(p.Locations[i])
	}
	for i := 0; i < int(RtdbConstUserintSize); i++ {
		rtn.userints[i] = C.int(p.UserInts[i])
	}
	for i := 0; i < int(RtdbConstUserrealSize); i++ {
		rtn.userreals[i] = C.float(p.UserReals[i])
	}

	return &rtn
}

// RtdbCalc 最大长度计算点扩展属性集
type RtdbCalc struct {
	// 全库唯一标识。0表示无效。
	ID PointID

	// 实时方程式。
	// 缺省值：空（NULL）；
	// 字节长度不要超出 \b RTDB_MAX_EQUATION_SIZE，长度超长的方程式将被拒绝设置入库，返回一个错误，避免错误的方程式进入系统，引发不安全因素。
	Equation string

	// 计算触发机制。枚举值参见 \b RTDB_TRIGGER。
	// 仅对"计算"类别标签点起作用，用于设置实时方程式服务对单个计算点的计算触发采用"事件触发"还是"周期触发"，
	// 对于"周期触发"以"事件触发"作为其先决判断条件，如果"事件触发"不满足，则不进行"周期触发"。
	Trigger RtdbTrigger

	// 计算结果时间戳参考, 枚举值参见 \b RTDB_TIME_COPY
	// 0: 表示采用计算时间作为计算结果时间戳；
	// 1: 表示采用输入标签点中的最晚时间作为计算结果时间戳；
	// 2: 表示采用输入标签点中的最早时间作为计算结果时间戳。
	TimeCopy RtdbTimeCopy

	// 对于“周期触发”的计算点，设定其计算周期，单位：秒
	Period int32
}

func NewRtdbCalc(equation string, trigger RtdbTrigger, timeCopy RtdbTimeCopy, period int32) *RtdbCalc {
	rtn := RtdbCalc{
		Equation: equation,
		Trigger:  trigger,
		TimeCopy: timeCopy,
		Period:   period,
	}
	return &rtn
}

func cToRtdbCalc(p *C.RTDB_MAX_CALC_POINT) *RtdbCalc {
	if p == nil {
		return nil
	}

	rtn := RtdbCalc{
		ID:       PointID(p.id),
		Equation: CCharArrayToString(&p.equation[0], int(C.RTDB_MAX_EQUATION_SIZE)),
		Trigger:  RtdbTrigger(p.trigger),
		TimeCopy: RtdbTimeCopy(p.timecopy),
		Period:   int32(p.period),
	}

	return &rtn
}

func goToCRtdbCalc(p *RtdbCalc) *C.RTDB_MAX_CALC_POINT {
	if p == nil {
		return nil
	}

	rtn := C.RTDB_MAX_CALC_POINT{}
	rtn.id = C.int(p.ID)
	GoStringToCCharArray(p.Equation, &rtn.equation[0], int(C.RTDB_MAX_EQUATION_SIZE))
	rtn.trigger = C.rtdb_byte(p.Trigger)
	rtn.timecopy = C.rtdb_byte(p.TimeCopy)
	rtn.period = C.int(p.Period)

	return &rtn
}

type RtdbSortFlag uint32

const (
	// RtdbSortFlagDescend 降序
	RtdbSortFlagDescend = RtdbSortFlag(C.RTDB_SORT_FLAG_DESCEND)

	// RtdbSortFlagCaseSensitive 大小写敏感
	RtdbSortFlagCaseSensitive = RtdbSortFlag(C.RTDB_SORT_FLAG_CASE_SENSITIVE)

	// RtdbSortFlagRecycled 用于回收站标签点排序
	RtdbSortFlagRecycled = RtdbSortFlag(C.RTDB_SORT_FLAG_RECYCLED)
)

type RtdbSearch int32

const (
	RtdbSearchAny            = RtdbSearch(C.RTDB_SEARCH_NULL)           //!< 不使用任何标签点属性作为搜索条件
	RtdbSearchNull           = RtdbSearch(C.RTDB_SEARCH_NULL)           //!< 不使用任何标签点属性作为搜索条件
	RtdbSearchCompdev        = RtdbSearch(C.RTDB_SEARCH_COMPDEV)        //!< 使用压缩偏差作为搜索条件
	RtdbSearchCompmax        = RtdbSearch(C.RTDB_SEARCH_COMPMAX)        //!< 最大压缩间隔
	RtdbSearchCompmin        = RtdbSearch(C.RTDB_SEARCH_COMPMIN)        //!< 最小压缩间隔
	RtdbSearchExcdev         = RtdbSearch(C.RTDB_SEARCH_EXCDEV)         //!< 例外偏差
	RtdbSearchExcmax         = RtdbSearch(C.RTDB_SEARCH_EXCMAX)         //!< 最大例外间隔
	RtdbSearchExcmin         = RtdbSearch(C.RTDB_SEARCH_EXCMIN)         //!< 最小例外间隔
	RtdbSearchSummary        = RtdbSearch(C.RTDB_SEARCH_SUMMARY)        //!< 是否加速
	RtdbSearchMirror         = RtdbSearch(C.RTDB_SEARCH_MIRROR)         //!< 是否镜像
	RtdbSearchCompress       = RtdbSearch(C.RTDB_SEARCH_COMPRESS)       //!< 是否压缩
	RtdbSearchStep           = RtdbSearch(C.RTDB_SEARCH_STEP)           //!< 是否阶跃
	RtdbSearchShutdown       = RtdbSearch(C.RTDB_SEARCH_SHUTDOWN)       //!< 是否停机补写
	RtdbSearchArchive        = RtdbSearch(C.RTDB_SEARCH_ARCHIVE)        //!< 是否归档
	RtdbSearchChanger        = RtdbSearch(C.RTDB_SEARCH_CHANGER)        //!< 修改人
	RtdbSearchCreator        = RtdbSearch(C.RTDB_SEARCH_CREATOR)        //!< 创建人
	RtdbSearchLowlimit       = RtdbSearch(C.RTDB_SEARCH_LOWLIMIT)       //!< 量程下限
	RtdbSearchHighlimit      = RtdbSearch(C.RTDB_SEARCH_HIGHLIMIT)      //!< 量程上限
	RtdbSearchTypical        = RtdbSearch(C.RTDB_SEARCH_TYPICAL)        //!< 典型值
	RtdbSearchChangedate     = RtdbSearch(C.RTDB_SEARCH_CHANGEDATE)     //!< 修改时间
	RtdbSearchCreatedate     = RtdbSearch(C.RTDB_SEARCH_CREATEDATE)     //!< 创建时间
	RtdbSearchDigits         = RtdbSearch(C.RTDB_SEARCH_DIGITS)         //!< 数值位数
	RtdbSearchCompdevpercent = RtdbSearch(C.RTDB_SEARCH_COMPDEVPERCENT) //!< 压缩偏差百分比
	RtdbSearchExcdevpercent  = RtdbSearch(C.RTDB_SEARCH_EXCDEVPERCENT)  //!< 例外偏差百分比
	RtdbSearchScan_begin     = RtdbSearch(C.RTDB_SEARCH_SCAN_BEGIN)     //!< 辅助作用，不能作为搜索条件
	RtdbSearchScan           = RtdbSearch(C.RTDB_SEARCH_SCAN)           //!< 是否采集
	RtdbSearchLocation1      = RtdbSearch(C.RTDB_SEARCH_LOCATION1)      //!< 设备位址1
	RtdbSearchLocation2      = RtdbSearch(C.RTDB_SEARCH_LOCATION2)      //!< 设备位址2
	RtdbSearchLocation3      = RtdbSearch(C.RTDB_SEARCH_LOCATION3)      //!< 设备位址3
	RtdbSearchLocation4      = RtdbSearch(C.RTDB_SEARCH_LOCATION4)      //!< 设备位址4
	RtdbSearchLocation5      = RtdbSearch(C.RTDB_SEARCH_LOCATION5)      //!< 设备位址5
	RtdbSearchUserint1       = RtdbSearch(C.RTDB_SEARCH_USERINT1)       //!< 自定义整数1
	RtdbSearchUserint2       = RtdbSearch(C.RTDB_SEARCH_USERINT2)       //!< 自定义整数2
	RtdbSearchUserreal1      = RtdbSearch(C.RTDB_SEARCH_USERREAL1)      //!< 自定义单精度浮点数1
	RtdbSearchUserreal2      = RtdbSearch(C.RTDB_SEARCH_USERREAL2)      //!< 自定义单精度浮点数2
	RtdbSearchScan_end       = RtdbSearch(C.RTDB_SEARCH_SCAN_END)       //!< 辅助作用，不能作为搜索条件
	RtdbSearchCalc_begin     = RtdbSearch(C.RTDB_SEARCH_CALC_BEGIN)     //!< 辅助作用，不能作为搜索条件
	RtdbSearchEquation       = RtdbSearch(C.RTDB_SEARCH_EQUATION)       //!< 方程式
	RtdbSearchTrigger        = RtdbSearch(C.RTDB_SEARCH_TRIGGER)        //!< 计算触发机制
	RtdbSearchTimecopy       = RtdbSearch(C.RTDB_SEARCH_TIMECOPY)       //!< 计算结果时间戳参考
	RtdbSearchPeriod         = RtdbSearch(C.RTDB_SEARCH_PERIOD)         //!< 计算周期
	RtdbSearchCalc_end       = RtdbSearch(C.RTDB_SEARCH_CALC_END)       //!< 辅助作用，不能作为搜索条件
)

func (rs RtdbSearch) Desc() string {
	switch rs {
	case RtdbSearchNull:
		return "不使用任何标签点属性作为搜索条件"
	case RtdbSearchCompdev:
		return "使用压缩偏差作为搜索条件"
	case RtdbSearchCompmax:
		return "最大压缩间隔"
	case RtdbSearchCompmin:
		return "最小压缩间隔"
	case RtdbSearchExcdev:
		return "例外偏差"
	case RtdbSearchExcmax:
		return "最大例外间隔"
	case RtdbSearchExcmin:
		return "最小例外间隔"
	case RtdbSearchSummary:
		return "是否统计加速"
	case RtdbSearchMirror:
		return "是否镜像"
	case RtdbSearchCompress:
		return "是否压缩"
	case RtdbSearchStep:
		return "是否阶跃"
	case RtdbSearchShutdown:
		return "是否停机补写"
	case RtdbSearchArchive:
		return "是否归档"
	case RtdbSearchChanger:
		return "修改人"
	case RtdbSearchCreator:
		return "创建人"
	case RtdbSearchLowlimit:
		return "量程下限"
	case RtdbSearchHighlimit:
		return "量程上限"
	case RtdbSearchTypical:
		return "典型值"
	case RtdbSearchChangedate:
		return "修改时间"
	case RtdbSearchCreatedate:
		return "创建时间"
	case RtdbSearchDigits:
		return "数值位数"
	case RtdbSearchCompdevpercent:
		return "压缩偏差百分比"
	case RtdbSearchExcdevpercent:
		return "例外偏差百分比"
	case RtdbSearchScan_begin:
		return "辅助作用，不能作为搜索条件"
	case RtdbSearchScan:
		return "是否采集"
	case RtdbSearchLocation1:
		return "设备位址1"
	case RtdbSearchLocation2:
		return "设备位址2"
	case RtdbSearchLocation3:
		return "设备位址3"
	case RtdbSearchLocation4:
		return "设备位址4"
	case RtdbSearchLocation5:
		return "设备位址5"
	case RtdbSearchUserint1:
		return "自定义整数1"
	case RtdbSearchUserint2:
		return "自定义整数2"
	case RtdbSearchUserreal1:
		return "自定义单精度浮点数1"
	case RtdbSearchUserreal2:
		return "自定义单精度浮点数2"
	case RtdbSearchScan_end:
		return "辅助作用，不能作为搜索条件"
	case RtdbSearchCalc_begin:
		return "辅助作用，不能作为搜索条件"
	case RtdbSearchEquation:
		return "方程式"
	case RtdbSearchTrigger:
		return "计算触发机制"
	case RtdbSearchTimecopy:
		return "计算结果时间戳参考"
	case RtdbSearchPeriod:
		return "计算周期"
	case RtdbSearchCalc_end:
		return "辅助作用，不能作为搜索条件"
	default:
		return "未知搜索条件"
	}
}

type RtdbTagIndex uint32

const (
	// RtdbTagIndexTag tag
	RtdbTagIndexTag = RtdbTagIndex(C.RTDB_TAG_INDEX_TAG)

	// RtdbTagIndexId id
	RtdbTagIndexId = RtdbTagIndex(C.RTDB_TAG_INDEX_ID)

	// RtdbTagIndexType type
	RtdbTagIndexType = RtdbTagIndex(C.RTDB_TAG_INDEX_TYPE)

	// RtdbTagIndexTable table
	RtdbTagIndexTable = RtdbTagIndex(C.RTDB_TAG_INDEX_TABLE)

	// RtdbTagIndexDesc desc
	RtdbTagIndexDesc = RtdbTagIndex(C.RTDB_TAG_INDEX_DESC)

	// RtdbTagIndexUnit unit
	RtdbTagIndexUnit = RtdbTagIndex(C.RTDB_TAG_INDEX_UNIT)

	// RtdbTagIndexArchive archive
	RtdbTagIndexArchive = RtdbTagIndex(C.RTDB_TAG_INDEX_ARCHIVE)

	// RtdbTagIndexDigits digits
	RtdbTagIndexDigits = RtdbTagIndex(C.RTDB_TAG_INDEX_DIGITS)

	// RtdbTagIndexShutdown shutdown
	RtdbTagIndexShutdown = RtdbTagIndex(C.RTDB_TAG_INDEX_SHUTDOWN)

	// RtdbTagIndexLowlimit lowlimit
	RtdbTagIndexLowlimit = RtdbTagIndex(C.RTDB_TAG_INDEX_LOWLIMIT)

	// RtdbTagIndexHighlimit highlimit
	RtdbTagIndexHighlimit = RtdbTagIndex(C.RTDB_TAG_INDEX_HIGHLIMIT)

	// RtdbTagIndexStep step
	RtdbTagIndexStep = RtdbTagIndex(C.RTDB_TAG_INDEX_STEP)

	// RtdbTagIndexTypical typical
	RtdbTagIndexTypical = RtdbTagIndex(C.RTDB_TAG_INDEX_TYPICAL)

	// RtdbTagIndexCompress compress
	RtdbTagIndexCompress = RtdbTagIndex(C.RTDB_TAG_INDEX_COMPRESS)

	// RtdbTagIndexCompdev compdev
	RtdbTagIndexCompdev = RtdbTagIndex(C.RTDB_TAG_INDEX_COMPDEV)

	// RtdbTagIndexCompdevpercent compdevpercent
	RtdbTagIndexCompdevpercent = RtdbTagIndex(C.RTDB_TAG_INDEX_COMPDEVPERCENT)

	// RtdbTagIndexComptimemax comptimemax
	RtdbTagIndexComptimemax = RtdbTagIndex(C.RTDB_TAG_INDEX_COMPTIMEMAX)

	// RtdbTagIndexComptimemin comptimemin
	RtdbTagIndexComptimemin = RtdbTagIndex(C.RTDB_TAG_INDEX_COMPTIMEMIN)

	// RtdbTagIndexExcdev excdev
	RtdbTagIndexExcdev = RtdbTagIndex(C.RTDB_TAG_INDEX_EXCDEV)

	// RtdbTagIndexExcdevpercent excdevpercent
	RtdbTagIndexExcdevpercent = RtdbTagIndex(C.RTDB_TAG_INDEX_EXCDEVPERCENT)

	// RtdbTagIndexExctimemax exctimemax
	RtdbTagIndexExctimemax = RtdbTagIndex(C.RTDB_TAG_INDEX_EXCTIMEMAX)

	// RtdbTagIndexExctimemin exctimemin
	RtdbTagIndexExctimemin = RtdbTagIndex(C.RTDB_TAG_INDEX_EXCTIMEMIN)

	// RtdbTagIndexClassof classof
	RtdbTagIndexClassof = RtdbTagIndex(C.RTDB_TAG_INDEX_CLASSOF)

	// RtdbTagIndexChangedate changedate
	RtdbTagIndexChangedate = RtdbTagIndex(C.RTDB_TAG_INDEX_CHANGEDATE)

	// RtdbTagIndexChanger changer
	RtdbTagIndexChanger = RtdbTagIndex(C.RTDB_TAG_INDEX_CHANGER)

	// RtdbTagIndexCreatedate createdate
	RtdbTagIndexCreatedate = RtdbTagIndex(C.RTDB_TAG_INDEX_CREATEDATE)

	// RtdbTagIndexCreator creator
	RtdbTagIndexCreator = RtdbTagIndex(C.RTDB_TAG_INDEX_CREATOR)

	// RtdbTagIndexMirror mirror
	RtdbTagIndexMirror = RtdbTagIndex(C.RTDB_TAG_INDEX_MIRROR)

	// RtdbTagIndexMs ms
	RtdbTagIndexMs = RtdbTagIndex(C.RTDB_TAG_INDEX_MS)

	// RtdbTagIndexFullname fullname
	RtdbTagIndexFullname = RtdbTagIndex(C.RTDB_TAG_INDEX_FULLNAME)

	// RtdbTagIndexSummary summary
	RtdbTagIndexSummary = RtdbTagIndex(C.RTDB_TAG_INDEX_SUMMARY)

	// RtdbTagIndexDatetimeformat datetimeformat
	RtdbTagIndexDatetimeformat = RtdbTagIndex(C.RTDB_TAG_INDEX_DATETIMEFORMAT)

	// RtdbTagIndexSource source
	RtdbTagIndexSource = RtdbTagIndex(C.RTDB_TAG_INDEX_SOURCE)

	// RtdbTagIndexScan scan
	RtdbTagIndexScan = RtdbTagIndex(C.RTDB_TAG_INDEX_SCAN)

	// RtdbTagIndexInstrument instrument
	RtdbTagIndexInstrument = RtdbTagIndex(C.RTDB_TAG_INDEX_INSTRUMENT)

	// RtdbTagIndexLocation1 location1
	RtdbTagIndexLocation1 = RtdbTagIndex(C.RTDB_TAG_INDEX_LOCATION1)

	// RtdbTagIndexLocation2 location2
	RtdbTagIndexLocation2 = RtdbTagIndex(C.RTDB_TAG_INDEX_LOCATION2)

	// RtdbTagIndexLocation3 location3
	RtdbTagIndexLocation3 = RtdbTagIndex(C.RTDB_TAG_INDEX_LOCATION3)

	// RtdbTagIndexLocation4 location4
	RtdbTagIndexLocation4 = RtdbTagIndex(C.RTDB_TAG_INDEX_LOCATION4)

	// RtdbTagIndexLocation5 location5
	RtdbTagIndexLocation5 = RtdbTagIndex(C.RTDB_TAG_INDEX_LOCATION5)

	// RtdbTagIndexUserint1 userint1
	RtdbTagIndexUserint1 = RtdbTagIndex(C.RTDB_TAG_INDEX_USERINT1)

	// RtdbTagIndexUserint2 userint2
	RtdbTagIndexUserint2 = RtdbTagIndex(C.RTDB_TAG_INDEX_USERINT2)

	// RtdbTagIndexUserreal1 userreal1
	RtdbTagIndexUserreal1 = RtdbTagIndex(C.RTDB_TAG_INDEX_USERREAL1)

	// RtdbTagIndexUserreal2 userreal2
	RtdbTagIndexUserreal2 = RtdbTagIndex(C.RTDB_TAG_INDEX_USERREAL2)

	// RtdbTagIndexEquation equation
	RtdbTagIndexEquation = RtdbTagIndex(C.RTDB_TAG_INDEX_EQUATION)

	// RtdbTagIndexTrigger trigger
	RtdbTagIndexTrigger = RtdbTagIndex(C.RTDB_TAG_INDEX_TRIGGER)

	// RtdbTagIndexTimecopy timecopy
	RtdbTagIndexTimecopy = RtdbTagIndex(C.RTDB_TAG_INDEX_TIMECOPY)

	// RtdbTagIndexPeriod period
	RtdbTagIndexPeriod = RtdbTagIndex(C.RTDB_TAG_INDEX_PERIOD)

	// RtdbTagIndexTimestamp timestamp
	RtdbTagIndexTimestamp = RtdbTagIndex(C.RTDB_TAG_INDEX_TIMESTAMP)

	// RtdbTagIndexValue value
	RtdbTagIndexValue = RtdbTagIndex(C.RTDB_TAG_INDEX_VALUE)

	// RtdbTagIndexQuality quality
	RtdbTagIndexQuality = RtdbTagIndex(C.RTDB_TAG_INDEX_QUALITY)
)

// RtdbSubscribeOption 订阅选项
type RtdbSubscribeOption int32

const (
	// RtdbSubscribeOptionAutoConn 自动重连
	RtdbSubscribeOptionAutoConn = RtdbSubscribeOption(C.RTDB_O_AUTOCONN)
)

// RtdbTagChangeReason 标签点变更原因，用于标签点订阅
type RtdbTagChangeReason int32

const (
	RtdbTagCreated  = RtdbTagChangeReason(C.RTDB_TAG_CREATED)
	RtdbTagUpdated  = RtdbTagChangeReason(C.RTDB_TAG_UPDATED)
	RtdbTagRemoved  = RtdbTagChangeReason(C.RTDB_TAG_REMOVED)
	RtdbTagRecoverd = RtdbTagChangeReason(C.RTDB_TAG_RECOVERD)
	RtdbTagPurged   = RtdbTagChangeReason(C.RTDB_TAG_PURGED)
	RtdbTabUpdated  = RtdbTagChangeReason(C.RTDB_TAB_UPDATED)
	RtdbTabRemoved  = RtdbTagChangeReason(C.RTDB_TAB_REMOVED)
)

// RtdbDataTypeField 自定义类型字段项
type RtdbDataTypeField struct {
	// 自定义类型的字段的名称，不要大于\b RTDB_TYPE_NAME_SIZE个字节
	Name string

	// 字段的类型,只支持 \b RTDB_TYPE 里的类型，不支持struct，union等组合类型。
	Type RtdbType

	// 字段类型的长度, \b RTDB_STRING、\b RTDB_BLOB等类型的具体的长度，基本类型本身的长度(基本类型可以忽略)，单位：字节。
	Length int32

	// 字段类型的描述，不要大于 \b RTDB_DESC_SIZE个字节
	Desc string
}

func goToCRtdbDataTypeField(field *RtdbDataTypeField) *C.RTDB_DATA_TYPE_FIELD {
	rtn := C.RTDB_DATA_TYPE_FIELD{}
	GoStringToCCharArray(field.Name, &rtn.name[0], int(C.RTDB_TYPE_NAME_SIZE))
	rtn._type = C.rtdb_int32(field.Type)
	rtn.length = C.rtdb_int32(field.Length)
	GoStringToCCharArray(field.Desc, &rtn.desc[0], int(C.RTDB_DESC_SIZE))
	return &rtn
}

func cToRtdbDataTypeField(field *C.RTDB_DATA_TYPE_FIELD) *RtdbDataTypeField {
	rtn := RtdbDataTypeField{
		Name:   CCharArrayToString(&field.name[0], int(C.RTDB_TYPE_NAME_SIZE)),
		Type:   RtdbType(field._type),
		Length: int32(field.length),
		Desc:   CCharArrayToString(&field.desc[0], int(C.RTDB_DESC_SIZE)),
	}
	return &rtn
}

// RtdbSyncRole 元数据同步角色
type RtdbSyncRole int8

const (
	// RtdbSyncRoleOffline 离线
	RtdbSyncRoleOffline = RtdbSyncRole(C.RTDB_SYNC_ROLE_OFFLINE)

	// RtdbSyncRoleUnsynced 未同步
	RtdbSyncRoleUnsynced = RtdbSyncRole(C.RTDB_SYNC_ROLE_UNSYNCED)

	// RtdbSyncRoleSyncing 同步中
	RtdbSyncRoleSyncing = RtdbSyncRole(C.RTDB_SYNC_ROLE_SYNCING)

	// RtdbSyncRoleSlave 备库
	RtdbSyncRoleSlave = RtdbSyncRole(C.RTDB_SYNC_ROLE_SLAVE)

	// RtdbSyncRoleMaster 主库
	RtdbSyncRoleMaster = RtdbSyncRole(C.RTDB_SYNC_ROLE_MASTER)
)

func (rs RtdbSyncRole) Desc() string {
	switch rs {
	case RtdbSyncRoleOffline:
		return "离线"
	case RtdbSyncRoleUnsynced:
		return "未同步"
	case RtdbSyncRoleSyncing:
		return "同步中"
	case RtdbSyncRoleSlave:
		return "备库"
	case RtdbSyncRoleMaster:
		return "主库"
	default:
		return "未知角色"
	}
}

// RtdbSyncStatus 元数据同步状态
type RtdbSyncStatus int8

const (
	// RtdbSyncStatusInit 正常
	RtdbSyncStatusInit = RtdbSyncStatus(C.RTDB_SYNC_STATUS_INIT)

	// RtdbSyncStatusStart 启动同步
	RtdbSyncStatusStart = RtdbSyncStatus(C.RTDB_SYNC_STATUS_START)

	// RtdbSyncStatusFile 同步文件
	RtdbSyncStatusFile = RtdbSyncStatus(C.RTDB_SYNC_STATUS_FILE)

	// RtdbSyncStatusCache 同步缓存
	RtdbSyncStatusCache = RtdbSyncStatus(C.RTDB_SYNC_STATUS_CACHE)
)

func (ss RtdbSyncStatus) Desc() string {
	switch ss {
	case RtdbSyncStatusInit:
		return "正常"
	case RtdbSyncStatusStart:
		return "启动同步"
	case RtdbSyncStatusFile:
		return "同步文件"
	case RtdbSyncStatusCache:
		return "同步缓存"
	default:
		return "未知同步状态"
	}
}

// RtdbSyncInfo 节点的元数据同步信息
type RtdbSyncInfo struct {
	Role     RtdbSyncRole
	Status   RtdbSyncStatus
	IP       uint32
	Version  uint64
	DataSize uint64
}

/////////////////////////////// 上面是结构定义 ////////////////////////////////////
/////////////////////////////// -- 躺平的分割线 -- ////////////////////////////////
/////////////////////////////// -- 躺平的分割线 -- ////////////////////////////////
/////////////////////////////// 下面是函数实现 ////////////////////////////////////

// RawRtdbGetApiVersionWarp 返回 ApiVersion 版本号
//
// output:
//   - ApiVersion 指的是 API库 的版本号
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_api_version_warp(rtdb_int32 *major, rtdb_int32 *minor, rtdb_int32 *beta)
func RawRtdbGetApiVersionWarp() (ApiVersion, error) {
	major, minor, beta := C.rtdb_int32(0), C.rtdb_int32(0), C.rtdb_int32(0)
	err := C.rtdb_get_api_version_warp(&major, &minor, &beta)
	version := ApiVersion{
		Major: int32(major),
		Minor: int32(minor),
		Beta:  int32(beta),
	}
	return version, RtdbError(err).GoError()
}

// RawRtdbSetOptionWarp 配置 API库 的行为参数，详见 RtdbApiOption 枚举
//
// input:
//   - optionType API库 的行为参数枚举
//   - value 每个 API库 行为参数枚举， 都可以附带一个 value 值对该行为参数进行调整
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_set_option_warp(rtdb_int32 type, rtdb_int32 value)
func RawRtdbSetOptionWarp(optionType RtdbApiOption, value int32) error {
	err := C.rtdb_set_option_warp(C.rtdb_int32(optionType), C.rtdb_int32(value))
	return RtdbError(err).GoError()
}

// RawRtdbCreateDatagramHandleWarp 创建数据流
//
// input:
//   - port 端口号
//   - remoteHost 对端IP地址
//
// output:
//   - DatagramHandle 数据流句柄
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_create_datagram_handle_warp(rtdb_int32 port, const char* remotehost, rtdb_datagram_handle* handle)
func RawRtdbCreateDatagramHandleWarp(port int32, remoteHost string) (DatagramHandle, error) {
	var handle C.rtdb_datagram_handle
	cRemoteHost := C.CString(remoteHost)
	defer C.free(unsafe.Pointer(cRemoteHost))
	err := C.rtdb_create_datagram_handle_warp(C.rtdb_int32(port), cRemoteHost, &handle)
	return DatagramHandle{handle: handle}, RtdbError(err).GoError()
}

// RawRtdbRemoveDatagramHandleWarp 删除数据流
//
// input:
//   - handle 数据流句柄
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_remove_datagram_handle_warp(rtdb_datagram_handle handle)
func RawRtdbRemoveDatagramHandleWarp(handle DatagramHandle) error {
	err := C.rtdb_remove_datagram_handle_warp(handle.handle)
	return RtdbError(err).GoError()
}

// RawRtdbRecvDatagramWarp 接收数据流
//
// input:
//   - handle  数据流句柄
//   - cacheLen 缓存大小，会创建对应大小的缓存，用于接收数据流返回的数据
//   - remoteAddr 对端IP地址
//   - timeout 超时时间(单位秒)
//
// output:
//   - []byte 返回数据流
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_recv_datagram_warp(char* message, rtdb_int32* message_len, rtdb_datagram_handle handle, char* remote_addr, rtdb_int32 timeout)
func RawRtdbRecvDatagramWarp(handle DatagramHandle, cacheLen int32, remoteAddr string, timeout int32) ([]byte, error) {
	message := make([]byte, cacheLen)
	messageLen := C.rtdb_int32(cacheLen)
	cRemoteAddr := C.CString(remoteAddr)
	defer C.free(unsafe.Pointer(cRemoteAddr))
	err := C.rtdb_recv_datagram_warp((*C.char)(unsafe.Pointer(&message[0])), &messageLen, handle.handle, cRemoteAddr, C.rtdb_int32(timeout))
	return message[0:messageLen], RtdbError(err).GoError()
}

// RawRtdbConnectWarp 建立同 RTDB 数据库的网络连接, 注意这里只是创建连接，并没有进行用户登陆
//
// input:
//   - hostname 数据库IP地址
//   - port 数据库端口号
//
// output:
//   - ConnectHandle 返回连接句柄
//
// raw_fn:
// - rtdb_error RTDBAPI_CALLRULE rtdb_connect_warp(const char *hostname, rtdb_int32 port, rtdb_int32 *handle)
func RawRtdbConnectWarp(hostname string, port int32) (ConnectHandle, error) {
	cHostname := C.CString(hostname)
	defer C.free(unsafe.Pointer(cHostname))
	cPort := C.rtdb_int32(port)
	cHandle := C.rtdb_int32(0)
	err := C.rtdb_connect_warp(cHostname, cPort, &cHandle)
	return ConnectHandle(cHandle), RtdbError(err).GoError()
}

// RawRtdbLoginWarp 以有效帐户登录
//
// input:
//   - handle 连接句柄
//   - user 用户名
//   - password 密码
//
// output:
//   - PrivGroup 用户权限
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_login_warp(rtdb_int32 handle, const char *user, const char *password, rtdb_int32 *priv)
func RawRtdbLoginWarp(handle ConnectHandle, user string, password string) (PrivGroup, error) {
	cUser := C.CString(user)
	defer C.free(unsafe.Pointer(cUser))
	cPassword := C.CString(password)
	defer C.free(unsafe.Pointer(cPassword))
	cPriv := C.rtdb_int32(0)
	err := C.rtdb_login_warp(C.rtdb_int32(handle), cUser, cPassword, &cPriv)
	return PrivGroup(cPriv), RtdbError(err).GoError()
}

// RawRtdbDisconnectWarp 断开同 RTDB 数据平台的连接
//
// input:
//   - handle 连接句柄
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_disconnect_warp(rtdb_int32 handle)
func RawRtdbDisconnectWarp(handle ConnectHandle) error {
	err := C.rtdb_disconnect_warp(C.rtdb_int32(handle))
	return RtdbError(err).GoError()
}

// RawRtdbConnectionCountWarp 获取 RTDB 服务器当前连接个数
//
// input:
//   - handle 连接句柄
//   - nodeNumber 单机模式下写0, 双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
//
// output:
//   - int32 返回连接个数
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_connection_count_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32 *count)
func RawRtdbConnectionCountWarp(handle ConnectHandle, nodeNumber int32) (int32, error) {
	count := C.rtdb_int32(0)
	err := C.rtdb_connection_count_warp(C.rtdb_int32(handle), C.rtdb_int32(nodeNumber), &count)
	return int32(count), RtdbError(err).GoError()
}

// RawRtdbGetDbInfo1Warp 获得字符串型数据库系统参数
//
// input:
//   - handle 连接句柄
//   - param 要取得的参数索引
//
// output:
//   - ParamString 参数索引对应的字符串
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_db_info1_warp(rtdb_int32 handle, rtdb_int32 index, char *str, rtdb_int32 size)
func RawRtdbGetDbInfo1Warp(handle ConnectHandle, param RtdbParam) (ParamString, error) {
	goStr := make([]byte, RtdbConstApiServerDescriptionLen)
	cStr := (*C.char)(unsafe.Pointer(&goStr[0]))
	err := C.rtdb_get_db_info1_warp(C.rtdb_int32(handle), C.rtdb_int32(param), cStr, C.rtdb_int32(RtdbConstApiServerDescriptionLen))
	rtn := C.GoString((*C.char)(unsafe.Pointer(&goStr[0])))
	return ParamString(rtn), RtdbError(err).GoError()
}

// RawRtdbGetDbInfo2Warp 获得整型数据库系统参数
//
// input:
//   - handle 连接句柄
//   - param 要取得的参数索引
//
// output:
//   - ParamInt 参数索引对应的数值
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_db_info2_warp(rtdb_int32 handle, rtdb_int32 index, rtdb_uint32 *value)
func RawRtdbGetDbInfo2Warp(handle ConnectHandle, param RtdbParam) (ParamInt, error) {
	value := C.rtdb_uint32(0)
	err := C.rtdb_get_db_info2_warp(C.rtdb_int32(handle), C.rtdb_int32(param), &value)
	return ParamInt(value), RtdbError(err).GoError()
}

// RawRtdbSetDbInfo1Warp 设置字符串型数据库系统参数
//
// input:
//   - handle 连接句柄
//   - param 要设置参数索引
//   - value 参数值
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_set_db_info1_warp(rtdb_int32 handle, rtdb_int32 index, const char *str)
func RawRtdbSetDbInfo1Warp(handle ConnectHandle, param RtdbParam, value ParamString) error {
	cValue := C.CString(string(value))
	defer C.free(unsafe.Pointer(cValue))
	err := C.rtdb_set_db_info1_warp(C.rtdb_int32(handle), C.rtdb_int32(param), cValue)
	return RtdbError(err).GoError()
}

// RawRtdbSetDbInfo2Warp 设置整型数据库系统参数
//
// input:
//   - handle 连接句柄
//   - index 要取得的参数索引
//   - value 参数数值
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_set_db_info2_warp(rtdb_int32 handle, rtdb_int32 index, rtdb_uint32 value)
func RawRtdbSetDbInfo2Warp(handle ConnectHandle, param RtdbParam, value ParamInt) error {
	err := C.rtdb_set_db_info2_warp(C.rtdb_int32(handle), C.rtdb_int32(param), C.rtdb_uint32(value))
	return RtdbError(err).GoError()
}

// RawRtdbGetConnectionsWarp 列出 RTDB 服务器的所有socket连接句柄, 注意这里指的是socket连接，区分于ConnectHandle
//
// input:
//   - handle 连接句柄
//   - nodeNumber 双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP, 单机模式下写0
//
// output:
//   - []SocketHandle socket连接数组
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_connections_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32 *sockets, rtdb_int32 *count)
func RawRtdbGetConnectionsWarp(handle ConnectHandle, nodeNumber int32) ([]SocketHandle, error) {
	connectionCount, err := RawRtdbGetDbInfo2Warp(handle, RtdbParamServerConnectionCount)
	if err != nil {
		return nil, err
	}
	cCount := C.rtdb_int32(connectionCount)
	sockets := make([]SocketHandle, int32(cCount))
	cSockets := (*C.rtdb_int32)(unsafe.Pointer(&sockets[0]))
	err2 := C.rtdb_get_connections_warp(C.rtdb_int32(handle), C.rtdb_int32(nodeNumber), cSockets, &cCount)
	return sockets[0:cCount], RtdbError(err2).GoError()
}

// RawRtdbGetOwnConnectionWarp 获取当前连接的socket句柄
//
// input:
//   - handle 连接句柄
//   - nodeNumber 双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP, 单机模式下写0
//
// output:
//   - SocketHandle socket连接
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_own_connection_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32* socket)
func RawRtdbGetOwnConnectionWarp(handle ConnectHandle, nodeNumber int32) (SocketHandle, error) {
	socket := C.rtdb_int32(0)
	err := C.rtdb_get_own_connection_warp(C.rtdb_int32(handle), C.rtdb_int32(nodeNumber), &socket)
	return SocketHandle(socket), RtdbError(err).GoError()
}

// RawRtdbGetConnectionInfoWarp 获取 RTDB 服务器指定连接的信息
// 备注： ipv6版本兼容此API，因此暂时注释掉此API
//
// input:
//   - handle 连接句柄
//   - nodeNumber 双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
//   - socket socket连接句柄
//
// output:
//   - RtdbHostConnectInfo 连接信息
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_connection_info_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32 socket, RTDB_HOST_CONNECT_INFO *info)
// func RawRtdbGetConnectionInfoWarp(handle ConnectHandle, nodeNumber int32, socket SocketHandle) (RtdbHostConnectInfo, error) {
// 	cInfo := C.RTDB_HOST_CONNECT_INFO{}
// 	err := C.rtdb_get_connection_info_warp(C.rtdb_int32(handle), C.rtdb_int32(nodeNumber), C.rtdb_int32(socket), &cInfo)
// 	goInfo := cToRtdbHostConnectInfo(&cInfo)
// 	return goInfo, RtdbError(err).GoError()
// }

// RawRtdbGetConnectionInfoIpv6Warp 获取 RTDB 服务器指定连接的ipv6版本
//
// input:
//   - handle 连接句柄
//   - nodeNumber 双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
//   - socket socket连接句柄
//
// output:
//   - RtdbHostConnectInfoIpv6 连接信息
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_connection_info_ipv6_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32 socket, RTDB_HOST_CONNECT_INFO_IPV6* info)
func RawRtdbGetConnectionInfoIpv6Warp(handle ConnectHandle, nodeNumber int32, socket SocketHandle) (RtdbHostConnectInfoIpv6, error) {
	cInfo := C.RTDB_HOST_CONNECT_INFO_IPV6{}
	err := C.rtdb_get_connection_info_ipv6_warp(C.rtdb_int32(handle), C.rtdb_int32(nodeNumber), C.rtdb_int32(socket), &cInfo)
	goInfo := cToRtdbHostConnectInfoIpv6(&cInfo)
	return goInfo, RtdbError(err).GoError()
}

// RawRtdbOsType 获取连接句柄所连接的服务器操作系统类型
//
// input:
//   - handle 连接句柄
//
// output:
//   - RtdbOsType 操作系统类型
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_linked_ostype_warp(rtdb_int32 handle, RTDB_OS_TYPE* ostype)
func RawRtdbOsType(handle ConnectHandle) (RtdbOsType, error) {
	osType := C.RTDB_OS_TYPE(C.RTDB_OS_INVALID)
	err := C.rtdb_get_linked_ostype_warp(C.rtdb_int32(handle), &osType)
	return RtdbOsType(osType), RtdbError(err).GoError()
}

// RawRtdbChangePasswordWarp 修改用户帐户口令
//
// input:
//   - handle 连接句柄
//   - user 已有帐户
//   - password 帐户新口令
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_change_password_warp(rtdb_int32 handle, const char *user, const char *password)
func RawRtdbChangePasswordWarp(handle ConnectHandle, user string, password string) error {
	cUser := C.CString(user)
	defer C.free(unsafe.Pointer(cUser))
	cPassword := C.CString(password)
	defer C.free(unsafe.Pointer(cPassword))
	err := C.rtdb_change_password_warp(C.rtdb_int32(handle), cUser, cPassword)
	return RtdbError(err).GoError()
}

// RawRtdbChangeMyPasswordWarp 用户修改自己帐户口令
//
// input:
//   - handle 连接句柄
//   - oldPwd 帐户原口令
//   - newPwd 帐户新口令
//
// raw_fn
//   - rtdb_error RTDBAPI_CALLRULE rtdb_change_my_password_warp(rtdb_int32 handle, const char *old_pwd, const char *new_pwd)
func RawRtdbChangeMyPasswordWarp(handle ConnectHandle, oldPwd string, newPwd string) error {
	cOldPwd := C.CString(oldPwd)
	defer C.free(unsafe.Pointer(cOldPwd))
	cNewPwd := C.CString(newPwd)
	defer C.free(unsafe.Pointer(cNewPwd))
	err := C.rtdb_change_my_password_warp(C.rtdb_int32(handle), cOldPwd, cNewPwd)
	return RtdbError(err).GoError()
}

// RawRtdbGetPrivWarp 获取连接权限
//
// input:
//   - handle 连接句柄
//
// output:
//   - PrivGroup 用户权限
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_priv_warp(rtdb_int32 handle, rtdb_int32 *priv)
func RawRtdbGetPrivWarp(handle ConnectHandle) (PrivGroup, error) {
	priv := C.rtdb_int32(0)
	err := C.rtdb_get_priv_warp(C.rtdb_int32(handle), &priv)
	return PrivGroup(priv), RtdbError(err).GoError()
}

// RawRtdbChangePrivWarp 修改用户帐户权限, 只有管理员有修改权限
//
// input:
//   - handle 连接句柄
//   - user 已有帐户
//   - priv 帐户权限
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_change_priv_warp(rtdb_int32 handle, const char *user, rtdb_int32 priv)
func RawRtdbChangePrivWarp(handle ConnectHandle, user string, priv PrivGroup) error {
	cUser := C.CString(user)
	defer C.free(unsafe.Pointer(cUser))
	err := C.rtdb_change_priv_warp(C.rtdb_int32(handle), cUser, C.rtdb_int32(priv))
	return RtdbError(err).GoError()
}

// RawRtdbAddUserWarp 添加用户帐户
//
// input:
//   - handle 连接句柄
//   - user 帐户
//   - password 帐户初始口令
//   - priv 帐户权限
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_add_user_warp(rtdb_int32 handle, const char *user, const char *password, rtdb_int32 priv)
func RawRtdbAddUserWarp(handle ConnectHandle, user string, password string, priv PrivGroup) error {
	cUser := C.CString(user)
	defer C.free(unsafe.Pointer(cUser))
	cPassword := C.CString(password)
	defer C.free(unsafe.Pointer(cPassword))
	err := C.rtdb_add_user_warp(C.rtdb_int32(handle), cUser, cPassword, C.rtdb_int32(priv))
	return RtdbError(err).GoError()
}

// RawRtdbRemoveUserWarp 删除用户帐户
//
// input:
//   - handle 连接句柄
//   - user 帐户
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_remove_user_warp(rtdb_int32 handle, const char *user)
func RawRtdbRemoveUserWarp(handle ConnectHandle, user string) error {
	cUser := C.CString(user)
	defer C.free(unsafe.Pointer(cUser))
	err := C.rtdb_remove_user_warp(C.rtdb_int32(handle), cUser)
	return RtdbError(err).GoError()
}

// RawRtdbLockUserWarp 启用或禁用用户, 只有管理员有启用禁用权限
//
// input:
//   - handle 连接句柄
//   - user 帐户名
//   - lock 是否禁用, ON(1)启动， OFF(0)禁用
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_lock_user_warp(rtdb_int32 handle, const char *user, rtdb_int8 lock)
func RawRtdbLockUserWarp(handle ConnectHandle, user string, lock Switch) error {
	cUser := C.CString(user)
	defer C.free(unsafe.Pointer(cUser))
	err := C.rtdb_lock_user_warp(C.rtdb_int32(handle), cUser, C.rtdb_int8(lock))
	return RtdbError(err).GoError()
}

// RawRtdbGetUsersWarp 获得所有用户
//
// input:
//   - handle 连接句柄
//
// output:
//   - []RtdbUserInfo 用户列表
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_users_warp(rtdb_int32 handle, rtdb_int32 *count, RTDB_USER_INFO *infos)
func RawRtdbGetUsersWarp(handle ConnectHandle) ([]RtdbUserInfo, error) {
	cCount := C.rtdb_int32(RtdbConstMaxUserCount)
	cInfos := make([]C.RTDB_USER_INFO, RtdbConstMaxUserCount)
	err := C.rtdb_get_users_warp(C.rtdb_int32(handle), &cCount, &cInfos[0])
	goInfos := make([]RtdbUserInfo, 0)
	for i := 0; i < int(cCount); i++ {
		goInfos = append(goInfos, cToRtdbUserInfo(&cInfos[i]))
	}
	return goInfos, RtdbError(err).GoError()
}

// RawRtdbAddBlacklistWarp 添加连接黑名单项
//
// input:
//   - handle  连接句柄
//   - addr 阻止连接段地址
//   - mask 阻止连接段子网掩码
//   - desc 阻止连接段的说明
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_add_blacklist_warp(rtdb_int32 handle, const char *addr, const char *mask, const char *desc)
func RawRtdbAddBlacklistWarp(handle ConnectHandle, addr string, mask string, desc string) error {
	cAddr := C.CString(addr)
	defer C.free(unsafe.Pointer(cAddr))
	cMask := C.CString(mask)
	defer C.free(unsafe.Pointer(cMask))
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))
	err := C.rtdb_add_blacklist_warp(C.rtdb_int32(handle), cAddr, cMask, cDesc)
	return RtdbError(err).GoError()
}

// RawRtdbUpdateBlacklistWarp 更新连接连接黑名单项
//
// input:
//   - handle 连接句柄
//   - oldAddr 原阻止连接段地址
//   - oldMask 原阻止连接段子网掩码
//   - newAddr 新的阻止连接段地址
//   - newMask 新的阻止连接段子网掩码
//   - newDesc 新的阻止连接段的说明，超过 511 字符将被截断
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_update_blacklist_warp(rtdb_int32 handle, const char *addr, const char *mask, const char *addr_new, const char *mask_new, const char *desc)
func RawRtdbUpdateBlacklistWarp(handle ConnectHandle, oldAddr string, oldMask string, newAddr string, newMask string, newDesc string) error {
	cOldAddr := C.CString(oldAddr)
	defer C.free(unsafe.Pointer(cOldAddr))
	cOldMask := C.CString(oldMask)
	defer C.free(unsafe.Pointer(cOldMask))
	cNewAddr := C.CString(newAddr)
	defer C.free(unsafe.Pointer(cNewAddr))
	cNewMask := C.CString(newMask)
	defer C.free(unsafe.Pointer(cNewMask))
	cNewDesc := C.CString(newDesc)
	defer C.free(unsafe.Pointer(cNewDesc))
	err := C.rtdb_update_blacklist_warp(C.rtdb_int32(handle), cOldAddr, cOldMask, cNewAddr, cNewMask, cNewDesc)
	return RtdbError(err).GoError()
}

// RawRtdbRemoveBlacklistWarp 删除连接黑名单项
//
// input:
//   - handle  连接句柄
//   - addr 阻止连接段地址
//   - mask 阻止连接段子网掩码
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_remove_blacklist_warp(rtdb_int32 handle, const char *addr, const char *mask)
func RawRtdbRemoveBlacklistWarp(handle ConnectHandle, addr string, mask string) error {
	cAddr := C.CString(addr)
	defer C.free(unsafe.Pointer(cAddr))
	cMask := C.CString(mask)
	defer C.free(unsafe.Pointer(cMask))
	err := C.rtdb_remove_blacklist_warp(C.rtdb_int32(handle), cAddr, cMask)
	return RtdbError(err).GoError()
}

// RawRtdbGetBlacklistWarp 获得连接黑名单
//
// input:
//   - handle 连接句柄
//
// output:
//   - []BlackList 黑名单列表
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_blacklist_warp(rtdb_int32 handle, char* const* addrs, char* const* masks, char* const* descs, rtdb_int32 *count)
func RawRtdbGetBlacklistWarp(handle ConnectHandle) ([]BlackList, error) {
	cAddrs := make([]*C.char, RtdbConstMaxBlacklistLen)
	for i := int32(0); i < int32(RtdbConstMaxBlacklistLen); i++ {
		cAddrs[i] = (*C.char)(C.CBytes(make([]byte, 32)))
	}
	defer func() {
		for i := int32(0); i < int32(RtdbConstMaxBlacklistLen); i++ {
			C.free(unsafe.Pointer(cAddrs[i]))
		}
	}()
	cgoAddrs := &cAddrs[0]

	cMakes := make([]*C.char, RtdbConstMaxBlacklistLen)
	for i := int32(0); i < int32(RtdbConstMaxBlacklistLen); i++ {
		cMakes[i] = (*C.char)(C.CBytes(make([]byte, 32)))
	}
	defer func() {
		for i := int32(0); i < int32(RtdbConstMaxBlacklistLen); i++ {
			C.free(unsafe.Pointer(cMakes[i]))
		}
	}()
	cgoMasks := &cMakes[0]

	cDescs := make([]*C.char, RtdbConstMaxBlacklistLen)
	for i := int32(0); i < int32(RtdbConstMaxBlacklistLen); i++ {
		cDescs[i] = (*C.char)(C.CBytes(make([]byte, 512)))
	}
	defer func() {
		for i := int32(0); i < int32(RtdbConstMaxBlacklistLen); i++ {
			C.free(unsafe.Pointer(cDescs[i]))
		}
	}()
	cgoDescs := &cDescs[0]

	cgoCount := C.rtdb_int32(RtdbConstMaxBlacklistLen)
	err := C.rtdb_get_blacklist_warp(C.rtdb_int32(handle), cgoAddrs, cgoMasks, cgoDescs, &cgoCount)

	rtn := make([]BlackList, 0)
	for i := int32(0); i < int32(cgoCount); i++ {
		rtn = append(rtn, BlackList{
			Addr: CCharArrayToString(cAddrs[i], 32),
			Mask: CCharArrayToString(cMakes[i], 32),
			Desc: CCharArrayToString(cDescs[i], 512),
		})
	}
	return rtn, RtdbError(err).GoError()
}

// RawRtdbAddAuthorizationWarp 添加信任连接段
//
// input:
//   - handle 连接句柄
//   - addr 信任连接段地址
//   - mask 信任连接段子网掩码。
//   - priv 信任连接段拥有的用户权限。
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_add_authorization_warp(rtdb_int32 handle, const char *addr, const char *mask, rtdb_int32 priv, const char *desc)
func RawRtdbAddAuthorizationWarp(handle ConnectHandle, addr string, mask string, desc string, priv PrivGroup) error {
	cAddr := C.CString(addr)
	defer C.free(unsafe.Pointer(cAddr))
	cMask := C.CString(mask)
	defer C.free(unsafe.Pointer(cMask))
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))
	err := C.rtdb_add_authorization_warp(C.rtdb_int32(handle), cAddr, cMask, C.rtdb_int32(priv), cDesc)
	return RtdbError(err).GoError()
}

// RawRtdbUpdateAuthorizationWarp 更新信任连接段
//
// input:
//   - handle 连接句柄
//   - oldAddr 原信任连接段地址
//   - oldMask 原信任连接段子网掩码
//   - newAddr 新的信任连接段地址
//   - newMask 新的信任连接段子网掩码
//   - newDesc 新的信任连接段的说明，超过 511 字符将被截断
//   - priv 新的信任连接段拥有的用户权限
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_update_authorization_warp(rtdb_int32 handle, const char *addr, const char *mask, const char *addr_new, const char *mask_new, rtdb_int32 priv, const char *desc)
func RawRtdbUpdateAuthorizationWarp(handle ConnectHandle, oldAddr string, oldMask string, newAddr string, newMask string, newDesc string, priv PrivGroup) error {
	cOldAddr := C.CString(oldAddr)
	defer C.free(unsafe.Pointer(cOldAddr))
	cOldMask := C.CString(oldMask)
	defer C.free(unsafe.Pointer(cOldMask))
	cNewAddr := C.CString(newAddr)
	defer C.free(unsafe.Pointer(cNewAddr))
	cNewMask := C.CString(newMask)
	defer C.free(unsafe.Pointer(cNewMask))
	cNewDesc := C.CString(newDesc)
	defer C.free(unsafe.Pointer(cNewDesc))
	err := C.rtdb_update_authorization_warp(C.rtdb_int32(handle), cOldAddr, cOldMask, cNewAddr, cNewMask, C.rtdb_int32(priv), cNewDesc)
	return RtdbError(err).GoError()
}

// RawRtdbRemoveAuthorizationWarp 删除信任连接段
//
// input:
//   - handle 连接句柄
//   - addr 信任连接段地址
//   - mask 信任连接段子网掩码
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_remove_authorization_warp(rtdb_int32 handle, const char *addr, const char *mask)
func RawRtdbRemoveAuthorizationWarp(handle ConnectHandle, addr string, mask string) error {
	cAddr := C.CString(addr)
	defer C.free(unsafe.Pointer(cAddr))
	cMask := C.CString(mask)
	defer C.free(unsafe.Pointer(cMask))
	err := C.rtdb_remove_authorization_warp(C.rtdb_int32(handle), cAddr, cMask)
	return RtdbError(err).GoError()
}

// RawRtdbGetAuthorizationsWarp 获得所有信任连接段
//
// input:
//   - handle 连接句柄
//
// output:
//   - []AuthorizationsList 白名单列表
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_authorizations_warp(rtdb_int32 handle, char* const* addrs, char* const* masks, rtdb_int32 *privs, char* const* descs, rtdb_int32 *count)
func RawRtdbGetAuthorizationsWarp(handle ConnectHandle) ([]AuthorizationsList, error) {
	cAddrs := make([]*C.char, RtdbConstMaxAuthCount)
	for i := int32(0); i < int32(RtdbConstMaxAuthCount); i++ {
		cAddrs[i] = (*C.char)(C.CBytes(make([]byte, 32)))
	}
	defer func() {
		for i := int32(0); i < int32(RtdbConstMaxAuthCount); i++ {
			C.free(unsafe.Pointer(cAddrs[i]))
		}
	}()
	cgoAddrs := &cAddrs[0]

	cMakes := make([]*C.char, RtdbConstMaxAuthCount)
	for i := int32(0); i < int32(RtdbConstMaxAuthCount); i++ {
		cMakes[i] = (*C.char)(C.CBytes(make([]byte, 32)))
	}
	defer func() {
		for i := int32(0); i < int32(RtdbConstMaxAuthCount); i++ {
			C.free(unsafe.Pointer(cMakes[i]))
		}
	}()
	cgoMasks := &cMakes[0]

	cDescs := make([]*C.char, RtdbConstMaxAuthCount)
	for i := int32(0); i < int32(RtdbConstMaxAuthCount); i++ {
		cDescs[i] = (*C.char)(C.CBytes(make([]byte, 512)))
	}
	defer func() {
		for i := int32(0); i < int32(RtdbConstMaxAuthCount); i++ {
			C.free(unsafe.Pointer(cDescs[i]))
		}
	}()
	cgoDescs := &cDescs[0]
	cgoCount := C.rtdb_int32(RtdbConstMaxAuthCount)

	privs := make([]PrivGroup, int32(RtdbConstMaxAuthCount))
	cgoPrivs := (*C.rtdb_int32)(unsafe.Pointer(&privs[0]))
	err := C.rtdb_get_authorizations_warp(C.rtdb_int32(handle), cgoAddrs, cgoMasks, cgoPrivs, cgoDescs, &cgoCount)

	rtn := make([]AuthorizationsList, 0)
	for i := int32(0); i < int32(cgoCount); i++ {
		rtn = append(rtn, AuthorizationsList{
			Addr: CCharArrayToString(cAddrs[i], 32),
			Mask: CCharArrayToString(cMakes[i], 32),
			Desc: CCharArrayToString(cDescs[i], 512),
			Priv: privs[i],
		})
	}
	return rtn, RtdbError(err).GoError()
}

// RawRtdbHostTimeWarp 获取 RTDB 服务器当前UTC时间
// 备注：32和64位时间戳，统一使用64位, 因此屏蔽32位时间戳
//
// input:
//   - handle       连接句柄
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_host_time_warp(rtdb_int32 handle, rtdb_int32 *hosttime)
// func RawRtdbHostTimeWarp(handle ConnectHandle) (DateTimeType, error) {
// 	hostTime := C.rtdb_int32(0)
// 	err := C.rtdb_host_time_warp(C.rtdb_int32(handle), &hostTime)
// 	return DateTimeType(hostTime), RtdbError(err).GoError()
// }

// RawRtdbHostTime64Warp 获取 RTDB 服务器当前UTC时间
//
// input:
//   - handle       连接句柄
//
// output:
//   - TimestampType 服务器当前UDC时间, 单位秒
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_host_time64_warp(rtdb_int32 handle, rtdb_timestamp_type* hosttime)
func RawRtdbHostTime64Warp(handle ConnectHandle) (TimestampType, error) {
	ts := C.rtdb_timestamp_type(0)
	err := C.rtdb_host_time64_warp(C.rtdb_int32(handle), &ts)
	return TimestampType(ts), RtdbError(err).GoError()
}

// RawRtdbFormatTimespanWarp 根据时间跨度值生成时间格式字符串, 如：输入10， 输出10s, 输入60，输出1n
//
// input:
//
//   - timespan 要处理的时间跨度秒数, 跨度单位如下，备注：这是遵循工业Pi数据库的标准, 和通用标准稍有不同
//     ?y    ?年, 1年 = 365日
//     ?m    ?月, 1月 = 30 日
//     ?d    ?日
//     ?h    ?小时
//     ?n    ?分钟
//     ?s    ?秒
//
// output:
//   - string 时间格式字符串
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_format_timespan_warp(char *str, rtdb_int32 timespan)
func RawRtdbFormatTimespanWarp(timespan int32) (string, error) {
	cgoStr := (*C.char)(C.CBytes(make([]byte, 512)))
	defer C.free(unsafe.Pointer(cgoStr))
	cgoDatetime := C.rtdb_int32(timespan)
	err := C.rtdb_format_timespan_warp(cgoStr, cgoDatetime)
	tStr := C.GoString(cgoStr)
	return tStr, RtdbError(err).GoError()
}

// RawRtdbParseTimespanWarp 根据时间格式字符串解析时间跨度值, 如：输入2n，输出120，表示2分钟
//
// input:
//   - 时间格式字符串，规则如下：
//     [time_span]
//     时间跨度部分可以出现多次，
//     可用的时间跨度代码及含义如下：
//     ?y            ?年, 1年 = 365日
//     ?m            ?月, 1月 = 30 日
//     ?d            ?日
//     ?h            ?小时
//     ?n            ?分钟
//     ?s            ?秒
//     例如："1d" 表示时间跨度为24小时。
//
// output:
//   - DateTimeType 时间跨度值
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_parse_timespan_warp(const char *str, rtdb_int32 *timespan)
func RawRtdbParseTimespanWarp(tStr string) (DateTimeType, error) {
	cStr := C.CString(tStr)
	defer C.free(unsafe.Pointer(cStr))
	ts := C.rtdb_int32(0)
	err := C.rtdb_parse_timespan_warp(cStr, &ts)
	return DateTimeType(ts), RtdbError(err).GoError()
}

// RawRtdbParseTimeWarp 根据时间格式字符串解析时间值
//
// input:
//   - str          字符串，输入，时间格式字符串，规则如下：
//     base_time [+|- offset_time]
//     其中 base_time 表示基本时间，有三种形式：
//     1. 时间字符串，如 "2010-1-1" 及 "2010-1-1 8:00:00"；
//     2. 时间代码，表示客户端相对时间；
//     可用的时间代码及含义如下：
//     td             当天零点
//     yd             昨天零点
//     tm             明天零点
//     mon            本周一零点
//     tue            本周二零点
//     wed            本周三零点
//     thu            本周四零点
//     fri            本周五零点
//     sat            本周六零点
//     sun            本周日零点
//     3. 星号('*')，表示客户端当前时间。
//     offset_time 是可选的，可以出现多次，
//     可用的时间偏移代码及含义如下：
//     [+|-] ?y            偏移?年, 1年 = 365日
//     [+|-] ?m            偏移?月, 1月 = 30 日
//     [+|-] ?d            偏移?日
//     [+|-] ?h            偏移?小时
//     [+|-] ?n            偏移?分钟
//     [+|-] ?s            偏移?秒
//     [+|-] ?ms           偏移?毫秒
//     例如："*-1d" 表示当前时刻减去24小时。
//
// output:
//   - TimestampType 时间戳，秒级
//   - SubtimeType 时间戳，亚秒级(亚秒指的是 毫秒、微妙、纳秒 之一， 需要根据当前标签点的时间戳精度来确定单位)
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_parse_time_warp(const char *str, rtdb_int64 *datetime, rtdb_int16 *ms)
func RawRtdbParseTimeWarp(tStr string) (TimestampType, SubtimeType, error) {
	cStr := C.CString(tStr)
	defer C.free(unsafe.Pointer(cStr))
	ts := C.rtdb_int64(0)
	ms := C.rtdb_int16(0)
	err := C.rtdb_parse_time_warp(cStr, &ts, &ms)
	return TimestampType(ts), SubtimeType(ms), RtdbError(err).GoError()
}

// RawRtdbFormatMessageWarp 获取 Rtdb API 调用返回值的简短描述(错误码对应的Desc)
//
// input:
//   - err 错误码
//
// output:
//   - name 函数名
//   - message 函数描述
//
// raw_fn:
//   - void RTDBAPI_CALLRULE rtdb_format_message_warp(rtdb_error ecode, char *message, char *name, rtdb_int32 size)
func RawRtdbFormatMessageWarp(err RtdbError) (string, string) {
	cgoErr := C.rtdb_error(err)
	cgoMessage := (*C.char)(C.CBytes(make([]byte, 10240)))
	defer C.free(unsafe.Pointer(cgoMessage))
	cgoName := (*C.char)(C.CBytes(make([]byte, 10240)))
	defer C.free(unsafe.Pointer(cgoName))
	cgoSize := C.rtdb_int32(10240)
	C.rtdb_format_message_warp(cgoErr, cgoMessage, cgoName, cgoSize)
	return C.GoString(cgoName), C.GoString(cgoMessage)
}

// RawRtdbJobMessageWarp 获取任务的简短描述
//
// input:
//   - jobID RTDB_HOST_CONNECT_INFO::job 字段所表示的最近任务的描述
//
// output:
//   - name Job名称
//   - desc Job描述
//
// raw_fn:
//   - void RTDBAPI_CALLRULE rtdb_job_message_warp(rtdb_int32 job_id, char *desc, char *name, rtdb_int32 size)
func RawRtdbJobMessageWarp(jobID int32) (string, string) {
	cgoDesc := (*C.char)(C.CBytes(make([]byte, 1024)))
	defer C.free(unsafe.Pointer(cgoDesc))
	cgoName := (*C.char)(C.CBytes(make([]byte, 1024)))
	defer C.free(unsafe.Pointer(cgoName))
	cgoSize := C.rtdb_int32(1024)
	cgoJob := C.rtdb_int32(jobID)
	C.rtdb_job_message_warp(cgoJob, cgoDesc, cgoName, cgoSize)
	return C.GoString(cgoName), C.GoString(cgoDesc)
}

// RawRtdbSetTimeoutWarp 设置连接超时时间
//
// input:
//   - handle   连接句柄
//   - socket   整型，输入，要设置超时时间的连接
//   - timeout  整型，输入，超时时间，单位为秒，0 表示始终保持
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_set_timeout_warp(rtdb_int32 handle, rtdb_int32 socket, rtdb_int32 timeout)
func RawRtdbSetTimeoutWarp(handle ConnectHandle, socket SocketHandle, timeout DateTimeType) error {
	err := C.rtdb_set_timeout_warp(C.rtdb_int32(handle), C.rtdb_int32(socket), C.rtdb_int32(timeout))
	return RtdbError(err).GoError()
}

// RawRtdbGetTimeoutWarp 获得连接超时时间
//
// input:
//   - handle 连接句柄
//   - sockt 要获取超时时间的连接
//
// output:
//   - DateTimeType 连接超时时间，单位秒
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_timeout_warp(rtdb_int32 handle, rtdb_int32 socket, rtdb_int32 *timeout)
func RawRtdbGetTimeoutWarp(handle ConnectHandle, socket SocketHandle) (DateTimeType, error) {
	timeout := C.rtdb_int32(0)
	err := C.rtdb_get_timeout_warp(C.rtdb_int32(handle), C.rtdb_int32(socket), &timeout)
	return DateTimeType(timeout), RtdbError(err).GoError()
}

// RawRtdbKillConnectionWarp 断开已知连接
//
// input:
//   - handle 连接句柄
//   - socket 要断开的连接
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_kill_connection_warp(rtdb_int32 handle, rtdb_int32 socket)
func RawRtdbKillConnectionWarp(handle ConnectHandle, socket SocketHandle) error {
	err := C.rtdb_kill_connection_warp(C.rtdb_int32(handle), C.rtdb_int32(socket))
	return RtdbError(err).GoError()
}

// RawRtdbGetLogicalDriversWarp 获得逻辑盘符
//
// input:
//   - handle 连接句柄
//
// output:
//   - []string 盘符数组
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_logical_drivers_warp(rtdb_int32 handle, char *drivers)
func RawRtdbGetLogicalDriversWarp(handle ConnectHandle) ([]string, error) {
	drives := make([]byte, 512)
	cDrives := (*C.char)(unsafe.Pointer(&drives[0]))
	err := C.rtdb_get_logical_drivers_warp(C.rtdb_int32(handle), cDrives)
	sDs := C.GoString(cDrives)
	rtn := make([]string, 0)
	for _, c := range sDs {
		rtn = append(rtn, string(c))
	}
	return rtn, RtdbError(err).GoError()
}

// RawRtdbOpenPathWarp 打开目录以便遍历其中的文件和子目录。
//
// input:
//   - handle 连接句柄
//   - dir 要打开的目录
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_open_path_warp(rtdb_int32 handle, const char *dir)
func RawRtdbOpenPathWarp(handle ConnectHandle, dir string) error {
	cDir := C.CString(dir)
	defer C.free(unsafe.Pointer(cDir))
	err := C.rtdb_open_path_warp(C.rtdb_int32(handle), cDir)
	return RtdbError(err).GoError()
}

// RawRtdbReadPathWarp 读取目录中的文件或子目录
// 备注：此函数返回的时间戳是 32位的，暂不实现，统一使用64位时间戳
//
// input:
//   - handle 连接句柄
//
// output:
//   - DirItem
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_read_path_warp(rtdb_int32 handle, char *path, rtdb_int16 *is_dir, rtdb_int32 *atime, rtdb_int32 *ctime, rtdb_int32 *mtime, rtdb_int64 *size)
// func RawRtdbReadPathWarp(handle ConnectHandle) (DirItem, error) {
// 	cgoHandle := C.rtdb_int32(handle)
// 	cgoPath := (*C.char)(C.CBytes(make([]byte, RtdbMaxPath)))
// 	defer C.free(unsafe.Pointer(cgoPath))
// 	cgoIsDir := C.rtdb_int16(0)
// 	cgoATime := C.rtdb_int32(0)
// 	cgoCTime := C.rtdb_int32(0)
// 	cgoMTime := C.rtdb_int32(0)
// 	cgoSize := C.rtdb_int64(0)
// 	err := C.rtdb_read_path(cgoHandle, cgoPath, &cgoIsDir, &cgoATime, &cgoCTime, &cgoMTime, &cgoSize)
//
// 	rtnPath := C.GoString(cgoPath)
// 	rtnIsDir := false
// 	if cgoIsDir > 0 {
// 		rtnIsDir = true
// 	}
// 	rtnATime := int32(cgoATime)
// 	rtnCTime := int32(cgoCTime)
// 	rtnMTime := int32(cgoMTime)
// 	rtnSize := int64(cgoSize)
//
// 	item := DirItem{
// 		Path:  rtnPath,
// 		IsDir: rtnIsDir,
// 		ATime: DateTimeType(rtnATime),
// 		CTime: DateTimeType(rtnCTime),
// 		MTime: DateTimeType(rtnMTime),
// 		Size:  rtnSize,
// 	}
//
// 	return item, RtdbError(err).GoError()
// }

// RawRtdbReadPath64Warp 读取目录中的文件或子目录
//
// input:
//   - handle 连接句柄
//
// output:
//   - DirItem 目录项
//
// err_code:
//   - 当返回值为 RteBatchEnd 时表示目录下所有子目录和文件已经遍历完毕。
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_read_path64_warp(rtdb_int32 handle, char* path, rtdb_int16* is_dir, rtdb_timestamp_type* atime, rtdb_timestamp_type* ctime, rtdb_timestamp_type* mtime, rtdb_int64* size)
func RawRtdbReadPath64Warp(handle ConnectHandle) (DirItem, error) {
	cgoHandle := C.rtdb_int32(handle)
	cgoPath := (*C.char)(C.CBytes(make([]byte, RtdbMaxPath)))
	defer C.free(unsafe.Pointer(cgoPath))
	cgoIsDir := C.rtdb_int16(0)
	cgoATime := C.rtdb_timestamp_type(0)
	cgoCTime := C.rtdb_timestamp_type(0)
	cgoMTime := C.rtdb_timestamp_type(0)
	cgoSize := C.rtdb_int64(0)
	err := C.rtdb_read_path64_warp(cgoHandle, cgoPath, &cgoIsDir, &cgoATime, &cgoCTime, &cgoMTime, &cgoSize)

	rtnIsDir := Switch(cgoIsDir)
	rtnPath := C.GoString(cgoPath)
	rtnATime := TimestampType(cgoATime)
	rtnCTime := TimestampType(cgoCTime)
	rtnMTime := TimestampType(cgoMTime)
	rtnSize := int64(cgoSize)

	item := DirItem{
		Path:  rtnPath,
		IsDir: rtnIsDir,
		ATime: rtnATime,
		CTime: rtnCTime,
		MTime: rtnMTime,
		Size:  rtnSize,
	}

	return item, RtdbError(err).GoError()
}

// RawRtdbClosePathWarp 关闭当前遍历的目录
//
// input:
//   - handle 连接句柄
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_close_path_warp(rtdb_int32 handle)
func RawRtdbClosePathWarp(handle ConnectHandle) error {
	err := C.rtdb_close_path_warp(C.rtdb_int32(handle))
	return RtdbError(err).GoError()
}

// RawRtdbMkdirWarp 建立目录
//
// input:
//   - handle 连接句柄
//   - dirName 新建目录的全路径
//
// output:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_mkdir_warp(rtdb_int32 handle, const char *dir)
func RawRtdbMkdirWarp(handle ConnectHandle, dirName string) error {
	cDirName := C.CString(dirName)
	defer C.free(unsafe.Pointer(cDirName))
	err := C.rtdb_mkdir_warp(C.rtdb_int32(handle), cDirName)
	return RtdbError(err).GoError()
}

// RawRtdbGetFileSizeWarp 获得指定服务器端文件的大小
//
// input:
//   - handle 连接句柄
//   - file 文件名
//
// output:
//   - int64 文件大小
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_file_size_warp(rtdb_int32 handle, const char *file, rtdb_int64 *size)
func RawRtdbGetFileSizeWarp(handle ConnectHandle, filePath string) (int64, error) {
	cFilePath := C.CString(filePath)
	defer C.free(unsafe.Pointer(cFilePath))
	cSize := C.rtdb_int64(0)
	err := C.rtdb_get_file_size_warp(C.rtdb_int32(handle), cFilePath, &cSize)
	return int64(cSize), RtdbError(err).GoError()
}

// RawRtdbReadFileWarp 读取服务器端指定文件的内容
//
// input:
//   - handle 连接句柄
//   - fileName 要读取内容的文件名
//   - pos 读取文件的起始位置
//
// output:
//   - []byte 读取出来的数据
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_read_file_warp(rtdb_int32 handle, const char *file, char *content, rtdb_int64 pos, rtdb_int64 *size)
func RawRtdbReadFileWarp(handle ConnectHandle, filePath string, pos int64, cacheSize int64) ([]byte, error) {
	cFilePath := C.CString(filePath)
	defer C.free(unsafe.Pointer(cFilePath))
	buf := make([]byte, cacheSize)
	cBuf := (*C.char)(unsafe.Pointer(&buf[0]))
	cSize := C.rtdb_int64(cacheSize)
	err := C.rtdb_read_file_warp(C.rtdb_int32(handle), cFilePath, cBuf, C.rtdb_int64(pos), &cSize)
	return buf[:int64(cSize)], RtdbError(err).GoError()
}

// RawRtdbGetMaxBlobLenWarp 取得数据库允许的blob与str类型测点的最大长度
//
// input:
//   - handle       连接句柄
//
// output:
//   - int32 blob/str最大长度
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_get_max_blob_len_warp(rtdb_int32 handle, rtdb_int32 *len)
func RawRtdbGetMaxBlobLenWarp(handle ConnectHandle) (int32, error) {
	cLen := C.rtdb_int32(0)
	err := C.rtdb_get_max_blob_len_warp(C.rtdb_int32(handle), &cLen)
	return int32(cLen), RtdbError(err).GoError()
}

// RawRtdbFormatQualityWarp 取得质量码对应的定义
//
// input:
//   - handle       连接句柄
//   - qualities    质量码数组
//     OPC质量码把8位分3部分定义：XX XXXX XX，对应着：品质位域 分状态位域 限定位域
//     品质位域：00（Bad）01（Uncertain）10（N/A）11（Good）
//     分状态位域：不同品质位域对应各自的分状态位域
//     限定位域：00（Not limited）01（Low limited）10（high limited）11（Constant）
//     三个域之间用逗号隔开，输出到definitions参数中，前面有有RTDB，OPC或者USER标识，说明标签点类别
//
// output:
//   - []string 质量码对应的说明
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_format_quality_warp(rtdb_int32 handle, rtdb_int32 *count, rtdb_int16 *qualities, rtdb_byte **definitions, rtdb_int32 *lens)
func RawRtdbFormatQualityWarp(handle ConnectHandle, qualities []Quality) ([]string, error) {
	cgoHandle := C.rtdb_int32(handle)
	count := len(qualities)
	cgoCount := C.rtdb_int32(count)
	cgoQualities := (*C.rtdb_int16)(unsafe.Pointer(&qualities[0]))
	definitions := make([]*C.rtdb_byte, count)
	defer func() {
		for i := 0; i < int(count); i++ {
			C.free(unsafe.Pointer(definitions[i]))
		}
	}()
	for i := 0; i < int(cgoCount); i++ {
		definitions[i] = (*C.rtdb_byte)(C.CBytes(make([]byte, 256)))
	}
	cgoDefinitions := (**C.rtdb_byte)(unsafe.Pointer(&definitions[0]))
	lens := make([]int32, count)
	cgoLens := (*C.rtdb_int32)(unsafe.Pointer(&lens[0]))
	err := C.rtdb_format_quality_warp(cgoHandle, &cgoCount, cgoQualities, cgoDefinitions, cgoLens)

	rtnDefinitions := make([]string, 0)
	for i := 0; i < int(cgoCount); i++ {
		bs := C.GoBytes(unsafe.Pointer(definitions[i]), 256)
		st := string(bs[3:lens[i]])
		rtnDefinitions = append(rtnDefinitions, st)
	}

	return rtnDefinitions, RtdbError(err).GoError()
}

// RawRtdbJudgeConnectStatusWarp 判断连接是否可用
//
// input:
//   - handle   连接句柄
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdb_judge_connect_status_warp(rtdb_int32 handle, rtdb_int8* change_connection GAPI_DEFAULT_VALUE(0), char* current_ip_addr GAPI_DEFAULT_VALUE(0), rtdb_int32 size GAPI_DEFAULT_VALUE(0))
func RawRtdbJudgeConnectStatusWarp(handle ConnectHandle) error {
	cgoHandle := C.rtdb_int32(handle)
	cgoChangeConnection := C.rtdb_int8(0)
	cgoAddr := (*C.char)(C.CBytes(make([]byte, 1024)))
	defer C.free(unsafe.Pointer(cgoAddr))
	err := C.rtdb_judge_connect_status_warp(cgoHandle, &cgoChangeConnection, cgoAddr, 1024)
	return RtdbError(err).GoError()
}

// RawRtdbFormatIpaddrWarp 将整形IP转换为字符串形式的IP
//
// input:
//   - ip 整数类型IP(IPv4)
//
// output:
//   - string 字符串型IP(IPv4)
//
// raw_fn:
//   - void RTDBAPI_CALLRULE rtdb_format_ipaddr_warp(rtdb_uint32 ip, char* ip_addr, rtdb_int32 size)
func RawRtdbFormatIpaddrWarp(ip uint32) string {
	addr := make([]byte, 128)
	cAddr := (*C.char)(unsafe.Pointer(&addr[0]))
	C.rtdb_format_ipaddr_warp(C.rtdb_uint32(ip), cAddr, 128)
	return C.GoString(cAddr)
}

// RawRtdbbAppendTableWarp 添加新表
//
// input:
//   - handle   连接句柄
//   - tableName 表名
//   - tableDesc 表描述
//
// output:
//   - RtdbTable 返回表
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_append_table_warp(rtdb_int32 handle, RTDB_TABLE *field)
func RawRtdbbAppendTableWarp(handle ConnectHandle, tableName, tableDesc string) (RtdbTable, error) {
	cgoHandle := C.rtdb_int32(handle)
	table := C.RTDB_TABLE{}
	GoStringToCCharArray(tableName, &table.name[0], int(C.RTDB_TAG_SIZE))
	GoStringToCCharArray(tableDesc, &table.desc[0], int(C.RTDB_DESC_SIZE))
	err := C.rtdbb_append_table_warp(cgoHandle, &table)
	return cToRtdbTable(&table), RtdbError(err).GoError()
}

// RawRtdbbRemoveTableByIdWarp 根据表 id 删除表及表中标签点
//
// input:
//   - handle        连接句柄
//   - id            整型，输入，表 id
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_remove_table_by_id_warp(rtdb_int32 handle, rtdb_int32 id)
func RawRtdbbRemoveTableByIdWarp(handle ConnectHandle, tableID TableID) error {
	err := C.rtdbb_remove_table_by_id_warp(C.rtdb_int32(handle), C.rtdb_int32(tableID))
	return RtdbError(err).GoError()
}

// RawRtdbbRemoveTableByNameWarp 根据表名删除表及表中标签点
//
// input:
//   - handle 连接句柄
//   - name 表名称
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_remove_table_by_name_warp(rtdb_int32 handle, const char *name)
func RawRtdbbRemoveTableByNameWarp(handle ConnectHandle, name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	err := C.rtdbb_remove_table_by_name_warp(C.rtdb_int32(handle), cName)
	return RtdbError(err).GoError()
}

// RawRtdbbTablesCountWarp 取得标签点表总数
//
// input:
//   - handle 连接句柄
//
// output:
//   - in32 标签点表总数
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_tables_count_warp(rtdb_int32 handle, rtdb_int32 *count)
func RawRtdbbTablesCountWarp(handle ConnectHandle) (int32, error) {
	cCount := C.rtdb_int32(0)
	err := C.rtdbb_tables_count_warp(C.rtdb_int32(handle), &cCount)
	return int32(cCount), RtdbError(err).GoError()
}

// RawRtdbbGetTablesWarp 取得所有标签点表的ID
//
// input:
//   - handle 连接句柄
//
// output:
//   - []TableID 获取表ID列表
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_tables_warp(rtdb_int32 handle, rtdb_int32 *ids, rtdb_int32 *count)
func RawRtdbbGetTablesWarp(handle ConnectHandle) ([]TableID, error) {
	count, err := RawRtdbbTablesCountWarp(handle)
	if err != nil {
		return nil, err
	}
	ids := make([]TableID, count)
	cgoIDs := (*C.rtdb_int32)(unsafe.Pointer(&ids[0]))
	cgoCount := C.rtdb_int32(count)
	e := C.rtdbb_get_tables_warp(C.rtdb_int32(handle), cgoIDs, &cgoCount)
	return ids[:cgoCount], RtdbError(e).GoError()
}

// RawRtdbbGetTableSizeByIdWarp 根据表 id 获取表中包含的标签点数量(大概数量, 包含被标记删除的点)
//
// input:
//   - handle 连接句柄
//   - tableID 表ID
//
// output:
//   - int32 表中标签点的数量
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_size_by_id_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 *size)
func RawRtdbbGetTableSizeByIdWarp(handle ConnectHandle, tableID TableID) (int32, error) {
	cSize := C.rtdb_int32(0)
	err := C.rtdbb_get_table_size_by_id_warp(C.rtdb_int32(handle), C.rtdb_int32(tableID), &cSize)
	return int32(cSize), RtdbError(err).GoError()
}

// RawRtdbbGetTableSizeByNameWarp 根据表名称获取表中包含的标签点数量(大概数量, 包含被标记删除的点)
//
// input:
//   - handle 连接句柄
//   - tableName 表名
//
// output:
//   - int32 表中标签点的数量
//
// raw_fn:
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_size_by_name_warp(rtdb_int32 handle, const char *name, rtdb_int32 *size)
func RawRtdbbGetTableSizeByNameWarp(handle ConnectHandle, tableName string) (int32, error) {
	cSize := C.rtdb_int32(0)
	cName := C.CString(tableName)
	defer C.free(unsafe.Pointer(cName))
	err := C.rtdbb_get_table_size_by_name_warp(C.rtdb_int32(handle), cName, &cSize)
	return int32(cSize), RtdbError(err).GoError()
}

// RawRtdbbGetTableRealSizeByIdWarp 根据表 id 获取表中实际包含的标签点数量(实际数量, 不含被删除的点)
//
// input:
//   - handle 连接句柄
//   - tableID 表ID
//
// output:
//   - int32 表中标签点的实际数量
//
// raw_fn:
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_real_size_by_id_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 *size)
func RawRtdbbGetTableRealSizeByIdWarp(handle ConnectHandle, tableID TableID) (int32, error) {
	cSize := C.rtdb_int32(0)
	err := C.rtdbb_get_table_real_size_by_id_warp(C.rtdb_int32(handle), C.rtdb_int32(tableID), &cSize)
	return int32(cSize), RtdbError(err).GoError()
}

// RawRtdbbGetTablePropertyByIdWarp 根据标签点表 id 获取表属性
//
// input:
//   - handle 连接句柄
//
// output:
//   - RtdbTable 表属性
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_property_by_id_warp(rtdb_int32 handle, RTDB_TABLE *field)
func RawRtdbbGetTablePropertyByIdWarp(handle ConnectHandle, tableID TableID) (RtdbTable, error) {
	table := C.RTDB_TABLE{}
	table.id = C.rtdb_int32(tableID)
	err := C.rtdbb_get_table_property_by_id_warp(C.rtdb_int32(handle), &table)
	return cToRtdbTable(&table), RtdbError(err).GoError()
}

// RawRtdbbGetTablePropertyByNameWarp 根据表名获取标签点表属性
//
// input:
//   - handle 连接句柄
//   - tableName 表名
//
// output:
//   - RtdbTable 表结构
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_table_property_by_name_warp(rtdb_int32 handle, RTDB_TABLE *field)
func RawRtdbbGetTablePropertyByNameWarp(handle ConnectHandle, tableName string) (RtdbTable, error) {
	table := C.RTDB_TABLE{}
	GoStringToCCharArray(tableName, &table.name[0], int(C.RTDB_TAG_SIZE))
	err := C.rtdbb_get_table_property_by_name_warp(C.rtdb_int32(handle), &table)
	return cToRtdbTable(&table), RtdbError(err).GoError()
}

// RawRtdbbInsertPointWarp 使用完整的属性集来创建单个标签点
// 备注：不实现，统一使用最大长度Calc
//
// *  \param handle 连接句柄
// *  \param base RTDB_POINT 结构，输入/输出，
// *       输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
// *  \param scan RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
// *  \param calc RTDB_CALC_POINT 结构，输入，计算标签点扩展属性集。
// *  \remark 如果新建的标签点没有对应的扩展属性集，可置为空指针。
// rtdb_error RTDBAPI_CALLRULE rtdbb_insert_point_warp(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_CALC_POINT *calc)
// func RawRtdbbInsertPointWarp() {}

// RawRtdbbInsertMaxPointWarp 使用最大长度的完整属性集来创建单个标签点
//
// input:
//   - handle 连接句柄
//   - base 基本点信息
//     输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
//   - scan 采集点信息，RtdbClass 非采集点的时候，可填nil
//   - calc 计算点信息, RtdbClass 非计算点的时候，可填nil
//   - 备注：如果新建的标签点没有对应的扩展属性集，可置为空指针。
//
// output:
//   - RtdbPoint 基本点信息
//   - RtdbScan 采集点信息
//   - RtdbCalc 计算点信息
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_insert_max_point_warp(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_MAX_CALC_POINT *calc)
func RawRtdbbInsertMaxPointWarp(handle ConnectHandle, base *RtdbPoint, scan *RtdbScan, calc *RtdbCalc) (*RtdbPoint, *RtdbScan, *RtdbCalc, error) {
	cBase := goToCRtdbPoint(base)
	cScan := goToCRtdbScan(scan)
	cCalc := goToCRtdbCalc(calc)
	err := C.rtdbb_insert_max_point_warp(C.rtdb_int32(handle), cBase, cScan, cCalc)
	return cToRtdbPoint(cBase), cToRtdbScan(cScan), cToRtdbCalc(cCalc), RtdbError(err).GoError()
}

// RawRtdbbRemovePointByIdWarp 根据 id 删除单个标签点
//
// input:
//   - handle 连接句柄
//   - id 标签点标识
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_remove_point_by_id_warp(rtdb_int32 handle, rtdb_int32 id)
func RawRtdbbRemovePointByIdWarp(handle ConnectHandle, id PointID) error {
	err := C.rtdbb_remove_point_by_id_warp(C.rtdb_int32(handle), C.rtdb_int32(id))
	return RtdbError(err).GoError()
}

// RawRtdbbRemovePointByNameWarp 根据标签点全名删除单个标签点
//
// input:
//   - handle 连接句柄
//   - tableDotTag  标签点全名称："表名.标签点名"
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_remove_point_by_name_warp(rtdb_int32 handle, const char *table_dot_tag)
func RawRtdbbRemovePointByNameWarp(handle ConnectHandle, tableDotTag string) error {
	cTableDotTag := C.CString(tableDotTag)
	defer C.free(unsafe.Pointer(cTableDotTag))
	err := C.rtdbb_remove_point_by_name_warp(C.rtdb_int32(handle), cTableDotTag)
	return RtdbError(err).GoError()
}

// RawRtdbbInsertMaxPointsWarp 使用最大长度的完整属性集来批量创建标签点
// 备注：暂不实现
//
// * [handle] 连接句柄
// * [count] count, 输入，base/scan/calc数组的长度；输出，成功的个数
// * [bases] RTDB_POINT 数组，输入/输出，
// * 输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
// * [scans] RTDB_SCAN_POINT 数组，输入，采集标签点扩展属性集。
// * [calcs] RTDB_MAX_CALC_POINT 数组，输入，计算标签点扩展属性集。
// * [errors] rtdb_error数组，输出，对应每个标签点的结果
// * 备注：如果新建的标签点没有对应的扩展属性集，可置为空指针。
// rtdb_error RTDBAPI_CALLRULE rtdbb_insert_max_points_warp(rtdb_int32 handle, rtdb_int32* count, RTDB_POINT* bases, RTDB_SCAN_POINT* scans, RTDB_MAX_CALC_POINT* calcs, rtdb_error* errors)
// func RawRtdbbInsertMaxPointsWarp(handle ConnectHandle, points []RtdbPoint, scans []RtdbScan, calcs []RtdbCalc) /*([]RtdbPoint, []RtdbScan, []RtdbCalc, []error, error)*/ {}

// RawRtdbbInsertBasePointWarp 使用最小的属性集来创建单个标签点
// 备注：不实现，统一使用最大长度
//
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
// func RawRtdbbInsertBasePointWarp() {}

// RawRtdbbInsertNamedTypePointWarp 使用完整的属性集来创建单个自定义数据类型标签点
//
// input:
//   - handle 连接句柄
//   - base 基本点信息
//   - scan 采集点信息
//   - name 自定义类型名称
//
// output:
//   - RtdbPoint 基本点信息
//   - RtdbScan 采集点信息
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_insert_named_type_point_warp(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, const char* name)
func RawRtdbbInsertNamedTypePointWarp(handle ConnectHandle, base *RtdbPoint, scan *RtdbScan, name string) (*RtdbPoint, *RtdbScan, error) {
	cBase := goToCRtdbPoint(base)
	cScan := goToCRtdbScan(scan)
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	err := C.rtdbb_insert_named_type_point_warp(C.rtdb_int32(handle), cBase, cScan, cName)
	return cToRtdbPoint(cBase), cToRtdbScan(cScan), RtdbError(err).GoError()
}

// RawRtdbbMovePointByIdWarp 根据 id 移动单个标签点到其他表
//
// input:
//   - handle 连接句柄
//   - pointID 标签点标识
//   - tableName 表名称
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_move_point_by_id_warp(rtdb_int32 handle, rtdb_int32 id, const char* dest_table_name)
func RawRtdbbMovePointByIdWarp(handle ConnectHandle, pointID PointID, tableName string) error {
	cTableName := C.CString(tableName)
	defer C.free(unsafe.Pointer(cTableName))
	err := C.rtdbb_move_point_by_id_warp(C.rtdb_int32(handle), C.rtdb_int32(pointID), cTableName)
	return RtdbError(err).GoError()
}

// RawRtdbbGetPointsPropertyWarp 批量获取标签点属性
// 备注：不实现, 统一使用最大长度
//
// * \param handle 连接句柄
// * \param count  整数，输入，表示标签点个数。
// * \param base   RTDB_POINT 结构数组，输入/输出，标签点基本属性列表，
// *                 输入时，id 字段指定需要得到属性的标签点，输出时，其它字段返回标签点属性值。
// * \param scan   RTDB_SCAN_POINT 结构数组，输出，采集标签点扩展属性列表
// * \param calc   RTDB_CALC_POINT 结构数组，输出，计算标签点扩展属性列表
// * \param errors 无符号整型数组，输出，获取标签属性的返回值列表，参考rtdb_error.h
// * \remark 用户须保证分配给 base、scan、calc、errors 的空间与 count 相符，
// *        扩展属性集 scan、calc 可为空指针，此时将不返回对应的扩展属性集。
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_points_property_warp(rtdb_int32 handle, rtdb_int32 count, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_CALC_POINT *calc, rtdb_error *errors)
// func RawRtdbbGetPointsPropertyWarp() {}

// RawRtdbbGetMaxPointsPropertyWarp 按最大长度批量获取标签点属性
//
// input:
//   - handle 连接句柄
//   - pointIDs 标签点ID列表
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_max_points_property_warp(rtdb_int32 handle, rtdb_int32 count, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_MAX_CALC_POINT *calc, rtdb_error *errors)
func RawRtdbbGetMaxPointsPropertyWarp(handle ConnectHandle, pointIDs []PointID) ([]RtdbPoint, []RtdbScan, []RtdbCalc, []error, error) {
	bases := make([]C.RTDB_POINT, len(pointIDs))
	for i, id := range pointIDs {
		bases[i].id = C.int(id)
	}
	scans := make([]C.RTDB_SCAN_POINT, len(pointIDs))
	calcs := make([]C.RTDB_MAX_CALC_POINT, len(pointIDs))
	errors := make([]C.rtdb_error, len(pointIDs))
	err := C.rtdbb_get_max_points_property_warp(C.rtdb_int32(handle), C.rtdb_int32(len(pointIDs)), &bases[0], &scans[0], &calcs[0], &errors[0])
	rtnBases := make([]RtdbPoint, 0)
	rtnScans := make([]RtdbScan, 0)
	rtnCalcs := make([]RtdbCalc, 0)
	rtnError := make([]error, 0)
	for i := 0; i < len(pointIDs); i++ {
		rtnBases = append(rtnBases, *cToRtdbPoint(&bases[i]))
		rtnScans = append(rtnScans, *cToRtdbScan(&scans[i]))
		rtnCalcs = append(rtnCalcs, *cToRtdbCalc(&calcs[i]))
		rtnError = append(rtnError, RtdbError(errors[i]).GoError())
	}
	return rtnBases, rtnScans, rtnCalcs, rtnError, RtdbError(err).GoError()
}

// RawRtdbbSearchWarp 搜索符合条件的标签点，使用标签点名时支持通配符
// 备注： 废弃，由于有数量限制，采用下面的分段搜索 (RawRtdbbSearchInBatchesWarp)
//
// input:
//   - handle 连接句柄
//   - tagMask 标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//   - tableMask 标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//   - source 数据源集合，字符串中的每个字符均表示一个数据源，空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
//   - unit 标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
//   - desc 标签点描述的子集，描述中包含该参数的标签点均满足条件，空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
//   - instrument 标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
//   - mode 搜索排序规则
//   - 注意：多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
//
// output:
//   - []PointID 返回标签点ID列表
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_search_warp(rtdb_int32 handle, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
// func RawRtdbbSearchWarp(handle ConnectHandle, tagMask, tableMask, source, unit, desc, instrument string, model RtdbSortFlag) ([]PointID, error) {
// 	if strings.TrimSpace(tagMask) == "" {
// 		tagMask = "*"
// 	}
// 	if strings.TrimSpace(tableMask) == "" {
// 		tableMask = "*"
// 	}
// 	cTagMask := C.CString(tagMask)
// 	defer C.free(unsafe.Pointer(cTagMask))
// 	cTableMask := C.CString(tableMask)
// 	defer C.free(unsafe.Pointer(cTableMask))
// 	cSource := C.CString(source)
// 	defer C.free(unsafe.Pointer(cSource))
// 	cUnit := C.CString(unit)
// 	defer C.free(unsafe.Pointer(cUnit))
// 	cDesc := C.CString(desc)
// 	defer C.free(unsafe.Pointer(cDesc))
// 	cInstrument := C.CString(instrument)
// 	defer C.free(unsafe.Pointer(cInstrument))
// 	cModel := C.rtdb_int32(model)
// 	count := C.rtdb_int32(1024)
// 	ids := make([]PointID, count)
// 	cIds := (*C.rtdb_int32)(unsafe.Pointer(&ids[0]))
//
// 	err := C.rtdbb_search_warp(C.rtdb_int32(handle), cTagMask, cTableMask, cSource, cUnit, cDesc, cInstrument, cModel, cIds, &count)
// 	return ids[:count], RtdbError(err).GoError()
// }

// RawRtdbbSearchInBatchesWarp 分批继续搜索符合条件的标签点，使用标签点名时支持通配符
//
// input:
//   - handle 连接句柄
//   - start 搜索起始位置。
//   - tagMask 标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//   - tableMask 标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//   - source 数据源集合，字符串中的每个字符均表示一个数据源，空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
//   - unit 标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
//   - desc 标签点描述的子集，描述中包含该参数的标签点均满足条件，空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
//   - instrument 标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
//   - mode 搜索后的排序模式
//   - 备注1：当搜索到的标签点数比提供的要小时，表示这是最后一批符合条件的标签点 (即全部搜索完毕)
//   - 备注2: 多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
//
// output:
//   - []PointID 返回ID列表
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_search_in_batches_warp(rtdb_int32 handle, rtdb_int32 start, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
func RawRtdbbSearchInBatchesWarp(handle ConnectHandle, start int32, tagMask, tableMask, source, unit, desc, instrument string, model RtdbSortFlag) ([]PointID, error) {
	if strings.TrimSpace(tagMask) == "" {
		tagMask = "*"
	}
	if strings.TrimSpace(tableMask) == "" {
		tableMask = "*"
	}
	cTagMask := C.CString(tagMask)
	defer C.free(unsafe.Pointer(cTagMask))
	cTableMask := C.CString(tableMask)
	defer C.free(unsafe.Pointer(cTableMask))
	cSource := C.CString(source)
	defer C.free(unsafe.Pointer(cSource))
	cUnit := C.CString(unit)
	defer C.free(unsafe.Pointer(cUnit))
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))
	cInstrument := C.CString(instrument)
	defer C.free(unsafe.Pointer(cInstrument))
	cCount := C.rtdb_int32(1024) // 固定1024，因为是批量获取的，这里就不指定了
	ids := make([]PointID, cCount)
	cIds := (*C.rtdb_int32)(unsafe.Pointer(&ids[0]))
	cModel := C.rtdb_int32(model)

	err := C.rtdbb_search_in_batches_warp(C.rtdb_int32(handle), C.rtdb_int32(start), cTagMask, cTableMask, cSource, cUnit, cDesc, cInstrument, cModel, cIds, &cCount)
	return ids[:cCount], RtdbError(err).GoError()
}

// RawRtdbbSearchExWarp 搜索符合条件的标签点，使用标签点名时支持通配符
//
// input:
//   - handle 连接句柄
//   - maxCount 最多返回maxCount个标签点
//   - tagMask 标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//   - tableMask 标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//   - source 数据源集合，字符串中的每个字符均表示一个数据源，空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
//   - unit 标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
//   - desc 标签点描述的子集，描述中包含该参数的标签点均满足条件，空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
//   - instrument 标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
//   - typeMask 标签点类型名称。缺省设置为空，长度不得超过 RTDB_TYPE_NAME_SIZE,内置的普通数据类型可以使用 bool、uint8、datetime等字符串表示，不区分大小写，支持模糊搜索。
//   - classOfMask 标签点的类别，缺省设置为-1，表示可以是任意类型的标签点，当使用标签点类型作为搜索条件时，必须是RTDB_CLASS枚举中的一项或者多项的组合。
//   - timeUnitMask 标签点的时间戳精度，缺省设置为-1，表示可以是任意时间戳精度，当使用此时间戳精度作为搜索条件时，timeunitmask的值可以为0或1，0表示时间戳精度为秒，1表示纳秒
//   - otherTypeMask 使用其他标签点属性作为搜索条件，缺省设置为0，表示不作为搜索条件，当使用此参数作为搜索条件时，othertypemaskvalue作为对应的搜索值，此参数的取值可以参考rtdb.h文件中的RTDB_SEARCH_MASK。
//   - otherTypeMaskValue 字符串，输入参数，当使用其他标签点属性作为搜索条件时，此参数作为对应的搜索值，缺省设置为0，表示不作为搜索条件，如果othertypemask的值为0，或者RTDB_SEARCH_NULL，则此参数被忽略, 当othertypemask对应的标签点属性为数值类型时，此搜索值只支持相等判断，当othertypemask对应的标签点属性为字符串类型时，此搜索值支持模糊搜索。
//   - mode 搜索结果排序模式
//   - 备注：多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
//
// output:
//   - []PointID 标签点列表
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_search_ex_warp(rtdb_int32 handle, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, const char *typemask, rtdb_int32 classofmask, rtdb_int32 timeunitmask, rtdb_int32 othertypemask, const char *othertypemaskvalue, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
func RawRtdbbSearchExWarp(handle ConnectHandle, maxCount int32, tagMask, tableMask, source, unit, desc, instrument, typeMask string, classOfMask RtdbType, timeUnitMask RtdbPrecision, otherTypeMask RtdbSearch, otherTypeMaskValue string, model RtdbSortFlag) ([]PointID, error) {
	if strings.TrimSpace(tagMask) == "" {
		tagMask = "*"
	}
	if strings.TrimSpace(tableMask) == "" {
		tableMask = "*"
	}
	cTagMask := C.CString(tagMask)
	defer C.free(unsafe.Pointer(cTagMask))
	cTableMask := C.CString(tableMask)
	defer C.free(unsafe.Pointer(cTableMask))
	cSource := C.CString(source)
	defer C.free(unsafe.Pointer(cSource))
	cUnit := C.CString(unit)
	defer C.free(unsafe.Pointer(cUnit))
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))
	cInstrument := C.CString(instrument)
	defer C.free(unsafe.Pointer(cInstrument))
	cTypeMask := C.CString(typeMask)
	defer C.free(unsafe.Pointer(cTypeMask))
	cClassOfMask := C.rtdb_int32(classOfMask)
	cTimeUnitMask := C.rtdb_int32(timeUnitMask)
	cOtherTypeMask := C.rtdb_int32(otherTypeMask)
	cOtherTypeMaskValue := C.CString(otherTypeMaskValue)
	defer C.free(unsafe.Pointer(cOtherTypeMaskValue))
	count := C.rtdb_int32(maxCount)
	ids := make([]PointID, count)
	cIds := (*C.rtdb_int32)(unsafe.Pointer(&ids[0]))
	cModel := C.rtdb_int32(model)
	err := C.rtdbb_search_ex_warp(C.rtdb_int32(handle), cTagMask, cTableMask, cSource, cUnit, cDesc, cInstrument, cTypeMask, cClassOfMask, cTimeUnitMask, cOtherTypeMask, cOtherTypeMaskValue, cModel, cIds, &count)
	return ids[:count], RtdbError(err).GoError()
}

// RawRtdbbSearchPointsCountWarp 搜索符合条件的标签点，获取标签点数，使用标签点名时支持通配符
//
// input:
//   - handle 连接句柄
//   - tagMask 标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//   - tableMask 标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//   - source 数据源集合，字符串中的每个字符均表示一个数据源，空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
//   - unit 标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
//   - desc 标签点描述的子集，描述中包含该参数的标签点均满足条件，空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
//   - instrument 标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
//   - typeMask 标签点类型名称。缺省设置为空，长度不得超过 RTDB_TYPE_NAME_SIZE,内置的普通数据类型可以使用 bool、uint8、datetime等字符串表示，不区分大小写，支持模糊搜索。
//   - classOfMask 标签点的类别，缺省设置为-1，表示可以是任意类型的标签点，当使用标签点类型作为搜索条件时，必须是RTDB_CLASS枚举中的一项或者多项的组合。
//   - timeUnitMask 标签点的时间戳精度，缺省设置为-1，表示可以是任意时间戳精度，当使用此时间戳精度作为搜索条件时，timeunitmask的值可以为0或1，0表示时间戳精度为秒，1表示纳秒
//   - otherTypeMask 使用其他标签点属性作为搜索条件，缺省设置为0，表示不作为搜索条件，当使用此参数作为搜索条件时，othertypemaskvalue作为对应的搜索值，此参数的取值可以参考rtdb.h文件中的RTDB_SEARCH_MASK。
//   - otherTypeMaskValue 字符串，输入参数，当使用其他标签点属性作为搜索条件时，此参数作为对应的搜索值，缺省设置为0，表示不作为搜索条件，如果othertypemask的值为0，或者RTDB_SEARCH_NULL，则此参数被忽略, 当othertypemask对应的标签点属性为数值类型时，此搜索值只支持相等判断，当othertypemask对应的标签点属性为字符串类型时，此搜索值支持模糊搜索。
//   - mode 搜索结果排序模式
//   - 备注：多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
//
// output:
//   - []PointID 标签点列表
//
// raw_fn:
// rtdb_error RTDBAPI_CALLRULE rtdbb_search_points_count_warp(rtdb_int32 handle, const char *tagmask, const char *tablemask, const char *source, const char *unit, const char *desc, const char *instrument, const char *typemask, rtdb_int32 classofmask, rtdb_int32 timeunitmask, rtdb_int32 othertypemask, const char *othertypemaskvalue, rtdb_int32 *count)
func RawRtdbbSearchPointsCountWarp(handle ConnectHandle, tagMask, tableMask, source, unit, desc, instrument, typeMask string, classOfMask RtdbType, timeUnitMask RtdbPrecision, otherTypeMask RtdbSearch, otherTypeMaskValue string) (int32, error) {
	if strings.TrimSpace(tagMask) == "" {
		tagMask = "*"
	}
	if strings.TrimSpace(tableMask) == "" {
		tableMask = "*"
	}
	cTagMask := C.CString(tagMask)
	defer C.free(unsafe.Pointer(cTagMask))
	cTableMask := C.CString(tableMask)
	defer C.free(unsafe.Pointer(cTableMask))
	cSource := C.CString(source)
	defer C.free(unsafe.Pointer(cSource))
	cUnit := C.CString(unit)
	defer C.free(unsafe.Pointer(cUnit))
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))
	cInstrument := C.CString(instrument)
	defer C.free(unsafe.Pointer(cInstrument))
	cTypeMask := C.CString(typeMask)
	defer C.free(unsafe.Pointer(cTypeMask))
	cClassOfMask := C.rtdb_int32(classOfMask)
	cTimeUnitMask := C.rtdb_int32(timeUnitMask)
	cOtherTypeMask := C.rtdb_int32(otherTypeMask)
	cOtherTypeMaskValue := C.CString(otherTypeMaskValue)
	defer C.free(unsafe.Pointer(cOtherTypeMaskValue))
	count := C.rtdb_int32(1024)
	err := C.rtdbb_search_points_count_warp(C.rtdb_int32(handle), cTagMask, cTableMask, cSource, cUnit, cDesc, cInstrument, cTypeMask, cClassOfMask, cTimeUnitMask, cOtherTypeMask, cOtherTypeMaskValue, &count)
	return int32(count), RtdbError(err).GoError()
}

// RawRtdbbUpdatePointPropertyWarp 更新单个标签点属性
// 备注： 废弃， 统一使用最大长度Calc
//
// input:
//   - handle 连接句柄
//   - base 基本标签点属性集。
//   - scan 采集标签点扩展属性集。
//   - calc 计算标签点扩展属性集。
//   - 备注：标签点由 base 参数的 id 字段指定，其中 id、table、type、millisecond 字段不能修改，changedate、changer、createdate、creator 字段由系统维护，
//   - 备注：其余字段均可修改，包括 classof 字段。输入参数中 scan、calc 可为空指针，对应的扩展属性将保持不变。
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_update_point_property_warp(rtdb_int32 handle, const RTDB_POINT *base, const RTDB_SCAN_POINT *scan, const RTDB_CALC_POINT *calc)
// func RawRtdbbUpdatePointPropertyWarp(handle ConnectHandle, base *RtdbPoint, scan *RtdbScan, calc *RtdbCalc) error {
// 	cBase := goToCRtdbPoint(base)
// 	cScan := goToCRtdbScan(scan)
// 	cCalc := goToCRtdbCalc(calc)
// 	err := C.rtdbb_update_point_property_warp(C.rtdb_int32(handle), cBase, cScan, cCalc)
// 	return RtdbError(err).GoError()
// }

// RawRtdbbUpdateMaxPointPropertyWarp 按最大长度更新单个标签点属性
//
// input:
//   - handle 连接句柄
//   - base 基本标签点属性集。
//   - scan 采集标签点扩展属性集。
//   - calc 计算标签点扩展属性集。
//   - 备注：标签点由 base 参数的 id 字段指定，其中 id、table、type、millisecond 字段不能修改，changedate、changer、createdate、creator 字段由系统维护，
//   - 备注：其余字段均可修改，包括 classof 字段。输入参数中 scan、calc 可为空指针，对应的扩展属性将保持不变。
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_update_max_point_property_warp(rtdb_int32 handle, const RTDB_POINT *base, const RTDB_SCAN_POINT *scan, const RTDB_MAX_CALC_POINT *calc)
func RawRtdbbUpdateMaxPointPropertyWarp(handle ConnectHandle, base *RtdbPoint, scan *RtdbScan, calc *RtdbCalc) error {
	cBase := goToCRtdbPoint(base)
	cScan := goToCRtdbScan(scan)
	cCalc := goToCRtdbCalc(calc)
	err := C.rtdbb_update_max_point_property_warp(C.rtdb_int32(handle), cBase, cScan, cCalc)
	return RtdbError(err).GoError()
}

// RawRtdbbFindPointsWarp 根据 "表名.标签点名" 格式批量获取标签点标识
// 备注：废弃，此函数已有扩展版，因此不实现， 参见RawRtdbbFindPointsExWarp
//
// input:
//   - handle 连接句柄
//   - count 输入时表示标签点个数 (即table_dot_tags、ids、types、classof、use_ms 的长度)，输出时表示找到的标签点个数
//   - table_dot_tags 输入，"表名.标签点名" 列表
//
// output:
//   - []PointID 标签点ID列表
//   - []RtdbType 标签点数值类型列表
//   - []RtdbClass 标签点种类列表
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_find_points_warp(rtdb_int32 handle, rtdb_int32 *count, const char* const* table_dot_tags, rtdb_int32 *ids, rtdb_int32 *types, rtdb_int32 *classof, rtdb_int16 *use_ms)
// func RawRtdbbFindPointsWarp(handle ConnectHandle, count int32, tableDotTags []string) ([]PointID, []RtdbType, []RtdbClass, error) {
// 	cCount := C.rtdb_int32(count)
// 	cTags := make([]*C.char, 0)
// 	for _, tag := range tableDotTags {
// 		cT := C.CString(tag)
// 		cTags = append(cTags, cT)
// 	}
// 	defer func() {
// 		for _, cT := range cTags {
// 			C.free(unsafe.Pointer(cT))
// 		}
// 	}()
// 	ccTags := &cTags[0]
// 	ids := make([]PointID, cCount)
// 	cIds := (*C.rtdb_int32)(unsafe.Pointer(&ids[0]))
// 	types := make([]RtdbType, cCount)
// 	cTypes := (*C.rtdb_int32)(unsafe.Pointer(&types[0]))
// 	classOfs := make([]RtdbClass, cCount)
// 	cClassOf := (*C.rtdb_int32)(unsafe.Pointer(&classofs[0]))
// 	useMs := make([]C.rtdb_int16, cCount)
// 	cUseMs := &useMs[0]
// 	err := C.rtdbb_find_points_warp(C.rtdb_int32(handle), &cCount, ccTags, cIds, cTypes, cClassOf, cUseMs)
// 	return ids[:cCount], types[:cCount], classOfs[:cCount], RtdbError(err).GoError()
// }

// RawRtdbbFindPointsExWarp 根据 "表名.标签点名" 格式批量获取标签点标识
//
// input:
//   - handle 连接句柄
//   - count 输入时表示标签点个数
//   - table_dot_tags 输入，"表名.标签点名" 列表
//
// output:
//   - []PointID 标签点标识列表, 返回 0 表示未找到
//   - []RtdbType 标签点数据类型
//   - []RtdbClass 标签点类别
//   - []PrecisionType 时间戳精度，
//   - []RtdbError 表示每个标签点的查询结果的错误码
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_find_points_ex_warp(rtdb_int32 handle, rtdb_int32* count, const char* const* table_dot_tags, rtdb_int32* ids, rtdb_int32* types, rtdb_int32* classof, rtdb_precision_type* precisions, rtdb_error* errors)
func RawRtdbbFindPointsExWarp(handle ConnectHandle, tableDotTags []string) ([]PointID, []RtdbType, []RtdbClass, []RtdbPrecision, []RtdbError, error) {
	if len(tableDotTags) == 0 {
		return nil, nil, nil, nil, nil, nil
	}

	count := len(tableDotTags)
	cCount := C.rtdb_int32(count)
	cHandle := C.rtdb_int32(handle)
	tableDotTagsList := make([]*C.char, count)
	for i := 0; i < count; i++ {
		tableDotTagsList[i] = C.CString(tableDotTags[i])
	}
	defer func() {
		for i := 0; i < count; i++ {
			C.free(unsafe.Pointer(tableDotTagsList[i]))
		}
	}()
	cTableDotTags := &tableDotTagsList[0]
	ids := make([]PointID, count)
	cIds := (*C.rtdb_int32)(unsafe.Pointer(&ids[0]))
	types := make([]RtdbType, count)
	cTypes := (*C.rtdb_int32)(unsafe.Pointer(&types[0]))
	classof := make([]RtdbClass, count)
	cClassof := (*C.rtdb_int32)(unsafe.Pointer(&classof[0]))
	precisions := make([]RtdbPrecision, count)
	cPrecisions := (*C.rtdb_precision_type)(unsafe.Pointer(&precisions[0]))
	errs := make([]RtdbError, count)
	cErrs := (*C.rtdb_error)(unsafe.Pointer(&errs[0]))

	err := C.rtdbb_find_points_ex_warp(cHandle, &cCount, cTableDotTags, cIds, cTypes, cClassof, cPrecisions, cErrs)
	return ids[:cCount], types[:cCount], classof[:cCount], precisions[:cCount], errs[:cCount], RtdbError(err).GoError()
}

// RawRtdbbSortPointsWarp 根据标签属性字段对标签点标识进行排序
//
// input:
//   - handle 连接句柄
//   - ids 标签点标识列表
//   - index 属性字段枚举，参见 RTDB_TAG_FIELD_INDEX，将根据该字段对 ID 进行排序。
//   - flag 标志位组合，参见 RTDB_TAG_SORT_FLAG 枚举
//   - 备注：当对可回收标签排序时，以下字段索引不可使用：RTDB_TAG_INDEX_TIMESTAMP、RTDB_TAG_INDEX_VALUE、RTDB_TAG_INDEX_QUALITY
//
// output:
//   - []PointID 标签点ID列表
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_sort_points_warp(rtdb_int32 handle, rtdb_int32 count, rtdb_int32 *ids, rtdb_int32 index, rtdb_int32 flag)
func RawRtdbbSortPointsWarp(handle ConnectHandle, ids []PointID, index RtdbTagIndex, flag RtdbSortFlag) ([]PointID, error) {
	err := C.rtdbb_sort_points_warp(C.rtdb_int32(handle), C.rtdb_int32(len(ids)), (*C.rtdb_int32)(unsafe.Pointer(&ids[0])), C.rtdb_int32(index), C.rtdb_int32(flag))
	return ids, RtdbError(err).GoError()
}

// RawRtdbbUpdateTableNameWarp 根据表 ID 更新表名称。
//
// input:
//   - handle 连接句柄
//   - id 表ID
//   - name 新的标签点表名称。
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_update_table_name_warp(rtdb_int32 handle, rtdb_int32 tab_id, const char *name)
func RawRtdbbUpdateTableNameWarp(handle ConnectHandle, id TableID, name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	err := C.rtdbb_update_table_name_warp(C.rtdb_int32(handle), C.rtdb_int32(id), cName)
	return RtdbError(err).GoError()
}

// RawRtdbbUpdateTableDescByIdWarp 根据表 ID 更新表描述。
//
// input:
//   - handle 连接句柄
//   - id 表ID
//   - desc 新的标签点表描述。
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_update_table_desc_by_id_warp(rtdb_int32 handle, rtdb_int32 tab_id, const char *desc)
func RawRtdbbUpdateTableDescByIdWarp(handle ConnectHandle, id TableID, desc string) error {
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))
	err := C.rtdbb_update_table_desc_by_id_warp(C.rtdb_int32(handle), C.rtdb_int32(id), cDesc)
	return RtdbError(err).GoError()
}

// RawRtdbbUpdateTableDescByNameWarp 根据表名称更新表描述。
//
// input:
//   - handle 连接句柄
//   - name 要修改表的名称。
//   - desc 新的表描述。
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_update_table_desc_by_name_warp(rtdb_int32 handle, const char *name, const char *desc)
func RawRtdbbUpdateTableDescByNameWarp(handle ConnectHandle, name string, desc string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))
	err := C.rtdbb_update_table_desc_by_name_warp(C.rtdb_int32(handle), cName, cDesc)
	return RtdbError(err).GoError()
}

// RawRtdbbRecoverPointWarp 恢复已删除标签点
//
// input:
//   - handle 连接句柄
//   - tableID 要将标签点恢复到的表标识
//   - point_id 待恢复的标签点标识
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_recover_point_warp(rtdb_int32 handle, rtdb_int32 table_id, rtdb_int32 point_id)
func RawRtdbbRecoverPointWarp(handle ConnectHandle, tableID TableID, pointID PointID) error {
	err := C.rtdbb_recover_point_warp(C.rtdb_int32(handle), C.rtdb_int32(tableID), C.rtdb_int32(pointID))
	return RtdbError(err).GoError()
}

// RawRtdbbPurgePointWarp 清除标签点
//
// input:
//   - handle 连接句柄
//   - id 要清除的标签点标识
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_purge_point_warp(rtdb_int32 handle, rtdb_int32 id)
func RawRtdbbPurgePointWarp(handle ConnectHandle, id PointID) error {
	err := C.rtdbb_purge_point_warp(C.rtdb_int32(handle), C.rtdb_int32(id))
	return RtdbError(err).GoError()
}

// RawRtdbbGetRecycledPointsCountWarp 获取可回收标签点数量
//
// input:
//   - handle 连接句柄
//   - count 可回收标签点的数量
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_points_count_warp(rtdb_int32 handle, rtdb_int32 *count)
func RawRtdbbGetRecycledPointsCountWarp(handle ConnectHandle) (int32, error) {
	count := C.rtdb_int32(0)
	err := C.rtdbb_get_recycled_points_count_warp(C.rtdb_int32(handle), &count)
	return int32(count), RtdbError(err).GoError()
}

// RawRtdbbGetRecycledPointsWarp 获取可回收标签点 id 列表
//
// input:
//   - handle 连接句柄
//   - count 最大返回标签点个数
//
// output:
//   - []RtdbPoint 标签点列表
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_points_warp(rtdb_int32 handle, rtdb_int32 *ids, rtdb_int32 *count)
func RawRtdbbGetRecycledPointsWarp(handle ConnectHandle, count int32) ([]RtdbPoint, error) {
	if count == 0 {
		return nil, nil
	}

	cCount := C.rtdb_int32(count)
	points := make([]RtdbPoint, cCount)
	err := C.rtdbb_get_recycled_points_warp(C.rtdb_int32(handle), (*C.rtdb_int32)(unsafe.Pointer(&points[0])), &cCount)

	return points[:cCount], RtdbError(err).GoError()
}

// RawRtdbbSearchRecycledPointsWarp 搜索符合条件的可回收标签点，使用标签点名时支持通配符
// 备注：不实现，统一使用分批搜索
//
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
// func RawRtdbbSearchRecycledPointsWarp() {}

// RawRtdbbSearchRecycledPointsInBatchesWarp 分批搜索符合条件的可回收标签点，使用标签点名时支持通配符
//
// input:
//   - handle 连接句柄
//   - start 搜索的起始位置。
//   - tagMask 标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
//   - tableMask 标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
//   - source 数据源集合，字符串中的每个字符均表示一个数据源，空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
//   - unit 标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
//   - desc 标签点描述的子集，描述中包含该参数的标签点均满足条件，空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
//   - instrument 标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
//   - mode 搜索结果排序模式
//
// output:
//   - []PointID 标签点ID列表
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_search_recycled_points_in_batches_warp(rtdb_int32 handle, rtdb_int32 start, const char *tagmask, const char *fullmask, const char *source, const char *unit, const char *desc, const char *instrument, rtdb_int32 mode, rtdb_int32 *ids, rtdb_int32 *count)
func RawRtdbbSearchRecycledPointsInBatchesWarp(handle ConnectHandle, start int32, tagMask, fullMask, source, unit, desc, instrument string, mode RtdbSortFlag) ([]PointID, error) {
	cTagMask := C.CString(tagMask)
	defer C.free(unsafe.Pointer(cTagMask))
	cFullMask := C.CString(fullMask)
	defer C.free(unsafe.Pointer(cFullMask))
	cSource := C.CString(source)
	defer C.free(unsafe.Pointer(cSource))
	cUnit := C.CString(unit)
	defer C.free(unsafe.Pointer(cUnit))
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))
	cInstrument := C.CString(instrument)
	defer C.free(unsafe.Pointer(cInstrument))
	ids := make([]PointID, 1024)
	cIds := (*C.rtdb_int32)(unsafe.Pointer(&ids[0]))
	cCount := C.rtdb_int32(len(ids))
	err := C.rtdbb_search_recycled_points_in_batches_warp(C.rtdb_int32(handle), C.rtdb_int32(start), cTagMask, cFullMask, cSource, cUnit, cDesc, cInstrument, C.rtdb_int32(mode), cIds, &cCount)
	return ids[:cCount], RtdbError(err).GoError()
}

// RawRtdbbGetRecycledPointPropertyWarp 获取可回收标签点的属性
// 备注： 不实现， 统一使用最大长度， RawRtdbbGetRecycledMaxPointPropertyWarp
//
// * \param handle   连接句柄
// * \param base     RTDB_POINT 结构，输入/输出，标签点基本属性。
// 输入时，由 id 字段指定要取得的可回收标签点。
// * \param scan     RTDB_SCAN_POINT 结构，输出，标签点采集扩展属性
// * \param calc     RTDB_CALC_POINT 结构，输出，标签点计算扩展属性
// * \remark scan、calc 可为空指针，对应的扩展信息将不返回。
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_point_property_warp(rtdb_int32 handle, RTDB_POINT *base, RTDB_SCAN_POINT *scan, RTDB_CALC_POINT *calc)
// func RawRtdbbGetRecycledPointPropertyWarp() {}

// RawRtdbbGetRecycledMaxPointPropertyWarp 按最大长度获取可回收标签点的属性
//
// input:
//   - handle 连接句柄
//   - id 标签点ID
//
// output:
//   - RtdbPoint 标签点基本属性
//   - RtdbScan 标签点采集属性
//   - RtdbCalc 标签点计算属性
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_max_point_property_warp(rtdb_int32 handle, RTDB_POINT* base, RTDB_SCAN_POINT* scan, RTDB_MAX_CALC_POINT* calc)
func RawRtdbbGetRecycledMaxPointPropertyWarp(handle ConnectHandle, id PointID) (*RtdbPoint, *RtdbScan, *RtdbCalc, error) {
	base := C.RTDB_POINT{}
	base.id = C.rtdb_int32(id)
	scan := C.RTDB_SCAN_POINT{}
	calc := C.RTDB_MAX_CALC_POINT{}
	err := C.rtdbb_get_recycled_max_point_property_warp(C.rtdb_int32(handle), &base, &scan, &calc)
	return cToRtdbPoint(&base), cToRtdbScan(&scan), cToRtdbCalc(&calc), RtdbError(err).GoError()
}

// RawRtdbbClearRecyclerWarp 清空标签点回收站
//
// input:
//   - handle 连接句柄
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_clear_recycler_warp(rtdb_int32 handle)
func RawRtdbbClearRecyclerWarp(handle ConnectHandle) error {
	err := C.rtdbb_clear_recycler_warp(C.rtdb_int32(handle))
	return RtdbError(err).GoError()
}

// RawRtdbbSubscribeTagsExWarp 标签点属性更改通知订阅
//
// input:
//   - handle 连接句柄
//   - options 订阅选项，参见枚举RTDB_OPTION
//   - param 用户参数，作为rtdbb_tags_change_ex的param参数
//   - callback rtdbb_tags_change_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
//   - 备注：当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
//   - 备注：作为event_type取值调用回掉函数后退出订阅。
//   - 备注：当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
//   - 备注：作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
//   - 备注：网络中断期间回掉函数调用频率为最少3秒
//   - 注意!!：用于订阅测点的连接句柄必需是独立的，不能再用来调用其它 api，
//
// callback:
//   - event_type 参数值含义如下：RTDB_E_DATA(标签点属性发生更改) RTDB_E_DISCONNECT(订阅客户端与数据库网络断开) RTDB_E_RECOVERY(订阅客户端与数据库网络及订阅恢复)
//   - handle 产生订阅回掉的连接句柄，调用rtdbb_subscribe_tags_ex时的handle参数
//   - param 用户自定义参数，调用rtdbb_subscribe_tags_ex时的param参数，由用户传入
//   - count event_type为RTDB_E_DATA时表示ids的数量, event_type为其它值时，count值为0
//   - ids event_type为RTDB_E_DATA时表示属性更改的标签点ID，数量由count指定, event_type为其它值时，ids值为NULL
//   - what event_type为RTDB_E_DATA时表示属性变更原因，参考RTDB_TAG_CHANGE_REASON, event_type为其它值时，what时值为0
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_subscribe_tags_ex_warp(rtdb_int32 handle, rtdb_uint32 options, void* param, rtdbb_tags_change_event_ex callback)
func RawRtdbbSubscribeTagsExWarp(handle ConnectHandle, options RtdbSubscribeOption, param unsafe.Pointer) error {
	err := C.rtdbb_subscribe_tags_ex_warp(C.rtdb_int32(handle), C.rtdb_uint32(options), param, (C.rtdbb_tags_change_event_ex)(unsafe.Pointer(C.goSubscribeTagsEx)))
	return RtdbError(err)
}

// RawRtdbbCancelSubscribeTagsWarp 取消标签点属性更改通知订阅
//
// input:
//   - handle 连接句柄
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_cancel_subscribe_tags_warp(rtdb_int32 handle)
func RawRtdbbCancelSubscribeTagsWarp(handle ConnectHandle) error {
	err := C.rtdbb_cancel_subscribe_tags_warp(C.rtdb_int32(handle))
	return RtdbError(err).GoError()
}

// RawRtdbbCreateNamedTypeWarp 创建自定义类型
//
// input:
//   - handle 连接句柄
//   - name 自定义类型的名称，类型的唯一标示,不能重复，长度不能超过RTDB_TYPE_NAME_SIZE，输入参数
//   - fields 自定义类型中包含的字段的属性，RTDB_DATA_TYPE_FIELD结构的数组，个数与field_count相等，RTDB_DATA_TYPE_FIELD中的length只对type为str或blob类型的数据有效。其他类型忽略
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_create_named_type_warp(rtdb_int32 handle, const char* name, rtdb_int32 field_count, const RTDB_DATA_TYPE_FIELD* fields, char desc[RTDB_DESC_SIZE])
func RawRtdbbCreateNamedTypeWarp(handle ConnectHandle, name string, fields []RtdbDataTypeField, desc string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))
	cFields := make([]C.RTDB_DATA_TYPE_FIELD, 0)
	for _, field := range fields {
		cFields = append(cFields, *goToCRtdbDataTypeField(&field))
	}
	err := C.rtdbb_create_named_type_warp(C.rtdb_int32(handle), cName, C.rtdb_int32(len(fields)), &cFields[0], cDesc)
	return RtdbError(err).GoError()
}

// RawRtdbbGetNamedTypesCountWarp 获取所有的自定义类型的总数
//
// input:
//   - handle 连接句柄
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_named_types_count_warp(rtdb_int32 handle, rtdb_int32* count)
func RawRtdbbGetNamedTypesCountWarp(handle ConnectHandle) (int32, error) {
	count := C.rtdb_int32(0)
	err := C.rtdbb_get_named_types_count_warp(C.rtdb_int32(handle), &count)
	return int32(count), RtdbError(err).GoError()
}

// RawRtdbbGetAllNamedTypesWarp 获取所有的自定义类型
//
// input:
//   - handle 连接句柄
//   - count 返回所有的自定义类型的总数
//
// output:
//   - []string 自定义类型名称数组
//   - []int32 自定义类型中子字段个数数组
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_all_named_types_warp(rtdb_int32 handle, rtdb_int32* count, char* name[RTDB_TYPE_NAME_SIZE], rtdb_int32* field_counts)
func RawRtdbbGetAllNamedTypesWarp(handle ConnectHandle, count int32) ([]string, []int32, error) {
	cgoHandle := C.rtdb_int32(handle)
	cgoCount := C.rtdb_int32(count)
	names := make([]*C.char, count)
	for i := 0; i < int(count); i++ {
		names[i] = (*C.char)(C.CBytes(make([]byte, RtdbConstTypeNameSize)))
	}
	defer func() {
		for i := 0; i < int(count); i++ {
			C.free(unsafe.Pointer(names[i]))
		}
	}()
	cgoNames := (**C.char)(unsafe.Pointer(&names[0]))
	counts := make([]int32, count)
	cgoCounts := (*C.rtdb_int32)(unsafe.Pointer(&counts[0]))
	err := C.rtdbb_get_all_named_types_warp(cgoHandle, &cgoCount, cgoNames, cgoCounts)
	goNames := make([]string, 0)
	for i := 0; i < int(count); i++ {
		goNames = append(goNames, C.GoString(names[i]))
	}
	return goNames, counts[:cgoCount], RtdbError(err)
}

// RawRtdbbGetNamedTypeWarp 获取自定义类型的所有字段
//
// input:
//   - handle 连接句柄
//   - name 自定义类型的名称
//   - field_count 返回name指定的自定义类型的字段个数，输入/输出参数
//
// output:
//   - []RtdbDataTypeField 自定义类型内部字段数组
//   - int32 自定义类型长度
//   - string 自定义类型描述
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_named_type_warp(rtdb_int32 handle, const char* name, rtdb_int32* field_count, RTDB_DATA_TYPE_FIELD* fields, rtdb_int32* type_size, char desc[RTDB_DESC_SIZE])
func RawRtdbbGetNamedTypeWarp(handle ConnectHandle, name string, fieldCount int32) ([]RtdbDataTypeField, int32, string, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cFieldCount := C.rtdb_int32(fieldCount)
	fields := make([]C.RTDB_DATA_TYPE_FIELD, cFieldCount)
	typeSize := C.rtdb_int32(0)
	desc := make([]C.char, int(C.RTDB_DESC_SIZE))
	err := C.rtdbb_get_named_type_warp(C.rtdb_int32(handle), cName, &cFieldCount, &fields[0], &typeSize, &desc[0])
	rtn := make([]RtdbDataTypeField, 0)
	for _, field := range fields[:cFieldCount] {
		rtn = append(rtn, *cToRtdbDataTypeField(&field))
	}
	return rtn, int32(typeSize), C.GoString(&desc[0]), RtdbError(err).GoError()
}

// RawRtdbbRemoveNamedTypeWarp 删除自定义类型
//
// input:
//   - handle 连接句柄
//   - name 自定义类型的名称
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_remove_named_type_warp(rtdb_int32 handle, const char* name, rtdb_int32 reserved GAPI_DEFAULT_VALUE(0))
func RawRtdbbRemoveNamedTypeWarp(handle int32, name string) error {
	cgoHandle := C.rtdb_int32(handle)
	cgoName := C.CString(name)
	defer C.free(unsafe.Pointer(cgoName))
	err := C.rtdbb_remove_named_type_warp(cgoHandle, cgoName, 0)
	return RtdbError(err).GoError()
}

// RawRtdbbGetNamedTypeNamesPropertyWarp 根据标签点id查询标签点所对应的自定义类型的名字和字段总数
//
// input:
//   - handle 连接句柄
//   - ids 标签点ID数组
//
// output:
//   - []string 自定义类型名称数组
//   - []int32 自定义类型中字段数量数组
//   - []RtdbError 错误数组
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_named_type_names_property_warp(rtdb_int32 handle, rtdb_int32 *count, rtdb_int32 *ids, char* const *named_type_names, rtdb_int32 *field_counts, rtdb_error *errors)
func RawRtdbbGetNamedTypeNamesPropertyWarp(handle ConnectHandle, ids []PointID) ([]string, []int32, []error, error) {
	cgoHandle := C.rtdb_int32(handle)
	count := len(ids)
	cgoCount := C.rtdb_int32(count)
	cgoIds := (*C.rtdb_int32)(unsafe.Pointer(&ids[0]))
	namedTypeNames := make([]*C.char, count)
	for i := 0; i < int(count); i++ {
		namedTypeNames[i] = (*C.char)(C.CBytes(make([]byte, 4096)))
	}
	defer func() {
		for i := 0; i < int(count); i++ {
			C.free(unsafe.Pointer(namedTypeNames[i]))
		}
	}()
	cgoNamedTypeNames := (**C.char)(unsafe.Pointer(&namedTypeNames[0]))
	counts := make([]int32, count)
	cgoCounts := (*C.rtdb_int32)(unsafe.Pointer(&counts[0]))
	errors := make([]RtdbError, count)
	cgoErrors := (*C.rtdb_error)(unsafe.Pointer(&errors[0]))
	err := C.rtdbb_get_named_type_names_property_warp(cgoHandle, &cgoCount, cgoIds, cgoNamedTypeNames, cgoCounts, cgoErrors)
	names := make([]string, 0)
	for i := 0; i < int(cgoCount); i++ {
		names = append(names, C.GoString(namedTypeNames[i]))
	}
	errs := make([]error, 0)
	for _, err := range errors[:cgoCount] {
		errs = append(errs, err.GoError())
	}
	return names[:cgoCount], counts[:cgoCount], errs[:cgoCount], RtdbError(err)
}

// RawRtdbbGetRecycledNamedTypeNamesPropertyWarp 根据回收站标签点id查询标签点所对应的自定义类型的名字和字段总数
//
// input:
//   - handle 连接句柄
//   - ids 标签点ID数组
//
// output:
//   - []string 自定义类型名称数组
//   - []int32 自定义类型中字段数量数组
//   - []RtdbError 错误数组
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_recycled_named_type_names_property_warp(rtdb_int32 handle, rtdb_int32 *count, rtdb_int32 *ids, char* const *named_type_names, rtdb_int32 *field_counts, rtdb_error *errors)
func RawRtdbbGetRecycledNamedTypeNamesPropertyWarp(handle ConnectHandle, ids []PointID) ([]string, []int32, []error, error) {
	cgoHandle := C.rtdb_int32(handle)
	count := len(ids)
	cgoCount := C.rtdb_int32(count)
	cgoIds := (*C.rtdb_int32)(unsafe.Pointer(&ids[0]))
	namedTypeNames := make([]*C.char, count)
	for i := 0; i < int(count); i++ {
		namedTypeNames[i] = (*C.char)(C.CBytes(make([]byte, 4096)))
	}
	defer func() {
		for i := 0; i < int(count); i++ {
			C.free(unsafe.Pointer(namedTypeNames[i]))
		}
	}()
	cgoNamedTypeNames := (**C.char)(unsafe.Pointer(&namedTypeNames[0]))
	fieldCounts := make([]int32, count)
	cgoCounts := (*C.rtdb_int32)(unsafe.Pointer(&fieldCounts[0]))
	errors := make([]RtdbError, count)
	cgoErrors := (*C.rtdb_error)(unsafe.Pointer(&errors[0]))
	err := C.rtdbb_get_recycled_named_type_names_property_warp(cgoHandle, &cgoCount, cgoIds, cgoNamedTypeNames, cgoCounts, cgoErrors)
	names := make([]string, 0)
	for i := 0; i < int(cgoCount); i++ {
		names = append(names, C.GoString(namedTypeNames[i]))
	}
	errs := make([]error, 0)
	for _, err := range errors[:cgoCount] {
		errs = append(errs, err.GoError())
	}
	return names[:cgoCount], fieldCounts[:cgoCount], errs[:cgoCount], RtdbError(err)
}

// RawRtdbbGetNamedTypePointsCountWarp 获取该自定义类型的所有标签点个数
//
// input:
//   - handle 连接句柄
//   - name 自定义类型的名称
//
// output:
//   - int32 自定义类型标签点个数
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_named_type_points_count_warp(rtdb_int32 handle, const char* name, rtdb_int32 *points_count)
func RawRtdbbGetNamedTypePointsCountWarp(handle ConnectHandle, name string) (int32, error) {
	cgoHandle := C.rtdb_int32(handle)
	cgoName := C.CString(name)
	defer C.free(unsafe.Pointer(cgoName))
	cgoCount := C.rtdb_int32(0)
	err := C.rtdbb_get_named_type_points_count_warp(cgoHandle, cgoName, &cgoCount)
	return int32(cgoCount), RtdbError(err).GoError()
}

// RawRtdbbGetBaseTypePointsCountWarp 获取该内置的基本类型的所有标签点个数
//
// input:
//   - handle 连接句柄
//   - rtdbType 数值类型
//
// output:
//   - int32 自定义类型标签点个数
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_get_base_type_points_count_warp(rtdb_int32 handle, rtdb_int32 type, rtdb_int32 *points_count)
func RawRtdbbGetBaseTypePointsCountWarp(handle ConnectHandle, rtdbType RtdbType) (int32, error) {
	cgoHandle := C.rtdb_int32(handle)
	cgoType := C.rtdb_int32(rtdbType)
	cgoCount := C.rtdb_int32(0)
	err := C.rtdbb_get_base_type_points_count_warp(cgoHandle, cgoType, &cgoCount)
	return int32(cgoCount), RtdbError(err).GoError()
}

// RawRtdbbModifyNamedTypeWarp 修改自定义类型名称,描述,字段名称,字段描述
//
// input:
//   - handle 连接句柄
//   - name 自定义类型的名称
//   - modifyName 要修改的自定义类型名称
//   - modifyDesc 要修改的自定义类型的描述
//   - modifyFieldName 要修改的自定义类型字段的名称
//   - modifyFieldDesc 要修改的自定义类型字段的描述
//
// raw_fn:
//   - rtdb_error RTDBAPI_CALLRULE rtdbb_modify_named_type_warp(rtdb_int32 handle, const char* name, const char* modify_name, const char* modify_desc, const char* modify_field_name[RTDB_TYPE_NAME_SIZE], const char* modify_field_desc[RTDB_DESC_SIZE], rtdb_int32 field_count)
func RawRtdbbModifyNamedTypeWarp(handle ConnectHandle, name string, modifyName string, modifyDesc string, fieldNames []string, fieldDescs []string) error {
	cgoHandle := C.rtdb_int32(handle)
	cgoName := C.CString(name)
	defer C.free(unsafe.Pointer(cgoName))
	cgoModifyName := C.CString(modifyName)
	defer C.free(unsafe.Pointer(cgoModifyName))
	cgoModifyDesc := C.CString(modifyDesc)
	defer C.free(unsafe.Pointer(cgoModifyDesc))
	cgoFieldLen := C.rtdb_int32(len(fieldNames))
	names := make([]*C.char, len(fieldNames))
	for i, n := range fieldNames {
		names[i] = C.CString(n)
	}
	defer func() {
		for _, n := range names {
			C.free(unsafe.Pointer(n))
		}
	}()
	cgoNames := (**C.char)(unsafe.Pointer(&names[0]))
	descs := make([]*C.char, len(fieldDescs))
	for i, d := range fieldDescs {
		descs[i] = C.CString(d)
	}
	defer func() {
		for _, d := range descs {
			C.free(unsafe.Pointer(d))
		}
	}()
	cgoDescs := (**C.char)(unsafe.Pointer(&descs[0]))
	err := C.rtdbb_modify_named_type_warp(cgoHandle, cgoName, cgoModifyName, cgoModifyDesc, cgoNames, cgoDescs, cgoFieldLen)
	return RtdbError(err).GoError()
}

// RawRtdbbGetMetaSyncInfoWarp 获取元数据同步信息
//
// input:
//   - handle 连接句柄
//   - node_number 双活节点id，1表示第一个节点，2表示第二个节点。0表示所有节点
//
// * \param sync_infos       RTDB_SYNC_INFO数组，输出参数，输出实际获取到的同步信息
// * \param errors           rtdb_error数组，输出参数，输出对应节点的错误信息
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_meta_sync_info_warp(rtdb_int32 handle, rtdb_int32 node_number, rtdb_int32* count, RTDB_SYNC_INFO* sync_infos, rtdb_error* errors)
func RawRtdbbGetMetaSyncInfoWarp(handle ConnectHandle, nodeNumber int32) {

}

// RawRtdbsGetSnapshots64Warp 批量读取开关量、模拟量快照数值
// * \param handle    连接句柄
// * \param count     整型，输入/输出，标签点个数，
// *                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
// *                    输出时表示成功获取实时值的标签点个数
// * \param ids       整型数组，输入，标签点标识列表
// * \param datetimes 整型数组，输出，实时数值时间列表,
// *                    表示距离1970年1月1日08:00:00的秒数
// * \param ms        短整型数组，输出，实时数值时间列表，
// *                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
// * \param values    双精度浮点型数组，输出，实时浮点型数值列表，
// *                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的快照值；否则为 0
// * \param states    64 位整型数组，输出，实时整型数值列表，
// *                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
// *                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的快照值；否则为 0
// * \param qualities 短整型数组，输出，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
// * \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
// *        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
// rtdb_error RTDBAPI_CALLRULE rtdbs_get_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities, rtdb_error* errors)
func RawRtdbsGetSnapshots64Warp() {}

// RawRtdbsPutSnapshots64Warp 批量写入开关量、模拟量快照数值
// * \param handle    连接句柄
// * \param count     整型，输入/输出，标签点个数，
// *                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
// *                    输出时表示成功写入实时值的标签点个数
// * \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
// *                    但它们的时间戳必需是递增的。
// * \param datetimes 整型数组，输入，实时数值时间列表,
// *                    表示距离1970年1月1日08:00:00的秒数
// * \param ms        短整型数组，输入，实时数值时间列表，
// *                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
// * \param values    双精度浮点型数组，输入，实时浮点型数值列表，
// *                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
// * \param states    64 位整型数组，输入，实时整型数值列表，
// *                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
// *                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
// * \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
// * \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
// *        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
// rtdb_error RTDBAPI_CALLRULE rtdbs_put_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors)
func RawRtdbsPutSnapshots64Warp() {}

// RawRtdbsPutSnapshots 批量写入开关量、模拟量快照数值
// * \param handle    连接句柄
// * \param count     整型，输入/输出，标签点个数，
// *                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
// *                    输出时表示成功写入实时值的标签点个数
// * \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
// *                    但它们的时间戳必需是递增的。
// * \param datetimes 整型数组，输入，实时数值时间列表,
// *                    表示距离1970年1月1日08:00:00的秒数
// * \param ms        短整型数组，输入，实时数值时间列表，
// *                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
// * \param values    双精度浮点型数组，输入，实时浮点型数值列表，
// *                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
// * \param states    64 位整型数组，输入，实时整型数值列表，
// *                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
// *                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
// * \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
// * \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
// *        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
// *        仅当输入时间戳与当前快照时间戳完全相等时，会替换当前快照的值和质量；
// *        其余情况下会调用 rtdbs_put_snapshots()
// rtdb_error RTDBAPI_CALLRULE rtdbs_fix_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors)
func RawRtdbsPutSnapshots() {}

// RawRtdbsBackSnapshots64Warp 批量回溯快照
// * \param handle    连接句柄
// * \param count     整型，输入/输出，标签点个数，
// *                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
// *                    输出时表示成功写入实时值的标签点个数
// * \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
// *
// * \param datetimes 整型数组，输入，实时数值时间列表,
// *                    表示距离1970年1月1日08:00:00的秒数
// * \param ms        短整型数组，输入，实时数值时间列表，
// *                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
// * \param values    双精度浮点型数组，输入，实时浮点型数值列表，
// *                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
// * \param states    64 位整型数组，输入，实时整型数值列表，
// *                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
// *                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
// * \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
// * \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
// *        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
// * 功能说明：
// *       批量将标签点的快照值vtmq改成传入的vtmq，如果传入的时间戳早于当前快照，会删除传入时间戳到当前快照的历史存储值。
// *       如果传入的时间戳等于或者晚于当前快照，什么也不做。
// rtdb_error RTDBAPI_CALLRULE rtdbs_back_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors)
func RawRtdbsBackSnapshots64Warp() {}

// RawRtdbsGetCoorSnapshots64Warp 批量读取坐标实时数据
// *
// * \param handle    连接句柄
// * \param count     整型，输入/输出，标签点个数，
// *                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
// *                    输出时表示成功获取实时值的标签点个数
// * \param ids       整型数组，输入，标签点标识列表
// * \param datetimes 整型数组，输出，实时数值时间列表,
// *                    表示距离1970年1月1日08:00:00的秒数
// * \param ms        短整型数组，输出，实时数值时间列表，
// *                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
// * \param x         单精度浮点型数组，输出，实时浮点型横坐标数值列表
// * \param y         单精度浮点型数组，输出，实时浮点型纵坐标数值列表
// * \param qualities 短整型数组，输出，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
// * \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
// *        本接口只对数据类型为 RTDB_COOR 的标签点有效。
// rtdb_error RTDBAPI_CALLRULE rtdbs_get_coor_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float32* x, rtdb_float32* y, rtdb_int16* qualities, rtdb_error* errors)
func RawRtdbsGetCoorSnapshots64Warp() {}

// RawRtdbsPutCoorSnapshots64Warp 批量写入坐标实时数据
// *
// * \param handle    连接句柄
// * \param count     整型，输入/输出，标签点个数，
// *                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
// *                    输出时表示成功获取实时值的标签点个数
// * \param ids       整型数组，输入，标签点标识列表
// * \param datetimes 整型数组，输入，实时数值时间列表,
// *                    表示距离1970年1月1日08:00:00的秒数
// * \param ms        短整型数组，输入，实时数值时间列表，
// *                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
// * \param x         单精度浮点型数组，输入，实时浮点型横坐标数值列表
// * \param y         单精度浮点型数组，输入，实时浮点型纵坐标数值列表
// * \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \param errors    无符号整型数组，输出，写入实时坐标数据的返回值列表，参考rtdb_error.h
// * \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
// *        本接口只对数据类型为 RTDB_COOR 的标签点有效。
// rtdb_error RTDBAPI_CALLRULE rtdbs_put_coor_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float32* x, const rtdb_float32* y, const rtdb_int16* qualities, rtdb_error* errors)
func RawRtdbsPutCoorSnapshots64Warp() {}

// RawRtdbsPutCoorSnapshots 批量写入坐标实时数据
// *
// * \param handle    连接句柄
// * \param count     整型，输入/输出，标签点个数，
// *                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
// *                    输出时表示成功获取实时值的标签点个数
// * \param ids       整型数组，输入，标签点标识列表
// * \param datetimes 整型数组，输入，实时数值时间列表,
// *                    表示距离1970年1月1日08:00:00的秒数
// * \param ms        短整型数组，输入，实时数值时间列表，
// *                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
// * \param x         单精度浮点型数组，输入，实时浮点型横坐标数值列表
// * \param y         单精度浮点型数组，输入，实时浮点型纵坐标数值列表
// * \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \param errors    无符号整型数组，输出，写入实时坐标数据的返回值列表，参考rtdb_error.h
// * \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
// *        本接口只对数据类型为 RTDB_COOR 的标签点有效。
// *        仅当输入时间戳与当前快照时间戳完全相等时，会替换当前快照的值和质量；
// *        其余情况下会调用 rtdbs_put_coor_snapshots()
// rtdb_error RTDBAPI_CALLRULE rtdbs_fix_coor_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float32* x, const rtdb_float32* y, const rtdb_int16* qualities, rtdb_error* errors)
func RawRtdbsPutCoorSnapshots() {}

// RawRtdbsGetBlobSnapshot64Warp 读取二进制/字符串实时数据
// *
// * \param handle    连接句柄
// * \param id        整型，输入，标签点标识
// * \param datetime  整型，输出，实时数值时间列表,
// *                    表示距离1970年1月1日08:00:00的秒数
// * \param ms        短整型，输出，实时数值时间列表，
// *                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
// * \param blob      字节型数组，输出，实时二进制/字符串数值
// * \param len       短整型，输出，二进制/字符串数值长度
// * \param quality   短整型，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
// rtdb_error RTDBAPI_CALLRULE rtdbs_get_blob_snapshot64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type* datetime, rtdb_subtime_type* subtime, rtdb_byte* blob, rtdb_length_type* len, rtdb_int16* quality)
func RawRtdbsGetBlobSnapshot64Warp() {}

// RawRtdbsGetBlobSnapshots64Warp 批量读取二进制/字符串实时数据
// * \param handle    连接句柄
// * \param count     整型，输入/输出，标签点个数，
// *                    输入时表示 ids、datetimes、ms、blobs、lens、qualities、errors 的长度，
// *                    输出时表示成功获取实时值的标签点个数
// * \param ids       整型数组，输入，标签点标识
// * \param datetimes 整型数组，输出，实时数值时间列表,
// *                    表示距离1970年1月1日08:00:00的秒数
// * \param ms        短整型数组，输出，实时数值时间列表，
// *                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
// * \param blobs     字节型指针数组，输出，实时二进制/字符串数值
// * \param lens      短整型数组，输入/输出，二进制/字符串数值长度，
// *                    输入时表示对应的 blobs 指针指向的缓冲区长度，
// *                    输出时表示实际得到的 blob 长度，如果 blob 的长度大于缓冲区长度，会被截断。
// * \param qualities 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
// * \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
// rtdb_error RTDBAPI_CALLRULE rtdbs_get_blob_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_byte* const* blobs, rtdb_length_type* lens, rtdb_int16* qualities, rtdb_error* errors)
func RawRtdbsGetBlobSnapshots64Warp() {}

// RawRtdbsPutBlobSnapshot64Warp 写入二进制/字符串实时数据
// * \param handle    连接句柄
// * \param id        整型，输入，标签点标识
// * \param datetime  整型，输入，实时数值时间列表,
// * 表示距离1970年1月1日08:00:00的秒数
// * \param ms        短整型，输入，实时数值时间列表，
// * 对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
// * \param blob      字节型数组，输入，实时二进制/字符串数值
// * \param len       短整型，输入，二进制/字符串数值长度，超过一个页大小数据将被截断。
// * \param quality   短整型，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
// rtdb_error RTDBAPI_CALLRULE rtdbs_put_blob_snapshot64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, const rtdb_byte* blob, rtdb_length_type len, rtdb_int16 quality)
func RawRtdbsPutBlobSnapshot64Warp() {}

// RawRtdbsPutBlobSnapshots64Warp 批量写入二进制/字符串实时数据
// * \param handle    连接句柄
// * \param count     整型，输入/输出，标签点个数，
// *                    输入时表示 ids、datetimes、ms、blobs、lens、qualities、errors 的长度，
// *                    输出时表示成功获取实时值的标签点个数
// * \param ids       整型数组，输入，标签点标识
// * \param datetimes 整型数组，输入，实时数值时间列表,
// *                    表示距离1970年1月1日08:00:00的秒数
// * \param ms        短整型数组，输入，实时数值时间列表，
// *                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
// * \param blobs     字节型指针数组，输入，实时二进制/字符串数值
// * \param lens      短整型数组，输入，二进制/字符串数值长度，
// *                    表示对应的 blobs 指针指向的缓冲区长度，超过一个页大小数据将被截断。
// * \param qualities 短整型数组，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
// * \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
// rtdb_error RTDBAPI_CALLRULE rtdbs_put_blob_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_byte* const* blobs, const rtdb_length_type* lens, const rtdb_int16* qualities, rtdb_error* errors)
func RawRtdbsPutBlobSnapshots64Warp() {}

// RawRtdbsGetDatetimeSnapshots64Warp 批量读取datetime类型标签点实时数据
// * \param handle    连接句柄
// * \param count     整型，输入/输出，标签点个数，
// *                    输入时表示 ids、datetimes、ms、dtvalues、dtlens、qualities、errors 的长度，
// *                    输出时表示成功获取实时值的标签点个数
// * \param ids       整型数组，输入，标签点标识
// * \param datetimes 整型数组，输出，实时数值时间列表,
// *                    表示距离1970年1月1日08:00:00的秒数
// * \param ms        短整型数组，输出，实时数值时间列表，
// *                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
// * \param dtvalues  字节型指针数组，输出，实时datetime数值
// * \param dtlens    短整型数组，输入/输出，datetime数值长度，
// *                    输入时表示对应的 dtvalues 指针指向的缓冲区长度，
// * \param qualities 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
// * \param type      短整型，输入，所有标签点的显示类型，如“yyyy-mm-dd hh:mm:ss.000”的type为1，默认类型1，
// *                    “yyyy/mm/dd hh:mm:ss.000”的type为2
// *                    如果不传type，则按照标签点属性显示，否则按照type类型显示
// * \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
// rtdb_error RTDBAPI_CALLRULE rtdbs_get_datetime_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_byte* const* dtvalues, rtdb_length_type* dtlens, rtdb_int16* qualities, rtdb_error* errors, rtdb_int16 type)
func RawRtdbsGetDatetimeSnapshots64Warp() {}

// RawRtdbsPutDatetimeSnapshots64Warp 批量插入datetime类型标签点数据
// * \param handle      连接句柄
// * \param count       整型，输入/输出，标签点个数，
// *                      输入时表示 ids、datetimes、ms、dtvalues、dtlens、qualities、errors的长度，
// *                      输出时表示成功写入的标签点个数
// * \param ids         整型数组，输入，标签点标识
// * \param datetimes   整型数组，输入，实时值时间列表
// *                      表示距离1970年1月1日08:00:00的秒数
// * \param ms          短整型数组，输入，实时数值时间列表，
// *                      对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
// * \param dtvalues    字节型指针数组，输入，datetime标签点的值
// * \param dtlens      短整型数组，输入，数值长度
// * \param qualities   短整型数组，输入，实时数值品质，，数据库预定义的品质参见枚举 RTDB_QUALITY
// * \param errors      无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
// * \remark 被接口只对数据类型 RTDB_DATETIME 的标签点有效。
// rtdb_error RTDBAPI_CALLRULE rtdbs_put_datetime_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_byte* const* dtvalues, const rtdb_length_type* dtlens, const rtdb_int16* qualities, rtdb_error* errors)
func RawRtdbsPutDatetimeSnapshots64Warp() {}

// RawRtdbsSubscribeSnapshotsEx64Warp 批量标签点快照改变的通知订阅
// *
// * \param handle         连接句柄
// * \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
// *                           输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
// * \param ids            整型数组，输入，标签点标识列表。
// * \param options        订阅选项
// *                           RTDB_O_AUTOCONN 自动重连
// * \param param          用户自定义参数
// * \param callback       rtdbs_snaps_event_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
// *                       当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
// *                       作为event_type取值调用回掉函数后退出订阅。
// *                       当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
// *                       作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
// *                       网络中断期间回掉函数调用频率为最少3秒
// *                       event_type参数值含义如下：
// *                         RTDB_E_DATA        标签点快照改变
// *                         RTDB_E_DISCONNECT  订阅客户端与数据库网络断开
// *                         RTDB_E_RECOVERY    订阅客户端与数据库网络及订阅恢复
// *                         RTDB_E_CHANGED     客户端修改订阅标签点信息，即通过rtdbs_change_subscribe_snapshots修改订阅信息的结果
// *                       handle 产生订阅回掉的连接句柄，调用rtdbs_subscribe_snapshots_ex时的handle参数
// *                       param  用户自定义参数，调用rtdbs_subscribe_snapshots_ex时的param参数
// *                       count  event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示ids，datetimes,values等的数量
// *                              event_type为其它值时，count值为0
// *                       ids    event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示快照改变的标签点ID，数量由count指定
// *                              event_type为其它值时，ids值为NULL
// *                       datetimes event_type为RTDB_E_DATA时表示快照时间，数量由count指定
// *                                 event_type为RTDB_E_CHANGED时表示changed_types，即通过rtdbs_change_subscribe_snapshots传入的changed_types
// *                                 event_type为其它值时，datetimes值为NULL
// *                       ms     event_type为RTDB_E_DATA时表示快照的毫秒，数量由count指定
// *                              event_type为其它值时，ms值为NULL
// *                       values event_type为RTDB_E_DATA时表示浮点数据类型快照值，数量由count指定
// *                              event_type为RTDB_E_CHANGED时，表示delta_values，即通过rtdbs_change_subscribe_snapshots传入的delta_values
// *                              event_type为其它值时，values值为NULL
// *                       states event_type为RTDB_E_DATA时表示整形数据类型快照值，数量由count指定
// *                              event_type为RTDB_E_CHANGED时，表示delta_states，通过rtdbs_change_subscribe_snapshots传入的delta_states
// *                              event_type为其它值时，states值为NULL
// *                       qualities event_type为RTDB_E_DATA时表示快照质量码，数量由count指定
// *                              event_type为其它值时，qualities值为NULL
// *                       errors event_type为RTDB_E_DATA时表示快照错误码，数量由count指定
// *                              event_type为RTDB_E_CHANGED时，表示修改结果对应的错误码，数量由count指定
// *                              event_type为其它值时，errors值为NULL
// * \param errors         无符号整型数组，输出，
// *                           写入实时数据的返回值列表，参考rtdb_error.h
// * \remark   用户须保证 ids、errors 的长度与 count 一致。
// *        用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
// *        否则返回 RtE_OTHER_SDK_DOING 错误。
// rtdb_error RTDBAPI_CALLRULE rtdbs_subscribe_snapshots_ex64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, rtdb_uint32 options, void* param, rtdbs_snaps_event_ex64 callback, rtdb_error* errors)
func RawRtdbsSubscribeSnapshotsEx64Warp() {}

// RawRtdbsSubscribeDeltaSnapshots64Warp 批量标签点快照改变的通知订阅
//   - \param handle         连接句柄
//   - \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
//   - 输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
//   - \param ids            整型数组，输入，标签点标识列表。
//   - \param delta_values   double型数组，输入，订阅浮点类型标签点的容差值，变化超过设置的容差值才会推送
//   - \param delta_values   整型数组，输入，订阅整型标签点的容差值，变化超过设置的容差值才会推送
//   - \param options        订阅选项
//   - RTDB_O_AUTOCONN 自动重连
//   - \param param          用户自定义参数
//   - \param callback       rtdbs_snaps_event_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
//   - 当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
//   - 作为event_type取值调用回掉函数后退出订阅。
//   - 当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
//   - 作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
//   - 网络中断期间回掉函数调用频率为最少3秒
//   - event_type参数值含义如下：
//   - RTDB_E_DATA        标签点快照改变
//   - RTDB_E_DISCONNECT  订阅客户端与数据库网络断开
//   - RTDB_E_RECOVERY    订阅客户端与数据库网络及订阅恢复
//   - RTDB_E_CHANGED     客户端修改订阅标签点信息，即通过rtdbs_change_subscribe_snapshots修改订阅信息的结果
//   - handle 产生订阅回掉的连接句柄，调用rtdbs_subscribe_snapshots_ex时的handle参数
//   - param  用户自定义参数，调用rtdbs_subscribe_snapshots_ex时的param参数
//   - count  event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示ids，datetimes,values等的数量
//   - event_type为其它值时，count值为0
//   - ids    event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示快照改变的标签点ID，数量由count指定
//   - event_type为其它值时，ids值为NULL
//   - datetimes event_type为RTDB_E_DATA时表示快照时间，数量由count指定
//   - event_type为RTDB_E_CHANGED时表示changed_types，即通过rtdbs_change_subscribe_snapshots传入的changed_types
//   - event_type为其它值时，datetimes值为NULL
//   - ms     event_type为RTDB_E_DATA时表示快照的毫秒，数量由count指定
//   - event_type为其它值时，ms值为NULL
//   - values event_type为RTDB_E_DATA时表示浮点数据类型快照值，数量由count指定
//   - event_type为RTDB_E_CHANGED时，表示delta_values，即通过rtdbs_change_subscribe_snapshots传入的delta_values
//   - event_type为其它值时，values值为NULL
//   - states event_type为RTDB_E_DATA时表示整形数据类型快照值，数量由count指定
//   - event_type为RTDB_E_CHANGED时，表示delta_states，通过rtdbs_change_subscribe_snapshots传入的delta_states
//   - event_type为其它值时，states值为NULL
//   - qualities event_type为RTDB_E_DATA时表示快照质量码，数量由count指定
//   - event_type为其它值时，qualities值为NULL
//   - errors event_type为RTDB_E_DATA时表示快照错误码，数量由count指定
//   - event_type为RTDB_E_CHANGED时，表示修改结果对应的错误码，数量由count指定
//   - event_type为其它值时，errors值为NULL
//   - \param errors         无符号整型数组，输出，
//   - 写入实时数据的返回值列表，参考rtdb_error.h
//   - \remark delta_values和delta_states可以为空指针，表示不设置容差值。 只有两个参数都不为空时，设置容差值才会生效。
//   - 用户须保证 ids、delta_values、delta_states、errors 的长度与 count 一致
//   - 用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
//   - 否则返回 RtE_OTHER_SDK_DOING 错误。
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_subscribe_delta_snapshots64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_float64* delta_values, const rtdb_int64* delta_states, rtdb_uint32 options, void* param, rtdbs_snaps_event_ex64 callback, rtdb_error* errors)
func RawRtdbsSubscribeDeltaSnapshots64Warp() {}

// RawRtdbsChangeSubscribeSnapshotsWarp 批量修改订阅标签点信息
//   - \param handle         连接句柄
//   - \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
//   - 输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
//   - \param ids            整型数组，输入，标签点标识列表。
//   - \param delta_values   double型数组，输入，订阅浮点类型标签点的容差值，变化超过设置的容差值才会推送
//   - \param delta_values   整型数组，输入，订阅整型标签点的容差值，变化超过设置的容差值才会推送
//   - \param changed_types  整型数组，输入，修改类型，参考RTDB_SUBSCRIBE_CHANGE_TYPE
//   - \param errors         异步调用，保留参数，暂时不启用
//   - \remark   用户须保证 ids、delta_values、delta_states、errors 的长度与 count 一致。
//   - 可以同时添加、修改、删除订阅的标签点信息，
//   - delta_values和delta_states，可以为空指针，为空，则表示不设置容差值，即写入新数据即推送
//   - 只有delta_values和delta_states都不为空时，设置的容差值才有效。
//   - 用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
//   - 否则返回 RtE_OTHER_SDK_DOING 错误。
//   - 此方法是异步方法，当网络中断等异常情况时，会通过方法的返回值返回错误，参考rtdb_error.h。
//   - 当方法返回值为RtE_OK时，表示已经成功发送给数据库，但是并没有等待修改结果。
//   - 数据库的修改结果，会异步通知给api的回调函数，通过rtdbs_snaps_event_ex的RTDB_E_CHANGED事件通知修改结果
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_change_subscribe_snapshots_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_float64* delta_values, const rtdb_int64* delta_states, const rtdb_int32* changed_types, rtdb_error* errors)
func RawRtdbsChangeSubscribeSnapshotsWarp() {}

// RawRtdbsCancelSubscribeSnapshotsWarp 取消标签点快照更改通知订阅
//   - \param handle    连接句柄
//
// rtdb_error RTDBAPI_CALLRULE rtdbs_cancel_subscribe_snapshots_warp(rtdb_int32 handle)
func RawRtdbsCancelSubscribeSnapshotsWarp() {}

// RawRtdbsGetNamedTypeSnapshot64Warp 获取自定义类型测点的单个快照
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
func RawRtdbsGetNamedTypeSnapshot64Warp() {}

// RawRtdbsGetNamedTypeSnapshots64Warp 批量获取自定义类型测点的快照
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
func RawRtdbsGetNamedTypeSnapshots64Warp() {}

// RawRtdbsPutNamedTypeSnapshot64Warp 写入单个自定义类型标签点的快照
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
func RawRtdbsPutNamedTypeSnapshot64Warp() {}

// RawRtdbsPutNamedTypeSnapshots64Warp 批量写入自定义类型标签点的快照
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
func RawRtdbsPutNamedTypeSnapshots64Warp() {}

// RawRtdbaGetArchivesCountWarp 获取存档文件数量
//   - \param handle    连接句柄
//   - \param count     整型，输出，存档文件数量
//
// rtdb_error RTDBAPI_CALLRULE rtdba_get_archives_count_warp(rtdb_int32 handle, rtdb_int32 *count)
func RawRtdbaGetArchivesCountWarp() {}

// RawRtdbaCreateRangedArchive64Warp 新建指定时间范围的历史存档文件并插入到历史数据库
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名。
// * \param begin      整数，输入，起始时间，距离1970年1月1日08:00:00的秒数
// * \param end        整数，输入，终止时间，距离1970年1月1日08:00:00的秒数
// * \param mb_size    整型，输入，文件兆字节大小，单位为 MB。
// rtdb_error RTDBAPI_CALLRULE rtdba_create_ranged_archive64_warp(rtdb_int32 handle, const char* path, const char* file, rtdb_timestamp_type begin, rtdb_timestamp_type end, rtdb_int32 mb_size)
func RawRtdbaCreateRangedArchive64Warp() {}

// RawRtdbaAppendArchiveWarp 追加磁盘上的历史存档文件到历史数据库。
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名，后缀名应为.rdf。
// * \param state      整型，输入，取值 RTDB_ACTIVED_ARCHIVE、RTDB_NORMAL_ARCHIVE、
// *                     RTDB_READONLY_ARCHIVE 之一，表示文件状态
// rtdb_error RTDBAPI_CALLRULE rtdba_append_archive_warp(rtdb_int32 handle, const char *path, const char *file, rtdb_int32 state)
func RawRtdbaAppendArchiveWarp() {}

// RawRtdbaRemoveArchiveWarp 从历史数据库中移出历史存档文件。
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名。
// rtdb_error RTDBAPI_CALLRULE rtdba_remove_archive_warp(rtdb_int32 handle, const char *path, const char *file)
func RawRtdbaRemoveArchiveWarp() {}

// RawRtdbaShiftActivedWarp 切换活动文件
//   - \param handle     连接句柄
//   - \remark 当前活动文件被写满时该事务被启动，
//   - 改变当前活动文件的状态为普通状态，
//   - 在所有历史数据存档文件中寻找未被使用过的
//   - 插入到前活动文件的右侧并改为活动状态，
//   - 若找不到则将前活动文件右侧的文件改为活动状态，
//   - 并将active_archive_指向该文件。该事务进行过程中，
//   - 用锁保证所有读写操作都暂停等待该事务完成。
//
// rtdb_error RTDBAPI_CALLRULE rtdba_shift_actived_warp(rtdb_int32 handle)
func RawRtdbaShiftActivedWarp() {}

// RawRtdbaGetArchivesWarp 获取存档文件的路径、名称、状态和最早允许写入时间。
//   - [handle]          连接句柄
//   - [paths]            字符串数组，输出，存档文件的目录路径，长度至少为 RTDB_PATH_SIZE。
//   - [files]            字符串数组，输出，存档文件的名称，长度至少为 RTDB_FILE_NAME_SIZE。
//   - [states]           整型数组，输出，取值 RTDB_INVALID_ARCHIVE、RTDB_ACTIVED_ARCHIVE、
//   - RTDB_NORMAL_ARCHIVE、RTDB_READONLY_ARCHIVE 之一，表示文件状态
//
// rtdb_error RTDBAPI_CALLRULE rtdba_get_archives_warp(rtdb_int32 handle, rtdb_int32* count, rtdb_path_string* paths, rtdb_filename_string* files, rtdb_int32 *states)
func RawRtdbaGetArchivesWarp() {}

// RawRtdbaGetArchivesInfoWarp 获取存档信息
//   - [handle]: in, 句柄
//   - [count]: out, 数量
//   - [paths]: out, 路径
//   - [files]: out, 文件
//   - [infos]: out, 存档信息
//   - [errors]: out, 错误
//
// rtdb_error RTDBAPI_CALLRULE rtdba_get_archives_info_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_path_string* const paths, const rtdb_filename_string* const files, RTDB_HEADER_PAGE *infos, rtdb_error* errors)
func RawRtdbaGetArchivesInfoWarp() {}

// RawRtdbaGetArchivesPerfDataWarp 获取存档的实时信息
//   - [handle]: in, 句柄
//   - [count]: out, 数量
//   - [paths]: out, 路径
//   - [files]: out, 文件
//   - [real_time_datas]: out, 实时数据
//   - [total_datas]: 总数
//   - [errors]: 错误
//
// rtdb_error RTDBAPI_CALLRULE rtdba_get_archives_perf_data_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_path_string* const paths, const rtdb_filename_string* const files, RTDB_ARCHIVE_PERF_DATA* real_time_datas, RTDB_ARCHIVE_PERF_DATA* total_datas, rtdb_error* errors)
func RawRtdbaGetArchivesPerfDataWarp() {}

// RawRtdbaGetArchivesStatusWarp 获取存档状态
//   - [handle]: in, 句柄
//   - [status]: out, 存档状态
//
// rtdb_error RTDBAPI_CALLRULE rtdba_get_archives_status_warp(rtdb_int32 handle, rtdb_error* status)
func RawRtdbaGetArchivesStatusWarp() {}

// RawRtdbaGetArchiveInfoWarp 获取存档文件及其附属文件的详细信息。
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名。
// * \param file_id    整型，输入，附属文件标识，0 表示获取主文件信息。
// * \param info       RTDB_HEADER_PAGE 结构，输出，存档文件信息
// rtdb_error RTDBAPI_CALLRULE rtdba_get_archive_info_warp(rtdb_int32 handle, const char *path, const char *file, rtdb_int32 file_id, RTDB_HEADER_PAGE *info)
func RawRtdbaGetArchiveInfoWarp() {}

// RawRtdbaUpdateArchiveWarp 修改存档文件的可配置项。
// * \param handle         连接句柄
// * \param path           字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file           字符串，输入，文件名。
// * \param rated_capacity 整型，输入，文件额定大小，单位为 MB。
// * \param ex_capacity    整型，输入，附属文件大小，单位为 MB。
// * \param auto_merge     短整型，输入，是否自动合并附属文件。
// * \param auto_arrange   短整型，输入，是否自动整理存档文件。
// * 备注: rated_capacity 与 ex_capacity 参数可为 0，表示不修改对应的配置项。
// rtdb_error RTDBAPI_CALLRULE rtdba_update_archive_warp(rtdb_int32 handle, const char *path, const char *file, rtdb_int32 rated_capacity, rtdb_int32 ex_capacity, rtdb_int16 auto_merge, rtdb_int16 auto_arrange)
func RawRtdbaUpdateArchiveWarp() {}

// RawRtdbaArrangeArchiveWarp 整理存档文件，将同一标签点的数据块存放在一起以提高查询效率。
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名。
// rtdb_error RTDBAPI_CALLRULE rtdba_arrange_archive_warp(rtdb_int32 handle, const char *path, const char *file)
func RawRtdbaArrangeArchiveWarp() {}

// RawRtdbaReindexArchiveWarp 为存档文件重新生成索引，用于恢复数据。
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名。
// rtdb_error RTDBAPI_CALLRULE rtdba_reindex_archive_warp(rtdb_int32 handle, const char *path, const char *file)
func RawRtdbaReindexArchiveWarp() {}

// RawRtdbaBackupArchiveWarp 备份主存档文件及其附属文件到指定路径
// * \param handle     连接句柄
// * \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// * \param file       字符串，输入，文件名。
// * \param dest       字符串，输入，备份目录路径，必须以"\"或"/"结尾。
// rtdb_error RTDBAPI_CALLRULE rtdba_backup_archive_warp(rtdb_int32 handle, const char *path, const char *file, const char *dest)
func RawRtdbaBackupArchiveWarp() {}

// RawRtdbaMoveArchiveWarp 将存档文件移动到指定目录
// *        [handle]     连接句柄
// *        [path]       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// *        [file]       字符串，输入，文件名。
// *        [dest]       字符串，输入，移动目录路径，必须以"\"或"/"结尾。
// rtdb_error RTDBAPI_CALLRULE rtdba_move_archive_warp(rtdb_int32 handle, const char *path, const char *file, const char *dest)
func RawRtdbaMoveArchiveWarp() {}

// RawRtdbaConvertIndexWarp 为存档文件转换索引格式。
// *        [handle]     连接句柄
// *        [path]       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
// *        [file]       字符串，输入，文件名。
// rtdb_error RTDBAPI_CALLRULE rtdba_convert_index_warp(rtdb_int32 handle, const char *path, const char *file)
func RawRtdbaConvertIndexWarp() {}

// RawRtdbaQueryBigJob64Warp 查询进程正在执行的后台任务类型、状态和进度
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
func RawRtdbaQueryBigJob64Warp() {}

// RawRtdbaCancelBigJobWarp 取消进程正在执行的后台任务
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
func RawRtdbaCancelBigJobWarp() {}

// RawRtdbhArchivedValuesCount64Warp 获取单个标签点在一段时间范围内的存储值数量.
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
//   - \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
//   - \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
//   - \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
//   - \param count         整型，输出，返回上述时间范围内的存储值数量
//   - \remark 由 datetime1、ms1 形成的时间可以大于 datetime2、ms2 表示的时间，
//   - 此时前者表示结束时间，后者表示起始时间。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_archived_values_count64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, rtdb_int32* count)
func RawRtdbhArchivedValuesCount64Warp() {}

// RawRtdbhArchivedValuesRealCount64Warp 获取单个标签点在一段时间范围内的真实的存储值数量.
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
//   - \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
//   - \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
//   - \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
//   - \param count         整型，输出，返回上述时间范围内的存储值数量
//   - \remark 由 datetime1、ms1 形成的时间可以大于 datetime2、ms2 表示的时间，
//   - 此时前者表示结束时间，后者表示起始时间。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_archived_values_real_count64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, rtdb_int32* count)
func RawRtdbhArchivedValuesRealCount64Warp() {}

// RawRtdbhGetArchivedValues64Warp 读取单个标签点一段时间内的储存数据
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param count         整型，输入/输出，
//   - 输入时表示 datetimes、ms、values、states、qualities 的长度；
//   - 输出时返回实际得到的数值个数
//   - \param datetimes     整型数组，输入/输出，
//   - 输入时第一个元素表示起始时间秒数，
//   - 最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
//   - 输出时表示对应的历史数值时间秒数。
//   - \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//   - 则输入时第一个元素表示起始时间纳秒，
//   - 最后一个元素表示结束时间纳秒；
//   - 输出时表示对应的历史数值时间纳秒。
//   - 否则忽略输入，输出时为 0。
//   - \param values        双精度浮点数数组，输出，历史浮点型数值列表
//   - 对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史存储值；否则为 0
//   - \param states        64 位整数数组，输出，历史整型数值列表，
//   - 对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//   - RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史存储值；否则为 0
//   - \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
//   - 在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
//   - 大于最后一个元素形成的时间，此时第一个元素表示结束时间，
//   - 最后一个元素表示开始时间。
//   - 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
func RawRtdbhGetArchivedValues64Warp() {}

// RawRtdbhGetArchivedValuesBackward64Warp 逆向读取单个标签点一段时间内的储存数据
//
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
func RawRtdbhGetArchivedValuesBackward64Warp() {}

// RawRtdbhGetArchivedCoorValues64Warp 读取单个标签点一段时间内的坐标型储存数据
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param count         整型，输入/输出，
//   - 输入时表示 datetimes、ms、x、y、qualities 的长度；
//   - 输出时返回实际得到的数值个数
//   - \param datetimes     整型数组，输入/输出，
//   - 输入时第一个元素表示起始时间秒数，
//   - 最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
//   - 输出时表示对应的历史数值时间秒数。
//   - \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//   - 则输入时第一个元素表示起始时间纳秒，
//   - 最后一个元素表示结束时间纳秒；
//   - 输出时表示对应的历史数值时间纳秒。
//   - 否则忽略输入，输出时为 0。
//   - \param x             单精度浮点型数组，输出，浮点型横坐标历史数值列表
//   - \param y             单精度浮点型数组，输出，浮点型纵坐标历史数值列表
//   - \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 用户须保证 datetimes、ms、x、y、qualities 的长度与 count 一致，
//   - 在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
//   - 大于最后一个元素形成的时间，此时第一个元素表示结束时间，
//   - 最后一个元素表示开始时间。
//   - 本接口只对数据类型为 RTDB_COOR 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_coor_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float32* x, rtdb_float32* y, rtdb_int16* qualities)
func RawRtdbhGetArchivedCoorValues64Warp() {}

// RawRtdbhGetArchivedCoorValuesBackward64Warp 逆向读取单个标签点一段时间内的坐标型储存数据
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param count         整型，输入/输出，
//   - 输入时表示 datetimes、ms、x、y、qualities 的长度；
//   - 输出时返回实际得到的数值个数
//   - \param datetimes     整型数组，输入/输出，
//   - 输入时第一个元素表示起始时间秒数，
//   - 最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
//   - 输出时表示对应的历史数值时间秒数。
//   - \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//   - 则输入时第一个元素表示起始时间纳秒，
//   - 最后一个元素表示结束时间纳秒；
//   - 输出时表示对应的历史数值时间纳秒。
//   - 否则忽略输入，输出时为 0。
//   - \param x             单精度浮点型数组，输出，浮点型横坐标历史数值列表
//   - \param y             单精度浮点型数组，输出，浮点型纵坐标历史数值列表
//   - \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 用户须保证 datetimes、ms、x、y、qualities 的长度与 count 一致，
//   - 在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
//   - 大于最后一个元素形成的时间，此时第一个元素表示结束时间，
//   - 最后一个元素表示开始时间。
//   - 本接口只对数据类型为 RTDB_COOR 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_coor_values_backward64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float32* x, rtdb_float32* y, rtdb_int16* qualities)
func RawRtdbhGetArchivedCoorValuesBackward64Warp() {}

// RawRtdbhGetArchivedValuesInBatches64Warp 开始以分段返回方式读取一段时间内的储存数据
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
//   - \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
//   - \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
//   - \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
//   - \param count         整型，输出，返回上述时间范围内的存储值数量
//   - \param batch_count   整型，输出，每次分段返回的长度，用于继续调用 rtdbh_get_next_archived_values 接口
//   - \remark 由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
//   - 此时前者表示结束时间，后者表示起始时间。
//   - 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_values_in_batches64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, rtdb_int32* count, rtdb_int32* batch_count)
func RawRtdbhGetArchivedValuesInBatches64Warp() {}

// RawRtdbhGetNextArchivedValues64Warp 分段读取一段时间内的储存数据
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param count         整形，输入/输出，
//   - 输入时表示 datetimes、ms、values、states、qualities 的长度；
//   - 输出时表示实际得到的存储值个数。
//   - \param datetimes     整型数组，输出，历史数值时间列表,
//   - 表示距离1970年1月1日08:00:00的秒数
//   - \param ms            短整型数组，输出，历史数值时间列表，
//   - 对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
//   - \param values        双精度浮点型数组，输出，历史浮点型数值列表，
//   - 对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的历史存储值；否则为 0
//   - \param states        64 位整型数组，输出，历史整型数值列表，
//   - 对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//   - RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的历史存储值；否则为 0
//   - \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 相符，
//   - 且 count 不能小于 rtdbh_get_archived_values_in_batches 接口中返回的 batch_count 的值，
//   - 当返回 RtE_BATCH_END 表示全部数据获取完毕。
//   - 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_next_archived_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
func RawRtdbhGetNextArchivedValues64Warp() {}

// RawRtdbhGetTimedValues64Warp 获取单个标签点的单调递增时间序列历史插值。
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param count         整型，输入，表示 datetimes、ms、values、states、qualities 的长度。
//   - \param datetimes     整型数组，输入，表示需要的单调递增时间列表，
//   - 为距离1970年1月1日08:00:00的秒数
//   - \param ms            短整型数组，输入，对于时间精度为纳秒的标签点，
//   - 表示需要的单调递增时间对应的纳秒值；否则忽略。
//   - \param values        双精度浮点型数组，输出，历史浮点型数值列表，
//   - 对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的历史插值；否则为 0
//   - \param states        64 位整型数组，输出，历史整型数值列表，
//   - 对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//   - RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的历史插值；否则为 0
//   - \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 相符，
//   - 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_timed_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 count, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
func RawRtdbhGetTimedValues64Warp() {}

// RawRtdbhGetTimedCoorValues64Warp 获取单个坐标标签点的单调递增时间序列历史插值。
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param count         整型，输入，表示 datetimes、ms、x、y、qualities 的长度。
//   - \param datetimes     整型数组，输入，表示需要的单调递增时间列表，
//   - 为距离1970年1月1日08:00:00的秒数
//   - \param ms            短整型数组，输入，对于时间精度为纳秒的标签点，
//   - 表示需要的单调递增时间对应的纳秒值；否则忽略。
//   - \param x             单精度浮点型数组，输出，浮点型横坐标历史插值数值列表
//   - \param y             单精度浮点型数组，输出，浮点型纵坐标历史插值数值列表
//   - \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 用户须保证 datetimes、ms、x、y、qualities 的长度与 count 相符，
//   - 本接口只对数据类型为 RTDB_COOR 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_timed_coor_values64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 count, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, rtdb_float32* x, rtdb_float32* y, rtdb_int16* qualities)
func RawRtdbhGetTimedCoorValues64Warp() {}

//	RawRtdbhGetInterpoValues64Warp 获取单个标签点一段时间内等间隔历史插值
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
func RawRtdbhGetInterpoValues64Warp() {}

// RawRtdbhGetIntervalValues64Warp 读取单个标签点某个时刻之后一定数量的等间隔内插值替换的历史数值
//
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
func RawRtdbhGetIntervalValues64Warp() {}

// RawRtdbhGetSingleValue64Warp 读取单个标签点某个时间的历史数据
//
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
func RawRtdbhGetSingleValue64Warp() {}

// RawRtdbhGetSingleCoorValue64Warp 读取单个标签点某个时间的坐标型历史数据
//
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
func RawRtdbhGetSingleCoorValue64Warp() {}

// RawRtdbhGetSingleBlobValue64Warp 读取单个标签点某个时间的二进制/字符串型历史数据
//
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
func RawRtdbhGetSingleBlobValue64Warp() {}

// RawRtdbhGetArchivedBlobValues64Warp 读取单个标签点一段时间的二进制/字符串型历史数据
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
func RawRtdbhGetArchivedBlobValues64Warp() {}

// RawRtdbhGetArchivedBlobValuesFilt64Warp 读取并模糊搜索单个标签点一段时间的二进制/字符串型历史数据
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
//   - \param filter        字符串，输入，支持通配符的模糊搜索字符串，多个模糊搜索的条件通过空格分隔，只针对string类型有效
//   - 当filter为空指针时，表示不进行过滤,
//   - 限制最大长度为RTDB_EQUATION_SIZE-1，超过此长度会返回错误
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
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_archived_blob_values_filt64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, const char* filter, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_length_type* lens, rtdb_byte* const* blobs, rtdb_int16* qualities)
func RawRtdbhGetArchivedBlobValuesFilt64Warp() {}

// RawRtdbhGetSingleDatetimeValue64Warp 读取单个标签点某个时间的datetime历史数据
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
func RawRtdbhGetSingleDatetimeValue64Warp() {}

// RawRtdbhGetArchivedDatetimeValues64Warp 读取单个标签点一段时间的时间类型历史数据
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
func RawRtdbhGetArchivedDatetimeValues64Warp() {}

// RawRtdbhPutArchivedDatetimeValues64Warp 写入批量标签点批量时间型历史存储数据
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
func RawRtdbhPutArchivedDatetimeValues64Warp() {}

// RawRtdbhSummaryDataWarp 获取单个标签点一段时间内的统计值。
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
func RawRtdbhSummaryDataWarp() {}

// RawRtdbhSummaryDataInBatchesWarp 分批获取单一标签点一段时间内的统计值
//
//	  *
//	- \param handle            连接句柄
//	- \param id                整型，输入，标签点标识
//	- \param count             整形，输入/输出，输入时表示 datatimes1、ms1、datatimes2、ms2、
//	- max_values、min_values、total_values、calc_avgs、power_avgs、errors 的长度，
//	- 即分段的个数；输出时表示成功取得统计值的分段个数。
//	- \param interval          64 位整型，输入，分段时间间隔，单位为纳秒。
//	- 如果为纳秒点，输入时间必须大于1纳秒，如果为秒级点，则必须大于1000000000纳秒。
//	- \param datetimes1        整型数组，输入/输出，输入时第一个元素表示起始时间秒数。
//	- 如果为 0，表示从存档中最早时间的数据开始进行统计。
//	- 输出时返回各个分段对应的最大值的时间秒数。
//	- \param ms1               短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	- 第一个元素表示起始时间对应的纳秒，
//	- 输出时返回各个分段对应的最大值的时间纳秒数；否则忽略，返回值为 0
//	- \param datetimes2        整型数组，输入/输出，输入时第一个元素表示结束时间秒数。
//	- 如果为 0，表示统计到存档中最近时间的数据为止。
//	- 输出时返回各个分段对应的最小值的时间秒数。
//	- \param ms2               短整型数组，如果 id 指定的标签点时间精度为纳秒，
//	- 第一个元素表示结束时间对应的纳秒，
//	- 输出时返回各个分段对应的最小值的时间纳秒数；否则忽略，返回值为 0
//	- \param max_values        双精度浮点型数组，输出，表示统计时间段内的最大数值。
//	- \param min_values        双精度浮点型数组，输出，表示统计时间段内的最小数值。
//	- \param total_values      双精度浮点型数组，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
//	- \param calc_avgs         双精度浮点型数组，输出，表示统计时间段内的算术平均值。
//	- \param power_avgs        双精度浮点型数组，输出，表示统计时间段内的加权平均值。
//	- \param errors            无符号整型数组，输出，表示各个分段取得统计值的返回值。
//	- \remark 由 datetimes1[0]、ms1[0] 表示的时间可以大于 datetimes2[0]、ms2[0] 表示的时间，
//	- 此时前者表示结束时间，后者表示起始时间。
//	- 如果输出的最大值或最小值的时间戳秒值为 0，
//	- 则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
//	- 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_summary_data_in_batches_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32* count, rtdb_int64 interval, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, RTDB_SUMMARY_DATA* summary_datas, rtdb_error* errors)
func RawRtdbhSummaryDataInBatchesWarp() {}

//	RawRtdbhGetPlotValues64Warp 获取单个标签点一段时间内用于绘图的历史数据
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
func RawRtdbhGetPlotValues64Warp() {}

// RawRtdbhGetCrossSectionValues64Warp 获取批量标签点在某一时间的历史断面数据
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
func RawRtdbhGetCrossSectionValues64Warp() {}

// RawRtdbhGetArchivedValuesFilt64Warp 读取单个标签点在一段时间内经复杂条件筛选后的历史储存值
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
func RawRtdbhGetArchivedValuesFilt64Warp() {}

// RawRtdbhGetIntervalValuesFilt64Warp 读取单个标签点某个时刻之后经复杂条件筛选后一定数量的等间隔内插值替换的历史数值
//
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param filter        字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
//   - 长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
//   - \param interval      整型，输入，插值时间间隔，单位为纳秒
//   - \param count         整型，输入，表示 datetimes、ms、values、states、qualities 的长度，
//   - 即需要的插值个数。
//   - \param datetimes     整型数组，输入/输出，
//   - 输入时第一个元素表示起始时间秒数；
//   - 输出时表示对应的历史数值时间秒数。
//   - \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//   - 则输入时第一个元素表示起始时间纳秒；
//   - 输出时表示对应的历史数值时间纳秒。
//   - 否则忽略输入，输出时为 0。
//   - \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
//   - 对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
//   - \param states        64 位整数数组，输出，整型历史插值数值列表，
//   - 对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//   - RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
//   - \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
//   - 在输入时，datetimes、ms 中至少应有一个元素用于表示起始时间。
//   - 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_interval_values_filt64_warp(rtdb_int32 handle, rtdb_int32 id, const char* filter, rtdb_int64 interval, rtdb_int32 count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
func RawRtdbhGetIntervalValuesFilt64Warp() {}

// RawRtdbhGetInterpoValuesFilt64Warp 获取单个标签点一段时间内经复杂条件筛选后的等间隔插值
//
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param filter        字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
//   - 长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
//   - \param count         整型，输入/输出，
//   - 输入时表示 datetimes、ms、values、states、qualities 的长度，
//   - 即需要的插值个数；输出时返回实际得到的插值个数
//   - \param datetimes     整型数组，输入/输出，
//   - 输入时第一个元素表示起始时间秒数，
//   - 最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
//   - 输出时表示对应的历史数值时间秒数。
//   - \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//   - 则输入时第一个元素表示起始时间纳秒，
//   - 最后一个元素表示结束时间纳秒；
//   - 输出时表示对应的历史数值时间纳秒。
//   - 否则忽略输入，输出时为 0。
//   - \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
//   - 对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
//   - \param states        64 位整数数组，输出，整型历史插值数值列表，
//   - 对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//   - RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
//   - \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
//   - 在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
//   - 大于最后一个元素形成的时间，此时第一个元素表示结束时间，
//   - 最后一个元素表示开始时间。
//   - 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_get_interpo_values_filt64_warp(rtdb_int32 handle, rtdb_int32 id, const char* filter, rtdb_int32* count, rtdb_timestamp_type* datetimes, rtdb_subtime_type* subtimes, rtdb_float64* values, rtdb_int64* states, rtdb_int16* qualities)
func RawRtdbhGetInterpoValuesFilt64Warp() {}

// RawRtdbhSummaryDataFiltWarp 获取单个标签点一段时间内经复杂条件筛选后的统计值
//
//   - \param handle            连接句柄
//   - \param id                整型，输入，标签点标识
//   - \param filter            字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
//   - 长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
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
// rtdb_error RTDBAPI_CALLRULE rtdbh_summary_data_filt_warp(rtdb_int32 handle, rtdb_int32 id, const char* filter, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, RTDB_SUMMARY_DATA* summary_data)
func RawRtdbhSummaryDataFiltWarp() {}

// RawRtdbhSummaryDataFiltInBatchesWarp 分批获取单一标签点一段时间内经复杂条件筛选后的统计值
//
//	  *
//	- \param handle            连接句柄
//	- \param id                整型，输入，标签点标识
//	- \param filter            字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
//	- 长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
//	- \param count             整形，输入/输出，输入时表示 datatimes1、ms1、datatimes2、ms2、
//	- max_values、min_values、total_values、calc_avgs、power_avgs、errors 的长度，
//	- 即分段的个数；输出时表示成功取得统计值的分段个数。
//	- \param interval          64 位整型，输入，分段时间间隔，单位为纳秒。
//	- \param datetimes1        整型数组，输入/输出，输入时第一个元素表示起始时间秒数。
//	- 如果为 0，表示从存档中最早时间的数据开始进行统计。
//	- 输出时返回各个分段对应的最大值的时间秒数。
//	- \param ms1               短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
//	- 第一个元素表示起始时间对应的纳秒，
//	- 输出时返回各个分段对应的最大值的时间纳秒数；否则忽略，返回值为 0
//	- \param datetimes2        整型数组，输入/输出，输入时第一个元素表示结束时间秒数。
//	- 如果为 0，表示统计到存档中最近时间的数据为止。
//	- 输出时返回各个分段对应的最小值的时间秒数。
//	- \param ms2               短整型数组，如果 id 指定的标签点时间精度为纳秒，
//	- 第一个元素表示结束时间对应的纳秒，
//	- 输出时返回各个分段对应的最小值的时间纳秒数；否则忽略，返回值为 0
//	- \param max_values        双精度浮点型数组，输出，表示统计时间段内的最大数值。
//	- \param min_values        双精度浮点型数组，输出，表示统计时间段内的最小数值。
//	- \param total_values      双精度浮点型数组，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
//	- \param calc_avgs         双精度浮点型数组，输出，表示统计时间段内的算术平均值。
//	- \param power_avgs        双精度浮点型数组，输出，表示统计时间段内的加权平均值。
//	- \param errors            无符号整型数组，输出，表示各个分段取得统计值的返回值。
//	- \remark 由 datetimes1[0]、ms1[0] 表示的时间可以大于 datetimes2[0]、ms2[0] 表示的时间，
//	- 此时前者表示结束时间，后者表示起始时间。
//	- 如果输出的最大值或最小值的时间戳秒值为 0，
//	- 则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
//	- 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_summary_data_filt_in_batches_warp(rtdb_int32 handle, rtdb_int32 id, const char* filter, rtdb_int32* count, rtdb_int64 interval, rtdb_timestamp_type datetime1, rtdb_subtime_type subtime1, rtdb_timestamp_type datetime2, rtdb_subtime_type subtime2, RTDB_SUMMARY_DATA* summary_datas, rtdb_error* errors)
func RawRtdbhSummaryDataFiltInBatchesWarp() {}

// RawRtdbhUpdateValue64Warp 修改单个标签点某一时间的历史存储值.
//
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param datetime      整型，输入，时间秒数
//   - \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
//   - 表示时间纳秒数；否则忽略。
//   - \param value         双精度浮点数，输入，浮点型历史数值
//   - 对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放新的历史值；否则忽略
//   - \param state         64 位整数，输入，整型历史数值，
//   - 对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//   - RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放新的历史值；否则忽略
//   - \param quality       短整型，输入，新的历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_update_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, rtdb_float64 value, rtdb_int64 state, rtdb_int16 quality)
func RawRtdbhUpdateValue64Warp() {}

// RawRtdbhUpdateCoorValue64Warp 修改单个标签点某一时间的历史存储值.
//
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
func RawRtdbhUpdateCoorValue64Warp() {}

// RawRtdbhRemoveValue64Warp 删除单个标签点某个时间的历史存储值
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param datetime      整型，输入，时间秒数
//   - \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
//   - 表示时间纳秒数；否则忽略。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_remove_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime)
func RawRtdbhRemoveValue64Warp() {}

// RawRtdbhRemoveValues64Warp 删除单个标签点一段时间内的历史存储值
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
func RawRtdbhRemoveValues64Warp() {}

// RawRtdbhPutSingleValue64Warp 写入单个标签点在某一时间的历史数据。
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param datetime      整型，输入，时间秒数
//   - \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
//   - 表示时间纳秒数；否则忽略。
//   - \param value         双精度浮点数，输入，浮点型历史数值
//   - 对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放历史值；否则忽略
//   - \param state         64 位整数，输入，整型历史数值，
//   - 对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//   - RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放历史值；否则忽略
//   - \param quality       短整型，输入，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//   - 如果 datetimes、ms 标识的数据已经存在，其值将被替换。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_single_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, rtdb_float64 value, rtdb_int64 state, rtdb_int16 quality)
func RawRtdbhPutSingleValue64Warp() {}

// RawRtdbhPutSingleCoorValue64Warp 写入单个标签点在某一时间的坐标型历史数据。
//   - \param handle              连接句柄
//   - \param id            整型，输入，标签点标识
//   - \param datetime      整型，输入，时间秒数
//   - \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
//   - 表示时间纳秒数；否则忽略。
//   - \param x             单精度浮点型，输入，横坐标历史数值
//   - \param y             单精度浮点型，输入，纵坐标历史数值
//   - \param quality       短整型，输入，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//   - 如果 datetimes、ms 标识的数据已经存在，其值将被替换。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_single_coor_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, rtdb_float32 x, rtdb_float32 y, rtdb_int16 quality)
func RawRtdbhPutSingleCoorValue64Warp() {}

// RawRtdbhPutSingleBlobValue64Warp 写入单个二进制/字符串标签点在某一时间的历史数据
//   - \param handle    连接句柄
//   - \param id        整型，输入，标签点标识
//   - \param datetime  整型，输入，数值时间列表,
//   - 表示距离1970年1月1日08:00:00的秒数
//   - \param ms        短整型，输入，历史数值时间，
//   - 对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
//   - \param blob      字节型数组，输入，历史二进制/字符串数值
//   - \param len       短整型，输入，二进制/字符串数值长度，超过一个页大小数据将被截断。
//   - \param quality   短整型，输入，历史数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//   - \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_single_blob_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, const rtdb_byte* blob, rtdb_length_type len, rtdb_int16 quality)
func RawRtdbhPutSingleBlobValue64Warp() {}

// RawRtdbhPutArchivedValues64Warp 写入批量标签点批量历史存储数据
//
//	  *
//	- \param handle        连接句柄
//	- \param count         整型，输入/输出，
//	- 输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
//	- 即历史值个数；输出时返回实际写入的数值个数
//	- \param ids           整型数组，输入，标签点标识，同一个标签点标识可以出现多次，
//	- 但它们的时间戳必需是递增的。
//	- \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
//	- \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
//	- 表示对应的历史数值时间纳秒；否则忽略。
//	- \param values        双精度浮点数数组，输入，浮点型历史数值列表
//	- 对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，表示相应的历史存储值；否则忽略
//	- \param states        64 位整数数组，输入，整型历史数值列表，
//	- 对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
//	- RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，表示相应的历史存储值；否则忽略
//	- \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
//	- \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
//	- \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致，
//	- 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
//	- 如果 datetimes、ms 标识的数据已经存在，其值将被替换。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_archived_values64_warp(rtdb_int32 handle, rtdb_int32* count, const rtdb_int32* ids, const rtdb_timestamp_type* datetimes, const rtdb_subtime_type* subtimes, const rtdb_float64* values, const rtdb_int64* states, const rtdb_int16* qualities, rtdb_error* errors)
func RawRtdbhPutArchivedValues64Warp() {}

// RawRtdbhPutArchivedCoorValues64Warp 写入批量标签点批量坐标型历史存储数据
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
func RawRtdbhPutArchivedCoorValues64Warp() {}

// RawRtdbhPutSingleDatetimeValue64Warp 写入单个datetime标签点在某一时间的历史数据
//
//	  *
//	- \param handle    连接句柄
//	- \param id        整型，输入，标签点标识
//	- \param datetime  整型，输入，数值时间列表,
//	- 表示距离1970年1月1日08:00:00的秒数
//	- \param ms        短整型，输入，历史数值时间，
//	- 对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
//	- \param blob      字节型数组，输入，历史datetime数值
//	- \param len       短整型，输入，datetime数值长度，超过一个页大小数据将被截断。
//	- \param quality   短整型，输入，历史数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
//	- \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_put_single_datetime_value64_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_timestamp_type datetime, rtdb_subtime_type subtime, const rtdb_byte* blob, rtdb_length_type len, rtdb_int16 quality)
func RawRtdbhPutSingleDatetimeValue64Warp() {}

// RawRtdbhPutArchivedBlobValues64Warp 写入批量标签点批量字符串型历史存储数据
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
func RawRtdbhPutArchivedBlobValues64Warp() {}

// RawRtdbhFlushArchivedValuesWarp 将标签点未写满的补历史缓存页写入存档文件中。
//   - \param handle        连接句柄
//   - \param id            整型，输入，标签点标识。
//   - \param count         整型，输出，缓存页中数据个数。
//   - \remark 补历史缓存页写满后会自动写入存档文件中，不满的历史缓存页也会写入文件，
//   - 但会有一个时间延迟，在此期间此段数据可能查询不到，为了及时看到补历史的结果，
//   - 应在结束补历史后调用本接口。
//   - count 参数可为空指针，对应的信息将不再返回。
//
// rtdb_error RTDBAPI_CALLRULE rtdbh_flush_archived_values_warp(rtdb_int32 handle, rtdb_int32 id, rtdb_int32 *count)
func RawRtdbhFlushArchivedValuesWarp() {}

// RawRtdbhGetSingleNamedTypeValue64Warp 读取单个自定义类型标签点某个时间的历史数据
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
func RawRtdbhGetSingleNamedTypeValue64Warp() {}

// RawRtdbhGetArchivedNamedTypeValues64Warp 连续读取自定义类型标签点的历史数据
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
func RawRtdbhGetArchivedNamedTypeValues64Warp() {}

// RawRtdbhPutSingleNamedTypeValue64Warp 写入自定义类型标签点的单个历史事件
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
func RawRtdbhPutSingleNamedTypeValue64Warp() {}

// RawRtdbhPutArchivedNamedTypeValues64Warp 批量补写自定义类型标签点的历史事件
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
func RawRtdbhPutArchivedNamedTypeValues64Warp() {}

// RawRtdbeComputeHistory64Warp 重算或补算批量计算标签点历史数据
//
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
func RawRtdbeComputeHistory64Warp() {}

// RawRtdbbGetEquationByFileNameWarp 根据文件名获取方程式
// *      [handle]   连接句柄
// *      [file_name] 输入，字符串，方程式路径
// *      [equation]  输出，返回的方程式长度最长为RTDB_MAX_EQUATION_SIZE-1
// *备注：用户调用时为equation分配的空间不得小于RTDB_MAX_EQUATION_SIZE
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_equation_by_file_name_warp(rtdb_int32 handle, const char* file_name, char equation[RTDB_MAX_EQUATION_SIZE])
func RawRtdbbGetEquationByFileNameWarp() {}

// RawRtdbbGetEquationByIdWarp 根ID径获取方程式
// * [handle]   连接句柄
// * [id]				输入，整型，方程式ID
// * [equation] 输出，返回的方程式长度最长为RTDB_MAX_EQUATION_SIZE-1
// * 备注：用户调用时为equation分配的空间不得小于RTDB_MAX_EQUATION_SIZE
// rtdb_error RTDBAPI_CALLRULE rtdbb_get_equation_by_id_warp(rtdb_int32 handle, rtdb_int32 id, char equation[RTDB_MAX_EQUATION_SIZE])
func RawRtdbbGetEquationByIdWarp() {}

// RawRtdbeGetEquationGraphCountWarp 根据标签点 id 获取相关联方程式键值对数量
//   - 参数：
//   - [handle]   连接句柄
//   - [id]       整型，输入，标签点标识
//   - [flag]     枚举，输入，获取的拓扑图的关系
//   - [count]    整型，输入，拓扑图键值对数量
//   - 备注：键值对为数据结构，存储方程式涉及到的各标签点ID、及其父ID等
//   - 具体参考rtdbe_get_equation_graph_datas
//
// rtdb_error RTDBAPI_CALLRULE rtdbe_get_equation_graph_count_warp(rtdb_int32 handle, rtdb_int32 id, RTDB_GRAPH_FLAG flag, rtdb_int32 *count)
func RawRtdbeGetEquationGraphCountWarp() {}

// RawRtdbeGetEquationGraphDatasWarp 根据标签点 id 获取相关联方程式键值对数据
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
func RawRtdbeGetEquationGraphDatasWarp() {}

// RawRtdbpGetPerfTagsCountWarp 获取Perf服务中支持的性能计数点的数量
//   - 参数：
//   - [handle]   连接句柄
//   - [count]    整型，输出，表示实际获取到的Perf服务中支持的性能计数点的数量
//
// rtdb_error RTDBAPI_CALLRULE rtdbp_get_perf_tags_count_warp(rtdb_int32 handle, int* count)
func RawRtdbpGetPerfTagsCountWarp() {}

// RawRtdbpGetPerfTagsInfoWarp 根据性能计数点ID获取相关的性能计数点信息
//   - 参数：
//   - [handle]   连接句柄
//   - [count]    整型，输入，输出
//   - 输入时，表示想要获取的性能计数点信息的数量，也表示tags_info，errors等的长度
//   - 输出时，表示实际获取到的性能计数点信息的数量
//   - [errors] 无符号整型数组，输出，获取性能计数点信息的返回值列表，参考rtdb_error.h
//   - 备注：用户须保证分配给 tags_info，errors 的空间与 count 相符
//
// rtdb_error RTDBAPI_CALLRULE rtdbp_get_perf_tags_info_warp(rtdb_int32 handle, rtdb_int32* count, RTDB_PERF_TAG_INFO* tags_info, rtdb_error* errors)
func RawRtdbpGetPerfTagsInfoWarp() {}

// RawRtdbpGetPerfValues64Warp 批量读取性能计数点的当前快照数值
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
func RawRtdbpGetPerfValues64Warp() {}
