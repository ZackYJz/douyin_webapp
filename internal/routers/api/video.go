package api

import (
	"go_webapp/global"
	"go_webapp/internal/service"
	"go_webapp/pkg/app"
	"go_webapp/pkg/errcode"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Video struct{}

func NewVideo() Video {
	return Video{}
}

func (v Video) PublishHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.VideoRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("PublishHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}

	svc := service.New(c)
	err := svc.PublishVideo(param)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(nil)
}

func (v Video) IndexListHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.IndexListRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("IndexListHandler.BindAndValid", zap.Error(errors))
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
	videoVOList, totalRows, err := svc.IndexList(param)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseList(videoVOList, totalRows)
}

func (v Video) VideoDetailHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.UserVideoRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("VideoDetailHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	videoVO, err := svc.GetVideoDetailById(param.UserId, param.VideoId)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(videoVO)
}

func (v Video) ToPrivateHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.UserVideoRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("ToPrivateHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	err := svc.ChangeVideoStatus(param.UserId, param.VideoId, global.YES)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(nil)
}

func (v Video) PrivateListHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.IndexListRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("PrivateListHandler.BindAndValid", zap.Error(errors))
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
	videoList, totalRows, err := svc.QueryMyVideoList(param, global.YES)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseList(videoList, totalRows)
}

func (v Video) ToPublicHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.UserVideoRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("ToPrivateHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	err := svc.ChangeVideoStatus(param.UserId, param.VideoId, global.NO)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(nil)
}

func (v Video) PublicListHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.IndexListRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("PrivateListHandler.BindAndValid", zap.Error(errors))
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
	videoList, totalRows, err := svc.QueryMyVideoList(param, global.NO)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseList(videoList, totalRows)
}

func (v Video) LikeHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.UserVideoRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("LikeHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	err := svc.LikeVideo(param.UserId, param.VideoId, param.PublisherId)
	if err != nil {
		response.ResponseError(err)
	}
	response.ResponseSuccess(nil)
}

func (v Video) UnLikeHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.UserVideoRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("UnLikeHandler.BindAndValid", zap.Error(errors))
		response.ResponseError(errcode.InvalidParams)
		return
	}
	svc := service.New(c)
	err := svc.UnLikeVideo(param.UserId, param.VideoId, param.PublisherId)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(nil)
}

func (v Video) LikedListHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.IndexListRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("LikedListHandler.BindAndValid", zap.Error(errors))
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
	videoVOList, totalRows, err := svc.MyLikedList(param)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseList(videoVOList, totalRows)
}

// TotalLikedCountsHandler 视频点赞总数
func (v Video) TotalLikedCountsHandler(c *gin.Context) {
	videoId := c.Query("vlogId")
	response := app.NewResponse(c)
	if videoId == "" {
		response.ResponseError(errcode.InvalidParams)
		return
	}

	svc := service.New(c)
	counts, err := svc.VideoTotalLikedCounts(videoId)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseSuccess(counts)
}

func (v Video) FollowVideoListHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.MyRelationVideoListRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("LikedListHandler.BindAndValid", zap.Error(errors))
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
	videoVOList, totalRows, err := svc.GetMyFollowVlogList(param)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseList(videoVOList, totalRows)
}

func (v Video) FriendVideoListHandler(c *gin.Context) {
	response := app.NewResponse(c)
	param := &service.MyRelationVideoListRequest{}
	valid, errors := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Error("LikedListHandler.BindAndValid", zap.Error(errors))
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
	videoVOList, totalRows, err := svc.GetMyFriendVlogList(param)
	if err != nil {
		response.ResponseError(err)
		return
	}
	response.ResponseList(videoVOList, totalRows)
}
