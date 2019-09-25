package Log

import (
	"fmt"
	"time"
)

const (
	timeFmt string = "2006-01-02 15:04:05.000000 Z0700"
)

func FmtPrintf(src string, params ...interface{}) {
	fmt.Printf(time.Now().Local().Format(timeFmt)+" "+src, params...)
	fmt.Println("\n")
}

func FmtPrintln(src string, params ...interface{}) {
	content := make([]interface{}, 0, len(params)+1)
	content = append(content, time.Now().Format(timeFmt)+" ")
	content = append(content, src)
	content = append(content, params...)
	fmt.Println(content...)
}
