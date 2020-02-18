package U_Stack

import (
	"common/stacktrace"
	"testing"
)

func stack_1(normal bool) {
	stack_2(normal)
}

func stack_2(normal bool) {
	stack_3(normal)
}

func stack_3(normal bool) {
	if normal {
		stacktrace.NormalStackLog()
	} else {
		stacktrace.RedStackLog()
	}
}

func TestWhiteStackLog(t *testing.T) {
	stack_1(true)
}

func TestRedStackLog(t *testing.T) {
	stack_1(false)
}
