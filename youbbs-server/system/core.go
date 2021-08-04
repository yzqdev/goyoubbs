package system

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"runtime"

	"github.com/gorilla/securecookie"
	"github.com/qiniu/api.v7/storage"
	"goyoubbs/util"
	"goyoubbs/youdb"
	"net/url"
	"strings"
)

type MainConf struct {
	HttpPort       int
	HttpsOn        bool
	Domain         string // 若启用https 则该domain 为注册的域名，eg: domain.com、www.domain.com
	HttpsPort      int
	PubDir         string
	ViewDir        string
	Youdb          string
	CookieSecure   bool
	CookieHttpOnly bool
	OldSiteDomain  string
	TLSCrtFile     string
	TLSKeyFile     string
}

type SiteConf struct {
	GoVersion         string
	MD5Sums           string
	Name              string
	Desc              string
	AdminEmail        string
	MainDomain        string // 上传图片后添加网址前缀, eg: http://domian.com 、http://234.21.35.89:8082
	MainNodeIds       string
	TimeZone          int
	HomeShowNum       int
	PageShowNum       int
	TagShowNum        int
	CategoryShowNum   int
	TitleMaxLen       int
	ContentMaxLen     int
	PostInterval      int
	CommentListNum    int
	CommentInterval   int
	Authorized        bool
	RegReview         bool
	CloseReg          bool
	AutoDataBackup    bool
	ResetCookieKey    bool // 重设cookie key （强迫重新登录）
	AutoGetTag        bool
	GetTagApi         string
	QQClientID        int
	QQClientSecret    string
	WeiboClientID     int
	WeiboClientSecret string // eg: "jpg,jpeg,gif,zip,pdf"
	UploadSuffix      string
	UploadImgOnly     bool
	UploadImgResize   bool
	UploadMaxSize     int
	UploadMaxSizeByte int64
	QiniuAccessKey    string
	QiniuSecretKey    string
	QiniuDomain       string
	QiniuBucket       string
	UpyunDomain       string
	UpyunBucket       string
	UpyunUser         string
	UpyunPw           string
	BaiduSubUrl       string
	BingSubUrl        string
}

type AppConf struct {
	Main *MainConf
	Site *SiteConf
}

type Application struct {
	Cf     *AppConf
	Db     *youdb.DB
	Sc     *securecookie.SecureCookie
	QnZone *storage.Zone
}

func LoadConfig() *viper.Viper {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath("config") // path to look for the config file in
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	return viper.GetViper()
}

func (app *Application) Init(c *viper.Viper, currentFilePath string) {

	mcf := &MainConf{}
	err := c.UnmarshalKey("Main", mcf)
	if err != nil {

		return
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

	scf := &SiteConf{}
	err2 := c.UnmarshalKey("Site", scf)
	if err2 != nil {
		return
	}
	scf.GoVersion = runtime.Version()
	fMd5, _ := util.HashFileMD5(currentFilePath)
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

	app.Cf = &AppConf{mcf, scf}
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
	//app.Sc.SetSerializer(securecookie.JSONEncoder{})

	log.Println("youdb Connect to", mcf.Youdb)
}

func (app *Application) Close() {
	defer app.Db.Close()
	log.Println("db cloded")
}
