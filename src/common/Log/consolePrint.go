package Log

// add by stefan

import (
	"common/public"
	"fmt"
	"time"
)

func FmtPrintf(src string, params ...interface{}) {
	var dst string
	if len(params) == 0 {
		dst = fmt.Sprintf(time.Now().Local().Format(public.CstTimeFmt)+" "+src) + "\n"
	} else {
		dst = fmt.Sprintf(time.Now().Local().Format(public.CstTimeFmt)+" "+src, params...) + "\n"
	}

	fmt.Println(dst)
}

func FmtPrintln(params ...interface{}) {
	content := make([]interface{}, 0, len(params)+1)
	content = append(content, time.Now().Format(public.CstTimeFmt))
	if len(params) > 0 {
		content = append(content, params...)
	}
	fmt.Println(content...)
}

func RetError(context string, params ...interface{}) error {
	return fmt.Errorf(context, params...)
}
