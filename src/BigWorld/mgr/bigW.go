package mgr

//add by stefan 20190610 14:27

import (
	. "common/tcpNet"
	"context"
	"sync"
	//"time"
	. "BigWorld/bigcache"
	"net"
	"github.com/golang/protobuf/proto"
	. "common/S2SMessage"
	"log"
	. "common/Define"
	"fmt"
)

type TBigWordMgr struct {
	host 	string
	ctx 	context.Context
	cancle 	context.CancelFunc
	wg  	sync.WaitGroup
	srcSvr  int32
	dstSvr  int32
	tcpsvr  *TcpServer
}

func NewBigWord(host string) *TBigWordMgr {
	return &TBigWordMgr{
		host: host,
	}
}

func (self *TBigWordMgr) Run() {
	self.ctx, self.cancle = context.WithCancel(context.Background())
	self.wg.Add(1)
	self.tcpsvr = NewTcpServer(self.host, self.srcSvr, self.dstSvr, self.Recv)
	self.tcpsvr.StartTcpServer()
}
//tcp net message recv call back.
func (self *TBigWordMgr) Recv(conn net.Conn, data []byte, len int) {
	var recv = &SS_BaseMessage_Req{}
	pm := proto.Unmarshal(data, recv)
	if pm != nil {
		log.Fatal("unmarshal message fail.")
		return
	}
	if recv.Dstid != int32(ERouteId_ER_BigWorld){
		fmt.Println("recv dst: ", recv.Dstid)
		return
	}
	//do something action...

	//respone messge.
	var rsp = &SS_BaseMessage_Rsp{}
	rsp.Ret = EMessageErr_Success
	brsp, err := proto.Marshal(rsp)
	if err != nil {
		fmt.Println("can not marshal base message.")
		return
	}
	wlen, err := conn.Write(brsp)
	if err != nil {
		fmt.Println("respone fail, err: ", err, wlen)
		return
	}
	// do another action...
	var outparams interface{} = nil
	self.cacheAction(recv.Data, outparams)
	if outparams != nil {
		self.send(recv.Srcid, outparams)
	}
}

func (self *TBigWordMgr) cacheAction(data []byte, outparams interface{}){
	var secDatas = &CacheOperation{}
	secpack := proto.Unmarshal(data, secDatas)
	if secpack != nil {
		log.Fatal("unmarshal message Data fail.")
		return
	}
	SelectOper("", secDatas, outparams)
}

func (self *TBigWordMgr) send(srcSvr int32, outparams interface{}){
	if _, ok := ERouteId_name[srcSvr]; !ok {
		return
	}
	switch ERouteId(srcSvr) {
	case ERouteId_ER_Game:

	case ERouteId_ER_DB:

	default:
		log.Fatal("invalid srcSvr: ", srcSvr)
	}
}

func (self *TBigWordMgr) loop(){
	defer self.wg.Done()
	for {
		select {
		case <-self.ctx.Done():
			self.Exit()
			return
		default:
			//...
		}
	}
}

func (self *TBigWordMgr) Exit(){
	self.tcpsvr.Exit()
	self.cancle()
	self.wg.Wait()
}