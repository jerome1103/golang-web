package itemDAO

import (
	"context"
	httputil "webserver/infrastructure/commons/utils"
	inits "webserver/infrastructure/init"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Item struct {
	Iid   string `bson:"iid"`
	Name  string `bson:"name"`
	Price string `bson:"price"`
	Text  string `bson:"text"`
}

var tableitem *mongo.Collection

func init() {
	tableitem = inits.GetinitObj().Mongoclient.Database("shopweb").Collection("item")
}

func FindItem(c *gin.Context) {
	//限制查找後返回資料的數量
	findOptions := options.Find()
	findOptions.SetLimit(50)
	//查詢條件
	filter := bson.D{{}}
	//資料查找
	cur, err := tableitem.Find(context.TODO(), filter, findOptions)
	if err != nil {
		httputil.HttpErrorJson(c, err, "find_db_fail", "find db fail.")
		return
	}
	//包裝成Item
	var items []Item
	if err = cur.All(context.TODO(), &items); err != nil {
		httputil.HttpErrorJson(c, err, "to_data_fail", "the data to item error.")
		return
	}
	//成功
	httputil.HttpSuccessJson(c, items)
}

func AddItem(ctx *gin.Context) {
	//擷取資料
	newItem := paramItemData(ctx)
	//擷取id
	newItem.Iid = ctx.PostForm("iid")
	//插入一條數據
	res, err := tableitem.InsertOne(context.TODO(), newItem)
	if err != nil {
		httputil.HttpErrorJson(ctx, err, "add_fail", "add failed.")
		return
	}
	httputil.HttpSuccessJson(ctx, res.InsertedID)
}

func BuyItem(c *gin.Context) {

}

func paramItemData(ctx *gin.Context) Item {
	id := ctx.PostForm("iid")
	name := ctx.PostForm("name")
	price := ctx.PostForm("price")
	text := ctx.DefaultPostForm("text", "無")
	return Item{id, name, price, text}
}
