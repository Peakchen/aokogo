package Log

// add by stefan

import (
	"common/aktime"
	"common/public"
	"fmt"
)

func FmtPrintf(src string, params ...interface{}) {
	var dst string
	if len(params) == 0 {
		dst = fmt.Sprintf(aktime.Now().Local().Format(public.CstTimeFmt)+" "+src) + "\n"
	} else {
		dst = fmt.Sprintf(aktime.Now().Local().Format(public.CstTimeFmt)+" "+src, params...) + "\n"
	}

	fmt.Println(dst)
}

func FmtPrintln(params ...interface{}) {
	content := make([]interface{}, 0, len(params)+1)
	content = append(content, aktime.Now().Format(public.CstTimeFmt))
	if len(params) > 0 {
		content = append(content, params...)
	}
	fmt.Println(content...)
}

func RetError(context string, params ...interface{}) error {
	return fmt.Errorf(context, params...)
}
