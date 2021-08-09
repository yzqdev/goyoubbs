package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"goyoubbs/model"
	"goyoubbs/util"
	"strings"
)

func (h *BaseHandler) SearchDetail(c *gin.Context) {
	userContext, exist := c.Get("user")
	if !exist {
		color.Danger.Println("失败了")
	}
	//查询用户组及该组的功能权限
	currentUser, ok := userContext.(model.User) //这个是类型推断,判断接口是什么类型
	if !ok {

		color.Danger.Println("断言失败")
	}
	if currentUser.Id == 0 {
		//http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	q := c.PostForm("q")

	qLow := strings.ToLower(q)

	db := h.App.Db
	scf := h.App.Cf.Site

	where := "title"
	if strings.HasPrefix(qLow, "c:") {
		where = "content"
		qLow = qLow[2:]
	}

	pageInfo := model.ArticleSearchList(db, where, qLow, scf.PageShowNum, scf.TimeZone)

	type pageData struct {
		PageData
		Q        string
		PageInfo model.ArticlePageInfo
	}

	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = qLow + " - " + scf.Name
	//evn.IsMobile = tpl == "mobile"

	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "category_detail"
	evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)

	evn.Q = qLow
	evn.PageInfo = pageInfo
	util.ErrJSON(c, 200, "success", evn)
}
