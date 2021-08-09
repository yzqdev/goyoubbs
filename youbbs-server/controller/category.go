package controller

import (
	"github.com/gin-gonic/gin"
	"goyoubbs/model"
	"goyoubbs/util"
	"goyoubbs/youdb"
	"net/http"
	"strconv"
)

func (h *BaseHandler) CategoryDetail(c *gin.Context) {
	btn, key, score := c.PostForm("btn"), c.PostForm("key"), c.PostForm("score")
	if len(key) > 0 {
		_, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			util.JSON(c, 400, "key", `{"retcode":400,"retmsg":"key type err"}`)
			return
		}
	}
	if len(score) > 0 {
		_, err := strconv.ParseUint(score, 10, 64)
		if err != nil {
			util.JSON(c, 400, "key", `{"retcode":400,"retmsg":"score type err"}`)
			return
		}
	}

	cid := c.PostForm("cid")
	_, err := strconv.Atoi(cid)
	if err != nil {
		util.JSON(c, 400, "key", `{"retcode":400,"retmsg":"cid type err"}`)
		return
	}

	cmd := "zrscan"
	if btn == "prev" {
		cmd = "zscan"
	}

	db := h.App.Db
	scf := h.App.Cf.Site
	cobj, err := model.CategoryGetById(db, cid)
	if err != nil {
		util.JSON(c, 400, "key", []byte(err.Error()))
		return
	}

	currentUser, _ := h.CurrentUser(w, r)

	if cobj.Hidden && currentUser.Flag < 99 {
		c.WriteHeader(http.StatusNotFound)
		util.JSON(c, 404, "key", "not found")
		return
	}
	cobj.Articles = db.Zget("category_article_num", youdb.I2b(cobj.Id)).Uint64()
	pageInfo := model.ArticleList(db, cmd, "category_article_timeline:"+cid, key, score, scf.HomeShowNum, scf.TimeZone)

	type pageData struct {
		PageData
		Cobj     model.Category
		PageInfo model.ArticlePageInfo
	}

	tpl := h.CurrentTpl(r)

	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = cobj.Name + " - " + scf.Name
	evn.Keywords = cobj.Name
	evn.Description = cobj.About
	evn.IsMobile = tpl == "mobile"

	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "category_detail"
	evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)

	evn.Cobj = cobj
	evn.PageInfo = pageInfo
	util.JSON(c, 200, "success", evn)
}
