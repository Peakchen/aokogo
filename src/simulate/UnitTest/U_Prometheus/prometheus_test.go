package U_Prometheus

import (
	"testing"
	"net/http"
    "common/Log"
    "time"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/shirou/gopsutil/mem"
)

func TestPrometheus(t *testing.T){
	//初始一个http handler
    http.Handle("/metrics", promhttp.Handler())

    //初始化一个容器
    diskPercent := prometheus.NewGaugeVec(prometheus.GaugeOpts{
            Name: "memeory_percent",
            Help: "memeory use percent",
        },
        []string {"percent"},
    )
    prometheus.MustRegister(diskPercent)

    // 启动web服务，监听9090端口
    go func() {
        Log.FmtPrintln("ListenAndServe at:localhost:9090")
        err := http.ListenAndServe("localhost:9090", nil)
        if err != nil {
            Log.Error("ListenAndServe: ", err)
        }
    }()

    //收集内存使用的百分比
    for {
        Log.FmtPrintln("start collect memory used percent!")
        v, err := mem.VirtualMemory()
        if err != nil {
			Log.FmtPrintln("get memeory use percent error:%s", err)
			continue
        }

        Log.FmtPrintln("get memeory use percent:", v.UsedPercent)
        diskPercent.WithLabelValues("usedMemory").Set(v.UsedPercent)
        time.Sleep(time.Second*5)
    }
}
