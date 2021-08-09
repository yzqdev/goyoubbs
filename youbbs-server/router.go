package main

import (
	"github.com/gin-gonic/gin"
	"goyoubbs/controller"
	"goyoubbs/middleware"
	"goyoubbs/system"
)

func NewRouter(e *gin.Engine) {
	app := &system.Application{}

	h := controller.BaseHandler{App: app}
	baseRouter := e.Group("")
	baseRouter.POST("/login", h.UserLoginPost)
	baseRouter.POST("/register", h.Register)

	baseRouter.GET("/qqlogin", h.QQOauthHandler)
	baseRouter.GET("/oauth/qq/callback", h.QQOauthCallback)
	baseRouter.GET("/wblogin", h.WeiboOauthHandler)
	baseRouter.GET("/oauth/wb/callback", h.WeiboOauthCallback)

	sp := e.Group("/api", middleware.JwtHandler())

	{
		sp.POST("/", h.ArticleHomeList)
		//sp.POST("/view",h.ViewAtTpl)
		sp.GET("/feed", h.FeedHandler)
		//sp.GET("/robots.txt", h.Robots)
		//sp.GET("/sitemap.xml", h.SiteMapHandler)
		//sp.GET("/captcha/*",captcha.Server(captcha.StdWidth, captcha.StdHeight))
		sp.GET("/n/:cid", h.CategoryDetail)
		sp.GET("/member/:uid", h.UserDetail)
		sp.GET("/tag/:tag", h.TagDetail)
		sp.GET("/search", h.SearchDetail)

		sp.GET("/logout", h.UserLogout)
		sp.GET("/notification", h.UserNotification)

		sp.GET("/topic/:aid", h.ArticleDetail)
		sp.POST("/topic/:aid", h.ArticleDetailPost)

		sp.GET("/setting", h.UserSetting)
		sp.POST("/setting", h.UserSettingPost)

		sp.GET("/newpost/:cid", h.ArticleAdd)
		sp.POST("/newpost/:cid", h.ArticleAddPost)

		sp.POST("/content/preview", h.ContentPreviewPost)
		sp.POST("/file/upload", h.FileUpload)

		sp.GET("/admin/post/edit/:aid", h.ArticleEdit)
		sp.POST("/admin/post/edit/:aid", h.ArticleEditPost)
		sp.GET("/admin/comment/edit/:aid/:cid", h.CommentEdit)
		sp.POST("/admin/comment/edit/:aid/:cid", h.CommentEditPost)
		sp.GET("/admin/user/edit/:uid", h.UserEdit)
		sp.POST("/admin/user/edit/:uid", h.UserEditPost)
		sp.GET("/admin/user/list", h.AdminUserList)
		sp.POST("/admin/user/list", h.AdminUserListPost)
		sp.GET("/admin/category/list", h.AdminCategoryList)
		sp.POST("/admin/category/list", h.AdminCategoryListPost)
		sp.GET("/admin/link/list", h.AdminLinkList)
		sp.POST("/admin/link/list", h.AdminLinkListPost)

		sp.GET("/:filepath", h.StaticFile)
	}

}
