package rpc

/*
	single func rpc process
	date: 20191203
	author: stefan
	version: 1.0
*/

import (
	"common/Log"
	"common/akNet"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Rpc"
	"encoding/json"
	"reflect"
)

type TRpcInfo struct {
	fun  reflect.Value  //存储函数
	args []reflect.Type //存储参数
}

var (
	_rpcMap = map[string]*TRpcInfo{}
)

/*
	@func: RegisterRpc
	@param1: name 函数名
	@param2: funcName 函数处理方法
*/
func RegisterRpc(name string, funcName interface{}) {
	f := reflect.ValueOf(funcName)
	if f.Kind() != reflect.Func {
		Log.Error("func name: ", name, " which is not func type.")
		return
	}

	params := []reflect.Type{}
	for pidx := 0; pidx < f.Type().NumIn(); pidx++ {
		params = append(params, f.Type().In(pidx))
	}

	_rpcMap[name] = &TRpcInfo{
		fun:  f,
		args: params,
	}
}

/*
	@func: SendRpcMsg 发送rpc消息
	@param1: session obj
	@param2: module, func, data
*/
func SendRpcMsg(session akNet.TcpSession, module, funcName string, data interface{}) {
	rsp := &MSG_Rpc.CS_Rpc_Req{}
	rsp.Rpcmodule = module
	rsp.Rpcfunc = funcName
	dst, err := json.Marshal(data)
	if err != nil {
		Log.Error("rpc msg marshal fail, info: ", module, funcName)
		return
	}
	rsp.Data = dst
	session.SendInnerMsg(session.GetIdentify(),
		uint16(MSG_MainModule.MAINMSG_RPC),
		uint16(MSG_Rpc.SUBMSG_CS_Rpc),
		rsp)
}

/*
	@func: onSingleRpc 处理单一函数rpc消息
	@param1: Rpcfunc 方法名
	@param2: data 数据
*/
func onSingleRpc(Rpcfunc string, data []byte) (succ bool, err error) {
	module := _rpcMap[Rpcfunc]
	if module == nil {
		Log.Error("can not find rpc module: ", Rpcfunc)
		return
	}

	dst := reflect.New(module.args[0].Elem()).Interface()
	err = json.Unmarshal(data, dst)
	if err != nil {
		Log.Error("unmarshal data fail, Rpcfunc: ", Rpcfunc)
		return
	}

	params := []reflect.Value{
		reflect.ValueOf(dst),
	}

	module.fun.Call(params)
	succ = true
	err = nil
	return
}
