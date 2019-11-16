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

func (this *TCache) Set(key string, data interface{}) {
	d := &TData{
		key: 		key,
		deadtime: 	time.Now().Unix()+this.td,
	}
	this.cl.PushBack(d)
	this.c.Store(key, data)
}

func (this *TCache) Get(key string) interface{}{
	val, ok := this.c.Load(key)
	if !ok {
		return nil
	} 
	return val
}

func (this *TCache) Remove(key string) {
	this.c.Delete(key)
}

func (this *TCache) Init(td int64, ctx context.Context) {
	this.td = td
	this.ctx = ctx
	if this.cl == nil {
		this.cl = list.New()
	}
}

func (this *TCache) Run(){
	this.wg.Add(1)
	go this.loopcheck()
}

func (this *TCache) exit(){
	this.wg.Wait()
}

func (this *TCache) loopcheck() {
	defer this.wg.Done()
	t := time.NewTicker(time.Duration(this.td))
	for{
		select{
		case <-this.ctx.Done():
			this.exit()
			return
		case <-t.C:
			if this.cl.Len() == 0 {
				break
			}
			e := this.cl.Front()
			data := e.Value.(*TData)
			if data.deadtime > time.Now().Unix() {
				break
			}
			this.cl.Remove(e)
			this.c.Delete(data.key)
		}
	}
}