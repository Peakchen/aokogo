package MgoConn

import (
	"common/Log"
	"common/RedisConn"
	. "common/public"
	"fmt"
	"reflect"
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
	server      string
}

func NewMgoConn(server, Username, Passwd, Host string) *AokoMgo {
	aokomogo := &AokoMgo{}
	aokomogo.UserName = Username
	aokomogo.Passwd = Passwd
	aokomogo.ServiceHost = Host
	aokomogo.server = server
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

func (this *AokoMgo) QueryAcc(usrName string, OutParam IDBCache) (err error, exist bool) {
	condition := bson.M{OutParam.SubModel() + "." + "username": usrName}
	return this.QueryByCondition(condition, OutParam)
}

func (this *AokoMgo) QueryOne(Identify string, OutParam IDBCache) (err error, exist bool) {
	condition := bson.M{"_id": Identify}
	return this.QueryByCondition(condition, OutParam)
}

func (this *AokoMgo) QueryByCondition(condition bson.M, OutParam IDBCache) (err error, exist bool) {
	session, err := this.GetSession()
	if err != nil {
		err = Log.RetError("get sesson err: %v.", err)
		return
	}

	s := session.Clone()
	defer s.Close()
	collection := s.DB(this.server).C(OutParam.MainModel())
	retQuerys := collection.Find(condition)
	count, ret := retQuerys.Count()
	if ret != nil || count == 0 {
		err = Log.RetError("[mgo] query data err: %v, %v.", ret, count)
		return
	}

	selectRet := retQuerys.Select(bson.M{OutParam.SubModel(): 1, "_id": 1}).Limit(1)
	if selectRet == nil {
		err = Log.RetError("[mgo] selectRet invalid, submodule: %v.", OutParam.SubModel())
		return
	}

	outval := reflect.MakeMap(reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(OutParam)))
	ret = selectRet.One(outval.Interface())
	if ret != nil {
		err = Log.RetError("[mgo] select one error: %v.", ret)
		return
	}

	retIdxVal := outval.MapIndex(reflect.ValueOf(OutParam.SubModel()))
	if !retIdxVal.IsValid() {
		err = Log.RetError("[mgo] outval MapIndex invalid.")
		return
	}

	reflect.ValueOf(OutParam).Elem().Set(retIdxVal.Elem())
	exist = true
	return
}

func (this *AokoMgo) QuerySome(Identify string, OutParam IDBCache) (err error) {
	session, err := this.GetSession()
	if err != nil {
		return err
	}

	s := session.Clone()
	defer s.Close()
	collection := s.DB(this.server).C(OutParam.MainModel())
	err = collection.Find(bson.M{"_id": Identify}).All(&OutParam)
	if err != nil {
		err = fmt.Errorf("Identify: %v, MainModel: %v, SubModel: %v, err: %v.\n", Identify, OutParam.MainModel(), OutParam.SubModel(), err)
		Log.Error("[QuerySome] err: %v.\n", err)
	}
	return
}

func (this *AokoMgo) InsertOne(Identify string, InParam IDBCache) (err error) {
	session, err := this.GetSession()
	if err != nil {
		return err
	}

	s := session.Clone()
	defer s.Close()
	collection := s.DB(this.server).C(InParam.MainModel())
	operAction := bson.M{"_id": Identify, InParam.SubModel(): InParam}
	err = collection.Insert(operAction)
	if err != nil {
		err = fmt.Errorf("Identify: %v, MainModel: %v, SubModel: %v, err: %v.\n", Identify, InParam.MainModel(), InParam.SubModel(), err)
		Log.Error("[SaveOne] err: %v.\n", err)
	}
	return
}

func (this *AokoMgo) SaveOne(Identify string, InParam IDBCache) (err error) {
	session, err := this.GetSession()
	if err != nil {
		return err
	}
	s := session.Clone()
	defer s.Close()
	collection := s.DB(this.server).C(InParam.MainModel())
	operAction := bson.M{InParam.SubModel(): InParam}
	err = collection.Update(bson.M{"_id": Identify}, operAction)
	if err != nil {
		err = fmt.Errorf("Identify: %v, MainModel: %v, SubModel: %v, err: %v.\n", Identify, InParam.MainModel(), InParam.SubModel(), err)
		Log.Error("[SaveOne] err: %v.\n", err)
	}
	return
}

func (this *AokoMgo) Save(redkey string, data interface{}) (err error) {
	session, err := this.GetSession()
	if err != nil {
		return err
	}
	s := session.Clone()
	defer s.Close()

	main, sub, key := RedisConn.ParseRedisKey(redkey)
	collection := s.DB(this.server).C(main)
	operAction := bson.M{sub: data}
	_, err = collection.UpsertId(key, operAction)
	if err != nil {
		err = fmt.Errorf("main: %v, sub: %v, key: %v, err: %v.\n", main, sub, key, err)
		Log.Error("[Save] err: %v.\n", err)
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
	c := s.DB(this.server).C(InParam.SubModel())
	err = c.EnsureIndex(mgo.Index{Key: idxs})
	return err
}
