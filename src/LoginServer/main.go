/*
* CopyRight(C) StefanChen e-mail:2572915286@qq.com
 */

package main

import (
	"LoginServer/logindefine"
	"common/Define"
	"common/Log"
	"common/tcpNet"
)

func main() {
	Log.FmtPrintln("start login server.")
	var (
		mapsvr map[int32][]int32 = map[int32][]int32{
			int32(Define.ERouteId_ER_ESG): []int32{int32(Define.ERouteId_ER_Login)},
		}
	)

	gameSvr := tcpNet.NewClient(Define.LoginServerHost,
		&mapsvr,
		logindefine.LoginMessageCallBack)

	gameSvr.Run()
	return
}
