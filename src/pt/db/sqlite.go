package db

import (
	"database/sql"
	"fmt"
	//"reflect"
	"common/tool"
)

type TSqliteDB struct {
	SqliteDB *sql.DB
}

var (
	_sqlitedb = &TSqliteDB{}
)

func InitDB(db *sql.DB) {
	_sqlitedb.init(db)
}

func (this *TSqliteDB) init(db *sql.DB) {
	this.SqliteDB = db
}

func (this *TSqliteDB) BatchUpdate(selectfields string, fieldsdata []interface{}) {
	stmt, err := this.SqliteDB.Prepare(selectfields)
	if err != nil {
		tool.MyFmtPrint_Error("db prepare insert, err: ", err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(fieldsdata...)
	if err != nil {
		tool.MyFmtPrint_Error("db Exec insert, err: ", err)
		return
	}
}

func (this *TSqliteDB) BatchUpdateDirect(sqlcmds string) (err error){
	err = this.UpdateDirect(sqlcmds)
	return
}

/*
	@purpose: insert,update,delete operations.
	@param1: selectfields string
	@param2: fieldsdata []interface{}
*/
func (this *TSqliteDB) Update(selectfields string, fieldsdata []interface{}) {
	stmt, err := this.SqliteDB.Prepare(selectfields)
	if err != nil {
		tool.MyFmtPrint_Error("db prepare insert, err: ", err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(fieldsdata...)
	if err != nil {
		tool.MyFmtPrint_Error("db Exec insert, err: ", err)
		return
	}
}

func (this *TSqliteDB) UpdateDirect(sqlcmds string) (err error){
	_, err = this.SqliteDB.Exec(sqlcmds)
	if err != nil {
		err = fmt.Errorf("db Exec insert, err: %v.", err)
		return
	}
	return
}

func GetSqliteDB()*TSqliteDB{
	return _sqlitedb
}

/*
	@purpose: query operation.
	@param1: sqlcmd string
*/
func (this *TSqliteDB) Query(sqlcmd string) (ret []([]*TQueryField), err error){
	ret = make([]([]*TQueryField), 0, 1000000)
	rows, err := this.SqliteDB.Query(sqlcmd)
	if err != nil {
		err = fmt.Errorf("db prepare insert, err: %v.", err)
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		err = fmt.Errorf("Columns: %v", err)
		return
	}

	//tool.MyFmtPrint_Info("Query Columns: ", cols)

	fieldsdata := make([]interface{}, len(cols))
	for i := range fieldsdata {
		var temp interface{}
		fieldsdata[i] = &temp
	}

	for rows.Next() {
		err = rows.Scan(fieldsdata...)
		if err != nil {
			err = fmt.Errorf("db scan field data fail, err: %v.", err)
			return
		}

		//tool.MyFmtPrint_Info("query fieldsdata: ", fieldsdata)
		columdata := []*TQueryField{}
		for idx, data := range fieldsdata{
			//tool.MyFmtPrint_Info("Field: ", cols[idx], ", data: ", *data.(*interface{}))
			columdata = append(columdata, &TQueryField{
				Field: cols[idx],
				Data: *data.(*interface{}),
			})
		}

		ret = append(ret, columdata)
	}

	if err = rows.Close(); err != nil {
		err = fmt.Errorf("error closing rows: %s", err)
		return
	}

	err = nil
	//tool.MyFmtPrint_Info("query success, ret: ", len(ret), ret)
	return
}
