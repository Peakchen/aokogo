package main

import (
	"common/Log"
	"simulate/AutoTest"
)

func main() {
	Log.FmtPrintf("main msg test.")
	AutoTest.Start()
}
