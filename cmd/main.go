package main

import (
	"go-chat/goChat_webSocket/conf"
	"go-chat/goChat_webSocket/router"
	"go-chat/goChat_webSocket/service"
)

func main() {
	conf.Init()
	go func() {
		service.Manager.Start()
	}()
	r := router.NewRouter()
	r.Run(conf.HttpPort)
}
