package service

import (
	"go_webapp/global"
	"go_webapp/internal/dao/redis"
	"go_webapp/internal/model"
	"go_webapp/pkg/errcode"
	"go_webapp/pkg/snowflake"
	"go_webapp/pkg/util"
	"strconv"

	"go.uber.org/zap"

	"gorm.io/gorm"
)

type VideoRequest struct {
	Id             string `json:"id"`
	PublisherId    string `json:"vlogerId" binding:"required"`
	Url            string `json:"url" binding:"required"`
	Cover          string `json:"cover" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	LikeCounts     int    `json:"likeCounts"`
	CommentsCounts int    `json:"commentsCounts"`
}

type IndexListRequest struct {
	UserId   string `form:"userId"`
	Search   string `form:"search"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

type UserVideoRequest struct {
	UserId      string `form:"userId"`
	VideoId     string `form:"vlogId"`
	PublisherId string `form:"vlogerId"`
}

type MyRelationVideoListRequest struct {
	UserId   string `form:"myId"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

func (s *Service) PublishVideo(param *VideoRequest) (error *errcode.Error) {
	id := snowflake.GenID()
	video := &model.Video{
		Model:          &gorm.Model{ID: uint(id)},
		Publisher:      util.StrTo(param.PublisherId).MustInt64(),
		Url:            param.Url,
		Cover:          param.Cover,
		Title:          param.Title,
		Width:          param.Width,
		Height:         param.Height,
		LikeCounts:     param.LikeCounts,
		CommentsCounts: param.CommentsCounts,
		IsPrivate:      global.NO,
	}
	err := s.dao.Insert(video)
	if err != nil {
		return errcode.VideoPublishFailed
	}
	return nil
}

func (s *Service) IndexList(param *IndexListRequest) ([]*model.VideoVO, int, *errcode.Error) {
	totalRows, err := s.dao.CountVideoIndexList()
	if err != nil {
		global.Logger.Error("CountVideoIndexList Error", zap.Error(err))
		return nil, 0, errcode.ServerError
	}
	videos, err := s.dao.GetVideoIndexList(param.Search, param.Page, param.PageSize)
	if err != nil {
		global.Logger.Error("GetVideoIndexList Error", zap.Error(err))
		return nil, 0, errcode.ServerError
	}

	var videoVOList []*model.VideoVO
	for _, v := range videos {
		//查询关注和喜欢
		publisherId, videoId := v.Publisher, v.ID
		publisher, err := s.dao.GetUserByUserId(strconv.FormatInt(publisherId, 10))
		if err != nil {
			global.Logger.Error("GetUserByUserId failed", zap.Error(err))
		}
		//额外需要查询的参数
		likeCounts, _ := s.VideoTotalLikedCounts(strconv.Itoa(int(videoId)))
		isFollow, _ := s.DoIFollowPublisher(param.UserId, strconv.FormatInt(publisherId, 10))
		isLike := redis.DoILikeVideo(param.UserId, strconv.Itoa(int(videoId)))
		commentCount, err := redis.GetVideoCommentCount(strconv.Itoa(int(videoId)))
		videoVOList = append(videoVOList, &model.VideoVO{
			ID:                 strconv.Itoa(int(v.ID)),
			Publisher:          strconv.FormatInt(v.Publisher, 10),
			PublisherFace:      publisher.Photo,
			PublisherName:      publisher.Nickname,
			Title:              v.Title,
			Url:                v.Url,
			Cover:              v.Cover,
			Width:              v.Width,
			Height:             v.Height,
			LikeCounts:         likeCounts,
			CommentsCounts:     util.StrTo(commentCount).MustInt(),
			IsPrivate:          v.IsPrivate,
			IsPlay:             false,
			DoIFollowPublisher: isFollow,
			DoILikeThisVideo:   isLike,
		})
	}
	return videoVOList, totalRows, nil
}

func (s *Service) GetVideoDetailById(userId, videoId string) (*model.VideoVO, *errcode.Error) {
	video, err := s.dao.GetVlogById(videoId)
	if err != nil {
		global.Logger.Error("GetVideoDetailById.GetVlogById Error", zap.Error(err))
		return nil, errcode.ServerError
	}
	publisherId := video.Publisher
	publisher, err := s.dao.GetUserByUserId(strconv.FormatInt(publisherId, 10))
	if err != nil {
		global.Logger.Error("GetVideoDetailById.GetUserByUserId Error", zap.Error(err))
		return nil, errcode.ServerError
	}
	//likeCounts, _ := s.VideoTotalLikedCounts(videoId)
	//isFollow, _ := s.DoIFollowPublisher(userId, strconv.FormatInt(publisherId, 10))
	isLike := redis.DoILikeVideo(userId, videoId)
	//commentCount, err := redis.GetVideoCommentCount(videoId))
	videoVO := &model.VideoVO{
		ID:                 strconv.Itoa(int(video.ID)),
		Publisher:          strconv.FormatInt(video.Publisher, 10),
		PublisherFace:      publisher.Photo,
		PublisherName:      publisher.Nickname,
		Title:              video.Title,
		Url:                video.Url,
		Cover:              video.Cover,
		Width:              video.Width,
		Height:             video.Height,
		LikeCounts:         video.LikeCounts,
		CommentsCounts:     video.CommentsCounts,
		IsPrivate:          video.IsPrivate,
		IsPlay:             false,
		DoIFollowPublisher: false,
		DoILikeThisVideo:   isLike,
	}
	return videoVO, nil
}

func (s *Service) ChangeVideoStatus(userId, videoId string, status int) *errcode.Error {
	err := s.dao.ChangeVideoStatus(userId, videoId, status)
	if err != nil {
		global.Logger.Error("ChangeVideoStatus.ChangeVideoStatus Error", zap.Error(err))
		return errcode.ServerError
	}
	return nil
}

func (s *Service) QueryMyVideoList(param *IndexListRequest, isPrivate int) ([]*model.VideoSimpleVO, int, *errcode.Error) {
	totalRows, err := s.dao.CountMyVideoList(param.UserId, isPrivate)
	if err != nil {
		global.Logger.Error("QueryMyVideoList.CountMyVideoList Error", zap.Error(err))
		return nil, 0, errcode.ServerError
	}
	videoList, err := s.dao.QueryMyVideoList(param.UserId, isPrivate, param.Page, param.PageSize)
	if err != nil {
		global.Logger.Error("QueryMyVideoList.QueryMyVideoList Error", zap.Error(err))
		return nil, 0, errcode.ServerError
	}
	return videoList, totalRows, nil
}

func (s *Service) LikeVideo(userId, videoId, publisher string) *errcode.Error {
	if redis.DoILikeVideo(userId, videoId) {
		return errcode.DuplicateLike
	}

	like := &model.Like{
		Model:   &gorm.Model{ID: uint(snowflake.GenID())},
		UserId:  util.StrTo(userId).MustInt64(),
		VideoId: util.StrTo(videoId).MustInt64(),
	}

	row, err := s.dao.InsertLike(like)
	if row == 0 || err != nil {
		global.Logger.Error("LikeVideo.InsertLike", zap.Error(err))
		return errcode.ServerError
	}

	err = redis.LikeVideo(userId, videoId, publisher)
	if err != nil {
		global.Logger.Error("LikeVideo.LikeVideo", zap.Error(err))
		return errcode.ServerError
	}

	//TODO 发送点赞消息
	return nil
}

func (s *Service) UnLikeVideo(userId, videoId, publisher string) *errcode.Error {
	if !redis.DoILikeVideo(userId, videoId) {
		return errcode.UnLike
	}

	like := &model.Like{
		UserId:  util.StrTo(userId).MustInt64(),
		VideoId: util.StrTo(videoId).MustInt64(),
	}

	_, err := s.dao.DeleteLike(like)
	if err != nil {
		global.Logger.Error("UnLikeVideo.DeleteLike", zap.Error(err))
		return errcode.ServerError
	}
	err = redis.UnLikeVideo(userId, videoId, publisher)
	if err != nil {
		global.Logger.Error("LikeVideo.LikeVideo", zap.Error(err))
		return errcode.ServerError
	}
	return nil
}

func (s *Service) MyLikedList(param *IndexListRequest) ([]*model.VideoVO, int, *errcode.Error) {
	totalRows, err := s.dao.CountMyLikedVideo(param.UserId)
	if err != nil {
		global.Logger.Error("MyLikedList.CountMyLikedVideo Error", zap.Error(err))
		return nil, 0, errcode.ServerError
	}
	videoList, err := s.dao.QueryMyLikedVideoList(param.UserId, param.Page, param.PageSize)
	if err != nil {
		global.Logger.Error("QueryMyVideoList.QueryMyVideoList Error", zap.Error(err))
		return nil, 0, errcode.ServerError
	}
	return videoList, totalRows, nil
}

func (s *Service) VideoTotalLikedCounts(videoId string) (int, *errcode.Error) {
	counts, err := redis.GetVideoLikedCounts(videoId)
	if err != nil {
		global.Logger.Error("VideoTotalLikedCounts.GetVideoLikedCounts", zap.Error(err))
		return 0, nil
	}
	return counts, nil
}

func (s *Service) GetMyFollowVlogList(param *MyRelationVideoListRequest) ([]*model.VideoVO, int, *errcode.Error) {
	totalRows, err := s.dao.CountMyFollowVideo(param.UserId)
	if err != nil {
		global.Logger.Error("MyLikedList.CountMyFollowVideo Error", zap.Error(err))
		return nil, 0, errcode.ServerError
	}
	videoList, err := s.dao.QueryMyFollowVideoList(param.UserId, param.Page, param.PageSize)
	if err != nil {
		global.Logger.Error("QueryMyVideoList.QueryMyFollowVideoList Error", zap.Error(err))
		return nil, 0, errcode.ServerError
	}
	for _, video := range videoList {
		videoId := video.ID
		video.DoIFollowPublisher = true
		video.DoILikeThisVideo = redis.DoILikeVideo(param.UserId, videoId)
		video.LikeCounts, _ = s.VideoTotalLikedCounts(videoId)
	}
	return videoList, totalRows, nil
}

func (s *Service) GetMyFriendVlogList(param *MyRelationVideoListRequest) ([]*model.VideoVO, int, *errcode.Error) {
	totalRows, err := s.dao.CountMyFriendVideo(param.UserId)
	if err != nil {
		global.Logger.Error("MyLikedList.CountMyLikedVideo Error", zap.Error(err))
		return nil, 0, errcode.ServerError
	}
	videoList, err := s.dao.QueryMyFriendVideoList(param.UserId, param.Page, param.PageSize)
	if err != nil {
		global.Logger.Error("QueryMyVideoList.QueryMyVideoList Error", zap.Error(err))
		return nil, 0, errcode.ServerError
	}
	for _, video := range videoList {
		videoId := video.ID
		video.DoIFollowPublisher = true
		video.DoILikeThisVideo = redis.DoILikeVideo(param.UserId, videoId)
		video.LikeCounts, _ = s.VideoTotalLikedCounts(videoId)
	}
	return videoList, totalRows, nil
}
