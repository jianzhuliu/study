package main

import (
	"flag"
	
	"gitee.com/jianzhuliu/tools/common"
	"gitee.com/jianzhuliu/tools/conf"
	"gitee.com/jianzhuliu/tools/logger"
)

func init(){
	//初始化定义相关目录，比如日志文件夹
	err := common.DefaultInit()
	if err != nil {
		common.ExitErr("common.DefaultInit()", err)
	}

	flag.Parse()
	
	//初始化数据库
	err = common.InitDb()
	if err != nil {
		common.ExitErr("common.InitDb()", err)
	}
}

func main(){
	db := conf.V_db
	
	var version string 
	err := db.QueryRow("select version()").Scan(&version)
	if err != nil {
		common.ExitErr("db.QueryRow()",err)
	}
	
	logger.PrintfWithTime("%s", version)
}