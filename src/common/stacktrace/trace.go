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
	return string(debug.Stack())
}

/*
	red code log print for panic question log.
*/
func RedStackLog() {
	debug.PrintStack()
}

func Catchcrash() {
	if r := recover(); r != nil {
		stacklog := NormalStackLog()
		Log.Error("catch recover: ", r, stacklog)
	}
}
