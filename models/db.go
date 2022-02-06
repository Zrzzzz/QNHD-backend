package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"qnhd/pkg/logging"
	"qnhd/pkg/setting"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB
var sqlDB *sql.DB

func Setup() {
	var err error
	database := setting.DatabaseSetting.Database
	user := setting.DatabaseSetting.User
	password := setting.DatabaseSetting.Password
	host := setting.DatabaseSetting.Host
	port := setting.DatabaseSetting.Port
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host,
		user, password, database, port)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "qnhd.",
			SingularTable: true,
		},
	})
	if err != nil {
		logging.Fatal("Fail to open database: %v", err)
	}
	sqlDB, _ = db.DB()
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
}
func Close() {
	sqlDB.Close()
}
