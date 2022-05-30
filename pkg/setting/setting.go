package setting

import (
	"log"
	"os"
	"time"

	"github.com/go-ini/ini"
)

type Server struct {
	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type App struct {
	JwtSecret       string
	TokenExpireTime int

	RuntimeRootPath string

	GinLogSavePath  string
	GormLogSavePath string
	LogSavePath     string
	LogSaveName     string
	LogFileExt      string
	TimeFormat      string

	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string
	ImagePrefixUrl string

	PageSize int

	WPYDomain    string
	WPYAppSecret string
	WPYAppKey    string

	WPYLoginAc string
	WPYLoginPw string

	YunPianAppKey string

	// 发贴间隔时间
	TimeLimit       int
	EnableTimeLimit bool
}

type Database struct {
	User     string
	Password string
	Host     string
	Database string
	Port     string
}
type Environment struct {
	DB_DEBUG     string
	QNHD_REFRESH string
	RELEASE      string
}

var ServerSetting = &Server{}
var AppSetting = &App{}
var DatabaseSetting = &Database{}
var EnvironmentSetting = &Environment{}

func setupEnvironment() {
	EnvironmentSetting.DB_DEBUG = os.Getenv("DB_DEBUG")
	EnvironmentSetting.QNHD_REFRESH = os.Getenv("QNHD_REFRESH")
	EnvironmentSetting.RELEASE = os.Getenv("RELEASE")
}

func Setup() {
	Cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}
	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo AppSetting err: %v", err)
	}
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize << 20

	err = Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo ServerSetting err: %v", err)
	}
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.ReadTimeout * time.Second

	err = Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo DatabaseSetting err: %v", err)
	}

	setupEnvironment()
}
