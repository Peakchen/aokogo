package tool

// add by stefan

import (
	"bytes"
	"github.com/gonutz/ide/w32"
	"os"
	"os/exec"
	"strconv"
	"strings"
	//"runtime"
	"common/Log"
	//"syscall"
)

func CheckPortUsed(port int) bool {
	args := []string{"cmd", "/c", "netstat", "-an", "|", "findstr", strconv.Itoa(port)}
	cmd := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		Log.Error(err.Error() + ": " + stderr.String())
		return false
	}

	return true
}

func CmdHide() {
	handler, err := os.OpenFile("hide.vbs", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	vbscontext := `Set ws = CreateObject("Wscript.Shell") 
				ws.run "cmd /c zxcad.exe",vbhide`

	handler.WriteString(vbscontext)
	handler.Close()

	args := []string{"cmd", "/c", "hide.vbs"}
	cmd := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		Log.Error(err.Error() + ": " + stderr.String())
		return
	}

	err = os.Remove("hide.vbs")
	if err != nil {
		Log.Error(err.Error())
		return
	}
}

func BatHide(param []string) {
	handler, err := os.OpenFile("hide.bat", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	batpath := "\"%~nx0 h\""
	batContext := "@echo off \r\n" +
		"if %1 == h goto begin \r\n" +
		//"start " + param[0] + " \r\n" +
		"mshta vbscriptcreateobject(wscript.shell).run(" + batpath + ",0) \r\n" +
		//"mshta vbscriptcreateobject(wscript.shell).run(%~nx0 h,0)(window.close)&&exit \r\n" +
		param[0] + " \r\n" +
		":begin \r\n" +
		"REM"

	handler.WriteString(batContext)
	handler.Close()

	args := []string{"cmd", "/c", "hide.bat", param[1]}
	cmd := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		Log.Error(err.Error() + ": " + stderr.String())
		return
	}

	err = os.Remove("hide.bat")
	if err != nil {
		Log.Error(err.Error())
		return
	}
}
