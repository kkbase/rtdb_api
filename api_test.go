package rtdb_api

import (
	"fmt"
	"testing"
)

const Hostname = "127.0.0.1"
const Port = 6327
const Username = "sa"
const Password = "golden"

// 获取 API库 版本
func TestRawRtdbGetApiVersion(t *testing.T) {
	apiVersion, err := RawRtdbGetApiVersionWarp()
	if err != nil {
		t.Error("获取版本号失败:", err)
	}
	fmt.Println("库版本号:", apiVersion)
}

// 设置 API库 基本参数
func TestRawRtdbSetOptionWarp(t *testing.T) {
	err := RawRtdbSetOptionWarp(RtdbApiOptionAutoReconn, 0)
	if err != nil {
		t.Error("设置 API库 基本参数失败：", err)
	}
}

// 测试登录、登出，涉及到3个 原始API
// - RawRtdbConnectWarp
// - RawRtdbLoginWarp
// - RawRtdbDisconnectWarp
// 创建 API库 和 数据库 之间的连接
// 使用 用户名、密码 登录数据库
// 断开 API库 和 数据库 之间的连接
func TestLoginAndLogout(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	priv, err := RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	fmt.Println("登录权限：", priv.Desc())
	err = RawRtdbDisconnectWarp(handle)
	if err != nil {
		t.Error("断开链接失败：", err)
		return
	}
}

// 测试服务器连接个数
func TestRawRtdbConnectionCountWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	count, err := RawRtdbConnectionCountWarp(handle, 0)
	if err != nil {
		t.Error("获取Count失败", err)
		return
	}
	fmt.Println("当前服务器连接个数: ", count)
}

func TestRawRtdbGetDbInfo1Warp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	param, err := RawRtdbGetDbInfo1Warp(handle, RtdbParamTableFile)
	if err != nil {
		t.Error("获取Str参数失败", err)
		return
	}
	fmt.Println(param)
}

func TestRawRtdbGetDbInfo2Warp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	param, err := RawRtdbGetDbInfo2Warp(handle, RtdbParamLicScanCount)
	if err != nil {
		t.Error("获取Str参数失败", err)
		return
	}
	fmt.Println(param)
}

func TestRawRtdbSetDbInfo1Warp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	// TODO
	// err = RawRtdbSetDbInfo1Warp(handle, RtdbParamAutoBackupPath, "/tmp/rtdb")
	// if err != nil {
	// 	t.Error("设置Str参数失败")
	// 	return
	// }
}

func TestRawRtdbSetDbInfo2Warp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	err = RawRtdbSetDbInfo2Warp(handle, RtdbParamArchiveIncreaseSize, 256)
	if err != nil {
		t.Error("设置Int参数失败")
		return
	}
}

func TestRawRtdbGetConnectionsWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	sockets, err := RawRtdbGetConnectionsWarp(handle, 0)
	if err != nil {
		t.Error("获取连接失败：", err)
		return
	}
	fmt.Println(sockets)
}

func TestRawRtdbGetOwnConnectionWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	socket, err := RawRtdbGetOwnConnectionWarp(handle, 0)
	if err != nil {
		t.Error("获取连接失败：", err)
	}
	fmt.Println(socket)
}

func TestRawRtdbGetConnectionInfoIpv6Warp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	socket, err := RawRtdbGetOwnConnectionWarp(handle, 0)
	if err != nil {
		t.Error("获取连接失败：", err)
		return
	}
	fmt.Println(socket)

	infoV6, err := RawRtdbGetConnectionInfoIpv6Warp(handle, 0, socket)
	if err != nil {
		t.Error("获取v6连接信息失败: ", err)
	}
	fmt.Println(infoV6)
}

func TestRawRtdbOsType(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	socket, err := RawRtdbGetOwnConnectionWarp(handle, 0)
	if err != nil {
		t.Error("获取连接失败：", err)
		return
	}
	fmt.Println(socket)

	osType, err := RawRtdbOsType(handle)
	if err != nil {
		t.Error("获取操作系统失败:", err)
		return
	}
	fmt.Println(osType.Desc())
}

func TestRawRtdbChangePasswordWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	err = RawRtdbChangePasswordWarp(handle, "sa", "golden")
	if err != nil {
		t.Error("修改密码失败：", err)
		return
	}
}

func TestRawRtdbChangeMyPasswordWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	err = RawRtdbChangeMyPasswordWarp(handle, "golden", "golden")
	if err != nil {
		t.Error("修改密码失败：", err)
		return
	}
}

func TestRawRtdbGetPrivWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	priv, err := RawRtdbGetPrivWarp(handle)
	if err != nil {
		t.Error("获取权限失败：", err)
		return
	}
	fmt.Println(priv.Desc())
}

func TestRawRawRtdbChangePrivWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	err = RawRtdbChangePrivWarp(handle, "sa", PrivGroupRtdbSA)
	if err != nil {
		t.Error("获取权限失败：", err)
		return
	}
}

func TestRawRtdbAddDelUserWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	err = RawRtdbAddUserWarp(handle, "t1", "golden", PrivGroupRtdbSA)
	if err != nil {
		t.Error("添加用户失败:", err)
		return
	}

	err = RawRtdbRemoveUserWarp(handle, "t1")
	if err != nil {
		t.Error("删除用户失败:", err)
		return
	}
}

func TestRawRtdbLockUserWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	err = RawRtdbAddUserWarp(handle, "t1", "golden", PrivGroupRtdbSA)
	if err != nil {
		t.Error("添加用户失败:", err)
		return
	}

	err = RawRtdbLockUserWarp(handle, "t1", true)
	if err != nil {
		t.Error("启用User失败：", err)
		return
	}

	err = RawRtdbRemoveUserWarp(handle, "t1")
	if err != nil {
		t.Error("删除用户失败:", err)
		return
	}
}

func TestRawRtdbGetUsersWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	users, err := RawRtdbGetUsersWarp(handle)
	if err != nil {
		t.Error("获取用户列表失败：", err)
		return
	}
	fmt.Println(users)
}

func TestBlackList(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	// 添加&查看
	err = RawRtdbAddBlacklistWarp(handle, "192.168.10.11", "255.255.255.0", "test desc")
	if err != nil {
		t.Error("添加黑名单失败：", err)
		return
	}
	bList, err := RawRtdbGetBlacklistWarp(handle)
	if err != nil {
		t.Error("获取黑名单失败：", err)
		return
	}
	fmt.Println(bList)

	// 修改&查看
	err = RawRtdbUpdateBlacklistWarp(handle, "192.168.10.11", "255.255.255.0", "192.168.10.11", "255.255.255.0", "test update")
	if err != nil {
		t.Error("更新黑名单失败：", err)
		return
	}
	bList, err = RawRtdbGetBlacklistWarp(handle)
	if err != nil {
		t.Error("获取黑名单失败：", err)
		return
	}
	fmt.Println(bList)

	// 删除&查看
	err = RawRtdbRemoveBlacklistWarp(handle, "192.168.10.11", "255.255.255.0")
	if err != nil {
		t.Error("删除黑名单:", err)
		return
	}
	bList, err = RawRtdbGetBlacklistWarp(handle)
	if err != nil {
		t.Error("获取黑名单失败：", err)
		return
	}
	fmt.Println(bList)
}

func TestAuthorizations(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err)
		return
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if err != nil {
		t.Error("登录失败:", err)
		return
	}
	defer func() { _ = RawRtdbDisconnectWarp(handle) }()

	// 添加&查看
	err = RawRtdbAddAuthorizationWarp(handle, "192.168.12.12", "255.255.255.0", "test desc", PrivGroupRtdbSA)
	if err != nil {
		t.Error("添加白名单失败", err)
		return
	}
	aList, err := RawRtdbGetAuthorizationsWarp(handle)
	if err != nil {
		t.Error("获取白名单失败", err)
		return
	}
	fmt.Println(aList)

	// 修改&查看
	err = RawRtdbUpdateAuthorizationWarp(handle, "192.168.12.12", "255.255.255.0", "192.168.12.12", "255.255.255.0", "test update", PrivGroupRtdbSA)
	if err != nil {
		t.Error("修改白名单失败", err)
		return
	}
	aList, err = RawRtdbGetAuthorizationsWarp(handle)
	if err != nil {
		t.Error("获取白名单失败", err)
		return
	}
	fmt.Println(aList)

	// 删除&查看
	err = RawRtdbRemoveAuthorizationWarp(handle, "192.168.12.12", "255.255.255.0")
	if err != nil {
		t.Error("删除白名单失败", err)
		return
	}
	aList, err = RawRtdbGetAuthorizationsWarp(handle)
	if err != nil {
		t.Error("获取白名单失败", err)
		return
	}
	fmt.Println(aList)
}
