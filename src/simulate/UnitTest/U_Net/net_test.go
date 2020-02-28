package U_Net

import (
	"common/Log"
	"net"
	"testing"
)

// parse web addr to normal ip.
func Test_1(t *testing.T) {
	addrs, err := net.LookupHost("www.baidu.com")
	if err != nil {
		Log.FmtPrintf("Err: %s", err.Error())
		return
	}

	for _, addr := range addrs {
		if net.ParseIP(addr) != nil {
			Log.FmtPrintf("A literal IP address, addr: %s.", addr)
		}
	}
}
