package rtdb_api

import (
	"fmt"
	"testing"
)

func TestRtdbGetApiVersion(t *testing.T) {
	major, minor, beta, err := RtdbGetApiVersion()
	fmt.Println(major, minor, beta, err)
}
