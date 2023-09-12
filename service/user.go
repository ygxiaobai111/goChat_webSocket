package service

import (
	"context"
	"go-chat/goChat_webSocket/model"
	"go-chat/goChat_webSocket/serializer"
)

type UserRegisterService struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func (service *UserRegisterService) Register() serializer.Response {
	var user model.User
	code := 200
	var count int64
	model.NewDBClient(context.Background()).Model(model.User{}).Where("user_name=?", service.UserName).First(&user).Count(&count)

	if count != 0 {
		code = 400
		return serializer.Response{
			Status: code,
			Msg:    "该用户名已被注册",
		}
	}
	user = model.User{

		UserName: service.UserName,
	}
	if err := user.SetPassword(service.Password); err != nil {
		code = 400
		return serializer.Response{
			Status: code,
			Msg:    "该用户名已被注册",
		}
	}

	model.NewDBClient(context.Background()).Create(&user)
	return serializer.Response{
		Status: 200,

		Msg: "注册成功",
	}
}
