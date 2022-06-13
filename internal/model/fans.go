package model

import "gorm.io/gorm"

type Fans struct {
	*gorm.Model
	PublisherId string
	FanId       string
	IsFriend    int
}

func (f Fans) TableName() string {
	return "tb_fans"
}
