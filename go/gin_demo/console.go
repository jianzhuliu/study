/*
后台入口
*/
package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"

	"gin_demo/conf"
	"gin_demo/ginlib"
	"gin_demo/routes"

	"github.com/jianzhuliu/tools/browser"
)

const (
	sign string = "console" //后台标识，用于日志
)

func init() {
	flag.Parse()
}

func getEngine() *gin.Engine {
	ginlib.DefaultSet(sign)

	r := gin.New()

	//定制中间件参数
	r.Use(gin.LoggerWithFormatter(ginlib.ParamsLogFormatter))
	r.Use(gin.CustomRecovery(ginlib.ParamsRecoveryFunc))

	//加载模板文件
	r.LoadHTMLGlob("templates/*/*")

	//静态文件
	r.Static("/static", "./static")

	return r
}

func main() {
	//创建 engine
	r := getEngine()

	//设置路由
	routes.SetRoutes(r)

	//浏览器自动开启

	addr := fmt.Sprintf("%s:%d", conf.V_host, conf.V_port)
	go browser.OpenWithNotice(addr)

	//启动
	r.Run(addr)
}
