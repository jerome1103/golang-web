package userdao

import (
	"context"
	"fmt"
	"webserver/infrastructure/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ctx context.Context
	cfg *config.Config
)

func Connect() {
	mongoconn := options.Client().ApplyURI(cfg.DBUri)
	//連線到Mongodb
	var err error
	mongoclient, err := mongo.Connect(ctx, mongoconn)
	if err != nil {
		panic(err)
	}
	//檢查連線
	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	//連線成功
	fmt.Println("MongoDB successfully connected...")
}
