package rtdb_api

import (
	"fmt"
	"testing"
)

const Hostname = "127.0.0.1"
const Port = 6327
const Username = "sa"
const Password = "golden"

func TestRtdbGetApiVersion(t *testing.T) {
	major, minor, beta, err := RtdbGetApiVersionWarp()
	fmt.Println(major, minor, beta, err)
}

func TestRtdbConnectWarp(t *testing.T) {
	handle, err := RtdbConnectWarp(Hostname, Port)
	if err != nil {
		fmt.Println("创建连接失败", err.Error())
		return
	}
	fmt.Println(handle)
}

func TestRtdbConnectionCountWarp(t *testing.T) {
	handle, err := RtdbConnectWarp(Hostname, Port)
	if err != nil {
		t.Error("创建连接失败", err.Error())
		return
	}
	count, err := RtdbConnectionCountWarp(handle, 0)
	if err != nil {
		t.Error("获取Count失败", err.Error())
		return
	}
	fmt.Println("Connect Count: ", count)
}

func TestHello(t *testing.T) {
	fmt.Println("243")
}
