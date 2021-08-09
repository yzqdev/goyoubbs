package main

import (
	"github.com/gin-gonic/gin"
	"goyoubbs/controller"
	"goyoubbs/goji/pat"
	"goyoubbs/middleware"
	"goyoubbs/system"
)

func NewRouter(e *gin.Engine) {
	app := &system.Application{}
	sp := e.Group("")
	h := controller.BaseHandler{App: app}
	sp.POST("/", h.ArticleHomeList)
	//sp.POST("/view",h.ViewAtTpl)
	sp.GET("/feed", h.FeedHandler)
	//sp.GET("/robots.txt", h.Robots)
	//sp.GET("/sitemap.xml", h.SiteMapHandler)
	//sp.GET("/captcha/*",captcha.Server(captcha.StdWidth, captcha.StdHeight))
	sp.GET("/n/:cid", h.CategoryDetail, middleware.JwtHandler())
	sp.GET("/member/:uid", h.UserDetail)
	sp.GET("/tag/:tag", h.TagDetail, middleware.JwtHandler())
	sp.GET("/search", h.SearchDetail, middleware.JwtHandler())

	sp.GET("/logout", h.UserLogout)
	sp.GET("/notification", h.UserNotification, middleware.JwtHandler())

	sp.GET("/topic/:aid", h.ArticleDetail, middleware.JwtHandler())
	sp.POST("/topic/:aid", h.ArticleDetailPost, middleware.JwtHandler())

	sp.GET("/setting", h.UserSetting, middleware.JwtHandler())
	sp.POST("/setting", h.UserSettingPost, middleware.JwtHandler())

	sp.GET("/newpost/:cid", h.ArticleAdd)
	sp.POST("/newpost/:cid", h.ArticleAddPost)

	sp.GET("/login", h.UserLogin)
	sp.POST("/login", h.UserLoginPost)
	sp.GET("/register", h.UserLogin)
	sp.POST("/register", h.UserLoginPost)

	sp.GET("/qqlogin", h.QQOauthHandler)
	sp.GET("/oauth/qq/callback", h.QQOauthCallback)
	sp.GET("/wblogin", h.WeiboOauthHandler)
	sp.GET("/oauth/wb/callback", h.WeiboOauthCallback)

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
