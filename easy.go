package rtdb_api

import "C"
import (
	"errors"
	"fmt"
	"strconv"
	"time"
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
func (c *RtdbConnect) DurationToString(duration time.Duration) (string, error) {
	durationStr, rte := RawRtdbFormatTimespanWarp(int32(duration.Seconds()))
	if !RteIsOk(rte) {
		return "", rte.GoError()
	}
	return durationStr, nil
}
