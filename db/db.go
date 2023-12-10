package db

import (
	"github.com/light-speak/lighthouse/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var db *gorm.DB

func Init() error {
	if db != nil {
		return nil
	}
	host := env.GetEnvString("DB_HOST", "localhost")
	port := env.GetEnvString("DB_PORT", "3306")
	user := env.GetEnvString("DB_USER", "root")
	password := env.GetEnvString("DB_PASSWORD", "")
	database := env.GetEnvString("DB_NAME", "lighthouse")

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"

	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   getLogger(),
	})
	if err != nil {
		return err
	}
	db = db_
	return nil
}

func GetDb() *gorm.DB {
	return db
}
func getLogger() logger.Interface {
	//TODO: 统一日志管理

	_level := env.GetEnvString("DB_LOG_LEVEL", "info")
	var level logger.LogLevel
	switch _level {
	case "info":
		level = logger.Info
		break
	case "error":
		level = logger.Error
		break
	}

	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}
