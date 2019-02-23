package test

import (
	"fmt"
	"testing"
	
)

func Test_tcp_1(t *testing.T){
	fmt.Println("[Test_tcp_1] start.")
	cmd := 65546
	c := uint16(cmd)
	a := uint16(cmd >> 16)
	b := uint16(cmd)

	fmt.Printf("a: %v, b: %v, c: %v.\n", a,b,c)
}

func init(){

}