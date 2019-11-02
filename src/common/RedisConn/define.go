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

package RedisConn

import (
	"time"
)

// pool Idl
const (
	IDle_Invalie = iota
	IDle_one
	IDle_two
	IDle_three
	IDle_four
	IDle_five
)

// idl timeout value
const (
	IDleTimeOut_invalid = iota
	//second
	IDleTimeOut_five_sec = 5 * time.Second
	IDleTimeOut_ten_sec  = 10 * time.Second
	//minute
	IDleTimeOut_one_min   = 60 * time.Second
	IDleTimeOut_two_min   = 120 * time.Second
	IDleTimeOut_three_min = 180 * time.Second
	IDleTimeOut_four_min  = 240 * time.Second
)

const (
	// common used MilliSecond
	MSec_one         = time.Millisecond
	MSec_ten         = 10 * time.Millisecond
	MSec_one_Hundred = 100 * time.Millisecond

	// common used Second
	Sec_five    = 5 * time.Second
	Sec_ten     = 10 * time.Second
	Sec_fifteen = 15 * time.Second
	Sec_twenty  = 20 * time.Second
	Sec_thirty  = 30 * time.Second
	Sec_fourty  = 40 * time.Second
	Sec_fifty   = 50 * time.Second
	Sec_sixty   = 60 * time.Second
)

type REDIS_INT32 int32

const (
	REDIS_SET_DEADLINE REDIS_INT32 = 600 //s

)
