package main

import (
	"ginblog/controller"
	"ginblog/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
)

func InitRouter(e *gin.Engine) {

	baseRouter := e.Group("")
	baseRouter.POST("/login", controller.UserLogin)
	baseRouter.POST("/register", controller.Register)

	//baseRouter.GET("/qqlogin", controller.QQOauthHandler)
	//baseRouter.GET("/oauth/qq/callback", controller.QQOauthCallback)
	//baseRouter.GET("/wblogin", controller.WeiboOauthHandler)
	//baseRouter.GET("/oauth/wb/callback", controller.WeiboOauthCallback)

	sp := e.Group("/api", middleware.JwtHandler())

	{
		color.Redln(sp)
		//sp.POST("/", controller.ArticleHomeList)
		////sp.POST("/view",controller.ViewAtTpl)
		//sp.GET("/feed", controller.FeedHandler)
		////sp.GET("/robots.txt", controller.Robots)
		////sp.GET("/sitemap.xml", controller.SiteMapHandler)
		////sp.GET("/captcha/*",captcha.Server(captcha.StdWidth, captcha.StdHeight))
		//sp.GET("/n/:cid", controller.CategoryDetail)
		//sp.GET("/member/:uid", controller.UserDetail)
		//sp.GET("/tag/:tag", controller.TagDetail)
		//sp.GET("/search", controller.SearchDetail)
		//
		//sp.GET("/logout", controller.UserLogout)
		//sp.GET("/notification", controller.UserNotification)
		//
		//sp.GET("/topic/:aid", controller.ArticleDetail)
		//sp.POST("/topic/:aid", controller.ArticleDetailPost)
		//
		//sp.GET("/setting", controller.UserSetting)
		//sp.POST("/setting", controller.UserSettingPost)
		//
		//sp.GET("/newpost/:cid", controller.ArticleAdd)
		//sp.POST("/newpost/:cid", controller.ArticleAddPost)
		//
		//sp.POST("/content/preview", controller.ContentPreviewPost)
		//sp.POST("/file/upload", controller.FileUpload)
		//
		//sp.GET("/admin/post/edit/:aid", controller.ArticleEdit)
		//sp.POST("/admin/post/edit/:aid", controller.ArticleEditPost)
		//sp.GET("/admin/comment/edit/:aid/:cid", controller.CommentEdit)
		//sp.POST("/admin/comment/edit/:aid/:cid", controller.CommentEditPost)
		//sp.GET("/admin/user/edit/:uid", controller.UserEdit)
		//sp.POST("/admin/user/edit/:uid", controller.UserEditPost)
		//sp.GET("/admin/user/list", controller.AdminUserList)
		//sp.POST("/admin/user/list", controller.AdminUserListPost)
		//sp.GET("/admin/category/list", controller.AdminCategoryList)
		//sp.POST("/admin/category/list", controller.AdminCategoryListPost)
		//sp.GET("/admin/link/list", controller.AdminLinkList)
		//sp.POST("/admin/link/list", controller.AdminLinkListPost)
		//
		//sp.GET("/:filepath", controller.StaticFile)
	}

}
