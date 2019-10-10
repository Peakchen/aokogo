package main

import (
	"common/Log"
	"simulate/M_login"
)

func main() {
	Log.FmtPrintf("main msg test.")
	// M_Server.TestServer(nil)
	// M_login.Testlogin(nil)
	M_login.LoginRun()
}

func init() {

}
