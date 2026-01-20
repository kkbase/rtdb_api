package rtdb_api

type RtdbConnect struct {
	HostName         string        // 服务端名称
	Port             int32         // 服务端端口
	UserName         string        // 用户名
	Password         string        // 密码
	ConnectHandle    ConnectHandle // 连接句柄
	SocketHandle     SocketHandle  // 套接字句柄
	ServerOsType     RtdbOsType    // 服务端操作系统类型
	StringBlobMaxLen int32         // 最大支持String/Blob长度
	Priv             PrivGroup     // 用户权限
}

// Login 登录数据库
func Login(hostName string, port int32, userName string, password string) (*RtdbConnect, error) {
	rtn := RtdbConnect{
		HostName: hostName,
		Port:     port,
		UserName: userName,
		Password: password,
	}
	cHandle, err := RawRtdbConnectWarp(rtn.HostName, rtn.Port)
	if err != nil {
		return nil, err
	}
	rtn.ConnectHandle = cHandle

	priv, err := RawRtdbLoginWarp(rtn.ConnectHandle, rtn.UserName, rtn.Password)
	if err != nil {
		return nil, err
	}
	rtn.Priv = priv

	return &rtn, nil
}

// Logout 登出数据库
func (c *RtdbConnect) Logout() error {
	err := RawRtdbDisconnectWarp(c.ConnectHandle)
	return err
}
