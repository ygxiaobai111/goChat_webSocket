package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/goChat_webSocket/api"
	"go-chat/goChat_webSocket/service"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	//该中间件用于程序恢复500错误码引起的恐慌，Default自带这两个
	//r.Use(gin.Recovery(), gin.Logger())
	v1 := r.Group("/")
	{
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(200, "pong!")
		})
		v1.POST("user/register", api.UserRegister)
		v1.GET("ws", service.Handler)
	}
	return r

}
