package M_Log

import (
	"common/Log"
	"testing"
)

func TestLogNormal(t *testing.T) {
	Log.FmtPrintln("test log: ", "yes")
	Log.Error("test error.")
}
