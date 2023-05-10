package inits

import (
	"context"
	"fmt"
	"webserver/infrastructure/config"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type InitObj struct {
	Server      *gin.Engine
	Ctx         context.Context
	Cfg         *config.Config
	Mongoclient *mongo.Client
	Redisclient *redis.Client
}

var (
	initObj *InitObj
	err     error
)

func GetinitObj() *InitObj {
	if initObj == nil {
		initObj = newInitObj()
	}
	return initObj
}

func SetInitObj(obj *InitObj) {
	initObj = obj
}

func newInitObj() *InitObj {
	initObj = &InitObj{}
	// 設定initObj中的Ctx
	initObj.Ctx = context.TODO()
	// 設定initObj中的Cfg
	setInitObjCfg()
	// 設定initObj中的Server
	setInitObjServer()
	// 設定initObj中的Mongoclient
	setInitObjMongoclient()
	// 設定initObj中的Redisclient
	setInitObjRedisclient()
	return initObj
}

// 讀取infrastructure/config/app.env
func setInitObjCfg() {
	if initObj.Cfg, err = config.LoadConfig("./infrastructure/config/"); err != nil {
		panic(err)
	}
}

// 使用Gin來建立Server
func setInitObjServer() {
	initObj.Server = gin.Default()
}

// 連線並驗證Mongodb是否正常連線
func setInitObjMongoclient() {
	mongoconn := options.Client().ApplyURI(initObj.Cfg.DBUri)
	// 新增一個MongodbClient並寫入initObj中的Mongodbclient參數
	if initObj.Mongoclient, err = mongo.Connect(initObj.Ctx, mongoconn); err != nil {
		panic(err)
	}
	// 檢查連線是否正常
	if err := initObj.Mongoclient.Ping(initObj.Ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	// 成功訊息
	fmt.Println("MongoDB successfully connected...")
}

// 連線並驗證Redis是否正常連線
func setInitObjRedisclient() {
	// 新增一個RedisClient並寫入initObj中的Redisclient參數
	initObj.Redisclient = redis.NewClient(&redis.Options{
		Addr: initObj.Cfg.RedisUri,
	})
	// 檢查連線是否正常
	if _, err := initObj.Redisclient.Ping(initObj.Ctx).Result(); err != nil {
		panic(err)
	}
	// 成功訊息
	fmt.Println("Redis client connected successfully...")
}
