package rtdb_api

import (
	"fmt"
	"testing"
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
