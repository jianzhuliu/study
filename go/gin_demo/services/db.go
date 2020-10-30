package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
	
	"gin_demo/db"
)

//展示数据库表列表
func Db(c *gin.Context) {
	dbname := c.Param("dbname")
	
	data, err := db.GetTableList(dbname)
	
	var msg string
	if err != nil {
		msg = err.Error()
	} 
	
	c.HTML(http.StatusOK, "db/index.tmpl", gin.H{"msg":msg, "dbname":dbname,"data":data})
}
