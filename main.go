package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"webserver/infrastructure/config"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// 創建稍後將重新分配的必要變數
var (
	cfg         *config.Config
	ctx         context.Context
	server      *gin.Engine
	mongoclient *mongo.Client
	redisclient *redis.Client
	tableitem   *mongo.Collection
	err         error
	zsetKey     = "rank"
)

type Item struct {
	Iid   int    `bson:"iid"`
	Name  string `bson:"name"`
	Price int    `bson:"price"`
	Text  string `bson:"text"`
}

// 將在"main"函數之前運行的初始化函數
func init() {
	// 加載 .env 變量
	cfg, err = config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	// 創建 context
	ctx = context.TODO()

	// 創建 MongoDB
	initMongodb()

	// 創建 Redis
	initRedis()
	//初始化排行榜
	initRank()
	// 創建 Gin 引擎
	server = gin.Default()
}

func initMongodb() {
	mongoconn := options.Client().ApplyURI(cfg.DBUri)
	//連線到Mongodb
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		panic(err)
	}
	//檢查連線
	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	//指定要操作的数据集
	// tableuser = mongoclient.Database("shopweb").Collection("user")
	tableitem = mongoclient.Database("shopweb").Collection("item")
	//成功訊息
	fmt.Println("MongoDB successfully connected...")
}

func initRedis() {
	//連線到Redis
	redisclient = redis.NewClient(&redis.Options{
		Addr: cfg.RedisUri,
	})
	//檢查連線
	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		panic(err)
	}
	//成功訊息
	fmt.Println("Redis client connected successfully...")
}

func initRank() {
	companys := []*redis.Z{
		{Score: 10, Member: "HERMES"},
		{Score: 20, Member: "Gucci"},
		{Score: 30, Member: "CHANEL"},
		{Score: 40, Member: "LV"},
		{Score: 50, Member: "Dior"},
		{Score: 60, Member: "Prada"},
	}
	err := redisclient.ZAdd(ctx, zsetKey, companys...).Err()
	if err != nil {
		panic(err)
	}
}

func main() {
	//完成後斷開與Mongodb連線
	defer mongoclient.Disconnect(ctx)

	router := server.Group("/api")

	router.GET("item", getItem)
	router.GET("item/:id", getItemByID)
	router.POST("item", addItem)
	router.PUT("item/:id", updateItem)
	router.DELETE("item/:id", deleteItem)
	router.GET("rank", getRank)
	router.PUT("rank", updateRank)

	log.Fatal(server.Run(":" + cfg.Port))
}

func getItem(ctx *gin.Context) {
	//限制查找後返回資料的數量
	findOptions := options.Find()
	findOptions.SetLimit(50)
	//查詢條件
	filter := bson.D{{}}
	//資料查找
	cur, err := tableitem.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	//包裝成Item
	var items []Item
	if err = cur.All(context.TODO(), &items); err != nil {
		log.Fatal(err)
	}
	//成功
	httpSuccessJson(ctx, items)
}

func getItemByID(ctx *gin.Context) {
	//擷取到值轉換成int
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "no_input", "message": "The input is not an integer."})
		return
	}
	//查詢條件
	filter := bson.D{primitive.E{Key: "iid", Value: id}}
	//單筆查詢
	var item Item
	if err = tableitem.FindOne(context.TODO(), filter).Decode(&item); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "no_data", "message": "Data not found."})
		return
	}
	//成功
	httpSuccessJson(ctx, item)
}

func addItem(ctx *gin.Context) {
	//擷取id
	id, err := strconv.Atoi(ctx.PostForm("iid"))
	if err != nil {
		httpErrorJson(ctx, err, "input_type_error", "price is not integer")
		return
	}
	//擷取資料
	newItem, errMsg, err := paramItemData(ctx, id)
	if err != nil {
		httpErrorJson(ctx, err, "input_type_error", errMsg)
		return
	}
	//插入一條數據
	res, err := tableitem.InsertOne(context.TODO(), newItem)
	if err != nil {
		httpErrorJson(ctx, err, "add_fail", "add failed.")
		return
	}
	httpSuccessJson(ctx, res.InsertedID)
}

func updateItem(ctx *gin.Context) {
	//擷取到值轉換成int
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "no_input", "message": "The input is not an integer."})
		return
	}
	//擷取資料
	newItem, errMsg, err := paramItemData(ctx, id)
	if err != nil {
		httpErrorJson(ctx, err, "input_type_error", errMsg)
		return
	}
	filter := bson.D{primitive.E{Key: "iid", Value: id}}
	update := bson.D{primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{Key: "name", Value: newItem.Name},
			primitive.E{Key: "price", Value: newItem.Price},
			primitive.E{Key: "text", Value: newItem.Text},
		},
	}}
	// updateOpts := options.Update().SetUpsert(true) // 设置upsert模式
	updateResult, err := tableitem.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		httpErrorJson(ctx, err, "update_error", errMsg)
		return
	}
	httpSuccessJson(ctx, updateResult)
}

func deleteItem(ctx *gin.Context) {
	// 擷取到值轉換成int
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httpErrorJson(ctx, err, "no_input", "The input is not an integer.")
		return
	}
	// 刪除一條數據
	filter := bson.D{primitive.E{Key: "iid", Value: id}}
	deleteResult, err := tableitem.DeleteOne(context.TODO(), filter)
	if err != nil {
		httpErrorJson(ctx, err, "fail_del", "delete fail.")
	}
	httpSuccessJson(ctx, deleteResult)
}

func paramItemData(ctx *gin.Context, id int) (Item, string, error) {
	name := ctx.PostForm("name")
	price, err := strconv.Atoi(ctx.PostForm("price"))
	if err != nil {
		return Item{}, "price is not integer", err
	}
	text := ctx.DefaultPostForm("text", "無")
	return Item{id, name, price, text}, "", nil
}

// 查詢排行榜
func getRank(ctx *gin.Context) {
	//取得前3名排行
	rank, err := redisclient.ZRangeWithScores(ctx, zsetKey, -3, -1).Result()
	if err != nil {
		httpErrorJson(ctx, err, "no_rank", "get rank fail.")
		return
	}
	//成功
	httpSuccessJson(ctx, rank)
}

func updateRank(ctx *gin.Context) {
	Score, err := strconv.ParseFloat(ctx.PostForm("score"), 64)
	if err != nil {
		httpErrorJson(ctx, err, "input_type_error", "score is not integer.")
		return
	}
	name := ctx.PostForm("name")
	companys := []*redis.Z{{Score: Score, Member: name}}
	err = redisclient.ZAdd(ctx, zsetKey, companys...).Err()
	if err != nil {
		fmt.Printf("redis ZAdd fail,error msg:%v\n", err)
	}
	//成功
	httpSuccessJson(ctx, companys)
}

// log顯示錯誤訊息
// 回傳Http狀態碼400與訊息
func httpErrorJson(ctx *gin.Context, err error, status string, msg string) {
	log.Println(err)
	ctx.JSON(http.StatusBadRequest, gin.H{"status": status, "message": msg})
}

// 回傳Http狀態碼200與訊息
func httpSuccessJson(ctx *gin.Context, msg any) {
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": msg})
}
