package service

import (
	"errors"
	"go_webapp/global"
	"go_webapp/internal/dao/redis"
	"go_webapp/internal/model"
	"go_webapp/pkg/errcode"
	"go_webapp/pkg/snowflake"

	"gorm.io/gorm"

	"go.uber.org/zap"
)

type FansRequest struct {
	MyId        string `form:"myId" binding:"required"`
	PublisherId string `form:"vlogerId"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pageSize"`
}

func (s *Service) Follow(myId, publisherId string) *errcode.Error {
	isFollow := redis.GetFansRelations(myId, publisherId)
	if isFollow == "1" {
		return errcode.DuplicateFollow
	}

	Me, err := s.dao.GetUserByUserId(myId)
	Him, _ := s.dao.GetUserByUserId(publisherId)
	if Me == nil || Him == nil {
		global.Logger.Error("Follow.dao.GetUserByUserId", zap.Error(err))
		return errcode.ServerError
	}
	newFans := &model.Fans{
		Model:       &gorm.Model{ID: uint(snowflake.GenID())},
		PublisherId: publisherId,
		FanId:       myId,
		IsFriend:    0,
	}
	fans, err := s.queryIsFans(publisherId, myId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newFans.IsFriend = global.NO
	} else if fans != nil {
		newFans.IsFriend = global.YES
		fans.IsFriend = global.YES
		err := s.dao.UpdateIsFriend(fans)
		if err != nil {
			global.Logger.Error("Follow.dao.UpdateIsFriend", zap.Error(err))
			return errcode.ServerError
		}
	} else {
		global.Logger.Error("queryIsFans", zap.Error(err))
		return errcode.ServerError
	}

	err = s.dao.DoFollow(newFans)
	//TODO 系统消息：关注
	if err != nil {
		global.Logger.Error("Follow.dao.DoFollow", zap.Error(err))
		return errcode.ServerError
	}
	err = redis.DoFollow(myId, publisherId)
	if err != nil {
		global.Logger.Error("Follow.redis.DoFollow", zap.Error(err))
		return errcode.ServerError
	}
	return nil
}

func (s *Service) UnFollow(myId, publisherId string) *errcode.Error {
	isFollow := redis.GetFansRelations(myId, publisherId)
	if isFollow != "1" {
		return errcode.UnFollow
	}
	fans, err := s.queryIsFans(myId, publisherId)
	if err != nil {
		global.Logger.Error("UnFollow.queryIsFans", zap.Error(err))
		return errcode.ServerError
	}
	//判断双方是否是朋友关系，如果是则需要取消双方的关系
	if fans != nil && fans.IsFriend == global.YES {
		fan, err := s.queryIsFans(publisherId, myId)
		if err != nil {
			global.Logger.Error("UnFollow.queryIsFans", zap.Error(err))
			return errcode.ServerError
		}
		// 抹除双方的朋友关系
		fan.IsFriend = global.NO
		err = s.dao.UpdateIsFriend(fan)
		if err != nil {
			global.Logger.Error("UnFollow.queryIsFans", zap.Error(err))
			return errcode.ServerError
		}
	}
	//删除自己的关注关系
	err = s.dao.DeleteFans(fans)
	if err != nil {
		global.Logger.Error("UnFollow.queryIsFans", zap.Error(err))
		return errcode.ServerError
	}
	err = redis.UnFollow(myId, publisherId)

	return nil
}

func (s *Service) DoIFollowPublisher(myId, publisherId string) (bool, *errcode.Error) {
	fans, err := s.queryIsFans(myId, publisherId)
	if fans == nil || err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if fans != nil || err == nil {
		return true, nil
	}
	return false, errcode.ServerError

}

func (s *Service) queryIsFans(myId, publisherId string) (*model.Fans, error) {
	return s.dao.QueryByTwoIds(myId, publisherId)
}

func (s *Service) GetFollowerList(param *FansRequest) ([]*model.FollowerVO, int, *errcode.Error) {
	rowCount, err := s.dao.CountMyFollower(param.MyId)
	if err != nil {
		return nil, 0, errcode.ServerError
	}
	followerList, err := s.dao.QueryMyFollowerList(param.MyId, param.Page, param.PageSize)
	if err != nil {
		return nil, 0, errcode.ServerError
	}
	for _, v := range followerList {
		v.IsFollowed = true
	}
	return followerList, int(rowCount), nil
}

func (s *Service) GetFansList(param *FansRequest) ([]*model.FansVO, int, *errcode.Error) {
	rowCount, err := s.dao.CountMyFans(param.MyId)
	if err != nil {
		return nil, 0, errcode.ServerError
	}
	fansList, err := s.dao.QueryMyFansList(param.MyId, param.Page, param.PageSize)
	if err != nil {
		return nil, 0, errcode.ServerError
	}
	//添加互粉的标记
	for _, v := range fansList {
		isFollow := redis.GetFansRelations(param.MyId, v.UserId)
		if isFollow == "1" {
			v.IsFriend = true
		}
	}
	return fansList, int(rowCount), nil
}
