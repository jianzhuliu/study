/*
公共常量配置入口
*/

package conf

import (
	"os"
)

const (
	C_PERM           os.FileMode = 0644                                  //默认文件权限
	C_OPEN_FILE_FLAG             = os.O_CREATE | os.O_APPEND | os.O_RDWR //打开文件默认模式
)
