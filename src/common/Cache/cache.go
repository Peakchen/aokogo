package Cache
// mutli server cache
// add by stefan
import (
	"sync"
	"time"
	"context"
	"container/list"
)

type TData struct {
	key 		string
	deadtime 	int64
}

type TCache struct {
	c 		sync.Map
	wg  	sync.WaitGroup
	//ckey 	chan string
	td 		int64
	cl		*list.List
	ctx 	context.Context
}

func (self *TCache) set(key string, data interface{}) {
	d := &TData{
		key: 		key,
		deadtime: 	time.Now().Unix()+self.td,
	}
	self.cl.PushBack(d)
	self.c.Store(key, data)
}

func (self *TCache) get(key string) interface{}{
	val, ok := self.c.Load(key)
	if !ok {
		return nil
	} 
	return val
}

func (self *TCache) init(td int64, ctx context.Context) {
	self.td = td
	self.ctx = ctx
	if self.cl == nil {
		self.cl = list.New()
	}
}

func (self *TCache) run(){
	self.wg.Add(1)
	go self.loopcheck()
}

func (self *TCache) exit(){
	self.wg.Wait()
}

func (self *TCache) loopcheck() {
	defer self.wg.Done()
	t := time.NewTicker(time.Duration(self.td))
	for{
		select{
		case <-self.ctx.Done():
			self.exit()
			return
		case <-t.C:
			if self.cl.Len() == 0 {
				break
			}
			e := self.cl.Front()
			data := e.Value.(*TData)
			if data.deadtime > time.Now().Unix() {
				break
			}
			self.cl.Remove(e)
			self.c.Delete(data.key)
		}
	}
}