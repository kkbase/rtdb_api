package rtdb_api

import "testing"

// 测试用户登录/登出
func TestLoginLogout(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()
}
