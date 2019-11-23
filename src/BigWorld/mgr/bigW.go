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

func (this *TBigWordMgr) Run() {
	this.ctx, this.cancle = context.WithCancel(context.Background())
	this.wg.Add(1)
	this.tcpsvr = NewTcpServer(this.host, this.srcSvr, this.dstSvr, this.Recv)
	this.tcpsvr.Run()
}
//tcp net message recv call back.
func (this *TBigWordMgr) Recv(conn net.Conn, data []byte, len int) {
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
	this.cacheAction(recv.Data, outparams)
	if outparams != nil {
		this.send(recv.Srcid, outparams)
	}
}

func (this *TBigWordMgr) cacheAction(data []byte, outparams interface{}){
	var secDatas = &CacheOperation{}
	secpack := proto.Unmarshal(data, secDatas)
	if secpack != nil {
		log.Fatal("unmarshal message Data fail.")
		return
	}
	SelectOper("", secDatas, outparams)
}

func (this *TBigWordMgr) send(srcSvr int32, outparams interface{}){
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

func (this *TBigWordMgr) loop(){
	defer this.wg.Done()
	for {
		select {
		case <-this.ctx.Done():
			this.Exit()
			return
		default:
			//...
		}
	}
}

func (this *TBigWordMgr) Exit(){
	this.tcpsvr.Exit()
	this.cancle()
	this.wg.Wait()
}