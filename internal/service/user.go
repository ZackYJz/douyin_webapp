package service

import (
	"errors"
	"go_webapp/global"
	"go_webapp/internal/dao/redis"
	"go_webapp/internal/model"
	"go_webapp/pkg/errcode"
	"go_webapp/pkg/jwt"
	"go_webapp/pkg/oss"
	"go_webapp/pkg/sms"
	"go_webapp/pkg/snowflake"
	"go_webapp/pkg/util"
	"mime/multipart"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseUserRequest struct {
	Username string `form:"username" binding:"required,min=4,max=8"`
	Password string `form:"password" binding:"required,min=8"`
}

type MobileCodeRequest struct {
	Mobile  string `form:"mobile" binding:"required,len=11"`
	SMSCode string `form:"smsCode" binding:"required"`
}

type MobileRequest struct {
	Mobile string `form:"mobile" binding:"required,len=11"`
}
type UserIDRequest struct {
	UserId string `json:"userId" form:"userId" binding:"required"`
}

type ModifyUserInfoRequest struct {
	UserId   string        `json:"id" binding:"required"`
	Nickname string        `json:"nickname"`
	Face     string        `json:"face"`
	Sex      int           `json:"sex"`
	Birthday util.JsonTime `json:"birthday"`
	Country  string        `json:"country"`
	Province string        `json:"province"`
	City     string        `json:"city"`
	District string        `json:"district"`
	Desc     string        `json:"description"`
	BgImg    string        `json:"bgImg"`
}

type ModifyImageRequest struct {
	UserId string `form:"userId" json:"userId" binding:"required"`
	Type   int    `form:"type" json:"type" binding:"required"`
}

var DEFAULT_FACE = "https://zackyj-typora.oss-cn-chengdu.aliyuncs.com/douyin/default_face.jpeg"

func (s *Service) SendSMS(param *MobileRequest) (err error) {
	//限制同一ip 发送验证码次数
	userIp := util.GetRealIp(s.c)
	err = redis.SetSMSLimit(userIp)
	if err != nil {
		return err
	}
	smsCode := util.GenerateSmsCode(6)
	global.Logger.Info("请求发送短信验证码", zap.String("mobile", param.Mobile), zap.String("ip", userIp), zap.String("code", smsCode))
	//TODO 开发时先不真的发短信，看日志就行，节约额度
	go sms.SendSMS(param.Mobile, smsCode)

	err = redis.SavePhoneSMS(param.Mobile, smsCode)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) LoginWithMobileCode(param MobileCodeRequest) (u *model.UserVO, error *errcode.Error) {
	//从 redis 获得验证码，校验是否匹配
	redisCode, err := redis.GetSMSCode(param.Mobile)
	if err != nil {
		return nil, errcode.MobileCodeError
	}
	if redisCode == "" || redisCode != param.SMSCode {
		return nil, errcode.MobileCodeError
	}
	//查询用户是否存在
	user, err := s.dao.GetUserByMobile(param.Mobile)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errcode.ServerError
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user, err = s.CreateUser(param.Mobile)
		if err != nil {
			return nil, errcode.ServerError
		}
	}
	//用户已存在(已创建），颁布 token, 缓存刷新 token
	token, err := jwt.GenToken(user.UserId)
	//TODO Token 刷新机制
	err = redis.SetRefreshToken(user.UserId, token)
	if err != nil {
		return nil, errcode.ServerError
	}
	err = redis.DeleteSMSCode(param.Mobile)
	if err != nil {
		return nil, errcode.ServerError
	}
	userVO := user.ToUserVO()
	userVO.Token = token
	return userVO, nil
}

func (s *Service) CreateUser(mobile string) (user *model.User, err error) {
	//构造 User实例
	u := model.User{}
	//生成 userid
	u.UserId = snowflake.GenID()
	u.Username = "手机用户" + mobile
	u.Mobile = mobile
	u.Photo = DEFAULT_FACE
	u.Sex = 3
	u.Birthday = time.Date(1970, time.November, 10, 23, 0, 0, 0, time.UTC)
	u.Nickname = "手机用户" + mobile
	u.Country = "中国"
	u.Province = ""
	u.City = ""
	u.District = ""
	u.Desc = "这家伙很懒，啥也没留下"
	//插入数据库
	err = s.dao.CreateUser(&u)
	return &u, err
}

func Login(param *BaseUserRequest) (err error) {
	return nil
}

func (s *Service) QueryInfo(userId string) (userVO *model.UserVO, error *errcode.Error) {
	user, err := s.dao.GetUserByUserId(userId)
	if err != nil {
		global.Logger.Error("GetUserByUserId Error", zap.Error(err))
		return nil, errcode.ServerError
	}
	userVO = user.ToUserVO()
	follow, fans, liked, err := redis.GetUserInfoNum(userId)
	if err != nil {
		global.Logger.Error("GetUserInfoNum Error", zap.Error(err))
		return nil, errcode.ServerError
	}
	userVO.MyFollowsCounts = follow
	userVO.MyFansCounts = fans
	userVO.TotalLikeMeCounts = liked

	return userVO, nil
}

func (s *Service) ModifyUserInfo(param *ModifyUserInfoRequest, typeNum int) (user *model.User, error *errcode.Error) {
	userDto := &model.UpdateUserDTO{
		UserId:   param.UserId,
		Nickname: param.Nickname,
		Photo:    param.Face,
		Sex:      param.Sex,
		Birthday: time.Time(param.Birthday),
		Country:  param.Country,
		Province: param.Province,
		City:     param.City,
		District: param.District,
		Desc:     param.Desc,
		BgImg:    param.BgImg,
	}
	user, err := s.dao.UpdateInfoByUserId(userDto, typeNum)
	if err != nil {
		global.Logger.Error("Service.ModifyUserInfo", zap.Error(err))
		return nil, errcode.ModifyUserInfoError
	}
	return user, nil
}

func (s *Service) ModifyImage(param *ModifyImageRequest, file multipart.File, fileHeader *multipart.FileHeader) (user *model.User, error *errcode.Error) {
	fileName := oss.GetFileName(fileHeader.Filename)
	if !oss.CheckContainExt(fileName) {
		return nil, errcode.NotAllowFileExts
	}
	if oss.CheckMaxSize(file) {
		return nil, errcode.NotAllowFileSize
	}

	newImageUrl := oss.UploadOss(fileHeader, fileName)

	userDto := &model.UpdateUserDTO{
		UserId: param.UserId,
	}
	if param.Type == global.FACE {
		userDto.Photo = newImageUrl
	} else if param.Type == global.BGIMG {
		userDto.BgImg = newImageUrl
	}
	user, err := s.dao.UpdateImageByUserId(userDto)
	if err != nil {
		global.Logger.Error("Service.ModifyUserInfo", zap.Error(err))
		return nil, errcode.ModifyUserInfoError
	}
	return user, nil
}
