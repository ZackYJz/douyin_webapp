package service

import (
	"go_webapp/global"
	"go_webapp/internal/dao/redis"
	"go_webapp/internal/model"
	"go_webapp/pkg/errcode"
	"go_webapp/pkg/snowflake"
	"go_webapp/pkg/util"

	"gorm.io/gorm"
)

type CommentRequest struct {
	Publisher string `json:"vlogerId" binding:"required"`
	FatherId  string `json:"fatherCommentId" `
	VideoId   string `json:"vlogId" binding:"required"`
	UserId    string `json:"commentUserId" binding:"required"`
	Content   string `json:"content" binding:"required,max=50"`
}

type CommentListRequest struct {
	VideoId  string `form:"vlogId" binding:"required"`
	UserId   string `form:"userId"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

type BaseCommentRequest struct {
	CommentUserId string `form:"commentUserId"`
	UserId        string `form:"userId"`
	Id            string `form:"commentId"`
	VideoId       string `form:"vlogId"`
}

func (s *Service) CreateComment(param *CommentRequest) (*model.CommentVO, *errcode.Error) {
	commentId := snowflake.GenID()
	newComment := &model.Comment{
		Model:      &gorm.Model{ID: uint(commentId)},
		Publisher:  util.StrTo(param.Publisher).MustInt64(),
		VideoId:    util.StrTo(param.VideoId).MustInt64(),
		FatherId:   util.StrTo(param.FatherId).MustInt64(),
		UserId:     util.StrTo(param.UserId).MustInt64(),
		Content:    param.Content,
		LikeCounts: 0,
	}

	err := s.dao.InsertComment(newComment)
	err = redis.AddVideoCommentCount(param.VideoId)
	if err != nil {
		return nil, errcode.ServerError
	}
	commentVO := newComment.GetVO()
	return commentVO, nil
}

func (s *Service) GetCommentList(param *CommentListRequest) ([]*model.CommentVO, int, *errcode.Error) {
	//先查总数
	rowCount, err := s.dao.CountComment(param.VideoId)
	if err != nil {
		return nil, 0, errcode.ServerError
	}
	commentList, err := s.dao.GetCommentVOList(param.VideoId, param.Page, param.PageSize)
	if err != nil {
		return nil, 0, errcode.ServerError
	}
	for _, v := range commentList {
		//评论点赞总数
		commentId := v.CommentId
		likeCounts, _ := redis.GetCommentLikedCount(commentId)
		v.LikeCounts = util.StrTo(likeCounts).MustInt()
		//当前是否点赞
		doILike, _ := redis.GetDoILikeComment(param.UserId, commentId)
		if "1" == doILike {
			v.IsLike = global.YES
		}
	}
	return commentList, int(rowCount), nil
}

func (s *Service) DeleteComment(param *BaseCommentRequest) *errcode.Error {
	err := s.dao.DeleteComment(param.Id)
	if err != nil {
		return errcode.ServerError
	}
	//减少评论数
	err = redis.DecrCommentNum(param.VideoId)
	if err != nil {
		return errcode.ServerError
	}
	return nil
}

func (s *Service) LikeComment(param *BaseCommentRequest) *errcode.Error {
	err := redis.LikeComment(param.UserId, param.Id)
	if err != nil {
		return errcode.ServerError
	}
	//TODO 点赞的系统消息
	return nil
}

func (s *Service) UnlikeComment(param *BaseCommentRequest) *errcode.Error {
	err := redis.UnlikeComment(param.UserId, param.Id)
	if err != nil {
		return errcode.ServerError
	}
	return nil
}
