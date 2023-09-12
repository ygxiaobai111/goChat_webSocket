package service

import (
	"context"
	"fmt"
	"go-chat/goChat_webSocket/conf"
	"go-chat/goChat_webSocket/model/ws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type SebdSortMsg struct {
	Content  string `json:"content"`
	Read     uint   `json:"read"`
	CreateAt int64  `json:"create_at"`
}

func InsertMsg(database, id, content string, read uint, expire int64) (err error) {
	//插入mongoDB
	collection := conf.MongoDBClient.Database(database).Collection(id) //没有这个id集合将自动创建
	comment := ws.Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      read,
	}
	_, err = collection.InsertOne(context.TODO(), comment)
	return
}
func FindMany(database, sendID, id string, time int64, pageSize int) (results []ws.Result, err error) {
	var resultMe []ws.Trainer  //id
	var resultYou []ws.Trainer // sendId

	filter := bson.D{} // 空文档表示无条件
	opts := options.Find()
	opts.SetSort(bson.D{{"startTime", -1}}) // 通过 startTime 倒序排序
	opts.SetLimit(int64(pageSize))

	db := conf.MongoDBClient.Database(database)
	sendIdCollection := db.Collection(sendID)
	// 限制返回 10 条结果
	cursor, _ := sendIdCollection.Find(context.TODO(), filter, opts)

	err = cursor.All(context.TODO(), &resultYou)
	if err != nil {
		log.Println(err)
	}
	idCollection := db.Collection(id)

	// 限制返回 10 条结果
	cursor, _ = idCollection.Find(context.TODO(), filter, opts)
	err = cursor.All(context.TODO(), &resultMe)
	if err != nil {
		log.Println(err)
	}

	results = AppendAndSort(resultMe, resultYou)
	return
}

func AppendAndSort(resultMe, resultYou []ws.Trainer) (results []ws.Result) {
	for _, r := range resultMe {
		sendSort := SebdSortMsg{ //构造返回信息
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{ //构造返回的所有内容，包括发送者
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "me",
		}
		results = append(results, result)
	}

	for _, r := range resultYou {
		sendSort := SebdSortMsg{ //构造返回信息
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{ //构造返回的所有内容，包括发送者
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "you",
		}

		results = append(results, result)
	}

	return results
}

// FindUnread
func FindUnread(database, sendID string, time int64, pageSize int) (results []ws.Result, err error) {
	var resultYou []ws.Trainer // sendId

	filter := bson.D{} // 空文档表示无条件
	opts := options.Find()
	opts.SetSort(bson.D{{"startTime", -1}}) // 通过 startTime 倒序排序
	opts.SetLimit(int64(pageSize))
	//TODO:添加条件，只获取未读消息
	db := conf.MongoDBClient.Database(database)
	sendIdCollection := db.Collection(sendID)

	// 限制返回 10 条结果
	cursor, _ := sendIdCollection.Find(context.TODO(), filter, opts)

	err = cursor.All(context.TODO(), &resultYou)
	if err != nil {
		log.Println(err)
	}
	results = FindUnreadAppend(resultYou)
	return results, nil
}
func FindUnreadAppend(resultYou []ws.Trainer) (results []ws.Result) {
	for _, r := range resultYou {
		sendSort := SebdSortMsg{ //构造返回信息
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{ //构造返回的所有内容，包括发送者
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "you",
		}

		results = append(results, result)
	}

	return results
}
