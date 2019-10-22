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

func (this *AokoMgo) NewDial() {

	for i := 1; i <= this.PoolCnt; i++ {
		session, err := this.getSession()
		if err != nil {
			Log.FmtPrintln(err)
			continue
		}
		this.chSessions <- session
	}
	return
}

func (this *AokoMgo) getSession() (session *mgo.Session, err error) {
	MdialInfo := &mgo.DialInfo{
		Addrs:        []string{this.ServiceHost},
		Username:     this.UserName,
		Password:     this.Passwd,
		Direct:       false,
		Timeout:      time.Second * 3,
		PoolLimit:    4096,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	session, err = mgo.DialWithInfo(MdialInfo)
	if err != nil {
		err = Log.RetError("mgo dial err: %v.\n", err)
		Log.Error("mgo dial err: %v.\n", err)
		return
	}

	err = session.Ping()
	if err != nil {
		err = Log.RetError("session ping out, err: %v.", err)
		Log.Error("session ping out, err: %v.", err)
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
	err = nil
	return
}

func (this *AokoMgo) Stop() {
	if this.session != nil {
		this.session.Close()
	}
}

func (this *AokoMgo) GetDB() *mgo.Session {
	if this.session == nil {
		this.NewDial()
	}

	return this.session.Clone()
}

func (this *AokoMgo) OnTimer2FlushDB() {
	reach := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-reach.C:
			// todo:
			this.FlushDB()
		default:
			// nothing...

		}
	}
}

func (this *AokoMgo) FlushDB() {

}

func (this *AokoMgo) GetSession() (sess *mgo.Session, err error) {
	select {
	case s, _ := <-this.chSessions:
		return s, nil
	case <-time.After(time.Duration(time.Second)):
	default:
	}
	return nil, fmt.Errorf("aoko mongo session time out and not get.")
}

func MakeMgoModel(Identify, MainModel, SubModel string) string {
	return MainModel + "." + SubModel + "." + Identify
}

func (this *AokoMgo) QueryOne(OutParam IDBCache) (err error) {
	session, err := this.GetSession()
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

func (this *AokoMgo) QuerySome(OutParam IDBCache) (err error) {
	session, err := this.GetSession()
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

func (this *AokoMgo) SaveOne(InParam IDBCache) (err error) {
	session, err := this.GetSession()
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

func (this *AokoMgo) EnsureIndex(InParam IDBCache, idxs []string) (err error) {
	session, err := this.GetSession()
	if err != nil {
		return err
	}
	s := session.Clone()
	defer s.Close()
	c := s.DB(InParam.MainModel()).C(InParam.SubModel())
	err = c.EnsureIndex(mgo.Index{Key: idxs})
	return err
}
