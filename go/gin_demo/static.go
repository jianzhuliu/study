/*
静态文件处理服务
*/
package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	host *string = flag.String("host", "127.0.0.1", "setting the host")
	port *int    = flag.Int("port", 3001, "setting the port")
)

func init() {
	flag.Parse()
}

func main() {
	addr := fmt.Sprintf("%s:%d", *host, *port)
	http.Handle("/static/", http.FileServer(http.Dir("./")))
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("[%s] start fail -- %v \n", addr, err.Error())
	}
}
