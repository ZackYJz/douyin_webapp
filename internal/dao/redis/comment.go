package redis

import "github.com/go-redis/redis"

func AddVideoCommentCount(videoId string) error {
	return rdb.Incr(VIDEO_COMMENT_COUNTS + ":" + videoId).Err()
}

func DecrCommentNum(videoId string) error {
	return rdb.Decr(VIDEO_COMMENT_COUNTS + ":" + videoId).Err()
}

func GetVideoCommentCount(videoId string) (string, error) {
	return rdb.Get(VIDEO_COMMENT_COUNTS + ":" + videoId).Result()
}

func GetCommentLikedCount(commentId string) (string, error) {
	return rdb.HGet(COMMENT_LIKED_COUNTS, commentId).Result()
}

func GetDoILikeComment(userId, commentId string) (string, error) {
	return rdb.HGet(USER_LIKE_COMMENT, userId+":"+commentId).Result()
}

func LikeComment(userId, commentId string) error {
	isLiked, err := rdb.HGet(USER_LIKE_COMMENT, userId+":"+commentId).Result()
	if isLiked == "1" {
		return nil
	}
	var eerr error
	if err == redis.Nil {
		pipeline := rdb.TxPipeline()
		rdb.HIncrBy(COMMENT_LIKED_COUNTS, commentId, 1)
		rdb.HSet(USER_LIKE_COMMENT, userId+":"+commentId, "1")
		_, eerr = pipeline.Exec()
	}

	return eerr
}

func UnlikeComment(userId, commentId string) error {
	_, err := rdb.HGet(USER_LIKE_COMMENT, userId+":"+commentId).Result()
	var eerr error
	if err != redis.Nil {
		pipeline := rdb.TxPipeline()
		rdb.HIncrBy(COMMENT_LIKED_COUNTS, commentId, -1)
		rdb.HDel(USER_LIKE_COMMENT, userId+":"+commentId)
		_, eerr = pipeline.Exec()
	}
	return eerr
}
