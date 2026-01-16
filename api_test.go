package rtdb_api

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

const Hostname = "127.0.0.1"
const Port = 6327
const Username = "sa"
const Password = "golden"

func TestNULL(t *testing.T) {}

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

	err = RawRtdbLockUserWarp(handle, "t1", OFF)
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

func TestTime(t *testing.T) {
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

	// hTime, err := RawRtdbHostTimeWarp(handle)
	// if err != nil {
	// 	t.Error("获取系统时间: ", err)
	// 	return
	// }
	// fmt.Println(hTime)

	hTime2, err := RawRtdbHostTime64Warp(handle)
	if err != nil {
		t.Error("获取系统时间：", err)
		return
	}
	fmt.Println(hTime2)

	tStr, err := RawRtdbFormatTimespanWarp(10)
	if err != nil {
		t.Error("获取跨度时间:", err)
		return
	}
	fmt.Println(tStr)

	d, err := RawRtdbParseTimespanWarp("2n")
	if err != nil {
		t.Error("获取跨度时间:", err)
		return
	}
	fmt.Println(d)

	ts, ms, err := RawRtdbParseTimeWarp("2010-1-1 8:00:00")
	if err != nil {
		t.Error("解析时间失败:", err)
		return
	}
	fmt.Println(ts, ms)
}

func TestRawRtdbFormatMessageWarp(t *testing.T) {
	name, message := RawRtdbFormatMessageWarp(RteCantLoadBase)
	fmt.Println(name, message)
}

func TestRawRtdbJobMessageWarp(t *testing.T) {
	name, desc := RawRtdbJobMessageWarp(1)
	fmt.Println(name, desc)
}

func TestRawRtdbGetSetTimeoutWarp(t *testing.T) {
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

	timeout, err := RawRtdbGetTimeoutWarp(handle, socket)
	fmt.Println(timeout)

	err = RawRtdbSetTimeoutWarp(handle, socket, timeout)
	if err != nil {
		t.Error("设置超时时间: ", err)
		return
	}
}

func TestRawRtdbKillConnectionWarp(t *testing.T) {
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

	err = RawRtdbKillConnectionWarp(handle, socket)
	if err != nil {
		t.Error("Kill套接字失败: ", err)
		return
	}
}

func TestRawRtdbGetLogicalDriversWarp(t *testing.T) {
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

	ds, err := RawRtdbGetLogicalDriversWarp(handle)
	if err != nil {
		t.Error("获取盘符：", err)
		return
	}
	fmt.Println(ds)
}

func TestDir(t *testing.T) {
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

	err = RawRtdbOpenPathWarp(handle, "/")
	if err != nil {
		t.Error("打开目录失败：", err)
		return
	}
	defer func() { _ = RawRtdbClosePathWarp(handle) }()

	for {
		dir, err := RawRtdbReadPath64Warp(handle)
		if err != nil {
			if errors.Is(err, RteBatchEnd) {
				break
			} else {
				t.Error("读取目录失败：", err)
				return
			}
		}
		fmt.Println(dir)
	}

	err = RawRtdbMkdirWarp(handle, "/tttAAA")
	if err != nil {
		t.Error("创建目录失败", err)
		return
	}

	s, err := RawRtdbGetFileSizeWarp(handle, "/opt/rtdb/bin/gstart")
	if err != nil {
		t.Error("获取size失败：", err)
		return
	}
	fmt.Println(s)

	data, err := RawRtdbReadFileWarp(handle, "/opt/rtdb/bin/gstart.ini.example", 0, 1024)
	if err != nil {
		t.Error("读取文件失败：", err)
		return
	}
	fmt.Println(string(data))
}

func TestRawRtdbGetMaxBlobLenWarp(t *testing.T) {
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

	l, err := RawRtdbGetMaxBlobLenWarp(handle)
	if err != nil {
		t.Error("获取Blob最大长度失败：", err)
		return
	}
	fmt.Println(l)
}

func TestRawRtdbFormatQualityWarp(t *testing.T) {
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

	ds, err := RawRtdbFormatQualityWarp(handle, []Quality{1, 2, 3})
	if err != nil {
		t.Error("获取质量码说明失败: ", err)
		return
	}
	fmt.Println(ds)
}

func TestRawRtdbJudgeConnectStatusWarp(t *testing.T) {
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

	err = RawRtdbJudgeConnectStatusWarp(handle)
	if err != nil {
		t.Error("获取连接状态失败:", err)
		return
	}
}

func TestRawRtdbFormatIpaddrWarp(t *testing.T) {
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

	addr := RawRtdbFormatIpaddrWarp(0x11221122)
	fmt.Println(addr)
}

func TestRawRtdbbAppendTableWarp(t *testing.T) {
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

	table, err := RawRtdbbAppendTableWarp(handle, "aaa", "aaa test")
	if err != nil {
		t.Error("添加表失败: ", err)
		return
	}

	tt, err := RawRtdbbGetTablePropertyByIdWarp(handle, table.ID)
	if err != nil {
		t.Error("获取表失败：", err)
		return
	}
	fmt.Println("表信息1：", tt)

	tt2, err := RawRtdbbGetTablePropertyByNameWarp(handle, table.Name)
	if err != nil {
		t.Error("获取表失败：", err)
		return
	}
	fmt.Println("表信息2：", tt2)

	err = RawRtdbbRemoveTableByIdWarp(handle, table.ID)
	if err != nil {
		t.Error("删除表失败：", err)
		return
	}

	table2, err := RawRtdbbAppendTableWarp(handle, "aaa", "aaa test")
	if err != nil {
		t.Error("添加表失败: ", err)
		return
	}

	count, err := RawRtdbbTablesCountWarp(handle)
	if err != nil {
		t.Error("获取表数量失败：", err)
		return
	}
	fmt.Println("表数量：", count)

	ids, err := RawRtdbbGetTablesWarp(handle)
	if err != nil {
		t.Error("获取表ID列表", err)
		return
	}
	fmt.Println("表ID列表: ", ids)

	err = RawRtdbbRemoveTableByNameWarp(handle, table2.Name)
	if err != nil {
		t.Error("删除表失败:", err)
		return
	}
}

func TestPointCount(t *testing.T) {
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

	table, err := RawRtdbbAppendTableWarp(handle, "aaa", "aaa test")
	if err != nil {
		t.Error("添加表失败: ", err)
		return
	}

	defer func() {
		err := RawRtdbbRemoveTableByIdWarp(handle, table.ID)
		if err != nil {
			t.Error("删除表失败：", err)
			return
		}
	}()

	pCount, err := RawRtdbbGetTableSizeByIdWarp(handle, table.ID)
	if err != nil {
		t.Error("获取Table中Point数量失败：", err)
		return
	}
	fmt.Println("point count by id: ", pCount)

	pCount2, err := RawRtdbbGetTableSizeByNameWarp(handle, table.Name)
	if err != nil {
		t.Error("获取Table中Point数量失败：", err)
		return
	}
	fmt.Println("point count by name: ", pCount2)

	pCount3, err := RawRtdbbGetTableRealSizeByIdWarp(handle, table.ID)
	if err != nil {
		t.Error("获取Table中Point数量失败：", err)
		return
	}
	fmt.Println("point count by name:", pCount3)
}

func TestAddPoint(t *testing.T) {
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

	table, err := RawRtdbbAppendTableWarp(handle, "aaa", "aaa test")
	if err != nil {
		t.Error("添加表失败: ", err)
		return
	}

	defer func() {
		err := RawRtdbbRemoveTableByIdWarp(handle, table.ID)
		if err != nil {
			t.Error("删除表失败：", err)
			return
		}
	}()

	base := NewDefaultPoint("ttt", RtdbTypeInt32, table.ID, RtdbClassBase, RtdbPrecisionMicro)
	base, scan, calc, err := RawRtdbbInsertMaxPointWarp(handle, base, nil, nil)
	if err != nil {
		t.Error("添加标签点失败：", err)
		return
	}
	fmt.Println("base: ", base)
	fmt.Println("scan: ", scan)
	fmt.Println("calc: ", calc)

	err = RawRtdbbRemovePointByIdWarp(handle, base.ID)
	if err != nil {
		t.Error("删除标签点失败: ", err)
		return
	}

	base, scan, calc, err = RawRtdbbInsertMaxPointWarp(handle, base, nil, nil)
	if err != nil {
		t.Error("添加标签点失败：", err)
		return
	}
	fmt.Println("base: ", base)
	fmt.Println("scan: ", scan)
	fmt.Println("calc: ", calc)

	bases, scans, calcs, errs, err := RawRtdbbGetMaxPointsPropertyWarp(handle, []PointID{base.ID})
	fmt.Println("bases: ", bases)
	fmt.Println("scans: ", scans)
	fmt.Println("calcs: ", calcs)
	fmt.Println("errs: ", errs)
	if err != nil {
		t.Error("批量获取标签点失败：", err)
		return
	}

	ids, err := RawRtdbbSearchInBatchesWarp(handle, 0, "", "", "", "", "", "", RtdbSortFlagDescend)
	if err != nil {
		t.Error("搜索标签点失败：", err)
		return
	}
	fmt.Println("搜索标签点：", ids)

	ids2, err := RawRtdbbSearchExWarp(handle, 1024, "", "", "", "", "", "", "", RtdbTypeAny, RtdbPrecisionAny, RtdbSearchAny, "", RtdbSortFlagDescend)
	if err != nil {
		t.Error("搜索标签点失败2：", err)
		return
	}
	fmt.Println("搜索标签点2：", ids2)

	count, err := RawRtdbbSearchPointsCountWarp(handle, "", "", "", "", "", "", "", RtdbTypeAny, RtdbPrecisionAny, RtdbSearchAny, "")
	if err != nil {
		t.Error("搜索标签点失败3：", err)
		return
	}
	fmt.Println("数量：", count)

	defer func() {
		err = RawRtdbbRemovePointByNameWarp(handle, table.Name+"."+base.Tag)
		if err != nil {
			t.Error("删除标签点失败: ", err)
			return
		}
	}()

	table2, err := RawRtdbbAppendTableWarp(handle, "bbb", "bbb test")
	if err != nil {
		t.Error("添加表失败: ", err)
		return
	}

	defer func() {
		err := RawRtdbbRemoveTableByIdWarp(handle, table2.ID)
		if err != nil {
			t.Error("删除表失败：", err)
			return
		}
	}()

	// err = RawRtdbbMovePointByIdWarp(handle, base.ID, table.Name)
	// if err != nil {
	// 	t.Error("移动标签点失败：", err)
	// 	return
	// }
	// time.Sleep(1 * time.Second)
}

func TestPoint2(t *testing.T) {
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

	// 清理表
	_ = RawRtdbbRemoveTableByNameWarp(handle, "aaa")
	time.Sleep(time.Second)

	table, err := RawRtdbbAppendTableWarp(handle, "aaa", "aaa test")
	if err != nil {
		t.Error("添加表失败: ", err)
		return
	}

	defer func() {
		err := RawRtdbbRemoveTableByIdWarp(handle, table.ID)
		if err != nil {
			t.Error("删除表失败：", err)
			return
		}
	}()

	base := NewDefaultPoint("ttt", RtdbTypeInt32, table.ID, RtdbClassBase, RtdbPrecisionMicro)
	base, scan, calc, err := RawRtdbbInsertMaxPointWarp(handle, base, nil, nil)
	if err != nil {
		t.Error("添加标签点失败：", err)
		return
	}
	fmt.Println("base: ", base)
	fmt.Println("scan: ", scan)
	fmt.Println("calc: ", calc)

	defer func() {
		err = RawRtdbbRemovePointByIdWarp(handle, base.ID)
		if err != nil {
			t.Error("删除标签点失败: ", err)
			return
		}
	}()

	ids, types, classes, precisions, errs, err := RawRtdbbFindPointsExWarp(handle, []string{"aaa.ttt"})
	if err != nil {
		t.Error("查找标签点失败", err)
		return
	}
	fmt.Println(ids, types, classes, precisions, errs)
}

func TestPerf(t *testing.T) {
	code := `
  PFT_CPU_USAGE_OF_LOGGER,            //!< 日志服务CPU使用
  PFT_MEM_BYTES_OF_LOGGER,            //!< 日志服务内存
  PFT_VMEM_BYTES_OF_LOGGER,           //!< 日志服务虚拟内存
  PFT_READ_BYTES_OF_LOGGER,           //!< 日志服务 I/O 读取字节
  PFT_WRITE_BYTES_OF_LOGGER,          //!< 日志服务 I/O 写入字节
  PFT_CPU_USAGE_OF_HISTORIAN,         //!< 历史数据服务CPU使用
  PFT_MEM_BYTES_OF_HISTORIAN,         //!< 历史数据服务内存
  PFT_VMEM_BYTES_OF_HISTORIAN,        //!< 历史数据服务虚拟内存
  PFT_READ_BYTES_OF_HISTORIAN,        //!< 历史数据服务 I/O 读取字节
  PFT_WRITE_BYTES_OF_HISTORIAN,       //!< 历史数据服务 I/O 写入字节
  PFT_CPU_USAGE_OF_SNAPSHOT,          //!< 快照数据服务CPU使用
  PFT_MEM_BYTES_OF_SNAPSHOT,          //!< 快照数据服务内存
  PFT_VMEM_BYTES_OF_SNAPSHOT,         //!< 快照数据服务虚拟内存
  PFT_READ_BYTES_OF_SNAPSHOT,         //!< 快照数据服务 I/O 读取字节
  PFT_WRITE_BYTES_OF_SNAPSHOT,        //!< 快照数据服务 I/O 写入字节
  PFT_CPU_USAGE_OF_EQUATION,          //!< 实时方程式服务CPU使用
  PFT_MEM_BYTES_OF_EQUATION,          //!< 实时方程式服务内存
  PFT_VMEM_BYTES_OF_EQUATION,         //!< 实时方程式服务虚拟内存
  PFT_READ_BYTES_OF_EQUATION,         //!< 实时方程式服务 I/O 读取字节
  PFT_WRITE_BYTES_OF_EQUATION,        //!< 实时方程式服务 I/O 写入字节
  PFT_CPU_USAGE_OF_BASE,              //!< 标签点信息服务CPU使用
  PFT_MEM_BYTES_OF_BASE,              //!< 标签点信息服务内存
  PFT_VMEM_BYTES_OF_BASE,             //!< 标签点信息服务虚拟内存
  PFT_READ_BYTES_OF_BASE,             //!< 标签点信息服务 I/O 读取字节
  PFT_WRITE_BYTES_OF_BASE,            //!< 标签点信息服务 I/O 写入字节
  PFT_CPU_USAGE_OF_SERVER,            //!< 网络服务CPU使用
  PFT_MEM_BYTES_OF_SERVER,            //!< 网络服务内存
  PFT_VMEM_BYTES_OF_SERVER,           //!< 网络服务虚拟内存
  PFT_READ_BYTES_OF_SERVER,           //!< 网络服务 I/O 读取字节
  PFT_WRITE_BYTES_OF_SERVER,          //!< 网络服务 I/O 写入字节
  PFT_ARV_ASYNC_QUEUE,                //!< 历史数据队列地址
  PFT_ARV_ASYNC_QUEUE_USAGE,          //!< 历史数据队列使用率
  PFT_ARVEX_ASYNC_QUEUE,              //!< 补历史数据队列地址
  PFT_ARVEX_ASYNC_QUEUE_USAGE,        //!< 补历史数据队列使用率
  PFT_EVENTS_INPUT_RATE,              //!< 普通事件入库速度（KB/秒）
  PFT_EVENTS_OUTPUT_RATE,             //!< 普通事件归档速度（KB/秒）
  PFT_FILL_IN_INPUT_RATE,             //!< 补历史事件入库速度（KB/秒）
  PFT_FILL_IN_OUTPUT_RATE,            //!< 补历史事件归档速度（KB/秒）
  PFT_ARV_CACHE_USAGE,                //!< 历史数据缓存使用率
  PFT_ARVEX_CACHE_USAGE,              //!< 补历史数据缓存使用率
  PFT_MIRROR_SNAPSHOTS_QUEUE,         //!< 快照数据的镜像队列地址
  PFT_MIRROR_SNAPSHOTS_QUEUE_USAGE,   //!< 快照数据的镜像队列使用率
  PFT_ARVEX_BLOB_ASYNC_QUEUE,         //!< str、blob补历史数据队列地址
  PFT_ARVEX_BLOB_ASYNC_QUEUE_USAGE,   //!< str、blob补历史数据队列使用率
  PFT_ARVEX_BLOB_CACHE_USAGE,         //!< str、blob补历史数据缓存使用率
  PFT_MIRROR_BUFFER_SIZE,             //!< 快照数据的镜像缓存文件
  PFT_CLUTTER_POOL_USAGE,             //!< 消息交换池利用率
  PFT_MAX_BLOCK_IN_CLUTTER_POOL,      //!< 消息交换池的最大可用额度
  PFT_ARV_ARCHIVED_TIME,              //!< 历史数据归档耗时
  PFT_ARVEX_ARCHIVED_TIME,            //!< 补历史数据归档耗时
  PFT_ARVEX_BLOB_ARCHIVED_TIME,       //!< str、blob补历史数据归档耗时
  PFT_ARV_ARCHIVED_PAGE_COUNT,        //!< 历史数据归档的数据页数量
  PFT_ARVEX_ARCHIVED_PAGE_COUNT,      //!< 补历史数据归档的数据页数量
  PFT_ARVEX_BLOB_ARCHIVED_PAGE_COUNT, //!< str、blob补历史数据归档的数据页数量
  PFT_MIRROR_ARV_VALUES_QUEUE,        //!< 补写历史数据的镜像队列地址
  PFT_MIRROR_ARV_VALUES_QUEUE_USAGE,  //!< 补写历史数据的镜像队列使用率
  PFT_MIRROR_ARV_BUFFER_SIZE,         //!< 补写历史数据的镜像缓存文件
  PFT_ARV_WRITE_COUNT,                //!< 历史数据归档写磁盘次数
  PFT_ARV_READ_COUNT,                 //!< 历史数据归档读磁盘次数
  PFT_ARV_WRITE_TIME,                 //!< 历史数据归档写磁盘时间
  PFT_ARV_READ_TIME,                  //!< 历史数据归档读磁盘时间
  PFT_ARV_INDEX_WRITE_COUNT,          //!< 历史数据归档写索引次数
  PFT_ARV_INDEX_READ_COUNT,           //!< 历史数据归档读索引次数
  PFT_ARV_INDEX_WRITE_TIME,           //!< 历史数据归档写索引时间
  PFT_ARV_INDEX_READ_TIME,            //!< 历史数据归档读索引时间
  PFT_ARV_ARC_LIST_LOCK_TIME,         //!< 历史数据归档列表锁时间
  PFT_ARV_ARC_LOCK_TIME,              //!< 历史数据归档文件锁时间
  PFT_ARV_INDEX_LOCK_TIME,            //!< 历史数据归档索引锁时间
  PFT_ARV_TOTAL_LOCK_TIME,            //!< 历史数据归档锁总时间
  PFT_ARV_WRITE_SIZE,                 //!< 历史数据归档写磁盘数据量
  PFT_ARV_READ_SIZE,                  //!< 历史数据归档读磁盘数据量
  PFT_ARV_WRITE_REAL_SIZE,            //!< 历史数据归档写磁盘有效数据量
  PFT_ARV_READ_REAL_SIZE,             //!< 历史数据归档读磁盘有效数据量
  PFT_ARVEX_WRITE_COUNT,              //!< 补历史数据归档写磁盘次数
  PFT_ARVEX_READ_COUNT,               //!< 补历史数据归档读磁盘次数
  PFT_ARVEX_WRITE_TIME,               //!< 补历史数据归档写磁盘时间
  PFT_ARVEX_READ_TIME,                //!< 补历史数据归档读磁盘时间
  PFT_ARVEX_INDEX_WRITE_COUNT,        //!< 补历史数据归档写索引次数
  PFT_ARVEX_INDEX_READ_COUNT,         //!< 补历史数据归档读索引次数
  PFT_ARVEX_INDEX_WRITE_TIME,         //!< 补历史数据归档写索引时间
  PFT_ARVEX_INDEX_READ_TIME,          //!< 补历史数据归档读索引时间
  PFT_ARVEX_ARC_LIST_LOCK_TIME,       //!< 补历史数据归档列表锁时间
  PFT_ARVEX_ARC_LOCK_TIME,            //!< 补历史数据归档文件锁时间
  PFT_ARVEX_INDEX_LOCK_TIME,          //!< 补历史数据归档索引锁时间
  PFT_ARVEX_TOTAL_LOCK_TIME,          //!< 补历史数据归档锁总时间
  PFT_ARVEX_WRITE_SIZE,               //!< 补历史数据归档写磁盘数据量
  PFT_ARVEX_READ_SIZE,                //!< 补历史数据归档读磁盘数据量
  PFT_ARVEX_WRITE_REAL_SIZE,          //!< 补历史数据归档写磁盘有效数据量
  PFT_ARVEX_READ_REAL_SIZE,           //!< 补历史数据归档读磁盘有效数据量
  PFT_PLOT_POOL_POINT_COUNT,          //!< 曲线缓存标签点数量
  PFT_PLOT_POOL_WEIGHTED_POINT_COUNT, //!< 曲线缓存权重点数量
  PFT_PLOT_POOL_TOTAL_MEM_SIZE,       //!< 曲线缓存总内存数
  PFT_PLOT_POOL_CACHED_HIT_PERCENT,   //!< 曲线缓存命中率
  PFT_OS_CPU_USAGE,              //!< 数据库所在操作系统的CPU使用率
  PFT_OS_MEM_SIZE,                    //!< 数据库所在操作系统的物理内存大小，单位MB
  PFT_OS_MEM_USAGE,                   //!< 数据库所在操作系统的物理内存使用率
  PFT_QUERY_POOL_WAIT_TASKS_SIZE,     //!< 查询线程池中等待执行的任务数
  PFT_MIRROR_ENQUEUE,                 //!< 镜像每秒入队的数量，单位字节
  PFT_MIRROR_OUTQUEUE,                //!< 镜像每秒出对的数量，单位字节
  PFT_MIRROR_SEND_CPRS,               //!< 镜像每秒压缩的数量，单位字节
  PFT_MIRROR_RECV_CPRS,               //!< 镜像每秒收到的压缩数量，单位字节
  PFT_MIRROR_RECV_UNCPRS,             //!< 镜像每秒解压缩的数量，单位字节
  PFT_MIRROR_CPRS_SPAN,               //!< 镜像报文每秒的压缩耗时总和，单位毫秒
  PFT_MIRROR_COMPRESS_RATE,			  //!< 1秒内镜像的压缩率
  PFT_API_CPRS_RATE,                  //!< API报文压缩率
  PFT_SERVER_CPRS_RATE,               //!< Server报文压缩率
  PFT_TAG_SUBSCRIBE_CUSTOMER_COUNT,                 //!< 标签点信息订阅客户端数量
  PFT_TAG_SUBSCRIBE_SEND_EVENT_COUNT,               //!< 标签点信息订阅发送事件数量
  PFT_SNAP_SUBSCRIBE_CUSTOMER_COUNT,                //!< 快照信息订阅客户端数量
  PFT_SNAP_SUBSCRIBE_SEND_EVENT_COUNT,              //!< 快照信息订阅发送事件数量
  PFT_SNAP_SUBSCRIBE_POINT_COUNT,                   //!< 快照信息订阅标签点数量
  PFT_CONNECT_SUBSCRIBE_CUSTOMER_COUNT,             //!< API监视订阅客户端数量
  PFT_CONNECT_SUBSCRIBE_SEND_EVENT_COUNT,           //!< API监视订阅发送事件数量
  PFT_NAMED_TYPE_CREATE_SUBSCRIBE_CUSTOMER_COUNT,   //!< 创建自定义类型订阅客户端数量
  PFT_NAMED_TYPE_CREATE_SEND_EVENT_COUNT,           //!< 创建自定义类型订阅发送事件数量
  PFT_NAMED_TYPE_REMOVE_SUBSCRIBE_CUSTOMER_COUNT,   //!< 删除自定义类型订阅客户端数量
  PFT_NAMED_TYPE_REMOVE_SEND_EVENT_COUNT,           //!< 删除自定义类型订阅发送事件数量
  PFT_DOUBLE_ACTIVE_SYNC_SEND_COUNT,                //!< 双活同步每秒同步发送的数据量
  PFT_DOUBLE_ACTIVE_SYNC_RECEIVE_COUNT,             //!< 双活同步每秒同步接授的数据量
  PFT_IS_RECEIVING_NORMAL_DATA_FROM_PEER,           //!< 双活系统正在接收普通类型数据
  PFT_IS_RECEIVING_BLOB_DATA_FROM_PEER,             //!< 双活系统正在接收blob类型数据
  PFT_REPLICATOR_BUFFER_BLOCK_COUNT,                //!< 双活本地的同步历史缓存还有多少数据块
  PFT_REPLICATOR_EX_BUFFER_BLOCK_COUNT,             //!< 双活本地的同步补历史缓存还有多少数据块
  PFT_REPLICATOR_BLOB_BUFFER_BLOCK_COUNT,           //!< 双活本地的同步历史缓存(blob string 数据)还有多少数据块
  PFT_REPLICATOR_BLOB_EX_BUFFER_BLOCK_COUNT,        //!< 双活本地的同步补历史缓存(blob string 数据)还有多少数据块
  PFT_SNAPSHOT_PUT_RATE,                            //!< 每秒写入快照记录数，单位 条
  PFT_SNAPSHOT_GET_RATE,                            //!< 每秒读取快照记录数，单位 条
  PFT_HISTORIAN_PUT_RATE,                           //!< 每秒写入历史记录数，单位 条
  PFT_HISTORIAN_GET_RATE,                           //!< 每秒读取历史记录数，单位 条
  PFT_HISTORIAN_WRITE_RECORD_COUNT,                 //!< 每秒写入历史数据块数
  PFT_HISTORIAN_READ_RECORD_COUNT,                  //!< 每秒读取历史数据块数
  PFT_SERVER_NETWORK_READ_BYTES,                    //!< 网络服务网络 IO 每秒读取字节数
  PFT_SERVER_NETWORK_WRITE_BYTES,                   //!< 网络服务网络 IO 每秒写入字节数
  PFT_END,                            //!< 信息数量
`

	lines := strings.Split(code, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		sp1 := strings.Split(line, ",")
		cName := strings.TrimSpace(sp1[0])

		sp2 := strings.Split(line, "//!<")
		cDesc := strings.TrimSpace(sp2[1])
		goName := SnakeToPascal(cName)
		fmt.Printf("// %s %s\n", goName, cDesc)
		fmt.Printf("%s = RtdbPerfTagID(C.%s)\n\n", goName, cName)
	}
}

func SnakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		p = strings.ToLower(p)
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}
