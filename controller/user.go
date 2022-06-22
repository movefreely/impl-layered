package controller

import (
	"github.com/gin-gonic/gin"
	"in-server/service"
	"in-server/util/middleware"
	"strconv"
)

var UserService = service.UserService{}

func UserRegister(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	nickname := ctx.PostForm("nickname")
	// 未获取到传输参数
	if password == "" || email == "" || nickname == "" {
		ctx.JSON(200, gin.H{
			"code":    10002,
			"message": "参数错误",
		})
		return
	}
	// 注册
	register, err := UserService.UserRegister(email, password, nickname)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    10004,
			"message": err.Error(),
		})
	} else {
		token, _ := middleware.CreateToken(email, strconv.Itoa(int(register.ID)))
		ctx.JSON(200, gin.H{
			"code":    10000,
			"message": "success",
			"token":   token,
			"data":    register,
		})
	}
}

func getUserInfo(ctx *gin.Context) {
	userID := ctx.PostForm("id")
	if userID == "" {
		ctx.JSON(200, gin.H{
			"code":    10002,
			"message": "参数错误",
		})
		return
	}
	id, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    10002,
			"message": "请输入正确的id",
		})
		return
	}
	user, err := UserService.UserInfo(uint64(id))
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    10004,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    10000,
			"message": "success",
			"data":    user,
		})
	}
}

// UserLogin 用户登录
func UserLogin(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	// 未获取到传输参数
	if password == "" || email == "" {
		ctx.JSON(200, gin.H{
			"code":    10003,
			"message": "参数错误",
		})
		return
	}
	// 登录
	login, err := UserService.UserLogin(email, password)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    10004,
			"message": err.Error(),
		})
	} else {
		token, _ := middleware.CreateToken(email, strconv.Itoa(int(login.ID)))
		ctx.JSON(200, gin.H{
			"code":    10000,
			"message": "success",
			"data":    login,
			"token":   token,
		})
	}
}

func GetToken(ctx *gin.Context) {
	email := ctx.PostForm("email")
	id := ctx.PostForm("id")
	// 未获取到传输参数
	if id == "" || email == "" {
		ctx.JSON(200, gin.H{
			"code":    10003,
			"message": "参数错误",
		})
		return
	}
	token, err := middleware.CreateToken(email, id)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    10004,
			"message": err.Error(),
		})
	}
	ctx.JSON(200, gin.H{
		"code":    10000,
		"message": "success",
		"data":    token,
	})
}

func TestToken(ctx *gin.Context) {
	email := ctx.Keys["email"].(string)
	id := ctx.Keys["id"].(string)
	ctx.JSON(200, gin.H{
		"code":    10000,
		"message": "success",
		"data":    "email:" + email + "id:" + id,
	})
}

func RefreshToken(ctx *gin.Context) {
	email := ctx.Keys["email"].(string)
	id := ctx.Keys["id"].(string)
	// 未获取到传输参数
	if id == "" || email == "" {
		ctx.JSON(200, gin.H{
			"code":    10003,
			"message": "参数错误",
		})
		return
	}
	token, err := middleware.CreateToken(email, id)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    10004,
			"message": err.Error(),
		})
	}
	ctx.JSON(200, gin.H{
		"code":    10000,
		"message": "success",
		"data":    token,
	})
}

func ChangeAvatar(ctx *gin.Context) {
	id := ctx.Keys["id"].(string)
	url := ctx.PostForm("url")
	if id == "" || url == "" {
		ctx.JSON(200, gin.H{
			"code":    10003,
			"message": "参数错误",
		})
		return
	}
	atoi, err := strconv.Atoi(id)

	err = UserService.ChangeAvatar(uint64(atoi), url)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    10004,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":    10000,
		"message": "success",
	})
}
