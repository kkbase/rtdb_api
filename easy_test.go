package rtdb_api

import (
	"fmt"
	"testing"
)

// 测试用户登录/登出
func TestLoginLogout(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	fmt.Println(conn.SyncInfos, conn.StringBlobMaxLen)
}

func TestRtdbConnect_GetClientVersion(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	version, err := conn.GetClientVersion()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(version)
}
