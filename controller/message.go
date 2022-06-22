package controller

import (
	"github.com/gin-gonic/gin"
	"in-server/service"
	"net/http"
	"strconv"
)

var msgService = &service.Message{}

func GetHistory(ctx *gin.Context) {
	userid := ctx.Keys["id"].(string)
	id, err := strconv.Atoi(userid)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 10004,
			"msg":  "id不合法",
		})
		return
	}
	message, err := msgService.NewestMessage(int64(id))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 10005,
			"msg":  "get history message error",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": message,
	})
}

func HistoryMessage(ctx *gin.Context) {
	userid := ctx.Keys["id"].(string)
	oneid, err := strconv.Atoi(userid)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 10004,
			"msg":  "id不合法",
		})
		return
	}
	otherId := ctx.PostForm("friendId")
	otherid, err := strconv.Atoi(otherId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 10004,
			"msg":  "id不合法",
		})
		return
	}
	page := ctx.PostForm("page")
	pageNum, err := strconv.Atoi(page)
	message, err := msgService.HistoryMessage(uint64(oneid), uint64(otherid), pageNum)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 10005,
			"msg":  "get history message error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": message,
	})
}
