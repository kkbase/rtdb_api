package rtdb_api

import "testing"

func TestLogin(t *testing.T) {
	Login(Hostname, Port, Username, Password)
}
