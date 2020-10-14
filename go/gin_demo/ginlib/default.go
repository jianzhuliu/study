/*
gin 相关公用设置
*/

package ginlib

import (
	"fmt"
	"gin_demo/conf"
	"gin_demo/lib"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
)

func DefaultSet(sign string) {
	//日志输出禁止输出颜色
	gin.DisableConsoleColor()

	//日志及错误输出配置
	f, err := os.OpenFile(filepath.Join(conf.V_pathLogs, fmt.Sprintf("%s_%d.log", sign, conf.V_port)), conf.C_OPEN_FILE_FLAG, conf.C_PERM)
	if err != nil {
		lib.ShowPanic("os.OpenFile() -- log", err)
	}

	ferror, err := os.OpenFile(filepath.Join(conf.V_pathLogs, fmt.Sprintf("%s_%d_error.log", sign, conf.V_port)), conf.C_OPEN_FILE_FLAG, conf.C_PERM)
	if err != nil {
		lib.ShowPanic("os.OpenFile() -- log", err)
	}

	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	gin.DefaultErrorWriter = io.MultiWriter(ferror, os.Stdout)
}
