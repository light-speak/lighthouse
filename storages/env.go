package storages

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/light-speak/lighthouse/logs"
	"github.com/light-speak/lighthouse/utils"
)

// StorageDriver 存储驱动类型
type StorageDriver string

const (
	DriverS3  StorageDriver = "s3"  // S3 兼容存储
	DriverCOS StorageDriver = "cos" // 腾讯云 COS
)

// StorageConfig 存储配置
type StorageConfig struct {
	Driver StorageDriver // 存储驱动

	// S3 配置
	S3 struct {
		Endpoint        string // 端点地址
		AccessKeyID     string // 访问密钥 ID
		SecretAccessKey string // 访问密钥
		UseSSL          bool   // 是否使用 SSL
		DefaultBucket   string // 默认存储桶
	}

	// COS 配置
	COS struct {
		SecretID      string // 密钥 ID
		SecretKey     string // 密钥
		Region        string // 地域
		DefaultBucket string // 默认存储桶
	}
}

var (
	config  *StorageConfig
	storage Storage
)

func init() {
	// 初始化默认配置
	config = &StorageConfig{
		Driver: DriverS3,
	}
	config.S3.Endpoint = "localhost:9000"
	config.S3.AccessKeyID = ""
	config.S3.SecretAccessKey = ""
	config.S3.UseSSL = false
	config.S3.DefaultBucket = "default"

	config.COS.SecretID = ""
	config.COS.SecretKey = ""
	config.COS.Region = "ap-beijing"
	config.COS.DefaultBucket = "default"

	// 加载 .env 文件
	if curPath, err := os.Getwd(); err == nil {
		err = godotenv.Load(filepath.Join(curPath, ".env"))
		if err != nil {
			log.Println("Error loading .env file:", err)
		}
	}

	// 读取存储驱动配置
	config.Driver = StorageDriver(utils.GetEnv("STORAGE_DRIVER", string(config.Driver)))

	// 读取 S3 配置
	config.S3.Endpoint = utils.GetEnv("S3_ENDPOINT", config.S3.Endpoint)
	config.S3.AccessKeyID = utils.GetEnv("S3_ACCESS_KEY", config.S3.AccessKeyID)
	config.S3.SecretAccessKey = utils.GetEnv("S3_SECRET_KEY", config.S3.SecretAccessKey)
	config.S3.UseSSL = utils.GetEnvBool("S3_USE_SSL", config.S3.UseSSL)
	config.S3.DefaultBucket = utils.GetEnv("S3_DEFAULT_BUCKET", config.S3.DefaultBucket)

	// 读取 COS 配置
	config.COS.SecretID = utils.GetEnv("COS_SECRET_ID", config.COS.SecretID)
	config.COS.SecretKey = utils.GetEnv("COS_SECRET_KEY", config.COS.SecretKey)
	config.COS.Region = utils.GetEnv("COS_REGION", config.COS.Region)
	config.COS.DefaultBucket = utils.GetEnv("COS_DEFAULT_BUCKET", config.COS.DefaultBucket)

	// 初始化存储实例
	initStorage()
}

// initStorage 根据配置初始化存储实例
func initStorage() {
	var err error

	switch config.Driver {
	case DriverS3:
		if config.S3.AccessKeyID == "" || config.S3.SecretAccessKey == "" {
			logs.Warn().Msg("S3 storage not configured properly, skipping initialization")
			return
		}

		s3Config := S3Config{
			Endpoint:        config.S3.Endpoint,
			AccessKeyID:     config.S3.AccessKeyID,
			SecretAccessKey: config.S3.SecretAccessKey,
			UseSSL:          config.S3.UseSSL,
		}
		storage, err = NewS3Storage(s3Config)
		if err != nil {
			logs.Error().Err(err).Msg("failed to initialize S3 storage")
			return
		}
		logs.Info().Msgf("S3 storage initialized successfully, bucket: %s", config.S3.DefaultBucket)

	case DriverCOS:
		if config.COS.SecretID == "" || config.COS.SecretKey == "" {
			logs.Warn().Msg("COS storage not configured properly, skipping initialization")
			return
		}

		cosConfig := COSConfig{
			SecretID:  config.COS.SecretID,
			SecretKey: config.COS.SecretKey,
			Region:    config.COS.Region,
		}
		storage, err = NewCOSStorage(cosConfig)
		if err != nil {
			logs.Error().Err(err).Msg("failed to initialize COS storage")
			return
		}
		logs.Info().Msg("COS storage initialized successfully")

	default:
		logs.Error().Str("driver", string(config.Driver)).Msg("unsupported storage driver")
	}
}

// GetStorage 获取存储实例
func GetStorage() (Storage, error) {
	if storage == nil {
		return nil, ErrStorageNotInitialized
	}
	return storage, nil
}

// GetConfig 获取存储配置
func GetConfig() *StorageConfig {
	return config
}

// GetDefaultBucket 获取默认存储桶名称
func GetDefaultBucket() string {
	switch config.Driver {
	case DriverS3:
		return config.S3.DefaultBucket
	case DriverCOS:
		return config.COS.DefaultBucket
	default:
		return "default"
	}
}
