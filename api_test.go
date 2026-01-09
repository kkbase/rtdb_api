package rtdb_api

import (
	"fmt"
	"testing"
)

func TestRtdbGetApiVersion(t *testing.T) {
	major, minor, beta, err := RtdbGetApiVersionWarp()
	fmt.Println(major, minor, beta, err)
}

func TestRtdbConnectWarp(t *testing.T) {
	handle, err := RtdbConnectWarp("localhost", 6327)
	if !err.IsOk() {
		fmt.Println("创建连接失败", err.Error())
		return
	}
	fmt.Println(handle)
}
