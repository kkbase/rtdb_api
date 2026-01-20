package rtdb_api

import "testing"

func TestLogin(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()
}
