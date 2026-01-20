package rtdb_api

type RtdbConnect struct {
	HostName         string        // 服务端名称
	Port             int16         // 服务端端口
	UserName         string        // 用户名
	Password         string        // 密码
	ConnectHandle    ConnectHandle // 连接句柄
	SocketHandle     SocketHandle  // 套接字句柄
	ServerOsType     RtdbOsType    // 服务端操作系统类型
	StringBlobMaxLen int32         // 最大支持String/Blob长度
	Priv             PrivGroup     // 用户权限
}

// Login 登录数据库
func Login(hostName string, Port int16, userName string, password string) *RtdbConnect {
	return &RtdbConnect{}
}

// Logout 登出数据库
func (c *RtdbConnect) Logout() {}
