package model

import (
	"context"
	"fmt"
	"time"

	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	log.Info().Msg(env.LighthouseConfig.App.Name)

	host := env.LighthouseConfig.Database.Host
	port := env.LighthouseConfig.Database.Port
	user := env.LighthouseConfig.Database.User
	password := env.LighthouseConfig.Database.Password
	dbName := env.LighthouseConfig.Database.Name

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai&timeout=10s", user, password, host, port, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   &DBLogger{LogLevel: logger.Info},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to connect database")
	}

	DB = db
}

func GetDB() *gorm.DB { return DB }

type DBLogger struct {
	LogLevel logger.LogLevel
}

func (l *DBLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *DBLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Info().Msgf(msg, data...)
}

func (l *DBLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Warn().Msgf(msg, data...)
}

func (l *DBLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Error().Msgf(msg, data...)
}

func (l *DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil && l.LogLevel >= logger.Error {
		log.Error().Err(err).Str("sql", sql).Int64("rows", rows).Msg("database error")
	} else if elapsed > 200*time.Millisecond && l.LogLevel >= logger.Warn {
		log.Warn().Str("sql", sql).Int64("rows", rows).Msg("database slow query")
	} else if l.LogLevel >= logger.Info {
		log.Info().Str("sql", sql).Int64("rows", rows).Msg("database query")
	}
}
