package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"github.com/gorilla/securecookie"
	"goyoubbs/cronjob"
	"goyoubbs/system"
	"goyoubbs/util"
	"goyoubbs/youdb"
	"net/url"
	"os"
	"runtime"
	"strings"
)

func GlobalData(c *gin.Context) {
	c := system.LoadConfig()
	mcf := &system.MainConf{}
	err := c.UnmarshalKey("Main", mcf)
	if err != nil {

		color.Redln(err)
	}
	// check domain
	if strings.HasPrefix(mcf.Domain, "http") {
		dm, err := url.Parse(mcf.Domain)
		if err != nil {
			log.Fatal("domain fmt err", err)
		}
		mcf.Domain = dm.Host
	} else {
		mcf.Domain = strings.Trim(mcf.Domain, "/")
	}

	scf := &system.SiteConf{}
	err2 := c.UnmarshalKey("Site", scf)
	if err2 != nil {
		return
	}
	scf.GoVersion = runtime.Version()
	fMd5, _ := util.HashFileMD5(os.Args[0])
	scf.MD5Sums = fMd5
	scf.MainDomain = strings.Trim(scf.MainDomain, "/")
	log.Println("MainDomain:", scf.MainDomain)
	if scf.TimeZone < -12 || scf.TimeZone > 12 {
		scf.TimeZone = 0
	}
	if scf.UploadMaxSize < 1 {
		scf.UploadMaxSize = 1
	}
	scf.UploadMaxSizeByte = int64(scf.UploadMaxSize) << 20

	app.Cf = &system.AppConf{mcf, scf}
	color.Redln("打开数据库")
	db, err := youdb.Open(mcf.Youdb)
	if err != nil {
		log.Fatalf("Connect Error: %v", err)
	}
	app.Db = db

	defer app.Db.Close()
	// set main node
	db.Hset("keyValue", []byte("main_category"), []byte(scf.MainNodeIds))

	var hashKey []byte
	var blockKey []byte
	if scf.ResetCookieKey {
		hashKey = securecookie.GenerateRandomKey(64)
		blockKey = securecookie.GenerateRandomKey(32)
		_ = db.Hmset("keyValue", []byte("hashKey"), hashKey, []byte("blockKey"), blockKey)
	} else {
		hashKey = append(hashKey, db.Hget("keyValue", []byte("hashKey")).Bytes()...)
		blockKey = append(blockKey, db.Hget("keyValue", []byte("blockKey")).Bytes()...)
		if len(hashKey) == 0 {
			hashKey = securecookie.GenerateRandomKey(64)
			blockKey = securecookie.GenerateRandomKey(32)
			_ = db.Hmset("keyValue", []byte("hashKey"), hashKey, []byte("blockKey"), blockKey)
		}
	}

	app.Sc = securecookie.New(hashKey, blockKey)
	// cron job
	cr := cronjob.BaseHandler{App: app}
	go cr.MainCronJob()

	staticPath := mcf.PubDir
	if len(staticPath) == 0 {
		staticPath = "static"
	}
	util.JSON(c, 200, "success", "data")
}
