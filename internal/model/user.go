package model

import (
	"go_webapp/pkg/snowflake"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	UserId   int64     `json:"user_id" gorm:"type:bigint(20);index:idx_user_id"`
	Username string    `json:"username" gorm:"type:varchar(64);index:idx_user_name;unique;default:''"`
	Password string    `json:"-" gorm:"type:varchar(64);default:''"`
	Mobile   string    `json:"mobile" gorm:"type:varchar(12)"`
	Nickname string    `json:"nickname" gorm:"type:varchar(20)"`
	Photo    string    `json:"face" gorm:"type:varchar(128)"`
	Sex      int       `json:"sex" gorm:"type:int(2);default:3"`
	Birthday time.Time `json:"birthday" gorm:"type:date"`
	Desc     string    `json:"description" gorm:"type:varchar(64);default:''"`
	BgImg    string    `json:"bgImg" gorm:"type:varchar(128);default:''"`
	Country  string    `json:"country" gorm:"type:varchar(32);default:''"`
	Province string    `json:"province" gorm:"type:varchar(32);default:''"`
	City     string    `json:"city" gorm:"type:varchar(32);default:''"`
	District string    `json:"district" gorm:"type:varchar(32);default:''"`
}

type UserVO struct {
	UserId   string    `json:"id"`
	Username string    `json:"imoocNum"`
	Mobile   string    `json:"mobile"`
	Nickname string    `json:"nickname"`
	Photo    string    `json:"face"`
	Sex      int       `json:"sex"`
	Birthday time.Time `json:"birthday"`
	Desc     string    `json:"description"`
	BgImg    string    `json:"bgImg"`
	Country  string    `json:"country"`
	Province string    `json:"province"`
	City     string    `json:"city" `
	District string    `json:"district"`

	Token string `json:"userToken"`

	MyFollowsCounts   int `json:"myFollowsCounts"`
	MyFansCounts      int `json:"myFansCounts"`
	TotalLikeMeCounts int `json:"totalLikeMeCounts"`
}

type UpdateUserDTO struct {
	UserId   string    `json:"id" binding:"required"`
	Nickname string    `json:"nickname"`
	Photo    string    `json:"face"`
	Sex      int       `json:"sex"`
	Birthday time.Time `json:"birthday"`
	Country  string    `json:"country"`
	Province string    `json:"province"`
	City     string    `json:"city"`
	District string    `json:"district"`
	Desc     string    `json:"description"`
	BgImg    string    `json:"bgImg"`
}

type FollowerVO struct {
	UserId     string `json:"vlogerId"`
	Nickname   string `json:"nickname"`
	Photo      string `json:"face"`
	IsFollowed bool   `json:"followed"`
}
type FansVO struct {
	UserId   string `json:"fanId"`
	Nickname string `json:"nickname"`
	Photo    string `json:"face"`
	IsFriend bool   `json:"friend"`
}

func (u User) ToUserVO() (userVO *UserVO) {
	return &UserVO{
		UserId:            strconv.FormatInt(u.UserId, 10),
		Username:          u.Username,
		Mobile:            u.Mobile,
		Nickname:          u.Nickname,
		Photo:             u.Photo,
		Sex:               u.Sex,
		Birthday:          u.Birthday,
		Desc:              u.Desc,
		BgImg:             u.BgImg,
		Country:           u.Country,
		Province:          u.Province,
		City:              u.City,
		District:          u.District,
		Token:             "",
		MyFollowsCounts:   0,
		MyFansCounts:      0,
		TotalLikeMeCounts: 0,
	}
}

func (u User) TableName() string {
	return "tb_user"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UserId = snowflake.GenID()
	return
}
