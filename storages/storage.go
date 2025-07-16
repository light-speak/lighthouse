package storages

import (
	"context"
	"errors"
	"io"
	"time"
)

// 定义存储相关错误
var (
	ErrStorageNotInitialized = errors.New("storage not initialized")
	ErrDriverNotSupported    = errors.New("storage driver not supported")
)

// Storage 统一存储接口
// 支持多种存储后端实现，如 S3、COS 等
type Storage interface {
	// Put 上传文件到指定存储桶
	// ctx: 上下文
	// bucket: 存储桶名称
	// key: 文件路径/键名
	// reader: 文件内容读取器
	Put(ctx context.Context, bucket string, key string, reader io.Reader) error

	// GetPresignedURL 获取预签名下载 URL
	// ctx: 上下文
	// bucket: 存储桶名称
	// key: 文件路径/键名
	// expiry: 过期时间
	GetPresignedURL(ctx context.Context, bucket string, key string, expiry time.Duration) (string, error)

	// GetPresignedPutURL 获取预签名上传 URL
	// ctx: 上下文
	// bucket: 存储桶名称
	// key: 文件路径/键名
	// expiry: 过期时间
	GetPresignedPutURL(ctx context.Context, bucket string, key string, expiry time.Duration) (string, error)

	// GetPublicURL 获取公开访问 URL
	// bucket: 存储桶名称
	// key: 文件路径/键名
	GetPublicURL(bucket string, key string) string
}
