package db

import (
	"net/http"
	"fmt"
	"common/message"
	"io/ioutil"
	//"encoding/json"
	"strconv"
	"reflect"
	"time"
	"common/tool"
)

// oper ret value
var (
	cstBatchUpdate = string("BatchUpdate")
	cstSingleUpdate = string("Update")
	cstSingleQuery = string("Query")
)

type TResponseCode struct {
	Ret interface{} `json:"Ret"`
	Code string `json:"Code"`
}

var (
	cstSuccess = string("success")
)

func parseHttpRequest(r *http.Request) (sqlcmds string, err error){
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("read http request fail, err: %v.", err)
		return
	}

	tool.MyFmtPrint_Info("recv data: ", string(data)) //data[1:len(data)-1]
	sqlcmds = string(data)
	return
}

func ErrResponse(w http.ResponseWriter, HeadValue, err string){
	err = "(\"" + "ret" + "\"" + CstPointformat + "\"" + err + "\")"
	message.HttpResponse(w, []*message.THttpResponseHead{
		&message.THttpResponseHead{
			HeadKey: "key",
			HeadValue: HeadValue,
		},
	}, message.StatusMethodNotAllowed, err)
}

func BatchUpdate(w http.ResponseWriter, r *http.Request) {
	tool.MyFmtPrint_Info("recv BatchUpdate message.")
	sqlcmds, err := parseHttpRequest(r)
	if err != nil {
		tool.MyFmtPrint_Error(err)
		ErrResponse(w, cstBatchUpdate, err.Error())
		return
	}

	db := GetSqliteDB()
	tx, err := db.SqliteDB.Begin()
	if err != nil {
		tool.MyFmtPrint_Error("begin event fail, err: ", err)
		ErrResponse(w, cstBatchUpdate, err.Error())
		return
	}

	err = db.BatchUpdateDirect(sqlcmds)
	if err != nil {
		tool.MyFmtPrint_Error("db update fail, err: ", err)
		ErrResponse(w, cstBatchUpdate, err.Error())
		return
	}

	err = tx.Commit()
	if err != nil{
		tool.MyFmtPrint_Error("db commit fail, err: ", err)
		ErrResponse(w, cstBatchUpdate, err.Error())
		return
	}

	succ := "(\"" + "ret" + "\"" + CstPointformat + "\"" + cstSuccess + "\")"
	message.HttpResponse(w, []*message.THttpResponseHead{
		&message.THttpResponseHead{
			HeadKey: "key",
			HeadValue: cstBatchUpdate,
		},
	}, message.StatusOK, succ)
}

func DBUpdate(w http.ResponseWriter, r *http.Request) {
	tool.MyFmtPrint_Info("recv DBUpdate message.")
	sqlcmds, err := parseHttpRequest(r)
	if err != nil {
		errStr := "parse http params fail"
		tool.MyFmtPrint_Error(err)
		ErrResponse(w, cstSingleUpdate, errStr)
		return
	}

	db := GetSqliteDB()
	tx, err := db.SqliteDB.Begin()
	if err != nil {
		errStr := "begin update event fail"
		tool.MyFmtPrint_Error(errStr + "err: ", err)
		ErrResponse(w, cstSingleUpdate, err.Error())
		return
	}

	err = db.UpdateDirect(sqlcmds)
	if err != nil {
		errStr := "db update fail"
		tool.MyFmtPrint_Error(errStr + "err: ", err)
		ErrResponse(w, cstSingleUpdate, errStr)
		return
	}

	err = tx.Commit()
	if err != nil{
		errStr := "db commit fail"
		tool.MyFmtPrint_Error(errStr + "err: ", err)
		ErrResponse(w, cstSingleUpdate, errStr)
		return
	}

	succ := "(\"" + "ret" + "\"" + CstPointformat + "\"" + cstSuccess + "\")"
	message.HttpResponse(w, []*message.THttpResponseHead{
		&message.THttpResponseHead{
			HeadKey: "key",
			HeadValue: cstSingleUpdate,
		},
	}, message.StatusOK, succ)

}

func DBQuery(w http.ResponseWriter, r *http.Request) {
	tool.MyFmtPrint_Info("recv DBQuery message.")
	tb:=time.Now().Unix()
	sqlcmds, err := parseHttpRequest(r)
	if err != nil {
		errStr := "parse http params fail"
		tool.MyFmtPrint_Error(err)
		ErrResponse(w, cstSingleQuery, errStr)
		return
	}

	db := GetSqliteDB()
	tx, err := db.SqliteDB.Begin()
	if err != nil {
		errStr := "db Transaction begin fail"
		tool.MyFmtPrint_Error("begin event fail, err: ", err)
		ErrResponse(w, cstSingleQuery, errStr)
		return
	}

	defer tx.Rollback()
	
	tool.MyFmtPrint_Info("DBQuery sqlcmds: ", sqlcmds)
	rets, err := db.Query(sqlcmds)
	if err != nil {
		errStr := "Query fail, sql: " + sqlcmds
		tool.MyFmtPrint_Error(err)
		ErrResponse(w, cstSingleQuery, errStr)
		return
	}

	if len(rets) == 0 {
		errStr := "Not any data can be queryed."
		ErrResponse(w, cstSingleQuery, errStr)
		return
	}

	tool.MyFmtPrint_Info("DBQuery query finish, rets: ", len(rets))
	var (
		pointTabledata = make([]string, len(rets))
	)
	
	if len(rets) < CstMinLimit {
		singleQueryAct(rets, 0, len(rets), &pointTabledata)
	}else{
		
		var (
			goroute int
		)

		for _, querylimit := range CstQueryLimit{
			if len(rets) >= querylimit.Min && len(rets) <= querylimit.Max {
				goroute = querylimit.GoRoute
				break
			}
		}

		tool.MyFmtPrint_Info("goroute: ", goroute)
		multiQueyAct(rets, goroute, &pointTabledata)
		
	}

	if len(pointTabledata) == 0 {
		return
	}

	var (
		dstPointTable string
	)

	for _, item := range pointTabledata {
		dstPointTable += item
	}
	te := time.Now().Unix()
	tool.MyFmtPrint_Info("spend time: ", te-tb)
	// post message.
	message.HttpResponse(w, []*message.THttpResponseHead{
		&message.THttpResponseHead{
			HeadKey: "key",
			HeadValue: cstSingleQuery,
		},
	}, message.StatusOK, dstPointTable)
	te2 := time.Now().Unix()
	tool.MyFmtPrint_Info("spend time: ", te-te2)

}

func singleQueryAct(ret []([]*TQueryField), begin int, end int, pointTabledata *[]string){
	for idx := begin; idx < end; idx++ {
		mapdata := ret[idx]
	//for idx, mapdata := range ret {
		(*pointTabledata)[idx] += "("
		for _, item := range mapdata {
			data := item.Data
			field := item.Field
			//tool.MyFmtPrint_Info("field: ", item.Field, ", data: ", data)
			if dst, ok := data.(string); ok {
				(*pointTabledata)[idx] += "(\"" + field + "\"" + CstPointformat + "\"" + dst + "\")"
			}else if dst, ok := data.(int32); ok {
				(*pointTabledata)[idx] += "(\"" + field + "\"" + CstPointformat + "\"" + strconv.Itoa(int(dst)) + "\")"
			}else if dst, ok := data.(int64); ok {
				(*pointTabledata)[idx] += "(\"" + field + "\"" + CstPointformat + "\"" + strconv.Itoa(int(dst)) + "\")"
			}else if dst, ok := data.(float64); ok {
				(*pointTabledata)[idx] += "(\"" + field + "\"" + CstPointformat + "\"" + strconv.FormatFloat(dst, 'E', -1, 64) + "\")"
			}else if dst, ok := data.(float32); ok {
				(*pointTabledata)[idx] += "(\"" + field + "\"" + CstPointformat + "\"" + strconv.FormatFloat(float64(dst), 'E', -1, 32) + "\")"
			}else if data == nil || reflect.TypeOf(data).String() == "<nil>" {
				//err = fmt.Errorf("invalid data type: %v.", reflect.TypeOf(data))
				//ErrResponse(err.Error())
				continue
			}
		}
		(*pointTabledata)[idx] += ")"
	}
}

func multiQueyAct(ret []([]*TQueryField), goroute int, pointTabledata *[]string){
	eachSep := len(ret) / goroute
	for i := 0 ;i < goroute; i++ {
		go singleQueryAct(ret, i*eachSep, eachSep*(i+1), pointTabledata)
	}

	time.Sleep(time.Duration(20*goroute)*time.Millisecond)
}