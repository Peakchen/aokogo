/*
* CopyRight(C) Stefan e-mail:2572915286@qq.com
 */

package main

import (
	"LoginServer/server"
	"common/Log"
)

func main() {
	Log.FmtPrintln("start login server.")

	server.StartServer()
	return
}
