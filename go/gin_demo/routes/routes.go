package routes

import (
	"gin_demo/services"
	
	"github.com/gin-gonic/gin"
)

//路由配置
func SetRoutes(r *gin.Engine) {
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "123456",
	}))

	authorized.GET("index", services.Index)
	authorized.GET("", services.Login)
	authorized.POST("login", services.PostLogin)
}
