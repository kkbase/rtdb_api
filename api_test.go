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
		t.Error("创建连接失败", err.Error())
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
		t.Error("创建连接失败", err.Error())
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
		t.Error("获取Count失败", err.Error())
		return
	}
	fmt.Println("当前服务器连接个数: ", count)
}

func TestRawRtdbGetDbInfo1Warp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err.Error())
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
		t.Error("创建连接失败", err.Error())
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
		t.Error("创建连接失败", err.Error())
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
		t.Error("创建连接失败", err.Error())
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
		t.Error("创建连接失败", err.Error())
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
		t.Error("创建连接失败", err.Error())
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
		t.Error("创建连接失败", err.Error())
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
