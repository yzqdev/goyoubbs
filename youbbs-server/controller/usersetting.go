package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gookit/color"
	"github.com/rs/xid"
	"goyoubbs/model"
	"goyoubbs/util"
	"goyoubbs/youdb"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (h *BaseHandler) UserSetting(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := h.CurrentUser(w, r)
	if currentUser.Id == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	type pageData struct {
		PageData
		Uobj model.User
		Now  int64
	}

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = h.App.Cf.Site
	evn.Title = "设置"
	evn.Keywords = ""
	evn.Description = ""
	evn.IsMobile = tpl == "mobile"
	evn.CurrentUser = currentUser

	evn.ShowSideAd = true
	evn.PageName = "user_setting"

	evn.Uobj = currentUser
	evn.Now = time.Now().UTC().Unix()

	h.SetCookie(w, "token", xid.New().String(), 1)
	h.Render(w, tpl, evn, "layout.html", "usersetting.html")
}

func (h *BaseHandler) UserSettingPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	token := h.GetCookie(r, "token")
	if len(token) == 0 {
		w.Write([]byte(`{"retcode":400,"retmsg":"token cookie missed"}`))
		return
	}

	currentUser, _ := h.CurrentUser(w, r)
	if currentUser.Id == 0 {
		w.Write([]byte(`{"retcode":401,"retmsg":"authored require"}`))
		return
	}

	// r.ParseForm() // don't use ParseForm !important
	act := r.FormValue("act")
	if act == "avatar" {

		r.ParseMultipartForm(32 << 20)

		file, _, err := r.FormFile("avatar")
		defer file.Close()

		buff := make([]byte, 512)
		file.Read(buff)
		if len(util.CheckImageType(buff)) == 0 {
			w.Write([]byte(`{"retcode":400,"retmsg":"unknown image format"}`))
			return
		}

		var imgData bytes.Buffer
		file.Seek(0, 0)

		if fileSize, err := io.Copy(&imgData, file); err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"read image data err ` + err.Error() + `"}`))
			return
		} else {
			if fileSize > 5360690 {
				w.Write([]byte(`{"retcode":400,"retmsg":"image size too much"}`))
				return
			}
		}

		img, err := util.GetImageObj(&imgData)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"fail to get image obj ` + err.Error() + `"}`))
			return
		}

		uid := strconv.FormatUint(currentUser.Id, 10)
		color.Redln(GetAppHome("/avatar/"))
		err = util.AvatarResize(img, 73, 73, GetAppHome("/avatar/")+uid+".jpg")
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"fail to resize avatar ` + err.Error() + `"}`))
			return
		}

		if currentUser.Avatar == "0" || len(currentUser.Avatar) == 0 {
			currentUser.Avatar = uid
			jb, _ := json.Marshal(currentUser)
			h.App.Db.Hset("user", youdb.I2b(currentUser.Id), jb)
		}

		http.Redirect(w, r, "/setting#2", http.StatusSeeOther)
		return
	}

	type recForm struct {
		Act       string `json:"act"`
		Email     string `json:"email"`
		Url       string `json:"url"`
		About     string `json:"about"`
		Password0 string `json:"password0"`
		Password  string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var rec recForm
	err := decoder.Decode(&rec)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"json Decode err:` + err.Error() + `"}`))
		return
	}
	defer r.Body.Close()

	recAct := rec.Act
	if len(recAct) == 0 {
		w.Write([]byte(`{"retcode":400,"retmsg":"missed act "}`))
		return
	}

	isChanged := false
	if recAct == "info" {
		currentUser.Email = rec.Email
		currentUser.Url = rec.Url
		currentUser.About = rec.About
		isChanged = true
	} else if recAct == "change_pw" {
		if len(rec.Password0) == 0 || len(rec.Password) == 0 {
			w.Write([]byte(`{"retcode":400,"retmsg":"missed args"}`))
			return
		}
		if currentUser.Password != rec.Password0 {
			w.Write([]byte(`{"retcode":400,"retmsg":"当前密码不正确"}`))
			return
		}
		currentUser.Password = rec.Password
		isChanged = true
	} else if recAct == "set_pw" {
		if len(rec.Password) == 0 {
			w.Write([]byte(`{"retcode":400,"retmsg":"missed args"}`))
			return
		}
		currentUser.Password = rec.Password
		isChanged = true
	}

	if isChanged {
		jb, _ := json.Marshal(currentUser)
		h.App.Db.Hset("user", youdb.I2b(currentUser.Id), jb)
	}

	type response struct {
		normalRsp
	}

	rsp := response{}
	rsp.Retcode = 200
	rsp.Retmsg = "修改成功"
	json.NewEncoder(w).Encode(rsp)
}
