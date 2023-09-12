package conf

import (
	"context"
	"go-chat/goChat_webSocket/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/ini.v1"
	"log"
	"strings"
)

var (
	HttpPort string

	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string

	MongoDBName string
	MongoDBAddr string
	MongoDBPwd  string
	MongoPort   string

	MongoDBClient *mongo.Client
)

func Init() {

	//本地读取环境变量

	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		panic(err)
	}
	LoadServer(file)
	LoadMysql(file)
	LoadRedis(file)
	LoadMongo(file)
	MongoDB()

	//mysql
	path := strings.Join([]string{DbUser, ":", DbPassword, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=true"}, "")
	model.Database(path)

}

func LoadServer(file *ini.File) {

	HttpPort = file.Section("service").Key("HttpPort").String()
}

func LoadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("DB").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassword = file.Section("mysql").Key("DbPassword").String()
	DbName = file.Section("mysql").Key("DbName").String()

}
func LoadRedis(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDbName").String()
	RedisAddr = file.Section("redis").Key("RedisDbName").String()
	RedisPw = file.Section("redis").Key("RedisDbName").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}
func LoadMongo(file *ini.File) {

	MongoDBName = file.Section("mongoDB").Key("MongoDBName").String()
	MongoDBAddr = file.Section("mongoDB").Key("MongoDBAddr").String()
	MongoDBPwd = file.Section("mongoDB").Key("MongoDBPwd").String()
	MongoPort = file.Section("mongoDB").Key("MongoPort").String()
}

func MongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://" + MongoDBAddr + ":" + MongoPort)
	var err error
	MongoDBClient, err = mongo.Connect(context.TODO(), clientOptions.SetAuth(
		options.Credential{
			Username: "root",
			Password: "123456",
		},
	))
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println("MongoDB Connect Successfully")
}
