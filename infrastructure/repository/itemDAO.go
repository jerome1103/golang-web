package itemDAO

import (
	"context"
	httputil "webserver/infrastructure/commons/utils"
	inits "webserver/infrastructure/init"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	// 限制查找後返回資料的數量
	findOptions := options.Find()
	findOptions.SetLimit(50)
	// 查詢條件
	filter := bson.D{{}}
	// 資料查找
	cur, err := tableitem.Find(context.TODO(), filter, findOptions)
	if err != nil {
		httputil.HttpErrorJson(c, err, "find_db_fail", "find db fail.")
		return
	}
	// 包裝成Item
	var items []Item
	if err = cur.All(context.TODO(), &items); err != nil {
		httputil.HttpErrorJson(c, err, "to_data_fail", "the data to item error.")
		return
	}
	// 成功
	httputil.HttpSuccessJson(c, items)
}

func FindItemByID(ctx *gin.Context) {
	//擷取到值轉換成int
	id := ctx.Param("id")
	//查詢條件
	filter := bson.D{primitive.E{Key: "iid", Value: id}}
	//單筆查詢
	var item Item
	if err := tableitem.FindOne(context.TODO(), filter).Decode(&item); err != nil {
		httputil.HttpErrorJson(ctx, err, "no_data", "Data not found.")
		return
	}

	//成功
	httputil.HttpSuccessJson(ctx, item)
}

func AddItem(ctx *gin.Context) {
	// 擷取資料
	newItem := paramItemData(ctx)
	// 插入一條數據
	res, err := tableitem.InsertOne(context.TODO(), newItem)
	if err != nil {
		httputil.HttpErrorJson(ctx, err, "add_fail", "add failed.")
		return
	}
	httputil.HttpSuccessJson(ctx, res.InsertedID)
}

func UpdateItem(ctx *gin.Context) {
	// 擷取資料
	newItem := paramItemData(ctx)
	// 擷取id
	newItem.Iid = ctx.Param("iid")
	// 過濾條件
	filter := bson.D{primitive.E{Key: "iid", Value: newItem.Iid}}
	// 包裝要更新的資料
	update := bson.D{primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{Key: "name", Value: newItem.Name},
			primitive.E{Key: "price", Value: newItem.Price},
			primitive.E{Key: "text", Value: newItem.Text},
		},
	}}
	// 設置upsert模式
	// 當開啟此行,若更新時沒有這筆資料,就進行新增
	// FIXME:未來預計把新增跟修改都呼叫此func
	// updateOpts := options.Update().SetUpsert(true)
	// 更新資料
	updateResult, err := tableitem.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		httputil.HttpErrorJson(ctx, err, "update_error", "update fail.")
		return
	}
	httputil.HttpSuccessJson(ctx, updateResult)
}

func DeleteItem(ctx *gin.Context) {
	// 擷取到值轉換成int
	id := ctx.Param("id")
	// 刪除一條數據
	filter := bson.D{primitive.E{Key: "iid", Value: id}}
	deleteResult, err := tableitem.DeleteOne(context.TODO(), filter)
	if err != nil {
		httputil.HttpErrorJson(ctx, err, "fail_del", "delete fail.")
	}
	httputil.HttpSuccessJson(ctx, deleteResult)
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
