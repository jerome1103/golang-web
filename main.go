package main

import (
	inits "webserver/infrastructure/init"
	itemDAO "webserver/infrastructure/repository"
)

/*
這裡儲存了以下內容
1. Server 由Gin產生的Server
2. Cfg 自定義Config檔案內容擷取出來的資料
3. Mongoclient Mongodb連線
4. Redisclient Redis連線
*/
var initObjs *inits.InitObj

func init() {
	//初始化initObjs參數
	initObjs = inits.GetinitObj()
}

func main() {
	router := initObjs.Server.Group(initObjs.Cfg.Group)
	router.POST("buy", itemDAO.BuyItem)
	router.GET("item", itemDAO.FindItem)
	router.GET("item/:id", itemDAO.FindItemByID)
	router.POST("item", itemDAO.AddItem)
	router.PUT("item/:id", itemDAO.UpdateItem)
	router.DELETE("item/:id", itemDAO.DeleteItem)

	initObjs.Server.Run(":" + initObjs.Cfg.Port)
}
