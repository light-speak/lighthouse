package databases

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	_ "time/tzdata"

	"github.com/light-speak/lighthouse/logs"
	"github.com/light-speak/lighthouse/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type LightDatabase struct {
	MainDB    *gorm.DB
	SlaveDBs  []*gorm.DB
	Completed bool
	Error     error
}

var LightDatabaseClient *LightDatabase

func init() {
	// 使用固定的 CST 时区代替加载 Asia/Shanghai
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = loc
	databaseConfig.Main.LogLevel = databaseConfig.LogLevel
	databaseConfig.Slave.LogLevel = databaseConfig.LogLevel

	// 初始化数据库连接，添加重试机制
	initDatabaseWithRetry(loc)

}

// initDatabaseWithRetry 初始化数据库连接，添加重试机制
func initDatabaseWithRetry(loc *time.Location) {
	maxRetries := 5
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		// 尝试初始化主库
		mainDB, err := initDB(databaseConfig.Main, loc)
		if err != nil {
			logs.Error().Err(err).Int("retry", i+1).Msg("main database init error, retrying...")

			// 如果已经是最后一次尝试，则设置错误状态并返回
			if i == maxRetries-1 {
				logs.Error().Err(err).Msg("main database init failed after maximum retries")
				LightDatabaseClient = &LightDatabase{
					Completed: false,
					MainDB:    nil,
					Error:     err,
				}
				return
			}

			// 等待一段时间后重试
			time.Sleep(retryInterval)
			continue
		}

		// 初始化从库
		var slaveDBs []*gorm.DB
		if databaseConfig.EnableSlave && len(databaseConfig.Slave.Hosts) > 0 {
			for _, host := range databaseConfig.Slave.Hosts {
				slaveConfig := *databaseConfig.Slave
				slaveConfig.Hosts = []string{host}
				slaveDB, err := initDB(&slaveConfig, loc)
				if err != nil {
					logs.Error().Err(err).Str("host", host).Msg("slave database init error")
					continue
				}
				slaveDBs = append(slaveDBs, slaveDB)
			}
		}

		// 如果没有从库，使用主库作为从库
		if len(slaveDBs) == 0 {
			slaveDBs = []*gorm.DB{mainDB}
		}

		LightDatabaseClient = &LightDatabase{
			MainDB:    mainDB,
			SlaveDBs:  slaveDBs,
			Completed: true,
		}

		logs.Info().Msg("database connection initialized successfully")
		return
	}
}

func initDB(config *DatabaseConfig, loc *time.Location) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4",
		config.User,
		config.Password,
		config.Hosts[0], // 使用第一个host
		config.Port,
		config.Name,
	)
	dsn += "&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
		Logger:                                   &DBLogger{LogLevel: config.LogLevel},
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
		NowFunc: func() time.Time {
			return time.Now().In(loc)
		},
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(utils.GetEnvInt("DB_MAX_IDLE_CONNS", 50))
	sqlDB.SetMaxOpenConns(utils.GetEnvInt("DB_MAX_OPEN_CONNS", 200))
	sqlDB.SetConnMaxLifetime(time.Duration(utils.GetEnvInt("DB_CONN_MAX_LIFETIME", 30)) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(utils.GetEnvInt("DB_CONN_MAX_IDLE_TIME", 5)) * time.Minute)
	return db, nil
}

// GetDB 获取主库连接
func (l *LightDatabase) GetDB(ctx context.Context) (*gorm.DB, error) {
	if l == nil {
		return nil, fmt.Errorf("database is not initialized")
	}
	if !l.Completed {
		return nil, fmt.Errorf("database is not completed, error: %v", l.Error)
	}

	return l.MainDB, nil
}

// GetSlaveDB 获取从库连接，实现负载均衡
func (l *LightDatabase) GetSlaveDB(ctx context.Context) (*gorm.DB, error) {
	if l == nil {
		return nil, fmt.Errorf("database is not initialized")
	}
	if !l.Completed {
		return nil, fmt.Errorf("database is not completed, error: %v", l.Error)
	}

	slaveDB := l.SlaveDBs[rand.Intn(len(l.SlaveDBs))]

	return slaveDB, nil
}

type DBLogger struct {
	LogLevel logger.LogLevel
}

func (l *DBLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *DBLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logs.Info().Msgf(msg, data...)
}

func (l *DBLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logs.Warn().Msgf(msg, data...)
}

func (l *DBLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logs.Error().Msgf(msg, data...)
}

func (l *DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil && l.LogLevel >= logger.Error {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		logs.Error().Err(err).Str("sql", sql).Int64("rows", rows).Msg("database error")
	} else if elapsed > 200*time.Millisecond && l.LogLevel >= logger.Warn {
		logs.Warn().Str("sql", sql).Int64("rows", rows).Msg("database slow query")
	} else if l.LogLevel >= logger.Info {
		logs.Info().Str("sql", sql).Int64("rows", rows).Msg("database query")
	}
}

// CloseConnections 提供一个方法用于安全地关闭数据库连接
// 应仅在确认不再需要使用数据库时调用，例如应用程序关闭时
func (l *LightDatabase) CloseConnections() {
	if l == nil || !l.Completed {
		return
	}

	if l.MainDB != nil {
		sqlDB, err := l.MainDB.DB()
		if err != nil {
			logs.Error().Err(err).Msg("error getting main DB connection while closing")
		} else {
			sqlDB.Close()
			logs.Info().Msg("main database connection closed")
		}
	}

	for i, slaveDB := range l.SlaveDBs {
		if slaveDB != nil {
			sqlDB, err := slaveDB.DB()
			if err != nil {
				logs.Error().Err(err).Int("slave_index", i).Msg("error getting slave DB connection while closing")
			} else {
				sqlDB.Close()
				logs.Info().Int("slave_index", i).Msg("slave database connection closed")
			}
		}
	}
}
