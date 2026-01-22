package rtdb_api

import (
	"fmt"
	"testing"
	"time"
)

// 用户登录/登出
func TestLoginLogout(t *testing.T) {
	// 登录
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	// 登出
	defer func() { _ = conn.Logout() }()

	fmt.Println(conn.SyncInfos, conn.StringBlobMaxLen)
}

// 获取客户端版本
func TestRtdbConnect_GetClientVersion(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 获取客户端版本
	version, err := conn.GetClientVersion()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(version)
}

// 设置客户端选项
func TestRtdbConnect_SetClientOption(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 设置客户端选项
	err = conn.SetClientOption(RtdbApiOptionAutoReconn, 0)
	if err != nil {
		t.Error(err)
	}
}

// 服务端选项
func TestRtdbConnect_GetSetServerOption(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 获取服务端选项
	opt, err := conn.GetServerOption(RtdbParamLockedPagesMem)
	if err != nil {
		t.Error("获取服务端选项失败", err)
		return
	}
	fmt.Println(opt.GetLiteralValue())

	// 设置服务端选项
	err = conn.SetServerOption(RtdbParamLockedPagesMem, *opt)
	if err != nil {
		t.Error("设置服务端选项失败", err)
		return
	}
}

// 获取当前用户的SocketInfo，获取所有用户的SocketInfo
func TestRtdbConnect_GetSocketInfo(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 获取自己的Socket信息
	ownInfo, err := conn.GetOwnSocketInfo()
	if err != nil {
		t.Error("获取自己的SocketInfo失败", err)
		return
	}
	fmt.Println(ownInfo)

	// 设置Socket超时时间
	err = conn.SetSocketTimeout(ownInfo[0], ownInfo[0].Timeout)
	if err != nil {
		t.Error("设置timeout失败", err)
		return
	}

	// 获取全部Socket信息
	allInfos, err := conn.GetSocketInfos()
	if err != nil {
		t.Error("获取所有SocketInfo失败", err)
		return
	}
	fmt.Println(allInfos)
}

// IP黑名单
func TestRtdbConnect_BlackList(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 添加黑名单
	err = conn.AddIpBlackList("192.168.123.123", "255.255.255.0", "add 123")
	if err != nil {
		t.Error("添加黑名单失败：", err)
		return
	}

	// 修改黑名单
	err = conn.UpdateIpBlackList("192.168.123.123", "255.255.255.0", "192.168.123.123", "255.255.255.0", "update 123")
	if err != nil {
		t.Error("修改黑名单失败：", err)
		return
	}

	// 获取黑名单
	bLists, err := conn.GetIpBlackLists()
	if err != nil {
		t.Error("获取黑名单失败：", err)
		return
	}
	bOk := false
	for _, b := range bLists {
		if b.Desc == "update 123" {
			bOk = true
			break
		}
	}
	if !bOk {
		t.Error("修改黑名单失败")
		return
	}

	err = conn.DeleteIpBlackList("192.168.123.123", "255.255.255.0")
	if err != nil {
		t.Error("删除黑名单失败：", err)
		return
	}
}

// IP白名单
func TestRtdbConnect_WhiteList(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 添加白名单
	err = conn.AddIpWhiteList("192.168.123.120", "255.255.255.0", "add 120", PrivGroupRtdbSA)
	if err != nil {
		t.Error("添加白名单失败：", err)
		return
	}

	// 修改白名单
	err = conn.UpdateIpWhiteList("192.168.123.120", "255.255.255.0", "192.168.123.120", "255.255.255.0", "update 120", PrivGroupRtdbSA)
	if err != nil {
		t.Error("修改白名单失败：", err)
		return
	}

	// 获取白名单
	wLists, err := conn.GetIpWhiteLists()
	if err != nil {
		t.Error("获取白名单失败：", err)
		return
	}
	wOk := false
	for _, w := range wLists {
		if w.Desc == "update 120" {
			wOk = true
			break
		}
	}
	if !wOk {
		t.Error("修改白名单失败")
		return
	}

	err = conn.DeleteIpWhiteList("192.168.123.120", "255.255.255.0")
	if err != nil {
		t.Error("删除白名单失败：", err)
		return
	}
}

// 用户
func TestRtdbConnect_User(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 添加用户
	err = conn.AddUser("test111", "122333", PrivGroupRtdbSA)
	if err != nil {
		t.Error("添加用户失败: ", err)
		return
	}

	// 修改用户密码
	err = conn.UpdatePassword("test111", "123123")
	if err != nil {
		t.Error("修改密码失败: ", err)
		return
	}

	// 验证密码是否修改成功
	conn2, err := Login(Hostname, Port, "test111", "123123")
	if err != nil {
		t.Error("登录用户失败", err)
		return
	}
	defer func() { _ = conn2.Logout() }()

	// 修改自己的密码
	err = conn2.UpdateOwnPassword("123123", "122333")
	if err != nil {
		t.Error("修改自己的密码失败：", err)
		return
	}

	// 获取连接权限
	priv, err := conn2.GetPriv()
	if err != nil {
		t.Error("获取权限失败：", err)
		return
	}
	if *priv != PrivGroupRtdbSA {
		t.Error("验证权限失败")
		return
	}

	// 设置连接权限
	err = conn2.SetPriv("test111", PrivGroupRtdbRO)
	if err != nil {
		t.Error("设置权限失败：", err)
		return
	}

	// 锁定用户
	err = conn.LockUser("test111", OFF)
	if err != nil {
		t.Error("锁定用户失败：", err)
		return
	}

	// 用户列表
	users, err := conn.GetUsers()
	if err != nil {
		t.Error("获取用户列表失败：", err)
		return
	}
	uOk := false
	// 验证
	for _, u := range users {
		if u.User == "test111" {
			uOk = true
			break
		}
	}
	if !uOk {
		t.Error("用户列表中不存在test111")
		return
	}

	// 删除用户
	err = conn.DeleteUser("test111")
	if err != nil {
		t.Error("删除用户失败：", err)
	}
}

// 自定义类型
func TestRtdbConnect_NamedType(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 创建自定义类型
	err = conn.AddNamedType(
		"abc",
		"abc desc",
		RtdbDataTypeField{
			Name:   "A",
			Type:   RtdbTypeReal64,
			Length: 0,
			Desc:   "A desc",
		}, RtdbDataTypeField{
			Name:   "B",
			Type:   RtdbTypeReal64,
			Length: 0,
			Desc:   "B desc",
		}, RtdbDataTypeField{
			Name:   "C",
			Type:   RtdbTypeReal64,
			Length: 0,
			Desc:   "C desc",
		})
	if err != nil {
		t.Error("添加自定义类型失败")
		return
	}

	// 删除自定义类型
	defer func() {
		err := conn.DeleteNamedType("abc")
		if err != nil {
			t.Error("删除自定义类型失败")
			return
		}
	}()

	// 获取自定义类型
	types, err := conn.GetNamedTypes()
	if err != nil {
		t.Error("获取列表失败")
		return
	}
	fmt.Println(types)

	// 更新自定义类型
	desc := "up abc desc"
	err = conn.UpdateNamedType("abc", nil, &desc, map[string]string{"A": "A up", "B": "B up", "C": "C up"})
	if err != nil {
		t.Error("更新列表失败")
		return
	}

	// 获取自定义类型
	typ, err := conn.GetNamedType("abc")
	if err != nil {
		t.Error("获取列表失败")
		return
	}
	fmt.Println(typ)
}

// 时间
func TestRtdbConnect_Time(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 获取服务端主机时间
	hostTime, err := conn.ServerHostTime()
	if err != nil {
		t.Error("获取服务端时间失败：", err)
		return
	}
	fmt.Println(hostTime)

	// 时间段转字符串
	dStr, err := conn.DurationToString(time.Second * 60)
	if err != nil {
		t.Error("时间段转换失败：", err)
		return
	}
	if dStr != "1n" {
		t.Error("不为1n")
	}
}
