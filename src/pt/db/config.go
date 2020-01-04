package db

import (
	"database/sql"
	//"fmt"
	"common/tool"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"strings"
)

var (
	ListenHost string
)

func LoadDB() {
	tool.MyFmtPrint_Info("begin open db.") 
	
	var (
		findedDB string
	)

	files, _ := ioutil.ReadDir("./")
    for _, f := range files {
		if strings.Contains(f.Name(), ".db") {
			findedDB = f.Name()
			break
		} 
	}
	
	if len(findedDB) == 0 {
		tool.MyFmtPrint_Info("can not find sqlite db file.")
		return
	}

	db, err := sql.Open("sqlite3", "./" + findedDB)
	if err != nil {
		tool.MyFmtPrint_Info("load sqlite db fail, err: ", err)
		return
	}

	if err := db.Ping(); err != nil{
		tool.MyFmtPrint_Info("sqlite db ping fail, err: ", err)
		return
	}

	tool.MyFmtPrint_Info("load success.")
	InitDB(db)
}

func SetListenHost(host string){
	ListenHost =  "http://"+ host
}