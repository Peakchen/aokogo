/*
* CopyRight(C) StefanChen e-mail:2572915286@qq.com
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/
package memoryMgr

import(
	"os"
	"log"
	"runtime/pprof"
	"os/signal"
	//"syscall"
)

const(
	LP_GOROUTINE = 1
)

var(
	m_mapRunInfo = map[int32]string{
		LP_GOROUTINE : "goroutine",
	}
)

func LoopWaitforSignal(){
	for{
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt,os.Kill)

	}
}

func LookupInfo(id int32){
	var lstr, ok = m_mapRunInfo[id]
	if !ok {
		log.Printf("look up err id: %d.", id)
		return
	}

	pprof.Lookup(lstr).WriteTo(os.Stdout, 2)
}