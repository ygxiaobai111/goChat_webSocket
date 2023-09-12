package api

import (
	"github.com/gin-gonic/gin"
	"go-chat/goChat_webSocket/service"
	"log"
)

func UserRegister(ctx *gin.Context) {
	var userRegisterService service.UserRegisterService
	if err := ctx.ShouldBind(&userRegisterService); err == nil {
		res := userRegisterService.Register()
		ctx.JSON(200, res)
	} else {
		ctx.JSON(400, ErrorResponse(err))
		log.Println(err)
	}
}
