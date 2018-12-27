package ExGateway

import (
	"fmt"
	"net"
	"sync"
	"common/define"
)


type TMultiWay struct{
	swait  sync.WaitGroup

}

var GMultiWay = &TMultiWay{}

func StartListen()error{
	tcpAddr, err := net.ResolveTCPAddr("ip4", szExternalServerHost)
	if err != nil {
		err := fmt.Errorf("[StartListen] net Resolve TCP Addr , fail to Resolve Addr: %v.", szExternalServerHost)
		return err
	}

	listenr, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil{
		err := fmt.Errorf("[StartListen] net start , fail to listen ip: %v.", szExternalServerHost)
		return err
	}

	GMultiWay.swait.Add(1)
	go AcceptLoop(listenr, &GMultiWay.swait)
	GMultiWay.swait.Wait()
	return nil
}

func AcceptLoop(lis *net.TCPListener, w *sync.WaitGroup){
	defer w.Done()
	for{
		conn, err := lis.AcceptTCP()
		if err != nil{
			fmt.Printf("[AcceptLoop] listen accept err: %v.\n", err)
			continue
		}

		go AcceptProc(conn)
	}
}

func AcceptProc(c net.Conn){
	defer c.Close()
	for {
		arrReadBuff := [128]byte{}
		rn, err := c.Read(arrReadBuff[:])
		if err != nil{
			fmt.Printf("[AcceptProc] read buff data err: %v.\n", err)
			continue
		}

		// data convert to ...
		TansportMessage(arrReadBuff[:rn])
	}
}

func TansportMessage(data []byte){
	
}