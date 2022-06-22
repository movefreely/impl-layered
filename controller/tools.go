package controller

import (
	"github.com/gin-gonic/gin"
	"in-server/service"
)

var tool = service.Tools{}

func IsAd(ctx *gin.Context) {
	message := ctx.PostForm("message")
	spam := tool.JudgeIsSpam(message)
	ctx.JSON(200, gin.H{
		"isSpam": spam,
	})
}
