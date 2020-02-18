// add by stefan
package main

import (
	"InnerGateway/client"
	"InnerGateway/server"
	"common/Log"
	"sync"
)

func startInnerGW() {
	var sw sync.WaitGroup
	sw.Add(2)
	go server.StartServer()
	go client.StartClient()
	sw.Wait()
}

func main() {
	Log.FmtPrintf("start InnerGateway.")
	startInnerGW()
}
