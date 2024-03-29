package AsyncLock

//add by stefan
//etcd version for low frequency and long duration.

import (
	"github.com/zieckey/etcdsync"
	"os"
)

//ip -> ip:port

var (
	_etcdmachines = []string{}
	_etcdlocks    = map[string]*etcdsync.Mutex{}
)

//machines are the ectd cluster addresses, such as: http://127.0.0.1:2379
func NewEtcdLock(machines []string) {
	_etcdmachines = machines
}

func AddEtcdLock(key, Name string) (succ bool) {
	lockid := key + ":" + Name
	m, err := etcdsync.New("/"+lockid, 10, _etcdmachines)
	if m == nil || err != nil {
		Log.FmtPrintln("etcdsync New failed.")
		return
	}

	m.SetDebugLogger(os.Stdout)
	err = m.Lock()
	if err != nil {
		Log.Error("etcdsync Lock failed.")
		return
	}

	succ = true
	return
}

func ReleaseEtcdLock(key, Name string) {
	lockKey := key + ":" + Name
	lock, exist := _etcdlocks[lockKey]
	if exist {
		lock.Unlock()
		delete(_etcdlocks, lockKey)
	}
}
