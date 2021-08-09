package controller

import (
	"encoding/json"
	"ginblog/model"
	"ginblog/utils"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gookit/color"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var SecretKey = []byte("9hUxqaGelNnCZaCW")

type NewJwtClaims struct {
	*model.User
	jwt.StandardClaims
}

type recForm struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	CaptchaId       string `json:"captchaId"`
	CaptchaSolution string `json:"captchaSolution"`
}

type response struct {
	normalRsp
}

func Register(c *gin.Context) {
	db := h.App.Db
	var rec recForm
	// register
	siteCf := h.App.Cf.Site
	if siteCf.QQClientID > 0 || siteCf.WeiboClientID > 0 {
		utils.JSON(c, 400, "err", gin.H{"retcode": 400, "retmsg": "请用QQ 或 微博一键登录"})
		return
	}
	if siteCf.CloseReg {
		utils.JSON(c, 400, "err", gin.H{"retcode": 400, "retmsg": "stop to new register"})
		return
	}
	var nameLow string
	if db.Hget("user_name2uid", []byte(nameLow)).State == "ok" {

		utils.JSON(c, 405, "err", gin.H{"retcode": 405, "retmsg": "name is exist", "newCaptchaId": "` + captcha.NewLen(2) + `"})
		return
	}

	userId, _ := db.Hincr("count", []byte("user"), 1)
	flag := 5
	if siteCf.RegReview {
		flag = 1
	}

	if userId == 1 {
		flag = 99
	}
	timeStamp := uint64(time.Now().UTC().Unix())
	uobj := model.User{
		Id:            userId,
		Name:          rec.Name,
		Password:      rec.Password,
		Flag:          flag,
		RegTime:       timeStamp,
		LastLoginTime: timeStamp,
		Session:       xid.New().String(),
	}

	uidStr := strconv.FormatUint(userId, 10)

	err := util.GenerateAvatar("male", rec.Name, 73, 73, GetAppHome("/avatar/")+uidStr+".jpg")
	if err != nil {
		uobj.Avatar = "0"
	} else {
		uobj.Avatar = uidStr
	}

	jb, _ := json.Marshal(uobj)
	db.Hset("user", youdb.I2b(uobj.Id), jb)
	db.Hset("user_name2uid", []byte(nameLow), youdb.I2b(userId))
	db.Hset("user_flag:"+strconv.Itoa(flag), youdb.I2b(uobj.Id), []byte(""))

}
func UserLoginPost(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=UTF-8")
	token := c.GetHeader("token")
	if len(token) == 0 {
		//w.Write([]byte(`{"retcode":400,"retmsg":"token cookie missed"}`))
		utils.JSON(c, 400, "error", gin.H{"retcode": 400, "retmsg": "token cookie missed"})

	}

	act := strings.TrimLeft(c.Request.Host, "/")

	var rec recForm

	if err := c.ShouldBind(rec); err != nil {
		panic(err)
	}

	if len(rec.Name) == 0 || len(rec.Password) == 0 {
		utils.JSON(c, 400, "error", gin.H{"retcode": 400, "retmsg": "name or pw is empty"})
		return
	}
	nameLow := strings.ToLower(rec.Name)
	if !util.IsUserName(nameLow) {
		utils.JSON(c, 405, "success", gin.H{"retcode": 400, "retmsg": "name fmt err"})
		return
	}

	if !captcha.VerifyString(rec.CaptchaId, rec.CaptchaSolution) {

		utils.JSON(c, 200, "error", gin.H{"retcode": 405, "retmsg": "验证码错误", "newCaptchaId": "` + captcha.NewLen(2) + `"})
		return
	}

	db := h.App.Db
	timeStamp := uint64(time.Now().UTC().Unix())
	result := &model.Result{
		Code:    200,
		Message: "登录成功n",
		Data:    nil,
	}
	if act == "login" {
		bn := "user_login_token"
		key := []byte(token + ":loginerr")
		if db.Zget(bn, key).State == "ok" {
			// todo
			//w.Write([]byte(`{"retcode":400,"retmsg":"name and pw not match"}`))
			//return
		}
		color.Redln(GetAppHome("/avatar/"))
		color.Redln("根据用户名获取用户")
		color.Redln(db)
		uobj, err := model.UserGetByName(db, nameLow)
		if err != nil {
			utils.JSON(c, 405, "error", gin.H{"retcode": 405, "retmsg": "json Decode err:` + err.Error() + `", "newCaptchaId": "` + captcha.NewLen(2) + `"})
			return
		}
		if uobj.Password != rec.Password {
			db.Zset(bn, key, uint64(time.Now().UTC().Unix()))
			utils.JSON(c, 405, "eror", gin.H{"retcode": 405, "retmsg": "name and pw not match", "newCaptchaId": "` + captcha.NewLen(2) + `"})
			return
		}
		expiresTime := time.Now().Unix() + int64(60*60*24)
		uobj.LastLoginTime = timeStamp
		stdClaims := jwt.StandardClaims{

			Audience:  "啊啊啊",             // 受众
			ExpiresAt: expiresTime,       // 失效时间
			Id:        "id",              // 编号
			IssuedAt:  time.Now().Unix(), // 签发时间
			Issuer:    "sqlU.Username",   // 签发人
			NotBefore: time.Now().Unix(), // 生效时间
			Subject:   "login",           // 主题
		}
		newClaims := NewJwtClaims{
			User:           &uobj,
			StandardClaims: stdClaims,
		}
		tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
		if token, err := tokenClaims.SignedString(SecretKey); err == nil {
			result.Message = "登录成功"
			result.Data = token
			result.Code = http.StatusOK
			jb, _ := json.Marshal(uobj)
			db.Hset("user", youdb.I2b(uobj.Id), jb)
			c.JSON(result.Code, result)
		} else {
			result.Message = "登录失败，请重新登陆"
			result.Code = http.StatusOK
			c.JSON(result.Code, gin.H{
				"result": result,
			})
		}

	}

	rsp := response{}
	rsp.Retcode = 200
}

func UserNotification(c *gin.Context) {
	userContext, exist := c.Get("user")
	if !exist {
		color.Danger.Println("失败了")
	}
	//查询用户组及该组的功能权限
	user, ok := userContext.(model.User) //这个是类型推断,判断接口是什么类型
	if !ok {

		color.Danger.Println("断言失败")
	}
	if user.Id == 0 {
		//http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	type pageData struct {
		PageData
		PageInfo model.ArticlePageInfo
	}

	db := h.App.Db
	scf := h.App.Cf.Site

	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = "站内提醒 - " + scf.Name
	//evn.IsMobile = tpl == "mobile"

	evn.CurrentUser = user
	evn.ShowSideAd = true
	evn.PageName = "user_notification"
	evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)

	evn.PageInfo = model.ArticleNotificationList(db, user.Notice, scf.TimeZone)

	// fix user.NoticeNum != len(evn.PageInfo.Items)
	if user.NoticeNum != len(evn.PageInfo.Items) {
		var newKeys []string
		for _, item := range evn.PageInfo.Items {
			newKeys = append(newKeys, strconv.FormatUint(item.Id, 10))
		}

		user.Notice = strings.Join(newKeys, ",")
		user.NoticeNum = len(newKeys)

		jb, _ := json.Marshal(user)
		db.Hset("user", youdb.I2b(user.Id), jb)

		evn.CurrentUser = user
	}
	utils.JSON(c, 200, "success", evn)

}

func (h *BaseHandler) UserLogout(c *gin.Context) {
	cks := []string{"SessionID", "QQUrlState", "WeiboUrlState", "token"}
	for _, k := range cks {
		color.Redln(k)
		//h.DelCookie(w, k)
	}
	//http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *BaseHandler) UserDetail(c *gin.Context) {
	act, btn, key, score := c.PostForm("act"), c.PostForm("btn"), c.PostForm("key"), c.PostForm("score")
	if len(key) > 0 {
		_, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			utils.JSON(c, 200, "success", []byte(`{"retcode":400,"retmsg":"key type err"}`))
			return
		}
	}
	if len(score) > 0 {
		_, err := strconv.ParseUint(score, 10, 64)
		if err != nil {
			utils.JSON(c, 200, "success", []byte(`{"retcode":400,"retmsg":"score type err"}`))
			return
		}
	}

	db := h.App.Db
	scf := h.App.Cf.Site

	uid := c.Param("uid")
	uidi, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		uid = model.UserGetIdByName(db, strings.ToLower(uid))
		if uid == "" {
			utils.JSON(c, 200, "success", []byte(`{"retcode":400,"retmsg":"uid type err"}`))
			return
		}
		//http.Redirect(w, r, "/member/"+uid, 301)
		return
	}

	cmd := "rscan"
	if btn == "prev" {
		cmd = "scan"
	}

	uobj := model.UserGetById(uidi)
	if err != nil {
		utils.JSON(c, 500, "error", []byte(err.Error()))
		return
	}
	userContext, exist := c.Get("user")
	if !exist {
		color.Danger.Println("失败了")
	}
	//查询用户组及该组的功能权限
	currentUser, ok := userContext.(model.User) //这个是类型推断,判断接口是什么类型
	if !ok {

		color.Danger.Println("断言失败")
	}

	if uobj.Hidden && currentUser.Flag < 99 {
		//w.WriteHeader(http.StatusNotFound)
		//w.Write([]byte(`{"retcode":404,"retmsg":"not found"}`))
		return
	}

	var pageInfo model.ArticlePageInfo

	if act == "reply" {
		tb := "user_article_reply:" + uid
		// pageInfo = model.UserArticleList(db, cmd, tb, key, h.App.Cf.Site.PageShowNum)
		pageInfo = model.ArticleList(db, "z"+cmd, tb, key, score, scf.PageShowNum, scf.TimeZone)
	} else {
		act = "post"
		tb := "user_article_timeline:" + uid
		pageInfo = model.UserArticleList(db, "h"+cmd, tb, key, scf.PageShowNum, scf.TimeZone)
	}

	type userDetail struct {
		model.User
		RegTimeFmt string
	}
	type pageData struct {
		PageData
		Act      string
		Uobj     userDetail
		PageInfo model.ArticlePageInfo
	}

	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = uobj.Name + " - " + scf.Name
	evn.Keywords = uobj.Name
	evn.Description = uobj.About
	//evn.IsMobile = tpl == "mobile"

	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "category_detail"
	evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)

	evn.Act = act
	evn.Uobj = userDetail{
		User:       uobj,
		RegTimeFmt: util.TimeFmt(uobj.RegTime, "2006-01-02 15:04", scf.TimeZone),
	}
	evn.PageInfo = pageInfo
	utils.JSON(c, 200, "success", evn)
}
