package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
	
	"gin_demo/db"
)

//首页
func Index(c *gin.Context) {
	data, err := db.GetDatabaseList()
	
	var msg string
	if err != nil {
		msg = err.Error()
	} 
	
	c.HTML(http.StatusOK, "index/index.tmpl", gin.H{"msg":msg, "data":data})
}
