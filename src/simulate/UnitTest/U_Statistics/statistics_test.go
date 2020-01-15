package U_Statistics

import (
	"testing"
	"common/ado/dbStatistics"
	"time"
)

func stack_1(){
	stack_2()
	stack_3()
}

func stack_2(){
	stack_3()
}

func stack_3() {
	dbStatistics.DBOperStatistics("1", "hello")
	dbStatistics.DBMsgStatistics("1",1,2)
}

func TestStatistics(t *testing.T){
	dbStatistics.InitDBStatistics()
	stack_1()
	time.Sleep(time.Duration(5)*time.Second)
	dbStatistics.DBStatisticsStop()
}