/*
* CopyRight(C) StefanChen e-mail:2572915286@qq.com
 */

package main

import (
	"LoginServer/server"
	"common/Config/serverConfig"
	"common/Log"
)

func main() {
	Log.FmtPrintln("start login server.")

	server.StartServer(serverConfig.GServerBaseConfig)
	return
}
