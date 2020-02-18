// add by stefan
package main

import (
	"DBServer/server"
	"common/Log"
)

func main() {
	Log.FmtPrintln("run db server.")
	server.StartDBServer()
}

func init() {

}
