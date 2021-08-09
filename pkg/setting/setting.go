package setting

import (
	"time"

	"qnhd/pkg/logging"

	"github.com/go-ini/ini"
)

var (
	Cfg          *ini.File
	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	OfficeEmail, OfficePass string
	EmailSmtp, EmailPort    string

	JwtSecret            string
	AdminName, AdminPass string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		logging.Fatal("Fail to parse 'conf/app.ini': %v", err)
	}
	LoadBase()
	LoadServer()
	LoadApp()
}
func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}
func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		logging.Fatal("Fail to get section 'server': %v", err)
	}
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}
func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		logging.Fatal("Fail to get section 'app': %v", err)
	}
	OfficeEmail = sec.Key("EMAIL").MustString("")
	OfficePass = sec.Key("PASSWORD").MustString("")
	EmailSmtp = sec.Key("SMTP").MustString("")
	EmailPort = sec.Key("PORT").MustString("")
	JwtSecret = sec.Key("JWTSECRET").MustString("")
	AdminName = sec.Key("ADMINNAME").MustString("")
	AdminPass = sec.Key("ADMINPASS").MustString("")

	if OfficeEmail == "" || OfficePass == "" || EmailSmtp == "" || EmailPort == "" {
		logging.Fatal("Failed to init App because lacking of email or password keyword")
	}

}
