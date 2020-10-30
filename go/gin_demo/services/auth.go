package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"crypto/md5"
	
	"gin_demo/db"
	"gin_demo/lib"
	"fmt"
)

//登录接口数据
type loginData struct {
	Host string `form:"host" binding:"required"`
	Port int `form:"port" binding:"required"`
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:`
	Token string `form:"t" binding:"required"`
}

//登录界面
func Login(c *gin.Context) {
	token := fmt.Sprintf("%x", md5.Sum([]byte("asdfasdf")))
	c.HTML(http.StatusOK, "auth/login.tmpl", gin.H{"token":token})
}

//登录验证
func PostLogin(c *gin.Context) {
	var data loginData
	var responseObj lib.ResponseObj 
	var err error
	if err = c.ShouldBind(&data); err != nil {
		responseObj.Code = 2
		responseObj.Msg = err.Error()
		responseObj.Json(c)
		
		return
	}

	err = db.Setup(data.Username,data.Password,data.Host,data.Port)

	if err != nil {
		responseObj.Code = 3
		responseObj.Msg = err.Error()
		responseObj.Json(c)
		
		return
	}
	
	responseObj.Code = 0
	responseObj.Msg = "成功登录"
	responseObj.Json(c)
}
