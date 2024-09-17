package db

import (
	"context"
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var db *gorm.DB

func Init() error {
	if db != nil {
		return nil
	}
	host := env.Getenv("DB_HOST", "localhost")
	port := env.Getenv("DB_PORT", "3306")
	user := env.Getenv("DB_USER", "root")
	password := env.Getenv("DB_PASSWORD", "")
	database := env.Getenv("DB_NAME", "lighthouse")

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"

	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   NewCustomLogger(), // 使用自定义的日志器
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

// CustomLogger 是一个自定义的 Gorm 日志器
type CustomLogger struct {
	LogLevel logger.LogLevel
}

// NewCustomLogger 创建一个新的 CustomLogger 实例
func NewCustomLogger() logger.Interface {
	_level := env.Getenv("DB_LOG_LEVEL", "INFO")
	var level logger.LogLevel
	switch _level {
	case "INFO":
		level = logger.Info
	case "ERROR":
		level = logger.Error
	default:
		level = logger.Info
	}

	return &CustomLogger{
		LogLevel: level,
	}
}

// LogMode 实现 logger.Interface 的 LogMode 方法
func (l *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 实现 logger.Interface 的 Info 方法
func (l *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		log.Info("[GORM INFO] "+msg, data...)
	}
}

// Warn 实现 logger.Interface 的 Warn 方法
func (l *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		log.Warn("[GORM WARN] "+msg, data...)
	}
}

// Error 实现 logger.Interface 的 Error 方法
func (l *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		log.Error("[GORM ERROR] "+msg, data...)
	}
}

// Trace 实现 logger.Interface 的 Trace 方法
func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil && l.LogLevel >= logger.Error {
		log.Error("[GORM TRACE] %s | %v | %d rows | %s", err, elapsed, rows, sql)
	} else if elapsed > 200*time.Millisecond && l.LogLevel >= logger.Warn {
		log.Warn("[GORM TRACE] SLOW SQL > %v | %d rows | %s", elapsed, rows, sql)
	} else if l.LogLevel >= logger.Info {
		log.Info("[GORM TRACE] %v | %d rows | %s", elapsed, rows, sql)
	}
}
