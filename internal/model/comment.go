package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	*gorm.Model
	Publisher  int64
	VideoId    int64
	FatherId   int64
	UserId     int64
	Content    string
	LikeCounts int
}

type CommentVO struct {
	Id                  string    `json:"id"`
	CommentId           string    `json:"commentId"`
	Publisher           string    `json:"vlogerId"`
	FatherId            string    `json:"fatherCommentId"`
	VideoId             string    `json:"vlogId"`
	UserId              string    `json:"commentUserId"`
	Nickname            string    `json:"commentUserNickname"`
	Face                string    `json:"commentUserFace"`
	Content             string    `json:"content"`
	LikeCounts          int       `json:"likeCounts"`
	ReplyedUserNickname string    `json:"replyedUserNickname"`
	CreateTime          time.Time `json:"createTime"`
	IsLike              int       `json:"isLike"`
}

func (f Comment) GetVO() *CommentVO {
	return &CommentVO{
		Id:         strconv.Itoa(int(f.ID)),
		Publisher:  strconv.FormatInt(f.Publisher, 10),
		FatherId:   strconv.FormatInt(f.FatherId, 10),
		VideoId:    strconv.FormatInt(f.VideoId, 10),
		UserId:     strconv.FormatInt(f.UserId, 10),
		Content:    f.Content,
		LikeCounts: f.LikeCounts,
		CreateTime: f.CreatedAt,
	}
}

func (f Comment) TableName() string {
	return "tb_comment"
}
