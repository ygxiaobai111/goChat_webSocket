package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-chat/goChat_webSocket/cache"
	"go-chat/goChat_webSocket/conf"
	"go-chat/goChat_webSocket/pkg/e"
	"log"
	"net/http"
	"strconv"
	"time"
)

const month = 60 * 60 * 24 * 30

type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string
}
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

// 广播
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// 用户管理
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = &ClientManager{
	Clients:    make(map[string]*Client), //最大连接数
	Broadcast:  make(chan *Broadcast),
	Reply:      make(chan *Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func CreateID(uid, toUid string) string {
	return uid + "->" + toUid
}
func Handler(ctx *gin.Context) {
	uid := ctx.Query("uid")
	toUid := ctx.Query("toUid")
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(ctx.Writer, ctx.Request, nil) //升级ws协议
	if err != nil {
		log.Println("err:", err)
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	//创建用户实例
	client := &Client{
		ID:     CreateID(uid, toUid), //1->2
		SendID: CreateID(toUid, uid), //2->1
		Socket: conn,
		Send:   make(chan []byte),
	}
	//用户注册到用户管理
	Manager.Register <- client
	go client.Read()
	go client.Write()
}

// 读取用户传入
func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		ctx := context.Background()
		c.Socket.PongHandler()
		sendMsg := new(SendMsg)
		//c.Socket.ReadMessage()
		err := c.Socket.ReadJSON(&sendMsg)
		if err != nil {
			log.Println("数据格式不正确", err)
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		if sendMsg.Type == 1 {
			//发送消息 1-》2
			r1, _ := cache.RedisClient.Get(ctx, c.ID).Result()     //1->2
			r2, _ := cache.RedisClient.Get(ctx, c.SendID).Result() // 2->1
			if r1 > "3" && r2 == "" {                              //一给二发消息，发了三条，但二没回，或没看见，就停止1发送
				replyMsg := ReplyMsg{
					From:    "",
					Code:    e.WebsocketError,
					Content: "达到上限",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			} else {
				cache.RedisClient.Incr(ctx, c.ID)
				_, _ = cache.RedisClient.Expire(ctx, c.ID, time.Hour*24*30*3).Result()
				//防止过快分离，建立连接三个月过期

			}
			//加入通道
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(sendMsg.Content), //发送过来的消息
			}

		} else if sendMsg.Type == 2 {
			// 获取历史消息
			timeT, err := strconv.Atoi(sendMsg.Content) //获取历史消息时这里传进来的是时间戳
			if err != nil {
				timeT = 999999
			}
			results, _ := FindMany(conf.MongoDBName, c.SendID, c.ID, int64(timeT), 10) //获取历史消息十条
			fmt.Println("id:", c.SendID, c.ID)
			if len(results) > 10 {
				results = results[10:]
			} else if len(results) == 0 {
				replyMsg := ReplyMsg{

					Code:    e.WebsocketError,
					Content: "上面没有消息了",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: result.Msg,
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}

		} else if sendMsg.Type == 3 { //获取所有未读消息
			// 获取历史消息
			timeT, err := strconv.Atoi(sendMsg.Content) //获取历史消息时这里传进来的是时间戳
			if err != nil {
				timeT = 999999
			}
			results, _ := FindUnread(conf.MongoDBName, c.SendID, int64(timeT), 10) //获取历史消息十条

			if len(results) > 10 {
				results = results[10:]
			} else if len(results) == 0 {
				replyMsg := ReplyMsg{

					Code:    e.WebsocketError,
					Content: "上面没有对方的未读消息了",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: result.Msg,
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		}

	}
}
func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			replyMsg := ReplyMsg{

				Code:    e.WebsocketSuccess,
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
