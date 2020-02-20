package AsyncLock

//add by stefan
// zookeeper version for low frequency and long duration.

import (
	"github.com/samuel/go-zookeeper/zk"
	"net"
	"time"
)

var (
	zkconn   *net.Conn
	_zklocks = map[string]*zk.Lock{}
)

//ip -> ip:port
func NewZKLock(ips []string) {
	var err error
	zkconn, _, err = zk.Connect(ips, time.Second) //default 1s
	if err != nil {
		panic(err)
	}
}

func AddZKLock(key, Name string) (succ bool) {
	lockKey := key + ":" + Name
	zl := zk.NewLock(zkconn, "/"+lockKey, zk.WorldACL(zk.PermAll))
	if err := zl.Lock(); err != nil {
		panic(err)
	}
	succ = true
	_zklocks[lockKey] = zl
	return
}

func ReleaseZKLock(key, Name string) {
	lockKey := key + ":" + Name
	lock, exist := _zklocks[lockKey]
	if exist {
		lock.Unlock()
		delete(_zklocks, lockKey)
	}
}
