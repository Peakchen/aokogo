// add by stefan

package main

import (
	"GameServer/server"
	"common/Log"
	//"log"
)

func init() {

}

func main() {
	Log.FmtPrintf("start gameServer.")

	server.StartServer()
	return
}
