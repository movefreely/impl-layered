package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"in-server/model"
	"in-server/service"
	"strconv"
	"time"
)

var fiendServer = service.FriendServer{}
var userServer = service.UserService{}

func AddFriend(ctx *gin.Context) {
	oneId := ctx.PostForm("selfId")
	otherId := ctx.PostForm("friendId")

	oneid, err := strconv.Atoi(oneId)
	if err != nil {
		return
	}
	otherid, err := strconv.Atoi(otherId)
	if err != nil {
		return
	}

	friend, err := fiendServer.AddFriend(uint64(oneid), uint64(otherid))
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 10004,
			"msg":  err.Error(),
		})
		return
	}
	//SendMsg()
	sendM(oneid, otherid, &model.Message{
		FromId:   uint64(oneid),
		ToId:     uint64(otherid),
		Type:     2,
		Content:  oneId + "请求添加你为好友",
		CreateAt: time.Now(),
	})
	ctx.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": friend,
	})
}

func FindUserByName(ctx *gin.Context) {
	name := ctx.PostForm("name")
	if name == "" {
		ctx.JSON(200, gin.H{
			"code": 10004,
			"msg":  "请传入正确的用户名",
		})
		return
	}
	userByName := fiendServer.FindUserByNickname(name)
	atoi, err := strconv.Atoi(name)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":      10000,
			"msg":       "success",
			"data_name": userByName,
			"data_id":   make([]int, 0),
		})
		return
	}
	userById, err := userServer.UserInfo(uint64(atoi))
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":      10000,
			"msg":       "success",
			"data_name": userByName,
			"data_id":   make([]int, 0),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":      10000,
		"msg":       "success",
		"data_name": userByName,
		"data_id":   userById,
	})
}

func AgreeToFriend(ctx *gin.Context) {
	oneId := ctx.PostForm("selfId")
	otherID := ctx.PostForm("friendId")

	flag, err := fiendServer.AgreeFriend(oneId, otherID)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 10004,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": flag,
	})
	oneid, err := strconv.Atoi(oneId)
	if err != nil {
		return
	}
	otherid, err := strconv.Atoi(otherID)
	if err != nil {
		return
	}
	sendM(oneid, otherid, &model.Message{
		FromId:   uint64(oneid),
		ToId:     uint64(otherid),
		Type:     0,
		Content:  "我们已经是好友了，快来聊天吧",
		CreateAt: time.Now(),
	})
}

func RejectToFriend(ctx *gin.Context) {
	oneId := ctx.PostForm("selfId")
	otherID := ctx.PostForm("friendId")

	flag, err := fiendServer.RejectFriend(oneId, otherID)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 10004,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": flag,
	})
}

func sendM(oneid, otherid int, msg *model.Message) {
	go func() {
		err := MsgDeal.UpdateMessage(msg)
		if err != nil {
			fmt.Println("in dispatch update message error, err = " + err.Error())
		}
	}()
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("转为JSON出错了")
		return
	}
	err = SendMsg(uint64(otherid), data, msg)
	if err != nil {
		fmt.Println("请求好友发送信息出错")
		return
	}
}
