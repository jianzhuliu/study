package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//首页
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index/index.tmpl", gin.H{"msg":"欢迎你！！！"})
}
