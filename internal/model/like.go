package model

import "gorm.io/gorm"

type Like struct {
	*gorm.Model
	UserId  int64
	VideoId int64
}

func (l Like) TableName() string {
	return "tb_like"
}
