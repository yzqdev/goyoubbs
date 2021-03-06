package controller

import (
	"encoding/json"
	"github.com/rs/xid"
	"goyoubbs/lib/qqOAuth"
	"goyoubbs/model"
	"goyoubbs/util"
	"goyoubbs/youdb"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *BaseHandler) QQOauthHandler(w http.ResponseWriter, r *http.Request) {
	scf := h.App.Cf.Site
	qq, err := qqOAuth.NewQQOAuth(strconv.Itoa(scf.QQClientID), scf.QQClientSecret, scf.MainDomain+"/oauth/qq/callback")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	// qqOAuth.Logging = true

	now := time.Now().UTC().Unix()
	qqUrlState := strconv.FormatInt(now, 10)[6:]

	urlStr, err := qq.GetAuthorizationURL(qqUrlState)
	if err != nil {
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	h.SetCookie(w, "QQUrlState", qqUrlState, 1)
	http.Redirect(w, r, urlStr, http.StatusSeeOther)
}

func (h *BaseHandler) QQOauthCallback(w http.ResponseWriter, r *http.Request) {
	qqUrlState := h.GetCookie(r, "QQUrlState")
	if len(qqUrlState) == 0 {
		w.Write([]byte(`qqUrlState cookie missed`))
		return
	}

	scf := h.App.Cf.Site
	qq, err := qqOAuth.NewQQOAuth(strconv.Itoa(scf.QQClientID), scf.QQClientSecret, scf.MainDomain+"/oauth/qq/callback")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	// qqOAuth.Logging = true

	code := r.FormValue("code")
	if code == "" {
		w.Write([]byte("Invalid code"))
		return
	}

	state := r.FormValue("state")
	if state != qqUrlState {
		w.Write([]byte("Invalid state"))
		return
	}

	token, err := qq.GetAccessToken(code)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	openid, err := qq.GetOpenID(token.AccessToken)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	timeStamp := uint64(time.Now().UTC().Unix())

	db := h.App.Db
	rs := db.Hget("oauth_qq", []byte(openid.OpenID))
	if rs.State == "ok" {
		// login
		obj := model.QQ{}
		json.Unmarshal(rs.Data[0], &obj)
		uobj, err := model.UserGetById(db, obj.Uid)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		sessionid := xid.New().String()
		uobj.LastLoginTime = timeStamp
		uobj.Session = sessionid
		jb, _ := json.Marshal(uobj)
		db.Hset("user", youdb.I2b(uobj.Id), jb)
		h.SetCookie(w, "SessionID", strconv.FormatUint(uobj.Id, 10)+":"+sessionid, 365)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	profile, err := qq.GetUserInfo(token.AccessToken, openid.OpenID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if profile.Ret != 0 {
		w.Write([]byte(profile.Message))
		return
	}

	// register

	siteCf := h.App.Cf.Site
	if siteCf.CloseReg {
		w.Write([]byte(`{"retcode":400,"retmsg":"stop to new register"}`))
		return
	}

	name := util.RemoveCharacter(profile.Nickname)
	name = strings.TrimSpace(strings.Replace(name, " ", "", -1))
	if len(name) == 0 {
		name = "qq"
	}
	nameLow := strings.ToLower(name)
	i := 1
	for {
		if db.Hget("user_name2uid", []byte(nameLow)).State == "ok" {
			i++
			nameLow = name + strconv.Itoa(i)
		} else {
			name = nameLow
			break
		}
	}

	userId, _ := db.Hincr("count", []byte("user"), 1)
	flag := 5
	if siteCf.RegReview {
		flag = 1
	}
	if userId == 1 {
		flag = 99
	}

	gender := "female"
	if profile.Gender == "???" {
		gender = "male"
	}

	uobj := model.User{
		Id:            userId,
		Name:          name,
		Flag:          flag,
		Gender:        gender,
		RegTime:       timeStamp,
		LastLoginTime: timeStamp,
		Session:       xid.New().String(),
	}

	uidStr := strconv.FormatUint(userId, 10)
	savePath := "static/avatar/" + uidStr + ".jpg"
	err = util.FetchAvatar(profile.Avatar, savePath, r.UserAgent())
	if err != nil {
		err = util.GenerateAvatar(gender, name, 73, 73, savePath)
	}
	if err != nil {
		uobj.Avatar = "0"
	} else {
		uobj.Avatar = uidStr
	}

	jb, _ := json.Marshal(uobj)
	db.Hset("user", youdb.I2b(uobj.Id), jb)
	db.Hset("user_name2uid", []byte(nameLow), youdb.I2b(userId))
	db.Hset("user_flag:"+strconv.Itoa(flag), youdb.I2b(uobj.Id), []byte(""))

	obj := model.QQ{
		Uid:    userId,
		Name:   name,
		Openid: openid.OpenID,
	}
	jb, _ = json.Marshal(obj)
	db.Hset("oauth_qq", []byte(openid.OpenID), jb)

	h.SetCookie(w, "SessionID", strconv.FormatUint(uobj.Id, 10)+":"+uobj.Session, 365)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
