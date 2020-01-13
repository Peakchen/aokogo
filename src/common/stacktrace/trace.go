package stacktrace

/*
	purpose: Stack trace for bug code question finding.
	date: 20200113 14:04
*/

import (
	"runtime/debug"
	"common/Log"
)

/*
	white code log print for normal log.
*/
func NormalStackLog()(buf string){
	buf = string(debug.Stack())
	Log.FmtPrintln("stack trace: ", buf)
	return
}

/*
	red code log print for panic question log. 
*/
func RedStackLog(){
	debug.PrintStack()
}

