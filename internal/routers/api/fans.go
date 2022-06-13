package api

import (
	"go_webapp/global"
	"go_webapp/internal/service"
	"go_webapp/pkg/app"
	"go_webapp/pkg/errcode"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Fans struct{}

func NewFans() Fans {
	return Fans{}
}

func (f Fans) FollowHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.FansRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("FollowHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	if param.PublisherId == "" || param.PublisherId == param.MyId {
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	err := svc.Follow(param.MyId, param.PublisherId)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(nil)
}

func (f Fans) CancelFollowHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.FansRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("FollowHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	err := svc.UnFollow(param.MyId, param.PublisherId)
	if err != nil {
		response.ResponseError(err)
	}
	response.ResponseSuccess(nil)
}

func (f Fans) DoIFollowPublisherHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.FansRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("DoIFollowPublisher.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	isFollow, err := svc.DoIFollowPublisher(param.MyId, param.PublisherId)
	if err != nil {
		response.ResponseError(err)
	}
	response.ResponseSuccess(isFollow)
}

func (f Fans) MyFollowListHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.FansRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("DoIFollowPublisher.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	if param.Page == 0 {
		param.Page = global.DEFAULT_START_PAGE
	}
	if param.PageSize == 0 {
		param.PageSize = global.DEFAULT_PAGE_SIZE
	}
	FollowerList, totalRows, err := svc.GetFollowerList(param)
	if err != nil {
		response.ResponseError(errcode.ServerError)
		return
	}
	response.ResponseList(FollowerList, totalRows)
}

func (f Fans) MyFansListHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.FansRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("DoIFollowPublisher.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	if param.Page == 0 {
		param.Page = global.DEFAULT_START_PAGE
	}
	if param.PageSize == 0 {
		param.PageSize = global.DEFAULT_PAGE_SIZE
	}
	FansList, totalRows, err := svc.GetFansList(param)
	if err != nil {
		response.ResponseError(errcode.ServerError)
		return
	}
	response.ResponseList(FansList, totalRows)
}
