package httputil

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// log顯示錯誤訊息
// 回傳Http狀態碼400與訊息
func HttpErrorJson(ctx *gin.Context, err error, status string, msg string) {
	log.Println(err)
	ctx.JSON(http.StatusBadRequest, gin.H{"status": status, "message": msg})
}

// 回傳Http狀態碼200與訊息
func HttpSuccessJson(ctx *gin.Context, msg any) {
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": msg})
}
