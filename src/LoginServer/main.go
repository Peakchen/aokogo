/*
* CopyRight(C) StefanChen e-mail:2572915286@qq.com
*/

package main

import (
	"log"
	"flag"
	"common/websockNet"
	"common/httpExService"
	"common/httpsExServiceTLS"
	"runtime"
	"pro/LoginServer/loginHandler"
)

var (
	addr = flag.String("addr", "localhost:8080", "http service address")
	port = flag.Int("port", 433, "http port.")
	isHttp = flag.Bool("ishttp", true, "just to get it which is http or not.")
	certfile = flag.String("tlscert", "", "it is a certfile.")
	keyfile = flag.String("tlskey","","about tls key file.")
)

func init(){
	runtime.GOMAXPROCS(1)
	log.Println("init login server.")
}

func main()  {
	log.Println("start login server.")
	flag.Parse()
	
	// todo: start http or https thread
	if *isHttp {
		httpExService.StartHttpService(*addr, loginHandler.DealWitchLoginHandler)
	}else{
		httpsExServiceTLS.StartHttpsServiceTLS(*port, *certfile, *keyfile, loginHandler.DealWitchLoginHandler)
	}

	// start websock server thread. 
	go websockNet.StartWebSockService(*addr)

}