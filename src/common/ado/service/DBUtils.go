package service

/*
Copyright (c) <year> <copyright holders>

"Anti 996" License Version 1.0 (Draft)

Permission is hereby granted to any individual or legal entity
obtaining a copy of this licensed work (including the source code,
documentation and/or related items, hereinafter collectively referred
to as the "licensed work"), free of charge, to deal with the licensed
work for any purpose, including without limitation, the rights to use,
reproduce, modify, prepare derivative works of, distribute, publish
and sublicense the licensed work, subject to the following conditions:

1. The individual or the legal entity must conspicuously display,
without modification, this License and the notice on each redistributed
or derivative copy of the Licensed Work.

2. The individual or the legal entity must strictly comply with all
applicable laws, regulations, rules and standards of the jurisdiction
relating to labor and employment where the individual is physically
located or where the individual was born or naturalized; or where the
legal entity is registered or is operating (whichever is stricter). In
case that the jurisdiction has no such laws, regulations, rules and
standards or its laws, regulations, rules and standards are
unenforceable, the individual or the legal entity are required to
comply with Core International Labor Standards.

3. The individual or the legal entity shall not induce, metaphor or force
its employee(s), whether full-time or part-time, or its independent
contractor(s), in any methods, to agree in oral or written form, to
directly or indirectly restrict, weaken or relinquish his or her
rights or remedies under such laws, regulations, rules and standards
relating to labor and employment as mentioned above, no matter whether
such written or oral agreement are enforceable under the laws of the
said jurisdiction, nor shall such individual or the legal entity
limit, in any methods, the rights of its employee(s) or independent
contractor(s) from reporting or complaining to the copyright holder or
relevant authorities monitoring the compliance of the license about
its violation(s) of the said license.

THE LICENSED WORK IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE COPYRIGHT HOLDER BE LIABLE FOR ANY CLAIM,
DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN ANY WAY CONNECTION WITH THE
LICENSED WORK OR THE USE OR OTHER DEALINGS IN THE LICENSED WORK.
*/

import (
	"common/ado"
	"common/public"
)

/*
	@func: Update
	@param1: Identify string     player only one
	@param2: data public.IDBCache  module need save
	@param3: Oper ado.EDBOperType  db operation
	purpose: first db save data then update cache.
*/
func (this *TDBProvider) Update(Identify string, data public.IDBCache, Oper ado.EDBOperType) (err error) {
	var cacheOper bool
	err, cacheOper = this.rconn.Update(Identify, data, Oper)
	if err != nil || cacheOper {
		return
	}

	err = this.mconn.SaveOne(Identify, data)
	if err != nil {
		return
	}

	return
}

/*
	@func: Insert
	@param1: Identify string     player only one
	@param2: data public.IDBCache  module need save
	purpose: first insert data to cache then update db.
*/
func (this *TDBProvider) Insert(Identify string, data public.IDBCache) (err error) {
	err = this.rconn.Insert(Identify, data)
	if err == nil {
		err = this.mconn.InsertOne(Identify, data)
	}
	return
}

/*
	@func: Get
	@param1: Identify string     player only one
	@param2: Output public.IDBCache  module need query
	purpose: first query from cache if not exist, then find from db.
*/
func (this *TDBProvider) Get(Identify string, Output public.IDBCache) (err error, exist bool) {
	err = this.rconn.Query(Identify, Output)
	if err != nil {
		err, exist = this.mconn.QueryOne(Identify, Output)
		//redis not exist, then update
		err, _ = this.rconn.Update(Identify, Output, ado.EDBOper_Update)
	} else {
		exist = true
	}
	return
}

/*
	@func: GetAcc
	@param1: Identify string     player only one
	@param2: Output public.IDBCache  module need query
	purpose: first query use accout from db.
*/
func (this *TDBProvider) GetAcc(usrName string, Output public.IDBCache) (err error, exist bool) {
	err, exist = this.mconn.QueryAcc(usrName, Output)
	return
}

/*
	@func: DBGetSome
	@param1: Identify string     player only one
	@param2: Output public.IDBCache  module need query
	purpose: first query from cache if not then db.
*/
func (this *TDBProvider) DBGetSome(Output public.IDBCache) (err error) {
	// has no func need.
	return nil
}
