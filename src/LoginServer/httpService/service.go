package httpService

import (
	"LoginServer/logindefine"
	"common/httpExService"
	"common/httpsExServiceTLS"
	"common/tcpWebNet"
	"flag"
	"log"
)

var (
	addr     = flag.String("addr", "localhost:8080", "http service address")
	port     = flag.Int("port", 433, "http port.")
	isHttp   = flag.Bool("ishttp", true, "just to get it which is http or not.")
	certfile = flag.String("tlscert", "", "it is a certfile.")
	keyfile  = flag.String("tlskey", "", "about tls key file.")
)

func start() {
	log.Println("start login server.")
	flag.Parse()

	// todo: start http or https thread
	if *isHttp {
		httpExService.StartHttpService(*addr, logindefine.DealWitchLoginHandler)
	} else {
		httpsExServiceTLS.StartHttpsServiceTLS(*port, *certfile, *keyfile, logindefine.DealWitchLoginHandler)
	}

	// start websock server thread.
	tcpWebNet.StartWebSockSvr(*addr)

}
