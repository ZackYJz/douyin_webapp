package routers

import (
	"go_webapp/internal/middleware"
	"go_webapp/internal/routers/api"
	"go_webapp/pkg/logger"

	_ "go_webapp/docs"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	//跨域中间件
	r.Use(middleware.Cors())
	//替换为 zapper 的日志中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//r.Use(middleware.AccessLog())
	//错误消息的国际化
	//r.Use(middleware.Translations())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	user := api.NewUser()
	video := api.NewVideo()
	fans := api.NewFans()
	comment := api.NewComment()

	api := r.Group("/douyin")
	api.GET("/token", user.GetToken)
	{
		/*
			测试模块
		*/
		demo := api.Group("/demo")
		demo.Use(middleware.JWTAuthMiddleware())
		{
			demo.GET("/test", user.Test)
			//demo.GET("/token", user.GetToken)
		}
		/*
			用户模块
		*/
		userGroup := api.Group("/passport")
		{
			userGroup.POST("/register", user.RegisterHandler)
			userGroup.POST("/login", user.MobileLoginHandler)
			//userGroup.POST("/getSMSCode", middleware.SMSCodeLimitMiddleware(), user.GetSMSCodeHandler)
			userGroup.POST("/getSMSCode", user.GetSMSCodeHandler)
			userGroup.POST("/logout", user.LogoutHandler)
		}
		/*
			用户信息模块
		*/
		userInfoGroup := api.Group("/userInfo")
		{
			userInfoGroup.POST("/modifyImage", middleware.JWTAuthMiddleware(), user.ModifyImageHandler)
			userInfoGroup.POST("/modifyUserInfo", middleware.JWTAuthMiddleware(), user.ModifyUserInfoHandler)
			userInfoGroup.GET("/query", user.QueryInfoHandler)
		}
		/*
			视频模块
		*/
		api.GET("/vlog/indexList", video.IndexListHandler)
		api.GET("/vlog/detail", video.VideoDetailHandler)
		api.POST("/vlog/totalLikedCounts", video.TotalLikedCountsHandler)
		videoGroup := api.Group("/vlog", middleware.JWTAuthMiddleware())
		{
			videoGroup.POST("/publish", video.PublishHandler)
			videoGroup.GET("/myPublicList", video.PublicListHandler)
			videoGroup.GET("/myPrivateList", video.PrivateListHandler)
			videoGroup.POST("/changeToPublic", video.ToPublicHandler)
			videoGroup.POST("/changeToPrivate", video.ToPrivateHandler)
			videoGroup.GET("/followList", video.FollowVideoListHandler)
			videoGroup.GET("/friendList", video.FriendVideoListHandler)
			videoGroup.POST("/like", video.LikeHandler)
			videoGroup.POST("/unlike", video.UnLikeHandler)
			videoGroup.GET("/myLikedList", video.LikedListHandler)
		}
		/*
			粉丝关注模块
		*/
		fansGroup := api.Group("/fans", middleware.JWTAuthMiddleware())
		{
			fansGroup.POST("/cancel", fans.CancelFollowHandler)
			fansGroup.POST("/follow", fans.FollowHandler)
			fansGroup.GET("/queryDoIFollowVloger", fans.DoIFollowPublisherHandler)
			fansGroup.GET("/queryMyFollows", fans.MyFollowListHandler)
			fansGroup.GET("/queryMyFans", fans.MyFansListHandler)
		}
		commentGroup := api.Group("/comment")
		{
			commentGroup.POST("/create", middleware.JWTAuthMiddleware(), comment.CreateCommentHandler)
			commentGroup.GET("/list", comment.CommentListHandler)
			commentGroup.GET("/counts", comment.CommentCountsHandler)
			commentGroup.DELETE("/delete", middleware.JWTAuthMiddleware(), comment.DeleteHandler)
			commentGroup.POST("/like", comment.LikeCommentHandler)
			commentGroup.POST("/unlike", comment.UnLikeCommentHandler)
		}
	}
	return r
}
