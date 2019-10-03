package Log

import (
	"fmt"
	"time"
)

const (
	timeFmt string = "2006-01-02 15:04:05.000000 Z0700"
)

func FmtPrintf(src string, params ...interface{}) {
	var dst string
	if len(params) == 0 {
		dst = fmt.Sprintf(time.Now().Local().Format(timeFmt)+" "+src) + "\n"
	} else {
		dst = fmt.Sprintf(time.Now().Local().Format(timeFmt)+" "+src, params...) + "\n"
	}

	fmt.Println(dst)
	WriteLog("[Info]", dst)
}

func FmtPrintln(params ...interface{}) {
	content := make([]interface{}, 0, len(params)+1)
	content = append(content, time.Now().Format(timeFmt)+" ")
	if len(params) > 0 {
		content = append(content, params...)
	}
	fmt.Println(content...)
	WriteLog("[Info]", "", content...)
}
