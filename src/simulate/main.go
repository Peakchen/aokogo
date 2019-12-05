package main

import (
	"common/Log"
	"simulate/AutoTest"
)

func main() {
	Log.FmtPrintf("main msg test.")
	// U_Server.TestServer(nil)
	// U_login.Testlogin(nil)
	//U_login.LoginRun()
	AutoTest.Run()
}
