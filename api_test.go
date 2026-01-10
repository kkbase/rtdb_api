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

func TestRawRtdbConnectWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		fmt.Println("创建连接失败", err.Error())
		return
	}
	fmt.Println(handle)
}

func TestRawRtdbConnectionCountWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err.Error())
		return
	}
	count, err := RawRtdbConnectionCountWarp(handle, 0)
	if err != nil {
		t.Error("获取Count失败", err.Error())
		return
	}
	fmt.Println("Connect Count: ", count)
}
