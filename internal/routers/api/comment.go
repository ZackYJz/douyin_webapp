package api

import (
	"go_webapp/global"
	"go_webapp/internal/dao/redis"
	"go_webapp/internal/service"
	"go_webapp/pkg/app"
	"go_webapp/pkg/errcode"
	"go_webapp/pkg/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Comment struct{}

func NewComment() Comment {
	return Comment{}
}

func (ct Comment) CreateCommentHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.CommentRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("DoIFollowPublisher.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	commentVO, err := svc.CreateComment(param)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(commentVO)
}

func (ct Comment) CommentCountsHandler(c *gin.Context) {
	videoId := c.Query("vlogId")
	response := app.NewResponse(c)
	if videoId == "" {
		response.ResponseError(errcode.InvalidParams)
		return
	}
	count, err := redis.GetVideoCommentCount(videoId)
	if err != nil {
		response.ResponseErrorString(err.Error())
		return
	}
	response.ResponseSuccess(util.StrTo(count).MustInt())
}

func (ct Comment) CommentListHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.CommentListRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("CommentListHandler.BindAndValid", zap.Error(errors))
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
	CommentVOList, totalRows, err := svc.GetCommentList(param)
	if err != nil {
		response.ResponseError(errcode.ServerError)
		return
	}
	response.ResponseList(CommentVOList, totalRows)
}

func (ct Comment) DeleteHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.BaseCommentRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("DeleteHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	err := svc.DeleteComment(param)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(nil)
}

func (ct Comment) LikeCommentHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.BaseCommentRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("DeleteHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	err := svc.LikeComment(param)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(nil)
}

func (ct Comment) UnLikeCommentHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.BaseCommentRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("DeleteHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	err := svc.UnlikeComment(param)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(nil)
}
