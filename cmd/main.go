package main

import (
	"go-chat/goChat_webSocket/conf"
	"go-chat/goChat_webSocket/router"
	"go-chat/goChat_webSocket/service"
	"log"
)

func main() {
	conf.Init()
	go func() {
		service.Manager.Start()
	}()

	err := service.LikeArticle("65098c8acab6ce73bf3665b4")
	log.Println("err:", err)

	r := router.NewRouter()
	r.Run(conf.HttpPort)

}
