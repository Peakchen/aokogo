package MgoService

import (
	"common/Log"
	. "common/public"
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type AokoMgo struct {
	session     *mgo.Session
	UserName    string
	Passwd      string
	ServiceHost string
	chSessions  chan *mgo.Session
	PoolCnt     int
}

func NewMgoConn(Username, Passwd, Host string) *AokoMgo {
	aokomogo := &AokoMgo{}
	aokomogo.UserName = Username
	aokomogo.Passwd = Passwd
	aokomogo.ServiceHost = Host
	//template set 10 session.
	aokomogo.PoolCnt = 10
	aokomogo.chSessions = make(chan *mgo.Session, aokomogo.PoolCnt)
	aokomogo.NewDial()
	return aokomogo
}

func (self *AokoMgo) NewDial() {
	MdialInfo := &mgo.DialInfo{
		Addrs:        []string{self.ServiceHost},
		Username:     self.UserName,
		Password:     self.Passwd,
		Direct:       false,
		Timeout:      time.Second * 3,
		PoolLimit:    4096,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	session, err := mgo.DialWithInfo(MdialInfo)
	if err != nil {
		Log.Error("mgo dial err: %v.\n", err)
		return
	}

	err = session.Ping()
	if err != nil {
		Log.Error("session ping out, err: ", err)
		return
	}

	session.SetMode(mgo.Monotonic, true)
	session.SetCursorTimeout(0)
	// focus on those selects.
	//http://www.mongoing.com/archives/1723
	Safe := &mgo.Safe{
		J:     true,       //true:写入落到磁盘才会返回|false:不等待落到磁盘|此项保证落到磁盘
		W:     1,          //0:不会getLastError|1:主节点成功写入到内存|此项保证正确写入
		WMode: "majority", //"majority":多节点写入|此项保证一致性|如果我们是单节点不需要这只此项
	}
	session.SetSafe(Safe)
	//session.SetSocketTimeout(time.Duration(5 * time.Second()))
	for i := 1; i <= self.PoolCnt; i++ {
		self.chSessions <- self.session.Copy()
	}
	return
}

func (self *AokoMgo) Stop() {
	if self.session != nil {
		self.session.Close()
	}
}

func (self *AokoMgo) GetDB() *mgo.Session {
	if self.session == nil {
		self.NewDial()
	}

	return self.session.Clone()
}

func (self *AokoMgo) OnTimer2FlushDB() {
	reach := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-reach.C:
			// todo:
			self.FlushDB()
		default:
			// nothing...

		}
	}
}

func (self *AokoMgo) FlushDB() {

}

func (self *AokoMgo) GetSession() (sess *mgo.Session, err error) {
	select {
	case s, _ := <-self.chSessions:
		return s, nil
	case <-time.After(time.Duration(time.Second)):
	default:
	}
	return nil, fmt.Errorf("aoko mongo session time out and not get.")
}

func MakeMgoModel(Identify, MainModel, SubModel string) string {
	return MainModel + "." + SubModel + "." + Identify
}

func (self *AokoMgo) QueryOne(OutParam IDBCache) (err error) {
	session, err := self.GetSession()
	if err != nil {
		return err
	}
	s := session.Clone()
	defer s.Close()
	collection := s.DB(OutParam.MainModel()).C(OutParam.SubModel())
	err = collection.Find(bson.M{"_id": OutParam.CacheKey()}).One(&OutParam)
	if err != nil {
		err = fmt.Errorf("CacheKey: %v, MainModel: %v, SubModel: %v, err: %v.\n", OutParam.CacheKey(), OutParam.MainModel(), OutParam.SubModel(), err)
		Log.Error("[QueryOne] err: %v.\n", err)
	}
	return
}

func (self *AokoMgo) QuerySome(OutParam IDBCache) (err error) {
	session, err := self.GetSession()
	if err != nil {
		return err
	}
	s := session.Clone()
	defer s.Close()
	collection := s.DB(OutParam.MainModel()).C(OutParam.SubModel())
	err = collection.Find(bson.M{"_id": OutParam.CacheKey()}).All(&OutParam)
	if err != nil {
		err = fmt.Errorf("CacheKey: %v, MainModel: %v, SubModel: %v, err: %v.\n", OutParam.CacheKey(), OutParam.MainModel(), OutParam.SubModel(), err)
		Log.Error("[QuerySome] err: %v.\n", err)
	}
	return
}

func (self *AokoMgo) SaveOne(InParam IDBCache) (err error) {
	session, err := self.GetSession()
	if err != nil {
		return err
	}
	s := session.Clone()
	defer s.Close()
	collection := s.DB(InParam.MainModel()).C(InParam.SubModel())
	operAction := bson.M{"$set": bson.M{InParam.SubModel(): InParam}}
	err = collection.Update(bson.M{"_id": InParam.CacheKey()}, operAction)
	if err != nil {
		err = fmt.Errorf("CacheKey: %v, MainModel: %v, SubModel: %v, err: %v.\n", InParam.CacheKey(), InParam.MainModel(), InParam.SubModel(), err)
		Log.Error("[SaveOne] err: %v.\n", err)
	}
	return
}

func (self *AokoMgo) EnsureIndex(InParam IDBCache, idxs []string) (err error) {
	session, err := self.GetSession()
	if err != nil {
		return err
	}
	s := session.Clone()
	defer s.Close()
	c := s.DB(InParam.MainModel()).C(InParam.SubModel())
	err = c.EnsureIndex(mgo.Index{Key: idxs})
	return err
}
