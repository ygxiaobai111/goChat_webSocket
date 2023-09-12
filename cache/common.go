package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gopkg.in/ini.v1"
	"strconv"
)

var (
	RedisClient *redis.Client
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

func init() {
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Println("redis config err", err)
	}
	LoadRedisData(file)
	Redis()
}

func LoadRedisData(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}

func Redis() {
	db, _ := strconv.ParseUint(RedisDbName, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPw,
		DB:       int(db),
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	RedisClient = client

}
