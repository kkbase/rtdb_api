package rtdb_api

import (
	"fmt"
	"testing"
)

const Hostname = "127.0.0.1"
const Port = 6327
const Username = "sa"
const Password = "golden"

func TestRawRtdbGetApiVersion(t *testing.T) {
	apiVersion, err := RawRtdbGetApiVersionWarp()
	if err != nil {
		t.Error("获取版本号失败:", err)
	}
	fmt.Println("库版本号:", apiVersion)
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
