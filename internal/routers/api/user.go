package api

import (
	"go_webapp/global"
	"go_webapp/internal/dao/redis"
	"go_webapp/internal/middleware"
	"go_webapp/internal/service"
	"go_webapp/pkg/app"
	"go_webapp/pkg/errcode"
	"go_webapp/pkg/jwt"
	"go_webapp/pkg/snowflake"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type User struct{}

func NewUser() User {
	return User{}
}

// @Tags 测试接口
// @Summary 测试
// @Produce  json
// @Router /douyin/demo/test [get]
func (u User) Test(c *gin.Context) {
	response := app.NewResponse(c)
	response.ResponseSuccess("注册接口访问 ok")
}

func (u User) GetToken(c *gin.Context) {
	response := app.NewResponse(c)
	token, err := jwt.GenToken(snowflake.GenID())
	if err != nil {
		global.Logger.Error("生成token失败", zap.Error(err))
		response.ResponseError(errcode.ServerError)
		return
	}
	response.ResponseSuccess(token)
}

// @Tags 注册接口
// @Summary 用户账号密码注册接口
// @Produce  json
// @Param BaseUserRequest body service.BaseUserRequest true "注册的用户名和密码"
// @Success 200 {object} model.User "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /douyin/user/register [post]
func (u User) RegisterHandler(c *gin.Context) {
	param := service.BaseUserRequest{}
	response := app.NewResponse(c)
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("RegisterHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	//业务逻辑处理

	//返回结果
	response.ResponseSuccess("测试注册 ok")
}

func (u User) LoginHandler(c *gin.Context) {
	param := service.BaseUserRequest{}
	response := app.NewResponse(c)
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("LoginHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.LoginFailed)
		return
	}
	if err := service.Login(&param); err != nil {
		global.Logger.Error("login failed", zap.String("username", param.Username), zap.Error(err))
		response.ResponseError(errcode.LoginFailed)
		return
	}
	response.ResponseSuccess("测试登录 ok")
}

func (u User) GetSMSCodeHandler(c *gin.Context) {
	param := service.MobileRequest{}
	response := app.NewResponse(c)
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("GetSMSCodeHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	err := svc.SendSMS(&param)
	if err != nil {
		global.Logger.Error("SendSMS Failed", zap.Error(err))
		response.ResponseError(errcode.GetSMSCodeError)
		return
	}

	response.ResponseSuccess(nil)
}

func (u User) MobileLoginHandler(c *gin.Context) {
	param := service.MobileCodeRequest{}
	response := app.NewResponse(c)
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("MobileLoginHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}

	svc := service.New(c)
	user, err := svc.LoginWithMobileCode(param)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(user)
}

func (u User) LogoutHandler(c *gin.Context) {
	response := app.NewResponse(c)
	err := redis.RemoveRefreshToken(c.GetInt64(middleware.ContextUserIDKey))
	if err != nil {
		return
	}
	response.ResponseSuccess(nil)
}

func (u User) QueryInfoHandler(c *gin.Context) {
	param := service.UserIDRequest{}
	response := app.NewResponse(c)
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("QueryInfoHandler.BindAndValid", zap.Error(errors))
		response.ResponseErrorString(errors.Error())
		return
	}

	svc := service.New(c)
	userVO, err := svc.QueryInfo(param.UserId)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(userVO)
}

func (u User) ModifyUserInfoHandler(c *gin.Context) {
	param := service.ModifyUserInfoRequest{}
	response := app.NewResponse(c)
	valid, errors := app.BindAndValid(c, &param)
	typeNum, _ := strconv.Atoi(c.Query("type"))
	if !valid {
		global.Logger.Error("ModifyUserInfoHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}

	svc := service.New(c)
	user, err := svc.ModifyUserInfo(&param, typeNum)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(user)
}

func (u User) ModifyImageHandler(c *gin.Context) {
	param := service.ModifyImageRequest{}
	response := app.NewResponse(c)
	//valid, errors := app.BindAndValid(c, &param)
	//if !valid {
	//	global.Logger.Error("ModifyImageHandler.BindAndValid", zap.Error(errors))
	//	response.ResponseErrorString(errors.Error())
	//	return
	//}
	param.UserId = c.Query("userId")
	param.Type, _ = strconv.Atoi(c.Query("type"))

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		global.Logger.Error("文件接收有误", zap.Error(err))
		response.ResponseErrorString("上传文件出错")
		return
	}
	defer file.Close()

	svc := service.New(c)
	user, e := svc.ModifyImage(&param, file, fileHeader)
	if e != nil {
		response.ResponseError(e)
		return
	}
	response.ResponseSuccess(user)
}
