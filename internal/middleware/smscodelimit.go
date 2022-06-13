package middleware

import (
	"go_webapp/global"
	"go_webapp/internal/dao/redis"
	"go_webapp/pkg/app"
	"go_webapp/pkg/errcode"
	"go_webapp/pkg/util"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SMSCodeLimitMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		ip := util.GetRealIp(c)
		global.Logger.Info("拦截短信发送", zap.String("ip", ip))
		has := redis.IsUserSMSHas(ip)
		if has {
			response := app.NewResponse(c)
			response.ResponseError(errcode.GetSMSCodeTooFast)
			c.Abort()
			return
		}

		c.Next()
	}
}
