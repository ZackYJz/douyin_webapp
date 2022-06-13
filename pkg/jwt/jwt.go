package jwt

import (
	"errors"
	"go_webapp/global"
	"time"

	"go.uber.org/zap"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID int64) (string, error) {
	// 创建一个我们自己的声明的数据
	c := Claims{
		userID, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * global.JWTSetting.Expire).Unix(), // 过期时间
			Issuer:    global.JWTSetting.Issuer,                                    // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象,加密并获得完整的编码后的字符串token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	//需要 []byte 类型，string类型会报错 key is of invalid type
	return token.SignedString([]byte(global.JWTSetting.Secret))
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (claims *Claims, err error) {
	// 解析token
	var token *jwt.Token
	claims = new(Claims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		global.Logger.Error("token parse error!", zap.Error(err))
		return
	}
	if !token.Valid { // 校验token
		err = errors.New("invalid token")
	}
	return
}

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return []byte(global.JWTSetting.Secret), nil
}

// RefreshToken 刷新AccessToken
//func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
//	//解析 refresh token 是否有效，无效直接返回
//	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
//		return
//	}
//
//	// 从旧access token中解析出claims数据
//	var claims Claims
//	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
//	v, _ := err.(*jwt.ValidationError)
//
//	// 当access token是过期错误 并且 refresh token没有过期时就创建一个新的access token
//	if v.Errors == jwt.ValidationErrorExpired {
//		return GenToken(claims.UserID)
//	}
//	return
//}
