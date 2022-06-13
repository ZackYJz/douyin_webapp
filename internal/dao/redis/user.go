package redis

import (
	"go_webapp/global"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

//
// SetSMSLimit
// @Description: 根据ip设置验证码限制
//
func SetSMSLimit(ip string) error {
	err := rdb.Set(SMSCODE+":"+ip, ip, time.Second*120).Err()
	if err != nil {
		global.Logger.Error("SetSMSLimit Error")
		return err
	}
	return nil
}

func IsUserSMSHas(ip string) bool {
	ipp, err := rdb.Get(SMSCODE + ":" + ip).Result()
	if err == redis.Nil || ipp == "" {
		return false
	} else if err != nil {
		return true
	}
	return true
}

//
// SavePhoneSMS
// @Description: 保存手机号-验证码 到 Redis,设置 60 秒过期时间
//
func SavePhoneSMS(mobile, code string) error {
	err := rdb.Set(SMSCODE+":"+mobile, code, time.Second*60).Err()
	if err != nil {
		global.Logger.Error("SavePhoneSMS Error")
		return err
	}
	return nil
}
func DeleteSMSCode(mobile string) error {
	err := rdb.Del(SMSCODE + ":" + mobile).Err()
	if err != nil {
		global.Logger.Error("DeleteSMSCode Error")
		return err
	}
	return nil
}

func GetSMSCode(mobile string) (code string, err error) {
	code, err = rdb.Get(SMSCODE + ":" + mobile).Result()
	if err != nil {
		global.Logger.Error("GetSMSCode Error")
		return "", err
	}
	return
}

func SetRefreshToken(userid int64, token string) (err error) {
	err = rdb.Set(USERTOKEN+":"+strconv.FormatInt(userid, 10), token, time.Hour*global.JWTSetting.Expire*2).Err()
	return
}

func RemoveRefreshToken(userid int64) (err error) {
	return rdb.Del(USERTOKEN + ":" + strconv.FormatInt(userid, 10)).Err()
}

func GetUserInfoNum(userId string) (myFollowsCounts, myFansCounts, myVideoLikedCounts int, err error) {
	//userIdStr := strconv.FormatInt(userid, 10)
	myFollowsCountsStr, err := rdb.Get(MY_FOLLOWS_COUNTS + ":" + userId).Result()
	if err != nil {
		myFollowsCounts = 0
		err = nil
	}
	myFollowsCounts, _ = strconv.Atoi(myFollowsCountsStr)
	myFansCountsStr, err := rdb.Get(MY_FANS_COUNTS + ":" + userId).Result()
	if err != nil {
		myFansCounts = 0
		err = nil
	}
	myFansCounts, _ = strconv.Atoi(myFansCountsStr)
	myVideoLikedCountsStr, err := rdb.Get(VLOGER_BE_LIKED_COUNTS + ":" + userId).Result()
	if err != nil {
		myFansCounts = 0
		err = nil
	}
	myVideoLikedCounts, _ = strconv.Atoi(myVideoLikedCountsStr)
	global.Logger.Info("用户首页数字", zap.String("id", userId),
		zap.String("关注数", myFollowsCountsStr),
		zap.String("粉丝数", myFansCountsStr),
		zap.String("获赞数", myVideoLikedCountsStr))
	return
}
