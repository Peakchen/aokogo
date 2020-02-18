// add by stefan

package main

import (
	. "BigWorld/define"
	. "BigWorld/mgr"
)

func init() {

}

func main() {
	bigw := NewBigWord(ConstBigWorldHost)
	bigw.Run()
}
