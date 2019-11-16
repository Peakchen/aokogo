package Cache

/*
	cache data not exists forever,
	which exist over deadline,
	then delete it and to get new data form bigword.
*/

import (
	. "common/Define"
	. "common/S2SMessage"
	. "common/tcpNet"
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

type TCacheMgr struct {
	wg        sync.WaitGroup
	c         *TCache //some one cache
	ctx       context.Context
	cancel    context.CancelFunc
	s         *TcpSession
	mapSvr    map[int32][]int32
	sessAlive bool
	cb        MessageCb
	// add cache obj ...
}

func (this *TCacheMgr) connect() error {
	c, err := net.Dial("tcp", ConstBigWordHost)
	if err != nil {
		fmt.Println("net dial err: ", err)
		return err
	}
	c.(*net.TCPConn).SetNoDelay(true)
	this.s = NewSession(ConstBigWordHost, c, this.ctx, this.mapSvr, this.cb)
	this.s.HandleSession()
	this.sessAlive = true
	return nil
}

func (this *TCacheMgr) run(srcSev, dstSvr int32, cb MessageCb) {
	this.ctx, this.cancel = context.WithCancel(context.Background())
	this.srcSvr = srcSev
	this.dstSvr = dstSvr
	this.cb = cb

	this.c = &TCache{}
	this.c.Init(ConstCacheOverTime, this.ctx)
	this.c.Run()
	this.connect()

	this.wg.Add(1)
	go this.loopc()
	// add ...
}

func (this *TCacheMgr) exit() {
	this.cancel()
	this.wg.Wait()
}

func (this *TCacheMgr) loopc() {
	defer this.wg.Done()
	t := time.NewTicker(time.Duration(ConstCacheUpdateTime))
	for {
		select {
		case <-this.ctx.Done():
			this.exit()
			return
		case <-t.C:
			//begin request data from bigword server.
			msg := &SS_BaseMessage_Req{
				Srcid: this.srcSvr,
				Dstid: int32(ERouteId_ER_BigWorld),
			}
			data, err := proto.Marshal(msg)
			if err == nil {
				fmt.Println("Marshal message fail.")
				return
			}
			this.s.SetSendCache(data)
		default:
			if !this.sessAlive {
				this.connect()
			}
		}
	}
}
