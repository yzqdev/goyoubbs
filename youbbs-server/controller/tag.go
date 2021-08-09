package controller

import (
	"github.com/gin-gonic/gin"
	"goyoubbs/goji/pat"
	"goyoubbs/model"
	"net/http"
	"strconv"
	"strings"
)

func (h *BaseHandler) TagDetail(c *gin.Context) {
	btn, key := c.PostForm("btn"), c.PostForm("key")
	if len(key) > 0 {
		_, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"key type err"}`))
			return
		}
	}

	tag := c.Param("tag")
	tagLow := strings.ToLower(tag)

	cmd := "hrscan"
	if btn == "prev" {
		cmd = "hscan"
	}

	db := h.App.Db
	scf := h.App.Cf.Site
	rs := db.Hscan("tag:"+tagLow, nil, 1)
	if rs.State != "ok" {
		//w.Write([]byte(`{"retcode":404,"retmsg":"not found"}`))
		return
	}

	currentUser, _ := h.CurrentUser(w, r)

	pageInfo := model.UserArticleList(db, cmd, "tag:"+tagLow, key, scf.PageShowNum, scf.TimeZone)

	type tagDetail struct {
		Name   string
		Number uint64
	}
	type pageData struct {
		PageData
		Tag      tagDetail
		PageInfo model.ArticlePageInfo
	}

	tpl := h.CurrentTpl(r)

	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = tag + " - " + scf.Name
	evn.Keywords = tag
	evn.Description = tag
	evn.IsMobile = tpl == "mobile"

	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "category_detail"
	evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)

	evn.Tag = tagDetail{
		Name:   tag,
		Number: db.Zget("tag_article_num", []byte(tagLow)).Uint64(),
	}
	evn.PageInfo = pageInfo

	h.Render(w, tpl, evn, "layout.html", "tag.html")
}
