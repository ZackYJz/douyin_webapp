package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取请求方法
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		//允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token,headerUserId,headerUserToken")
		//允许的请求方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			//数据体无内容，返回 204 状态
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
