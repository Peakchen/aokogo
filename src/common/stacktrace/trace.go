package stacktrace

/*
	purpose: Stack trace for bug code question finding.
	date: 20200113 14:04
*/

import (
	"common/Log"
	"runtime/debug"
)

/*
	white code log print for normal log.
*/
func NormalStackLog() (stacklog string) {
	stacklog = string(debug.Stack())
	Log.FmtPrintln("stack trace: ", stacklog)
	return
}

/*
	red code log print for panic question log.
*/
func RedStackLog() {
	debug.PrintStack()
}
