/*
* CopyRight(C) StefanChen e-mail:2572915286@qq.com
*/

package main

import(
	"flag"
	"Pro/common/tcpsockNet"
)

var addr = flag.String("addr", "localhost:17000", "http service address")

func main(){
	tcpsockNet.startTcpServer(*addr)
}
