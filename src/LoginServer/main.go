/*
* CopyRight(C) StefanChen e-mail:2572915286@qq.com
 */

package main

import (
	"ExternalGateway/SessionMgr"
	"LoginServer/LogicMsg"
	"common/Define"
	"common/Log"
	"common/tcpNet"
)

func main() {
	Log.FmtPrintln("start login server.")

	gameSvr := tcpNet.NewClient(Define.LoginServerHost,
		Define.ERouteId_ER_Login,
		Define.ERouteId_ER_Login,
		Define.ERouteId_ER_ESG,
		LogicMsg.LoginMessageCallBack,
		nil,
		SessionMgr.GClient2ServerSession)

	gameSvr.Run()
	return
}
