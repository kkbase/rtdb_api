package rtdb_api

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
	cHandle, err := RawRtdbConnectWarp(rtn.HostName, rtn.Port)
	if err != nil {
		return nil, err
	}
	rtn.ConnectHandle = cHandle

	// 登录数据库
	priv, err := RawRtdbLoginWarp(rtn.ConnectHandle, rtn.UserName, rtn.Password)
	if err != nil {
		return nil, err
	}
	rtn.Priv = priv

	// 获取元信息
	infos, errs, err := RawRtdbbGetMetaSyncInfoWarp(rtn.ConnectHandle, 0)
	if err != nil {
		return nil, err
	}
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	rtn.SyncInfos = infos

	// 获取套接字句柄
	for i := range infos {
		sHandle, err := RawRtdbGetOwnConnectionWarp(rtn.ConnectHandle, int32(i+1))
		if err != nil {
			return nil, err
		}
		rtn.SocketHandles = append(rtn.SocketHandles, sHandle)
	}

	// 获取服务器操作系统类型
	osType, err := RawRtdbOsType(rtn.ConnectHandle)
	if err != nil {
		return nil, err
	}
	rtn.ServerOsType = osType

	// 获取String/Blob最大长度
	maxLen, err := RawRtdbGetMaxBlobLenWarp(rtn.ConnectHandle)
	if err != nil {
		return nil, err
	}
	rtn.StringBlobMaxLen = maxLen

	return &rtn, nil
}

// Logout 登出数据库
func (c *RtdbConnect) Logout() error {
	err := RawRtdbDisconnectWarp(c.ConnectHandle)
	return err
}

// GetClientVersion 获取客户端版本
func (c *RtdbConnect) GetClientVersion() (ApiVersion, error) {
	version, err := RawRtdbGetApiVersionWarp()
	return version, err
}

// SetClientOption 设置客户端参数
func (c *RtdbConnect) SetClientOption(option RtdbApiOption, value int32) error {
	err := RawRtdbSetOptionWarp(option, value)
	return err
}
