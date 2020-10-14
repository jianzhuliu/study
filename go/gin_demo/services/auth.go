package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//登录接口数据
type loginData struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:`
}

//登录界面
func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/login.tmpl", gin.H{})
}

//登录验证
func PostLogin(c *gin.Context) {
	var data loginData
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 2,
			"msg":  err.Error(),
		})

		return
	}

	code := 3
	msg := "账号或者密码错误，请重试"
	if data.Username == "admin" && data.Password == "123456" {
		code = 0
		msg = "成功登录"
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}
