package db

import (
	"github.com/light-speak/lighthouse/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() error {
	host := env.GetEnvString("DB_HOST", "localhost")
	port := env.GetEnvString("DB_PORT", "3006")
	user := env.GetEnvString("DB_USER", "root")
	password := env.GetEnvString("DB_PASSWORD", "")
	database := env.GetEnvString("DB_NAME", "lighthouse")

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db = db_
	return nil
}

func GetDb() *gorm.DB {
	return db
}
