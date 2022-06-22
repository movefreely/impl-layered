package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"net/http"
	"strings"
	"time"
)

var (
	Secret = "123#111" //密钥
)

type JWTClaims struct {
	jwt.StandardClaims
	ID    string `json:"id"`
	Email string `json:"email"`
}

func CreateToken(email, ID string) (string, error) {
	expireTime := time.Now().Add(2 * time.Hour) //过期时间
	nowTime := time.Now()                       //当前时间
	claims := JWTClaims{
		Email: email,
		ID:    ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间戳
			IssuedAt:  nowTime.Unix(),    //当前时间戳
			Issuer:    "cutim",           //颁发者签名
			Subject:   "userToken",       //签名主题
		},
	}
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenStruct.SignedString([]byte(Secret))
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*JWTClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return []byte(Secret), nil
		})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 10002,
				"msg":  "token is empty",
			})
			c.Abort()
			return
		}
		token = strings.Replace(token, "Bearer ", "", -1)
		//fmt.Println(token)
		// 按.分割字符串
		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// 解析token
		mc, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 2005,
				"msg":  "token解析失败",
			})
			c.Abort()
			return
		}
		//token超时
		if time.Now().Unix() > mc.ExpiresAt {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 0, "msg": "token过期"})
			c.Abort()
			return
		}
		c.Set("email", mc.Email)
		c.Set("id", mc.ID)
		c.Next()
	}
}
