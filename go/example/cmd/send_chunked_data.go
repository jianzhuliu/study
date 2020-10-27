package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"gitee.com/jianzhuliu/tools/browser"
	"gitee.com/jianzhuliu/tools/common"
	"gitee.com/jianzhuliu/tools/conf"
	"gitee.com/jianzhuliu/tools/logger"
)

func init() {
	//初始化定义相关目录，比如日志文件夹
	err := common.DefaultInit("send_chunked_data")
	if err != nil {
		panic(fmt.Sprint("common.InitConfPath()", "--", err.Error()))
	}

	flag.Parse()
}

func main() {
	addr := fmt.Sprintf("%s:%d", conf.V_host, conf.V_port)

	http.HandleFunc("/", handleIndex)

	//另起 goroutine 用于浏览器自动打开
	go browser.OpenWithNotice(addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.PrintfWithTime("start serve[%s] failed,err=%v", addr, err.Error())
	}
}

func handleIndex(rw http.ResponseWriter, r *http.Request) {
	header := rw.Header()
	header.Set("Transfer-Encoding", "chunked")
	header.Set("Content-Type", "text/html")

	rw.WriteHeader(http.StatusOK)

	rw.Write(common.StringToBytes("<html><body>"))
	rw.(http.Flusher).Flush()

	for i := 0; i < 5; i++ {
		rw.Write(common.StringToBytes(fmt.Sprintf("<h1>%d</h1>", i)))
		rw.(http.Flusher).Flush()
		time.Sleep(time.Duration(1) * time.Second)
	}

	rw.Write(common.StringToBytes("</body></html>"))
	rw.(http.Flusher).Flush()
}
