package rtdb_api

// #cgo CFLAGS: -DPNG_DEBUG=1 -I./cinclude
// #cgo CXXFLAGS: -std=c++11
// #include <stdlib.h>
// #include "api.h"
import "C"
import (
	_ "embed"
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

type RtdbError uint32

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

// RtdbSetOptionWarp 配置 api 行为参数，参见枚举 \ref RTDB_API_OPTION
// \param [in] type  选项类别
// \param [in] value 选项值
// \return rtdb_error
// \remark 选项设置后在下一次调用 api 时才生效
func RtdbSetOptionWarp(optionType RtdbApiOption, value int32) RtdbError {
	err := C.rtdb_set_option_warp(C.rtdb_int32(optionType), C.rtdb_int32(value))
	return RtdbError(err)
}

type DatagramHandle struct {
	handle C.rtdb_datagram_handle
}

// RtdbCreateDatagramHandleWarp 创建数据流
// * \param [in] in 端口
// * \param [out] remotehost 对端地址
// * \param [out] handle 数据流句柄
// * \return rtdb_error
// * \remark 创建数据流 (备注：C代码中没文档，Go这边补的)
func RtdbCreateDatagramHandleWarp(port int32, remoteHost string) (DatagramHandle, RtdbError) {
	var handle C.rtdb_datagram_handle
	cRemoteHost := C.CString(remoteHost)
	defer C.free(unsafe.Pointer(cRemoteHost))
	err := C.rtdb_create_datagram_handle_warp(C.rtdb_int32(port), cRemoteHost, &handle)
	return DatagramHandle{handle: handle}, RtdbError(err)
}

// RtdbRemoveDatagramHandleWarp 删除数据流
// * \param [in] handle 数据流句柄
// * \return rtdb_error
// * \remark 删除数据流 (备注：C代码中没文档，Go这边补的)
func RtdbRemoveDatagramHandleWarp(handle DatagramHandle) RtdbError {
	err := C.rtdb_remove_datagram_handle_warp(handle.handle)
	return RtdbError(err)
}

// RtdbRecvDatagramWarp 接收数据流
// * \param [in] cacheLen 缓存大小(会分配一个[]byte)
// * \param [in] handle 数据流句柄
// * \param [in] remote_addr 对端地址
// * \param [in] timeout 超时时间
// * \return rtdb_error
// * \remark 接收数据流 (备注：C代码中没文档，Go这边补的)
func RtdbRecvDatagramWarp(cacheLen int32, handle DatagramHandle, remoteAddr string, timeout int32) ([]byte, RtdbError) {
	message := make([]byte, cacheLen)
	messageLen := C.rtdb_int32(cacheLen)
	cRemoteAddr := C.CString(remoteAddr)
	defer C.free(unsafe.Pointer(cRemoteAddr))
	err := C.rtdb_recv_datagram_warp((*C.char)(unsafe.Pointer(&message[0])), &messageLen, handle.handle, cRemoteAddr, C.rtdb_int32(timeout))
	return message[0:messageLen], RtdbError(err)
}
