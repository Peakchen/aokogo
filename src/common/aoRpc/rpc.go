package aoRpc

// add by stefan 20190614 19:49
// add aorpc for between server and server conmunication. 
import (
	"sync"
	"context"
	"time"
	"fmt"
	"reflect"
	"container/list"
)

const (
	rpcdealline int32 = 60*5 // five min ove time
	ActChanMaxSize  int = 1000 // act call params
)

type TModelAct struct {
	actid	string
	modf  	reflect.Value
	params 	[]reflect.Value
}

type TActRet struct {
	actid	string
	rets 	[]reflect.Value
}

type TAorpc struct {
	models      map[string]interface{}
	wg     		sync.WaitGroup
	ctx 		context.Context
	cancel		context.CancelFunc
	acts 		*list.List
	actchan		chan *TModelAct
	mutex		sync.Mutex
	retchan     chan *TActRet
}

var Aorpc *TAorpc = nil
func init(){
	Aorpc = &TAorpc{}
	Aorpc.Init()
}

func (self *TAorpc) Init(){
	self.models = map[string]interface{}{}
	self.acts = list.New()
}

func (self *TAorpc) Run(){
	self.ctx, self.cancel = context.WithCancel(context.Background())
	self.wg.Add(2)
	go self.loop()
	go self.loopAct()
}

/*
	take model and func witch func in params,   
*/
func (self *TAorpc) Call(key, modelname, funcName string, ins []interface{}, outs []interface{})(error){
	m, ok := self.models[modelname]
	if !ok {
		return fmt.Errorf("can not find model, input model name: %v.", modelname)
	}
	v := reflect.ValueOf(m)
	f := v.MethodByName(funcName)
	rv := []reflect.Value{}
	for _, in := range ins {
		rv = append(rv, reflect.ValueOf(in))
	}
	//f.Call(rv)
	actkey := key+":"+modelname+":"+funcName
	self.actchan <- &TModelAct{
		actid:	actkey,
		modf: 	f,
		params: rv,
	}
	var twg sync.WaitGroup
	twg.Add(1)
	go self.loopRet(actkey, outs, &twg)
	twg.Wait()
	return nil
}

func (self *TAorpc) loopRet(actkey string, outs []interface{}, twg *sync.WaitGroup){
	t := time.NewTicker(time.Duration(rpcdealline))
	for {
		select {
		case ar := <-self.retchan:
			if ar.actid == actkey {
				for i, ret := range ar.rets {
					reflect.ValueOf(outs[i]).Set(ret)
				}
				twg.Done()
				return
			}
		case <-t.C:
			// beyond return time, then return nothing.
			twg.Done()
		}
	}
}

func (self *TAorpc) loop(){
	defer self.wg.Done()
	//t := time.NewTicker(time.Duration(rpcdealline))
	for {
		select {
		case <-self.ctx.Done():
			self.Exit()
			return
		//case <-t.C:
		
		case act := <-self.actchan:
			if act == nil {
				return
			}
			if self.acts.Len() >= ActChanMaxSize {
				fmt.Println("has enough acts in chan.")
				return
			} 
			self.acts.PushBack(act)
		}
	}
}

func (self *TAorpc) loopAct(){
	defer self.wg.Done()
	for {
		if self.acts.Len() == 0 {
			continue
		}
		self.mutex.Lock()
		e := self.acts.Front()
		act := e.Value.(*TModelAct)
		if act == nil {
			fmt.Println("act value invalid: ", e.Value)
			continue
		}
		mrts := act.modf.Call(act.params)

		self.retchan <- &TActRet{
			actid:	act.actid,
			rets:	mrts,
		}
		self.acts.Remove(e)
		self.mutex.Unlock()
	}
}

func (self *TAorpc) Exit(){
	self.cancel()
	self.wg.Wait()
}

func Register(name string, model interface{}) {
	_, ok := Aorpc.models[name]
	if ok {
		return
	}
	Aorpc.models[name] = model
}

