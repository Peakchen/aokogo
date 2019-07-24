package timer

import (
	"common/RedisService"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// add by stefan 20190715 19:39
// purpose: each model timer call back refresh data.

type TAokoCallBackParam struct {
	cb       interface{} //call back func
	params   interface{} //func params
	interval int32       //timer interval
	times    int32       //func call times
}

type TDataPack struct {
	Key  string
	Data string
}

type TAokoTimer struct {
	tmo 	[]string //
	conn	*RedisService.RedisConn.Conn
}

var (
	GAokoTimer = &TAokoTimer{}
)

func NewTimer(ctx context.Context, wg *sync.WaitGroup, c *RedisService.RedisConn) {
	GAokoTimer.conn = c.Conn
	GAokoTimer.tmo = []string{}
	wg.Add(1)
	go self.loop(ctx, wg)

}

func (self *TAokoTimer) Register(key, name string, model interface{}) {
	cbdata := &TAokoCallBackParam{
		cb: model,
	}
	bydata, err := json.Marshal(cbdata)
	if err != nil {
		fmt.Println("register marshal fail: ", err)
		return
	}
	self.tmo = append(self.tmo, name)
	pack = &TDataPack{
		Key:  key,
		Data: string(bydata),
	}
	bypack, err := json.Marshal(pack)
	if err != nil {
		fmt.Println("register marshal fail: ", err)
		return
	}
	// key + value
	_, err := self.conn.Do("RPUSH", name, bypack...)
	if err != nil {
		Log.Error("RPUSH data: %v, err: %v.\n", Ret, err)
		return
	}
}

func (self *TAokoTimer) loop(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		for _, name := range self.tmo{
			data, err := self.conn.Do("LPOP", name)
			if err != nil {
				Log.Error("[Save] SETNX data: %v, err: %v.\n", data, err)
				return
			}
			if len(string(data)) == 0{
				return
			}
			info := &TDataPack{}
			if err := json.Unmarshal(data, info); err != nil{
				log.Error("")
				return
			}

		}
	}
}

func (self *TAokoTimer) handler() {

}

func (self *TAokoTimer) exit(wg *sync.WaitGroup) {
	wg.Wait()
}
