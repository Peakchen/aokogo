/*
* CopyRight(C) StefanChen e-mail:2572915286@qq.com
 */

package main

import (
	"flag"
	"fmt"

	//"log"
	"common/tcpWebNet"
)

func init() {

}

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	fmt.Println("start ExternalWebGateway.")
	flag.Parse()

	// start websock server.
	tcpWebNet.StartWebSockSvr(*addr)
}
