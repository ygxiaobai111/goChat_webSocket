package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go-chat/goChat_webSocket/conf"
	"go-chat/goChat_webSocket/pkg/e"
)

func (manager *ClientManager) Start() {
	for {
		fmt.Println("---监听管道")
		select {
		case conn := <-Manager.Register:
			fmt.Printf("有新连接 %v\n ", conn.ID)
			Manager.Clients[conn.ID] = conn //将连接放入用户管理
			replyMsg := ReplyMsg{

				Code:    e.WebsocketSuccess,
				Content: "已经连接到服务器",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
		case conn := <-Manager.Unregister:
			fmt.Printf("连接失败%s", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok { //查看连接中是否有内容
				replyMsg := ReplyMsg{

					Code:    e.WebsocketError,
					Content: "连接中断",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}

			//有人发消息
		case brodcast := <-Manager.Broadcast: //1->2
			message := brodcast.Message
			sendId := brodcast.Client.SendID //2->1 二给一的链接
			flag := false                    //默认对方不在线
			//Manager.Clients是用户连接表
			for id, conn := range Manager.Clients {
				if id != sendId { //如果在该列表找到，对方就是在线
					continue
				}
				select {
				//将消息放入通道
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}
			id := brodcast.Client.ID //发送方id 1-》2
			if flag {                //对方在线
				replyMsg := ReplyMsg{

					Code:    e.WebsocketSuccess,
					Content: "对方在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = brodcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)        //广播给自己，告诉用户对方在线应答
				err := InsertMsg(conf.MongoDBName, id, string(message), 1, int64(3*month)) //1是已读的意思
				if err != nil {
					fmt.Println("inset Err", err)
				}
			} else {
				fmt.Println("对方不在线")
				replyMsg := ReplyMsg{

					Code:    e.WebsocketSuccess,
					Content: "对方不在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = brodcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)        //广播给自己，告诉用户对方在线应答
				err := InsertMsg(conf.MongoDBName, id, string(message), 0, int64(3*month)) //0是未读的意思
				if err != nil {
					fmt.Println("inset Err", err)
				}
			}

		}

	}
}
