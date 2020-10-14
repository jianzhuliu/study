/*
配置文件初始化配置
*/

package conf

import (
	"os"
	"flag"
	"path/filepath"

	"gin_demo/lib"
)

func init() {
	//目录配置
	curPath, err := os.Getwd()
	if err != nil {
		lib.ShowPanic("os.Getwd()", err)
	}

	V_pathRoot = curPath
	V_pathLogs = filepath.Join(V_pathRoot, "logs")

	os.MkdirAll(V_pathLogs, C_PERM)
	
	flag.StringVar(&V_host,"host","127.0.0.1","setting the host")
	flag.IntVar(&V_port, "port", 3000, "setting the port")
}
