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
	initObjs.Server.POST("/buy", itemDAO.BuyItem)
	initObjs.Server.GET("/item", itemDAO.FindItem)
	initObjs.Server.POST("/item", itemDAO.AddItem)
	initObjs.Server.Run(":" + initObjs.Cfg.Port)
}
