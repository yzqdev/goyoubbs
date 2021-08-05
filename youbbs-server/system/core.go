package system

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/viper"
	"goyoubbs/youdb"
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
