package redis

const (
	PREFIX    = "douyin:"
	SMSCODE   = PREFIX + "user_smscode"
	USERTOKEN = PREFIX + "user_refresh_token"

	MY_FOLLOWS_COUNTS = PREFIX + "my_follows_counts"
	MY_FANS_COUNTS    = PREFIX + "my_fans_counts"

	VLOGER_BE_LIKED_COUNTS = PREFIX + "publisher_liked_counts"
	VLOG_BE_LIKED_COUNTS   = PREFIX + "video_liked_counts"

	// USER_LIKE_VIDEO 用户点赞视频
	USER_LIKE_VIDEO = PREFIX + "user_like_video"

	// FANS_TO_VLOGGER_FOLLOW 粉丝关注博主关系
	FANS_TO_VLOGGER_FOLLOW = PREFIX + "fans_follow"

	// VIDEO_COMMENT_COUNTS 视频评论数
	VIDEO_COMMENT_COUNTS = PREFIX + "video_comment_counts"

	// COMMENT_LIKED_COUNTS 评论的点赞数
	COMMENT_LIKED_COUNTS = PREFIX + "comment_liked_counts"

	// USER_LIKE_COMMENT 用户点赞视频
	USER_LIKE_COMMENT = PREFIX + "user_like_comment"
)
