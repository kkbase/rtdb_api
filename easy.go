package rtdb_api

import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	// MaxBlockSize 单次读取文件最大大小
	MaxBlockSize = 5 * 1024 * 1024

	// MaxFileSize 允许读取文件最大大小
	MaxFileSize = 512 * 1024 * 1024
)

// ServerOption 服务端配置
type ServerOption struct {
	IsString     bool
	StringOption ParamString
	IntOption    ParamInt
}

// NewServerOption 新建服务端类型（通过字面值新建服务端配置, 会自动推断配置类型是String或Int）
func NewServerOption(option string) ServerOption {
	intOption, err := strconv.Atoi(option)
	if err != nil {
		return ServerOption{StringOption: ParamString(option), IsString: true}
	} else {
		return ServerOption{IntOption: ParamInt(intOption), IsString: false}
	}
}

// NewStringServerOption 新建String类型服务端配置
func NewStringServerOption(option ParamString) ServerOption {
	return ServerOption{StringOption: option, IsString: true}
}

// NewIntServerOption 新建Int类型服务端配置
func NewIntServerOption(option ParamInt) ServerOption {
	return ServerOption{IntOption: option, IsString: false}
}

// GetString 获取String类型配置，如果配置为Int类型则会报错
func (o *ServerOption) GetString() (ParamString, error) {
	if o.IsString {
		return o.StringOption, nil
	} else {
		return "", errors.New("配置为Int类型")
	}
}

// GetInt 获取Int类型配置，如果配置为String类型则会报错
func (o *ServerOption) GetInt() (ParamInt, error) {
	if o.IsString {
		return 0, errors.New("配置为String类型")
	} else {
		return o.IntOption, nil
	}
}

// GetLiteralValue 获取字面值，无论是String还是Int都会转换为字符串，方便前端显示
func (o *ServerOption) GetLiteralValue() string {
	if o.IsString {
		return string(o.StringOption)
	} else {
		return strconv.Itoa(int(o.IntOption))
	}
}

// SocketInfo Socket基本信息
type SocketInfo struct {
	SocketHandle SocketHandle // Socket句柄
	IpAddr       string       // IP地址
	Port         int32        // 端口号
	JobId        int32        // 连接最近处理的任务编号
	JobTime      DateTimeType // 最近处理任务的时间
	ConnectTime  DateTimeType // 客户端连接时间
	Timeout      DateTimeType // 连接超时时间
	Client       string       // 连接的客户端主机名称
	Process      string       // 连接的客户端程序名
	User         string       // 登录的用户
}

func getSocketInfo(handle ConnectHandle, nodeNumber int32, socket SocketHandle) (*SocketInfo, error) {
	connInfo, rte := RawRtdbGetConnectionInfoIpv6Warp(handle, nodeNumber, socket)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	ipAddr := connInfo.IpAddr6
	if ipAddr == "" {
		ipAddr = fmt.Sprintf("%d.%d.%d.%d", byte(connInfo.IpAddr>>24), byte(connInfo.IpAddr>>16), byte(connInfo.IpAddr>>8), byte(connInfo.IpAddr))
	}
	timeout, rte := RawRtdbGetTimeoutWarp(handle, socket)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	info := SocketInfo{
		SocketHandle: socket,
		IpAddr:       ipAddr,
		Port:         int32(connInfo.Port),
		JobId:        connInfo.Job,
		JobTime:      connInfo.JobTime,
		ConnectTime:  connInfo.ConnectTime,
		Timeout:      timeout,
		Client:       connInfo.Client,
		Process:      connInfo.Process,
		User:         connInfo.User,
	}
	return &info, nil
}

// NamedType 自定义类型
type NamedType struct {
	Name   string              // 自定义类型名称
	Fields []RtdbDataTypeField // 字段列表
	Desc   string              // 自定义类型描述
	Length int32               // 自定义类型长度(所有字段长度的累加和)
}

// ValueType 数值类型
type ValueType string

// 基本数值类型
const (
	// ValueTypeBool 布尔类型
	ValueTypeBool = ValueType("bool")

	// ValueTypeUint8 无符号8位整数
	ValueTypeUint8 = ValueType("uint8")

	// ValueTypeInt8 有符号8位整数
	ValueTypeInt8 = ValueType("int8")

	// ValueTypeChar 单字节字符
	ValueTypeChar = ValueType("char")

	// ValueTypeUint16 无符号16位整数
	ValueTypeUint16 = ValueType("uint16")

	// ValueTypeInt16 有符号16位整数
	ValueTypeInt16 = ValueType("int16")

	// ValueTypeUint32 无符号32位整数
	ValueTypeUint32 = ValueType("uint32")

	// ValueTypeInt32 有符号32位整数
	ValueTypeInt32 = ValueType("int32")

	// ValueTypeUint64 无符号64位整数
	ValueTypeUint64 = ValueType("uint64")

	// ValueTypeInt64 有符号64位整数
	ValueTypeInt64 = ValueType("int64")

	// ValueTypeFloat16 16位浮点数
	ValueTypeFloat16 = ValueType("float16")

	// ValueTypeFloat32 32位浮点数
	ValueTypeFloat32 = ValueType("float32")

	// ValueTypeFloat64 64位浮点数
	ValueTypeFloat64 = ValueType("float64")

	// ValueTypeCoor 二维坐标
	ValueTypeCoor = ValueType("coor")

	// ValueTypeString 字符串
	ValueTypeString = ValueType("string")

	// ValueTypeBlob 数据块
	ValueTypeBlob = ValueType("blob")

	// ValueTypeDatetime 时间
	ValueTypeDatetime = ValueType("datetime")

	// ValueTypeFp16 16位定点数
	ValueTypeFp16 = ValueType("fp16")

	// ValueTypeFp32 32位定点数
	ValueTypeFp32 = ValueType("fp32")

	// ValueTypeFp64 64位定点数
	ValueTypeFp64 = ValueType("fp64")
)

func (vt ValueType) ToRawType() (RtdbType, string) {
	switch vt {
	case ValueTypeBool:
		return RtdbTypeBool, "bool"
	case ValueTypeUint8:
		return RtdbTypeUint8, "uint8"
	case ValueTypeInt8:
		return RtdbTypeInt8, "int8"
	case ValueTypeChar:
		return RtdbTypeChar, "char"
	case ValueTypeUint16:
		return RtdbTypeUint16, "uint16"
	case ValueTypeInt16:
		return RtdbTypeInt16, "int16"
	case ValueTypeUint32:
		return RtdbTypeUint32, "uint32"
	case ValueTypeInt32:
		return RtdbTypeInt32, "int32"
	case ValueTypeInt64:
		return RtdbTypeInt64, "int64"
	case ValueTypeFloat16:
		return RtdbTypeReal16, "float16"
	case ValueTypeFloat32:
		return RtdbTypeReal32, "float32"
	case ValueTypeFloat64:
		return RtdbTypeReal64, "float64"
	case ValueTypeCoor:
		return RtdbTypeCoor, "coor"
	case ValueTypeString:
		return RtdbTypeString, "string"
	case ValueTypeBlob:
		return RtdbTypeBlob, "blob"
	case ValueTypeDatetime:
		return RtdbTypeDatetime, "datetime"
	case ValueTypeFp16:
		return RtdbTypeFp16, "fp16"
	case ValueTypeFp32:
		return RtdbTypeFp32, "fp32"
	case ValueTypeFp64:
		return RtdbTypeFp64, "fp64"
	default:
		return RtdbTypeNamedT, string(vt)
	}
}

func FromRawType(typ RtdbType, namedTypeName string) ValueType {
	switch typ {
	case RtdbTypeBool:
		return ValueTypeBool
	case RtdbTypeUint8:
		return ValueTypeUint8
	case RtdbTypeInt8:
		return ValueTypeInt8
	case RtdbTypeChar:
		return ValueTypeChar
	case RtdbTypeUint16:
		return ValueTypeUint16
	case RtdbTypeInt16:
		return ValueTypeInt16
	case RtdbTypeUint32:
		return ValueTypeUint32
	case RtdbTypeInt32:
		return ValueTypeInt32
	case RtdbTypeInt64:
		return ValueTypeInt64
	case RtdbTypeReal16:
		return ValueTypeFloat16
	case RtdbTypeReal32:
		return ValueTypeFloat32
	case RtdbTypeReal64:
		return ValueTypeFloat64
	case RtdbTypeCoor:
		return ValueTypeCoor
	case RtdbTypeString:
		return ValueTypeString
	case RtdbTypeBlob:
		return ValueTypeBlob
	case RtdbTypeNamedT:
		return ValueType(namedTypeName)
	case RtdbTypeDatetime:
		return ValueTypeDatetime
	case RtdbTypeFp16:
		return ValueTypeFp16
	case RtdbTypeFp32:
		return ValueTypeFp32
	case RtdbTypeFp64:
		return ValueTypeFp64
	default:
		panic("分支不可达")
	}
}

// PointClass 点类型
type PointClass int32

const (
	// PointBase 基本点
	PointBase = PointClass(RtdbClassBase)

	// PointScan 采集点
	PointScan = PointClass(RtdbClassBase | RtdbClassScan)

	// PointCalc 计算点
	PointCalc = PointClass(RtdbClassBase | RtdbClassCalc)

	// PointScanCalc 计算采集点
	PointScanCalc = PointClass(RtdbClassBase | RtdbClassScan | RtdbClassCalc)
)

// IsScan 是否为采集点
func (pc PointClass) IsScan() bool {
	return pc&PointScan != 0
}

// IsCalc 是否为计算点
func (pc PointClass) IsCalc() bool {
	return pc&PointCalc != 0
}

// PointInfo 点属性
type PointInfo struct {
	// 核心配置
	ID        PointID       // 标签点ID
	TableID   TableID       // 当前标签点所属表ID
	Name      string        // 标签点名称
	ValueType ValueType     // 数值类型
	Class     PointClass    // 标签点类别
	Precision RtdbPrecision // 时间戳精度

	// 基本点配置
	Desc           string     // 标签点描述
	Unit           string     // 工程单位
	Archive        Switch     // 是否存档
	Digits         int16      // 数值位数
	Shutdown       Switch     // 停机状态字
	LowLimit       float32    // 量程下限
	HighLimit      float32    // 量程上限
	Step           Switch     // 是否阶跃
	Typical        float32    // 典型值
	Compress       Switch     // 是否压缩
	CompDev        float32    // 压缩偏差
	CompDevPercent float32    // 压缩偏差百分比
	CompTimeMax    int32      // 最大压缩间隔
	CompTimeMin    int32      // 最短压缩间隔
	ExcDev         float32    // 例外偏差
	ExcDevPercent  float32    // 例外偏差百分比
	ExcTimeMax     int32      // 最大例外间隔
	ExcTimeMin     int32      // 最短例外间隔
	Mirror         RtdbMirror // 镜像收发控制
	Summary        Switch     // 统计加速

	// 采集点配置，仅采集点有效
	Source     string                         // 数据源
	Scan       Switch                         // 是否采集
	Instrument string                         // 设备标签
	Locations  [RtdbConstLocationsSize]int32  // 共包含五个设备位址
	UserInts   [RtdbConstUserintSize]int32    // 共包含两个自定义整数
	UserReals  [RtdbConstUserrealSize]float32 // 共包含两个自定义单精度浮点数

	// 计算点配置, 仅计算点有效
	Equation string       // 实时方程式
	Trigger  RtdbTrigger  // 计算触发机制
	TimeCopy RtdbTimeCopy // 计算结果时间戳参考
	Period   int32        // 触发周期

	// 只读属性
	NamedType   NamedType    // 自定义类型结构, 仅自定义类型有效
	TableDotTag string       // 标签点全名，格式为“表名称.标签点名称”
	ChangeDate  DateTimeType // 标签点属性最后一次被修改的时间
	Changer     string       // 标签点属性最后一次被修改的用户名
	CreateDate  DateTimeType // 标签点被创建的时间
	Creator     string       // 标签点创建者的用户名
}

// NewPointInfo 新建标签点属性, 备注：只需填写必要属性，其他都是默认，需要时可自行设置
//
// input:
//   - name 点名
//   - tableId 表ID
//   - valueType 数值类型
//   - class 点类型
//   - precision 点时间戳精度
//   - unit 点单位
//   - desc 点描述
func NewPointInfo(name string, tableId TableID, valueType ValueType, class PointClass, precision RtdbPrecision, unit, desc string) *PointInfo {
	return &PointInfo{
		Name:           name,
		ValueType:      valueType,
		TableID:        tableId,
		Class:          class,
		Unit:           unit,
		Desc:           desc,
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
		Mirror:         RtdbMirrorPointOff,
		Summary:        OFF,
		Precision:      precision,
	}
}

// SetLimit 设置量程上下限
//
// input:
//   - lowLimit 量程上限
//   - highLimit 量程下限
//   - typical 典型值(默认值)
func (p *PointInfo) SetLimit(lowLimit float32, highLimit float32, typical float32) {
	p.LowLimit = lowLimit
	p.HighLimit = highLimit
	p.Typical = typical
}

// SetStoreDisplay 设置存储显示
//
// input:
//   - archive 是否存档
//   - digits 数值显示位数
//   - shutdown 是否停机补写
//   - step 是否阶跃
//   - mirror 镜像配置
//   - summary 是否开启统计加速
func (p *PointInfo) SetStoreDisplay(archive Switch, digits int16, shutdown Switch, step Switch, mirror RtdbMirror, summary Switch) {
	p.Archive = archive
	p.Digits = digits
	p.Shutdown = shutdown
	p.Step = step
	p.Mirror = mirror
	p.Summary = summary
}

// SetCompress 设置压缩
//
// input:
//   - compress 是否压缩
//   - compDev 压缩偏差
//   - compDevPercent 压缩偏差百分比
//   - compTimeMax 最大压缩间隔
//   - compTimeMin 最小压缩间隔
func (p *PointInfo) SetCompress(compress Switch, compDev float32, compDevPercent float32, compTimeMax int32, compTimeMin int32) {
	p.Compress = compress
	p.CompDev = compDev
	p.CompDevPercent = compDevPercent
	p.CompTimeMax = compTimeMax
	p.CompTimeMin = compTimeMin
}

// SetException 设置例外偏差
//
// input:
//   - excDev 例外偏差
//   - excDevPercent 例外偏差百分比
//   - excTimeMax 最大例外间隔
//   - excTimeMin 最短例外间隔
func (p *PointInfo) SetException(excDev float32, excDevPercent float32, excTimeMax int32, excTimeMin int32) {
	p.ExcDev = excDev
	p.ExcDevPercent = excDevPercent
	p.ExcTimeMax = excTimeMax
	p.ExcTimeMin = excTimeMin
}

// SetScan 设置采集点属性
//
// input:
//   - source 数据源
//   - scan 是否采集
//   - instrument 设备标签
//   - locations 共包含五个设备地址
//   - userInts 共包含两个自定义整数
//   - userReals 共包含两个自定义单精度浮点数
func (p *PointInfo) SetScan(
	source string, scan Switch, instrument string, locations [RtdbConstLocationsSize]int32,
	userInts [RtdbConstUserintSize]int32, userReals [RtdbConstUserrealSize]float32,
) {
	p.Class |= PointClass(RtdbClassScan)
	p.Source = source
	p.Scan = scan
	p.Instrument = instrument
	p.Locations = locations
	p.UserInts = userInts
	p.UserReals = userReals
}

// SetCalc 设置计算点
//
// input:
//   - equation 实时方程式
//   - trigger 计算触发机制
//   - timeCopy 计算结果时间戳参考
//   - period 触发周期
func (p *PointInfo) SetCalc(equation string, trigger RtdbTrigger, timeCopy RtdbTimeCopy, period int32) {
	p.Class |= PointClass(RtdbClassCalc)
	p.Equation = equation
	p.Trigger = trigger
	p.TimeCopy = timeCopy
	p.Period = period
}

// PointInfoToRaw 点信息转换为Raw点属性表
func PointInfoToRaw(info *PointInfo) (*RtdbPoint, *RtdbScan, *RtdbCalc, string) {
	rtdbType, namedTypeName := info.ValueType.ToRawType()
	milliSecond := int8(0)
	if info.Precision != RtdbPrecisionSecond {
		milliSecond = 1
	}
	base := &RtdbPoint{
		Tag:            info.Name,
		Type:           rtdbType,
		Table:          info.TableID,
		Desc:           info.Desc,
		Unit:           info.Unit,
		Archive:        info.Archive,
		Digits:         info.Digits,
		Shutdown:       info.Shutdown,
		LowLimit:       info.LowLimit,
		HighLimit:      info.HighLimit,
		Step:           info.Step,
		Typical:        info.Typical,
		Compress:       info.Compress,
		CompDev:        info.CompDev,
		CompDevPercent: info.CompDevPercent,
		CompTimeMax:    info.CompTimeMax,
		CompTimeMin:    info.CompTimeMin,
		ExcDev:         info.ExcDev,
		ExcDevPercent:  info.ExcDevPercent,
		ExcTimeMin:     info.ExcTimeMin,
		ExcTimeMax:     info.ExcTimeMax,
		Class:          RtdbClass(info.Class),
		Mirror:         info.Mirror,
		Summary:        info.Summary,
		Precision:      info.Precision,
		MilliSecond:    milliSecond,
	}
	scan := (*RtdbScan)(nil)
	if info.Class.IsScan() {
		scan = &RtdbScan{
			Source:     info.Source,
			Scan:       info.Scan,
			Instrument: info.Instrument,
			Locations:  info.Locations,
			UserInts:   info.UserInts,
			UserReals:  info.UserReals,
		}
	}
	calc := (*RtdbCalc)(nil)
	if info.Class.IsCalc() {
		calc = &RtdbCalc{
			Equation: info.Equation,
			Trigger:  info.Trigger,
			TimeCopy: info.TimeCopy,
			Period:   info.Period,
		}
	}
	return base, scan, calc, namedTypeName
}

// PointInfoFromRaw 点属性表转换为点信息
func PointInfoFromRaw(handle ConnectHandle, base *RtdbPoint, scan *RtdbScan, calc *RtdbCalc, isRecycled bool) (*PointInfo, error) {
	typ := (*NamedType)(nil)
	if base.Type == RtdbTypeNamedT {
		if !isRecycled {
			names, counts, rtes, rte := RawRtdbbGetNamedTypeNamesPropertyWarp(handle, []PointID{base.ID})
			if !RteIsOk(rte) {
				return nil, rte.GoError()
			}
			if !RteIsOk(rtes[0]) {
				return nil, rte.GoError()
			}
			fields, tLen, desc, rte := RawRtdbbGetNamedTypeWarp(handle, names[0], counts[0])
			typ = &NamedType{Name: names[0], Fields: fields, Desc: desc, Length: tLen}
		} else {
			names, counts, rtes, rte := RawRtdbbGetRecycledNamedTypeNamesPropertyWarp(handle, []PointID{base.ID})
			if !RteIsOk(rte) {
				return nil, rte.GoError()
			}
			if !RteIsOk(rtes[0]) {
				return nil, rte.GoError()
			}
			fields, tLen, desc, rte := RawRtdbbGetNamedTypeWarp(handle, names[0], counts[0])
			typ = &NamedType{Name: names[0], Fields: fields, Desc: desc, Length: tLen}
		}
	}
	namedTypeName := ""
	if typ != nil {
		namedTypeName = typ.Name
	}
	info := &PointInfo{
		ID:             base.ID,
		TableID:        base.Table,
		Name:           base.Tag,
		ValueType:      FromRawType(base.Type, namedTypeName),
		Class:          PointClass(base.Class),
		Precision:      base.Precision,
		Desc:           base.Desc,
		Unit:           base.Unit,
		Archive:        base.Archive,
		Digits:         base.Digits,
		Shutdown:       base.Shutdown,
		LowLimit:       base.LowLimit,
		HighLimit:      base.HighLimit,
		Step:           base.Step,
		Typical:        base.Typical,
		Compress:       base.Compress,
		CompDev:        base.CompDev,
		CompDevPercent: base.CompDevPercent,
		CompTimeMax:    base.CompTimeMax,
		CompTimeMin:    base.CompTimeMin,
		ExcDev:         base.ExcDev,
		ExcDevPercent:  base.ExcDevPercent,
		ExcTimeMax:     base.ExcTimeMax,
		ExcTimeMin:     base.ExcTimeMin,
		Mirror:         base.Mirror,
		Summary:        base.Summary,
		TableDotTag:    base.TableDotTag,
		ChangeDate:     base.ChangeDate,
		Changer:        base.Changer,
		CreateDate:     base.CreateDate,
		Creator:        base.Creator,
	}
	if typ != nil {
		info.NamedType = *typ
	}
	if scan != nil {
		info.Source = scan.Source
		info.Scan = scan.Scan
		info.Instrument = scan.Instrument
		info.Locations = scan.Locations
		info.UserInts = scan.UserInts
		info.UserReals = scan.UserReals
	}
	if calc != nil {
		info.Equation = calc.Equation
		info.Trigger = calc.Trigger
		info.TimeCopy = calc.TimeCopy
		info.Period = calc.Period
	}

	return info, nil
}

////////////////////////////////////////////////
//////////////////上面是一些结构//////////////////
////////////////////摆烂的分隔线/////////////////
/////////////////下面是RtdbConnect函数///////////
////////////////////////////////////////////////

type RtdbConnect struct {
	HostIp           string         // 服务端名称
	Port             int32          // 服务端端口
	UserName         string         // 用户名
	Password         string         // 密码
	ConnectHandle    ConnectHandle  // 连接句柄
	Priv             PrivGroup      // 用户权限
	SyncInfos        []RtdbSyncInfo // 元数据信息
	SocketHandles    []SocketHandle // 套接字句柄
	ServerOsType     RtdbOsType     // 服务端操作系统类型
	StringBlobMaxLen int32          // 最大支持String/Blob长度
}

// Login 登录数据库
//
// input:
//   - hostIp 主机IP
//   - port 端口
//   - userName 用户名
//   - password 密码
//
// output:
//   - RtdbConnect(conn) 返回数据库连接
func Login(hostIp string, port int32, userName string, password string) (*RtdbConnect, error) {
	rtn := RtdbConnect{
		HostIp:   hostIp,
		Port:     port,
		UserName: userName,
		Password: password,
	}

	// 连接数据库
	cHandle, rte := RawRtdbConnectWarp(rtn.HostIp, rtn.Port)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	rtn.ConnectHandle = cHandle

	// 登录数据库
	priv, rte := RawRtdbLoginWarp(rtn.ConnectHandle, rtn.UserName, rtn.Password)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	rtn.Priv = priv

	// 获取元信息
	infos, errs, rte := RawRtdbbGetMetaSyncInfoWarp(rtn.ConnectHandle, 0)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	for _, rte := range errs {
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
	}
	rtn.SyncInfos = infos

	// 获取套接字句柄
	for i := range infos {
		sHandle, rte := RawRtdbGetOwnConnectionWarp(rtn.ConnectHandle, int32(i+1))
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		rtn.SocketHandles = append(rtn.SocketHandles, sHandle)
	}

	// 获取服务器操作系统类型
	osType, rte := RawRtdbOsType(rtn.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	rtn.ServerOsType = osType

	// 获取String/Blob最大长度
	maxLen, rte := RawRtdbGetMaxBlobLenWarp(rtn.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	rtn.StringBlobMaxLen = maxLen

	return &rtn, nil
}

// Logout 登出数据库
func (c *RtdbConnect) Logout() error {
	rte := RawRtdbDisconnectWarp(c.ConnectHandle)
	return rte.GoError()
}

// GetClientVersion 获取客户端版本
//
// output:
//   - ApiVersion(version) 客户端版本
func (c *RtdbConnect) GetClientVersion() (*ApiVersion, error) {
	version, rte := RawRtdbGetApiVersionWarp()
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return &version, rte.GoError()
}

// SetClientOption 设置客户端参数
//
// input:
//   - option: 客户端参数选项
//   - value: 客户端参数值
func (c *RtdbConnect) SetClientOption(option RtdbApiOption, value int32) error {
	rte := RawRtdbSetOptionWarp(option, value)
	return rte.GoError()
}

// GetServerOption 获取服务端参数
//
// input:
//   - param 服务端参数选项
//
// output:
//   - ServerOption(option) 服务端参数值
func (c *RtdbConnect) GetServerOption(param RtdbParam) (*ServerOption, error) {
	if param.IsStringParam() {
		opt, rte := RawRtdbGetDbInfo1Warp(c.ConnectHandle, param)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		return &ServerOption{StringOption: opt, IsString: true}, nil
	} else {
		opt, rte := RawRtdbGetDbInfo2Warp(c.ConnectHandle, param)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		return &ServerOption{IntOption: opt, IsString: false}, nil
	}
}

// SetServerOption 设置服务端参数
//
// input:
//   - param 服务端参数选项
//   - option 服务端参数值
func (c *RtdbConnect) SetServerOption(param RtdbParam, option ServerOption) error {
	if param.IsStringParam() {
		strOpt, err := option.GetString()
		if err != nil {
			return err
		}
		rte := RawRtdbSetDbInfo1Warp(c.ConnectHandle, param, strOpt)
		return rte.GoError()
	} else {
		intOpt, err := option.GetInt()
		if err != nil {
			return err
		}
		rte := RawRtdbSetDbInfo2Warp(c.ConnectHandle, param, intOpt)
		return rte.GoError()
	}
}

// GetSocketInfos 获取服务端SocketInfo列表，单机服务端返回一个SocketInfo列表，双活服务端返回两个SocketInfo列表
//
// output:
//   - [][]SocketInfo(infos) Socket信息列表
func (c *RtdbConnect) GetSocketInfos() ([][]SocketInfo, error) {
	if len(c.SyncInfos) == 1 { /* 单机,返回一个Socket列表 */
		count, rte := RawRtdbConnectionCountWarp(c.ConnectHandle, 0)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		sockets, rte := RawRtdbGetConnectionsWarp(c.ConnectHandle, 0, count)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}

		infos := make([]SocketInfo, 0)
		for _, socket := range sockets {
			info, err := getSocketInfo(c.ConnectHandle, 0, socket)
			if err != nil {
				return nil, err
			}
			infos = append(infos, *info)
		}
		return [][]SocketInfo{infos}, nil
	} else { /* 双活,返回两个Socket列表 */
		count1, rte := RawRtdbConnectionCountWarp(c.ConnectHandle, 1)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		sockets1, rte := RawRtdbGetConnectionsWarp(c.ConnectHandle, 1, count1)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		infos1 := make([]SocketInfo, 0)
		for _, socket := range sockets1 {
			info, err := getSocketInfo(c.ConnectHandle, 1, socket)
			if err != nil {
				return nil, err
			}
			infos1 = append(infos1, *info)
		}

		count2, rte := RawRtdbConnectionCountWarp(c.ConnectHandle, 2)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		sockets2, rte := RawRtdbGetConnectionsWarp(c.ConnectHandle, 2, count2)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		infos2 := make([]SocketInfo, 0)
		for _, socket := range sockets2 {
			info, err := getSocketInfo(c.ConnectHandle, 2, socket)
			if err != nil {
				return nil, err
			}
			infos2 = append(infos2, *info)
		}

		return [][]SocketInfo{infos1, infos2}, nil
	}
}

// GetOwnSocketInfo 获取当前连接的SocketInfo，单机服务端返回一个SocketInfo，双活服务端返回两个SocketInfo
//
// output:
//   - []Socket Socket信息
func (c *RtdbConnect) GetOwnSocketInfo() ([]SocketInfo, error) {
	if len(c.SyncInfos) == 1 { /* 单机,返回一个Socket句柄 */
		socket, rte := RawRtdbGetOwnConnectionWarp(c.ConnectHandle, 0)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		info, err := getSocketInfo(c.ConnectHandle, 0, socket)
		if err != nil {
			return nil, err
		}
		return []SocketInfo{*info}, nil
	} else { /* 双活,返回两个Socket句柄 */
		socket1, rte := RawRtdbGetOwnConnectionWarp(c.ConnectHandle, 1)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		info1, err := getSocketInfo(c.ConnectHandle, 1, socket1)
		if err != nil {
			return nil, err
		}
		socket2, rte := RawRtdbGetOwnConnectionWarp(c.ConnectHandle, 2)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		info2, err := getSocketInfo(c.ConnectHandle, 2, socket2)
		if err != nil {
			return nil, err
		}
		return []SocketInfo{*info1, *info2}, nil
	}
}

// SetSocketTimeout 设置Socket超时时间
//
// input:
//   - info Socket信息结构
//   - timeout 超时时间
func (c *RtdbConnect) SetSocketTimeout(info SocketInfo, timeout DateTimeType) error {
	rte := RawRtdbSetTimeoutWarp(c.ConnectHandle, info.SocketHandle, timeout)
	return rte.GoError()
}

// KillSocket 断开Socket
//
// input:
//   - info Socket信息结构
func (c *RtdbConnect) KillSocket(info SocketInfo) error {
	rte := RawRtdbKillConnectionWarp(c.ConnectHandle, info.SocketHandle)
	return rte.GoError()
}

// AddIpBlackList 添加IP黑名单项
//
// input:
//   - address 阻止连接段地址
//   - mask 阻止连接段子网掩码
//   - desc 阻止连接段的说明
func (c *RtdbConnect) AddIpBlackList(address string, mask string, desc string) error {
	rte := RawRtdbAddBlacklistWarp(c.ConnectHandle, address, mask, desc)
	return rte.GoError()
}

// UpdateIpBlackList 更新连接黑名单项
//
// input:
//   - oldAddr 原黑名单地址
//   - oldMask 原黑名单掩码
//   - newAddr 新黑名单地址
//   - newMask 新黑名单掩码
//   - newDesc 新黑名单描述
func (c *RtdbConnect) UpdateIpBlackList(oldAddr string, oldMask string, newAddr string, newMask string, newDesc string) error {
	rte := RawRtdbUpdateBlacklistWarp(c.ConnectHandle, oldAddr, oldMask, newAddr, newMask, newDesc)
	return rte.GoError()
}

// DeleteIpBlackList 删除连接黑名单项
//
// input:
//   - addr 黑名单地址
//   - mask 黑名单掩码
func (c *RtdbConnect) DeleteIpBlackList(addr string, mask string) error {
	rte := RawRtdbRemoveBlacklistWarp(c.ConnectHandle, addr, mask)
	return rte.GoError()
}

// GetIpBlackLists 获得连接黑名单列表
//
// output:
//   - []BlackList(lists) 连接黑名单列表
func (c *RtdbConnect) GetIpBlackLists() ([]BlackList, error) {
	lists, rte := RawRtdbGetBlacklistWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return lists, nil
}

// AddIpWhiteList 添加连接白名单
//
// input:
//   - addr 连接白名单地址
//   - mask 连接白名单掩码
//   - desc 连接白名单描述
//   - priv 连接白名单权限
func (c *RtdbConnect) AddIpWhiteList(addr string, mask string, desc string, priv PrivGroup) error {
	rte := RawRtdbAddAuthorizationWarp(c.ConnectHandle, addr, mask, desc, priv)
	return rte.GoError()
}

// UpdateIpWhiteList 更新连接白名单
//
// input:
//   - oldAddr 原连接白名单地址
//   - oldMask 原连接白名单掩码
//   - newAddr 新连接白名单地址
//   - newMask 新连接白名单掩码
//   - newDesc 新连接白名单描述
//   - newPriv 新连接白名单权限
func (c *RtdbConnect) UpdateIpWhiteList(oldAddr string, oldMask string, newAddr string, newMask string, newDesc string, newPriv PrivGroup) error {
	rte := RawRtdbUpdateAuthorizationWarp(c.ConnectHandle, oldAddr, oldMask, newAddr, newMask, newDesc, newPriv)
	return rte.GoError()
}

// DeleteIpWhiteList 删除白名单
//
// input:
//   - addr 连接白名单地址
//   - mask 连接白名单掩码
func (c *RtdbConnect) DeleteIpWhiteList(addr string, mask string) error {
	rte := RawRtdbRemoveAuthorizationWarp(c.ConnectHandle, addr, mask)
	return rte.GoError()
}

// GetIpWhiteLists 获取连接白名单列表
//
// output:
//   - []AuthorizationsList(lists)
func (c *RtdbConnect) GetIpWhiteLists() ([]AuthorizationsList, error) {
	lists, rte := RawRtdbGetAuthorizationsWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return lists, nil
}

// UpdatePassword 修改用户密码
//
// input:
//   - user 用户名
//   - password 用户密码
func (c *RtdbConnect) UpdatePassword(user string, password string) error {
	rte := RawRtdbChangePasswordWarp(c.ConnectHandle, user, password)
	return rte.GoError()
}

// UpdateOwnPassword 修改自己的密码
//
// input:
//   - oldPwd 旧密码
//   - newPwd 新密码
func (c *RtdbConnect) UpdateOwnPassword(oldPwd string, newPwd string) error {
	rte := RawRtdbChangeMyPasswordWarp(c.ConnectHandle, oldPwd, newPwd)
	return rte.GoError()
}

// GetPriv 获取连接权限
//
// output:
//   - PrivGroup(priv) 用户权限
func (c *RtdbConnect) GetPriv() (*PrivGroup, error) {
	priv, rte := RawRtdbGetPrivWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return &priv, nil
}

// SetPriv 设置连接权限
//
// input:
//   - user 用户名
//   - priv 用户权限
func (c *RtdbConnect) SetPriv(user string, priv PrivGroup) error {
	rte := RawRtdbChangePrivWarp(c.ConnectHandle, user, priv)
	if RteIsOk(rte) && c.UserName == user {
		c.Priv = priv
	}
	return rte.GoError()
}

// AddUser 添加用户
//
// input:
//   - user 用户名
//   - password 用户密码
//   - priv 用户权限
func (c *RtdbConnect) AddUser(user string, password string, priv PrivGroup) error {
	rte := RawRtdbAddUserWarp(c.ConnectHandle, user, password, priv)
	return rte.GoError()
}

// DeleteUser 删除用户
//
// input:
//   - user 用户名
func (c *RtdbConnect) DeleteUser(user string) error {
	rte := RawRtdbRemoveUserWarp(c.ConnectHandle, user)
	return rte.GoError()
}

// LockUser 锁定用户
//
// input:
//   - user 用户名
//   - lock 是否锁定
func (c *RtdbConnect) LockUser(user string, lock Switch) error {
	rte := RawRtdbLockUserWarp(c.ConnectHandle, user, lock)
	return rte.GoError()
}

// GetUsers 获取用户列表
//
// output:
//   - []RtdbUserInfo(users) 用户列表
func (c *RtdbConnect) GetUsers() ([]RtdbUserInfo, error) {
	users, rte := RawRtdbGetUsersWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return users, nil
}

// AddNamedType 创建自定义类型
//
// input:
//   - name 自定义类型名称
//   - fields 自定义类型字段列表
//   - desc 自定义类型描述
func (c *RtdbConnect) AddNamedType(name string, desc string, fields ...RtdbDataTypeField) error {
	rte := RawRtdbbCreateNamedTypeWarp(c.ConnectHandle, name, desc, fields...)
	return rte.GoError()
}

// DeleteNamedType 删除自定义类型
//
// input:
//   - name 自定义类型的名称
func (c *RtdbConnect) DeleteNamedType(name string) error {
	rte := RawRtdbbRemoveNamedTypeWarp(c.ConnectHandle, name)
	return rte.GoError()
}

// GetNamedType 获取自定义类型
//
// output:
//   - NamedType 自定义类型
func (c *RtdbConnect) GetNamedType(name string) (*NamedType, error) {
	types, err := c.GetNamedTypes()
	if err != nil {
		return nil, err
	}

	for _, typ := range types {
		if typ.Name == name {
			return &typ, nil
		}
	}

	return nil, errors.New("未知自定义类型")
}

// GetNamedTypes 获取自定义类型列表
//
// output:
//   - []NamedType(types) 自定义类型列表
func (c *RtdbConnect) GetNamedTypes() ([]NamedType, error) {
	count, rte := RawRtdbbGetNamedTypesCountWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	names, fieldCounts, rte := RawRtdbbGetAllNamedTypesWarp(c.ConnectHandle, count)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}

	types := make([]NamedType, count)
	for i := 0; i < len(names); i++ {
		fields, length, desc, rte := RawRtdbbGetNamedTypeWarp(c.ConnectHandle, names[i], fieldCounts[i])
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		types = append(types, NamedType{
			Name:   names[i],
			Fields: fields,
			Length: length,
			Desc:   desc,
		})
	}
	return types, nil
}

// UpdateNamedType 修改自定义类型
//
// input:
//   - name 自定义类型的名称
//   - modifyName 要修改的 自定义类型名称
//   - modifyDesc 要修改的 自定义类型的描述
//   - modifyFields 要修改的 字段名称<->字段描述
func (c *RtdbConnect) UpdateNamedType(name string, modifyName *string, modifyDesc *string, modifyFields map[string]string) error {
	fieldNames := make([]string, 0)
	fieldDescs := make([]string, 0)
	for name, desc := range modifyFields {
		fieldNames = append(fieldNames, name)
		fieldDescs = append(fieldDescs, desc)
	}
	rte := RawRtdbbModifyNamedTypeWarp(c.ConnectHandle, name, modifyName, modifyDesc, fieldNames, fieldDescs)
	return rte.GoError()
}

// ServerHostTime 服务端主机时间
func (c *RtdbConnect) ServerHostTime() (*time.Time, error) {
	datetime, rte := RawRtdbHostTime64Warp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	hostTime := time.Unix(int64(datetime), 0)
	return &hostTime, nil
}

// DurationToString 时间段转字符串, 这个是服务端的时间段字符串格式，和通用时间段字符串有区别, 具体如下：
//
//	?y    ?年, 1年 = 365日
//	?m    ?月, 1月 = 30 日
//	?d    ?日
//	?h    ?小时
//	?n    ?分钟
//	?s    ?秒
//
// input:
//   - duration 时间段
//
// output:
//   - string(字符串格式时间段)
func (c *RtdbConnect) DurationToString(duration time.Duration) (string, error) {
	durationStr, rte := RawRtdbFormatTimespanWarp(int32(duration.Seconds()))
	if !RteIsOk(rte) {
		return "", rte.GoError()
	}
	return durationStr, nil
}

// StringToDuration 字符串转时间段, 这个是服务端的时间段字符串格式，和通用时间段字符串有区别, 具体如下：
//
//	?y    ?年, 1年 = 365日
//	?m    ?月, 1月 = 30 日
//	?d    ?日
//	?h    ?小时
//	?n    ?分钟
//	?s    ?秒
//
// input:
//   - strDuration 字符串类型时间段
//
// output:
//   - time.Duration(duration) 时间段
func (c *RtdbConnect) StringToDuration(strDuration string) (time.Duration, error) {
	duration, rte := RawRtdbParseTimespanWarp(strDuration)
	if !RteIsOk(rte) {
		return 0, rte.GoError()
	}
	return time.Second * time.Duration(duration), nil
}

// StringToTime 字符串转时间戳
//
//	其中 base_time 表示基本时间，有三种形式：
//	1. 时间字符串，如 "2010-1-1" 及 "2010-1-1 8:00:00"；
//	2. 时间代码，表示客户端相对时间；
//	可用的时间代码及含义如下：
//	td             当天零点
//	yd             昨天零点
//	tm             明天零点
//	mon            本周一零点
//	tue            本周二零点
//	wed            本周三零点
//	thu            本周四零点
//	fri            本周五零点
//	sat            本周六零点
//	sun            本周日零点
//	3. 星号('*')，表示客户端当前时间。
//	offset_time 是可选的，可以出现多次，
//	可用的时间偏移代码及含义如下：
//	[+|-] ?y            偏移?年, 1年 = 365日
//	[+|-] ?m            偏移?月, 1月 = 30 日
//	[+|-] ?d            偏移?日
//	[+|-] ?h            偏移?小时
//	[+|-] ?n            偏移?分钟
//	[+|-] ?s            偏移?秒
//	[+|-] ?ms           偏移?毫秒
//	例如："*-1d" 表示当前时刻减去24小时。
//
// input:
//   - strTime 字符串类型时间戳
//
// output:
//   - time.Time(timestamp) 时间戳
func (c *RtdbConnect) StringToTime(strTime string) (*time.Time, error) {
	datetime, subtime, rte := RawRtdbParseTimeWarp(strTime)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	goTime := time.Unix(int64(datetime), int64(subtime))
	return &goTime, nil
}

// GetDriveLetterList 获取盘符列表, windows平台是C、D、E、F这些盘符，linux平台是 / 盘符
//
// output:
//   - []string(litters) 盘符列表
func (c *RtdbConnect) GetDriveLetterList() ([]string, error) {
	letters, rte := RawRtdbGetLogicalDriversWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return letters, nil
}

// GetDirItemList 获取目录项列表
//
// input:
//   - dir 目录路径
//
// output:
//   - []DirItem(items) 目录项列表
func (c *RtdbConnect) GetDirItemList(dir string) ([]DirItem, error) {
	rte := RawRtdbOpenPathWarp(c.ConnectHandle, dir)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	defer func() {
		_ = RawRtdbClosePathWarp(c.ConnectHandle)
	}()

	items := make([]DirItem, 0)
	for {
		item, rte := RawRtdbReadPath64Warp(c.ConnectHandle)
		if !RteIsOk(rte) {
			if errors.Is(rte, RteBatchEnd) {
				break
			} else {
				return nil, rte.GoError()
			}
		}
		items = append(items, item)
	}
	return items, nil
}

// CreateDir 创建目录
//
// input:
//   - path 目录路径
func (c *RtdbConnect) CreateDir(path string) error {
	rte := RawRtdbMkdirWarp(c.ConnectHandle, path)
	return rte.GoError()
}

// ReadFile 读取文件
//
// input:
//   - path 文件路径
//
// output:
//   - []byte(data) 文件内容
func (c *RtdbConnect) ReadFile(path string) ([]byte, error) {
	size, rte := RawRtdbGetFileSizeWarp(c.ConnectHandle, path)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	if size > MaxFileSize {
		return nil, errors.New("当前文件大小超出允许读取长度")
	}

	buf := bytes.NewBuffer(nil)
	for i := 0; i < int(size); i += MaxBlockSize {
		data, rte := RawRtdbReadFileWarp(c.ConnectHandle, path, int64(i*MaxBlockSize), MaxBlockSize)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		_, err := buf.Write(data)
		if err != nil {
			return nil, fmt.Errorf("写入缓存失败：%v", err)
		}
	}
	return buf.Bytes(), nil
}

// CreateTable 创建表
//
// input:
//   - name 表名
//   - desc 表描述
//
// output:
//   - RtdbTable(table) 返回表
func (c *RtdbConnect) CreateTable(name string, desc string) (*RtdbTable, error) {
	table, rte := RawRtdbbAppendTableWarp(c.ConnectHandle, name, desc)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return &table, nil
}

// DeleteTable 删除表
//
// input:
//   - id 表ID
func (c *RtdbConnect) DeleteTable(id TableID) error {
	rte := RawRtdbbRemoveTableByIdWarp(c.ConnectHandle, id)
	return rte.GoError()
}

func (c *RtdbConnect) GetTable(id TableID) (*RtdbTable, error) {
	table, rte := RawRtdbbGetTablePropertyByIdWarp(c.ConnectHandle, id)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return &table, nil
}

// GetTables 获取表列表
//
// output:
//   - []RtdbTable(tables) 表列表
func (c *RtdbConnect) GetTables() ([]RtdbTable, error) {
	count, rte := RawRtdbbTablesCountWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	ids, rte := RawRtdbbGetTablesWarp(c.ConnectHandle, count)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	tables := make([]RtdbTable, 0)
	for _, id := range ids {
		table, rte := RawRtdbbGetTablePropertyByIdWarp(c.ConnectHandle, id)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		tables = append(tables, table)
	}
	return tables, nil
}

// UpdateTableName 更新表名
// input:
//   - id 表ID
//   - name 表名
func (c *RtdbConnect) UpdateTableName(id TableID, name string) error {
	rte := RawRtdbbUpdateTableNameWarp(c.ConnectHandle, id, name)
	return rte.GoError()
}

// UpdateTableDesc 更新表描述
//
// input:
//   - id 表ID
//   - desc 表描述
func (c *RtdbConnect) UpdateTableDesc(id TableID, desc string) error {
	rte := RawRtdbbUpdateTableDescByIdWarp(c.ConnectHandle, id, desc)
	return rte.GoError()
}

// CreatePoint 创建点
//
// input:
//   - info 输入点信息
//
// output:
//   - PointInfo(info) 输出点信息
func (c *RtdbConnect) CreatePoint(info *PointInfo) (*PointInfo, error) {
	base, scan, calc, tName := PointInfoToRaw(info)
	if base.Type == RtdbTypeNamedT {
		if tName == "" {
			return nil, errors.New("点数值类型为RtdbTypeNamedT, 此时NamedTypeName不能为空")
		}
		_, err := c.GetNamedType(tName)
		if err != nil {
			return nil, err
		}
		base, scan, rte := RawRtdbbInsertNamedTypePointWarp(c.ConnectHandle, base, scan, tName)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		return PointInfoFromRaw(c.ConnectHandle, base, scan, nil, false)
	} else {
		base, scan, calc, rte := RawRtdbbInsertMaxPointWarp(c.ConnectHandle, base, scan, calc)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		return PointInfoFromRaw(c.ConnectHandle, base, scan, calc, false)
	}
}

// DeletePoint 删除点
//
// input:
//   - id 点ID
func (c *RtdbConnect) DeletePoint(id PointID) error {
	rte := RawRtdbbRemovePointByIdWarp(c.ConnectHandle, id)
	return rte.GoError()
}

// GetPoints 批量获取标签点
//
// input:
//   - ids 标签点列表
//
// output:
//   - []PointInfo(infos) 标签点属性列表
func (c *RtdbConnect) GetPoints(ids []PointID) ([]*PointInfo, []error, error) {
	bases, scans, calcs, rtes, rte := RawRtdbbGetMaxPointsPropertyWarp(c.ConnectHandle, ids)
	if !RteIsOk(rte) {
		return nil, nil, rte.GoError()
	}
	errs := RtdbErrorListToErrorList(rtes)
	infos := make([]*PointInfo, 0)
	for i := 0; i < len(ids); i++ {
		info, err := PointInfoFromRaw(c.ConnectHandle, &bases[i], &scans[i], &calcs[i], false)
		if err != nil {
			errs[i] = err
		}
		infos = append(infos, info)
	}
	return infos, errs, nil
}

// GetPoint 获取点
//
// input:
//   - id 点ID
//
// output:
//   - PointInfo(info) 返回点信息
func (c *RtdbConnect) GetPoint(id PointID) (*PointInfo, error) {
	bases, scans, calcs, rtes, rte := RawRtdbbGetMaxPointsPropertyWarp(c.ConnectHandle, []PointID{id})
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	for _, rte := range rtes {
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
	}
	return PointInfoFromRaw(c.ConnectHandle, &bases[0], &scans[0], &calcs[0], false)
}
