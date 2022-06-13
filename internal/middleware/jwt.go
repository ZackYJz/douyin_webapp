package middleware

import (
	"go_webapp/pkg/app"
	"go_webapp/pkg/errcode"
	j "go_webapp/pkg/jwt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey = "userID"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			token   string
			errCode = errcode.Success
		)
		if s, exist := c.GetQuery("headerUserToken"); exist {
			token = s
		} else {
			token = c.GetHeader("headerUserToken")
			//global.Logger.Debug("接收到token来自header", zap.String("token", token))
		}
		if token == "" {
			errCode = errcode.NoToken
		} else {
			claims, err := j.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					errCode = errcode.UnauthorizedTokenTimeout
				default:
					errCode = errcode.UnauthorizedTokenError
				}
			}
			c.Set(ContextUserIDKey, claims.UserID)
		}

		if errCode != errcode.Success {
			response := app.NewResponse(c)
			response.ResponseError(errCode)
			c.Abort()
			return
		}

		c.Next()
	}
}
