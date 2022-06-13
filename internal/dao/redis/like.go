package redis

import (
	"go_webapp/pkg/util"

	"github.com/go-redis/redis"
)

func LikeVideo(userId, videoId, publisher string) error {
	pipe := rdb.TxPipeline()
	//给视频发布者获赞数 + 1
	rdb.Incr(VLOGER_BE_LIKED_COUNTS + ":" + publisher)
	//给视频获赞数 + 1
	rdb.Incr(VLOG_BE_LIKED_COUNTS + ":" + videoId)
	//保存点赞关系
	rdb.Set(USER_LIKE_VIDEO+":"+userId+":"+videoId, "1", 0)
	_, err := pipe.Exec()
	return err
}

func UnLikeVideo(userId, videoId, publisher string) error {
	pipe := rdb.TxPipeline()
	//给视频发布者获赞数 - 1
	rdb.Decr(VLOGER_BE_LIKED_COUNTS + ":" + publisher)
	//给视频获赞数 - 1
	rdb.Decr(VLOG_BE_LIKED_COUNTS + ":" + videoId)
	//删除点赞关系
	rdb.Del(USER_LIKE_VIDEO + ":" + userId + ":" + videoId)
	_, err := pipe.Exec()
	return err
}

func GetVideoLikedCounts(videoId string) (int, error) {
	count, err := rdb.Get(VLOG_BE_LIKED_COUNTS + ":" + videoId).Result()
	if err != nil {
		//global.Logger.Error("GetVideoLikedCounts", zap.Error(err))
		return 0, nil
	}
	return util.StrTo(count).MustInt(), nil
}

func DoILikeVideo(myId string, videoId string) bool {
	result, err := rdb.Get(USER_LIKE_VIDEO + ":" + myId + ":" + videoId).Result()
	if result == "1" {
		return true
	} else if err == redis.Nil {
		return false
	} else {
		return false
	}
}
