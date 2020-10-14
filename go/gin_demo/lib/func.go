/*
公用函数配置
*/

package lib 

import (
	"fmt"
)

//统一启动错误时panic
func ShowPanic(msg string, err error){
	panic(fmt.Sprintf("[%s] err:%s", msg, err.Error()))
}
