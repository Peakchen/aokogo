package tool

import (
	"fmt"
)

var (
	cstinfoprint bool = true
	cstErrprint bool = true
	CstHideConsole bool = false
)

func MyFmtPrint_Info(args ...interface{}) {
	if cstinfoprint {
		fmt.Println(args...)
	}
}

func MyFmtPrint_Error(args ...interface{}) {
	if cstErrprint {
		fmt.Println(args...)
	}
}