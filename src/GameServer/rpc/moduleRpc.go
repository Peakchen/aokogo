package rpc

/*
	module rpc process
	date: 20191203
	author: stefan
	version: 1.0
*/

import (
	"common/Log"
	"encoding/json"
	"reflect"
)

type TModuleRpcInfo struct {
	module interface{}
	funcs  map[string]*TRpcInfo
}

var (
	_moduleRpcMap = map[string]*TModuleRpcInfo{}
)

/*
	@func: RegisterModuleRpc 注册模块rpc消息
	@param1: name 模块名
	@param2: moduleName 模块对象
*/
func RegisterModuleRpc(name string, moduleName interface{}) {
	module := reflect.ValueOf(moduleName)
	if module.Kind() != reflect.Ptr {
		Log.Error("module name: ", name, " which is not func type.")
		return
	}

	funcs := map[string]*TRpcInfo{}
	for pidx := 0; pidx < module.Type().NumMethod(); pidx++ {
		f := module.Type().Method(pidx)
		params := []reflect.Type{}
		for pidx := 0; pidx < f.Func.Type().NumIn(); pidx++ {
			params = append(params, f.Func.Type().In(pidx))
		}
		funcs[f.Name] = &TRpcInfo{
			fun:  f.Func,
			args: params,
		}
	}

	_moduleRpcMap[name] = &TModuleRpcInfo{
		module: moduleName,
		funcs:  funcs,
	}
}

/*
	@func: onModuleRpcProcess 处理模块rpc消息
	@param1: moduleName 模块名
	@param2: Rpcfunc rpc方法名
	@param3: data 数据
*/
func onModuleRpcProcess(moduleName, Rpcfunc string, data []byte) (succ bool, err error) {
	Log.FmtPrintf("rpc process, rpc module: %v, func: %v.", moduleName, Rpcfunc)
	rpcdata := _moduleRpcMap[moduleName]
	if rpcdata == nil {
		Log.Error("can not find rpc module: ", moduleName)
		return
	}

	f, exist := rpcdata.funcs[Rpcfunc]
	if !exist {
		return
	}

	funcobj := reflect.ValueOf(rpcdata.module).MethodByName(Rpcfunc)
	if funcobj.IsNil() || !funcobj.IsValid() {
		Log.Error("it is module: ", rpcdata.module, ", not find method: ", Rpcfunc)
		return
	}

	dst := reflect.New(f.args[1].Elem()).Interface()
	err = json.Unmarshal(data, dst)
	if err != nil {
		Log.Error("unmarshal data fail, Rpcfunc: ", Rpcfunc)
		return
	}

	params := []reflect.Value{
		reflect.ValueOf(dst),
	}

	funcobj.Call(params)
	succ = true
	err = nil
	return
}
