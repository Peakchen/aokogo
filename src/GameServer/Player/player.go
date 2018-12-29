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

package Player

import(
	"github.com/gorilla/websocket"
)

type playMsg struct{
	ldbid int64
}

var playsession = make(map[int64]*websocket.Conn)

func (p *playMsg) push(ldbid int64, c *websocket.Conn){
	playsession[ldbid] = c
}

func (p *playMsg) pop(ldbid int64){
	delete(playsession, ldbid)
}