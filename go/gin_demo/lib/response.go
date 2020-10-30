package lib 

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

//统一回报数据
type ResponseObj struct {
	Code int `json:"code"`
	Msg string  `json:"msg"`
	Data interface{} `json:"data"`
}

func (obj ResponseObj) Json(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"code": obj.Code,
		"msg":  obj.Msg,
		"data": obj.Data,
	})
} 