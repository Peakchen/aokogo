package main

// add by stefan

import (
	"common/Log"
	"simulate/AutoTest"
)

func main() {
	Log.FmtPrintf("main msg test.")
	AutoTest.Start()
}
