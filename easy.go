package rtdb_api

import "errors"

type RtdbConnect struct {
	HostName         string         // 服务端名称
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
func Login(hostName string, port int32, userName string, password string) (*RtdbConnect, error) {
	rtn := RtdbConnect{
		HostName: hostName,
		Port:     port,
		UserName: userName,
		Password: password,
	}

	// 连接数据库
	cHandle, rte := RawRtdbConnectWarp(rtn.HostName, rtn.Port)
	if !errors.Is(rte, RteOk) {
		return nil, rte.GoError()
	}
	rtn.ConnectHandle = cHandle

	// 登录数据库
	priv, rte := RawRtdbLoginWarp(rtn.ConnectHandle, rtn.UserName, rtn.Password)
	if !errors.Is(rte, RteOk) {
		return nil, rte.GoError()
	}
	rtn.Priv = priv

	// 获取元信息
	infos, errs, rte := RawRtdbbGetMetaSyncInfoWarp(rtn.ConnectHandle, 0)
	if !errors.Is(rte, RteOk) {
		return nil, rte.GoError()
	}
	for _, rte := range errs {
		if !errors.Is(rte, RteOk) {
			return nil, rte.GoError()
		}
	}
	rtn.SyncInfos = infos

	// 获取套接字句柄
	for i := range infos {
		sHandle, rte := RawRtdbGetOwnConnectionWarp(rtn.ConnectHandle, int32(i+1))
		if !errors.Is(rte, RteOk) {
			return nil, rte.GoError()
		}
		rtn.SocketHandles = append(rtn.SocketHandles, sHandle)
	}

	// 获取服务器操作系统类型
	osType, rte := RawRtdbOsType(rtn.ConnectHandle)
	if !errors.Is(rte, RteOk) {
		return nil, rte.GoError()
	}
	rtn.ServerOsType = osType

	// 获取String/Blob最大长度
	maxLen, rte := RawRtdbGetMaxBlobLenWarp(rtn.ConnectHandle)
	if !errors.Is(rte, RteOk) {
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
func (c *RtdbConnect) GetClientVersion() (*ApiVersion, error) {
	version, rte := RawRtdbGetApiVersionWarp()
	if !errors.Is(rte, RteOk) {
		return nil, rte.GoError()
	}
	return &version, rte.GoError()
}

// SetClientOption 设置客户端参数
func (c *RtdbConnect) SetClientOption(option RtdbApiOption, value int32) error {
	rte := RawRtdbSetOptionWarp(option, value)
	return rte.GoError()
}
