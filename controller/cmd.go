package controller

import (
	"github.com/gin-gonic/gin"
	"in-server/service"
	"in-server/util/middleware"
)

func Entrance() {
	service.InitBayes()
	router := gin.Default()
	router.Use(middleware.Cors())
	router.POST("/register", UserRegister)
	router.POST("/login", UserLogin)
	//router.POST("/token", GetToken)
	router.POST("/isad", IsAd)
	router.GET("/ws", Chat)
	verify := router.Group("/verify").Use(middleware.JWTAuthMiddleware())
	{
		verify.POST("/validation", TestToken)
		verify.POST("/refresh", RefreshToken)
		verify.POST("/message-list", GetHistory)
		verify.POST("/user-info", getUserInfo)
		verify.POST("/add-friend", AddFriend)
		verify.POST("/find-user", FindUserByName)
		verify.POST("/agree-friend", AgreeToFriend)
		verify.POST("/refuse-friend", RejectToFriend)
		verify.POST("/history-message", HistoryMessage)
		verify.POST("/change-avatar", ChangeAvatar)
	}
	_ = router.Run("0.0.0.0:8080")
}
