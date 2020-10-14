package ginlib

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var timeFormat = "2006-01-02 15:04:05.000000000"
var ParamsLogFormatter gin.LogFormatter = func(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.ClientIP,
		param.TimeStamp.Format(timeFormat),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

var ParamsRecoveryFunc gin.RecoveryFunc = func(c *gin.Context, recovered interface{}) {
	if err, ok := recovered.(string); ok {
		c.String(http.StatusInternalServerError, "err:%v", err)
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}
