package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
	
	"gin_demo/db"
)

//展示数据库表信息
func Table(c *gin.Context) {
	dbname := c.Param("dbname")
	tblname := c.Param("tblname")
	
	data, err := db.GetTableDesc(dbname,tblname)
	
	var msg string
	if err != nil {
		msg = err.Error()
	} 
	
	c.HTML(http.StatusOK, "table/desc.tmpl", gin.H{"msg":msg,"data":data})
}
