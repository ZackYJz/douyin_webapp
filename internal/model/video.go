package model

import (
	"strconv"

	"gorm.io/gorm"
)

type Video struct {
	*gorm.Model
	Publisher      int64
	Url            string
	Cover          string
	Title          string
	Width          int
	Height         int
	LikeCounts     int
	CommentsCounts int
	IsPrivate      int
}

type VideoSimpleVO struct {
	ID             string `json:"vlogId"`
	Publisher      string `json:"vlogerId"`
	Title          string `json:"content"`
	Url            string `json:"url"`
	Cover          string `json:"cover"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	LikeCounts     int    `json:"likeCounts"`
	CommentsCounts int    `json:"commentsCounts"`
	IsPrivate      int    `json:"isPrivate"`
}

type VideoVO struct {
	ID                 string `json:"vlogId"`
	Publisher          string `json:"vlogerId"`
	PublisherFace      string `json:"vlogerFace"`
	PublisherName      string `json:"vlogerName"`
	Title              string `json:"content"`
	Url                string `json:"url"`
	Cover              string `json:"cover"`
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	LikeCounts         int    `json:"likeCounts"`
	CommentsCounts     int    `json:"commentsCounts"`
	IsPrivate          int    `json:"isPrivate"`
	IsPlay             bool   `json:"isPlay"`
	DoIFollowPublisher bool   `json:"doIFollowVloger"`
	DoILikeThisVideo   bool   `json:"doILikeThisVlog"`
}

func (v Video) GetSimpleVO() *VideoSimpleVO {
	return &VideoSimpleVO{
		ID:             strconv.Itoa(int(v.ID)),
		Publisher:      strconv.FormatInt(v.Publisher, 10),
		Title:          v.Title,
		Url:            v.Url,
		Cover:          v.Cover,
		Width:          v.Width,
		Height:         v.Height,
		LikeCounts:     v.LikeCounts,
		CommentsCounts: v.CommentsCounts,
		IsPrivate:      v.IsPrivate,
	}
}

func (v Video) TableName() string {
	return "tb_video"
}

func (v VideoSimpleVO) TableName() string {
	return "tb_video"
}
