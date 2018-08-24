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

package RedisEx

import(
	"time"
)

type RedisInfo struct{
	nServerPort	int32
	szRedisServerHost string
	szServerUUID	string
}

// pool Idl
const(
	IDle_Invalie = iota
	IDle_one
	IDle_two
	IDle_three
	IDle_four
	IDle_five
)

// idl timeout value
const(
	IDleTimeOut_invalid = iota
	//second
	IDleTimeOut_five_sec = 5*time.Second
	IDleTimeOut_ten_sec = 10*time.Second
	//minute
	IDleTimeOut_one_min = 60*time.Second
	IDleTimeOut_two_min = 120*time.Second
	IDleTimeOut_three_min = 180*time.Second
	IDleTimeOut_four_min = 240*time.Second

)

const(
	// common used MilliSecond
	MSec_one = time.Millisecond
	MSec_ten = 10*time.Millisecond
	MSec_one_Hundred = 100*time.Millisecond

	// common used Second 
	Sec_five = 5*time.Second
	Sec_ten = 10*time.Second
	Sec_fifteen = 15*time.Second
	Sec_twenty = 20*time.Second
	Sec_thirty = 30*time.Second
	Sec_fourty = 40*time.Second
	Sec_fifty = 50*time.Second
	Sec_sixty = 60*time.Second
)