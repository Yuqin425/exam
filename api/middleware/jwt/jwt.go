package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

type MyClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 定义Secret
var mySecret = []byte("夏天夏天悄悄过去")

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}

// 定义JWT的过期时间
const TokenExpireDuration = time.Hour * 2

// GenToken 生成access token 和 refresh token
func GenToken(userID uint64, username string) (aToken, rToken string, err error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		userID,   // 自定义字段
		username, // 自定义字段
		jwt.StandardClaims{ // JWT规定的7个官方字段
			ExpiresAt: time.Now().Add(
				time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(), // 过期时间
			Issuer: "onlineJudge", // 签发人
		},
	}
	// 加密并获得完整的编码后的字符串token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)

	// refresh token 不需要存任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * 30).Unix(), // 过期时间
		Issuer:    "onlineJudge",                           // 签发人
	}).SignedString(mySecret)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return
}

func ParseToken(tokenString string) (claims *MyClaims, err error) {
	// 解析token
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return
	}
	if !token.Valid { // 校验token
		err = errors.New("invalid token")
	}
	return
}

// RefreshToken 刷新AccessToken
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// refresh token无效直接返回
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}

	// 从旧access token中解析出claims数据	解析出payload负载信息
	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)

	// 当access token是过期错误 并且 refresh token没有过期时就创建一个新的access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserID, claims.Username)
	}
	return
}